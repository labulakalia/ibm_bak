# java.util.concurrent，第 2 部分
并发编程意味着更智慧地工作，而不是更困难地工作

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-5things5/)

Ted Neward, Alex Theedom

发布: 2010-07-09

* * *

并发 Collections 提供了线程安全、经过良好调优的数据结构，简化了并发编程。然而，在一些情形下，开发人员需要更进一步，思考如何调节和/或限制线程执行。由于 `java.util.concurrent` 的总体目标是简化多线程编程，您可能希望该包包含同步实用程序，而它确实包含。

本文是 [第 1 部分](/zh/articles/j-5things4/) 的延续，将介绍几个比核心语言原语（监视器）更高级的同步结构，但它们还未包含在 Collection 类中。一旦您了解了这些锁和门的用途，使用它们将非常直观。

##### 关于本系列

您觉得自己懂 Java 编程？事实是，大多数开发人员都只领会到了 Java 平台的皮毛，所学也只够应付工作。在本 [本系列](/zh/series/5-things-you-didnt-know-about/) 中，Ted Neward 深度挖掘 Java 平台的核心功能，揭示一些鲜为人知的事实，帮助您解决最棘手的编程困难。

## 1\. Semaphore

在一些企业系统中，开发人员经常需要限制未处理的特定资源请求（线程/操作）数量，事实上，限制有时候能够提高系统的吞吐量，因为它们减少了对特定资源的争用。尽管完全可以手动编写限制代码，但使用 Semaphore 类可以更轻松地完成此任务，它将帮您执行限制，如清单 1 所示：

##### 清单 1\. 使用 Semaphore 执行限制

```
import java.util.*;import java.util.concurrent.*;

public class SemApp
{
    public static void main(String[] args)
    {
        Runnable limitedCall = new Runnable() {
            final Random rand = new Random();
            final Semaphore available = new Semaphore(3);
            int count = 0;
            public void run()
            {
                int time = rand.nextInt(15);
                int num = count++;

                try
                {
                    available.acquire();

                    System.out.println("Executing " +
                        "long-running action for " +
                        time + " seconds... #" + num);

                    Thread.sleep(time * 1000);

                    System.out.println("Done with #" +
                        num + "!");

                    available.release();
                }
                catch (InterruptedException intEx)
                {
                    intEx.printStackTrace();
                }
            }
        };

        for (int i=0; i<10; i++)
            new Thread(limitedCall).start();
    }
}

```

Show moreShow more icon

即使本例中的 10 个线程都在运行（您可以对运行 `SemApp` 的 Java 进程执行 `jstack` 来验证），但只有 3 个线程是活跃的。在一个信号计数器释放之前，其他 7 个线程都处于空闲状态。（实际上， `Semaphore` 类支持一次获取和释放多个 _permit_ ，但这不适用于本场景。）

## 2\. CountDownLatch

如果 `Semaphore` 是允许一次进入一个（这可能会勾起一些流行夜总会的保安的记忆）线程的并发性类，那么 `CountDownLatch` 就像是赛马场的起跑门栅。此类持有所有空闲线程，直到满足特定条件，这时它将会一次释放所有这些线程。

##### 清单 2\. CountDownLatch：让我们去赛马吧！

```
import java.util.*;
import java.util.concurrent.*;

class Race
{
    private Random rand = new Random();

    private int distance = rand.nextInt(250);
    private CountDownLatch start;
    private CountDownLatch finish;

    private List<String> horses = new ArrayList<String>();

    public Race(String... names)
    {
        this.horses.addAll(Arrays.asList(names));
    }

    public void run()
        throws InterruptedException
    {
        System.out.println("And the horses are stepping up to the gate...");
        final CountDownLatch start = new CountDownLatch(1);
        final CountDownLatch finish = new CountDownLatch(horses.size());
        final List<String> places =
            Collections.synchronizedList(new ArrayList<String>());

        for (final String h : horses)
        {
            new Thread(new Runnable() {
                public void run() {
                    try
                    {
                        System.out.println(h +
                            " stepping up to the gate...");
                        start.await();

                        int traveled = 0;
                        while (traveled < distance)
                        {
                            // In a 0-2 second period of time....
                            Thread.sleep(rand.nextInt(3) * 1000);

                            // ... a horse travels 0-14 lengths
                            traveled += rand.nextInt(15);
                            System.out.println(h +
                                " advanced to " + traveled + "!");
                        }
                        finish.countDown();
                        System.out.println(h +
                            " crossed the finish!");
                        places.add(h);
                    }
                    catch (InterruptedException intEx)
                    {
                        System.out.println("ABORTING RACE!!!");
                        intEx.printStackTrace();
                    }
                }
            }).start();
        }

        System.out.println("And... they're off!");
        start.countDown();

        finish.await();
        System.out.println("And we have our winners!");
        System.out.println(places.get(0) + " took the gold...");
        System.out.println(places.get(1) + " got the silver...");
        System.out.println("and " + places.get(2) + " took home the bronze.");
    }
}

public class CDLApp
{
    public static void main(String[] args)
        throws InterruptedException, java.io.IOException
    {
        System.out.println("Prepping...");

        Race r = new Race(
            "Beverly Takes a Bath",
            "RockerHorse",
            "Phineas",
            "Ferb",
            "Tin Cup",
            "I'm Faster Than a Monkey",
            "Glue Factory Reject"
            );

        System.out.println("It's a race of " + r.getDistance() + " lengths");

        System.out.println("Press Enter to run the race....");
        System.in.read();

        r.run();
    }
}

```

Show moreShow more icon

