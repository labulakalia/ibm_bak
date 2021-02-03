# Java 15 新特性概述
快速了解 Java 15 正式版本

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/the-new-features-of-java-15/)

李林锋

发布: 2021-01-07

* * *

Java 15 如期于 2020 年 9 月 15 日正式发布，此次更新是继半年前 Java 14 这一大版本发布之后的又一次常规版本更新，自 2017 年发布 Java 9 以来，Java 版本发布基本上都是按照敏捷的开发节奏来发布，由 2017 年之前的每三年一个版本，转变为现在的每半年一个版本，并且一直保持着。在这 2020 年发布的第二个版本版中，主要带来了 ZGC 增强、instanceof 增强、Shenandoah 回收器更新等方面的改动、增强和新特性。

本文主要介绍 Java 15 中的主要新特性，跟大家一起快速了解、学习 Java 15 带来了哪些不一样的体验和变化。

## Edwards-Curve 数字签名算法

Java 15 中引入了 Edwards-Curve 数据签名算法（也称为爱德华曲线算法，EdDSA），EdDSA 是一种实现了 RFC 8032 标准化方案的椭圆曲线数字加密签名算法，相比 Java 版本中现有的 ECDSA 方案的实现，EdDAS 具有更好的性能和更高的安全性，并且已经在一些加密库中有相关实现和支持，如：OpenSSL、BoringSSL 等。

Java 15 中的爱德华曲线算法，主要是通过添加新的 Signature、KeyFactory 和 KeyPairGenerator 等类或方法来实现的，具体用法如下：

##### 清单 1\. 生成密钥并签名

```
KeyPairGenerator kpg = KeyPairGenerator.getInstance("Ed25519");
KeyPair kp = kpg.generateKeyPair();
// algorithm is pure Ed25519
Signature sig = Signature.getInstance("Ed25519");
sig.initSign(kp.getPrivate());
sig.update(msg);
byte[] s = sig.sign();

```

Show moreShow more icon

上述代码中，先使用 Ed25519 机制生成密钥，并使用生成的密钥对消息进行签名。其中前面提到的 Ed25519，是一种基于爱德华曲线构建的 EdDSA 签名机制，其具有安全高效的特点，并且据实验室数据，在 128bit 安全强度条件下，Ed25519 签名机制能够达到 10 万/秒的签名速度和 7 万/秒的验签速度。

##### 清单 2\. 生成公钥

```
KeyFactory kf = KeyFactory.getInstance("EdDSA");
boolean xOdd = ...
BigInteger y = ...
NamedParameterSpec paramSpec = new NamedParameterSpec("Ed25519");
EdECPublicKeySpec pubSpec = new EdECPublicKeySpec(paramSpec, new EdPoint(xOdd, y));
PublicKey pubKey = kf.generatePublic(pubSpec);

```

Show moreShow more icon

上面代码，使用 EdDSA 算法获取密钥 Factory，然后根据 Ed25519 机制生成公钥。

## 封闭类（预览）

Java 15 引入了封闭类的概念来增强 Java 的语法特性，封闭类或接口主要能够限制被哪些类或接口扩展或实现，以此来达到开发人员对所定义的类或接口被继承或实现的管控，是一种比访问修饰符更具申明性的限制类使用的方式。

有了这一个特性，意味着以后遇到一个类，不是随便想继承就能够继承，想实现就能够实现了，需要得到定义者的允许才能行。

封闭类对其允许的子类有下面三种约束：

1. 封闭类及被其允许的子类需同属一模块，并且如果声明在未命名的模块中，则需要属于同一包；
2. 每个被其允许的子类都需要直接扩展密封的类；
3. 每个被其允许的子类都需要选择一个修饰符来描述其是如何继承自其父类发起的密封性。

##### 清单 3\. 封闭类举例

```
package com.example.geometry;

public abstract sealed class Shape
    permits Circle, Rectangle, Square {...}

public final class Circle extends Shape {...}

public sealed class Rectangle extends Shape
    permits TransparentRectangle, FilledRectangle {...}
public final class TransparentRectangle extends Rectangle {...}
public final class FilledRectangle extends Rectangle {...}

public non-sealed class Square extends Shape {...}

```

