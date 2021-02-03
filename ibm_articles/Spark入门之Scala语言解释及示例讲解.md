# Spark 入门之 Scala 语言解释及示例讲解
安装、运行 Scala 编译器，快速掌握 Scala 语言

**标签:** Java,分析

[原文链接](https://developer.ibm.com/zh/articles/os-cn-spark-scala/)

周 明耀

发布: 2015-06-23

* * *

Scala 语言衍生自 Funnel 语言。Funnel 语言尝试将函数式编程和 Petri 网结合起来，而 Scala 的预期目标是将面向对象、函数式编程和强大的类型系统结合起来，同时让人要能写出优雅、简洁的代码。本文希望通过一系列 Java 与 Scala 语言编写的相同程序代码的对比，让读者能够尽快地熟悉 Scala 语言。

## 安装 Scala 并调试

首先，我们需要 [从官方网站下载](http://www.scala-lang.org/downloads) 最新的 Scala 运行包，把下载的文件上传到 Linux 服务器并解压，然后进入解压后目录的 bin 目录，进入 Scala 编译器环境，如清单 1 所示。

##### 清单 1\. 进入 Scala 编译器

```
[root@localhost:4 bin]# ./scala
Welcome to Scala version 2.11.6 (OpenJDK 64-Bit Server VM, Java 1.7.0_65).
Type in expressions to have them evaluated.
Type :help for more information.

scala>

```

Show moreShow more icon

清单 1 显示我们使用的是 64 位操作系统，JDK1.7。

在正式讲解 Scala 之前，我们先来简单了解一下它。Scala 是一种解释性语言，可以直接翻译，如清单 2 所示，我们让 1 加上 3，编译器会直接输出整型 4。

##### 清单 2\. 整数相加

```
scala> 1+3
res0: Int = 4

```

Show moreShow more icon

清单 2 输出的 res0 表示变量名，Int 表示类型，输出值是 4，注意 Scala 是强类型语言，必须定义类型，但是 Scala 会帮助您判断数据类型。清单 2 所定义的变量 res0，我们可以直接操作它，如清单 3 所示。

##### 清单 3\. 变量乘法

```
scala> res0*3
res1: Int = 12

```

Show moreShow more icon

清单 3 里面解释器又自动输出一个变量 res1，注意 Scala 的所有变量都是对象，接下来我们在清单 4 所示程序里面把两个变量相乘。

##### 清单 4\. 变量相乘

```
scala> res0*res1
res2: Int = 48

```

Show moreShow more icon

##### 清单 5\. 输出文本

```
scala> println("hello world!")
hello world!

```

Show moreShow more icon

注意，这里由于 println 是 Scala 预定义导入的类，所以可以直接使用，其他非预定义的类，需要手动导入。

如果想要像执行 Shell 文件一样执行 Scala 程序，可以编写.scala 文件，如清单 6 所示。

##### 清单 6\. Scala 程序

```
[root@localhost:4 bin]# cat hello.scala
println("Hello, world, from a script!")
[root@localhost:4 bin]# ./scala hello.scala
Hello, world, from a script!

```

Show moreShow more icon

通过上面简单的介绍，读者应该可以上手写代码了，我们进入 Scala 简介章节。

## Scala 简介

Scala 是一种把面向对象和函数式编程理念加入到静态类型语言中的语言，可以把 Scala 应用在很大范围的编程任务上，无论是小脚本或是大系统都是用 Scala 实现。Scala 运行在标准的 Java 平台上，可以与所有的 Java 库实现无缝交互。可以用来编写脚本把 Java 控件链接在一起。

函数式编程有两种理念做指导，第一种理念是函数是第一类值。在函数式语言中，函数也是值，例如整数和字符串，它们的地位相同。您可以把函数当作参数传递给其他函数，当作结果从函数中返回或保存在变量里。也可以在函数里定义其他函数，就好像在函数里定义整数一样。函数式编程的第二个主要理念是程序的操作符应该把输入值映射到输出值而不是就地修改数据。

Scala 程序会被编译为 JVM 的字节码。它们的执行期性能通常与 Java 程序一致。Scala 代码可以调用 Java 方法，访问 Java 字段，继承自 Java 类和实现 Java 接口。实际上，几乎所有 Scala 代码都极度依赖于 Java 库。

Scala 极度重用了 Java 类型，Scala 的 Int 类型代表了 Java 的原始整数类型 int，Float 代表了 float，Boolean 代表了 boolean，数组被映射到 Java 数组。Scala 同样重用了许多标准 Java 库类型。例如，Scala 里的字符串文本是 Java.lang.String，而抛出的异常必须是 java.lang.Throwable 的子类。

## Scala 编程

Scala 的语法避免了一些束缚 Java 程序的固定写法。例如，Scala 里的分号是可选的，且通常不写。Scala 语法里还有很多其他地方省略了。例如，如何写类及构造函数，清单 7 所示分别采用 Java 和 Scala。

##### 清单 7\. 构造函数写法

```
//Java 代码
class MyClass {
private int index;
private String name;
public MyClass(int index, String name) {
this.index = index; this.name = name;
}
}
//Scala 代码
class MyClass(index: Int, name: String)

```

Show moreShow more icon

根据清单 7 所示代码，Scala 编译器将制造有两个私有字段的类，一个名为 index 的 int 类型和一个叫做 name 的 String 类型，还有一个用这些变量作为参数获得初始值的构造函数。

Scala 可以通过让您提升您设计和使用的接口的抽象级别来帮助您管理复杂性。例如，假设您有一个 String 变量 name，您想弄清楚是否 String 包含一个大写字符。清单 8 所示代码分别采用 Java 和 Scala。

##### 清单 8\. 验证是否大写字符

```
// Java 代码
boolean nameHasUpperCase = false;
for (int i = 0; i < name.length(); ++i) {
if (Character.isUpperCase(name.charAt(i))) {
nameHasUpperCase = true; break;
}
}
//Scala 代码
val nameHasUpperCase = name.exists(_.isUpperCase)

```

Show moreShow more icon

清单 8 所示代码，Java 代码把字符串看作循环中逐字符递进的底层级实体。

Scala 有两种类型的变量，val 和 var。val 变量的值只能初始化一次，再次赋值会发生错误，var 和 Java 的变量相同，可以随时修改。val 是函数式编程的风格，变量一旦赋值就不要再做修改。

清单 9 所示代码定义了变量并操作变量。

##### 清单 9\. Scala 变量操作

```
scala> val message = "hellp world"
message: String = hellp world
scala> val test = "1"
test: String = 1

scala> test ="2"
<console>:8: error: reassignment to val
test ="2"
^
scala> var test1="1"
test1: String = 1

scala> test1="2"
test1: String = 2

```

Show moreShow more icon

清单 9 所示代码定义了变量 message、test，并对 test 重新赋值，由于 val 类型的变量是一次性的，所以抛出错误。var 类型的变量可以重新赋值，并输出新值。字符串支持多行定义，按回车后自动会换行，如清单 10 所示。

##### 清单 10\. Scala 定义多行字符

```
scala> val multiline=
| "try multiple line"
multiline: String = try multiple line

scala> println(multiline)
try multiple line

```

Show moreShow more icon

在 Scala 里，定义方法采用 def 标示符，示例代码如清单 11 所示。

##### 清单 11\. 定义方法

```
scala> def max(x: Int, y: Int): Int = if(x < y) y else x
max: (x: Int, y: Int)Int
scala> max(3,8)
res0: Int = 8

```

Show moreShow more icon

清单 11 所示代码定义了方法 Max，用于比较传入的两个参数的大小，输出较大值。

函数的定义用 def 开始。每个函数参数后面必须带前缀冒号的类型标注，因为 Scala 编译器没办法推断函数参数类型。清单 11 所定义的函数如图 1 所示，分别对函数体内的每一个元素列出了用途。

##### 图 1\. 函数定义解释

![图 1. 函数定义解释](../ibm_articles_img/os-cn-spark-scala_images_img001.jpg)

如清单 12 所示的 Java 和 Scala 代码，我们定义了一个函数 greet，调用该函数会打印出 “Good Moring！”。

##### 清单 12\. 定义函数

```
//Java 代码
public static void main(String[] args){
JavaScala.greet();
}

public static void greet(){
System.out.println("Good Morning!");
}
//Scala 代码
scala> def greet()=println("Good Morning!")
greet: ()Unit
scala> greet();
Good Morning!

```

Show moreShow more icon

上例定义了 greet() 函数，编译器回应 greet 是函数名，空白的括号说明函数不带参数，Unit 是 greet 的结果类型。Unit 的结果类型指的是函数没有返回有用的值。Scala 的 Unit 类型接近于 Java 的 void 类型，而且实际上 Java 里每一个返回 void 的方法都被映射为 Scala 里返回 Unit 的方法。

注意，离开 Scala 编译器可以用:quit 或:q 命令。

与 Java 一样, 可以通过 Scala 的名为 args 的数组获得传递给 Scala 脚本的命令行参数。Scala 里，数组以零开始，通过在括号里指定索引访问一个元素。所以 Scala 里数组 steps 的第一个元素是 steps(0)，而不是 Java 里的 steps[0]。清单 13 所示代码编写了一个 Scala 文件，定义读入第一个参数。

##### 清单 13\. 定义 Main 函数参数

```
[root@localhost:4 bin]# ./scala hello.scala zhoumingyao
Hello, world, from a script!zhoumingyao
[root@localhost:4 bin]# cat hello.scala
println("Hello, world, from a script!"+args(0))
//Java 代码
System.out.println("Hello, world, from a script!"+args(0));

```

Show moreShow more icon

当我们需要执行循环的时候，While 是一个不错的选择。清单 14 所示是 While 的实现。

##### 清单 14\. While 循环

```
[root@localhost:4 bin]# cat hello.scala
var i = 0;
while(i < args.length){
println(args(i))
i += 1
}
[root@localhost:4 bin]# ./scala hello.scala hello world ! this is zhoumingyao
hello
world
!
this
is
zhoumingyao

```

Show moreShow more icon

上面的 While 循环读取输入的参数，直到参数读取完毕。

Scala 里可以使用 new 实例化对象或类实例，通过把加在括号里的对象传递给实例的构造器的方式来用值参数化实例。例如，清单 15 所示代码的 Scala 程序实例化 java.math.BigInteger，实例化字符串数组。

##### 清单 15\. 实例化参数

```
val big = new java.math.BigInteger("12345")
val greetStrings = new Array[String](3)
greetStrings(0) = "Hello"
greetStrings(1) = ", "
greetStrings(2) = "world!\n"
for(i <- 0 to 2)
print(greetStrings(i))
[root@localhost:4 bin]# ./scala hello.scala
Hello, world!

```

Show moreShow more icon

从技术上讲，Scala 没有操作符重载，因为它根本没有传统意义上的操作符。取而代之的是，诸如+，-，\*和/这样的字符可以用来做方法名。

Scala 的 List 是不可变对象序列。List[String] 包含的仅仅是 String。Scala 的 List、Scala.List。不同于 Java 的 java.util.List，总是不可变的，在 Java 中是可变的。

val oneTwoThree = List(1, 2, 3)

List “：：：” 的方法实现叠加功能，程序如清单 16 所示。

##### 清单 16\. 叠加方式代码

```
[root@localhost:4 bin]# cat hello.scala
val oneList = List(1,2)
val twoList = List(3,4)
val combinedList = oneList ::: twoList
println(oneList + " and " + twoList +" and " + combinedList)
[root@localhost:4 bin]# ./scala hello.scala
List(1, 2) and List(3, 4) and List(1, 2, 3, 4)

```

Show moreShow more icon

初始化新 List 的方法是把所有元素用 cons 操作符串联起来，Nil 作为最后一个元素。

##### 清单 17\. Nil 方式

```
val oneTwoThree = 1 :: 2 :: 3 :: Nil
println(oneTwoThree)
[root@localhost:4 bin]# ./scala hello*
List(1, 2, 3)

```

Show moreShow more icon

##### 表格 1\. 方法列表

方法名方法作用List() 或 Nil空 ListList(“Cool”, “tools”, “rule)创建带有三个值”Cool”，”tools”和”rule”的新 List[String]val thrill = “Will”::”fill”::”until”::Nil创建带有三个值”Will”，”fill”和”until”的新 List[String]List(“a”, “b”) ::: List(“c”, “d”)叠加两个列表（返回带”a”，”b”，”c”和”d”的新 List[String]）thrill(2)返回在 thrill 列表上索引为 2（基于 0）的元素（返回”until”）thrill.count(s => s.length == 4)计算长度为 4 的 String 元素个数（返回 2）thrill.drop(2)返回去掉前 2 个元素的 thrill 列表（返回 List(“until”)）thrill.dropRight(2)返回去掉后 2 个元素的 thrill 列表（返回 List(“Will”)）thrill.exists(s => s == “until”)判断是否有值为”until”的字串元素在 thrill 里（返回 true）thrill.filter(s => s.length == 4)依次返回所有长度为 4 的元素组成的列表（返回 List(“Will”, “fill”)）thrill.forall(s => s.endsWith(“1”))辨别是否 thrill 列表里所有元素都以”l”结尾（返回 true）thrill.foreach(s => print(s))对 thrill 列表每个字串执行 print 语句（”Willfilluntil”）thrill.foreach(print)与前相同，不过更简洁（同上）thrill.head返回 thrill 列表的第一个元素（返回”Will”）thrill.init返回 thrill 列表除最后一个以外其他元素组成的列表（返回 List(“Will”, “fill”)）thrill.isEmpty说明 thrill 列表是否为空（返回 false）thrill.last返回 thrill 列表的最后一个元素（返回”until”）thrill.length返回 thrill 列表的元素数量（返回 3）thrill.map(s => s + “y”)返回由 thrill 列表里每一个 String 元素都加了”y”构成的列表（返回 List(“Willy”, “filly”, “untily”)）thrill.mkString(“, “)用列表的元素创建字串（返回”will, fill, until”）thrill.remove(s => s.length == 4)返回去除了 thrill 列表中长度为 4 的元素后依次排列的元素列表（返回 List(“until”)）thrill.reverse返回含有 thrill 列表的逆序元素的列表（返回 List(“until”, “fill”, “Will”)）thrill.sort((s, t) => s.charAt(0).toLowerCase < t.charAt(0).toLowerCase)返回包括 thrill 列表所有元素，并且第一个字符小写按照字母顺序排列的列表（返回 List(“fill”, “until”, “Will”)）thrill.tail返回除掉第一个元素的 thrill 列表（返回 List(“fill”, “until”)）

另一种容器对象是元组 (tuple)，与列表一样，元组是不可变得。但与列表不同，元组可以包含不同类型的元素。元组的用处，如果您需要在方法里返回多个对象。实例化一个装有一些对象的新元组，只要把这些对象放在括号里，并用逗号分隔即可。一旦实例化一个元组，可以用点号、下划线和一个基于 1 的元素索引访问它。如清单 18 所示。

##### 清单 18\. 元祖代码

```
val pair = (99, "Luftballons","whawo")
println(pair._1)
println(pair._2)
println(pair._3)
[root@localhost:4 bin]# ./scala hello.scala
99
Luftballons
Whawo

```

Show moreShow more icon

##### 清单 19\. Set 操作代码

```
import scala.collection.mutable.Set
val movieSet = Set("Hitch", "Poltergeist")
movieSet += "Shrek"
println(movieSet)
import scala.collection.immutable.HashSet
val hashSet = HashSet("Tomatoes", "Chilies")
println(hashSet + "Coriander")
[root@localhost:4 bin]# ./scala hello.scala
Set(Poltergeist, Shrek, Hitch)
Set(Chilies, Tomatoes, Coriander)

```

Show moreShow more icon

Map 是 Scala 里另一种有用的集合类。Map 的类继承机制看上去和 Set 的很像。scala.collection 包里面有一个基础 Map 特质和两个子特质 Map：可变的 Map 在 scala.collection.mutable 里，不可变的在 scala.collection.immutable 里。

##### 清单 19\. Map 操作代码

```
import scala.collection.mutable.Map
val treasureMap = Map[Int, String]()
treasureMap += (1 -> "Go to island.")
treasureMap += (2 -> "Find big X on ground.")
treasureMap += (3 -> "Dig.")
println(treasureMap(2))
输出：
Find big X on ground.

```

Show moreShow more icon

如果我们尝试从文件按行读取内容，代码可以如清单 20 所示。

##### 清单 20\. 读取文件

```
[root@localhost:2 bin]# cat hello.scala
import scala.io.Source
if (args.length > 0) {
for (line <- Source.fromFile(args(0)).getLines)
println(line.length + " " + line)
} else
Console.err.println("Please enter filename")
运行命令: [root@localhost:2 bin]# ./scala hello.scala hello.scala
输出结果:
23 import scala.io.Source
23 if (args.length > 0) {
52 for (line <- Source.fromFile(args(0)).getLines)
42 println(line.length + " " + line)
7 } else
48 Console.err.println("Please enter filename")

```

Show moreShow more icon

清单 20 所示脚本从 scala.io 引入类 Source，然后检查是否命令行里定义了至少一个参数。表达式 Source.fromFile(args(0))，尝试打开指定的文件并返回一个 Source 对象。函数返回 Iterator[String]，在每个枚举里提供一行包含行结束符的信息。

类是对象的蓝图我们在很多场合都需要使用类，举例来说，您定义了 ChecksumAccumulator 类并给它一个叫做 sum 的 var 字段，然后再实例化两次。代码如清单 21 所示。

##### 清单 21\. 类的定义和使用

```
class ChecksumAccumulator { var sum = 0 }
val acc = new ChecksumAccumulator val csa = new ChecksumAccumulator

```

Show moreShow more icon

注意，Public 是 Scala 的缺省访问级别，Scala 比 Java 更面向对象的一个方面是 Scala 没有静态成员。替代品是 Scala 有单例对象 Singleton object。除了用 object 关键字替换了 class 关键字以外，单例对象的定义看上去就像是类定义，清单 22 是单例对象的定义方法，Java 类似代码请见已发布的本文作者编写的《设计模式第一部分：单例模式》一文。

##### 清单 22\. 单例方法

```
import scala.collection.mutable.Map
object ChecksumAccumulator {
private val cache = Map[String, Int]()
def calculate(s: String): Int =
if (cache.contains(s))//检查缓存是否存在
cache(s)//缓存存在，直接返回映射里面的值
else {//缓存不存在
val acc = new ChecksumAccumulator
for (c <- s)//对传入字符串的每个字符循环一次
acc.add(c.toByte)//字符转换成 Byte
val cs = acc.checksum()
cache += (s -> cs)
cs
}
}

```

Show moreShow more icon

上面的单例对象有一个 calculate 方法，用来计算所带的 String 参数中字符的校验和。它还有一个私有字段 Cache，一个缓存之前计算过的校验和的可变映射。这里我们使用了缓存例子来说明带有域的单例对象。像这样的缓存是通过内存换计算时间的方式做到性能的优化。通常意义上说，只有遇到了缓存能解决的性能问题时，才可能用到这样的例子，而且应该使用弱映射（weak map），如 scala.Collection.jcl 的 WeakHashMap，这样如果内存稀缺的话，缓存里的条目就会被垃圾回收机制回收掉。

要执行 Scala 程序，一定要提供一个有 main 方法的孤立单例对象名，main 方法带有一个 Array[String] 的参数，结果类型为 Unit。任何拥有合适签名的 main 方法的单例对象都可以用来作为程序的入口点。

##### 清单 23\. main 函数

```
import ChecksumAccumulator.calculate
object Summer {
def main(args: Array[String]) {
for (arg <- args)
println(arg + ": " + calculate(arg))
}
}

```

Show moreShow more icon

Scala 隐式引用了包 java.lang、scala 的成员、Predef 的单例对象。Predef 包括 println 和 assert 等等。清单 21 和 23 所示代码里，无论 ChecksumAccumulator.scala 还是 Summer.scala 都不是脚本，因为他们是以定义结束的。反过来说，脚本必然以一个结果表达式结束。因此如果您尝试以脚本方式执行 Summer.scala，Scala 解释器将会报错说 Summer.scala 不是以结果表达式结束的（当然前提是您没有在 Summer 对象定义之后加上任何您自己的表达式）。正确的做法是，您需要用 Scala 编译器真正地编译这些文件，然后执行输出的类文件。其中一种方式是使用 scalac，Scala 的基本编译器。输入$ scalac ChecksumAccumulator.scala Summer.scala 命令会编译您的源代码，每次编译器启动时，都要花一些时间扫描 jar 文件内容，并在即使您提交的是新的源文件也需要查看之前完成其他初始化工作。

因此，Scala 的发布包里还包括了一个叫做 fsc（快速 Scala 编译器）的 Scala 编译器后台服务：daemon。您可以这样使用： $ fsc ChecksumAccumulator.scala Summer.scala 第一次执行 fsc 时，会创建一个绑定在您计算机端口上的本地服务器后台进程。然后它就会把文件列表通过端口发送给后台进程去编译，后台进程完成编译。下一次您执行 fsc 时，后台进程就已经在运行了，于是 fsc 将只是把文件列表发给后台进程，它会立刻开始编译文件。使用 fsc，您只需要在第一次等待 Java 运行时环境的启动。如果想停止 fsc 后台进程，可以执行 fsc -shutdown 来关闭。

不论执行 scalac 还是 fsc 命令，都将创建 Java 类文件，然后您可以用 Scala 命令，就像之前的例子里调用解释器那样运行它。不过，不是像前面每个例子里那样把包含了 Scala 代码的带有.scala 扩展名的文件交给它解释执行，而是采用这样的方式，$ scala Summer of love。

本文对 Scala 语言的基础做了一些解释，由于篇幅所限，所以下一篇文章里会针对 Spark 附带的示例代码、Spark 源代码中出现的 Scala 代码进行解释。

## 结束语

通过本文的学习，读者了解了如何下载、部署 Scala。此外，通过编写 Scala 与 Java 相同功能的程序，让 Java 程序员可以快速掌握 Scala 语言，为后面的 Spark 源代码分析文章做知识准备。目前市面上发布的 Spark 中文书籍对于初学者来说大多较为难读懂，作者力求推出一系列 Spark 文章，让读者能够从实际入手的角度来了解 Spark。后续除了应用之外的文章，还会致力于基于 Spark 的系统架构、源代码解释等方面的文章发布。