# 深入剖析 z/VM 虚拟机管理程序
管理程序，虚拟化和云

**标签:** 云计算

[原文链接](https://developer.ibm.com/zh/articles/cl-hypervisorcompare-zvm/)

Bhanu P Tholeti

发布: 2011-12-27

* * *

##### 关于本系列

本系列文章首先介绍虚拟机管理程序类型和系统虚拟化的背景知识，然后介绍 5 个虚拟机管理程序的功能、它们的部署流程，以及您可能遇到的管理问题。

- [简介](http://www.ibm.com/developerworks/cn/cloud/library/cl-hypervisorcompare/)
- [PowerVM](http://www.ibm.com//developerworks/cn/cloud/library/cl-hypervisorcompare-powervm/)
- [VMware ESX Server](http://www.ibm.com//developerworks/cn/cloud/library/cl-hypervisorcompare-vmwareesx/)
- [Xen](http://www.ibm.com//developerworks/cn/cloud/library/cl-hypervisorcompare-xen/)
- [KVM](http://www.ibm.com//developerworks/cn/cloud/library/cl-hypervisorcompare-kvm/)
- z/VM

可以使用本系列文章作为理解虚拟机管理程序在云中的虚拟化过程中所扮演的角色的一个简单起点，或者您可以参阅本系列的各篇文章，帮助您确定哪个虚拟机管理程序最能满足您的云计算需求。

## 预备知识

z/VM® 虚拟机管理程序旨在帮助将大型机技术的价值扩大到整个企业，集成应用程序和数据，同时提供极高水平的可用性、安全性和操作简单性。

z/VM 虚拟化技术设计用于在单个运行其他 System z® 操作系统（比如 z/OS® ）的单一大型机上运行成百上千个 Linux® 服务器，或者用作大规模、仅限 Linux 的企业服务器解决方案。

z/VM V6.1 和 z/VM V5.4 也可帮助提高生产力，在相同的 System z 服务器上承载非 Linux 工作负载，比如 z/OS、z/VSE® 和 z/TPF，或者用作大规模的企业服务器解决方案。

z/VM 支持 Linux、z/OS、z/OS.e、Transaction Processing Facility (TPF) 和 z/VSE。z/VM 还支持将 z/VM 作为 Guest 操作系统。

## 功能

z/VM 基础产品包括以下组件和工具：

- **控制程序 (CP)：** CP 是一个虚拟机管理程序和真实机器资源管理器。
- **Conversational Monitor System (CMS)** CMS 提供了一个高容量的应用程序和交互式用户环境，还提供了 z/VM 文件系统。
- **TCP/IP for z/VM：** TCP/IP for z/VM 提供了对 TCP/IP 网络环境的支持。
- **高级程序间通信/虚拟机 (APPC/VM) 虚拟远程通信访问方法 (VTAM) 支持 (AVS)：** AVS 在 IBM Systems Network Architecture (SNA) 网络中提供了连接性。
- **Dump Viewing Facility：** 它是以交互方式诊断 z/VM 系统问题的一款工具。
- **Group Control System (GCS)：** GCS 是一个虚拟机管理程序，提供了多任务服务，支持本机 SNA 网络。
- **z/VM 硬件配置定义 (HCD) 和硬件配置管理器 (HCM)：** HCD 和 HCM 提供了一个完善的 I/O 配置管理环境。
- **z/VM 语言环境：** 语言环境为使用 C/C++、COBOL 或 PL/I 编写的 z/VM 应用程序提供了运行时环境。
- **开放系统适配器/支持工具 (OSA/SF)** OSA/SF 是一个自定义 OSA 硬件功能操作模式的工具。
- **REXX/VM：** REXX/VM 包含处理 REXX 编程语言的解释器。
- **透明服务访问工具 (TSAF)：** TSAF 在一组 z/VM 系统中提供了通信服务，而没有使用 VTAM。
- **Virtual Machine Serviceability Enhancements Staged/Extended (VMSES/E)：** VMSES/E 提供了一个安装和维护 z/VM 和其他支持的产品的工具套件。

z/VM 还提供了以下可选功能：

- **Data Facility Storage Management Subsystem for VM (DFSMS/VM)：** DFSMS/VM 控制数据和存储资源。
- **Directory Maintenance Facility for z/VM (DirMaint)：** DirMaint 提供了管理 z/VM 用户目录的交互式工具。
- **Performance Toolkit for VM：** Performance Toolkit 提供了分析 z/VM 和 Linux 性能数据的工具。
- **Resource Access Control Facility (RACF) Security Server for z/VM：** RACF 通过控制对安装程序的访问，为安装程序提供了数据安全保护。
- **Remote Spooling Communications Subsystem (RSCS) Networking for z/VM：** RSCS 支持用户向网络中的其他用户发送消息、命令、文件和作业。

## 部署虚拟化

要部署 z/VM 虚拟化：

- 创建逻辑分区。
- 在一个或多个逻辑分区中安装和配置 z/VM。
- 创建虚拟机。
- 安装和配置 Guest 操作系统。
- 为虚拟系统配置虚拟网络。

## 管理您的虚拟机

z/VM 通过 IBM Systems Director 管理虚拟机，将 IBM Systems Director 用作平台管理基础，以便支持与 Tivoli 的集成和第三方管理平台。使用 IBM Systems Director，您可以：

- 自动化数据中心操作。
- 统一 IBM 服务器、存储和网络设备的管理。
- 简化物理和虚拟平台资源的管理。
- 降低操作复杂性，提供 IT 系统的关系和健康状态视图。

您甚至可以获取整个数据中心的实际能源使用的单一视图。

## 选择 z/VM

优点：

- 能够将每个 LPAR 虚拟化为数百个或更多虚拟机。
- 能够虚拟化处理器、内存、I/O 和网络资源。
- 动态配置处理器、内存、I/O 和网络资源。
- 最大限度地提高资源利用，实现很高的系统利用率和先进的动态资源分配。
- 体验先进的系统管理和核算工具。

缺点：

- 您可能需要高度熟练、经过大型机知识培训的 IT 专业人员来维护它。

本文翻译自： [Dive into the z/VM hypervisor](https://developer.ibm.com/articles/cl-hypervisorcompare-zvm/)（2011-09-24）