Show moreShow more icon

如上图所示，每个被其允许的子类，只能使用下面三种修饰符 final、sealed 和 non-sealed 中的一种，一个类不能既定义为 sealed（封闭子类）又定义为 final（隐含无子类），只能使用上述一种修饰符修饰。

同时，因为新增加了封闭类，java.lang.Class 类中将新增下面两个 public 方法：

- java.lang.constant.ClassDesc[] getPermittedSubclasses()
- boolean isSealed()

方法 getPermittedSubclasses() 返回一个包含 java.lang.constant.ClassDesc 对象的数组，如果该对象是封闭类的，则表示该类的所有允许的子类；如果是非封闭类，则返回一个空数组。

如果给定的类或接口是密封的，则 isSealed() 方法返回 true，否则返回 false；

## 禁用、弃用偏向锁

由于将偏向锁用于传统同步优化，维护成本比较高，Java 15 开始禁用、弃用了一直以来使用的偏向锁。所谓偏向锁就是 HotSpot 虚拟机中的一种同步优化技术，主要用于减少线程之间锁竞争时的同步开销，其主要目标是避免线程竞争锁资源时需要执行 CAS 操作的性能消耗。轻量级锁每次申请、释放锁都至少需要一次 CAS，但偏向锁只有初始化时需要进行一次 CAS 操作。

偏向锁假定只有第一个申请锁的线程使用锁（后面其他的线程不会再来申请锁），因此，只需要在 Mark Word 中 CAS 记录 owner（本质上也是更新，只是初始值为空），如果记录成功，则偏向锁获取成功，记录锁状态为偏向锁，当前线程等于 owner 就可以零成本的直接获得锁；否则，说明有其他线程竞争，将偏向锁撤销，上升为轻量级锁。

处理过多线程同步问题的技术人员，一般都会接触锁，按照类型，从宏观上，可以分为悲观锁和乐观锁。

乐观锁，是一种乐观思想的体现，认为使用场景中主要读操作比较多，而写操作会比较少，发生并发写的可能性比较低，每次去读数据的时候都认为数据不会被修改，所以不会上锁，但是在更新的时候会判断一下在此期间别人有没有去更新这个数据，采取在写时先读出当前版本号，然后加锁操作（比较跟上一次的版本号，如果一样则更新），如果失败则要重复读-比较-写的操作。

Java 中的乐观锁基本都是通过 CAS 操作实现的，CAS 是一种更新的原子操作，比较当前值跟传入值是否一样，一样则更新，否则失败。

悲观锁，就是一种悲观思想，即认为写操作较多，遇到并发写操作的可能性比较高，每次去读取数据的时候都认为数据会被修改，所以每次在读、写数据的时候都会加上锁，这样别人想读写这个数据就会 block 直到获取到锁。Java 语言中的悲观锁就是 synchronized，而 AQS 框架（全称 AbstractQueuedSynchronizer，抽象队列同步 ）下的锁则是先尝试 CAS 乐观锁去获取锁，获取不到，才会转换为悲观锁，如 RetreenLock。

主要可以分为偏向锁、轻量级锁、自旋锁、重量级锁等。

Java 中的偏向锁 (Biased Locking) 是自 Java 6 开始引入的一项多线程优化技术。偏向锁，顾名思义，它会偏向于第一个访问锁的线程，如果程序在运行过程中，同步锁只有一个线程访问，不存在多线程争用的情况，则线程是不需要触发同步的，这种情况下，就会给线程加一个偏向锁。 如果在运行过程中，遇到了其他线程抢占锁，则持有偏向锁的线程会被挂起，JVM 会消除它身上的偏向锁，将锁恢复为标准的轻量级锁。

轻量级锁是由偏向锁升级而来的，偏向锁如果运行在只有一个线程进入同步块的情况下，当第二个线程加入到锁竞争的时候，偏向锁就会升级为轻量级锁；

