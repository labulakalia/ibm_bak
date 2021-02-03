# 深入解析 Monad
函数式编程中最具有神秘色彩的概念

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-understanding-functional-programming-5/)

成 富

发布: 2018-12-03

* * *

在本系列的前四篇文章中对函数式编程进行了多方位的介绍。本文将着重介绍函数式编程中一个重要而又复杂的概念：Monad。一直以来，Monad 都是函数式编程中最具有神秘色彩的概念。正如 JSON 格式的提出者 Douglas Crockford 所指出的，Monad 有一种魔咒，一旦你真正理解了它的意义，就失去了解释给其他人的能力。本文尝试深入解析 Monad 这一概念。由于 Monad 的概念会涉及到一些数学理论，可能读起来会比较枯燥。本文侧重在 Monad 与编程相关的方面，并结合 Java 示例代码来进行说明。

## 范畴论

要解释 Monad，就必须提到范畴论（Category Theory）。范畴（category）本身是一个很简单的概念。一个范畴由对象（object）以及对象之间的箭头（arrow）组成。范畴的核心是组合，体现在箭头的组合性上。如果从对象 A 到对象 B 有一个箭头，从对象 B 到对象 C 也有一个箭头，那么必然有一个从对象 A 到对象 C 的箭头。从 A 到 C 的这个箭头，就是 A 到 B 的箭头和 B 到 C 的箭头的组合。这种组合的必然存在性，是范畴的核心特征。以专业术语来说，箭头被称为态射（morphisms）。范畴中对象和箭头的概念可以很容易地映射到函数中。类型可以作为范畴中的对象，把函数看成是箭头。如果有一个函数 f 的参数类型是 A，返回值类型是 B，那么这个函数是从 A 到 B 的态射；另外一个函数 g 的参数类型是 B，返回值类型是 C，这个函数是从 B 到 C 的态射。可以把 f 和 g 组合起来，得到一个新的从类型 A 到类型 C 的函数，记为 g ∘f，也就是从 A 到 C 的态射。这种函数的组合方式是必然存在的。

一个范畴中的组合需要满足两个条件：

- 组合必须是传递的（associative）。如果有 3 个态射 f、g 和 h 可以按照 h∘g∘f 的顺序组合，那么不管是 g 和 h 先组合，还是 f 和 g 先组合，所产生的结果都是一样的。
- 对于每个对象 A，都有一个作为组合基本单元的箭头。这个箭头的起始和终止都是该对象 A 本身。当该箭头与从对象 A 起始或结束的其他箭头组合时，得到的结果是原始的箭头。以函数的概念来说，这个函数称为恒等函数（identity function）。在 Java 中，这个函数由 `Function.identity()` 表示。

从编程的角度来说，范畴论的概念要求在设计时应该考虑对象的接口，而不是具体的实现。范畴论中的对象非常的抽象，没有关于对象的任何定义。我们只知道对象上的箭头，而对于对象本身则一无所知。对象实际上是由它们之间的相互组合关系来定义的。

范畴的概念虽然抽象，实际上也很容易找到现实的例子。最直接的例子是从有向图中创建出范畴。对于有向图中的每个节点，首先添加一个从当前节点到自身的箭头。然后对于每两条首尾相接的边，添加一条新的箭头连接起始和结束节点。如此反复，就得到了一个范畴。

范畴中的对象和态射的概念很抽象。从编程的角度来说，我们可以找到更好的表达方式。在程序中，讨论单个的对象实例并没有意义，更重要的是对象的类型。在各种编程语言中，我们已经认识了很多类型，包括 int、long、double 和 char 等。类型可以看成是值的集合。比如 bool 类型就只有两个值 true 和 false，int 类型包含所有的整数。类型的值可以是有限的，也可以是无限的。比如 String 类型的值是无限的。编程语言中的函数其实是从类型到类型的映射。对于参数超过 1 个的函数，总是可以使用柯里化来转换为只有一个参数的函数。

