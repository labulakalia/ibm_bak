# 使用 Java 实现 Comet 风格的 Web 应用
实现 Servlet 3.0 规范

**标签:** Java,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/wa-cometjava/)

Michael Galpin

发布: 2009-07-22

* * *

## 开始

在本文中，我将展示如何使用各种不同的 Java 技术构建一些简单的 Comet 风格的 Web 应用程序。读者对 Java servlet、Ajax 和 JavaScript 应该有一定的了解。我们将考察 Tomcat 和 Jetty 中一些支持 Comet 的特性，因此需要使用这两个产品的最新版本。本文使用 Tomcat 6.0.14 和 Jetty 6.1.14。另外还需要一个支持 Java 5 或更高版本的 JDK。本文使用 JDK 1.5.0-16。此外还需要看看 Jetty 7 的预发布版，因为它实现了 Servlet 3.0 规范，我们将在本文中研究该规范。请参阅 参考资料 ，找到下载链接。

## 理解 Comet

您可能已经听说过 Comet，因为它最近受到了一定的关注。Comet 有时也称反向 Ajax 或服务器端推技术（server-side push）。其思想很简单：将数据直接从服务器推到浏览器，而不必等到浏览器请求数据。听起来简单，但是如果熟悉 Web 应用程序，尤其是 HTTP 协议，那么您就会知道，这绝不简单。实现 Comet 风格的 Web 应用程序，同时保证在浏览器和服务器上的可伸缩性，这只是在最近几年才成为可能。在本文的后面，我们将看看一些流行的 Java Web 服务器如何支持可伸缩的 Comet 架构，但首先我们来看看为什么要创建 Comet 应用程序，以及用于实现它们的常见设计模式。

### 使用 Comet 的动机

HTTP 协议的成功毋庸置疑。它是 Internet 上大部分信息交换的基础。然而，它也有一些局限性。特别是，它是无状态、单向的协议。请求被发送到 Web 服务器，服务器处理请求并发回一个响应 — 仅此而已。请求必须由客户机发出，而服务器则只能在对请求的响应中发送数据。这至少会影响很多类型的 Web 应用程序的实用性。典型的例子就是聊天程序。另外还有一些例子，例如比赛的比分、股票行情或电子邮件程序。

HTTP 的这些局限性也是它取得一定成功的原因。请求/响应周期使它成为了经典的模型，即每个连接使用一个线程。只要能够快速为请求提供服务，这种方法就有巨大的可伸缩性。每秒钟可以处理大量的请求，只需使用少量的服务器就可以处理很大数量的用户。对于很多经典的 Web 应用程序，例如内容管理系统、搜索应用程序和电子商务站点等等而言，这非常适合。在以上任何一种 Web 应用程序中，服务器提供用户请求的数据，然后关闭连接，并释放那个线程，使之可以为其他请求服务。如果提供初始数据之后仍可能存在交互，那么将连接保持为打开状态，因此线程就不能释放出来，服务器也就不能为很多用户服务。

但是，如果想在对请求做出响应并发送初始数据之后，仍然保持与用户的交互呢？在 Web 早期，这一点常使用 meta 刷新实现。这将自动指示浏览器在指定秒数之后重新装载页面，从而支持简陋的轮询（polling）。这不仅是一种糟糕的用户体验，而且通常效率非常低下。如果没有新的数据要显示在页面上呢？这时不得不重新呈现同样的页面。如果对页面的更改很少，并且页面的大部分没有变化呢？同样，不管是否有必要，都得重新请求和获取页面上的一切内容。

Ajax 的发明和流行改变了上述状况。现在，服务器可以异步通信，因此不必重新请求整个页面。现在可以进行增量式的更新。只需使用 `XMLHttpRequest` 轮询服务器。这项技术通常被称作 Comet。这项技术存在一些变体，每种变体具有不同的性能和可伸缩性。我们来看看这些不同风格的 Comet。

### Comet 风格

