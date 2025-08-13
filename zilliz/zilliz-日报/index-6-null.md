### 提交信息
- 提交: `694506ad997d894c751dc2f51f73bb04abbf2b38`
- 作者: Yinwei Li
- 标题: 重构 RTreeIndex 与 FieldData 处理

### 目标
- 统一空值统计口径，去掉 Wrapper 层的冗余空值计数，改由 `RTreeIndex` 基于 `null_offset_` 计算总行数。
- 修复/增强几何字段的可空数据填充与有效位图同步。
- 索引参数读取改为强类型获取，减少字符串解析带来的错误。
- 清理无用字段/依赖，简化实现与维护成本。

### 关键改动
- RTreeIndex
  - `internal/core/src/index/RTreeIndex.cpp/.h`
    - 总行数计算改为：`wrapper_->count() + null_offset_.size()`（移除对 `wrapper_->null_count()` 的依赖）。
    - `Build()/Load()` 完成后，`total_num_rows_` 按上述口径更新。
- RTreeIndexWrapper
  - `internal/core/src/index/RTreeIndexWrapper.cpp/.h`
    - 删除空值计数逻辑：移除 `null_rows_` 字段、`null_count()` 接口。
    - 批量加载数据流中不再累加空值计数；`finish()` 写出的 `.meta.json` 不再包含 `null_rows`；`load()` 不再读取该字段。
- FieldData（几何）
  - `internal/core/src/common/FieldDataInterface.h`
    - `FieldDataGeometryImpl::FillFieldData`：当字段可空且存在有效位图时，将源数组的 null bitmap 精确拷贝到内部 `valid_data_`，保证几何字段的可空性状态正确下沉。
  - `internal/core/src/common/FieldData.cpp`
    - 修复 `FillFieldData(...)` 数组分支的计数传参，始终使用 `element_count`，避免此前误用导致的数据长度错误。
- 索引参数读取
  - `internal/core/src/index/Utils.cpp`
    - `GetFillFactorFromConfig` 改为直接获取 `double`；`GetIndexCapacityFromConfig`/`GetLeafCapacityFromConfig` 改为直接获取 `uint32_t`，不再做字符串解析。
- 单测与依赖
  - `internal/core/unittest/test_rtree_index.cpp`
    - 构建用例中显式设置 `fillFactor/indexCapacity/leafCapacity/rv`，与强类型读取保持一致。
  - `pkg/go.sum`
    - 移除未使用的 Go 依赖条目（如 `goid`、`go-deadlock`）。

### 行为变化
- 统一空值统计：总行数与 `IsNull/IsNotNull` 依据 `null_offset_`，与其他标量索引保持一致且可序列化持久化。
- `.meta.json` 不再写入/读取 `null_rows`；向前读取旧文件时该字段被忽略，不影响加载。
- 几何字段可空数据在填充后，内部有效位图与来源一致，避免后续查询/索引阶段的可空性偏差。
- 索引参数配置读取更健壮（类型安全），减少运行期解析错误。

### 风险与兼容性
- 若历史代码/脚本依赖 `null_rows` 元数据字段，将不再可用；但索引加载与行数统计不受影响。
- 需要确保上游配置以正确类型传入（double/uint32_t），否则断言会触发。

### 建议与后续
- 在索引构建/加载日志中增加 `null_offset_` 规模与排序校验，便于排障。
- 对历史 `.meta.json` 中的 `null_rows` 提示废弃信息（仅日志）。

- 本次提交聚焦于空值统计口径统一、几何字段可空处理修复、索引配置强类型化，以及清理无用实现，提升了稳定性与可维护性。