
在main函数中,`RunMilvus`启动milvus

调用cmd/milvus/milvus.go `c.execute`执行

cmd/milvus/run.go中execute中,使用roles.Run()启动milvus组件

cmd/roles/roles.go Run()中,调用runStreamingNode()启动StreamingNode,在此函数中,通过向runComponent()传入NewStreamingNode()函数来完成StreamingNode的创建

在cmd/components NewStreamingNode()中,通过调用distributed的streamingnode下NewServer服务来新建一个StreamingNode server,至此StreamingNode server成功创建,并将这个Server返回给runComponent()函数中的role对象

下面是启动流程:

此时runComponent还没有执行完成,role对象通过role.Run(),来真正调用streamingNode的Run()方法(internal/distributed/streamingnode/service.go)

在Run()中,首先调用init()函数完成一系列初始化,包括etcd服务,chunkManager(storage层读写操作)等,其中也包括了grpc服务的启动,在initGRPCServer中,设置了一系列的参数,然后注册了一个state服务.之后通过一系列的build函数,给服务器设置了一些参数,然后调用Build()来完成构建.

在Build()中,首先由builder创建了一个真正的server,随后调用server的init()函数

在server的init()中,通过initService()函数,来创建了streamingnode的handlerService和managerService,并最终注册了server的gprc服务.

对于其中的handlerService,有Produce和Consume两个方法,分别在一个load node上输出数据和消费数据.

Produce:首先新建了一个ProduceServer,随后调用execute来启动了这个server,在execuet中,开启一个go routine recvloop启动一个无限循环用于处理来自client的信息,这里包含了处理WAL请求的逻辑.当recvloop收到produce请求时,调用handleProducer,它接受一个ProduceMessageRequest类型的参数,这个参数中包含了一个Message类型,名为Message的成员,***而这个成员中的Payload,正是我们序列化后的数据存储的地方***.

回到handleProduce,它先进行一个检查工作,然后使用payload等创建了一个msg变量,然后向WAL服务传递这个msg,WAL开始进行异步写入.

而sendloop则是根据WAL和等待队列等的状态,来向client发送处理结果.

Consume:Produce负责将数据写入WAL,而consume负责从wal读取数据并返回给client.首先新建一个consumeserver然后execute.在create的过程中,consumeserver中有一个scanner字段获得了从wal读取的能力.然后同样启动两个循环,其中sendloop是主循环,不断地从scanner中读取msg,如果收到了client的读取请求,就会向client发送从WAL中读取到的数据.

在两个服务注册完成并启动后,build结束.返回到internal/distributed/StreamingNode/service.go/Run()

接下来Run()继续运行,调用start()函数来启动streamingnode server,其中启动了grpcserver,在启动过程中启动一个协程,用于监听特定listenner上的信息.随后再进行一些额外的工作后,streamingnode server正式启动.

然后是如何将从wal读取的数据写入对象存储:

同样是在roles.go的Run()中,其实在上一节所说的执行runstreamingnode之前,先执行了streaming.init(),在这个init()函数中,调用newWALAccesser()来创建了一个walaccesser,也就是后来消费WAL数据的client.

在newWALAccesser()中,创建了streamingcoordclient和handlerclient两个客户端,随后返回init中注册并启动

在调用newhandlerclient来创建handlerclient时,调用了又调用了createproducer和createconsumer两个函数分别注册了两个服务.

Produce:负责将数据写入WAL,它启动sendloop和recvloop,前者将produce message送给server处理,后者接受server返回的处理结果

Consumer:只启动一个recvloop(),不断接受从服务端返回的结果,并反序列化其中的数据,创建newImmutableMsg,对于insert消息,调用Handle来处理

之后消息被路由,然后交由pipeline处理

此处insert消息是由querynode处理的,querynode服务端提供了watchdmchannel服务,在这个服务中,调用了一个关键函数:ConsumeMsgStream,这个函数对消息进行了路由转发

随后返回watchdmchannel函数,调用pipeline.start()启动了pipeline,在启动时,开启一个work协程,不断接受msg,并将其传递给inputChannel,随后process,其中将消息传递给node的Operate进行处理,
这里调用的是insert_node的Operate函数,其中调用了addInsertData函数,这个函数中又调用了TransferInsertMsgToInsertRecord,将msg消息转换为适合存储的格式,在这个函数中,根据不同的数据类型进行了相应的处理.(internal/storage/utils)

随后返回,调用MergeFieldData,根据不同的数据结构,将数据追加到特定segment的fieldsdata中

在addInsertData函数结束后,返回Operate,再调用delegator.ProcessInsert,寻找growing状态的segment,如果没有就创建一个,随后调用其上的Insert函数,至此我们已经到达了数据写入的关键步骤.

