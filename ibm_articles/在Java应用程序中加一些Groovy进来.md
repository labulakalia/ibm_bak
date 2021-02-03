# 在 Java 应用程序中加一些 Groovy 进来
嵌入简单的、易于编写的脚本，从而利用 Groovy 的简单性

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-pg05245/)

Andrew Glover

发布: 2005-06-13

* * *

如果您一直在阅读这个系列，那么您应该已经看到有各种各样使用 Groovy 的有趣方式，Groovy 的主要优势之一就是它的生产力。Groovy 代码通常要比 Java 代码更容易编写，而且编写起来也更快，这使得它有足够的资格成为开发工作包中的一个附件。在另一方面，正如我在这个系列中反复强调的那样，Groovy 并不是 —— 而且也不打算成为 —— Java 语言的替代。所以，这里存在的问题是，能否把 Groovy 集成到 Java 的编程实践中？或者说这样做有什么用？ _什么时候_ 这样做有用？

这个月，我将尝试回答这个问题。我将从一些熟悉的事物开始，即从如何将 Groovy 脚本编译成与 Java 兼容的类文件开始，然后进一步仔细研究 Groovy 的编译工具（`groovyc`）是如何让这个奇迹实现的。了解 Groovy 在幕后做了什么是在 Java 代码中使用 Groovy 的第一步。

注意，本月示例中演示的一些编程技术是 `Groovlets` 框架和 Groovy 的 `GroovyTestCase` 的核心，这些技术我在前面的文章中已经讨论过。

##### 关于本系列

把任何工具集成到自己的开发实践的关键就是知道什么时候使用它，而什么时候应当把它留在箱子里。脚本语言能够成为工具箱中极为强大的附件，但只在将它恰当应用到合适场景时才这样。为此， _实战 Groovy_ 的一系列文章专门探索了 Groovy 的实际应用，并告诉您什么时候应用它们，以及如何成功地应用它们。

## 天作之合？

