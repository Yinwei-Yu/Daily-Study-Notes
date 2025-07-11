
## internal/storage 简介

`internal/storage` 模块为 Milvus 提供了一个**数据存储的抽象层**，通过 `ChunkManager` 接口对外提供与文件操作相关的各种功能接口。

### 对外提供的核心接口（ChunkManager）

在 `internal/storage/types.go` 中定义了如下接口：

#### 文件操作
- **Write(ctx, filePath, content)** - 将数据写入存储系统  
- **Read(ctx, filePath)** - 从存储中读取数据  
- **MultiWrite(ctx, contents)** - 批量写入多个文件  
- **MultiRead(ctx, filePaths)** - 批量读取多个文件  
- **Exist(ctx, filePath)** - 判断文件是否存在  
- **Remove(ctx, filePath)** - 删除指定文件  
- **MultiRemove(ctx, filePaths)** - 批量删除多个文件  
- **RemoveWithPrefix(ctx, prefix)** - 根据前缀批量删除文件  

#### 文件管理
- **Size(ctx, filePath)** - 获取文件大小  
- **Path(ctx, filePath)** - 获取文件完整路径  
- **RootPath()** - 获取存储系统的根路径  
- **Reader(ctx, filePath)** - 获取文件读取器  
- **ReadAt(ctx, filePath, offset, length)** - 读取文件特定偏移和长度的内容  
- **Mmap(ctx, filePath)** - 使用内存映射方式读取文件  

这些接口由不同的实现类来完成，支持本地和云端等多种存储方式。

---

## 实现类说明

1. **LocalChunkManager** - 本地文件系统存储实现
2. **RemoteChunkManager** - 云存储通用封装（支持 MinIO、AWS S3、Azure、GCP）
3. **MinioChunkManager** - MinIO 对象存储专用实现
4. **AzureChunkManager** - Azure Blob 存储实现
5. **GcpNativeChunkManager** - Google Cloud Storage 实现
6. **OpenDALChunkManager** - 基于 OpenDAL 的统一存储接口封装

其中，`LocalChunkManager` 完整实现了 `ChunkManager` 接口的所有方法，后几个则是在RemoteChunkManager里做了统一的封装

---

## 谁使用了这些接口？

### QueryNodeV2 模块
用于 segment 的加载、数据检索与缓存：
- **segments/segment.go** - Segment 核心操作逻辑  
- **segments/segment_loader.go** - 从存储中加载 segment 数据  
- **segments/segment_interface.go** - Segment 接口定义  
- **segments/utils.go** - 存储操作工具函数  

### Delegator 组件
处理数据转发与裁剪：
- **delegator/delegator.go** - 主要逻辑控制  
- **delegator/delegator_data.go** - 数据管理  
- **delegator/segment_pruner.go** - 基于存储统计信息进行 segment 裁剪  
- **delegator/scalar_pruner.go** - 标量字段裁剪  
- **delegator/delta_forward.go** - Delta 数据转发  
- **delegator/buffered_forwarder.go** - 缓冲数据转发  

### Pipeline 组件
数据处理流水线：
- **pipeline/insert_node.go** - 插入数据处理节点  
- **pipeline/delete_node.go** - 删除数据处理节点  
- **pipeline/embedding_node.go** - 向量嵌入处理节点  
- **pipeline/type.go** - 流水线类型定义  

### FlushCommon 模块
负责数据刷新到磁盘：
- **syncmgr/sync_manager.go** - 数据同步主控  
- **syncmgr/pack_writer.go** - 数据打包与写入  
- **syncmgr/pack_writer_v2.go** - 第二代打包写入器  
- **syncmgr/storage_serializer.go** - 数据序列化  
- **syncmgr/task.go** - 同步任务管理  

- **writebuffer/insert_buffer.go** - 插入数据缓冲区  
- **writebuffer/delta_buffer.go** - Delta 数据缓冲区  
- **writebuffer/stats_buffer.go** - 统计信息缓冲区  
- **writebuffer/segment_buffer.go** - Segment 缓冲区  
- **writebuffer/l0_write_buffer.go** - L0 写入缓冲区  

- **metacache/meta_cache.go** - 元数据缓存主控  
- **metacache/segment.go** - Segment 元数据  
- **metacache/actions.go** - 缓存操作行为  
- **metacache/bm25_stats.go** - BM25 统计信息  

- **pipeline/data_sync_service.go** - 数据同步服务  
- **pipeline/flow_graph_message.go** - 流图消息定义  

### StreamingNode 模块
流式节点相关：
- **server/server.go** - 主服务逻辑  
- **server/builder.go** - 服务构建器  
- **server/resource/resource.go** - 资源管理  
- **server/flusher/flusherimpl/flusher_components.go** - Flusher 组件  

- **WAL (Write-Ahead Log)**  
  - **server/wal/recovery/segment_recovery_info_test.go** - 恢复测试  
  - **server/wal/interceptors/shard/shards/segment_alloc_worker.go** - 分段分配工作器  

