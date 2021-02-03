# express-no-stress 生成器：使用 Express.js 和 Swagger 构建优秀的 API
了解如何使用 Node.js 构建 API 应用程序

**标签:** API 管理

[原文链接](https://developer.ibm.com/zh/articles/cl-express-no-stress-generator/)

Carmine DiMascio

发布: 2017-08-02

* * *

在这篇短文中，我将介绍如何使用 [Node.js](https://nodejs.org/en/) 以及一流的技术和方法搭建一个 API 应用程序。我将介绍 [express-no-stress](https://github.com/cdimascio/generator-express-no-stress) 生成器，可用它来快速搭建一个构建于 [Express.js](http://expressjs.com/) 之上的新 API 应用程序。搭建的这个应用程序提供了结构化的日志记录、API 请求验证、交互式 API 文档、环境驱动的配置、一个简单的构建/转译过程等。该应用程序也适合部署在像 [IBM Cloud](https://cloud.ibm.com/?cm_sp=ibmdev-_-developer-articles-_-cloudreg) 这样的现代云平台上。

为什么这很重要？如今，软件通常作为一项服务在网络上提供，就是这些软件推动 API 经济发展。API 通过提供其他人想要的功能来创造价值：它们帮助开发人员自由创新、推动变革，无缝连接不同领域，最终提供新价值。

在构建 API 时，满足使用它们的开发人员的需求也很重要。API 必须值得信赖、可扩展、直观且易用。它们还必须最大限度减少维护所需的时间和成本。理想情况下，API 应用程序还应：

- 容易设置和自动化
- 容易在执行环境之间移植
- 适合部署到云
- 而且只需极少的系统管理和监视工作

幸运的是， [express-no-stress](https://github.com/cdimascio/generator-express-no-stress) 生成器创建的应用程序满足这些要求。express-no-stress 是我启动的一个微型项目 — 一个搭建一致的、全功能的 REST API 的 Yeoman 生成器。express-no-stress 结合了构建优秀应用程序所需的工具。它只需几秒即可生成，使您能将精力集中在重要事务上 — 编写应用程序代码。

express-no-stress 包含：

- [**Express.js**](http://expressjs.com/) — 一个基于 Node.js 的快速、开放、最简化的 Web 框架。
- [**Babel.js**](http://babeljs.io/) — 一个帮助您使用如今的下一代 JavaScript 的编译器。
- [**Bunyan**](https://github.com/trentm/node-bunyan) — 一个基于 Node.js 的简单且快速的 (JSON) 结构化日志模块。
- [**dotenv**](https://www.npmjs.com/package/dotenv) — 从 .env 加载环境变量。（12 因素原则 III）。
- [**Backpack**](https://github.com/palmerhq/backpack) — 一个最简化的 Node.js 构建系统。
- [**Swagger**](http://swagger.io/) — RESTful API 的一种简单但强大的表示。
- [**SwaggerUI**](http://swagger.io/swagger-ui/) — 用符合 Swagger 标准的 API 动态生成美观的文档和沙箱。
- **Cloud Ready** — 部署到任何云；提供了一个 [CloudFoundry](https://www.cloudfoundry.org/)/IBM Cloud 示例（参见 [自述文件](https://github.com/cdimascio/generator-express-no-stress)）。

要查看搭建的现成应用程序的实际效果，请访问 [https://express-no-stress-scaffolded-app.mybluemix.net](https://express-no-stress-scaffolded-app.mybluemix.net)。

**备注：** 搭建的这个应用程序仅提供了最简单的示例 API，需要您按自己想要的方式实现 API 代码。

要开始使用 express-no-stress，请访问 [GitHub 存储库](https://github.com/cdimascio/generator-express-no-stress)。

如果您迫不及待想使用它，可以运行：

```
npm install -g yo generator-express-no-stress
yo express-no-stress myapp

```

Show moreShow more icon

express-no-stress Yeoman 生成器使得通过 Node.js 和 Express 快速创建 Web API 变得很轻松。它是一个值得装入您的工具箱中的有用工具 — 非常适合生产项目和编外项目。如果您在微服务环境中工作（或者即使您不在其中工作），express-no-stress 有助于确保新 Node.js API 被一致地搭建并使用一组通用的可靠工具。搜出来看看！

本文翻译自： [express-no-stress: Build awesome APIs with Express.js and Swagger](https://developer.ibm.com/articles/cl-express-no-stress-generator/)（2018-09-12）