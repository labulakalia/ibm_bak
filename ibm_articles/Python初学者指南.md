# Python 初学者指南
花 5 分钟快速掌握这门科学计算语言

**标签:** Python

[原文链接](https://developer.ibm.com/zh/articles/os-beginners-guide-python/)

Steven Frankel

发布: 2018-04-17

* * *

您是否正在寻找一种容易学习的编程语言来帮助完成科学工作？不必舍近求远，Python 就能办到。我们将介绍开始使用这种简单的编程语言需要了解的一些基本概念，并展示如何使用 Python 执行从运行代数计算到根据数据生成图形输出的所有操作。

## 科学计算回顾

科学计算涉及到使用计算机解决科学问题。具体来讲，可以使用计算机来解方程。从单一非线性方程（ _求根_ ）到线性代数方程组（ _数字线性代数_ ），再到解决非线性偏微分方程组（ _计算物理学_ ）。

过去，解决这些问题的数字算法是用诸如 C/C++ 和 Fortran 之类的语言编写的 — 现在仍是如此。那么，Python 适用于哪些地方呢？Python 非常适合用来快速实现和测试新（或旧）算法，并将多个物理学代码编组到一起（美国的顶级实验室经常这么做）。Python 易于使用，学起来很有趣，而且非常强大。那么您还在等什么？让我们开始吧！

## 下载 Python

Python 可广泛应用在所有运行 Linux 或 macOS 操作系统的计算机上。您甚至可以在 iPad 上使用 [Pythonista](http://omz-software.com/pythonista/) 应用程序来运行 Python。您还可以从 [Python.org](http://www.python.org/) 下载一个适合 Windows 的版本。但是如果您要执行科学计算（即使您不执行科学计算），建议您下载并安装 [Anaconda](https://www.anaconda.com/) 。Anaconda 提供了一个完整的 Python 安装，还提供了许多用于科学计算的不错的包（我将它们称为 _模块_ ）。它还能够轻松访问集成开发环境 Spyder。

## Python 乐意为您效劳

安装 Anaconda 后，您可以单击 Anaconda Navigator 的图标，并开始享受乐趣。在右下角的窗口中，您会看到一个命令提示符。将鼠标放在此提示符右侧，就可以开始输入 Python 命令。如果采用学习新编程语言的传统路线，可以先输入 `print("Hello World!")` ，然后按下 **Return** 。

可以使用命令提示符输入一个或多个命令来快速测试代码段，或者为您的工作生成输出。在涉及大量代码时，最好单独生成并保存一个程序文件（稍后会详细介绍）。

另一个选项（至少在 Linux 和 macOS 上）是打开终端窗口，并在命令提示符上输入 `Python` 。这样做会启动 Python 命令提示符，而且您可以开始输入命令并运行 Python 代码。如果在终端窗口中输入 `idle` ，您会获得一个包含 Idle Python 编辑器的新窗口 — 可以方便地编写新 Python 脚本并使用强大的 `F5` 命令运行它们。

## 这一名称的含义是什么？

您已安装 Python 并知道如何开始输入命令，现在，您可以继续执行一些数学和科学操作。解方程的计算机编程涉及到使用变量和处理代表这些变量的数字。例如，在 Python 中，可以通过在命令提示符上输入以下命令来定义一个变量：

```
>>> x0 = 1.5
>>> x1 = 2.0

```

Show moreShow more icon

恭喜！您刚才同时创建了两个名为 `x0` 和 `x1` 的新变量，而且分别为它们分配了值 `1.5` 和 `2.0` 。要查看实际效果，只需输入以下内容：

```
>>> x0,x1

```

Show moreShow more icon

或者可以调用 `print` 函数：

```
>>> print (x0,x1)

```

Show moreShow more icon

您不需要将这些变量声明为实数（浮点数）或整数，因为 Python 是一种动态输入语言；它根据赋给变量的值来动态确定变量类型。

## 计算机和代数

您为两个数字分配了两个变量，所以您可以立即对它们执行简单的计算机代数运算；您可以按您的意愿对它们执行加减乘除。这是计算机最在行的事情。要查看此代数运算的运行效果，可在命令提示符上输入以下命令：

```
>>> yp = x0 + x1
>>> ym = x1 - x0
>>> yt = x0*x1
>>> yd = x1/x0
>>> print(yp,ym,yt,yd)

```

Show moreShow more icon

您现在已在正式执行科学计算。

## 计算机和逻辑

如果计算机只能执行代数运算，那么它们对科学计算的影响将是有限的。事实上，计算机还擅长处理逻辑，这使得创建复杂程序成为可能。您可能很熟悉 _if this, then that_ ( [IFTTT](https://ifttt.com/)) 逻辑。这与我讲的主题并不完全相符，但很接近。我指的是程序流控制，或者仅在某些条件下执行一行或一段（多行）代码，在其他条件下执行其他行或其他段代码的能力。要了解此话的真正含义，请输入以下命令：

```
>>> x1 = 2.0
>>> if x1 < 0:
        x = 0
        print("Negative switched to zero")
    elif x1 == 0:
        print ("Zero")
    else:
        print ("x1 is positive")

```

Show moreShow more icon

该代码是一个 `if` 代码段的示例，其中 `elif` 是 _else if_ 的缩写，如果前面两段（或您需要的多段）代码测试失败，则会执行 `else` 。 要获得更详细的解释，请查阅 Python 文档中的 [更多控制流工具](https://docs.python.org/2/tutorial/controlflow.html) 。

科学计算中的许多算法背后的强大功能，都与对不同数据多次执行同一段代码的能力相关。在这里，循环就派上了用场。考虑下面这段代码，它使用 Python 内置函数 `range` 生成一个包含 10 个整数的列表，从 0 开始：

```
>>> x0 = 1.5
>>> for i in range(10):
        x1 = x0 + i

```

Show moreShow more icon

此代码执行计算 x1 = x0 + i 10 次，从 i=0 开始，到 i=9 结束。

## 您的函数是什么？

有了函数，就可以执行将大型编程任务分解为少量子任务或 _函数_ 的重要编程过程。 Python 拥有内置函数和外部库，我稍后会解释它们。您也可以构建自己的函数。可以使用 Python 关键字 `def` 创建函数，如下所示，名为 `f` 的函数接收输入变量 `x` ，并返回对所编写的代数表达式执行计算后获得的值：

```
>>> def f(x):
        return x**3 - 5*x**2 + 5*x + 1
>>> x0 = 1
>>> print ("f(x0) = ", f(x0))

```

Show moreShow more icon

要创建一个函数来计算上面的函数的分析衍生版本，请输入以下内容：

```
>>> def fp(x):
        return 3*x**2 - 10*x + 5
>>> x1 = 2
>>> print ("The derivative of f(x) is: ", fp(2))

```

Show moreShow more icon

## 将此代码存档

目前，您已在命令提示符上输入了所有 Python 命令，这适合简短的、可任意使用的代码。但是，如果您在开发一个大型项目，或者您想保存您的工作供以后使用，建议创建一个新 Python 文件或 _脚本_ 。可以在终端窗口中使用您最喜欢的文本编辑器完成此任务。例如，使用 vi，可以输入 `vi newton.py` 来创建名为 _newton_ 且以 .py 作为文件扩展名的空白文本文件，让每个人（以及具体来讲计算机）都知道这是一个 Python 文件。然后，保持文件处于打开状态，您可以开始输入您的 Python 命令。

**备注：** Python 使用空格来指示代码段。标准约定是使用 4 个空格来缩进一个新代码段，比如组成一个函数的代码或一个 `if-then` 代码块。

编写程序的另一个重要方面是 _注解_ — 这些代码行会告诉阅读您的文件的人该脚本的用途。单行注解以井号 (`#`) 开头；要创建多行注解，可在它们前面加一个反斜杠后跟 # (`\#`)，并在它们后面加上 `#\` 。输入您的代码后，保存它并退出编辑器。然后从终端窗口命令行输入 `python newton.py` 来运行该代码（假设您与所保存的文件在同一目录中）。

在科学计算中，将一个问题或任务分解为更小的问题或 _子任务_ 通常是个不错的主意。在 Python 中，这些子任务称为 _模块_ 。模块是包含定义和语句的附加 Python 文件（具有文件扩展名 .py）。也有一些预先构建的模块。您可以在自己的程序中使用任何模块，只需使用 `import` 关键字导入它即可。例如，math 模块包含基本数学函数，比如正弦和余弦函数；可通过关键字 `import math` 来使用它们。

## 将科学计算能力导入 Python 中

您希望在 Python 中使用的两个强大的科学计算模块是 [NumPy](http://www.numpy.org/) 和 [SciPy](http://www.scipy.org/) 。NumPy 包含许多强大的特性，但在这里尤为重要的是，它能够创建分配给单个变量的相同数据类型的数字集合（称为 _数组_ ）。它还包含广泛的线性代数、傅里叶变换和随机数功能。SciPy 是一个包罗万象的科学计算生态系统，它包含 NumPy 和其他大量包，比如 matplotlib，我们将在下一节中介绍它。下面的代码提供了如何导入 NumPy 并在一个代码段中使用它的简单示例：

```
import numpy as np
/* Now we create and assign a list of integers 0 through 9 to x[0] through x[9], effectively creating a one-dimensional NumPy array */

x = np.linspace(10)
print(x)

```

Show moreShow more icon

## 使用 Python 生成图形输出

利用科学计算生成数据的有效图形输出，是理解和传达结果的关键。Python 中实现这个重要目标的标准包是 [matplotlib](https://matplotlib.org/) 模块。可以看到，该模块很容易访问和使用：

```
>>> import numpy as np
>>> import matplotlib.pyplot as plt
>>> x = np.arange(0, 5, 0.1)
>>> y = np.sin(x)
>>> plt.plot(x,y)
>>> plt.show()

```

Show moreShow more icon

matplotlib 模块提供了一些命令来控制行类型、颜色和样式，并保存您的图表。

## 后续行动

网络上有数百个网站可帮助您了解 Python 和它在科学计算中的作用。您还可以通过以下这些优秀资源进一步学习：

- 我最喜欢的资源之一是 [用于科学计算的 Python 的简介](https://engineering.ucsb.edu/~shell/che210d/python.pdf) 。
- 正在寻找优秀的图书？请查阅 Anthony Scopatz 编写的 [Effective Computation in Physics: Field Guide to Research with Python](https://www.amazon.com/Effective-Computation-Physics-Research-Python/dp/1491901535) 。
- 最近的一篇实用性的参考文章是 Jake VanderPlas 编写的 [A Whirlwind Tour of Python](https://www.oreilly.com/learning/a-whirlwind-tour-of-python) ，来自他的图书 [Python Data Science Handbook: Essential Tools for Working with Data](https://www.amazon.com/Python-Data-Science-Handbook-Essential/dp/1491912057) 。
- 在您取得更多进步或希望超越 Python 中的科学计算时，推荐阅读 Dan Bader 编写的 [Python Tricks: A Buffet of Awesome Python Features](https://www.amazon.com/Python-Tricks-Buffet-Awesome-Features/dp/1775093301) 。

编码愉快！

本文翻译自： [Beginner’s guide to Python](https://developer.ibm.com/articles/os-beginners-guide-python/)（2018-02-27）