在Insert函数中,进一步调用segCore的Insert函数(internal/util/segcore/segment.go/Insert()),这个函数进行了go与c++层的转化,首先将一些元信息转化为C类型,接着调用了c++的Insert函数(),这个函数的声明在internal/core/src/segcore/segment_c.h:

```c
CStatus
Insert(CSegmentInterface c_segment,
int64_t reserved_offset,
int64_t size,
const int64_t* row_ids,
const uint64_t* timestamps,
const uint8_t* data_info,
const uint64_t data_info_len);
```

它的实现中调用了一个GrowingSegment的Insert函数,这个函数将数据写入concurrentvector,只是写入了memory.

另外,这个Insert函数中调用了append_field_data,这里面对不同的数据类型进行了处理.

然后是数据最终刷入对象存储的步骤:

省略从服务启动的步骤,直接到syncmgr的Run函数启动,这个函数中根据storage版本的不同,调用不同的bulkpackwriter,返回一个writer对象,然后调用这个writer对象上的Write方法,这个函数中分别调用了:

1. writeInserts
2. writeStats
3. writeDelta
4. writeBM25Stasts

这些函数内部都调用writelog,而writelog中调用storage的chunkManager提供的Write函数

这个函数根据不同的写入对象,最终处理写入请求,如调用三方服务的api等.至此,整个从streamingnode启动到数据写入的链路已清晰.

还有几个细节问题:

1. 为什么数据从WAL中读出后要先insert到segment里管理,然后再刷入存储,不是直接刷进存储层?
2. 为什么要对数据进行这么多层的包装,传给这么多组件来处理?
3. 为什么要将数据存储的功能交给streamingnode来做?StreamingNode相对datanode有什么优势?
4. 还有具体的代码细节,没有时间再看明白了.


# 整理版本

以下是对原始描述的整理与优化，保留原有技术流程不变，仅对语言结构、术语规范和逻辑表达进行了清晰化处理：

---

## Milvus StreamingNode 启动及数据写入链路详解

### 一、StreamingNode 启动流程

1. **main 函数中调用 `RunMilvus` 启动 Milvus**
   - 实际调用路径为：`cmd/milvus/milvus.go` 中的 `c.Execute()`。

2. **在 `cmd/milvus/run.go` 的 `execute` 方法中**
   - 调用 `roles.Run()` 来启动 Milvus 各组件。

3. **`cmd/roles/roles.go` 中的 `Run()` 方法**
   - 调用 `runStreamingNode()` 来启动 StreamingNode。
   - 在该函数内部，通过将 `NewStreamingNode()` 作为参数传入 `runComponent()`，完成 StreamingNode 的创建。

4. **`cmd/components/streamingnode.go` 中的 `NewStreamingNode()`**
   - 调用 `distributed/streamingnode/service.go` 中的 `NewServer()` 方法，创建一个 StreamingNode Server。
   - 此时，StreamingNode Server 创建成功，并赋值给 `runComponent` 中的 `role` 对象。

5. **`runComponent` 继续执行，调用 `role.Run()`**
   - 最终调用 `internal/distributed/streamingnode/service.go` 中的 `Run()` 方法。

6. **`Run()` 方法执行步骤如下**：
   - 首先调用 `init()` 完成初始化，包括 etcd 连接、chunkManager 初始化（storage 层读写操作）、gRPC 服务启动等。
   - 在 `initGRPCServer()` 中设置 gRPC 参数并注册状态服务。
   - 接着调用一系列 `build` 方法完成构建过程。

7. **构建阶段 (`Build()`) 执行如下**：
   - 由 builder 构建真正的 server 实例。
   - 调用 server 的 `init()` 方法。
     - 在 `init()` 中调用 `initService()`，创建 handlerService 和 managerService。
     - 并最终注册 server 的 gRPC 服务。

8. **HandlerService 提供两个核心方法**：
   - **Produce**：用于向 StreamingNode 写入数据。
   - **Consume**：用于从 StreamingNode 消费数据。

---

### 二、Produce 数据写入流程

- **Produce 流程**：
  1. 新建 ProduceServer 并调用其 `Execute()` 方法。
  2. 在 `Execute()` 中开启一个 goroutine `recvLoop`，进入无限循环监听客户端请求。
  3. 当收到 `ProduceRequest` 时，调用 `handleProducer` 方法。
  4. `handleProducer` 中提取出请求中的 `Message`，其中包含序列化后的 payload 数据。
  5. 将 payload 数据封装为 msg，发送至 WAL（Write-Ahead Log）服务进行异步写入。

- **SendLoop**：
  - 根据 WAL 状态及等待队列情况，向客户端发送写入结果反馈。