类型和函数可以分别与范畴中的对象和态射相对应。范畴中的对象是类型，而态射则是函数。类型的作用在于限定了范畴中态射可以组合的方式，也就是函数的组合方式。只有一个函数的返回值类型与另一个函数的参数类型匹配时，这两个函数才能并肯定可以组合。这也就满足了范畴的定义。

之前讨论的函数都是纯函数，不含任何副作用。而在实际的编程中，是离不开副作用的。纯函数适合于描述计算，但是没办法描述输出字符串到控制台或是写数据到文件这样的副作用。Monad 的作用正是解决了如何描述副作用的问题。实际上，纯粹的函数式编程语言 Haskell 正是用 Monad 来处理描述 IO 等基于副作用的操作。在介绍 Monad 之前，需要先说明 Functor。

## Functor

Functor 是范畴之间的映射。对于两个范畴 A 和 B，Functor F 把范畴 A 中的对象映射到范畴 B 中。Functor 在映射时会保留对象之间的连接关系。如果范畴 A 中存在从对象 a 到对象 b 的态射，那么 a 和 b 经过 Functor F 在范畴 B 中的映射值 F a 和 F b 之间也存在着态射。同样的，态射之间的组合关系，以及恒等态射都会被保留。所以说 Functor 不仅是范畴中对象之间的映射，也是态射之间的映射。如果一个 Functor 从一个范畴映射到自己，称为 endofunctor。

前面提到过，编程语言中的范畴中的对象是类型，而态射是函数。因此，这样的 endofunctor 是从类型到类型的映射，同时也是函数到函数的映射。我们首先看一个具体的 Functor ：Option。Option 的定义很简单，Java 标准库和 Vavr 中都有对应的类。不过我们这里讨论的 Option 与 Java 中的 Optional 类有很大不同。Option 本身是一个类型构造器，使用时需要提供一个类型，所得到的结果是另外一个新的类型。这里可以与 Java 中的泛型作为类比。Option 有两种可能的值：Some 和 None。Some 表示对应类型的一个值，而 None 表示没有值。对于一个从 a 到 b 的映射 f，可以很容易地找到与之对应的使用 Option 的映射。该映射把 None 对应到 None，而把 `f(Some a)` 映射到 `Some f(a)`。

## Monad

Monad 本身也是一种 Functor。Monad 的目的在于描述副作用。

### 函数的副作用与组合方式

清单 1 给出了一个简单的函数 `increase`。该函数的作用是返回输入的参数加 1 之后的值。除了进行计算之外，还通过 `count++` 来修改一个变量的值。这行语句的出现，使得函数 `increase` 不再是纯函数，每次调用都会对外部环境造成影响。

##### 清单 1\. 包含副作用的函数

```
int count = 0;

int increase(int x) {
count++;
return x + 1;
}

```

Show moreShow more icon

清单 1 中的函数 `increase` 可以划分成两个部分：产生副作用的 `count++`，以及剩余的不产生副作用的部分。如果可以通过一些转换，把副作用从函数 `increase` 中剥离出来，那么就可以得到另外一个纯函数的版本 `increase1`，如清单 2 所示。对函数 `increase1` 来说，我们可以把返回值改成一个 Vavr 中的 `Tuple2<Integer, Integer>` 类型，分别包含函数原始的返回值 `x + 1` 和在 `counter` 上增加的增量值 1。通过这样的转换之后，函数 `increase1` 就变成了一个纯函数。

##### 清单 2\. 转换之后的纯函数版本

```
Tuple2<Integer, Integer> increase1(int x) {
return Tuple.of(x + 1, 1);
}

```

Show moreShow more icon

在经过这样的转换之后，对于函数 `increase1` 的调用方式也发生了变化，如清单 3 所示。递增之后的值需要从 Tuple2 中获取，而 `count` 也需要通过 Tuple2 的值来更新。

##### 清单 3\. 调用转换之后的纯函数版本

```
int x = 0;
Tuple2<Integer, Integer> result = increase1(x);
x = result._1;
count += result._2;

```

Show moreShow more icon

我们可以采用同样的方式对另外一个相似的函数 `decrease` 做转换，如清单 4 所示。

