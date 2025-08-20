好的，这是一份关于 Milvus 中地理空间索引查询崩溃问题的详细技术报告。

---

## 技术报告：Milvus 地理空间索引查询崩溃问题分析与解决方案

**报告日期：** 2025年8月18日

### 1. Bug发现过程

#### 1.1. 发现情境描述

在对 Milvus 进行地理空间索引（RTREE）功能测试时，发现了一个导致查询服务崩溃的Bug。

**操作步骤：**

1.  **准备数据：** 使用 `tests/python_client/geo-test/parallel/data_preparation.py` 脚本，创建名为 `geo_index_benchmark` 的 Milvus Collection，并插入 1,000,000 条地理空间（POLYGON）数据。同时，脚本会为 `geo` 字段构建 RTREE 索引，并为 `vec` 字段构建 IVF\_FLAT 索引。为了后续调试，我们还在 `data_preparation.py` 中为 `id` 字段（后改为 `int_field`）创建了 `INVERTED` 索引。

    *   **环境：**
        *   操作系统：Linux (Ubuntu 5.15.0-151-generic)
        *   Milvus 版本：基于提供的代码路径 `/home/zilliz/milvus`，推测为 Milvus 的 C++ 核心组件。
        *   Python 客户端：`pymilvus`。

2.  **首次查询：** 在 `data_preparation.py` 脚本执行完毕后，立即运行 `tests/python_client/geo-test/parallel/query_benchmark.py` 脚本进行地理空间查询。此时，所有查询均能正常执行，系统运行稳定。

3.  **等待与再次查询：** 保持 Milvus 服务运行。等待约 2 分钟后（这个时间段足以让 Milvus 后台的 Compaction 机制触发并可能重构 Segment 及其索引），再次运行 `query_benchmark.py` 脚本。

4.  **Bug复现：** 此次查询执行时，系统立即报告错误，Milvus QueryNode 服务崩溃。该Bug复现频率稳定，只要触发 Compaction 并重新加载索引，问题就会出现。

#### 1.2. 发现时的现象及错误信息

当 Bug 发生时，`query_benchmark.py` 脚本在执行 `client.query` 操作时捕获到 `MilvusException`，并打印出详细的错误信息。

**错误日志片段：**

```log
2025-08-18 02:25:15,584 [ERROR][handler]: RPC error: [query], <MilvusException: (code=65535, message=fail to Query on QueryNode 2: worker(2) query failed: Operator::GetOutput failed for [Operator:PhyFilterBitsNode, plan node id: 654] : Assert "(bitset.size() == need_process_rows_)"  at /home/zilliz/milvus/internal/core/src/exec/operator/FilterBitsNode.cpp:98
)>, <Time:{'RPC start': '2025-08-18 02:25:09.056768', 'RPC error': '2025-08-18 02:25:15.584532'}> (decorators.py:140)
2025-08-18 02:25:15,584 [ERROR][query]: Failed to query collection: geo_index_benchmark (milvus_client.py:482)
Error during query execution: <MilvusException: (code=65535, message=fail to Query on QueryNode 2: worker(2) query failed: Operator::GetOutput failed for [Operator:PhyFilterBitsNode, plan node id: 654] : Assert "(bitset.size() == need_process_rows_)"  at /home/zilliz/milvus/internal/core/src/exec/operator/FilterBitsNode.cpp:98
```

**Milvus 服务端日志关键行（首次发现问题时）：**

```log
I20250818 03:05:27.666975 2950778 FilterBitsNode.cpp:98] [SERVER][GetOutput][MILVUS_FUTURE_C][]LiYinwei:bitset size: 1007616
I20250818 03:05:27.667153 2950778 FilterBitsNode.cpp:99] [SERVER][GetOutput][MILVUS_FUTURE_C][]LiYinwei:need_process_rows_: 1000000
```

这表明 `FilterBitsNode.cpp` 的断言 `Assert "(bitset.size() == need_process_rows_)"` 失败，其中 `bitset.size()` 为 1,007,616，而 `need_process_rows_` 为 1,000,000。

### 2. 调试过程

为了定位导致 `FilterBitsNode` 断言失败的根本原因，我们采取了以下调试步骤：

#### 2.1. 初步定位与 `bitset.size()` 及 `need_process_rows_` 来源分析

*   **目的：** 了解断言中两个变量的含义和来源。
*   **步骤：**
    *   **`need_process_rows_`：** 通过代码搜索发现，`need_process_rows_` 在 `PhyFilterBitsNode` 构造函数中初始化自 `query_context_->get_active_count()`。进一步追溯 `QueryContext::get_active_count()`，确认它返回的是当前 Segment 的逻辑活跃行数。在我们的例子中，这个值是数据插入的总行数 1,000,000。
    *   **`bitset.size()`：** 定位到 `FilterBitsNode::GetOutput()` 函数，发现 `bitset` 是通过循环调用 `exprs_->Eval()` 来累积的。每次 `Eval` 返回一个 `ColumnVector`，其内部的 `TargetBitmap` (`col_vec`) 的大小会被 `append` 到 `bitset` 中。
