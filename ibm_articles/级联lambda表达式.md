# 级联 lambda 表达式
可重用的函数有助于让代码变得非常简短，但是会不会过于简短呢？

**标签:** Java,Java 平台

[原文链接](https://developer.ibm.com/zh/articles/j-java8idioms9/)

Venkat Subramaniam

发布: 2017-11-29

* * *

**关于本系列**

Java 8 是自 Java 语言诞生以来进行的一次最重大更新  包含了非常丰富的新功能，您可能想知道从何处开始着手了解它。在 [本系列](/zh/series/java-8-idioms/) 中，作者兼教师 Venkat Subramaniam 提供了一种惯用的 Java 8 编程方法：这些简短的探索会激发您反思您认为理所当然的 Java 约定，同时逐步将新方法和语法集成到您的程序中。

在函数式编程中，函数既可以接收也可以返回其他函数。函数不再像传统的面向对象编程中一样，只是一个对象的 _工厂_ 或 _生成器_ ，它也能够创建和返回另一个函数。返回函数的函数可以变成 _级联 lambda 表达式_ ，特别值得注意的是代码非常简短。尽管此语法初看起来可能非常陌生，但它有自己的用途。本文将帮助您认识级联 lambda 表达式，理解它们的性质和在代码中的用途。

## 神秘的语法

您是否看到过类似这样的代码段？

```
x -> y -> x > y

```

Show moreShow more icon

如果您很好奇”这到底是什么意思？”，那么您并不孤单。对于不熟悉使用 lambda 表达式编程的开发人员，此语法可能看起来像货物正从快速行驶的卡车上一件件掉下来一样。

幸运的是，我们不会经常看到它们，但理解如何创建级联 lambda 表达式和如何在代码中理解它们会大大减少您的受挫感。

## 高阶函数

在谈论级联 lambda 表达式之前，有必要首先理解如何创建它们。对此，我们需要回顾一下高阶函数（已在 [本系列第 1 篇文章](/zh/articles/j-java8idioms1/) 中介绍）和它们在 _函数分解_ 中的作用，函数分解是一种将复杂流程分解为更小、更简单的部分的方式。

首先，考虑区分高阶函数与常规函数的规则：

常规函数

- 可以接收对象
- 可以创建对象
- 可以返回对象

高阶函数

- 可以接收函数
- 可以创建函数
- 可以返回函数

开发人员将匿名函数或 lambda 表达式传递给高阶函数，以让代码简短且富于表达。让我们看看这些高阶函数的两个示例。

### 示例 1：一个接收函数的函数

在 Java™ 中，我们使用函数接口来引用 lambda 表达式和方法引用。下面这个函数接收一个对象和一个函数：

```
public static int totalSelectedValues(List<Integer> values,
Predicate<Integer> selector) {

return values.stream()
    .filter(selector)
    .reduce(0, Integer::sum);
}

```

Show moreShow more icon

`totalSelectedValues` 的第一个参数是集合对象，而第二个参数是 `Predicate` 函数接口。 因为参数类型是函数接口 (`Predicate`)，所以我们现在可以将一个 lambda 表达式作为第二个参数传递给 `totalSelectedValues` 。例如，如果我们想仅对一个 `numbers` 列表中的 _偶数值_ 求和，可以调用 `totalSelectedValues` ，如下所示：

```
totalSelectedValues(numbers, e -> e % 2 == 0);

```

Show moreShow more icon

假设我们现在在 `Util` 类中有一个名为 `isEven` 的 `static` 方法。在此情况下，我们可以使用 `isEven` 作为 `totalSelectedValues` 的参数，而不传递 lambda 表达式：

```
totalSelectedValues(numbers, Util::isEven);

```

Show moreShow more icon

作为规则，只要一个函数接口显示为一个函数的参数的 _类型_ ，您看到的就是一个高阶函数。

### 示例 2：一个返回函数的函数

函数可以接收函数、lambda 表达式或方法引用作为参数。同样地，函数也可以返回 lambda 表达式或方法引用。在此情况下，返回类型将是函数接口。

让我们首先看一个创建并返回 `Predicate` 来验证给定值是否为奇数的函数：

```
public static Predicate<Integer> createIsOdd() {
Predicate<Integer> check = (Integer number) -> number % 2 != 0;
return check;
}

```

Show moreShow more icon

为了返回一个函数，我们必须提供一个函数接口作为返回类型。在本例中，我们的函数接口是 `Predicate` 。尽管上述代码在语法上是正确的，但它可以更加简短。 我们使用类型引用并删除临时变量来改进该代码：

```
public static Predicate<Integer> createIsOdd() {
return number -> number % 2 != 0;
}

```

Show moreShow more icon

这是使用的 `createIsOdd` 方法的一个示例：

```
Predicate<Integer> isOdd = createIsOdd();

isOdd.test(4);

```

Show moreShow more icon

请注意，在 `isOdd` 上调用 `test` 会返回 `false` 。我们也可以在 `isOdd` 上使用更多值来调用 `test` ；它并不限于使用一次。

## 创建可重用的函数

现在您已大体了解高阶函数和如何在代码中找到它们，我们可以考虑使用它们来让代码更加简短。

设想我们有两个列表 `numbers1` 和 `numbers2` 。假设我们想从第一个列表中仅提取大于 50 的数，然后从第二个列表中提取大于 50 的值并 _乘以 2_ 。

可通过以下代码实现这些目的：

```
List<Integer> result1 = numbers1.stream()
.filter(e -> e > 50)
.collect(toList());

List<Integer> result2 = numbers2.stream()
.filter(e -> e > 50)
.map(e -> e * 2)
.collect(toList());

```

Show moreShow more icon

此代码很好，但您注意到它很冗长了吗？我们对检查数字是否大于 50 的 lambda 表达式使用了两次。 我们可以通过创建并重用一个 `Predicate` ，从而删除重复代码，让代码更富于表达：

```
Predicate<Integer> isGreaterThan50 = number -> number > 50;

List<Integer> result1 = numbers1.stream()
.filter(isGreaterThan50)
.collect(toList());

List<Integer> result2 = numbers2.stream()
.filter(isGreaterThan50)
.map(e -> e * 2)
.collect(toList());

```

Show moreShow more icon

通过将 lambda 表达式存储在一个引用中，我们可以重用它，这是我们避免重复 lambda 表达式的方式。如果我们想跨方法重用 lambda 表达式，也可以将该引用放入一个单独的方法中，而不是放在一个局部变量引用中。

现在假设我们想从列表 `numbers1` 中提取大于 25、50 和 75 的值。我们可以首先编写 3 个不同的 lambda 表达式：

```
List<Integer> valuesOver25 = numbers1.stream()
.filter(e -> e > 25)
.collect(toList());

List<Integer> valuesOver50 = numbers1.stream()
.filter(e -> e > 50)
.collect(toList());

List<Integer> valuesOver75 = numbers1.stream()
.filter(e -> e > 75)
.collect(toList());

```

Show moreShow more icon

尽管上面每个 lambda 表达式将输入与一个不同的值比较，但它们做的事情完全相同。如何以较少的重复来重写此代码？

## 创建和重用 lambda 表达式

尽管上一个示例中的两个 lambda 表达式相同，但上面 3 个表达式稍微不同。创建一个返回 `Predicate` 的 `Function` 可以解决此问题。

首先，函数接口 `Function<T, U>` 将一个 `T` 类型的输入转换为 `U` 类型的输出。例如，下面的示例将一个给定值转换为它的平方根：

```
Function<Integer, Double> sqrt = value -> Math.sqrt(value);

```

Show moreShow more icon

在这里，返回类型 `U` 可以很简单，比如 `Double` 、 `String` 或 `Person` 。或者它也可以更复杂，比如 `Consumer` 或 `Predicate` 等另一个函数接口。

在本例中，我们希望一个 `Function` 创建一个 `Predicate` 。所以代码如下：

```
Function<Integer, Predicate<Integer>> isGreaterThan = (Integer pivot) -> {
Predicate<Integer> isGreaterThanPivot = (Integer candidate) -> {
    return candidate > pivot;
};

return isGreaterThanPivot;
};

```

Show moreShow more icon

引用 `isGreaterThan` 引用了一个表示 `Function<T, U>` — 或更准确地讲表示 `Function<Integer, Predicate<Integer>>` 的 lambda 表达式。输入是一个 `Integer` ，输出是一个 `Predicate<Integer>` 。

在 lambda 表达式的主体中（外部 `{}` 内），我们创建了另一个引用 `isGreaterThanPivot` ，它包含对另一个 lambda 表达式的引用。这一次，该引用是一个 `Predicate` 而不是 `Function` 。最后，我们返回该引用。

`isGreaterThan` 是一个 lambda 表达式的引用，该表达式在调用时返回 _另一个_ lambda 表达式 — 换言之，这里隐藏着一种 lambda 表达式级联关系。

现在，我们可以使用新创建的外部 lamba 表达式来解决代码中的重复问题：

```
List<Integer> valuesOver25 = numbers1.stream()
.filter(isGreaterThan.apply(25))
.collect(toList());

List<Integer> valuesOver50 = numbers1.stream()
.filter(isGreaterThan.apply(50))
.collect(toList());

List<Integer> valuesOver75 = numbers1.stream()
.filter(isGreaterThan.apply(75))
.collect(toList());

```

Show moreShow more icon

在 `isGreaterThan` 上调用 `apply` 会返回一个 `Predicate` ，后者然后作为参数传递给 `filter` 方法。

尽管整个过程非常简单（作为示例），但是能够抽象为一个函数对于谓词更加复杂的场景来说尤其有用。

## 保持简短的秘诀

我们已从代码中成功删除了重复的 lambda 表达式，但 `isGreaterThan` 的定义看起来仍然很杂乱。幸运的是，我们可以组合一些 Java 8 约定来减少杂乱，让代码更简短。

我们首先重构以下代码：

```
Function<Integer, Predicate<Integer>> isGreaterThan = (Integer pivot) -> {
Predicate<Integer> isGreaterThanPivot = (Integer candidate) -> {
    return candidate > pivot;
};

return isGreaterThanPivot;
};

```

Show moreShow more icon

可以使用类型引用来从外部和内部 lambda 表达式的参数中删除类型细节：

```
Function<Integer, Predicate<Integer>> isGreaterThan = (pivot) -> {
Predicate<Integer> isGreaterThanPivot = (candidate) -> {
    return candidate > pivot;
};

return isGreaterThanPivot;
};

```

Show moreShow more icon

目前，我们从代码中删除了两个单词，改进不大。

接下来，我们删除多余的 `()`，以及外部 lambda 表达式中不必要的临时引用：

```
Function<Integer, Predicate<Integer>> isGreaterThan = pivot -> {
return candidate -> {
    return candidate > pivot;
};
};

```

Show moreShow more icon

代码更加简短了，但是仍然看起来有些杂乱。

可以看到内部 lambda 表达式的主体只有一行，显然 `{}` 和 `return` 是多余的。让我们删除它们：

```
Function<Integer, Predicate<Integer>> isGreaterThan = pivot -> {
return candidate -> candidate > pivot;
};

```

Show moreShow more icon

现在可以看到，外部 lambda 表达式的主体 _也_ 只有一行，所以 `{}` 和 `return` 在这里也是多余的。在这里，我们应用最后一次重构：

```
Function<Integer, Predicate<Integer>> isGreaterThan =
pivot -> candidate -> candidate > pivot;

```

Show moreShow more icon

现在可以看到 — 这是我们的级联 lambda 表达式。

## 理解级联 lambda 表达式

我们通过一个适合每个阶段的重构过程，得到了最终的代码 – 级联 lambda 表达式。在本例中，外部 lambda 表达式接收 `pivot` 作为参数，内部 lambda 表达式接收 `candidate` 作为参数。内部 lambda 表达式的主体同时使用它收到的参数 (`candidate`) 和来自外部范围的参数。也就是说，内部 lambda 表达式的主体同时依靠它的参数和它的 _词法范围_ 或 _定义范围_ 。

级联 lambda 表达式对于编写它的人非常有意义。但是对于读者呢？

看到一个只有一个向右箭头 (`->`) 的 lambda 表达式时，您应该知道您看到的是一个匿名函数，它接受参数（可能是空的）并执行一个操作或返回一个结果值。

看到一个包含两个向右箭头 (`->`) 的 lambda 表达式时，您看到的也是一个匿名函数，但它接受参数（可能是空的）并返回另一个 lambda 表达式。返回的 lambda 表达式可以接受它自己的参数或者可能是空的。它可以执行一个操作或返回一个值。它甚至可以返回另一个 lambda 表达式，但这通常有点大材小用，最好避免。

大体上讲，当您看到两个向右箭头时，可以将第一个箭头右侧的所有内容视为一个黑盒：一个由外部 lambda 表达式返回的 lambda 表达式。

## 结束语

级联 lambda 表达式不是很常见，但您应该知道如何在代码中识别和理解它们。当一个 lambda 表达式返回另一个 lambda 表达式，而不是接受一个操作或返回一个值时，您将看到两个箭头。这种代码非常简短，但可能在最初遇到时非常难以理解。但是，一旦您学会识别这种函数式语法，理解和掌握它就会变得容易得多。

本文翻译自： [Cascading lambdas](https://developer.ibm.com/articles/j-java8idioms9/)（2017-11-07）