##### 清单 4\. 函数 decrease 及其纯函数版本

```
int decrease(int x) {
count++;
return x - 1;
}

Tuple2<Integer, Integer> decrease1(int x) {
return Tuple.of(x - 1, 1);
}

```

Show moreShow more icon

不过需要注意的是，经过这样的转换之后，函数的组合方式发生了变化。对于之前的 `increase` 和 `decrease` 函数，可以直接组合，因为它们的参数和返回值类型是匹配的，如类似 `increase(decrease(x))` 或是 `decrease(increase(x))` 这样的组合方式。而经过转换之后的 `increase1` 和 `decrease1`，由于返回值类型改变，`increase1` 和 `decrease1` 不能按照之前的方式进行组合。函数 `increase1` 的返回值类型与 `decrease1` 的参数类型不匹配。对于这两个函数，需要另外的方式来组合。

在清单 5 中，`compose` 方法把两个类型为 `Function<Integer, Tuple2<Integer, Integer>>` 的函数 func1 和 func2 进行组合，返回结果是另外一个类型为 `Function<Integer, Tuple2<Integer, Integer>>` 的函数。在进行组合时，Tuple2 的第一个元素是实际需要返回的结果，按照纯函数组合的方式来进行，也就是把 func1 调用结果的 Tuple2 的第一个元素作为输入参数来调用 func2。Tuple2 的第二个元素是对 count 的增量。需要把这两个增量相加，作为 `compose` 方法返回的 Tuple2 的第二个元素。

##### 清单 5\. 函数的组合方式

```
Function<Integer, Tuple2<Integer, Integer>> compose(
    Function<Integer, Tuple2<Integer, Integer>> func1,
    Function<Integer, Tuple2<Integer, Integer>> func2) {
return x -> {
    Tuple2<Integer, Integer> result1 = func1.apply(x);
    Tuple2<Integer, Integer> result2 = func2.apply(result1._1);
    return Tuple.of(result2._1, result1._2 + result2._2);
};
}

```

Show moreShow more icon

清单 6 中的 `doCompose` 函数对 `increase1` 和 `decrease1` 进行组合。对于一个输入 `x`，由于 `increase1` 和 `decrease1` 的作用相互抵消，得到的结果是值为 `(x, 2)` 的对象。

##### 清单 6\. 函数组合示例

```
Tuple2<Integer, Integer> doCompose(int x) {
return compose(this::increase1, this::decrease1).apply(x);
}

```

Show moreShow more icon

可以看到，`doCompose` 函数的输入参数和返回值类型与 `increase1` 和 `decrease1` 相同。所返回的结果可以继续使用 `doCompose` 函数来与其他类型相同的函数进行组合。

### Monad 的定义

