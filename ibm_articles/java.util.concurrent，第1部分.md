# java.util.concurrent，第 1 部分
通过并发 Collections 进行多线程编程

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-5things4/)

Ted Neward, Alex Theedom

更新: 2010-07-07 \| 发布: 2010-07-01

* * *

##### 关于本系列

您觉得自己懂 Java 编程？事实是，大多数开发人员都只领会到了 Java 平台的皮毛，所学也只够应付工作。在 [本系列](/zh/series/5-things-you-didnt-know-about/) 中，Ted Neward 深度挖掘 Java 平台的核心功能，揭示一些鲜为人知的事实，帮助您解决最棘手的编程困难。

Concurrent Collections 是 Java™ 5 的巨大附加产品，但是在关于注解和泛型的争执中很多 Java 开发人员忽视了它们。此外（或者更老实地说），许多开发人员避免使用这个数据包，因为他们认为它一定很复杂，就像它所要解决的问题一样。

事实上，`java.util.concurrent` 包含许多类，能够有效解决普通的并发问题，无需复杂工序。阅读本文，了解 `java.util.concurrent` 类，比如 `CopyOnWriteArrayList` 和 `BlockingQueue` 如何帮助您解决多线程编程的棘手问题。

## 1\. TimeUnit

尽管 _本质上_ 不是 Collections 类，但 `java.util.concurrent.TimeUnit` 枚举让代码更易读懂。使用 `TimeUnit` 将使用您的方法或 API 的开发人员从毫秒的 “暴政” 中解放出来。

`TimeUnit` 包括所有时间单位，从 `MILLISECONDS` 和 `MICROSECONDS` 到 `DAYS` 和 `HOURS` ，这就意味着它能够处理一个开发人员所需的几乎所有的时间范围类型。同时，因为在列举上声明了转换方法，在时间加快时，将 `HOURS` 转换回 `MILLISECONDS` 甚至变得更容易。

## 2\. CopyOnWriteArrayList

创建数组的全新副本是过于昂贵的操作，无论是从时间上，还是从内存开销上，因此在通常使用中很少考虑；开发人员往往求助于使用同步的 `ArrayList` 。然而，这也是一个成本较高的选择，因为每当您跨集合内容进行迭代时，您就不得不同步所有操作，包括读和写，以此保证一致性。

这又让成本结构回到这样一个场景：需多读者都在读取 `ArrayList` ，但是几乎没人会去修改它。

`CopyOnWriteArrayList` 是个巧妙的小宝贝，能解决这一问题。它的 Javadoc 将 `CopyOnWriteArrayList` 定义为一个 “`ArrayList` 的线程安全变体，在这个变体中所有易变操作（添加，设置等）可以通过复制全新的数组来实现”。

集合从内部将它的内容复制到一个没有修改的新数组，这样读者访问数组内容时就不会产生同步成本（因为他们从来不是在易变数据上操作）。

本质上讲， `CopyOnWriteArrayList` 很适合处理 `ArrayList` 经常让我们失败的这种场景：读取频繁，但很少有写操作的集合，例如 JavaBean 事件的 `Listener` s。

## 3\. BlockingQueue

`BlockingQueue` 接口表示它是一个 `Queue` ，意思是它的项以先入先出（FIFO）顺序存储。在特定顺序插入的项以相同的顺序检索 — 但是需要附加保证，从空队列检索一个项的任何尝试都会阻塞调用线程，直到这个项准备好被检索。同理，想要将一个项插入到满队列的尝试也会导致阻塞调用线程，直到队列的存储空间可用。

