# JVM 垃圾回收器工作原理及使用实例介绍
掌握基本的 JVM 垃圾回收器设计原理及使用规范并通过实验来找出最适合自己的优化方案

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-lo-jvmgarbagecollection/)

周 明耀

发布: 2015-04-08

* * *

## 垃圾收集基础

Java 语言的一大特点就是可以进行自动垃圾回收处理，而无需开发人员过于关注系统资源，例如内存资源的释放情况。自动垃圾收集虽然大大减轻了开发人员的工作量，但是也增加了软件系统的负担。

拥有垃圾收集器可以说是 Java 语言与 C++语言的一项显著区别。在 C++语言中，程序员必须小心谨慎地处理每一项内存分配，且内存使用完后必须手工释放曾经占用的内存空间。当内存释放不够完全时，即存在分配但永不释放的内存块，就会引起内存泄漏，严重时甚至导致程序瘫痪。

以下列举了垃圾回收器常用的算法及实验原理：

- 引用计数法 (Reference Counting)

引用计数器在微软的 COM 组件技术中、Adobe 的 ActionScript3 种都有使用。

引用计数器的实现很简单，对于一个对象 A，只要有任何一个对象引用了 A，则 A 的引用计数器就加 1，当引用失效时，引用计数器就减 1。只要对象 A 的引用计数器的值为 0，则对象 A 就不可能再被使用。

引用计数器的实现也非常简单，只需要为每个对象配置一个整形的计数器即可。但是引用计数器有一个严重的问题，即无法处理循环引用的情况。因此，在 Java 的垃圾回收器中没有使用这种算法。

一个简单的循环引用问题描述如下：有对象 A 和对象 B，对象 A 中含有对象 B 的引用，对象 B 中含有对象 A 的引用。此时，对象 A 和对象 B 的引用计数器都不为 0。但是在系统中却不存在任何第 3 个对象引用了 A 或 B。也就是说，A 和 B 是应该被回收的垃圾对象，但由于垃圾对象间相互引用，从而使垃圾回收器无法识别，引起内存泄漏。

- 标记-清除算法 (Mark-Sweep)

标记-清除算法将垃圾回收分为两个阶段：标记阶段和清除阶段。一种可行的实现是，在标记阶段首先通过根节点，标记所有从根节点开始的较大对象。因此，未被标记的对象就是未被引用的垃圾对象。然后，在清除阶段，清除所有未被标记的对象。该算法最大的问题是存在大量的空间碎片，因为回收后的空间是不连续的。在对象的堆空间分配过程中，尤其是大对象的内存分配，不连续的内存空间的工作效率要低于连续的空间。

- 复制算法 (Copying)

将现有的内存空间分为两快，每次只使用其中一块，在垃圾回收时将正在使用的内存中的存活对象复制到未被使用的内存块中，之后，清除正在使用的内存块中的所有对象，交换两个内存的角色，完成垃圾回收。

如果系统中的垃圾对象很多，复制算法需要复制的存活对象数量并不会太大。因此在真正需要垃圾回收的时刻，复制算法的效率是很高的。又由于对象在垃圾回收过程中统一被复制到新的内存空间中，因此，可确保回收后的内存空间是没有碎片的。该算法的缺点是将系统内存折半。

Java 的新生代串行垃圾回收器中使用了复制算法的思想。新生代分为 eden 空间、from 空间、to 空间 3 个部分。其中 from 空间和 to 空间可以视为用于复制的两块大小相同、地位相等，且可进行角色互换的空间块。from 和 to 空间也称为 survivor 空间，即幸存者空间，用于存放未被回收的对象。

在垃圾回收时，eden 空间中的存活对象会被复制到未使用的 survivor 空间中 (假设是 to)，正在使用的 survivor 空间 (假设是 from) 中的年轻对象也会被复制到 to 空间中 (大对象，或者老年对象会直接进入老年带，如果 to 空间已满，则对象也会直接进入老年代)。此时，eden 空间和 from 空间中的剩余对象就是垃圾对象，可以直接清空，to 空间则存放此次回收后的存活对象。这种改进的复制算法既保证了空间的连续性，又避免了大量的内存空间浪费。

- 标记-压缩算法 (Mark-Compact)

复制算法的高效性是建立在存活对象少、垃圾对象多的前提下的。这种情况在年轻代经常发生，但是在老年代更常见的情况是大部分对象都是存活对象。如果依然使用复制算法，由于存活的对象较多，复制的成本也将很高。

