# 传统 for 循环的函数式替代方案
3 个消除复杂迭代中的麻烦的新方法

**标签:** Java,Java 平台

[原文链接](https://developer.ibm.com/zh/articles/j-java8idioms3/)

Venkat Subramaniam

发布: 2017-06-05

* * *

尽管 `for` 循环包含许多可变部分，但许多开发人员仍非常熟悉它，并会不假思索地使用它。从 Java™ 8 开始，我们有多个强大的新方法可帮助简化复杂迭代。在本文中，您将了解如何使用 `IntStream` 方法 `range`、`iterate` 和 `limit` 来迭代范围和跳过范围中的值。您还将了解新的 `takeWhile` 和 `dropWhile` 方法（即将在 Java 9 中引入）。

**关于本系列**

Java 8 是自 Java 语言诞生以来进行的一次最重大更新  包含了非常丰富的新功能，您可能想知道从何处开始着手了解它。在 [本系列](/zh/series/java-8-idioms/) 中，作者兼教师 Venkat Subramaniam 提供了一种惯用的 Java 8 编程方法：这些简短的探索会激发您反思您认为理所当然的 Java 约定，同时逐步将新方法和语法集成到您的程序中。

## for 循环的麻烦

在 Java 语言的第 1 个版本中就开始引入了传统的 `for` 循环，它的更简单的变体 `for-each` 是在 Java 5 中引入的。大部分开发人员更喜欢使用 `for-each` 执行日常迭代，但对于迭代一个范围或跳过范围中的值等操作，他们仍会使用 `for` 。

`for` 循环非常强大，但它包含太多可变部分。甚至在打印 `get set` 提示的最简单任务中，也可以看出这一点：

##### 清单 1\. 完成一个简单任务的复杂代码

```
System.out.print("Get set...");
for(int i = 1; i < 4; i++) {
    System.out.print(i + "...");
}

```

Show moreShow more icon

在清单 1 中，我们从 1 开始循环处理索引变量 `i` ，将它限制到小于 4 的值。请注意， `for` 循环需要我们告诉循环是递增的。在本例中，我们还选择了前递增而不是后递增。

清单 1 中没有太多代码，但比较繁琐。Java 8 提供了一种更简单、更优雅的替代方法： `IntStream` 的 `range` 方法。以下是打印清单 1 中的相同 `get set` 提示的 `range` 方法：

##### 清单 2\. 完成一个简单任务的简单代码

```
System.out.print("Get set...");
IntStream.range(1, 4)
    .forEach(i -> System.out.print(i + "..."));

```

Show moreShow more icon

在清单 2 中，我们看到并没有显著减少代码量，但降低了它的复杂性。这样做有两个重要原因：

1. 不同于 `for` ， `range` 不会强迫我们初始化某个可变变量。
2. 迭代会自动执行，所以我们不需要像循环索引一样定义增量。

在语义上，最初的 `for` 循环中的变量 `i` 是一个 _可变变量_ 。理解 `range` 和类似方法的价值对理解该设计的结果很有帮助。

## 可变变量与参数

`for` 循环中定义的变量 `i` 是单个变量，它会在每次对循环执行迭代时发生改变。 `range` 示例中的变量 `i` 是拉姆达表达式的参数，所以它在每次迭代中都是一个全新的变量。这是一个细微区别，但决定了两种方法的不同。以下示例有助于阐明这一点。

清单 3 中的 `for` 循环想在一个内部类中使用索引变量：

##### 清单 3\. 在内部类中使用索引变量

```
ExecutorService executorService = Executors.newFixedThreadPool(10);

      for(int i = 0; i < 5; i++) {
        int temp = i;

        executorService.submit(new Runnable() {
          public void run() {
            //If uncommented the next line will result in an error
            //System.out.println("Running task " + i);
            //local variables referenced from an inner class must be final or effectively final

            System.out.println("Running task " + temp);
          }
        });
      }

      executorService.shutdown();

```

Show moreShow more icon

我们有一个匿名的内部类实现了 `Runnable` 接口。我们想在 `run` 方法中访问索引变量 `i` ，但编译器不允许这么做。

作为此限制的解决办法，我们可以创建一个局部临时变量，比如 `temp` ，它是索引变量的一个副本。每次新的迭代都会创建变量 `temp` 。在 Java 8 以前，我们需要将该变量标记为 `final` 。从 Java 8 开始，可以将它视为实际的最终结果，因为我们不会再更改它。无论如何，由于事实上索引变量是一个在迭代中改变的变量， `for` 循环中就会出现这个额外变量。

现在尝试使用 `range` 函数解决同一个问题。

##### 清单 4\. 在内部类中使用拉姆达参数

```
ExecutorService executorService = Executors.newFixedThreadPool(10);

      IntStream.range(0, 5)
        .forEach(i ->
          executorService.submit(new Runnable() {
            public void run() {
              System.out.println("Running task " + i);
            }
          }));

      executorService.shutdown();

```

Show moreShow more icon

在作为一个参数被拉姆达表达式接受后，索引变量 `i` 的语义与循环索引变量有所不同。与清单 3 中手动创建的 `temp` 非常相似，这个 `i` 参数在每次迭代中都表现为一个全新的变量。它是 _实际最终变量_ ，因为我们不会在任何地方更改它的值。因此，我们可以直接在内部类的上下文中使用它 — 且不会有任何麻烦。

因为 `Runnable` 是一个函数接口，所以我们可以轻松地将匿名的内部类替换为拉姆达表达式，比如：

##### 清单 5\. 将内部类替换为拉姆达表达式

```
IntStream.range(0, 5)
       .forEach(i ->
         executorService.submit(() -> System.out.println("Running task " + i)));

```

Show moreShow more icon

显然，对于相对简单的迭代，使用 `range` 代替 `for` 具有一定优势，但 `for` 的特殊价值体现在于它能处理更复杂的迭代场景。让我们看看 `range` 和其他 Java 8 方法孰优孰劣。

## 封闭范围

创建 `for` 循环时，可以将索引变量封闭在一个范围内，比如：

##### 清单 6\. 一个具有封闭范围的 for 循环

```
for(int i = 0; i <= 5; i++) {

```

Show moreShow more icon

索引变量 `i` 接受值 `0` 、 `1` 、…… `5` 。无需使用 `for` ，我们可以使用 `rangeClosed` 方法。在本例中，我们告诉 `IntStream` 将最后一个值限制在该范围内：

##### 清单 7\. rangeClosed 方法

```
IntStream.rangeClosed(0, 5)

```

Show moreShow more icon

迭代此范围时，我们会获得包含边界值 5 在内的值。

## 跳过值

对于基本循环， `range` 和 `rangeClosed` 方法是 `for` 的更简单、更优雅的替代方法，但是如果想跳过一些值该怎么办？在这种情况下， `for` 对前期工作的需求使该运算变得非常容易。在清单 8 中， `for` 循环在迭代期间快速跳过两个值：

##### 清单 8\. 使用 for 跳过值

```
int total = 0;
for(int i = 1; i <= 100; i = i + 3) {
total += i;
}

```

Show moreShow more icon

清单 8 中的循环在 1 到 100 内对每次读到的第三个值作求和计算 — 这种复杂运算可使用 `for` 轻松完成。能否也使用 `range` 解决此问题？

首先，可以考虑使用 `IntStream` 的 `range` 方法，再结合使用 `filter` 或 `map` 。但是，所涉及的工作比使用 `for` 循环要多。一种更可行的解决方案是结合使用 `iterate` 和 `limit` ：

##### 清单 9\. 使用 limit 的迭代

```
IntStream.iterate(1, e -> e + 3)
.limit(34)
.sum()

```

Show moreShow more icon

`iterate` 方法很容易使用；它只需获取一个初始值即可开始迭代。作为第二参数传入的拉姆达表达式决定了迭代中的下一个值。这类似于清单 8，我们将一个表达式传递给 `for` 循环来递增索引变量的值。但是，在本例中有一个 _陷阱_ 。不同于 `range` 和 `rangeClosed` ，没有参数来告诉 `iterate` 方法何时停止迭代。如果我们没有限制该值，迭代会一直进行下去。

如何解决这个问题？

我们对 1 到 100 之间的值感兴趣，而且想从 1 开始跳过两个值。稍加运算，即可确定给定范围中有 34 个符合要求的值。所以我们将该数字传递给 `limit` 方法。

此代码很有效，但过程太复杂：提前执行数学运算不那么有趣，而且它限制了我们的代码。如果我们决定跳过 3 个值而不是 2 个值，该怎么办？我们不仅需要更改代码，结果也很容易出错。我们需要有一个更好的方法。

## takeWhile 方法

Java 9 中即将引入的 `takeWhile` 是一个新方法，它使得执行有限制的迭代变得更容易。使用 `takeWhile` ，可以直接表明 _只要满足想要的条件_ ，迭代就应该继续执行。以下是使用 `takeWhile` 实现清单 9 中的迭代的代码。

##### 清单 10\. 有条件的迭代

```
IntStream.iterate(1, e -> e + 3)
     .takeWhile(i -> i <= 100) //available in Java 9
     .sum()

```

Show moreShow more icon

无需将迭代限制到预先计算的次数，我们使用提供给 `takeWhile` 的条件，动态确定何时终止迭代。与尝试预先计算迭代次数相比，这种方法简单得多，而且更不容易出错。

与 `takeWhile` 方法相反的是 `dropWhile` ，它跳过满足给定条件前的值，这两个方法都是 JDK 中非常需要的补充方法。 `takeWhile` 方法类似于 break，而 `dropWhile` 则类似于 continue。从 Java 9 开始，它们将可用于任何类型的 `Stream` 。

## 逆向迭代

与正向迭代相比，逆向迭代同样非常简单，无论使用传统的 `for` 循环还是 `IntStream` 。

以下是一个逆向的 `for` 循环迭代：

##### 清单 11\. 使用 for 的逆向迭代

```
for(int i = 7; i > 0; i--) {

```

Show moreShow more icon

`range` 或 `rangeClosed` 中的第一个参数不能大于第二个参数，所以我们无法使用这两种方法来执行逆向迭代。但可以使用 `iterate` 方法：

##### 清单 12\. 使用 iterate 的逆向迭代

```
IntStream.iterate(7, e -> e - 1)
     .limit(7)

```

Show moreShow more icon

将一个拉姆达表达式作为参数传递给 `iterate` 方法，该方法对给定值进行递减，以便沿相反方向执行迭代。我们使用 `limit` 函数指定我们希望在逆向迭代期间看到总共多少个值。如有必要，还可以使用 `takeWhile` 和 `dropWhile` 方法来动态调整迭代流。

## 结束语

尽管传统 `for` 循环非常强大，但它有些过于复杂。Java 8 和 Java 9 中的新方法可帮助简化迭代，甚至是简化复杂的迭代。方法 `range` 、 `iterate` 和 `limit` 的可变部分较少，这有助于提高代码效率。这些方法还满足了 Java 的一个长期以来的要求，那就是局部变量必须声明为 final，然后才能从内部类访问它。将一个可变索引变量更换为实际的 final 参数只有很小的语义差别，但它减少了大量垃圾变量。最终您会得到更简单、更优雅的代码。

本文翻译自： [Functional alternatives to the traditional for loop](https://developer.ibm.com/articles/j-java8idioms3/)（2017-04-30）