`BlockingQueue` 干净利落地解决了如何将一个线程收集的项”传递”给另一线程用于处理的问题，无需考虑同步问题。Java Tutorial 的 Guarded Blocks 试用版就是一个很好的例子。它构建一个单插槽绑定的缓存，当新的项可用，而且插槽也准备好接受新的项时，使用手动同步和 `wait()` / `notifyAll()` 在线程之间发信。（详见 [Guarded Blocks 实现](http://java.sun.com/docs/books/tutorial/essential/concurrency/guardmeth.html) 。）

尽管 Guarded Blocks 教程中的代码有效，但是它耗时久，混乱，而且也并非完全直观。退回到 Java 平台较早的时候，没错，Java 开发人员不得不纠缠于这种代码；但现在是 2010 年 — 情况难道没有改善？

清单 1 显示了 Guarded Blocks 代码的重写版，其中我使用了一个 `ArrayBlockingQueue` ，而不是手写的 `Drop` 。

##### 清单 1\. BlockingQueue

```
import java.util.*;
import java.util.concurrent.*;

class Producer
    implements Runnable
{
    private BlockingQueue<String> drop;
    List<String> messages = Arrays.asList(
        "Mares eat oats",
        "Does eat oats",
        "Little lambs eat ivy",
        "Wouldn't you eat ivy too?");

    public Producer(BlockingQueue<String> d) { this.drop = d; }

    public void run()
    {
        try
        {
            for (String s : messages)
                drop.put(s);
            drop.put("DONE");
        }
        catch (InterruptedException intEx)
        {
            System.out.println("Interrupted! " +
                "Last one out, turn out the lights!");
        }
    }
}

class Consumer
    implements Runnable
{
    private BlockingQueue<String> drop;
    public Consumer(BlockingQueue<String> d) { this.drop = d; }

    public void run()
    {
        try
        {
            String msg = null;
            while (!((msg = drop.take()).equals("DONE")))
                System.out.println(msg);
        }
        catch (InterruptedException intEx)
        {
            System.out.println("Interrupted! " +
                "Last one out, turn out the lights!");
        }
    }
}

public class ABQApp
{
    public static void main(String[] args)
    {
        BlockingQueue<String> drop = new ArrayBlockingQueue(1, true);
        (new Thread(new Producer(drop))).start();
        (new Thread(new Consumer(drop))).start();
    }
}

```

Show moreShow more icon

`ArrayBlockingQueue` 还体现了”公平” — 意思是它为读取器和编写器提供线程先入先出访问。这种替代方法是一个更有效，但又冒穷尽部分线程风险的政策。（即，允许一些读取器在其他读取器锁定时运行效率更高，但是您可能会有读取器线程的流持续不断的风险，导致编写器无法进行工作。）

##### 注意 Bug！

顺便说一句，如果您注意到 Guarded Blocks 包含一个重大 bug，那么您是对的  如果开发人员在 `main()` 中的 `Drop` 实例上同步，会出现什么情况呢？

`BlockingQueue` 还支持接收时间参数的方法，时间参数表明线程在返回信号故障以插入或者检索有关项之前需要阻塞的时间。这么做会避免非绑定的等待，这对一个生产系统是致命的，因为一个非绑定的等待会很容易导致需要重启的系统挂起。

## 4\. ConcurrentMap

`Map` 有一个微妙的并发 bug，这个 bug 将许多不知情的 Java 开发人员引入歧途。`ConcurrentMap` 是最容易的解决方案。

当一个 `Map` 被从多个线程访问时，通常使用 `containsKey()` 或者 `get()` 来查看给定键是否在存储键/值对之前出现。但是即使有一个同步的 `Map`，线程还是可以在这个过程中潜入，然后夺取对 `Map` 的控制权。问题是，在对 `put()` 的调用中，锁在 `get()` 开始时获取，然后在可以再次获取锁之前释放。它的结果是个竞争条件：这是两个线程之间的竞争，结果也会因谁先运行而不同。

如果两个线程几乎同时调用一个方法，两者都会进行测试，调用 put，在处理中丢失第一线程的值。幸运的是，`ConcurrentMap` 接口支持许多附加方法，它们设计用于在一个锁下进行两个任务：`putIfAbsent()`，例如，首先进行测试，然后仅当键没有存储在 `Map` 中时进行 put。

## 5\. SynchronousQueues

根据 Javadoc，`SynchronousQueue` 是个有趣的东西：

> \_这是一个阻塞队列，其中，每个插入操作必须等待另一个线程的对应移除操作，反之亦然。一个同步队列不具有任何内部容量，甚至不具有 1 的容量。

本质上讲，`SynchronousQueue` 是之前提过的 `BlockingQueue` 的又一实现。它给我们提供了在线程之间交换单一元素的极轻量级方法，使用 `ArrayBlockingQueue` 使用的阻塞语义。在清单 2 中，我重写了 [清单 1](#清单-1-blockingqueue) 的代码，使用 `SynchronousQueue` 替代 `ArrayBlockingQueue`：

##### 清单 2\. SynchronousQueue

```
import java.util.*;
import java.util.concurrent.*;

class Producer
    implements Runnable
{
    private BlockingQueue<String> drop;
    List<String> messages = Arrays.asList(
        "Mares eat oats",
        "Does eat oats",
        "Little lambs eat ivy",
        "Wouldn't you eat ivy too?");

    public Producer(BlockingQueue<String> d) { this.drop = d; }

    public void run()
    {
        try
        {
            for (String s : messages)
                drop.put(s);
            drop.put("DONE");
        }
        catch (InterruptedException intEx)
        {
            System.out.println("Interrupted! " +
                "Last one out, turn out the lights!");
        }
    }
}

class Consumer
    implements Runnable
{
    private BlockingQueue<String> drop;
    public Consumer(BlockingQueue<String> d) { this.drop = d; }

    public void run()
    {
        try
        {
            String msg = null;
            while (!((msg = drop.take()).equals("DONE")))
                System.out.println(msg);
        }
        catch (InterruptedException intEx)
        {
            System.out.println("Interrupted! " +
                "Last one out, turn out the lights!");
        }
    }
}

public class SynQApp
{
    public static void main(String[] args)
    {
        BlockingQueue<String> drop = new SynchronousQueue<String>();
        (new Thread(new Producer(drop))).start();
        (new Thread(new Consumer(drop))).start();
    }
}

```

Show moreShow more icon

实现代码看起来几乎相同，但是应用程序有额外获益：`SynchronousQueue` 允许在队列进行一个插入，只要有一个线程等着使用它。

在实践中，`SynchronousQueue` 类似于 Ada 和 CSP 等语言中可用的 “会合通道”。这些通道有时在其他环境中也称为 “连接”，这样的环境包括 .NET（见参考资源）。

## 结束语

当 Java 运行时知识库提供便利、预置的并发性时，为什么还要苦苦挣扎，试图将并发性导入到您的 Collections 类？本系列的 [下一篇文章](/zh/articles/j-5things5/))将会进一步探讨 `java.util.concurrent` 名称空间的内容。

本文翻译自： [java.util.concurrent, Part 1](https://developer.ibm.com/articles/j-5things4/)（2010-05-18）