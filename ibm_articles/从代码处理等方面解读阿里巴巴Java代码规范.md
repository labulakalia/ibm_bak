# 从代码处理等方面解读阿里巴巴 Java 代码规范
阿里巴巴的 Java 代码规范深度解读的下篇

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/deconding-code-specification-part-2/)

周明耀

发布: 2018-01-17

* * *

## 前言

2017 年阿里云栖大会，阿里发布了针对 Java 程序员的《阿里巴巴 Java 开发手册（终极版）》，这篇文档作为阿里数千位 Java 程序员的经验积累呈现给公众，并随之发布了适用于 Eclipse 和 Intellim 的代码检查插件。为了能够深入了解 Java 程序员编码规范，也为了深入理解为什么阿里这样规定，是否规定有误，本文以阿里发布的这篇文档作为分析起源，扩大范围至业界其他公司的规范，例如谷歌、FaceBook、微软、百度、华为，并搜索网络上技术大牛发表的技术文章，深入理解每一条规范的设计背景和目标。

由于解读文章仅有两篇，所以按照阿里的篇幅权重分为上篇仅针对 Java 语言本身的编码规约，下篇包含日志管理、异常处理、单元测试、MySQL 规范、工程规范等方面内容进行解读。本文是下篇，主要针对编码规约部分进行解读，注意，本文所附代码限于篇幅可能并不完整，也可能由于机器不同，存在运行结果不完全一致的情况，请读者见谅。

## 异常日志

### 异常处理

#### 不要捕获 RuntimeException

阿里强制规定 Java 类库中的 RuntimeException 可以通过预先检查进行规避，而不应该通过 catch 来处理，例如 IndexOutOfBoundsException、NullPointerException 等。

##### 我的理解

RuntimeException，也被称为运行时异常，通常是由于代码中的 bug 引起的，正确的处理方式是去检查代码，通过添加数据长度判断，判断对象是否为空等方法区规避，而不是靠捕获来规避这种异常。

#### 事务中的异常需要回滚

阿里强制规定有 try 块放到了事务代码中，catch 异常后，如果需要回滚事务，一定要注意手动回滚事务。

##### 我的理解

try catch 代码块中对异常的处理，可能会遗漏事务的一致性，当事务控制不使用其他框架管理时，事务需要手动回滚。实际使用如果引入第三方的框架对事务进行管理，比如 Spring，则根据第三方框架的实际实现情况，确定是否有必要手动回滚。当第三方事务管理框架本身就会对于异常进行抛出时需要做事务回滚。例如 Spring 在@Transactional 的 annotation 注解下，会默认开启运行时异常事务回滚。

#### 不能在 finally 块中使用 return

阿里强制要求 finally 块中不使用 return，因为执行完该 return 后方法结束执行，不会再执行 try 块中的 return 语句。

##### 我的理解

我们来看一个示例，代码如清单 1 所示。

##### 清单 1 示例代码

```
public class demo{
    public static void main(String[] args){
        System.out.println(m_1());
    }
    public int m_1(){
        int i = 10;
        try{
            System.out.println("start");
            return i += 10;
        }catch(Exception e){
            System.out.println("error:"+e);
        }finally{
            if(i>10){
            System.out.println(i);
    }
     System.out.println("finally");
     return 50;
        }
    }
}

```

Show moreShow more icon

输出如下清单 2 所示。

##### 清单 2 清单 1 程序运行输出

```
start
20
finally
50

```

Show moreShow more icon

对此现象可以通过反编译 class 文件很容易找到原因，如清单 3 所示。

##### 清单 3 反编译文件

```
public class demo{
    public static void main(String[] args){
    System.out.println(m_1());
}
    public int m_1(){
        int i = 10;
        try{
             System.out.println("start");
             return i;
        }catch(Exception e){
             System.out.println("error:"+e);
        }finally{
             if(i>10){
             System.out.println(i);
        }
             System.out.println("finally");
             i = 50;
        }
             return i;
        }
    }
}

```

Show moreShow more icon