标记-压缩算法是一种老年代的回收算法，它在标记-清除算法的基础上做了一些优化。也首先需要从根节点开始对所有可达对象做一次标记，但之后，它并不简单地清理未标记的对象，而是将所有的存活对象压缩到内存的一端。之后，清理边界外所有的空间。这种方法既避免了碎片的产生，又不需要两块相同的内存空间，因此，其性价比比较高。

- 增量算法 (Incremental Collecting)

在垃圾回收过程中，应用软件将处于一种 CPU 消耗很高的状态。在这种 CPU 消耗很高的状态下，应用程序所有的线程都会挂起，暂停一切正常的工作，等待垃圾回收的完成。如果垃圾回收时间过长，应用程序会被挂起很久，将严重影响用户体验或者系统的稳定性。

增量算法的基本思想是，如果一次性将所有的垃圾进行处理，需要造成系统长时间的停顿，那么就可以让垃圾收集线程和应用程序线程交替执行。每次，垃圾收集线程只收集一小片区域的内存空间，接着切换到应用程序线程。依次反复，直到垃圾收集完成。使用这种方式，由于在垃圾回收过程中，间断性地还执行了应用程序代码，所以能减少系统的停顿时间。但是，因为线程切换和上下文转换的消耗，会使得垃圾回收的总体成本上升，造成系统吞吐量的下降。

- 分代 (Generational Collecting)

根据垃圾回收对象的特性，不同阶段最优的方式是使用合适的算法用于本阶段的垃圾回收，分代算法即是基于这种思想，它将内存区间根据对象的特点分成几块，根据每块内存区间的特点，使用不同的回收算法，以提高垃圾回收的效率。以 Hot Spot 虚拟机为例，它将所有的新建对象都放入称为年轻代的内存区域，年轻代的特点是对象会很快回收，因此，在年轻代就选择效率较高的复制算法。当一个对象经过几次回收后依然存活，对象就会被放入称为老生代的内存空间。在老生代中，几乎所有的对象都是经过几次垃圾回收后依然得以幸存的。因此，可以认为这些对象在一段时期内，甚至在应用程序的整个生命周期中，将是常驻内存的。如果依然使用复制算法回收老生代，将需要复制大量对象。再加上老生代的回收性价比也要低于新生代，因此这种做法也是不可取的。根据分代的思想，可以对老年代的回收使用与新生代不同的标记-压缩算法，以提高垃圾回收效率。

从不同角度分析垃圾收集器，可以将其分为不同的类型。

1. 按线程数分，可以分为串行垃圾回收器和并行垃圾回收器。串行垃圾回收器一次只使用一个线程进行垃圾回收；并行垃圾回收器一次将开启多个线程同时进行垃圾回收。在并行能力较强的 CPU 上，使用并行垃圾回收器可以缩短 GC 的停顿时间。

2. 按照工作模式分，可以分为并发式垃圾回收器和独占式垃圾回收器。并发式垃圾回收器与应用程序线程交替工作，以尽可能减少应用程序的停顿时间；独占式垃圾回收器 (Stop the world) 一旦运行，就停止应用程序中的其他所有线程，直到垃圾回收过程完全结束。

3. 按碎片处理方式可分为压缩式垃圾回收器和非压缩式垃圾回收器。压缩式垃圾回收器会在回收完成后，对存活对象进行压缩整理，消除回收后的碎片；非压缩式的垃圾回收器不进行这步操作。

4. 按工作的内存区间，又可分为新生代垃圾回收器和老年代垃圾回收器。


可以用以下指标评价一个垃圾处理器的好坏。

吞吐量：指在应用程序的生命周期内，应用程序所花费的时间和系统总运行时间的比值。系统总运行时间=应用程序耗时+GC 耗时。如果系统运行了 100min，GC 耗时 1min，那么系统的吞吐量就是 (100-1)/100=99%。

垃圾回收器负载：和吞吐量相反，垃圾回收器负载指来记回收器耗时与系统运行总时间的比值。

停顿时间：指垃圾回收器正在运行时，应用程序的暂停时间。对于独占回收器而言，停顿时间可能会比较长。使用并发的回收器时，由于垃圾回收器和应用程序交替运行，程序的停顿时间会变短，但是，由于其效率很可能不如独占垃圾回收器，故系统的吞吐量可能会较低。

垃圾回收频率：指垃圾回收器多长时间会运行一次。一般来说，对于固定的应用而言，垃圾回收器的频率应该是越低越好。通常增大堆空间可以有效降低垃圾回收发生的频率，但是可能会增加回收产生的停顿时间。

反应时间：指当一个对象被称为垃圾后多长时间内，它所占据的内存空间会被释放。

