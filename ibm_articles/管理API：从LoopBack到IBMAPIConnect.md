# 管理 API：从 LoopBack 到 IBM API Connect
在 API Connect 中导入并管理 LoopBack API

**标签:** API 管理,IBM API Connect,Loopback,Node.js

[原文链接](https://developer.ibm.com/zh/articles/manage-your-apis-from-loopback-to-ibm-api-connect/)

Deepak Rajamohan

发布: 2020-09-23

* * *

自 [LoopBack 4](https://loopback.io/) 采用 OpenAPI 3.0 规范之后，我们的团队临时创建了不同的查询格式、编码样式和其他功能。LoopBack 4 用户不必为其 API 手工编写规范代码，而是从控制器中提供的端点注释中自动生成代码。

一般而言，API 管理就是管理 API 生命周期。您可以将从 LoopBack 应用程序生成的 OpenAPI 规范导入到任何 API 管理产品，您可以使用这些产品管理和保护 API 并使用扩展来扩充 API。

[IBM API Connect](https://www.ibm.com/cn-zh/cloud/api-connect) 是一种集成 API 管理解决方案，包含 API 网关、开发者门户网站和分析功能。借助其提供的选项，可使用组合件编辑器 GUI 轻松地将网关策略扩展配置到 API 规范中。

本文展示了如何将 LoopBack 与 API Connect 集成，还展示了如何在 API Connect 中导入并管理从 LoopBack 应用程序创建的 API。

## LoopBack 中的 OpenAPI 3.0 支持

[@loopback/openapi-v3](https://github.com/strongloop/loopback-next/tree/master/packages/openapi-v3) 程序包添加了用于支持 OpenAPI 扩展的 `OASEnhancer` 接口，并将控制器的元数据转换为 OpenAPI 3.0 规范。因此，从 LoopBack 应用程序生成的 JSON 或 YAML 格式的 OpenAPI 规范，可以通过额外的规范（如规范扩展和供应商扩展）加以增强。

[@loopback/apiconnect](https://github.com/strongloop/loopback-next/tree/master/extensions/apiconnect) 程序包对 LoopBack 进行了扩展，使其能够与 IBM API Connect 集成。它随附的 `ApiConnectComponent` 中添加了 `OASEnhancer` 扩展，以将 `x-ibm-configuration` 提供给 LoopBack 应用程序生成的 OpenAPI 规范。

## IBM API Connect 中的 OpenAPI 3.0 支持

[IBM API Connect](https://www.ibm.com/cn-zh/cloud/api-connect) 现在支持 OpenAPI 3.0 规范，只是存在一些 [限制](https://www.ibm.com/support/knowledgecenter/en/SSMNED_2018/com.ibm.apic.toolkit.doc/rapic_oai3_support.html)。您可使用 API Connect CLI 命令，在已知 API Connect 命令行语法的情况下，它支持 OpenAPI 3.0 管理其生命周期，这与 OpenAPI 2.0 类似。

## ToDo 示例应用程序的演示

现在，我们来看一下如何将 LoopBack 应用程序创建的 REST API 导入到 API Connect 中进行 API 管理。出于演示目的，我将使用 [ToDo 示例](https://github.com/strongloop/loopback-next/tree/master/examples/todo) 作为起始点。

### 第 1 步：在本地安装 ToDo 应用

下载 ToDo 示例并安装 `@loopback/apiconnect`

```lang-sh
$ lb4 example todo
$ cd loopback4-example-todo
$ npm i --save @loopback/apiconnect

```

Show moreShow more icon

打开 `/src` 文件夹中的 `application.ts` 文件，并按如下所示进行编辑：

1. 添加以下 import 语句：





    ```lang-ts
         import {ApiConnectComponent, ApiConnectBindings, ApiConnectSpecOptions} from '@loopback/apiconnect';

    ```





    Show moreShow more icon

2. 将以下内容添加到 `TodoListApplication` 类的 `constructor()` 函数：





    ```lang-ts
         const apiConnectOptions: ApiConnectSpecOptions = {
           targetUrl: 'http://todo-app.cloud.url/apic$(request.path)',
         };
         this
           .configure(ApiConnectBindings.API_CONNECT_SPEC_ENHANCER)
           .to(apiConnectOptions);
         this.component(ApiConnectComponent);

    ```





    Show moreShow more icon


上面提到的 targetURL `todo-app.cloud.url` 必须替换为 ToDo 示例应用在部署到云平台之后的 URL。

### 第 2 步：生成 API 规范

使用 `npm start` 在本地启动应用程序。在浏览器中使用以下 URL 生成原始规范：`http://localhost:3000/openapi.json`。在本地将此规范保存为 `todo-api.json`。

保存 JSON 文件后，需要手动执行一些步骤。

1. 确保 `info.x-ibm-name` 中的名称不包含特殊字符。可能需要将现有名称重命名为 `loopback-example-todo`。

2. 将 `servers` url 更改为 `/`，即：





    ```
    "servers": [
    {
         "url": "/"
    }
    ],

    ```





    Show moreShow more icon

    这将确保门户网站 UI 能够在测试环境中正确呈现端点。


### 第 3 步：构建并部署 LoopBack 应用程序

1. 构建 LoopBack 应用程序并创建 Docker 镜像。





    ```lang-sh
        npm run build
        docker build -t loopback4-example-todo .

    ```





    Show moreShow more icon

2. 标记 Docker 镜像并将其推送到公共或专用存储库。

     例如，对于 [IBM Container Registry](https://cloud.ibm.com/docs/Registry?topic=registry-getting-started#gs_registry_images_pushing)，您会按如下所示标记和推送镜像：





    ```lang-sh
        docker tag loopback4-example-todo us.icr.io/loopback4-example-todo:1.2.0
        docker push us.icr.io/loopback/loopback4-example-todo:1.2.0

    ```





    Show moreShow more icon

3. 部署 Docker 镜像


现在您已经有了 Docker 镜像，就可以将其部署到任何云平台产品。如果您有 IBM Cloud 帐户，可以 [在默认的名称空间中快速部署以上镜像](https://cloud.ibm.com/docs/containers?topic=containers-images)。

## 第 4 步：将 API 规范导入到 API Connect

遵循此 [演示](https://developer.ibm.com/apiconnect/2019/10/30/manage-and-enforce-openapi-v3-oai-v3/)，以了解如何将 OpenAPI 3.0 逐步导入到 IBM API Connect。

另外，您可查看如何 [使用 CLI 将产品发布到 API Connect](https://www.ibm.com/support/knowledgecenter/SSMNED_2018/com.ibm.apic.cliref.doc/apic_products_publish.html)。

```lang-sh
apic products:publish todo-api.json --org orgname --catalog catalog --server servername

```

Show moreShow more icon

## 结束语

本文演示了如何在 LoopBack 中添加 OpenAPI 3.0 支持，以及 API Connect OpenAPI 增强程序如何简化将 OpenAPI 3.0 规范导入到 API Connect 的操作。希望您可以遵循这些步骤，将 LoopBack 应用程序部署到 Docker 镜像中，并将 API 导入到 API Connect 中。

## 参与进来

LoopBack 的成功取决于您。感谢您的一贯支持和积极参与，让 LoopBack 带来更好、更有意义的 API 创建体验。您可以通过以下方式加入我们并帮助改进项目：

- [报告问题](https://github.com/strongloop/loopback-next/issues)。
- [贡献](https://github.com/strongloop/loopback-next/blob/master/docs/CONTRIBUTING.md) 代码和文档。
- [针对某个 “good first issue”（比较容易修复的问题）创建拉取请求](https://github.com/strongloop/loopback-next/labels/good%20first%20issue)。
- [加入](https://join.slack.com/t/loopbackio/shared_invite/zt-8lbow73r-SKAKz61Vdao~_rGf91pcsw) LoopBack Slack 社区。

本文翻译自： [Manage your APIs: From LoopBack to IBM API Connect](https://developer.ibm.com/articles/manage-your-apis-from-loopback-to-ibm-api-connect/)（2020-08-07）