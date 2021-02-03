# 关于 Java Database Connectivity 您不知道的 5 件事
提升您和 JDBC API 的关系

**标签:** Java,数据库

[原文链接](https://developer.ibm.com/zh/articles/j-5things10/)

Ted Neward, Alex Theedom

发布: 2010-09-21

* * *

##### 关于本系列

您觉得自己懂 Java 编程？事实是，大多数开发人员都只领会到了 Java 平台的皮毛，所学也只够应付工作。在 [本系列](/zh/series/5-things-you-didnt-know-about/) 中，Ted Neward 深度挖掘 Java 平台的核心功能，揭示一些鲜为人知的事实，帮助您解决最棘手的编程困难。

目前，许多开发人员把 Java Database Connectivity (JDBC) API 当作一种数据访问平台，比如 Hibernate 或 SpringMany。然而 JDBC 在数据库连接中不仅仅充当后台角色。对于 JDBC，您了解的越多，您的 RDBMS 交互效率就越高。

在本期 [_5 件事_ 系列](/zh/series/5-things-you-didnt-know-about/) 中，我将向您介绍几种 JDBC 2.0 到 JDBC 4.0 中新引入的功能。设计时考虑到现代软件开发所面临的挑战，这些新特性支持应用程序可伸缩性，并提高开发人员的工作效率 — 这是现代 Java 开发人员面临的两个最常见的挑战。

## 1\. 标量函数

不同的 RDBMS 实现对 SQL 和/或增值特性（目的是让程序员的工作更为简单）提供不规则的支持。例如，众所周知，SQL 支持一个标量运算 `COUNT()`，返回满足特定 SQL 过滤规则的行数（更确切地说，是 `WHERE` 谓词）。除此之外，修改 SQL 返回的值是很棘手的 — 想要从数据库获取当前日期和时间会使 JDBC 开发人员、甚至最有耐心的程序员发疯（甚至是心力憔悴）。

于是，JDBC 规范针对不同的 RDBMS 实现通过标量函数提供一定程度的隔离/改写。JDBC 规范包括一系列受支持的操作，JDBC 驱动程序应该根据特定数据库实现的需要进行识别和改写。因此，对于一个支持返回当前日期和/或时间的数据库，时间查询应当如清单 1 那样简单：

##### 清单 1\. 当前时间？

```
Connection conn = ...; // get it from someplace
Statement stmt = conn.createStatement();
ResultSet rs = stmt.executeQuery("{CURRENT_DATE()}");

```

Show moreShow more icon

JDBC API 识别的标量函数完整列表在 JDBC 规范附录中给出（见参考资源），但是给定的驱动程序或数据库可能不支持完整列表。您可以使用从 `Connection` 返回的 `DatabaseMetaData` 对象来获取给定 JDBC 支持的函数，如清单 2 所示：

##### 清单 2\. 能为我提供什么？

```
Connection conn = ...; // get it from someplace
DatabaseMetaData dbmd = conn.getMetaData();

```

Show moreShow more icon

标量函数列表是从各种 `DatabaseMetaData` 方法返回的一个逗号分隔的 `String` 。例如，所有数值标量由 `getNumericFunctions()` 调用列出，在结果上执行一个 `String.split()` — _瞧！_ — 即刻出现 `equals()-testable` 列表。

## 2\. 可滚动 ResultSets

创建一个 `Connection` 对象，并用它来创建一个 `Statement`，这在 JDBC 中是最常用的。提供给 `SQL SELECT` 的 `Statement` 返回一个 `ResultSet`。然后，通过一个 `while` 循环（和 `Iterator` 没什么不同）得到 `ResultSet`，直到 `ResultSet` 为空，循环体从左到右的每次提取一列 。

这整个操作过程是如此普遍，近乎神圣：它这样做只是因为它应该这样做。唉！实际上这是完全没必要的。

### 引入可滚动 ResultSet

许多开发人员没有意识到，在过去的几年中 JBDC 已经有了相当大的增强，尽管这些增强在新版本中已经有所反映。 第一次重大增强是在 JDBC 2.0 中，发生在使用 JDK 1.2 期间。写这篇文章时，JDBC 已经发展到了 JDBC 4.0。

JDBC 2.0 中一个有趣的增强（尽管常常被忽略）是 `ResultSet` 的滚动功能，这意味着您可以根据需要前进或者后退，或者两者均可。这样做需要一点前瞻性，然而 — JDBC 调用必须指出在创建 `Statement` 时需要一个可以滚动的 `ResultSet` 。

##### 验证 ResultSet 类型

如果您怀疑一个驱动程序事实上可能不支持可滚动的 `ResultSet`s，不管 `DatabaseMetaData` 中是如何写的，您都要调用 `getType()` 来验证 `ResultSet` 类型。当然，如果您是个偏执的人，您可能也不相信 `getType()` 的返回值。可以这样说，如果 `getType()` 隐瞒关于 `ResultSet` 的返回值，它们确实 _是_ 要吃定您。

如果底层 JDBC 驱动程序支持滚动，一个可滚动的 `ResultSet` 将从那个 `Statement` 返回。但是在请求它之前最好弄清楚驱动程序是否支持可滚动性。您可以通过 `DatabaseMetaData` 对象探询滚动性，如上所述，这个对象可从任何 `Connection` 中获取。

一旦您有了一个 `DatabaseMetaData` 对象，一个对 `getJDBCMajorVersion()` 的调用将会确定驱动程序是否支持 JDBC 规范，至少是 JDBC 2.0 规范。当然一个驱动程序可能会隐瞒它对给定规范的支持程度，因此为了安全起见，用期望得到的 `ResultSet` 类型调用 `supportsResultSetType()` 方法。（在 `ResultSet` 类上它是一个常量；稍后我们将对其每个值进行讨论。）

##### 清单 3\. 可以滚动？

```
int JDBCVersion = dbmd.getJDBCMajorVersion();
boolean srs = dbmd.supportsResultSetType(ResultSet.TYPE_SCROLL_INSENSITIVE);
if (JDBCVersion > 2 || srs == true)
{
    // scroll, baby, scroll!
}

```

Show moreShow more icon

### 请求一个可滚动的 ResultSet

假设您的驱动程序回答 “是”（如果不是，您需要一个新的驱动程序或数据库），您可以通过传递两个参数到 `Connection.createStatement()` 调用来请求一个可滚动的 `ResultSet` ，如清单 4 所示：

##### 清单 4\. 我想要滚动！

```
Statement stmt = con.createStatement(
                       ResultSet.TYPE_SCROLL_INSENSITIVE,
                       ResultSet.CONCUR_READ_ONLY);
ResultSet scrollingRS = stmt.executeQuery("SELECT * FROM whatever");

```

Show moreShow more icon

在调用 `createStatement()` 时，您必须特别小心，因为它的第一个和第二个参数都是 `int` 的。（在 Java 5 之前我们不能使用枚举类型！）任何 `int` 值（包括错误的常量）对 `createStatement()` 都有效。

第一个参数，指定 `ResultSet` 中期望得到的 “可滚动性”，应该是以下 3 个值之一：

- `ResultSet.TYPE_FORWARD_ONLY`：这是默认的，是我们了解且喜欢的流水游标。
- `ResultSet.TYPE_SCROLL_INSENSITIVE`：这个 `ResultSet` 支持向后迭代以及向前迭代，但是，如果数据库中的数据发生变化，`ResultSet` 将不能反映出来。这个可滚动的 `ResultSet` 可能是最常用到的类型。
- `ResultSet.TYPE_SCROLL_SENSITIVE`：所创建的 `ResultSet` 不但支持双向迭代，而且当数据库中的数据发生变化时还为您提供一个 “实时” 数据视图。

第二个参数在下一个技巧中介绍，稍等片刻。

### 定向滚动

当您从 `Statement` 获取一个 `ResultSet` 后，通过它向后滚动只需调用 `previous()`，即向后滚动一行，而不是向前，就像 `next()` 那样。您也可以调用 `first()` 返回到 `ResultSet` 开头，或者调用 `last()` 转到 `ResultSet` 的末尾，或者…您自己拿主意。

`relative()` 和 `absolute()` 方法也是很有用的：前者移动指定数量的行（如果是正数则向前移动，是负数则向后移动），后者移动 `ResultSet` 中指定数量的行，不管游标在哪。当然，目前行数是由 `getRow()` 获取的。

如果您打算通过调用 `setFetchDirection()` 在一个特定方向进行一些滚动，可以通过指定方向来帮助 `ResultSet`。（无论向哪个方向滚动，对于 `ResultSet` 都可行，但是预先知道滚动方向可以优化其数据检索。）

## 3\. 可更新的 ResultSets

JDBC 不仅仅支持双向 `ResultSet`，也支持就地更新 `ResultSet`。这就是说，与其创建一个新 SQL 语句来修改目前存储在数据库中的值，您只需要修改保存在 `ResultSet` 中的值，之后该值会被自动发送到数据库中该行所对应的列。

请求一个可更新的 `ResultSet` 类似于请求一个可滚动的 `ResultSet` 的过程。事实上，在此您将为 `createStatement()` 使用第二个参数。您不需要为第二个参数指定 `ResultSet.CONCUR_READ_ONLY` ，只需要发送 `ResultSet.CONCUR_UPDATEABLE` 即可，如 清单 5 所示：

##### 清单 5\. 我想要一个可更新的 ResultSet

```
Statement stmt = con.createStatement(
                       ResultSet.TYPE_SCROLL_INSENSITIVE,
                       ResultSet.CONCUR_UPDATEABLE);
ResultSet scrollingRS = stmt.executeQuery("SELECT * FROM whatever");

```

Show moreShow more icon

假设您的驱动程序支持可更新光标（这是 JDBC 2.0 规范的另一个特性，这是大多数 “现实” 数据库所支持的），您可以更新 `ResultSet` 中任何给定的值，方法是导航到该行并调用它的一个 `update...()` 方法（如清单 6 所示），如同 `ResultSet` 的 `get...()` 方法。在 `ResultSet` 中 `update...()` 对于实际的列类型是超负荷的。因此要更改名为 “`PRICE` ” 的浮点列，调用 `updateFloat("PRICE")` 。然而，这样做只能更新 `ResultSet` 中的值。为了将该值插入支持它的数据库中，可以调用 `updateRow()` 。如果用户改变调整价格的想法，调用 `cancelRowUpdates()` 可以停止所有正在进行的更新。

##### 清单 6\. 一个更好的方法

```
Statement stmt = con.createStatement(
                       ResultSet.TYPE_SCROLL_INSENSITIVE,
                       ResultSet.CONCUR_UPDATEABLE);
ResultSet scrollingRS =
    stmt.executeQuery("SELECT * FROM lineitem WHERE id=1");
scrollingRS.first();
scrollingRS.udpateFloat("PRICE", 121.45f);
// ...
if (userSaidOK)
    scrollingRS.updateRow();
else
    scrollingRS.cancelRowUpdates();

```

Show moreShow more icon

JDBC 2.0 不只支持更新。如果用户想要添加一个全新的行，不需要创建一个新 `Statement` 并执行一个 `INSERT` ，只需要调用 `moveToInsertRow()` ，为每个列调用 `update...()` ，然后调用 `insertRow()` 完成工作。如果没有指定一个列值，数据库会默认将其看作 `SQL NULL` （如果数据库模式不允许该列为 `NULL` ，这可能触发 `SQLException` ）。

当然，如果 `ResultSet` 支持更新一行，也必然支持通过 `deleteRow()` 删除一行。

差点忘了强调一点，所有这些可滚动性和可更新性都适用于 `PreparedStatement` （通过向 `prepareStatement()` 方法传递参数），由于一直处于 SQL 注入攻击的危险中，这比一个规则的 `Statement` 好很多。

## 4\. Rowsets

既然所有这些功能在 JDBC 中大约有 10 年了，为什么大多数开发人员仍然迷恋向前滚动 `ResultSet` 和不连贯访问？

罪魁祸首是可伸缩性。保持最低的数据库连接是支持大量用户通过 Internet 访问公司网站的关键。因为滚动和/或更新 `ResultSet` 通常需要一个开放的网络连接，而许多开发人员通常不（或不能）使用这些连接。

幸好，JDBC 3.0 引入另一种解决方案让您同样可以做很多之前使用 `ResultSet` 方可以做的事情，而不需要数据库连接保持开放状态。

从概念上讲，`Rowset` 本质上是一个 `ResultSet`，但是它支持连接模型或断开模型，您所需要做的是创建一个 `Rowset`，将其指向一个 `ResultSet`，当它完成自我填充之后，将其作为一个 `ResultSet`，如清单 7 所示：

##### 清单 7\. Rowset 取代 ResultSet

```
Statement stmt = con.createStatement(
                       ResultSet.TYPE_SCROLL_INSENSITIVE,
                       ResultSet.CONCUR_UPDATEABLE);
ResultSet scrollingRS = stmt.executeQuery("SELECT * FROM whatever");
if (wantsConnected)
    JdbcRowSet rs = new JdbcRowSet(scrollingRS); // connected
else
    CachedRowSet crs = new CachedRowSet(scrollingRS); disconnected

```

Show moreShow more icon

JDBC 还附带了 5 个 `Rowset` 接口 “实现”（也就是扩展接口）。 `JdbcRowSet` 是一个连接的 `Rowset` 实现；其余 4 个是断开的：

- `CachedRowSet` 只是一个断开的 `Rowset`.
- `WebRowSet` 是 `CachedRowSet` 的一个子集，知道如何将其结果转换成 XML，并再次转换回来。
- `JoinRowSet` 是一个 `WebRowSet` ，知道如何形成一个 `SQL JOIN` ，而无需连接到数据库。
- `FilteredRowSet` 是一个 `WebRowSet` ，知道如何更进一步过滤传递回来的数据，而不需要连接到数据库。

`Rowsets` 是完整的 JavaBeans，意味着它们支持侦听类事件，因此，如果需要，也可以捕捉、检查并执行对 `Rowset` 的任何修改。事实上，如果 `Rowset` 有自己的 `Username`、 `Password`、 `URL` 和 `DatasourceName` 属性集（这意味着它将使用 `DriverManager.getConnection()` 创建一个连接）或者 `Datasource` 属性集（这很可能由 JNDI 获取），它甚至能管理对数据库的全部操作。然后，您可以在 `Command` 属性中指定要执行的 SQL，调用 `execute()`，然后处理结果 — 不需要更多的工作。

通常，`Rowset` 实现是由 JDBC 驱动程序提供的，因此实际的名称和/或包由您所使用的 JDBC 驱动程序决定。从 Java 5 开始 `Rowset` 实现已经是标准版本（standard distribution）的一部分了，因此您只需要创建一个 `...RowsetImpl()`，然后让其运行。（如果您的驱动程序不能提供一个参考实现，Sun 提供了一个，参见 参考资料 部分的链接。）

## 5\. 批量更新

尽管 `Rowset` 很实用，但有时候也不能满足您的需求，您可能需要返回来直接编写 SQL 语句。在这种情况下，特别是当您面对一大堆工作时，您就会很感激批量更新功能，可在一个网络往返行程中在数据库中执行多条 SQL 语句。

要确定 JDBC 驱动程序是否支持批量更新，快速调用 `DatabaseMetaData.supportsBatchUpdates()` 可产生一个明示支持与否的布尔值。在支持批量更新时（由一些非 `SELECT` 标示），所有任务逐个排队然后在某一瞬间同时得到更新，如清单 8 所示：

##### 清单 8\. 让数据库进行批量更新！

```
conn.setAutoCommit(false);

PreparedStatement pstmt = conn.prepareStatement("INSERT INTO lineitems VALUES(?,?,?,?)");
pstmt.setInt(1, 1);
pstmt.setString(2, "52919-49278");
pstmt.setFloat(3, 49.99);
pstmt.setBoolean(4, true);
pstmt.addBatch();

// rinse, lather, repeat

int[] updateCount = pstmt.executeBatch();
conn.commit();
conn.setAutoCommit(true);

```

Show moreShow more icon

默认必须调用 `setAutoCommit()`，驱动程序会试图交付提供给它的每条语句。除此之外，其余代码都是简单易懂的：使用 `Statement` 或 `PreparedStatement` 进行常见 SQL 操作，但是不调用 `execute()`，而调用 `executeBatch()`，排队等候调用而不是立即发送。

准备好各种语句之后，在数据库中使用 `executeBatch()` 触发所有的语句，这将返回一组整型值，每个值保存同样的结果，好像使用了 `executeUpdate()` 一样。

在批量处理的一条语句发生错误的情况下，如果驱动程序不支持批量更新，或者批处理中的一条语句返回 `ResultSet`，驱动程序将抛出一个 `BatchUpdateException`。有时候，在抛出一个异常之后，驱动程序可能试着继续执行语句。JDBC 规范不能授权某一行为，因此您应该事先试用驱动程序，这样就可以确切地知道它是如何工作的。（当然，您要执行单元测试，确保在错误成为问题之前发现它，对吧？）

## 结束语

作为 Java 开发的一个主题，JDBC API 是每个开发人员应该熟知的，就像您的左右手那样。有趣的是，在过去的几年中，许多开发人员并不了解 API 的增强功能，因此，他们错失了本文所讲到的省时技巧。

当然，您是否决定使用 JDBC 的新功能取决于您自己。需要考虑的一个关键因素是您所使用的系统的可伸缩性。对伸缩性的要求越高，对数据库的使用就越受限制，因此而减少的网络流量就会越多。 `Rowset` 、标量调用和批量更新将会是给您带来帮助的益友。另外，尝试可滚动和可更新的 `ResultSet`（这不像 `Rowset` 那样耗内存），并度量可伸缩性。这可能没您想象的糟糕。

[_5 件事_ 系列](/zh/series/5-things-you-didnt-know-about) 的下一期主题是: 命令行标志。

本文翻译自： [Java Database Connectivity](https://developer.ibm.com/articles/j-5things10/)（2010-08-10）