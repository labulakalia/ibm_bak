# Java 常用工具
Java 常用工具，如解析、计时和声音

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-5things12/)

Ted Neward, Alex Theedom

发布: 2010-10-18

* * *

##### 关于本系列

您觉得自己懂 Java 编程？事实是，大多数开发人员都只领会到了 Java 平台的皮毛，所学也只够应付工作。在 [本系列](/zh/series/5-things-you-didnt-know-about/) 中，Ted Neward 深度挖掘 Java 平台的核心功能，揭示一些鲜为人知的事实，帮助您解决最棘手的编程困难。

很多年前，当我还是高中生的时候，我曾考虑以小说作家作为我的职业追求，我订阅了一本 _Writer’s Digest_ 杂志。我记得其中有篇专栏文章，是关于 “太小而难以保存的线头”，专栏作者描述厨房储物抽屉中放满了无法分类的玩意儿。这句话我一直铭记在心，它正好用来描述本文的内容，本系列的最后一篇（至少目前是这样）。

Java 平台就充满了这样的 “线头” — 有用的命令行工具和库，大多数 Java 开发人员甚至都不知道，更别提使用了。其中很多无法划分到之前的 [_5 件事_ 系列](/zh/series/5-things-you-didnt-know-about/) 的编程分类中，但不管怎样，尝试一下：有些说不定会在您编程的厨房抽屉中占得一席之地。

## 1\. StAX

在千禧年左右，当 XML 第一次出现在很多 Java 开发人员面前时，有两种基本的解析 XML 文件的方法。SAX 解析器实际是由程序员对事件调用一系列回调方法的大型状态机。DOM 解析器将整个 XML 文档加入内存，并切割成离散的对象，它们连接在一起形成一个树。该树描述了文档的整个 XML Infoset 表示法。这两个解析器都有缺点：SAX 太低级，无法使用，DOM 代价太大，尤其对于大的 XML 文件 — 整个树成了一个庞然大物。

幸运的是，Java 开发人员找到第三种方法来解析 XML 文件，通过对文档建模成 “节点”，它们可以从文档流中一次取出一个，检查，然后处理或丢弃。这些 “节点” 的 “流” 提供了 SAX 和 DOM 的中间地带，名为 “Streaming API for XML”，或者叫做StAX。（此缩写用于区分新的 API 与原来的 SAX 解析器，它与此同名。）StAX 解析器后来包装到了 JDK 中，在 `javax.xml.stream` 包。

使用 StAX 相当简单：实例化 `XMLEventReader`，将它指向一个格式良好的 XML 文件，然后一次 “拉出” 一个节点（通常用 `while` 循环），查看。例如，在清单 1 中，列举出了 Ant 构造脚本中的所有目标：

##### 清单 1\. 只是让 StAX 指向目标

```
import java.io.*;
import javax.xml.namespace.QName;
import javax.xml.stream.*;
import javax.xml.stream.events.*;
import javax.xml.stream.util.*;

public class Targets
{
    public static void main(String[] args)
        throws Exception
    {
        for (String arg : args)
        {
            XMLEventReader xsr =
                XMLInputFactory.newInstance()
                    .createXMLEventReader(new FileReader(arg));
            while (xsr.hasNext())
            {
                XMLEvent evt = xsr.nextEvent();
                switch (evt.getEventType())
                {
                    case XMLEvent.START_ELEMENT:
                    {
                        StartElement se = evt.asStartElement();
                        if (se.getName().getLocalPart().equals("target"))
                        {
                            Attribute targetName =
                                se.getAttributeByName(new QName("name"));
                            // Found a target!
                            System.out.println(targetName.getValue());
                        }
                        break;
                    }
                    // Ignore everything else
                }
            }
        }
    }
}

```

Show moreShow more icon

StAX 解析器不会替换所有的 SAX 和 DOM 代码。但肯定会让某些任务容易些。尤其对完成不需要知道 XML 文档整个树结构的任务相当方便。

请注意，如果事件对象级别太高，无法使用，StAX 也有一个低级 API 在 `XMLStreamReader` 中。尽管也许没有阅读器有用，StAX 还有一个 `XMLEventWriter` ，同样，还有一个 `XMLStreamWriter` 类用于 XML 输出。

