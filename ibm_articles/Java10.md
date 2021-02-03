# Java 10
局部变量类型推断是一个充满争议的热门话题，但 Java 10 为 JVM 中的垃圾收集和容器感知带来了一些受欢迎的改变。

**标签:** Java,Java 平台,OpenJDK

[原文链接](https://developer.ibm.com/zh/articles/j-5things17/)

Alex Theedom

发布: 2018-05-16

* * *

**关于本系列**

您觉得自己了解 Java 编程？事实是，大多数开发人员只是了解 Java 平台的皮毛，所学知识也仅够应付工作。在这个 [连载系列](/zh/series/5-things-you-didnt-know-about/) 中，Java 技术侦探们将深度挖掘 Java 平台的核心功能，揭示一些可帮助您解决最棘手的编程挑战的技巧和诀窍。

Java™ 开发人员习惯于为某个新 Java 版本等待数年，但是新的高频率版本发布节奏改变了这一情况。在 Java 9 出现短短 6 个月后，Java 10 就面世了。再等 6 个月，我们将迎来 Java 11。一些开发人员可能认为这种发布节奏过快，但这一新节奏标志着一种期待已久的改变。

跟它的版本号一样，Java 10 提供了 10 个新特性，本文将介绍我认为最重要的 5 个新特性（可以在 [Open JDK 的 10 个项目](http://openjdk.java.net/projects/jdk/10/) 页面上查看所有新特性）。

[获取代码](https://github.com/readlearncode/Java-10-Sampler)

## 1.Java 的新版本发布节奏

过去，JDK 版本发布节奏受重大新特性所驱动。以最近为例，Java 8 以 lambda 表达式和流的形式引入了函数式编程，Java 9 引入了模块化 Java 系统。每个新版本都被热切期待，而一些小的修复补丁则常常被搁在一边，等待更大组件的最终确定。Java 的演变落后于其他语言。

新的高频率节奏促使 Java 在有较小的改进时就发布新版本。在发布日之前将已准备就绪的特性包含在新版本内，未准备就绪的特性可以计划包含在 6 个月后的下一个版本中。这个新周期中的第一个 Java 版本是 Java 9，它于 2017 年 10 月推出。Java 10 于 2018 年 3 月发布，Java 11 预计将在 2018 年 9 月发布。

作为新节奏的一部分，Oracle 表示它对每个主要版本的支持仅持续到下一个主要版本推出之前。当 Java 11 发布时，Oracle 将停止支持 Java 10。如果开发人员想要确保其 Java 版本受到支持，则需要每隔 6 个月就迁移到一个主要版本。不想要或不需要如此频繁地进行迁移的开发人员，可以使用 LTS（长期支持）版本，该版本每 3 年更新一次。对当前的 LTS 版本 (Java 8) 的支持将截止于今年秋天推出 Java 11 时。

## 2.局部变量类型推断

目前为止，局部变量类型推断是 Java 10 中最受人瞩目的特性。在经过激烈的争议之后，此特性才被引入 JDK 10 中，它允许编译器推断局部变量的类型，而不是要求程序员明确指定变量类型。

清单 1 展示了在 Java 10 之前是如何定义 `String` 变量类型的。

##### 清单 1\. 声明并分配一个 String 变量

```
String name = "Alex";

```

Show moreShow more icon

清单 2 给出了 Java 10 中定义的相同 `String` 变量。

##### 清单 2\. 一个使用局部变量类型推断功能定义的 String 变量

```
var name = "Alex";

```

Show moreShow more icon

可以看到，唯一的区别在于使用了保留类型名称 `var` 。使用右侧的表达式，编译器可以推断出变量 `name` 的类型为 `String` 。

这看起来似乎很简单，那么让我们来看一个更加复杂的示例。如果变量被分配给某个方法调用的返回值，结果会怎样呢？在这种情况下，编译器可以从方法的返回类型来推断变量类型，如清单 3 所示。

##### 清单 3\. 一个从返回类型推断出的 String 变量

```
var name = getName(); String getName(){ return "Alex"; }

```

Show moreShow more icon

### 使用局部变量类型

顾名思义，局部变量类型推断特性只能用于局部变量。不能将它用于定义实例或类变量，也不能在方法参数或返回类型中使用它。但是，在可以从迭代器推断类型的经典和增强的 for 循环中， _可以_ 使用 `var` ，如清单 4 所示。

##### 清单 4\. 在循环中使用 var

```
for(var book : books){} for(var i = 0; i < 10; i++){}

```

Show moreShow more icon

使用此类型的最明显理由是减少代码中的冗余。看看清单 5 中的示例。

##### 清单 5\. 长类型名称使得代码变长

```
String message = "Incy wincy spider..."; StringReader reader = new StringReader(message); StreamTokenizer tokenizer = new StreamTokenizer(reader);

```

Show moreShow more icon

请注意我们使用保留类型名称 `var` 重写清单 5 后发生的情况。

##### 清单 6\. var 类型减少了代码冗余

```
var message = "Incy wincy spider..."; var reader = new StringReader(message); var tokenizer = new StreamTokenizer(reader);

```

Show moreShow more icon

清单 6 中的类型声明是垂直排列的，而且该类型仅在每条语句中提过一次，它在构造函数调用的右侧。想象对 Java 框架中常见的长类名使用此类型的一些好处。

### 局部变量类型的问题

#### 1\. var 让类型变得模糊

您已了解到 `var` 可以改进代码可读性，但是另一方面，它也有可能让类型变得模糊。看看清单 7 中的示例。

##### 清单 7\. 返回类型不清楚

```
var result = searchService.retrieveTopResult();

```

Show moreShow more icon

在清单 7 中，我们必须猜测返回类型。让读者猜测正发生的事件的代码更加难以维护。

#### 2\. var 无法与 lambda 结合使用

在用于 lambda 表达式时，类型推断的效果不是很好，主要由于缺乏可供编译器使用的类型信息。清单 8 中的 lambda 表达式不会执行编译。

##### 清单 8\. 类型信息不足

```
Function<String, String> quotify = m -> "'" + message + "'"; var quotify = m -> "'" + message + "'";

```

Show moreShow more icon

在清单 8 中，右侧表达式中没有足够的类型信息供编译器推断变量类型。Lambda 语句必须始终声明一个显式类型。

#### 3\. var 无法与菱形运算符结合使用

在用于菱形运算符时，类型推断的效果也不好。看看清单 9 中的示例。

##### 清单 9\. 结合使用菱形运算符与 var

```
var books = new ArrayList<>();

```

Show moreShow more icon

**自行尝试**

要自行尝试局部变量类型推断，您需要 [下载 JDK 10](http://jdk.java.net/10/) 和一个支持它的 IDE。IntelliJ 的 EAP ( [Early Access Program](https://www.jetbrains.com/idea/nextversion/)) 版本提供了这项支持。下载并安装这些工具后，您可以首先检查与本文配套的 GitHub 存储库。您可以在这里找到局部变量类型推断的示例。

在清单 9 中， `books` 引用的 `ArrayList` 的参数类型是什么？您可能知道，您希望 `ArrayList` 存储一个图书列表，但编译器无法推断出这一点。相反，编译器只会执行它 _能_ 执行的操作，即推断一个由 `Object` 类型参数化的 `ArrayList` ： `ArrayList<Object>()` 。

替代方案是，在右侧表达式的菱形运算符中指定类型。然后，可以让编译器根据该信息推断变量类型，如清单 10 所示。否则，您必须采用传统方式来显式声明该变量： `List<Book> books` 。实际上，您可能更喜欢此选项，因为它允许您指定一个抽象类型并编程到 `List` 接口中：

##### 清单 10\. 指定类型

```
var books = new ArrayList<Book>();

```

Show moreShow more icon

## 3.增加、删除和弃用的特性

### 删除

Java 10 删除了一些工具：

- 命令行工具 `javah` ，但您可以使用 `javac -h` 代替它。
- 命令行选项 `-X:prof` ，但您可以使用 `jmap` 工具来访问探查信息。
- `policytool` 。

一些由于 Java 1.2 被永久删除而被标记为弃用的 API。这些 API 包括 `java.lang.SecurityManager.inCheck` 字段和以下方法：

- `java.lang.SecurityManager.classDepth(java.lang.String)`
- `java.lang.SecurityManager.classLoaderDepth()`
- `java.lang.SecurityManager.currentClassLoader()`
- `java.lang.SecurityManager.currentLoadedClass()`
- `java.lang.SecurityManager.getInCheck()`
- `java.lang.SecurityManager.inClass(java.lang.String)`
- `java.lang.SecurityManager.inClassLoader()`
- `java.lang.Runtime.getLocalizedInputStream(java.io.InputStream)`
- `java.lang.Runtime.getLocalizedOutputStream(java.io.OutputStream)`

### 弃用

JDK 10 也弃用了一些 API。 `java.security.acl` 包被标记为弃用， `java.security` 包中的各种相关类（ `Certificate` 、 `Identity` 、 `IdentityScope` 、 `Singer` 、 `auth.Policy` ）也是如此。此外， `javax.management.remote.rmi.RMIConnectorServer` 类中的 `CREDENTIAL_TYPES` 也被标记为弃用。 `java.io.FileInputStream` 和 `java.io.FileOutputStream` 中的 `finalize()` 方法已被标记为弃用。 `java.util.zip.Deflater` / `Inflater` / `ZipFile` 类中的 `finalize()` 方法也是如此。

### 增加和包含的特性

作为 Oracle JDK 与 Open JDK 持续合作的一部分，Open JDK 现在包含 Oracle JDK 中提供的根证书颁发机构的一部分。这些颁发机构包括 Java Flight Recorder 和 Java Mission Control。此外，在 `java.text` 、 `java.time` 和 `java.util` 包中的合适位置，JDK 10 还添加了对 BCP 的 47 种语言标记的 Unicode 扩展的增强支持。另一个新特性允许在线程上执行回调，而不执行全局 VM 安全点。这使得终止各个线程变得可行且经济，而不需要终止所有线程或全部不终止。

## 4.改进的容器感知

如果您部署到像 Docker 这样的容器中，那么此特性特别适合您。JVM 现在能感知到它在一个容器中运行，并向该容器查询可供使用的处理器数量，而不是向主机操作系统查询。也可以从外部连接到在容器中运行的 Java 进程，这使得监控 JVM 进程变得更容易。

以前，JVM 对它的容器一无所知，并会向主机操作系统询问活动 CPU 的数量。在某些情况下，这会导致向 JVM 报告过多的资源，致使多个容器在同一个操作系统上运行时出现问题。在 Java 10 中，您可以配置一个容器来使用主机操作系统的部分 CPU，JVM 能够确定被使用的 CPU 数量。也可以使用 `-XX:ActiveProcessorCount` 标志，显式指定容器化的 JVM 能看到的处理器数量。

## 5.应用程序类数据共享

此特性的目的是改善各运行之间的 JVM 启动时间，让多个 JVM 能够运行相同的代码，同时还能减少内存占用。这是通过在 JVM 之间共享关于类的元数据来实现的。JVM 的第一次运行会收集并记录有关它加载的类的数据。然后，它会让数据文件可用于其他 JVM 和该 JVM 的后续运行，并节省 JVM 初始化流程中的时间和资源。类数据共享实际上已存在了一段时间，但仅限于系统类。现在此特性已扩展，以包含所有的应用程序类。

## 结束语

Java 10 中最热门的特性显然是新的保留类型名称 `var` 。它能够阐明并简化代码，但是，如果不小心使用，它可能让含义和意图变得模糊。IDE 可以在类型不清楚时帮助识别类型，但不会将所有代码都读入到 IDE 中。我们通常在 GitHub 存储库、调试器或代码评审工具中在线读取代码。对于使用这个新特性的开发人员，建议考虑代码对于未来读者和维护者的可读性。

Java 新的高频率版本发布节奏是一个受欢迎的改变。它强制要求发布已在发布日期前准备好的特性，为延迟的特性留出较短的周转时间，延迟到下一个版本进行发布。新周期会加快 Java 的进度，开发人员不需要为已开发并搁在一边的特性而等待数年。对于从一个主要版本到下一个主要版本的较短周期的支持，人们给与了合理的关注，但 LTS 应该有助于缓解这一问题。另一个风险是版本疲劳，因为开发人员厌倦了不断迁移。尽管如此，总的来说，我认为这是一个积极的举措，它将帮助 Java 在未来一段时间内保持活力和发展。

本文翻译自： [Java 10](https://developer.ibm.com/articles/j-5things17/)（2018-04-17）