首先 i+=10;被分离为单独的一条语句，其次 return 50;被加在 try 和 catch 块的结尾，”覆盖”了 try 块中原有的返回值。如果我们在 finally 块中没有 return，则 finally 块中的赋值语句不会改变最后的返回结果，代码如清单 4 所示。

##### 清单 4 finally 块示例代码

```
public class demo{
    public static void main(String[] args){
        System.out.println(m_1());
}
    public static int m_1(){
        int i = 10;
        try{
            System.out.println("start");
            return i += 10;
        }catch(Exception e){
            System.out.println("error:"+e);
        }finally{
            if(i>10){
            System.out.println(i);
        }
            System.out.println("finally");
            i = 50;
        }
            return i;
    }
}

```

Show moreShow more icon

输出结果如清单 5 所示。

##### 清单 5 清单 4 程序运行输出

```
start
finally
10

```

Show moreShow more icon

### 日志规约

#### 不可直接使用日志系统

阿里强制规定应用中不可直接使用日志系统（Log4j、Logback）中的 API，而应依赖使用日志框架 SLF4J 中的 API，使用门面模式的日志框架，有利于维护和各个类的日志处理方式统一。

##### 我的理解

SLF4J 即简单日志门面模式，不是具体的日志解决方案，它只服务于各种各样的日志系统。在使用 SLF4J 时不需要指定哪个具体的日志系统，只需要将使用到的具体日志系统的配置文件放到类路径下去。

正例如清单 6 所示。

##### 清单 6 正例代码

```

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
public class HelloWorld{
private static final Logger logger =
    LoggerFactory.getLogger(HelloWorld.class);
public static void main(String[] args){
logger.info("please use SLF4J,rather than logback or log4j");
}
}

```

Show moreShow more icon

反例如清单 7 所示。

##### 清单 7 反例代码

```

import org.apache.log4j.Logger;
public class HelloWorld{
    private static final Logger logger =
    LoggerFactory.getLogger(HelloWorld.class);
    public static void main(String[] args){
        logger.info("please use SLF4J,rather than logback or log4j");
    }
}

```

Show moreShow more icon

#### 日志文件保留时间

阿里强制规定日志文件至少保存 15 天，因为有些异常具备以”周”为频次发生的特点。

##### 我的理解

日志保留时间推荐 15 天以上，但是保留时间也不宜过长，一般不超过 21 天，否则造成硬盘空间的浪费。对于一些长周期性执行的逻辑，可以根据实际情况调整该保存时间，同时也需要保证日志能够监控到关键的应用。

对于长周期执行的逻辑，可以使用特定的 appender，并使用不同的日志清理规则，如时间、大小等。如一月执行一次的定时任务，可以将日志输出到新的日志文件，然后通过大小限定的规则进行清理，并不一定要使用时间清理的逻辑。

### 安全规约

#### 权限控制校验

阿里强制要求对于隶属于用户个人的页面或者功能必须进行权限控制校验。

##### 我的理解

涉及到对于数据的增删改查，必须有权限的控制和校验，要有一个黑白名单的控制，不能依赖于前台页面的简单控制，后台要有对于完整的权限控制的实现。这样就能尽可能地防治数据的错误修改。

#### 用户传入参数校验

阿里强制要求用户请求传入的任何参数必须做有效校验。

##### 我的理解

对于用户输入的任何参数，前端页面上都必须要做一定的有效性校验，并且在数据发送至服务器的时候在页面上给出验证结果提示，那么在用户请求传入的任务参数，后台同样也要对其有效性进行验证，防止前端页面未能过滤或者暂时无法验证的错误参数。忽略参数的验证会导致的问题很多，page size 过大会导致内存溢出、SQL 溢出等，只有验证才能尽可能地减少这些问题的出现，进而减少错误的排查几率。

### 单元测试

#### 单元测试应该自动执行

阿里强制单元测试应该是全自动执行的，并且非交互式的。测试框架通常是定期执行的，执行过程必须完全自动化才有意义。输出结果需要人工检查的测试不是一个号的单元测试。单元测试中不准使用 System.out 来进行人肉验证，必须使用 assert 来验证。

##### 我的理解