## 2\. ServiceLoader

Java 开发人员经常希望将使用和创建组件的内容区分开来。这通常是通过创建一个描述组件动作的接口，并使用某种中介创建组件实例来完成的。很多开发人员使用 Spring 框架来完成，但还有其他的方法，它比 Spring 容器更轻量级。

`java.util` 的 `ServiceLoader` 类能读取隐藏在 JAR 文件中的配置文件，并找到接口的实现，然后使这些实现成为可选择的列表。例如，如果您需要一个私仆（personal-servant）组件来完成任务，您可以使用清单 2 中的代码来实现：

##### 清单 2\. IPersonalServant

```
public interface IPersonalServant
{
    // Process a file of commands to the servant
    public void process(java.io.File f)
        throws java.io.IOException;
    public boolean can(String command);
}

```

Show moreShow more icon

`can()` 方法可让您确定所提供的私仆实现是否满足需求。清单 3 中的 `ServiceLoader` 的 `IPersonalServant` 列表基本上满足需求：

##### 清单 3\. IPersonalServant 行吗？

```
import java.io.*;
import java.util.*;

public class Servant
{
    public static void main(String[] args)
        throws IOException
    {
        ServiceLoader<IPersonalServant> servantLoader =
            ServiceLoader.load(IPersonalServant.class);

        IPersonalServant i = null;
        for (IPersonalServant ii : servantLoader)
            if (ii.can("fetch tea"))
                i = ii;

        if (i == null)
            throw new IllegalArgumentException("No suitable servant found");

        for (String arg : args)
        {
            i.process(new File(arg));
        }
    }
}

```

Show moreShow more icon

假设有此接口的实现，如清单 4：

##### 清单 4\. Jeeves 实现了 IPersonalServant

```
import java.io.*;

public class Jeeves
    implements IPersonalServant
{
    public void process(File f)
    {
        System.out.println("Very good, sir.");
    }
    public boolean can(String cmd)
    {
        if (cmd.equals("fetch tea"))
            return true;
        else
            return false;
    }
}

```

Show moreShow more icon

剩下的就是配置包含实现的 JAR 文件，让 `ServiceLoader` 能识别 — 这可能会非常棘手。JDK 想要 JAR 文件有一个 `META-INF/services` 目录，它包含一个文本文件，其文件名与接口类名完全匹配 — 本例中是 `META-INF/services/IPersonalServant`。接口类名的内容是实现的名称，每行一个，如清单 5：

##### 清单 5\. META-INF/services/IPersonalServant

```
Jeeves   # comments are OK

```

Show moreShow more icon

幸运的是，Ant 构建系统（自 1.7.0 以来）包含一个对 `jar` 任务的服务标签，让这相对容易，见清单 6：

##### 清单 6\. Ant 构建的 IPersonalServant

```
<target name="serviceloader" depends="build">
        <jar destfile="misc.jar" basedir="./classes">
            <service type="IPersonalServant">
                <provider classname="Jeeves" />
            </service>
        </jar>
    </target>

```

Show moreShow more icon

这里，很容易调用 `IPersonalServant` ，让它执行命令。然而，解析和执行这些命令可能会非常棘手。这又是另一个 “小线头”。

## 3\. Scanner

有无数 Java 工具能帮助您构建解析器，很多函数语言已成功构建解析器函数库（解析器选择器）。但如果要解析的是逗号分隔值文件，或空格分隔文本文件，又怎么办呢？大多数工具用在此处就过于隆重了，而 `String.split()` 又不够。（对于正则表达式，请记住一句老话：” 您有一个问题，用正则表达式解决。那您就有两个问题了。”）

Java 平台的 `Scanner` 类会是这些类中您最好的选择。以轻量级文本解析器为目标，`Scanner` 提供了一个相对简单的 API，用于提取结构化文本，并放入强类型的部分。想象一下，如果您愿意，一组类似 DSL 的命令（源自 Terry Pratchett _Discworld_ 小说）排列在文本文件中，如清单 7：

##### 清单 7\. Igor 的任务

```
fetch 1 head
fetch 3 eye
fetch 1 foot
attach foot to head
attach eye to head
admire

```

Show moreShow more icon

