一、提交概览
- 仓库：Yinwei-Yu/milvus
- 关联分支（参考用户提供链接）：feature-rtree-li_yinwei
- 提交哈希：5d396d06ba789a4bed976396b45a2c5779b8c640
- 作者/提交者：Yinwei Li <yinwei.li@zilliz.com>
- 提交时间（UTC）：2025-08-13T11:54:28Z
- 父提交：cba865f25a168faa2bdae65fcd3b6d35bc499069
- 提交标题与说明（原文摘要）：
  - 增强 Geometry 和 GISFunctionFilterExpr 的性能与缓存
  - 为 Geometry 类新增轻量构造器：可在不拷贝的情况下解析 WKB，优化短生命周期只读对象的内存使用
  - 在 PhyGISFunctionFilterExpr 中对粗粒度候选 bitmap 进行缓存，减少重复查询、提升评估效率
  - 更新评估逻辑以复用缓存结果，提升 index segment 评估过程中的性能

二、变更内容解读（基于提交信息的技术分析）
1) Geometry 类轻量构造器（WKB 零拷贝解析）
- 目标：
  - 支持针对 WKB（Well-Known Binary）格式的轻量级解析，避免将外部缓冲复制到内部存储，降低内存分配与拷贝开销，尤其适用于短生命周期、只读场景（如过滤评估管线中的临时对象）。
- 典型技术要点（推测实现方向）：
  - 新增一个接受“指针+长度”或等价轻量视图（如 span）的构造器或工厂函数；
  - 直接在外部缓冲上完成 WKB 头与几何体结构的解析（包含大小端标志、几何类型、坐标序列等）；
  - 对象内部仅保存对外部缓冲的只读引用，不拥有内存所有权，析构时不释放；
  - 明确对象只读属性，避免任何需要写入的 API 在这种模式下被调用；
  - 在文档/注释中强调外部缓冲在对象存活期必须保持有效（生命周期约束）。
- 性能收益预期：
  - 避免 memcpy/分配，降低 GC/allocator 压力；
  - 在批量过滤或算子链路中，短周期构造/销毁几何对象的成本显著下降。
- 潜在风险与注意事项：
  - 生命周期管理：一旦外部缓冲失效，对象内部引用将悬空；
  - 并发访问：若外部缓冲被其他线程写入或释放，可能导致数据竞争或崩溃；
  - 端序/对齐：WKB 的大小端标志与平台字节序处理要严谨，解析需显式转换；
  - API 约束：对外暴露的接口需清晰区分“持有型（owning）”与“观测型（view）”。

2) PhyGISFunctionFilterExpr 的粗候选 bitmap 缓存
- 目标：
  - 将 GIS 函数过滤的“粗粒度候选集”（coarse candidate）以 bitmap 形式缓存起来，避免对同一谓词/输入反复计算，提升整体过滤评估效率。
- 典型技术要点（推测实现方向）：
  - 缓存键：应至少包含“索引 segment 标识 + 具体 GIS 函数谓词（含参数、几何常量/WKB、空间参考等）+ 版本/时间戳（保障一致性）”；
  - 缓存值：bitmap（位图）或等价的压缩位集，表示可能命中的候选行；
  - 写入与复用：在首次评估得到粗候选 bitmap 后写入缓存；后续评估在缓存命中时直接复用；
  - 并发/一致性：需要线程安全的插入/查询，或按分区/segment 局部化缓存；当底层数据或索引版本变化时必须失效；
  - 内存控制：设置容量上限、LRU 等策略；提供统计指标（命中率、大小、逐 segment 占用）。
- 性能收益预期：
  - 避免重复的粗筛计算，尤其在相同查询/相似查询频繁执行或并发度高时效果显著；
  - 降低下游精筛（精确几何计算）需要处理的数据量。
- 潜在风险与注意事项：
  - 失效策略：segment compaction/merge、增删改数据、索引重建、TTL/多版本读取等都可能要求失效；
  - 键规范化：谓词等价性判定需要严格（常量几何体、坐标精度、SRID、缓冲区内容等），避免“看似相同”但键不同带来缓存碎片；
  - 决定缓存层级：是节点内进程级缓存、每 segment 缓存，还是算子实例级缓存；影响复用率与隔离性。

