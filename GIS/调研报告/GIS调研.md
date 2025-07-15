## 背景

Milvus 是一款高性能的向量数据库，广泛应用于图像检索、推荐系统、语义搜索等领域。随着 LBS（基于位置服务）、多模态检索、无人驾驶等场景日益融合**空间感知**与**语义检索**需求，用户越来越需要能够在“地理上下文”中进行向量搜索。

目前 Milvus 仅支持数值和向量类型字段，**缺乏对地理信息的原生理解与索引支持**。这限制了它在新兴空间计算应用中的能力。

## 需求

以下是一个典型的地理信息与向量检索结合的使用场景:

*在空间约束下进行向量搜索*

> 我想搜索离我1公里内，且语义上最相似的商品、商店、图片。”

示例行业:

- 本地服务（如饿了么、Uber）
- 地图内容推荐（如高德地图推荐 POI）
- 电商场景（LBS广告、附近同款检索）
- 安防监控（给定位置附近相似人脸追踪

除了与向量检索结合，支持地理信息系统也可以满足诸多常规的基于地理信息分析的需求。例如热力图分析，规划运输道路，通过统计分析规划市场选址等。

有时用户还可能根据不同的使用场景，选用不同的坐标系来对地理数据进行不同角度的分析。且由于地理数据通常来源于第三方平台（地图，遥感，测试等），数据格式多样，最少需要支持WKT，WKB，后续可按需添加GeoJson，GeoHash等数据格式。但坐标系和数据格式转换功能并不是一个数据库应该提供的能力，数据库应该专注于数据的增删改查等，所以没有为milvus内核提供这两个功能。

基于上述需求，需要为milvus提供地理信息数据结构的支持。包括其核心数据类型的定义，多种查询算子，索引优化等。下面从几个方面来详细说明可行的实施方案。

## 目标

### 数据类型

地理空间数据类型（Geo Spatial DataType）是用于描述地理空间信息的数据结构。在OGC制定的SFA标准中，有以下常见的几何类型：

1. 点（Point）：表示一个二维坐标，通常根据比例尺的不同而代表不同的对象
2. 线串（LineString）：由两个及以上的点组成的有序集合，常用来代表河流，道路等
3. 多边形（Polygon）：表示一个平面区域，可以有“洞”
4. 多点（MutilPoint）：多个点的集合
5. 多线串（MutilLineString）：多个线串的集合
6. 多多边形（MutilPolygon）：多个多边形的集合
7. 几何集合（GeometryCollection）：以上所有几何体构成的集合

这些数据类型的输入输出在SFA标准中有两种表示方式：WKT（Well Known Text）和WKB（Well Known Binary）。前者是人类可阅读格式如：`POINT (0 0)`，`MULTILINESTRING ((0 0， 1 1)， (2 2， 3 3))`，后者是二进制格式，用于高效存储。

在milvus中，我们不需要手动定义上述数据类型，而是可以引入第三方库，在milvus中包装第三方库提供的数据结构，即可对几何数据进行统一的处理并对外提供统一的处理接口。

常见的第三方库有GEOS，GDAL等。其中GEOS（Geometry Engine - Open Source） 是一个遵循 OGC Simple Feature Access（SFA）标准的开源几何计算库，其内部已经定义了一套完整的、标准化的几何数据类型和操作接口 ，是 PostGIS， QGIS等地理信息数据库的核心依赖，且其十分轻量，不需要额外的依赖，非常适合嵌入到milvus内核当中。而GDAL虽然功能更为强大，但是其中包含了许多与光栅处理等和地理信息无关的内容，且依赖项繁多，构建复杂，比较臃肿，不适合嵌入到milvus内核中，而更适合作为数据预处理工具使用。所以这里我们选择抛弃原pr中集成GDAL库的做法，更改为使用GEOS作为地理信息系统的核心依赖。

### 坐标系统

常用的坐标系统有WGS 84，即国际经纬度标准；Web Mercator，网络地图标准等。

GEOS底层是无坐标系的，仅对几何图形做在笛卡尔集上的运算。而用户输入的WKT信息可能是基于经纬度信息，在尺度足够大的情况下，会出现球面坐标和平面坐标直接转换导致的误差过大问题，因此，我们需要引入PROJ库来在数据进入GEOS层处理之前，先进行坐标系的转换。

