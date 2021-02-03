# OpenShift 漫游指南
在使用 Red Hat OpenShift on IBM Cloud 时，只需记住一点：不要害怕！

**标签:** Kubernetes,Red Hat OpenShift on IBM Cloud,容器,微服务

[原文链接](https://developer.ibm.com/zh/articles/hitchhikers-guide-openshift/)

Anton McConville, [Yan Koyfman](https://developer.ibm.com/zh/profiles/koyfman), [Olaph Wagoner](https://developer.ibm.com/zh/profiles/mwagone)

发布: 2019-11-27

* * *

本文档提供关于 OpenShift 的术语表、链接和主题指南。这是面向开发者的指南，由使用 Red Hat® OpenShift® on IBM Cloud™ 的开发者所编写。

OpenShift 被 [描述](https://github.com/openshift) 为“对开发者和运维人员友好的 Kubernetes 发行版”。它可以在 Kubernetes 上运行（OpenShift 的先前版本使用与此不同的机制处理容器编排）。OpenShift 提供各类工具，帮助开发者和运维人员运行容器化的工作负载。在幕后，OpenShift 由 [Origin Kubernetes Distribution (OKD)](https://www.okd.io/) 提供技术支持，这包括 Kubernetes 和其他开源项目，如 Docker 和 lstio。

以下各节按字母顺序排列。我们以 _《银河系漫游指南》_ 系列作者 Douglas Adams 的名义，向您建议：”不要害怕！”

![OpenShift 漫游指南](../ibm_articles_img/hitchhikers-guide-openshift_images_card.jpg)

## 应用程序

随着 OpenShift 的发展，“应用程序”逐渐成为一个 [被过度使用的词汇](https://docs.openshift.com/enterprise/3.0/whats_new/applications.html#whats-new-applications)。虽然之前是一个具体的概念，但在 OpenShift 领域，它不再表示单个基础对象。而在控制台和命令行中，“应用程序”作为一种对工作负载某些特性的便捷分组而存在。`oc new-app` CLI 可从现有镜像或源代码位置创建多个组件，比如 Deployment 和 ImageStream，如果指定了端口，还会创建一个服务和路由配置。

## CI/CD 与 Jenkins 或 Tekton

持续集成和持续交付 (CI/CD) 的实施方法有很多。作为设置持续交付的常用应用程序， [Jenkins](https://jenkins.io/) 以 OpenShift 中经认证的容器形式提供。您可以用它来构建包，运行单元和集成测试，以及部署镜像。

用于构建管道的另一个新型开源工具 [Tekton](https://github.com/tektoncd/pipeline)，提供了一种可执行许多相同操作的云原生方法。

## 命令行工具

- **buildah**： [buildah](https://docs.openshift.com/container-platform/3.11/crio/crio_runtime.html) 是用于构建基于 CRI-O 的镜像的工具，非常类似于 Podman。buildah 基本上是由通过 Podman 提供的构建命令超集所组成, 支持更细粒度控制镜像创建过程。

- **kubectl**： [kubectl](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands) 是用于控制 Kubernetes 集群的标准命令行工具。因为 OpenShift 3.x 及以上版本都是基于 Kubernetes 的，所以 `kubectl` 可以在任何 OpenShift 集群上使用。

- **oc**： [oc](https://docs.openshift.com/container-platform/3.10/cli_reference/differences_oc_kubectl.html) 是指 OpenShift 客户端 CLI，可用于处理作为一级对象（包括项目、应用程序、路由和 ImageStream）的 OpenShift 原生构造。因为 OpenShift 在 Kubernetes 上添加了这些元素，所以需要使用 `oc` 来连接特定于 OpenShift 的特性。

- **odo**： [odo](https://developers.redhat.com/blog/2019/08/14/openshift-development-with-interactive-odo/) 表示“OpenShift do”，是一个简化常见操作的命令行工具。它的目标受众是开发者（与运维人员相对），支持他们快速部署代码并不断迭代。

- **s2i**： [s2i](https://github.com/openshift/source-to-image) 是一个命令行工具，用于组合来自 GitHub 源代码存储库的构建器镜像。输出的是可运行的 Docker 镜像。构建器镜像就像是一个带有内嵌脚本的模板，用于获取源代码，并将其编译成可运行的应用程序。


## 镜像流

[镜像流](https://docs.openshift.com/container-platform/3.11/architecture/core_concepts/builds_and_image_streams.html) 是一个抽象，允许 OpenShift 从公共镜像注册表部署应用程序，同时在新镜像版本进入注册表时进行动态部署。您可以配置构建和部署以监视镜像流，并在新镜像版本可用后自动更新镜像流。

## 内部镜像注册表

OpenShift 的另一个与众不同之处就是其内置的 [镜像注册表](https://docs.openshift.com/container-platform/3.11/install_config/registry/index.html)。为什么要使用内部注册表？除了将镜像部署到 Docker Hub 或另一个在线注册表，它提供了另外一种选择。内部 OpenShift 注册表允许集群内的多个项目访问注册表，通过基于角色的访问控制 (RBAC) 实现细粒度安全性。须注意，如果删除 OpenShift 集群，那么也会删除内部注册表中存储的所有镜像。

### Kubernetes

[Kubernetes](https://kubernetes.io/) 是一个开源容器编排系统，用于自动进行应用部署、扩展和管理。它最初是由 Google 设计，现在由云原生计算基金会负责维护。

## Minishift

与 [Minikube](https://github.com/kubernetes/minikube) 很相似，但前者用于在本地机器上设置 Kubernetes 集群以供测试，而 [Minishift](https://github.com/minishift/minishift) 则用于运行 OpenShift 的开发实例。它可以在虚拟机内运行单节点集群，使用 [libmachine](https://github.com/docker/machine/tree/master/libmachine) 进行配置，并使用 [OKD](https://www.okd.io/) 为集群本身提供支持。

## OKD

OpenShift 由 [Origin Kubernetes Distribution (OKD)](https://www.okd.io/) 提供技术支持，这包括 Kubernetes 和其他开源项目，如 Docker 和 lstio。

## OpenShift 开发者控制台

OpenShift 与基础 Kubernetes 不同的另一个特性就是它的 [开发者控制台](https://docs.openshift.com/container-platform/3.11/getting_started/developers_console.html)，其中包含了丰富的功能。OpenShift Web 控制台为 OpenShift 环境提供中央控制点。Web 控制台由以下主要视图构成：

- **集群控制台**：OpenShift [Web 控制台](https://docs.openshift.com/container-platform/3.5/architecture/infrastructure_components/web_console.html) 包含一个对管理员操作更为友好的视图。它具有一个整体呈现 OpenShift 集群的全局视图，涵盖底层硬件的详细信息。您可以通过内置的 Grafana 或 Promethues 仪表板查看集群使用指标，也可以管理路由和入口。

- **应用程序控制台**：在应用程序控制台上，您可以创建和删除应用及部署，也可以管理部署。一站式查看项目中部署的所有 Pod 的日志，以及查看构建历史记录。

- **服务目录**：在服务目录上，您可以直接从控制台部署数据库和中间件应用程序。控制台中的服务预先经过认证，能够与 OpenShift 配合使用。对于 GitHub 存储库中的许多语言和源代码，您都可以通过基础镜像来创建 Pod。此外，您还可以创建 Jenkins 管道。


## Operator

[Operator](https://www.openshift.com/learn/topics/operators) 的作用非常类似于包：它们是安装特定软件的一种方式。然而，在 OpenShift 中，Operator 通过不断检查该软件的状态（比如版本号），并纠正与预期配置的任何偏差，可以提高软件部署质量。

## Podman

[Podman](https://podman.io/) 是无守护程序容器引擎，用于在 Linux 系统上开发、管理和运行 OCI 容器。您可以作为 root 用户或以无根模式运行容器。简而言之：`alias docker=podman`。

## 项目

OpenShift 中的 [项目](https://docs.openshift.com/container-platform/3.11/architecture/core_concepts/projects_and_users.html) 类似于 Kubernetes 中的名称空间概念，但它们比名称空间出现得更早，后者获取了它们的许多当前特性。从逻辑上讲，项目在 OpenShift 集群中支持多租户模式，允许多个团队在单个集群中安全地运行隔离的工作负载，避免了配置不必要的 OpenShift 实例引起的额外开销和成本。集群管理员可以为用户提供项目创建许可，然后以单个项目为基础控制访问权限。此外，可以通过一组预定义的对象创建模板，然后用作为新项目的基础。

## 路由

[路由](https://docs.openshift.com/container-platform/3.11/architecture/networking/routes.html) 是创建应用程序的方法，或者更确切地讲，是一种特定服务（可通过 URL 连接外部环境）。配置完集群后，主机名由管理员（或许就是你自己！）指定。

## 安全性

因为 OpenShift 旨在成为 Kubernetes 平台，供安全意识强的成熟企业用于支持云迁移工作，所以它包含的多个特性都旨在确保运行多租户容器化工作负载的安全。其他安全层 – _深度防御_ – 可以阻止“未知风险”领域中的漏洞，例如： [https://blog.openshift.com/openshift-protects-against-nasty-container-exploit/](https://blog.openshift.com/openshift-protects-against-nasty-container-exploit/)。

- 默认情况下，容器不以 root 用户身份运行。相反，它们会被动态分配一个用户 ID。
- [安全上下文约束](https://docs.openshift.com/container-platform/3.11/admin_guide/manage_scc.html) 用于控制可以运行的用户 pod 和它们可以访问的资源。
- 安装 OpenShift 后，默认启用 SELinux。

## Source-to-Image (S2I)

[OpenShift Source-to-Image](https://docs.openshift.com/enterprise/3.0/architecture/core_concepts/builds_and_image_streams.html#source-build) (S2I) 是用于构建可复制的 Docker 镜像的一种工具。它通过将应用程序源代码注入到 Docker 镜像中，并组成一个新的 Docker 镜像，从而生成可供运行的镜像。新镜像整合了基础镜像（构建器）和构建的源代码。可以通过 `docker run` 命令使用该镜像。S2I 支持增量构建，这些构建会重复使用之前下载的依赖项和之前构建的工件。

## 模板

[OpenShift 模板](https://docs.openshift.com/container-platform/3.7/dev_guide/templates.html) 通常被称为是一种方法，我们可以通过这种方法为 OpenShift Web 控制台填充快速启动应用程序和其他内容。不过，它们也是非常强大的工具，使用得当可以成为基础架构即代码解决方案的构件块，用于管理集群和应用程序状态的众多方面。

## 结束语

我们希望本指南对您的 OpenShift 之旅有所帮助。关于所使用的教程、视频和 Code Pattern，可参阅 [Red Hat OpenShift on IBM Cloud](https://developer.ibm.com/cn/collections/openshift-on-ibm/) 页面。

我们期待听到关于《OpenShift 漫游指南》一系列内容的反馈。如果您有任何补充或改进建议，请联系 [Anton McConville（Twitter 账号）](https://twitter.com/antonmc)。《银河系漫游指南》采用众包形式，我们的指南理应效仿！

本文翻译自： [A hitchhiker’s guide to OpenShift](https://developer.ibm.com/articles/hitchhikers-guide-openshift/)（2019-11-08）