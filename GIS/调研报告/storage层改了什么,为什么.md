### 1. data_codec.go

**功能**：负责插入数据的序列化/反序列化（InsertData <-> binlog blob），即将内存中的数据结构转为可持久化的二进制格式，或反向操作。 **geo相关改动**：

- 支持新的 geo 数据类型（如 DataType_Geometry），在序列化/反序列化流程中增加对 geo 字段的处理分支。
- 在序列化时，针对 geo 字段调用专门的 PayloadWriter/Reader 方法（如 AddOneGeometryToPayload）。 **原因**：geo 字段的二进制格式和普通标量/向量字段不同，需要专门的序列化/反序列化逻辑。

---

### 2. data_sorter.go

**功能**：对插入数据按主键或时间戳等进行排序，保证数据有序。 **geo相关改动**：

- 在排序时，增加对 geo 字段的兼容处理，确保 geo 字段不会因类型不支持而导致排序失败。 **原因**：geo 字段可能作为主键或参与排序，需保证排序逻辑对所有字段类型兼容。

---

### 3. insert_data.go

**功能**：定义 InsertData 结构体，管理一批插入数据的内存表示。 **geo相关改动**：

- InsertData.Data 字段支持 geo 类型的 FieldData。
- NewInsertData/NewFieldData 等工厂方法支持 geo 字段的初始化。 **原因**：InsertData 需能容纳 geo 字段，便于后续序列化、排序等操作。

---

### 4. payload.go

**功能**：定义 PayloadWriter/Reader 接口及实现，负责字段级别的二进制读写。 **geo相关改动**：

- PayloadWriter/Reader 增加 AddOneGeometryToPayload、GetGeometryFromPayload 等接口。
- 支持 geo 字段的写入、读取、长度计算等。 **原因**：geo 字段的二进制格式与普通字段不同，需专门接口处理。

---

### 5. payload_reader/writer

**功能**：实现 payload.go中的接口，具体负责将 geo 字段写入 Parquet/Arrow 格式，或从中读取。 **geo相关改动**：

- 在 switch-case 结构中增加 geo 类型的分支。
- 处理 geo 字段的特殊二进制格式（如 WKB/WKT）。 **原因**：geo 字段的存储和读取方式与标量/向量不同，需定制实现。

---

### 6. print_binlog_test.go

**功能**：用于调试/测试，打印 binlog 文件内容，便于人工检查。 **geo相关改动**：

- 打印 geo 字段时，支持解析和展示 WKB/WKT 格式。
- 测试用例增加 geo 字段的覆盖。 **原因**：便于开发者调试 geo 字段的持久化和读取正确性。

---

### 7. serde.go

**功能**：负责 Arrow/Parquet 与内存结构的序列化/反序列化（SerDe）。 **geo相关改动**：

- 在 serdeEntry、compositeRecord 等结构和方法中增加 geo 字段的处理。
- 支持 geo 字段的 Arrow 类型映射和序列化。 **原因**：geo 字段需要和 Arrow/Parquet 互转，便于高效存储和读取。

---

### 8.  utils/_test.go

**功能**：存放通用工具函数和测试。 **geo相关改动**：

- 工具函数支持 geo 字段的校验、转换等。
- 测试用例覆盖 geo 字段的各种场景。 **原因**：保证 geo 字段在各种通用操作中不会出错。

## 总结

**改动核心**：所有涉及数据结构、序列化、排序、读写、测试的地方，都增加了对 geo 字段的识别和处理分支，保证 geo 字段能像普通字段一样被正确存储、读取、打印和测试。

**改动原因**：geo 字段的二进制格式、内存结构和常规字段不同，只有在所有数据流转链路中都支持，才能保证 geo 数据的全流程可用和高效。

