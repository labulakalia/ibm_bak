# Java 8 的 Lambda 表达式和流处理
从 JSR 335 出发对 Lambda 表达式进行了深入的介绍

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-understanding-functional-programming-3/)

成 富

发布: 2018-12-03

* * *

在本系列的前两篇文章中，已经对 [函数式编程的思想](https://developer.ibm.com/zh/articles/j-understanding-functional-programming-1/) 和 [函数式编程的重要概念](https://developer.ibm.com/zh/articles/j-understanding-functional-programming-2/) 做了介绍。本文将介绍 Java 平台本身对函数式编程的支持，着重介绍 Lambda 表达式和流（Stream）。

## Lambda 表达式

当提到 Java 8 的时候，Lambda 表达式总是第一个提到的新特性。Lambda 表达式把函数式编程风格引入到了 Java 平台上，可以极大的提高 Java 开发人员的效率。这也是 Java 社区期待已久的功能，已经有很多的文章和图书讨论过 Lambda 表达式。本文则是基于官方的 JSR 335（Lambda Expressions for the Java Programming Language）来从另外一个角度介绍 Lambda 表达式。

### 引入 Lambda 表达式的动机

我们先从清单 1 中的代码开始谈起。该示例的功能非常简单，只是启动一个线程并输出文本到控制台。虽然该 Java 程序一共有 9 行代码，但真正有价值的只有其中的第 5 行。剩下的代码全部都是为了满足语法要求而必须添加的冗余代码。代码中的第 3 到第 7 行，使用 `java.lang.Runnable` 接口的实现创建了一个新的 `java.lang.Thread` 对象，并调用 `Thread` 对象的 `start` 方法来启动它。`Runnable` 接口是通过一个匿名内部类实现的。

##### 清单 1\. 传统的启动线程的方式

```
public class OldThread {
public static void main(String[] args) {
new Thread(new Runnable() {
     public void run() {
       System.out.println("Hello World!");
     }
}).start();
}
}

```

Show moreShow more icon

从简化代码的角度出发，第 3 行和第 7 行的 `new Runnable()` 可以被删除，因为接口类型 `Runnable` 可以从类 `Thread` 的构造方法中推断出来。第 4 和第 6 行同样可以被删除，因为方法 `run` 是接口 `Runnable` 中的唯一方法。把第 5 行代码作为 `run` 方法的实现不会出现歧义。把第 3，4，6 和 7 行的代码删除掉之后，就得到了使用 Lambda 表达式的实现方式，如清单 2 所示。只用一行代码就完成了清单 1 中 5 行代码完成的工作。这是令人兴奋的变化。更少的代码意味着更高的开发效率和更低的维护成本。这也是 Lambda 表达式深受欢迎的原因。

##### 清单 2\. 使用 Lambda 表 达式启动线程

```
public class LambdaThread {
public static void main(String[] args) {
    new Thread(() -> System.out.println("Hello World!")).start();
}
}

```

Show moreShow more icon

简单来说，Lambda 表达式是创建匿名内部类的语法糖（syntax sugar）。在编译器的帮助下，可以让开发人员用更少的代码来完成工作。

### 函数式接口

在对 [清单 1](#清单-1-传统的启动线程的方式) 的代码进行简化时，我们定义了两个前提条件。第一个前提是要求接口类型，如示例中的 `Runnable`，可以从当前上下文中推断出来；第二个前提是要求接口中只有一个抽象方法。如果一个接口仅有一个抽象方法（除了来自 Object 的方法之外），它被称为函数式接口（functional interface）。函数式接口的特别之处在于其实例可以通过 Lambda 表达式或方法引用来创建。Java 8 的 `java.util.function` 包中添加了很多新的函数式接口。如果一个接口被设计为函数式接口，应该添加 `@FunctionalInterface` 注解。编译器会确保该接口确实是函数式接口。当尝试往该接口中添加新的方法时，编译器会报错。

### 目标类型

Lambda 表达式没有类型信息。一个 Lambda 表达式的类型由编译器根据其上下文环境在编译时刻推断得来。举例来说，Lambda 表达式 `() -> System.out.println("Hello World!")` 可以出现在任何要求一个函数式接口实例的上下文中，只要该函数式接口的唯一方法不接受任何参数，并且返回值是 `void`。这可能是 `Runnable` 接口，也可能是来自第三方库或应用代码的其他函数式接口。由上下文环境所确定的类型称为目标类型。Lambda 表达式在不同的上下文环境中可以有不同的类型。类似 Lambda 表达式这样，类型由目标类型确定的表达式称为多态表达式（poly expression）。

Lambda 表达式的语法很灵活。它们的声明方式类似 Java 中的方法，有形式参数列表和主体。参数的类型是可选的。在不指定类型时，由编译器通过上下文环境来推断。Lambda 表达式的主体可以返回值或 `void`。返回值的类型必须与目标类型相匹配。当 Lambda 表达式的主体抛出异常时，异常的类型必须与目标类型的 `throws` 声明相匹配。

由于 Lambda 表达式的类型由目标类型确定，在可能出现歧义的情况下，可能有多个类型满足要求，编译器无法独自完成类型推断。这个时候需要对代码进行改写，以帮助编译器完成类型推断。一个常见的做法是显式地把 Lambda 表达式赋值给一个类型确定的变量。另外一种做法是显示的指定类型。

在清单 3 中，函数式接口 A 和 B 分别有方法 `a` 和 `b`。两个方法 `a` 和 `b` 的类型是相同的。类 `UseAB` 的 `use` 方法有两个重载形式，分别接受类 A 和 B 的对象作为参数。在方法 `targetType` 中，如果直接使用 `() -> System.out.println("Use")` 来调用 `use` 方法，会出现编译错误。这是因为编译器无法推断该 Lambda 表达式的类型，类型可能是 A 或 B。这里通过显式的赋值操作为 Lambda 表达式指定了类型 A，从而可以编译通过。

##### 清单 3\. 可能出现歧义的目标类型

```
public class LambdaTargetType {

@FunctionalInterface
interface A {
    void a();
}

@FunctionalInterface
interface B {
    void b();
}

class UseAB {
    void use(A a) {
      System.out.println("Use A");
    }

    void use(B b) {
      System.out.println("Use B");
    }
}

void targetType() {
    UseAB useAB = new UseAB();
    A a = () -> System.out.println("Use");
    useAB.use(a);
}
}

```

Show moreShow more icon

### 名称解析

在 Lambda 表达式的主体中，经常需要引用来自包围它的上下文环境中的变量。Lambda 表达式使用一个简单的策略来处理主体中的名称解析问题。Lambda 表达式并没有引入新的命名域（scope）。Lambda 表达式中的名称与其所在上下文环境在同一个词法域中。Lambda 表达式在执行时，就相当于是在包围它的代码中。在 Lambda 表达式中的 this 也与包围它的代码中的含义相同。在清单 4 中，Lambda 表达式的主体中引用了来自包围它的上下文环境中的变量 `name`。

##### 清单 4\. Lambda 表 达式中的名称解析

```
public void run() {
String name = "Alex";
new Thread(() -> System.out.println("Hello, " + name)).start();
}

```

Show moreShow more icon

需要注意的是，可以在 Lambda 表达式中引用的变量必须是声明为 `final` 或是实际上 `final（effectively final）`的。实际上 `final` 的意思是变量虽然没有声明为 `final`，但是在初始化之后没有被赋值。因此变量的值没有改变。

## 流

Java 8 中的流表示的是元素的序列。流中的元素可能是对象、int、long 或 double 类型。流作为一个高层次的抽象，并不关注流中元素的来源或是管理方式。流只关注对流中元素所进行的操作。当流与函数式接口和 Lambda 表达式一同使用时，可以写出简洁高效的数据处理代码。下面介绍几个与流相关的基本概念。

### 顺序执行和 并行执行

流的操作可以顺序执行或并行执行, 后者可以获得比前者更好的性能。但是如果实现不当，可能由于数据竞争或无用的线程同步，导致并行执行时的性能更差。一个流是否会并行执行，可以通过其方法 `isParallel()` 来判断。根据流的创建方式，一个流有其默认的执行方式。可以使用方法 `sequential()` 或 `parallel()` 来将其执行方式设置为顺序或并行。

### 相遇顺序

一个流的相遇顺序（encounter order）是流中的元素被处理时的顺序。流根据其特征可能有，也可能没有一个确定的相遇顺序。举例来说，从 ArrayList 创建的流有确定的相遇顺序；从 HashSet 创建的流没有确定的相遇顺序。大部分的流操作会按照流的相遇顺序来依次处理元素。如果一个流是无序的，同一个流处理流水线在多次执行时可能产生不一样的结果。比如 Stream 的 `findFirst()` 方法获取到流中的第一个元素。如果在从 ArrayList 创建的流上应用该操作，返回的总是第一个元素；如果是从 HashSet 创建的流，则返回的结果是不确定的。对于一个无序的流，可以使用 `sorted` 操作来排序；对于一个有序的流，可以使用 `unordered()` 方法来使其无序。

### Spliterator

所有的流都是从 Spliterator 创建出来的。Spliterator 的名称来源于它所支持的两种操作：split 和 iterator。Spliterator 可以看成是 Iterator 的并行版本，允许通过对流中元素分片的方式来切分数据源。使用其 `tryAdvance` 方法来顺序遍历元素，也可以使用 `trySplit` 方法来创建一个新的 Spliterator 对象在新划分的数据集上工作。Spliterator 还提供了 `forEachRemaining` 方法进行批量顺序遍历。可以使用 `estimateSize` 方法来查询可能会遍历的元素数量。一般的做法是先使用 trySplit 切分数据源。当元素数量足够小时，使用 `forEachRemaining` 来对分片中的全部元素进行处理。这也是典型的分治法的思路。

每个 Spliterator 可以有一系列不同的特征，可以通过 `characteristics` 方法来查询。一个 Spliterator 具备的特征取决于其数据源和元素。所有可用的特征如下所示：

- CONCURRENT：表明数据源可以安全地由多个线程进行修改，而无需额外的同步机制。
- DISTINCT：表明数据源中的元素是唯一的，不存在重复元素。
- IMMUTABLE：表明数据源是不可变的， 无法进行修改操作。
- NONNULL：表明数据源中不存在 null 元素。
- ORDERED：表明数据源中的元素有确定的相遇顺序。
- SIZED：表明数据源中的元素的数量是确定的。
- SORTED：表明数据源中的元素是有序的。
- SUBSIZED：表明使用 trySplit 切分出来的子数据源也有 SIZED 和 SUBSIZED 的特征。

Spliterator 需要绑定到流之后才能遍历其中的元素。不同的 Spliterator 实现可能有不同的绑定时机。如果一个 Spliterator 是延迟绑定的，那么只有在进行首次遍历、首次切分或首次查询大小时，才会绑定到流上；反之，它会在创建时或首次调用任何方法时绑定到流上。绑定时机的重要性在于，在绑定之前对流所做的修改，在 Spliterator 遍历时是可见的。延迟绑定可以提供最大限度的灵活性。

### 有状态和无状态操作

流操作可以是有状态或无状态的。当一个有状态的操作在处理一个元素时，它可能需要使用处理之前的元素时保留的信息；无状态的操作可以独立处理每个元素，举例来说：

- `distinct` 和 `sorted` 是有状态操作的例子。`distinct` 操作从流中删除重复元素，它需要记录下之前已经遇到过的元素来确定当前元素是否应该被删除。`sorted` 操作对流进行排序，它需要知道所有元素来确定当前元素在排序之后的所在位置。
- `filter` 和 `map` 是无状态操作的例子。`filter` 操作在进行过滤时只需要看当前元素即可。`map` 操作可以独立转换当前元素。一般来说，有状态操作的运行代价要高于无状态操作，因为需要额外的空间保存中间状态信息。

`Stream<T>` 是表示流的接口，`T` 是流中元素的类型。对于原始类型的流，可以使用专门的类 `IntStream`、`LongStream` 和 `DoubleStream`。

### 流水线

在对流进行处理时，不同的流操作以级联的方式形成处理流水线。一个流水线由一个源（source），0 到多个中间操作（intermediate operation）和一个终结操作（terminal operation）完成。

- 源：源是流中元素的来源。Java 提供了很多内置的源，包括数组、集合、生成函数和 I/O 通道等。
- 中间操作：中间操作在一个流上进行操作，返回结果是一个新的流。这些操作是延迟执行的。
- 终结操作：终结操作遍历流来产生一个结果或是副作用。在一个流上执行终结操作之后，该流被消费，无法再次被消费。

流的处理流水线在其终结操作运行时才开始执行。

#### 源

Java 8 支持从不同的源中创建流。Stream.of 方法可以使用给定的元素创建一个顺序流。使用 `java.util.Arrays` 的静态方法可以从数组中创建流，如清单5 所示。

##### 清单 5\. 从数组中创建流

```
Arrays.stream(new String[] {"Hello", "World"})
.forEach(System.out::println);
// 输出"Hello\nWorld"到控制台

int sum = Arrays.stream(new int[] {1, 2, 3})
.reduce((a, b) -> a + b)
.getAsInt();
// "sum"的值是"6"

```

Show moreShow more icon

接口 Collection 的默认方法 `stream()` 和 `parallelStream()` 可以分别从集合中创建顺序流和并行流，如清单 6 所示。

##### 清单 6\. 从集合中创建流

```
List<String> list = new ArrayList<>();
list.add("Hello");
list.add("World");
list.stream()
.forEach(System.out::println);
// 输出 Hello 和 World

```

Show moreShow more icon

#### 中间操作

流中间操作在应用到流上，返回一个新的流。下面列出了常用的流中间操作：

- `map`：通过一个 Function 把一个元素类型为 T 的流转换成元素类型为 R 的流。
- `flatMap`：通过一个 Function 把一个元素类型为 T 的流中的每个元素转换成一个元素类型为 R 的流，再把这些转换之后的流合并。
- `filter`：过滤流中的元素，只保留满足由 Predicate 所指定的条件的元素。
- `distinct`：使用 `equals` 方法来删除流中的重复元素。
- `limit`：截断流使其最多只包含指定数量的元素。
- `skip`：返回一个新的流，并跳过原始流中的前 N 个元素。
- `sorted`：对流进行排序。
- `peek`：返回的流与原始流相同。当原始流中的元素被消费时，会首先调用 `peek` 方法中指定的 Consumer 实现对元素进行处理。
- `dropWhile`：从原始流起始位置开始删除满足指定 Predicate 的元素，直到遇到第一个不满足 Predicate 的元素。
- `takeWhile`：从原始流起始位置开始保留满足指定 Predicate 的元素，直到遇到第一个不满足 Predicate 的元素。

在 [清单 7](_清单7流的中间处理操作) 中，第一段代码展示了 `flatMap` 的用法，第二段代码展示了 `takeWhile` 和 `dropWhile` 的用法。

##### 清单 7\. 中间操作示例

```
Stream.of(1, 2, 3)
    .map(v -> v + 1)
    .flatMap(v -> Stream.of(v * 5, v * 10))
    .forEach(System.out::println);
//输出 10，20，15，30，20，40

Stream.of(1, 2, 3)
    .takeWhile(v -> v <  3)
    .dropWhile(v -> v <  2)
    .forEach(System.out::println);
//输出 2

```

Show moreShow more icon

#### 终结操作

终结操作产生最终的结果或副作用。下面是一些常见的终结操作。

`forEach` 和 `forEachOrdered` 对流中的每个元素执行由 Consumer 给定的实现。在使用 `forEach` 时，并没有确定的处理元素的顺序；`forEachOrdered` 则按照流的相遇顺序来处理元素，如果流有确定的相遇顺序的话。

`reduce` 操作把一个流约简成单个结果。约简操作可以有 3 个部分组成：

- 初始值：在对元素为空的流进行约简操作时，返回值为初始值。
- 叠加器：接受 2 个参数的 `BiFunction`。第一个参数是当前的约简值，第二个参数是当前元素，返回结果是新的约简值。
- 合并器：对于并行流来说，约简操作可能在流的不同部分上并行执行。合并器用来把部分约简结果合并为最终的结果。

在清单 8 中，第一个 `reduce` 操作是最简单的形式，只需要声明叠加器即可。初始值是流的第一个元素；第二个 `reduce` 操作提供了初始值和叠加器；第三个 `reduce` 操作声明了初始值、叠加器和合并器。

##### 清单 8\. reduce 操 作示例

```
Stream.of(1, 2, 3).reduce((v1, v2) -> v1 + v2)
    .ifPresent(System.out::println);
// 输出 6

int result1 = Stream.of(1, 2, 3, 4, 5)
    .reduce(1, (v1, v2) -> v1 * v2);
System.out.println(result1);
// 输出 120

int result2 = Stream.of(1, 2, 3, 4, 5)
    .parallel()
    .reduce(0, (v1, v2) -> v1 + v2, (v1, v2) -> v1 + v2);
System.out.println(result2);
// 输出 15

```

Show moreShow more icon

Max 和 min 是两种特殊的约简操作，分别求得流中元素的最大值和最小值。

对于一个流，操作 `allMatch`、`anyMatch` 和 `nonMatch` 分别用来检查是否流中的全部元素、任意元素或没有元素满足给定的条件。判断的条件由 Predicate 指定。

操作 `findFirst` 和 `findAny` 分别查找流中的第一个或任意一个元素。两个方法的返回值都是 `Optional` 对象。当流为空时，返回的是空的 `Optional` 对象。如果一个流没有确定的相遇顺序，那么 `findFirst` 和 `findAny` 的行为在本质上是相同的。

操作 `collect` 表示的是另外一类的约简操作。与 `reduce` 不同在于，`collect` 会把结果收集到可变的容器中，如 List 或 Set。收集操作通过接口 `java.util.stream.Collector` 来实现。Java 已经在类 `Collectors` 中提供了很多常用的 `Collector` 实现。

第一类收集操作是收集到集合中，常见的方法有 `toList()`、`toSet()` 和 `toMap()` 等。第二类收集操作是分组收集，即使用 `groupingBy` 对流中元素进行分组。分组时对流中所有元素应用同一个 `Function`。具有相同结果的元素被分到同一组。分组之后的结果是一个 `Map`，`Map` 的键是应用 `Function` 之后的结果，而对应的值是属于该组的所有元素的 List。在清单 9 中，流中的元素按照字符串的第一个字母分组，所得到的 Map 中的键是 A、B 和 D，而 A 对应的 List 值中包含了 Alex 和 Amy 两个元素，B 和 D 所对应的 List 值则只包含一个元素。

##### 清单 9\. 收集器 groupingBy 示 例

```
final Map<Character, List<String>> names = Stream.of("Alex", "Bob", "David", "Amy")
    .collect(Collectors.groupingBy(v -> v.charAt(0)));
System.out.println(names);

```

Show moreShow more icon

第三类的 `joining` 操作只对元素类型为 `CharSequence` 的流使用，其作用是把流中的字符串连接起来。清单 10 中把字符串流用”, “进行连接。

##### 清单 10\. 收集器 joining 示 例

```
String str = Stream.of("a", "b", "c")
.collect(Collectors.joining(", "));
System.out.println(str);

```

Show moreShow more icon

第四类的 `partitioningBy` 操作的作用类似于 `groupingBy`，只不过分组时使用的是 Predicate，也就是说元素最多分成两组。所得到结果的 Map 的键的类型是 Boolean，而值的类型同样是 List。

还有一些收集器可以进行数学计算，不过只对元素类型为 int、long 或 double 的流可用。这些数学计算包括：

- `averagingDouble`、`averagingInt` 和 `averagingLong` 计算流中元素的平均值。
- `summingDouble`、`summingInt` 和 `summingLong` 计算流中元素的和。
- `summarizingDouble`、`summarizingInt` 和 `summarizingLong` 对流中元素进行数学统计，可以得到平均值、数量、和、最大值和最小值。

清单 11 展示了这些数学计算相关的收集器的用法。

##### 清单 11\. 与数学计算相关的收集器

```
double avgLength = Stream.of("hello", "world", "a")
    .collect(Collectors.averagingInt(String::length));
System.out.println(avgLength);

final IntSummaryStatistics statistics = Stream.of("a", "b", "cd")
    .collect(Collectors.summarizingInt(String::length));
System.out.println(statistics.getAverage());
System.out.println(statistics.getCount());

```

Show moreShow more icon

Stream 中还有其他实用的操作，限于篇幅不能全部介绍。相关的用法可以查看 API 文档。

## 结束语

Java 8 引入的 Lambda 表达式和流处理是可以极大提高开发效率的重要特性。每个 Java 开发人员都应该熟练掌握它们的使用。本文从 JSR 335 出发对 Lambda 表达式进行了深入的介绍，同时也对流的特征和操作进行了详细说明。 [下一篇](https://developer.ibm.com/zh/articles/j-understanding-functional-programming-4/) 文章将对 Java 平台上流行的函数式编程库 Vavr 进行介绍。