Ajax 的出现使 Comet 成为可能。HTTP 的单向性质可以有效地加以规避。实际上有一些不同的方法可以绕过这一点。您可能已经猜到，支持 Comet 的最容易的方式是 _轮询（poll）_ 。使用 `XMLHttpRequest` 向服务器发出调用，返回后，等待一段固定的时间（通常使用 JavaScript 的 `setTimeout` 函数），然后再次调用。这是一项非常常见的技术。例如，大多数 webmail 应用程序就是通过这种技术在电子邮件到达时显示电子邮件的。

这项技术有优点也有缺点。在这种情况下，您期望快速返回响应，就像任何其他 Ajax 请求一样。在请求之间必须有一段暂停。否则，连续不断的请求会冲垮服务器，并且这种情况下显然不具有可伸缩性。这段暂停使应用程序产生一个延时。暂停的时间越长，服务器上的新数据就需要越多的时间才能到达客户机。如果缩短暂停时间，又将重新面临冲垮服务器的风险。但是另一方面，这显然是最简单的实现 Comet 的方式。

现在应该指出，很多人认为轮询并不属于 Comet。相反，他们认为 Comet 是对轮询的局限性的一个解决方案。最常见的 “真正的” Comet 技术是轮询的一种变体，即长轮询（long polling）。轮询与长轮询之间的主要区别在于服务器花多长的时间作出响应。长轮询通常将连接保持一段较长的时间 — 通常是数秒钟，但是也可能是一分钟甚至更长。当服务器上发生某个事件时，响应被发送并随即关闭，轮询立即重新开始。

长轮询相对于一般轮询的优点在于，数据一旦可用，便立即从服务器发送到客户机。请求可能等待较长的时间，期间没有任何数据返回，但是一旦有了新的数据，它将立即被发送到客户机。因此没有延时。如果您使用过基于 Web 的聊天程序，或者声称 “实时” 的任何程序，那么它很可能就是使用了这种技术。

长轮询有一种变体，这是第三种风格的 Comet。这通常被称为流（streaming）。按照这种风格，服务器将数据推回客户机，但是不关闭连接。连接将一直保持开启，直到过期，并导致重新发出请求。 `XMLHttpRequest` 规范表明，可以检查 `readyState` 的值是否为 3 或 Receiving（而不是 4 或 Loaded），并获取正从服务器 “流出” 的数据。和长轮询一样，这种方式也没有延时。当服务器上的数据就绪时，该数据被发送到客户机。这种方式的另一个优点是可以大大减少发送到服务器的请求，从而避免了与设置服务器连接相关的开销和延时。不幸的是， `XMLHttpRequest` 在不同的浏览器中有很多不同的实现。这项技术只能在较新版本的 Mozilla Firefox 中可靠地使用。对于 Internet Explorer 或 Safari，仍需使用长轮询。

至此，您可能会想，长轮询和流都有一个很大的问题。请求需要在服务器上存在一段较长的时间。这打破了每个请求使用一个线程的模型，因为用于一个请求的线程一直没有被释放。更糟糕的是，除非要发回数据，否则该线程一直处于空闲状态。这显然不具有可伸缩性。幸运的是，现代 Java Web 服务器有很多方式可以解决这个问题。

## Java 中的 Comet

现在有很多 Web 服务器是用 Java 构建的。一个原因是 Java 有一个丰富的本地线程模型。因此实现典型的每个连接一个线程的模型便非常简单。该模型对于 Comet 不大适用，但是，Java 对此同样有解决的办法。为了有效地处理 Comet，需要非阻塞 IO，Java 通过它的 NIO 库提供非阻塞 IO。两种最流行的开源服务器 Apache Tomcat 和 Jetty 都利用 NIO 增加非阻塞 IO，从而支持 Comet。然而，这两种服务器中的实现却各不相同。我们来看看 Tomcat 和 Jetty 对 Comet 的支持。

### Tomcat 和 Comet

对于 Apache Tomcat，要使用 Comet，主要需要做两件事。首先，需要对 Tomcat 的配置文件 server.xml 稍作修改。默认情况下启用的是更典型的同步 IO 连接器。现在只需将它切换成异步版本，如清单 1 所示。

##### 清单 1\. 修改 Tomcat 的 server.xml

