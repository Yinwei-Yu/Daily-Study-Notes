# Milvus GIS功能开发与修改指导报告

## 1. 数据类型与存储层

### 1.1 新增GeoSpatial数据类型
- **C++:**
  - 修改 `internal/core/src/common/Types.h`，添加 `DataType::GEOMETRY`。
  - 修改 `FieldData.cpp`、`Array.h` 等，支持WKB格式的序列化/反序列化。
- **Go:**
  - 修改 `internal/types/schema.go`，添加GeoSpatial类型。
  - 修改 `schema.pb.go`，同步proto定义。
- **Proto:**
  - 修改 `internal/core/src/pb/schema.proto`，添加GeoSpatial类型。

**开发建议：**
- 参考现有的String/Array类型实现，保持接口一致性。
- WKB格式存储，便于与GDAL/OGR库对接。

---

## 2. Geometry核心类与GDAL集成

### 2.1 Geometry类实现
- 新建 `internal/core/src/common/Geometry.h/.cpp`，封装GDAL/OGR对象。
- 支持WKT/WKB互转，空间关系计算（equals, contains, intersects等）。

**开发建议：**
- 使用智能指针管理GDAL对象，防止内存泄漏。
- 提供批量转换和计算接口，便于后续优化。

---

## 3. 表达式与查询解析

### 3.1 SQL/DSL语法扩展
- 修改 `internal/parser/planparserv2/Plan.g4`，添加ST_Contains等GIS函数Token。
- 重新生成解析器代码。

### 3.2 解析器访问者实现
- 修改 `parser_visitor.go`，将SQL中的GIS表达式解析为内部表达式树。
- 验证字段类型，解析WKT参数。

### 3.3 协议层定义
- 修改 `internal/proto/plan.proto`，新增 `GISFunctionFilterExpr` 消息类型。
- 支持表达式树的序列化/反序列化。

**开发建议：**
- 参考JSON/Array相关表达式的实现。
- 保证表达式树结构的可扩展性。

---

## 4. 查询执行与表达式系统

### 4.1 表达式系统扩展
- 修改 `internal/core/src/expr/ITypeExpr.h`，新增 `GISFunctioinFilterExpr`。
- 支持空间关系操作符和Geometry参数。

### 4.2 执行器实现
- 新建 `internal/core/src/exec/expression/GISFunctionFilterExpr.cpp/h`，实现 `PhyGISFunctionFilterExpr`。
- 实现Eval方法，批量计算空间关系。

### 4.3 查询计划生成
- 修改 `internal/core/src/query/PlanProto.cpp`，支持GIS表达式的反序列化和树构建。

**开发建议：**
- 优先实现批量处理逻辑，提升性能。
- 添加详细的错误处理和日志。

---

## 5. 存储与数据导入导出

### 5.1 存储层支持
- 修改 `internal/storage/data_codec.go`，支持WKB序列化/反序列化。
- 修改 `insert_data.go`，支持GeoSpatial类型数据插入。

### 5.2 数据校验与转换
- 修改 `internal/proxy/validate_util.go`，实现WKT校验、WKT<->WKB转换。
- 支持空值、默认值处理。

**开发建议：**
- 充分复用go-geom库的WKT/WKB转换能力。
- 保证数据一致性和健壮性。

---

## 6. 空间索引支持

### 6.1 索引算法实现
- 新建 `internal/core/src/indexbuilder/GeometryIndexBuilder.cpp/h`，实现R-Tree索引。
- 修改 `internal/core/src/index/`，支持索引的创建、删除、管理。

### 6.2 查询优化
- 查询执行器支持基于索引的范围查询。
- 优化索引使用和性能统计。

**开发建议：**
- 参考PostGIS/Oracle Spatial的索引实现思路。
- 预留接口支持多种空间索引。

---

## 7. API与客户端支持

### 7.1 Go SDK
- 修改 `client/column/geometry.go`，支持Geometry列类型。
- 实现WKT/WKB转换API。

### 7.2 REST/gRPC接口
- 修改 `internal/http/handlers/`，添加GIS相关API端点。
- 完善参数校验和文档。

**开发建议：**
- 保持API风格与现有类型一致。
- 提供详细的错误信息和用法示例。

---

## 8. 测试与工具

### 8.1 单元测试
- 为Geometry类、表达式、执行器、索引等编写单元测试。

### 8.2 集成测试
- 设计端到端测试场景，覆盖插入、查询、索引、错误处理等。

### 8.3 性能测试
- 编写大数据量插入和查询的性能测试脚本。

**开发建议：**
- 测试用例覆盖常见空间关系和边界情况。
- 自动化测试与CI集成。

---

## 9. 文档与部署

### 9.1 用户文档
- 编写GIS数据类型、SQL语法、API使用、最佳实践等文档。

### 9.2 开发者文档
- 编写架构设计、API参考、扩展开发、故障排查等文档。

### 9.3 部署与演示
- 更新Docker配置和构建脚本，准备部署文档和演示材料。

---

## 10. 其他注意事项
- 关注GDAL库版本兼容性，提前测试。
- 代码实现优先保证健壮性和可维护性。
- 及时同步文档和代码变更。
- 定期代码审查和性能分析。

---

**本报告为Milvus GIS功能开发的全链路指导，建议结合详细开发计划分阶段推进。** 