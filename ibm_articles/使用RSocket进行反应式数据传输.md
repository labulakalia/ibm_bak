# 使用 RSocket 进行反应式数据传输
微服务架构中的数据传输协议

**标签:** Java,Spring,反应式系统,微服务,消息传递

[原文链接](https://developer.ibm.com/zh/articles/j-using-rsocket-for-reactive-data-transfer/)

成 富

发布: 2019-10-14

* * *

在微服务架构中，不同服务之间通过应用协议进行数据传输。典型的传输方式包括基于 HTTP 协议的 REST 或 SOAP API 和基于 TCP 字节流的 gRPC 等。HTTP 协议的优势在于其广泛的适用性，有非常多的服务器和客户端实现的支持，但 HTTP 协议本身比较简单，只支持请求-响应模式。gRPC 等基于 TCP 的协议使用二进制字节流传输，保证了传输的效率。不过 gRPC 基于 HTTP/2，不支持其他传输层实现。HTTP/2 协议使用了二进制字节流，但是并没有改变 HTTP/1 协议已有的语义，更多的只是改变了传输时的格式。已有应用协议的问题在于单一的交互模式，只支持请求-响应模式，而此模式对于很多应用场景来说是不合适的。典型的例子是消息推送，以 HTTP 协议为例，如果客户端需要获取最新的推送消息，就必须使用轮询。客户端不停的发送请求到服务器来检查更新，这无疑造成了大量的资源浪费。虽然服务器发送事件（Server-Sent Events，SSE）可以用来推送消息，不过 SSE 是一个简单的文本协议，仅提供有限的功能。此外，WebSocket 可以进行双向数据传输，不过它没有提供应用层协议支持。请求-响应模式的另外一个问题是，如果某个请求的响应时间过长，会阻塞之后的其他请求的处理。RSocket 协议的出现，很好的解决了已有协议的这些问题。

## RSocket 介绍

RSocket 是一个 OSL 七层模型中 5/6 层的协议，是 TCP/IP 之上的应用层协议。RSocket 可以使用不同的底层传输层，包括 TCP、WebSocket 和 Aeron。TCP 适用于分布式系统的各个组件之间交互，WebSocket 适用于浏览器和服务器之间的交互，Aeron 是基于 UDP 协议的传输方式，这就保证了 RSocket 可以适应于不同的场景。使用 RSocket 的应用层实现可以保持不变，只需要根据系统环境、设备能力和性能要求来选择合适的底层传输方式即可。RSocket 作为一个应用层协议，可以很容易在其基础上定义应用自己的协议。此外，RSocket 使用二进制格式，保证了传输的高效，节省带宽。而且，通过基于反应式流语义的流控制，RSocket 保证了消息传输中的双方不会因为请求的压力过大而崩溃。

### RSocket 交互模式

RSocket 支持四种不同的交互模式，见表 1。

##### 表 1\. RSocket 支持的四种交互模式

**模式****说明**请求-响应（request/response）这是最典型也最常见的模式。发送方在发送消息给接收方之后，等待与之对应的响应消息。请求-响应流（request/stream）发送方的每个请求消息，都对应于接收方的一个消息流作为响应。发后不管（fire-and-forget）发送方的请求消息没有与之对应的响应。通道模式（channel）在发送方和接收方之间建立一个双向传输的通道。

下面介绍 RSocket 协议的具体内容。

### RSocket 帧

RSocket 协议在传输时使用帧（frame）来表示单个消息。每个帧中包含的可能是请求内容、响应内容或与协议相关的数据。一个应用消息可能被切分成多个片段（fragment）以包含在一个帧中。根据底层传输协议的不同，一个表示帧长度的字段可能是必须的。由于 TCP 协议没有提供帧支持，所以 RSocket 的帧长度字段是必须的。对于提供了帧支持的传输协议，RSocket 帧只是简单的封装在传输层消息中；对于没有提供帧支持的传输协议，每个 RSocket 帧之前都需要添加一个 24 字节的字段表示帧长度。

#### RSocket 帧的内容

在每个 RSocket 帧中，最起始的部分是帧头部，包括 31 字节的流标识符，6 字节的帧类型和 10 字节的标志位。RSocket 协议中的流（stream）表示的是一个操作的单元。每个流有自己唯一的标识符。流标识符由发送方生成。值为 0 的流标识符表示与连接相关的操作。客户端的流标识符从 1 开始，每次递增 2；服务器端的流标识符从 2 开始，每次递增 2。在帧头部之后的内容与帧类型相关。

#### RSocket 帧的类型

RSocket 中定义了不同类型的帧，见表 2。

##### 表 2\. RSocket 中的帧类型

**帧类型****说明**SETUP由客户端发送来建立连接，流标识符为 0。REQUEST\_RESPONSE请求-响应模式中的请求内容。REQUEST\_FNF发后不管模式中的请求内容。REQUEST\_STREAM请求-响应流模式中的请求内容。REQUEST\_CHANNEL通道模式中的建立通道的请求。发送方只能发一个 REQUEST\_CHANNEL 帧。REQUEST\_N基于反应式流语义的流控制，相当于反应式流中的 request(n)。PAYLOAD表示负载的帧。通过不同的标志位来表示状态。NEXT 标志位表示接收到流中的数据，COMPLETE 标志位表示流结束。ERROR表示连接层或应用层错误。CANCEL取消当前请求。KEEPALIVE表示连接层或应用层错误。KEEPALIVE启用租约模式。METADATA\_PUSH异步元数据推送。RESUME在连接中断后恢复传输。RESUME\_OK恢复传输成功。EXT帧类型扩展。

RSocket 中的负载分成元数据和数据两种，二者可以使用不同的编码方式。元数据是可选的。帧头部有标志位指示帧中是否包含元数据。某些特定类型的帧可以添加元数据。

#### RSocket 帧交互

根据不同的交互模式，发送方和接收方之间有不同的帧交互，下面是几个典型的帧交互示例：

- 在请求-响应模式中，发送方发送 REQUEST\_RESPONSE 帧，接收方发送 PAYLOAD 帧并设置 COMPLETE 标志位。
- 在发后不管模式中，发送方发送 REQUEST\_FNF 帧。
- 在请求-响应流模式中，发送方发送 REQUEST\_STREAM 帧，接收方发送多个 PAYLOAD 帧。设置 COMPLETE 标志位的 PAYLOAD 帧表示流结束。
- 在通道模式中，发送方发送 REQUEST\_CHANNEL 帧。发送方和接收方都可以发送 PAYLOAD 帧给对方。设置 COMPLETE 标志位的 PAYLOAD 帧表示其中一方的流结束。

除了发后不管模式之外，其余模式中的接收方都可以通过 ERROR 帧或 CANCEL 帧来结束流。

### 流控制

RSocket 使用 Reactive Streams 语义来进行流控制（flow control），也就是 request(n)模式。流的发送方通过 REQUEST\_N 帧来声明它允许接收方发送的 PAYLOAD 帧的数量。REQUEST\_N 帧一旦发出就不能收回，而且所产生的效果是累加的。比如，发送方发送 request(2)和 request(3)帧之后，接收方允许发送 5 个 PAYLOAD 帧。

除了基于 Reactive Streams 语义的流程控制之外，RSocket 还可以使用租约模式。租约模式只是限定了在某个时间段内，请求者所能发送的最大请求数量。

## Java 实现

RSocket 提供了不同语言的实现，包括 Java、Kotlin、JavaScript、Go、.NET 和 C++ 等。对 Java 项目来说，只需要添加相应的 Maven 依赖即可。RSocket 的 Java 实现库都在 Maven 分组 io.rsocket 中。其中常用的库包括核心功能库 rsocket-core 和表 3中列出的传输层实现。本文中使用的版本是 1.0.0-RC3。

##### 表 3\. RSocket 的传输层实现

**Maven 实现库名称****底层实现****支持协议**rsocket-transport-nettyReactor NettyTCP 和 WebSocketrsocket-transport-akkaAkkaTCP 和 WebSocketrsocket-transport-aeronAeronUDP

在代码清单 1 中，Maven 项目中添加了 RSocket 相关的依赖和基于 Reactor Netty 的传输层实现。

##### 清单 1\. 添加 RSocket 相关的 Maven 依赖

```
<dependency>
    <groupId>io.rsocket</groupId>
    <artifactId>rsocket-core</artifactId>
    <version>1.0.0-RC3</version>
</dependency>
<dependency>
    <groupId>io.rsocket</groupId>
    <artifactId>rsocket-transport-netty</artifactId>
    <version>1.0.0-RC3</version>
</dependency>

```

Show moreShow more icon

下面介绍 RSocket 中不同模式的使用方式。

### 请求-响应模式

首先从最常用的请求-响应模式开始介绍 RSocket 的用法。代码清单 2 给出了一个请求-响应模式的 RSocket 服务器和客户端的示例。RSocketFactory 类用来创建 RSocket 服务器和客户端。

##### 清单 2\. 请求-响应模式示例

```
import io.rsocket.AbstractRSocket;
import io.rsocket.Payload;
import io.rsocket.RSocket;
import io.rsocket.RSocketFactory;
import io.rsocket.transport.netty.client.TcpClientTransport;
import io.rsocket.transport.netty.server.TcpServerTransport;
import io.rsocket.util.DefaultPayload;
import reactor.core.publisher.Mono;

public class RequestResponseExample {

public static void main(String[] args) {
    RSocketFactory.receive()
        .acceptor(((setup, sendingSocket) -> Mono.just(
            new AbstractRSocket() {
              @Override
              public Mono<Payload> requestResponse(Payload payload) {
                return Mono.just(DefaultPayload.create("ECHO >> " + payload.getDataUtf8()));
              }
            }
        )))
        .transport(TcpServerTransport.create("localhost", 7000)) //指定传输层实现
        .start() //启动服务器
        .subscribe();

    RSocket socket = RSocketFactory.connect()
        .transport(TcpClientTransport.create("localhost", 7000)) //指定传输层实现
        .start() //启动客户端
        .block();

    socket.requestResponse(DefaultPayload.create("hello"))
        .map(Payload::getDataUtf8)
        .doOnNext(System.out::println)
        .block();

    socket.dispose();
}
}

```

Show moreShow more icon

`RSocketFactory.receive()` 方法返回用来创建服务器的 `ServerRSocketFactory` 类的对象。 `ServerRSocketFactory` 的 `acceptor()` 方法的参数是 `SocketAcceptor` 接口，该接口只有一个方法 `Mono<RSocket> accept(ConnectionSetupPayload setup, RSocket sendingSocket)` 。该方法的参数 `setup` 表示的是 `SETUP` 帧的负载内容，而 `sendingSocket` 表示的是发送请求的 `RSocket` 对象。该方法返回的 `Mono<RSocket>` 对象包含的是处理请求的 `RSocket` 对象。代码中使用 `Lambda` 表达式实现了 `SocketAcceptor` 接口。在代码中， `SocketAcceptor` 的 `accept()` 方法返回的是抽象类 `AbstractRSocket` 的匿名实现，只实现了其中的 `requestResponse()` 方法。具体的请求处理逻辑是在请求数据内容上添加 `"ECHO >> "` 前缀。接下来的 `transport()` 方法指定 `ServerTransport` 接口的实现作为 `RSocket` 底层的传输层实现（这里使用 `TcpServerTransport` 类的 `create()` 方法创建在 `localhost` 上 `7000` 端口的 `TCP` 服务器端）。再通过 `start()` 方法可以得到表示服务器实例的 `Mono` 对象。最后的 `subscribe()` 方法用来触发整个启动过程。

`RSocketFactory.connect()` 方法用来创建 `RSocket` 客户端，返回 `ClientRSocketFactory` 类的对象。接下来的 `transport()` 方法指定传输层 `ClientTransport` 实现 （这里通过 `TcpClientTransport.create()` 方法创建到服务器的 `TCP` 连接）。再通过其 `start()` 方法可以得到 `Mono<RSocket>` 对象。最后调用 `block()` 方法等待客户端启动并返回 `RSocket` 对象。

`RSocket` 对象用来与服务器端进行交互。 `RSocket` 类的 `requestResponse()` 方法发送 `Payload` 接口表示的负载并等待响应。该方法的返回值是表示响应的 `Mono<Payload>` 对象。对于返回的响应，示例中只是简单地输出到控制台。 `DefaultPayload.create()` 方法可以简单地创建 `Payload` 对象。 `RSocket` 类的 `dispose()` 方法用来销毁该对象。

### 请求-响应流模式

请求-响应流模式的用法与请求-响应模式很相似。代码清单 3 给出了请求-响应流模式的示例。服务器端的 `AbstractRSocket` 类的实现覆写了 `requestStream()` 方法。对于每个请求的 Payload 对象，都需要返回一个表示响应流的 `Flux<Payload>` 对象。这里的实现逻辑是把请求数据的字符串变成包含单个字符的流。客户端的 RSocket 对象使用 `requestStream()` 来发送请求，得到的是 `Flux<Payload>` 对象。

##### 清单 3\. 请求-响应流模式示例

```
public class RequestStreamExample {

public static void main(String[] args) {
    RSocketFactory.receive()
        .acceptor(((setup, sendingSocket) -> Mono.just(
            new AbstractRSocket() {
              @Override
              public Flux<Payload> requestStream(Payload payload) {
                return Flux.fromStream(payload.getDataUtf8().codePoints().mapToObj(c -> String.valueOf((char) c))
                    .map(DefaultPayload::create));
              }
            }
        )))
        .transport(TcpServerTransport.create("localhost", 7000))
        .start()
        .subscribe();

    RSocket socket = RSocketFactory.connect()
        .transport(TcpClientTransport.create("localhost", 7000))
        .start()
        .block();

    socket.requestStream(DefaultPayload.create("hello"))
        .map(Payload::getDataUtf8)
        .doOnNext(System.out::println)
        .blockLast();

    socket.dispose();
}
}

```

Show moreShow more icon

### 发后不管模式

发后不管模式的用法和之前的两种模式也是相似的。在代码清单 4 中， `AbstractRSocket` 类的实现覆写了 `fireAndForget()` 方法，对于请求的 Payload 对象，只需要返回 `Mono<Void>` 对象即可。客户端 RSocket 对象使用 `fireAndForget()` 方法发送请求。在发后不管模式中，由于发送方不需要等待接收方的响应，因此当程序结束时，服务器端并不一定接收到了请求。

##### 清单 4\. 发后不管模式示例

```
public class FireAndForgetExample {

public static void main(String[] args) {
    RSocketFactory.receive()
        .acceptor(((setup, sendingSocket) -> Mono.just(
            new AbstractRSocket() {
              @Override
              public Mono<Void> fireAndForget(Payload payload) {
                System.out.println("Receive: " + payload.getDataUtf8());
                return Mono.empty();
              }
            }
        )))
        .transport(TcpServerTransport.create("localhost", 7000))
        .start()
        .subscribe();

    RSocket socket = RSocketFactory.connect()
        .transport(TcpClientTransport.create("localhost", 7000))
        .start()
        .block();

    socket.fireAndForget(DefaultPayload.create("hello")).block();
    socket.fireAndForget(DefaultPayload.create("world")).block();

    socket.dispose();
}
}

```

Show moreShow more icon

### 通道模式

通道模式同样以相似的方式实现。在代码清单 5 中，服务器端的 `AbstractRSocket` 类的实现覆写了 `requestChannel()` 方法，对于请求的 `Publisher<Payload>` 对象，返回 `Flux<Payload>` 对象。客户端 RSocket 对象使用 `requestChannel()` 方法发送请求并处理响应。

##### 清单 5\. 通道模式示例

```
public class RequestChannelExample {

public static void main(String[] args) {
    RSocketFactory.receive()
        .acceptor(((setup, sendingSocket) -> Mono.just(
            new AbstractRSocket() {
              @Override
              public Flux<Payload> requestChannel(Publisher<Payload> payloads) {
                return Flux.from(payloads).flatMap(payload ->
                    Flux.fromStream(
                        payload.getDataUtf8().codePoints().mapToObj(c -> String.valueOf((char) c))
                            .map(DefaultPayload::create)));
              }
            }
        )))
        .transport(TcpServerTransport.create("localhost", 7000))
        .start()
        .subscribe();

    RSocket socket = RSocketFactory.connect()
        .transport(TcpClientTransport.create("localhost", 7000))
        .start()
        .block();

    socket.requestChannel(Flux.just("hello", "world", "goodbye").map(DefaultPayload::create))
        .map(Payload::getDataUtf8)
        .doOnNext(System.out::println)
        .blockLast();

    socket.dispose();
}
}

```

Show moreShow more icon

## Spring 集成

Spring 框架提供了对 RSocket 的集成，作为 spring-messaging 模块支持的一种消息传输方式。对于 Spring Boot 应用来说，只需要添加对 `spring-boot-starter-rsocket` 的依赖即可。本文使用的 Spring Boot 是 2.2.0.M6 版本，对应于 Spring 框架的 5.2.0.RC2 版本。代码清单 6 给出了 Spring Boot 应用中使用 RSocket 需要添加的 Maven 依赖。

##### 清单 6\. Spring Boot 的 RSocket 依赖

```
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-rsocket</artifactId>
    <version>2.2.0.M6</version>
</dependency>
<dependency>
<groupId>io.projectreactor</groupId>
    <artifactId>reactor-test</artifactId>
<version>3.3.0.RC1</version>
    <scope>test</scope>
</dependency>

```

Show moreShow more icon

### 消息处理控制器

在 Spring 应用的 `application.properties` 文件中，添加 `spring.rsocket.server.port=7100` 来设置内置 RSocket 服务器的端口。代码清单 7 中的 `EchoController` 是代码清单 2 中的请求-响应模式的实现。方法 `echo()` 接收 String 类型的请求，并返回 Mono作为响应。 `@MessageMapping("echo")` 注解指定了所处理消息的目的地。与 Spring 集成之后，RSocket 的使用变得很简洁。

##### 清单 7\. Spring 中使用 RSocket 的请求-响应模式

```
@Controller
public class EchoController {

@MessageMapping("echo")
public Mono<String> echo(String input) {
    return Mono.just("ECHO >> " + input);
}
}

```

Show moreShow more icon

### 单元测试

下面添加测试用例来进行测试。代码清单 8 中的 `AbstractTest` 类是一个抽象的 RSocket 测试类。其中包含的 `createRSocketRequester()` 方法用来创建发送 RSocket 请求的 `RSocketRequester` 对象。 `RSocketRequester.Builder` 类用来创建 `RSocketRequester` 对象，其中 `dataMimeType()` 方法指定负载中数据的 MIME 类型， `connect()` 方法指定连接的传输层实现。

##### 清单 8\. 抽象的 RSocket 测试类

```
import org.springframework.messaging.rsocket.RSocketRequester;
import org.springframework.util.MimeTypeUtils;

abstract class AbstractTest {

@Value("${spring.rsocket.server.port}")
private int serverPort;
@Autowired
private RSocketRequester.Builder builder;

RSocketRequester createRSocketRequester() {
    return builder.dataMimeType(MimeTypeUtils.TEXT_PLAIN)
        .connect(TcpClientTransport.create(serverPort)).block();
}
}

```

Show moreShow more icon

代码清单 9 中的 `EchoServerTest` 类测试代码清单 7 中的 `EchoController` 。在测试中，首先使用 `createRSocketRequester()` 方法来创建 `RSocketRequester` 对象，再使用 `route()` 方法来指定消息的目的地，最后使用 `data()` 方法指定负载中的数据。 `retrieveMono()` 方法用来发送请求。通过 StepVerifier 类来验证返回的 `Mono<String>` 对象。

##### 清单 9\. 测试 EchoController

```
@SpringBootTest
class EchoServerTest extends AbstractTest {

@Test
@DisplayName("Test echo server")
void testEcho() {
    RSocketRequester requester = createRSocketRequester();
    Mono<String> response = requester.route("echo")
        .data("hello")
        .retrieveMono(String.class);
    StepVerifier.create(response)
        .expectNext("ECHO >> hello")
        .expectComplete()
        .verify();
}

}

```

Show moreShow more icon

### 不同消息处理模式

除了请求-响应模式之外，其他模式也可以在 Spring 中使用。代码清单 10 中的 `StringSplitController` 类使用了请求-响应流模式。

##### 清单 10\. Spring 中使用 RSocket 的请求-响应流模式

```
@Controller
public class StringSplitController {
@MessageMapping("stringSplit")
public Flux<String> stringSplit(String input) {
    return Flux.fromStream(input.codePoints().mapToObj(c -> String.valueOf((char) c)));
}
}

```

Show moreShow more icon

代码清单 11 是代码清单 10 对应的测试用例。

##### 清单 11\. 测试 StringSplitController

```
@SpringBootTest
public class StringSplitTest extends AbstractTest {
@Test
@DisplayName("Test string split")
void testStringSplit() {
    RSocketRequester requester = createRSocketRequester();
    Flux<String> response = requester.route("stringSplit")
        .data("hello")
        .retrieveFlux(String.class);

    StepVerifier.create(response)
        .expectNext("h", "e", "l", "l", "o")
        .expectComplete()
        .verify();
}
}

```

Show moreShow more icon

## WebSocket 集成

如果需要使用 WebSocket 作为传输层实现，只需要替换 RSocket 服务器和客户端使用的传输层 Java 类即可，其他代码并不需要改动。在代码清单 12 中， `WebsocketServerTransport` 类是 WebSocket 服务器端实现，而 `WebsocketClientTransport` 类是客户端的实现。除了 `transport()` 方法使用的参数不同之外，其他的代码都与代码 [清单 2\. 请求-响应模式示例](#清单-2-请求-响应模式示例) 相同。

##### 清单 12\. 使用 WebSocket 的请求-响应模式示例

```
public class WebSocketRequestResponseExample {

public static void main(String[] args) {
    RSocketFactory.receive()
        .acceptor(((setup, sendingSocket) -> Mono.just(
            new AbstractRSocket() {
              @Override
              public Mono<Payload> requestResponse(Payload payload) {
                return Mono.just(DefaultPayload.create("ECHO >> " + payload.getDataUtf8()));
              }
            }
        )))
        .transport(WebsocketServerTransport.create("localhost", 7000))
        .start()
        .subscribe();

    RSocket socket = RSocketFactory.connect()
        .transport(WebsocketClientTransport.create("localhost", 7000))
        .start()
        .block();

    socket.requestResponse(DefaultPayload.create("hello"))
        .map(Payload::getDataUtf8)
        .doOnNext(System.out::println)
        .block();

    socket.dispose();
}
}

```

Show moreShow more icon

对于使用 Spring WebFlux 的应用，可以配置使用 RSocket 作为 WebSocket 服务器端实现。在代码清单 13 中的 Spring 配置中，在路径 `"/ws"` 上启用了基于 RSocket 的 WebSocket 支持。

##### 清单 13\. Spring WebFlux 中使用 RSocket 作为 WebSocket 实现

```
@Configuration
@EnableWebFlux
public class Config {

@Bean
RSocketWebSocketNettyRouteProvider rSocketWebsocketRouteProvider(
      RSocketMessageHandler messageHandler) {
    return new RSocketWebSocketNettyRouteProvider("/ws",
        messageHandler.responder());
}

static class RSocketWebSocketNettyRouteProvider implements NettyRouteProvider {

    private final String mappingPath;

    private final SocketAcceptor socketAcceptor;

    RSocketWebSocketNettyRouteProvider(String mappingPath, SocketAcceptor socketAcceptor) {
      this.mappingPath = mappingPath;
      this.socketAcceptor = socketAcceptor;
    }

    @Override
    public HttpServerRoutes apply(HttpServerRoutes httpServerRoutes) {
      ServerTransport.ConnectionAcceptor acceptor = RSocketFactory.receive()
          .acceptor(this.socketAcceptor)
          .toConnectionAcceptor();
      return httpServerRoutes.ws(this.mappingPath, WebsocketRouteTransport.newHandler(acceptor));
    }

}

}

```

Show moreShow more icon

由于 RSocket 有自己的二进制协议，在浏览器端的实现需要使用 RSocket 提供的 JavaScript 客户端与服务器端交互。在 Web 应用中使用 RSocket 提供的 NodeJS 模块 `rsocket-websocket-client` 即可。本文 [下载示例代码](#下载示例代码) 中包含一个 React 应用，连接到使用 RSocket 的 WebSocket 服务器并发送消息。代码清单 14中的 MessageService 类负责与服务器交互。首先创建一个 JavaScript 实现的 RSocket 客户端中的 `RSocketClient` 对象，并使用 `RSocketWebSocketClient` 作为传输层实现。在连接成功之后，使用 `RSocketClient` 的 `requestStream()` 方法发送消息并处理响应。这里的处理逻辑是调用提供的消息回调方法 `messageCallback` 。

##### 清单 14\. 连接使用 RSocket的WebSocket 服务器的 JavaScript 示例

```
import {RSocketClient, MAX_STREAM_ID} from 'rsocket-core';
import RSocketWebSocketClient from 'rsocket-websocket-client';

export default class MessageService {
constructor(connectCallback, messageCallback) {
    this._client = new RSocketClient({
      setup: {
        keepAlive: 10000,
        dataMimeType: 'text/plain',
        metadataMimeType: 'text/plain',
      },
      transport: new RSocketWebSocketClient({
        // eslint-disable-next-line no-restricted-globals
        url: `ws://${location.host}/ws`
      }),
    });
    this._connectCallback = connectCallback;
    this._messageCallback = messageCallback;
}

