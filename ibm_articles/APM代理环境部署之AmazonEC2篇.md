# APM 代理环境部署之 Amazon EC2 篇
全方位管理传统 IT、虚拟化、云、和混合环境

**标签:** DevOps

[原文链接](https://developer.ibm.com/zh/articles/d-proxy-environment-deployment-apm/)

牛凤枝, 陈佳蕴, 龙丽娟

发布: 2018-03-07

* * *

面对企业中越来越复杂且相互依赖的 IT 环境和应用。IBM Cloud Application Performance Management(APM)不仅是一个监控工具，更是可视化，可控性和自动化的解决方案。它是您身边的职能管家，全方位管理传统 IT、虚拟化、云、和混合环境。可以精确地发现并解决应用和 IT 基础架构中的性能问题。

## APM 架构介绍

APM 使用代理和数据收集器来收集监控主机的数据。代理和数据收集器将数据传送给 APM Server，并有 APM Server 将数据整合进 APM 控制台。请参照以下架构图：

##### 图 1\. APM 架构图

![APM 架构图](../ibm_articles_img/d-proxy-environment-deployment-apm_images_image001.png)

### 数据收集

代理和数据收集器监控系统、子系统或者应用并且收集数据。一个代理或数据收集器和单一资源（一个系统或者应用）进行交互，大多情况下是安装在系统或资源所在计算机或虚拟机上。例如，Linux OS 代理收集 Linux 主机上操作系统的性能参数。而 WebSphere Applications 代理监控 WebSphere applications servers 性能参数。

你可以给关键性能参数(KPIs)设置阀值。如果某个指数总是在阀值上下波动，就会产生报警。APM Server 做过适当的配置后，就可以把这些事件传送给 Netcool/OMNIbus 探测器给 Tivoli EIF 或者 SNMP Server. 利用警报通知配置 email 及时通知管理员。

代理和数据收集器被预先配置以便和 APM Server 通讯。

### 代理或数据收集器和服务器之间的通讯

代理和数据收集器会和 APM Server 建立 HTTPS 通讯。代理就是这个连接的客户端。

代理或数据收集器每 60 秒会发送一次数据给 APM Server. 例如即使此刻没有用户在查看数据，操作系统代理也会每 60 秒会发送例如 CPU 和内存利用数据给 APM Server。

### 数据存储

代理或数据收集器根据数据类型的不同每一分钟到八分钟不等推送给 APM Server, 事物监控数据也会一并发过去。APM Server 会保留这些数据 8 天。但是经过裁剪的数据会保存的久一些。被保存的数据称为历史数据。APM Server 把历史数据显示成图或表，从而帮助你分析环境的变化趋势。

### 可扩展性

依 APM 环境而定，最多连接 4000 个代理。

### 集成

APM 可以通过配置来集成其他的产品或组件以便与它们可以和 APM Server 通讯。

可以用于集成的产品或组件包括：Netcool/OMNIbus, Tivoli Monitoring, Operations Analytics–Log Analysis, Operations Analytics–Predictive Insights, and IBM Alert Notification.

Agent Builder 是一个可以用来创建定制化代理的组件。

### APM Server 安装

由于本篇着重介绍 APM 代理的安装配置。所以关于 APM Server 的安装请参阅 [APM Server 安装文档](https://www.ibm.com/support/knowledgecenter/SSHLNR_8.1.4/com.ibm.pm.doc/install/install_server_intro.htm) 。

## Amazon EC2 代理简介

亚马逊弹性计算云 Amazon Elastic Compute Cloud(EC2)是一个让使用者可以租用云端电脑运行所需应用的系统。EC2 借由提供 Web 服务的方式让使用者可以弹性地运行自己的 Amazon 机器映像档，使用者将可以在这个虚拟机器上运行任何自己想要的软件或应用程序。提供可调整的云计算能力。它旨在使开发者的网络规模计算变得更加容易。

那么实时了解 Amazon EC2 实例运行状态就变得非常重要。

Amazon EC2 代理将为您提供一个中心点以便您检测 Amazon EC2 实例的运行状况、可用性和性能。该代理将显示一套全面的指标以帮助您对 EC2 环境迅速作出决策。这些指标包括：CPU 使用率、弹性块存储(ESB)使用、网络的使用，亚马逊 Web 服务(AWS)维护更新和磁盘性能等。

APM 有数以十计的监控代理，Amazon EC2 代理是其中之一。本篇仅以 Amazon EC2 代理为例介绍其安装和配置。

## Amazon EC2 测试环境部署

在部署测试环境前，我们先熟悉以下几个 AWS 基本概念：

- AWS 区域和可用区
- AMI 和 EC2 实例
- AWS 安全凭证

### AWS 区域和可用区

##### 图 2\. AWS 区域和可用区

![AWS 区域和可用区](../ibm_articles_img/d-proxy-environment-deployment-apm_images_image002.png)

我们要部署一台机器在云上。首先需要选定区域和可用区。AWS 在全球有 11 个区域（这个数据可能后续有所变化），区域指全球大的地理都市圈里若干数据中心的集合。区域内部资源是相互连通的，区域与区域在物理上是相互孤立的。一个区域内部有若干可用区, 可用区之间是光纤连接。可用区之间传输带宽很大延迟很低，所以如果您有好几台机器分别部署在不同的可用区，那么您可以认为它们就在同一个机房，肩并肩靠在一起。可用区内部又有若干的数据中心。AWS 是通过区域和可用区双重保障为您的应用程序提供高可用性。

### AMI 和 EC2 实例

按照上述选择完区域和可用区后，接下来需要选择一个 AMI. AMI 决定了机器是什么样子的，是 Windows 还是 Linux 等。

如果按照以前的传统方式，当拿到一台服务器，我们需要装系统、打补丁和配置服务器等。这些无差别的重复劳动。AWS 提供了 AMI (Amazon Machine Image)帮您记录系统镜像。当然也可以自己做一个 AMI, 甚至把网站也配好。通过 AMI 可以做到预配置和动态加载之间的灵活平衡。准备好 AMI 就可以进行大规模的复制。

接下来选择 EC2 实例。实例决定了这台机器的性能。

##### 图 3\. EC2 实例

![EC2 实例](../ibm_articles_img/d-proxy-environment-deployment-apm_images_image003.png)

您可以从单个 AMI 启动不同的类型的实例。实例类型实际上决定了主机的硬件。每个实例类型提供不同的计算和内存功能。您可以根据实例上运行的应用程序或软件所需的内存和计算能力，来选择一个实例类型。

### AWS 安全凭证

在选择完 EC2 实例后，需要对实例进行网络、存储、安全组和防火墙进行配置。这些就绪后，管理员需要密钥对管理员才能访问主机。

还有一种叫访问键，它由访问密钥 ID (例如, AKIAJ34YEBHY3AS5LBIQ) 和一个秘密访问密钥 (例如, 3AD3So3anwUCM9gdVSkMXUhJMjg4PTxTCPiBbTPI) 组成。如果使用 AWS SDK、REST 或查询 API 操作, 则可以使用访问键对您对 AWS 所做的编程请求进行签名。AWS SDK 使用您的访问密钥为您签名请求, 这样您就不必处理签名过程。您还可以手动对请求进行签名。访问密钥 ID 和秘密访问密钥在我们后续对 Amazon EC2 代理配置时需要用到。

## Amazon EC2 测试环境搭建

在此我们简单部署一个 Amazon EC2 测试环境，以便可以清晰的通过 Amazon EC2 代理监控到数据。测试环境的搭建需要先具备以下条件:

- 申请 AWS 账户，请参照：AWS [创建新账户](https://portal.aws.amazon.com/billing/signup#/start) ，按照步骤申请即可。
- 包含有被监控 EC2 实例的 AWS 区域名称
- 访问 AWS 区域的安全凭证（访问密钥 ID 和秘密访问密钥）

测试环境搭建步骤如下:

- 登陆 Amazon 管理控制台
- 从 AMI 创建 Amazon Linux 实例
- 在 Amazon 管理控制台找到您的实例
- 登陆您的实例

### 登陆 Amazon 管理控制台

不再赘述，请参照：AWS [创建新账户](https://portal.aws.amazon.com/billing/signup#/start) ，按照步骤创建完账户后登陆即可。

### 从 AMI 创建 Amazon Linux 实例

登陆到 AWS 后，右上角可以选择区域。

##### 图 4\. 选择区域

![选择区域](../ibm_articles_img/d-proxy-environment-deployment-apm_images_image004.png)

选择完区域后，AWS 会自动分配可用区。单击 AWS 控制台首页 Amazon EC2 图标。选择启动实例：

##### 图 5\. 启动实例

![启动实例](../ibm_articles_img/d-proxy-environment-deployment-apm_images_image005.png)

选择 Amazon Linux AMI:

##### 图 6\. 选择 AMI 类型

![选择 AMI 类型](../ibm_articles_img/d-proxy-environment-deployment-apm_images_image006.png)

选择实例类型：

在此我们选择免费的实例类型供测试使用

##### 图 7\. 选择实例类型

![选择实例类型](../ibm_articles_img/d-proxy-environment-deployment-apm_images_image007.png)

简单起见，之后不做更复杂的配置了。点击”Review and Launch”。

下面建立新的密钥对，它被用来连接您的实例

##### 图 8\. 下载密钥对

![下载密钥对](../ibm_articles_img/d-proxy-environment-deployment-apm_images_image008.png)

完成以上步骤后，就可以启动实例了。

##### 图 9\. 查看实例

![查看实例](../ibm_articles_img/d-proxy-environment-deployment-apm_images_image009.png)

关于如何使用密钥对访问实例，不在本文讨论范围内，请参照 [AWS 相关文档](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AccessingInstances.html?icmpid=docs_ec2_console) 。

## Amazon EC2 代理环境部署

### 系统需求

请参照 [系统需求](https://www.ibm.com/support/knowledgecenter/SSHLNR_8.1.4/com.ibm.pm.doc/install/planning_requirements.htm) 来安排您的资源。Amazon EC2 代理支持远程监控。在此我们选择 Windows Sever 2008 R2 Enterprise 64bit 来安装我们的监控代理监控刚刚创建的 Amazon EC2 服务。

### 前提条件

在部署 Amazon EC2 代理环境前，我们需要先部署 Amazon EC2 测试环境(上述文章已经完成这部分的搭建)。

### 代理安装包的预配置

这部分在 APM Server 安装过程已经完成。我们只需要从 Server 上把安装包下载下来即可。也可以通过以下步骤重新配置：

登陆 APM Server 主机

在/opt/ibm/ccm（默认安装目录）下运行 “make\_configuration\_packages.sh”

##### 图 10\. 创建配置包文件

![创建配置包文件](../ibm_articles_img/d-proxy-environment-deployment-apm_images_image010.png)

然后在同一目录下运行 “configure\_agent\_images.sh”

##### 图 11\. 代理安装包预配置

![代理安装包预配置](../ibm_articles_img/d-proxy-environment-deployment-apm_images_image011.png)

### 代理的安装

经过上述步骤，经过预配置的安装包在/opt/ibm/ccm/depot 目录下。我们需要把这个安装包下载到 Windows 机器上并解压。

安装部分很简单，双击 “installAPMAgents.bat” 一步步往下装即可，此处不再赘述。

### 代理的配置

此处重点讲一下代理的配置。代理安装完成后，点击”开始->所有程序->IBM Monitoring agents”启动 Tivoli Enterprise Monitoring Services 管理面板。

##### 图 12\. TEMS 管理面板

![TEMS 管理面板](../ibm_articles_img/d-proxy-environment-deployment-apm_images_image012.png)

右键点击”Monitoring Agent for Amazon EC2″

##### 图 13\. Amazon EC2 配置步骤-1

![Amazon EC2 配置步骤-1](../ibm_articles_img/d-proxy-environment-deployment-apm_images_image013.png)

选择”使用预设值配置”并输入实例名称

##### 图 14\. Amazon EC2 配置步骤-2

![Amazon EC2 配置步骤-2](../ibm_articles_img/d-proxy-environment-deployment-apm_images_image014.png)

##### 图 15\. Amazon EC2 配置步骤-3

![Amazon EC2 配置步骤-3](../ibm_articles_img/d-proxy-environment-deployment-apm_images_image015.png)

此处需输入访问密钥 ID 和一个秘密访问密钥及区域名称（这些信息均在 Amazon EC2 测试环境部署章节获取）点击”OK”保存配置并启动代理。

接下来我们需要登录 APMUI 来验证代理是否和 APM Server 连接成功并呈现监控数据。登录 APMUI 方式：https://:9443 默认的用户名密码是 apmadmin/apmpass。

##### 图 16\. Amazon EC2 实例监控数据

![Amazon EC2 实例监控数据](../ibm_articles_img/d-proxy-environment-deployment-apm_images_image016.png)

可以在此处做个对照，这和我们上述创建的实例一致，在这里是条目二。

如果想了解这个实例的详细数据，可以点击实例条目查看。

##### 图 17\. 查看实例详情

![查看实例详情](../ibm_articles_img/d-proxy-environment-deployment-apm_images_image017.png)

从这里我们可以查看例如：实例的描述信息、CPU 利用率、平均网络利用率等监控数据。可以登录 AWS 控制台查看数据以做比较。

##### 图 18\. AWS 控制台实例数据

![AWS 控制台实例数据](../ibm_articles_img/d-proxy-environment-deployment-apm_images_image018.png)

## 结束语

至此，我们详细介绍了 Amazon EC2 及 APM Amazon EC2 代理的安装和配置。相信大家已经对 APM 监控架构有了初步的认识。

APM 有数以十计的代理帮助我们监控企业 IT 环境和应用的运行状态。后续我们也将推出系列文章讲述一些大家关心和时下比较流行的服务监控代理的安装、配置。希望这些文章对对 APM 感兴趣的开发和测试人员有所帮助。

## 参考资源 (resources)

- [AWS 官网](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/concepts.html): 了解并学习 Amazon EC2
- [User Guide](https://www.ibm.com/support/knowledgecenter/SSHLNR_8.1.4/com.ibm.pm.doc/install/overview_product.htm): 学习 APM 官方文档