*   **结果：** 问题在于 `FilterBitsNode` 期望处理 1,000,000 行，但从 `Eval` 链路上累积得到的位图最终大小却为 1,007,616。

#### 2.2. 根因探索 - 阶段1：R-Tree 索引计数膨胀假设

*   **目的：** 怀疑 R-Tree 索引在构建时，由于 Compaction 导致的数据源问题，其内部维护的条目数量（`rtree_->getStatistics()->getNumberOfData()`）超过了 Segment 的实际逻辑行数。
*   **步骤：**
    *   检查 `RTreeIndex<T>::Count()` 实现，确认它依赖 `wrapper_->count()`，而 `wrapper_->count()` 依赖 `rtree_->getStatistics()->getNumberOfData()`。
    *   分析 `RTreeIndexWrapper::bulk_load_from_field_data` 和 `BulkLoadDataStream::getNext()`，猜测 `absolute_offset_` 在跳过无效数据（例如，已删除行、WKB 解析失败的行）时仍然递增，导致 `libspatialindex` 索引了实际上无效或已删除的数据，从而使得 `getNumberOfData()` 返回一个虚高的值。
*   **结果：** 这是最初的猜测，认为 `1,007,616` 是 R-Tree 内部的错误计数。

#### 2.3. 根因探索 - 阶段2：精确化问题定位到 `EvalForIndexSegment` 内部批处理

*   **目的：** 验证 R-Tree 索引计数假设，并进一步缩小问题范围。
*   **步骤：**
    *   在 `RTreeIndex.cpp` 中 `RTreeIndex<T>::Query(const DatasetPtr& dataset)` 函数的 `TargetBitmap res(this->Count())` 之后，添加日志打印 `res.size()`。
        ```cpp
        // internal/core/src/index/RTreeIndex.cpp
        500|    LOG_INFO("LiYinwei:Query res.size(): {}", res.size());
        ```
    *   重新运行测试，观察日志。
*   **结果：** 日志显示 `I20250818 06:41:30.432663 ... RTreeIndex.cpp:500] ...LiYinwei:Query res.size(): 1000000`。
    **这推翻了阶段1的假设！** 这表明 R-Tree 索引在查询阶段返回的粗筛位图大小是正确的 1,000,000，并没有膨胀。问题必然发生在 `GISFunctionFilterExpr::EvalForIndexSegment()` 内部，以及它如何向 `FilterBitsNode` 传递数据。
