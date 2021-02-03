# IOCP 浅析
从 IOCP 的基本概念和思想出发比较 IOCP 与传统 Server/Client 实现的利弊

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-lo-iocp/)

郭仁祥

发布: 2012-10-22

* * *

## 什么是 IOCP

随着计算机技术，尤其是网络技术的飞速发展，如今的程序开发人员不仅仅局限于基于单机运行或单一线程的应用程序的开发。服务器端 / 客户端模式的最显著的特点是一个服务器端应用能同时为多个客户端提供服务。而传统的服务器端 / 客户端模式通常为每一个客户端创建一个独立的线程，这种实现方式在客户端数量不多的情况下问题不大，但对于类似于铁路网络订票这样的瞬间客户数量巨大的系统来说，效率极端低下。这是因为一方面创建新线程操作系统开销较大，另一方面同时有许多线程处于运行状态，操作系统内核需要花费大量时间进行上下文切换，并没有在线程执行上花更多的时间。

因此，微软在 Winsocket2 中引入了 IOCP（Input/Output Completion Port）模型。IOCP 是 Input/Output Completion Port（I/O 完成端口）的简称。简单的说，IOCP 是一种高性能的 I/O 模型，是一种应用程序使用线程池处理异步 I/O 请求的机制。Java7 中对 IOCP 有了很好的封装，程序员可以非常方便的时候经过封装的 channel 类来读写和传输数据。

## 同步 / 异步，阻塞 / 非阻塞

所谓同步，就是在发出一个功能调用时，在没有得到结果之前，该调用就不返回。按照这个定义，其实绝大多数函数或方法都是同步调用。

异步的概念和同步相对。当一个异步过程调用发出后，调用者不能立刻得到结果。实际处理这个调用的部件在完成后，通过状态、通知和回调来通知调用者。

通俗来讲，完成一件事再去做另外一件事就是同步，而一起做两件或者两件以上的事情就是异步了。

拿一个服务器与客户端通信的例子来说。

如果是同步：

Client 发送一条请求消息给 Server，这个时候 Client 就会等待 Server 处理该请求。这段时间内 Client 只有等待直到 Server 回复响应信息给 Client。Client 只有收到该响应信息后，才能发起第二条请求消息。这样无疑大大降低了系统的性能。

而如果是异步：

Client 发送一条请求消息给 Server，Client 并不等待 Server 的处理结果，而是继续发送第二条甚至更多的请求消息。Server 会将这些请求都存入队列，逐条处理，并将处理完的结果回复给 Client。这样一来，Client 就可以不用等待，效率大大提高。

阻塞调用是指调用结果返回之前，当前线程会被挂起。函数或方法只有在得到结果之后才会返回。阻塞和同步有点类似，但是同步调用的时候线程还是处于激活状态，而阻塞时线程会被挂起。

非阻塞调用和阻塞的概念相对应，指在不能立刻得到结果之前，该函数或方法不会阻塞当前线程而是立刻返回。

##### 清单 1\. 传统的网络应用代码

```
try {
     ServerSocket server = new ServerSocket(9080);
     while (true) {
         Socket client = server.accept();
         new Thread(new SocketHandle(client)).start();
     }
} catch (IOException e) {
     e.printStackTrace();
}

```

Show moreShow more icon

相信只要写过网络应用程序的朋友，应该对这样的结构再熟悉不过了。Accept 后线程被挂起，等待一个客户发出请求，而后创建新线程来处理请求。当新线程处理客户请求时，起初的线程循环回去等待另个客户请求。在这个并发模型中，对每个客户都创建了一个线程。其优点在于等待请求的线程只需要做很少的工作，而大部分的时间，该线程在休眠，因为 recv 处于阻塞状态。如前文所述，创建线程的开销远远大于程序员的预计，尤其是在并发量巨大的情况下，这种传统的并发模型效率极端低下。

