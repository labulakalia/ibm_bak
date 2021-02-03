# Java 中的一种更轻松的函数式编程途径
以声明方式思考如何在 Java 程序中采用函数方法

**标签:** Java,Java 平台

[原文链接](https://developer.ibm.com/zh/articles/j-java8idioms1/)

Venkat Subramaniam

更新: 2017-06-13 \| 发布: 2017-05-05

* * *

Java 开发人员早已习惯了命令式编程和面向对象的编程，因为 Java 语言从第一个版本开始就支持这些格式。在 Java 8 中，我们获得了一组强大的新的函数特性和语法。函数式编程已有数十年的历史，而且与面向对象的编程相比，函数式编程通常更简洁、更具表达力、更不容易出错，而且更容易并行化。所以在 Java 程序中引入函数特性是有充分理由的。尽管如此，函数式的编程需要对代码的设计方式进行一些改变。

**关于本系列**

Java 8 是自 Java 语言诞生以来进行的一次最重大更新  包含了非常丰富的新功能，您可能想知道从何处开始着手了解它。在 [本系列](/zh/series/java-8-idioms/) 中，作者兼教师 Venkat Subramaniam 提供了一种惯用的 Java 8 编程方法：这些简短的探索会激发您反思您认为理所当然的 Java 约定，同时逐步将新方法和语法集成到您的程序中。

我发现声明式思考（而不是命令式思考）可以简化向更加函数化编程的过渡。在新的 Java 8 习惯用语系列的第一篇文章中，我将解释命令式、声明式和函数式编程之间的区别和共性。然后，将展示如何使用声明式思考方式逐步将函数方法集成到 Java 程序中。

## 命令式格式

经过针对命令式编程培训的开发人员已经习惯了告诉程序做什么和该如何做。这是一个简单示例：

##### findNemo in the imperative style

```
import java.util.*;

public class FindNemo {
public static void main(String[] args) {
    List<String> names =
      Arrays.asList("Dory", "Gill", "Bruce", "Nemo", "Darla", "Marlin", "Jacques");

    findNemo(names);
}

public static void findNemo(List<String> names) {
    boolean found = false;
    for(String name : names) {
      if(name.equals("Nemo")) {
        found = true;
        break;
      }
    }

    if(found)
      System.out.println("Found Nemo");
    else
      System.out.println("Sorry, Nemo not found");
}
}

```

Show moreShow more icon

`findNemo()` 方法首先初始化一个可变的 _flag_ 变量，也称为 _垃圾（garbage）变量_ 。开发人员通常会随便给这些变量命名，比如 `f`、 `t` 或 `temp` ，这些名称表明了我们对这些变量不应存在的一般态度。在本例中，该变量名为 `found` 。

接下来，该程序会在给定的 `names` 列表中循环，一次处理一个元素。它检查获得的名称是否与它寻找的值（在本例中为 `Nemo`）相等。如果值匹配，那么它会将 flag 变量设置为 `true`，并通知控制流”跳出”循环。

因为这是一个命令式编程的程序 — 许多 Java 开发人员最熟悉的格式 — 所以您需要定义程序的每一步：告诉它迭代每个元素，比较值，设置 flag 变量，然后跳出循环。命令式格式为您提供了完全的控制权，这有时是件好事。而另一方面，您需要执行所有工作。在许多情况下，可以减少工作量来提高效率。

## 声明式格式

声明式编程意味着，您仍会告诉程序要做什么，但将实现细节留给底层的函数库。让我们看看用声明式格式重新编写来自 [findNemo in the imperative style](#findnemo-in-the-imperative-style) 的 `findNemo` 方法时，会发生什么：

##### findNemo in the declarative style

```
public static void findNemo(List<String> names) {
if(names.contains("Nemo"))
    System.out.println("Found Nemo");
else
    System.out.println("Sorry, Nemo not found");
}

```

Show moreShow more icon

首先请注意，这个版本中没有任何垃圾变量。您也没有将精力浪费在对集合的循环处理上，而是使用内置的 `contains()` 方法来完成工作。您仍然需要告诉程序要做什么 — 检查集合是否包含我们寻找的值 — 但将细节留给底层的方法。

在命令式示例中，您控制着迭代，而且程序完全按照要求来操作。在声明式版本中，您无需关心工作如何完成，只要它完成即可。 `contains()` 的实现可能不同，但只要结果符合预期，您可能就会很开心。花更少的精力获得同样的结果。

训练自己以声明式编程思考，这将大大简化向 Java 中的函数式编程的过渡。为什么呢？因为函数式编程是以声明式为基础而建立的。声明式思考试图逐步从命令式编程过渡到函数式编程。

## 函数式格式

尽管函数式格式的编程始终是声明式的，但简单地使用声明式编程并不等于函数式编程。这是因为函数式编程合并了声明式方法与高阶函数。图 1 直观地展示了命令式、声明式和函数式编程之间的关系。

##### 命令式、声明式和函数式编程的联系

### Java 中的高阶函数

在 Java 中，要将对象传递给方法，在方法内创建对象，并从方法中返回对象。可以对函数执行同样的操作。也就是说，可以将函数传递给方法，在方法内创建函数，并从方法返回函数。

在此上下文中， _方法_ 是类的一部分 — 静态或实例 — 但函数对于方法而言是本地函数，不能特意与类或实例关联。可以接收、创建或返回函数的函数或方法被视为 _高阶函数_ 。

## 函数式编程示例

采用新编程的格式需要改变您思考程序的方式。可以通过简单的示例来实践这一过程，并逐步建立更复杂的程序。

##### A Map in the imperative style

```
import java.util.*;

public class UseMap {
public static void main(String[] args) {
    Map<String, Integer> pageVisits = new HashMap<>();

    String page = "https://agiledeveloper.com";

    incrementPageVisit(pageVisits, page);
    incrementPageVisit(pageVisits, page);

    System.out.println(pageVisits.get(page));
}

public static void incrementPageVisit(Map<String, Integer> pageVisits, String page) {
    if(!pageVisits.containsKey(page)) {
       pageVisits.put(page, 0);
    }

    pageVisits.put(page, pageVisits.get(page) + 1);
}
}

```

Show moreShow more icon

在 [A Map in the imperative style](#a-map-in-the-imperative-style) 中， `main()` 函数创建了一个 `HashMap` ，其中包含一个网站的页面访问次数。同时，每次访问给定页面， `incrementPageVisit()` 方法都会增加计数。我们将重点查看此方法。

`incrementPageVisit()` 方法是使用命令式格式编写的：它的职责是递增给定页面的计数，将该计数存储在 `Map` 中。该方法不知道给定页面是否有计数，所以它首先会检查是否存在计数。如果不存在，那么它会插入一个 “0” 作为该页面的计数。然后获得该计数，递增它，并将新值存储在 `Map` 中。

声明式思考要求您将此方法的设计思路从 “如何做” 转变为 “做什么”。当调用方法 `incrementPageVisit()` 时，您希望将给定页面的计数初始化为 1 或将运行值递增 1。这就是 _做什么_ 。

因为您在执行声明式编程，所以下一步是扫描 JDK 库，查找 `Map` 接口中可以完成您的目标的方法 — 也就是说，寻找一个知道 _如何_ 完成给定任务的内置方法。

事实证明， `merge()` 方法能完美实现您的目标。清单 4 使用新的声明式方法修改来自 [A Map in the imperative style](#a-map-in-the-imperative-style) 的 `incrementPageVisit()` 方法。但是，在本例中，并不仅仅通过选择更智能的方法采用更加声明式的格式编程；因为 `merge()` 是一个高阶函数，所以新代码实际上是一个不错的函数式编程示例：

##### A Map in the functional style

```
public static void incrementPageVisit(Map<String, Integer> pageVisits, String page) {
    pageVisits.merge(page, 1, (oldValue, value) -> oldValue + value);
}

```

Show moreShow more icon

在清单 4 中， `page` 作为第一个参数传递给 `merge()` ：该键的值应该更新。第二个参数是给该键分配的初始值， \_如果`Map` 中不存在该键（在本例中，该值为 “1”）。第三个参数（一个拉姆达表达式）接收 map 中该键对应的值作为其参数，并且将该值作为变量传递给 merge 方法中的第二个参数。 这个拉姆达表达式返回的是它的参数的和，实际上就是递增计数。

比较 [A Map in the functional style](#a-map-in-the-functional-style) 中的 `incrementPageVisit()` 方法中的一行代码与 [A Map in the imperative style](#a-map-in-the-imperative-style) 中的多行代码。 [A Map in the functional style](#a-map-in-the-functional-style) 中的程序是一个函数式编程示例，由此可见声明式考虑问题有助于我们的跳跃性思维。

## 结束语

在 Java 程序中采用函数式方法和语法有许多好处：代码简洁，更富于表达，不易出错，更容易并行化，而且通常比面向对象的代码更容易理解。我们面临的挑战在于将思维方式从命令式编程 — 绝大多数开发人员都熟悉它 — 转变为声明式思考。

尽管函数式编程不那么容易或直观，但您可以通过学习关注您想要程序实现的 _目的_ 而不是关注您希望它执行的 _方式_，从而实现思维上的巨大飞跃。通过允许底层函数库管理执行代码，您将逐步直观地了解高阶函数，它们是函数式编程的构建基块。

本文翻译自： [An easier path to functional programming in Java](https://developer.ibm.com/articles/j-java8idioms1/)（2017-02-15）