在本系列中以前的文章中，当我介绍如何 [用 Groovy 测试普通 Java 程序](http://www.ibm.com/developerworks/java/library/j-pg11094/?S_TACT=105AGX52&S_CMP=cn-a-j) 的时候，您可能已经注意到一些奇怪的事：我 _编译了_ 那些 Groovy 脚本。实际上，我将 groovy 单元测试编译成普通的 Java .class 文件，然后把它们作为 Maven 构建的一部分来运行。

这种编译是通过调用 `groovyc` 命令进行的，该命令将 Groovy 脚本编译成普通的 Java 兼容的 .class 文件。例如，如果脚本声明了一个类，那么调用 `groovyc` 会生成至少三个 .class 。文件本身会遵守标准的 Java 规则：.class 文件名称要和声明的类名匹配。

作为示例，请参见清单 1，它创建了一个简单的脚本，脚本声明了几个类。然后，您自己就可以看出 `groovyc` 命令生成的结果：

##### 清单 1\. Groovy 中的类声明和编译

```
package com.vanward.groovy
class Person {
fname
lname
age
address
contactNumbers
String toString(){

    numstr = new StringBuffer()
    if (contactNumbers != null){
     contactNumbers.each{
       numstr.append(it)
       numstr.append(" ")
     }
    }
    "first name: " + fname + " last name: " + lname +
    " age: " + age + " address: " + address +
    " contact numbers: " + numstr.toString()
}
}
class Address {
street1
street2
city
state
zip
String toString(){
"street1: " + street1 + " street2: " + street2 +
     " city: " + city + " state: " + state + " zip: " + zip
}
}
class ContactNumber {
type
number
String toString(){
"Type: " + type + " number: " + number
}
}
nums = [new ContactNumber(type:"cell", number:"555.555.9999"),
    new ContactNumber(type:"office", number:"555.555.5598")]
addr = new Address(street1:"89 Main St.", street2:"Apt #2",
    city:"Utopia", state:"VA", zip:"34254")
pers = new Person(fname:"Mollie", lname:"Smith", age:34,
    address:addr, contactNumbers:nums)
println pers.toString()

```

Show moreShow more icon

在清单 1 中，我声明了三个类 —— `Person` 、 `Address` 和 `ContactNumber` 。之后的代码根据这些新定义的类型创建对象，然后调用 `toString()` 方法。迄今为止，Groovy 中的代码还非常简单，但现在来看一下清单 2 中 `groovyc` 产生什么样的结果：

##### 清单 2\. groovyc 命令生成的类

```
aglover@12d21 /cygdrive/c/dev/project/target/classes/com/vanward/groovy
$ ls -ls
total 15
4 -rwxrwxrwx+ 1 aglover  user   3317 May  3 21:12 Address.class
3 -rwxrwxrwx+ 1 aglover  user   3061 May  3 21:12 BusinessObjects.class
3 -rwxrwxrwx+ 1 aglover  user   2815 May  3 21:12 ContactNumber.class
1 -rwxrwxrwx+ 1 aglover  user   1003 May  3 21:12
Person$_toString_closure1.class
4 -rwxrwxrwx+ 1 aglover  user   4055 May  3 21:12 Person.class

```

Show moreShow more icon

哇！ _五个_.class 文件！我们了解 `Person` 、 `Address` 和 `ContactNumber` 文件的意义，但是其他两个文件有什么作用呢？

研究发现， `Person$_toString_closure1.class` 是 `Person` 类的 `toString()` 方法中发现的闭包的结果。它是 `Person` 的一个内部类，但是 `BusinessObjects.class` 文件是怎么回事 —— 它可能是什么呢？

对 [清单 1\. Groovy 中的类声明和编译](#清单-1-groovy-中的类声明和编译) 的深入观察指出：我在脚本主体中编写的代码（声明完三个类之后的代码）变成一个 .class 文件，它的名称采用的是脚本名称。在这个例子中，脚本被命名为 `BusinessObjects.groovy` ，所以，类定义中没有包含的代码被编译到一个名为 `BusinessObjects` 的 .class 文件。

### 反编译

反编译这些类可能会非常有趣。由于 Groovy 处于代码顶层，所以生成的 .java 文件可能相当巨大；不过，您应当注意的是 Groovy 脚本中声明的类（如 `Person` ） 与类之外的代码（比如 `BusinessObjects.class` 中找到的代码）之间的区别。在 Groovy 文件中定义的类完成了 `GroovyObject` 的实现，而在类之外定义的代码则被绑定到一个扩展自 `Script` 的类。

例如，如果研究由 BusinessObjects.class 生成的 .java 文件，可以发现：它定义了一个 `main()` 方法和一个 `run()` 方法。不用惊讶， `run()` 方法包含我编写的、用来创建这些对象的新实例的代码，而 `main()` 方法则调用 `run()` 方法。

这个细节的全部要点再一次回到了：对 Groovy 的理解越好，就越容易把它集成到 Java 程序中。有人也许会问：”为什么我要这么做呢？”好了，我们想说您用 Groovy 开发了一些很酷的东西；那么如果能把这些东西集成到 Java 程序中，那不是很好吗？

只是为了讨论的原因，我首先 _试图_ 用 Groovy 创建一些有用的东西，然后我再介绍把它嵌入到普通 Java 程序中的各种方法。

## 再制作一个音乐 Groovy

我热爱音乐。实际上，我的 CD 收藏超过了我计算机图书的收藏。多年以来，我把我的音乐截取到不同的计算机上，在这个过程中，我的 MP3 收藏乱到了这样一种层度：只是表示品种丰富的音乐目录就有一大堆。

最近，为了让我的音乐收藏回归有序，我采取了第一步行动。我编写了一个快速的 Groovy 脚本，在某个目录的 MP3 收藏上进行迭代，然后把每个文件的详细信息（例如艺术家、专辑名称等）提供给我。脚本如清单 3 所示：

##### 清单 3\. 一个非常有用的 Groovy 脚本

```
package com.vanward.groovy
import org.farng.mp3.MP3File
import groovy.util.AntBuilder
class Song {

mp3file
Song(String mp3name){
    mp3file = new MP3File(mp3name)
}
getTitle(){
    mp3file.getID3v1Tag().getTitle()
}
getAlbum(){
    mp3file.getID3v1Tag().getAlbum()
}
getArtist(){
    mp3file.getID3v1Tag().getArtist()
}
String toString(){
    "Artist: " + getArtist() + " Album: " +
      getAlbum() + " Song: " + getTitle()
}
static getSongsForDirectory(sdir){
    println "sdir is: " + sdir
    ant = new AntBuilder()
    scanner = ant.fileScanner {
       fileset(dir:sdir) {
         include(name:"**/*.mp3")
       }
    }
    songs = []
    for(f in scanner){
      songs << new Song(f.getAbsolutePath())
    }
    return songs
}
}
songs = Song.getSongsForDirectory(args[0])
songs.each{
println it
}

```

Show moreShow more icon

正如您所看到的，脚本非常简单，对于像我这样的人来说特别有用。而我要做的全部工作只是把某个具体的目录名传递给它，然后我就会得到该目录中每个 MP3 文件的相关信息（艺术家名称、歌曲名称和专辑） 。

现在让我们来看看，如果要把这个干净的脚本集成到一个能够通过数据库组织音乐甚至播放 MP3 的普通 Java 程序中，我需要做些什么。

## Class 文件是类文件

正如前面讨论过的，我的第一个选项可能只是用 `groovyc` 编译脚本。在这个例子中，我期望 `groovyc` 创建 _至少_ 两个 .class 文件 —— 一个用于 `Song` 类，另一个用于 `Song` 声明之后的脚本代码。

实际上， `groovyc` 可能创建 5 个 .class 文件。这是与 `Songs.groovy` 包含三个闭包有关，两个闭包在 `getSongsForDirectory()` 方法中，另一个在脚本体中，我在脚本体中对 `Song` 的集合进行迭代，并调用 `println` 。

因为 .class 文件中有三个实际上是 Song.class 和 Songs.class 的内部类，所以我只需要把注意力放在两个 .class 文件上。Song.class 直接映射到 Groovy 脚本中的 `Song` 声明，并实现了 `GroovyObject` ，而 Songs.class 则代表我在定义 `Song` 之后编写的代码，所以也扩展了 `Script` 。

此时此刻，关于如何把新编译的 Groovy 代码集成到 Java 代码，我有两个选择：可以通过 Songs.class 文件中的 `main()` 方法运行代码 （因为它扩展了 `Script` ），或者可以将 Song.class 包含到类路径中，就像在 Java 代码中使用其他对象一样使用它。

## 变得更容易些

通过 `java` 命令调用 Songs.class 文件非常简单，只要您记得把 Groovy 相关的依赖关系和 Groovy 脚本需要的依赖关系包含进来就可以。把 Groovy 需要的类全都包含进来的最简单方法就是把包含全部内容的 Groovy 可嵌入 jar 文件添加到类路径中。在我的例子中，这个文件是 groovy-all-1.0-beta-10.jar。要运行 Songs.class，需要记得包含将要用到的 MP3 库（jid3lib-0.5.jar>），而且因为我使用 `AntBuilder` ，所以我还需要在类路径中包含 `Ant` 。清单 4 把这些放在了一起：

##### 清单 4\. 通过 Java 命令行调用 Groovy

```
c:\dev\projects>java -cp  ./target/classes/;c:/dev/tools/groovy/
groovy-all-1.0-beta-10.jar;C:/dev/tools/groovy/ant-1.6.2.jar;
C:/dev/projects-2.0/jid3lib-0.5.jar
com.vanward.groovy.Songs c:\dev09\music\mp3s
Artist: U2 Album: Zooropa Song: Babyface
Artist: James Taylor Album: Greatest Hits Song: Carolina in My Mind
Artist: James Taylor Album: Greatest Hits Song: Fire and Rain
Artist: U2 Album: Zooropa Song: Lemon
Artist: James Taylor Album: Greatest Hits Song: Country Road
Artist: James Taylor Album: Greatest Hits Song: Don't Let Me
Be Lonely Tonight
Artist: U2 Album: Zooropa Song: Some Days Are Better Than Others
Artist: Paul Simon Album: Graceland Song: Under African Skies
Artist: Paul Simon Album: Graceland Song: Homeless
Artist: U2 Album: Zooropa Song: Dirty Day
Artist: Paul Simon Album: Graceland Song: That Was Your Mother

```

Show moreShow more icon

## 把 Groovy 嵌入 Java 代码

虽然命令行的解决方案简单有趣，但它并不是所有问题的最终解决方案。如果对更高层次的完善感兴趣，那么可能将 MP3 歌曲工具直接导入 Java 程序。在这个例子中，我想导入 Song.class ，并像在 Java 语言中使用其他类那样使用它。类路径的问题与上面相同 ：我需要确保包含了 _uber-Groovy_ jar 文件、 `Ant` 和 jid3lib-0.5.jar 文件。在清单 5 中，可以看到如何将 Groovy MP3 工具导入简单的 Java 类中：

##### 清单 5\. 嵌入的 Groovy 代码

```
package com.vanward.gembed;
import com.vanward.groovy.Song;
import java.util.Collection;
import java.util.Iterator;
public class SongEmbedGroovy{
public static void main(String args[]) {
Collection coll = (Collection)Song.getSongsForDirectory
    ("C:\\music\\temp\\mp3s");
for(Iterator it = coll.iterator(); it.hasNext();){
    System.out.println(it.next());
}
}
}

```

Show moreShow more icon

## Groovy 类加载器

就在您以为自己已经掌握全部的时候，我要告诉您的是，还有更多在 Java 语言中使用 Groovy 的方法。除了通过直接编译把 Groovy 脚本集成到 Java 程序中的这个选择之外，当我想直接嵌入脚本时，还有其他一些选择。

例如，我可以用 Groovy 的 `GroovyClassLoader` ， _动态地_ 加载一个脚本并执行它的行为，如清单 6 所示：

##### 清单 6\. GroovyClassLoader 动态地加载并执行 Groovy 脚本

```
package com.vanward.gembed;
import groovy.lang.GroovyClassLoader;
import groovy.lang.GroovyObject;
import groovy.lang.MetaMethod;
import java.io.File;
public class CLEmbedGroovy{
public static void main(String args[]) throws Throwable{

ClassLoader parent = CLEmbedGroovy.class.getClassLoader();
GroovyClassLoader loader = new GroovyClassLoader(parent);

Class groovyClass = loader.parseClass(
    new File("C:\\dev\\groovy-embed\\src\\groovy\\
      com\\vanward\\groovy\\Songs.groovy"));

GroovyObject groovyObject = (GroovyObject)
    groovyClass.newInstance();

Object[] path = {"C:\\music\\temp\\mp3s"};
groovyObject.setProperty("args", path);
Object[] argz = {};

groovyObject.invokeMethod("run", argz);

}
}

```

Show moreShow more icon

##### Meta，宝贝

如果您属于那群疯狂的人中的一员，热爱反射，喜欢利用它们能做的精彩事情，那么您将热衷于 Groovy 的 `Meta` 类。就像反射一样，使用这些类，您可以发现 `GroovyObject` 的各个方面（例如它的方法），这样就可以实际地 _创建_ 新的行为并执行它。而且，这是 Groovy 的核心 —— 想想运行脚本时它将如何发威吧！

注意，默认情况下，类加载器将加载与脚本名称对应的类 —— 在这个例子中是 Songs.class，而 _不是_ Song.class>。因为我（和您）知道 Songs.class 扩展了 Groovy 的 Script 类，所以不用想也知道接下来要做的就是执行 `run()` 方法。

您记起，我的 Groovy 脚本也依赖于运行时参数。所以，我需要恰当地配置 `args` 变量，在这个例子中，我把第一个元素设置为目录名。

## 更加动态的选择

对于使用编译好的类，而且，通过类加载器来动态加载 `GroovyObject` 的替代，是使用 Groovy 优美的 `GroovyScriptEngine` 和 `GroovyShell` 动态地执行 Groovy 脚本。

把 `GroovyShell` 对象嵌入普通 Java 类，可以像类加载器所做的那样动态执行 Groovy 脚本。除此之外，它还提供了大量关于控制脚本运行的选项。在清单 7 中，可以看到 `GroovyShell` 嵌入到普通 Java 类的方式：

##### 清单 7\. 嵌入 GroovyShell

```
package com.vanward.gembed;
import java.io.File;
import groovy.lang.GroovyShell;
public class ShellRunEmbedGroovy{
public static void main(String args[]) throws Throwable{

String[] path = {"C:\\music\\temp\\mp3s"};
GroovyShell shell = new GroovyShell();
shell.run(new File("C:\\dev\\groovy-embed\\src\\groovy\\
    com\\vanward\\groovy\\Songs.groovy"),
    path);
}
}

```

Show moreShow more icon

可以看到，运行 Groovy 脚本非常容易。我只是创建了 `GroovyShell` 的实例，传递进脚本名称，然后调用 `run()` 方法。

还可以做其他事情。如果您喜欢，那么也可以得到自己脚本的 `Script` 类型的 `GroovyShell` 实例。使用 `Script` 类型，您就可以传递进一个 `Binding` 对象，其中包含任何运行时值，然后再继续调用 `run()` 方法，如清单 8 所示：

##### 清单 8\. 有趣的 GroovyShell

```
package com.vanward.gembed;
import java.io.File;
import groovy.lang.Binding;
import groovy.lang.GroovyShell;
import groovy.lang.Script;
public class ShellParseEmbedGroovy{
public static void main(String args[]) throws Throwable{
GroovyShell shell = new GroovyShell();
Script scrpt = shell.parse(
    new File("C:\\dev\\groovy-embed\\src\\groovy\\
      com\\vanward\\groovy\\Songs.groovy"));

Binding binding = new Binding();
Object[] path = {"C:\\music\\temp\\mp3s"};
binding.setVariable("args",path);
scrpt.setBinding(binding);

scrpt.run();
}
}

```

Show moreShow more icon

## Groovy 的脚本引擎

`GroovyScriptEngine` 对象动态运行脚本的时候，非常像 `GroovyShell` 。区别在于：对于 `GroovyScriptEngine` ，您可以在实例化的时候给它提供一系列目录，然后让它根据要求去执行多个脚本，如清单 9 所示：

##### 清单 9\. GroovyScriptEngine 的作用

```
package com.vanward.gembed;
import java.io.File;
import groovy.lang.Binding;
import groovy.util.GroovyScriptEngine;
public class ScriptEngineEmbedGroovy{
public static void main(String args[]) throws Throwable{

String[] paths = {"C:\\dev\\groovy-embed\\src\\groovy\\
    com\\vanward\\groovy"};
GroovyScriptEngine gse = new GroovyScriptEngine(paths);
Binding binding = new Binding();
Object[] path = {"C:\\music\\temp\\mp3s"};
binding.setVariable("args",path);

gse.run("Songs.groovy", binding);
gse.run("BusinessObjects.groovy", binding);
}
}

```

Show moreShow more icon

在清单 9 中，我向实例化的 `GroovyScriptEngine` 传入了一个数组，数据中包含我要处理的路径，然后创建大家熟悉的 `Binding` 对象，然后再执行仍然很熟悉的 `Songs.groovy` 脚本。只是为了好玩，我还执行了 `BusinessObjects.groovy` 脚本，您或许还能回忆起来，它在开始这次讨论的时候出现过。

## Bean 脚本框架

最后，当然并不是最不重要的，是来自 Jakarta 的古老的 Bean 脚本框架（ Bean Scripting Framework —— BSF）。BSF 试图提供一个公共的 API，用来在普通 Java 应用程序中嵌入各种脚本语言（包括 Groovy）。这个标准的、但是有争议的最小公因子方法，可以让您毫不费力地将 Java 应用程序嵌入 Groovy 脚本。

还记得前面的 `BusinessObjects` 脚本吗？在清单 10 中，可以看到 BSF 可以多么容易地让我把这个脚本插入普通 Java 程序中：

##### 清单 10\. BSF 开始工作了

```
package com.vanward.gembed;
import org.apache.bsf.BSFManager;
import org.codehaus.groovy.runtime.DefaultGroovyMethods;
import java.io.File;
import groovy.lang.Binding;
public class BSFEmbedGroovy{
public static void main(String args[]) throws Exception {
String fileName = "C:\\dev\\project\\src\\groovy\\
    com\\vanward\\groovy\\BusinessObjects.groovy";
//this is required for bsf-2.3.0
//the "groovy" and "gy" are extensions
BSFManager.registerScriptingEngine("groovy",
    "org.codehaus.groovy.bsf.GroovyEngine", new
      String[] { "groovy" });
BSFManager manager = new BSFManager();
//DefaultGroovyMethods.getText just returns a
//string representation of the contents of the file
manager.exec("groovy", fileName, 0, 0,
    DefaultGroovyMethods.getText(new File(fileName)));
}
}

```

Show moreShow more icon

## 结束语

如果在本文中有一件事是清楚的话，那么只件事就是 Groovy 为了 Java 代码内部的重用提供了一堆选择。从把 Groovy 脚本编译成普通 Java .class 文件，到动态地加载和运行脚本，这些问题需要考虑的一些关键方面是灵活性和耦合。把 Groovy 脚本编译成普通 .class 文件是 _使用_ 您打算嵌入的功能的最简单选择，但是动态加载脚本可以使添加或修改脚本的行为变得更容易，同时还不必在编译上牺牲时间。（当然，这个选项只是在接口不变的情况下才有用。）

把脚本语言嵌入普通 Java 不是每天都发生，但是机会确实不断出现。这里提供的示例把一个简单的目录搜索工具嵌入到基于 Java 的应用程序中，这样 Java 应用程序就可以很容易地变成 MP3 播放程序或者其他 MP3 播放工具。虽然我 _可以_ 用 Java 代码重新编写 MP3 文件搜索器，但是我不需要这么做：Groovy 极好地兼容 Java 语言，而且，我很有兴趣去摆弄所有选项！

本文翻译自： [Practically Groovy, Stir some Groovy into your Java apps](https://www.ibm.com/developerworks/java/library/j-pg05245/)（2005-05-24）