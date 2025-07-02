在proxy分配了insert等数据插入操作后,会将数据分配给streamingcoord
streamingcoord负责消息的分配和广播
streamingnode接收到广播后,负责WAL和OSS的写入操作

## 架构

### 1. `internal/streamingcoord` - 流式协调器

**主要功能**：
- **消息分配**：负责将insert消息分配给合适的streamingnode
- **负载均衡**：通过balancer组件实现负载均衡
- **广播服务**：提供消息广播功能
- **资源管理**：管理streamingnode资源

**核心组件**：
- `server/balancer/` - 负载均衡器，负责分配channel到streamingnode
- `server/broadcaster/` - 广播器，负责消息广播
- `server/service/` - 服务层，提供gRPC服务接口
- `client/` - 客户端，用于与streamingnode通信

### 2. `internal/streamingnode` - 流式节点

**主要功能**：
- **WAL写入**：实现Write-Ahead Log，确保数据持久化
- **OSS写入**：将数据写入对象存储服务
- **消息处理**：处理produce和consume请求
- **数据同步**：通过flusher组件同步数据到存储

**核心组件**：
- `server/walmanager/` - WAL管理器，管理WAL生命周期
- `server/wal/` - WAL实现，支持多种存储后端（Kafka、Pulsar、RocketMQ等）
- `server/flusher/` - 数据刷新器，负责将WAL数据同步到OSS
- `server/service/` - 服务层，提供Handler和Manager服务

## 整个项目中与StreamingNode相关的文件列表

### 核心服务文件
1. `internal/streamingnode/server/server.go` - StreamingNode服务器主文件
2. `internal/streamingnode/server/builder.go` - 服务器构建器
3. `internal/streamingnode/server/service/handler.go` - 消息处理器服务
4. `internal/streamingnode/server/service/manager.go` - 管理器服务

### WAL相关文件
5. `internal/streamingnode/server/walmanager/manager_impl.go` - WAL管理器实现
6. `internal/streamingnode/server/walmanager/wal_lifetime.go` - WAL生命周期管理
7. `internal/streamingnode/server/walmanager/wal_state.go` - WAL状态管理
8. `internal/streamingnode/server/wal/wal.go` - WAL接口定义
9. `internal/streamingnode/server/wal/scanner.go` - WAL扫描器
10. `internal/streamingnode/server/wal/README.md` - WAL文档

### 消息处理文件
11. `internal/streamingnode/server/service/handler/producer/produce_server.go` - 生产者服务器
12. `internal/streamingnode/server/service/handler/consumer/` - 消费者相关文件
13. `internal/streamingnode/client/handler/` - 客户端处理器
14. `internal/streamingnode/client/manager/` - 客户端管理器

### 数据同步文件
15. `internal/streamingnode/server/flusher/flusherimpl/wal_flusher.go` - WAL刷新器
16. `internal/streamingnode/server/flusher/flusherimpl/flusher_components.go` - 刷新器组件
17. `internal/streamingnode/server/flusher/flusherimpl/msg_handler_impl.go` - 消息处理器实现

### StreamingCoord相关文件
18. `internal/streamingcoord/server/server.go` - StreamingCoord服务器
19. `internal/streamingcoord/server/builder.go` - 协调器构建器
20. `internal/streamingcoord/server/balancer/` - 负载均衡器
21. `internal/streamingcoord/server/broadcaster/` - 广播器
22. `internal/streamingcoord/server/service/` - 协调器服务
23. `internal/streamingcoord/client/client.go` - 协调器客户端

### 协议和工具文件
24. `pkg/proto/streamingpb/streaming.pb.go` - 流式协议定义
25. `pkg/proto/streamingpb/streaming_grpc.pb.go` - gRPC服务定义
26. `pkg/streaming/util/types/streaming_node.go` - 流式节点类型定义
27. `pkg/streaming/walimpls/` - WAL实现库

### 集成和启动文件
28. `cmd/components/streaming_node.go` - StreamingNode启动组件
29. `internal/distributed/streaming/streaming.go` - 分布式流式服务
30. `internal/proxy/task_insert_streaming.go` - Proxy中的流式插入任务

### 配置和工具文件
31. `pkg/util/paramtable/component_param.go` - 组件参数配置
32. `pkg/util/typeutil/type.go` - 类型定义
33. `pkg/metrics/streaming_service_metrics.go` - 流式服务指标
34. `internal/util/streamingutil/` - 流式工具库

## 最核心的文件

### 1. **WAL核心文件**
- `internal/streamingnode/server/walmanager/manager_impl.go` - WAL管理器核心实现
- `internal/streamingnode/server/wal/wal.go` - WAL接口定义
- `internal/streamingnode/server/wal/README.md` - WAL架构文档

### 2. **消息处理核心文件**
- `internal/streamingnode/server/service/handler/producer/produce_server.go` - 生产者核心逻辑
- `internal/streamingnode/server/service/handler.go` - 消息处理器服务

### 3. **数据同步核心文件**
- `internal/streamingnode/server/flusher/flusherimpl/wal_flusher.go` - WAL到OSS的同步逻辑
- `internal/streamingnode/server/flusher/flusherimpl/flusher_components.go` - 数据同步组件

### 4. **服务架构核心文件**
- `internal/streamingnode/server/server.go` - StreamingNode服务器主文件
- `internal/streamingcoord/server/server.go` - StreamingCoord服务器主文件
- `internal/proxy/task_insert_streaming.go` - Proxy中的流式插入入口

### 5. **协议和配置核心文件**
- `pkg/proto/streamingpb/streaming_grpc.pb.go` - gRPC服务协议
- `pkg/streaming/util/types/streaming_node.go` - 核心类型定义

这些文件构成了Milvus流式架构的核心，实现了从SDK insert调用到WAL写入和OSS写入的完整数据流。****