# 深度学习框架 – IBM AIX 上支持 ONNX 的 Pytorch 和 Caffe2
通过样本用例了解如何在 AIX 上安装 Pytorch 或 Caffe2

**标签:** IBM AIX,人工智能,数据科学,系统

[原文链接](https://developer.ibm.com/zh/articles/pytorchcaffe2-with-onnx-on-aix/)

Kavana N Bhat, Shajith Chandran

发布: 2019-08-28

* * *

## 简介

人工智能 (AI) 正在改变各个行业及业务职能，从而导致人们对 AI、其子领域以及相关领域（例如，机器学习和数据科学）的兴趣不断提升。IBM® AIX® 根据其客户不断变化的需求持续进行调整。为此，我们提供适用于 AIX 的 PyTorch 和 Caffe2 开源机器学习库源代码以供试验。

PyTorch 从 Caffe2 和 ONNX 中提取出面向生产的模块化功能，并将其与 PyTorch 现有的、注重研究的灵活设计结合起来，为各种 AI 项目提供快速无缝的路径（从确定研究原型到生产部署）。目前已提供具有所有这些功能的稳定版本 PyTorch 1.0.1，以供在 AIX 上进行试验。

本教程探讨了如何在 AIX 7.2 上构建或安装 PyTorch 和 Caffe2，并将其用于不同的机器学习 (ML) 和深度学习 (DL) 用例。此外，还探讨了将小尾数 (LE) 格式的可用 ONNX 模型转换为大尾数 (BE) 格式以在 AIX 系统上运行的方法。

## 前提条件

[适用于 Linux 应用程序的 AIX 工具箱](https://www.ibm.com/support/pages/aix-toolbox-linux-applications-overview) 包含专为 IBM AIX 系统构建的开源和 GNU 软件集合。所需的必备软件可使用 [“`yum`”安装工具](https://developer.ibm.com/articles/configure-yum-on-aix/) 来安装。 任何从属 Python 程序包均可使用 pip 命令来安装。应先安装 Python3、gcc 和 pip 程序包，然后才能构建 Protobuf、ONNX、PyTorch 或 Caffe2。请参阅 [在 IBM AIX 上配置 YUM 和创建本地存储库](https://developer.ibm.com/articles/configure-yum-on-aix/)，以获取更多相关信息。

## 使用 ONNX 安装 PyTorch 和 Caffe2

执行以下步骤以使用 ONNX 安装 PyTorch 或 Caffe2：

1. 设置以下编译器或链接器环境变量以采用 64 位模式进行构建：





    ```
    #export PATH=$PATH:/opt/freeware/bin
    #export CXXFLAGS="-mvsx -maix64"
    #export CFLAGS="-mvsx -maix64"
    #export CXX='g++ -L/opt/freeware/lib/pthread/ppc64 -lstdc++ -pthread'
    #export LDFLAGS='-latomic -lpthread -Wl,-bbigtoc'
    #export CC="gcc -L/opt/freeware/lib/pthread/ppc64 -pthread"
    #export OBJECT_MODE=64

    ```





    Show moreShow more icon

2. 构建协议缓冲区：





    ```
    #yum install binutils libtool autoconf automake

    #git clone https://github.com/aixoss/protobuf.git -baix_build (This is the Protobuf version 3.6.0 that is being used with Pytorch 1.0.1 with fixes to build on AIX )

    #cd protobuf; git submodule update --init –recursive

    #./autogen.sh

    #./configure

    #make install

    ```





    Show moreShow more icon

3. 构建 ONNX：





    ```
    #pip install numpy

    #git clone https://github.com/aixoss/onnx.git -baix_build (ONNX 1.3.0 with AIX changes)

    #cd onnx; git submodule update --init –recursive

    #python setup.py install

    ```





    Show moreShow more icon

4. 构建 Pytorch/Caffe2：





    ```
    #bash

    #git clone https://github.com/aixoss/pytorch.git -bv1.0.1

    #yum install lapack lapack-devel openblas

    #pip install setuptools six future PyYAML numpy protobuf

    #cd pytorch; source ./aix_setup.sh

    #USE_DISTRIBUTED=OFF python setup.py install

    ```





    Show moreShow more icon


## Caffe2 中的样本 RNN 用例

PyTorch 存储库中的 caffe2/examples/char\_rnn.py 中提供了一个样本 Caffe2 递归神经网络，其中网络能够长期学习并保留记忆，同时展现逐步改进的能力。此脚本不仅可以 _学习_ 英语、语法和拼写，还可以分辨任何给定文本源中所使用的结构和文句上的细微差别。

您需要输入以下命令才能运行此用例：

```
#export LIBPATH=/opt/freeware/lib64/python3.7/site-packages/torch/lib:/opt/freeware/lib/pthread/ppc64

#cd caffe2/python/examples

#wget https://caffe2.ai/static/datasets/shakespeare.txt

#python char_rnn.py --train_data shakespeare.txt

```

Show moreShow more icon

**样本输出**

```
# python char_rnn.py --train_data shakespeare.txt
WARNING:root:This caffe2 python run does not have GPU support.Will run in CPU only mode.
Input has 62 characters.Total input size: 99993
DEBUG:char_rnn:Start training
DEBUG:char_rnn:Training model
Characters Per Second: 2917
Iterations Per Second: 116
---------- Iteration 500 ----------
Khest h
the' the pins wind the thot no ghemgrtnt enlsswe hlear now
An monbamthe, and Othy anschake dod ind hot uchene woe peerC ht unl wind I'thee wame bybun et ane toXy suINend ans thesd tho
:sy unC piFe mnve al sak ufe Qhe one th moI
dralmepend, py thef the wit, th wine there

I mandinge

XusX mwave tothent,
the sth wee bor ndse; f forlHfblr wule wits heand noun bee,
And hn the the to: ml: wy thy to wotd w, woTly mnoll therft wot the.
FSRe goft.
Oine the wine whepyonq

M.bachy,
?nd male nout ma

DEBUG:char_rnn:Loss since last report: 70.67837668657303
DEBUG:char_rnn:Smooth loss: 90.0106901566045
Characters Per Second: 2768
Iterations Per Second: 110

```

Show moreShow more icon

## 使用 ONNX 将模型从 PyTorch 传输到 Caffe2 的用例

此示例演示了如何使用 PyTorch 和 Caffe2 通过分析来预测信用卡违约情况。它使用 PyTorch 为信用卡违约情况生成预测模型，将此模型保存到 ONNX 中，并使用 Caffe2 来加载已保存的 ONNX 模型进行在线评分。

```
#export LIBPATH=/opt/freeware/lib64/python3.7/site-packages/torch/lib:/opt/freeware/lib/pthread/ppc64

#git clone https://github.com/aixoss/onnx-example.git

#cd onnx-example/pytorch

#python Train.py------> Uses the credit card default data to generate the ONNX model in Pytorch
[1,    20] loss: 0.682
[1,    40] loss: 0.575
....
[50,   120] loss: 0.079
[50,   140] loss: 0.086
Correct Prediction:  138
Total Samples 140

#cd ../caffe2

#python load_model.py ------> (Uses Caffe2 to load the onnx model for inferencing
WARNING:root:This caffe2 python run does not have GPU support.Will run in CPU only mode.
Correct:  663  Total:  700

```

Show moreShow more icon

## 将小尾数 ONNX 模型转换为大尾数模型

ONNX 是一种用于表示深度学习模型的开放格式，旨在实现不同 DL 框架之间的互操作性。基本上，用户可以在一种框架中创建或训练模型，并将其部署到不同框架中，用于进行推断。在内部，ONNX 模型以 Protobuf 格式来表示。这使它独立于框架，并可改进互操作性。但模型文件中的部分元素采用原始格式保存，这就破坏了采用不同尾数法的系统之间的互操作性。即，在小尾数系统上生成的 ONNX 模型无法在大尾数系统上运行。
为解决此问题，我们开发了一种工具，用于使 ONNX 模型实现小尾数与大尾数之间的双向转换。这样可帮助 AIX 客户（大尾数）将 Linux 系统（小尾数）上训练的模型加载到 AIX 上，以便进行推断。我们通过 [ONNX Model Zoo](https://github.com/aixoss/onnx-convertor-le-be.git) 中提供的著名 DL 模型对此工具进行了测试。

可访问以下地址获取此工具的代码： [https://github.com/aixoss/onnx-convertor-le-be.git](https://github.com/aixoss/onnx-convertor-le-be.git)

此代码需要根据自述文件进行编译，并且此工具可用于对模型采用的尾数法进行转换。为了对来自 ONNX Model Zoo 的含有 Protobuf 格式样本数据的部分模型进行测试，样本数据同样需要进行尾数法转换。在以上 Git 存储库中同样提供了用于对样本数据进行尾数法转换的工具，以及有关工具使用方法的说明。

```
git clone https://github.com/aixoss/onnx-convertor-le-be.git

```

Show moreShow more icon

编译工具以进行模型尾数法转换：

```
/usr/bin/g++ -DONNX_NAMESPACE=onnx -DONNX_API= -maix64 -L/opt/freeware/lib/gcc/powerpc-ibm-aix7.2.0.0/8.1.0/pthread -lstdc++ -pthread -I /usr/local/include -mvsx -maix64 -O3 -DNDEBUG -fPIC -std=gnu++11 onnx.pb.cc onnx_conv.cc -I. -L. -lprotobuf -o onnx_conv

```

Show moreShow more icon

编译工具以进行样本数据尾数法转换：

```
/usr/bin/g++ -DONNX_NAMESPACE=onnx -DONNX_API= -maix64 -L/opt/freeware/lib/gcc/powerpc-ibm-aix7.2.0.0/8.1.0/pthread -lstdc++ -pthread -I /usr/local/include -mvsx -maix64 -O3 -DNDEBUG -fPIC -std=gnu++11 onnx.pb.cc tensor_conv.cc -I. -L. -lprotobuf -o tensor_conv

```

Show moreShow more icon

将 LE ONNX 模型（可从以下地址获取： [https://github.com/onnx/models](https://github.com/onnx/models)）转换为 BE 格式：

```
./onnx_conv ./model.onnx ./model.onnx.be

```

Show moreShow more icon

将样本数据转换为大尾数格式：

```
./tensor_conv ./test_data_set_0/input_0.pb ./test_data_set_0/input_0.pb.be
./tensor_conv ./test_data_set_0/output_0.pb ./test_data_set_0/input_0.pb.be

```

Show moreShow more icon

完成模型和样本数据转换后，即可在 AIX 上使用该模型。在位于 [https://github.com/onnx/models#others](https://github.com/onnx/models#others) 处的 ONNX Model Zoo 存储库中提供了样本推断代码。

## 结束语

本教程解释了如何在 AIX 上构建并安装 PyTorch 和 Caffe2，并探讨了 AIX 上的此生态系统所需的许多其他程序包（例如，protobuf、ONNX 和其他 Python 程序包）。相关用例为在 AIX 上使用这些框架提供了有用的示例。

本文翻译自： [Deep learning frameworks – Pytorch and Caffe2 with ONNX support on IBM AIX](https://developer.ibm.com/articles/pytorchcaffe2-with-onnx-on-aix/)（2019-07-22）