# Milvus Geometry 模块设计文档

## 一、设计目标

构建一个**统一的几何数据类型**，用于在 Milvus 向量数据库中支持对空间数据（GIS 数据）的插入、存储、索引和查询操作。该类型需满足以下需求：

- 支持标准几何格式（WKT / WKB）
    
- 支持多种几何类型（点、线、多边形等）
    
- 支持空间参考系统（SRID）与坐标转换
    
- 封装空间计算与空间关系判断
    
- 支持几何验证与修复，保证数据正确性
    
- 提供序列化与可视化接口，便于调试与网络传输
    

---

## 二、系统架构

```
+--------------------------+
|      Geometry 类         |
+--------------------------+
| - GEOSContextHandle_t    |
| - GEOSGeometry*          |
| - std::unique_ptr<WKB>   |
| - size_t wkb_size_       |
| - int srid_              |
+--------------------------+
| + Construct(WKT/WKB)     |
| + to_wkt_string()        |
| + to_wkb_string()        |
| + GetArea(), Length()    |
| + SpatialOps (==, ∩, ∪)  |
| + SRID Transform()       |
| + isValid(), repair()    |
+--------------------------+
```

---

## 三、核心设计组件

### 1. 枚举类型

|枚举名|说明|
|---|---|
|`GeometryType`|几何对象类型（点、线等）|
|`WkbByteOrder`|WKB 字节序（大端/小端）|
|`WkbGeometryType`|WKB 中表示的几何类型|

---

### 2. 内部数据结构

```cpp
struct WkbHeader {
    uint8_t byte_order;
    uint32_t geometry_type;
    uint32_t srid;
};
```

用于解析 WKB 头部以提取类型与 SRID。

---

### 3. 核心成员变量

|成员名|类型|功能描述|
|---|---|---|
|`geos_context_`|`GEOSContextHandle_t`|管理 GEOS 上下文|
|`geos_geometry_`|`GEOSGeometry*`|GEOS 几何对象指针|
|`wkb_data_`|`std::unique_ptr<unsigned char[]>`|几何的 WKB 表示|
|`wkb_size_`|`size_t`|WKB 数据大小|
|`srid_`|`int`|空间参考系标识，默认 WGS84 (4326)|

---

## 四、功能设计

### 1. 构造函数

- `Geometry(const std::string& wkt, int srid = 4326)`
    
- `Geometry(const void* wkb, size_t size, int srid = 4326)`
    
- 支持 copy/move 构造与赋值
    

自动完成：

- GEOS 上下文初始化
    
- 解析或构造 GEOSGeometry
    
- 转为 WKB 存储
    
- 几何合法性验证与必要修复
    

---

### 2. 序列化接口

|方法|返回类型|说明|
|---|---|---|
|`to_wkt_string()`|`std::string`|获取 WKT 表示|
|`to_wkb_string()`|`std::string`|获取 WKB 二进制串|
|`operator std::string`|`std::string`|用作默认文本表示|

---

### 3. 几何属性查询

|方法|类型|描述|
|---|---|---|
|`GetArea()`|`double`|面积（适用于 Polygon）|
|`GetLength()`|`double`|长度（适用于 LineString）|
|`GetCentroid()`|`Geometry`|几何中心点|
|`GetEnvelope()`|`Geometry`|包络盒（Bounding box）|

---

### 4. 空间关系判断（布尔操作）

|方法|描述|
|---|---|
|`equals()`|是否完全相同|
|`touches()`|是否接触|
|`overlaps()`|是否重叠|
|`crosses()`|是否交叉|
|`contains()`|是否包含|
|`within()`|是否被包含|
|`intersects()`|是否相交|
|`disjoint()`|是否不相交|

---

### 5. 空间几何操作（几何输出）

|方法|返回类型|描述|
|---|---|---|
|`intersection()`|`Geometry`|交集|
|`union_()`|`Geometry`|并集（注意方法名带 `_`）|
|`difference()`|`Geometry`|差集|

---

### 6. 坐标系统支持

| 方法                           | 描述                   |
| ---------------------------- | -------------------- |
| `SetSRID(int)`               | 设置空间参考系统             |
| `GetSRID()`                  | 获取当前 SRID            |
| `transform(int target_srid)` | 使用 PROJ 进行坐标变换（EPSG） |
|                              |                      |

---

### 7. 合法性校验与修复

|方法|返回类型|描述|
|---|---|---|
|`isValid()`|`bool`|是否合法 GEOS 几何|
|`repair()`|`Geometry`|返回修复后的合法几何|

---

## 五、错误处理机制

使用 `AssertInfo()` 和 `PanicInfo()` 宏函数进行健壮性检查与调试。

- 所有构造/操作均在非法输入时抛出异常或返回空对象
    
- 所有资源释放集中在 `cleanup()` 中完成