3) 评估逻辑更新以复用缓存结果
- 目标：
  - 在 index segment 的评估路径中优先尝试缓存命中，直接得到粗候选位图，将下游计算限制在更小的候选集合上；
- 典型流程（推测）：
  - 构造查询谓词键 → 查询缓存 → 命中则返回粗候选 → 进入精筛；
  - 未命中则按旧流程计算粗候选 → 存入缓存 → 进入精筛；
  - 在多谓词组合（AND/OR/NOT）场景中，可能将多个缓存位图进行按位运算以获得合成候选集；
- 质量要求：
  - 与旧逻辑功能一致、结果一致（等价性验证）；
  - 错误与降级路径清晰：缓存异常或 miss 不影响正确性，只影响性能。

三、兼容性与对外接口影响
- 源码/ABI 影响：
  - 如果 Geometry 对外新增了构造器或工厂方法，通常是源兼容的；若变更了已有构造语义或成员布局，需评估 ABI 稳定性（C++ 动态库场景）；
  - PhyGISFunctionFilterExpr 若新增可选缓存开关或策略配置，属于行为增强，默认开启需谨慎评估风险。
- 配置/参数：
  - 可能新增开关项：是否启用粗候选缓存、容量/过期策略、统计上报；
  - 建议默认值审慎设置，并提供运维侧可观测性。

四、性能与资源评估建议
- 基准测试场景：
  - 小/中/大三档数据规模；不同几何类型（点/线/面/多面等）；不同谓词（Contains/Within/Intersects/Touches 等）；
  - 冷启动（缓存空）与热身后（缓存热）分别测；
  - 单查询重复与多查询混合负载；
  - 单 segment 与多 segment 并行。
- 指标建议：
  - 端到端延迟、QPS、P95/P99；
  - 粗筛与精筛各阶段耗时拆分；
  - 缓存命中率、插入/淘汰次数、内存占用；
  - Geometry WKB 解析耗时与内存分配次数（对比旧实现）。
- 目标预期（示例性）：
  - 热缓存命中场景下，粗筛阶段耗时显著下降，整体查询延迟可见优化；
  - WKB 轻量解析减少内存 copy 与分配，CPU 指标与 GC/allocator 压力下降。

五、测试计划建议
- 单元测试
  - Geometry 轻量构造：
    - 多种 WKB（不同端序/类型/环/洞等）解析正确性；
    - 非法/截断 WKB 的错误处理；
    - 外部缓冲失效场景的防御（文档与断言）；
  - 缓存逻辑：
    - 键等价性判定（同构谓词/几何常量应命中）；
    - 并发访问安全（多线程 get/put）；
    - 失效策略（索引版本变更、数据变更、配置变化）；
- 集成/端到端测试
  - 在典型查询工作负载下，对比开启/关闭缓存的结果一致性与性能；
  - 多谓词组合与布尔逻辑下的位图合并正确性；
- 回归/稳定性
  - 长时间压力测试，观察内存占用、缓存膨胀、碎片化；
  - 崩溃注入/异常路径（缓存失败、内存不足）是否优雅降级。
  
六、运维与可观测性
- 日志与指标
  - 缓存命中/未命中/插入/淘汰计数与速率；
  - 逐 segment 缓存占用与上限控制；
  - 解析错误与异常路径埋点；
- 配置管理
  - 开关、容量、TTL/LRU 等可配置；
  - 动态调整能力与热更新可行性。

七、风险清单与缓解
- 生命周期风险（Geometry 轻量构造）
  - 缓解：清晰的注释/文档；在 debug 模式下增加断言；对不安全用法进行封装限制；
- 一致性风险（缓存）
  - 缓解：缓存键包含索引版本/拉链时间；对数据变更事件订阅进行失效；提供强制失效操作；
- 行为回退
  - 提供配置开关以便在故障时回退到“无缓存/拷贝解析”模式；
- 性能回退监控
  - 若命中率低或带来额外内存压力，及时报警与策略调整。

八、文档与示例
- 开发者文档
  - Geometry 新增构造器的使用示例、限制与注意事项；
  - PhyGISFunctionFilterExpr 缓存设计、键规范、失效策略与指标；
- 用户/运维文档
  - 新增配置项说明与调优指南；
  - 典型排障（缓存不命中、内存超限、性能异常）。