这条原则比较容易理解。单元测试是整个系统的最小测试单元，针对的是一个类中一个方法的测试，如果这些测试的结果需要人工校验是否正确，那么对于验证人来说是一项痛苦而且耗时的工作。另外，单元测试作为系统最基本的保障，需要在修改代码、编译、打包过程中都会运行测试用例，保障基本功能，自动化的测试是必要条件。其实自动化测试不仅是单元测试特有的，包括集成测试、系统测试等，都在慢慢地转向自动化测试，以降低测试的人力成本。

#### 单元测试应该是独立的

阿里强制保持单元测试的独立性。为了保证单元测试稳定可靠且便于维护，单元测试用例之间决不能互相调用，也不能依赖执行的先后次序。反例：method2 需要依赖 method1 的执行，将执行结果作为 method2 的输入。

##### 我的理解

单元测试作为系统的最小测试单元，主要目的是尽可能早地测试编写的代码，降低后续集成测试期间的测试成本，以及在运行测试用例的时候能够快速地定位到对应的代码段并解决相关问题。

我们假设这么一个场景，method1 方法被 10 个其他 method 方法调用，如果 10 个 method 方法的测试用例都需要依赖 method1，那么当 methdo1 被修改导致运行出错的情况下，会导致 method1 以及依赖它的 10 个 method 的所有测试用例报错，这样就需要排查这 11 个方法到底哪里出了问题，这与单元测试的初衷不符，也会大大的增加排查工作量，所以单元测试必须是独立的，不会因为受到外部修改（这里的修改包括了依赖方法的修改以及外部环境的修改），编写单元测试时遇到的这类依赖可以使用 mock 来模拟输入和期望的返回，这样所以来的方法内部逻辑的变更就不会影响到外部的实现。

#### BCDE 原则

阿里推荐编写单元测试代码遵守 BCDE 原则，以保证被测试模块的交付质量。

##### 我的理解

BCDE 原则逐一解释如下：

B（Border）：确保参数边界值均被覆盖。

例如：对于数字，测试负数、0、正数、最小值、最大值、NaN（非数字）、无穷大值等。对于字符串，测试空字符串、单字符、非 ASCII 字符串、多字节字符串等。对于集合类型，测试空、第一个元素、最后一个元素等。对于日期，测试 1 月 1 日、2 月 29 日、12 月 31 日等。被测试的类本身也会暗示一些特定情况下的边界值。对于边界情况的测试一定要详尽。

C（Connect）：确保输入和输出的正确关联性。

例如，测试某个时间判断的方法 boolean inTimeZone(Long timeStamp)，该方法根据输入的时间戳判断该事件是否存在于某个时间段内，返回 boolean 类型。如果测试输入的测试数据为 Long 类型的时间戳，对于输出的判断应该是对于 boolean 类型的处理。如果测试输入的测试数据为非 Long 类型数据，对于输出的判断应该是报错信息是否正确。

D（Design）：任务程序的开发包括单元测试都应该遵循设计文档。

E（Error）：单元测试包括对各种方法的异常测试，测试程序对异常的响应能力。

除了这些解释之外，《单元测试之道（Java 版）》这本书里面提到了关于边界测试的 CORRECT 原则：

一致性（Conformance）：值是否符合预期格式（正常的数据），列出所有可能不一致的数据，进行验证。

有序性（Ordering）：传入的参数的顺序不同的结果是否正确，对排序算法会产生影响，或者是对类的属性赋值顺序不同会不会产生错误。

区间性（Range）：参数的取值范围是否在某个合理的区间范围内。

引用/耦合性（Reference）：程序依赖外部的一些条件是否已满足。前置条件：系统必须处于什么状态下，该方法才能运行。后置条件，你的方法将会保证哪些状态发生改变。

存在性（Existence）：参数是否真的存在，引用为 Null，String 为空，数值为 0 或者物理介质不存在时，程序是否能正常运行。

基数性（Cardinality）：考虑以”0-1-N 原则”，当数值分别为 0、1、N 时，可能出现的结果，其中 N 为最大值。

时间性（Time）：相对时间指的是函数执行的依赖顺序，绝对时间指的是超时问题、并发问题。

#### 建表的是与否规则