堆分配：不同的垃圾回收器对堆内存的分配方式可能是不同的。一个良好的垃圾收集器应该有一个合理的堆内存区间划分。

## JVM 垃圾回收器分类

- 新生代串行收集器

串行收集器主要有两个特点：第一，它仅仅使用单线程进行垃圾回收；第二，它独占式的垃圾回收。

在串行收集器进行垃圾回收时，Java 应用程序中的线程都需要暂停，等待垃圾回收的完成，这样给用户体验造成较差效果。虽然如此，串行收集器却是一个成熟、经过长时间生产环境考验的极为高效的收集器。新生代串行处理器使用复制算法，实现相对简单，逻辑处理特别高效，且没有线程切换的开销。在诸如单 CPU 处理器或者较小的应用内存等硬件平台不是特别优越的场合，它的性能表现可以超过并行回收器和并发回收器。在 HotSpot 虚拟机中，使用-XX：+UseSerialGC 参数可以指定使用新生代串行收集器和老年代串行收集器。当 JVM 在 Client 模式下运行时，它是默认的垃圾收集器。一次新生代串行收集器的工作输出日志类似如清单 1 信息 (使用-XX:+PrintGCDetails 开关) 所示。

##### 清单 1\. 一次新生代串行收集器的工作输出日志

```
[GC [DefNew: 3468K->150K(9216K), 0.0028638 secs][Tenured:
1562K->1712K(10240K), 0.0084220 secs] 3468K->1712K(19456K),
[Perm : 377K->377K(12288K)],
0.0113816 secs] [Times: user=0.02 sys=0.00, real=0.01 secs]

```

Show moreShow more icon

它显示了一次垃圾回收前的新生代的内存占用量和垃圾回收后的新生代内存占用量，以及垃圾回收所消耗的时间。

- 老年代串行收集器

老年代串行收集器使用的是标记-压缩算法。和新生代串行收集器一样，它也是一个串行的、独占式的垃圾回收器。由于老年代垃圾回收通常会使用比新生代垃圾回收更长的时间，因此，在堆空间较大的应用程序中，一旦老年代串行收集器启动，应用程序很可能会因此停顿几秒甚至更长时间。虽然如此，老年代串行回收器可以和多种新生代回收器配合使用，同时它也可以作为 CMS 回收器的备用回收器。若要启用老年代串行回收器，可以尝试使用以下参数：-XX:+UseSerialGC: 新生代、老年代都使用串行回收器。

##### 清单 2\. 一次老年代串行收集器的工作输出日志

```
Heap
def new generation total 4928K, used 1373K [0x27010000, 0x27560000, 0x2c560000)
eden space 4416K, 31% used [0x27010000, 0x27167628, 0x27460000)
from space 512K, 0% used [0x27460000, 0x27460000, 0x274e0000)
to space 512K, 0% used [0x274e0000, 0x274e0000, 0x27560000)
tenured generation total 10944K, used 0K [0x2c560000, 0x2d010000, 0x37010000)
the space 10944K, 0% used [0x2c560000, 0x2c560000, 0x2c560200, 0x2d010000)
compacting perm gen total 12288K, used 376K [0x37010000, 0x37c10000, 0x3b010000)
the space 12288K, 3% used [0x37010000, 0x3706e0b8, 0x3706e200, 0x37c10000)
ro space 10240K, 51% used [0x3b010000, 0x3b543000, 0x3b543000, 0x3ba10000)
rw space 12288K, 55% used [0x3ba10000, 0x3c0ae4f8, 0x3c0ae600, 0x3c610000)

```

Show moreShow more icon

如果使用-XX:+UseParNewGC 参数设置，表示新生代使用并行收集器，老年代使用串行收集器，如清单 3 所示。

##### 清单 3\. 一次串并行收集器混合使用的工作输出日志

```
Heap
par new generation total 4928K, used 1373K [0x0f010000, 0x0f560000, 0x14560000)
eden space 4416K, 31% used [0x0f010000, 0x0f167620, 0x0f460000)
from space 512K, 0% used [0x0f460000, 0x0f460000, 0x0f4e0000)
to space 512K, 0% used [0x0f4e0000, 0x0f4e0000, 0x0f560000)
tenured generation total 10944K, used 0K [0x14560000, 0x15010000, 0x1f010000)
the space 10944K, 0% used [0x14560000, 0x14560000, 0x14560200, 0x15010000)
compacting perm gen total 12288K, used 2056K [0x1f010000, 0x1fc10000, 0x23010000)
the space 12288K, 16% used [0x1f010000, 0x1f2121d0, 0x1f212200, 0x1fc10000)
No shared spaces configured.

```

Show moreShow more icon

