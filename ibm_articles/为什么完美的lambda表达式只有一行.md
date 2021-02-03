# 为什么完美的 lambda 表达式只有一行
使代码更容易阅读、测试和重用的单行 lambda 表达式

**标签:** Java,Java 平台

[原文链接](https://developer.ibm.com/zh/articles/j-java8idioms6/)

Venkat Subramaniam

发布: 2017-08-16

* * *

**关于本系列**

Java 8 是自 Java 语言诞生以来进行的一次最重大更新  包含了非常丰富的新功能，您可能想知道从何处开始着手了解它。在 [本系列](/zh/series/java-8-idioms/) 中，作者兼教师 Venkat Subramaniam 提供了一种惯用的 Java 8 编程方法：这些简短的探索会激发您反思您认为理所当然的 Java 约定，同时逐步将新方法和语法集成到您的程序中。

目前您已在本系列中了解到，函数组合的一个主要好处是它会获得富于表达的代码。编写简短的 lambda 表达式是实现这一表达能力的关键，但通常说起来容易做起来难。本文会加深您目前对创建单行 lambda 表达式的各个方面的了解。通过学习函数组合的结构和好处，您很快就会掌握完美的 lambda 表达式， — 一个仅短短一行的表达式。

## 编写 lambda 表达式的两种方法

众所周知，lambda 表达式是匿名函数，它们天生就很简洁。普通的函数或方法通常有 4 个元素：

- 一个名称
- 返回类型
- 参数列表
- 主体

在这里可以看到，lambda 表达式只有这 4 元素中的最后两个：

```
(parameter list) -> body

```

Show moreShow more icon

`"->”` 将参数列表与函数主体分离，旨在对给定参数进行处理。函数的主体可能是一个表达式或一条语句。下面给出了一个示例：

```
(Integer e) -> e * 2

```

Show moreShow more icon

在此代码中，主体只有一行：一个返回给定参数两次的表达式。信噪比很高，没有分号，也不需要 `return` 关键字。这就是一个理想的 lambda 表达式。

### 多行 lambda 表达式

在 Java™ 中，lambda 表达式的主体也可能是一个复杂的表达式或声明语句；也就是说，一个 lambda 表达式包含多行。在这种情况下，分号必不可少。如果 lambda 表达式返回一个结果，也会需要 `return` 关键字。下面给出了一个示例：

```
(Integer e) -> {
double sqrt = Math.sqrt(e);
double log = Math.log(e);

return sqrt + log;
}

```

Show moreShow more icon

本示例中的 lambda 表达式返回了 `sqrt` 和给定参数的 `log` 的和。因为主体包含多行，所以括号 (`{}`)、分号 (`;`) 和 `return` 关键字都是必需的。

如果感觉好像 Java 因为我们编写多行 lambda 表达式而惩罚我们， — 或许我们应该接受这样的暗示。

## 函数组合的强大功能

函数式编码风格利用了函数组合的表达能力。比较两段代码时，很容易看出富于表达的好处。第一段代码是用命令式风格编写的：

```
int result = 0;
for(int e : values) {
if(e > 3 && e % 2 == 0) {
    result = e * 2;
    break;
}
}

```

Show moreShow more icon

现在考虑用函数式风格编写的相同代码：

```
int result = values.stream()
.filter(e -> e > 3)
.filter(e -> e % 2 == 0)
.map(e -> e * 2)
.findFirst()
.orElse(0);

```

Show moreShow more icon

两段代码获得了相同的结果。在命令式代码中，我们需要读入 `for` 循环，按照分支和中断来跟随流程。第二段代码使用了函数组合，更容易阅读一些。因为它是从上往下执行的，所以我们只需要传递该代码一次。

本质上，第二段代码读起来像是一个问题陈述： _给定一些值，仅选择大于 3 的值。从这些值中，仅选择偶数值，并将它们乘以 2。最后，挑选第一个结果。如果没有任何值存在，则返回 0。_

此代码不仅优雅，而且它的 _工作量并不比_ 命令式代码多。得益于 `Stream` 的惰性计算能力，这里没有浪费计算资源。

函数组合的表达能力很大程度上依赖于每个 lambda 表达式的简洁性。如果您的 lambda 表达式包含多行（甚至 [两行可能都太多](https://twitter.com/venkat_s/status/611119147586600960) ），您可能没有理解函数式编程的关键点。

## 充满危险的长 lambda 表达式

要更好地理解编写简短的 lambda 表达式的好处，可考虑反面情况：一个包含多行代码的杂乱 lambda 表达式：

```
System.out.println(
values.stream()
.mapToInt(e -> {
    int sum = 0;
    for(int i = 1; i <= e; i++) {
      if(e % i == 0) {
        sum += i;
      }
    }

    return sum;
})
.sum());

```

Show moreShow more icon

尽管此代码是用函数式风格编写的，但它丢失了函数式编程的优点。让我们考虑一下原因何在。

#### 1.难以读懂

好代码应该易于读懂。此代码需要绞尽脑汁才能读懂：很难找到不同部分的开头和结尾。

#### 2.用途不明

好代码读起来应该像一个故事，而不是像一个字谜。像这样冗长的、无特色的代码隐藏了它的具体用途，会耗费读者的时间和精力。将这段代码包装在一个命名函数中可以使其模块化，同时也可以通过相关的名称揭示它的用途。

#### 3.代码质量差

无论您的代码有何用途，您可能都希望在某个时候重用它。这段代码的逻辑已嵌入在 lambda 表达式中，后者又以参数形式传递给另一个函数 `mapToInt` 。如果我们在程序的其他某个地方需要该代码，我们可能忍不住重写它，这会引起代码库中的不一致性。或者，我们也可以复制并粘贴该代码。两种选项都会得到好代码或高品质的软件。

#### 4.难以测试

代码始终依靠键入的内容进行操作，而且不一定是我们打算执行的操作，所以这代表着必须测试任何非平凡代码。如果 lambda 表达式中的代码无法用作一个单元，则无法对它执行单元测试。您可以运行集成测试，但这无法取代单元测试，尤其是在代码执行重要工作时。

#### 5.代码覆盖范围小

一位学习 Java 8 课程的学员最近表示他们讨厌 lambda 表达式。在我问为什么时，他们向我展示了一位同事的作品，其中包含的 lambda 表达式运行了数百行代码。嵌入在参数中的 Lambda 表达式无法轻松地作为单元提取出来，而且许多表达式在覆盖范围报告中显示为红色。由于一无所知，该团队很难假设这些代码在工作正常。

## 使用 Lambda 作为粘合代码

解决所有这些问题的方法是让您的 lambda 表达式高度简洁。作为第一个且非常有用的步骤，避免在 lambda 表达式中使用括号。考虑如何使用此技术轻松地重写前面杂乱的 lambda 表达式：

```
System.out.println(
values.stream()
.mapToInt(e -> sumOfFactors(e))
.sum());

```

Show moreShow more icon

此代码很简洁，尽管它不完整。该代码也具有很高的可读性。它的用途是： _给定一组值，将该列表转换为每个数的因数之和，然后计算所获得的集合的和。_ 显式命名以前包含在括号中的代码主体，这会让此代码变得更容易阅读和理解。函数管道现在既整洁又容易理解。

不需要解字谜，因为此代码直接表明了它的用途。计算因数之和的代码已模块化为一个名为 `sumOfFactors` 的单独方法，该方法可以重用。因为它是一个单独方法，所以对它的逻辑执行单元测试也很容易。因为此代码如此容易测试，所以您可以确保良好的代码覆盖范围。

简言之，曾经杂乱的 lambda 表达式现在成为了 _粘合代码_ — 它没有承担大量责任，只是将命名函数粘合到 `mapToInt` 函数。

### 使用方法引用进行调优

就像在 [本系列的其他地方](https://www.ibm.com/developerworks/cn/java/j-java8idioms5/index.html) 看到的一样，可以通过将 lambda 表达式替换为方法引用，让上述代码更富于表达（其中 `sumOfFactors` 是一个名为 `Sample` 的类的方法）：

```
System.out.println(
values.stream()
.mapToInt(Sample::sumOfFactors)
.sum());

```

Show moreShow more icon

这是重写后的 `sumOfFactors` 方法：

```
public static int sumOfFactors(int number) {
return IntStream.rangeClosed(1, number)
    .filter(i -> number % i == 0)
    .sum();
}

```

Show moreShow more icon

现在它是一个简短的方法。该方法中的 lambda 表达式也很简洁：只有一行，没有过多的繁杂过程或噪音。

## 结束语

简短的 lambda 表达式能提高代码可读性，这是函数式编程的重要好处之一。包含多行的 lambda 表达式具有相反的效果，会让代码变得杂乱且难以阅读。多行 lambda 表达式还难以测试和重用，这可能导致重复工作和代码质量差。幸运的是，通过将多行 lambda 表达式的主体转移到一个命名函数中，然后从 lambda 表达式内调用该函数，这样很容易避免这些问题。我也推荐尽可能将 lambda 表达式替换为方法引用。

简言之，我推荐避免多行 lambda 表达式，除非是为了演示它们为什么不好。

本文翻译自： [Why the perfect lambda expression is just one line](https://developer.ibm.com/articles/j-java8idioms6/)（2017-08-02）