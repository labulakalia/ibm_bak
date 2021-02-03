# Deeplearning4j 入门
探究此基于 Python 的深度学习库

**标签:** 人工智能

[原文链接](https://developer.ibm.com/zh/articles/cc-get-started-deeplearning4j/)

Vinay Rao

发布: 2018-02-14

* * *

Eclipse Deeplearning4J (DL4J) 是包含深度学习工具和库的框架，专为充分利用 Java™ 虚拟机 (JVM) 而编写。它具有为 Java 和 Scala 语言编写的分布式深度学习库，并且内置集成了 Apache Hadoop 和 Spark。 [Deeplearning4j](https://deeplearning4j.org/) 有助于弥合使用 Python 语言的数据科学家和使用 Java 语言的企业开发人员之间的鸿沟，从而简化了在企业大数据应用程序中部署深度学习的过程。

DL4J 可在分布式 CPU 和图形处理单元 (GPU) 上运行。社区版本和企业版本均已面市。

## Eclipse Deeplearning4J 框架

**Skymind**

Skymind 是总部位于旧金山的人工智能 (AI) 初创企业，由 DL4J 首席开发人员 Adam Gibson 联合他人一起创办。Skymind 销售面向 DL4J 生态系统的商业支持服务和培训服务。此外，该公司的 Skymind Intelligence Layer 平台可填补 Python 应用程序与企业 JVM 之间的空白。

IBM PowerAI 团队已将 DL4J 移植到 PowerAI 上。随后，该团队与 Skymind 协作，将 DL4J 融入 IBM POWER8 架构，包括使用 NVLink 提供 NVIDIA GPU 支持。

DL4J 是由来自旧金山和东京的一群开源贡献者协作开发的。2014 年末，他们将其发布为 Apache 2.0 许可证下的开源框架。主要是作为一种平台来使用，通过这种平台来部署商用深度学习算法。创立于 2014 年的 Skymind 是 DL4J 的商业支持机构。

2017 年 10 月，Skymind 加入了 Eclipse 基金会，并且将 DL4J 贡献给开源 Java Enterprise Edition 库生态系统。有了 Eclipse 基金会的支持，人们能够更加肯定 DL4J 项目必将得到妥善监管，同时也确保为商业开发提供合适的开源许可证。Java AI 开发人员已将 DL4J 视为成熟且安全的框架，因此，这些新建立的伙伴关系将吸引企业在商业领域使用 DL4J。

### 企业大数据应用程序中的深度学习

数据科学家使用 Python 来开发深度学习算法。相比之下，企业大数据应用程序倾向于使用 Java 平台。因此，为填补缺口并在大数据应用程序中部署深度学习，DL4J 的开发人员必须对若干解决方案进行创新。

[Keras](https://keras.io) 应用程序编程接口 (API) 规范的采用，有助于从其他框架（例如，TensorFlow、Caffe、Microsoft® Cognitive Toolkit (CNTK) 和 Theano）导入深度学习模型。Keras API 可通过 JVM 语言（例如，Java、Scala、Clojure 乃至 Kotlin）来访问，从而使深度学习模型可供 Java 开发人员使用。

不熟悉 Keras？阅读教程： [“Keras 入门”](https://www.ibm.com/developerworks/cn/cognitive/library/cc-get-started-keras/index.html)

在 JVM 上运行深度学习高性能计算负载时，将面临着诸多挑战。内存管理和垃圾回收等 Java 功能可能会影响性能，使用较大内存时尤为如此。DL4J 绕过了其中部分限制。

### Deeplearning4j 框架中的库

_Deeplearing4j_ 既指的是框架软件分发版，又指的是框架中的特定库。DL4J 框架包含以下库：

- [Deeplearning4j 库](#deeplearning4j-库)
- [ND4J 库](#nd4j-库)
- [Datavec 库](#datavec-库)
- [libnd4j 库](#libnd4j-库)
- [RL4J 库](#rl4j)
- [Jumpy 库](#jumpy-库)
- [Arbiter 库](#arbiter-库)

#### Deeplearning4j 库

Deeplearning4j 库实际上是神经网络平台。它包含各种工具，用于配置神经网络和构建计算图形。开发人员使用此库来构建由数据管道和 Spark 集成的神经网络模型。

除核心库外，DL4J 库还包含许多其他库，用于实现特定功能：

- **deeplearning4j-core。** deeplearning4j-core 库包含了运行 DL4J 所需的全部功能，例如，用户界面 (UI)。它还具有构建神经网络所需的各种工具和实用程序。
- **deeplearning4j-cuda。** deeplearning4j-cuda 库支持 DL4J 在使用 NVIDIA CUDA® 深度神经网络库 (CuDNN) 的 GPU 上运行。此库支持标准化以及卷积神经网络和递归神经网络。
- **deeplearning4j-graph。** deeplearning4j-graph 库执行图形处理来构建 DeepWalk 中所使用的图形矢量化模型，DeepWalk 是一个无监督学习算法，用于学习图形中每个顶点的矢量表示法。您可以使用这些学到的矢量表示法对图形中的相似数据进行分类、分群或搜索。
- **deeplearning4j-modelimport。** deeplearning4j-modelimport 库从 Keras 导入模型，Keras 又可从 Theano、TensorFlow、Caffe 和 CNTK 导入模型。这是关键的 DL4J 库，支持将模型从其他框架导入 DL4J。
- **deeplearning4j-nlp-parent。** deeplearning4j-nlp-parent 库支持将 DL4J 与外部自然语言处理 (NLP) 插件和工具相集成。此接口遵循非结构化信息管理架构 (UIMA)，后者最初是由 IBM 开发的用于内容分析的开放标准。此库还包含适用于英语、中文、日语和韩语的文本分析。
- **deeplearning4j-nlp。** deeplearning4j-nlp 库是 NLP 工具（如 Word2Vec 和 Doc2Vec）的集合。Word2Vec 是一个用于处理文本的双层神经网络。Word2Vec 对于在“矢量空间”中对相似词语矢量进行分组十分有用。Doc2Vec 是 Word2Vec 的一个扩展，用于学习将标签与词语相关联，而不是将不同词语关联起来。
- **deeplearning4j-nn。** deeplearning4j-nn 库是核心库的精简版本，减少了依赖关系。它使用构建器模式来设置超参数，同时配置多层网络，支持使用设计模式在 Java 中构造神经网络。
- **deeplearning4j-scaleout。** deeplearning4j-scaleout 库是各种库的集合，适用于配备 Amazon Web Services 服务器以及封装 Spark 并行代码，以便在多达 96 核的常规服务器（而不是 Spark）上运行。此库还有助于在 Spark 以及包含 Kafka 和其他视频分析流选项的 Spark 上配置 NLP。
- **deeplearning4j-ui-parent。** deeplearning4j-ui-parent 库实际上是 DL4J UI，包含神经网络训练启发式方法和可视化工具。

#### ND4J 库

N-Dimensional Arrays for Java (ND4J) 是科学计算 `C++` 库，类似于 Python 的 NumPy。它支持 JVM 上运行的多种语言，例如，Java、Scala、Clojure 和 Kotlin。您可以使用 ND4J 来执行线性代数或操作矩阵。ND4J 可与 Hadoop 或 Spark 进行集成并由此实现扩展，同时可在分布式 CPU 和 GPU 上运行。

Java AI 开发人员可以使用 ND4J 在 Java 中定义 _N_ 维数组，这使其能够在 JVM 上执行张量运算。ND4J 使用 JVM 外部的“堆外”内存来存储张量。JVM 仅保存指向此外部内存的指针，Java 程序通过 Java 本机接口 (JNI) 将这些指针传递至 ND4J `C++` 后端代码。此结构配合来自本机代码（例如，基本线性代数子程序 (BLAS) 和 CUDA 库）的张量使用时，可提供更佳的性能。ND4J 与 Spark 集成，并且可使用不同后端在 CPU 或 GPU 上运行。Scala API ND4S 有助于实现这种集成。在稍后部分中讨论 DL4J 如何使用硬件加速时，将再次探讨 ND4J 架构。

#### Datavec 库

DataVec 库的主要功能是将数据格式化为张量。通过这种预处理，神经网络能够访问和使用数据。DataVec 执行抽取、转换和加载 (ETL) 操作，同时还支持通过一系列数据转换连接到各种数据源和输出张量。DL4J 支持众多数据类型，包括图像、逗号分隔值 (CSV)、属性关联文件格式 (ARFF) 和纯文本。DL4J 还支持 Apache Camel 集成。

#### libnd4j 库

Libnd4j 是一个纯 `C++` 库，支持 ND4J 访问属于 BLAS 和 Intel Math Kernel Library 的线性代数函数。它与 JavaCPP 开源库紧密结合运行。JavaCPP 不属于 DL4J 框架项目，但支持此代码的开发人员是相同的。

#### Jumpy 库

Jumpy 是一个 Python 库，支持 NumPy 无需移动数据即可使用 ND4J 库。此库实现了针对 NumPy 和 [Pyjnius](https://pyjnius.readthedocs.io/en/latest/) 的包装器。MLlib 或 PySpark 开发人员可以使用 Jumpy，以便在 JVM 上使用 NumPy 数组。

#### Arbiter 库

DL4J 使用此工具来自动调优神经网络。可使用诸如网格搜索、随机搜索和贝叶斯方法之类的各种方法来优化具有许多超参数的复杂模型，从而提高性能。

## Deeplearning4J 的优势

DL4J 具有众多优势。让我们来了解一下其中的三大优势。

### Python 可与 Java、Scala、Clojure 和 Kotlin 实现互操作性。

Python 为数据科学家所广泛采用，而大数据编程人员则在 Hadoop 和 Spark 上使用 Java 或 Scala 来开展工作。DL4J 填补了之间的鸿沟，开发人员因而能够在 Python 与 JVM 语言（例如，Java、Scala、Clojure 和 Kotlin）之间迁移。

通过使用 [Keras API](https://keras.io/)，DL4J 支持从其他框架（例如，TensorFlow、Caffe、Theano 和 CNTK）迁移深度学习模型。甚至有人建议将 DL4J 作为 Keras 官方贡献的后端之一。

### 分布式处理

DL4J 可在最新分布式计算平台（例如，Hadoop 和 Spark）上运行，并且可使用分布式 CPU 或 GPU 实现加速。通过使用多个 GPU，DL4J 可以实现与 Caffe 相媲美的性能。DL4J 也可以在许多云计算平台（包括 IBM Cloud）上运行。

### 并行处理

DL4J 包含单线程选项和分布式多线程选项。这种减少迭代次数的方法可在集群中并行训练多个神经网络。因此，DL4J 非常适合使用微服务架构来设计应用程序。

## Eclipse DL4J 应用程序

DL4J 具有各种应用程序，从图形处理到 NLP，皆涵盖在内。借助内置数据预处理和矢量化工具，DL4J 能够异常灵活地处理许多不同的数据格式。Keras API 还使其更便于使用来自其他框架的预先训练模型。

典型的 DL4J 应用程序包括：

- 安全性应用程序，例如欺诈检测和网络入侵检测
- 在客户关系管理、广告推送和客户忠诚度及维系方面使用的推荐系统
- 面向物联网和其他流数据的回归和预测分析
- 传统面部识别和图像识别应用程序
- 语音搜索和语音转录应用程序
- 针对硬件或工业应用程序的预防性诊断和异常检测

## 支持 DL4J 的平台

理论上，任何支持 JVM 且运行 Java V1.7 或更高版本的 64 位平台均可支持 DL4J。DL4J 使用 ND4J 来访问受支持的 GPU。ND4J 则依赖于其他软件，例如 CUDA 和 CuDNN。

实际上，商用 DL4J 需要生产级 Java 平台。IBM Open Platform for Apache Hadoop and Apache Spark 已在 PowerAI 上完成了 DL4J 框架认证。Cloudera 也已在其 Enterprise Data Hub CDH5 上完成了 DL4J 认证，Hortonworks 同样也在其 Data Platform\\HDP2.4 上完成了认证。

有关安装 ND4J 和其他必备软件的更多信息，请查看 ND4J 入门文档。

### 从源代码构建 DL4J

在其他框架中，使用预先构建的二进制文件来安装框架有时更为方便。但对于 DL4J，最好是从源代码构建和安装 DL4J，以便确保正确处理多种依赖关系。安装为多步骤的复杂过程，适合生产级安装。Apache Maven V3.3.9 或更高版本可对构建和依赖关系进行管理。对于集成开发环境，可以使用 IntelliJ IDEA 或 Eclipse。

专家建议使用“uber-jar”方法来构建 DL4J 应用程序。此方法支持使用 .rpm 或 .deb 包在“uber-jar”内部分发依赖关系，并将部署与开发隔离开来。

有关使用 DL4J 构建应用程序的更多信息，请阅读 [快速入门指南](https://deeplearning4j.org/docs/latest/deeplearning4j-quickstart)。

### Power 架构上的 DL4J

DL4J 团队官方在 DL4J 存储库中为 IBM Power® 维护 DL4J。Maven 负责管理构建和安装过程。ND4J 具有两个用于 Power 架构的本机平台后端：第一个在 POWER8 CPU 上运行，第二个则使用 NVLink 互连在 NVIDIA GPU 上运行。

ND4J CPU 后端 `nd4j-native-platform` 运行 OpenBLAS 库的优化版本。为实现加速，ND4J CPU 后端在 POWER8 处理器上使用矢量/标量浮点单元。

ND4J GPU 后端 `nd4j-cuda-8.0-platform` 可在使用 NVLink 接口与 POWER8 处理器互连的 NVIDIA GPU 上运行 CUDA 或 cuDNN。

有关在 Power 架构上安装 DL4J 的更多信息，请访问 [https://deeplearning4j.org/](https://deeplearning4j.org/)。

**注：** 在 POWER8 平台上，首选的 Java 级别为 V8。为 POWER8 而调优的唯一 Java 7 发行版为 Java 7.1 SR1 或更高版本。

## DL4J 如何使用硬件加速？

DL4J 分别依靠 ND4J 和 ND4S 这两种特定于平台的后端来使用硬件加速。ND4J 和 ND4S Java 与 Scala API 用于封装 BLAS 库，例如，Jblas、Netlib-blas 和 Jcublas。

ND4J 具有两种级别的运算。高级运算包括卷积、快速傅立叶变换、各种损失函数、变换（例如，sigmoid 变换或 tanh 变换）和约简。通过 ND4J API 调用时，BLAS 会实施低级运算。BLAS 运算包括矢量加法、标量乘法、点乘、线性组合和矩阵乘法。

ND4J 在特定于平台架构的后端上运行。CPU 后端带有 _nd4j-native_ 前缀，GPU 后端带有 _nd4j-cuda-_ 前缀，其中 CUDA 版本可能是 7.5 或 8.0 等。无论使用何种后端，ND4J API 都保持不变。

最后，ND4J 实现特定于 BLAS 的数据缓冲区，用于存储 BLAS 处理的数组和原始数据字节。根据后端，此存储抽象层具有不同的实现。JVM 与后端之间通过 JNI 进行通信。

## 结束语

Deeplearning4j 是深度学习工具和库框架，专为充分利用 Java 虚拟机而编写。DL4J 框架旨在用于在生产级服务器上部署商用深度学习算法。它具有为 Java 和 Scala 而编写的分布式深度学习库，并且内置了与 Hadoop 和 Spark 的集成。DL4J 可在分布式 CPU 和 GPU 上运行，提供了社区版本和企业版本。Skymind 作为其商业支持机构，为企业版本提供专业的服务和支持。

本文翻译自： [Get started with Deeplearning4j](https://developer.ibm.com/articles/cc-get-started-deeplearning4j/)（2017-12-18）