九、后续工作与建议
- 在获取到完整 diff 后，补充：
  - 逐文件改动列表（新增/修改/删除）、函数级差异、关键代码片段；
  - 接口签名变化与调用方影响矩阵；
- 引入基准测试脚本与公用数据集，纳入 CI 性能基线对比；
- 将缓存指标纳入监控看板，并设置合理阈值与告警；
- 评估与既有筛选/索引管线的进一步融合（例如与 R-Tree/索引粗筛阶段的协同）。

十、当前报告的已知限制
- 本次未能获取提交的具体文件差异与行级改动，因而“逐文件/逐行”的解读部分以设计与实现推断为主，未包含代码片段与精确行变更。待能访问到完整变更内容后，可将本报告升级为精准的代码级审阅版本。

——
结论：该提交聚焦两点性能优化与工程化增强——Geometry 的 WKB 轻量化解析（降内存与 CPU 开销）与 GIS 粗候选位图缓存（降重复计算）。这两项在高并发、重复查询和复杂几何谓词场景下将带来显著性能收益。需重点关注生命周期管理、缓存一致性与失效策略、并发安全及观测指标完善，以保障收益可控且稳定落地。


# **gemini**

### 提交信息概述

此提交旨在通过以下两个主要方面的改进来增强 `Geometry` 类和 `GISFunctionFilterExpr` 的性能和缓存机制：

1.  **为 `Geometry` 类添加轻量级构造函数**: 引入一个新构造函数，用于在解析 WKB (Well-Known Binary) 数据时避免不必要的内存复制。这对于生命周期短、只读的 `Geometry` 对象特别有用，例如在空间谓词（如 `equals`, `intersects` 等）评估过程中。
2.  **在 `PhyGISFunctionFilterExpr` 中实现粗略候选位图缓存**: 对 `PhyGISFunctionFilterExpr` 的评估逻辑进行了修改，以缓存粗略的候选位图，从而减少重复的 R-Tree 查询并提高评估效率。评估逻辑已更新为利用这些缓存结果，从而提升索引段评估期间的性能。

### 文件更改详细内容

#### 1. `internal/core/src/common/Geometry.h`

**更改摘要**: 为 `Geometry` 类添加了一个新的构造函数。

**详细更改**:

*   **新增构造函数**:
    ```cpp
    32,6 +32,21 @@ class Geometry {
        to_wkb_internal();
    }

    // lightweight constructor: parse wkb but **do not** copy it back to internal buffer,
    // which avoids an extra malloc + memcpy. Suitable for short-lived, read-only objects
    // where we only need spatial predicates (equals/intersects/...).
    explicit Geometry(const void* wkb, size_t size, bool copy_wkb) {
        OGRGeometry* geometry = nullptr;
        OGRGeometryFactory::createFromWkb(wkb, nullptr, &geometry, size);
        AssertInfo(geometry != nullptr,
                   "failed to construct geometry from wkb data");
        geometry_.reset(geometry);
        size_ = size;
        if (copy_wkb) {
            to_wkb_internal();
        }
    }

    explicit Geometry(const char* wkt) {
        OGRGeometry* geometry = nullptr;
        OGRGeometryFactory::createFromWkt(wkt, nullptr, &geometry);
    ```
    *   这个新的构造函数接受 `const void* wkb` (WKB数据指针), `size_t size` (WKB数据大小) 和 `bool copy_wkb` (是否复制WKB数据到内部缓冲区)。
    *   当 `copy_wkb` 为 `false` 时，它避免了额外的 `malloc` 和 `memcpy` 操作，从而优化了内存使用，适用于那些只用于空间谓词判断、不需要持久化 WKB 数据的临时 `Geometry` 对象。

#### 2. `internal/core/src/exec/expression/GISFunctionFilterExpr.cpp`

**更改摘要**: 引入了粗略候选位图的缓存机制，并修改了 `EvalForIndexSegment` 函数的评估逻辑以利用缓存。

**详细更改**:

