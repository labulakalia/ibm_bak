# 通过为 Python 提速来助力科学研究
优化 Python 代码以用于研究

**标签:** Python

[原文链接](https://developer.ibm.com/zh/articles/ba-accelerate-python/)

Vinay Rao

发布: 2018-05-16

* * *

科学家通过将复杂科学问题分解为更简单的数据库来加以解决。然后，对每个区块应用最佳工具来获得解决方案。解决每组问题的最佳方法可能有所不同，因此每个问题可能都需要不同的工具。

## 从这里开始

我建议阅读本人之前的博客文章” [为何应使用 Python 进行科学研究](https://developer.ibm.com/dwblog/2018/use-python-for-scientific-research)”，我在文中创建了使用 Python 进行科学研究的案例，不是用于替代 Fortran 或 C/C++，而是作为包装器或各种专用编程模块和加速硬件之间的桥梁。

Python 非常适用于数据科学、机器学习和深度学习，所有这些都作为解决科学问题的工具，正日渐普及。Python 的优势在于它整合了多种解决问题的方法。通过将所有解决问题的工具都整合到一个容器中，Python 是一个非常棒的工具包。

相比之下，Fortran 在数字运算方面拥有最出色的性能。科学家使用 Fortran 为科学研究编写数值算法程序。Python 则为科学家提供了一种包装专用程序的有效方法，由此能够通过通用应用层轻松访问这些程序。此外，Python 还可以使用合适的后端，比如 CPU 后端、GP-GPU 后端或量子处理后端，加快这些工具的速度（例如，整合深度学习和量子计算）。

在本文中，我们将讨论如何优化 Python，助力科学研究。具体来说，我们将介绍：

- 使用 Python 时遇到的障碍以及克服障碍的技巧
- Python 并发性和可扩展性
- Python 编译器
- 加快运行 Python 的硬件速度
- Python 的主要实现

## 克服进入科学研究领域时遇到的种种障碍

Python 很受开发人员的欢迎。民意调查将其列入 10 大编程语言，如这篇 IEEE Spectrum [报告](https://spectrum.ieee.org/computing/software/the-2017-top-programming-languages) 所示。尽管很受欢迎，但 Python 在核心科学与技术领域并没有得到广泛使用。原因是担心 Python 的性能，这无可厚非，尤其是科学家们。

他们继续依赖于 Fortran、C 或 C++，而不使用 Python。一些科学家探索类似 Julia 等从头开始设计的新语言，旨在克服 Python 的缺陷，同时保留 Fortran 的性能。请记住，我们并不是要用 Python 取代 Fortran 或 C/C++，这有助于缩小讨论的范围。因此，在这一点上，需要正确看待问题，将我们的注意力集中在阻止 Python 进入科学研究领域的真正障碍以及如何克服这些障碍上。

### 性能

纯粹原生的 Python 程序运行非常缓慢。使用 C/C++ 或 Fortran 进行混合编程可以加快性能。在”Python 可扩展性”部分，我会探讨一些混合编程的方法。

诸如 NumPy 和 SciPy 等程序包，通过使用特定平台、经过优化的数学库，加快线性代数和矩阵运算，这些数学库包括使用了硬件加速的 Linear Algebra Package (LAPACK) 和 Basic Linear Algebra Subprograms (BLAS) 等。有关这方面的更多信息，请参阅”Python 硬件加速”部分。

### 并发性

Python 在使用线程时，为了性能会牺牲简单性和稳健性。Python 线程是实际操作系统线程的简单抽象。CPython 拥有 [全局解释器锁 (GIL)](https://wiki.python.org/moin/GlobalInterpreterLock) ，将执行限制为一次一个线程。这种互斥锁也会限制 Python 利用多核 CPU 的能力。所有基于 C 语言的 Python 扩展都受制于这种限制。

一些编译器可以围绕一段代码显式释放 GIL。因此，扩展可以在将控件转交给外部库时释放 GIL，并在控件返回至 Python 代码时重新获取。在”Python 并发性”部分，我会再次讨论这个概念。

### 碎片化

Python2 不仅会在 Python3 的首次发布之后继续存在十年，而且仍比 Python3 的速度要快。Python 3.6 是 Python3 目前最快的版本，早期报告显示，Python 3.7 将比 Python 3.6 更快，但是仍比不上 Python 2.7，该版本预计是 Python2 的最后一个版本。另一个问题在于 PyPy，它是最有前景的 CPython 性能提升替代项之一，目前尚不支持 Python3。

Python3 从很多方面提升了性能，包括 GIL 变更（自 Python 3.2 以上版本）、新的 [concurrent.futures](https://docs.python.org/3/library/concurrent.futures.html) 模块以及众多语言增强功能。所有这些提升以及 Python2 即将进入维护方式这一事实，促使 Python3 被选为继续向前发展的 Python 版本。

### 设计理念

Python 的设计将生产力置于原始性能之前，这使得 Python 备受开发人员追捧。Python 实现者已经克服或规避了 Python 中的许多性能缺陷，但人们仍对此存有各种担忧，这主要是因为缺乏介绍这些解决方案的良好文档。在某种程度上，这是因为一些话题非常深奥，探讨的是 Python 内部原理而非 Python 应用开发。

本文概述 Python 优化的一些重要方面，介绍如何加快 Python 性能，使其适用于科学研究计算。

## Python 并发性

按照设计，Python 将并发性问题留给了操作系统，而只为操作系统机制提供了一个简单的包装器。Python 字节码直接在单核处理器上运行。一次只执行一个线程。通过获取 GIL，执行中的线程会将其他线程锁在外面。这种简单稳健的设计也使其更容易编写扩展。

当部分线程为 I/O 密集型时，多线程性能保持良好，但是当所有线程为 CPU 密集型时，性能会变差，这对于科学计算来说是个大问题，因为大部分任务都是 CPU 密集型任务。摆脱 GIL 看似微不足道，但做起来并不简单。到目前为止，从 CPython 移除 GIL 的尝试还不是很成功。大部分尝试不是会影响单线程操作模式的性能，就是会破坏与扩展的可兼容性，这两种情况都不是理想的结果。

一种替代方法就是生成进程而不是线程。Python 利用多处理模块实现这项功能。每个进程都有自己的 GIL，因此互不影响。在 Python3 中，这项功能包含在新的 concurrent.futures 模块中。这种方法会有一些开销，但是这样 Python 可以使用多核 CPU 和 CPU 集群。

正如我们前面所见，一些编译器可以围绕一段代码显式释放 GIL。因此，扩展可以在将控件转交给外部库时释放 GIL，并在控件返回至 Python 代码时重新获取。这项功能允许 Python 切换到 C 代码，这种代码可以处理多线程或多处理任务。Python 中的科学编程包（比如 NumPy 和 SciPy）就使用这种方法。

## Python 可扩展性

Python 具有高可扩展性，存在许多使用 C 语言或 Fortran 编写扩展的方法。必要时，Python 代码可以直接将这些扩展作为子例程来调用。这部分讨论用于构建扩展的一些主要编译器（绝对不是完整列表）。

### Cython

[Cython](http://cython.org/) （不同于 CPython）既是指一种语言，也是指一种编译器。Cython 语言是添加了 C 语言语法的 Python 语言的超集。Cython 可以在代码段或完整函数中显式释放 GIL。变量和类属性上的 C 类型声明以及对 C 函数的调用都使用 C 语法。其余部分代码则使用 Python 语法。通过这个混合的 Cython 代码，Cython 编译器可生成高效的 C 代码。任何定期优化的 C/C++ 编译器都可以编译此 C 代码，从而高度优化扩展的运行时代码，性能接近于原生的 C 代码性能。

### Numba

[Numba](https://numba.pydata.org/) 是一个动态、即时 (JIT) 且可感知 NumPy 的 Python 编译器。Numba 使用 LLVM 编译器基础架构，生成优化的机器代码和从 Python 调用代码的包装器。与 Cython 不同，编码使用常规的 Python 语言。Numba 可读取来自装饰器中所嵌入注释的类型信息，并优化代码。对于使用 NumPy 数据结构的程序，比如数组以及许多数学函数，它可以实现与 C 或 Fortran 语言类似的性能。NumPy 对线性代数和矩阵函数使用硬件加速，利用 LAPACK 和 BLAS 提供额外加速，大大提升了性能，参见 IBM 博客文章 [C、Julia、Python、Numba 和 Cython 在 LU 因式分解方面的速度比较](https://www.ibm.com/developerworks/community/blogs/jfp/entry/A_Comparison_Of_C_Julia_Python_Numba_Cython_Scipy_and_BLAS_on_LU_Factorization?lang=en) 。

除 CPU 以外，Numba 还能够使用 GP-GPU 后端。Anaconda, Inc. 是 Python 某个主要发行版的幕后公司，该公司还开发了 Numba 和商业版的 Numba Pro。

### Fortran to Python Interface Generator

Fortran to Python Interface Generator (F2Py) 起初为一个独立的程序包，现在包含在 NumPy 中。F2Py 支持 Python 调用以 Fortran 编写的数值例程，就好像它们是另一个 Python 模块一样。因为 Python 解释器无法理解 Fortran 源代码，所以 F2Py 以动态库文件格式将 Fortran 编译为本机代码，这是一种共享对象，包含具有 Python 模块接口的函数。因此，Python 可以直接将这些函数作为子例程来调用，以原生 Fortran 代码的速度和性能来执行。

## Python 编译器

[CPython](https://github.com/python/cpython) 是标准的内置 Python 编译器。与前一部分讨论过的编译器不同，CPython 是真正使用 C 语言实现的 Python 解释器。在交互模式下，CPython 就像任何经典的解释器一样，逐行执行指令。但是在脱机模式下，CPython 以透明方式将 Python 代码转换为字节码，位于内存或磁盘上的缓存文件夹中。

Python 字节码是 Python 源代码的低级别、与平台无关的转换。在这方面，字节码不同于编译某种语言程序所产生的机器代码，比如 C 或 Fortran 程序。Python 字节码在 Python 虚拟机 (PVM) 上执行，这种虚拟机是 Python 运行时引擎，也是 CPython 的组成部分。这种 PVM 只是一个解释器循环，用于解释和执行字节码。因此，Python 程序的运行速度不如其他程序，比如 C 程序。

这个事实是引发人们对 Python 产生太多困惑和误解的根源。标准的 Python 库的确比 C 语言要慢，但是 Python 包可以利用 Fortran、C/C++ 和硬件加速功能来增强性能。以上所有讨论均关于 Python 的标准实现，但是还有许多 CPython 的替代实现，Python 的某些发行版可能会重新包装，甚至是替换 CPython。

以下是一系列主要的 CPython 替代方案：

- [IronPython](http://ironpython.net/) ：这个 Python 实现与 [Microsoft .NET Framework](https://docs.microsoft.com/en-us/dotnet/) 集成，目的是在 Microsoft [公共语言运行时](https://docs.microsoft.com/en-us/dotnet/standard/clr) (CLR) 上运行。这是少数不实施 GIL 的 Python 实现之一。
- [PyPy](https://pypy.org/) ：这是 Python 的新款 JIT 编译器。PyPy 的运行速度比 CPython 更快，使用的内存更少，但是它也实施了 GIL，因此存在的限制与 CPython 相同。 [PyPy 的实验分支](http://doc.pypy.org/en/latest/stm.html) 实现了去掉 GIL 的 [软件事务内存](http://conference.scipy.org/proceedings/scipy2017/pdfs/dillon_niederhut.pdf) (STM)。
- [Jython](http://www.jython.org/) ：该编译器在 Java™ 虚拟机 (JVM) 上运行 Python，是没有实施 GIL 的另一个 Python 实现。Jython 真正实现了多线程，使用 [Java Foreign Function Interface](https://github.com/jnr/jffi) (JFFI) 支持扩展。支持 C Foreign Function Interface (CFFI)，C 扩展仍是路线图中的项目。
- [Stackless Python](https://github.com/stackless-dev/stackless/wiki) ：该 CPython 增强版本支持微线程。据实施者所述，微线程或”threadlets”的运行速度比操作系统线程至少快一个数量级，扩展性能也更好。和普遍观点不同的是，Stackless Python 没有实施 GIL，因此无法利用操作系统或硬件处理器的任何多核特性。标准 Python 通过使用 [greenlets](https://greenlet.readthedocs.io/en/latest/) 模块也能获得多线程优势，后者是这种开发分支的衍生产品。

## 硬件加速和 Python

Python 可以在许多不同的后端上运行，从单个 CPU 到量子计算机皆在范畴之内。因为 Python 是高级语言，所以硬件加速需要借助外部硬件驱动程序和库。在这个部分，我会概述一些广泛使用的硬件后端，介绍 Python 如何使用这些硬件后端以及量子计算等一些新后端。

### 大部分机器上的 CPU

大部分机器要么是本地桌面，要么是基于云的远程服务器，全都采用多核 CPU：

- **单核 CPU** ：Python 字节码直接在单核处理器上运行。一次只执行一个线程。通过获取 GIL，执行中的线程会将其他线程锁在外面。这是运行 Python 的基础环境。
- **多核 CPU** ：处理模块或 concurrent.futures 模块是 Python 利用多核 CPU、实现适度性能提升的首选方法。

### CPU 集群/协处理器

可以通过多种方法来利用超级计算机上的 CPU 集群或通常出现在高性能计算工作站中的 CPU 协处理器。最受欢迎的方法就是实现消息传递接口 (MPI)，比如 OpenMPI。针对 MPI 的 Python 绑定位于 MPI4Py 模块中。MPI 在分布式内存机器上最有效。当使用 MPI4Py 时，在 Python 中实现并行处理即成为可能。

### 图形芯片组上的 GP-GPU

根据使用的是 NVIDIA 还是其他 GPU 芯片组，可以使用以下两个模块之一：

- NVIDIA [PyCUDA](https://developer.nvidia.com/pycuda) ：该模块将 NVIDIA CUDA 映射到 Python 上，因此 Python 可以利用 NVIDIA GPU 芯片组上的 GP-GPU 编程。
- [PyOpenCL](https://pypi.python.org/pypi/pyopencl) ：该模块允许 Python 访问 [OpenCL](https://www.khronos.org/opencl/) API，支持 Python 使用 AMD 和 Intel 等 GPU 芯片组供应商提供的 GP-GPU 后端。

### Google Tensor Processing 后端

目前，Python 访问 [Tensor Processing Unit](https://cloud.google.com/blog/big-data/2017/05/an-in-depth-look-at-googles-first-tensor-processing-unit-tpu) (TPU) 后端的唯一方法就是通过使用 [TensorFlow](https://developer.ibm.com/articles/cc-get-started-tensorflow/) 框架。这么做的好处在于可以将张量（n 维数组）与 TensorFlow Python API 结合使用。TPU 是矩阵处理器而不是向量处理器，每个循环可能执行数十万条指令，而无需访问内存。

## 量子计算后端

诸如 IBM [量子处理器芯片](https://spectrum.ieee.org/tech-talk/computing/hardware/ibm-edges-closer-to-quantum-supremacy-with-50qubit-processor) 等量子计算平台是相对较新的平台，支持科学研究所需的高端科学计算。理论上，传统超级计算机可能需要数百万年时间才能解决的某些种类的问题，量子计算机只需数小时即可解决。复杂分子建模是量子计算机擅长的另一个应用领域。

IBM 将这种强大的新资源作为云后端提供，包含在 IBM Q™ Experience 中，用于开展学术和科学研究。 [IBM Q Experience](https://quantumexperience.ng.bluemix.net/qx/experience) 支持研究人员和学术机构通过 IBM Cloud™ 服务访问 5 量子位和 16 量子位的 IBM 量子计算机。 [IBM Q](https://www.research.ibm.com/ibm-q/) ™ 是 20 位 IBM 量子计算机的商业实现，也是一种 IBM Cloud 服务。针对在线后端实现 Python API 客户端的 SDK 支持访问系统。

IBMQuantumExperience Python 包是官方 API 客户端，用于在 Python 中使用 IBM Q Experience。

## Python 的主要发行版

现在，您可能已经意识到优化或加速 Python 的确并不简单，最好留给专家来处理。为此，不要从标准发行版开始，避免投入大量时间进行优化，最好是一开始就使用专为科学研究定制的 Python 主要发行版之一，以下列出了这些发行版：

- [ActiveState ActivePython](http://www.activestate.com/activepython/) ：这个 Python 商用发行版包含了大部分的主要科学计算模块。同时也提供了免费的社区版本。
- [Anaconda Python](https://www.anaconda.com/) ：这是最受欢迎的 Python 商用发行版之一。它包含易用的图形开发环境管理器，能够管理来自不同开发渠道的多个 Python 版本。包括 Intel、Microsoft 和 IBM 在内的大部分主要供应商，已经将其产品集成到 Anaconda 发行版中。同时也提供了免费的社区版本。
- Enthought Canopy：这个商用发行版也提供免费精简版本，有超过 450 个科学计算模块可供使用。

## Power 架构上的 Python

IBM 与 Anaconda, Inc. 联手，将 Anaconda Python 发行版带入 IBM® POWER® 平台，比如 IBM POWER8® 和 POWER9™。该发行版集成了用于机器学习和深度学习的 IBM PowerAI 软件发行版。它也可以使用与该平台集成的 NVIDIA Tesla Pascal P100 GPU 加速器的 NVIDIA NVLink 高速接口，提升深度学习和分析应用的性能。

## 后续步骤

- [Python 编程初学者指南](https://developer.ibm.com/articles/os-beginners-guide-python/)

本文翻译自： [Accelerating Python for scientific research](https://developer.ibm.com/articles/ba-accelerate-python/)（2018-05-16）