### 查询算子

常见的查询算子包括：

#### 1. 拓扑关系（Topological Relationships）

这些函数用于判断两个几何对象之间是否存在某种拓扑关系（如相交、包含、覆盖等），基于的是空间对象的形状和位置之间的逻辑关系。

| 函数名                   | 描述                                                                   |
| --------------------- | -------------------------------------------------------------------- |
| `ST_Contains(A， B)`   | 判断 **A 是否完全包含 B**，并且 A 和 B 的内部至少有一个共同点。例如：一个城市是否包含一个公园。              |
| `ST_Covers(A， B)`     | 判断 **A 是否覆盖 B**，即 B 的所有点都在 A 内部或边界上。                                 |
| `ST_Crosses(A， B)`    | 判断 **A 和 B 是否部分相交**，但不是完全包含。例如：两条交叉的道路。                              |
| `ST_Disjoint(A， B)`   | 判断 **A 和 B 是否没有交集**，即两者没有任何公共点。                                      |
| `ST_Equals(A， B)`     | 判断 **A 和 B 是否表示相同的几何对象**，即具有相同的坐标序列和类型。                              |
| `ST_Intersects(A， B)` | 判断 **A 和 B 是否有至少一个公共点**，是最常用的空间关系判断函数。                               |
| `ST_Overlaps(A， B)`   | 判断 **A 和 B 是否重叠**，即它们具有相同的维度（如都是面或都是线），并且部分区域重合，但彼此都不完全包含对方。         |
| `ST_Touches(A， B)`    | 判断 **A 和 B 是否仅在边界上接触**，但内部不相交。例如：相邻地块的边界。                            |
| `ST_Within(A， B)`     | 判断 **A 是否完全位于 B 内部**，并且 A 和 B 的内部至少有一个公共点。是 `ST_Contains(B， A)` 的别名。 |

#### 2. 距离关系（Distance Relationships）

这些函数用于判断两个几何对象之间的距离是否小于某个阈值，通常用于近邻搜索、缓冲区查询等场景。

| 函数名                                                       | 描述                                                                 |
| --------------------------------------------------------- | ------------------------------------------------------------------ |
| `ST_DWithin(geomA， geomB， distance)`                      | 判断 **geomA 和 geomB 是否在指定距离范围内**，常用于查找“附近的点”、“附近的城市”等。              |
| `ST_PointInsideCircle(point， center_x， center_y， radius)` | 判断 **一个点是否位于以 (center_x， center_y) 为圆心、radius 为半径的圆内**。适用于圆形缓冲区查询。 |
#### 3. 数据测量（Measurement Functions）

| 函数名                                         | 描述                                |
| ------------------------------------------- | --------------------------------- |
| `ST_Area(geometry)`                         | 返回多边形几何的面积。                       |
| `ST_Distance(geometry1， geometry2)`         | 返回两个几何之间的最小欧几里得距离（二维平面）。          |
| `ST_DistanceSphere(geography1， geography2)` | 使用球面模型计算两个地理坐标（经纬度）之间的最短距离（单位：米）。 |

#### 4. 聚合

