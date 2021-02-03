# Java Interface 是常量存放的最佳地点吗？
Java Interface 作为常量存放的最佳地点会产生的问题

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/l-java-interface/)

bright

发布: 2003-04-21

* * *

## 前言

由于 java interface 中声明的字段在编译时会自动加上 static final 的修饰符，即声明为常量。因而 interface 通常是存放常量的最佳地点。然而在 java 的实际应用时却会产生一些问题。

问题的起因有两个，第一，是我们所使用的常量并不是一成不变的，而是相对于变量不能赋值改变。例如我们在一个工程初期定义常量∏＝3.14，而由于计算精度的提高我们可能会重新定义∏＝3.14159，此时整个项目对此常量的引用都应该做出改变。第二，java 是动态语言。与 c++之类的静态语言不同,java 对一些字段的引用可以在运行期动态进行，这种灵活性是 java 这样的动态语言的一大优势。也就使得我们在 java 工程中有时部分内容的改变不用重新编译整个项目，而只需编译改变的部分重新发布就可以改变整个应用。

讲了这么多，你还不知道我要说什么吗？好，我们来看一个简单的例子：

有一个 interface A，一个 class B，代码如下：

```
//file A.java
public interface A{
    String name = "bright";
}
//file B.java
public class B{
    public static void main(String[] args){
        System.out.println("Class A's name = " + A.name);
    }
}

```

Show moreShow more icon

够简单吧，好，编译 A.java 和 B.java。

运行，输入 java B，显然结果如下：

```
Class A's name = bright

```

Show moreShow more icon

我们现在修改 A.java 如下：

```
//file A.java
public interface A{
    String name = "bright sea";
}

```

Show moreShow more icon

编译 A.java 后重新运行 B class，输入 java B，注意：结果如下

```
Class A's name = bright

```

Show moreShow more icon

为什么不是”Class A’s name = bright sea”？让我们使用 jdk 提供的反编译工具 javap 反编译 B.class 看个究竟，输入：javap -c B ，结果如下：

```
Compiled from B.java
public class B extends java.lang.Object {
    public B();
    public static void main(java.lang.String[]);
}
Method B()
0 aload_0
1 invokespecial #1 <Method java.lang.Object()>
4 return
Method void main(java.lang.String[])
0 getstatic #2 <Field java.io.PrintStream out>
3 ldc #3 <String "Class A's name = bright">
5 invokevirtual #4 <Method void println(java.lang.String)>
8 return

```

Show moreShow more icon

注意到标号 3 的代码了吗？由于引用了一个 static final 的字段，编译器已经将 interface A 中 name 的内容编译进了 class B 中，而不是对 interface A 中的 name 的引用。因此除非我们重新编译 class B，interface A 中 name 发生的变化无法在 class B 中反映。如果这样去做那么 java 的动态优势就消失殆尽。

解决方案，有两种解决方法。

第一种方法是不再使用常量，将所需字段放入 class 中声明，并去掉 final 修饰符。但这种方法存在一定的风险，由于不再是常量着因而在系统运行时有可能被其他类修改其值而发生错误，也就违背了我们设置它为常量的初衷，因而不推荐使用。

第二种方法，将常量放入 class 中声明，使用 class 方法来得到此常量的值。为了保持对此常量引用的简单性，我们可以使用一个静态方法。我们将 A.java 和 B.java 修改如下：

```
//file A.java
public class A{
    private static final String name = "bright";
    public static String getName(){
        return name;
    }
}
//file B.java
public class B{
    public static void main(String[] args){
        System.out.println("Class A's name = " + A.getName());
    }
}

```

Show moreShow more icon

同样我们编译 A.java 和 B.java。运行 class B，输入 java B，显然结果如下：

Class A’s name = bright

现在我们修改 A.java 如下：

```
//file A.java
public class A{
    private static final String name = "bright";
    public static String getName(){
        return name;
    }
}

```

Show moreShow more icon

我们再次编译 A.java 后重新运行 B class，输入 java B：结果如下

```
Class A's name = bright sea

```

Show moreShow more icon

终于得到了我们想要的结果，我们可以再次反编译 B.class 看看 class B 的改变，输入：

javap -c B,结果如下：

```
Compiled from B.java
public class B extends java.lang.Object {
    public B();
    public static void main(java.lang.String[]);
}
Method B()
0 aload_0
1 invokespecial #1 <Method java.lang.Object()>
4 return
Method void main(java.lang.String[])
0 getstatic #2 <Field java.io.PrintStream out>
3 new #3 <Class java.lang.StringBuffer>
6 dup
7 invokespecial #4 <Method java.lang.StringBuffer()>
10 ldc #5 <String "Class A's name = ">
12 invokevirtual #6 <Method java.lang.StringBuffer append(java.lang.String)>
15 invokestatic #7 <Method java.lang.String getName()>
18 invokevirtual #6 <Method java.lang.StringBuffer append(java.lang.String)>
21 invokevirtual #8 <Method java.lang.String toString()>
24 invokevirtual #9 <Method void println(java.lang.String)>
27 return

```

Show moreShow more icon

注意标号 10 至 15 行的代码，class B 中已经变为对 A class 的 getName()方法的引用，当常量 name 的值改变时我们只需对 class A 中的常量做修改并重新编译，无需编译整个项目工程我们就能改变整个应用对此常量的引用，即保持了 java 动态优势又保持了我们使用常量的初衷，因而方法二是一个最佳解决方案。