
### . [task_index.go](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-sandbox/workbench/workbench.html)

**功能**：负责索引相关任务（如创建、删除、描述索引等）的调度与参数校验。 **geo相关改动**：

- 支持 geometry 字段的索引创建、参数校验。例如在解析索引参数、判断字段类型时，增加对 [schemapb.DataType_Geometry](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-sandbox/workbench/workbench.html) 的分支。
- 在自动索引、参数包装等流程中，兼容 geometry 字段的特殊参数（如空间索引类型、空间参数）。 **原因**：geo 字段的索引类型、参数与普通向量/标量不同，需单独处理。

---

### 2. [task_insert_test.go](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-sandbox/workbench/workbench.html)

**功能**：插入任务相关的单元测试。 **geo相关改动**：

- 测试用例中增加 geometry 字段的插入、对齐校验、数据生成等。
- 验证 geometry 字段能被正确插入、校验、对齐（如 `CheckAligned` 测试 geometry 字段）。 **原因**：保证插入流程对 geo 字段的支持和健壮性。

---

### 3. [task_query.go](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-sandbox/workbench/workbench.html)

**功能**：负责查询任务（如 Query、Retrieve）的调度、参数处理、结果组装。 **geo相关改动**：

- 查询参数、输出字段、schema 解析等流程中，兼容 geometry 字段。
- 查询结果组装时，能正确处理 geometry 字段的返回（如 WKB 格式）。 **原因**：geo 字段的查询、结果返回格式与普通字段不同，需全链路兼容。

---

### 4. [task_search.go](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-sandbox/workbench/workbench.html)

**功能**：负责向量/混合检索任务的调度、参数处理、结果组装。 **geo相关改动**：

- 支持 geometry 字段的检索参数、schema 解析、结果组装。
- 检索结果中能正确返回 geometry 字段数据。 **原因**：geo 字段可能作为检索目标或输出字段，需保证检索链路兼容。

---

### 5. [task_test.go](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-sandbox/workbench/workbench.html)

**功能**：综合性任务相关的单元测试。 **geo相关改动**：

- 增加 geometry 字段相关的任务测试，如索引、检索、插入、schema 兼容性等。
- 验证 geometry 字段在各种任务中的全流程可用性。 **原因**：保证所有任务类型对 geo 字段的端到端支持。

---

### 6. [validate_util.go](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-sandbox/workbench/workbench.html)

**功能**：字段、数据、结果等的通用校验工具。 **geo相关改动**：

- 增加 [validateGeometryFieldSearchResult](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-sandbox/workbench/workbench.html) 等函数，专门校验 geometry 字段的合法性（如 WKB 格式、空间数据有效性）。
- 在通用校验流程中，增加对 geometry 字段的分支处理。 **原因**：geo 字段的校验逻辑与普通字段不同，需专门处理，防止无效空间数据流入系统。

---

## 总结

**核心改动**：所有涉及 schema 解析、参数校验、数据插入、检索、索引、结果组装、单元测试等流程的地方，都增加了对 geometry/geo 字段的识别和专门处理分支，保证 geo 字段能像普通字段一样被正确支持。