如果使用-XX:+UseParallelGC 参数设置，表示新生代和老年代均使用并行回收收集器。如清单 4 所示。

##### 清单 4\. 一次老年代并行回收器的工作输出日志

```
[Full GC [Tenured: 1712K->1699K(10240K), 0.0071963 secs] 1712K->1699K(19456K),
      [Perm : 377K->372K(12288K)], 0.0072393 secs] [Times: user=0.00 sys=0.00,
      real=0.01 secs]

```

Show moreShow more icon

它显示了垃圾回收前老年代和永久区的内存占用量，以及垃圾回收后老年代和永久区的内存使用量。

- 并行收集器

并行收集器是工作在新生代的垃圾收集器，它只简单地将串行回收器多线程化。它的回收策略、算法以及参数和串行回收器一样。

并行回收器也是独占式的回收器，在收集过程中，应用程序会全部暂停。但由于并行回收器使用多线程进行垃圾回收，因此，在并发能力比较强的 CPU 上，它产生的停顿时间要短于串行回收器，而在单 CPU 或者并发能力较弱的系统中，并行回收器的效果不会比串行回收器好，由于多线程的压力，它的实际表现很可能比串行回收器差。

开启并行回收器可以使用参数-XX:+UseParNewGC，该参数设置新生代使用并行收集器，老年代使用串行收集器。

##### 清单 5\. 设置参数-XX:+UseParNewGC 的输出日志

```
[GC [ParNew: 825K->161K(4928K), 0.0155258 secs][Tenured: 8704K->661K(10944K),
0.0071964 secs] 9017K->661K(15872K),
[Perm : 2049K->2049K(12288K)], 0.0228090 secs] [Times: user=0.01 sys=0.00, real=0.01 secs]
Heap
par new generation total 4992K, used 179K [0x0f010000, 0x0f570000, 0x14560000)
eden space 4480K, 4% used [0x0f010000, 0x0f03cda8, 0x0f470000)
from space 512K, 0% used [0x0f470000, 0x0f470000, 0x0f4f0000)
to space 512K, 0% used [0x0f4f0000, 0x0f4f0000, 0x0f570000)
tenured generation total 10944K, used 8853K [0x14560000, 0x15010000, 0x1f010000)
the space 10944K, 80% used [0x14560000, 0x14e057c0, 0x14e05800, 0x15010000)
compacting perm gen total 12288K, used 2060K [0x1f010000, 0x1fc10000, 0x23010000)
the space 12288K, 16% used [0x1f010000, 0x1f213228, 0x1f213400, 0x1fc10000)
No shared spaces configured.

```

Show moreShow more icon

设置参数-XX:+UseConcMarkSweepGC 可以要求新生代使用并行收集器，老年代使用 CMS。

##### 清单 6\. 设置参数-XX:+UseConcMarkSweepGC 的输出日志

```
[GC [ParNew: 8967K->669K(14784K), 0.0040895 secs] 8967K->669K(63936K),
0.0043255 secs] [Times: user=0.00 sys=0.00, real=0.00 secs]
Heap
par new generation total 14784K, used 9389K [0x03f50000, 0x04f50000, 0x04f50000)
eden space 13184K, 66% used [0x03f50000, 0x047d3e58, 0x04c30000)
from space 1600K, 41% used [0x04dc0000, 0x04e67738, 0x04f50000)
to space 1600K, 0% used [0x04c30000, 0x04c30000, 0x04dc0000)
concurrent mark-sweep generation total 49152K, used 0K [0x04f50000, 0x07f50000, 0x09f50000)
concurrent-mark-sweep perm gen total 12288K, used 2060K [0x09f50000, 0x0ab50000, 0x0df50000)

```

Show moreShow more icon

并行收集器工作时的线程数量可以使用-XX:ParallelGCThreads 参数指定。一般，最好与 CPU 数量相当，避免过多的线程数影响垃圾收集性能。在默认情况下，当 CPU 数量小于 8 个，ParallelGCThreads 的值等于 CPU 数量，大于 8 个，ParallelGCThreads 的值等于 3+[5\*CPU\_Count]/8]。以下测试显示了笔者笔记本上运行 8 个线程时耗时最短，本人笔记本是 8 核 IntelCPU。

##### 清单 7\. 设置为 8 个线程后 GC 输出

