# 针对开发者的 IBM MQ 下载项
开发者使用 IBM MQ 之前需要的所有下载项

**标签:** IBM MQ

[原文链接](https://developer.ibm.com/zh/articles/mq-downloads/)

Richard J. Coppen

发布: 2020-07-29

* * *

如果您准备开始编程，直接跳至 [开发资源](#languages-apis-and-protocols)。

有关 IBM MQ 概念的简要概述，阅读 “ [IBM MQ 基础知识](/zh/articles/mq-fundamentals)”一文。

在开发 IBM MQ 应用程序之前，您需要一个可与应用程序进行交互的队列或主题。IBM MQ 队列和主题都托管在队列管理器上。队列管理器是用于托管队列和主题的服务器。您的应用程序将作为客户端连接到 IBM MQ。

如果您是 MQ 管理员，则可以查阅 Knowledge Center 中的 [IBM MQ 产品文档](https://www.ibm.com/support/knowledgecenter/SSFKSJ_latest/com.ibm.mq.pro.doc/q002610_.htm)，然后也是在 Knowledge Center 中访问全部的 [IBM MQ 产品下载项](https://www.ibm.com/support/knowledgecenter/SSFKSJ/com.ibm.mq.helphome.doc/mq_downloads_admins.htm)。

要开发 IBM MQ 应用程序，您只需要：

- [MQ 客户端库](#get-mq-client-libraries)
- [一个队列管理器](#get-a-queue-manager-mq-server)

## 获取队列管理器（MQ 服务器）

您可以在许多不同的环境（包括不同的云）上运行队列管理器（MQ 服务器）。

您可以下载最新的免费版 IBM MQ Advanced for Developers（不适用于生产环境），其中包括具有默认队列和主题的队列管理器，以便您快速入门。

按照这些教程中的说明来开始使用队列管理器：

- 在 [容器](/zh/tutorials/mq-connect-app-queue-manager-containers/) 中。
- 在 [IBM Cloud](/zh/tutorials/mq-connect-app-queue-manager-cloud/) 中。
- 在各种操作系统上： [Linux/Ubuntu](/zh/tutorials/mq-connect-app-queue-manager-ubuntu/) 或 [Windows](/zh/tutorials/mq-connect-app-queue-manager-windows/)。对于 MacOS，使用 [MQ on Containers](/zh/tutorials/mq-connect-app-queue-manager-containers/)。

要从 IBM Fix Central 下载 MQ 服务器或 MQ 客户端，您需要使用 IBMid 来登录，因此务必先 [创建 IBMid](https://www.ibm.com/account/reg/us-en/signup?formid=urx-19776)。 另外，如果您在 Fix Central 中遇到问题，尝试使用其他浏览器。

如果您知道如何操作，或者需要其他版本或其他 IBM MQ 服务器安装程序，参阅 [面向管理员的 IBM MQ 下载项](https://www.ibm.com/support/knowledgecenter/SSFKSJ/com.ibm.mq.helphome.doc/mq_downloads_admins.htm) 页面。

## 获取 MQ 客户端库

要开发 MQ 应用程序，您需要一个队列管理器和以下三项：

1. 适用于您的语言和平台的 MQ 客户端库
2. 编译器（在 Windows 上，使用 Microsoft Visual Studio。在 MacOS 上，使用 XCode。在 Linux 上，使用 GCC。对于 Java，SDK 中已经包含了编译器。）
3. 适用于您的语言的 SDK

MQ 允许您选择不同的 API 和协议。参阅下一部分中的表格，以查看哪种语言可与哪种 API 和协议一起使用。

要从 IBM Fix Central 下载 MQ 服务器或 MQ 客户端，您需要使用 IBMid 来登录，因此务必先 [创建 IBMid](https://www.ibm.com/account/reg/us-en/signup?formid=urx-19776)。 另外，如果您在 Fix Central 中遇到问题，尝试使用其他浏览器。

您可以使用我们的一种可再分发客户端 (redist) 在部署应用程序的位置运行该应用程序。只需下载可再分发的客户端库即可。然后，您必须将应用程序与用于目标环境的可再分发库打包在一起。在 Knowledge Center 内的 [IBM MQ 文档](https://www.ibm.com/support/knowledgecenter/SSFKSJ_9.1.0/com.ibm.mq.ins.doc/q122882_.htm) 中了解如何执行此操作（并务必阅读关于再分发 MQ 客户端的 [许可信息](https://www.ibm.com/support/knowledgecenter/SSFKSJ_9.1.0/com.ibm.mq.pro.doc/q134260_.htm)）。

### 语言、API 和协议

您可以根据所使用的语言来选择 API 和协议。API 的功能和易用程度各不相同，因此在选择 API 时应该考虑您应用程序的需求。协议通常对应用程序是隐藏的，因此，除非您已经知道要使用特定协议，否则不需要担心协议的问题。

在给定的 API 和协议组合下，您可以执行的操作会有所不同，因为并非所有的 API 和协议都提供完全相同的消息传递功能。

除了 IBM 支持的库外，MQ 应用程序开发者还可以使用许多第三方库和开源库。参阅表格下方适用于您的语言的注释和段落。

语言API协议[Java](#java)JMSMQ[Java](#java)MQTTMQTT (1)[Java](#java)REST (2)HTTP/HTTPS[Java](#java)MQ Light (3)AMQP 1.0 (4)[C# (.NET)](#cprime)XMSMQ[C# (.NET)](#cprime)MQ Object OrientedMQ[Python](#python)MQ Light (3)AMQP 1.0 (4)[Python](#python)MQI（带有 pymqi） (5)MQ[Python](#python)REST (2)HTTP/HTTPS[Node.js](#node-js)MQIMQ[Node.js](#node-js)REST (2)HTTP/HTTPS[Node.js](#node-js)MQ Light (3)AMQP 1.0 (4)[Ruby](#ruby)MQ Light (3)AMQP 1.0 (4)[C++](#c-plus-plus)XMSMQ[C++](#c-plus-plus)MQIMQ[C++](#c-plus-plus)MQ Object OrientedMQ[C](#c-lang)MQIMQ[C](#c-lang)XMSMQ[COBOL](#cobol)MQIMQ[HLASM（高级汇编器）](https://www.ibm.com/support/knowledgecenter/SSFKSJ_latest/com.ibm.mq.dev.doc/q023760_.htm)MQIMQ[PL/1](https://www.ibm.com/support/knowledgecenter/SSFKSJ_latest/com.ibm.mq.dev.doc/q023790_.htm)MQIMQ[pTAL](https://www.ibm.com/support/knowledgecenter/SSFKSJ_latest/com.ibm.mq.dev.doc/q114500_.htm)MQIMQ[RPG – IBMi](https://www.ibm.com/support/knowledgecenter/SSFKSJ_latest/com.ibm.mq.dev.doc/q023780_.htm)MQIMQ

(1) MQTT 协议非常适合开发 IoT 应用程序，或者在带宽受限且需要较小消息开销的情况下使用。IBM MQ 支持 MQTT 3.1.1 协议。可以使用 [Eclipse Paho 客户端](https://www.eclipse.org/paho/) 来连接到 IBM MQ。

(2) 消息传递 REST API 可用于与队列进行交互，并且从 [V9.1.5 开始，可用于发布主题](https://www.ibm.com/support/knowledgecenter/SSFKSJ_9.1.0/com.ibm.mq.pro.doc/q133690_.htm#q133690___restapiPub)。有关更多信息，参阅 “ [IBM MQ 消息传递 REST API 入门](/zh/tutorials/mq-develop-mq-rest-api)”教程和 IBM MQ Knowledge Center 中的“ [使用 REST API 进行消息传递](https://www.ibm.com/support/knowledgecenter/SSFKSJ_latest/com.ibm.mq.dev.doc/q130940_.htm)”。

(3) 在我们的 [MQ Light GitHub 代码库](https://github.com/mqlight) 中提供了多种语言版本的 MQ Light 客户端。

(4) AMQP 是一种开放标准的线路层协议，在使用开源库的开发者中很受欢迎。IBM MQ 仅支持一部分 AMQP1.0 功能，并用于增强 MQ Light API。

(5) [pymqi](https://dsuch.github.io/pymqi/) 是一个受欢迎的第三方库。除了 [pymqi 提供的示例](https://dsuch.github.io/pymqi/examples.html) 外，另参阅 [GitHub 代码库中的 Python 样本](https://github.com/ibm-messaging/mq-dev-patterns/tree/master/Python)。

#### Java

Java（适用于 Java 的 MQ 类）和 Java 消息服务（适用于 Java 的 MQ JMS 类）现在都是很受欢迎的选择，它们适合于下列开发者：打算开发要在应用程序服务器（如 Liberty 服务器概要文件）上运行的应用程序的开发者，或者希望使用 Spring 进行开发的开发者。它们还适用于开发独立的 Java SE 和 Java EE 应用程序。适用于 Java 的 MQ 类具有与原生 MQI API 相似的 API，因此您可以开发与 MQI 更紧密结合的应用程序。 Java 消息服务 (JMS) 是一种开放标准 API，可用于开发与特定消息传递提供者紧密耦合的消息传递应用程序。适用于 Java 的 MQ JMS 类都实现了 JMS 标准。

我们提供了跨平台的库（包括“适用于 Java 的 MQ 类”和“适用于 Java 的 MQ JMS 类”），这些库打包在 com.ibm.allclient jar 文件中。

您可以选择使用以下方式来获取 MQ Java/JMS 库。

包WindowsLinuxMacPull pkg (Maven)[com.ibm.mq.allclient.jar](https://ibm.biz/mq-jms-allclient-jar)[com.ibm.mq.allclient.jar](https://ibm.biz/mq-jms-allclient-jar)[com.ibm.mq.allclient.jar](https://ibm.biz/mq-jms-allclient-jar)Redist (grab & go) MQ 下载[IBM-MQC-Redist-Java.zip](https://ibm.biz/915-IBM-MQC-Redist-Javazip)[IBM-MQC-Redist-Java.zip](https://ibm.biz/915-IBM-MQC-Redist-Javazip)[IBM-MQC-Redist-Java.zip](https://ibm.biz/915-IBM-MQC-Redist-Javazip)Redist (Fix Central)[IBM-MQC-Redist-Java](https://ibm.biz/915-IBM-MQC-Redist-Java)[IBM-MQC-Redist-Java](https://ibm.biz/915-IBM-MQC-Redist-Java)不适用Client (Fix Central)[IBM-MQ-Install-Java-All](https://ibm.biz/915-IBM-MQ-Install-Java-All)[IBM-MQ-Install-Java-All](https://ibm.biz/915-IBM-MQ-Install-Java-All)不适用

对于 JMS，您还需要：

- [JMS jar](https://mvnrepository.com/artifact/javax.jms/javax.jms-api/2.0.1)
- [Java 8 或 Java 11 JDK（从 V9.1.5 开始）](https://adoptopenjdk.net/archive.html)

对于 Spring Boot，使用 [适用于 Java 的 Spring Boot 启动器](https://github.com/ibm-messaging/mq-jms-spring)。

对于在 Java EE 上运行的应用程序，使用 [MQ Resource Adaptor](http://www.ibm.com/support/fixcentral/swg/quickorder?parent=ibm%7EWebSphere&product=ibm/WebSphere/WebSphere+MQ&release=9.1.4&platform=Windows+64-bit,+x86&function=fixId&fixids=9.1.5.0-IBM-MQ-Java-InstallRA&includeSupersedes=0&source=fc) 代替 `allclient` 库。

要使用 Java 和 REST，参阅 “ [IBM MQ 消息传递 REST API 入门](/zh/tutorials/mq-develop-mq-rest-api)”教程和 IBM MQ Knowledge Center 中的“ [使用 REST API 进行消息传递](https://www.ibm.com/support/knowledgecenter/SSFKSJ_latest/com.ibm.mq.dev.doc/q130940_.htm)”。

要将 Java 与 MQTT 一起使用，参阅 [Eclipse Paho 客户端](https://www.eclipse.org/paho/index.php?page=downloads.php)。

要将 Java 与 MQ Light 一起使用，参阅 [MQ Light GitHub 代码库](https://github.com/mqlight) 和 [Maven](https://search.maven.org/artifact/com.ibm.mqlight/mqlight-api/1.0.2016062300/jar) 中的 Java 客户端。

要了解如何使用 MQ 类进行开发，阅读 Knowledge Center 中的 [IBM MQ 文档](https://www.ibm.com/support/knowledgecenter/SSFKSJ_9.1.0/com.ibm.mq.dev.doc/q118320_.htm)。

适用于 Java 的 MQ 类已经很稳定，这意味着除了修复缺陷和确保满足系统需求外，不会再添加任何新功能。

#### C\# (.NET)

我们提供了 XMS（这是一种适用于 .NET 的类似于 JMS 的 API），还提供了一个基于 MQI 的低级别 API，称为“适用于 .NET 的 MQ 类”。

选择客户端、可再分发客户端或 NuGet 包（这样就万事俱备了）：

包WindowsLinuxMacPull pkg C# (.NET) XMS[IBMXMSDotnetClient](http://ibm.biz/mq-dotnet-xms-nuget)[IBMXMSDotnetClient](http://ibm.biz/mq-dotnet-xms-nuget)不适用Pull pkg C# (.NET)[IBMMQDotnetClient](http://ibm.biz/mq-dotnet-nuget)[IBMMQDotnetClient](http://ibm.biz/mq-dotnet-nuget)不适用Redist (grab & go) MQ 下载[IBM-MQC-Redist-Win64.zip](https://ibm.biz/915-IBM-MQC-Redist-Win64zip)[IBM-MQC-Redist-LinuxX64.tar.gz](https://ibm.biz/915-IBM-MQC-Redist-LinuxX64targz)不适用Redist (Fix Central)[IBM-MQC-Redist-Win64](https://ibm.biz/915-IBM-MQC-Redist-Win64)[IBM-MQC-Redist-LinuxX64](https://ibm.biz/915-IBM-MQC-Redist-LinuxX64)不适用Client (Fix Central)[IBM-MQC-Win64](https://ibm.biz/915-IBM-MQC-Win64)[IBM-MQC-LinuxX64](https://ibm.biz/915-IBM-MQC-LinuxX64) 或 [IBM-MQC-UbuntuLinuxX64](https://ibm.biz/915-IBM-MQC-UbuntuLinuxX64)不适用

要开始开发 .NET 应用程序，参阅我们的 GitHub 代码库中的 [IBM MQ XMS 样本](https://github.com/ibm-messaging/mq-dev-patterns/tree/master/dotnet)。 要了解如何开发 .NET 应用程序阅读 Knowledge Center 中的 [IBM MQ 文档](https://www.ibm.com/support/knowledgecenter/SSFKSJ_9.1.0/com.ibm.mq.dev.doc/q029250_.htm)。

#### Node.js

我们通过 [npm](https://www.npmjs.com/package/ibmmq) 和 [GitHub 代码库](https://github.com/ibm-messaging/mq-mqi-nodejs) 提供了适用于 MQI API 的 Node.js 包装器库，该库可与适用于您操作系统的 MQ 客户端库结合使用。获取适用于您操作系统的客户端或可再分发客户端：

包WindowsLinuxMacRedist (grab & go) MQ 下载[IBM-MQC-Redist-Win64.zip](https://ibm.biz/915-IBM-MQC-Redist-Win64zip)[IBM-MQC-Redist-LinuxX64.tar.gz](https://ibm.biz/915-IBM-MQC-Redist-LinuxX64targz)[MacOS Toolkit](https://ibm.biz/mq-mac-toolkit)Redist (Fix Central)[IBM-MQC-Redist-Win64](https://ibm.biz/915-IBM-MQC-Redist-Win64)[IBM-MQC-Redist-LinuxX64](https://ibm.biz/915-IBM-MQC-Redist-LinuxX64)[MacOS Toolkit](https://ibm.biz/mq-mac-toolkit)Client (Fix Central)[IBM-MQC-Win64](https://ibm.biz/915-IBM-MQC-Win64)[IBM-MQC-LinuxX64](https://ibm.biz/915-IBM-MQC-LinuxX64) 或 [IBM-MQC-UbuntuLinuxX64](https://ibm.biz/915-IBM-MQC-UbuntuLinuxX64)[MacOS Toolkit](https://ibm.biz/mq-mac-toolkit)

要使用 MQ Light API 来开发 Node.js 应用程序，参阅我们的 [GitHub 代码库](https://github.com/mqlight/nodejs-mqlight) 和 [npm 包](https://www.npmjs.com/package/mqlight)。

要开始开发 Node.js 应用程序，参阅我们的 GitHub 代码库中的 [IBM MQ Node.js 样本](https://github.com/ibm-messaging/mq-dev-patterns/tree/master/Node.js)。

#### Golang

IBM 提供了适用于 MQI API 的开源 [Golang 包装器库](https://github.com/ibm-messaging/mq-golang)，该库可与适用于您操作系统的 MQ 客户端库结合使用。 另参阅适用于 Go 程序的开源 [类似于 JMS 的包装器](https://github.com/ibm-messaging/mq-golang-jms20)。

包WindowsLinuxMacRedist (grab & go) MQ 下载[IBM-MQC-Redist-Win64.zip](https://ibm.biz/915-IBM-MQC-Redist-Win64zip)[IBM-MQC-Redist-LinuxX64.tar.gz](https://ibm.biz/915-IBM-MQC-Redist-LinuxX64targz)[MacOS Toolkit](https://ibm.biz/mq-mac-toolkit)Redist (Fix Central)[IBM-MQC-Redist-Win64](https://ibm.biz/915-IBM-MQC-Redist-Win64)[IBM-MQC-Redist-LinuxX64](https://ibm.biz/915-IBM-MQC-Redist-LinuxX64)[MacOS Toolkit](https://ibm.biz/mq-mac-toolkit)Client (Fix Central)[IBM-MQC-Win64](https://ibm.biz/915-IBM-MQC-Win64)[IBM-MQC-LinuxX64](https://ibm.biz/915-IBM-MQC-LinuxX64) 或 [IBM-MQC-UbuntuLinuxX64](https://ibm.biz/915-IBM-MQC-UbuntuLinuxX64)[MacOS Toolkit](https://ibm.biz/mq-mac-toolkit)

要开始开发 Golang 应用程序，参阅我们的 GitHub 代码库中的 [IBM MQ Golang 样本](https://github.com/ibm-messaging/mq-dev-patterns/tree/master/Go)。

#### Python

您可以将 [pymqi 开源包装器库](https://dsuch.github.io/pymqi/#installation) 与适用于您操作系统的 MQ 客户端库结合使用。

包WindowsLinuxMacRedist (grab & go) MQ 下载[IBM-MQC-Redist-Win64.zip](https://ibm.biz/915-IBM-MQC-Redist-Win64zip)[IBM-MQC-Redist-LinuxX64.tar.gz](https://ibm.biz/915-IBM-MQC-Redist-LinuxX64targz)[MacOS Toolkit](https://ibm.biz/mq-mac-toolkit)Redist (Fix Central)[IBM-MQC-Redist-Win64](https://ibm.biz/915-IBM-MQC-Redist-Win64)[IBM-MQC-Redist-LinuxX64](https://ibm.biz/915-IBM-MQC-Redist-LinuxX64)[MacOS Toolkit](https://ibm.biz/mq-mac-toolkit)Client (Fix Central)[IBM-MQC-Win64](https://ibm.biz/915-IBM-MQC-Win64)[IBM-MQC-LinuxX64](https://ibm.biz/915-IBM-MQC-LinuxX64) 或 [IBM-MQC-UbuntuLinuxX64](https://ibm.biz/915-IBM-MQC-UbuntuLinuxX64)[MacOS Toolkit](https://ibm.biz/mq-mac-toolkit)

[pymqi](https://dsuch.github.io/pymqi/) 是一个受欢迎的第三方库。除了 [pymqi 提供的示例](https://dsuch.github.io/pymqi/examples.html) 外，另参阅 [此 GitHub 代码库中的 Python 样本](https://github.com/ibm-messaging/mq-dev-patterns/tree/master/Python)。

#### Ruby

要开发 Ruby 应用程序，参阅我们的 [MQ Light GitHub 代码库](https://github.com/mqlight/ruby-mqlight) 和 [MQ Light Ruby Gem](https://rubygems.org/gems/mqlight)。

#### C

C 涵盖了最全面的 MQI 内容，对于要构建高度优化的 IBM MQ 应用程序的开发者而言非常受欢迎。

包WindowsLinuxMacClient (Fix Central)[IBM-MQC-Win64](https://ibm.biz/915-IBM-MQC-Win64)[IBM-MQC-LinuxX64](https://ibm.biz/915-IBM-MQC-LinuxX64) 或 [IBM-MQC-UbuntuLinuxX64](https://ibm.biz/915-IBM-MQC-UbuntuLinuxX64)[MacOS Toolkit](https://ibm.biz/mq-mac-toolkit)Redist (Fix Central)[IBM-MQC-Redist-Win64](https://ibm.biz/915-IBM-MQC-Redist-Win64)[IBM-MQC-Redist-LinuxX64](https://ibm.biz/915-IBM-MQC-Redist-LinuxX64)[MacOS Toolkit](https://ibm.biz/mq-mac-toolkit)Redist (grab & go) MQ 下载[IBM-MQC-Redist-Win64.zip](https://ibm.biz/915-IBM-MQC-Redist-Win64zip)[IBM-MQC-Redist-LinuxX64.tar.gz](https://ibm.biz/915-IBM-MQC-Redist-LinuxX64targz)[MacOS Toolkit](https://ibm.biz/mq-mac-toolkit)

#### C++

C++ 开发者可以直接与 MQI 进行交互，还可以使用 XMS，后者提供了类似于 JMS 的消息传递抽象。

也可以将 MQ Object Oriented 库与 C++ 结合使用，但是这些库 [在 IBM MQ V7 中已经很稳定](https://www-01.ibm.com/common/ssi/cgi-bin/ssialias?subtype=ca&infotype=an&appname=iSource&supplier=897&letternum=ENUS209-245)（参阅“规划信息”下的内容）。

包WindowsLinuxMacClient (Fix Central)[IBM-MQC-Win64](https://ibm.biz/915-IBM-MQC-Win64)[IBM-MQC-LinuxX64](https://ibm.biz/915-IBM-MQC-LinuxX64) 或 [IBM-MQC-UbuntuLinuxX64](https://ibm.biz/915-IBM-MQC-UbuntuLinuxX64)[MacOS Toolkit](https://ibm.biz/mq-mac-toolkit)Redist (Fix Central)[IBM-MQC-Redist-Win64](https://ibm.biz/915-IBM-MQC-Redist-Win64)[IBM-MQC-Redist-LinuxX64](https://ibm.biz/915-IBM-MQC-Redist-LinuxX64)[MacOS Toolkit](https://ibm.biz/mq-mac-toolkit)Redist (grab & go) MQ 下载[IBM-MQC-Redist-Win64.zip](https://ibm.biz/915-IBM-MQC-Redist-Win64zip)[IBM-MQC-Redist-LinuxX64.tar.gz](https://ibm.biz/915-IBM-MQC-Redist-LinuxX64targz)[MacOS Toolkit](https://ibm.biz/mq-mac-toolkit)

#### COBOL

在许多平台上，您可以用 COBOL 编写应用程序来与运行在 IBM z/OS 大型机系统和 IBMi 系统上的队列管理器进行交互。

包WindowsLinuxMacClient (Fix Central)[IBM-MQC-Win64](https://ibm.biz/915-IBM-MQC-Win64)[IBM-MQC-LinuxX64](https://ibm.biz/915-IBM-MQC-LinuxX64) 或 [IBM-MQC-UbuntuLinuxX64](https://ibm.biz/915-IBM-MQC-UbuntuLinuxX64)[MacOS Toolkit](https://ibm.biz/mq-mac-toolkit)Redist (Fix Central)[IBM-MQC-Redist-Win64](https://ibm.biz/915-IBM-MQC-Redist-Win64)[IBM-MQC-Redist-LinuxX64](https://ibm.biz/915-IBM-MQC-Redist-LinuxX64)[MacOS Toolkit](https://ibm.biz/mq-mac-toolkit)Redist (grab & go) MQ 下载[IBM-MQC-Redist-Win64.zip](https://ibm.biz/915-IBM-MQC-Redist-Win64zip)[IBM-MQC-Redist-LinuxX64.tar.gz](https://ibm.biz/915-IBM-MQC-Redist-LinuxX64targz)[MacOS Toolkit](https://ibm.biz/mq-mac-toolkit)

#### 最新客户端

此页面上链接的所有客户端都是最新的持续交付 (CD) 版本客户端。我们还提供了长期服务 (LTS) 版本客户端。

如果需要使用其他版本进行开发，从以下 Fix Central 链接中选择合适的客户端：

- [最新的 CD 9.x.x 版本客户端 – 与此页面上的其他链接的版本相同](https://ibm.biz/ibm-mq-clients-all-cd)
- [最新的 LTS 9.1 版本客户端](https://ibm.biz/ibm-mq-clients-all-lts_91)
- [最新的 LTS 9.0 版本客户端](https://ibm.biz/ibm-mq-clients-all-lts_90)

要了解有关 IBM MQ 版本类型的更多信息，参阅 [IBM MQ Knowledge Center](https://www.ibm.com/support/knowledgecenter/SSFKSJ_latest/com.ibm.mq.pla.doc/q004715_.htm) 和 [IBM 支持中针对 LTS 和 CD 的 FAQ 页面](https://www.ibm.com/support/pages/ibm-mq-faq-long-term-support-and-continuous-delivery-releases)。

本文翻译自： [IBM MQ Downloads for developers](https://developer.ibm.com/articles/mq-downloads/)（2020-06-23）