注意，在 [清单 2](#清单-2-countdownlatch：让我们去赛马吧！) 中， `CountDownLatch` 有两个用途：首先，它同时释放所有线程，模拟马赛的起点，但随后会设置一个门闩模拟马赛的终点。这样，”主” 线程就可以输出结果。 为了让马赛有更多的输出注释，可以在赛场的 “转弯处” 和 “半程” 点，比如赛马跨过跑道的四分之一、二分之一和四分之三线时，添加 `CountDownLatch` 。

## 3\. Executor

[清单 1](#清单-1-使用-semaphore-执行限制) 和 [清单 2](#清单-2-countdownlatch：让我们去赛马吧！) 中的示例都存在一个重要的缺陷，它们要求您直接创建 `Thread` 对象。这可以解决一些问题，因为在一些 JVM 中，创建 `Thread` 是一项重量型的操作，重用现有 `Thread` 比创建新线程要容易得多。而在另一些 JVM 中，情况正好相反： `Thread` 是轻量型的，可以在需要时很容易地新建一个线程。当然，如果 Murphy 拥有自己的解决办法（他通常都会拥有），那么您无论使用哪种方法对于您最终将部署的平台都是不对的。

JSR-166 专家组在一定程度上预测到了这一情形。Java 开发人员无需直接创建 `Thread` ，他们引入了 `Executor` 接口，这是对创建新线程的一种抽象。如清单 3 所示， `Executor` 使您不必亲自对 `Thread` 对象执行 `new` 就能够创建新线程：

##### 清单 3\. Executor

```
Executor exec = getAnExecutorFromSomeplace();
exec.execute(new Runnable() { ... });

```

Show moreShow more icon

使用 `Executor` 的主要缺陷与我们在所有工厂中遇到的一样：工厂必须来自某个位置。不幸的是，与 CLR 不同，JVM 没有附带一个标准的 VM 级线程池。

`Executor` 类 _实际上_ 充当着一个提供 `Executor` 实现实例的共同位置，但它只有 `new` 方法（例如用于创建新线程池）；它没有预先创建实例。所以您可以自行决定是否希望在代码中创建和使用 `Executor` 实例。（或者在某些情况下，您将能够使用所选的容器/平台提供的实例。）

### ExecutorService 随时可以使用

尽管不必担心 `Thread` 来自何处，但 `Executor` 接口缺乏 Java 开发人员可能期望的某种功能，比如结束一个用于生成结果的线程并以非阻塞方式等待结果可用。（这是桌面应用程序的一个常见需求，用户将执行需要访问数据库的 UI 操作，然后如果该操作花费了很长时间，可能希望在它完成之前取消它。）

对于此问题，JSR-166 专家创建了一个更加有用的抽象（ `ExecutorService` 接口），它将线程启动工厂建模为一个可集中控制的服务。例如，无需每执行一项任务就调用一次 `execute()` ， `ExecutorService` 可以接受一组任务并返回一个表示每项任务的未来结果的 _未来列表_ 。

## 4\. ScheduledExecutorServices

尽管 `ExecutorService` 接口非常有用，但某些任务仍需要以计划方式执行，比如以确定的时间间隔或在特定时间执行给定的任务。这就是 `ScheduledExecutorService` 的应用范围，它扩展了 `ExecutorService` 。

如果您的目标是创建一个每隔 5 秒跳一次的 “心跳” 命令，使用 `ScheduledExecutorService` 可以轻松实现，如清单 4 所示：

##### 清单 4\. ScheduledExecutorService 模拟心跳

```
import java.util.concurrent.*;

public class Ping
{
    public static void main(String[] args)
    {
        ScheduledExecutorService ses =
            Executors.newScheduledThreadPool(1);
        Runnable pinger = new Runnable() {
            public void run() {
                System.out.println("PING!");
            }
        };
        ses.scheduleAtFixedRate(pinger, 5, 5, TimeUnit.SECONDS);
    }
}

```

Show moreShow more icon

这项功能怎么样？不用过于担心线程，不用过于担心用户希望取消心跳时会发生什么，也不用明确地将线程标记为前台或后台；只需将所有的计划细节留给 `ScheduledExecutorService` 。

顺便说一下，如果用户希望取消心跳， `scheduleAtFixedRate` 调用将返回一个 `ScheduledFuture` 实例，它不仅封装了结果（如果有），还拥有一个 `cancel` 方法来关闭计划的操作。

## 5\. Timeout 方法

为阻塞操作设置一个具体的超时值（以避免死锁）的能力是 `java.util.concurrent` 库相比起早期并发特性的一大进步，比如监控锁定。

这些方法几乎总是包含一个 `int`/`TimeUnit` 对，指示这些方法应该等待多长时间才释放控制权并将其返回给程序。它需要开发人员执行更多工作 — 如果没有获取锁，您将如何重新获取？ — 但结果几乎总是正确的：更少的死锁和更加适合生产的代码。（关于编写生产就绪代码的更多信息，请参见 参考资料 中 Michael Nygard 编写的 _Release It!_ 。）

## 结束语

`java.util.concurrent` 包还包含了其他许多好用的实用程序，它们很好地扩展到了 Collections 之外，尤其是在 `.locks` 和 `.atomic` 包中。深入研究，您还将发现一些有用的控制结构，比如 `CyclicBarrier` 等。

与 Java 平台的许多其他方面一样，您无需费劲地查找可能非常有用的基础架构代码。在编写多线程代码时，请记住本文讨论的实用程序和 [上一篇文章](/zh/articles/j-5things4/) 中讨论的实用程序。

下一次，我们将进入一个新的主题： [关于 Jars 您不知道的 5 件事](/zh/articles/j-5things6/)。

本文翻译自： [java.util.concurrent, Part 2](https://developer.ibm.com/articles/j-5things5/)（2010-06-01）