*   **新增预取和缓存逻辑**:
    ```cpp
    156,119 +156,145 @@ PhyGISFunctionFilterExpr::EvalForIndexSegment() {
         ds->Set(milvus::index::OPERATOR_TYPE, expr_->op_);
         ds->Set(milvus::index::MATCH_VALUE, expr_->geometry_.to_wkb_string());

    +    /* ------------------------------------------------------------------
    +     * Prefetch: if coarse results are not cached yet, run a single R-Tree
    +     * query for all index chunks and cache their coarse bitmaps.
    +     * ------------------------------------------------------------------*/
    +    if (!coarse_cached_) {
    +        coarse_cache_.resize(num_index_chunk_);
    +        coarse_valid_cache_.resize(num_index_chunk_);
    +
    +        for (size_t cid = 0; cid < num_index_chunk_; ++cid) {
    +            const Index& idx_ref =
    +                segment_->chunk_scalar_index<std::string>(field_id_, cid);
    +            auto* idx_ptr = const_cast<Index*>(&idx_ref);
    +
    +            auto coarse = idx_ptr->Query(ds);
    +            coarse_cache_[cid] = std::move(coarse);
    +
    +            auto valid = idx_ptr->IsNotNull();
    +            coarse_valid_cache_[cid] = std::move(valid);
    +        }
    +        coarse_cached_ = true;
    +    }
    +
         TargetBitmap batch_result;
         TargetBitmap batch_valid;
         int processed_rows = 0;

    -    for (size_t i = current_index_chunk_; i < num_index_chunk_; i++) {
    -        // 1) fetch index for this chunk and run coarse query
    -        const Index& index_ref =
    -            segment_->chunk_scalar_index<std::string>(field_id_, i);
    -        Index* index_ptr = const_cast<Index*>(&index_ref);
    -        auto coarse = index_ptr->Query(ds);
    -        auto chunk_valid = index_ptr->IsNotNull();
    +    for (size_t i = current_index_chunk_; i < num_index_chunk_; ++i) {
    +        // 1) Build and cache refined bitmap for this chunk (coarse + exact)
    +        if (cached_index_chunk_id_ != static_cast<int64_t>(i)) {
    +            // Reuse segment-level coarse cache directly
    +            auto& coarse = coarse_cache_[i];
    +            auto& chunk_valid = coarse_valid_cache_[i];
    ```
    *   在 `EvalForIndexSegment` 函数的开始部分，增加了一个预取逻辑。如果 `coarse_cached_` 为 `false`（表示粗略结果尚未缓存），则会遍历所有的索引块 (`num_index_chunk_`)，对每个块执行一次 R-Tree 查询 (`idx_ptr->Query(ds)`) 来获取粗略的候选位图 (`coarse`) 和非空位图 (`valid`)。
    *   这些位图被存储在 `coarse_cache_` 和 `coarse_valid_cache_` 中，并将 `coarse_cached_` 设置为 `true`，以避免后续重复查询。

