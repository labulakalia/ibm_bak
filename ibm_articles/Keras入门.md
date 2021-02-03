# Keras 入门
探究此基于 Python 的深度学习库

**标签:** 人工智能

[原文链接](https://developer.ibm.com/zh/articles/cc-get-started-keras/)

Vinay Rao

发布: 2018-02-28

* * *

Keras 是基于 Python 的深度学习库，不同于其他深度学习框架。 [Keras](https://keras.io/) 充当神经网络的高级 API 规范。它既可作为用户界面，也可扩展它在其中运行的其他深度学习框架后端的功能。

Keras 起初是作为学术界热门 Theano 框架的简化前端。此后，Keras API 成为了 Google TensorFlow 的一部分。Keras 正式支持 Microsoft Cognitive Toolkit (CNTK)、Deeplearning4J，不久之后还将支持 Apache MXNet。

鉴于获得了广泛的支持，Keras 作为框架间迁移工具的地位已不可撼动。开发人员不仅可以移植深度学习神经网络算法和模型，还可以移植预先训练的网络和权重。

## 关于 Keras

**_Keras_ 名称的起源**

Chollet 创建了 Keras 作为开放式神经电子智能机器人操作系统 (ONEIROS) 机器人研究项目的神经网络的 API。 _ONEIROS_ 这一名称是对古希腊史诗《 _奥德赛_》的致意，在这部史诗中，神话人物 _Oneiroi_（ _Oneiros_ 的单数形式）为人类指明了两条进入梦境的路：一条路穿过宏伟的象牙之门进入噩梦，另一条路则穿过低矮的兽角之门，最终呈现一片神圣的景象。Keras 在希腊语中意为 _角_，这个名称非常合适，因为 Keras API 旨在为与神经网络协同使用提供一条捷径。

Keras 是开源 Python 包，由麻省理工学院 (MIT) 许可发行，由 François Chollet、Google、Microsoft 和其他贡献者共同持有该软件的部分版权。

Keras 前端支持在研究中快速构建神经网络模型的原型。此 API 易于学习和使用，并且具有易于在框架间移植模型的附加优势。

由于 Keras 的独立性，使用时无需与运行它的后端框架进行交互。Keras 具有自己的图形数据结构，用于定义计算图形：它不依靠底层后端框架的图形数据结构。此方法使您可以免于学习对后端框架进行编程，正因如此，Google 已选择将 Keras API 添加到其 TensorFlow 核心。

本文将概述 Keras，包括此框架的优势、支持的平台、安装注意事项以及支持的后端。

## Keras 的优势

为何要使用 Keras？它具有多种优势，包括：

- **更加良好的深度学习应用程序用户体验 (UX)。** Keras API 对用户友好。此 API 经过精心设计、面向对象且灵活易用，因而改善了用户体验。研究人员无需使用可能十分复杂的后端即可定义全新深度学习模型，从而实现了更简洁明了的代码。
- **无缝 Python 集成。** Keras 是本机 Python 包，能够轻松访问整个 Python 数据科学生态系统。例如，Python Scikit-learn API 也可以使用 Keras 模型。熟悉后端（如 TensorFlow）的开发人员同样可以使用 Python 来扩展 Keras。
- **大型的可移植工作主体和强大的知识库。** 目前，研究人员已将 Keras 与 Theano 后端结合使用了一段时间。由此产生了大型工作主体和强大的社区知识库，可供深度学习开发人员轻松地从 Theano 后端移植到 TensorFlow 后端。甚至还可以在后端之间移植权重，这意味着经过预先训练的模型只需稍作调整就可轻松切换后端。Keras 和 Theano 研究仍与 TensorFlow 和其他后端紧密相关。此外，Keras 还免费提供了许多学习资源、文档和代码样本。

## Keras 应用程序

借助拟合生成器、数据预处理和实时数据扩充等 Keras 功能，开发人员只需少量训练数据集便可训练出强大的图像分类器。Keras 随附预先经过训练的内置图像分类器模型，包括： [Inception-ResNet-v2](https://github.com/fchollet/keras/blob/master/keras/applications/inception_resnet_v2.py)、 [Inception-v3](https://github.com/fchollet/keras/blob/master/keras/applications/inception_v3.py)、 [MobileNet](https://github.com/fchollet/keras/blob/master/keras/applications/mobilenet.py)、 [ResNet-50](https://github.com/fchollet/keras/blob/master/keras/applications/resnet50.py) 、 [VGG16](https://github.com/fchollet/keras/blob/master/keras/applications/vgg16.py)、 [VGG19](https://github.com/fchollet/keras/blob/master/keras/applications/vgg19.py) 和 [Xception](https://github.com/fchollet/keras/blob/master/keras/applications/xception.py)。

注：由于这些模型的来源各不相同，因此有若干不同的许可证用于控制这些模型的权重使用情况。

借助 Keras，只需几行代码即可定义复杂模型。Keras 尤其适合用于通过小型训练数据集来训练卷积神经网络。虽然 Keras 在图像分类应用程序中已获得了更广泛的使用，它同样也适用于文本和语音的自然语言处理 (NLP) 应用程序。

## 哪些平台支持 Keras？

支持 Python 开发环境的平台同时也能支持 Keras。正式构建测试是在 Python V2.7x 和 V3.5 上运行的，但与 Keras 结合使用的后端需要特定平台才能访问支持的图形处理单元 (GPU)。大部分后端取决于其他软件，例如，NVIDIA® CUDA® 工具包和 CUDA 深度神经网络库 (cuDNN)。

TensorFlow 是 Keras 的缺省后端，但它还支持 Theano 和 CNTK 后端。Apache MXNet 的支持工作还在进行当中，同时也为 Keras 提供了 R 接口。许多供应商都已将 Keras API 移植到其深度学习产品中，由此而能够导入 Keras 模型。例如，基于 Java® 的后端 Eclipse Deeplearning4j，便能够导入 Keras 模型。此外，Scala 包装器也可用于 Keras。因此，Keras 平台支持成为了一个争议点。更重要的是确保目标平台支持您所选的 Keras 后端。

有关哪些平台支持 TensorFlow 的更多信息，请参阅 [TensorFlow 入门](https://www.ibm.com/developerworks/cn/cognitive/library/cc-get-started-tensorflow/index.html)。有关哪些平台支持 Theano 的详细信息，请阅读 [Theano 文档](http://deeplearning.net/software/theano/install.html)。有关哪些平台支持 CNTK 的更多信息，请参阅 [CNTK 文档](https://docs.microsoft.com/en-us/cognitive-toolkit/Setup-CNTK-on-your-machine)。

### 可选依赖关系

Keras 通过使用开源 Hierarchical Data Format 5 (HDF5) 二进制格式来管理数据。因此，它需要使用 HDF5 及其 h5py Python 包装器，才能将 Keras 模型保存至磁盘。Keras 通过使用开源 GraphViz DOT 格式来绘制图形。因此，它需要使用 GraphViz 及其 pydot Python 包装器，才能直观显示数据。Keras GPU 支持还需要使用 cuDNN 库。

### 从源代码构建 Keras

由于 Keras 是一个纯 Python 包，因此没有理由从源代码进行构建。Keras 不含任何特定于平台的后端代码。强烈建议改为从 Python Package Index (PyPI) 安装 Keras。

## Keras 安装注意事项

上文中已提到，Keras 可在支持 Python 开发环境的任何平台上运行。这足以训练和测试大部分简单示例和教程。大部分专家强烈建议，针对研究或商业开发等应用程序使用高性能计算 (HPC) 平台。

由于 Keras 使用第三方后端，因此无任何安装注意事项。后端将负责执行硬件加速。总之，安装 Keras 后端的开发人员应考虑以下因素和选项：

- [处理器和内存需求](#处理器和内存需求)
- [虚拟机选项](#虚拟机选项)
- [Docker 安装选项](#docker-安装选项)
- [云安装选项](#云安装选项)

### 处理器和内存需求

深度学习算法为计算密集型算法，至少需要一个具有矢量扩展的快速多核 CPU。此外，一个或多个支持高端 CUDA 的 GPU 卡是深度学习环境的标配。

深度学习进程通过使用共享内存中的缓冲区相互进行通信。因此，分配的内存应已足够。大部分专家由此还建议采用较大的 CPU 和 GPU RAM，因为从性能和能源使用的角度来看，内存传输非常昂贵。更大的 RAM 可避免这些操作。

### 虚拟机选项

用于深度学习的虚拟机 (VM) 目前最适合有许多核心的以 CPU 为中心的硬件。因为主机操作系统控制着物理 GPU，所以在 VM 上实现 GPU 加速很复杂。主要有两种方法：

- GPU 直通：

    - 仅适用于 1 类管理程序，比如 Citrix Xen、VMware ESXi、Kernel Virtual Machine 和 IBM® Power® 。
    - 基于 CPU、芯片集、管理程序和操作系统的特定组合，直通方法的开销可能会有所不同。通常，对于最新一代的硬件，开销要低得多。
    - 一种给定的管理程序-操作系统组合仅支持特定的 NVIDIA GPU 卡。
- GPU 虚拟化：

    - 各大主要 GPU 供应商均支持，包括 NVIDIA GRID™ 、AMD MxGPU 和 Intel® Graphics Virtualization Technology。
    - 最新版本在特定的较新 GPU 上支持开放计算语言 (OpenCL)。在大部分主要后端（包括 TensorFlow）上，不存在正式的 OpenCL 支持。
    - 最新版的 NVIDIA GRID 可以在特定的较新 GPU 上支持 CUDA 和 OpenCL。

### Docker 安装选项

在 Docker 容器或 Kubernetes 集群中运行 Keras 存在诸多优势。Keras 存储库包含一个 Docker 文件，具有针对 Mac OS X 和 Ubuntu 的 CUDA 支持。此映像支持 Theano 或 TensorFlow 后端。使用 Docker 的优势主要在于，后端可以访问物理 GPU 核心（设备）并在其中运行。

### 云安装选项

在云服务上运行 Keras 时有许多选项。Keras 可用于在一个供应商生态系统上训练模型，但只需稍作调整即可在另一个供应商生态系统上用于生产部署。

- **IBM Cloud®** 数据科学和数据管理为 Python 环境提供了 Jupyter Notebook 和 Spark。Keras 和 TensorFlow 已预先安装。IBM Cloud 上的 Kubernetes 集群可运行 Keras 和 TensorFlow Docker 映像。
- **Google Cloud**：Google 提供了一些机器实例，它们能访问特定区域的 1、4 或 8 个 NVIDIA GPU 设备。同时还可以在容器化的 GPU 支持的 Jupyter Notebook 上运行 Keras 和 TensorFlow
- **Amazon Web Services**：Amazon 提供了 Amazon Web Services 深度学习 Amazon Machine Image (AMI)，以及可选的 NVIDIA GPU 支持，后者能在各种各样的 Amazon Elastic Compute Cloud 实例上运行。已预先安装 Keras、TensorFlow 和其他深度学习框架。AMI 可以支持多达 64 个 CPU 核心和多达 8 个 NVIDIA GPU (K80)。
- **Microsoft Azure**：您可以在 Microsoft 数据科学虚拟机系列的 Microsoft Azure 机器实例上以 CNTK 后端安装 Keras，仅使用 CPU 或者包含最多四个 K80 GPU 都可。

## 将 Keras 用作其他框架的 API

Keras 各层及模型均与纯 TensorFlow 张量完全兼容；因此，Keras 为 TensorFlow 提供了良好的模型定义附件。您甚至可以同时使用 Keras 和其他 TensorFlow 库。Keras 现已成为 TensorFlow 核心的正式组成部分。有关详细信息，请阅读此 [博客文章](https://blog.keras.io/keras-as-a-simplified-interface-to-tensorflow-tutorial.html)。

从 TensorFlow 后端切换至其他某个正式支持的 Keras 后端十分简单，只需在 JavaScript 对象表示法 (JSON) 配置文件中进行一项更改即可。有关详细信息，请参阅 [Keras 文档](https://keras.io/backend)。

目前，可以使用 Keras 作为以下这些框架的 API：

- **Keras 与 Theano。** 最近淘汰的 Theano 是 Keras 最先支持的后端，后被 TensorFlow 所取代。TensorFlow 支持大部分 Theano 开发的 Keras 模型。要使用 GPU 运行 Theano 后端，请遵循 [此文档](https://keras.io/getting-started/faq/#how-can-i-run-keras-on-gpu) 中有关 Theano 的部分来操作。
- **Keras 与 CNTK。** Keras 对 Microsoft Cognitive Toolkit (CNTK) 后端的支持仍处于测试阶段。您可以阅读 [Microsoft 文档](https://docs.microsoft.com/en-us/cognitive-toolkit/Using-CNTK-with-Keras)，了解更多详细信息和资料。
- **Keras 与 Deeplearning4j。** Deeplearing4j 可使用其 `deeplearing4j-modelimport` 模块来导入大部分 Keras 模型。目前，Deeplearning4j 可支持导入有关层次、损失、激活、初始化程序、正则化项、约束条件、度量标准和优化程序的模型信息。
- **Keras 与 Apache MXNet。** Keras 对 Apache MXNet 后端的支持仍处于早期测试阶段。这是由 Distributed (Deep) Machine Learning Community 主导的工作。此后端正逐步变为另一个正式支持的 Keras 后端。此 [GitHub 存储库](https://github.com/dmlc/keras) 中提供了该后端的代码。

## 结束语

Keras 不同于其他深度学习框架。按照设计，它旨在成为神经网络建模的 API 规范。它可作为用户界面，也可扩展它在其中运行的其他深度学习框架后端的功能。

Keras API 已成为 Google TensorFlow 的一部分。Keras 同时还正式支持 CNTK、Deeplearning4j，很快就会支持 Apache MXNet。

由于这一广泛的支持，Keras 已成为了实现框架间迁移的实际工具。开发人员不仅可以移植深度学习神经网络算法和模型，还可以移植预先训练的网络和权重。

本文翻译自： [Get started with Keras](https://developer.ibm.com/articles/cc-get-started-keras/)（2017-12-17）