而当多个线程竞争轻量锁时，等待轻量锁的线程不会阻塞，而是会一直自旋的对待，这里就是自旋锁，尝试获取锁的线程，在没有获得锁的时候，不被挂起，而转而去执行一个空循环，即自旋。在若干个自旋后，如果还没有获得锁，则才被挂起，获得锁，则执行代码。

自旋锁原理非常简单，如果持有锁的线程能在很短时间内释放锁资源，那么那些等待竞争锁的线程就不需要做内核态和用户态之间的切换进入阻塞挂起状态，它们只需要等一等（自旋），等持有锁的线程释放锁后即可立即获取锁，这样就能避免用户线程和内核切换的消耗。

重量级锁，最为常见的是 synchronized 关键字所代表的锁机制，synchronized 会导致争用不到锁的线程进入阻塞状态，所以说它是 Java 语言中一个重量级的同步操纵，也被称为重量级锁，在 Java 5 开始，为了解决 synchronized 重量级锁带来的性能问题，引入了轻量锁与偏向锁，默认启用了自旋锁，这三类都属于乐观锁。

碍于篇幅，只能对锁进行上述简短介绍，不同的锁有其不同的特点，每种锁只有在其特定的场景下，才会有出色的表现，Java 语言中没有哪种锁能够在所有情况下都能表现出非常出色的效率，引入这么多锁的原因就是为了应对不同情况下的性能问题。

## Record 类型（预览）

Java 14 富有建设性的将 Record 类型作为预览特性而引入，而 Java 15 中基于反馈做出改进之后，进一步发布为预览版本。Record 类型允许在代码中使用紧凑的语法形式来声明类，而这些类能够作为不可变数据类型的封装持有者。Record 这一特性主要用在特定领域的类上；与枚举类型一样，Record 类型是一种受限形式的类型，主要用于存储、保存数据，并且没有其它额外自定义行为的场景下。

在以往开发过程中，被当作数据载体的类对象，在正确声明定义过程中，通常需要编写大量的无实际业务、重复性质的代码，其中包括：构造函数、属性调用、访问以及 equals() 、hashCode()、toString()等方法，这里引入的 Record 类型，其效果有些类似 Lombok 的 @Data 注解、Kotlin 中的 data class，但是又不尽完全相同，它们的共同点都是类的部分或者全部可以直接在类头中定义、描述，并且这个类只用于存储数据而已。对于 Record 类型，具体可以用下面代码来说明：

##### 清单 4\. Record 类型定义

```
public record Person(String name, int age) {
    public static String address;

    public String getName() {
        return name;
    }
}

```

Show moreShow more icon

对上述代码进行编译，然后反编译之后可以看到如下结果：

##### 清单 5\. Record 类型反编译结果

```
public final class Person extends java.lang.Record {
    private final java.lang.String name;
    private final java.lang.String age;

    public Person(java.lang.String name, java.lang.String age) { /* compiled code */ }

    public java.lang.String getName() { /* compiled code */ }

    public java.lang.String toString() { /* compiled code */ }

    public final int hashCode() { /* compiled code */ }

    public final boolean equals(java.lang.Object o) { /* compiled code */ }

    public java.lang.String name() { /* compiled code */ }

    public java.lang.String age() { /* compiled code */ }
}

```

Show moreShow more icon

根据反编译结果，可以得出，当用 Record 来声明一个类时，该类将自动拥有下面特征：

1. 拥有一个构造方法
2. 获取成员属性值的方法：name()、age()
3. hashCode() 方法和 euqals() 方法
4. toString() 方法
5. 类对象和属性被 final 关键字修饰，不能被继承，类的示例属性也都被 final 修饰，不能再被赋值使用。
6. 还可以在 Record 声明的类中定义静态属性、方法和示例方法。注意，不能在 Record 声明的类中定义示例字段，类也不能声明为抽象类等。

可以看到，该预览特性提供了一种更为紧凑的语法来声明类，并且可以大幅减少定义类似数据类型时所需的重复性代码。

引入的 Record 这种新的类型，将在 java.lang.Class 中引入了下面两个新方法：

##### 清单 6\. Record 新引入至 Class 中的方法

- RecordComponent[] getRecordComponents()
- boolean isRecord()

