# 以简单的方式消除 Java 的冗长
Lombok 是一种 Java 实用工具

**标签:** 

[原文链接](https://developer.ibm.com/zh/articles/os-lombok/)

Brian Carey

发布: 2010-04-11

* * *

## 何为 Lombok？

Lombok 是一种 Java Archive (JAR) 文件，可用来消除 Java 代码的冗长。

我们看这样一个例子，一个标准的 Java bean。一个典型的 Java bean 一般具有几个属性。每个属性具有一个 accessor 和 mutator（getter 和 setter）。通常还会有一个 `toString()` 方法、一个 `equals()` 方法和一个 `hashCode()` 方法。

初看上去，其中可预见的冗余就已经非常多了。如果每个属性都具有一个 getter 和 setter，并且通常如此，那么又何必详细说明呢？

让我们来看看 Lombok。为了消除代码行，Lombok 使用注释来标识类和 Java 代码块。在前述的那个 Java bean 示例中，所有的 getter、setter 以及其他三个方法都是在编译时被暗示并包括进来的。

而且更好的是如果您使用的是 Eclipse 或 IBM® WebSphere® Studio Application Developer（如果还没用的话，建议最好使用），您就可以将 Lombok 集成到 Java 项目并即刻获得开发时结果。换言之，Eclipse 编译器可以立即识别所暗指的 getters/setters，而其他 Java 代码则可引用这些方法。

最直接的好处当然是代码行的减少，这真的很棒。并且，如果有一个特定的 getter 或 setter 需要特别的注意，那么您就不必为了找到这个特定的 getter 或 setter 而遍历数十行代码。代码也会更为简洁并且冗余也少了。

Lombok 还让您得以简化代码的其他部分 — 不仅仅是 Java bean。比如，还可以减少 try/catch/finally 块内以及同步方法内的冗余代码。

现在，我们来看看在您自己的开发环境中如何能实现上述目的。

## 安装 Lombok

要进行安装，本文假设您使用的是 Eclipse 或 WebSphere Studio Application Developer。如果不是，您仍可使用 Lombok；但是不能享用开发时的种种益处。不过，您仍然可以享用编译时的益处。

首先，打开您的浏览器并将 URL 指向 `http://projectlombok.org/` 。

在撰写本文之时，用这个 URL 打开的页面的右上角会出现一个很大的单词。这个单词是 “Download”。单击该单词并开始下载 lombok.jar。此文件无需解压缩，而从其他站点下载的文件中，99% 都需要解压缩。

下载此文件后，需要执行这个 JAR 文件。在您的操作系统中打开一个提示符，进入到安装了 lombok.jar 的那个目录，并键入 `java -jar lombok.jar` 。

以上假设在您的路径内已经有 Java Runtime Environment (JRE)。如果没有，需要添加它。如果要了解如何添加，可以参考针对您的具体操作系统的相关文档。

如果您使用的是 Microsoft® Windows® ，那么还可以双击这个 lombok.jar 图标。同样地，您必须能够从您的图形用户界面（GUI）执行 JAR。

不管采取何种方式，应该最终都能看到一个 Lombok 安装屏幕。该屏幕会提问 Eclipse 或 WebSphere Studio Application Developer 可执行文件位于何处。它的默认位置有可能是正确的。但有时可能需要更改这个默认位置。

单击 **Install/Update** ，Lombok 会被迅速并入 Eclipse 开发环境。如果已经运行了 Eclipse，那么就需要关闭它并重启。

## 使用 Lombok

现在，就可以在 Eclipse 或 WebSphere Studio Application Developer 内开始使用 Lombok 了。请参考清单 1 内的代码。

##### 清单 1\. Java bean 的一个良好开端

```
public class Lure {
    private String name;
    private int size;
    private String color;
    private String style;
}

```

Show moreShow more icon

以上是一个简单的 Java bean 的典型开始。从这里，可以为每个属性添加 getters 和 setters。然后再添加一个 `equals()` 方法、一个 `toString()` 方法和一个 `hashCode()` 方法。

有了 Lombok，您无需自己完成上述操作。相反，您只需添加一个注释： `@Data` 。

没错，就这么简单。清单 2 中包括了 `@Data` 。

##### 清单 2\. Java bean 的一个更好的开端

```
import lombok.Data
public @Data class Lure {
    private String name;
    private int size;
    private String color;
    private String style;
}

```

Show moreShow more icon

不过请记住，只有当 lombok.jar 位于您的构建路径且 lombok.Data 被导入到这个 Java 类时，上述代码才会奏效。

如果在 Eclipse 或 WebSphere Studio Application Developer 内查看这个类的概要（通常位于屏幕上这个类的右侧），就能看到这些方法会被自动添加到这个 `Lure` 类。

若不能立即看到这个概要，可以单击 Eclipse 内的 Window 菜单，然后选择 **Show View** 。从所出现的弹出菜单中，选择 **Outline** ，它应该出现在屏幕的右侧。强制显示类的概要的热键组合是 **Alt+Shift+Q** ，然后是 **O** 。

如果您编写了另一个类来实例化 `Lure` ，您将能立刻拥有对 `Lure` 所暗指的方法（比如 `getName()` 或 `setSize()` ）的访问。您还能拥有对 `equals()` 、 `hashCode()` 和 `toString()` 的访问。很棒，对吧？

如果您使用的不是 Eclipse 或 WebSphere Studio Application Developer，那么所暗指的这些方法添加只有在实际编译这些代码时才能被认可。所以虽然在没有 Eclipse 或 WebSphere Studio Application Developer 时仍可以使用 Lombok，但 Lombok 最初的设计目的就是与 Eclipse 或 WebSphere Studio Application Developer 相集成。

在生成 getter/setter 方法时，Lombok 遵从传统的标准。所有这些方法名都以 `get` 或 `set` 开头并且属性名都是大写的。当然，如果属性是一个 Boolean，情况例外。在这种情况下，getter 以 `is` 开始，而非 `get` 。这是 Java bean 的一种标准实践。

现在，假设有一个 Java bean 对您的一个 getter 具有特殊要求。在清单 2 的例子中， `getStyle()` 可能返回颜色和大小的组合。在这种情况下，可以按自己的意愿编写 `getStyle()` 方法的代码。Lombok 检查您的代码并且不会基于这个属性创建其自己的 `getStyle` 版本。

又假设，您有一个 getter 方法不想公开。为此，Lombok 让您可以输入一个附加参数。清单 3 给出了一个定制的修饰符（modifier）。

##### 清单 3\. 一个定制的修饰符

```
private String name;
    @Getter(AccessLevel.PROTECTED) private int size;
    private String color;
    private String style;

```

Show moreShow more icon

在本例中， `getSize()` 方法将不会被公开。它具有一个受保护的修饰符，所以它只对派生子类可用并且在 `Lure` 类本身的内部。

您可能并不总是想接受 Lombok 为您提供的其他默认值。比如， `toString()` 方法会列出类名以及所有的属性名和值，中间以逗号分割。这个列表出现在类名的旁边。

比如，假设在记录这个 `Lure` 类时，您并不关心颜色。为了更改 `toString()` 的默认设置，需要使用 `ToString` 注释。

##### 清单 4\. 修改 `toString()`

```
@ToString(exclude="color")
public @Data class Lure {
    private String name;
    private int size;
    private String color;
    private String style;
}

```

Show moreShow more icon

若输出一个实例化了的 `Lure` 类，它应该看上去类似于：

```
Lure(name=Wishy-Washy, size=1, style=trolling)

```

Show moreShow more icon

注意到颜色没有被包括？这是因为您之前用注释告诉过 Lombok 不包括颜色。

您还可以修改 `equals()` 和 `hashCode()` 方法该如何被处理。清单 5 很直白，不需要过多解释。

##### 清单 5\. 修改 `hashCode()`

```
@EqualsAndHashCode(exclude="style")
public @Data class Lure {
    private String name;
    private int size;
    private String color;
    private String style;
}

```

Show moreShow more icon

在本例中，当 `equals()` 和 `hashCode()` 方法生成时， `style` 属性并没有被包括。

## 其他特性

您是不是也一直非常痛恨编写 try/catch/finally 块呢？我是这样的。幸运的是，有了 Lombok，您无需这么做了。这也是 Lombok 消除 Java 冗余的另一种方式。为了消除 try/catch/finally 块的冗余，只需使用 `@Cleanup` 注释。参见清单 6。

##### 清单 6\. 使用 `@Cleanup` 注释

```
public static void main(String[] args) throws IOException {
    @Cleanup InputStream in = new FileInputStream(args[0]);
    @Cleanup OutputStream out = new FileOutputStream(args[1]);
    //write file code goes here
}

```

Show moreShow more icon

上述代码较我们通常在标准 Java 代码内看到的整洁了很多。请注意您还是需要抛出由被调用代码捕获的异常（在本例中，为 `IOException` ）。

清单 6 中的这个代码块不仅消除了 try/catch/finally 块，而且还关闭了开放流。如果您处理的对象使用一个方法而不是 `close()` 来释放资源，那么就需要用一个带附加说明的注释调用该方法。比如， `@Cleanup("relinquish")` 。

Lombok 还可以减少同步方法所需的代码的冗余。很自然，这是用 `@Synchronized` 方法实现的。

##### 清单 7\. 使用 `@Synchronized` 注释

```
@Synchronized
private int foo() {
    //some magic done here
    return 1;
}

```

Show moreShow more icon

在本例中，Lombok 会自动创建一个名为 `$lock` 的实例对象，并会针对该对象同步方法 `foo()` 。

如果用 `@Synchronized` 注释的这个方法是静态的，那么 Lombok 就会创建一个名为 `$LOCK` 的类对象，并会针对该对象同步这个方法。

您还可以指定一个对象用以通过一个附加参数进行显式的锁定。比如， `@Synchronized("myObject")` 会针对对象 `myObject` 同步这个方法。在这种情况下，必须显式地定义它。

## 结束语

使用 Lombok，可以实现所有应用程序开发人员都竭尽全力实现的一个目标：消除冗余。

您还可以让您的代码可读性更好。在 Java bean 内寻找 “特殊”（即不遵循典型的标准）的具有大量属性的 getter 和 setter 方法将更为简便。这是因为只有这些特殊的 getter/setter 方法是需要被实际编码的。

Lombok 有助于代码的整洁、效率的提高以及冗余的减少。为何不在您自己的环境内尝试一下呢？

Brian Carey