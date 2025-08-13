### 改动概览
- 目标：提升 R-Tree 索引的加载鲁棒性、构建流程正确性、索引元数据持久化能力，并显著完善单测覆盖面。
- 影响范围：`RTreeIndex` 加载/构建流程、`RTreeIndexWrapper` 生命周期管理与元数据写入、单测场景扩充。

### 关键代码改动
- `internal/core/src/index/RTreeIndex.cpp`
  - 加载流程优化：
    - 从缓存到本地的索引文件中，明确优先选择 `.dat` 或 `.idx` 的任一作为基路径，找不到则回退到 `.meta.json`，仍无则使用首个路径，避免误选非数据文件。
  - 构建流程修正：
    - 将 `InitForBuildIndex()` 前移，确保先创建 `wrapper_` 后再设置参数（`fillFactor/indexCapacity/leafCapacity/rv`），修复潜在空指针调用。
- `internal/core/src/index/RTreeIndexWrapper.cpp`
  - 新增元数据持久化：
    - 在 `finish()` 完成后写入 `<base>.meta.json`，内容包含 `index_id/variant/fill_factor/index_capacity/leaf_capacity/dimension`，便于后续可靠加载。
  - `finish()` 幂等保护：
    - 增加 `finished_` 标志，重复调用直接跳过；兼容 `rtree_ == nullptr` 的情况，防止二次释放。
  - 加载改进：
    - `load()` 时若存在 `.meta.json` 则读取其中 `index_id`，用以加载正确的 R-Tree，默认回退到 `0`；替代原固定 ID 的做法。
  - 其他：
    - 构造器保存 `index_id_`，并引入 `nlohmann/json` 和文件写入逻辑。
- `internal/core/src/index/RTreeIndexWrapper.h`
  - 新增成员：`finished_`（幂等控制）、`index_id_`（持久化/加载用）。
- `internal/core/unittest/test_rtree_index.cpp`
  - 单测大幅增强，覆盖如下场景：
    - 基本 E2E：构建→上传→加载；混合文件名/全路径加载；大数据量（10k）稳定性验证。
    - 异常路径：空输入、非法 WKB、仅 `.idx` / 仅 `.dat` / 仅 `.meta.json`、远端文件不存在等。
    - 配置与元数据：
      - 通过 `insert_files` 流程构建，校验 `.meta.json` 中参数持久化（`fill_factor/index_capacity/leaf_capacity/dimension`）。
      - 非法 `rv` 变体触发异常。
    - API 未实现的接口统一抛错校验（`Serialize/Load(BinarySet)/In/Range/...`）。
  - 新增测试工具方法：
    - `CreateWkbFromWkt`（WKT→WKB）、`WriteGeometryInsertFile`（通过 `InsertData` 序列化生成远端 parquet 并写入存储）。

### 行为与质量提升
- 加载更稳健：可识别多种文件组合，避免选到 `.meta.json` 导致加载失败；支持路径/文件名混用。
- 构建更正确：先初始化 `wrapper_` 再设置参数，解决参数设置时机问题。
- 元数据可追溯：构建完成写 `.meta.json`，加载时利用其中 `index_id`，减少因硬编码 ID 造成的失败。
- 幂等安全：`finish()` 可重复调用，不影响资源释放顺序与一致性。
- 单测覆盖广：覆盖 E2E、异常、参数、规模等关键路径，回归风险下降。

### 风险与兼容性
- 新增依赖：使用 `nlohmann::json` 持久化/读取元数据；需确保构建环境包含该依赖。
- 文件命名与后缀：远端索引文件可能带分片后缀（如 `_0`），本次加载逻辑通过扩展匹配和回退策略规避选择错误基路径的风险。
- 仍未实现：`Serialize/Load(BinarySet)/In/NotIn/Range/QueryCandidates/Query` 等 API 维持未实现并抛错，行为明确。

### 建议与后续
- 在加载路径中对 “仅 `.idx` 或仅 `.dat`” 明确报错信息，指导用户上传完整索引文件集。
- 扩展实际空间查询能力（`QueryCandidates/Query`），将 `.meta.json` 中参数用于运行期优化。
- 增加元数据版本与校验字段，增强跨版本兼容与告警能力。

- 本次提交主要提升了 R-Tree 索引的加载可靠性与构建流程正确性，新增 `.meta.json` 元数据落盘与幂等 finish 机制，同时补齐了大量单测覆盖，显著增强了可维护性与可诊断性。