*   **进一步调试 `EvalForIndexSegment`：**
    *   在 `GISFunctionFilterExpr::EvalForIndexSegment()` 的 `for` 循环内部，以及 `ProcessIndexOneChunk` 函数内部的关键位置添加日志，用于追踪 `batch_result` 的大小、`processed_rows`、以及 `size` 的变化。
        *   在 `Expr.h` 的 `ProcessIndexOneChunk` 中：
            ```cpp
            // internal/core/src/exec/expression/Expr.h
            816|        LOG_INFO("LiYinwei:ProcessIndexOneChunk result size: {}", result.size());
            817|        LOG_INFO("LiYinwei:ProcessIndexOneChunk valid_result size: {}", valid_result.size());
            ```
        *   在 `GISFunctionFilterExpr.cpp` 的 `EvalForIndexSegment` 循环前后：
            ```cpp
            // internal/core/src/exec/expression/GISFunctionFilterExpr.cpp
            331|    LOG_INFO("LiYinwei:EvalForIndexSegment loop - i: {}, processed_rows_before_append: {}, size_from_ProcessIndexOneChunk: {}, batch_size_: {}", i, processed_rows, size, batch_size_);
            346|    LOG_INFO("LiYinwei:EvalForIndexSegment loop - processed_rows_after_append_and_check: {}", processed_rows + size);
            342|    LOG_INFO("LiYinwei:EvalForIndexSegment returning batch_result size: {}", batch_result.size());
            ```
    *   **分析日志（以一次 `EvalForIndexSegment` 调用为例）：**
        ```log
        I20250818 07:21:39.892235 ... Expr.h:816] ...LiYinwei:ProcessIndexOneChunk result size: 576  // 第一次 ProcessIndexOneChunk 调用
        I20250818 07:21:39.892307 ... GISFunctionFilterExpr.cpp:331] ...EvalForIndexSegment loop - i: 0, processed_rows_before_append: 0, size_from_ProcessIndexOneChunk: 576, batch_size_: 8192
        I20250818 07:21:39.892333 ... GISFunctionFilterExpr.cpp:346] ...EvalForIndexSegment loop - processed_rows_after_append_and_check: 576
        // 注意：这里 EvalForIndexSegment 没有立即返回，因为 576 + 576 < 8192，for 循环继续
        I20250818 07:21:39.964669 ... Expr.h:816] ...LiYinwei:ProcessIndexOneChunk result size: 8192  // 第二次 ProcessIndexOneChunk 调用
        I20250818 07:21:39.964831 ... GISFunctionFilterExpr.cpp:331] ...EvalForIndexSegment loop - i: 1, processed_rows_before_append: 576, size_from_ProcessIndexOneChunk: 7616, batch_size_: 8192
        I20250818 07:21:39.964864 ... GISFunctionFilterExpr.cpp:353] ...EvalForIndexSegment returning batch_result size: 8192 // EvalForIndexSegment 返回
        I20250818 07:21:39.965375 ... FilterBitsNode.cpp:98] ...LiYinwei:bitset size: 1007616 // FilterBitsNode 最终收到导致崩溃
        ```
    *   **分析结果：** 确认 `GISFunctionFilterExpr::EvalForIndexSegment()` 在一次调用中，其内部的 `for` 循环处理了多个 `chunk` (`i=0` 和 `i=1`)。当 `FilterBitsNode` 期望一个较小的尾部批次（例如 `576` 行）时，`EvalForIndexSegment` 内部的 `if (processed_rows + size >= batch_size_)` 条件未能精确地在 `576` 行处停止。由于 `processed_rows + size` 仍小于 `batch_size_`，循环继续，并从下一个内部 `chunk` 中多读取了 `7616` 行（`8192 - 576`），导致返回给 `FilterBitsNode` 的 `ColumnVector` 总大小为 `8192`。当 `FilterBitsNode` 累积到总行数时，这个多出来的 `7616` 导致了 `bitset.size()` 的膨胀，最终触发断言。

#### 2.4. 对比分析：`TermExpr` 为何未崩溃？

*   **目的：** 既然 `GISFunctionFilterExpr` 和 `TermExpr` 都使用了类似的 `ProcessIndexChunks` 批处理逻辑，为何只有 GIS 查询崩溃？
*   **步骤：**
    *   修改 `data_preparation.py`，添加一个非 PK 的 `int_field` 并为其创建 `INVERTED` 索引，以确保 `TermExpr` 走索引路径。
        ```python
        # In data_preparation.py
        schema.add_field("int_field", DataType.INT64)
        # In generate_batch:
        batch.append({"int_field": idx % 100})
        # Renamed build_id_index to build_int_field_index and target int_field
        ```
    *   创建 `test_term_expr_query.py` 脚本，查询 `int_field`。
    *   在 Milvus 服务器日志中，观察 `TermExpr` 相关的 `AssertInfo`。
*   **结果：**
    *   日志显示 `TermExpr` 确实会遇到 `ProcessIndexOneChunk result size: 576` 的情况，表明其内部也存在非完整批次的场景。
    *   关键是 `TermExpr::ExecVisitorImplForIndex()` 在调用 `ProcessIndexChunks` 后，有一个 `AssertInfo(res->size() == real_batch_size, ...)` 的校验。
    *   **分析结论：** `TermExpr` 之所以没有崩溃，是因为这个 `AssertInfo` 在 **Release 构建** 中不会导致程序终止，而只会记录日志（如果条件不满足）。它强制要求 `ProcessIndexChunks` 返回的 `ColumnVector` 的大小精确地等于 `GetNextBatchSize()` 计算出的 `real_batch_size`。这表明 `ProcessIndexChunks` 在实践中是能够达到这种精确度的，或者 `TermExpr` 的场景没有触发 `GISFunctionFilterExpr` 那样的多余填充问题。而 `GISFunctionFilterExpr` 缺乏这个最终的精确化步骤，导致其返回的 `ColumnVector` 大小不准确，最终触发了 `FilterBitsNode` 中更严格的 `Assert`。

### 3. 解决方案

#### 3.1. 最终解决Bug的方案

Bug 的根本原因在于 `GISFunctionFilterExpr::EvalForIndexSegment()` 在处理 Segment 尾部数据时，其内部的批处理和循环退出逻辑不够精确，导致返回给 `FilterBitsNode` 的 `ColumnVector` 大小超出了 `FilterBitsNode` 所期望的 `real_batch_size`。

**修改代码行：**

在 `internal/core/src/exec/expression/GISFunctionFilterExpr.cpp` 文件的 `EvalForIndexSegment()` 函数的末尾，`return` 语句之前，添加一个强制截断的逻辑。

