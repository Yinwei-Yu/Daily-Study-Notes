### 1. chunkSegmentSealedImpl.cpp

**功能**：  
实现了 [ChunkedSegmentSealedImpl](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-sandbox/workbench/workbench.html)，负责 Sealed Segment 的数据加载、索引、删除、查询等核心逻辑。

**geo 相关改动**：

- 在 `load_field_data_common`、[generate_interim_index](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-sandbox/workbench/workbench.html) 等方法中，增加了对 geo 类型（如 DataType::GEOMETRY）的分支处理，确保地理空间数据能被正确加载、索引和管理。
- 支持 geo 字段的 mmap 加载、chunk 组织、索引等。

**原因**：  
Sealed Segment 需要支持 geo 字段的高效加载和索引，保证 geo 数据能与其他类型一样被管理和检索。

---

### 2. ConcunrentVector.cpp

**功能**：  
实现并发向量（ConcurrentVector），用于高效存储和访问分块数据，支持多种数据类型。

**geo 相关改动**：

- 在 [set_data_raw](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-sandbox/workbench/workbench.html) 等方法中，增加了对 geo 类型的处理分支（如 DataType::GEOMETRY、DataType::ARRAY），支持 geo 数据的批量写入和 chunk 组织。

**原因**：  
geo 字段的数据结构与普通标量/向量不同，需要专门的批量写入和分块存储逻辑。

---

### 3. [InsertRecord.h](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-sandbox/workbench/workbench.html)

**功能**：  
定义了 [InsertRecord](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-sandbox/workbench/workbench.html) 结构体，管理插入数据的主键、时间戳、字段数据等，支持 Growing/Sealed 两种模式。

**geo 相关改动**：

- 在 `append_field_meta`、[append_data](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-sandbox/workbench/workbench.html) 等模板方法中，增加了对 geo 类型的分支，确保 geo 字段能被正确注册和存储。
- 支持 geo 字段的主键映射、数据插入、分块管理等。

**原因**：  
插入记录需要支持 geo 字段的注册、数据插入和主键映射，保证 geo 数据能参与增量写入和检索。

---

### 4. [SegmentGrowingImpl.cpp](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-sandbox/workbench/workbench.html) / [SegmentGrowingImpl.h](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-sandbox/workbench/workbench.html)

**功能**：  
实现了 Growing Segment 的数据插入、加载、查询等逻辑。

**geo 相关改动**：

- 在 [Insert](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-sandbox/workbench/workbench.html)、`load_field_data_common`、`append_field_meta` 等流程中，增加了 geo 类型的处理分支。
- 支持 geo 字段的批量插入、数据填充、分块管理等。

**原因**：  
Growing Segment 需要支持 geo 字段的动态写入和分块管理，保证 geo 数据能被实时插入和查询。

---

### 5. SegmentSealedImpl.cpp

**功能**：  
实现了 Sealed Segment 的部分接口，主要用于数据的只读管理和查询。

**geo 相关改动**：

- 在数据加载、主键检索等流程中，增加了 geo 类型的分支处理。

**原因**：  
Sealed Segment 需要支持 geo 字段的只读检索和主键映射。

---

### 6. Utils.cpp / [Utils.h](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-sandbox/workbench/workbench.html)

**功能**：  
提供 segcore 层的通用工具函数。

**geo 相关改动**：

- 可能增加了 geo 类型相关的辅助函数，便于统一处理 geo 字段。

**原因**：  
便于 geo 字段在 segcore 层的统一处理和复用。