其中 getRecordComponents() 方法返回一组 java.lang.reflect.RecordComponent 对象组成的数组，java.lang.reflect.RecordComponent 也是一个新引入类，该数组的元素与 Record 类中的组件相对应，其顺序与在记录声明中出现的顺序相同，可以从该数组中的每个 RecordComponent 中提取到组件信息，包括其名称、类型、泛型类型、注释及其访问方法。

而 isRecord() 方法，则返回所在类是否是 Record 类型，如果是，则返回 true。

## ZGC: 一种可伸缩、低延迟的垃圾回收器（正式版）

ZGC 是最初在 Java 11 中引入，同时在后续几个 Java 版本中，根据收到的积极反馈，不断进行改进的一款基于内存 Region，同时使用了内存读屏障、染色指针和内存多重映射等技术，并且以可伸缩、低延迟为目标的内存垃圾回收器，不过在 Java 15 之前版本中，都是以预览特性来介绍的。

此次 Java 15 版本中，ZGC 终于可以作为一个可以正式使用的垃圾回收器使用。

在 Java 15 版本中只需要使用命令行便可以开启 ZGC 功能，具体需要添加如下 JVM 参数：

`-XX:+UseZGC`

相比 Java 14 中，不再需要使用参数 `-XX:+UnlockExperimentalVMOptions`。

## 文本块（正式版）

自 Java 13 引入了文本块来解决多行文本的问题，文本块主要以三重双引号开头，并同样以三重双引号结尾终止，它们之间的任何内容都被解释为文本块字符串的一部分，包括换行符，避免了对大多数转义序列的需要，并且它仍然是普通的 java.lang.String 对象。文本块可以在 Java 中能够使用字符串的任何地方进行使用，而与编译后的代码没有区别，还增强了 Java 程序中的字符串可读性。并且通过这种方式，可以更直观地表示字符串，可以支持跨越多行，而且不会出现转义的视觉混乱，将可以广泛提高 Java 类程序的可读性和可写性。Java 14 在 Java 13 引入的文本块的基础之上，新加入了两个转义符以解决新特性预览过程中反馈的一些问题，而 Java 15 中，该特性终于变成正式版可以正常使用。

在 Java 14 版本中，新加入的两个转义符，分别是：\ 和 \\s，这两个转义符分别表达涵义如下：

- \ ：行终止符，主要用于阻止插入换行符；
- \\s：表示一个空格。可以用来避免末尾的白字符被去掉。

在 Java 13 之前，多行字符串写法为：

##### 清单 7\. 多行字符串写法

```
String literal = "Lorem ipsum dolor sit amet, consectetur adipiscing " +
                 "elit, sed do eiusmod tempor incididunt ut labore " +
                 "et dolore magna aliqua.";

```

Show moreShow more icon

而在 Java 15 之后，上述内容可以用下面方式来表述：

##### 清单 8\. 多行文本块加上转义符的写法

```
String text = """
              Lorem ipsum dolor sit amet, consectetur adipiscing \
              elit, sed do eiusmod tempor incididunt ut labore \
              et dolore magna aliqua.\
              """;

```

Show moreShow more icon

上述两种写法，text 实际还是只有一行内容。

对于转义符：\\s，用法如下，能够保证下列文本每行正好都是六个字符长度：

##### 清单 9\. 多行文本块加上转义符的写法

```
String colors = """
             red  \s
            green\s
             blue \s
             """;

```

Show moreShow more icon

同时，使用这两个转义符，能够简化跨多行字符串编码问题，通过转义符，能够避免对换行等特殊字符串进行转移，从而简化代码编写，同时也增强了使用 String 来表达 HTML、XML、SQL 或 JSON 等格式字符串的编码可读性，且易于维护。

同时 Java 15 中对文本块的改进，还包括对 String 进行了方法扩展：

- stripIndent()：用于从文本块中去除空白字符
- translateEscapes()：用于翻译转义字符
- formatted(Object… args)：用于格式化

## Shenandoah：一种低停顿的垃圾回收器（正式版）

