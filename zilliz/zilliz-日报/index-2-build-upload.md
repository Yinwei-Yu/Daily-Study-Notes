### 提交概览
- **主题**: 新增 R-Tree 标量索引的构建、上传与加载能力
- **目的**: 为 GIS 几何类型引入磁盘持久化的 R-Tree 索引管线（构建→上传→加载），并配套单测验证
- **影响范围**: `index` 栈（新增 `RTreeIndex`）、`RTreeIndexWrapper` 能力增强、索引参数元信息、标量索引类型枚举、单元测试

### 关键改动
- 新增
  - `internal/core/src/index/RTreeIndex.cpp`：R-Tree 索引实现，支持 Build/Upload/Load、UT 专用构建
  - `internal/core/src/index/RTreeIndex.h`：接口定义与 `ScalarIndex` 集成，显式实例化 `std::string`
  - `internal/core/unittest/test_rtree_index.cpp`：端到端单测（构建→上传→加载；仅文件名加载路径）
- 修改
  - `internal/core/src/index/RTreeIndexWrapper.cpp/.h`
    - 新增方法：`count()`、`set_rtree_variant()`、`set_fill_factor()`、`set_index_capacity()`、`set_leaf_capacity()`
    - 构造与生命周期日志增强；销毁顺序注释修正（去除多余空格）
  - `internal/core/src/index/Meta.h`
    - 新增配置键：`R_TREE_VARIANT_KEY("rv")`、`FILL_FACTOR_KEY("fillFactor")`、`INDEX_CAPACITY_KEY("indexCapacity")`、`LEAF_CAPACITY_KEY("leafCapacity")`
  - `internal/core/src/index/ScalarIndex.h`
    - 新增 `ScalarIndexType::RTREE` 及 `ToString` 映射

### 功能详述
- **构建 Build**
  - 暂存目录：`/tmp/milvus/rtree-index/<indexIdentifier>`，构建前校验空目录
  - 从 `insert_files` 读取原始 WKB 字段数据，逐条 `add_geometry` 写入 R-Tree
  - 支持索引参数：
    - `rv`: `RSTAR` | `QUADRATIC` | `LINEAR`
    - `fillFactor`(默认0.8)、`indexCapacity`(默认100)、`leafCapacity`(默认100)
  - 构建完成后调用 `finish()` 刷盘，记录总行数
- **上传 Upload**
  - 遍历构建目录，逐文件登记到 `DiskFileManager`（远端路径与文件大小）
  - 汇总生成 `IndexStats`（返回远端索引文件列表与总大小）
- **加载 Load**
  - 接收 `index_files`（支持全路径或仅文件名）；若仅文件名，则自动拼接远端前缀
  - 通过 `DiskFileManager` 将远端索引缓存到本地，再由 `RTreeIndexWrapper` 加载
- **统计 Count**
  - 基于 `SpatialIndex::IStatistics::getNumberOfData()` 返回数据条目数
- **生命周期与清理**
  - 析构时释放 `wrapper_` 并清理临时目录

### 使用示例
- 构建
```cpp
nlohmann::json cfg;
cfg["insert_files"] = std::vector<std::string>{/* raw data files */};
cfg["rv"] = "RSTAR";
cfg["fillFactor"] = 0.8;
cfg["indexCapacity"] = 100;
cfg["leafCapacity"] = 100;
rtree.Build(cfg);
auto stats = rtree.Upload({});
```
- 加载
```cpp
nlohmann::json cfg;
cfg["index_files"] = stats->GetIndexFiles(); // 或仅文件名列表
milvus::tracer::TraceContext ctx;
rtree.Load(ctx, cfg);
```

### 单测覆盖
- `Build_Upload_Load`
  - 使用 WKB POINT 造数，验证构建计数=2、上传产物非空、加载后计数=2
- `Load_WithFileNamesOnly`
  - 仅提供文件名（无目录）加载，验证自动前缀拼接与计数=2
- 涉及组件：`FileManagerContext`、`DiskFileManagerImpl`、`MemFileManagerImpl`、临时路径与本地存储

### 非功能与限制
- 尚未实现：`Serialize()`、`Load(BinarySet)`、`In/NotIn/Range/QueryCandidates/Query` 等通用检索接口（均返回未实现）
- 仅支持 `std::string` 字段（WKB 几何）
- 参数异常与目录状态异常均带断言/报错与日志

### 兼容性与风险
- 新增 `ScalarIndexType::RTREE` 不影响既有类型；`ToString` 已覆盖
- 构建过程依赖本地临时目录与 `DiskFileManager`/`MemFileManager` 协作
- `RTreeIndexWrapper` 新配置需与底层 `SpatialIndex` 行为保持一致，参数错误会报错

### 后续建议
- 实现 GIS 空间关系算子检索（`QueryCandidates/Query`）
- 支持序列化/反序列化 `BinarySet`
- 完善空值处理与更多几何类型 UT
- Long-running 构建与异常恢复策略、指标采集

- 主要变化点：新增 `RTreeIndex` 及端到端构建/上传/加载流程；`Wrapper` 增强；索引参数元信息与类型枚举扩展；单测覆盖两条关键路径。