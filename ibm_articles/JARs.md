# JARs
Java Archive 不仅仅是一堆类

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-5things6/)

Ted Neward, Alex Theedom

发布: 2010-07-22

* * *

对于大多数 Java 开发人员来说，JAR 文件及其 “近亲” WAR 和 EAR 都只不过是漫长的 Ant 或 Maven 流程的最终结果。标准步骤是将一个 JAR 复制到服务器（或者，少数情况下是用户机）中的合适位置，然后忘记它。

事实上，JAR 能做的不止是存储源代码，您应该了解 JAR 还能做什么，以及如何进行。在这一期的 _[5 件事](/zh/series/5-things-you-didnt-know-about/)_ 系列中，将向您展示如何最大限度地利用 Java Archive 文件（有时候也可是 WAR 和 EAR），特别是在部署时。

由于有很多 Java 开发人员使用 Spring（因为 Spring 框架给传统的 JAR 使用带来一些特有的挑战），这里有几个具体技巧用于在 Spring 应用程序中处理 JAR 。

##### 关于本系列

您觉得自己懂 Java 编程？事实是，大多数开发人员都只领会到了 Java 平台的皮毛，所学也只够应付工作。在 [本系列](/zh/series/5-things-you-didnt-know-about/) 中，Ted Neward 深度挖掘 Java 平台的核心功能，揭示一些鲜为人知的事实，帮助您解决最棘手的编程困难。

我将以一个标准 Java Archive 文件产生过程的简单示例开始，这将作为以下技巧的基础。

## 把它放在 JAR 中

通常，在源代码被编译之后，您需要构建一个 JAR 文件，使用 `jar` 命令行实用工具，或者，更常用的是 Ant `jar` 任务将 Java 代码（已经被包分离）收集到一个单独的集合中，过程简洁易懂，我不想在这做过多的说明，稍后将继续说明如何构建 JAR。现在，我只需要存档 `Hello`，这是一个独立控制台实用工具，对于执行打印消息到控制台这个任务十分有用。如清单 1 所示：

##### 清单 1\. 存档控制台实用工具

```
package com.tedneward.jars;

public class Hello
{
    public static void main(String[] args)
    {
        System.out.println("Howdy!");
    }
}

```

Show moreShow more icon

`Hello` 实用工具内容并不多，但是对于研究 JAR 文件却是一个很有用的 “脚手架”，我们先从执行此代码开始。

## 1\. JAR 是可执行的

.NET 和 C++ 这类语言一直是 OS 友好的，只需要在命令行（`helloWorld.exe`）引用其名称，或在 GUI shell 中双击它的图标就可以启动应用程序。然而在 Java 编程中，启动器程序 — `java` — 将 JVM 引导入进程中，我们需要传递一个命令行参数（`com.tedneward.Hello`）指定想要启动的 `main()` 方法的类。

这些附加步骤使使用 Java 创建界面友好的应用程序更加困难。不仅终端用户需要在命令行输入所有参数（终端用户宁愿避开），而且极有可能使他或她操作失误以及返回一个难以理解的错误。

这个解决方案使 JAR 文件 “可执行”，以致 Java 启动程序在执行 JAR 文件时，自动识别哪个类将要启动。我们所要做的是，将一个入口引入 JAR 文件清单文件（ `MANIFEST.MF` 在 JAR 的 `META-INF` 子目录下），像这样：

##### 清单 2\. 展示入口点！

```
Main-Class: com.tedneward.jars.Hello

```

Show moreShow more icon

这个清单文件只是一个名值对。因为有时候清单文件很难处理回车和空格，然而在构建 JAR 时，使用 Ant 来生成清单文件是很容易的。在清单 3 中，使用 Ant `jar` 任务的 `manifest` 元素来指定清单文件：

##### 清单 3\. 构建我的入口点！

```
<target name="jar" depends="build">
        <jar destfile="outapp.jar" basedir="classes">
            <manifest>
                <attribute name="Main-Class" value="com.tedneward.jars.Hello" />
            </manifest>
        </jar>
    </target>

```