connect() {
    this._client.connect().then(
      (socket) => {
        this._socket = socket;
        this._connectCallback(null);
      },
      (error) => this._connectCallback(error),
    );
}

send(message) {
    this._socket.requestStream({
      data: message,
    }).subscribe({
      onNext: (payload) => this._messageCallback(null, payload.data),
      onError: (error) => this._messageCallback(error),
      onSubscribe: (_subscription) => _subscription.request(MAX_STREAM_ID),
    });
}
}

```

Show moreShow more icon

## RSocket 进阶

下面介绍 RSocket 相关的高级话题。

### 调试

由于 RSocket 使用二进制协议，所以调试 RSocket 应用消息比 HTTP/1 协议要复杂一些。从 RSocket 帧的二进制内容无法直接得知帧的含义。需要辅助工具来解析二进制格式消息。对 Java 应用来说，只需要把日志记录器 `io.rsocket.FrameLogger` 设置为 `DEBUG` 级别，就可以看到每个 RSocket 帧的内容。代码清单 15 给出了类型为 `REQUEST_RESPONSE` 帧的调试信息。除此之外，还可以使用 Wireshark 工具及其 RSocket 插件。

##### 清单 15\. FrameLogger 输出的调试信息

```
2019-09-17 05:02:59.906 DEBUG 15532 --- [actor-tcp-nio-1]
io.rsocket.FrameLogger                   : sending ->
Frame => Stream ID: 1 Type: REQUEST_RESPONSE Flags: 0b100000000 Length: 23
Metadata:
         +-------------------------------------------------+
         |  0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f |
