# 将 Spring 与 IBM Cloud 和 IBM 软件结合使用
利用 Spring 启动器将 IBM 技术融入到您的应用程序

**标签:** Java,Spring

[原文链接](https://developer.ibm.com/zh/articles/use-spring-with-ibm-cloud-and-ibm-software/)

Ozzy Osborne

更新: 2018-08-19 \| 发布: 2018-08-17

* * *

## 什么是 Spring？

Spring 是一个流行的 Java 企业级应用开发框架。它已经突破最初提供的功能，发展成一个使用和围绕 Spring 框架构建的整个项目家族。Spring家族项目解决广泛的技术问题，支持使用 Spring 功能构建整体解决方案。

## 什么是 Spring Boot？

“Spring Boot” 是一个 Spring 项目，通过一套丰富的约定（约定优于配置）和注释（配置即代码）取代老旧冗余的 XML 配置，从而简化了开发流程。“Spring 启动器” 通过配置注入和自动组件初始化，为各种技术的使用提供进一步定向援助（opinionated assistance）。通过采用“约定优于配置”的方法，可以最大限度减少开发人员使用新的 API 或服务所需的工作量。Spring Boot 应用还能打包为单个可执行的 JAR 文件，以便于测试和部署。

## IBM 在用 Spring 做什么？

IBM 正在积极为 IBM 软件和 IBM Cloud 服务创建各种 Spring 启动器，帮助 Spring 开发人员轻松地将 IBM 技术整合到他们的应用中。

目前存在的 Spring 启动器有：

- [Watson](https://github.com/watson-developer-cloud/spring-boot-starter)
- [MQ/JMS](https://github.com/ibm-messaging/mq-jms-spring)
- [OpenLiberty](https://openliberty.io/blog/2017/11/29/liberty-spring-boot.html)
- [Cloudant](https://github.com/cloudant-labs/cloudant-spring)

IBM 还在创建初学者工具包，帮助开发人员快速创建 Spring 应用并部署到 [IBM Cloud](https://www.ibm.com/cloud/) 上。

初学者工具包为您的云原生应用开发奠定了基础。生成的每个初学者工具包都包含所选服务、Dockerfile 以及其他部署元数据的配置和依赖关系，并且预置了监控应用健康状况的工具。

初学者工具包含有：

- “ [Java Spring Basic](https://cloud.ibm.com/developer/appservice/starter-kits/java-web-app-with-spring?cm_sp=ibmdev-_-developer-articles-_-cloudreg)” – 一个使用 Spring Boot 和 Tomcat 且具备 Kubernetes/CloudFoundry/Pipeline 支持的 Web 应用。
- “ [Java Spring Backend](https://cloud.ibm.com/developer/appservice/starter-kits/java-bff-example-with-spring?cm_sp=ibmdev-_-developer-articles-_-cloudreg)” – 一个使用 Spring Boot 和 Open API 的 Backend for Frontend (BFF, 服务于前端的后端) Java 应用。
- “ [Java Spring Microservice](https://cloud.ibm.com/developer/appservice/starter-kits/java-microservice-with-spring?cm_sp=ibmdev-_-developer-articles-_-cloudreg)” – 一个使用 Spring Boot 的 Java 微服务应用。
- 您也可以创建一个为选定的 IBM Cloud 服务预先配置的定制 Spring 应用。

本文翻译自： [Use Spring with IBM Cloud and IBM software](https://developer.ibm.com/articles/use-spring-with-ibm-cloud-and-ibm-software/) (2018-08-17）