```
<!-- This is the usual Connector, comment it out and add the NIO one -->
<!-- Connector URIEncoding="utf-8" connectionTimeout="20000" port="8084"
protocol="HTTP/1.1" redirectPort="8443"/ -->
<Connector connectionTimeout="20000" port="8080" protocol="org.apache.
coyote.http11.Http11NioProtocol" redirectPort="8443"/>

```

Show moreShow more icon

这使 Tomcat 可以处理更多的并发连接，但需要说明的是，其中大多数连接有很多时间都处于空闲状态。利用这一点的最容易的方式是创建一个实现 `org.apache.catalina.CometProcessor` 接口的 servlet。这显然是 Tomcat 特有的一个接口。清单 2 显示了一个这样的例子。

##### 清单 2\. Tomcat Comet servlet

```
public class TomcatWeatherServlet extends HttpServlet implements CometProcessor {

    private MessageSender messageSender = null;
    private static final Integer TIMEOUT = 60 * 1000;

    @Override
    public void destroy() {
        messageSender.stop();
        messageSender = null;

    }

    @Override
    public void init() throws ServletException {
        messageSender = new MessageSender();
        Thread messageSenderThread =
                new Thread(messageSender, "MessageSender[" + getServletContext()
.getContextPath() + "]");
        messageSenderThread.setDaemon(true);
        messageSenderThread.start();

    }

    public void event(final CometEvent event) throws IOException, ServletException {
        HttpServletRequest request = event.getHttpServletRequest();
        HttpServletResponse response = event.getHttpServletResponse();
        if (event.getEventType() == CometEvent.EventType.BEGIN) {
            request.setAttribute("org.apache.tomcat.comet.timeout", TIMEOUT);
            log("Begin for session: " + request.getSession(true).getId());
            messageSender.setConnection(response);
            Weatherman weatherman = new Weatherman(95118, 32408);
            new Thread(weatherman).start();
        } else if (event.getEventType() == CometEvent.EventType.ERROR) {
            log("Error for session: " + request.getSession(true).getId());
            event.close();
        } else if (event.getEventType() == CometEvent.EventType.END) {
            log("End for session: " + request.getSession(true).getId());
            event.close();
        } else if (event.getEventType() == CometEvent.EventType.READ) {
            throw new UnsupportedOperationException("This servlet does not accept
data");
        }

    }
}

```

Show moreShow more icon

`CometProcessor` 接口要求实现 `event` 方法。这是用于 Comet 交互的一个生命周期方法。Tomcat 将使用不同的 `CometEvent` 实例调用。通过检查 `CometEvent` 的 `eventType` ，可以判断正处在生命周期的哪个阶段。当请求第一次传入时，即发生 `BEGIN` 事件。 `READ` 事件表明数据正在被发送，只有当请求为 `POST` 时才需要该事件。遇到 `END` 或 `ERROR` 事件时，请求终止。

在清单 2 的例子中，servlet 使用一个 `MessageSender` 类发送数据。这个类的实例是在 servlet 的 init 方法中在其自身的线程中创建，并在 servlet 的 destroy 方法中销毁的。清单 3 显示了 `MessageSender` 。

##### 清单 3\. `MessageSender`

```
private class MessageSender implements Runnable {

    protected boolean running = true;
    protected final ArrayList<String> messages = new ArrayList<String>();
    private ServletResponse connection;
    private synchronized void setConnection(ServletResponse connection){
        this.connection = connection;
        notify();
    }
    public void send(String message) {
        synchronized (messages) {
            messages.add(message);
            log("Message added #messages=" + messages.size());
            messages.notify();
        }
    }
    public void run() {
        while (running) {
            if (messages.size() == 0) {
                try {
                    synchronized (messages) {
                        messages.wait();
                    }
                } catch (InterruptedException e) {
                    // Ignore
                }
            }
            String[] pendingMessages = null;
            synchronized (messages) {
                pendingMessages = messages.toArray(new String[0]);
                messages.clear();
            }
            try {
                if (connection == null){
                    try{
                        synchronized(this){
                            wait();
                        }
                    } catch (InterruptedException e){
                        // Ignore
                    }
                }
                PrintWriter writer = connection.getWriter();
                for (int j = 0; j < pendingMessages.length; j++) {
                    final String forecast = pendingMessages[j] + "<br>";
                    writer.println(forecast);
                    log("Writing:" + forecast);
                }
                writer.flush();
                writer.close();
                connection = null;
                log("Closing connection");
            } catch (IOException e) {
                log("IOExeption sending message", e);
            }
        }
    }
}

```