*   **精确细化逻辑的修改**:
    ```cpp
    195,83 +221,80 @@ class PhyGISFunctionFilterExpr : public SegmentExpr {
                         if (ok) {
                             refined.set(pos);
                         }
                     }
                 }
    -        } else {
    -            // growing: std::string values
    -            auto span = segment_->chunk_data<std::string>(field_id_, i);
    -            for (size_t pos = 0; pos < coarse.size(); ++pos) {
    -                if (!coarse[pos]) {
    -                    continue;
    -                }
    -                const auto& wkb = span[pos];
    -                Geometry left(wkb.data(), wkb.size());
    -                bool ok = false;
    -                switch (expr_->op_) {
    -                    case proto::plan::GISFunctionFilterExpr_GISOp_Equals:
    -                        ok = left.equals(expr_->geometry_);
    -                        break;
    -                    case proto::plan::GISFunctionFilterExpr_GISOp_Touches:
    -                        ok = left.touches(expr_->geometry_);
    -                        break;
    -                    case proto::plan::GISFunctionFilterExpr_GISOp_Overlaps:
    -                        ok = left.overlaps(expr_->geometry_);
    -                        break;
    -                    case proto::plan::GISFunctionFilterExpr_GISOp_Crosses:
    -                        ok = left.crosses(expr_->geometry_);
    -                        break;
    -                    case proto::plan::GISFunctionFilterExpr_GISOp_Contains:
    -                        ok = left.contains(expr_->geometry_);
    -                        break;
    -                    case proto::plan::GISFunctionFilterExpr_GISOp_Intersects:
    -                        ok = left.intersects(expr_->geometry_);
    -                        break;
    -                    case proto::plan::GISFunctionFilterExpr_GISOp_Within:
    -                        ok = left.within(expr_->geometry_);
    -                        break;
    -                    default:\
    -                        PanicInfo(\
    -                            NotImplemented, "unknown GIS op : {}", expr_->op_);
    +            } else {  // Growing segment
    +                auto span = segment_->chunk_data<std::string>(field_id_, i);
    +                for (size_t pos = 0; pos < coarse.size(); ++pos) {
    +                    if (!coarse[pos])
    +                        continue;
    +
    +                    const auto& wkb = span[pos];
    +                    Geometry left(wkb.data(), wkb.size(), false);\
    +                    bool ok = false;
    +                    switch (expr_->op_) {
    +                        case proto::plan::GISFunctionFilterExpr_GISOp_Equals:\
    +                            ok = left.equals(expr_->geometry_);\
    +                            break;\
    +                        case proto::plan::GISFunctionFilterExpr_GISOp_Touches:\
    +                            ok = left.touches(expr_->geometry_);\
    +                            break;\
    +                        case proto::plan::GISFunctionFilterExpr_GISOp_Overlaps:\
    +                            ok = left.overlaps(expr_->geometry_);\
    +                            break;\
    +                        case proto::plan::GISFunctionFilterExpr_GISOp_Crosses:\
    +                            ok = left.crosses(expr_->geometry_);\
    +                            break;\
    +                        case proto::plan::GISFunctionFilterExpr_GISOp_Contains:\
    +                            ok = left.contains(expr_->geometry_);\
    +                            break;\
    +                        case proto::plan::\
    +                            GISFunctionFilterExpr_GISOp_Intersects:\
    +                            ok = left.intersects(expr_->geometry_);\
    +                            break;\
    +                        case proto::plan::GISFunctionFilterExpr_GISOp_Within:\
    +                            ok = left.within(expr_->geometry_);\
    +                            break;\
    +                        default:\
    +                            PanicInfo(NotImplemented,\
    +                                      "unknown GIS op : {}",\
    +                                      expr_->op_);
+                    }\
+                    if (ok) {\
+                        refined.set(pos);\
+                    }
                 }\
-                if (ok) {\
-                    refined.set(pos);\
+            }\
+
+            // Cache refined result for reuse by subsequent batches
+            cached_index_chunk_id_ = i;\
+            cached_index_chunk_res_ = std::move(refined);\
+            // No need to copy valid bitmap into member; use coarse_valid_cache_[i] directly l
+ater
+        }\
+
+        // 2) Append this chunk\'s cached results into current batch window
+        const auto& chunk_valid_ref = coarse_valid_cache_[i];\
+
+        auto size = ProcessIndexOneChunk(batch_result,\
+                                         batch_valid,\
+                                         i,\
+                                         cached_index_chunk_res_,\
+                                         chunk_valid_ref,\
+                                         processed_rows);\
+
+        if (processed_rows + size >= batch_size_) {\
+            current_index_chunk_ = i;\
                 }
             }\
    -
    -        // 3) append this chunk\'s refined result to batch result according to batch window
    -        auto data_pos =
    -            i == current_index_chunk_ ? current_index_chunk_pos_ : 0;
    -        auto size = std::min(
    -            std::min(size_per_chunk_ - data_pos, batch_size_ - processed_rows),
    -            int64_t(refined.size()));
    -
    -        batch_result.append(refined, data_pos, size);
    -        batch_valid.append(chunk_valid, data_pos, size);
    -
    -        if (processed_rows + size >= batch_size_) {
    -            current_index_chunk_ = i;
    ```
    *   在精确细化部分，创建 `Geometry` 对象时，现在使用新的轻量级构造函数 `Geometry left(wkb_view.data(), wkb_view.size(), false);` 或 `Geometry left(wkb.data(), wkb.size(), false);`，其中 `copy_wkb` 参数设置为 `false`，以避免不必要的内存复制。
    *   在处理完一个索引块的精确细化后，会将细化后的结果 (`refined`) 缓存到 `cached_index_chunk_res_` 中，并将 `cached_index_chunk_id_` 更新为当前块的 ID。
    *   后续批处理会直接使用缓存的细化结果 (`cached_index_chunk_res_`) 和粗略非空位图 (`coarse_valid_cache_[i]`)，通过 `ProcessIndexOneChunk` 函数将结果追加到 `batch_result` 和 `batch_valid` 中。这避免了对已处理块的重复精确细化。
    *   移除了之前直接在循环内部进行 `Query` 和 `IsNotNull` 的调用，而是改为从缓存中获取。

