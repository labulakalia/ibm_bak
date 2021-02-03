# 开始通过 PHP 使用 ODBC
使用 PHP 实现通用数据库连接性的初学者指南

**标签:** PHP,Web 开发,数据库

[原文链接](https://developer.ibm.com/zh/articles/os-php-odbc/)

Daniel Lewis

发布: 2011-09-05

* * *

PHP 是动态网站开发最常使用的编程语言之一。PHP 相当强大和有效，并且还十分简单，对于初学者，因为该语言的灵活性，它将是一个不错的语言学习选择。

单独从语言角度而言，PHP 是一门不错的语言（尤其在与 XHTML 结合使用时）。但是，大多数应用程序需要一种简便的数据存储方法，通常由数据库（比如 MySQL 或 PostgreSQL）实现。为了链接到数据存储系统，连接器需要允许使用 `mysql_query()` 和 `pg_query()` 之类的函数。大多数情况下，这种方法是有效的，但是偶尔地，因为业务需求，可能会要求将数据存储在更容易被人操作的电子表格之类的地方。在这些情况下，需要使用不同类型的连接器来连接数据。

如果您想将 SQL 发送到非特定数据库系统，并想使用 SQL 处理您配置的任何数据存储，该怎么做？因为开放数据库连接（Open Database Connectivity, ODBC）已经创建，并且已经安装了正确的连接器，所以您甚至可以访问 Microsoft® Excel® 、CSV 和其他数据文件类型。ODBC 是一个连接器，它使得 PHP 开发变得 “与数据库连接器无关”。它对数据库（比如 MySQL、PostgreSQL、SQLite、Microsoft SQL Server® 、IBM® DB2® 、Sybase、OpenLink Virtuoso、FileMaker 和 Microsoft Office® Access® ）使用像 `odbc_query()` 这样的函数。还可以将 ODBC 用于 CSV 和 Excel 电子表格，具体取决于正确的 ODBC 驱动程序设置。

## ODBC 是什么？

ODBC 是一个软件驱动程序系统，用于连接编程语言与数据存储。ODBC 是一个免费的开放源码系统，出现于 1992 年，它试图通过编程语言和数据库查询访问（SQL 标准化）来标准化连接方法，比如功能和配置。

ODBC 的作用是充当接口或连接器，它具有双重设计目标：首先，对于 ODBC 系统，它充当的是编程语言系统，其次，对于数据存储系统，它充当的是 ODBC 系统。所以，ODBC 需要一个 “对 ODBC 而言是编程语言” 的驱动程序（例如 PHP-ODBC 库）和一个 “对数据存储系统而言是 ODBC” 的驱动程序（比如 MySQL-ODBC 库）。除了 ODBC 系统本身之外，ODBC 还可以处理数据源的配置，允许数据源和编程语言之间存在模糊性。

## 如何使 PHP 与 SQL 相互适应？

PHP 是一种编程语言，通常用作服务器端语言，用来加速动态网站的发展。因为是一种动态的、弱类型语言，PHP 非常灵活。许多开发人员都非常熟悉 PHP，因为它受到 C 编程语言的影响。PHP 是一个免费的开放源码编程语言，出现于 1995 年，您可以通过连接器对数据库使用 PHP，以生成 XHTML 和 HTML，然后在 Web 浏览器中呈现内容。

SQL 是用于询问数据存储的一种跨平台语言，主要用于关系数据库，但它有一些过程性的、面向对象的对象关系扩展。现代的 SQL 实现可以在 MySQL、PostgreSQL、SQLite、DB2（商业和 Express-C 版本）、Microsoft SQL Server、OpenLink Virtuoso、FileMaker 和 Microsoft Access 中找到，所有这些都可以通过连接系统（ODBC）使用编程语言（比如 PHP）来实现。

## 设置 ODBC

让我们来看一下如何将典型的 Linux-Apache-PHP-MySQL (LAMP) 环境转换为灵活的 Linux-Apache-PHP-ODBC (LAPO) 环境。在 Linux® 上，有两个针对 ODBC 驱动程序的常规选项：iODBC 和 unixODBC。这两套驱动程序各有自己的优缺点，但它们都能处理不同的数据库集。我选择使用 iODBC，因为它在连接到 Web 编程语言（比如 PHP 和 Ruby）时通常表现得很强大，此外，它在处理 ODBC 友好的数据库（比如 MySQL 和 OpenLink Virtuoso）时很稳定。但这只是一种选择，您可能想调查自己的效率要求。除了内部的细微差别，从连接的角度来看，iODBC 和 unixODBC 的使用方式与编程语言（PHP 功能是相同的）和数据库（例如，MySQL 将不受影响）是相同的。

iODBC 和 unixODBC 都可以在 Linux Software Package Managers 中使用。例如，在 Debian、Ubuntu 或 Linux Mint 命令行中，都可以运行 `sudo apt-get install iodbc` 。

##### 其他环境中的 ODBC

Apple Mac OS X 会随 iODBC 一起预安装，因此没必要再安装任何东西。Windows 操作系统也有自己的 ODBC 驱动程序。

### 使用 ODBC 连接数据库

如果您还没有安装数据库系统，比如 MySQL 或 PostgreSQL，请安装一个。然后安装 “ODBC 到数据库” 的连接器。此连接器因数据库的不同而各异，但是，以 MySQL 为例，可以通过安装从 MySQL 网站下载的与操作系统相关的驱动程序来安装连接器。

Linux 发行版支持 `apt` ，您可以从控制台运行以下命令： `sudo apt-get install libmyodbc` 。

您必须配置您的 ODBC 客户端，以便通过运行程序（比如 iodbcadm-gtk）来安装数据库。您还可以手动编辑 iODBC 文件。（在 Linux 中，该文件通常位于 /etc/iodbc.ini 中。）

### 示例：使用 ODBC 连接 PHP

接下来，您必须安装 PHP ODBC 驱动程序。通过将 iODBC 或 unixODBC 添加到 PHP 编译脚本中（非常复杂），或者通过安装 PHP-ODBC 库，都可以实现 PHP ODBC 驱动程序的安装。在基于 apt 的发行版中，可以使用以下命令： `sudo apt-get install php5-odbc` 。

可通过在交互模式下运行 PHP (`php -a`) 来测试流量。这将打开 PHP 交互控制台，在那里您可以使用与清单 1 中的示例类似的方式进行交互。

##### 清单 1\. 命令行 ODBC 连接

```
php > $conn = odbc_connect(
"DRIVER={MySQL ODBC 3.51 Driver};Server=localhost;Database=phpodbcdb",
"username", "password");
php > $sql = "SELECT 1 as test";
php > $rs = odbc_exec($conn,$sql);
php > odbc_fetch_row($rs);
php > echo "\nTest\n—--\n” . odbc_result($rs,"test") . "\n";

Test
----
1
php > odbc_close($conn);
php > exit;

```

Show moreShow more icon

让我们来分析一下清单 1 中的代码：

1. 使用 PHP 中的 `odbc_connect()` 函数建立一个连接。该函数将 ODBC 连接器字符串、用户名称和密码作为参数。连接器字符串应与 odbc.ini 文件匹配，以确保它与预先安排的相符。
2. 将一个变量实例化为字符串，该字符串代表了您的 SQL 语句。
3. 使用 `odbc_exec` 函数执行该 SQL 语句，此函数将接受您的连接和您的 SQL 字符串，并返回一个结果集。
4. 使用 `odbc_fetch_row()` 仅从结果集中提取一行结果，该函数接受作为参数的结果集。这是一个迭代函数，这意味着如果您再次调用它，会得到结果集中的下一个结果（依此类推，直到结果集中没有结果，如果没有结果，则返回 false）。
5. 使用函数 `odbc_result()` ，该函数接受结果集和列名称（字符串形式），并返回行迭代程序所指向的行中的单元值。
6. 使用 `odbc_close` 函数关闭 ODBC 连接，该函数接受连接本身。
7. 通过发送 `exit` 命令退出 PHP 交互模式。

这方法很有用，但在 Web 应用程序级别上不常使用。如果您想在客户端/服务器样式的 Web 浏览模式下测试流量，则需要安装一台 Web 服务器，比如 Apache 或 Lighttpd。（如果运行的是 Linux 系统，则必须确保提供了用于 Web 服务的 PHP 模块，否则，PHP 将无法运行。）

清单 2 展示了通过 Web 服务器执行此操作时使用的相同技术。PHP 代码类似于 [清单 1](#清单-1-命令行-odbc-连接) 中使用的代码，但它通过 XHTML 而不是命令行导出结果。

##### 清单 2\. 基于 XHTML 的 ODBC 连接的示例

```
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN"
    "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
    <head>
        <title>PHP and ODBC: XHTML Example 1</title>
        <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    </head>
    <body>
    <?php
        $conn = odbc_connect(
          "DRIVER={MySQL ODBC 3.51 Driver};Server=localhost;Database=phpodbcdb",
          "username", "password");
        if (!($conn)) {
          echo "<p>Connection to DB via ODBC failed: ";
          echo odbc_errormsg ($conn );
          echo "</p>\n";
        }

        $sql = "SELECT 1 as test";
        $rs = odbc_exec($conn,$sql);
        echo "<table><tr>";
        echo "<th>Test</th></tr>";
        while (odbc_fetch_row($rs))
         {
         $result = odbc_result($rs,"test");
         echo "<tr><td>$result</td></tr>";
        }
        odbc_close($conn);
        echo "</table>";
    ?>
    </body>
</html>

```

Show moreShow more icon

该清单在 [清单 1](#清单-1-命令行-odbc-连接) 的基础上添加了一些内容：现在可以完全将 `odbc_fetch_row()` 函数用作迭代函数，只需将它放在 `while` 循环中即可。这意味着，如果 SQL 稍微有点复杂，并且要查询一个表来查找多个结果，那么该函数会在所呈现的 XHTML 表中创建一个新行。

有许多类型的 XHTML 和 HTML，它们配备有各种数量的浏览器支持，且易于使用。 [清单 2](#清单-2-基于-xhtml-的-odbc-连接的示例) 生成了标准化的 XHTML 1.0 Strict，它是 XHTML 用来开发以数据为中心的、强大的、跨浏览器的文档的最佳形式之一。

## PHP-ODBC 编程

ODBC 函数有 4 种主要类型：用于连接、查询、提取数据和错误报告的函数。查询函数能够处理标准化数据库事务，以便创建、读取、更新和删除数据（通称 _CRUD_ 操作）。

### 连接函数

每个已开始的进程都必须有一个完结；否则，就会导致内存和处理器问题。所以，您要确保已经关闭了数据库连接。

如您已经看到的， `odbc_connect()` 函数接受 ODBC 友好的链接字符串、数据库用户名称和相关密码。它返回一个您可以在整个 PHP 程序中使用的连接对象。以下代码显示了一个示例：

```
$connection = odbc_connect($connection_string, $username, $password);

```

Show moreShow more icon

在前面的示例中还可以看见， `odbc_close()` 函数接受了一个连接对象，并终止了与 ODBC 和数据库的通信。我要强调的是，您必须关闭您的链接，否则会有过多的连接到您的数据库的连接，此外，您还必须重启您的数据库管理系统，在更糟糕的情况下，您甚至需要重启机器。以下是该函数的运行方式： `odbc_close($connection);` 。

### 查询函数

前面曾使用过 `odbc_exec()` 函数，它接受了一个连接对象和一个 SQL 字符串，执行该函数后，会返回一个结果集对象。该结果集对象是一个复杂对象，它通常位于数据库管理系统的存储器中，并且只能通过与之交互的函数辨认。 `odbc_exec()` 行与以下代码类似： `$resultset = odbc_exec($connection, $sql);` 。

在将未知变量注入 SQL 中时， `odbc_prepare()` 和 `odbc_execute` 函数非常有用。 `odbc_prepare()` 函数为数据库管理系统准备了一个 SQL 语句，然后 `odbc_execute()` 函数会在变量中发送该语句。这意味着它比使用 PHP 建立一串 SQL 字符串并通过 `odbc_exec()` 发送 SQL 语句更强大、更安全、更有效。将这些函数放在一起时，它们看起来如下所示：

```
$resultset = odbc_prepare($connection, $sql);
$success = odbc_execute($resultset, $variables);

```

Show moreShow more icon

清单 3 是一个很好的示例，创建它是为了根据位置和出生日期变量来搜索用户表中的用户。请注意 SQL 字符串中的问号 (`?`)，它表示 `odbc_execute()` 函数的串行数组中定义的变量。

##### 清单 3\. 使用 prepare 和 execute 命令实现 SQL 变量注入

```
$location = "London";
$mindateofbirth = time() - 567648000; /* i.e. 18 years ago from now */
$resultset = odbc_prepare(
        $connection,
        "SELECT * FROM user WHERE location = ? AND dateofbirth <= ?"
      );
$success = odbc_execute($resultset, array($location, $mindateofbirth));

```

Show moreShow more icon

### 提取函数

`odbc_fetch_row()` 函数接受来自某个查询的结果集，并将迭代器指针从一行转向下一行。此函数常常与 `odbc_result()` 函数结合使用，以提取各种单元格的值：

```
odbc_fetch_row($resultset);

```

Show moreShow more icon

在前面的示例中， `odbc_result()` 函数接受了一个 `$resultset` 和一个列名称字符串，并返回某个单元格的值。此函数可以与 `odbc_fetch_row()` 函数结合使用，以指向结果集中的特定行：

```
$value = odbc_result($resultset,"columnname");

```

Show moreShow more icon

`odbc_fetch_array()` 函数在某些地方类似于用来从查询结果集中提取数据的迭代函数。但是，在这里，它返回了一个代表行的数组， 并使用列名称作为键，使用单元格值作为值：

```
$rowarray = odbc_fetch_array($resultset);

```

Show moreShow more icon

与 `odbc_fetch_array()` 函数类似， `odbc_fetch_object()` 提取代表行的面向对象的结构。它将列名称作为对象属性，并将单元格值作为属性值：

```
$rowobject = odbc_fetch_object($resultset);

```

Show moreShow more icon

此函数在打印 HTML 形式的结果集时非常有用。它只是简单地呈现结果，但在原型制作或调试时很有用：

```
odbc_result_all($resultset);

```

Show moreShow more icon

`odbc_num_fields()` 函数是一个相当不错的函数，它只接受结果集，并会告诉您结果集中的行数：

```
$numberofrows = odbc_num_rows($resultset);

```

Show moreShow more icon

### 问题解决和调试

PHP ODBC 有两个得心应手的函数，其中一个是 `odbc_error()` ，如果发生错误，它会返回错误代码，如果没有发生错误，则返回 false；另一个函数是 `odbc_errormsg()` ，它返回用户友好的消息。您可以组合使用这两个函数，从而形成一个简单的错误消息序列：

```
if (odbc_error()) {
    echo "I've found a problem: " . odbc_errormsg($conn);
}

```

Show moreShow more icon

如果在开发的时候出错，您会获得另一个提示，不要害怕向导致问题的行附近添加打印语句，当然，在显示其他行时，系统会为您提供删除这些 “调试行” 的权利。请注意下面的 PHP 函数，它通常会在关键时候为您伸出援手:

```
print_r($variable);

```

Show moreShow more icon

这个简单的函数接受任何变量，并将它显示在屏幕上。变量可以像一个整数或字符串那样简单，也可以像多维数组或对象那样复杂。

## 通用连接性

可以考虑构建一个可在任何地方（这方面的例子可以包括 Drupal、WordPress 或 Joomla）部署的 Web 应用程序。这些应用程序通常使用数据库（比如 MySQL）以及特定于数据库的函数（比如 `mysql_connect()` ）来构建，然后，您可以小心地重构它们，以便通过更改函数（比如 `pg_connect()` ）来处理其他数据库。在使用 ODBC 时，这是一个冗余体验，因为在初始化应用程序时已经完成配置，所以对于数据库系统而言，ODBC 函数具有双重作用。

但是，有一件事需要注意，尽管所有的数据库管理系统都共享一个标准化的 SQL，但它们有时也会提供它们自己的扩展。这也是将现有 PHP-MySQL、PHP-PostgreSQL 或 PHP-MS-SQL 应用程序转换成 PHP-ODBC 应用程序有点棘手的原因。因此，当从头开始构建应用程序时，必须小心使用严格标准化的 SQL（在多数情况下，扩展的 SQL 很常见）。

如上所述，使用 ODBC 连接到电子表格是有可能的。在使用数据库时，必须使用连接器实现上述连接。有许多工具可用来解决此问题，其中一些工具是开放源码的，但许多是专有的。Microsoft Office for Windows 就是这样一个例子，它为 Excel 电子表格提供了 ODBC 连接器，在通过 ODBC 处理电子表格时，它可能显得有点笨拙，可是，如果您将简单的电子表格转换成数据库表，ODBC 可能会解决许多麻烦。在建立电子表格 ODBC 连接之后，可以或多或少地将它看作是一个数据库连接，这种情况同样适用于 ODBC PHP 函数，但要使用类似于 SQL 的语言，而在 Excel 中，则会使用标准电子表格公式。

## 对链接数据和语义 Web 的影响

链接数据移动看起来就像是通过 Web 连接数据。这样做有很多好处：主要是让机器了解数据的特定元素，对于用户，可以使他们能够更容易地找到信息。链接数据移动使用了语义 Web 中先前已经存在的标准（比如资源描述框架和 Web Ontology Language）和 Internet/Web 的标准化（如 HTTP 和 OpenID）。您可能慢慢开始了解到，链接数据连接方法有些类似于 ODBC，因为通过已建立的连接，URI 有些类似于连接字符串，而语义 Web 查询语言（Semantic Web Query Language, SPARQL）有些类似于 SQL。

可以将这个理论再扩展一下：链接数据有些类似于 ODBC，它可能会建立到链接数据存储（比如 “三重存储”）的 ODBC 连接，并将 SPARQL 查询发送 到 ODBC 连接。使用 OpenLink Virtuoso 时可能出现这种情况，在这种情况下，允许您通过标准 ODBC 连接来建立连接。

## 许可和 ODBC

iODBC 需要具有 GNU 通用公共许可证（General Public License, GPL）和 Berkeley 软件开发开放源码许可证的双重许可才可以使用。UnixODBC 也需要具有 GPL 开放源码许可证的许可才可以使用。这意味着，无论您使用这两个库开发什么软件，该软件不必是开放源码的，可以是专有的。Microsoft 的 ODBC 驱动程序也可以成为专有软件的一部分，但会受到 Microsoft Office 软件（在这里是 Access 数据库和 Excel 电子表格）和 SQL Server（在这里是 SQL Server 数据库）的许可协议的约束。

## 结束语

ODBC 是一项伟大的技术，可最大限度地提高通用连接性。ODBC 提高了效率，允许您扩展应用程序，以便处理新形式的数据，比如基于 Web 的链接数据。但它自身也有一些缺点：为了实现通用连接，您必须仔细选择构建 SQL 查询的方式，因为只能在所有数据库管理系统之间使用所有可用 SQL 命令的子集。我们希望本文能为您提供通过 ODBC 使用 PHP 编程语言处理数据库所需的所有知识。