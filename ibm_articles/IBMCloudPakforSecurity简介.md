# IBM Cloud Pak for Security 简介
集成安全工具来深入了解混合多云环境中的威胁

**标签:** 云计算

[原文链接](https://developer.ibm.com/zh/articles/introduction-to-ibm-cloud-pak-for-security/)

Russ Warren

发布: 2020-07-29

* * *

[IBM Cloud Pak for Security](https://www.ibm.com/cn-zh/products/cloud-pak-for-security) 提供了一个平台，可帮助企业安全团队更快地集成现有的安全工具，以便其更深入地了解混合多云环境中的威胁。它将使用与基础架构无关且可在任何位置运行的通用操作环境。您可以快速搜索威胁、协调操作并自动进行响应，而无需移动数据。Cloud Pak 可满足以下关键业务需求：

- 简化并加快调查。利用联合搜索，您可以使用现有的安全工具来调查整个组织中的威胁和感染指标 (IOC)。根据您的威胁情报来源或 IBM 来源，发现并分析这些洞察。 [进一步了解 Data Explorer](https://www.ibm.com/cn-zh/products/cloud-pak-for-security/data-explorer)。
- 快速彻底响应威胁。编排和自动化可帮助您稳妥地应对网络安全事件。通过自动执行任务并对任务划分优先级以及在团队之间开展协作，找到并修复威胁。 [进一步了解 Resilient](https://www.ibm.com/cn-zh/products/cloud-pak-for-security/resilient)。
- 对相关威胁信息划分优先级并在环境中扫描以查找威胁。 [进一步了解 IBM Security Threat Intelligence Insights](https://www.ibm.com/cn-zh/products/cloud-pak-for-security/threat-intelligence-insights)。

## 如何与 IBM Cloud Pak for Security 集成

IBM Cloud Pak for Security 将连接到第三方工具和数据源，包括多个安全信息和事件管理 (SIEM) 工具、端点检测系统、威胁情报服务以及身份和云存储库。您还可以为环境中的任何工具或本地数据库构建定制连接器。IBM Cloud Pak for Security 提供了多个有用的选项，例如，联合安全数据，以及提供要与其他数据源关联的资产信息。

### 通过 Universal Data Insights 服务和 STIX-shifter 来联合安全数据

利用联合搜索，您可以调查和分析整个公司中的安全洞察，而无需移动数据。从共享数据源中查询或读取安全数据的任何应用程序都可以使用 Universal Data Insights 服务 API。它与称为 [STIX-shifter](https://developer.ibm.com/technologies/security/projects/stix-shifter/) 的可扩展开源软件开发包 (SDK) 集成，后者是通过 OASIS Open Cybersecurity Alliance 提供的。

### 提供要与其他数据源关联的资产信息

整合来自各种安全工具和 IT 工具的资产和风险信息，以便识别安全漏洞并更好地了解整体安全状况。Connect Assets and Risks (CAR) 连接器提供了从各种数据源（如 SIEM 工具、威胁检测产品和端点产品）提取并上传的资产信息。要了解公司的环境和风险状况，就需要整合资产信息。

## 结束语

Cloud Pak for Security 通过这种集成方法为客户提供了很多有价值的服务。通过提供一个安全平台来整合来自多个源的安全信息，您的客户可以：

- 进一步发掘现有安全工具的价值。
- 提高分析人员的工作效率，使其能够完成更多工作。
- 联合数据而不必移动数据，让数据留在原来的位置。无需额外的数据湖。
- 通过在一个屏幕中搜索各种不同的数据，更快地发现隐藏的威胁。
- 通过将数据复制到数据湖来降低隐私风险。
- 利用预先构建的集成，避免在内部构建成本高昂的产品集成。

了解如何根据教程 _[使用 STIX-shifter 交换网络威胁情报](https://developer.ibm.com/tutorials/exchange-cyber-threat-intelligence-with-stix-shifter/)_ 来开发新的 STIX-shifter 适配器。

本文翻译自： [Introduction to IBM Cloud Pak for Security](https://developer.ibm.com/articles/introduction-to-ibm-cloud-pak-for-security/)（2020-06-10）