```cpp
// internal/core/src/exec/expression/GISFunctionFilterExpr.cpp
// In PhyGISFunctionFilterExpr::EvalForIndexSegment() method

// ... existing code ...
        processed_rows += size;
    } // End of for loop

    // CRITICAL FIX: Ensure the returned ColumnVector exactly matches the real_batch_size
    // This handles the case where the loop might have accumulated slightly more
    // due to chunking/batching logic not perfectly aligning with `real_batch_size`.
    if (batch_result.size() > real_batch_size) {
        LOG_WARN("LiYinwei:EvalForIndexSegment: Truncating batch_result from {} to {} to match real_batch_size.", batch_result.size(), real_batch_size);
        batch_result.resize(real_batch_size);
        batch_valid.resize(real_batch_size);
    }

    LOG_INFO("LiYinwei:EvalForIndexSegment returning batch_result size: {}", batch_result.size());
    return std::make_shared<ColumnVector>(std::move(batch_result),
                                          std::move(batch_valid));
}
```

#### 3.2. 方案原理和依据

1.  **问题核心：** `FilterBitsNode` 期望其接收到的 `ColumnVector`（从 `GISFunctionFilterExpr` 返回）的 `size` 能够精确地匹配 `query_context_->get_active_count()` （即 `need_process_rows_`）的剩余部分。
2.  **`EvalForIndexSegment` 的不足：** 尽管 `GISFunctionFilterExpr::EvalForIndexSegment()` 内部的 `for` 循环尝试通过 `if (processed_rows + size >= batch_size_) { break; }` 来控制批次大小，但在处理到 Segment 尾部且需要返回一个非完整 `batch_size_` 的批次时，它的逻辑不够精确。例如，当只需要 576 行时，由于 `576 + 576 < 8192` (batch\_size)，它会继续循环，并从下一个内部 `chunk` 中额外读取 7616 行，导致返回的 `ColumnVector` 总大小变为 8192。
3.  **修复原理：** 新增的 `if (batch_result.size() > real_batch_size) { batch_result.resize(real_batch_size); }` 语句作为一道“安全网”或“校准器”。它强制确保 `EvalForIndexSegment` 函数在返回其结果 `ColumnVector` 之前，将其内部累积的 `batch_result` 的大小精确地截断为 `GetNextBatchSize()` 所计算出的 `real_batch_size`。
4.  **依据：** 这个修复方案借鉴了 `TermExpr` 链路的成功经验。`TermExpr` 通过其 `AssertInfo(res->size() == real_batch_size, ...)` 间接保证了返回位图的精确性。此方案在 `GISFunctionFilterExpr` 层面直接实现了这种精确化，满足了下游 `FilterBitsNode` 的严格要求。

#### 3.3. 方案的副作用和潜在风险

1.  **副作用：**
    *   **性能影响（微乎其微）：** 增加了一个 `if` 判断和可能的 `resize` 操作。对于大多数正常批次（大小等于 `batch_size_`）而言，`if` 条件不满足，几乎没有开销。只有在需要截断的尾部批次时，才会发生实际的 `resize`。这种操作通常是高效的，对整体查询性能影响可以忽略不计。
    *   **数据丢弃（预期行为）：** 截断操作会“丢弃”掉那些超出了 `real_batch_size` 的计算结果。但这些结果本身就是多余的，因为 `FilterBitsNode` 在当前调用中并不需要它们。因此，这不是功能上的副作用，而是校正行为。

2.  **潜在风险：**
    *   **症状级别修复：** 尽管解决了崩溃问题，但此方案仍被认为是症状级别的修复。它没有修改 `GISFunctionFilterExpr::EvalForIndexSegment()` 内部 `for` 循环中导致过度填充的根本逻辑（即 `if (processed_rows + size >= batch_size_)` 退出条件在处理 Segment 尾部批次时的不精确性）。如果未来 `GetNextBatchSize()` 的计算或 `ProcessIndexOneChunk` 的行为发生变化，可能需要更深入地优化内部循环的退出条件。
    *   **依赖 `real_batch_size` 的准确性：** 该方案的前提是 `GetNextBatchSize()` 能够准确计算出当前 `FilterBitsNode` 真正需要多少行。如果 `GetNextBatchSize()` 本身存在 Bug 导致计算错误，此修复也无法纠正。然而，根据目前调试，`GetNextBatchSize()` 似乎是可靠的。
    *   **进一步的测试：** 在合并此修复后，需要对所有 GIS 查询类型（`st_equals`, `st_touches`, `st_intersects` 等）、不同数据分布、不同数据量、以及不同 Compaction 策略下的 Milvus 进行全面的回归测试和性能基准测试，确保修复的稳定性和无其他新引入的问题。

---