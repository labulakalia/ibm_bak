# 操作技巧：将 Spark 中的文本转换为 Parquet 以提升性能
Spark 进阶

**标签:** 分析

[原文链接](https://developer.ibm.com/zh/articles/ba-parquet-for-spark-sql/)

JESSE CHEN

发布: 2016-01-19

* * *

列式存储布局（比如 Parquet）可以加速查询，因为它只检查所有需要的列并对它们的值执行计算，因此只读取一个数据文件或表的小部分数据。Parquet 还支持灵活的压缩选项，因此可以显著减少磁盘上的存储。

如果您在 HDFS 上拥有基于文本的数据文件或表，而且正在使用 Spark SQL 对它们执行查询，那么强烈推荐将文本数据文件转换为 Parquet 数据文件，以实现性能和存储收益。当然，转换需要时间，但查询性能的提升在某些情况下可能达到 30 倍或更高，存储的节省可高达 75%！

已有文章介绍使用 Parquet 存储为 BigSQL、Hive 和 Impala 带来类似的性能收益，本文将介绍如何编写一个简单的 Scala 应用程序，将现有的基于文本的数据文件或表转换为 Parquet 数据文件，还将展示给 Spark SQL 带来的实际存储节省和查询性能提升。

## 让我们转换为 Parquet 吧！

Spark SQL 提供了对读取和写入 Parquet 文件的支持，能够自动保留原始数据的模式。Parquet 模式通过 Data Frame API，使数据文件对 Spark SQL 应用程序 “不言自明”。当然，Spark SQL 还支持读取已存储为 Parquet 的现有 Hive 表，但您需要配置 Spark，以便使用 Hive 的元存储来加载所有信息。在我们的示例中，不涉及 Hive 元存储。

以下 Scala 代码示例将读取一个基于文本的 CSV 表，并将它写入 Parquet 表：

```
def convert(sqlContext: SQLContext, filename: String, schema: StructType, tablename: String) {
      // import text-based table first into a data frame
      val df = sqlContext.read.format("com.databricks.spark.csv").
        schema(schema).option("delimiter", "|").load(filename)
      // now simply write to a parquet file
      df.write.parquet("/user/spark/data/parquet/"+tablename)
}

// usage exampe -- a tpc-ds table called catalog_page
schema= StructType(Array(
          StructField("cp_catalog_page_sk",        IntegerType,false),
          StructField("cp_catalog_page_id",        StringType,false),
          StructField("cp_start_date_sk",          IntegerType,true),
          StructField("cp_end_date_sk",            IntegerType,true),
          StructField("cp_department",             StringType,true),
          StructField("cp_catalog_number",         LongType,true),
          StructField("cp_catalog_page_number",    LongType,true),
          StructField("cp_description",            StringType,true),
          StructField("cp_type",                   StringType,true)))
convert(sqlContext,
          hadoopdsPath+"/catalog_page/*",
          schema,
          "catalog_page")

```

Show moreShow more icon

上面的代码将会读取 hadoopdsPath+”/catalog\_page/\* 中基于文本的 CSV 文件，并将转换的 Parquet 文件保存在 /user/spark/data/parquet/ 下。此外，转换的 Parquet 文件会在 gzip 中自动压缩，因为 Spark 变量 spark.sql.parquet.compression.codec 已在默认情况下设置为 gzip。您还可以将压缩编解码器设置为 uncompressed、snappy 或 lzo。

## 转换 1 TB 数据将花费多长时间？

50 分钟，在一个 6 数据节点的 Spark v1.5.1 集群上可达到约 20 GB/分的吞吐量。使用的总内存约为 500GB。HDFS 上最终的 Parquet 文件的格式为：

```
...
/user/spark/data/parquet/catalog_page/part-r-00000-9ff58e65-0674-440a-883d-256370f33c66.gz.parquet
/user/spark/data/parquet/catalog_page/part-r-00001-9ff58e65-0674-440a-883d-256370f33c66.gz.parquet
...

```

Show moreShow more icon

### 存储节省

以下 Linux 输出显示了 TEXT 和 PARQUET 在 HDFS 上的大小比较：

```
% hadoop fs -du -h -s /user/spark/hadoopds1000g
    897.9 G  /user/spark/hadoopds1000g
    % hadoop fs -du -h -s /user/spark/data/parquet
    231.4 G  /user/spark/data/parquet

```

Show moreShow more icon

1 TB 数据的存储节省了将近 75%！

### 查询性能提升

Parquet 文件是自描述性的，所以保留了模式。要将 Parquet 文件加载到 DataFrame 中并将它注册为一个 temp 表，可执行以下操作：

```
val df = sqlContext.read.parquet(filename)
      df.show
      df.registerTempTable(tablename)

```

Show moreShow more icon

要对比性能，然后可以分别对 TEXT 和 PARQUET 表运行以下查询（假设所有其他 tpc-ds 表也都已转换为 Parquet）。您可以利用 spark-sql-perf 测试工具包来执行查询测试。举例而言，现在来看看 TPC-DS 基准测试中的查询 #76，

```
("q76", """
            | SELECT
            |    channel, col_name, d_year, d_qoy, i_category, COUNT(*) sales_cnt,
            |    SUM(ext_sales_price) sales_amt
            | FROM(
            |    SELECT
            |        'store' as channel, ss_store_sk col_name, d_year, d_qoy, i_category,
            |        ss_ext_sales_price ext_sales_price
            |    FROM store_sales, item, date_dim
            |    WHERE ss_store_sk IS NULL
            |      AND ss_sold_date_sk=d_date_sk
            |      AND ss_item_sk=i_item_sk
            |    UNION ALL
            |    SELECT
            |        'web' as channel, ws_ship_customer_sk col_name, d_year, d_qoy, i_category,
            |        ws_ext_sales_price ext_sales_price
            |    FROM web_sales, item, date_dim
            |    WHERE ws_ship_customer_sk IS NULL
            |      AND ws_sold_date_sk=d_date_sk
            |      AND ws_item_sk=i_item_sk
            |    UNION ALL
            |    SELECT
            |        'catalog' as channel, cs_ship_addr_sk col_name, d_year, d_qoy, i_category,
            |        cs_ext_sales_price ext_sales_price
            |    FROM catalog_sales, item, date_dim
            |    WHERE cs_ship_addr_sk IS NULL
            |      AND cs_sold_date_sk=d_date_sk
            |      AND cs_item_sk=i_item_sk) foo
            | GROUP BY channel, col_name, d_year, d_qoy, i_category
            | ORDER BY channel, col_name, d_year, d_qoy, i_category
            | limit 100

```

Show moreShow more icon

查询时间如下：

```
TIME               TEXT     PARQUET
            Query time (sec)    698          21

```

Show moreShow more icon

查询 76 的查询时间从将近 12 分钟加速到不到半分钟，提高了 30 倍！

本文翻译自： [How-to: Convert Text to Parquet in Spark to Boost Performance](https://developer.ibm.com/hadoop/2015/12/03/parquet-for-spark-sql/)（2016-01-19）