阿里强制要求如果遇到需要表达是与否的概念时，必须使用 is\_xxx 的方法命令，数据类型是 unsigned tinyint，1 表示是，0 表示否。

说明：任务字段如果为非负数，必须是 unsigned。

正例：表达逻辑删除的字段名 is\_deleted，1 表示删除，0 表示未删除。

##### 我的理解

命名使用 is\_xxx 第一个好处是比较清晰的，第二个好处是使用者根据命名就可以知道这个字段的取值范围，也方便做参数验证。

类型使用 unsigned 的好处是如果只存整数，unsigned 类型有更大的取值范围，可以节约磁盘和内存使用。

对于表的名字，MySQL 社区有自己推荐的命名规范：

1. 表包含多个英文单词时，需要用下划线进行单词分割，这一点类似于 Java 类名的命名规范，例如 master\_schedule、security\_user\_permission；
2. 由于 InnoDB 存储引擎本身是针对操作系统的可插拔设计的，所以原则上所有的表名组成全部由小写字母组成；
3. 不允许出现空格，需要分割一律采用下划线；
4. 名字不允许出现数字，仅包含英文字母；
5. 名字需要总长度少于 64 个字符。

#### 数据类型精度考量

阿里强制要求存放小数时使用 decimal，禁止使用 float 和 double。

说明：float 和 double 在存储的时候，存在精度损失的问题，很可能在值的比较时，得到不正确的结果。如果存储的数据范围超过 decimal 的范围，建议将数据拆成整数和小数分开存储。

##### 我的理解

我们先来看看各个精度的范围。

Float：浮点型，4 字节数 32 位，表示数据范围-3.4E38~3.4E38

Double：双精度型，8 字节数 64 位，表示数据范围-1.7E308~1.7E308

Decimal：数字型，16 字节数 128 位，不存在精度损失，常用于银行账目计算

在精确计算中使用浮点数是非常危险的，在对精度要求高的情况下，比如银行账目就需要使用 Decimal 存储数据。

实际上，所有涉及到数据存储的类型定义，都会涉及数据精度损失问题。Java 的数据类型也存在 float 和 double 精度损失情况，阿里没有指出这条规约，就全文来说，这是一个比较严重的规约缺失。

Joshua Bloch（著名的 Effective Java 书作者）认为，float 和 double 这两个原生的数据类型本身是为了科学和工程计算设计的，它们本质上都采用单精度算法，也就是说在较宽的范围内快速获得精准数据值。但是，需要注意的是，这两个原生类型都不保证也不会提供很精确的值。单精度和双精度类型特别不适用于货币计算，因为不可能准确地表示 0.1（或者任何其他十的负幂）。

举个例子，如清单 8 所示。

##### 清单 8 示例代码

```
float calnUM1;
double calNum2;
calNum1 = (float)(1.03-.42);
calNum2 = 1.03-.42;
System.out.println("calNum1="+ calNum1);
System.out.println("calNum2="+ calNum2);
System.out.println(1.03-.42);
calNum1 = (float)(1.00-9*.10);
calNum2 = 1.00-9*.10;
System.out.println("calNum1="+ calNum1);
System.out.println("calNum2="+ calNum2);
System.out.println(1.00-9*.10);

```

Show moreShow more icon

输出结果如清单 9 所示。

##### 清单 9 清单 8 示例代码运行输出结果

```
calNum1=0.61
calNum2=0.6100000000000001
0.6100000000000001
calNum1=0.1
calNum2=0.09999999999999998
0. 09999999999999998

```

Show moreShow more icon

从上面的输出结果来看，如果寄希望于打印时自动进行四舍五入，这是不切实际的。

我们再来看一个实际的例子。假设你有 1 块钱，现在每次购买蛋糕的价格都会递增 0.10 元，为我们一共可以买几块蛋糕。口算一下，应该是 4 块（因为 0.1+0.2+0.3+0.4=1.0），我们写个程序验证看看，如清单 10 所示。

##### 清单 10 示例代码