您，或者是本例中称为 `Igor` 的私仆，能轻松使用 `Scanner` 解析这组违法命令，如清单 8 所示：

##### 清单 8\. Igor 的任务，由 Scanner 解析

```
import java.io.*;
import java.util.*;

public class Igor
    implements IPersonalServant
{
    public boolean can(String cmd)
    {
        if (cmd.equals("fetch body parts"))
            return true;
        if (cmd.equals("attach body parts"))
            return true;
        else
            return false;
    }
    public void process(File commandFile)
        throws FileNotFoundException
    {
        Scanner scanner = new Scanner(commandFile);
        // Commands come in a verb/number/noun or verb form
        while (scanner.hasNext())
        {
            String verb = scanner.next();
            if (verb.equals("fetch"))
            {
                int num = scanner.nextInt();
                String type = scanner.next();
                fetch (num, type);
            }
            else if (verb.equals("attach"))
            {
                String item = scanner.next();
                String to = scanner.next();
                String target = scanner.next();
                attach(item, target);
            }
            else if (verb.equals("admire"))
            {
                admire();
            }
            else
            {
                System.out.println("I don't know how to "
                    + verb + ", marthter.");
            }
        }
    }

    public void fetch(int number, String type)
    {
        if (parts.get(type) == null)
        {
            System.out.println("Fetching " + number + " "
                + type + (number > 1 ? "s" : "") + ", marthter!");
            parts.put(type, number);
        }
        else
        {
            System.out.println("Fetching " + number + " more "
                + type + (number > 1 ? "s" : "") + ", marthter!");
            Integer currentTotal = parts.get(type);
            parts.put(type, currentTotal + number);
        }
        System.out.println("We now have " + parts.toString());
    }

    public void attach(String item, String target)
    {
        System.out.println("Attaching the " + item + " to the " +
            target + ", marthter!");
    }

    public void admire()
    {
        System.out.println("It'th quite the creathion, marthter");
    }

    private Map<String, Integer> parts = new HashMap<String, Integer>();
}

```

Show moreShow more icon

假设 `Igor` 已在 `ServantLoader` 中注册，可以很方便地将 `can()` 调用改得更实用，并重用前面的 `Servant` 代码，如清单 9 所示：

##### 清单 9\. Igor 做了什么

```
import java.io.*;
import java.util.*;

public class Servant
{
    public static void main(String[] args)
        throws IOException
    {
        ServiceLoader<IPersonalServant> servantLoader =
            ServiceLoader.load(IPersonalServant.class);

        IPersonalServant i = null;
        for (IPersonalServant ii : servantLoader)
            if (ii.can("fetch body parts"))
                i = ii;

        if (i == null)
            throw new IllegalArgumentException("No suitable servant found");

        for (String arg : args)
        {
            i.process(new File(arg));
        }
    }
}

```

Show moreShow more icon

真正 DSL 实现显然不会仅仅打印到标准输出流。我把追踪哪些部分、跟随哪些部分的细节留待给您（当然，还有忠诚的 `Igor` ）。

## 4\. Timer

`java.util.Timer` 和 `TimerTask` 类提供了方便、相对简单的方法可在定期或一次性延迟的基础上执行任务：

##### 清单 10\. 稍后执行

```
import java.util.*;

public class Later
{
    public static void main(String[] args)
    {
        Timer t = new Timer("TimerThread");
        t.schedule(new TimerTask() {
            public void run() {
                System.out.println("This is later");
                System.exit(0);
            }
        }, 1 * 1000);
        System.out.println("Exiting main()");
    }
}

```

Show moreShow more icon

`Timer` 有许多 `schedule()` 重载，它们提示某一任务是一次性还是重复的，并且有一个启动的 `TimerTask` 实例。 `TimerTask` 实际上是一个 `Runnable` （事实上，它实现了它），但还有另外两个方法： `cancel()` 用来取消任务， `scheduledExecutionTime()` 用来返回任务何时启动的近似值。

请注意 `Timer` 却创建了一个非守护线程在后台启动任务，因此在清单 10 中我需要调用 `System.exit()` 来取消任务。在长时间运行的程序中，最好创建一个 `Timer` 守护线程（使用带有指示守护线程状态的参数的构造函数），从而它不会让 VM 活动。

