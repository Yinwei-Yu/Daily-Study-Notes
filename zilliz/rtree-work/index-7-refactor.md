### 提交信息
- 提交: `cba865f25a168faa2bdae65fcd3b6d35bc499069`
- 作者: Yinwei Li
- 标题: 更新 RTreeIndex 配置处理与测试用例

### 变更概览
- 目标：统一索引配置读取为“字符串输入”，内部再转换为目标类型；同步调整单测与统计逻辑，清理多余依赖。
- 影响范围：`RTreeIndex`、`Utils` 配置解析、UT 期望值与构建参数。

### 关键改动
- 配置读取
  - `internal/core/src/index/Utils.cpp`
    - `GetFillFactorFromConfig`：从 `string` 读取并用 `std::stod` 转换为 `double`。
    - `GetIndexCapacityFromConfig`、`GetLeafCapacityFromConfig`：从 `string` 读取并用 `std::stoi` 转换为 `uint32_t`。
- RTreeIndex
  - `internal/core/src/index/RTreeIndex.cpp`
    - `Build`：改用上述强一致的配置读取函数；时序仍为先 `InitForBuildIndex()` 再设参数。
  - `internal/core/src/index/RTreeIndex.h`
    - `Count()`：当未构建完成时，使用 `wrapper_->count() + null_offset_.size()` 统计总行数（继续对齐空值管理方案）。
    - 头文件清理：移除无用包含（`Index.h`、`StringIndex.h`）。
- 单元测试
  - `internal/core/unittest/test_rtree_index.cpp`
    - 构建配置全部改为字符串形式：`fillFactor/indexCapacity/leafCapacity` 使用 `"0.8"/"50"/"50"` 等。
    - `.meta.json` 断言改为字符串值：`"fill_factor"`, `"index_capacity"`, `"leaf_capacity"` 比对为字符串（与当前产出保持一致）。

### 行为变化
- 配置输入更宽容：JSON 中数值配置以字符串形式传入更稳定，内部统一转换，减少类型不一致导致的解析错误。
- 行数统计口径保持一致：未完成构建时也正确计入 `null_offset_`，与加载路径对齐。

### 风险与兼容性
- 若外部仍以数值类型传入配置键，需确保上层调用按字符串提供或在进入本层前转换；否则会触发断言。
- 单测期望 `.meta.json` 中相关字段为字符串；请确保产出一致（部署环境需与当前实现一致）。

### 建议
- 在用户入口层统一将索引构建配置规范化（全部字符串），避免类型歧义。
- 在日志中打印配置解析前后的值与类型，便于排障。

- 本次提交将索引配置读取统一为字符串并内部转换，测试用例与元数据校验随之调整，同时保持 R-Tree 行数统计对空值的正确计入与头文件清理，提升了健壮性与一致性。

