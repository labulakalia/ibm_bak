# 使用 Docker 容器快速开发物联网应用程序
使用基于容器的虚拟化来开发、测试、部署和更新 IoT 应用程序

**标签:** Docker,IoT

[原文链接](https://developer.ibm.com/zh/articles/iot-docker-containers/)

Anna Gerber

更新: 2017-07-07 \| 发布: 2015-08-25

* * *

物联网正在快速发展，它是智慧设备的高度互联网络，这些设备包括环境传感器、健康跟踪器、家用电器和工业设备等。到 2020 年，预计有 200 亿个设备将连入网络，这超过 PC、智能手机和平板电脑加起来的数量的两倍。开发人员正在快速开始为 IoT 创建应用程序，而使用容器可在不同方面为他们带来帮助。

容器是一种轻量型的虚拟化方法，开发人员可使用该方法快速、大规模地开发、测试、部署和更新 IoT 应用程序。许多 Web 和移动应用程序开发人员使用虚拟机管理程序（比如 VirtualBox）来运行虚拟机 (VM)，在一个跨平台的开发、测试和部署工作流中虚拟化物理硬件。

基于容器的虚拟化（有时称为操作系统级虚拟化）要轻量得多。每个容器作为隔离的用户空间实例在一个共享的主机操作系统内核上运行。尽管操作系统是共享的，但各个容器拥有独立的虚拟网络接口、独立的流程空间和不同的文件系统。可使用实现资源隔离的控制组来向这些容器分配系统资源，比如 RAM。在基于虚拟机管理程序的虚拟化中，每个 VM 运行自己的操作系统，这会增加系统资源的使用，与此相比，容器使用的磁盘和内存资源少得多。

[**Docker**](https://www.docker.com/whatisdocker/) 是 Linux 上基于容器的虚拟化的开放平台。通过 Docker，可以更快速、更轻松地构建容器，并且几乎可以在任何地方部署这些容器：在私有云或公共云中、在本地 VM 中或者在物理硬件（包括物联网设备）上。IBM Containers 是基于 Docker 和 Kubernetes 的 IBM Cloud 功能部件，用于在 IBM Cloud 平台上交付和部署容器化应用程序。

## 设置 Docker 以用于物联网开发

物联网应用程序面对的是各种各样的设备平台。在原型设计和开发的早期阶段，开发物联网项目时，通常采用的是基于通用微控制器的开发板或单板计算机 (SBCS)，比如 Raspberry Pi。

要为 Raspberry PI 设置 Docker，您将需要一个满足以下条件的 Raspberry PI：具备 wifi（例如，Raspberry Pi Zero W 或 Raspberry Pi 3），包含写入 Raspbian Jesse 镜像的 microSD 卡，且启用了 SSH。（ [观看这个 IoT 实践视频，了解如何设置 Raspberry Pi](https://developer.ibm.com/tv/setting-up-your-raspberry-pi-without-a-keyboard-mouse-or-monitor/)。）您可以通过在 microSD 卡的引导分区上创建一个名为“ssh”的文件来启用 SSH。可以通过 SSH（默认密码为 raspberry）从 Mac 或 PC 以无头模式来访问 Raspberry Pi：

`ssh pi@raspberrypi.local`

在 raspbian 命令行中，可以通过运行以下脚本来安装 Docker：

`$ curl -sSL get.docker.com |sh`

您需要将默认用户 (pi) 添加到 docker 组以避免权限问题，然后将用户 (su) 切换到 pi 以使更改生效，或者重新引导 Raspberry pi。

```

$ sudo usermod ‑aG docker pi
$ su pi

```

Show moreShow more icon

一旦在 Raspberry Pi 上安装了 Docker，就可以运行 docker ps 来确认 Docker 正在运行。

```

$ docker ps

CONTAINER ID  IMAGE  COMMAND  CREATED  STATUS  PORTS  NAMES

```

Show moreShow more icon

## 使用容器、镜像和注册表

每个 Docker 容器包含一个或多个运行的进程。Docker 容器从一个镜像启动，该镜像指定了以下元素：

- 要在容器内运行的应用程序的配置信息
- 它的依赖项，比如库和共享二进制文件

Docker 使用一种统一的文件系统，将镜像存储为一系列层。层会在构建过程中缓存，这使衍生镜像的构建既快捷又高效。

这种分层方法还意味着 Docker 镜像很小且可移植，通过将它们发布到公共或私有注册表，可以轻松地共享这些镜像。 [Docker Hub](https://hub.docker.com/) 是 Docker 镜像存储库的最大注册表；它列出了超过 100,000 个存储库，并托管了许多私有存储库。如果您正在使用 IBM Cloud 平台，那么支持从 Docker Hub 中提取公共镜像，并且还提供托管的私有镜像存储库，这些存储库允许您在组织内共享私有镜像。

如果使用一种流行的开源框架或服务，您可能能够在 Docker Hub 中找到一个预先构建的公共镜像，许多镜像由与这些项目相关的开源社区来维护。Docker Hub Web 界面有一个较小的官方存储库列表，这是一个精挑细选的存储库列表，包含由 Docker 团队对已知安全漏洞进行了测试的镜像。

避免使用旧的 Docker 镜像，即使这些镜像的来源可靠，因为许多旧镜像仅用于归档和测试用途。如果旧镜像仍得到维护，它们会在 Docker Hub 中被标记为受支持。公共注册表中的大多数镜像都是社区贡献的，可能没有详细备案或维护。怀疑有问题时，可使用一个官方镜像作为基础来创建自己的镜像。

可完成以下步骤来创建镜像：

1. 从注册表拉入一个基础镜像
2. 交互式地运行一系列命令
3. 提交结果作为新镜像

例如，下面的命令列表从 Docker Hub 获取官方的 Ubuntu version 14.04 镜像，从该镜像启动容器，运行一个命令来将 Git 安装到该容器中，然后列出创建的容器的 ID。

##### 创建一个镜像

```
$ docker pull ubuntu:14.04
$ docker run ubuntu apt-get install -y git
$ docker ps -l

```

Show moreShow more icon

可发出下面这条命令来将此容器 (cb23e345fde0) 保存为新镜像（名为 demo/git）： `$ docker commit cb23e345fde0 demo/git` 。

要使用 IBM Container Extensions 工具在本地运行 Docker 命令，可将 `docker` 命令更改为 `ice --local` 。例如，使用此命令： `$ ice --local pull ubuntu:14.04`

如果需要运行不止一两个命令来设置应用程序环境，您需要创建一个 Dockerfile。Dockerfile 是指定用于构建镜像的一组指令的文本文件。下面这个代码清单中的 Dockerfile 创建了一个运行 Nginx Web 服务器的基础镜像。

##### 创建一个运行 Nginx Web 服务器的基础镜像的 Dockerfile

```
FROM ubuntu:14.04
RUN apt-get -y update && apt-get install -y nginx
RUN mkdir -p /var/www/html
EXPOSE 80
EXPOSE 443
WORKDIR /etc/nginx
CMD ["nginx"]

```

Show moreShow more icon

每个 Dockerfile 中的第一条指令是 `FROM` 指令，这表明基础镜像。Dockerfile 中的每个后续指令存储为一个层：

- `RUN` 指令运行 `apt-get` 等 Linux 命令。
- `ADD` 和 `COPY` 指令将应用程序文件添加到容器中。
- `EXPOSE` 指令打开端口。
- `ENV` 指令配置环境变量。

最后，每个 Dockerfile 包含 `WORKDIR` 、 `ENTRYPOINT` 或 `CMD` 指令，用来指定在该容器启动时，如何和在何处运行该应用程序。

运行 Docker 容器需要极少的开销，所以可将您的应用程序分解为一组在不同的容器中运行的服务，这些容器可在运行时按名称来链接。例如，您可将应用程序拆分为一个运行 Node.js 应用程序的容器，该容器链接到另一个运行 Redis 键值存储的容器。

## 使用容器维护一致的 IoT 开发环境

为物联网设备开发的嵌入式应用程序以后可以升级，从而在自定义原型开发板上运行，最终在生产设备上运行。您可能会开始使用现成的硬件，比如 Raspberry Pi，然后评估一系列潜在的微控制器和片上系统 (SoC) 设备，设备驱动程序需求各不相同，有时甚至不兼容。每个设备或经过修正的设备可能还需要不同版本和配置的开发工具包，用于设备闪存、监控和通信。

容器可用于捕获适用于每个修正设备的已知开发环境，并在开发者团队之间共享此环境。例如，如果您的团队正在使用几种类型的 Arduino 兼容开发板，您可以创建一个包含 Arduino 命令行工具包的镜像和一个用于串行通信的哑终端仿真程序，作为您的基线开发镜像。然后，您可以使用这个基本 Arduino 开发镜像为每个需要特定自定义驱动程序的开发板创建新的变体镜像。

Docker 的分层文件系统在这种情况下有效地利用了空间，因为在创建新镜像时只存储了唯一的层。Docker history 命令显示组成镜像的各个层的列表，包括创建每个层的指令和层大小，例如：

```
IMAGE           CREATED         CREATED BY              SIZE
3ca43a9a4095    20 hours ago    apt-get install -y git  37.64 MB

```

Show moreShow more icon

层会被缓存，所以您可快速地尝试镜像的新变体。例如，您可能希望向驱动程序或开发工具应用更新。层还使您能够轻松地回滚到镜像的早期版本，在某处出错时尝试一种不同的方法。

如果对开发环境的更改取得成功，可将新镜像推送到您的私有团队注册表，以快速将更改传播给团队的其余成员。每个团队成员拉入更新的镜像时，他们已缓存了大部分层，所以 Docker 通常只需几秒到几分钟时间运行创建新层的指令，即可让新镜像正常运行。与推送和拉入 VM 镜像快照或为每个变体从头构建完整镜像相比，层节省了存储空间，而且更重要的是节省了开发人员的时间。

## 将容器部署到 IoT 设备

Docker 化应用程序的魅力在于，在您构建了一个镜像之后，几乎可以在任何地方发送和运行此镜像。如果您的物联网设备运行 Linux，那么可以直接在设备上部署容器。尽管此类物联网设备上可用的系统资源有限，但部署 Docker 容器是可行的，因为它们的运行时开销几乎为零。您甚至可以在设备上运行多个容器，例如，在您希望并行运行应用程序的不同版本以进行比较的情况下，就可以执行此操作。

Docker 需要一个支持内核名称空间和控制组的现代 Linux 内核，所以您首选的物联网设备可能还没有合适的基本镜像。随着合适的 SoC 的价格不断下降（例如，9 美元的 C.H.I.P. 开发板），预计会开发越来越多的支持 Linux 的物联网设备；这样一来，就会为这些平台提供更广泛的可供选择的 Linux 发行版。

Docker Hub 包含了许多基于 ARM 架构的镜像，适用于基于 ARM 架构的数码爱好者喜欢的 SBC，包括 [Raspberry Pi](https://www.raspberrypi.org/)、Orange Pi 和 [BeagleBone Black](http://beagleboard.org/BLACK)。DockerHub 上的 arm32v7、arm32v6 和弃用的 [armhf](https://hub.docker.com/u/armhf/) 镜像是官方镜像，是为了支持 DockerHub 的多架构而专门为 ARM 架构所构建的。DockerHub 上的 [resin](https://hub.docker.com/u/resin/) 组织也提供了第三方基本镜像，旨在用于 resin.io 和 resinOS。

物联网开发者面临的挑战是，要持续更新互联设备网络，并在可能并不可靠的低带宽无线连接上运行。此外，他们还必须考虑如何维护其中许多设备收集的高度个人化的敏感数据的安全性。Docker 容器有助于提高物联网的安全性，它们通过 Docker 1.10 中引入的用户名称空间以及 Docker 1.13 和更高版本中提供的 [Docker 密钥](https://docs.docker.com/engine/swarm/secrets/) 支持 [隔离](https://www.rtinsights.com/docker-containers-for-the-iot/)，这些 Docker 密钥可用来存储和加密密钥或者访问令牌（用于身份验证以及设备之间或上游与云服务之间的安全通信）。

物联网设备可能只是偶尔连接，为了节省电力，有些设备设置为按不同的时间表进行休眠。其他物联网设备可能通过使用 Mesh 网络拓扑进行连接，在这种拓扑中，在任何给定时间，只有网络的某些部分是可以访问的。当网络连接中断时，基于推送的更新可能会失败，而应用不完整或可能损坏的更新造成的危险是设备可能进入不一致的状态或无法操作。

Docker 为这种更新问题提供了可行的解决方案。对于请求拉取最新版本应用程序镜像的设备，仅通过空中下载技术向其发送镜像的差异部分，而不是整个镜像。完成基于差异的更新速度更快，这将减少设备需要连接的时间，降低故障的可能性，从而减少对低带宽网络的压力。这样就可以更频繁地应用更新。

有些物联网设备允许通过物理按钮或小触摸屏进行有限的用户交互，但是对于许多当代物联网设备来说，用户交互的主要手段是使用移动应用程序。人们越来越多地采用 Docker 容器来加速移动应用程序的自动化测试、持续集成和交付，这一点也就不足为奇了。

## 集成 IoT 设备与云

将应用程序部署到智能设备和移动设备只是物联网发展之路上的一部分。在物联网设备和网关上运行的容器中，可以执行边缘分析，但是，许多当代物联网设备还会将传感器数据和事件发布到云服务以供进一步处理。

云服务经常采用容器。但是，容器是短暂的，因此任何写入运行的数据库容器中文件系统的数据都应该被视为临时的；也就是说，如果关闭容器，数据将会丢失。从某种意义上说，使用容器会迫使您养成开发无状态云应用程序的良好习惯。任何要持久化的数据都应该使用数据卷来存储。数据卷在容器关闭后持续存在，并且可以在多个容器之间共享。

随着互联设备数量的增加，物联网的基于云的应用程序将需要不断扩展才能处理所生成的数据量。幸运的是，我们可以利用围绕 Docker 构建的不断壮大的编排工具生态系统来开发可伸缩的云应用程序，这包括 Docker Machine、Swarm 和 Compose 工具。这些工具可以在许多云平台上使用，包括 EC2 Container Services、Microsoft Azure 和 IBM Cloud。IBM Cloud 目前还支持运行多组容器实现负载平衡和故障转移。

## 结束语

为了满足物联网应用程序的预期需求，开发者需要采用相关工具和实践，使他们能够快速开发物联网应用程序和服务，以便在智能设备、移动设备和云端运行。Docker 容器能够在广泛的设备上随时随地发送和运行应用程序，运行时开销最小，支持自动化，并且具有分层文件系统（实现轻量级可移植镜像和快速镜像构建），因而是物联网开发者的绝佳工具。

本文翻译自： [Rapidly develop Internet of Things apps with Docker Containers](https://developer.ibm.com/articles/iot-docker-containers/)（2017-07-07）