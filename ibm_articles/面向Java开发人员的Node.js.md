# 面向 Java 开发人员的 Node.js
用于 Web 应用程序的轻量级的、事件驱动的 I/O

**标签:** Java,Node.js

[原文链接](https://developer.ibm.com/zh/articles/j-nodejs/)

Andrew Glover

发布: 2011-12-28

* * *

在过去的几年中，JavaScript 已经成为 Web 应用程序开发的幕后英雄。这样的情景让那些习惯将 JavaScript 称为 “玩具式语言” 的软件开发人员大跌眼镜。尽管有其他更流行的语言（这里是说开发人员忠实拥护的那些），而 JavaScript 作为标准且浏览器中立的脚本语言的地位从未动摇过。在客户端 Web 开发方面，它可能是世界上运用最广泛的语言。

##### Java 开发人员的 JavaScript

JavaScript 是现代 Java 开发人员的必备工具，而且并不难学。 [Java 开发 2.0：面向 Java 开发人员的 JavaScript](http://www.ibm.com/developerworks/cn/java/j-javadev2-18/) Andrew Glover，他将会介绍用来构建一流的 Web 应用程序所需的语法，其中包括 JavaScript 变量、类型、函数和类。

JavaScript 在服务器端脚本也占有一席之地，而且影响还会越来越大。虽然过去 JavaScript 在服务器端做过一些尝试，但并没有什么大的收获，除了 Node.js 或 Node。

Node 专门用来协助开发人员构建可扩展的网络程序，它是服务器端的编程环境，但几乎重建了整个 JavaScript 以服务于新一代的开发人员。对很多 Java 开发人员而言，Node 最大的吸引力在于软件并发的新方法。尽管 Java 平台在并发方式上还会不断进步（预计在 Java 7 和 Java 8 中会有大幅提升），但事实是 Node 满足了一个需求，它以更轻量级方式进行处理，这是众多 Java 开发人员想要的。就像客户端脚本使用 JavaScript 一样，在 Node 环境下编辑服务器端脚本非常棒，因为它非常有效，而且它能够在很多 Java 开发人员遇到挑战的领域中发挥作用。

在本文中，我将为您介绍服务器端脚本的革命，那就是 Node。首先我们从整体结构上总揽全局，看看是什么让 Node 与众不同，然后会演示如何快速构建一个使用 MongoDB 持久化数据的可扩展性 Web 应用程序。您将会亲身体验 Node 是多么有趣，并且看到用它来组建一个有效的 Web 应用程序有多快。

## Node 的事件驱动并发

Node 是建立在 Google 的 V8 JavaScript 引擎基础上的一个可扩展、事件驱动 I/O 环境。Google V8 实际上在执行之前就将 JavaScript 编译到本机机器代码中，从而达到快速的运行时性能，实际上这与 JavaScript 关系不大。Node 本身能够让您快速构建速度极快、高度并发的网络应用程序。

_事件驱动 I/O_ 对 Java 开发人员来说可能有些陌生，但并非是全新的概念。与 Java 平台中使用的多线程编程模型不同，Node 处理并发的方法是单线程加上 _事件循环_ 。Node 结构可以确保不会不堵塞和异步 I/O 传输。那些通常会引起堵塞的调用如等待数据库查询结果等在 Node 中不会发生。Node 应用程序不是等待代价高昂的 I/O 活动来完成操作，而是会发出一个回调（roll back）。当资源返回后，会异步调用关联的回调。

##### 为什么选择 Node.js？

Java 平台处理并发的方式帮助确立了其在企业开发方面的领导地位，这很难发生改变。诸如 Netty（以及 Gretty；参阅资源部分的框架和诸如 NIO 和 `java.util.concurrent` 的核心库，使得 JVM 成为处理并发的首选。Node 的特别之处就在于它是专门为应对并发编程而设计的现代开发环境。Node 的事件驱动编程方式意味着，您无需添加额外的库来处理并发性，这对关注多核硬件的开发人员来说是个好消息。

并发性在 Node 程序中 _照样工作_ 。如果在 Java 平台上运行之前的场景，我可能会选择复杂且耗时的方法，从传统的线程到 Java NIO 中较新的库，甚至是改善和更新的 `java.util.concurrent` 包。尽管 Java 的并发很强大，但却难以理解，还要用它来编码！相比而言，Node 的回调机制是嵌在语言内部的；您不需要添加诸如 `synchronized` 这样额外的结构使它工作。Node 的并发模型极其简单，让广大的开发人员都能够接受。

Node 的 JavaScript 语法也为您减少敲键的次数，节省大量时间。只需少量代码，您就能构建快速的、可扩展的 Web 应用程序，并且它能处理大量的并发连接。当然您也可以在 Java 平台上实现，但将会需要更多代码行和一大堆库和结构。如果您担心对新的编程环境不熟悉，那么大可不必：如果已经了解一些 JavaScript，那么 Node 很容易学习，这一点我可以向您保证。

## 开始使用 Node

如前所述，Node 很容易上手，而且有很多好的在线教程可以为您提供帮助。在本文中（以及我的 [Java 开发 2.0：面向 Java 开发人员的 JavaScript](http://www.ibm.com/developerworks/cn/java/j-javadev2-18/) 中），我重点关注的是帮助 Java 开发人员理解 Node 的优势。我不会从标准的 “Hello, world” Web 服务器应用程序讲起，我会直接讲解一个有实际意义的应用程序：将它想象成在 Node 上构建的 Foursquare 网站。

**安装 Node** 需要您遵循特定平台的说明；如果您使用的是 UNIX 一类的平台，如 OSX，那么我推荐您使用 [Node Version Manager](https://github.com/creationix/nvm) 或 NVM，它能处理安装正确版本 Node 的细节信息。无论哪种情况，请现在就 [下载并安装 Node](http://nodejs.org/#download) 。

我们还使用了第三方库来构建此应用程序，因此您可能还想 [安装 NPM](http://npmjs.org/) ，它是 Node 的包管理器。NPM 允许您指定项目的版本依赖关系，这样就可以下载并包含在您的构建路径中。NPM 在很多方面很像 Java 平台的 Maven，或者是 Ruby 的 Bundler。

## Node Express

Node 在 Web 开发人员心目中占据一席之地，一是因为它对并发的处理能力，二是它是按 Web 开发的需求构建的。最流行的一个第三方 Node 工具是轻量级 Web 开发框架 Express，我们将会用它来开发应用程序（参阅 参考资料 以了解关于 Express 的更多内容）。

Express 有很多特性，包括复杂路由、动态模板视图（参见 Node 框架 _du jour:_ Jade）和内容协商。Express 是非常轻量级的，没有内嵌的 ORM 或类似东西来加重其负担。在这一方面，Express 无法与 Rails、Grails 或其他有完整堆栈的 Web 框架相比。

安装并使用 Express 的一个简单方法是通过 NPM `package.json` 文件将其声明为依赖关系，如清单 1 所示。该文件与 Maven 的 `pom.xml` 或 Bundler 的 `Gemfile` 类似，但它是 JSON 格式的。

##### 清单 1.NPM 的 package.json 文件

```
{
"name":"magnus-server",
"version":"0.0.1",
"dependencies":{
    "express":"2.4.6"
}
}

```

Show moreShow more icon

在 [清单 1.NPM 的 package.json 文件](#清单-1-npm-的-package-json-文件) 中，我赋予 Node 项目一个名称 (magnus-server) 和一个版本号 (0.0.1)。我还将 Express 2.4.6 版本声明为一个依赖关系。NPM 的一个好处就是它会获取 Express 所有的传递依赖关系，迅速加载 Express 所需的其他所有第三方 Node 库。

通过 `package.json` 定义项目依赖关系之后，就可以通过在命令行输入 `npm install` 来安装所需的包。您会看到 NPM 安装 Express 和依赖关系，如 _connect_ 、 _mime_ 等等。

### 编写网络应用程序

我们通过创建一个 JavaScript 文件来编写示例应用程序；我将其命名为 web.js，但实际上可以随便起一个名字。在您最喜欢的编辑器或 IDE 中打开此文件；例如，您可以使用 Eclipse JavaScript 插件 JSDT（参阅 参考资料 ）。

在此文件中，添加清单 2 中的代码：

##### Magnus Server：第一步

```
var express = require('express');

var app = express.createServer(express.logger());

app.put('/', function(req, res) {
res.contentType('json');
res.send(JSON.stringify({ status:"success" }));
});

var port = process.env.PORT || 3000;

app.listen(port, function() {
console.log("Listening on " + port);
});

```

Show moreShow more icon

这一小段代码发挥了很大的作用，因此我从头讲起。首先，如果想要在 Node 中使用第三方库，就要用到 `require` 短语；在 [Magnus Server：第一步](#magnus-server：第一步) 中，我们 _需要_ Express 框架，并通过 `express` 变量获得其句柄。下一步，我们通过 `createServer` 调用创建一个应用程序实例，它会创建一个 HTTP 服务器。

然后我们通过 `app.put` 定义一个端点。在本例中，我们定义一个 HTTP `PUT` 作为在应用程序根部 (`/`) 监听的 HTTP 方法。 `put` 调用有两个参数：路由和调用路由时相应的回调。第二个参数是在运行时端点 `/` 被启动时的调用函数。记住，此回调就是 Node 所谓的事件驱动或 _事件_ I/O，即异步调用回调函数。终端可以同时处理大量请求而不必手动创建线程。

作为端点定义的一部分，我们创建处理 `/` 的 `PUT` 逻辑。为了简单起见，我们将响应类型设置为 JSON，然后发送一个简单的 JSON 文档： `({"status":"success"})` 。请注意这里的 `stringify` 方法恰到好处，它会接收哈希表然后将其转换成 JSON 格式。

##### JavaScript 和 JSON

JSON 和 JavaScript 是同胞兄弟，这种关系延续到了 Node 中。在 Node 应用程序中解析 JSON 无需特别的库或结构；您可以使用与对象图类似的逻辑调用。简而言之，Node 对待 JSON 就像自有的类型一样，这就使得编写基于 JSON 的 Web 应用程序变得非常简单。

下一步，创建用来表示应用程序所监听的端口的变量；可以通过获取 `PORT` 环境变量或直接设置为 3000 来完成。最后，通过调用 `listen` 方法来启动此应用程序。我们再次传入一个回调函数，它将会在应用程序启动并运行至将消息打印到控制台时被调用（本例中，是 `standard out` ）。

### 试一下！

这个完美的应用程序会对所有 `PUT` 作出响应，因此只要在命令行输入 `node web.js` 即可运行。如果您想要进一步测试该应用程序，我建议您下载 WizTools.org 的 [RESTClient](http://code.google.com/p/rest-client/) 。有了 RESTClient，您就可以通过对 `http://localhost:3000` 执行 HTTP `PUT` 以快速测出 Magnus Server 是否工作正常。如果正常，您会看到表示执行成功的 JSON 响应。（参阅 参考资料 了解更多有关安装和使用 RESTClient 的信息。）

## 在 Express 中处理 JSON

JavaScript 和 JSON 紧密关联，这使得在 Express 中管理 JSON 变得非常简单。在这一节，我们将在 [Magnus Server：第一步](#magnus-server：第一步) 的主应用程序中添加一些代码，从而能获取传入的 JSON 文档并将其打印到 `standard out` 。然后，我们将所有内容持久化为 MongoDB 实例。

传入的文档与清单 3 类似（请注意，简单起见，我省略了位置信息）：

##### Freddie Fingers 的免费午餐！

```
{
"deal_description":"free food at Freddie Fingers",
"all_tags":"free,burgers,fries"
}

```

Show moreShow more icon

清单 4 添加了解析传入文档的功能：

##### 用 Express 解析 JSON

```
app.use(express.bodyParser());

app.put('/', function(req, res) {
var deal = req.body.deal_description;
var tags = req.body.all_tags;

console.log("deal is :"  + deal + " and tags are " + tags);

res.contentType('json');
res.send(JSON.stringify({ status:"success" }));
});

```

Show moreShow more icon

请注意， [用 Express 解析 JSON](#用-express-解析-json) 中包含了一行代码指引 Express 使用 `bodyParser` 。这会让我们能轻松地（我是说 _轻松地_ ）获取传入的 JSON 文档的属性。

`put` 回调函数中的代码是用来获取传入文档的 `deal_description` 和 `all_tags` 属性值。请注意，我们获得请求文档的每个元素时是多么轻松：在本例中， `req.body.deal_description` 获取 `deal_description` 的值。

### 测试一下！

你也可以对此实现进行测试。关闭 `magnus-server` 实例并重启，然后使用一个 HTTP `PUT` 将一个 JSON 文档提交到 Express 应用程序中。首先，您会看到一个成功响应。其次，Express 会把您提交的值发送到 `standard out` 。通过使用我的 Freddie Fingers 文档，我获取了输出结果。

```
deal is : free food at Freddie Fingers and tags are free, burgers, fries.

```

Show moreShow more icon

## Node 的持久性

我们已经拥有了一个工作正常的应用程序，它能接收和解析 JSON 文档，并返回响应。现在要做的就是增加一些持久性逻辑。由于我偏爱 MongoDB（参见 参考资料 ），因此我选择通过 MongoDB 实例持久化数据。为了让处理过程更简单一些，我们将会用到第三方库 [Mongolian DeadBeef](https://github.com/marcello3d/node-mongolian) ，我们将用它来存储传入 JSON 文档的值。

Mongolian DeadBeef 是众多的用于 Node 的 MongoDB 库中的一个。我之所以选择它是因为我觉得它名字很有趣，还因为它对本机 MongoDB 驱动的监控功能让我非常满意。

现在，您已经知道，使用 Mongolian DeadBeef 的第一步是升级 `package.json` 文件，如清单 5 所示：

##### 添加 JSON 解析

```
{
"name":"magnus-server",
"version":"0.0.1",
"dependencies":{
    "express":"2.4.6",
    "mongolian":"0.1.12"
}
}

```

Show moreShow more icon

由于我们将会使用 MongoDB 数据存储，因此需要通过运行 `npm install` 来更新项目的硬件依赖关系。为了提高 Mongolian DeadBeef MongoDB 驱动的性能，我们还需要安装本机 C++ bson 解析器，NPM 能指导我们完成这一过程。

在开始使用 Mongolian DeadBeef 之前，在现有的实现中再添加一个 `require` ，然后连接到所需的 MongoDB 实例上（如清单 6 所示）。在本示例中，我们将连接到由 MongoHQ 托管的实例上，它是 MongoDB 的云供应者。

##### 将 Mongolian DeadBeef 添加到 magnus-server

```
var mongolian = require("mongolian");
var db = new mongolian("mongo://a_username:a_password@flume.mongohq.com:23034/magnus");

```

Show moreShow more icon

在 `PUT` 回调函数中，我们持久化来自传入的 JSON 文档的值，如清单 7 所示：

##### 添加 Mongolian insert 逻辑

```
app.put('/', function(req, res) {
var deal = req.body.deal_description;
var tags = req.body.all_tags;

db.collection("deals").insert({
     deal: deal,
     deal_tags: tags.split(",")
})

res.contentType('json');
res.send(JSON.stringify({ status:"success" }));
});

```

Show moreShow more icon

仔细看，您会发现 `insert` 语句和 MongoDB shell 中的 insert 一样。这不是巧合，MongoDB 的 shell 也使用了 JavaScript！因此，我们能够轻松地持久化拥有两个字段的文档： `deal` 和 `deal_tags` 。请注意，我们通过对 `tags` 字符串使用 `split` 方法将 `deal_tags` 设置为数组。

### 可以测试吗？（当然可以！）

如果您想要测试一下（有谁不想呢？），那么重启实例，再发送一个 JSON 文档，然后再检查 MongoDB 中的 `deals` 集合。您应该会看到一个与您发送的 JSON 文档几乎一模一样的文档。

##### 重启实例

```
{
deal:"free food at Freddie Fingers",
deal_tags:["free", "burgers", "fries"],
_id:"4e73ff3a41258b7423000001"
}

```

Show moreShow more icon

## 结束语（成功！）

如果您以为我偷懒，这么快就结束这篇 Node.js 简介，那么我要告诉您：确实结束了！虽然只写了 20 多行代码，但已经足以构成一个完整、稳定的应用程序，这就是 Node 的美妙之处。编写和理解代码都非常容易，而且异步回调函数使其功能非常强大。应用程序编写完成后，就可以部署在任意个 PaaS 提供程序上，从而达到最大的可扩展性。

参见参考资源部分了解本文中所讨论技术的更多信息，包括 Node.js、 MongoDB 和 PaaS 选项，如 Google App Engine、Amazon 的 Elastic Beanstalk 和 Heroku。