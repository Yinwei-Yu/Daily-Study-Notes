### 提交信息
- 提交: `fa77ea5120e4ab49ff854f5c7a906bafae1665e7`
- 作者: Yinwei Li
- 标题: 实现 GIS 表达式基于索引的评估，增强 RTreeIndex 功能

### 目标与范围
- 为 GIS 过滤表达式引入“索引段执行”路径：先用 R-Tree 粗筛，再用几何库做精确裁剪。
- 为 `RTreeIndex` 增强空值处理、序列化/反序列化索引附属数据（null offsets），完善查询能力。
- 增加端到端与多场景单测，验证表达式执行与 R-Tree 查询的一致性与正确性。

### 关键改动
- 表达式执行层
  - `internal/core/src/exec/expression/GISFunctionFilterExpr.cpp/.h`
    - 启用索引模式：`Eval()` 中当 `is_index_mode_` 为真时，调用新实现的 `EvalForIndexSegment()`
    - 新增 `EvalForIndexSegment()`：
      - 构造查询数据集（`OPERATOR_TYPE` + `MATCH_VALUE` WKB）
      - 逐 chunk 获取标量索引：先 `Query()` 得到 R-Tree 粗筛候选，再用原始几何做精确判断（`equals/touches/overlaps/crosses/contains/intersects/within`）
      - 封装批次结果为 `ColumnVector`（结果位图 + 有效位图），与批量窗口对齐

- R-Tree 索引层
  - `internal/core/src/index/RTreeIndex.cpp/.h`
    - Load/Build
      - 加载时解析并提取 `index_null_offset`（支持切片：`INDEX_FILE_SLICE_META` + 分片文件，或单文件），填充到 `null_offset_`
      - 总行数 = R-Tree 非空数 + null 数；`Count()`/`Load()` 都按此统计
      - 构建时单次遍历收集 null 行偏移；批量构建后设置 `total_num_rows_`
    - 序列化/上传
      - 实现 `Serialize()`：将 `null_offset_` 以 `index_null_offset` 写入 `BinarySet`，并 `Disassemble`
      - `Upload()` 同时登记磁盘索引文件和内存产物（null offsets），汇总到 `IndexStats`
    - 查询能力
      - 实现 `QueryCandidates(op, wkb, out)`：解析 query WKB 为 OGR 几何，调用 wrapper 的 `query_candidates`
      - 实现 `Query(dataset)`：根据 `OPERATOR_TYPE` 和 WKB 先做 R-Tree 粗筛，返回候选位图（精确裁剪在表达式层完成）
    - 空值查询
      - 实现 `IsNull()`/`IsNotNull()`：基于 `null_offset_` 构造位图；读写由 `folly::SharedMutex` 保护
    - 线程安全/状态
      - 新增成员：`folly::SharedMutexWritePriority mutex_`、`std::vector<size_t> null_offset_`
      - `Count()` 增加 `null_offset_.size()` 计入总量

  - `internal/core/src/index/RTreeIndexWrapper.cpp`
    - 批量数据流 `BulkLoadDataStream` 小改（去除多余注释）；核心逻辑不变（跳过 null 与无效 WKB）

- 单元测试
  - `internal/core/unittest/test_expr.cpp`
    - 调整期望：Within 示例中包含更多几何（0、1、3）
  - `internal/core/unittest/test_rtree_index.cpp`
    - 设置 `field_schema.data_type = Geometry` 以匹配几何索引场景
    - 移除已实现 API 的“未实现抛错”用例
    - 新增与增强：
      - 粗筛/精确查询用例：`Equals/Intersects/Within/Touches/Contains/Crosses/Overlaps` 场景
      - 端到端表达式执行：构造 `Sealed` 段，加载 R-Tree 索引，使用 `ExecPlanNodeVisitor` 执行 GIS 过滤表达式并断言
      - 保持已有构建/上传/加载、元数据、异常路径（仅 `.idx`/`.dat`/`.meta.json`、远端缺失等）用例的通过

### 行为变化
- **表达式执行**：当存在几何索引时，走“粗筛（R-Tree）+ 精确（GDAL/GEOS）”组合路径；提升性能并保证准确性。
- **空值处理**：R-Tree 索引现在能正确标识/统计空值行，`IsNull/IsNotNull/Count` 与其他标量索引保持一致。
- **索引产物**：`Upload()` 会包含 `index_null_offset` 内存产物，`Load()` 自动读取（含切片情况），保证统计一致性。
- **查询接口**：`Query()` 提供粗筛候选；精确裁剪上移至表达式层，使职责更清晰。

### 兼容性与风险
- 新增依赖/接口：
  - 使用 `common/Slice.h` 的 `INDEX_FILE_SLICE_META` 与 `Disassemble` 处理切片内存产物
  - `RTreeIndex` 的 `Serialize/IsNull/IsNotNull/Query/QueryCandidates` 均已实现，相关旧用例已调整
- 数据一致性：
  - 需要确保生成与上传 `index_null_offset` 的流程在分布式/切片场景下可靠；本次已覆盖切片元信息路径
- 表达式精确裁剪依赖 GDAL/GEOS 几何判断；不同版本边界处理可能存在差异，测试内已尽量规避不一致断言

### 建议与后续
- 为粗筛/精确裁剪增加统计指标（命中率、过滤率）与日志采样，便于性能分析
- 为 `Query()` 提供可选“粗+精”一体接口，用于非表达式栈复用
- 在 `.meta.json` 中记录 null 统计与切片布局，增强可观测性与排障能力

- 主要收益：GIS 表达式支持索引加速路径；R-Tree 索引具备完备的空值管理与可序列化能力；查询能力贯通“粗+精”，并通过丰富单测验证端到端正确性。