```
[GC [ParNew: 8967K->676K(14784K), 0.0036983 secs] 8967K->676K(63936K),
0.0037662 secs] [Times: user=0.00 sys=0.00, real=0.00 secs]
Heap
par new generation total 14784K, used 9395K [0x040e0000, 0x050e0000, 0x050e0000)
eden space 13184K, 66% used [0x040e0000, 0x04963e58, 0x04dc0000)
from space 1600K, 42% used [0x04f50000, 0x04ff9100, 0x050e0000)
to space 1600K, 0% used [0x04dc0000, 0x04dc0000, 0x04f50000)
concurrent mark-sweep generation total 49152K, used 0K [0x050e0000, 0x080e0000, 0x0a0e0000)
concurrent-mark-sweep perm gen total 12288K, used 2060K [0x0a0e0000, 0x0ace0000, 0x0e0e0000)

```

Show moreShow more icon

##### 清单 8\. 设置为 128 个线程后 GC 输出

```
[GC [ParNew: 8967K->664K(14784K), 0.0207327 secs] 8967K->664K(63936K),
0.0208077 secs] [Times: user=0.03 sys=0.00, real=0.02 secs]

```

Show moreShow more icon

##### 清单 9\. 设置为 640 个线程后 GC 输出

```
[GC [ParNew: 8967K->667K(14784K), 0.2323704 secs] 8967K->667K(63936K),
0.2324778 secs] [Times: user=0.34 sys=0.02, real=0.23 secs]

```

Show moreShow more icon

##### 清单 10\. 设置为 1280 个线程后 GC 输出

```
Error occurred during initialization of VM
Too small new size specified

```

Show moreShow more icon

- 新生代并行回收 (Parallel Scavenge) 收集器

新生代并行回收收集器也是使用复制算法的收集器。从表面上看，它和并行收集器一样都是多线程、独占式的收集器。但是，并行回收收集器有一个重要的特点：它非常关注系统的吞吐量。

新生代并行回收收集器可以使用以下参数启用：

-XX:+UseParallelGC:新生代使用并行回收收集器，老年代使用串行收集器。

-XX:+UseParallelOldGC:新生代和老年代都是用并行回收收集器。

##### 清单 11\. 设置为 24 个线程后 GC 输出

```
Heap
PSYoungGen total 4800K, used 893K [0x1dac0000, 0x1e010000, 0x23010000)
eden space 4160K, 21% used [0x1dac0000,0x1db9f570,0x1ded0000)
from space 640K, 0% used [0x1df70000,0x1df70000,0x1e010000)
to space 640K, 0% used [0x1ded0000,0x1ded0000,0x1df70000)
ParOldGen total 19200K, used 16384K [0x13010000, 0x142d0000, 0x1dac0000)
object space 19200K, 85% used [0x13010000,0x14010020,0x142d0000)
PSPermGen total 12288K, used 2054K [0x0f010000, 0x0fc10000, 0x13010000)
object space 12288K, 16% used [0x0f010000,0x0f2119c0,0x0fc10000)

```

Show moreShow more icon

新生代并行回收收集器可以使用以下参数启用：

-XX:+MaxGCPauseMills:设置最大垃圾收集停顿时间，它的值是一个大于 0 的整数。收集器在工作时会调整 Java 堆大小或者其他一些参数，尽可能地把停顿时间控制在 MaxGCPauseMills 以内。如果希望减少停顿时间，而把这个值设置得很小，为了达到预期的停顿时间，JVM 可能会使用一个较小的堆 (一个小堆比一个大堆回收快)，而这将导致垃圾回收变得很频繁，从而增加了垃圾回收总时间，降低了吞吐量。

-XX:+GCTimeRatio：设置吞吐量大小，它的值是一个 0-100 之间的整数。假设 GCTimeRatio 的值为 n，那么系统将花费不超过 1/(1+n) 的时间用于垃圾收集。比如 GCTimeRatio 等于 19，则系统用于垃圾收集的时间不超过 1/(1+19)=5%。默认情况下，它的取值是 99，即不超过 1%的时间用于垃圾收集。

除此之外，并行回收收集器与并行收集器另一个不同之处在于，它支持一种自适应的 GC 调节策略，使用-XX:+UseAdaptiveSizePolicy 可以打开自适应 GC 策略。在这种模式下，新生代的大小、eden 和 survivor 的比例、晋升老年代的对象年龄等参数会被自动调整，以达到在堆大小、吞吐量和停顿时间之间的平衡点。在手工调优比较困难的场合，可以直接使用这种自适应的方式，仅指定虚拟机的最大堆、目标的吞吐量 (GCTimeRatio) 和停顿时间 (MaxGCPauseMills)，让虚拟机自己完成调优工作。

##### 清单 12\. 新生代并行回收收集器 GC 输出