这个类没什么神奇的，但它确实能帮助我们对后台启动的程序的目的了解得更清楚。它还能节省一些 `Thread` 代码，并作为轻量级 `ScheduledExecutorService` （对于还没准备好了解整个 `java.util.concurrent` 包的人来说）。

## 5\. JavaSound

尽管在服务器端应用程序中不常出现，但 sound 对管理员有着有用的 “被动” 意义 — 它是恶作剧的好材料。尽管它很晚才出现在 Java 平台中，JavaSound API 最终还是加入了核心运行时库，封装在 `javax.sound *` 包 — 其中一个包是 MIDI 文件，另一个是音频文件示例（如普遍的 .WAV 文件格式）。

JavaSound 的 “hello world” 是播放一个片段，如清单 11 所示：

##### 清单 11\. 再放一遍，Sam

```
public static void playClip(String audioFile)
{
    try
    {
        AudioInputStream inputStream =
            AudioSystem.getAudioInputStream(
                this.getClass().getResourceAsStream(audioFile));
        DataLine.Info info =
            new DataLine.Info( Clip.class, audioInputStream.getFormat() );
        Clip clip = (Clip) AudioSystem.getLine(info);
        clip.addLineListener(new LineListener() {
                public void update(LineEvent e) {
                    if (e.getType() == LineEvent.Type.STOP) {
                        synchronized(clip) {
                            clip.notify();
                        }
                    }
                }
            });
        clip.open(audioInputStream);

        clip.setFramePosition(0);

        clip.start();
        synchronized (clip) {
            clip.wait();
        }
        clip.drain();
        clip.close();
    }
    catch (Exception ex)
    {
        ex.printStackTrace();
    }
}

```

Show moreShow more icon

大多数还是相当简单（至少 JavaSound 一样简单）。第一步是创建一个文件的 `AudioInputStream` 来播放。为了让此方法尽量与上下文无关，我们从加载类的 `ClassLoader` 中抓取文件作为 `InputStream` 。（ `AudioSystem` 还需要一个 `File` 或 `String` ，如果提前知道声音文件的具体路径。）一旦完成， `DataLine.Info` 对象就提供给 `AudioSystem` ，得到一个 `Clip` ，这是播放音频片段最简单的方法。（其他方法提供了对片段更多的控制 — 例如获取一个 `SourceDataLine` — 但对于 “播放” 来说，过于复杂）。

这里应该和对 `AudioInputStream` 调用 `open()` 一样简单。（”应该” 的意思是如果您没遇到下节描述的错误。）调用 `start()` 开始播放， `drain()` 等待播放完成， `close()` 释放音频线路。播放是在单独的线程进行，因此调用 `stop()` 将会停止播放，然后调用 `start()` 将会从播放暂停的地方重新开始；使用 `setFramePosition(0)` 重新定位到开始。

### 没声音？

JDK 5 发行版中有个讨厌的小错误：在有些平台上，对于一些短的音频片段，代码看上去运行正常，但就是 … 没声音。显然媒体播放器在应该出现的位置之前触发了 `STOP` 事件。（见 参考资料 一节中错误页的链接。）

这个错误 “无法修复”，但解决方法相当简单：注册一个 `LineListener` 来监听 `STOP` 事件，当触发时，调用片段对象的 `notifyAll()` 。然后在 “调用者” 代码中，通过调用 `wait()` 等待片段完成（还调用 `notifyAll()` ）。在没出现错误的平台上，这些错误是多余的，在 Windows® 及有些 Linux® 版本上，会让程序员 “开心” 或 “愤怒”。

## 结束语

现在您都了解了，厨房里的工具。我知道很多人已清楚了解我此处介绍的工具，而我的职业经验告诉我，很多人将从这篇介绍文章，或者说是对长期遗忘在凌乱的抽屉中的小工具的提示中受益。

我在此系列中做个简短的中断，让别人能加入分享他们各自领域的专业经验。但别担心，我还会回来的，无论是 [本系列](/zh/series/5-things-you-didnt-know-about/) 还是其他领域的新的 _5 件事_。在那之前，我鼓励您一直探索 Java 平台，去发现那些能让编程更高效的宝石。

本文翻译自： [Everyday Java tools](https://developer.ibm.com/articles/j-5things12/)（2010-09-14）