# Db2 9 中的 pureXML：怎样查询您的 XML 数据？
灵活性使任务更容易

**标签:** 数据库

[原文链接](https://developer.ibm.com/zh/articles/dm-0606nicola/)

Matthias Nicola, Fatma Ozcan

发布: 2006-06-15

* * *

## 简介

Db2 中的 pureXML 支持为管理 XML 数据提供了高效且通用的功能。Db2 以 XML 数据自身固有的分层格式存储和处理这些数据，避免因为将 XML 存储为 CLOB 中的文本或将它映射为关系表而导致的性能和灵活性限制。与仅使用 XML 的数据库不同，Db2 V9 还提供了关系型数据与 XML 数据在数据库中的无缝集成 —— 甚至是表的某一行中的集成。这样的灵活性表现在语言支持中，使您可访问关系型数据、XML 数据，或者同时访问这两种数据。您可以通过以下四种可选方案中的任一种查询 XML：

- 普通 SQL（不包含 XQuery）
- SQL/XML，即嵌入了 XQuery 的 SQL
- XQuery 作为独立语言（不包含 SQL）
- 嵌入了 SQL 的 XQuery

关于使用 XQuery 和 SQL/XML 查询 XML 数据的介绍，请参阅 developerWorks 中的 前几期文章 ，”用 SQL 查询 Db2 XML 数据” 及 “使用 XQuery 查询 Db2 XML 数据”。本文假设您已经熟悉这两篇文章中介绍的基本概念。请注意，XPath 是 XQuery 的一种子语言，因此我们提到 XQuery 时，也暗中包含 XPath 语言。如果您使用过 Db2 XML Extender 中的 XSLT 样式表或位置路径，那么您应该已经了解 XPath。在很多时候，使用 XPath 足以提取 XML 值或表示 XML 谓词，因此，即便您还不熟悉 XQuery 的所有其他特性，也可以开始使用 XPath。

Db2 使您能够利用所有这些可选方案最大限度地提高生产力、使查询适应应用程序的需求。本文将说明的问题如下：

- 这四种可选方案的关键特征是什么？各有哪些优缺点？
- 您应该在哪种情况下选择哪种方案？

让我们先给出一个高度总结，然后再详细研究各个可选方案的细节和特定实例。

## 总结与指导原则

您可以在普通 XQuery、SQL/XML 或具有内置 SQL 的 XQuery 中表达许多查询。在特定情况下，您可能会发现其中之一能够比其他方案更为直观地表达您的查询逻辑。一般而言，查询 XML 的 “正确” 途径需要在 “逐个处理” 的基础上加以选择，需要考虑应用程序的具体需求和特征。但我们可总结出以下指导原则。

- **不带任何 XQuery 或 XPath 的普通 SQL** 仅对全文档检索以及整个文档的插入、删除、更新操作有用。文档的选择必须基于同一表中的非 XML 列。
- **带有嵌入在 SQL 内的 XQuery 或 XPath 语句的 SQL/XML** 提供了最广泛的功能性和最少的局限性。您可在 XML 列上表示谓词、提取文档片段、向 XQuery 表达式传递参数标记、使用全文本搜索、SQL 级聚集与分组，您还可以用一种灵活的方式将关系型数据与 XML 数据进行联合和连接。这种方案可很好地服务于绝大多数应用程序。即便您不是立即需要利用所有这些优势，可能仍然会考虑选择这种方案，从而使您的选择能够应对未来的扩展。
- **XQuery** 是一种强大的查询语言，专为查询 XML 数据而设计。同样，如果您的应用程序只需查询和操纵 XML 数据，且不涉及任何关系型数据，那么 XQuery 也是一种极为出色的选择方案。此方案有时可能较为简单直观。此外，如果您正从一个仅使用 XML 的数据库移植到 Db2 9，且已有 XQuery，那么您很可能愿意继续使用 XQuery。
- **嵌入了 SQL 的 XQuery** 在您希望利用关系型谓词和索引同时又想利用全文本搜索预先过滤随后将作为 XQuery 输入的 XML 列中的文档时不失为明智之选。嵌入在 XQuery 中的 SQL 允许您在 XML 列上运行外部函数。但若您需要执行带有分组和聚集的数据分析查询，那么 SQL/XML 是更好的选择。

无论您选择在一条语句中怎样结合 SQL 和 XQuery，Db2 都使用一种混合编译器来为整个查询生成及优化一种执行规划 —— 不会导致查询执行的性能损失。

下表总结了查询 XML 数据的四种不同可选方案的各自优点。

##### 表 1\. 总结

普通 SQLSQL/XML普通 XQuery嵌入了 SQL/XML 的 XQueryXML 谓词–++++++关系谓词++++–+XML 及关系谓词–++–++将 XML 与关系型相连接–++–++将 XML 与 XML 相连接–+++++转换 XML 数据–o++++插入、更新和删除++++––参数标记+++––全文本搜索+++–++XML 聚集与分组–++oo函数调用++++–++

在上表中，”-” 表示给定语言不支持某项特性；”+” 表示支持此特性，但存在更有效或更便捷的方式；”++” 表示给定语言极为适合表现该特性；最后，”o” 表示尽管可表现此特性，但从某种程度上来说，效果非常糟糕或者效率很低。

现在，让我们来定义一些示例数据和表，以查看具体的查询示例。

## 示例表和数据

在讨论 XML 查询可选方案的过程中，我们使用了以下三个表。dept 表有两列，列名分别为 unitID 和 deptdoc。dept 表的每一个行都描述了一家虚构企业的一个部门。unitID 列标识包含部门的单位（一个单位可能包含多个部门），deptdoc 列包含一个列出部门中的员工的 XML 文档。project 表只有一列，列名为 projectDoc，类型为 XML。projectDoc 表中的每一个行都包含一个描述特定项目的 XML 文档。一个项目可能会涉及到多个部门。为阐明混合关系型查询、XML 查询和连接，我们还提供了一个纯粹的关系型表 —— unit，其中列出了各单位的名称、主管等信息。一个单位可包含多个部门。

```
create table dept( unitID char(8), deptdoc xml)

```

Show moreShow more icon

unitIDdeptdocWWPR`<dept deptID="PR27"> <employee id="901"> <name>Jim Qu</name> <phone>408 555 1212</phone> </employee> <employee id="902"> <name>Peter Pan</name> <office>216</office> </employee> </dept>`WWPR`<dept deptID="V15"> <employee id="673"> <name>Matt Foreman</name> <phone>416 891 7301</phone> <poffice>216</office> </employee> <description>This dept supports sales world wide</description> </dept>`S-USE………

```
create table project(projectDoc xml)

```

Show moreShow more icon

projectDOC`<project ID="P0001"> <name>Hello World</name> <deptID>PR27</deptID> <manager>Peter Pan</manager> </project>``<project ID="P0009"> <name>New Horizon</name> <deptID>PR27</deptID> <deptID>V15</deptID> <manager>Matt Foreman</manager> <description>This project is brand new</description> </project>`

```
create table unit( unitID char(8) primary key not null, name char(20), manager char(20),...)

```

Show moreShow more icon

unitID名称主管…WWPRWorldwide MarketingJim Qu…S-USESales US East CoastTom Jones……………

## 普通 SQL

您可使用不带任何 XPath 或 XQuery 的普通 SQL，在不使用 XML 谓词的情况下读取整个文档。只要应用程序能识别出 XML 文档，从而基于同一表中的关系谓词进行全文档检索，那么这种方法就是一种非常简单易行的方法。例如，查询 1 检索 “WWPR” 这个单位的所有 department 文档：

##### 查询 1

```
select deptdoc
from dept
where unitID = 'WWPR';

```

Show moreShow more icon

同样，查询 2 返回 Jim Qu 所管理的所有部门的 XML 数据。

##### 查询 2

```
select deptdoc
from dept d, unit u
where d.unitID = u.unitID and u.manager= 'Jim Qu';

```

Show moreShow more icon

显而易见的缺陷就是您无法表示关于 XML 数据本身的谓词，也无法仅检索 XML 文档的片段。在我们的示例中，普通 SQL 不足以实现仅选择部门 PR27 且仅返回员工姓名这项任务。

若您的查询不需要在 XML 上应用谓词，并且总是返回整个 XML 文档，那么普通 SQL 足以满足您的需要。在本例中，将在 VARCHAR 或 CLOB 列中存储 XML，对于全文档插入及检索操作而言，这将给您带来一定的性能优势。

普通 SQL 没有使用 Db2 V9 支持的任何 SQL/XML 或 XQuery，但它依然允许在查询中使用全文本搜索条件。通过 Db2 Net Search Extender ，您可以在 XML 列上创建全文本索引来支持文本搜索 —— 从基本的关键字搜索到高级的 37 种语言的词干提取（stemming）、辞典（thesaurus）和模糊搜索。您还可以通过 path 表达式将文本搜索限定为仅搜索特定文档部分。以文本-索引查找为基础，查询 3 返回 /dept/description 目录下所有包含字符串 “sales” 的 department 文档：

##### 查询 3

```
select deptdoc
from dept
where CONTAINS(deptdoc,'SECTION("/dept/description") "sales" ')=1;

```

Show moreShow more icon

## SQL/XML（嵌入 SQL 中的 XQuery/XPath）

SQL/XML 是 SQL 语言标准的一部分，定义了一种新的 XML 数据类型以及多种查询、构造、验证和转换 XML 数据的函数。Db2 V8 已加入了许多 SQL/XML 发布函数，使用户能够使用 XMLELEMENT、XMLATTRIBUTE、XMLFOREST、XMLAGG 和其他一些函数将关系型数据构造为 XML。

Db2 9 添加了新的 SQL/XML 查询函数，包括 XMLQUERY、XMLTABLE 函数以及 XMLEXISTS 谓词。这些函数允许用户在 SQL 语句中嵌入 XQuery 或简单的 XPath 表达式。

如查询 4 所示，XMLQUERY 函数通常用在 select 子句中，用于从 XML 列中提取 XML 片段，而 XMLEXISTS 通常在 where 子句中使用，用于在 XML 数据上应用谓词。

##### 查询 4

```
select unitID, XMLQUERY('for $e in $d/dept/employee return $e/name/text()'
                   passing  d.deptdoc as "d")
from dept d
where  unitID LIKE 'WW%' and
    XMLEXISTS('$d/dept[@deptID = "V15"]' passing d.deptdoc as "d");

```

Show moreShow more icon

这个示例查询使用 XMLEXISTS 来选择 “WW” 部门 V15，并应用 XMLQUERY 来返回该 department 文档中的所有员工姓名。查询结果如下：

部门职员WWPRMatt Foreman

此查询还突出展示了以一种整合的方式使用 SQL/XML 查询 XML 数据和关系型数据的方法。select 子句从关系型列和 XML 列中检索数据，而 where 子句包含关系型谓词和 XML 语句。Db2 9.5 可以同时使用 XML 索引和关系索引来评估这些语句并最大化查询性能。

在 Db2 9.5 中，通过在 XMLEXISTS 和 XMLQUERY 函数中省略 “passing” 子句，可以更简单地编写这个查询。如果只将 XML 列 “deptdoc” 传递到 XQuery 中，那么可以在 XQuery 表达式中将该列引用为 $ deptdoc，而不使用 “passing” 子句。您可以在查询 5 中看到这一点。

##### 查询 5

```

select unitID, XMLQUERY('for $e in $DEPTDOC/dept/employee return $e/name/text()')
from dept d
where  unitID LIKE 'WW%' and
       XMLEXISTS('$DEPTDOC/dept[@deptID = "V15"]');
```

Show moreShow more icon

我们可以使用 XMLTABLE 函数来表示相同的查询，如查询 6 所示。在这种格式中，我们指定条件来限制输入数据并提取我们感兴趣的输出值。在查询 6 中，XMLTABLE 函数中的 XQuery 表达式标识在部门 V15 中工作的员工，COLUMNS 子句 (“name/text()”) 中的路径表达式返回他们的名称。它的输出与查询 4 和 5 相同。

##### 查询 6

```

select d.unitID, T.name
from dept d, XMLTABLE('$d/dept[@deptID="V15"]/employee' passing d.deptdoc as "d"
                       COLUMNS
                          name varchar(50) path 'name/text()'  ) as T
where unitID LIKE 'WW%';
```

Show moreShow more icon

### SQL/XML 的优点

SQL/XML 方案具有以下优点：

- 如果您有一个现有 SQL 应用程序，需要在各处逐步添加一些 XML 功能，那么 SQL/XML 具有一定的优势。
- 如果您是 SQL 的忠实爱好者，并且希望以 SQL 为首选语言 —— 因为它是您和您的团队最熟悉的语言，那么 SQL/XML 有优势。
- 如果您的查询需要同时从关系型列和 XML 列中返回数据，那么 SQL/XML 有优势。
- 如果您的查询需要如查询 3 所示的全文本搜索条件，那么 SQL/XML 有优势。
- 如果您希望将结果作为集合返回，且将丢失的 XML 元素表示为空，那么 SQL/XML 有优势。
- 如果您希望使用参数标记，那么 SQL/XML 有优势，因为 Db2 V9 XQuery 不支持外部参数。XMLQUERY、XMLTABLE 和 XMLEXISTS 中的传递机制允许您将一个 SQL 参数标记作为变量（$x）传递到嵌入式 XQuery 表达式中：

##### 查询 7

```
select deptdoc
from dept
where XMLEXISTS('$d/dept[@deptID =  $x]'
                passing deptdoc as "d", cast(? as char(8)) as "x");

```

Show moreShow more icon

- 对于那些需要整合关系型数据和 XML 数据的应用程序来说，SQL/XML 有优势。它提供了连接 XML 数据与关系型数据的最简便途径。以下示例选择拥有 unit 表中作为主管列出的那些员工的所有部门的 unit ID。此查询执行了一个关系型值（unit.manager）和 XML 值（//employee/name）之间的连接：

##### 查询 8

```
select u.unitID
from dept d, unit u
where XMLEXISTS('$d//employee[name = $m]'
              passing d.deptdoc as "d", u.manager as "m");

```

Show moreShow more icon

为完成此连接，我们将 unit manger 传递到 XMLEXISTS 谓词之中，这样实际连接条件就是一个 XQuery 谓词。反之，我们可将 XML 文档中的员工姓名提取出来，传递到 SQL 上下文中，从而使连接条件成为一个 SQL 语句：

##### 查询 9

```
select u.unitID
from dept d, unit u
where u.manager = XMLCAST(XMLQUERY('$d//employee/name '
                      passing d.deptdoc as "d") as char(20));

```

Show moreShow more icon

通常，查询 8 要优于查询 9，因为 XMLCAST 函数只需要一个输入值。但在我们的示例中，对于那些有着多名符合条件的员工的部门，此查询将失败。查询 9 使用了 XMLCAST，对于各文档只需进行一次 XML 值连接非常有用，因为它允许在 unit.manager 上使用关系型索引。此索引无法在查询 8 中使用，因为其连接条件并非关系型谓词，而是一个 XQuery 语句。

- SQL/XML 在 XML 分组和聚集方面有优势。XQuery 语言不提供显式的 group-by 结构。尽管可利用自连接在 XQuery 中表示分组和聚集，但那非常糟糕。作为示例，让我们计算一下按办公室分组的员工数量，换句话说，就是各办公室中的员工数量。查询 10 展示了普通 XQuery 中的实现方式。查询 10 中的 db2-fn:xmlcolumn 函数提供了 Db2 中的 XML 数据访问。它接受作为参数传入的 XML 列名，然后返回存储在该列中的 XML 值序列。使用 SQL/XML 函数（如 XMLTABLE 或 XMLQUERY）更简单，它们可使您轻松从 XML 列中提取数据项，然后使用熟悉的 SQL 概念在此基础上表示分组和聚集。查询 11 返回与查询 10 相同的逻辑结果，但使用的是 SQL 的 group by 子句。

##### 查询 10

```
XQUERY
for $o in distinct-values(db2-fn:xmlcolumn("DEPT.DEPTDOC")/dept/employee/office)
let $emps := db2-fn:xmlcolumn("DEPT.DEPTDOC")/dept/employee[office/text()=$o]
return
      <result><office>{$o}</office><cnt>{count($emps)}</cnt></result>;

Result:
<result><office>216</office><cnt>2</cnt></result>

```

Show moreShow more icon

##### 查询 11

```
select X.office, count(X.emp)
from dept, XMLTABLE ('$d/dept/employee' passing deptdoc as "d"
     COLUMNS
       emp         VARCHAR(30)   PATH  ' name',
       office     INTEGER          PATH  ' office ') as X
GROUP BY X.office;

Result:
216   2
  -     1

```

Show moreShow more icon

在查询 11 中，XMLTABLE 函数从带有 “emp” 和 “office” 这两列的表形式的所有文档中提取 /dept/employee/name 和 /dept/employee/office。对该表使用 SQL 的 group by 函数和聚集函数通常比使用生成相同结果的普通 XQuery 更有效。

请注意，在查询 11 中多获得了一行，这是因为 SQL 的 group by 函数还为 NULL 值生成一个组，我们的示例表中有一名员工没有办公室信息。查询 10 没有为该员工生成一行，因为 for 循环遍历的是不同的 office 值，这些值不包括丢失的办公室信息。

### SQL/XML 的缺点

- 对于将 XML 文档转换成另一个 XML 文档，SQL/XML 并非总是最佳选择。在仅需处理 XML 数据时，使用独立通常更合适，也更直观。
- 在 SQL/XML 中表示两个 XML 列 —— 或者更常见的情况，两个 XML 值之间的连接非常不便，用普通 XQuery 通常更高效、更直观。例如，查询 12 连接 dept 和 project 表的 XML 列，返回执行任一项目的员工。

##### 查询 12

```
select XMLQUERY('$d/dept/employee' passing d.deptdoc as "d")
from dept d, project p
where XMLEXISTS('$d/dept[@deptID=$p/project/deptID] '
                  passing d.deptdoc as "d", p.projectDoc as "p");

```

Show moreShow more icon

在查询 12 中，我们将每一个 dept 和 project 文档放入 XMLEXISTS 中传递给 XQuery，并在那里表示连接条件。您可能会发现，在普通 XQuery 中（查询 14）编写此类连接更为简单高效。

## XQuery 作为独立语言

让我们暂时回过头来，问一个问题：什么是 XQuery，我们又为什么需要使用 XQuery？正如 SQL 是为关系型数据模型设计的查询语言一样，XQuery 是专为查询 XML 数据设计的语言。由于 XML 数据可能与关系型数据有很大的差异，因此我们需要一种专门的语言来高效处理 XML 数据。关系型数据是平面的、高度结构化、强类型化、无序的，而 XML 数据是有序、嵌套的、层次化的、可选类型化的，并且往往不规则且是不完全结构化的。SQL 无法处理这样的情况，但 XQuery 专门为此设计。具体而言，XQuery 设计用于遍历 XML 文档树，并提取 XML 片段，而且还包括一些表达式，这些表达式用于创建、操纵和遍历 XML 项序列及构建新的 XML 数据。

IBM 已扩展了 Db2 的所有主要应用程序编程接口（API），以支持将 XQuery 作为 SQL 那样的一流语言。其中包括 CLI/ODBC、嵌入式 SQL、JDBC 和 .NET。因此，Db2 命令行处理器也支持 XQuery。您可按原样提交 XQuery，但需要使用 XQUERY 关键字启动它们，以通知 Db2 使用 XQuery 解析器，如下例所示：

##### 查询 13

```
XQUERY
for $dept in db2-fn:xmlcolumn("DEPT.DEPTDOC")/dept
where $dept/@deptID="PR27"
return $dept/employee/name;

```

Show moreShow more icon

查询 13 遍历各 department 文档中的 “dept” 元素，并返回在部门 PR27 工作的员工姓名。

### XQuery 的优点

- 对于仅涉及 XML、不需要（或不希望）使用 SQL 或关系结构的应用程序来说，XQuery 很适用。
- 对于从仅包含 XML 的数据库到 Db2 V9 的移植来说，XQuery 很适用。现有 XQuery 很可能仅需做出小小的修改即可在 Db2 中运行。例如，Db2 中的 XQuery 的数据输入来自 db2-fn:xmlcolumn() 函数，而其他数据库可能称之为 collection()。在这种情况下，仅需重新对函数进行命名即可。
- 如果查询结果需要嵌入在新构造的 XML 文档（此文档与数据库中其他 XML 文档不同）中（并以这个新文档的形式返回），那么 XQuery 有优势。
- XQuery 非常适合表示两个 XML 文档之间的连接以及联合 XML 值。例如，使用 XQuery，您可以将查询 11 中的连接以一种更为直接、高效的方式表示出来，如查询 14 所示。

##### 查询 14

```
XQUERY
for $dept in db2-fn:xmlcolumn("DEPT.DEPTDOC")/dept
     for $proj in db2-fn:xmlcolumn("PROJECT.PROJECTDOC")/project
where $dept/@deptID = $proj/deptID
return $dept/employee;

```

Show moreShow more icon

### XQuery 的缺点

- 使用普通 XQuery，您无法利用 Db2 Net Search Extender（NSE）所提供的全文本搜索功能。要进行全文本搜索，则需要包含 SQL。
- 普通 XQuery 不允许您调用 SQL 用户定义的函数（UDF）或者用 C 语言或 Java 编写的外部 UDF。
- 目前，Db2 不提供调用带参数标记的独立 XQuery 的方法。要向 XQuery 传递参数，则必须使用 SQL/XML 将一个参数标记（”?”）转换为指定变量。例如，在查询 13 中，您可能希望使用一个问号（?）作为 SQL 风格的参数标记，来取代字面值 “PR27”。但那是无效查询。

    在查询 7 中，可以看到 SQL/XML 允许您将一个 SQL 参数标记作为变量传递到嵌入式 XQuery 表达式中。SQL/XML 查询的 select 子句中通常有一个 XMLQUERY 函数，用于提取 XML 文档的某些部分；where 子句中通常会有一个 XMLEXISTS 谓词，用于过滤满足条件的文档。如果您希望在一条 FLWOR 表达式中表示整个查询逻辑，而不是使用两个单独的 XQuery 调用（XMLQUERY 中的一个调用和 XMLEXISTS 中的一个调用），并且希望依然使用参数标记，那么可以考虑将查询 13 重写为查询 15 所示形式：


##### 查询 15

```
values( XMLQUERY(
             ' for $dept in db2-fn:xmlcolumn("DEPT.DEPTDOC")/dept
                       where $dept/@deptID = $z
                       return $dept/employee/name'
                passing cast(? as char(8)) as "z" ) );

```

Show moreShow more icon

查询 15 是一个 SQL/XML 语句，因为它是 SQL fullselect 的 values 子句。values 子句通过指定结果表中的各列表达式来返回一个值表。在查询 15 中，结果表为 XML 类型的单行单列表，XMLQUERY 函数生成输出列的值。所有员工名都将在一个 XML 值中返回给客户机。查询 15 中的 FLOWR 表达式几乎与查询 13 中的表达式相同，惟一的不同在于查询 15 包含一个外部变量（$z），此变量被作为参数传递到 XMLQUERY 函数中。

## 嵌入了 SQL 的 XQuery

普通 XQuery 允许且只允许访问 XML 数据。如果只处理 XML 数据，那么这种方案非常好，但如果应用程序需要利用两种语言和数据模型的全部力量，综合访问 XML 数据和关系型数据，那么普通 XQuery 不足以胜任。使用本文前面讨论过的 SQL/XML —— 嵌入了 XQuery 的 SQL —— 可完成此任务。相反地，嵌入了 SQL 的 XQuery 为您带来了额外的可能性。除前文列出的优缺点之外，还有以下几个重要的方面：

### 嵌入了 SQL 的 XQuery 的优点

- 如果希望以关系列上的条件为依据，仅处理 XML 文档的子集，那么嵌入了 SQL 的 XQuery 非常适用。特别是您可以使用关系型谓词来约束特定 XQuery 的输入。为此，Db2 提供了另外一个输入函数 db2-fn:sqlquery，以便在 XQuery 中调用 SQL 查询。此函数接受 SQL SELECT 语句，并返回 XML 列作为输出。例如，查询 15 不考虑 XML 列的 deptdoc 中的所有文档，它有一个嵌入式 SQL 语句，通过连接到应用了谓词的 unit 表来预先过滤 XML 文档：

##### 查询 16

```
XQUERY
for $emp in db2-fn:sqlquery("select deptdoc
                     from dept d, unit u
                    where d.unitID=u.unitID and
                    u.manager = 'Jim Qu'")/dept/employee
where $emp/office = 216
return $emp/name;

```

Show moreShow more icon

两个表的 unitID 列上的常规关系索引以及 unit 表的 manager 列上的常规关系索引有助于提高嵌入式 SQL 查询的速度。Db2 9.5甚至可以同时使用 XML 和关系索引，比如嵌入式 SQL 语句的关系索引加上 XML 语句 $emp/office = 216 的 XML 索引。

- 嵌入了 SQL 的 XQuery 允许您利用全文本搜索，您可在嵌入式 SQL 语句的 where 子句中使用文本搜索函数 “contains”。在查询 17 中，XQuery 中的 SQL 从 dept 表中选择文档，其中文字 “sales” 出现在 /dept/description 中的某处。全文本索引迅速找到这些 XML 文档，随后将它们输入到从这些文档中提取所有员工姓名的 FLWOR 表达式中。就我们的示例表而言，查询 18 返回相同的结果，但是使用 SQL/XML 表示法表示的。

##### 查询 17

```
XQUERY
for $emp in db2-fn:sqlquery("
             select deptdoc from dept
             where CONTAINS(deptdoc, 'SECTION(""/dept/description"") ""sales"" ')=1
                                       ")//employee
return $emp/name;

```

Show moreShow more icon

##### 查询 18

```
select XMLQUERY('$d//employee/name' passing deptdoc as "d")
from dept
where CONTAINS(deptdoc,'SECTION("/dept/description") "sales" ')=1;

```

Show moreShow more icon

- 嵌入了 SQL 的 XQuery 对于那些需要整合关系型数据和 XML 数据的应用程序非常有用。您可以通过组合方式查询 XML 数据与关系型数据。这同样适用于 SQL/XML。但您可能会发现，在 SQL/XML 中利用 XMLEXISTS 函数连接 XML 值和关系值更轻松。请比较查询 19 和查询 20。查询 19 返回其员工中包含单位主管的那些部门的 deptID。嵌入式 SQL 语句将 unit 表中的主管姓名强制转换为 XML 类型（本例中为 XML 字符串），然后将其提供给 XQuery。XQuery 的 where 子句包含连接条件，这些条件用于主管姓名与员工姓名元素的比较。查询 20 用 SQL/XML 表示法表示了同样的连接。

##### 查询 19

```
XQUERY
for $m in db2-fn:sqlquery('select XMLCAST(u.manager as XML) from unit u')
for $d in db2-fn:xmlcolumn("DEPT.DEPTDOC")/dept
where $d/employee/name = $m
return $d/data(@deptID);

```

Show moreShow more icon

##### 查询 20

```
select XMLQUERY('$d/dept/data(@deptID)' passing  d.deptdoc as "d")
from dept d, unit u
where XMLEXISTS('$d//employee[name = $m]'
              passing d.deptdoc as "d", u.manager as "m");

```

Show moreShow more icon

- 从查询 19 中可以看出，嵌入了 SQL 的 XQuery 对于将关系型数据提供给 XQuery 来说非常有用。这允许您联合及合并 XML 数据和关系型数据。查询 21 构造了一个 result 文档，其中包含单位和部门信息。部门信息是从 XML 列的 deptdoc 中检索的 XML 文档。单位信息来自关系型表 unit。嵌入式 SQL 语句使用 SQL/XML 发布函数来构造 XML 元素 “Unit”，该元素带有三个子元素，子元素的值来自 unit 表的关系列，即 unitID、name 和 manager 这三列。

##### 查询 21

```
XQUERY
     let $d := db2-fn:sqlquery("select deptdoc from dept where unitID = 'WWPR' ")
     let $u := db2-fn:sqlquery("select XMLELEMENT(NAME ""Unit"",
                         XMLFOREST(unitID, name, manager))
                from unit where unitID = 'WWPR' " )
     return <result>
           <units>{$u}</units>
          <department>{$d}</department>
          </result>;

```

Show moreShow more icon

查询 21 的输出结果形式如下所示：

```
<result>
           <units><unit>
                     <UNITID> WWPR </UNITID>
                     <NAME> World Wide Markeing</NAME>
                     <MANAGER> Jim Qu </MANAGER>
                  </unit>
                  <unit>
                     <UNITID> WWPR </UNITID>
                     <NAME> ... </NAME>
                     <MANAGER> ... </MANAGER>
                 </unit>
              ....
         </units>
           <department>
               <dept deptID="PR27">
                    <employee id="901">
                    <name>Jim Qu</name>
                    <phone>408 555 1212</phone>
                  </employee>
                  <employee id="902">
                    <name>Peter Pan</name>
                    <office>216</office>
                  </employee>
               </dept>
           </department>
      </result>

```

Show moreShow more icon

- 嵌入了 SQL 的 XQuery 比较出色，因为嵌入式 SQL 语句可包含对用户定义的函数（UDF）的调用。UDF 在许多 IT 组织内都得到广泛应用，用于封装关键业务逻辑和简化应用程序开发需求。这非常重要，因为普通 XQuery 不能调用 UDF。

### 嵌入了 SQL 的 XQuery 的缺点

- 尽管 db2-fn:sqlquery 允许在 XQuery 内嵌入 SQL 语句，但目前还不能在嵌入式 SQL 语句中使用参数标记。
- 目前，db2-fn:sqlquery 不允许从 XQuery 向 SQL 中传递参数。因此在结合 XML 数据和关系型数据方面，SQL/XML 更强大。例如，嵌入了 SQL 的 XQuery 无法像 SQL/XML 那样用于表示关系列和 XML 列之间的连接。

## XML 查询结果

您需要注意的一件事就是：根据编写具体查询的方式不同，Db2 可能会交付不同格式的查询结果。例如，普通 XQuery 返回结果集形式的项（例如元素或文档片段），一项作为一行，即便多个项来自数据库中的同一文档（行）也是如此。另一方面，在使用 XMLQUERY 函数时，SQL/XML 可能会在一行中返回多个项，而不是返回多个独立的行。在某些情况下，您会认为这很令人满意，但在某些情况下又并非如此，具体取决于您的特定应用程序。让我们来看看示例。

查询 22 和查询 23 均请求获得 dept 文档中的员工姓名。查询 22 以 SQL/XML 方式表示，而查询 23 是用 XQuery 编写的。

##### 查询 22

```
select XMLQUERY('$d/dept/employee/name' passing deptdoc as "d")
from dept;

```

Show moreShow more icon

##### 查询 23

```
XQUERY db2-fn:xmlcolumn("DEPT.DEPTDOC")/dept/employee/name;

```

Show moreShow more icon

查询 22 为每个 dept 文档都返回一行，每个行都包含在该部门中工作的员工姓名：

Jim Qu Peter Pan Matt Foreman

而查询 23 将每一个员工姓名作为独立的一行返回：

Jim Qu Peter Pan Matt Foreman

查询 23 的结果通常更易于为应用程序所用，即：一次处理一个 XML 值。但在这个示例中，您并不知道哪些 name 元素来自同一个 department 文档。查询 22 的输出结果保留了这一信息，但应用程序利用其结果较为困难，因为应用程序可能需要分离这些结果行，以便分别访问每个 name 元素。如果应用程序使用 XML 解析器从 Db2 中获取每个 XML 结果行，那么解析器将拒绝查询 22 的第一个结果行，因为它并非格式良好的文档（缺少单独的根元素）。为解决此问题，您可向查询 22 中添加一个 XMLELEMENT 构造函数，从而添加一个根元素，如查询 24 所示：

##### 查询 24

```
Select XMLELEMENT(name "employees",
    XMLQUERY('$d/dept/employee/name'   passing d.deptdoc as "d"))
from dept d;

```

Show moreShow more icon

查询结果发生了改变，每个结果行都是一个格式良好的 XML 文档：

Jim Qu Peter Pan Matt Foreman….

回忆一下，查询 15 使用了一个 SQL values 子句和 XMLQUERY 函数来允许向 XQuery 传递参数。但查询 15 的输出结果为单行，其中包含所有员工的姓名。如果您更希望使各员工的姓名自成一行，同时还需要使用参数标记，那么可以借助 XMLTABLE 函数，如查询 25 所示：

##### 查询 25

```
select X.*
from dept d, XMLTABLE('for $dept in $d/dept
                       where $dept/@deptID = $z
                       return $dept/employee/name'
                       passing d.deptdoc as "d", cast(? as char(10)) as "z"
            COLUMNS
                       "name"   XML  PATH ".") as X ;

```

Show moreShow more icon

## 结束语

Db2 V9 为查询 XML 数据提供了一组丰富的选项。您将根据应用程序的需求和特征选择最适合您的方案。如果需要结合关系型数据和 XML 数据，那么 SQL/XML 在绝大多数情况下都是最佳选择。具体而言，SQL/XML 方案允许您为 XML 数据使用参数标记。如果仅涉及 XML 的应用程序，那么独立 XQuery 是一种功能强大的方案，嵌入式 SQL 可进一步增强它，支持您使用全文本搜索和 UDF 调用。本文所讨论的示例帮助您作出了明智的决策。您可使用我们的查询模式作为起点，编写自己的查询，查询您自己的 XML 数据。

一个通用的指导原则就是根据实际需要为查询引入恰如其分的复杂性。例如，使用一个带有内置 XQuery 的嵌入式 SQL 语句的 XQuery 是完全可能的。但是，我们的经验显示，表达所需的查询逻辑通常并不需要两种语言的嵌套层次超过一层。因此，我们建议仅在 SQL 中使用一层嵌入式 XQuery，反之亦然。

## 致谢

感谢 Don Chamberlin、Bert van der Linden、Cindy Saracco、Jan-Eike Michels 和 Hardeep Singh 对本文的审阅。

本文翻译自： [pureXML in Db2 9: Which way to query your XML data?](https://developer.ibm.com/articles/dm-0606nicola/)（2006-06-08）