```
//错误的方式
double funds1 = 1.00;
int itemsBought = 0;
for(double price = .10;funds>=price;price+=.10){
    funds1 -=price;
    itemsBought++;
}
    System.out.println(itemsBought+" items boughts.");
    System.out.println("Changes:"+funds1);
    //正确的方式
    final BigDecimal TEN_CENTS = new BigDecimal(".10");
    itemsBought = 0;
    BigDecimal funds2 = new BigDecimal("1.00");
for(BigDecimal price = TEN_CENTS;funds2.compareTo(price)>0;price =
    price.add(TEN_CENTS)){
    fund2 = fund2.substract(price);
    itemsBought++;
    }
    System.out.println(itemsBought+" items boughts.");
    System.out.println("Changes:"+funds2);

```

Show moreShow more icon

运行输出如清单 11 所示。

##### 清单 11 清单 10 示例代码运行输出结果

```
3 items boughts.
Changes:0.3999999999999999
4 items boughts.
Changes:0.00

```

Show moreShow more icon

这里我们可以看到使用了 BigDecimal 解决了问题，实际上 int、long 也可以解决这类问题。采用 BigDecimal 有一个缺点，就是使用过程中没有原始数据这么方便，效率也不高。如果采用 int 方式，最好不要在有小数点的场景下使用，可以在 100、10 这样业务场景下选择使用。

#### 使用 Char

阿里强制要求如果存储的字符串长度几乎相等，使用 Char 定长字符串类型。

##### 我的理解

我这里不讨论 MySQL，而是聊聊另一种主流关系型数据库-PostgreSQL。在 PostgreSQL 中建议使用 varchar 或者 text，而不是 char，这是因为它们之间没有性能区别，但是 varchar、text 能支持动态的长度调整，存储空间也更节省。

在 PostgreSQL 官方文档中记录了这两种类型的比较，如下所示：

SQL 定义了两种基本的字符类型：character varying(n)和 character(n)，这里的 n 是一个正整数。两种类型都可以存储最多 n 个字符的字符串（没有字节）。试图存储更长的字符串到这些类型的字段里会产生一个错误，除非超出长度的字符都是空白，这种情况下该字符串将被截断为最大长度。如果要存储的字符串比申明的长度短，类型为 character 的数值将会用空白填满，而类型为 character varying 的数值就不会填满。

如果我们明确地把一个数值转换成 character varying(n)或 character(n)，那么超长的数值将被截断成 n 个字符，且不会抛出错误。这也是 SQL 标准的要求。

varchar(n)和 char(n)分别是 character varying(n)和 character(n)的别名，没有申明长度的 character 等于 character(1)；如果不带长度说明可以使用 character varying，那么该类型接受任何长度的字符串。

另外，PostgreSQL 提供 text 类型，它可以存储任何长度的字符串。尽管类型 text 不是 SQL 标准，但是许多其他 SQL 数据库系统也有它。

Character 类型的数值物理上都用空白填充到指定的长度 n，并且以这种方式存储。不过，填充的空白是暂时无语义的。在比较两个 character 值的时候，填充的空白都不会被关注，在转换成其他字符串类型的时候，character 值里面的空白会被删除。请注意，在 character varying 和 text 数值的，结尾的空白是有语义的，并且当使用模式匹配时，如 LIKE，使用正则表达式。

一个简短的字符串（最多 126 个字节）的存储要求是 1 个字节加上实际的字符串，其中包括空格填充的 character。更长的字符串有 4 个字节的开销，而不是 1。长的字符串将会自动被系统压缩，因此在磁盘上的物理需求可能会更少些。更长的数值也会存储在数据表里面，这样它们就不会干扰对短字段值的快速访问。不管怎样，允许存储的最长字符串大概是 1GB。允许在数据类型声明中出现的 n 的最大值比这个还小。修改该行也没有什么意义，因为在多字节编码下字符和字节的数目可能差别很大。如果你想存储没有特定上限的长字符串，那么使用 text 或没有长度声明的 character varying，而不要选择一个任意长度限制。

从性能上分析，character(n)通常是最慢的，在大多数情况下，应该使用 text 或者 character varying。

### 工程结构

#### 应用分层

##### 服务间依赖关系

阿里推荐默认上层依赖于下层，箭头关系表示可直接依赖，如：开放接口层可以依赖于 Web 层，也可以直接依赖于 Service 层。

