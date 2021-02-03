# 实战 Jetty
一个非常方便使用的 Web 服务器

**标签:** Java,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/wa-lo-jetty/)

周登朋

发布: 2008-04-30

* * *

## 特性简介

### 易用性

易用性是 Jetty 设计的基本原则，易用性主要体现在以下几个方面：

1. 通过 XML 或者 API 来对 Jetty 进行配置；
2. 默认配置可以满足大部分的需求；
3. 将 Jetty 嵌入到应用程序当中只需要非常少的代码；

### 可扩展性

在使用了 Ajax 的 Web 2.0 的应用程序中，每个连接需要保持更长的时间，这样线程和内存的消耗量会急剧的增加。这就使得我们担心整个程序会因为单个组件陷入瓶颈而影响整个程序的性能。但是有了 Jetty：

1. 即使在有大量服务请求的情况下，系统的性能也能保持在一个可以接受的状态。
2. 利用 Continuation 机制来处理大量的用户请求以及时间比较长的连接。

另外 Jetty 设计了非常良好的接口，因此在 Jetty 的某种实现无法满足用户的需要时，用户可以非常方便地对 Jetty 的某些实现进行修改，使得 Jetty 适用于特殊的应用程序的需求。

### 易嵌入性

Jetty 设计之初就是作为一个优秀的组件来设计的，这也就意味着 Jetty 可以非常容易的嵌入到应用程序当中而不需要程序为了使用 Jetty 做修改。从某种程度上，你也可以把 Jetty 理解为一个嵌入式的Web服务器。

## 部署应用程序

将自己的应用程序部署到 Jetty 上面是非常简单的，首先将开发好的应用程序打成 WAR 包放到 Jetty 的 Webapps 目录下面。然后用如下的命令来启动 Jetty 服务器： `Java –jar start.jar` ， 在启动服务器后。我们就可以访问我们的应用程序了，Jetty 的默认端口是 8080，WAR 的名字也就是我们的应用程序的 Root Context。例如一个典型的 URL 就是：`http://127.0.0.1:8080/sample/index.jsp` 。

## 如何将 Jetty 嵌入到程序当中

将 Jetty 嵌入到程序当中是非常简单的， 如 代码 1 所示：首先我们创建一个 Server 对象， 并设置端口为 8080，然后为这个 Server 对象添加一个默认的 Handler。接着我们用配置文件 jetty.xml 对这个 server 进行设置，最后我们使用方法 `server.start()` 将 Server 启动起来就可以了。从这段代码可以看出，Jetty 是非常适合用于作为一个组件来嵌入到我们的应用程序当中的，这也是 Jetty 的一个非常重要的特点。

##### 清单 1\. 代码片断

```
public class JettyServer {

    public static void main(String[] args) {
        Server server = new Server(8080);
        server.setHandler(new DefaultHandler());
        XmlConfiguration configuration = null;
        try {
            configuration = new XmlConfiguration(
                new FileInputStream("C:/development/Jetty/jetty-6.1.6rc0/etc/jetty.xml"));
        } catch (FileNotFoundException e1) {
            e1.printStackTrace();
        } catch (SAXException e1) {
            e1.printStackTrace();
        } catch (IOException e1) {
            e1.printStackTrace();
        }

        try {
            configuration.configure(server);
            server.start();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}

```

Show moreShow more icon

接下来我们分析一下 Jetty Server 是如何启动的。首先我们注意到 Server 类，这个类实际上继承了 HttpServer, 当启动 Jetty 服务器的时候，就是说，在 Jetty 根目录下的命令行下如果输入 `java -jar start.jar etc/jetty.xml`，注意这里有一个配置文件 jetty.xml 做为运行参数，这个参数也可以是其它的配置文件，可以是多个 XML 配置文件，其实这个配置文件好比我们使用 Struts 时的 struts-config.xml 文件，将运行 Server 需要用到的组件写在里面，比如上一节中 HttpServer 的配置需要的组件类都可以写在这个配置文件中。按上述方法启动 Jetty Server 时，就会调用 Server 类里面的 `main` 方法，这个入口方法首先会构造一个 Server 类实例（其实也就构造了一个 HttpServer），创建实例的过程中就会构造 XmlConfiguration 类的对象来读取参数配置文件，之后再由这个配置文件产生的 XmlConfiguration 对象来配置这个 Server，配置过程其实是运用了 Java 的反射机制，调用 Server 的方法并传入配置文件中所写的参数来向这个 Server 添加 HttpListener，HttpContext，HttpHandler，以及 Web Application（对应于我们的 Web 应用）。

## Jetty 的 Continuation 机制

讨论 Jetty 的 Continuation 机制，首先需要提到 Ajax 技术，Ajax 技术是当前开发 Web 应用的非常热门的技术，也是 Web 2.0 的一个重要的组成部分。Ajax 技术中的一个核心对象是 XMLHttpRequest 对象，这个对象支持异步请求，所谓异步请求即是指当客户端发送一个请求到服务器的时候，客户端不必一直等待服务器的响应。这样就不会造成整个页面的刷新，给用户带来更好的体验。而当服务器端响应返回时，客户端利用一个 Javascript 函数对返回值进行处理，以更新页面上的部分元素的值。但很多时候这种异步事件只是在很小一部分的情况下才会发生，那么怎么保证一旦服务器端有了响应之后客户端马上就知道呢，我们有两种方法来解决这个问题，一是让浏览器每隔几秒请求服务器来获得更改，我们称之为 **轮询** 。二是服务器维持与浏览器的长时间的连接来传递数据，长连接的技术称之为 **Comet** 。