| 函数名                                                | 描述                   |
| -------------------------------------------------- | -------------------- |
| [ST_Union](https://postgis.net/docs/ST_Union.html) | 将多个几何对象合并为一个无重叠的几何对象 |

#### 5. 几何属性访问器（Geometry Accessors）

|函数名|参数|描述|
|---|---|---|
|`ST_IsValid`|geometries|判断几何是否有效|
|`ST_IsSimple`|geometries|判断几何是否简单（无自相交）|
|`ST_GeometryType`|geometries|返回几何类型（如 POINT， LINESTRING， POLYGON 等）|
|`ST_NPoints`|geo_arr|返回几何中的顶点数量|

#### 6.  特定几何属性访问器

| 几何类型               | 函数名                             | 参数         | 描述                       |
| ------------------ | ------------------------------- | ---------- | ------------------------ |
| **点（Point）**       | `ST_X(geometry)`                | point      | 返回点的 X 坐标                |
|                    | `ST_Y(geometry)`                | point      | 返回点的 Y 坐标                |
| **线串（LineString）** | `ST_Length(geometry)`           | linestring | 返回线串的长度                  |
|                    | `ST_StartPoint(geometry)`       | linestring | 返回线串的第一个点（作为 POINT 类型）   |
|                    | `ST_EndPoint(geometry)`         | linestring | 返回线串的最后一个点（作为 POINT 类型）  |
|                    | `ST_NPoints(geometry)`          | linestring | 返回线串中包含的坐标点数量            |
| **多边形（Polygon）**   | `ST_Area(geometry)`             | polygon    | 返回多边形的面积                 |
|                    | `ST_NRings(geometry)`           | polygon    | 返回多边形中的环数量（外环 + 内环）      |
|                    | `ST_ExteriorRing(geometry)`     | polygon    | 返回多边形的外环（作为 LINESTRING）  |
|                    | `ST_InteriorRingN(geometry， n)` | polygon， n | 返回第 n 个内环（作为 LINESTRING） |
|                    | `ST_Perimeter(geometry)`        | polygon    | 返回所有环的总周长（外环 + 所有内环）     |

### 索引优化

为加快查询速度，需要为地理信息数据建立索引。在标准的数据库中，无论是标量还是向量，都是根据被索引列的值来构建索引的，但是对于空间索引却不同，因为数据库并不能直接索引几何字段的值，也就是几何对象本身。因此，需要索引几何对象的**边界范围框**。下面这个例子来自于PostGIS文档：

![_images/bbox.png](https://postgis.net/workshops/zh_Hans/postgis-intro/_images/bbox.png)

上图中，和黄色星星相交的线的数量是 **1**，即红色那条线。但是与黄色框相交的范围框有红色和蓝色，共 2 个。

数据库求解 “什么线与黄色星相交” 这个问题，是先用空间索引求解 “什么范围框与黄色范围框相交” 这个问题的（速度非常快），然后才是 “什么线与黄色的星星相交”。上述过程仅对于第一次测试的空间要素而言。

对于数量庞大的数据表，这种索引先行，然后局部精确计算的 “两遍法” 可以在根本上减少查询计算量。

常见的索引方式有：R-Tree，Quadtree等，目前阶段我们先支持R-Tree索引。

对于R-Tree索引，我们可以选择使用libspatialindex库作为第三方依赖，它不仅支持R-Tree索引，还支持其他类型的空间索引，方便后期进行扩展。此外，它还支持索引的持久化存储。然而，由于R-Tree索引建立需要使用到边界框，而GEOS库不支持边界框的数据格式，所以我们在建立索引时还需要一个中间层来给几何对象包装上一个边界框，幸运的是，GEOS库对每个几何对象都提供了`getEnvelope`函数来获得边界框的坐标，我们只需要简单地使用坐标构造出边界框即可。

## 用户使用案例：pymilvus+restful

### pymilvus

```python
from pymilvus import MilvusClient，DataType
collection_name = "geo_point_collection"
milvus_client = MilvusClient("http://localhost:19530")

schema = MilvusClient.create_schema(
	auto_id = False，
	enable_dynamic_field = False，
)

# create geo fields
schema.add_field(name = "id"，datatype = DataType.INT64，is_primary = True)
schema.add_field(name = "name"，datatype = DataType.VARCHAR，max_length = 255)
schema.add_field(name = "location"，datatype = DataType.GEOMETRY)

# create collection
milvus_client.create_collection(collection_name，schema)

# insert 
data =[
	{"id": 1001，"name": "Shop A"，"location": "POINT(116.4 39.9)"}，
	{"id": 1002，"name": "Shop B"，"location": "POINT(116.5 39.8)"}，
	{"id": 1003，"name": "Shop C"，"location": "POINT(116.6 39.7)"}
]

milvus_client.insert(collection_name，data)
# query
# 1. spatial relationship
  
# The usage of spatial relationship querys are like this:
# ST_XXX({field_name}， {wkt_string})
# where {field_name} is the name of the field that you want to query，
# and {wkt_string} is the wkt string of the geometry.
# So the results of the query are the specific geometry objects in the field that meet the spatial relationship to a given geometry
# Include:
# ST_Contains，ST_Within，ST_Covers，ST_Intersects
# ST_Disjoint，ST_Equals，ST_Crosses，ST_Overlaps，ST_Touches

# ST_Within
within_wkt = "POLYGON((116.4 39.9，116.5 39.9，116.6 39.9，116.4 39.9))"
results = milvus_client.query(
	collection_name=collection_name，
	filter=f"ST_Within(location， '{within_wkt}')"，
	output_fields=["name"，"location"]
)
print(results)

  
# ST_Covers
covers_wkt = "POLYGON((116.4 39.9，116.5 39.9，116.6 39.9，116.4 39.9))"
results = milvus_client.query(
	collection_name=collection_name，
	filter=f"ST_Covers(location， '{covers_wkt}')"，
	output_fields=["name"，"location"]
)
print(results)

# 2. distance relationship
# The usage of distance relationship querys are like this:
# ST_DWithin
point_wkt = "POINT(116.5 39.9)"
distance = 10000 # meters
results = milvus_client.query(
	collection_name=collection_name，
	filter=f"ST_DWithin(location， '{point_wkt}'， {distance})"，
	output_fields=["name"，"location"]
)
print(results)

# ST_Distance
target_point = "POINT(116.5 39.9)"
results = milvus_client.query(
	collection_name=collection_name，
	filter=f"ST_Distance(location， '{target_point}') < 10000"，
	output_fields=["name"，"ST_Distance(location， '{target_point}')"] # Here we may need to support the calculation of the result
)
print(results)

# 3. Others
# These queries need to do some calcutions in filter or output_fields like ST_Distance above，so we need to support this function.
# But it may be not easy，and confilict with the current implementation of the filter and output_fields.

# 4. hybrid query
# We may sometimes use a specific condition to query and get some caculation results，consider the following case:
# Hybrid query example: filter by coverage condition and calculate area for qualifying polygons

coverage_wkt = "POLYGON((116.3 39.8，116.7 39.8，116.7 40.0，116.3 40.0，116.3 39.8))"
results = milvus_client.query(
collection_name=collection_name，
	filter=f"ST_Covers(location， '{coverage_wkt}')"， # Filter polygons that cover the specified area，here the location may not be a point，but can be a polygon
	output_fields=["name"， "location"， "ST_Area(location)"] # Calculate area for qualifying polygons
)

print("Polygons covering the specified area and their areas:")
for result in results:
	print(f"Name: {result['name']}， Area: {result['ST_Area(location)']} square meters")
```

### restful

根据pymilvus的示例即可大致推断出使用方法.
```restful
curl --request POST \
    --url "${CLUSTER_ENDPOINT}/v2/vectordb/entities/query" \
    --header 'accept: application/json' \
    --header 'content-type: application/json' \
    -d '{
        "dbName": "_default"，
        "collectionName": "geo_point_collection"，
        "filter": "ST_Within(location， ''POLYGON((116.4 39.9，116.5 39.9，116.6 39.9，116.4 39.9))'')"，
        "outputFields": ["name"， "location"]
    }'
```

对查找结果进行计算:

```restful
curl --request POST \
    --url "${CLUSTER_ENDPOINT}/v2/vectordb/entities/query" \
    --header 'accept: application/json' \
    --header 'content-type: application/json' \
    -d '{
        "dbName": "_default"，
        "collectionName": "geo_point_collection"，
        "filter": "ST_Distance(location， ''POINT(116.5 39.9)'') < 10000"，
        "outputFields": ["name"， "ST_Distance(location， ''POINT(116.5 39.9)'')"]
    }'
```

结合上述可能的使用场景，在查询和删除的场景中，针对filter和output_field字段，需要支持查询时计算和返回计算结果。

## 实现方案

### 依赖安装

我们需要的第三方依赖有：

1. geos：核心依赖，用于提供几何对象的封装和所有查询算子
2. proj：坐标转换库，用于将用户输入坐标转换为geos使用的平面坐标
3. libspatialindex：索引库，用于构建索引，也直接提供了基于索引的查询功能

### 数据插入

 在Milvus中，数据插入的操作是以`shard`为单位来完成的。用户可以选择为每个collection使用多少个shard，每个shard会被映射到一个virtual channel，每个virtual channel被分配给一个physical channel，一个physical channel可以被分配多个virtual channel。每个physical channel被绑定到一个StreamingNode。当数据在Proxy层验证后，包裹该数据的Msg会被proxy拆分为多个package，并根据特定规则分发到不同的shard中，然后数据会被送到该shard对应的pchannel中，传递给SN。

SN给每个package都打上一个时间戳来建立整体的操作顺序，它还会进行一些负载均衡的操作。SN首先将数据写入WAL，并将WAL划分为多个segment。当某个segment的WAL记录被milvus处理完后，会触发刷新操作，数据被最终写入对象存储。这些操作都是由StreamingNode来完成的。

首先需要增加相关proto定义中的数据类型，现在的代码中已经完成了。

接下来按照从上到下的顺序为geo数据结构添加支持：

1. 首先需要添加proxy解析层对geometry字段的验证，需要修改internal/proxy/validate.go的内容，可以参考原pr的实现思路
2. 之后的数据打包并发放到StreamingNode等链路不需要修改，因为都包装成了统一的fielddata进行传输
3. 内存中的数据结构，涉及到internal/storage的相关内容
4. c++核心层：需要新建一个新的Geometry类，封装WKB数据，并提供必要的接口，在本阶段先不提供查询算子功能。还有一些其他的与数据类型相关的内容，如Type，Array，Fielddata等.

对于其中涉及到的Segment，有GrowingSegment和SealedSegment两种segment需要处理，分别涉及:SegmentGrowingImpl.cpp，SegmentSealedImpl.cpp

综上，需要修改的文件有:

proxy层:validate_util.go
strorage层:
1. data_codec.go:负责序列化和反序列化
2. data_sorter.go:兼容geo数据类型
3. insert_data.go:添加geo数据支持和相关工厂函数
4. payload.go:geo的二进制读写支持
5. payload_writer/reader:相关接口
6. serde.go:与Arrow，Parquet的相关转化
7. util，print_binlog，test:相关工具和测试

c++ core层:
common:
1. Types.h:数据类型枚举
2. Geometry.h:新建文件，负责geo数据类型的定义，对象的构造，序列化，相关算子等.
3. fieldData.h/cpp/interface.h:管理字段数据的内存表示，填充，访问.需要支持对geo数据的处理
4. array.h:新增对geo的支持
5. chunk.h/writer:实现geo的二进制块支持
6. vectorTrait.h:增加geo数据的判断逻辑，和json等地位相同

segcore:
1. chunkSegmentSealedImpl.cpp:负责Sealed Segment的数据加载，索引，删除等，需要在一些case里加入对geo类型的处理
2. ConcunrentVector.cpp:这里是insert数据的存储位置，需要添加geo数据支持
3. InsertRecord.h:在append_field_meta等方法中添加geo支持
4. SegmentGrowingImpl.cpp/h: 在Insert等方法中添加geo支持
5. SegmentSealedImpl.go:在数据加载等case中添加geo支持
6. Utils.cpp:相关工具的geo支持

storage:
1. Event.cpp:序列化部分添加geo支持
2. Util:在相关case中添加geo支持

对于Geometry.h的核心设计，初步的想法如下:

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
| + GetArea()， Length()    |
| + SpatialOps (==， ∩， ∪)  |
| + SRID Transform()       |
| + isValid()， repair()    |
+--------------------------+
```

算子部分没有写全，成员变量部分分别表示当前geo操作的上下文，geo对象，wkb数据，数据大小，以及坐标类型。
### 查询

Milvus中每个Collection会被分为多个segment。Segment的类型有两种：Growing和Sealed，前者代表该segment还可以继续写入数据，后者代表该segment是只读的。在Segment的大小达到一定规模或者经过一段时间后，growing segment会被刷写到磁盘中，成为sealed segment。

在Milvus的运作过程中，StreamingNode加载了growing segment而QueryNode加载sealed segment。当一个查询请求到来时，会被代理并行分发给所有持有该数据shard的SN，每个SN都会产生自己的查询Plan，查询它的growing data，并同时和QueryNode进行通信，查询sealed segment中的数据。

为了支持查询需要修改的部分较少，主要涉及到查询表达式的解析:

1. Plan.g4:这是最关键的一个文件，需要给它加上对上述查询算子的描述
2. internal/core/src/expr/exec/expression/Expr.cpp: 表达式求值等，需要增加对geo的支持
3. 为GIS查询算子增加表达式节点，就像JSON的一系列处理函数

此外，还需要给milvus额外新加一个对查询结果进行计算的功能，例如:
```python
results = milvus_client.query(
collection_name=collection_name，
	filter=f"ST_Covers(location， '{coverage_wkt}')"， 
	output_fields=["name"， "location"， "ST_Area(location)"] 
)
```
此处的我们需要获得符合查询条件的location的面积，由于用户希望直接看到面积结果，所以我们应该直接返回经过ST_Area处理后的面积数据，而不是location field的原始数据。

为此，我们还需要修改:

1. 扩展SQL解析层以识别聚合函数:parser_visitor.go
2. 为聚合函数添加新的proto定义:plan.proto
3. 扩展output_field字段的处理能力:internal/proxy/util.go
4. createPlan添加聚合函数支持用于处理需要计算的字段:task_query.go
5. c++层处理引擎需要新增功能:internal/core/src/exec/aggregation
6. 结果处理:internal/proxy/default_limit_reducer.go
7. 计算算子具体实现

由于第二个功能较为复杂，所以我们这部分可以分成两个阶段来做:第一阶段只支持常规情况下的expr，第二阶段再对output_field做聚合函数支持。

### 索引

Milvus中的索引建立由DataNode来完成。
当SDK客户端发送建立索引的请求后，Milvus服务端并不会立刻进行索引的建立，而是将此请求写入日志，并通过一个channel将建立索引的信息发送给Datacoord。
Datacoord会监听该channel，当它发现该请求后，会创建一个任务然后至于调度队列中，当任务被调度后，发送任务到一个DataNode。该DataNode将要建立索引的数据从对象存储中加载到内存，然后建立索引，再将数据和索引写入对象存储。需要注意的是，collection中有多个segment，DataNode会为每一个已经flush的segment建立单独的索引。
由于R-Tree是动态索引，所以可以考虑自动索引构建，而不是让用户手动构建索引。现在有几个方案：
1. 在插入数据时构建索引：实现起来较为简单，但是在bulk insert时会造成严重的性能下降
2. 在flush操作时建立索引：在streaming node进行flush操作时，自动为数据建立索引

为了支持Geometry类型的索引，大体需要修改：

3. 参数检查阶段（Go 层）
`pkg/util/typeutil/schema.go`：添加 `IsGeometryType` 判断函数 
`internal/proxy/task_index.go` ：更新 `parseIndexParams` 支持 Geo 类型 
`pkg/util/paramtable/autoindex_param.go`：在 `AutoIndexConfig` 中添加 Geo 类型索引配置项
2. C++ Core 层索引工厂与实现
`internal/core/src/index/IndexFactory.cpp` 增加对 R-Tree 或 Geo 类型索引的支持
`internal/util/indexparamcheck/index_type.go` 更新索引类型校验逻辑
`internal/cpre/src/query/ScalarIndex.h` 添加 Scalar Index 的辅助函数
`internal/core/src/segcore/FieldIndexing.cpp` 更新字段索引构建逻辑 

## 排期

| 功能                    | 日期        |
| --------------------- | --------- |
| 正确安装依赖，并构建milvus      | 7-14      |
| 数据插入支持，能够持久化到对象存储     | 7-15～7-18 |
| 正确查询到数据，并输出geometry字段 | 7-21～7-25 |
| 支持聚合函数，对输出字段进行计算      | 7-28～8-1  |
| 能够正确构建索引              | 8-4～8-8   |
| 能够在索引上进行搜索，并通过对比测试来验证 | 8-11～8-15 |
| 代码测试与性能测试             | 8-18～8-22 |
| 项目文档                  | 8-25～8-29 |
