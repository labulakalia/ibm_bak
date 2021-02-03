# 用 Groovy 进行 JDBC 编程
用 GroovySql 构建下一个报告应用程序

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/java-j-pg01115/)

Andrew Glover

发布: 2005-01-11

* * *

在 _实战 Groovy_ 系列的前几期中，您已经了解了 Groovy 的一些非常优美的特性。在 [第 1 期](/zh/articles/j-pg11094/) 中，学习了如何用 Groovy 对普通的 Java™ 代码进行更简单、更迅速的单元测试。在 [第 2 期](/zh/articles/j-pg12144/) 中，看到了 Groovy 能够给 Ant 构建带来的表现能力。这一次您会发现 Groovy 的另一个实际应用，即如何用它迅速地构建基于 SQL 的报告应用程序。

脚本语言对于迅速地构建报告应用程序来说是典型的优秀工具，但是构建这类应用程序对于 Groovy 来说特别容易。Groovy 轻量级的语法可以消除 Java 语言的 JDBC 的一些繁冗之处，但是它真正的威力来自闭包，闭包很优雅地把资源管理的责任从客户机转移到框架本身，因此更容易举重若轻。

在本月的文章中，我先从 GroovySql 的概述开始，然后通过构建一个简单的数据报告应用程序，向您演示如何将它们投入工作。为了从讨论中得到最大收获，您应当熟悉 Java 平台上的 JDBC 编程。您可能还想回顾 [上个月对 Groovy 中闭包的介绍](//zh/articles/j-pg12144/) ，因为它们在这里扮演了重要的角色。但是，本月关注的重点概念是迭代（iteration），因为迭代器在 Groovy 对 JDBC 的增强中扮演着重要角色。所以，我将从 Groovy 中迭代器的概述开始。

## 进入迭代器

迭代是各种编程环境中最常见、最有用的技术。 _迭代器_ 是某种代码助手，可以让您迅速地访问任何集合或容器中的数据，每次一个数据。Groovy 把迭代器变成隐含的，使用起来更简单，从而改善了 Java 语言的迭代器概念。在清单 1 中，您可以看到使用 Java 语言打印 `String` 集合的每个元素需要做的工作。

##### 清单 1\. 普通 Java 代码中的迭代器

```
import java.util.ArrayList;
import java.util.Collection;
import java.util.Iterator;
public class JavaIteratorExample {
public static void main(String[] args) {
     Collection coll = new ArrayList();
     coll.add("JMS");
     coll.add("EJB");
     coll.add("JMX");
     for(Iterator iter = coll.iterator(); iter.hasNext();){
        System.out.println(iter.next());
     }
}
}

```

Show moreShow more icon

在清单 2 中，您可以看到 Groovy 如何简化了我的工作。在这里，我跳过了 `Iterator` 接口，直接在集合上使用类似迭代器的方法。而且， Groovy 的迭代器方法接受闭包，每个迭代中都会调用闭包。清单 2 显示了前面基于 Java 语言的示例用 Groovy 转换后的样子。

##### 清单 2\. Groovy 中的迭代器

```
class IteratorExample1{
static void main(args) {
     coll = ["JMS", "EJB", "JMX"]
     coll.each{ item | println item }
}
}

```

Show moreShow more icon

您可以看到，与典型的 Java 代码不同，Groovy 在允许我传递进我需要的行为的同时，控制了特定于迭代的代码。使用这个控制，Groovy 干净漂亮地把资源管理的责任从我手上转移到它自己身上。让 Groovy 负责资源管理是极为强大的。它还使编程工作更加容易，从而也就更迅速。

##### 关于本系列

把任何一个工具集成进您的开发实践的关键是，知道什么时候使用它而什么时候把它留在箱子中。脚本语言可以是您的工具箱中极为强大的附件，但是只有在恰当地应用到适当的场景时才是这样。

## GroovySql 简介

Groovy 的 SQL 魔力在于一个叫做 GroovySql 的精致的 API。使用闭包和迭代器，GroovySql 干净漂亮地把 JDBC 的资源管理职责从开发人员转移到 Groovy 框架。这么做之后，就消除了 JDBC 编程的繁琐，从而使您可以把注意力放在查询和查询结果上。

如果您忘记了普通的 Java JDBC 编程有多麻烦，我会很高兴提醒您！在清单 3 中，您可以看到用 Java 语言进行的一个简单的 JDBC 编程示例。

##### 清单 3\. 普通 Java 的 JDBC 编程

```
import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;
public class JDBCExample1 {
public static void main(String[] args) {
    Connection con = null;
    Statement stmt = null;
    ResultSet rs = null;
    try{
      Class.forName("org.gjt.mm.mysql.Driver");
      con = DriverManager.getConnection("jdbc:mysql://localhost:3306/words",
           "words", "words");
      stmt = con.createStatement();
      rs = stmt.executeQuery("select * from word");
      while (rs.next()) {
        System.out.println("word id: " + rs.getLong(1) +
            " spelling: " + rs.getString(2) +
            " part of speech: " + rs.getString(3));
      }
    }catch(SQLException e){
      e.printStackTrace();
    }catch(ClassNotFoundException e){
      e.printStackTrace();
    }finally{
      try{rs.close();}catch(Exception e){}
      try{stmt.close();}catch(Exception e){}
      try{con.close();}catch(Exception e){}
}
}
}

```

Show moreShow more icon

哇。清单 3 包含近 40 行代码，就是为了查看表中的内容！如果用 GroovySql，您猜要用多少行？如果您猜超过 10 行，那么您就错了。请看清单 4 中，Groovy 替我处理底层资源，从而非常漂亮地让我把注意力集中在手边的任务 —— 执行一个简单的查询。

##### 清单 4\. 欢迎使用 GroovySql ！

```
import groovy.sql.Sql
class GroovySqlExample1{
static void main(args) {
    sql = Sql.newInstance("jdbc:mysql://localhost:3306/words", "words",
           "words", "org.gjt.mm.mysql.Driver")
    sql.eachRow("select * from word"){ row |
       println row.word_id + " " + row.spelling + " " + row.part_of_speech
    }
}
}

```

Show moreShow more icon

真不错。只用了几行，我就编码出与清单 3 相同的行为，不用关闭 `Connection` ，也不用关闭 `ResultSet` ，或者在 JDBC 编程中可以找到的任何其他熟悉的重要特性。这是多么激动人心的事情啊，并且还是如此容易。现在让我来仔细介绍我是如何做到的。

## 执行简单的查询

在 [清单 4\. 欢迎使用 GroovySql ！](#清单-4-欢迎使用-groovysql-！) 的第一行中，我创建了 Groovy 的 `Sql` 类的实例，用它来连接指定的数据库。在这个例子中，我创建了 `Sql` 实例，指向在我机器上运行的 MySQL 数据库。到现在为止都非常基本，对么？真正重要的是一下部分，迭代器和闭包一两下就显示出了它们的威力。

请把 `eachRow` 方法当成传进来的查询生成的结果上的迭代器。在底层，您可以看到返回了 JDBC `ResultSet` 对象，它的内容被传递进 `for` 循环。所以，每个迭代都要执行我传递进去的闭包。如果在数据库中找到的 `word` 表只有三行，那么闭包就会执行三次 —— 打印出 `word_id` 、 `spelling` 和 `part_of_speech` 的值。

如果将等式中我指定的变量 `row` 去掉，而使用 Groovy 的一个隐含变量 `it` （它恰好就是迭代器的实例），代码可以进一步简化。如果我这样做，那么前面的代码就可以写成清单 5 所示的这样：

##### 清单 5\. GroovySql 中的 `it` 变量

```
import groovy.sql.Sql
class GroovySqlExample1{
static void main(args) {
    sql = Sql.newInstance("jdbc:mysql://localhost:3306/words", "words",
           "words", "org.gjt.mm.mysql.Driver")
    sql.eachRow("select * from word"){ println it.spelling +  " ${it.part_of_speech}"}
}
}

```

Show moreShow more icon

在这个代码中，我可以删除 `row` 变量，用 `it` 代替。而且，我还能在 `String` 语句中引用 `it` 变量，就像我在 `${it.part_of_speech}` 中所做的那样。

## 执行更复杂的查询

前面的例子相当简单，但是 GroovySql 在处理更复杂的数据操纵查询（例如 `insert` 、 `update` 和 `delete` 查询）时，也是非常可靠的。 对于这些查询，您没有必要用迭代器，所以 Groovy 的 `Sql` 对象另外提供了 `execute` 和 `executeUpdate` 方法。这些方法让人想起普通的 JDBC `statement` 类，它也有 `execute` 和 `executeUpdate` 方法。

在清单 6 中，您看到一个简单的 `insert` ，它再次以 `${}` 语法使用变量替换。这个代码只是向 `word` 表插入一个新行。

##### 清单 6\. 用 GroovySql 进行插入

```
wid = 999
spelling = "Nefarious"
pospeech = "Adjective"
sql.execute("insert into word (word_id, spelling, part_of_speech)
values (${wid}, ${spelling}, ${pospeech})")

```

Show moreShow more icon

Groovy 还提供 `execute` 方法的一个重载版本，它接收一列值，这些值与查询中发现的 `?` 元素对应。在清单 7 中，我简单地查询了 `word` 表中的某个行。在底层，GroovySql 创建了普通 Java 语言 `java.sql.PreparedStatement` 的一个实例。

##### 清单 7\. 用 GroovySql 创建 PreparedStatement 的实例

```
val = sql.execute("select * from word where word_id = ?", [5])

```

Show moreShow more icon

更新的方式基本相同，也使用 `executeUpdate` 方法。还请注意，在清单 8 中 `executeUpdate` 方法接收一列值，与查询中的 `?` 元素对应。

##### 清单 8\. 用 GroovySql 进行更新

```
nid = 5
spelling = "Nefarious"
sql.executeUpdate("update word set word_id = ? where spelling = ?", [nid, spelling])

```

Show moreShow more icon

删除实际上与插入相同，当然，语法不同，如清单 9 所示。

##### 清单 9\. 用 GroovySql 进行删除

```
sql.execute("delete from word where word_id = ?" , [5])

```

Show moreShow more icon

## 简化数据操纵

任何想简化 JDBC 编程的 API 或工具最好有一些好的数据操纵特性，在这一节中，我要向您再介绍三个。

### 数据集 (DataSet)

构建于 GroovySql 简单性的基础之上，GroovySql 支持 `DataSet` 类型的概念，这基本上是数据库表的对象表示。使用 `DataSet` ，您可以在行中遍历，也可以添加新行。实际上，用数据集是方便地表示表格的公共数据集合的方式。

但是，目前 GroovySql `DataSet` 类型的不足之处是它们没有代表关系；它们只是与数据库表的一对一映射。在清单 10 中，我创建了一个来自 `word` 表的 `DataSet` 。

##### 清单 10\. 用 GroovySql 创建数据集

```
import groovy.sql.Sql
class GroovyDatasetsExample1{
static void main(args) {
    sql = Sql.newInstance("jdbc:mysql://localhost:3306/words", "words",
          "words", "org.gjt.mm.mysql.Driver")
    words = sql.dataSet("word")
    words.each{ word |
     println word.word_id + " " + word.spelling
    }
    words.add(word_id:"9999", spelling:"clerisy", part_of_speech:"Noun")
}
}

```

Show moreShow more icon

您可以看到，GroovySql 的 `DataSet` 类型可以容易地用 `each` 方法对表的内容进行遍历，容易地用 `add` 方法添加新行， `add` 方法接受一个 `map` 表示需要的数据。

### 使用存储过程和负索引

存储过程调用和负索引（negative indexing）可能是数据操纵的重要方面。GroovySql 使存储过程调用简单得就像在 `Sql` 类上使用 `call` 方法一样。对于负索引， GroovySql 提供了自己增强的 `ResultSet` 类型，它工作起来非常像 Groovy 中的 _collections_ 。例如，如果您想获取结果集中的最后一个项目，您可以像清单 11 所示的那样做：

##### 清单 11\. 用 GroovySql 进行负索引

```
sql.eachRow("select * from word"){ grs |
println "-1  = " + grs.getAt(-1) //prints spelling
println "2  = " + grs.getAt(2) //prints spelling
}

```

Show moreShow more icon

您在清单 11 中可以看到，提取结果集的最后一个元素非常容易，只要用 -1 作索引就可以。如果想试试，我也可以用索引 2 访问同一元素。

这些例子非常简单，但是它们能够让您很好地感觉到 GroovySql 的威力。我现在要用一个演示目前讨论的所有特性的实际例子来结束本月的课程。

## 编写一个简单的报告应用程序

报告应用程序通常要从数据库拖出信息。在典型的业务环境中，可能会要求您编写一个报告应用程序，通知销售团队当前的 Web 销售情况，或者让开发团队对系统某些方面（例如系统的数据库）的性能进行日常检测。

为了继续这个简单的例子，假设您刚刚部署了一个企业范围的 Web 应用程序。当然，因为您在编写代码时还（用 Groovy）编写了充足的单元测试，所以它运行得毫无问题；但是您还是需要生成有关数据库状态的报告，以便调优。您想知道客户是如何使用应用程序的，这样才能发现性能问题并解决问题。

通常，时间约束限制了您在这类应用程序中能够使用的提示信息的数量。但是您新得到的 GroovySql 知识可以让您轻而易举地完成这个应用程序，从而有时间添加更多您想要的特性。

### 细节

在这个例子中，您的目标数据库是 MySQL，它恰好支持用查询发现状态信息这一概念。以下是您有兴趣的状态信息：

- 运行时间。
- 处理的全部查询数量。
- 特定查询的比例，例如 `insert` 、 `update` 和 `select` 。

用 GroovySql 从 MySQL 数据库得到这个信息太容易了。由于您正在为开发团队构建状态信息，所以您可能只是从一个简单的命令行报告开始，但是您可以在后面的迭代中把报告放在 Web 上。这个报告例子的用例看起来可能像这样：

1.连接到我们的应用程序的活动数据库。2.发布 `show status` 查询并捕获：a. 运行时间b. 全部查询数c. 全部 `insert` 数d. 全部 `update` 数e. 全部 `select` 数3.使用这些数据点，计算：a. 每分钟查询数b. 全部 `insert` 查询百分比c. 全部 `update` 查询百分比d. 全部 `select` 查询百分比

在清单 12 中，您可以看到最终结果：一个将会报告所需数据库统计信息的应用程序。代码开始的几行获得到生产数据库的连接，接着是一系列 `show status` 查询，让您计算每分钟的查询数，并按类型把它们分开。请注意像 `uptime` 这样的变量如何在定义的时候就创建。

##### 清单 12\. 用 GroovySql 进行数据库状态报告

```
import groovy.sql.Sql
class DBStatusReport{
static void main(args) {
     sql = Sql.newInstance("jdbc:mysql://yourserver.anywhere/tiger", "scott",
        "tiger", "org.gjt.mm.mysql.Driver")
     sql.eachRow("show status"){ status |
        if(status.variable_name == "Uptime"){
           uptime =  status[1]
        }else if (status.variable_name == "Questions"){
           questions =  status[1]
        }
     }
     println "Uptime for Database: " + uptime
     println "Number of Queries: " + questions
     println "Queries per Minute = "
         + Integer.valueOf(questions)/Integer.valueOf(uptime)
     sql.eachRow("show status like 'Com_%'"){ status |
        if(status.variable_name == "Com_insert"){
           insertnum =  Integer.valueOf(status[1])
        }else if (status.variable_name == "Com_select"){
           selectnum =  Integer.valueOf(status[1])
        }else if (status.variable_name == "Com_update"){
           updatenum =  Integer.valueOf(status[1])
       }
    }
    println "% Queries Inserts = " + 100 * (insertnum / Integer.valueOf(uptime))
    println "% Queries Selects = " + 100 * (selectnum / Integer.valueOf(uptime))
    println "% Queries Updates = " + 100 * (updatenum / Integer.valueOf(uptime))
    }
}

```

Show moreShow more icon

## 结束语

在 _实战 Groovy_ 本月的这一期文章中，您看到了 GroovySql 如何简化 JDBC 编程。这个干净漂亮的 API 把闭包和迭代器与 Groovy 轻松的语法结合在一起，有助于在 Java 平台上进行快速数据库应用程序开发。最强大的是，GroovySql 把资源管理任务从开发人员转移到底层的 Groovy 框架，这使您可以把精力集中在更加重要的查询和查询结果上。但是不要只记住我这句话。下次如果您被要求处理 JDBC 的麻烦事，那时可以试试小小的 GroovySql 魔力。然后给我发封电子邮件告诉我您的体会。

在 _实战 Groovy_ 下一月的文章中，我要介绍 Groovy 模板框架的细节。您会发现，用这个更加聪明的框架创建应用程序的视图组件简直就是小菜一碟。