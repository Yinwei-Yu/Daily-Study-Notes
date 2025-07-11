我们以一个简单的例子出发,来搞清楚数据查询链路,即这个查询语句怎么被proxy转发,其中的各个字段又是如何被处理的,尤其是expr表达式是如何解析,并转化为底层代码的查询条件的,然后结果又是如何返回的:

```python
client.query(
	"colleciton_name",
	expr="id>10",
	output_field=["id","vector"]
)
```


当sdk发送查询请求后,被proxy分发到dq队列,在进行process时,同样先进行preExcute:

这里会先解析一些与query相关的参数,比如limit等,我们关心的关键部分在于expr的解析:

首先调用了CreateRetrievePlan,在这个函数里,调用ParseExpr对表达式进行解析,它调用handleExpr来进行实际的解析工作,这个函数又分别调用handleInternal来进行AST解析,然后将AST转化为protobuf

先来看handleInternal:

首先尝试从缓存中拿数据,没有的话:

1. 调用converHanToASCII:将汉字转化为utf8编码
2. 调用NewInputStream:创建ANTLR输入流,此处即调用了第三方库的服务
3. 调用getLexer:创建解析器
4. 调用parser.Expr():对expr进行解析->这个函数就是根据Plan.g4文件产生的
5. 缓存

在parser.Expr()中,调用expr(_ p int),接下来就是ANTLR自动实现的逻辑了,我们无需关心,只需要知道表达式会被正确解析就好了,这里的关键是在Plan.g4,它规定了我们表达式解析的方法.

在handleInternal返回后,我们有了一个解析好的AST,接着调用ast上的Accept方法,将ast解析并包装成protobuf,之后解析完毕,返回preExcute函数,然后对这个消息序列化.在Excute阶段,调用queryShard函数,这个函数中调用queryNode client的Query函数,并传入query请求,至此,完成解析后的真正的query开始了.

在Query的Query函数中(internal/querynodev2/services.go),收到来自dq channnel的请求,并对每个请求启动一个协程来处理

这个处理过程中涉及到queryChannel函数,他会路由并转发这个请求.然后调用shardDelegator的Query函数来处理,它会启动子任务来处理这些查询请求,其中worker.QuerySegments(internal/querynodev2/services.go)真正起到了查询的作用.

它会新创建一个query task,然后加入调度器,接着等待执行完成.

在这个任务被调度后,调用Execute(internal/querynodev2/tasks/query_task.go),这个函数中首先从序列化的信息中解析出retrievePlan(涉及到c++层),然后调用segments的Retrieve方法

在Retrieve中,先获得segments,然后调用retrieveOnSegments,进行真正的检索工作.这个函数里定义了一个retriever的函数类型变量,然后在每个segments上执行它,调用LocalSegment.retrieve方法,这个方法里真正调用了c++层的检索方法


在c++层,根据表达式执行查询计划,然后返回一个bitset_holder,用01表示是否符合结果,然后包装结果并返回.

在整个查询过程中,涉及到的主要是表达式部分,因此我们需要修改的地方主要有:

1. Plan.g4,增加geo查询表达式的支持
2. c++parser层,增加对geo表达式的解析处理

主要涉及到:

1. plan.g4
2. internal/core/src/expr/exec/expression/Expr.cpp + gis解析器
3. ITypeExpr.h:


## 整理

## 🧭 总体数据流概览

```plaintext
client.query(...) → Proxy → DQ队列 → PreExecute → CreateRetrievePlan
→ ParseExpr → ANTLR解析expr → AST → AST.Accept() → Protobuf Plan
→ Serialize Plan → Execute → QueryShard → QueryNode.Query()
→ QueryChannel → shardDelegator → worker.QuerySegments
→ task.Execute → Retrieve → retrieveOnSegments → LocalSegment.retrieve
→ C++检索逻辑执行 → bitset_holder → 封装结果 → 返回
```

### 1. **Python SDK 发出请求**

```python
client.query(
	"collection_name",
	expr="id > 10",
	output_fields=["id", "vector"]
)
```

生成一个 gRPC 请求，核心字段：
→ 请求转发到 Proxy

### 2. **Proxy 层：表达式预处理**

#### 🔹 文件路径：

`internal/proxy/query.go`

#### 🔹 函数路径：

```go
PreExecuteQuery()
  → createRetrievePlan()
    → ParseExpr()
      → handleExpr()
        → handleInternal() // 核心入口
```

#### 🔹 handleInternal 做了什么：

1. **缓存命中判断**
    
2. **调用 ANTLR 解析表达式**
    - `convertHanToASCII()`：转义中文为 ASCII
    - `NewInputStream()`：构造输入流
    - `getLexer()`：构建词法分析器
    - `parser.Expr()`：调用由 ANTLR 根据 `Plan.g4` 生成的 `Expr()` 函数来构造 AST
    - **关键依赖：`Plan.g4` 文件决定了语法结构**

> ✅ 表达式语法解析完成，生成 AST（ANTLR Tree）

---

### 3. **AST 转化为 Protobuf Plan**

#### 🔹 文件路径：

`internal/core/src/expr/plan_converter.cpp`  
通过调用 `ast.Accept(visitor)`：

```cpp
auto plan = ast.Accept(visitor);
```

#### 🔹 visitor 的作用：

- 把 AST 转换为逻辑查询计划（Protobuf Plan 格式）
- 表达式变成逻辑节点（如：CompareNode、TermNode、RangeNode）

### 4. **Proxy → QueryNode 分发任务**

#### 🔹 文件路径：

- `internal/querynodev2/services.go`
- `QueryNode.Query()` 监听来自 Proxy 的请求，解包 Plan 后，分发任务

---

### 5. **执行阶段：Query Task**

#### 🔹 Query Task 启动流程：

```go
worker.QuerySegments() → 创建 queryTask → scheduler 调度
task.Execute() → segment.Retrieve()
```

#### 🔹 目标：

- 根据 Plan 选出数据片段（segment）
- 在 segment 上执行匹配操作

---

### 6. **C++ 层执行表达式计划**

#### 🔹 文件路径：

- `internal/core/src/query/Retrieve.cpp`
- `internal/core/src/expr/expression/Expr.cpp`：表达式执行逻辑
- `internal/core/include/expr/ITypeExpr.h`：表达式接口定义
- `internal/core/src/expr/exec/`：具体的执行器目录

#### 🔹 表达式执行：

- Plan 会转为布尔表达式树
- 在每条记录上执行布尔表达式 → 生成 bitset
- bitset 表示是否满足查询条件
- 满足的记录 → 返回