+--------+-------------------------------------------------+----------------+
|00000000| fe 00 00 05 04 65 63 68 6f                      |.....echo       |
+--------+-------------------------------------------------+----------------+
Data:
         +-------------------------------------------------+
         |  0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f |
+--------+-------------------------------------------------+----------------+
|00000000| 68 65 6c 6c 6f                                  |hello           |
+--------+-------------------------------------------------+----------------+

```

Show moreShow more icon

### 负载零拷贝

在代码清单 2 的 `AbstractRSocket` 类的实现中，对于接收的 Payload 对象是直接使用的。这是因为 RSocket 默认对请求的负载进行了拷贝。这样的做法在实现上虽然简单，但会带来性能上的损失，增加响应时间。为了提高性能，可以通过 `ServerRSocketFactory` 类或 `ClientRSocketFactory` 类的 `frameDecoder()` 方法来指定 PayloadDecoder 接口的实现。 `PayloadDecoder.ZERO_COPY` 是内置提供的零拷贝实现类。当使用了负载零拷贝之后，负载的内容不再被拷贝。需要通过 Payload 对象的 `release()` 方法来手动释放负载对应的内存，否则会造成内存泄漏。如果使用 Spring Boot 提供的 RSocket 支持， `PayloadDecoder.ZERO_COPY` 默认已经被启用，并由 Spring 负责相应的内存释放。

## 结束语

RSocket 协议支持四种不同的交互模式，适用于不同类型的应用场景。相对于传统的 HTTP 协议来说，RSocket 提供了更多的灵活性。RSocket 可以使用 TCP、WebSocket 和 Aeron 作为传输层实现，简化了应用层协议的开发。此外，RSocket 提供的基于反应式流语义的流控制，保证了应用之间消息传递的健壮性。在分布式系统集成中，RSocket 是很好的选择。

## 下载示例代码

本文示例代码可以 [在 GitHub 上下载](https://github.com/VividcodeIO/rsocket-starter)。