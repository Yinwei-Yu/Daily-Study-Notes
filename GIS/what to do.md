### 1. **表达式与查询解析相关模块**

- **parser/planparserv2/parser_visitor.go**  
  新增了对ST_Equals、ST_Touches、ST_Overlaps、ST_Crosses、ST_Contains、ST_Intersects、ST_Within等空间关系函数的解析支持，将SQL中的GIS表达式解析为内部表达式树。

- **internal/proto/plan.proto**  
  新增了`GISFunctionFilterExpr`消息类型，定义了空间关系操作符（Equals、Touches、Overlaps、Crosses、Contains、Intersects、Within）和WKT参数，作为表达式树的节点。

- **internal/core/src/exec/expression/GISFunctionFilterExpr.h/cpp**  
  新增了C++端的物理表达式节点`PhyGISFunctionFilterExpr`，实现了空间关系算子的批量计算逻辑，直接调用GDAL/OGR库。

- **internal/core/src/expr/ITypeExpr.h**  
  新增了`GISFunctioinFilterExpr`表达式类型，作为表达式树中的节点。

- **internal/core/src/query/PlanProto.cpp**  
  新增了对`GISFunctionFilterExpr`的反序列化和表达式树构建逻辑。

---

### 2. **几何数据类型与存储相关模块**

- **internal/core/src/common/Geometry.h**  
  新增了`Geometry`类，封装GDAL/OGR的几何对象，支持WKT/WKB互转和空间关系计算。

- **internal/core/src/common/Array.h、FieldDataInterface.h、ChunkWriter.cpp、Utils.cpp、SegmentSealedImpl.cpp**  
  增加了对`DataType::GEOMETRY`的支持，包括数据的存储、序列化、反序列化、内存管理等。

- **internal/proxy/validate_util.go**  
  新增了对几何数据的校验、WKT与WKB互转、插入和查询时的格式转换。

- **client/column/geometry.go**  
  Go SDK端增加了对几何数据的封装和处理。

---

### 3. **SQL/DSL接口与API层**

- **SQL解析层**  
  支持用户通过SQL/DSL直接调用ST_Equals等空间关系函数。

- **API接口**  
  支持插入和查询Geometry类型字段，自动完成WKT/WKB转换。

---

### 4. **依赖与构建系统**
- **CMake/Conan等依赖管理**  
  增加了GDAL/OGR等地理信息库的依赖。

---

### 5. **测试与工具**
- **单元测试、数据生成工具**  
  增加了几何数据的生成、插入、查询相关的测试用例和工具。

---

## **总结**

作者为了添加GIS系统，主要涉及**表达式解析、查询执行、数据类型与存储、API接口、依赖管理、测试工具**等多个核心模块的修改。  
这些改动实现了从SQL到C++内核的全链路GIS表达式支持，并在数据流、存储、校验等各环节打通了Geometry类型的处理能力。