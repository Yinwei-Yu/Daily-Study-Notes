**Parquet 格式（Parquet File Format）** 是一种**列式存储格式（Columnar Storage Format）**，广泛用于大数据处理系统中，比如 Apache Hadoop、Apache Spark、Presto、Hive、Impala 等。

## 📌 什么是 Parquet？

Parquet 是一个开源的、与语言和平台无关的二进制文件格式，设计目标是高效地存储和读取大规模数据集。它支持嵌套结构的数据，并且具备良好的压缩和编码能力，非常适合用于分析型查询（OLAP）场景。

---

## 🔍 主要特点

| 特性 | 描述 |
|------|------|
| **列式存储（Columnar Storage）** | 数据按列存储，适合只访问部分字段的查询，提高 I/O 效率 |
| **高效的压缩和编码** | 支持多种编码方式（如 RLE、字典编码等），节省存储空间并提升性能 |
| **Schema 演化（Schema Evolution）** | 支持添加新列、修改字段类型等操作，不影响旧数据 |
| **支持复杂嵌套结构** | 如 List、Map、Struct 等结构，适用于 JSON、Avro 等数据模型 |
| **跨平台兼容性强** | 被主流大数据工具（Spark、Flink、Hive、Presto、BigQuery 等）广泛支持 |
| **可切分（Splittable）** | 文件可以被分割成多个块进行分布式处理 |

---

## 🧱 Parquet 文件结构（简化版）

一个 Parquet 文件由以下几个主要部分组成：

1. **Row Group（行组）**
   - 将数据按行划分成多个“行组”
   - 每个 Row Group 包含若干行
2. **Column Chunk（列块）**
   - 每个 Row Group 内部又按列切割成 Column Chunk
3. **Page（页）**
   - Column Chunk 再分为多个 Page，每个 Page 可以独立解码
4. **Metadata（元数据）**
   - 包括 Schema、统计信息（min/max/count）、编码方式等

---

## 🧪 示例：Parquet vs CSV/JSON

假设你有如下数据：

| id | name     | age | gender |
|----|----------|-----|--------|
| 1  | Alice    | 25  | F      |
| 2  | Bob      | 30  | M      |
| 3  | Charlie  | 28  | M      |

- 如果你只想查询 `age` 和 `gender`：
  - CSV/JSON 需要读取整行数据；
  - Parquet 只需要读取这两个列的数据，节省大量 I/O。

---

## 🛠️ 常见使用场景

- **ETL 处理**：将原始数据转换为 Parquet 存储，提升后续查询效率。
- **数据湖（Data Lake）**：作为底层存储格式，配合 AWS S3、Azure Blob、HDFS 等使用。
- **数仓建模**：在 Hive、Iceberg、Delta Lake 中作为默认或推荐存储格式。
- **BI 分析**：在 Presto、Trino、ClickHouse 等中加速报表生成。
- **机器学习特征工程**：快速加载高维特征列，跳过不相关字段。

---

## 💡 Parquet 与其他格式对比

| 格式       | 是否列式 | 压缩率 | 易读性 | 兼容性 | 适用场景         |
|------------|-----------|---------|--------|--------|------------------|
| Parquet    | ✅        | 高      | ❌     | 高     | OLAP、数据仓库、数据湖 |
| ORC        | ✅        | 高      | ❌     | 高     | Hive 生态         |
| Avro       | ❌（行式）| 中      | ❌     | 高     | 日志、序列化数据 |
| JSON       | ❌        | 低      | ✅     | 高     | 小数据、调试     |
| CSV        | ❌        | 低      | ✅     | 高     | 导入导出         |
| XML        | ❌        | 低      | ✅     | 低     | 传统系统         |

---

## 🧰 如何创建和查看 Parquet 文件？

### 使用 Python（PyArrow 或 Pandas）

```python
import pandas as pd

# 写入 Parquet
df = pd.DataFrame({'id': [1, 2, 3], 'name': ['Alice', 'Bob', 'Charlie']})
df.to_parquet('data.parquet')

# 读取 Parquet
df = pd.read_parquet('data.parquet')
```

### 使用 Spark（Scala / PySpark）

```python
# 写入
df.write.parquet("output_path")

# 读取
df = spark.read.parquet("input_path")
```

### 查看内容（命令行）

使用 `parquet-tools`（Java 工具）：

```bash
parquet-tools cat data.parquet
parquet-tools schema data.parquet
```

---

## 📚 相关项目 / 生态

- **Apache Arrow**：内存中的列式数据结构，常与 Parquet 结合使用。
- **Delta Lake / Iceberg / Hudi**：基于 Parquet 的表格式（Table Format），提供事务、版本控制等功能。
- **AWS Glue / Athena**：支持直接查询 Parquet 格式数据。
- **Apache Drill**：支持原生查询 Parquet 文件。

---

## ✅ 总结

Parquet 是现代大数据生态中非常重要的存储格式，特别适合：

- 列式访问为主的分析查询；
- 对性能和存储成本敏感的大规模数据集；
- 需要 Schema 灵活性和演化能力的场景。

如果你正在做数据湖、数仓或 ETL 工作，**Parquet 是一个几乎必选的格式**。