Shenandoah 垃圾收集器与 ZGC 非常类似，最初都是在 Java 11 中引入，在后续几个 Java 版本中作为预览特性引入，在 Red Hat 系统上面已经经过广泛测试，并根据收到的积极反馈，不断进行改进的一款内存垃圾回收器。此次 Java 15 版本中，这款内存回收器，终于作为正式版本发布，可以正式使用。

Shenandoah 是一款并发的垃圾收集器，是跟 ZGC 一样的另外一款面向低停顿的垃圾收集器，不过不同之处在于 ZGC 是基于内存读屏障、染色指针和内存多重映射等技术来实现，而 Shenandoah 回收器是基于 BrooksForwarding Pointer 技术来实现。

Shenandoah 主要目标是实现虚拟机 99.9%停顿小于 10ms，并且暂停时间与堆大小无关。

Shenandoah 回收器的主要有下面九个工作阶段：

- Init Mark 并发标记的初始化阶段，它为并发标记准备堆和应用线程，然后扫描根节点集合。这是整个回收器工作生命周期第一次停顿，这个阶段主要工作是根节点集合扫描，所以停顿时间主要取决于根节点集合的大小。
- Concurrent Marking 并发标记阶段，贯穿整个堆，以根节点集合为起点，跟踪可达的所有对象。这个阶段和应用程序一起运行，即并发（concurrent）。这个阶段的持续时间主要取决于存活对象的数量，以及堆中对象图的结构。在这个阶段，应用依然可以分配新的数据，所以在并发标记阶段，堆内存使用率会上升。
- Final Mark 标记阶段，清空所有待处理的标记/更新队列，重新扫描根节点集合，结束并发标记。这个阶段还会搞明白需要被清理（evacuated）的区域（即垃圾收集集合），并且通常为下一阶段做准备。最终标记是整个内存回收周期的第二个停顿阶段，这个阶段的部分工作能在并发预清理阶段完成，这个阶段最耗时的还是清空队列和扫描根节点集合。
- Concurrent Cleanup 并发回收阶段，回收垃圾区域，这些区域是指并发标记后，探测不到任何存活的对象的区域。
- Concurrent Evacuation 并发整理阶段，从垃圾回收集合中拷贝存活的对象到其他的区域中，这里是 Shenandoah 垃圾回收器有别于其他内存回收器关键点。这个阶段能和应用程序一起运行，所以应用程序依然可以继续分配内存，这个阶段持续时间主要取决于选中的垃圾收集集合大小（比如整个堆划分 128 个区域，如果有 16 个区域被选中，其耗时肯定超过 8 个区域被选中）。
- Init Update Refs 初始化更新引用阶段，它除了确保所有内存回收线程和应用线程已经完成并发整理阶段，以及为下一阶段内存回收做准备以外，其他什么都没有做。这是整个内存回收周期中，第三次停顿，也是时间最短的一次。
- Concurrent Update References 并发更新引用阶段，再次遍历整个堆，更新那些在并发整理阶段被移动的对象的引用。这也是有别于其他 GC 特性之处，这个阶段持续时间主要取决于堆中对象的数量，和对象图结构无关，因为这个过程是线性扫描堆。这个阶段是和应用一起并发运行的。
- Final Update Refs 通过再次更新现有的根节点集合完成更新引用阶段，这个阶段会回收内存收集集合中的区域，因为现在的堆已经没有对这些区域中的对象的引用。这是整个内存回收周期最后一个阶段，它的持续时间主要取决于根节点集合的大小。
- Concurrent Cleanup 并发清理阶段，回收那些现在没有任何引用的区域集合。

在 Java 15 版本中只需要使用命令行便可以开启 Shenandoah 垃圾收集器的功能，具体需要添加如下 JVM 参数：

`-XX:+UseShenandoahGC`

不用再添加参数 `-XX:+UnlockExperimentalVMOptions`。

## 结束语

Java 15 版本的发布带来了一些新特性、功能实用性的增强和实验预览功能的转正，预览特性的转正，更好的方便大家在以后的学习、生活工作中能够正式使用这些新特性。碍于篇幅原因，本文仅针对其中对使用人员影响较大的以及其中的一些主要的特性做了学习和介绍，如有兴趣，您还可以自行下载相关代码，继续深入研究。