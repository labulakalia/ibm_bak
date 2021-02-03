# Java 知道您的类型
学习如何在 lambda 表达式中使用类型推断，掌握改进参数命名的技巧

**标签:** Java,Java 平台

[原文链接](https://developer.ibm.com/zh/articles/j-java8idioms8/)

Venkat Subramaniam

发布: 2017-11-14

* * *

**关于本系列**

Java 8 是自 Java 语言诞生以来进行的一次最重大更新  包含了非常丰富的新功能，您可能想知道从何处开始着手了解它。在 [本系列](/zh/series/java-8-idioms/) 中，作者兼教师 Venkat Subramaniam 提供了一种惯用的 Java 8 编程方法：这些简短的探索会激发您反思您认为理所当然的 Java 约定，同时逐步将新方法和语法集成到您的程序中。

Java™ 8 是第一个支持类型推断的 Java 版本，而且它仅对 lambda 表达式支持此功能。在 lambda 表达式中使用类型推断具有强大的作用，它将帮助您做好准备以应对未来的 Java 版本，在今后的版本中还会将类型推断用于变量等更多可能。这里的诀窍在于恰当地命名参数，相信 Java 编译器会推断出剩余的信息。

大多数时候，编译器完全能够推断类型。在它无法推断出来的时候，就会报错。

了解 lambda 表达式中的类型推断的工作原理，至少查看一个无法推断类型的示例。即使如此，也有解决办法。

## 显式类型和冗余

假设您询问某个人“您叫什么名字？”，他会回答“我名叫约翰”。这种情况经常发生，但简单地说“约翰”会更高效。您需要的只是一个名称，所以该句子的剩余部分都是多余的。

不幸的是，我们总是在代码中做这类多余的事情。Java 开发人员可以使用 `forEach` 迭代并输出某个范围内的每个值的双倍值，如下所示：

```
IntStream.rangeClosed(1, 5)
.forEach((int number) -> System.out.println(number * 2));

```

Show moreShow more icon

`rangeClosed` 方法生成一个从 1 到 5 的 `int` 值流。lambda 表达式的唯一职责就是接收一个名为 `number` 的 `int` 参数，使用 `PrintStream` 的 `println` 方法输出该值的双倍值。从语法上讲，该 lambda 表达式没有错，但类型细节有些冗余。

## Java 8 中的类型推断

当您从某个数字范围中提取一个值时，编译器知道该值的类型为 `int` 。不需要在代码中显式声明该值，尽管这是目前为止的约定。

在 Java 8 中，我们可以丢弃 lambda 表达式中的类型，如下所示：

```
IntStream.rangeClosed(1, 5)
.forEach((number) -> System.out.println(number * 2));

```

Show moreShow more icon

由于 Java 是静态类型语言，它需要在编译时知道所有对象和变量的类型。在 lambda 表达式的参数列表中省略类型并不会让 Java 更接近动态类型语言。但是，添加适当的类型推断功能会让 Java 更接近其他静态类型语言，比如 Scala 或 Haskell。

## 信任编译器

如果您在 lambda 表达式的一个参数中省略类型，Java 需要通过上下文细节来推断该类型。

返回到上一个示例，当我们在 `IntStream` 上调用 `forEach` 时，编译器会查找该方法来确定它采用的参数。 `IntStream` 的 `forEach` 方法期望使用函数接口 `IntConsumer` ，该接口的抽象方法 `accept` 采用了一个 `int` 类型的参数并返回 `void` 。

如果在参数列表中指定了该类型，编译器将会确认该类型符合预期。

如果省略该类型，编译器会推断出预期的类型 — 在本例中为 `int` 。

无论是您提供类型还是编译器推断出该类型，Java 都会在编译时知道 lambda 表达式参数的类型。要测试这种情况，可以在 lambda 表达式中引入一个错误，同时省略参数的类型：

```
IntStream.rangeClosed(1, 5)
.forEach((number) -> System.out.println(number.length() * 2));

```

Show moreShow more icon

编译此代码时，Java 编译器会返回以下错误：

```
Sample.java:7: error: int cannot be dereferenced
      .forEach((number) -> System.out.println(number.length() * 2));
                                                    ^
1 error

```

Show moreShow more icon

编译器知道名为 `number` 的参数的类型。它报错是因为它无法使用点运算符解除对某个 `int` 类型的变量的引用。可以对对象执行此操作，但不能对 `int` 变量这么做。

## 类型推断的好处

在 lambda 表达式中省略类型有两个主要好处：

- 键入的内容更少。无需输入类型信息，因为编译器自己能轻松确定该类型。
- 代码杂质更少 — `(number)` 比 `(int number)` 简单得多。

此外，一般来讲，如果我们仅有一个参数，省略类型意味着也可以省略 `()` ，如下所示：

```
IntStream.rangeClosed(1, 5)
.forEach(number -> System.out.println(number * 2));

```

Show moreShow more icon

请注意，您 _将_ 需要为采用多个参数的 lambda 表达式添加括号。

## 类型推断和可读性

lambda 表达式中的类型推断违背了 Java 中的常规做法，在常规做法中，会指定每个变量和参数的类型。尽管一些开发人员辩称 Java 指定类型的约定让代码变得更可读、更容易理解，但我认为这种偏好反映出一种习惯而不是必要性。

以一个包含一系列转换的函数管道为例：

```
List<String> result =
cars.stream()
    .map((Car c) -> c.getRegistration())
    .map((String s) -> DMVRecords.getOwner(s))
    .map((Person o) -> o.getName())
    .map((String s) -> s.toUpperCase())
    .collect(toList());

```

Show moreShow more icon

在这里，我们首先提供了一组 `Car` 实例和相关的注册信息。我们获取每辆车的车主和车主姓名，并将该姓名转换为大写。最后，将结果放入一个列表中。

这段代码中的每个 lambda 表达式都为其参数指定了一个类型，但我们为参数使用了单字母变量名。这在 Java 中很常见。但这种做法不合适，因为它丢弃了特定于域的上下文。

我们可以做得比这更好。让我们看看使用更强大的参数名重写代码后发生的情况：

```
List<String> result =
cars.stream()
    .map((Car car) -> car.getRegistration())
    .map((String registration) -> DMVRecords.getOwner(registration))
    .map((Person owner) -> owner.getName())
    .map((String name) -> name.toUpperCase())
    .collect(toList());

```

Show moreShow more icon

这些参数名包含了特定于域的信息。我们没有使用 `s` 来表示 `String` ，而是指定了特定于域的细节，比如 `registration` 和 `name` 。类似地，我们没有使用 `p` 或 `o` ，而是使用 `owner` 表明 `Person` 不只是一个人，还是这辆车的车主。

这个示例中的每个 lambda 表达式都比它所取代的表达式更好。在读取 lambda 表达式（例如 `(Person owner) -> owner.getName()` ）时，我们知道我们获得了车主的姓名，而不只是随便某个人的姓名。

### 命名参数

Scala 和 TypeScript 等一些语言更加重视参数名而不是类型。在 Scala 中，我们在定义类型之前定义参数，例如通过编写：

```
def getOwner(registration: String)

```

Show moreShow more icon

而不是：

```
def getOwner(String registration)

```

Show moreShow more icon

类型和参数名都很有用，但在 Scala 中，参数名更重要一些。我们用 Java 编写 lambda 表达式时，也可以考虑这一想法。请注意我们在 Java 中的车辆注册示例中丢弃类型细节和括号时发生的情况：

```
List<String> result =
cars.stream()
    .map(car -> car.getRegistration())
    .map(registration -> DMVRecords.getOwner(registration))
    .map(owner -> owner.getName())
    .map(name -> name.toUpperCase())
    .collect(toList());

```

Show moreShow more icon

因为我们添加了描述性的参数名，所以我们没有丢失太多上下文，而且显式类型（现在是冗余内容）已悄然消失。结果是我们获得了更干净、更朴实的代码。

## 类型推断的局限性

尽管使用类型推断可以提高效率和可读性，但这种技术并不适用于所有场合。在某些情况下，完全无法使用类型推断。幸运的是，您可以依靠 Java 编译器来获知何时出现这种情况。

我们首先看一个测试编译器并获得成功的示例，然后看一个测试失败的示例。最重要的是，在两种情况下，都能够相信编译器会按期望方式工作。

### 扩展类型推断

在我们的第一个示例中，假设我们想创建一个 `Comparator` 来比较 `Car` 实例。我们首先需要一个 `Car` 类：

```
class Car {
public String getRegistration() { return null; }
}

```

Show moreShow more icon

接下来，我们将创建一个 `Comparator` ，以便基于 `Car` 实例的注册信息对它们进行比较：

```
public static Comparator<Car> createComparator() {
return comparing((Car car) -> car.getRegistration());
}

```

Show moreShow more icon

用作 `comparing` 方法的参数的 lambda 表达式在其参数列表中包含了类型信息。我们知道 Java 编译器非常擅长类型推断，那么让我们看看在省略参数类型的情况下会发生什么，如下所示：

```
public static Comparator<Car> createComparator() {
return comparing(car -> car.getRegistration());
}

```

Show moreShow more icon

`comparing` 方法采用了 1 个参数。它期望使用 `Function<? super T, ? extends U>` 并返回 `Comparator<T>` 。因为 `comparing` 是 `Comparator<T>` 上的一个静态方法，所以编译器目前没有关于 `T` 或 `U` 可能是什么的线索。

为了解决此问题，编译器稍微扩展了推断范围，将范围扩大到传递给 `comparing` 方法的参数之外。它观察我们是如何处理调用 `comparing` 的结果的。根据此信息，编译器确定我们仅返回该结果。接下来，它看到由 `comparing` 返回的 `Comparator<T>` 又作为 `Comparator<Car>` 由 `createComparator` 返回 。

_注意了！_ 编译器现在已明白我们的意图：它推断应该将 `T` 绑定到 `Car` 。根据此信息，它知道 lambda 表达式中的 `car` 参数的类型应该为 `Car` 。

在这个例子中，编译器必须执行一些额外的工作来推断类型，但它成功了。接下来，让我们看看在提高挑战难度，让编译器达到其能力极限时，会发生什么。

### 推断的局限性

首先，我们在前一个 `comparing` 调用后面添加了一个新调用。在本例中，我们还为 lambda 表达式的参数重新引入显式类型：

```
public static Comparator<Car> createComparator() {
return comparing((Car car) -> car.getRegistration()).reversed();
}

```

Show moreShow more icon

借助显式类型，此代码没有编译问题，但现在让我们丢弃类型信息，看看会发生什么：

```
public static Comparator<Car> createComparator() {
return comparing(car -> car.getRegistration()).reversed();
}

```

Show moreShow more icon

如您下面所见，进展并不顺利。Java 编译器抛出了错误：

```
Sample.java:21: error: cannot find symbol
    return comparing(car -> car.getRegistration()).reversed();
                               ^
symbol:   method getRegistration()
location: variable car of type Object
Sample.java:21: error: incompatible types: Comparator<Object> cannot be converted to Comparator<Car>
    return comparing(car -> car.getRegistration()).reversed();
                                                           ^
2 errors

```

Show moreShow more icon

像上一个场景一样，在包含 `.reversed()` 之前，编译器会询问我们将如何处理调用 `comparing(car -> car.getRegistration())` 的结果。在上一个示例中，我们以 `Comparable<Car>` 形式返回结果，所以编译器能推断出 `T` 的类型为 `Car` 。

但在修改过后的版本中，我们将传递 `comparable` 的结果作为调用 `reversed()` 的目标。 `comparable` 返回 `Comparable<T>` ， `reversed()` 没有展示任何有关 `T` 的可能含义的额外信息。根据此信息，编译器推断 `T` 的类型肯定是 `Object` 。遗憾的是，此信息对于该代码而言并不足够，因为 `Object` 缺少我们在 lambda 表达式中调用的 `getRegistration()` 方法。

类型推断在这一刻失败了。在这种情况下，编译器实际上需要一些信息。类型推断会分析参数、返回元素或赋值元素来确定类型，但在上下文提供的细节不足时，编译器就会达到其能力极限。

### 能否采用方法引用作为补救措施？

在我们放弃这种特殊情况之前，让我们尝试另一种方法：不使用 lambda 表达式，而是尝试使用方法引用：

```
public static Comparator<Car> createComparator() {
return comparing(Car::getRegistration).reversed();
}

```

Show moreShow more icon

编译器对此解决方案非常满意。它在方法引用中使用 `Car::` 来推断类型。

## 结束语

Java 8 为 lambda 表达式的参数引入了有限的类型推断能力，在未来的 Java 版本中，会将类型推断扩展到局部变量。现在应该学会省略类型细节并信任编译器，这有助于您轻松步入未来的 Java 环境。

依靠类型推断和适当命名的参数，编写简明、更富于表达且更少杂质的代码。只要您相信编译器能自行推断出类型，就可以使用类型推断。仅在您确定编译器确实需要您的帮助的情况下提供类型细节。

本文翻译自： [Java knows your type](https://developer.ibm.com/articles/j-java8idioms8/)（2017-10-11）