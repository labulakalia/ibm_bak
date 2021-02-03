# PyTorch 入门
了解这个原生 Python 包如何在 Python 中重新设计和实现 Torch

**标签:** 人工智能

[原文链接](https://developer.ibm.com/zh/articles/cc-get-started-pytorch/)

Vinay Rao

发布: 2018-03-21

* * *

PyTorch 是根据经过改良的 Berkeley Software Distribution 许可发布的一个开源 Python 包。Facebook、Idiap Research Institute、New York University (NYU) 和 NEC Labs America 拥有 [PyTorch](http://pytorch.org/) 的版权。尽管 Python 是数据科学的首选语言，但 PyTorch 是一个相对较新的深度学习领域。

本文将概述 PyTorch 系统并进一步介绍以下主题：

- PyTorch 的背景
- 使用 PyTorch 的好处
- 典型的 PyTorch 应用
- 支持的平台
- 硬件加速
- 安装考虑因素和选项

## PyTorch 概述

神经网络算法通常计算一个损失函数的峰值或谷值，主要使用一个梯度下降函数来实现此目标。在 Torch（PyTorch 的前代产品）中，Twitter 所贡献的 Torch Autograd 包可以计算梯度函数。Torch Autograd 基于 Python Autograd。

Autograd 和 Chainer 等 Python 包都使用了一种称为 _[基于磁带的自动微分 (tape-based auto-differentiation)](http://pytorch.org/docs/master/autograd.html)_ 的技术来计算梯度。因此，这些包显著影响和启发了 PyTorch 的设计。

**用于研究的 PyTorch**

PyTorch 的主要用途是研究。Facebook 使用 PyTorch 进行创新研究，并将 Caffe2 用于生产用途。一种名为 _[Open Neural Network Exchange](https://onnx.ai/)_ 的新格式允许用户在 PyTorch 和 Caffe2 之间转换模型，减少研究与生产之间的滞后时间。

基于磁带的自动微分的工作方式类似于一个磁带录音机，因为它会记录执行的操作，然后回放它们来计算梯度 – 这种方法也称为\_ [反向模式自动微分](https://arxiv.org/pdf/1509.07164.pdf)。\_PyTorch Autograd 拥有此函数的最快实现之一。

使用此特性，PyTorch 用户可通过任意方式调节其神经网络，而没有任何开销或滞后惩罚。因此，与在大部分著名框架中不同，PyTorch 用户可以动态构建图表，而且该框架的速度和灵活性促进了新的深度学习算法的研究和开发。

此性能一定程度上源于 PyTorch 核心中使用的模块化设计。PyTorch 将 CPU 和图形处理单元 (GPU) 的大部分张量和神经网络后端实现为单独的、精简的、基于 `C` 的模块，集成了数学加速库来提高速度。

PyTorch 与 Python 是无缝集成的，而且特意使用了命令式编码风格。此外，该设计保留了可扩展性，该特性使得其基于 Lua 的前代产品得以流行。用户可以使用 `C/C` \+\+ 借助一个基于 `C` Foreign Function Interface (cFFI) for Python 的扩展应用编程接口 (API) 来设计程序。

## 背景

在深入剖析 PyTorch 的功能和它带来的优势之前，让我们了解一些背景知识。

### Torch 是什么？

Torch 是一个用于机器学习和科学计算的模块化开源库。Torch 最初是 NYU 的研究人员为学术研究而开发的。该库通过对 LuaJIT 编译器的利用提高了性能，而且基于 `C` 的 NVIDIA CUDA 扩展使得 Torch 能够利用 GPU 加速。

许多开发人员使用 Torch 作为受 GPU 支持的 NumPy 替代方案；其他开发人员使用它来开发深度学习算法。Torch 得以闻名，源于 Facebook 和 Twitter 对它的使用。Google DeepMind 人工智能 (AI) 项目最初使用的是 Torch，随后切换到了 TensorFlow。

### Lua 和 LuaJIT 是什么？

Lua 是一种支持多种编程模型的轻量型脚本语言；它源自于应用程序的可扩展性。Lua 非常紧凑，而且是用 `C` 编写的，这使它能够在资源受限的嵌入式平台上运行。巴西的 Pontifical Catholic University of Rio de Janeiro 于 1993 年首次介绍了 Lua。

LuaJIT 是一种即时 (JIT) 编译器，采用特定于平台的优化来提高 Lua 的性能。它还扩展并增强了标准 Lua 的 `C` 应用编程接口 (API)。

## PyTorch 是什么？

PyTorch 有两种变体。最初，Hugh Perkins 开发了 `pytorch` 作为基于 LuaJIT 的 Torch 框架的 Python 包装器。但是，本文讨论的 PyTorch 变体是一种全新的开发。不同于旧的变体，PyTorch 不再使用 Lua 语言和 LuaJIT。相反，它是一个原生的 Python 包。

PyTorch 在 Python 中重新设计和实现了 Torch，同时在后端代码中共享了相同的核心 `C` 库。PyTorch 开发人员调节了其后端代码来高效运行 Python。他们还保留了基于 GPU 的硬件加速，以及让基于 Lua 的 Torch 在研究人员中得以流行的可扩展特性。

## PyTorch 的优势

PyTorch 有众多优势。让我们看看一些主要优势。

### 动态计算图表

大部分使用计算图表的深度学习框架都会在运行时之前生成和分析图表。相反，PyTorch 在运行时期间使用反向模式自动微分来构建图表。因此，对模型的随意更改不会增加运行时间延迟和重建模型的开销。

PyTorch 拥有反向模式自动微分的最快实现之一。除了更容易调试之外，动态图表还使得 PyTorch 能够处理可变长度输入和输出，这在对文本和语音的自然语言处理 (NLP) 中特别有用。

### 精简的后端

PyTorch 没有使用单一后端，而是对 CPU、GPU 和不同的功能特性使用了单独的后端。例如，针对 CPU 的张量后端为 TH，而针对 GPU 的张量后端为 THC。类似地，针对 CPU 和 GPU 的神经网络后端分别是 THNN 和 THCUNN。单个后端得以获得精简的代码，这些代码高度关注在特定类型的处理器上以高内存效率运行的特定任务。单独后端的使用使得将 PyTorch 部署在资源受限的系统上变得更容易，比如嵌入式应用程序中使用的系统。

### Python 优先的方法

尽管 PyTorch 是 Torch 的衍生产品，但它被特意设计为一个原生的 Python 包。它不是一种 Python 语言绑定，而被作为 Python 的一个完整部分。PyTorch 将它的所有功能都构建为 Python 类。因此，PyTorch 代码能够与 Python 函数和其他 Python 包无缝集成。

### 命令式编程风格

因为对程序状态的直接更改会触发计算，所以代码执行不会延迟，而且会生成简单的代码，这避免了可能扰乱代码执行方式的大量异步执行。处理数据结构时，这种风格不但直观，而且容易调试。

### 高度可扩展

用户可以使用 `C/C` \+\+ 来编程，只需使用一个基于 cFFI for Python 并针对 CPU 进行了编译的扩展 API，或者使用 CUDA for GPU 操作即可。此特性有助于针对新的试验用途而扩展 PyTorch，使它对研究用途很有吸引力。例如，PyTorch 音频扩展允许加载音频文件。

## 典型的 PyTorch 应用

Torch 和 PyTorch 共享相同的后端代码，而且文献中通常在基于 Lua 的 Torch 与 PyTorch 之间存在许多混淆。因此，除非查看时间轴，否则很难区分二者。

例如，Google DeepMind AI 项目在切换到 TensorFlow 之前使用了 Torch。因为切换是在 PyTorch 出现之前发生的，所以不能将它视为 PyTorch 应用的示例。Twitter 是 Torch 的贡献者，现在使用 TensorFlow 和 PyTorch 来调节它在时间轴上的排名算法。

Torch 和 PyTorch 都被 Facebook 广泛使用来研究文本和语音的 NLP。Facebook 发布了许多开源 PyTorch 项目，包括与以下方面相关的项目：

- 聊天机器人
- 机器翻译
- 文本搜索
- 文本到语音转换
- 图像和视频分类

PyTorch Torchvision 包使得用户能够访问模型架构和预先训练的流行的图像分类模型，比如 AlexNet、VGG 和 ResNet。

由于灵活的、可扩展的、模块化的设计，PyTorch 没有将您限制于特定的模型或应用。可以使用 PyTorch 取代 NumPy，或者实现机器学习和深度学习算法。要获得关于应用和贡献的模型的更多信息，请参阅 [https://github.com/pytorch/examples](https://github.com/pytorch/examples) 上的 PyTorch Examples 页面。

## 支持 PyTorch 的平台

PyTorch 需要 64 位处理器、操作系统和 Python 发行版才能正常工作。为了访问受支持的 GPU，PyTorch 需要依靠 CUDA 等其他软件。在发表本文时，最新的 PyTorch 版本为 0.2.0\_4。 [支持 PyTorch 的操作系统](#支持-pytorch-的操作系统) 给出了预先构建的 PyTorch 二进制文件的可用性，以及对此版本的 GPU 支持。

## 支持 PyTorch 的操作系统

操作系统Python 2.7Python 3.5Python 3.664 位 Linux®是是是_是否支持 CPU 和 GPU？_都支持都支持都支持Mac OS X是是是_是否支持 CPU 和 GPU？_仅支持 CPU仅支持 CPU仅支持 CPU64 位 Windows®否否否_是否支持 CPU 和 GPU？_否否否

**备注：**

- 不要使用主要 Anaconda 通道进行安装。使用来自 PyTorch 维护者 soumith 的通道，以确保支持更高版本的 CUDA 及恰当优化的 CPU 和 GPU 后端，以及支持 Mac OS X。
- 在 Mac OS X 上，可以通过从源代码构建 PyTorch 来启用 CUDA GPU。
- PyTorch 没有正式支持 Windows。Python 3.x 上的 Windows 支持是由社区创建的，而且是试验性的。虽然提供了 CPU 支持，但 GPU 支持是局部的，因为该构建版本仅部分支持 CUDA。

要获得有关更多信息，请参阅 [https://github.com/tylergenter/pytorch/blob/master/README\_Windows.md](https://github.com/tylergenter/pytorch/blob/master/README_Windows.md) 上的 README 文件。

### 从源代码构建 PyTorch

PyTorch 构建版本使用 CMake 进行构建管理。要从源代码构建 PyTorch 并确保拥有高质量的 Basic Linear Algebra Subprograms (BLAS) 库（Intel® Math Kernel Library 或简称 MKL），包维护者推荐使用 Anaconda Python 发行版。无论使用哪个操作系统，此发行版都可以确保拥有一个受控的编译器版本。

要获得有关细节，请参阅 [https://github.com/pytorch/pytorch#from-source](https://github.com/pytorch/pytorch#from-source)。

### PPC64le 架构上的 PyTorch

PyTorch 没有正式支持 PPC64le 架构，比如 IBM® POWER8® ，但您可以使用一个 Linux 构建 recipe 从源代码构建 PyTorch。这些构建版本仍是试验性的，而且没有通过所有测试，尤其是在启用 CUDA 时。此外，数学内核库和线性代数库没有得到全面优化。

## PyTorch 如何使用硬件加速？

PyTorch 没有使用单个后端，而对 CPU、GPU 和不同的任务使用了多个不同的后端模块。此方法意味着每个模块都是精简和重点明确的，并针对任务进行了高度优化。这些模块的性能高度依赖于基础数学和线性代数库（MKL 和 BLAS）。

### CPU

因为 Anaconda Python 针对 Intel 架构调节了 MKL，所以 PyTorch 在 Intel CPU 上的性能最佳。Intel 为其以 CPU 为中心的高性能计算 (HPC) 架构（比如 Intel Xeon® 和 Xeon Phi™ 处理器系列）提供了 MKL-Deep Neural Network 原语。

PyTorch 也可以使用 OpenBLAS。在 IBM PowerPC® 机器上，IBM 提供了 [Math Acceleration Subsystem (MASS) for Linux](https://www-01.ibm.com/software/awdtools/mass/linux/mass-linux.html)。您可以构建支持 MASS 的 OpenBLAS。IBM 最近还宣布将 Anaconda 集成到其基于 IBM Power Systems™ S822LC 平台的 HPC 平台上。因此，PyTorch 可能在 IBM 系统上也有不错的性能。

### GPU

PyTorch 仅支持 NVIDIA GPU 卡。在 GPU 上，PyTorch 使用了 NVIDIA CUDA Deep Neural Network (CuDNN) 库，这个 GPU 加速库是为深度学习算法而设计的。Open Computing Language (OpenCL) 支持不在 PyTorch 的发展路线图上，但基于 Lua 的 Torch 为该语言提供了一些有限的支持。目前也没有社区开展移植到 OpenCL 的工作。

PyTorch 在 CUDA 开发曲线上仍然有些落后。PyTorch 使用了 CUDA V8，而其他许多深度学习库已在使用 CUDA V9。早期的 PyTorch 版本都基于 CUDA 7 和 7.5。因此，PyTorch 用户无法利用最新的 NVIDIA 显卡。通过 CUDA 9 从源代码构建，目前会生成混杂的结果和不一致的行为。

## 安装考虑因素

通常，PyTorch 可以在任何支持 64 位 Python 开发环境的平台上运行。这足以训练和测试大部分简单示例和教程。但是，对于研究或专业开发，大部分专家都强烈建议使用 HPC 平台。

### 处理器和内存需求

深度学习算法为计算密集型算法，至少需要一个具有矢量扩展的快速多核 CPU。此外，一个或多个支持高端 CUDA 的 GPU 是深度学习环境的标配。

PyTorch 中的处理器使用共享内存中的缓冲区来相互通信，所以分配的内存必须足够用于此用途。大部分专家因此还建议采用较大的 CPU 和 GPU RAM，因为从性能和能源使用的角度来看，内存传输非常昂贵。更大的 RAM 可避免这些操作。

### 虚拟机选项

用于深度学习的虚拟机 (VM) 目前最适合包含多核的以 CPU 为中心的硬件。因为主机操作系统控制着物理 GPU，所以在 VM 上实现 GPU 加速具有挑战性。尽管如此，仍有两种主要方法可实现此目的：

- GPU 直通：

    - GPU 直通仅适用于 1 类管理程序，比如 Citrix Xen、VMware ESXi、Kernel Virtual Machine 和 IBM Power。
    - 直通的开销因为 CPU、芯片组、管理程序和操作系统的具体组合而各不相同；一般而言，在最新一代硬件上的开销要小得多。
    - 一种给定的管理程序-操作系统组合仅支持特定的 NVIDIA GPU 卡。
- GPU 虚拟化：

    - 所有主要 GPU 供应商 – NVIDIA GRID、AMD MxGPU 和 Intel Graphics Virtualization Technology –g (GVT -g) – 都支持 GPU 虚拟化。
    - 最新版本在特定的较新 GPU 卡上支持 OpenCL。（没有提供对 PyTorch 的 OpenCL 支持）。
    - 最新版的 GRID 支持用于特定的较新 GPU 卡的 CUDA 和 OpenCL。

### Docker 安装选项

在 Docker 容器或 Kubernetes 集群中运行 PyTorch 有很多优势。PyTorch 源代码包含一个 Dockerfile，后者进而包含 CUDA 支持。Docker Hub 也包含一个预先构建的运行时 Docker 映像。使用 Docker 的优势主要在于，PyTorch 模型可以访问物理 GPU 核心（设备）并在其上运行。此外，NVIDIA 有一个用于 Docker Engine 的名为 _nvidia-docker_ 的工具，该工具可以在 GPU 上运行 PyTorch Docker 映像。

### 云安装选项

最后，PyTorch 提供了多个云安装选项：

- **Google Cloud Platform。** Google 提供了一些机器实例，它们能访问特定区域的 1、4 或 8 个 NVIDIA GPU 设备。同时还可以在容器化的支持 GPU 的 Jupyter Notebook 上运行 PyTorch。
- **Amazon Web Services (AWS)。** Amazon 提供了 AWS Deep Learning Amazon Machine Image (AMI)，以及可选的 NVIDIA GPU 支持，后者能在各种各样的 Amazon Elastic Compute Cloud 实例上运行。我们已预先安装了 PyTorch、Caffe2 和其他深度学习框架。AMI 可以支持多达 64 个 CPU 核心和多达 8 个 NVIDIA GPU (K80)。
- **Microsoft® Azure® Cloud。** 您可以在 Azure 机器实例的 Microsoft 数据科学虚拟机系列上设置 PyTorch，以便仅使用 CPU 或者包含最多 4 个 K80 GPU 。
- **IBM Cloud® 数据科学和数据管理。** IBM 提供了一个包含 Jupyter Notebook 和 Spark 的 Python 环境。可以使用 Python3 内核来安装 PyTorch；首先需要确保 `conda` 是最新的，然后使用以下命令安装 PyTorch：

    `conda install pytorch torchvision -c soumith`

- **IBM Cloud Kubernetes 集群。** IBM Cloud 上的 Kubernetes 集群可以运行 PyTorch Docker 映像。

本文翻译自： [Get started with PyTorch](https://developer.ibm.com/articles/cc-get-started-pytorch/)（2018-01-18）