```
Heap
PSYoungGen total 4800K, used 893K [0x1dac0000, 0x1e010000, 0x23010000)
eden space 4160K, 21% used [0x1dac0000,0x1db9f570,0x1ded0000)
from space 640K, 0% used [0x1df70000,0x1df70000,0x1e010000)
to space 640K, 0% used [0x1ded0000,0x1ded0000,0x1df70000)
PSOldGen total 19200K, used 16384K [0x13010000, 0x142d0000, 0x1dac0000)
object space 19200K, 85% used [0x13010000,0x14010020,0x142d0000)
PSPermGen total 12288K, used 2054K [0x0f010000, 0x0fc10000, 0x13010000)
object space 12288K, 16% used [0x0f010000,0x0f2119c0,0x0fc10000)

```

Show moreShow more icon

它也显示了收集器的工作成果，也就是回收前的内存大小和回收后的内存大小，以及花费的时间。

- 老年代并行回收收集器

老年代的并行回收收集器也是一种多线程并发的收集器。和新生代并行回收收集器一样，它也是一种关注吞吐量的收集器。老年代并行回收收集器使用标记-压缩算法，JDK1.6 之后开始启用。

使用-XX:+UseParallelOldGC 可以在新生代和老生代都使用并行回收收集器，这是一对非常关注吞吐量的垃圾收集器组合，在对吞吐量敏感的系统中，可以考虑使用。参数-XX:ParallelGCThreads 也可以用于设置垃圾回收时的线程数量。

清单 13 是设置线程数量为 100 时老年代并行回收收集器输出日志：

##### 清单 13\. 老年代并行回收收集器设置 100 线程时 GC 输出

```
Heap
PSYoungGen total 4800K, used 893K [0x1dac0000, 0x1e010000, 0x23010000)
eden space 4160K, 21% used [0x1dac0000,0x1db9f570,0x1ded0000)
from space 640K, 0% used [0x1df70000,0x1df70000,0x1e010000)
to space 640K, 0% used [0x1ded0000,0x1ded0000,0x1df70000)
ParOldGen total 19200K, used 16384K [0x13010000, 0x142d0000, 0x1dac0000)
object space 19200K, 85% used [0x13010000,0x14010020,0x142d0000)
PSPermGen total 12288K, used 2054K [0x0f010000, 0x0fc10000, 0x13010000)
object space 12288K, 16% used [0x0f010000,0x0f2119c0,0x0fc10000)

```

Show moreShow more icon

- CMS 收集器

与并行回收收集器不同，CMS 收集器主要关注于系统停顿时间。CMS 是 Concurrent Mark Sweep 的缩写，意为并发标记清除，从名称上可以得知，它使用的是标记-清除算法，同时它又是一个使用多线程并发回收的垃圾收集器。

CMS 工作时，主要步骤有：初始标记、并发标记、重新标记、并发清除和并发重置。其中初始标记和重新标记是独占系统资源的，而并发标记、并发清除和并发重置是可以和用户线程一起执行的。因此，从整体上来说，CMS 收集不是独占式的，它可以在应用程序运行过程中进行垃圾回收。

根据标记-清除算法，初始标记、并发标记和重新标记都是为了标记出需要回收的对象。并发清理则是在标记完成后，正式回收垃圾对象；并发重置是指在垃圾回收完成后，重新初始化 CMS 数据结构和数据，为下一次垃圾回收做好准备。并发标记、并发清理和并发重置都是可以和应用程序线程一起执行的。

CMS 收集器在其主要的工作阶段虽然没有暴力地彻底暂停应用程序线程，但是由于它和应用程序线程并发执行，相互抢占 CPU，所以在 CMS 执行期内对应用程序吞吐量造成一定影响。CMS 默认启动的线程数是 (ParallelGCThreads+3)/4),ParallelGCThreads 是新生代并行收集器的线程数，也可以通过-XX:ParallelCMSThreads 参数手工设定 CMS 的线程数量。当 CPU 资源比较紧张时，受到 CMS 收集器线程的影响，应用程序的性能在垃圾回收阶段可能会非常糟糕。

由于 CMS 收集器不是独占式的回收器，在 CMS 回收过程中，应用程序仍然在不停地工作。在应用程序工作过程中，又会不断地产生垃圾。这些新生成的垃圾在当前 CMS 回收过程中是无法清除的。同时，因为应用程序没有中断，所以在 CMS 回收过程中，还应该确保应用程序有足够的内存可用。因此，CMS 收集器不会等待堆内存饱和时才进行垃圾回收，而是当前堆内存使用率达到某一阈值时，便开始进行回收，以确保应用程序在 CMS 工作过程中依然有足够的空间支持应用程序运行。