Show moreShow more icon

这个类基本上是样板代码，与 Comet 没有直接的关系。但是，有两点要注意。这个类含有一个 `ServletResponse` 对象。回头看看清单 2 中的 `event` 方法，当事件为 `BEGIN` 时，response 对象被传入到 `MessageSender` 中。在 MessageSender 的 `run` 方法中，它使用 `ServletResponse` 将数据发送回客户机。注意，一旦发送完所有排队等待的消息后，它将关闭连接。这样就实现了长轮询。如果要实现流风格的 Comet，那么需要使连接保持开启，但是仍然刷新数据。

回头看清单 2 可以发现，其中创建了一个 `Weatherman` 类。正是这个类使用 `MessageSender` 将数据发送回客户机。这个类使用 Yahoo RSS feed 获得不同地区的天气信息，并将该信息发送到客户机。这是一个特别设计的例子，用于模拟以异步方式发送数据的数据源。清单 4 显示了它的代码。

##### 清单 4\. Weatherman

```
private class Weatherman implements Runnable{

    private final List<URL> zipCodes;
    private final String YAHOO_WEATHER = "http://weather.yahooapis.com/forecastrss?p=";

    public Weatherman(Integer... zips) {
        zipCodes = new ArrayList<URL>(zips.length);
        for (Integer zip : zips) {
            try {
                zipCodes.add(new URL(YAHOO_WEATHER + zip));
            } catch (Exception e) {
                // dont add it if it sucks
            }
        }
    }

public void run() {
       int i = 0;
       while (i >= 0) {
           int j = i % zipCodes.size();
           SyndFeedInput input = new SyndFeedInput();
           try {
               SyndFeed feed = input.build(new InputStreamReader(zipCodes.get(j)
.openStream()));
               SyndEntry entry = (SyndEntry) feed.getEntries().get(0);
               messageSender.send(entryToHtml(entry));
               Thread.sleep(30000L);
           } catch (Exception e) {
               // just eat it, eat it
           }
           i++;
       }
}

    private String entryToHtml(SyndEntry entry){
        StringBuilder html = new StringBuilder("<h2>");
        html.append(entry.getTitle());
        html.append("</h2>");
        html.append(entry.getDescription().getValue());
        return html.toString();
    }
}

```

Show moreShow more icon

这个类使用 Project Rome 库解析来自 Yahoo Weather 的 RSS feed。如果需要生成或使用 RSS 或 Atom feed，这是一个非常有用的库。此外，这个代码中只有一个地方值得注意，那就是它产生另一个线程，用于每过 30 秒钟发送一次天气数据。最后，我们再看一个地方：使用该 servlet 的客户机代码。在这种情况下，一个简单的 JSP 加上少量的 JavaScript 就足够了。清单 5 显示了该代码。

##### 清单 5\. 客户机 Comet 代码

```
<%@page contentType="text/html" pageEncoding="UTF-8"%>
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN"
"http://www.w3.org/TR/html4/loose.dtd">

<html>
    <head>
        <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
        <title>Comet Weather</title>
        <SCRIPT TYPE="text/javascript">
            function go(){
                var url = "http://localhost:8484/WeatherServer/Weather"
                var request =  new XMLHttpRequest();
                request.open("GET", url, true);
                request.setRequestHeader("Content-Type","application/x-javascript;");
                request.onreadystatechange = function() {
                    if (request.readyState == 4) {
                        if (request.status == 200){
                            if (request.responseText) {
                                document.getElementById("forecasts").innerHTML =
request.responseText;
                            }
                        }
                        go();
                    }
                };
                request.send(null);
            }
        </SCRIPT>
    </head>
    <body>
        <h1>Rapid Fire Weather</h1>
        <input type="button" onclick="go()" value="Go!"></input>
        <div id="forecasts"></div>
    </body>
</html>

```

