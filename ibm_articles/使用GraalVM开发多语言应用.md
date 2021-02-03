# 使用 GraalVM 开发多语言应用
多语言应用开发利器

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-use-graalvm-to-run-polyglot-apps/)

成 富

发布: 2019-11-12

* * *

随着应用开发复杂度的不断提高，越来越多的应用选择在开发时使用多种语言。这是由于不同语言及其平台都有其各自所擅长处理的领域。在后端开发中，除了 Java 和 JVM 平台语言（如 Kotlin、Scala 和 Groovy 等）之外，其他语言也有各自的用武之地。JavaScript 语法灵活，在 NodeJS 平台上可以使用海量的第三方库；Python 在数据处理、机器学习和人工智能等方面有强大的优势；Ruby 在 Web 应用开发中以高效著称；R 语言则提供了对统计分析、数据挖掘和绘图的良好支持。在应用开发中，应该根据实际的需求来选择最合适的语言，期望用单一语言来解决所有问题是不现实的，会产生事倍功半的结果。

每个语言都需要自己的运行时平台提供支持。这就意味着多语言应用在运行时需要安装多种不同的平台。除了 Java 应用所需的 Java 虚拟机之外，还需要 NodeJS、Python 和 Ruby 等。除了安装这些语言之外，更复杂的问题是如何在多个语言之间实现互操作。Java 平台从 1.6 开始引入了脚本语言支持（JSR 223），后来转变为 Java 9 中的 `java.scripting` 模块。这是 Java 平台在多语言应用开发上的一种尝试。不过，Java 平台中的脚本引擎仅提供了简单的互操作功能，只支持在 Java 语言中运行其他脚本语言，并且在运行时依赖 Nashorn、JRuby 和 Jython 等第三方库的支持。

## 认识 GraalVM

GraalVM 采用了不同的实现方式，用一个单一的虚拟机来对不同语言提供支持。GraalVM 底层是基于 OpenJDK 的 Java HotSpot 虚拟机。Graal 编译器是基于 JVM 编译器接口（JVM Compiler Interface）实现的即时（just-in-time，JIT）编译器，用来把 Java 字节代码转换成本地机器代码。Java 和其他 JVM 语言（如 Groovy、Kotlin 和 Scala 等）的源代码在编译成 Java 字节代码之后，直接运行在 GraalVM 的 Java 虚拟机之上。GraalVM 中包含用来创建其他语言实现的 Truffle 框架。GraalVM 对 JavaScript、Ruby、Python、R 和 C/C++语言的支持，都是基于 Truffle 框架来实现的。Truffle 框架是一个开放的语言实现框架。其他语言都可以通过 Truffle 框架运行在 GraalVM 上，甚至是应用本身私有的语言。

## GraalVM 使用