这个回收阈值可以使用-XX:CMSInitiatingOccupancyFraction 来指定，默认是 68。即当老年代的空间使用率达到 68%时，会执行一次 CMS 回收。如果应用程序的内存使用率增长很快，在 CMS 的执行过程中，已经出现了内存不足的情况，此时，CMS 回收将会失败，JVM 将启动老年代串行收集器进行垃圾回收。如果这样，应用程序将完全中断，直到垃圾收集完成，这时，应用程序的停顿时间可能很长。因此，根据应用程序的特点，可以对-XX:CMSInitiatingOccupancyFraction 进行调优。如果内存增长缓慢，则可以设置一个稍大的值，大的阈值可以有效降低 CMS 的触发频率，减少老年代回收的次数可以较为明显地改善应用程序性能。反之，如果应用程序内存使用率增长很快，则应该降低这个阈值，以避免频繁触发老年代串行收集器。

标记-清除算法将会造成大量内存碎片，离散的可用空间无法分配较大的对象。在这种情况下，即使堆内存仍然有较大的剩余空间，也可能会被迫进行一次垃圾回收，以换取一块可用的连续内存，这种现象对系统性能是相当不利的，为了解决这个问题，CMS 收集器还提供了几个用于内存压缩整理的算法。

-XX:+UseCMSCompactAtFullCollection 参数可以使 CMS 在垃圾收集完成后，进行一次内存碎片整理。内存碎片的整理并不是并发进行的。-XX:CMSFullGCsBeforeCompaction 参数可以用于设定进行多少次 CMS 回收后，进行一次内存压缩。

-XX:CMSInitiatingOccupancyFraction 设置为 100，同时设置-XX:+UseCMSCompactAtFullCollection 和-XX:CMSFullGCsBeforeCompaction，日志输出如下：

##### 清单 14.CMS 垃圾回收器 GC 输出

```
[GC [DefNew: 825K->149K(4928K), 0.0023384 secs][Tenured: 8704K->661K(10944K),
0.0587725 secs] 9017K->661K(15872K),
[Perm : 374K->374K(12288K)], 0.0612037 secs] [Times: user=0.01 sys=0.02, real=0.06 secs]
Heap
def new generation total 4992K, used 179K [0x27010000, 0x27570000, 0x2c560000)
eden space 4480K, 4% used [0x27010000, 0x2703cda8, 0x27470000)
from space 512K, 0% used [0x27470000, 0x27470000, 0x274f0000)
to space 512K, 0% used [0x274f0000, 0x274f0000, 0x27570000)
tenured generation total 10944K, used 8853K [0x2c560000, 0x2d010000, 0x37010000)
the space 10944K, 80% used [0x2c560000, 0x2ce057c8, 0x2ce05800, 0x2d010000)
compacting perm gen total 12288K, used 374K [0x37010000, 0x37c10000, 0x3b010000)
the space 12288K, 3% used [0x37010000, 0x3706db10, 0x3706dc00, 0x37c10000)
ro space 10240K, 51% used [0x3b010000, 0x3b543000, 0x3b543000, 0x3ba10000)
rw space 12288K, 55% used [0x3ba10000, 0x3c0ae4f8, 0x3c0ae600, 0x3c610000)

```

Show moreShow more icon

- G1 收集器 (Garbage First)

G1 收集器的目标是作为一款服务器的垃圾收集器，因此，它在吞吐量和停顿控制上，预期要优于 CMS 收集器。

与 CMS 收集器相比，G1 收集器是基于标记-压缩算法的。因此，它不会产生空间碎片，也没有必要在收集完成后，进行一次独占式的碎片整理工作。G1 收集器还可以进行非常精确的停顿控制。它可以让开发人员指定当停顿时长为 M 时，垃圾回收时间不超过 N。使用参数-XX:+UnlockExperimentalVMOptions –XX:+UseG1GC 来启用 G1 回收器，设置 G1 回收器的目标停顿时间：-XX:MaxGCPauseMills=20,-XX:GCPauseIntervalMills=200。

## 收集器对系统性能的影响

通过清单 15 所示代码运行 1 万次循环，每次分配 512\*100B 空间，采用不同的垃圾回收器，输出程序运行所消耗的时间。

##### 清单 15\. 性能测试代码

```
import java.util.HashMap;

public class GCTimeTest {
static HashMap map = new HashMap();

public static void main(String[] args){
long begintime = System.currentTimeMillis();
for(int i=0;i<10000;i++){
if(map.size()*512/1024/1024>=400){
map.clear();//保护内存不溢出
System.out.println("clean map");
}
byte[] b1;
for(int j=0;j<100;j++){
b1 = new byte[512];
map.put(System.nanoTime(), b1);//不断消耗内存
}
}
long endtime = System.currentTimeMillis();
System.out.println(endtime-begintime);
}
}

```

