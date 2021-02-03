# 使用 IoT 平台简化 IoT 应用程序的开发
帮助您理解为什么应该使用 IoT 平台的指南

**标签:** IoT

[原文链接](https://developer.ibm.com/zh/articles/iot-lp101-why-use-iot-platform/)

Anna Gerber

更新: 2020-01-31 \| 发布: 2017-08-02

* * *

##### 学习路径: 开始使用物联网开发

本文是 IoT 101 学习路径的一部分，是物联网开发人员的快速入门指南。

- [物联网概念和技能](/zh/articles/iot-key-concepts-skills-get-started-iot)
- [物联网硬件指南](/zh/articles/iot-lp101-best-hardware-devices-iot-project)
- [物联网网络指南](/zh/articles/iot-lp101-connectivity-network-protocols/)
-  物联网平台（本文）
- [教程：构建一个简单的家庭自动化系统](/zh/tutorials/iot-lp101-get-started-develop-iot-home-automation/)

开发和部署 IoT 系统很复杂 — 涉及到协调各种各样的连接设备、网络、云服务，以及移动和 Web 应用程序。

IoT 平台提供了中间件来连接和管理硬件设备，以及它们使用面向用户的移动和 Web 应用程序收集的数据。这些 IoT 平台通常设计为能在设备、应用程序和服务之间建立安全通信。它们还使得开发和部署 IoT 应用程序变得更快、更容易。IoT 平台被设计为可扩展和可靠的平台。例如，由于 IoT 系统实现了容错，所以它们在发生设备、网络或云服务故障时会继续运行。通过利用标准和提供具有详细记录的 API 和 SDK，IoT 平台可帮助您将完全不同的设备、应用程序和服务集成到 IoT 系统中，让您能根据需求变化来灵活地自定义和替换组件。

物联网正在快速演化，所以您需要快速且高效地开发 IoT 应用程序，以便最大限度缩短上市准备时间并保持竞争力。尽管每个 IoT 计划都具有自己的一组独特需求，但是对于连接设备的管理，以及它们生成的数据的传输、存储和分析，大部分 IoT 项目都有一些共同需求。IoT 平台解决了这些基本需求，它们不仅提供了通用的、可自定义的工具和服务，还提供了 [抽象](https://developer.ibm.com/iotplatform/2017/01/17/iot-platform-vs-iot-infrastructure/)，您可以在这些抽象的基础上开发自定义 IoT 解决方案。

通过探索 IoT 平台的重要功能，分析最顶级的多个用途的端到端 IoT 平台，您会意识到为什么应该采用 IoT 平台来快速开发 IoT 解决方案，并快速从 IoT 数据中获取最大价值。

## IoT 平台功能

考虑以下重要的 IoT 平台功能：

- 设备管理
- 数据通信协议
- 数据存储
- 规则和分析
- 快速应用程序开发和部署
- 集成
- 安全性
- 开发和维护使用 IoT 平台构建的解决方案的成本。

每项功能的相对重要性取决于您的具体用例。例如，考虑一个小型家庭自动化系统。设备管理不是主要问题，因为要管理的设备很少。由于家庭自动化用例是一个非常被动的用例，所以不太需要 IoT 平台提供数据存储和数据分析功能来存储和分析历史传感器读数。但是，IoT 平台必须支持快速可靠的数据通信协议，以确保没有错过传感器检测到的事件，并能得到近实时的响应。但是在工业 IoT 应用程序中，涉及到数千个传感器，它们监视着工厂设备，以检测设备故障和预测何时对设备执行维护，这时 IoT 平台的设备管理和分析功能显然重要得多。

### 设备管理

要了解 IoT 设备管理的详细信息，请阅读我的文章“ [为您的下一个 IoT 项目选择最佳硬件](https://www.ibm.com/developerworks/cn/iot/library/iot-lp101-best-hardware-devices-iot-project/index.html)”。

IoT 平台的一项核心功能就是管理 IoT 硬件设备。设备管理会在每个 IoT 设备的整个生命周期中发生，包括配备新设备，监视和维护操作设备，以及最终在设备达到使用期限时让它们退役。IoT 平台提供了管理设备注册、配置和无线更新的功能。借助 IoT 平台，您可以收到设备状态警报，监视设备健康或使用情况的统计信息。IoT 平台还会在发生故障时提供远程调试。资产管理也属于设备管理的范畴，包括列出和搜索连接设备，查询和管理设备元数据的能力。

寻找对设备不可知并能与现有企业资源规划 (ERP) 系统集成的 IoT 平台。

### 数据通信协议

要了解 IoT 网络协议的详细信息，请参阅我的文章“ [连接物联网中的所有事物](https://www.ibm.com/developerworks/cn/iot/library/iot-lp101-connectivity-network-protocols/index.html)”。

IoT 设备、网关以及基于云的应用程序和服务之间的安全可靠的通信，对于实现远程管理、更新连接设备和传输数据是不可或缺的。尽管一些数据可存储在本地并通过边缘分析加以处理，但从传感器和其他 IoT 设备收集的一些数据需要传送到上游的云服务或应用程序，以便实现聚合、进一步处理或可视化。一个 IoT 系统中通常连接着数千个设备，所以许多 IoT 平台都包含消息代理服务，使设备和网关能够低延迟地大范围发送和接收消息。这些消息代理服务通常会使用标准通信协议，比如 MQTT、CoAP 或 XMPP。一些平台还支持使用 Web 套接字来实现实时通信。

### 数据存储

鉴于数十亿个连接设备正在以前所未有的速度生成数据，将海量数据存储在何处是选择 IoT 平台时要考虑的另一个方面。数据湖（Data lakes）可用来存储各种不同格式的数据，包括原始数据、结构化数据、半结构化数据或非结构化数据。这些数据湖通常构建于 NoSQL 或 Hadoop 技术之上，以实现最高的灵活性和可扩展性。数据湖和 IoT 平台提供的其他数据存储通常高度可用，支持并发读写，而且包含可配置的索引来提高 IoT 应用程序访问和查询数据时的性能。

### 规则和分析

以可靠方式存储数据后，接下来要考虑的是如何从 IoT 数据中获取洞察。许多 IoT 平台提供了包括实时处理和批处理在内的分析功能。规则引擎证明了实时数据分析生成可操作的洞察的能力，可以使用规则对用于触发操作的条件进行编码。一些平台提供了立即可用的分析应用程序，但是支持对内置报告、分析仪表板和数据可视化进行自定义的平台是更灵活的选择。

### 快速应用程序开发和部署

加快特定用途应用程序的开发和交付，是采用 IoT 平台的主要推动因素之一。因此，合理的做法是选择一个提供了代码样板生成工具、模板或小部件来支持快速开发，并提供了工具来简化和自动化部署的 IoT 平台。

### 集成

由于各种各样的连接设备和 IoT 相关服务已投入生产，更别提未来几年预计上线的数十亿个设备，所以互操作性是 IoT 中的一个重要考虑因素。鉴于 IoT 演化速度如此之快，让 IoT 系统顺应未来就显得至关重要。IoT 平台实现此互操作性的方法是采用一些标准，并在底层技术之上提供额外的抽象层和标准化层。例如，如果升级 IoT 设备来使用具有更高分辨率的传感器，在事物发生变化时，您的应用程序被破坏的可能性就更低。IoT 平台可以通过 SDK 和 API（其中 REST API 尤为常见）与其他平台、设备、Web 服务、工具和应用程序集成。

### 安全性

安全性是任何 IoT 解决方案设计的一个重要组成部分。安全性包括保护设备和网络通信，以及实施云和应用程序级别的安全保护。选择一个提供了身份验证和授权服务的 IoT 平台，而且该平台最好集成了用于身份管理（比如 LDAP）的现有系统，以及采用了设备身份验证标准（比如 X.509）的现有系统。要实现安全通信，IoT 平台应采用 TLS 等标准来加密内容，确保数据完整性。IoT 平台也可以对用户、设备、网关和服务实施基于角色的访问控制，并提供安全监控和审计工具。

### 成本

加速 IoT 解决方案的开发是采用 IoT 平台的一个主要优势，但 IoT 平台也可以帮助降低开发成本。与开发和维护定制的内部解决方案相比，持续成本（比如访问托管 IoT 平台的订阅费或自托管平台的托管成本）通常更少。许多 IoT 平台提供了免费试用期或较低的使用门槛，您可以利用这一点几乎没有前期开支地开始使用 IoT 平台，但是一定要考虑到这些成本会在初期试用结束后发生变化。

## 比较 IoT 平台

一些供应商提供了数百个 IoT 平台。一些平台非常专业化，例如仅迎合某个垂直市场，比如 GE Predix ( [https://www.predix.io/](https://www.predix.io/)) 就是一个专门针对工业 IoT 的 PaaS IoT 平台。其他 IoT 平台仅专注于提供您的 IoT 系统可能需要的功能子集，而不是端到端 IoT 平台解决方案，例如 Cisco and SAS Edge-to-Enterprise IoT Analytics Platform 仅专注于 IoT 分析。最终，您可能选择采用多个 IoT 平台。

以下是 5 个适合许多不同应用的流行的多用途的端到端 IoT 平台，可以使用它们作为起点：

- **[IBM Watson IoT](https://www.ibm.com/internet-of-things/solutions/iot-platform/watson-iot-platform/)** >

    IBM Watson IoT 构建于 IBM 的 Bluemix PaaS 平台之上，是一个成熟的、对开发人员友好的 IoT 平台，实时分析和认知计算是它的强项。
    该平台还支持设备管理、身份验证和授权、安全数据存储和通信、标准通信协议（比如 MQTT 和 HTTPS）、规则、REST API 和 SDK。根据市场研究公司 IDC 的研究， [Watson IoT 产品组合和 Watson IoT Platform 被评为“平台的 IoT 平台”](https://www.ibmjournal.com/hubfs/Internet_of_Things_content/The_Platform_of_platforms.pdf?t=1475517407639)。

- **[Amazon Web Services IoT](https://aws.amazon.com/iot-platform/how-it-works/)**

    AWS IoT 是一个高度可扩展的 IoT 平台。它包含一个 SDK、身份验证和授权、设备注册表、一个使用 MQTT、WebSocket 或 HTTP 与设备通信的设备网关，以及一个与 DynamoDB 等现有 AWS 服务集成的规则引擎。AWS IoT 的一个重要特性是“设备卷影 (device shadow)”，这是每个包含最后已知状态的设备的持久虚拟版本。

- **[Microsoft Azure IoT Suite](https://www.microsoft.com/en-au/internet-of-things/azure-iot-suite)**

    Microsoft Azure IoT Suite 是另一个综合 IoT 平台，它使用 Azure IoT Hub 支持设备管理和孪生设备（比如设备卷影），还支持标准通信协议（包括 MQTT、AMQP 和 HTTP）、安全存储、规则，以及分析引擎（支持预测分析和数据可视化）。

- **[ThingWorx](https://www.thingworx.com/)**

    PTC ThingWorx 是一个企业级 IoT 平台，它支持模型驱动的快速应用程序开发。特性包括设备管理、应用程序建模，支持标准协议（包括 MQTT、AMQP、XMPP、CoAP 和 WebSocket）、预测分析，以及 REST API 和 SDK 以提供集成支持。

- **[Kaa](https://www.kaaproject.org/)**

    Kaa 是一个免费的开源 IoT 平台，它是依据 Apache 2.0 许可来发布的，这使得该平台可以自行托管。Kaa 包含用于 Java、C++ 和 C 的 REST API 和 SDK，具有设备管理、数据收集、配置管理、通知、负载平衡和数据分析功能。


## 结束语

面对如此多可供选择的不同 IoT 平台，选择一个平台乍一看似乎有点困难。您可以考虑自己对每种平台的功能的需求，选择与您的业务最匹配的平台。采用 IoT 平台可以简化、加速和改善您的下一个 IoT 项目的开发。通过使用 IoT 平台，您可以专注于您的用例所独有的需求，快速向客户提供价值。

本文翻译自： [Streamlining the development of your IoT applications by using an IoT platform](https://developer.ibm.com/articles/iot-lp101-why-use-iot-platform/)（2020-01-31）