#### 3. `internal/core/src/exec/expression/GISFunctionFilterExpr.h`

**更改摘要**: 为 `PhyGISFunctionFilterExpr` 类添加了新的成员变量，用于支持粗略候选位图的缓存。

**详细更改**:

*   **新增缓存相关成员变量**:
    ```cpp
    56,6 +56,17 @@ class PhyGISFunctionFilterExpr : public SegmentExpr {

      private:
         std::shared_ptr<const milvus::expr::GISFunctionFilterExpr> expr_;
    +
    +    /*
    +     * Segment-level cache: run a single R-Tree Query for all index chunks to
    +     * obtain coarse candidate bitmaps. Subsequent batches reuse these cached
    +     * results to avoid repeated ScalarIndex::Query calls per chunk.
    +     */
    +    bool coarse_cached_ =
    +        false;  // whether coarse results have been prefetched once
    +    std::vector<TargetBitmap>
    +        coarse_cache_;  // per-chunk coarse candidate bitmap
    +    std::vector<TargetBitmap> coarse_valid_cache_;  // per-chunk not-null bitmap
     };
     }  //namespace exec
     }  // namespace milvus
    ```
    *   `coarse_cached_` (类型: `bool`, 默认值: `false`): 标志位，指示粗略结果是否已经被预取和缓存。
    *   `coarse_cache_` (类型: `std::vector<TargetBitmap>`): 存储每个索引块的粗略候选位图。
    *   `coarse_valid_cache_` (类型: `std::vector<TargetBitmap>`): 存储每个索引块的非空位图。

#### 4. `internal/core/src/index/RTreeIndexWrapper.cpp`

**更改摘要**: 移除了 `RTreeIndexWrapper::query_candidates` 函数中对 `GISFunctionFilterExpr_GISOp_Contains` 操作的特殊处理。

**详细更改**:

*   **移除 `GISOp_Contains` 的特殊处理**:
    ```cpp
    379,9 +379,6 @@ RTreeIndexWrapper::query_candidates(proto::plan::GISFunctionFilterExpr_GISO\np op,\

         // Perform query based on operation type
         switch (op) {
    -        case proto::plan::GISFunctionFilterExpr_GISOp_Contains:
    -            rtree_->containsWhatQuery(query_region, visitor);
    -            break;
             default:
                 // For all GIS operations, we use intersection query as coarse filtering
                 // The exact geometric relationship will be checked in the refinement phase
    ```
    *   原先对 `GISFunctionFilterExpr_GISOp_Contains` 操作有单独的 `rtree_->containsWhatQuery` 调用。
    *   现在，所有 GIS 操作都统一使用 `intersection query` 作为粗略过滤，精确的几何关系将在后续的细化阶段进行检查。这简化了 `RTreeIndexWrapper` 的逻辑，并将精确细化完全交给 `GISFunctionFilterExpr` 处理，与缓存机制更好地集成。

### 结论和影响

这个提交对 Milvus 的地理空间索引查询性能带来了显著的提升。

*   **内存优化**: `Geometry` 类中的轻量级构造函数减少了不必要的内存复制，特别是在处理大量短生命周期几何对象时。
*   **性能提升**: `GISFunctionFilterExpr` 中的粗略候选位图缓存机制，避免了对相同索引块的重复 R-Tree 查询。这使得在索引段评估过程中，可以更快地获取粗略的候选集，然后进行精确的几何关系判断，从而大大提高了查询效率。
*   **代码简化**: `RTreeIndexWrapper` 的修改使得所有 GIS 操作的粗略过滤逻辑更加统一，将精确细化完全委派给 `GISFunctionFilterExpr`，使得整体架构更清晰。

这些改进将有助于在处理地理空间数据时，特别是在需要频繁进行空间查询的场景下，提供更优异的性能表现。