---

### 三、Consume 数据消费流程

- **Consume 流程**：
  1. 新建 ConsumeServer 并调用其 `Create()` 方法。
  2. `Create()` 中 scanner 字段获得从 WAL 读取的能力。
  3. 启动 `sendLoop` 循环，持续从 scanner 读取消息。
  4. 若有客户端发起读取请求，则将读取到的消息返回给客户端。

---

### 四、WAL 数据写入对象存储流程

#### 1. **StreamingNode 初始化阶段**

- 在 `roles.Run()` 中，`runStreamingNode()` 之前会调用 `streaming.Init()`。
- 在 `Init()` 中调用 `newWALAccesser()` 创建 WAL Accesser。

#### 2. **newWALAccesser() 执行内容**：

- 创建 `streamingCoordClient` 和 `handlerClient` 两个客户端。
- 注册并启动相关服务。

#### 3. **HandlerClient 创建过程**：

- 调用 `createProducer()` 和 `createConsumer()` 分别注册生产者和消费者服务。

##### Producer 功能：

- 发送 loop：将 produce message 发送给 server 处理。
- 接收 loop：接收 server 返回的处理结果。

##### Consumer 功能：

- 开启 recvLoop，不断接收 server 返回的数据。
- 反序列化后创建 `ImmutableMsg`。
- 对于 insert 类型消息，调用 `Handle()` 进行处理。

#### 4. **消息路由与 pipeline 处理**：

- Insert 消息被路由至 QueryNode。
- QueryNode 提供 `WatchDmChannel` 服务。
  - 在该服务中调用 `ConsumeMsgStream()`，实现消息转发。
  - 返回后调用 `pipeline.Start()` 启动 pipeline。
    - 启动 work 协程，不断接收 msg。
    - 将 msg 传递给 inputChannel。
    - 经过 process 流程，调用 node 的 `Operate()` 方法。

#### 5. **Insert Node 的 Operate 方法**：

- 调用 `addInsertData()`。
  - 内部调用 `TransferInsertMsgToInsertRecord()`，根据数据类型进行转换。
  - 调用 `MergeFieldData()`，将数据追加到 segment 的 fieldsData 中。

- 返回后调用 `delegator.ProcessInsert()`：
  - 查找处于 growing 状态的 segment。
  - 若不存在则新建一个 segment。
  - 调用其 `Insert()` 方法。

#### 6. **Segment Insert 实现**：

- 调用 `segCore.Insert()`（位于 `internal/util/segcore/segment.go`）。
  - 将元信息转为 C 类型。
  - 调用 C++ 层 `C_SegmentInterface_Insert()` 函数（定义于 `internal/core/src/segcore/segment_c.h`）。
  - 实现类为 `GrowingSegment`，数据被插入到内存中的 `concurrent_vector`。

---

### 五、数据刷盘流程

- **SyncMgr 启动**：
  - 在服务启动过程中，syncmgr 的 `Run()` 函数被调用。
  - 根据 storage 版本选择不同的 `BulkPackWriter`。
  - 创建 writer 对象并调用其 `Write()` 方法。

- **Write 方法执行步骤**：

1. `writeInserts`
2. `writeStats`
3. `writeDelta`
4. `writeBM25Stats`

- 所有上述函数内部均调用 `writeLog()`。
- `writeLog()` 调用 `chunkManager.Write()`，完成实际写入。
- `chunkManager` 根据目标存储类型（如 S3、MinIO、本地文件系统等），调用相应 API 完成数据持久化。

---

### 六、关键问题分析（待深入）

1. **为什么数据需要先写入 segment，再刷入对象存储？**
   - segment 是内存中的数据结构，便于快速访问和查询。
   - 刷盘前需进行格式转换、压缩、索引生成等预处理。

2. **为何设计多层包装与组件流转？**
   - 为了实现模块解耦、功能隔离、扩展性增强。
   - 支持多种数据类型、存储引擎、计算节点的灵活组合。

3. **为何由 StreamingNode 负责数据写入？与 DataNode 的区别？**
   - StreamingNode 更强调流式处理能力，支持实时写入与消费。
   - DataNode 更侧重离线批量处理与持久化。
   - StreamingNode 更适合实时场景下的数据缓冲与分发。

4. **代码细节尚待进一步理解**
   - 包括各组件之间的通信机制、gRPC 接口定义、WAL 的具体实现方式等。

---

### 总结

整个数据从 StreamingNode 接收、写入 WAL、经过 QueryNode 处理、最终写入对象存储的过程，涉及多个模块间的协作。其设计体现了 Milvus 对实时性、可扩展性和高可用性的追求。后续可通过阅读源码进一步深入理解各组件的交互细节与性能优化策略。