Show moreShow more icon

现在用户在执行 JAR 文件时需要做的就是通过 `java -jar outapp.jar` 在命令行上指定其文件名。就 GUI shell 来说，双击 JAR 文件即可。

## 2\. JAR 可以包括依赖关系信息

似乎 `Hello` 实用工具已经展开，改变实现的需求已经出现。Spring 或 Guice 这类依赖项注入（DI）容器可以为我们处理许多细节，但是仍然有点小问题：修改代码使其含有 DI 容器的用法可能导致清单 4 所示的结果，如：

##### 清单 4\. Hello、Spring world！

```
package com.tedneward.jars;

import org.springframework.context.*;
import org.springframework.context.support.*;

public class Hello
{
    public static void main(String[] args)
    {
        ApplicationContext appContext =
            new FileSystemXmlApplicationContext("./app.xml");
        ISpeak speaker = (ISpeak) appContext.getBean("speaker");
        System.out.println(speaker.sayHello());
    }
}

```

Show moreShow more icon

##### 关于 Spring 的更多信息

这个技巧假定您熟悉依赖项注入和 Spring 框架。

由于启动程序的 `-jar` 选项将覆盖 `-classpath` 命令行选项中的所有内容，因此运行这些代码时，Spring 必须是在 `CLASSPATH`和\_ 环境变量中。幸运的是，JAR 允许在清单文件中出现其他的 JAR 依赖项声明，这使得无需声明就可以隐式创建 CLASSPATH，如清单 5 所示：

##### 清单 5\. Hello、Spring CLASSPATH!

```
<target name="jar" depends="build">
        <jar destfile="outapp.jar" basedir="classes">
            <manifest>
                <attribute name="Main-Class" value="com.tedneward.jars.Hello" />
                <attribute name="Class-Path"
                    value="./lib/org.springframework.context-3.0.1.RELEASE-A.jar
                      ./lib/org.springframework.core-3.0.1.RELEASE-A.jar
                      ./lib/org.springframework.asm-3.0.1.RELEASE-A.jar
                      ./lib/org.springframework.beans-3.0.1.RELEASE-A.jar
                      ./lib/org.springframework.expression-3.0.1.RELEASE-A.jar
                      ./lib/commons-logging-1.0.4.jar" />
            </manifest>
        </jar>
    </target>

```

Show moreShow more icon

注意 `Class-Path` 属性包含一个与应用程序所依赖的 JAR 文件相关的引用。您可以将它写成一个绝对引用或者完全没有前缀。这种情况下，我们假设 JAR 文件同应用程序 JAR 在同一个目录下。

不幸的是， `value` 属性和 Ant `Class-Path` 属性必须出现在同一行，因为 JAR 清单文件不能处理多个 `Class-Path` 属性。因此，所有这些依赖项在清单文件中必须出现在一行。当然，这很难看，但为了使 `java -jar outapp.jar` 可用，还是值得的！

## 3\. JAR 可以被隐式引用

如果有几个不同的命令行实用工具（或其他的应用程序）在使用 Spring 框架，可能更容易将 Spring JAR 文件放在公共位置，使所有实用工具能够引用。这样就避免了文件系统中到处都有 JAR 副本。Java 运行时 JAR 的公共位置，众所周知是 “扩展目录” ，默认位于 `lib/ext` 子目录，在 JRE 的安装位置之下。

JRE 是一个可定制的位置，但是在一个给定的 Java 环境中很少定制，以至于可以完全假设 `lib/ext` 是存储 JAR 的一个安全地方，以及它们将隐式地用于 Java 环境的 `CLASSPATH` 上。

## 4\. Java 6 允许类路径通配符

为了避免庞大的 `CLASSPATH` 环境变量（Java 开发人员几年前就应该抛弃的）和/或命令行 `-classpath` 参数，Java 6 引入了 _类路径通配符_ 的概念。与其不得不启动参数中明确列出的每个 JAR 文件，还不如自己指定 `lib/*` ，让所有 JAR 文件列在该目录下（不递归），在类路径中。