GraalVM 有社区版和企业版两种版本可选。社区版可以从 [GitHub](https://github.com/oracle/graal/releases) 上下载，企业版则需要从 Oracle 官方网站下载。本文使用的是 GraalVM 社区版 19.2.0.1。该 GraalVM 版本基于 JDK 8u222 版本。GraalVM 提供了 Linux、macOS X 和 Windows 版本，其中 Windows 版本支持是实验性的。除了在本地上安装 GraalVM 之外，还可以使用 Docker 来运行。GraalVM 对应的 Docker 镜像名称是 `oracle/graalvm-ce` 。

GraalVM 包含一些核心组件和可选组件。核心组件包括 Java HotSpot 虚拟机和 Graal 编译器、NodeJS 运行时、LLVM 解释器、GraalVM 多语言支持 API 和 GraalVM 更新器等。可选组件包括 GraalVM 原生镜像支持、Python、Ruby 和 R 语言的解释器，以及 LLVM 工具链等。核心组件包含在安装包之中，可选组件需要通过 GraalVM 更新器安装。

### GraalVM 更新器

GraalVM 安装之后的 `bin` 目录中的命令行工具 `gu` 用来管理可选组件。使用 `gu available` 命令可以列出来所有可用的组件。表 1 列出了 GraalVM 所有的可选组件。

##### 表 1\. GraalVM 可用的可选组件

**组件标识符****版本****组件名称****来源**llvm-toolchain19.2.0.1LLVM.org toolchaingithub.comnative-image19.2.0.1Native Imagegithub.compython19.2.0.1Graal.Pythongithub.comR19.2.0.1FastRgithub.comruby19.2.0.1TruffleRubygithub.com

再使用 `gu install` 命令来安装对应组件，如 `gu install ruby` 安装 Ruby 语言解释器。使用 `gu list` 命令可以列出当前已安装组件。当不需要某个组件时，可以使用 `gu remove` 命令来删除，如 `gu remove ruby` 来删除 Ruby 语言解释器。在组件安装之后，会在 `bin` 目录添加语言相关的命令行工具。比如，安装 Ruby 语言解释器之后会添加 `ruby` 、 `irb` 、 `rake` 、 `rdoc` 、 `gem` 和 `bundler` 等命令。

## 多语言开发

对于 Java 和 JVM 语言来说，可以直接使用 GraalVM 中的 Java 编译器进行编译，使用 Java 虚拟机来运行。除了 Java 和 JVM 语言之外的其他语言，都是基于 Truffle 框架实现的，可以使用多语言 API 来进行互操作。这种互操作不仅限于 JVM 语言和其他语言之间，这些语言之间也能进行互操作。比如，可以在 JavaScript 代码中调用 Ruby 语言。除了 Java 和 JVM 语言之外，其他基于 Truffle 框架实现的语言都称为访客语言（guest language）。在 Java 应用中，使用 GraalVM 的多语言开发功能需要用到 GraalVM 的 Java SDK。该 SDK 的 JAR 文件 `graal-sdk.jar` 已经包含在 GraalVM 的 `jre/lib/boot` 目录中，可以直接使用。在开发中，可以把该 JAR 文件作为 JDK 的一部分添加到 IDE 中，以保证代码可以在 IDE 中编译。也可以使用 Maven 添加对 GraalVM SDK 的依赖，如代码清单 1 所示。

##### 清单 1\. 添加 GraalVM SDK 的 Maven 依赖

```
<dependency>
<groupId>org.graalvm.sdk</groupId>
<artifactId>graal-sdk</artifactId>
<version>19.2.0.1</version>
<scope>provided</scope>
</dependency>

```

Show moreShow more icon

### 在 Java 中调用 JavaScript

代码清单 2 给出了运行 JavaScript 代码的 Java 示例。

##### 清单 2\. 在 Java 中调用 JavaScript

```
import org.graalvm.polyglot.Context;
import org.graalvm.polyglot.Value;

public class JsEvaluation {

public static void main(String[] args) {
    try (Context context = Context.create()) {
      Value date = context.eval("js", "new Date().toString()");
      System.out.println(date.asString());
    }
}
}

```

Show moreShow more icon

`Context` 对象表示的是访客语言执行环境的上下文。该对象是 Java 代码与访客语言代码交互的桥梁，也限制了访客语言代码执行时的权限和所能使用的资源。 `Context.create()` 方法使用默认的配置创建一个 `Context` 对象。 `Context` 对象的 `eval()` 方法执行指定语言的一段代码，并返回 `Value` 对象作为结果。 `eval()` 方法的第一个参数是语言的标识符。每个语言都有自己唯一的标识符， `js` 表示的是 JavaScript。 `Value` 对象表示的是多语言值。 `Value` 类提供的方法可以对这样的值以语言无关的方式来访问。 `Value` 对象所表示的具体值与所对应的代码相关。代码清单 2 中的 `date` 表示的是字符串，可以使用 `asString()` 方法把它转换成 Java 中的 `String` 对象。

### 绑定对象

在执行语言代码时，可以通过绑定对象在代码运行的顶层域中设置值。代码清单 3 使用 Context 提供的构建器来创建 `Context` 对象，其中 `allowAllAccess(true)` 的作用是允许访客语言获得对宿主环境的完全访问。在使用 `Context` 对象的 `getBindings()` 方法得到指定语言的绑定对象之后，再通过 `putMember()` 方法来设置值。这里设置的是一个 Java 中的 `HashMap` 对象。在 JavaScript 代码中，通过绑定对象设置的值作为全局变量出现，因此可以直接使用。

##### 清单 3\. 使用绑定对象

```
public class JsBindings {

public static void main(String[] args) {
    try (Context context = Context.newBuilder().allowAllAccess(true).build()) {
      Value bindings = context.getBindings("js");
      Map<String, String> data = new HashMap<>();
      data.put("name", "Test");
      bindings.putMember("data", data);
      Value content = context.eval("js", "data.get('name')");
      System.out.println(content.asString());
    }
}
}

```

Show moreShow more icon

### 限制访问权限

使用 `allowAllAccess(true)` 的 `Context` 对象的访问权限过高，只适用于信任的访客语言代码。通常的做法是只允许访客代码访问宿主环境中的特定值。在 Java 中，使用 `@HostAccess.Export` 注解来声明开放的方法和字段。在代码清单 4 中， `Pointer` 类的 `getX()` 、 `getY()` 和 `move()` 方法都添加了 `@Export` 注解。因此可以在 JavaScript 代码中使用 `pointer.move()` 方法。

##### 清单 4\. 使用@HostAccess.Export 注解开放的方法

```
public class JsHostAccess {

public static void main(String[] args) {
    try (Context context = Context.create()) {
      Pointer pointer = new Pointer(0, 0);
      context.getBindings("js").putMember("pointer", pointer);
      Pointer movedPointer = context.eval("js", "pointer.move(1, 1)").as(Pointer.class);
      System.out.println(movedPointer);
    }
}

public static class Pointer {
    private final int x;
    private final int y;

    Pointer(int x, int y) {
      this.x = x;
      this.y = y;
    }

    @Export
    public int getX() {
      return x;
    }

    @Export
    public int getY() {
      return y;
    }

    @Export
    public Pointer move(int deltaX, int deltaY) {
      return new Pointer(x + deltaX, y + deltaY);
    }

    @Override
    public String toString() {
      return String.format("[%d, %d]", x, y);
    }
}
}

```

Show moreShow more icon

### 在 JavaScript 中调用 Java

在访客语言中也可以调用 Java。代码清单 5 中的 JavaScript 代码使用了 Java 中的 `java.time.Instant` 类。 `Java.type()` 方法可以加载指定的 Java 类。Java 类的静态方法可以直接调用，也可以使用 `new` 来创建对象。

##### 清单 5\. 在 JavaScript 中调用 Java

```
var Instant = Java.type('java.time.Instant');
var ChronoUnit = Java.type('java.time.temporal.ChronoUnit');
print('Now: ' + Instant.now());
var nowPlus30Days = Instant.now().plus(30, ChronoUnit.DAYS);
print('Now +30 days:' + nowPlus30Days);

```

Show moreShow more icon

### 其他语言的多语言开发

除了 Java 之外，其他语言也有相应的多语言 API。在 JavaScript 语言代码中，可以通过全局对象 `Polyglot` 来使用该 API。代码清单 6 中的 JavaScript 代码文件 `polyglot.js` 使用 `Polyglot.eval()` 方法来执行 Ruby 代码。

##### 清单 6\. 在 JavaScript 中调用 Ruby

```
var array = Polyglot.eval('ruby', '[1, "two", 3.0]');
print(array[1]);

```

Show moreShow more icon

代码清单 6 中的 `polyglot.js` 文件需要使用 `node --jvm --polyglot polyglot.js` 命令来运行。其中参数 `--jvm` 表示表示在 Java 虚拟机上运行 JavaScript 代码，而 `--polyglot` 表示启用多语言 API。

### 完整的 NodeJS 示例

代码清单 7 给出了一个完整的基于 NodeJS 的多语言应用。该应用是一个计算幂值的 API。该 HTTP API 接受两个参数 `base` 和 `exponent` ，并利用 Java 中的 `java.math.BigInteger` 类来进行实际的求值。

##### 清单 7\. 完整的 NodeJS 多语言示例

```
const http = require('http')
const url = require('url')
const port = 3000
const BigInteger = Java.type('java.math.BigInteger')

const requestHandler = (request, response) => {
try {
    const {query} = url.parse(request.url, true)
    const {base, exponent} = query
    const result = (new BigInteger(base)).pow(parseInt(exponent, 10))
    response.end(result.toString())
} catch (e) {
    console.error(e)
    response.end(e.message)
}
}

const server = http.createServer(requestHandler)

server.listen(port, (err) => {
if (err) {
    return console.log('Error', err)
}

console.log(`server is listening on ${port}`)
})

```

Show moreShow more icon

该 JavaScript 文件的名称为 `pow.js` 。在运行时，只需要用 GraalVM 提供的 node 命令即可。具体的命令是 `node --jvm --polyglot pow.js` 。除了使用语言本身的命令行工具运行代码之外，还可以使用 `polyglot` 命令来运行，如 `polyglot --jvm pow.js` 。使用 `polyglot` 命令的好处是可以指定多个不同语言的文件，如 `polyglot --jvm app.js app.rb app.py` 。

安装完其他语言的可选组件之后，如果在使用多语言 API 时出现语言找不到的错误，则说明需要重新构建多语言的镜像。GraalVM 提供的 `rebuild-images` 命令可以重新构建该镜像。当使用 `rebuild-images polyglot` 命令重新构建镜像之后，新增的语言应该可以被找到。运行 `polyglot --help` 命令可以显示出当前可用的语言。

## 构建原生镜像

运行 Java 应用都需要 Java 虚拟机的支持。即便是最简单的输出 `Hello World` 的 Java 应用也需要一个完整的 Java 虚拟机才能运行。Java 虚拟机除了占据硬盘和内存空间之外，也使得 Java 应用的分发和部署变得复杂。GraalVM 提供的原生镜像（native image）功能可以把 Java 代码预先编译（ahead-of-time）成单独的可执行文件，称为原生镜像。该可执行文件包括了应用本身的代码、所依赖的第三方库和 JDK 本身。该执行文件并不是运行在 Java 虚拟机之上，而是名为 Substrate 的虚拟机。与运行在传统的 Java 虚拟机上相比，原生镜像运行时的启动速度更快，所耗费的内存资源更少。以基于 GraalVM 的 Quarkus 库为例，使用 Quarkus 和 GraalVM 的 REST 应用的启动时间仅为 13 毫秒，占用内存仅 13MB；而传统的基于 Java 虚拟机的 REST 应用的启动时间需要 4.3 秒，占用内存为 140MB。

构建原生镜像需要使用 GraalVM 提供的生成工具。该工具是 GraalVM 的可选组件，可以通过 `gu install native-image` 命令来安装。安装完成之后会在 `bin` 目录中出现 `native-image` 命令行工具。该生成器会处理应用中所有的 Java 类及其依赖，来确定哪些类和方法是应用运行时需要的。GraalVM 编译器会以预先编译的方式把这些需要的代码编译到原生可执行文件中。

代码清单 8 中是最简单的输出 “Hello World” 的 Java 应用的源代码 `HelloWorld.java` 。在使用 `javac` 命令编译该文件之后，得到对应的 `HelloWorld.class` 文件。接着使用 `native-image HelloWorld hello` 命令可以从该 class 文件中产生名为 `hello` 的可执行文件。最终产生的 `hello` 文件仅有 2.2MB，节省了非常多的磁盘空间。

<h5 id=”清单-8-输出-” hello-world”-的-java-应用>清单 8. 输出 “Hello World” 的 Java 应用

```
public class HelloWorld {

public static void main(String[] args) {
    System.out.println("Hello World");
}
}

```

Show moreShow more icon

如果应用中使用了第三方库，可以使用 `native-image` 命令的 `-cp` 参数来指定查找相应库的目录。还可以使用 `-jar` 参数来从 JAR 文件中创建镜像。

如果 Java 应用中使用了 GraalVM 的多语言 API，则需要额外的配置来生成原生镜像。代码清单 9 给出了对应于 [清单 2](#清单-2-在-java-中调用-javascript) 的 `JsEvaluation` 类的原生镜像的生成命令。其中 `--language:js` 表示启用对 JavaScript 的多语言支持， `--initialize-at-build-time` 表示在构建过程中进行初始化。

##### 清单 9\. 生成使用多语言 API 的原生镜像

```
native-image --language:js --initialize-at-build-time -cp . JsEvaluation

```

Show moreShow more icon

GraalVM 的原生镜像功能通过静态代码分析来计算出运行时需要用到的全部 Java 类，因此不支持在运行时动态加载 Java 类。如果你的应用依赖于在运行时动态加载 Java 类的能力，那么该应用无法以 GraalVM 原生镜像的方式来运行。除此之外，有些 Java 虚拟机上的功能在使用原生镜像运行时的支持也是存在差别的。这是因为 GraalVM 原生镜像使用的 Substrate 虚拟机的功能不同于 Java HotSpot 虚拟机。以 Java 反射 API 为例，在构建原生镜像的过程中，静态分析会找出对反射相关 API 的调用，并确保所需的 Java 类被包括进来。当在代码中使用 `Class.forName()` 方法来加载类时，如果给出的类名是常量，对应的 Java 类会被添加。不过静态分析有它的局限性。如果 `Class.forName()` 方法的参数是一个变量，就无法包含所实际对应的类，从而在运行时产生错误。对于这样的情况，可以通过配置文件来指定运行时需要的类和方法。

除了反射 API 的限制之外，Java 7 中引入的 `invokedynamic` 指令和方法句柄也无法在 GraalVM 原生镜像中使用。Java 对象的 `finalizer` 也是不支持的。由于没有动态类加载，Java 的安全管理器（security manager）实际上是不需要的。Java 提供的管理和调试功能也是不支持的，包括 JVMTI、JMX 和其他虚拟机接口。这些是 Java 应用运行在普通 Java 虚拟机和以 GraalVM 原生镜像运行时的区别。因为这些差异，并不是所有应用都可以直接编译成原生镜像，可能需要添加相关的配置。不同第三方库对原生镜像的支持也不尽相同。Spring 框架计划在 5.3 版本中提供对原生镜像的完整支持。在原生镜像产生之后，需要进行大量的测试以确保不出现问题。考虑到 GraalVM 原生镜像所能带来的巨大的性能提升，对已有的应用进行改动以支持 GraalVM 原生镜像的付出是值得的。

## GraalVM 相关工具

下面对 GraalVM 提供的相关工具进行介绍。

### 调试器

GraalVM 对基于 Truffle 框架的语言都提供了远程调试的功能。GraalVM 实现了 Chrome DevTools 协议，可以通过 Chrome 开发者工具连接到 GraalVM 上运行的代码进行调试。只需要在运行时添加 `--inspect` 参数，就可以启动 GraalVM 内置的调试服务。进行调试的 URL 地址会被输出到控制台，只需要打开 Chrome 浏览器访问该地址即可。可以在 Chrome 开发者工具中打开源文件并添加断点进行调试。

### 性能评估

GraalVM 对于基于 Truffle 框架的语言提供了性能评估功能，可以分析 CPU 和内存的使用情况，找出性能瓶颈。GraalVM 提供了三种不同的工具，见表 2。其中 `--memtracer` 目前是实验性功能，需要添加额外参数 `--experimental-options` 启用该功能。

##### 表 2\. GraalVM 提供的性能评估功能

**启用参数****功能**–cpusampler记录 CPU 执行时间–cputracer记录语句执行次数–memtracer记录内存分配

当代码运行结束之后，收集的性能统计数据会被输出到控制台。代码清单 10 给出了启用表 2 中所有性能评估功能的命令行示例。

##### 清单 10\. GraalVM 性能评估示例

```
node --jvm --polyglot --cpusampler --cputracer --experimental-options --memtracer  pow.js

```

Show moreShow more icon

### VisualVM

GraalVM 提供了 VisualVM 的升级版。除了已有的功能之外，GraalVM 的 VisualVM 还可以查看 JavaScript、Python、Ruby 和 R 语言的内存堆信息、对象信息和线程信息等。不过该功能只在 GraalVM 的企业版中可用。

## 结束语

多语言应用可以充分利用不同语言和平台的优势，更加简洁高效的完成任务。GraalVM 在一个单一虚拟机上提供了对 Java 和 JVM 语言、JavaScript、Python、Ruby、R 和 C/C++等语言的支持，还可以使用 Truffle 框架进行扩展。GraalVM 的原生镜像功能可以创建启动迅速和占用资源更少的可执行文件。总言之，GraalVM 是开发多语言应用的极佳选择。

## 下载示例代码

[获取代码](https://github.com/VividcodeIO/graalvm-starter)