解决这个问题的方法之一就是 IOCP，说白了 IOCP 就是一个消息队列。我们设想一下，如果事先开好 N 个线程，让它们 hold 住，将所有用户的请求都投递到一个消息队列中去。让后这 N 个线程逐一从消息队列中去取出消息并加以处理。这样一来，就可以避免对没有用户请求都开新线程，不仅减少了线程的资源，也提高了线程的利用率。

## IOCP 实现的基本步骤

那么 IOCP 完成端口模型又是怎样实现的呢？首先我们创建一个完成端口 CreateIOCompletionPort，然后再创建一个或多个工作线程，并指定它们到这个完成端口上去读取数据。再将远程连接的套接字句柄关联到这个完成端口。工作线程调用 getQueuedCompletionStatus 方法在关联到这个完成端口上的所有套接字上等待 I/O 的完成，再判断完成了什么类型的 I/O，然后接着发出 WSASend 和 WSARecv，并继续下一次循环阻塞在 getQueuedCompletionStatus。

具体的说，一个完成端口大概的处理流程包括：

1. 创建一个完成端口；





    ```
    Port port = createIoCompletionPort(INVALID_HANDLE_VALUE, 0, 0, fixedThreadCount());

    ```





    Show moreShow more icon

2. 创建一个线程 ThreadA；

3. ThreadA 线程循环调用 GetQueuedCompletionStatus 方法来得到 I/O 操作结果，这个方法是一个阻塞方法；





    ```
    While(true){
         getQueuedCompletionStatus(port, ioResult);
    }

    ```





    Show moreShow more icon

4. 主线程循环调用 accept 等待客户端连接上来；

5. 主线程 accept 返回新连接建立以后，把这个新的套接字句柄用 CreateIoCompletionPort 关联到完成端口，然后发出一个异步的 Read 或者 Write 调用，因为是异步函数，Read/Write 会马上返回，实际的发送或者接收数据的操作由操作系统去做。





    ```
    if (handle != 0L) {
         createIoCompletionPort(handle, port, key, 0);
    }

    ```





    Show moreShow more icon

6. 主线程继续下一次循环，阻塞在 accept 这里等待客户端连接。

7. 操作系统完成 Read 或者 Write 的操作，把结果发到完成端口。
8. ThreadA 线程里的 GetQueuedCompletionStatus() 马上返回，并从完成端口取得刚完成的 Read/Write 的结果。
9. 在 ThreadA 线程里对这些数据进行处理 ( 如果处理过程很耗时，需要新开线程处理 )，然后接着发出 Read/Write，并继续下一次循环阻塞在 GetQueuedCompletionStatus() 这里。

## Java 中异步 I/O 和 IOCP

其实，JDK 中的 IOCP 的实现与前文中所描述的实现思想基本相同，唯一的区别是 JDK 中封装了相关的内部实现。

Iocp 类：该类使用默认的访问修饰符，因此只能被同包（sun.nio.ch）的其他类访问。同时，创建 IOCP 端口的 createIoCompletionPort 方法也被封装在 Iocp 类的构造函数中，套接字的关联也不再需要调用 createIoCompletionPort 而是直接调用 Iocp 类的 associate 方法。

Iocp.CompletionStatus 类：Iocp 的内部类 CompletionStatus 用于存取完成端口的状态信息。

Iocp.EventHandlerTask 类：Iocp 的内部类 EventHandlerTask 用来将 I/O 操作的 Handle 关联到完成端口，循环的获取完成端口的状态，根据操作系统返回的完成端口状态从任务队列中提取下一个要执行的任务并调用执行。

Iocp 的调用类：sun.nio.ch 包中的 WindowsAsynchronousSocketChannelImpl、WindowsAsynchronousServerSocketChannelImpl 和 WindowsAsynchronousFileChannelImpl 等类使用了 Iocp 类的相关方法，达到多并发量时的高性能 I/O 访问。

Native 方法：Iocp 类中的许多内部方法都是通过调用 JNI 方法实现的。

## 实例分析