Show moreShow more icon

该代码只是在用户单击 **Go** 按钮时开始长轮询。注意，它直接使用 `XMLHttpRequest` 对象，所以这在 Internet Explorer 6 中将不能工作。您可能需要使用一个 Ajax 库解决浏览器差异问题。除此之外，惟一需要注意的是回调函数，或者为请求的 `onreadystatechange` 函数创建的闭包。该函数粘贴来自服务器的新的数据，然后重新调用 go 函数。

现在，我们看过了一个简单的 Comet 应用程序在 Tomcat 上是什么样的。有两件与 Tomcat 密切相关的事情要做：一是配置它的连接器，二是在 servlet 中实现一个特定于 Tomcat 的接口。您可能想知道，将该代码 “移植” 到 Jetty 有多大难度。接下来我们就来看看这个问题。

### Jetty 和 Comet

Jetty 服务器使用稍微不同的技术来支持 Comet 的可伸缩的实现。Jetty 支持被称作 continuations 的编程结构。其思想很简单。请求先被暂停，然后在将来的某个时间点再继续。规定时间到期，或者某种有意义的事件发生，都可能导致请求继续。当请求被暂停时，它的线程被释放。

可以使用 Jetty 的 `org.mortbay.util.ajax.ContinuationSupport` 类为任何 `HttpServletRequest` 创建 `org.mortbay.util.ajax.Continuation` 的一个实例。这种方法与 Comet 有很大的不同。但是，continuations 可用于实现逻辑上等效的 Comet。清单 6 显示清单 2 中的 weather servlet “移植” 到 Jetty 后的代码。

##### 清单 6\. Jetty Comet servlet

```
public class JettyWeatherServlet extends HttpServlet {
    private MessageSender messageSender = null;
    private static final Integer TIMEOUT = 5 * 1000;
    public void begin(HttpServletRequest request, HttpServletResponse response)
throws IOException, ServletException {
        request.setAttribute("org.apache.tomcat.comet", Boolean.TRUE);
        request.setAttribute("org.apache.tomcat.comet.timeout", TIMEOUT);
        messageSender.setConnection(response);
        Weatherman weatherman = new Weatherman(95118, 32408);
        new Thread(weatherman).start();
    }
    public void end(HttpServletRequest request, HttpServletResponse response)
throws IOException, ServletException {
        synchronized (request) {
            request.removeAttribute("org.apache.tomcat.comet");
            Continuation continuation = ContinuationSupport.getContinuation
(request, request);
            if (continuation.isPending()) {
                continuation.resume();
            }
        }
    }
    public void error(HttpServletRequest request, HttpServletResponse response)
throws IOException, ServletException {
        end(request, response);
    }
    public boolean read(HttpServletRequest request, HttpServletResponse response)
throws IOException, ServletException {
        throw new UnsupportedOperationException();
    }
    @Override
    protected void service(HttpServletRequest request, HttpServletResponse response)
throws IOException, ServletException {
        synchronized (request) {
            Continuation continuation = ContinuationSupport.getContinuation
(request, request);
            if (!continuation.isPending()) {
                begin(request, response);
            }
            Integer timeout = (Integer) request.getAttribute
("org.apache.tomcat.comet.timeout");
            boolean resumed = continuation.suspend(timeout == null ? 10000 :
timeout.intValue());

            if (!resumed) {
                error(request, response);
            }
        }
    }
    public void setTimeout(HttpServletRequest request, HttpServletResponse response,
int timeout) throws IOException, ServletException,
            UnsupportedOperationException {
        request.setAttribute("org.apache.tomcat.comet.timeout", new Integer(timeout));
    }
}

```

Show moreShow more icon