大家很容易就能发现轮询方式的主要缺点是产生了大量的传输浪费。因为可能大部分向服务器的请求是无效的，也就是说客户端等待发生的事件没有发生，如果有大量的客户端的话，那么这种网络传输的浪费是非常厉害的。特别是对于服务器端很久才更新的应用程序来讲，比如邮件程序，这种浪费就更是巨大了。并且对 Server 端处理请求的能力也相应提高了要求。如果很长时间才向 Server 端发送一次请求的话，那么客户端就不能的得到及时的响应。

如果使用 Comet 技术的话，客户端和服务器端必须保持一个长连接，一般情况下，服务器端每一个 Servlet 都会独占一个线程，这样就会使得服务器端有很多线程同时存在，这在客户端非常多的情况下也会对服务器端的处理能力带来很大的挑战。

Jetty 利用 Java 语言的非堵塞 I/O 技术来处理并发的大量连接。 Jetty 有一个处理长连接的机制：一个被称为 **Continuations** 的特性。利用 Continuation 机制，Jetty 可以使得一个线程能够用来同时处理多个从客户端发送过来的异步请求，下面我们通过一个简化的聊天程序的服务器端的代码来演示不使用 Continuation 机制和使用 Continuation 的差别。

##### 清单 2\. Continuation 机制

```
public class ChatContinuation extends HttpServlet{

    public void doPost(HttpServletRequest request, HttpServletResponse response){
        postMessage(request, response);
    }

    private void postMessage(HttpServletRequest request, HttpServletResponse response)
    {
        HttpSession session = request.getSession(true);
        People people = (People)session.getAttribute(session.getId());
        if (!people.hasEvent())
        {
            Continuation continuation =
                ContinuationSupport.getContinuation(request, this);
            people.setContinuation(continuation);
            continuation.suspend(1000);
        }
        people.setContinuation(null);
        people.sendEvent(response);
    }
}

```

Show moreShow more icon

大家注意到，首先获取一个 Continuation 对象，然后把它挂起 1 秒钟，直到超时或者中间被 `resume` 函数唤醒位置，这里需要解释的是，在调用完 `suspend` 函数之后，这个线程就可处理其他的请求了，这也就大大提高了程序的并发性，使得长连接能够获得非常好的扩展性。

如果我们不使用 Continuation 机制，那么程序就如 清单 3 所示：

##### 清单 3\. 不使用 Continuation 机制

```
public class Chat extends HttpServlet{
    public void doPost(HttpServletRequest request, HttpServletResponse response){
        postMessage(request, response);
    }

    private void postMessage(HttpServletRequest request, HttpServletResponse response)
    {
        HttpSession session = request.getSession(true);

        People people = (People)session.getAttribute(session.getId());

        while (!people.hasEvent())
        {
            try {
                Thread.sleep(1000);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        }
        people.setContinuation(null);
        people.sendEvent(response);
    }
}

```

Show moreShow more icon

大家注意到在等待事件发生的时间里，线程被挂起，直到所等待的事件发生为止，但在等待过程中，这个线程不能处理其他请求，这也就造成了在客户端非常多的情况下服务器的处理能力跟不上的情况。下面我们解释一下 Jetty 的 Continuation 的机制是如何工作的。

为了使用 Continuatins，Jetty 必须配置为使用它的 `SelectChannelConnector` 处理请求。这个 connector 构建在 java.nio API 之上，允许它维持每个连接开放而不用消耗一个线程。当使用 `SelectChannelConnector` 时， `ContinuationSupport.getContinuation()` 提供一个 `SelectChannelConnector.RetryContinuation` 实例（但是，您必须针对 `Continuation` 接口编程）。当在 `RetryContinuation` 上调用 `suspend()` 时，它抛出一个特殊的运行时异常 — `RetryRequest` ，该异常传播到 servlet 外并且回溯到 filter 链，最后被 `SelectChannelConnector` 捕获。但是不会发送一个异常响应给客户端，而是将请求维持在未决 Continuations 队列里，则 HTTP 连接保持开放。这样，用来服务请求的线程返回给 ThreadPool，然后又可以用来服务其他请求。暂停的请求停留在未决 Continuations 队列里直到指定的过期时间，或者在它的 Continuation 上调用 `resume()` 方法。当任何一个条件触发时，请求会重新提交给 servlet（通过 filter 链）。这样，整个请求被”重播”直到 RetryRequest 异常不再抛出，然后继续按正常情况执行。

## Jetty 的安全性

为了防止任何人都有权限去关闭一个已经开启的 Jetty 服务器， 我们可以通过在启动 Jetty 服务器的时候指定参数来进行控制，使得用户必须提供密码才能关闭 Jetty 服务器，启动 Jetty 服务器的命令如下所示：

```
java -DSTOP.PORT=8079 -DSTOP.KEY=mypassword -jar start.jar

```

Show moreShow more icon

这样，用户在停止 Jetty 服务器的时候，就必须提供密码 “mypassword”。

## 结束语

Jetty 是一个非常方便使用的 Web 服务器，它的特点在于非常小，很容易嵌入到我们的应用程序当中，而且针对 Web 2.0 的 Ajax 技术进行了特别的优化，这也使得我们的使用 Ajax 的应用程序可以拥有更好的性能。