Show moreShow more icon

使用参数-Xmx512M -Xms512M -XX:+UseParNewGC 运行代码，输出如下：

clean map 8565

cost time=1655

使用参数-Xmx512M -Xms512M -XX:+UseParallelOldGC –XX:ParallelGCThreads=8 运行代码，输出如下：

clean map 8798

cost time=1998

使用参数-Xmx512M -Xms512M -XX:+UseSerialGC 运行代码，输出如下：

clean map 8864

cost time=1717

使用参数-Xmx512M -Xms512M -XX:+UseConcMarkSweepGC 运行代码，输出如下：

clean map 8862

cost time=1530

上面例子的 GC 输出可以看出，采用不同的垃圾回收机制及设定不同的线程数，对于代码段的整体执行时间有较大的影响。需要读者有针对性地选用适合自己代码段的垃圾回收机制。

## GC 相关参数总结

(1) 与串行回收器相关的参数

-XX:+UseSerialGC:在新生代和老年代使用串行回收器。

-XX:+SurvivorRatio:设置 eden 区大小和 survivor 区大小的比例。

-XX:+PretenureSizeThreshold:设置大对象直接进入老年代的阈值。当对象的大小超过这个值时，将直接在老年代分配。

-XX:MaxTenuringThreshold:设置对象进入老年代的年龄的最大值。每一次 Minor GC 后，对象年龄就加 1。任何大于这个年龄的对象，一定会进入老年代。

(2) 与并行 GC 相关的参数

-XX:+UseParNewGC: 在新生代使用并行收集器。

-XX:+UseParallelOldGC: 老年代使用并行回收收集器。

-XX:ParallelGCThreads：设置用于垃圾回收的线程数。通常情况下可以和 CPU 数量相等。但在 CPU 数量比较多的情况下，设置相对较小的数值也是合理的。

-XX:MaxGCPauseMills：设置最大垃圾收集停顿时间。它的值是一个大于 0 的整数。收集器在工作时，会调整 Java 堆大小或者其他一些参数，尽可能地把停顿时间控制在 MaxGCPauseMills 以内。

-XX:GCTimeRatio:设置吞吐量大小，它的值是一个 0-100 之间的整数。假设 GCTimeRatio 的值为 n，那么系统将花费不超过 1/(1+n) 的时间用于垃圾收集。

-XX:+UseAdaptiveSizePolicy:打开自适应 GC 策略。在这种模式下，新生代的大小，eden 和 survivor 的比例、晋升老年代的对象年龄等参数会被自动调整，以达到在堆大小、吞吐量和停顿时间之间的平衡点。

(3) 与 CMS 回收器相关的参数

-XX:+UseConcMarkSweepGC: 新生代使用并行收集器，老年代使用 CMS+串行收集器。

-XX:+ParallelCMSThreads: 设定 CMS 的线程数量。

-XX:+CMSInitiatingOccupancyFraction:设置 CMS 收集器在老年代空间被使用多少后触发，默认为 68%。

-XX:+UseFullGCsBeforeCompaction:设定进行多少次 CMS 垃圾回收后，进行一次内存压缩。

-XX:+CMSClassUnloadingEnabled:允许对类元数据进行回收。

-XX:+CMSParallelRemarkEndable:启用并行重标记。

-XX:CMSInitatingPermOccupancyFraction:当永久区占用率达到这一百分比后，启动 CMS 回收 (前提是-XX:+CMSClassUnloadingEnabled 激活了)。

-XX:UseCMSInitatingOccupancyOnly:表示只在到达阈值的时候，才进行 CMS 回收。

-XX:+CMSIncrementalMode:使用增量模式，比较适合单 CPU。

(4) 与 G1 回收器相关的参数

-XX:+UseG1GC：使用 G1 回收器。

-XX:+UnlockExperimentalVMOptions:允许使用实验性参数。

-XX:+MaxGCPauseMills:设置最大垃圾收集停顿时间。

-XX:+GCPauseIntervalMills:设置停顿间隔时间。

(5) 其他参数

-XX:+DisableExplicitGC: 禁用显示 GC。

## 结束语

通过本文的学习，读者可以掌握基本的 JVM 垃圾回收器设计原理及使用规范。基于笔者多年的工作经验，没有哪一条优化是可以照本宣科的，它一定是基于您对 JVM 垃圾回收器工作原理及自身程序设计有一定了解前提下，通过大量的实验来找出最适合自己的优化方案。