# 如何通过 Oozie 在 YARN 上运行 Spark 作业
了解 Oozie

**标签:** 分析

[原文链接](https://developer.ibm.com/zh/articles/run-spark-job-yarn-oozie/)

developerWorks 中国网站编辑组

发布: 2016-01-26

* * *

Apache Oozie 是一个用于管理 Apache Hadoop 作业的工作流调度程序。Oozie 将多个作业按顺序组合到一个逻辑工作单元中，作为操作的有向非循环图 (DAG)。Oozie 可靠、可伸缩、可扩展且与 Hadoop 堆栈紧密集成，使用 YARN 作为其架构中心。它开箱即用地提供了多种 Hadoop 作业类型，比如 Java map-reduce、Pig、Hive、Sqoop 和 DistCp，以及特定于系统的作业，比如 Java 程序和 shell 脚本。

Apache Spark 是一个快速的内存型数据处理引擎，它拥有简洁而又富于表达的 API，使您能够高效地执行流处理、机器学习或需要对数据集的快速迭代式访问的 SQL 工作负载。Hadoop 的基于 YARN 的架构为 Spark 共享通用的集群和数据集提供了基础。

Oozie 4.2.0 中提供了一种新的操作类型，该类型将 Spark 作业编织到您的工作流中。工作流等待 Spark 作业完成之后再继续执行下一个操作。本文将展示如何使用新的 Spark 操作在 IBM Open Platform with Apache Hadoop (IOP) 4.1 上运行 Spark 作业。

考虑一个简单的字数统计应用程序，它在一个文本文件集中创建一种文字分布。这个应用程序（在 Spark Java API 中编写）可用作您的工作流中的 Spark 作业，是一个不错例子。下面的列表大体确定了 Spark 驱动程序必须执行的操作：

1. 读取输入的一组文本文档。
2. 统计每个字出现的次数。
3. 按字数排序，并以 CSV 格式和降序输出结果。

以下各节将介绍如何使用 Oozie 在 YARN 上调度和启动这个 Spark 应用程序。本文末尾会给出一个完整的程序清单。

1. **创建一个工作流定义：workflow.xml**

    下面这个简单的工作流定义了执行一个 Spark 作业的配置方法：





    ```
    <workflow-app xmlns='uri:oozie:workflow:0.5' name='SparkWordCount'>
    <start to='spark-node' />
    <action name='spark-node'>
    <spark xmlns="uri:oozie:spark-action:0.1">
    <job-tracker>${jobTracker}</job-tracker>
    <name-node>${nameNode}</name-node>
    <prepare>
    <delete path="${nameNode}/user/${wf:user()}/${examplesRoot}/output-data"/>
    </prepare>
    <master>${master}</master>
    <name>Spark-Wordcount</name>
    <class>com.ibm.biginsights.oozie.examples.WordCountSparkMain</class>
    <jar>${nameNode}/user/${wf:user()}/${examplesRoot}/lib/examples-1.0.jar</jar>
    <spark-opts>–conf spark.driver.extraJavaOptions=-Diop.version=4.1.0.0</spark-opts>
    <arg>${nameNode}/user/${wf:user()}/${examplesRoot}/input-data</arg>
    <arg>${nameNode}/user/${wf:user()}/${examplesRoot}/output-data</arg>
    </spark>
    <ok to="end" />
    <error to="fail" />
    </action>
    <kill name="fail">
    <message>Workflow failed, error
    message[${wf:errorMessage(wf:lastErrorNode())}]
    </message>
    </kill>
    <end name='end' />
    </workflow-app>

    ```





    Show moreShow more icon

    一些元素定义如下：


    - **prepare** 元素指定一个要在启动作业之前删除或创建的路径列表。这些路径必须以 hdfs://host\_name:port\_number 开头。
    - **master** 元素指定 Spark Master 的 URL；例如 spark://host:port、mesos://host:port、yarn-cluster、yarn-master 或 local。对于 Spark on YARN 模式，在 master 元素中指定的 yarn-client 或 yarn-cluster。在这个示例中，master=yarn-cluster。
    - **name** 元素指定 Spark 应用程序的名称。
    - **class** 元素指定 Spark 应用程序的主要类。
    - **jar** 元素指定一个逗号分隔的 JAR 文件列表。
    - **spark-opts** 元素（如果存在）包含一个可通过指定 ‘-conf key=value’ 传递给 Spark 驱动程序的 Spark 配置选项列表。
    - **arg** 元素包含可传递给 Spark 应用程序的参数。



    有关 Oozie 中的 Spark XML 模式的详细信息，请参阅 [https://oozie.apache.org/docs/4.2.0/DG\_SparkActionExtension.html](https://oozie.apache.org/docs/4.2.0/DG_SparkActionExtension.html) 。

2. **创建一个 Oozie 作业配置：job.properties**





    ```
    nameNode=hdfs://nn:8020
    jobTracker=rm:8050
    master=yarn-cluster
    queueName=default
    examplesRoot=spark-example
    oozie.use.system.libpath=true
    oozie.wf.application.path=${nameNode}/user/${user.name}/${examplesRoot}

    ```





    Show moreShow more icon

3. **创建一个 Oozie 应用程序目录**

    创建一个包含工作流定义和资源的应用程序目录结构，如下面的示例所示：





    ```
    +-~/spark-example/
    +-job.properties
    +-workflow.xml
    +-lib/
    +-example-1.0.jar

    ```





    Show moreShow more icon

    example-1.0.jar 文件包含 Spark 应用程序。

4. **下载 spark-assembly.jar 文件**

    浏览 HDFS NameNode 用户界面，下载 spark-assembly.jar 文件。

    1. 从 Ambari 控制台中选择 **HDFS** ，然后选择 **Quick Links** –\> **NameNode UI** 。
    2. 单击 **Utilities** –\> **Browse the file system** 。
    3. 在 Hadoop 文件资源管理器中，导航到 /iop/apps/4.1.0.0/spark/jars，选择 spark-assembly.jar，单击 **Download** 并保存该文件。
    4. 将下载的 spark-assembly.jar 文件转移到 lib 目录，这会得到以下目录结构：





        ```
        +-~/spark-example/
            +-job.properties
            +-workflow.xml
            +-lib/
            +-example-1.0.jar
            +-spark-assembly.jar

        ```





        Show moreShow more icon
5. **将应用程序复制到 HDFS**

    将 spark-example/ 目录复制到 HDFS 中的用户 HOME 目录。确保 HDFS 中的 spark-example 位置与 job.properties 中的 oozie.wf.application.path 值匹配。





    ```
    $ hadoop fs -put spark-example spark-example

    ```





    Show moreShow more icon

6. **运行示例作业**

    运行以下命令来提交 Oozie 作业：





    ```
    $cd ~/spark-example
    $oozie job -oozie http://oozie-host:11000/oozie -config ./job.properties –run
    job: 0000012-151103233206132-oozie-oozi-W

    ```





    Show moreShow more icon

    检查工作流作业状态：





    ```
    $ oozie job –oozie http://oozie-host:11000/oozie -info 0000012-151103233206132-oozie-oozi-W

    Job ID : 0000012-151103233206132-oozie-oozi-W
    ————————————————————————————————————————————
    Workflow Name : SparkWordCount
    App Path : hdfs://bdvs1211.svl.ibm.com:8020/user/root/spark-example
    Status : SUCCEEDED
    Run : 0
    User : root
    Group : –
    Created : 2015-11-04 15:19 GMT
    Started : 2015-11-04 15:19 GMT
    Last Modified : 2015-11-04 15:23 GMT
    Ended : 2015-11-04 15:23 GMT
    CoordAction ID: –

    Actions
    ————————————————————————————————————————————
    ID Status Ext ID Ext Status Err Code
    ————————————————————————————————————————————
    0000012-151103233206132-oozie-oozi-W@:start: OK – OK –
    0000012-151103233206132-oozie-oozi-W@spark-node OK job_1446622088718_0022 SUCCEEDED –
    0000012-151103233206132-oozie-oozi-W@end OK – OK –
    ————————————————————————————————————————————

    ```





    Show moreShow more icon

    **完整的 Java 程序**





    ```
    public static void main(String[] args) {
    if (args.length < 2) {
    System.err.println("Usage: WordCountSparkMain <file> <file>");
    System.exit(1);
    }
    String inputPath = args[0];
    String outputPath = args[1];
    SparkConf sparkConf = new SparkConf().setAppName("Word count");
    try (JavaSparkContext ctx = new JavaSparkContext(sparkConf)) {
    JavaRDD<String> lines = ctx.textFile(inputPath, 1);
    JavaRDD<String> words = lines.flatMap(new FlatMapFunction<String,
    String>() {
    private static final long serialVersionUID = 1L;
    public Iterable<String> call(String sentence) {
    List<String> result = new ArrayList<>();
    if (sentence != null) {
    String[] words = sentence.split(" ");
    for (String word : words) {
    if (word != null && word.trim().length() > 0) {
    result.add(word.trim().toLowerCase());
    }
    }
    }
    return result;
    }
    });
    JavaPairRDD<String, Integer> pairs = words.mapToPair(new
    PairFunction<String, String, Integer>() {
    private static final long serialVersionUID = 1L;
    public Tuple2<String, Integer> call(String s) {
    return new Tuple2<>(s, 1);
    }
    });

    JavaPairRDD<String, Integer> counts = pairs.reduceByKey(new
    Function2<Integer, Integer, Integer>() {
    private static final long serialVersionUID = 1L;
    public Integer call(Integer a, Integer b) {
    return a + b;
    }
    }, 2);
    JavaPairRDD<Integer, String> countsAfterSwap = counts.mapToPair(new
    PairFunction<Tuple2<String, Integer>, Integer, String>() {
    private static final long serialVersionUID = 2267107270683328434L;
    @Override
    public Tuple2<Integer, String> call(Tuple2<String, Integer> t)
    throws Exception {
    return new Tuple2<>(t._2, t._1);
    }
    });
    countsAfterSwap = countsAfterSwap.sortByKey(false);
    counts = countsAfterSwap.mapToPair(new PairFunction<Tuple2<Integer,
    String>, String, Integer>() {
    private static final long serialVersionUID = 2267107270683328434L;
    @Override
    public Tuple2<String, Integer> call(Tuple2<Integer, String> t)
    throws Exception {
    return new Tuple2<>(t._2, t._1);
    }
    });
    JavaRDD<String> results = counts.map(new Function<Tuple2<String,
    Integer>, String>() {
    @Override
    public String call(Tuple2<String, Integer> v1) throws Exception {
    return String.format("%s,%s", v1._1, Integer.toString(v1._2));
    }
    });
    results.saveAsTextFile(outputPath);
    }
    }
    }

    ```





    Show moreShow more icon