从上节描述中，您或许会意识到，Java7 中对 IOCP 进行了很好的封装。我们可以直接使用 JDK 中众多的 ChannelImpl 类，这些类的读写操作都应用了 Iocp 类。当然，您也可以将自定义的逻辑实现类放置在 sun.nio.ch 包中以重用 Iocp 类，甚至还可以自己编写和实现 Iocp 类。本节中，我们通过一个 Log 传输的具体实例来看 Iocp 以及应用类在实际应用中的使用。

假设现在有建立了连接的服务器端和客户端，服务器端执行得到的大量 Log 信息需要传输到客户端的 console 中。这些信息的传输必须是异步的，也就是说服务器端不等待传输的完成而是继续执行并产生更多的 Log 信息。同时，Log 信息在 channel 上传输的时候又必须是同步的，这样才能保证 Log 信息的时间次序保持一致。

首先，服务器端监听端口，建立与客户端通信的 channel。

##### 清单 2\. 建立 channel

```
final AsynchronousServerSocketChannel server = AsynchronousServerSocketChannel.open();
server.bind(new InetSocketAddress(port));
server.accept(
new Session(null), new CompletionHandler<AsynchronousSocketChannel, Session>()
{
     public void completed(AsynchronousSocketChannel ch,Session session) {
         Session s = new Session(ch);
         server.accept(s, this);
     }

     public void failed(Throwable e, Session session) {
     }
});

```

Show moreShow more icon

服务器端主线程并不直接将 Log 信息通过 channel 输出到客户端的 console 上。而是调用 writeMessage 方法将执行得到的 log 信息写入事先准备好的 messageQueue 里，并继续业务逻辑的执行。

##### 清单 3\. 写入 log 到 Queue

```
public synchronized void flush() throws IOException {
    int size = buffer.position();
    Message msg = new Message(this.msgType, new byte[size]);
    buffer.rewind();
    final Object mutex = new Object();
    buffer.get(msg.body);
    buffer.rewind();
    session.writeMsg(msg);
}

Public void writeMessage(Message msg) {
    synchronized(messageQueue) {
        messageQueue.offer(msg);
    }
}

```

Show moreShow more icon

同时，另一个线程负责不断从 messageQueue 中读取 Log 信息，顺序写到 channel 上传输给客户端 console。WriteCompletionHandler 是 CompletionHandler 接口的自定义实现类，用来控制每次 channel 上的写操作必须完成才返回，避免混淆输出的 Log 次序。

##### 清单 4\. 从 Queue 中读取并传输 Log 至客户端

```
Thread thread=new Thread(new Runnable(){
    public void run(){
        Message msg = null;
        synchronized(messageQueue){
            msg = messageQueue.poll();
            ByteBuffer byteBuffer = msg.pasteByteBuffer();
            synchronized(canWrite){
                channel.write(byteBuffer, new WriteCompletionHandler(msg));
            }
        }
    }
}).start();

```

Show moreShow more icon

这样一来，我们能很轻松的使用 channelImpl 类来网络应用，而且如果您深入到这些 channelImpl 类的内部去，您就会看到许多逻辑的实现都依赖了底层的 sun.nio.ch.Iocp 类。

## 小结

本文从传统的 Server/Client 实现出发，分析了传统的网络应用开发的弊端。进而介绍了 IOCP 的基本概念和思想，比较 IOCP 与传统 Server/Client 实现的利弊。分析了 IOCP 实现的基本步骤和 JDK 中异步 I/O 和 IOCP 的具体实现，并以一个 Log 集中处理的实例分析和介绍了 IOCP 思想在具体实践中的应用。

## 相关主题

- [Wikipedia](https://en.wikipedia.org/wiki/Input/output_completion_port)
- [Input/Output Completion Port](https://docs.microsoft.com/en-us/windows/win32/fileio/i-o-completion-ports)
- [NIO.2 入门，第 1 部分：异步通道 API](https://www.ibm.com/developerworks/cn/java/j-nio2-1/)
- [Servlet 3.0 实战：异步 Servlet 与 Comet 风格应用程序](https://www.ibm.com/developerworks/cn/java/j-lo-comet/)