不幸的是，类路径通配符不适用于之前提到的 `Class-Path` 属性清单入口。但是这使得它更容易启动 Java 应用程序（包括服务器）开发人员任务，例如 code-gen 工具或分析工具。

## 5\. JAR 有的不只是代码

Spring，就像许多 Java 生态系统一样，依赖于一个描述构建环境的配置文件，前面提到过，Spring 依赖于一个 app.xml 文件，此文件同 JAR 文件位于同一目录 — 但是开发人员在复制 JAR 文件的同时忘记复制配置文件，这太常见了！

一些配置文件可用 sysadmin 进行编辑，但是其中很大一部分（例如 Hibernate 映射）都位于 sysadmin 域之外，这将导致部署漏洞。一个合理的解决方案是将配置文件和代码封装在一起 — 这是可行的，因为 JAR 从根本上来说就是一个 “乔装的” ZIP 文件。 当构建一个 JAR 时，只需要在 Ant 任务或 `jar` 命令行包括一个配置文件即可。

JAR 也可以包含其他类型的文件，不仅仅是配置文件。例如，如果我的 `SpeakEnglish` 部件要访问一个属性文件，我可以进行如下设置，如清单 6 所示：

##### 清单 6\. 随机响应

```
package com.tedneward.jars;

import java.util.*;

public class SpeakEnglish
    implements ISpeak
{
    Properties responses = new Properties();
    Random random = new Random();

    public String sayHello()
    {
        // Pick a response at random
        int which = random.nextInt(5);

        return responses.getProperty("response." + which);
    }
}

```

Show moreShow more icon

可以将 `responses.properties` 放入 JAR 文件，这意味着部署 JAR 文件时至少可以少考虑一个文件。这只需要在 JAR 步骤中包含 responses.properties 文件即可。

当您在 JAR 中存储属性之后，您可能想知道如何将它取回。如果所需要的数据与 JAR 文件在同一位置，正如前面的例子中提到的那样，不需要费心找出 JAR 文件的位置，使用 `JarFile` 对象就可将其打开。相反，可以使用类的 `ClassLoader` 找到它，像在 JAR 文件中寻找 “资源” 那样，使用 `ClassLoader getResourceAsStream()` 方法，如清单 7 所示：

##### 清单 7\. ClassLoader 定位资源

```
package com.tedneward.jars;

import java.util.*;

public class SpeakEnglish
    implements ISpeak
{
    Properties responses = new Properties();
    // ...

    public SpeakEnglish()
    {
        try
        {
            ClassLoader myCL = SpeakEnglish.class.getClassLoader();
            responses.load(
                myCL.getResourceAsStream(
                    "com/tedneward/jars/responses.properties"));
        }
        catch (Exception x)
        {
            x.printStackTrace();
        }
    }

    // ...
}

```

Show moreShow more icon

您可以按照以上步骤寻找任何类型的资源：配置文件、审计文件、图形文件，等等。几乎任何文件类型都能被捆绑进 JAR 中，作为一个 `InputStream` 获取（通过 `ClassLoader` ），并通过您喜欢的方式使用。

## 结束语

本文涵盖了关于 JAR 大多数开发人员所不知道的 5 件最重要的事 — 至少基于历史，有据可查。注意，所有的 JAR 相关技巧对于 WAR 同样可用，一些技巧（特别是 `Class-Path` 和 `Main-Class` 属性）对于 WAR 来说不是那么出色，因为 servlet 环境需要全部目录，并且要有一个预先确定的入口点，但是，总体上来看这些技巧可以使我们摆脱 “好的，开始在该目录下复制……” 的模式，这也使得他们部署 Java 应用程序更为简单。

本系列的下一个主题是： [关于 Java 应用程序性能监视您不知道的 5 件事](/zh/articles/j-5things7/)。

## 下载示例代码

[j-5things6-src.zip](http://public.dhe.ibm.com/software/dw/java/j-5things6-src.zip): 本文样例代码

本文翻译自： [JARs](https://developer.ibm.com/articles/j-5things6/)（2010-06-15）