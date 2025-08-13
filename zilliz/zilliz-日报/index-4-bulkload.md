### 提交信息
- 提交: `f325feb943bf8f22327df40d32cced8d3ea98de7`
- 作者: Yinwei Li
- 标题: 为 R-Tree 引入批量构建（bulk load），增强几何处理与日志

### 改动目的
- 引入 R-Tree 批量加载能力以提升构建性能（基于 `SpatialIndex::RTree::createAndBulkLoadNewRTree`）。
- 统一并增强构建/加载时的健壮性和可追踪性（日志优化、变体降级策略）。
- 补充/调整单测覆盖批量加载、空值与异常 WKB 等边界场景。

### 主要改动
- R-Tree 索引实现
  - `internal/core/src/index/RTreeIndex.cpp`
    - 构建路径：
      - `BuildWithFieldData(...)` 改为默认使用批量加载：调用 `wrapper_->bulk_load_from_field_data(field_datas, schema_.nullable())`，不再逐条 `add_geometry`。
      - 初始化顺序与参数设置保持正确（先 `InitForBuildIndex()` 再设置参数）。
    - 日志微调：去除部分冗余 `LOG_INFO`。
- 包装器与批量加载
  - `internal/core/src/index/RTreeIndexWrapper.h/.cpp`
    - 新增 API：`bulk_load_from_field_data(...)`，参数为 `std::vector<std::shared_ptr<FieldDataBase>>` 与 `nullable`。
    - 内部实现：
      - 定义 `BulkLoadDataStream` 实现 `SpatialIndex::IDataStream`，从 FieldData 按行生成 R-Tree `Region` 数据项：
        - 跳过空值（根据 `nullable` 或 FieldData 自身的 nullable/valid 位图）。
        - 跳过无法解析的 WKB（OGR 解析失败）。
        - 行偏移（offset）与绝对偏移严格递增，保持与原始数据行对齐。
      - 通过 `createAndBulkLoadNewRTree` 一次性构建，记录 `index_id_`。
    - 动态插入路径改造：
      - 将 R-Tree 的创建移动至 `add_geometry(...)` 内的“惰性初始化”；若尚未创建，才 `createNewRTree`，以区分“批量构建”与“逐条插入”两种模式。
    - 变体处理：
      - `set_rtree_variant("QUADRATIC")` 现在会降级为 `RSTAR` 并打印警告（此前允许设置 `QUADRATIC`）。
    - 其他：
      - 引入 `common/FieldDataInterface.h`；保持析构默认实现。
- 单元测试与构建
  - `internal/core/unittest/CMakeLists.txt`
    - 确保 `test_rtree_index.cpp` 纳入测试编译目标。
  - `internal/core/unittest/test_rtree_index.cpp`
    - 新增/调整用例：
      - 新增 `Build_BulkLoad_Nulls_And_BadWKB`：验证批量加载下对空值与不合法 WKB 的过滤，加载后计数正确（本用例中 5 条输入，1 条故意截断为非法，期望计数为 4）。
      - `Build_ConfigAndMetaJson` 用例参数调整：`rv` 改为 `RSTAR`，`indexCapacity` 设为 32，`leafCapacity` 设为 64，并同步断言 `.meta.json` 中的持久化参数。
    - 小幅增加调试输出，便于追踪测试流程。

### 行为变化
- 构建默认采用“批量加载”，大幅提升构建性能；逐条插入路径仍保留于 `add_geometry`，用于 UT/非批量场景，但当前 `BuildWithFieldData` 不再回退到逐条插入。
- `QUADRATIC` 变体将被警告并降级为 `RSTAR`（兼容性更强）；非法变体仍抛出异常。
- 对空值与无效 WKB 的处理更加一致：批量与逐条插入两路径均跳过无效几何，同时保持行偏移一致。

### 风险与兼容性
- 变体变更：`QUADRATIC` 不再被原样使用，可能影响依赖该策略的场景，但采用降级策略、并输出告警。
- 批量路径下的异常处理：若批量构建异常将记录错误；当前未在 `BuildWithFieldData` 中自动回退到逐条插入（可在未来通过配置引入开关）。
- 依赖 OGR/WKB 解析质量；非法数据将被过滤计数，符合预期但需确保上游数据质量。

### 建议与后续
- 提供配置开关允许在特定场景禁用批量加载，回退为逐条插入以便更细粒度控制或调试。
- 在 `.meta.json` 中增加“构建模式”（bulk/dynamic）与统计信息（有效/跳过条目数），增强可观测性。
- 扩展更多几何类型与复杂对象的 UT，覆盖多批次 FieldData 混合场景。

- 本次提交引入 R-Tree 批量构建能力，默认启用以提升性能；完善了对空值与异常 WKB 的处理；对索引变体进行了更稳健的策略调整，并通过新增/调整单测覆盖关键路径，提升了稳定性与可维护性。