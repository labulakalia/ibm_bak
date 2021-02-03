# 使用 Vavr 进行函数式编程
Vavr 让函数式编程变得更加简洁易用

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-understanding-functional-programming-4/)

成 富

发布: 2018-12-03

* * *

在本系列的 [上一篇](https://www.ibm.com/developerworks/cn/java/j-understanding-functional-programming-3/index.html) 文章中对 Java 平台提供的 Lambda 表达式和流做了介绍。受限于 Java 标准库的通用性要求和二进制文件大小，Java 标准库对函数式编程的 API 支持相对比较有限。函数的声明只提供了 Function 和 BiFunction 两种，流上所支持的操作的数量也较少。为了更好地进行函数式编程，我们需要第三方库的支持。Vavr 是 Java 平台上函数式编程库中的佼佼者。

Vavr 这个名字对很多开发人员可能比较陌生。它的前身 Javaslang 可能更为大家所熟悉。Vavr 作为一个标准的 Java 库，使用起来很简单。只需要添加对 `io.vavr:vavr` 库的 Maven 依赖即可。Vavr 需要 Java 8 及以上版本的支持。本文基于 Vavr 0.9.2 版本，示例代码基于 Java 10。

## 元组

元组（Tuple）是固定数量的不同类型的元素的组合。元组与集合的不同之处在于，元组中的元素类型可以是不同的，而且数量固定。元组的好处在于可以把多个元素作为一个单元传递。如果一个方法需要返回多个值，可以把这多个值作为元组返回，而不需要创建额外的类来表示。根据元素数量的不同，Vavr 总共提供了 Tuple0、Tuple1 到 Tuple8 等 9 个类。每个元组类都需要声明其元素类型。如 `Tuple2<String, Integer>` 表示的是两个元素的元组，第一个元素的类型为 String，第二个元素的类型为 Integer。对于元组对象，可以使用 `_1`、`_2` 到 `_8` 来访问其中的元素。所有元组对象都是不可变的，在创建之后不能更改。

元组通过接口 Tuple 的静态方法 of 来创建。元组类也提供了一些方法对它们进行操作。由于元组是不可变的，所有相关的操作都返回一个新的元组对象。在 清单 1 中，使用 Tuple.of 创建了一个 Tuple2 对象。Tuple2 的 map 方法用来转换元组中的每个元素，返回新的元组对象。而 apply 方法则把元组转换成单个值。其他元组类也有类似的方法。除了 `map` 方法之外，还有 map1、map2、map3 等方法来转换第 N 个元素；update1、update2 和 update3 等方法用来更新单个元素。

##### 清单 1\. 使用元组

```
Tuple2<String, Integer> tuple2 = Tuple.of("Hello", 100);
Tuple2<String, Integer> updatedTuple2 = tuple2.map(String::toUpperCase, v -> v * 5);
String result = updatedTuple2.apply((str, number) -> String.join(", ",
str, number.toString()));
System.out.println(result);

```

Show moreShow more icon

虽然元组使用起来很方便，但是不宜滥用，尤其是元素数量超过 3 个的元组。当元组的元素数量过多时，很难明确地记住每个元素的位置和含义，从而使得代码的可读性变差。这个时候使用 Java 类是更好的选择。

## 函数

Java 8 中只提供了接受一个参数的 Function 和接受 2 个参数的 BiFunction。Vavr 提供了函数式接口 Function0、Function1 到 Function8，可以描述最多接受 8 个参数的函数。这些接口的方法 `apply` 不能抛出异常。如果需要抛出异常，可以使用对应的接口 CheckedFunction0、CheckedFunction1 到 CheckedFunction8。

Vavr 的函数支持一些常见特征。

### 组合

函数的组合指的是用一个函数的执行结果作为参数，来调用另外一个函数所得到的新函数。比如 f 是从 x 到 y 的函数，g 是从 y 到 z 的函数，那么 `g(f(x))`是从 x 到 z 的函数。Vavr 的函数式接口提供了默认方法 `andThen` 把当前函数与另外一个 Function 表示的函数进行组合。Vavr 的 Function1 还提供了一个默认方法 compose 来在当前函数执行之前执行另外一个 Function 表示的函数。

在清单 2 中，第一个 function3 进行简单的数学计算，并使用 andThen 把 function3 的结果乘以 100。第二个 function1 从 String 的 `toUpperCase` 方法创建而来，并使用 `compose` 方法与 Object 的 `toString` 方法先进行组合。得到的方法对任何 Object 先调用 `toString`，再调用 `toUpperCase`。

##### 清单 2\. 函数的组合

```
Function3< Integer, Integer, Integer, Integer> function3 = (v1, v2, v3)
-> (v1 + v2) * v3;
Function3< Integer, Integer, Integer, Integer> composed =
function3.andThen(v -> v * 100);
int result = composed.apply(1, 2, 3);
System.out.println(result);
// 输出结果 900

Function1< String, String> function1 = String::toUpperCase;
Function1< Object, String> toUpperCase = function1.compose(Object::toString);
String str = toUpperCase.apply(List.of("a", "b"));
System.out.println(str);
// 输出结果[A, B]

```

Show moreShow more icon

### 部分应用

在 Vavr 中，函数的 apply 方法可以应用不同数量的参数。如果提供的参数数量小于函数所声明的参数数量（通过 `arity()` 方法获取），那么所得到的结果是另外一个函数，其所需的参数数量是剩余未指定值的参数的数量。在清单 3 中，Function4 接受 4 个参数，在 `apply` 调用时只提供了 2 个参数，得到的结果是一个 Function2 对象。

##### 清单 3\. 函数的部分应用

```
Function4< Integer, Integer, Integer, Integer, Integer> function4 =
(v1, v2, v3, v4) -> (v1 + v2) * (v3 + v4);
Function2< Integer, Integer, Integer> function2 = function4.apply(1, 2);
int result = function2.apply(4, 5);
System.out.println(result);
// 输出 27

```

Show moreShow more icon

### 柯里化方法

使用 `curried` 方法可以得到当前函数的柯里化版本。由于柯里化之后的函数只有一个参数，`curried` 的返回值都是 Function1 对象。在清单 4 中，对于 function3，在第一次的 `curried` 方法调用得到 Function1 之后，通过 `apply` 来为第一个参数应用值。以此类推，通过 3 次的 `curried` 和 `apply` 调用，把全部 3 个参数都应用值。

##### 清单 4\. 函数的柯里化

```
Function3<Integer, Integer, Integer, Integer> function3 = (v1, v2, v3)
-> (v1 + v2) * v3;
int result =
function3.curried().apply(1).curried().apply(2).curried().apply(3);
System.out.println(result);

```

Show moreShow more icon

### 记忆化方法

使用记忆化的函数会根据参数值来缓存之前计算的结果。对于同样的参数值，再次的调用会返回缓存的值，而不需要再次计算。这是一种典型的以空间换时间的策略。可以使用记忆化的前提是函数有引用透明性。

在清单 5 中，原始的函数实现中使用 `BigInteger` 的 `pow` 方法来计算乘方。使用 `memoized` 方法可以得到该函数的记忆化版本。接着使用同样的参数调用两次并记录下时间。从结果可以看出来，第二次的函数调用的时间非常短，因为直接从缓存中获取结果。

##### 清单 5\. 函数的记忆化

```
Function2<BigInteger, Integer, BigInteger> pow = BigInteger::pow;
Function2<BigInteger, Integer, BigInteger> memoized = pow.memoized();
long start = System.currentTimeMillis();
memoized.apply(BigInteger.valueOf(1024), 1024);
long end1 = System.currentTimeMillis();
memoized.apply(BigInteger.valueOf(1024), 1024);
long end2 = System.currentTimeMillis();
System.out.printf("%d ms -> %d ms", end1 - start, end2 - end1);

```

Show moreShow more icon

注意，`memoized` 方法只是把原始的函数当成一个黑盒子，并不会修改函数的内部实现。因此，`memoized` 并不适用于直接封装本系列 [第二篇](https://www.ibm.com/developerworks/cn/java/j-understanding-functional-programming-2/index.html) 文章中用递归方式计算斐波那契数列的函数。这是因为在函数的内部实现中，调用的仍然是没有记忆化的函数。

## 值

Vavr 中提供了一些不同类型的值。

### Option

Vavr 中的 `Option` 与 Java 8 中的 `Optional` 是相似的。不过 Vavr 的 Option 是一个接口，有两个实现类 `Option.Some` 和 `Option.None`，分别对应有值和无值两种情况。使用 `Option.some` 方法可以创建包含给定值的 `Some` 对象，而 `Option.none` 可以获取到 `None` 对象的实例。`Option` 也支持常用的 `map`、`flatMap` 和 `filter` 等操作，如清单 6 所示。

##### 清单 6\. 使用 Option 的示例

```
Option<String> str = Option.of("Hello");
str.map(String::length);
str.flatMap(v -> Option.of(v.length()));

```

Show moreShow more icon

### Either

`Either` 表示可能有两种不同类型的值，分别称为左值或右值。只能是其中的一种情况。`Either` 通常用来表示成功或失败两种情况。惯例是把成功的值作为右值，而失败的值作为左值。可以在 `Either` 上添加应用于左值或右值的计算。应用于右值的计算只有在 `Either` 包含右值时才生效，对左值也是同理。

在清单 7 中，根据随机的布尔值来创建包含左值或右值的 `Either` 对象。`Either` 的 `map` 和 `mapLeft` 方法分别对右值和左值进行计算。

##### 清单 7\. 使用 Either 的示例

```
import io.vavr.control.Either;
import java.util.concurrent.ThreadLocalRandom;

public class Eithers {

private static ThreadLocalRandom random =
ThreadLocalRandom.current();

public static void main(String[] args) {
    Either<String, String> either = compute()
        .map(str -> str + " World")
        .mapLeft(Throwable::getMessage);
    System.out.println(either);
}

private static Either<Throwable, String> compute() {
    return random.nextBoolean()
        ? Either.left(new RuntimeException("Boom!"))
        : Either.right("Hello");
}
}

```

Show moreShow more icon

### Try

`Try` 用来表示一个可能产生异常的计算。`Try` 接口有两个实现类，`Try.Success` 和 `Try.Failure`，分别表示成功和失败的情况。`Try.Success` 封装了计算成功时的返回值，而 `Try.Failure` 则封装了计算失败时的 `Throwable` 对象。Try 的实例可以从接口 `CheckedFunction0`、`Callable`、`Runnable` 或 `Supplier` 中创建。`Try` 也提供了 `map` 和 `filter` 等方法。值得一提的是 `Try` 的 `recover` 方法，可以在出现错误时根据异常进行恢复。

在清单 8 中，第一个 `Try` 表示的是 `1/0` 的结果，显然是异常结果。使用 `recover` 来返回 1。第二个 `Try` 表示的是读取文件的结果。由于文件不存在，`Try` 表示的也是异常。

##### 清单 8\. 使用 Try 的示例

```
Try<Integer> result = Try.of(() -> 1 / 0).recover(e -> 1);
System.out.println(result);

Try<String> lines = Try.of(() -> Files.readAllLines(Paths.get("1.txt")))
    .map(list -> String.join(",", list))
    .andThen((Consumer<String>) System.out::println);
System.out.println(lines);

```

Show moreShow more icon

### Lazy

`Lazy` 表示的是一个延迟计算的值。在第一次访问时才会进行求值操作，而且该值只会计算一次。之后的访问操作获取的是缓存的值。在清单 9 中，`Lazy.of` 从接口 `Supplier` 中创建 `Lazy` 对象。方法 `isEvaluated` 可以判断 `Lazy` 对象是否已经被求值。

##### 清单 9\. 使用 Lazy 的示例

```
Lazy<BigInteger> lazy = Lazy.of(() ->
BigInteger.valueOf(1024).pow(1024));
System.out.println(lazy.isEvaluated());
System.out.println(lazy.get());
System.out.println(lazy.isEvaluated());

```

Show moreShow more icon

## 数据结构

Vavr 重新在 Iterable 的基础上实现了自己的集合框架。Vavr 的集合框架侧重在不可变上。Vavr 的集合类在使用上比 Java 流更简洁。

Vavr 的 Stream 提供了比 Java 中 Stream 更多的操作。可以使用 `Stream.ofAll` 从 Iterable 对象中创建出 Vavr 的 Stream。下面是一些 Vavr 中添加的实用操作：

- `groupBy`：使用 Fuction 对元素进行分组。结果是一个 Map，Map 的键是分组的函数的结果，而值则是包含了同一组中全部元素的 Stream。
- `partition`：使用 Predicate 对元素进行分组。结果是包含 2 个 Stream 的 Tuple2。Tuple2 的第一个 Stream 的元素满足 Predicate 所指定的条件，第二个 Stream 的元素不满足 Predicate 所指定的条件。
- `scanLeft` 和 `scanRight`：分别按照从左到右或从右到左的顺序在元素上调用 Function，并累积结果。
- `zip`：把 Stream 和一个 Iterable 对象合并起来，返回的结果 Stream 中包含 Tuple2 对象。Tuple2 对象的两个元素分别来自 Stream 和 Iterable 对象。

在清单 10 中，第一个 `groupBy` 操作把 Stream 分成奇数和偶数两组；第二个 `partition` 操作把 Stream 分成大于 2 和不大于 2 两组；第三个 `scanLeft` 对包含字符串的 Stream 按照字符串长度进行累积；最后一个 `zip` 操作合并两个流，所得的结果 Stream 的元素数量与长度最小的输入流相同。

##### 清单 10\. Stream 的使用示例

```
Map<Boolean, List<Integer>> booleanListMap = Stream.ofAll(1, 2, 3, 4, 5)
    .groupBy(v -> v % 2 == 0)
    .mapValues(Value::toList);
System.out.println(booleanListMap);
// 输出 LinkedHashMap((false, List(1, 3, 5)), (true, List(2, 4)))

Tuple2<List<Integer>, List<Integer>> listTuple2 = Stream.ofAll(1, 2, 3, 4)
    .partition(v -> v > 2)
    .map(Value::toList, Value::toList);
System.out.println(listTuple2);
// 输出 (List(3, 4), List(1, 2))

List<Integer> integers = Stream.ofAll(List.of("Hello", "World", "a"))
    .scanLeft(0, (sum, str) -> sum + str.length())
    .toList();
System.out.println(integers);
// 输出 List(0, 5, 10, 11)

List<Tuple2<Integer, String>> tuple2List = Stream.ofAll(1, 2, 3)
    .zip(List.of("a", "b"))
    .toList();
System.out.println(tuple2List);
// 输出 List((1, a), (2, b))

```

Show moreShow more icon

Vavr 提供了常用的数据结构的实现，包括 List、Set、Map、Seq、Queue、Tree 和 TreeMap 等。这些数据结构的用法与 Java 标准库的对应实现是相似的，但是提供的操作更多，使用起来也更方便。在 Java 中，如果需要对一个 List 的元素进行 map 操作，需要使用 stream 方法来先转换为一个 Stream，再使用 map 操作，最后再通过收集器 `Collectors.toList` 来转换回 List。而在 Vavr 中，List 本身就提供了 map 操作。清单 11 中展示了这两种使用方式的区别。

##### 清单 11\. Vavr 中数据结构的用法

```
List.of(1, 2, 3).map(v -> v + 10); //Vavr
java.util.List.of(1, 2, 3).stream()
.map(v -> v + 10).collect(Collectors.toList()); //Java 中 Stream

```

Show moreShow more icon

## 模式匹配

在 Java 中，我们可以使用 switch 和 case 来根据值的不同来执行不同的逻辑。不过 switch 和 case 提供的功能很弱，只能进行相等匹配。Vavr 提供了模式匹配的 API，可以对多种情况进行匹配和执行相应的逻辑。在清单 12 中，我们使用 Vavr 的 Match 和 Case 替换了 Java 中的 switch 和 case。Match 的参数是需要进行匹配的值。Case 的第一个参数是匹配的条件，用 Predicate 来表示；第二个参数是匹配满足时的值。`$(value)` 表示值为 value 的相等匹配，而 `$()` 表示的是默认匹配，相当于 switch 中的 default。

##### 清单 12\. 模式匹配的示例

```
String input = "g";
String result = Match(input).of(
    Case($("g"), "good"),
    Case($("b"), "bad"),
    Case($(), "unknown")
);
System.out.println(result);
// 输出 good

```

Show moreShow more icon

在清单 13 中，我们用 `$(v -> v > 0)` 创建了一个值大于 0 的 Predicate。这里匹配的结果不是具体的值，而是通过 `run` 方法来产生副作用。

##### 清单 13\. 使用模式匹配来产生副作用

```
int value = -1;
Match(value).of(
    Case($(v -> v > 0), o -> run(() -> System.out.println("> 0"))),
    Case($(0), o -> run(() -> System.out.println("0"))),
    Case($(), o -> run(() -> System.out.println("< 0")))
);
// 输出<  0

```

Show moreShow more icon

## 结束语

当需要在 Java 平台上进行复杂的函数式编程时，Java 标准库所提供的支持已经不能满足需求。Vavr 作为 Java 平台上流行的函数式编程库，可以满足不同的需求。本文对 Vavr 提供的元组、函数、值、数据结构和模式匹配进行了详细的介绍。 [下一篇](https://www.ibm.com/developerworks/cn/java/j-understanding-functional-programming-5/index.html) 文章将介绍函数式编程中的重要概念 Monad。