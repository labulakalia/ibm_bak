# JVM 优化经验总结
通过实例及对应输出解释的形式让您对于 JVM 优化有一个初步认识

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-lo-jvm-optimize-experience/)

周 明耀

发布: 2015-06-25

* * *

## 开始之前

Java 虚拟机有自己完善的硬件架构, 如处理器、堆栈、寄存器等，还具有相应的指令系统。JVM 屏蔽了与具体操作系统平台相关的信息，使得 Java 程序只需生成在 Java 虚拟机上运行的目标代码 (字节码), 就可以在多种平台上不加修改地运行。Java 虚拟机在执行字节码时，实际上最终还是把字节码解释成具体平台上的机器指令执行。

注意：本文仅针对 JDK7、HotSPOT Java 虚拟机，对于 JDK8 引入的 JVM 新特性及其他 Java 虚拟机，本文不予关注。

我们以一个例子开始这篇文章。假设你是一个普通的 Java 对象，你出生在 Eden 区，在 Eden 区有许多和你差不多的小兄弟、小姐妹，可以把 Eden 区当成幼儿园，在这个幼儿园里大家玩了很长时间。Eden 区不能无休止地放你们在里面，所以当年纪稍大，你就要被送到学校去上学，这里假设从小学到高中都称为 Survivor 区。开始的时候你在 Survivor 区里面划分出来的的 “From”区，读到高年级了，就进了 Survivor 区的 “To”区，中间由于学习成绩不稳定，还经常来回折腾。直到你 18 岁的时候，高中毕业了，该去社会上闯闯了。于是你就去了年老代，年老代里面人也很多。在年老代里，你生活了 20 年 (每次 GC 加一岁)，最后寿终正寝，被 GC 回收。有一点没有提，你在年老代遇到了一个同学，他的名字叫爱德华 (慕光之城里的帅哥吸血鬼)，他以及他的家族永远不会死，那么他们就生活在永生代。

