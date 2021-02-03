# IBM Cloud Pak for Data 简介
了解 IBM Cloud Pak for Data 的基础知识

**标签:** IBM Cloud Pak for Data,分析,数据库,数据科学,数据管理,机器学习

[原文链接](https://developer.ibm.com/zh/articles/intro-to-cloud-pak-for-data/)

[Scott Dangelo](https://developer.ibm.com/zh/profiles/scott.dangelo), Clarinda Mascarenhas

发布: 2019-10-06

* * *

**本文已纳入 [IBM Cloud Pak for Data 快速入门学习路径](https://developer.ibm.com/zh/series/cloud-pak-for-data-learning-path/)**。

级别主题类型**100****[IBM Cloud Pak for Data 简介](https://developer.ibm.com/zh/articles/intro-to-cloud-pak-for-data)****文章**101[利用 Data Virtualization 实现 Db2 Warehouse 数据虚拟化](https://developer.ibm.com/zh/tutorials/virtualizing-db2-warehouse-data-with-data-virtualization)教程201[利用 Data Refinery 实现数据可视化](https://developer.ibm.com/zh/tutorials/data-visualization-with-data-refinery)教程202[使用 Watson Knowledge Catalog 查找、准备和理解数据](https://developer.ibm.com/zh/tutorials/find-prepare-and-understand-data-with-watson-knowledge-catalog)教程301A[借助内置 Notebook 的 Watson Machine Learning 进行数据分析、建模以及部署](https://developer.ibm.com/zh/patterns/data-analysis-model-building-and-deploying-with-wml)Pattern301B[使用 AutoAI 自动构建模型](https://developer.ibm.com/zh/tutorials/automate-model-building-with-autoai)教程301C[使用 IBM SPSS Modeler 快捷构建预测机器学习模型](https://developer.ibm.com/zh/tutorials/build-an-ai-model-visually-with-spss-modeler-flow)教程401[利用 Watson OpenScale 监视模型](https://developer.ibm.com/zh/patterns/watson-openscale-with-watson-machine-learning-engine-on-icp4d)Pattern

对于众多行业而言，AI 之旅都是一项长期战略，并且才刚刚起步。在此学习路径中，我们将研究电信公司案例。我们将查看数据收集流程，这些数据可能驻留在多个云中，采用各种数据库格式，并且具有不同的访问控制需求。在电信行业，我们将演示如何利用可视化工具和各种其他工具来组织数据。接下来，我们将查看客户流失率案例，并创建一个机器学习模型来帮助我们预测电信客户流失风险。最后，我们将分析电信业的机器学习模型部署情况，查看模型的性能、可解释性和公平性。

[IBM Cloud Pak for Data](https://www.ibm.com/products/cloud-pak-for-data) 是一个预先集成的统一数据和 AI 平台，在 [RedHat OpenShift Container 平台](https://www.openshift.com/products/container-platform) 上以本机方式运行。通过开放式可扩展云原生平台交付服务，以便收集、组织和分析数据。它采用单一界面借助内置监管功能执行端到端的分析。它还支持并管控端到端 AI 工作流程。

### 收集数据

- 无需迁移即可直接访问所有数据，从数据源提供安全保障。
- 连接至所有数据，消除数据孤岛。

### 组织数据

- 打造随时可用于业务的可信分析基础，可简化数据准备、策略、安全与合规流程。
- 管控并自动处理数据和 AI 生命周期。

### 分析数据

- 构建、部署和管理 AI 及机器学习功能，这些功能可在整个组织内以一致的方式进行扩展。

### 注入 AI

- 在整个业务中以透明且可信的方式有效运作 AI。

- 随时随地灵活运行，避免供应商锁定。


IBM Cloud Pak for Data 可提供规范方法来加速 AI 之旅：AI 云梯专为帮助客户实现业务数字化转型而开发，不论客户身处 AI 之旅的任何阶段，都能助一臂之力。IBM Cloud Pak 将所有关键云、数据和 AI 功能合而为一，构成容器化的微服务，进而在多云平台上交付 AI 云梯。

## 产品预演

IBM Cloud Pak for Data 可帮助您释放数据价值，并为 AI 创建信息架构。本产品预演提供分步演示，展示如何通过可扩展的 Kubernetes 平台来收集、组织、分析数据并注入 AI 功能。

## 架构

IBM Cloud Pak for Data 由在多节点 IBM Cloud Private 集群上运行的多种预配置的微服务组成。这些微服务支持您连接到自己的数据源，以便可以通过单一 Web 应用程序对数据进行编目、管控、浏览、剖析、变换和分析。

IBM Cloud Pak for Data 使用 RedHat OpenShift 部署在一个多节点 Kubernetes 集群上。虽然您可以将 IBM Cloud Pak for Data 部署在 3 节点集群上，但强烈建议您将自己的生产环境部署在一个至少具有 6 个节点的集群上，以便提高性能和集群稳定性，并便于扩展集群来支持工作负载增长需求。

## 结束语

本文介绍了 IBM Cloud Pak for Data、部分术语和概念、产品预演以及架构概述。本文纳入 [IBM Cloud Pak for Data 快速入门学习路径](https://developer.ibm.com/zh/series/cloud-pak-for-data-learning-path/)。要继续学习本系列并了解有关 IBM Cloud Pak for Data 的更多信息，可阅读下一个教程 [利用 Data Virtualization 实现 Db2 Warehouse 数据虚拟化](https://developer.ibm.com/zh/tutorials/virtualizing-db2-warehouse-data-with-data-virtualization)。

本文翻译自： [Introduction to IBM Cloud Pak for Data](https://developer.ibm.com/articles/intro-to-cloud-pak-for-data/)（2019-10-06）。