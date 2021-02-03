# Node.js 学习路径概述
获取 Node.js 教程系列的简要概述

**标签:** Node.js

[原文链接](https://developer.ibm.com/zh/articles/learn-node-unit-1-overview-nodejs-learning-path/)

J Steven Perry

发布: 2018-10-11

* * *

## 简介

您是否想在服务器上运行 JavaScript？本教程正是为您而打造。不过，通过这一学习路径，您将会了解到 Node.js 不仅仅是“服务器上的 JavaScript”。

作为一个主题，Node 可谓博大而精深。选择有限数量的主题来介绍如此庞大的研究对象是一项艰巨的任务。因为无法预测每个新的 Node 开发人员需要什么，所以我就问自己：“在开始 Node 旅程之前，我想知道什么？”

最终便产生了这个学习路径。希望您能喜欢。

## 为什么选择 Node.js

### Node 十分流行

除了作为 [全球最流行的编程语言](https://www.tiobe.com/tiobe-index/) 之一外，JavaScript 不仅功能强大而且还易于学习（但我不会撒谎，它很难精通）。除了其他方面，Node 还是服务器上的 JavaScript。

作为一个平台，Node.js 是现存最流行的服务器平台之一，也是 [MEAN 堆栈](https://www.codingdojo.com/what-is-the-mean-stack) 的支柱。
这意味着，如果您了解 Node，您的工作前景将会很美好。

### Node 十分强大

Node.js 使用非阻塞 I/O 模型和异步编程风格。虽然 JavaScript 是一种单线程语言，但对于表现良好的 Node 应用而言，这不会产生重大问题。JavaScript Web 开发人员习惯于使用回调、Promise 和新的异步/等待语法，在浏览器中进行异步编程。Node 可为服务器带来这一体验。

Node 还具有争夺可扩展性的秘密武器：事件循环。通过使用异步编程技术，结合在后台不同线程中处理 I/O 请求的事件循环，Node 可以实现 [荒谬模式](https://www.wired.com/2015/07/teslas-new-ludicrous-mode-makes-model-s-supercar/) 可扩展性（这是不容忽视的一点）。

### Node 是一个社区

Node.js 社区是一个互动、活跃且开放的社区。在 Node 社区中通过 [npm 注册表](https://www.npmjs.com) 实现代码共享是很常见的事情，您可以在此找到许多要在应用中使用的代码，以及配套的文档和源代码。

想要 [参与](https://nodejs.org/en/get-involved/)？您可以通过多种方式为 Node 社区中已广泛开展的工作做出贡献。

## 此学习路径的结构

此学习路径分为两个主要部分：

- 第一部分：学习 Node.js – 第 2-9 单元

    - 安装软件 – 第 2 单元
    - Node.js 编程模型 – 第 3-6 单元
    - npm 和包 – 第 7-8 单元
    - 测试 – 第 9 单元
- 第二部分：应用所学知识：Node.js 实践 – 第 10-14 单元


## 学习路径的源代码

Node.js 学习路径设有一个 [GitHub 代码库](https://github.com/jstevenperry/IBM-Developer/tree/master/Node.js)，其中包含每个单元的每个示例，您可以自行运行这些示例并按照相关资料操作。

此外，第一部分中的大多数单元都配有小测试和一些编程练习。每个编程练习的答案都包含在 GitHub 存储库中，以便在您遇到困难时给予帮助。

## 前提条件

为了充分利用此学习路径，您应该熟悉 JavaScript。如果您不熟悉 JavaScript，但有使用 C++、Java、C#、PHP、Python 等其他编程语言的经验，那么应该也没什么问题。

如果您根本没有编程经验，那么可能很难完成这个学习路径。但幸运的是，IBM Developer 上有很多 [其他很棒的资源](https://developer.ibm.com/zh/technologies/)，有助于您做好准备以完成此学习路径。

您必须能够在计算机上安装软件。同时还应该能够熟练地在您特定的平台上使用命令行。如果您使用的是 MacOS（就像我一样），这里指的就是终端窗口。在 Windows 上是指命令提示符，在 Linux 上则是指命令行。

现在，让我们来谈谈学习路径本身。

## 第一部分：学习 Node

您将学习如何使用 Node.js。您将安装完成学习路径所需的软件，然后学习 Node.js 的理论，学习如何编写 Node 应用代码以及如何编写测试。

以下是每个单元的简短摘要。

### 第 2 单元：安装 Node.js 和 npm 软件

本教程将展示如何安装 `node` 程序本身以及用于管理 Node 项目的 `npm`。

在本教程中，我将展示在计算机上获取 Node 软件的三种方法。活跃的 Node开发者社区意味着它在不断发展。例如，在编写此学习路径花费的 10 周时间里，Node 经历了五次修订（在我写这篇文章时，从 10.4.1 升级到 10.9.0）。

您必须平衡 Node 安装过程的简便程度与日后升级的难度。因此，我会按照安装过程的简便程度介绍安装选项，这恰好与升级的容易程度相反。我为您提供了必要的信息，帮助您选择合适您的安装方法。

- 从 nodejs.org 下载 – 下载适用于您平台的安装包，并将其安装在计算机上。
- [Homebrew](https://brew.sh)（仅限 MacOS）- 下载并安装 Homebrew，然后使用 `brew install node` 安装 Node
- [Node Version Manager (nvm)](https://github.com/creationix/nvm) – 下载并安装 nvm，然后使用 ‘nvm install\` 安装 Node，以便安装您需要的 Node 版本

您可以观看视频，我将在其中指导您完成每种类型的安装。

在完成第 2 单元后，您就可以编写 Node 程序了。

[阅读第 2 单元](/zh/tutorials/learn-nodejs-installing-node-nvm-and-vscode/)

### 第 3 单元： Node.js 之旅

为此，在第三个教程中，我将简要介绍 Node 的所有“移动部件”：

- Node 运行时，包括：
    - Node API：JavaScript 实用程序，如文件和网络 I/O，以及其他许多实用程序，如加密和压缩
    - Node 核心：一组实现 Node API 的 JavaScript 模块。（显然，一些模块依赖于 libuv 和其他 C++ 代码，但这是一个实现细节）。
    - JavaScript 引擎：Chrome 的 V8 引擎，一种快速的 JavaScript 到机器代码编译器，用于加载、优化和运行 JavaScript 代码
    - 事件循环：使用名为 libuv 的事件驱动的非阻塞 I/O 库实现，使其实现轻量级和高效（且可扩展）
- Read-Eval-Print Loop (REPL)
- 非阻塞 I/O 模型
- npm 生态系统

但是，这并不是所有理论。您编写的代码将在 REPL 中运行，并展示 Node 的非阻塞 I/O 模型。

您可以通过视频查看 REPL 实践，甚至还可以通过小测试来检验自己对单元中资料的理解程度。

作为 Node 程序员，如果对 Node 架构和非阻塞 I/O 模型的了解越深入，您就越有可能技高一筹。在完成本 Node 教程后，您将会基本了解 Node 的后台工作方式，知道所有移动部件如何协同工作才能使 Node 应用正常运行。

[阅读第 3 单元](/zh/tutorials/learn-nodejs-tour-node/)

### 第 4 单元：Node.js 概念

第 4 单元将深入介绍 Node.js 概念，通过比较和对比两种不同类型的业务：银行和咖啡店，您可以了解有关 Node 非阻塞 I/O 模型和异步编程风格的更多信息。要想成为一名高效的 Node 程序员，就应时刻牢记这些概念。

然后，您将仔细查看Node 的模块系统：如何使用和创建模块、如何使用模块，以及何时将它们用于同步和异步编程。

我将向您说明如何使用 Chrome V8 概要分析器查找代码中的性能热点，您将在此对同步和异步编程示例进行概要分析，这些示例说明了 Node 非阻塞 I/O 模型的强大功能以及常规的异步编程风格。

在第 4 单元的最后，将会概述在整个学习路径中广泛使用的 Promise。

完成本教程后，您将能够很好地掌握非阻塞 I/O、异步编程、Chrome V8 引擎和负载测试。

您可以查看此单元随附的视频，我在视频中介绍了 Node 模块的工作方式、非阻塞 I/O 的示例以及如何在 Node 应用上执行负载测试。另外，还提供了一个小测试来检验您的理解程度，以及一些编程练习，以便您可以编写代码。

[阅读第 4 单元](/zh/tutorials/learn-nodejs-node-basic-concepts/)

### 第 5 单元：Node 事件循环

Node 可扩展性的关键在于事件循环，它支持非阻塞 I/O 模型。

基于先前单元中的概念，您将进一步了解什么是事件循环、使用事件循环的原因以及如何对定时器、I/O、process.nextTick() 和 Promise 使用回调来绘制出事件循环的阶段。我谈论的其他主题包括：

- Node 事件以及如何创建一个定制事件
- 流以及如何使用它们来绘制出事件循环的阶段
- 如何将自变量传递至计时器回调

完成本教程之后，您将能够大体上全面掌握 Node 内部结构，特别是事件循环。

我在第 5 单元中提供了一个视频，帮助说明此单元中出现的一些概念。另外，还提供了一个小测试来帮助您评估学习进度，以及一些编程练习，以便您也可以编写代码。

[阅读第 5 单元](/zh/tutorials/learn-nodejs-the-event-loop/)

### 第 6 单元：您的第一个 Node.js 应用

在第 6 单元中，您将加入一个模拟的 Node 项目团队，在此完成编写应用核心，即帮助用户创建购物清单。本单元向您介绍作为真实项目的开发者应如何使用 Node。

别担心，您不必从头开始编写整个应用。在此特定场景中，您将进入部分创建的项目并完成该项目。

在该场景中，您将完成一组 REST 服务，支持贵公司已出售给客户的购物清单应用。随着学习路径的逐步推进，您将不断完善购物清单应用，了解功能测试，知道如何将数据加载到数据库中，以及如何在 VSCode 调试器中调试代码。

第 6 单元以视频结束，我会在这个视频中向您介绍如何设置项目。我还会带您浏览一下代码，介绍如何运行功能测试，进而确定应用何时按照规范运行。如果顺利通过功能测试，您的代码便可正常运行！

另外，还有一个小测试来检验您的理解程度，并且作为练习，您需要完成该应用（不必担心，如果您遇到困难，我会提供解决方案）。

在完成本教程后，您将编写完成第一个 Node.js 应用，并准备好深入了解 Node 编程知识。

[阅读第 6 单元](/zh/tutorials/learn-nodejs-your-first-node-application/)

### 第 7 单元：Node Package Manager (npm)

包管理是 Node 编程不可或缺的一部分，因此有必要单开一个单元。npm 不仅仅是包管理器程序 (`npm`)，还是文档、包和整个生态系统。在第 7 单元中，我会向您介绍所有这些内容及其他更多内容。

务必要参加小测试来检验您的理解程度。

完成本教程后，您将能够很好地掌握 npm。

[阅读第 7 单元](/zh/tutorials/learn-nodejs-node-package-manager/)

### 第 8 单元：Node 依赖管理

在第 8 单元中，您将学习如何管理依赖。

Node 依赖管理的核心是 `package.json`，它是 Node 项目的清单。它包括：

- 项目元数据，如名称、版本等
- _SemVer_（Semantic Versioning，语义化版本控制）格式的依赖信息
- 其他众多内容

我将展示如何创建一个新的 Node 项目，其中 `npm` 工具将为您创建一个骨架 `project.json`。

您还将了解语义化版本控制 ( [SemVer](https://semver.org))，明白为何它对 Node 依赖管理如此重要，以及为何必须理解它才能指定您的依赖关系，进而让 `npm` 正确管理这些依赖关系。

第 8 单元包含了一个视频，我会在该视频中介绍 `package.json` 文件的必需元素以及 semver 在依赖管理中的作用。另外，还提供了一个小测试来检验您对第 8 单元概念的理解程度。

在完成本教程后，您将能够全面了解 Node 依赖管理，以及 `package.json` 在 Node 项目中的作用。

[阅读第 8 单元](/zh/tutorials/learn-nodejs-manage-packages-in-your-project)

### 第 9 单元：测试

测试是软件开发的基础部分。在第 9 单元中，我会向您介绍几个工具，帮助您测试 Node 应用以捕获任何 bug：

- [Mocha](https://mochajs.org) – 测试框架
- [Chai](http://www.chaijs.com) – 断言库
- [Sinon](https://sinonjs.org) – 测试替代库

通过结合使用这三个工具，可帮助您构建出色的无瑕软件。

您还将学习如何在代码上运行名为 _linter_ 的工具，用于分析可能出现的错误，并确保您的代码遵循一致的行业标准编码样式。

第 9 单元包含了一个视频和一个小测试。

在完成本教程后，您就可以应用从本学习路径的前 9 个单元学到的 Node 知识。

[阅读第 9 单元](/zh/tutorials/learn-nodejs-unit-testing-in-nodejs/)

## 第二部分：应用所学知识

在第二部分，您可以亲自实践所学知识，并解决作为一名专业的 Node 开发者肯定会遇到的一些问题。

以下是第二部分中每个单元的简短摘要。

### 第 10 单元：日志记录

应用日志记录是一个需要清晰策略的领域。`console.log()` 适合业余爱好者，但专业的 Node 开发者应使用日志记录包。

在第 10 单元中，您将学习如何为 Node 应用使用两个非常流行的日志记录包：

- [Winston](https://www.npmjs.com/package/winston)
- [Log4js](https://www.npmjs.com/package/log4js)

我之所以挑选了这两个日志记录包，是因为作为一名专业的 Node.js 开发者您可能会用到它们。当然还有其他一些日志记录包，我们先着重关注这两个包，并在此基础上积累相关经验。

该单元还附带一个视频，您可以观看视频，了解这两个日志记录器的实际运用。

[阅读第 10 单元](/zh/tutorials/learn-nodejs-winston/)

### 第 11 单元：ExpressJS

[ExpressJS](https://expressjs.com/) 是一个简约的高性能开源 Web 框架，可能是用于 Node.js 应用的最流行的 Web 框架。由于良好地证明了其性能和流行程度，Express 也是 [LoopBack](https://loopback.io/) 等其他流行框架的支柱。

在第 11 单元中，您将使用我在第 6 单元中放置在购物清单应用的 REST 服务上的 Express Web 前端。您将了解 Express 中间件、路由以及使用 [Pug](https://pugjs.org/api/getting-started.html)（以前称为 Jade）查看模板处理。

在使用 Express 中间件时，您将学习使用 [`express-validator`](https://github.com/express-validator/express-validator) 进行数据验证和清理的技术。

该单元还附带一个视频，您可以观看视频，了解购物清单 UI 的实际运行情况。

[阅读第 11 单元](/zh/tutorials/learn-nodejs-expressjs/)

### 第 12 单元：MongoDB

[MongoDB](https://www.mongodb.com) 是全球最受欢迎的 NoSQL 数据库之一，因此，作为一名专业的 Node 开发者，您很可能会用到它（或类似它的 NoSQL 数据库）。

在本教程中，您将再次使用购物清单应用，只是这次记录系统是 MongoDB。我会向您简要介绍 MongoDB，展示如何在 MacOS、Windows 和 Linux 上安装此数据库，并介绍购物清单应用的数据访问组件（Data Access Components）如何使用 MongoDB。

此外，我还会额外介绍一下 [Mongoose](https://mongoosejs.com)，这是一个非常流行的 MongoDB 对象数据映射 (ODM) 实用程序。请务必查看该单元源码中包含的 Mongoose 源文件。作为练习，您可以随意使用它们，如果您觉得自己能够接受挑战，也可以将 Mongoose 插入购物清单应用中。

第 12 单元还提供了一个视频，我在视频中简要介绍了 MongoDB，以及一些 MongoDB 术语，并且还展示了 MongoDB API 的一些实际运用。

[阅读第 12 单元](/zh/tutorials/learn-nodejs-mongodb/)

### 第 13 单元：调试和概要分析

有时，我们编写的代码会有一些难以跟踪的错误，这时源码级调试就会派上用场。

在其他一些时候，代码虽然在功能上是完美的，但却无法良好执行（或扩展），概要分析器（profiler）这时就会派上用场。

对于任何开发者来说，能够执行源码级调试和概要分析都是一项非常不错的技能。

在第 13 单元中，我将向您介绍如何使用 VSCode 源码级调试器逐行走查代码。学习如何检查变量，设置监视表达式，甚至是如何调试 Mocha 测试。

我还会向您介绍如何使用 [Apache Bench](https://httpd.apache.org/docs/2.4/programs/ab.html)（我会在第 4 单元中介绍）和名为 [clinic](https://github.com/nearform/node-clinic) 的概要分析工具来查找代码中的热点。

概要分析是 Node 开发者必备的技能。为什么？异步编程并不容易。有时，这会导致单凭查看代码，很难（或根本不可能）找出代码中的热点（即性能瓶颈）。您 _必须_ 运行概要分析器。

这里还提供一个视频展示如何设置 VSCode 来执行源码级调试，以及如何使用 clinic doctor 和 bubble 概要分析器来查找性能瓶颈。

[阅读第 13 单元](/zh/tutorials/learn-nodejs-debugging-and-profiling-node-applications/)

### 第 14 单元：Cloudant

Cloudant 是 [IBM Cloud](https://cloud.ibm.com/catalog/services/cloudant?cm_sp=ibmdev-_-developer-articles-_-cloudreg) 中提供的高性能、高可用性、以文档为中心的 NoSQL 数据库即服务 (DBaaS) 产品。

在与 Node 非阻塞 I/O 模型和异步编程模型（如微服务）配合使用的应用程序类型中，NoSQL 数据库非常受欢迎。

##### 免费试用 IBM Cloud

利用 [IBM Cloud Lite](https://cloud.ibm.com/registration?cm_sp=ibmdev-_-developer-articles-_-cloudreg) 快速轻松地构建您的下一个应用程序。您的免费帐户从不过期，而且您会获得 256 MB 的 Cloud Foundry 运行时内存和包含 Kubernetes 集群的 2 GB 存储空间。 [了解所有细节](https://developer.ibm.com/dwblog/2017/building-with-ibm-watson/) 并确定如何开始。

大多数企业都开始了云迁移进程，因为托管自己的服务器不仅成本高昂，还会使您面临可用性问题。通过使用云计算，您只要为所需的云服务买单，因此就不会因无力承受的昂贵硬件而望洋兴叹。随着越来越多的公司迁移到云端，您需要确保知道如何创建可在云上扩展的应用。

第 14 单元即介绍了如何将购物清单应用的数据迁移到 IBM Cloud 中。您将学习如何设置自己的 Cloudant 服务，并且将能够使用完成的完整工作应用代码。

这里也提供一个视频，向您介绍如何设置 Cloudant 服务，帮助您熟悉 Cloudant 仪表板，了解如何通过仪表板直接使用数据库。

[阅读第 14 单元](/zh/tutorials/learn-nodejs-node-with-cloudant-dbaas/)

## 结束语

希望您对本学习路径所涵盖的内容已有了大致的了解，并且已经迫不及待地想要一显身手。

## 视频

在该视频中，我将指导您完成学习路径，解释您在这一过程中将要学习的内容。

本文翻译自： [Overview of Node.js Learning Path](https://developer.ibm.com/articles/learn-node-unit-1-overview-nodejs-learning-path/)（2018-10-11）