之前的文章 [《JVM 垃圾回收器工作原理及使用实例介绍》](https://developer.ibm.com/zh/articles/j-lo-JVMGarbageCollection/) 中已经介绍过年轻代、年老代、永生代，本文主要讲讲如何运用这些区域，为系统性能提供更好的帮助。本文不再重复这些概念，直接进入主题。

## 如何将新对象预留在年轻代

众所周知，由于 Full GC 的成本远远高于 Minor GC，因此某些情况下需要尽可能将对象分配在年轻代，这在很多情况下是一个明智的选择。虽然在大部分情况下，JVM 会尝试在 Eden 区分配对象，但是由于空间紧张等问题，很可能不得不将部分年轻对象提前向年老代压缩。因此，在 JVM 参数调优时可以为应用程序分配一个合理的年轻代空间，以最大限度避免新对象直接进入年老代的情况发生。清单 1 所示代码尝试分配 4MB 内存空间，观察一下它的内存使用情况。

##### 清单 1\. 相同大小内存分配

```
public class PutInEden {
public static void main(String[] args){
byte[] b1,b2,b3,b4;//定义变量
b1=new byte[1024*1024];//分配 1MB 堆空间，考察堆空间的使用情况
b2=new byte[1024*1024];
b3=new byte[1024*1024];
b4=new byte[1024*1024];
}
}

```

Show moreShow more icon

使用 JVM 参数-XX:+PrintGCDetails -Xmx20M -Xms20M 运行清单 1 所示代码，输出如清单 2 所示。

##### 清单 2\. 清单 1 运行输出

```
[GC [DefNew: 5504K->640K(6144K), 0.0114236 secs] 5504K->5352K(19840K),
0.0114595 secs] [Times: user=0.02 sys=0.00, real=0.02 secs]
[GC [DefNew: 6144K->640K(6144K), 0.0131261 secs] 10856K->10782K(19840K),
0.0131612 secs] [Times: user=0.02 sys=0.00, real=0.02 secs]
[GC [DefNew: 6144K->6144K(6144K), 0.0000170 secs][Tenured: 10142K->13695K(13696K),
0.1069249 secs] 16286K->15966K(19840K), [Perm : 376K->376K(12288K)],
0.1070058 secs] [Times: user=0.03 sys=0.00, real=0.11 secs]
[Full GC [Tenured: 13695K->13695K(13696K), 0.0302067 secs] 19839K->19595K(19840K),
[Perm : 376K->376K(12288K)], 0.0302635 secs] [Times: user=0.03 sys=0.00, real=0.03 secs]
[Full GC [Tenured: 13695K->13695K(13696K), 0.0311986 secs] 19839K->19839K(19840K),
[Perm : 376K->376K(12288K)], 0.0312515 secs] [Times: user=0.03 sys=0.00, real=0.03 secs]
[Full GC [Tenured: 13695K->13695K(13696K), 0.0358821 secs] 19839K->19825K(19840K),
[Perm : 376K->371K(12288K)], 0.0359315 secs] [Times: user=0.05 sys=0.00, real=0.05 secs]
[Full GC [Tenured: 13695K->13695K(13696K), 0.0283080 secs] 19839K->19839K(19840K),
[Perm : 371K->371K(12288K)], 0.0283723 secs] [Times: user=0.02 sys=0.00, real=0.01 secs]
[Full GC [Tenured: 13695K->13695K(13696K), 0.0284469 secs] 19839K->19839K(19840K),
[Perm : 371K->371K(12288K)], 0.0284990 secs] [Times: user=0.03 sys=0.00, real=0.03 secs]
[Full GC [Tenured: 13695K->13695K(13696K), 0.0283005 secs] 19839K->19839K(19840K),
[Perm : 371K->371K(12288K)], 0.0283475 secs] [Times: user=0.03 sys=0.00, real=0.03 secs]
[Full GC [Tenured: 13695K->13695K(13696K), 0.0287757 secs] 19839K->19839K(19840K),
[Perm : 371K->371K(12288K)], 0.0288294 secs] [Times: user=0.03 sys=0.00, real=0.03 secs]
[Full GC [Tenured: 13695K->13695K(13696K), 0.0288219 secs] 19839K->19839K(19840K),
[Perm : 371K->371K(12288K)], 0.0288709 secs] [Times: user=0.03 sys=0.00, real=0.03 secs]
[Full GC [Tenured: 13695K->13695K(13696K), 0.0293071 secs] 19839K->19839K(19840K),
[Perm : 371K->371K(12288K)], 0.0293607 secs] [Times: user=0.03 sys=0.00, real=0.03 secs]
[Full GC [Tenured: 13695K->13695K(13696K), 0.0356141 secs] 19839K->19838K(19840K),
[Perm : 371K->371K(12288K)], 0.0356654 secs] [Times: user=0.01 sys=0.00, real=0.03 secs]
Heap
def new generation total 6144K, used 6143K [0x35c10000, 0x362b0000, 0x362b0000)
eden space 5504K, 100% used [0x35c10000, 0x36170000, 0x36170000)
from space 640K, 99% used [0x36170000, 0x3620fc80, 0x36210000)
to space 640K, 0% used [0x36210000, 0x36210000, 0x362b0000)
tenured generation total 13696K, used 13695K [0x362b0000, 0x37010000, 0x37010000)
the space 13696K, 99% used [0x362b0000, 0x3700fff8, 0x37010000, 0x37010000)
compacting perm gen total 12288K, used 371K [0x37010000, 0x37c10000, 0x3b010000)
the space 12288K, 3% used [0x37010000, 0x3706cd20, 0x3706ce00, 0x37c10000)
ro space 10240K, 51% used [0x3b010000, 0x3b543000, 0x3b543000, 0x3ba10000)
rw space 12288K, 55% used [0x3ba10000, 0x3c0ae4f8, 0x3c0ae600, 0x3c610000)

```

Show moreShow more icon

清单 2 所示的日志输出显示年轻代 Eden 的大小有 5MB 左右。分配足够大的年轻代空间，使用 JVM 参数 -XX:+PrintGCDetails -Xmx20M -Xms20M-Xmn6M 运行清单 1 所示代码，输出如清单 3 所示。

##### 清单 3\. 增大 Eden 大小后清单 1 运行输出

```
[GC [DefNew: 4992K->576K(5568K), 0.0116036 secs] 4992K->4829K(19904K),
0.0116439 secs] [Times: user=0.02 sys=0.00, real=0.02 secs]
[GC [DefNew: 5568K->576K(5568K), 0.0130929 secs] 9821K->9653K(19904K),
0.0131336 secs] [Times: user=0.02 sys=0.00, real=0.02 secs]
[GC [DefNew: 5568K->575K(5568K), 0.0154148 secs] 14645K->14500K(19904K),
0.0154531 secs] [Times: user=0.00 sys=0.01, real=0.01 secs]
[GC [DefNew: 5567K->5567K(5568K), 0.0000197 secs][Tenured: 13924K->14335K(14336K),
0.0330724 secs] 19492K->19265K(19904K), [Perm : 376K->376K(12288K)],
0.0331624 secs] [Times: user=0.03 sys=0.00, real=0.03 secs]
[Full GC [Tenured: 14335K->14335K(14336K), 0.0292459 secs] 19903K->19902K(19904K),
[Perm : 376K->376K(12288K)], 0.0293000 secs] [Times: user=0.03 sys=0.00, real=0.03 secs]
[Full GC [Tenured: 14335K->14335K(14336K), 0.0278675 secs] 19903K->19903K(19904K),
[Perm : 376K->376K(12288K)], 0.0279215 secs] [Times: user=0.03 sys=0.00, real=0.03 secs]
[Full GC [Tenured: 14335K->14335K(14336K), 0.0348408 secs] 19903K->19889K(19904K),
[Perm : 376K->371K(12288K)], 0.0348945 secs] [Times: user=0.05 sys=0.00, real=0.05 secs]
[Full GC [Tenured: 14335K->14335K(14336K), 0.0299813 secs] 19903K->19903K(19904K),
[Perm : 371K->371K(12288K)], 0.0300349 secs] [Times: user=0.01 sys=0.00, real=0.02 secs]
[Full GC [Tenured: 14335K->14335K(14336K), 0.0298178 secs] 19903K->19903K(19904K),
[Perm : 371K->371K(12288K)], 0.0298688 secs] [Times: user=0.03 sys=0.00, real=0.03 secs]
Exception in thread "main" java.lang.OutOfMemoryError: Java heap space[Full GC [Tenured:
14335K->14335K(14336K), 0.0294953 secs] 19903K->19903K(19904K),
[Perm : 371K->371K(12288K)], 0.0295474 secs] [Times: user=0.03 sys=0.00, real=0.03 secs]
[Full GC [Tenured
: 14335K->14335K(14336K), 0.0287742 secs] 19903K->19903K(19904K),
[Perm : 371K->371K(12288K)], 0.0288239 secs] [Times: user=0.03 sys=0.00, real=0.03 secs]
[Full GC [Tenuredat GCTimeTest.main(GCTimeTest.java:16)
: 14335K->14335K(14336K), 0.0287102 secs] 19903K->19903K(19904K),
[Perm : 371K->371K(12288K)], 0.0287627 secs] [Times: user=0.03 sys=0.00, real=0.03 secs]
Heap
def new generation total 5568K, used 5567K [0x35c10000, 0x36210000, 0x36210000)
eden space 4992K, 100% used [0x35c10000, 0x360f0000, 0x360f0000)
from space 576K, 99% used [0x36180000, 0x3620ffe8, 0x36210000)
to space 576K, 0% used [0x360f0000, 0x360f0000, 0x36180000)
tenured generation total 14336K, used 14335K [0x36210000, 0x37010000, 0x37010000)
the space 14336K, 99% used [0x36210000, 0x3700ffd8, 0x37010000, 0x37010000)
compacting perm gen total 12288K, used 371K [0x37010000, 0x37c10000, 0x3b010000)
the space 12288K, 3% used [0x37010000, 0x3706ce28, 0x3706d000, 0x37c10000)
ro space 10240K, 51% used [0x3b010000, 0x3b543000, 0x3b543000, 0x3ba10000)
rw space 12288K, 55% used [0x3ba10000, 0x3c0ae4f8, 0x3c0ae600, 0x3c610000)

```

Show moreShow more icon

通过清单 2 和清单 3 对比，可以发现通过设置一个较大的年轻代预留新对象，设置合理的 Survivor 区并且提供 Survivor 区的使用率，可以将年轻对象保存在年轻代。一般来说，Survivor 区的空间不够，或者占用量达到 50%时，就会使对象进入年老代 (不管它的年龄有多大)。清单 4 创建了 3 个对象，分别分配一定的内存空间。

##### 清单 4\. 不同大小内存分配

```
public class PutInEden2 {
public static void main(String[] args){
byte[] b1,b2,b3;
b1=new byte[1024*512];//分配 0.5MB 堆空间
b2=new byte[1024*1024*4];//分配 4MB 堆空间
b3=new byte[1024*1024*4];
b3=null; //使 b3 可以被回收
b3=new byte[1024*1024*4];//分配 4MB 堆空间
}
}

```

Show moreShow more icon

使用参数-XX:+PrintGCDetails -Xmx1000M -Xms500M -Xmn100M -XX:SurvivorRatio=8 运行清单 4 所示代码，输出如清单 5 所示。

##### 清单 5\. 清单 4 运行输出

```
Heap
def new generation total 92160K, used 11878K [0x0f010000, 0x15410000, 0x15410000)
eden space 81920K, 2% used [0x0f010000, 0x0f1a9a20, 0x14010000)
from space 10240K, 99% used [0x14a10000, 0x1540fff8, 0x15410000)
to space 10240K, 0% used [0x14010000, 0x14010000, 0x14a10000)
tenured generation total 409600K, used 86434K [0x15410000, 0x2e410000, 0x4d810000)
the space 409600K, 21% used [0x15410000, 0x1a878b18, 0x1a878c00, 0x2e410000)
compacting perm gen total 12288K, used 2062K [0x4d810000, 0x4e410000, 0x51810000)
the space 12288K, 16% used [0x4d810000, 0x4da13b18, 0x4da13c00, 0x4e410000)
No shared spaces configured.

```

Show moreShow more icon

清单 5 输出的日志显示，年轻代分配了 8M，年老代也分配了 8M。我们可以尝试加上-XX:TargetSurvivorRatio=90 参数，这样可以提高 from 区的利用率，使 from 区使用到 90%时，再将对象送入年老代，运行清单 4 代码，输出如清单 6 所示。

##### 清单 6\. 修改运行参数后清单 4 输出

```
Heap
def new generation total 9216K, used 9215K [0x35c10000, 0x36610000, 0x36610000)
eden space 8192K, 100% used [0x35c10000, 0x36410000, 0x36410000)
from space 1024K, 99% used [0x36510000, 0x3660fc50, 0x36610000)
to space 1024K, 0% used [0x36410000, 0x36410000, 0x36510000)
tenured generation total 10240K, used 10239K [0x36610000, 0x37010000, 0x37010000)
the space 10240K, 99% used [0x36610000, 0x3700ff70, 0x37010000, 0x37010000)
compacting perm gen total 12288K, used 371K [0x37010000, 0x37c10000, 0x3b010000)
the space 12288K, 3% used [0x37010000, 0x3706cd90, 0x3706ce00, 0x37c10000)
ro space 10240K, 51% used [0x3b010000, 0x3b543000, 0x3b543000, 0x3ba10000)
rw space 12288K, 55% used [0x3ba10000, 0x3c0ae4f8, 0x3c0ae600, 0x3c610000)

```

Show moreShow more icon

如果将 SurvivorRatio 设置为 2，将 b1 对象预存在年轻代。输出如清单 7 所示。

##### 清单 7\. 再次修改运行参数后清单 4 输出

```
Heap
def new generation total 7680K, used 7679K [0x35c10000, 0x36610000, 0x36610000)
eden space 5120K, 100% used [0x35c10000, 0x36110000, 0x36110000)
from space 2560K, 99% used [0x36110000, 0x3638fff0, 0x36390000)
to space 2560K, 0% used [0x36390000, 0x36390000, 0x36610000)
tenured generation total 10240K, used 10239K [0x36610000, 0x37010000, 0x37010000)
the space 10240K, 99% used [0x36610000, 0x3700fff0, 0x37010000, 0x37010000)
compacting perm gen total 12288K, used 371K [0x37010000, 0x37c10000, 0x3b010000)
the space 12288K, 3% used [0x37010000, 0x3706ce28, 0x3706d000, 0x37c10000)
ro space 10240K, 51% used [0x3b010000, 0x3b543000, 0x3b543000, 0x3ba10000)
rw space 12288K, 55% used [0x3ba10000, 0x3c0ae4f8, 0x3c0ae600, 0x3c610000)

```

Show moreShow more icon

## 如何让大对象进入年老代

我们在大部分情况下都会选择将对象分配在年轻代。但是，对于占用内存较多的大对象而言，它的选择可能就不是这样的。因为大对象出现在年轻代很可能扰乱年轻代 GC，并破坏年轻代原有的对象结构。因为尝试在年轻代分配大对象，很可能导致空间不足，为了有足够的空间容纳大对象，JVM 不得不将年轻代中的年轻对象挪到年老代。因为大对象占用空间多，所以可能需要移动大量小的年轻对象进入年老代，这对 GC 相当不利。基于以上原因，可以将大对象直接分配到年老代，保持年轻代对象结构的完整性，这样可以提高 GC 的效率。如果一个大对象同时又是一个短命的对象，假设这种情况出现很频繁，那对于 GC 来说会是一场灾难。原本应该用于存放永久对象的年老代，被短命的对象塞满，这也意味着对堆空间进行了洗牌，扰乱了分代内存回收的基本思路。因此，在软件开发过程中，应该尽可能避免使用短命的大对象。可以使用参数-XX:PetenureSizeThreshold 设置大对象直接进入年老代的阈值。当对象的大小超过这个值时，将直接在年老代分配。参数-XX:PetenureSizeThreshold 只对串行收集器和年轻代并行收集器有效，并行回收收集器不识别这个参数。

##### 清单 8\. 创建一个大对象

```
public class BigObj2Old {
public static void main(String[] args){
byte[] b;
b = new byte[1024*1024];//分配一个 1MB 的对象
}
}

```

Show moreShow more icon

使用 JVM 参数-XX:+PrintGCDetails –Xmx20M –Xms20MB 运行，可以得到清单 9 所示日志输出。

##### 清单 9\. 清单 8 运行输出

```
Heap
def new generation total 6144K, used 1378K [0x35c10000, 0x362b0000, 0x362b0000)
eden space 5504K, 25% used [0x35c10000, 0x35d689e8, 0x36170000)
from space 640K, 0% used [0x36170000, 0x36170000, 0x36210000)
to space 640K, 0% used [0x36210000, 0x36210000, 0x362b0000)
tenured generation total 13696K, used 0K [0x362b0000, 0x37010000, 0x37010000)
the space 13696K, 0% used [0x362b0000, 0x362b0000, 0x362b0200, 0x37010000)
compacting perm gen total 12288K, used 374K [0x37010000, 0x37c10000, 0x3b010000)
the space 12288K, 3% used [0x37010000, 0x3706dac8, 0x3706dc00, 0x37c10000)
ro space 10240K, 51% used [0x3b010000, 0x3b543000, 0x3b543000, 0x3ba10000)
rw space 12288K, 55% used [0x3ba10000, 0x3c0ae4f8, 0x3c0ae600, 0x3c610000)

```

Show moreShow more icon

可以看到该对象被分配在了年轻代，占用了 25%的空间。如果需要将 1MB 以上的对象直接在年老代分配，设置-XX:PetenureSizeThreshold=1000000，程序运行后输出如清单 10 所示。

##### 清单 10\. 修改运行参数后清单 8 输出

```
Heap
def new generation total 6144K, used 354K [0x35c10000, 0x362b0000, 0x362b0000)
eden space 5504K, 6% used [0x35c10000, 0x35c689d8, 0x36170000)
from space 640K, 0% used [0x36170000, 0x36170000, 0x36210000)
to space 640K, 0% used [0x36210000, 0x36210000, 0x362b0000)
tenured generation total 13696K, used 1024K [0x362b0000, 0x37010000, 0x37010000)
the space 13696K, 7% used [0x362b0000, 0x363b0010, 0x363b0200, 0x37010000)
compacting perm gen total 12288K, used 374K [0x37010000, 0x37c10000, 0x3b010000)
the space 12288K, 3% used [0x37010000, 0x3706dac8, 0x3706dc00, 0x37c10000)
ro space 10240K, 51% used [0x3b010000, 0x3b543000, 0x3b543000, 0x3ba10000)
rw space 12288K, 55% used [0x3ba10000, 0x3c0ae4f8, 0x3c0ae600, 0x3c610000)

```

Show moreShow more icon

清单 10 里面可以看到当满 1MB 时进入到了年老代。

## 如何设置对象进入年老代的年龄

堆中的每一个对象都有自己的年龄。一般情况下，年轻对象存放在年轻代，年老对象存放在年老代。为了做到这点，虚拟机为每个对象都维护一个年龄。如果对象在 Eden 区，经过一次 GC 后依然存活，则被移动到 Survivor 区中，对象年龄加 1。以后，如果对象每经过一次 GC 依然存活，则年龄再加 1。当对象年龄达到阈值时，就移入年老代，成为老年对象。这个阈值的最大值可以通过参数-XX:MaxTenuringThreshold 来设置，默认值是 15。虽然-XX:MaxTenuringThreshold 的值可能是 15 或者更大，但这不意味着新对象非要达到这个年龄才能进入年老代。事实上，对象实际进入年老代的年龄是虚拟机在运行时根据内存使用情况动态计算的，这个参数指定的是阈值年龄的最大值。即，实际晋升年老代年龄等于动态计算所得的年龄与-XX:MaxTenuringThreshold 中较小的那个。清单 11 所示代码为 3 个对象申请了若干内存。

##### 清单 11\. 申请内存

```
public class MaxTenuringThreshold {
public static void main(String args[]){
byte[] b1,b2,b3;
b1 = new byte[1024*512];
b2 = new byte[1024*1024*2];
b3 = new byte[1024*1024*4];
b3 = null;
b3 = new byte[1024*1024*4];
}
}

```

Show moreShow more icon

参数设置为：-XX:+PrintGCDetails -Xmx20M -Xms20M -Xmn10M -XX:SurvivorRatio=2

运行清单 11 所示代码，输出如清单 12 所示。

##### 清单 12\. 清单 11 运行输出

```
[GC [DefNew: 2986K->690K(7680K), 0.0246816 secs] 2986K->2738K(17920K),
0.0247226 secs] [Times: user=0.00 sys=0.02, real=0.03 secs]
[GC [DefNew: 4786K->690K(7680K), 0.0016073 secs] 6834K->2738K(17920K),
0.0016436 secs] [Times: user=0.00 sys=0.00, real=0.00 secs]
Heap
def new generation total 7680K, used 4888K [0x35c10000, 0x36610000, 0x36610000)
eden space 5120K, 82% used [0x35c10000, 0x36029a18, 0x36110000)
from space 2560K, 26% used [0x36110000, 0x361bc950, 0x36390000)
to space 2560K, 0% used [0x36390000, 0x36390000, 0x36610000)
tenured generation total 10240K, used 2048K [0x36610000, 0x37010000, 0x37010000)
the space 10240K, 20% used [0x36610000, 0x36810010, 0x36810200, 0x37010000)
compacting perm gen total 12288K, used 374K [0x37010000, 0x37c10000, 0x3b010000)
the space 12288K, 3% used [0x37010000, 0x3706db50, 0x3706dc00, 0x37c10000)
ro space 10240K, 51% used [0x3b010000, 0x3b543000, 0x3b543000, 0x3ba10000)
rw space 12288K, 55% used [0x3ba10000, 0x3c0ae4f8, 0x3c0ae600, 0x3c610000)

```

Show moreShow more icon

更改参数为-XX:+PrintGCDetails -Xmx20M -Xms20M -Xmn10M -XX:SurvivorRatio=2 -XX:MaxTenuringThreshold=1，运行清单 11 所示代码，输出如清单 13 所示。

##### 清单 13\. 修改运行参数后清单 11 输出

```
[GC [DefNew: 2986K->690K(7680K), 0.0047778 secs] 2986K->2738K(17920K),
0.0048161 secs] [Times: user=0.00 sys=0.00, real=0.00 secs]
[GC [DefNew: 4888K->0K(7680K), 0.0016271 secs] 6936K->2738K(17920K),
0.0016630 secs] [Times: user=0.00 sys=0.00, real=0.00 secs]
Heap
def new generation total 7680K, used 4198K [0x35c10000, 0x36610000, 0x36610000)
eden space 5120K, 82% used [0x35c10000, 0x36029a18, 0x36110000)
from space 2560K, 0% used [0x36110000, 0x36110088, 0x36390000)
to space 2560K, 0% used [0x36390000, 0x36390000, 0x36610000)
tenured generation total 10240K, used 2738K [0x36610000, 0x37010000, 0x37010000)
the space 10240K, 26% used [0x36610000, 0x368bc890, 0x368bca00, 0x37010000)
compacting perm gen total 12288K, used 374K [0x37010000, 0x37c10000, 0x3b010000)
the space 12288K, 3% used [0x37010000, 0x3706db50, 0x3706dc00, 0x37c10000)
ro space 10240K, 51% used [0x3b010000, 0x3b543000, 0x3b543000, 0x3ba10000)
rw space 12288K, 55% used [0x3ba10000, 0x3c0ae4f8, 0x3c0ae600, 0x3c610000)

```

Show moreShow more icon

清单 13 所示，第一次运行时 b1 对象在程序结束后依然保存在年轻代。第二次运行前，我们减小了对象晋升年老代的年龄，设置为 1。即，所有经过一次 GC 的对象都可以直接进入年老代。程序运行后，可以发现 b1 对象已经被分配到年老代。如果希望对象尽可能长时间地停留在年轻代，可以设置一个较大的阈值。

## 稳定的 Java 堆 VS 动荡的 Java 堆

一般来说，稳定的堆大小对垃圾回收是有利的。获得一个稳定的堆大小的方法是使-Xms 和-Xmx 的大小一致，即最大堆和最小堆 (初始堆) 一样。如果这样设置，系统在运行时堆大小理论上是恒定的，稳定的堆空间可以减少 GC 的次数。因此，很多服务端应用都会将最大堆和最小堆设置为相同的数值。但是，一个不稳定的堆并非毫无用处。稳定的堆大小虽然可以减少 GC 次数，但同时也增加了每次 GC 的时间。让堆大小在一个区间中震荡，在系统不需要使用大内存时，压缩堆空间，使 GC 应对一个较小的堆，可以加快单次 GC 的速度。基于这样的考虑，JVM 还提供了两个参数用于压缩和扩展堆空间。

-XX:MinHeapFreeRatio 参数用来设置堆空间最小空闲比例，默认值是 40。当堆空间的空闲内存小于这个数值时，JVM 便会扩展堆空间。

-XX:MaxHeapFreeRatio 参数用来设置堆空间最大空闲比例，默认值是 70。当堆空间的空闲内存大于这个数值时，便会压缩堆空间，得到一个较小的堆。

当-Xmx 和-Xms 相等时，-XX:MinHeapFreeRatio 和-XX:MaxHeapFreeRatio 两个参数无效。

##### 清单 14\. 堆大小设置

```
import java.util.Vector;

public class HeapSize {
public static void main(String args[]) throws InterruptedException{
Vector v = new Vector();
while(true){
byte[] b = new byte[1024*1024];
v.add(b);
if(v.size() == 10){
v = new Vector();
}
Thread.sleep(1);
}
}
}

```

Show moreShow more icon

清单 14 所示代码是测试-XX:MinHeapFreeRatio 和-XX:MaxHeapFreeRatio 的作用，设置运行参数为-XX:+PrintGCDetails -Xms10M -Xmx40M -XX:MinHeapFreeRatio=40 -XX:MaxHeapFreeRatio=50 时，输出如清单 15 所示。

##### 清单 15\. 修改运行参数后清单 14 输出

```
[GC [DefNew: 2418K->178K(3072K), 0.0034827 secs] 2418K->2226K(9920K),
0.0035249 secs] [Times: user=0.00 sys=0.00, real=0.03 secs]
[GC [DefNew: 2312K->0K(3072K), 0.0028263 secs] 4360K->4274K(9920K),
0.0029905 secs] [Times: user=0.00 sys=0.00, real=0.03 secs]
[GC [DefNew: 2068K->0K(3072K), 0.0024363 secs] 6342K->6322K(9920K),
0.0024836 secs] [Times: user=0.00 sys=0.00, real=0.03 secs]
[GC [DefNew: 2061K->0K(3072K), 0.0017376 secs][Tenured: 8370K->8370K(8904K),
0.1392692 secs] 8384K->8370K(11976K), [Perm : 374K->374K(12288K)],
0.1411363 secs] [Times: user=0.00 sys=0.02, real=0.16 secs]
[GC [DefNew: 5138K->0K(6336K), 0.0038237 secs] 13508K->13490K(20288K),
0.0038632 secs] [Times: user=0.00 sys=0.00, real=0.03 secs]

```

Show moreShow more icon

改用参数：-XX:+PrintGCDetails -Xms40M -Xmx40M -XX:MinHeapFreeRatio=40 -XX:MaxHeapFreeRatio=50，运行输出如清单 16 所示。

##### 清单 16\. 再次修改运行参数后清单 14 输出

```
[GC [DefNew: 10678K->178K(12288K), 0.0019448 secs] 10678K->178K(39616K),
0.0019851 secs] [Times: user=0.00 sys=0.00, real=0.03 secs]
[GC [DefNew: 10751K->178K(12288K), 0.0010295 secs] 10751K->178K(39616K),
0.0010697 secs] [Times: user=0.00 sys=0.00, real=0.02 secs]
[GC [DefNew: 10493K->178K(12288K), 0.0008301 secs] 10493K->178K(39616K),
0.0008672 secs] [Times: user=0.00 sys=0.00, real=0.02 secs]
[GC [DefNew: 10467K->178K(12288K), 0.0008522 secs] 10467K->178K(39616K),
0.0008905 secs] [Times: user=0.00 sys=0.00, real=0.02 secs]
[GC [DefNew: 10450K->178K(12288K), 0.0008964 secs] 10450K->178K(39616K),
0.0009339 secs] [Times: user=0.00 sys=0.00, real=0.01 secs]
[GC [DefNew: 10439K->178K(12288K), 0.0009876 secs] 10439K->178K(39616K),
0.0010279 secs] [Times: user=0.00 sys=0.00, real=0.02 secs]

```

Show moreShow more icon

从清单 16 可以看出，此时堆空间的垃圾回收稳定在一个固定的范围。在一个稳定的堆中，堆空间大小始终不变，每次 GC 时，都要应对一个 40MB 的空间。因此，虽然 GC 次数减小了，但是单次 GC 速度不如一个震荡的堆。

## 增大吞吐量提升系统性能

吞吐量优先的方案将会尽可能减少系统执行垃圾回收的总时间，故可以考虑关注系统吞吐量的并行回收收集器。在拥有高性能的计算机上，进行吞吐量优先优化，可以使用参数：

```
java –Xmx3800m –Xms3800m –Xmn2G –Xss128k –XX:+UseParallelGC
–XX:ParallelGC-Threads=20 –XX:+UseParallelOldGC

```

Show moreShow more icon

–Xmx3800m –Xms3800m：设置 Java 堆的最大值和初始值。一般情况下，为了避免堆内存的频繁震荡，导致系统性能下降，我们的做法是设置最大堆等于最小堆。假设这里把最小堆减少为最大堆的一半，即 1900m，那么 JVM 会尽可能在 1900MB 堆空间中运行，如果这样，发生 GC 的可能性就会比较高；

-Xss128k：减少线程栈的大小，这样可以使剩余的系统内存支持更多的线程；

-Xmn2g：设置年轻代区域大小为 2GB；

–XX:+UseParallelGC：年轻代使用并行垃圾回收收集器。这是一个关注吞吐量的收集器，可以尽可能地减少 GC 时间。

–XX:ParallelGC-Threads：设置用于垃圾回收的线程数，通常情况下，可以设置和 CPU 数量相等。但在 CPU 数量比较多的情况下，设置相对较小的数值也是合理的；

–XX:+UseParallelOldGC：设置年老代使用并行回收收集器。

## 尝试使用大的内存分页

CPU 是通过寻址来访问内存的。32 位 CPU 的寻址宽度是 0~0xFFFFFFFF ，计算后得到的大小是 4G，也就是说可支持的物理内存最大是 4G。但在实践过程中，碰到了这样的问题，程序需要使用 4G 内存，而可用物理内存小于 4G，导致程序不得不降低内存占用。为了解决此类问题，现代 CPU 引入了 MMU（Memory Management Unit 内存管理单元）。MMU 的核心思想是利用虚拟地址替代物理地址，即 CPU 寻址时使用虚址，由 MMU 负责将虚址映射为物理地址。MMU 的引入，解决了对物理内存的限制，对程序来说，就像自己在使用 4G 内存一样。内存分页 (Paging) 是在使用 MMU 的基础上，提出的一种内存管理机制。它将虚拟地址和物理地址按固定大小（4K）分割成页 (page) 和页帧 (page frame)，并保证页与页帧的大小相同。这种机制，从数据结构上，保证了访问内存的高效，并使 OS 能支持非连续性的内存分配。在程序内存不够用时，还可以将不常用的物理内存页转移到其他存储设备上，比如磁盘，这就是大家耳熟能详的虚拟内存。

在 Solaris 系统中，JVM 可以支持 Large Page Size 的使用。使用大的内存分页可以增强 CPU 的内存寻址能力，从而提升系统的性能。

```
java –Xmx2506m –Xms2506m –Xmn1536m –Xss128k –XX:++UseParallelGC
–XX:ParallelGCThreads=20 –XX:+UseParallelOldGC –XX:+LargePageSizeInBytes=256m

```

Show moreShow more icon

–XX:+LargePageSizeInBytes：设置大页的大小。

过大的内存分页会导致 JVM 在计算 Heap 内部分区（perm, new, old）内存占用比例时，会出现超出正常值的划分，最坏情况下某个区会多占用一个页的大小。

## 使用非占有的垃圾回收器

为降低应用软件的垃圾回收时的停顿，首先考虑的是使用关注系统停顿的 CMS 回收器，其次，为了减少 Full GC 次数，应尽可能将对象预留在年轻代，因为年轻代 Minor GC 的成本远远小于年老代的 Full GC。

```
java –Xmx3550m –Xms3550m –Xmn2g –Xss128k –XX:ParallelGCThreads=20
–XX:+UseConcMarkSweepGC –XX:+UseParNewGC –XX:+SurvivorRatio=8 –XX:TargetSurvivorRatio=90
–XX:MaxTenuringThreshold=31

```

Show moreShow more icon

–XX:ParallelGCThreads=20：设置 20 个线程进行垃圾回收；

–XX:+UseParNewGC：年轻代使用并行回收器；

–XX:+UseConcMarkSweepGC：年老代使用 CMS 收集器降低停顿；

–XX:+SurvivorRatio：设置 Eden 区和 Survivor 区的比例为 8:1。稍大的 Survivor 空间可以提高在年轻代回收生命周期较短的对象的可能性，如果 Survivor 不够大，一些短命的对象可能直接进入年老代，这对系统来说是不利的。

–XX:TargetSurvivorRatio=90：设置 Survivor 区的可使用率。这里设置为 90%，则允许 90%的 Survivor 空间被使用。默认值是 50%。故该设置提高了 Survivor 区的使用率。当存放的对象超过这个百分比，则对象会向年老代压缩。因此，这个选项更有助于将对象留在年轻代。

–XX:MaxTenuringThreshold：设置年轻对象晋升到年老代的年龄。默认值是 15 次，即对象经过 15 次 Minor GC 依然存活，则进入年老代。这里设置为 31，目的是让对象尽可能地保存在年轻代区域。

## 结束语

通过本文的学习，读者了解了如何将新对象预留在年轻代、如何让大对象进入年老代、如何设置对象进入年老代的年龄、稳定的 Java 堆 VS 动荡的 Java 堆、增大吞吐量提升系统性能、尝试使用大的内存分页、使用非占有的垃圾回收器等主题，通过实例及对应输出解释的形式让读者对于 JVM 优化有一个初步认识。如其他文章相同的观点，没有哪一条优化是固定不变的，读者需要自己判断、实践后才能找到正确的道路。