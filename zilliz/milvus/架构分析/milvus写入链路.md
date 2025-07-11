
milvus的总体架构是一个分布式存储系统，数据的写入，存储，查询分别由不同的节点完成。来自sdk端的不同的请求需要先经过proxy代理解析然后分发不同的节点进行进一步处理。本文重点解析从sdk发送写入数据请求，即insert后，整个操作链路是怎样的。

写入数据属于DML，涉及到proxy层和StreamingNode。

当sdk端发送insert插入数据请求时，不妨以一段的代码为例：

```python
rng = np.random.default_rng(seed=19530)
rows = [
    {"id": 1, "embeddings": rng.random((1, dim))[0], "a": 100, "title": "t1"},
    {"id": 2, "embeddings": rng.random((1, dim))[0], "b": 200, "title": "t2"},
    {"id": 3, "embeddings": rng.random((1, dim))[0], "c": 300, "title": "t3"},
    {"id": 4, "embeddings": rng.random((1, dim))[0], "d": 400, "title": "t4"},
    {"id": 5, "embeddings": rng.random((1, dim))[0], "e": 500, "title": "t5"},
    {"id": 6, "embeddings": rng.random((1, dim))[0], "f": 600, "title": "t6"},
]

insert_result = milvus_client.insert(collection_name, rows)
```

接下来我们分析这段代码会如何被milvus处理,发起插入请求,将数据写入WAL存储,然后最终写入对象存储.

## sdk包装

首先,这段代码的数据被pymilvus sdk内部包装为proto定义的InsertRequest:

![[../../../附件/截屏2025-07-08 14.44.53.png]]

在请求中,我们的数据被包装到fields_data中:

![[../../../附件/截屏2025-07-08 14.45.52.png]]

这是对所有fielddata的封装,其中的几个字段中

1. type是数据的类型,定义了milvus中的所有数据类型,目前已经包括了GEO数据类型
2. field通过oneof关键字指定了数据的类型只能有三种:标量,向量,和一个递归的字段,具体解释如下:

- scalars (ScalarField): 用于存储标量类型的数据。ScalarField 内部也是一个 oneof，根据具体的数据类型（如 Int64, String, Bool 等）包含一个对应类型的数组（例如 LongArray、StringArray）。所以 scalars 字段存储的是一列标量值。

- vectors (VectorField): 用于存储向量类型的数据。VectorField 内部包含了向量的维度 (dim) 和具体的向量数据（如 float_vector, binary_vector 等）。

- struct_arrays (StructArrayField): 用于存储 ArrayOfStruct (结构体数组) 类型的数据。StructArrayField 的定义是 repeated FieldData fields = 1;。这是一个递归定义。这意味着一个结构体数组字段的数据，是由一个 FieldData 列表组成的，列表中的每个 FieldData 元素对应结构体中的一个子字段（sub-field）。

举个例子: 假设有一个字段是结构体数组，每个结构体包含 {"name": "xxx", "age": yy}。那么 struct_arrays 里面就会包含两个 FieldData 对象：
- 第一个 FieldData 用来存所有 name 的值（一个字符串数组）。
- 第二个 FieldData 用来存所有 age 的值（一个整数数组）。

## 请求解析和处理

在sdk包装了inset请求后,通过grpc连接将这个请求发送给milvus的proxy服务器.

服务器接收到请求并解析到是insert请求后,会调用internal/proxy/impl.go的insert函数:

```go
// internal/proxy/grpc_service.go
func (s *Server) Insert(ctx context.Context, request *milvuspb.InsertRequest) (*milvuspb.MutationResult, error) {
    // ... (初始设置和日志记录)

    // 1. 创建一个 InsertTask，这是处理插入请求的核心逻辑单元
    insertTask := &insertTask{
        // ... (初始化任务所需上下文、请求信息等)
    }
    // 2. 将任务放入一个队列中等待执行
    err := s.sched.dmQueue.Enqueue(insertTask)
    // ... (错误处理)
    // 3. 等待任务执行完成并返回结果
    err = insertTask.WaitToFinish() //在这里阻塞并等待执行完成
    if err != nil {
        return nil, err
    }
    return insertTask.result, nil
}
```

在任务Enqueue后,insert本身通过waittofinish阻塞等待任务完成,此时会有一个内部的taskscheduler来处理这个任务

taskscheduler在启动的时候,会并发启动四个队列,维护不同类型的任务

![[../../../附件/截屏2025-07-08 15.34.04.png]]

manipulationloop函数内部有一个无限循环,不断地从dmQueue中拿出任务并执行:

![[../../../附件/Pasted image 20250708153851.png]]

在process中,会对一个任务分别调用preExcute和Excute,前者完成一些验证,元数据校验,分区处理等工作,后者是任务执行的核心,insert到此处才开始正式执行,前面的所有任务都是为了应对大量并发请求而做的优化处理.这里调用Excute会根据具体任务的不同,调用其特定的excute函数.

Excute会做这些事:

- 获取 DML Channel: 调用 it.chMgr.getVChannels(collID) 来获取应该将数据写入到哪些 DML Channel。一个集合通常对应多个 DML Channel 以支持并发写入。
- 分配 SegmentID: Proxy 会调用 segIDAssigner.GetSegmentID(...) 来为这批数据分配一个或多个 Segment ID。segIDAssigner 内部会与 Data Coord 通信，告诉它：“我有一批新数据要写，请给我一个地方（Segment）写。” Data Coord 会返回可用的 Segment ID。
- 重打包数据 (repackInsertData): 根据分配到的 Segment ID，原始的 InsertMsg 会被重新组织成多个小的 InsertMsg，每个消息对应一个 Segment。这样做是为了将数据分发到不同的 Segment.
- 发送到消息队列: 最后，调用 stream.Produce(ctx, msgPack) 将打包好的 MsgPack（包含一个或多个 InsertMsg）发送到消息队列中

## 数据写入

在streamingnode链路中,streamingnode负责将数据写入WAL和对象存储