### 其他模块
包括元数据存储、导入工具、Mock 类等。

---

## 使用这些接口的主要目的

### 1. **ChunkManager 接口**
大多数组件通过该接口进行：
- 文件操作（Read, Write, Exist, Remove）  
- 多文件操作（MultiRead, MultiWrite）  
- 文件遍历（WalkWithPrefix）  

### 2. **InsertData 结构体**
广泛用于：
- 插入数据操作  
- 数据转换与处理  
- 高效内存数据管理  

### 3. **PayloadReader / PayloadWriter**
用于：
- 数据序列化与反序列化  
- Parquet 格式处理  
- Arrow 格式集成  

### 4. **统计信息与布隆过滤器**
用于：
- 主键统计（PrimaryKeyStats）  
- 布隆过滤器优化查找效率  
- 字段统计用于查询优化  

### 5. **Binlog 操作**
用于：
- 二进制日志读写  
- 事件处理  
- 数据恢复与重放  

---

## 附录:internal/storage 下各组件及其作用

| 文件名            | 功能描述                                             |
| -------------- | ------------------------------------------------ |
| **types.go**   | 定义核心接口与类型如 `ChunkManager`, `FieldID`, `UniqueID` |
| **factory.go** | 提供构造函数，创建不同类型的 ChunkManager 实例                   |
| **payload.go** | 提供序列化/反序列化接口，后续需增加 Geo 类型支持                      |

### 数据管理
- **insert_data.go** - 插入数据结构定义，包含 InsertData 及所有字段类型（Bool、Int8-64、Float、Double、String、Array、JSON、Vector 等）  
- **delta_data.go** - 删除操作结构定义（DeltaData, DeleteData）  
- **data_codec.go** - 插入数据的编解码逻辑（InsertCodec, DeleteCodec）  
- **index_data_codec.go** - 索引文件的编解码逻辑（IndexFileBinlogCodec, IndexCodec）  

### Binlog 系统
- **binlog_reader.go** - 二进制日志读取器  
- **binlog_writer.go** - 二进制日志写入器（插入、删除、索引文件）  
- **binlog_util.go** - 解析 binlog 路径中的 segment ID  
- **print_binlog.go** - debug 工具，打印 binlog 内容  

### 事件系统
- **event_writer.go** - 事件写入器（InsertBinlogWriter, DeleteBinlogWriter 等）  
- **event_reader.go** - 从 binlog 中读取事件  
- **event_header.go** - 事件头结构及读写逻辑  
- **event_data.go** - 不同事件类型的数据结构定义  

### Payload 处理
- **payload_writer.go** - 使用 NativePayloadWriter 写入 Parquet 格式  
- **payload_reader.go** - 从 Parquet 格式读取数据  
- **serde.go** - 使用 Arrow 格式的序列化/反序列化工具  
- **serde_events.go** - 事件数据的序列化/反序列化逻辑  
- **serde_events_v2.go** - V2 版本增强事件序列化逻辑  

### 数据处理与工具
- **utils.go** - 文件操作、数据转换、字段合并、内存管理等工具函数  
- **rw.go** - 支持版本控制和分块读取的 blob 读写  
- **sort.go** - 数据排序工具（Sort, MergeSort）  
- **data_sorter.go** - 按 row ID 排序 insert 数据  

### 字段与值管理
- **field_value.go** - ScalarFieldValue 和 VectorFieldValue 接口及其实现  
- **primary_key.go** - 主键接口及其实现（Int64/Varchar）  
- **field_stats.go** - 字段统计信息管理（最小值、最大值、布隆过滤器）  
- **partition_stats.go** - 分区级别统计信息（SegmentStats）  

### 统计与索引
- **stats.go** - 主键统计（PrimaryKeyStats）、布隆过滤器、BM25 统计  
- **pk_statistics.go** - 主键统计与布隆过滤器操作  
- **index_data_codec.go** - 索引文件编解码逻辑  

### Schema 与 Arrow 集成
- **schema.go** - Arrow schema 转换工具  
- **arrow_util.go** - Arrow 记录构建与数据类型转换工具  

### 性能与调试
- **unsafe.go** - 高性能内存操作（用于关键数据读取）  
- **print_binlog_test.go** - binlog 打印功能的单元测试  

### 存储后端实现
- **local_chunk_manager.go** - 本地文件系统存储实现  
- **remote_chunk_manager.go** - 远程存储通用封装  
- **minio_object_storage.go** - MinIO 存储实现  
- **gcp_native_object_storage.go** - GCP 存储实现  
- **azure_object_storage.go** - Azure Blob 存储实现  

---

## 总结

`internal/storage` 是 Milvus 的核心数据管理模块，提供了跨本地与云平台的统一存储接口。它不仅支持基本的文件读写操作，还集成了复杂的数据结构管理、序列化/反序列化、Binlog 机制、统计信息收集、索引处理等功能，构成了 Milvus 数据持久化与高效查询的基础能力。