这里最需要注意的是，该结构与 Tomcat 版本的代码非常类似。 `begin` 、 `read` 、 `end` 和 `error` 方法都与 Tomcat 中相同的事件匹配。该 servlet 的 service 方法被覆盖为在请求第一次进入时创建一个 continuation 并暂停该请求，直到超时时间已到，或者发生导致它重新开始的事件。上面没有显示 `init` 和 `destroy` 方法，因为它们与 Tomcat 版本是一样的。该 servlet 使用与 Tomcat 相同的 `MessageSender` 。因此不需要修改。注意 `begin` 方法如何创建 Weatherman 实例。对这个类的使用与 Tomcat 版本中也是完全相同的。甚至客户机代码也是一样的。只有 servlet 有更改。虽然 servlet 的变化比较大，但是与 Tomcat 中的事件模型仍是一一对应的。

希望这足以鼓舞人心。虽然完全相同的代码不能同时在 Tomcat 和 Jetty 中运行，但是它是非常相似的。当然，JavaEE 吸引人的一点是可移植性。大多数在 Tomcat 中运行的代码，无需修改就可以在 Jetty 中运行，反之亦然。因此，毫不奇怪，下一个版本的 Java Servlet 规范包括异步请求处理（即 Comet 背后的底层技术）的标准化。 我们来看看这个规范：Servlet 3.0 规范。

## Servlet 3.0 规范

在此，我们不深究 Servlet 3.0 规范的全部细节，只看看 Comet servlet 如果在 Servlet 3.0 容器中运行，可能会是什么样子。注意 “可能” 二字。该规范已经发布公共预览版，但在撰写本文之际，还没有最终版。因此，清单 7 显示的是遵从公共预览规范的一个实现。

##### 清单 7\. Servlet 3.0 Comet

```
@WebServlet(asyncSupported=true, asyncTimeout=5000)
public class WeatherServlet extends HttpServlet {
private MessageSender messageSender;

// init and destroy are the same as other
    @Override
    protected void doGet(HttpServletRequest request, HttpServletResponse response)
    throws ServletException, IOException {
        AsyncContext async = request.startAsync(request, response);
        messageSender.setConnection(async);
        Weatherman weatherman = new Weatherman(95118, 32444);
        async.start(weatherman);;
    }
}

```

Show moreShow more icon

值得高兴的是，这个版本要简单得多。平心而论，如果不遵从 Tomcat 的事件模型，在 Jetty 中可以有类似的实现。这种事件模型似乎比较合理，很容易在 Tomcat 以外的容器（例如 Jetty）中实现，只是没有相关的标准。

回头看看清单 7，注意它的标注声明它支持异步处理，并设置了超时时间。 `startAsync` 方法是 `HttpServletRequest` 上的一个新方法，它返回新的 `javax.servlet.AsyncContext` 类的一个实例。注意， `MessageSender` 现在传递 `AsynContext` 的引用，而不是 `ServletResponse` 的引用。在这里，不应该关闭响应，而是调用 `AsyncContext` 实例上的 complete 方法。还应注意， `Weatherman` 被直接传递到 `AsyncContext` 实例的 `start` 方法。这样将在当前 `ServletContext` 中开始一个新线程。

而且，尽管与 Tomcat 或 Jetty 相比都有较大的不同，但是修改相同风格的编程来处理 Servlet 3.0 规范提议的 API 并不是太难。还应注意，Jetty 7 是为实现 Servlet 3.0 而设计的，目前处于 beta 状态。但是，在撰写本文之际，它还没有实现该规范的最新版本。

## 结束语

Comet 风格的 Web 应用程序可以为 Web 带来全新的交互性。它为大规模地实现这些特性带来一些复杂的挑战。但是，领先的 Java Web 服务器正在为实现 Comet 提供成熟、稳定的技术。在本文中，您看到了 Tomcat 和 Jetty 上当前风格的 Comet 的不同点和相似点，以及正在进行的 Servlet 3.0 规范的标准化。Tomcat 和 Jetty 使如今构建可伸缩的 Comet 应用程序成为可能，并且明确了未来面向 Servlet 3.0 标准化的升级路线。

## 下载 Weather Server 源代码

[WeatherServer.zip](http://download.boulder.ibm.com/ibmdl/pub/software/dw/web/wa-cometjava/WeatherServer.zip)