现在回到函数 `increase` 和 `decrease`。从范畴论的角度出发，我们考虑下面一个范畴。该范畴中的对象仍然是 int 和 bool 等类型，但是其中的态射不再是简单的如 `increase` 和 `decrease` 这样的函数，而是把这些函数通过类似从 `increase` 到 `increase1` 这样的方式转换之后的函数。范畴中的态射必须是可以组合的，而这些函数的组合是通过调用类似 `doCompose` 这样的函数完成的。这样就满足了范畴的第一条原则。而第二条原则也很容易满足，只需要把参数 x 的值设为 0，就可以得到组合的基本单元。由此可以得出，我们定义了一个新的范畴，而这个范畴就叫做 Kleisli 范畴。每个 Kleisli 范畴所使用的函数转换方式是独特的。 [清单 2](#清单-2-转换之后的纯函数版本) 中的示例使用 Tuple2 来保存 `count` 的增量。与之对应的，Kleisli 范畴中对态射的组合方式也是独特的，类似清单 6 中的 doCompose 函数。

在对 Kleisli 范畴有了一个直观的了解之后，就可以对 Monad 给出一个形式化的定义。给定一个范畴 C 和 `endofunctor m`，与之相对应的 Kleisli 范畴中的对象与范畴 C 相同，但态射是不同的。K 中的两个对象 a 和 b 之间的态射，是由范畴 C 中的 a 到 m(b) 的态射来实现的。注意，Kleisli 范畴 K 中的态射箭头是从对象 a 到对象 b 的，而不是从对象 a 到 m(b)。如果存在一种传递的组合方式，并且每个对象都有组合单元箭头，也就是满足范畴的两大原则，那么这个 endofunctor m 就叫做 Monad。

一个 Monad 的定义中包含了 3 个要素。在定义 Monad 时需要提供一个类型构造器 M 和两个操作 unit 和 bind：

- 类型构造器的作用是从底层的类型中创建出一元类型（monadic type）。如果 M 是 Monad 的名称，而 t 是数据类型，则 M t 是对应的一元类型。
- `unit` 操作把一个普通值 t 通过类型构造器封装在一个容器中，所产生的值的类型是 M t。`unit` 操作也称为 `return` 操作。`return` 操作的名称来源于 Haskell。不过由于 `return` 在很多编程语言中是保留关键词，用 unit 做名称更为合适。
- `bind` 操作的类型声明是 `(M t)→(t→M u)→(M u)`。该操作接受类型为 M t 的值和类型为 `t → M u` 的函数来对值进行转换。在进行转换时，`bind` 操作把原始值从容器中抽取出来，再应用给定的函数进行转换。函数的返回值是一个新的容器值 M u。M u 可以作为下一次转换的起点。多个 `bind` 操作可以级联起来，形成处理流水线。

如果只看 Monad 的定义，会有点晦涩难懂。实际上 [清单 2](#清单-2-转换之后的纯函数版本) 中的示例就是一种常见的 Monad，称为 Writer Monad。下面我们结合 Java 代码来看几种常见的 Monad。

### Writer Monad

[清单 2](#清单-2-转换之后的纯函数版本) 展示了 Writer Monad 的一种用法，也就是累积 `count` 的值。实际上，Writer Monad 的主要作用是在函数调用过程中收集辅助信息，比如日志信息或是性能计数器等。其基本的思想是把副作用中对外部环境的修改聚合起来，从而把副作用从函数中分离出来。聚合的方式取决于所产生的副作用。 [清单 2](#清单-2-转换之后的纯函数版本) 中的副作用是修改计算器 count，相应的聚合方式是累加计数值。如果副作用是产生日志，相应的聚合方式是连接日志记录的字符串。聚合方式是每个 Writer Monad 的核心。对于聚合方式的要求和范畴中对于态射的要求是一样，也就是必须是传递的，而且有组合的基本单元。在 [清单 5](#清单-5-函数的组合方式) 中，聚合方式是 Integer 类型的相加操作，是传递的；同时也有基本单元，也就是加零。

下面对 Writer Monad 进行更加形式化的说明。Writer Monad 除了其本身的类型 T 之外，还有另外一个辅助类型 W，用来表示聚合值。对类型 W 的要求是前面提到的两点，也就是存在传递的组合操作和基本单元。Writer Monad 的 unit 操作比较简单，返回的是类型 T 的值 t 和类型 W 的基本单元。而 bind 操作则需要分别转换类型 T 和 W 的值。对于 T 的值，按照 Monad 自身的定义来转换；而对于 W 的值，则使用该类型的传递操作来聚合值。聚合的结果作为转换之后的新的 W 的值。

清单 7 中是记录日志的 Writer Monad 的实例。该 Monad 自身的类型使用 Java 泛型类型 T 来表示，而辅助类型是 `List<String>`，用来保存记录的日志。`List<String>` 满足作为辅助类型的要求。`List<String>` 上的相加操作是传递的，也存在作为基本单元的空列表。LoggingMonad 中的 `unit` 方法返回传入的值 value 和空列表。`bind` 方法的第一个参数是 `LoggingMonad<T1>` 类型，作为变换的输入；第二个参数是 F`unction<T1, LoggingMonad<T2>>` 类型，用来把类型 T1 转换成新的 `LoggingMonad<T2>` 类型。辅助类型 `List<String>` 中的值通过列表相加的方式进行组合。方法 `pipeline` 表示一个处理流水线，对于一个输入 Monad，依次应用指定的变换，得到最终的结果。在使用示例中，`LoggingMonad` 中封装的是 Integer 类型，第一个转换把值乘以 4，第二个变换把值除以 2。每个变换都记录自己的日志。在运行流水线之后，得到的结果包含了转换之后的值和聚合的日志。

##### 清单 7\. 记录日志的 Monad

```
public class LoggingMonad<T> {

private final T value;
private final List<String> logs;

public LoggingMonad(T value, List<String> logs) {
    this.value = value;
    this.logs = logs;
}

@Override
public String toString() {
    return "LoggingMonad{" +
        "value=" + value +
        ", logs=" + logs +
        '}';
}

public static <T> LoggingMonad<T> unit(T value) {
    return new LoggingMonad<>(value, List.of());
}

public static <T1, T2> LoggingMonad<T2> bind(LoggingMonad<T1> input,
      Function<T1, LoggingMonad<T2>> transform) {
    final LoggingMonad<T2> result = transform.apply(input.value);
    List<String> logs = new ArrayList<>(input.logs);
    logs.addAll(result.logs);
    return new LoggingMonad<>(result.value, logs);
}

public static <T> LoggingMonad<T> pipeline(LoggingMonad<T> monad,
      List<Function<T, LoggingMonad<T>>> transforms) {
    LoggingMonad<T> result = monad;
    for (Function<T, LoggingMonad<T>> transform : transforms) {
      result = bind(result, transform);
    }
    return result;
}

public static void main(String[] args) {
    Function<Integer, LoggingMonad<Integer>> transform1 =
        v -> new LoggingMonad<>(v * 4, List.of(v + " * 4"));
    Function<Integer, LoggingMonad<Integer>> transform2 =
        v -> new LoggingMonad<>(v / 2, List.of(v + " / 2"));
    final LoggingMonad<Integer> result =
pipeline(LoggingMonad.unit(8),
        List.of(transform1, transform2));
    System.out.println(result); // 输出为 LoggingMonad{value=16,
logs=[8 * 4, 32 / 2]}
}
}

```

Show moreShow more icon

### Reader Monad

Reader Monad 也被称为 Environment Monad，描述的是依赖共享环境的计算。Reader Monad 的类型构造器从类型 T 中创建出一元类型 `E → T`，而 E 是环境的类型。类型构造器把类型 T 转换成一个从类型 E 到 T 的函数。Reader Monad 的 unit 操作把类型 T 的值 t 转换成一个永远返回 t 的函数，而忽略类型为 E 的参数；`bind` 操作在转换时，在所返回的函数的函数体中对类型 T 的值 t 进行转换，同时保持函数的结构不变。

清单 8 是 Reader Monad 的示例。`Function<E, T>` 是一元类型的声明。ReaderMonad 的 `unit` 方法返回的 Function 只是简单的返回参数值 value。而 `bind` 方法的第一个参数是一元类型 `Function<E, T1>`，第二个参数是把类型 T1 转换成 `Function<E, T2>` 的函数，返回值是另外一个一元类型 `Function<E, T2>`。bind 方法的转换逻辑首先通过 `input.apply(e)` 来得到类型为 T1 的值，再使用 `transform.apply` 来得到类型为 `Function<E, T2>>` 的值，最后使用 `apply(e)` 来得到类型为 T2 的值。

##### 清单 8\. Reader Monad 示例

```
public class ReaderMonad {

public static <T, E> Function<E, T> unit(T value) {
    return e -> value;
}

public static <T1, T2, E> Function<E, T2> bind(Function<E, T1>
input, Function<T1, Function<E, T2>> transform) {
    return e -> transform.apply(input.apply(e)).apply(e);
}

public static void main(String[] args) {
    Function<Environment, String> m1 = unit("Hello");
    Function<Environment, String> m2 = bind(m1, value -> e ->
e.getPrefix() + value);
    Function<Environment, Integer> m3 = bind(m2, value -> e ->
e.getBase() + value.length());
    int result = m3.apply(new Environment());
    System.out.println(result);
}
}

```

Show moreShow more icon

清单 8 中使用的环境类型 `Environment` 如清单 9 所示，其中有两个方法 `getPrefix` 和 `getBase` 分别返回相应的值。清单 8 的 `m1` 是值为 `Hello` 的单元类型，`m2` 使用了 `Environment` 的 `getPrefix` 方法进行转换，而 `m3` 使用了 `getBase` 方法进行转换，最终输出的结果是 107。因为字符串 `Hello` 在添加了前缀 `$$` 之后的长度是 7，与 100 相加之后的值是 107。

##### 清单 9\. 环境类型

```
public class Environment {

public String getPrefix() {
    return "$$";
}

public int getBase() {
    return 100;
}
}

```

Show moreShow more icon

### State Monad

State Monad 可以在计算中附加任意类型的状态值。State Monad 与 Reader Monad 相似，只是 State Monad 在转换时会返回一个新的状态对象，从而可以描述可变的环境。State Monad 的类型构造器从类型 T 中创建一个函数类型，该函数类型的参数是状态对象的类型 S，而返回值包含类型 S 和 T 的值。State Monad 的 unit 操作返回的函数只是简单地返回输入的类型 S 的值；`bind` 操作所返回的函数类型负责在执行时传递正确的状态对象。

清单 10 给出了 State Monad 的示例。State Monad 使用元组 `Tuple2<T, S>` 来保存计算值和状态对象，所对应的一元类型是 `Function<S, Tuple2<T, S>>` 表示的函数。`unit` 方法所返回的函数只是简单地返回输入状态对象。`bind` 方法的转换逻辑使用 `input.apply(s)` 得到 T1 和 S 的值，再用得到的 S 值调用 `transform`。

##### 清单 10\. State Monad 示例

```
public class StateMonad {

public static <T, S> Function<S, Tuple2<T, S>> unit(T value) {
    return s -> Tuple.of(value, s);
}

public static <T1, T2, S> Function<S, Tuple2<T2, S>>
bind(Function<S, Tuple2<T1, S>> input,
      Function<T1, Function<S, Tuple2<T2, S>>> transform) {
    return s -> {
      Tuple2<T1, S> result = input.apply(s);
      return transform.apply(result._1).apply(result._2);
    };
}

public static void main(String[] args) {
    Function<String, Function<String, Function<State, Tuple2<String,
State>>>> transform =
        prefix -> value -> s -> Tuple
            .of(prefix + value, new State(s.getValue() +
value.length()));

    Function<State, Tuple2<String, State>> m1 = unit("Hello");
    Function<State, Tuple2<String, State>> m2 = bind(m1,
transform.apply("1"));
    Function<State, Tuple2<String, State>> m3 = bind(m2,
transform.apply("2"));
    Tuple2<String, State> result = m3.apply(new State(0));
    System.out.println(result);
}
}

```

Show moreShow more icon

State Monad 中使用的状态对象如清单 11 所示。`State` 是一个包含值 value 的不可变对象。清单 10 中的 m1 封装了值 `Hello`。`transform` 方法用来从输入的字符串前缀 `prefix` 中创建转换函数。转换函数会在字符串值上添加给定的前缀，同时会把字符串的长度进行累加。转换函数每次都返回一个新的 `State` 对象。转换之后的结果中字符串的值是 `21Hello`，而 `State` 对象中的 value 为 11，是字符串 `Hello` 和 `1Hello` 的长度相加的结果。

##### 清单 11\. 状态对象

```
public class State {

private final int value;

public State(final int value) {
    this.value = value;
}

public int getValue() {
    return value;
}

@Override
public String toString() {
    return "State{" +
        "value=" + value +
        '}';
}
}

```

Show moreShow more icon

## 结束语

作为本系列的最后一篇文章，本文对函数式编程中的重要概念 Monad 做了详细的介绍。本文从范畴论出发，介绍了使用 Monad 描述函数副作用的动机和方式，以及 Monad 的定义。本文还对常见的几种 Monad 进行了介绍，并添加了相应的 Java 代码。