##### 我的理解

《软件架构模式》一书中介绍了分层架构思想：

分层架构是一种很常见的架构模式，它也被叫做 N 层架构。这种架构是大多数 Java EE 应用的实际标准。许多传统 IT 公司的组织架构和分层模式十分的相似，所以它很自然地成为大多数应用的架构模式。

分层架构模式里的组件被分成几个平行的层次，每一层都代表了应用的一个功能（展示逻辑或者业务逻辑）。尽管分层架构没有规定自身要分成几层几种，大多数的结构都分成四个层次，即展示层、业务层、持久层和数据库层。业务层和持久层有时候可以合并成单独的一个业务层，尤其是持久层的逻辑绑定在业务层的组件当中。因此，有一些小的应用可能只有三层，一些有着更复杂的业务的大应用可能有五层甚至更多的层。

分层架构中的每一层都有着特定的角色和职能。举个例子，展示层负责所有的界面展示以及交互逻辑，业务层负责处理请求对应的业务。架构里的层次是具体工作的高度抽象，它们都是为了实现某种特定的业务请求。比如说展示层并不关心如何得到用户数据，它只需在屏幕上以特定的格式展示信息。业务层并不关心要展示在屏幕上的用户数据格式，也不关心这些用户数据从哪里来，它只需要从持久层得到数据，执行与数据有关的相应业务逻辑，然后把这些信息传递给展示层。

分层架构的一个突出特性地组件间关注点分离。一个层中的组件只会处理本层的逻辑。比如说，展示层的组件只会处理展示逻辑，业务层中的组件只会去处理业务逻辑。因为有了组件分离设计方式，让我们更容易构造有效的角色和强力的模型，这样应用变得更好开发、测试、管理和维护。

#### 服务器

##### 高并发服务器 time\_wait

阿里推荐高并发服务器建议调小 TCP 协议的 time\_wait 超时时间。

说明：操作系统默认 240 秒后才会关闭处于 time\_wait 状态的连接，在高并发访问下，服务器端会因为处于 time\_wait 的连接数太多，可能无法建立新的连接，所以需要在服务器上调小此等待值。

正例：在 Linux 服务器上通过变更/etc/sysctl.conf 文件去修改该缺省值（秒）：net.ipv4.tcp\_fin\_timeout=30

##### 我的理解

服务器在处理完客户端的连接后，主动关闭，就会有 time\_wait 状态。TCP 连接是双向的，所以在关闭连接的时候，两个方向各自都需要关闭。先发 FIN 包的一方执行的是主动关闭，后发 FIN 包的一方执行的是被动关闭。主动关闭的一方会进入 time\_wait 状态，并且在此状态停留两倍的 MSL 时长。

主动关闭的一方收到被动关闭的一方发出的 FIN 包后，回应 ACK 包，同时进入 time\_wait 状态，但是因为网络原因，主动关闭的一方发送的这个 ACK 包很可能延迟，从而触发被动连接一方重传 FIN 包。极端情况下，这一去一回就是两倍的 MSL 时长。如果主动关闭的一方跳过 time\_wait 直接进入 closed，或者在 time\_wait 停留的时长不足两倍的 MSL，那么当被动关闭的一方早于先发出的延迟包达到后，就可能出现类似下面的问题：

1. 旧的 TCP 连接已经不存在了，系统此时只能返回 RST 包

2. 新的 TCP 连接被建立起来了，延迟包可能干扰新的连接


不管是哪种情况都会让 TCP 不再可靠，所以 time\_wait 状态有存在的必要性。

修改 net.ipv4.tcp\_fin\_timeout 也就是修改了 MSL 参数。

## 结束语

本文主要介绍了异常日志、安全规约、单元测试、MySQL 数据库、工程结构等五部分关于编码规约的要求。由于篇幅有限，本专题仅提供上、下两篇文章，因此本文（专题下篇）仅覆盖了阿里代码规范的少数内容，更多内容请咨询本文作者。

## 参考资源

- 参考文档《阿里巴巴 Java 开发手册（又名阿里巴巴 Java 代码规约）》。
- 参考书籍 《Effective Java Second Edition》Joshua Bloch。