该 PR（https://github.com/milvus-io/milvus/pull/43661）主要完成了 Milvus 2.5 对地理空间（Geospatial）数据类型和 GIS（地理信息系统）函数的支持，具体包括以下几个方面：

1. **创建与描述地理空间类型的集合**：可以在 Milvus 集合中定义和描述 geospatial 类型字段。
2. **插入地理空间数据**：支持将地理空间数据插入到 insert binlog。
3. **加载包含地理空间数据的分片到内存**：分片（segment）中的地理空间数据可以正确加载到内存进行后续处理。
4. **查询与搜索支持地理空间数据展示**：查询和搜索结果可以显示 geospatial 数据。
5. **GIS 函数支持（如 ST_EQUALS）**：在查询中支持使用部分 GIS 空间关系函数进行过滤。

**主要实现细节**：
- 在 C++ 和 Go 层新增 Geospatial 类型，定义了相应的数据结构和接口。
- 引入地理空间数据处理库（C++ 使用 GDAL，Go 使用 go-geom）。
- 修改协议接口，支持 geospatial 数据的序列化和反序列化。
- 数据流转中，客户端与 proxy 以 WKT 格式交互，proxy 转为 WKB 格式进行后续处理，并在数据写入、分片封装、加载、缓存管理等流程中支持 geospatial 类型。
- 查询操作目前支持单列的地理空间数据的简单过滤（空间关系），并实现了查询表达式的解析和执行（目前仅支持暴力搜索）。
- 客户端支持地理空间数据的输入与测试（pymilvus 有相应修改）。

该 PR 还引入了大量代码变更（5700+ 行新增，101 个文件变更），主要是为 geospatial 类型及相关 GIS 功能的支持做底层和接口层的适配。

### 1. 创建与描述地理空间类型的集合

Milvus 的集合（Collection）Schema 支持自定义字段类型，通过 `Schema` 和 `FieldMeta`，可以定义 geospatial 类型字段：

```c++ name=internal/core/src/common/Schema.h
class Schema {
    // 添加字段（支持地理空间类型）
    void AddField(const FieldName& name, const FieldId id, DataType data_type, bool nullable, std::optional<DefaultValueType> default_value);
    void AddField(const FieldName& name, const FieldId id, DataType data_type, DataType element_type, bool nullable);
    // 支持多种类型，包括 VECTOR_ARRAY（地理空间向量数组）、标量等
};
```

`FieldMeta` 中有地理空间相关的扩展结构体和参数，例如：

```c++ name=internal/core/src/common/FieldMeta.h
class FieldMeta {
    // 地理空间、向量等信息
    struct VectorInfo {
        int64_t dim_;
        std::optional<knowhere::MetricType> metric_type_;
    }
}
```

Schema 允许通过类型和参数构造地理空间字段（如 WKT/WKB 格式的空间向量），并且支持在 Collection 层定义这些字段。

---

### 2. 插入地理空间数据，写入 binlog

Milvus 的 binlog 机制支持各种字段类型的数据，包括地理空间数据。插入数据时，通过如下流程：

```c++ name=internal/core/unittest/test_utils/storage_test_utils.h
inline LoadFieldDataInfo PrepareInsertBinlog(int64_t collection_id, int64_t partition_id, int64_t segment_id, const GeneratedData& dataset, const ChunkManagerPtr cm, std::string mmap_dir_path = "", std::vector<int64_t> excluded_field_ids = {}) {
    auto insert_data = std::make_shared<InsertData>(payload_reader);
    FieldDataMeta field_data_meta{collection_id, partition_id, segment_id, field_id};
    insert_data->SetFieldDataMeta(field_data_meta);
    auto serialized_insert_data = insert_data->serialize_to_remote_file();
    cm->Write(file, serialized_insert_data.data(), serialized_insert_size);
}
```

每个字段（包括地理空间字段）对应一个 binlog 文件，按列存储。插入接口支持 geospatial 类型数据的序列化与写入，binlog 文件格式可扩展以支持空间数据。

---

### 3. 加载包含地理空间数据的分片到内存

分片（segment）加载时，需要将 geospatial 字段的数据正确读入内存，支持 mmap 优化：

```c++ name=internal/core/src/segcore/ChunkedSegmentSealedImpl.cpp
void ChunkedSegmentSealedImpl::load_field_data_internal(const LoadFieldDataInfo& load_info) {
    for (auto& [id, info] : load_info.field_infos) {
        auto field_id = FieldId(id);
        auto field_meta = schema_->operator[](field_id);
        auto column = MakeChunkedColumnBase(data_type, std::move(translator), field_meta);
        load_field_data_common(field_id, column, num_rows, data_type, info.enable_mmap, false);
    }
}
```

加载 geospatial 字段时，支持不同的数据类型、分片格式和 mmap，最终将数据结构映射到内存并可用于后续查询和搜索。

---

### 4. 查询与搜索支持地理空间数据展示

查询和搜索会根据 Schema 返回 geospatial 字段数据，相关代码示例：

```c++ name=internal/core/src/query/SearchOnGrowing.cpp
void SearchOnGrowing(const segcore::SegmentGrowingImpl& segment, const SearchInfo& info, const void* query_data, int64_t num_queries, Timestamp timestamp, const BitsetView& bitset, SearchResult& search_result) {
    auto& schema = segment.get_schema();
    auto& field = schema[vecfield_id];
    auto data_type = field.get_data_type();
    AssertInfo(IsVectorDataType(data_type), "[SearchOnGrowing]Data type isn't vector type");
    ...
    // 查询结果支持 geospatial 字段展示
}
```

结果数据的构建会将 geospatial 字段按格式填充到返回结构，客户端可以直接展示地理空间数据。

---

### 5. GIS 函数支持（如 ST_EQUALS）

在查询过滤和表达式处理阶段，部分 GIS 空间关系函数（如 ST_EQUALS）被支持：

```c++ name=internal/core/src/query/Utils.h
template <typename T, typename U>
inline bool Match(const T& x, const U& y, OpType op) {
    // 支持空间关系运算，如 ST_EQUALS, ST_INTERSECTS
}

template <>
inline bool Match<std::string>(const std::string& str, const std::string& val, OpType op) {
    // 字符串空间关系比较
}
```

在查询计划和执行时，空间运算（如 ST_EQUALS）可用于对 geospatial 字段进行过滤和条件判断，暴力或索引方式均可。

---

**总结**  
Milvus 2.5 的地理空间数据支持，底层通过 Schema/FieldMeta 定义类型，binlog 支持序列化与写入，segment 加载适配分片存储，查询/搜索流程完整支持 geospatial 字段的数据展示，并在表达式中实现了空间关系函数的基本支持（如 ST_EQUALS），代码实现涉及 C++/Go 多层接口和数据流转。

如需具体某一实现细节或代码段的进一步讲解，可以继续追问。