# Ratpack：构建简单高效的 HTTP 微服务
Ratpack 的概念、基本的 HTTP 请求处理、异步和非阻塞的实现方式、响应式流处理和测试等

**标签:** Java,Web 开发,微服务

[原文链接](https://developer.ibm.com/zh/articles/j-lo-ratpack-http-microservice/)

成富

发布: 2016-03-14

* * *

微服务（microservice）是目前流行的软件架构模式。在微服务架构中，每个服务只实现特定的功能。不同的服务的实现是相互独立的，可以使用不同的数据库，甚至是不同的编程语言。在处理一个特定的请求时，可能需要多个微服务的参与和交互。微服务之间的通讯方式一般是 TCP/HTTP，采用的数据传输格式可以是 Protocol Buffer 或 JSON。通过 HTTP 的方式传输 JSON 数据，是一种常见的微服务集成方式。这要求每个微服务实现都暴露 HTTP 接口。暴露 HTTP 服务并不是一件复杂的事情，已经有非常多的开源框架和库可供利用。从简单的 Servlet 实现，到复杂的 JAX-RS 规范或 Spring MVC 框架。已有的这些框架和库的问题在于它们并不是为了实现微服务而设计的。它们的设计目标是开发传统的复杂 Web 应用。因此这些框架所包含的内容很多，运行起来的内存开销较大。Ratpack 是一个创建 HTTP 应用的轻量级库，可以作为创建 HTTP 微服务的基础。本文将对 Ratpack 进行详细介绍。

## Ratpack 简介

Ratpack 基于流行的网络开发库 Netty 来开发，在 Netty 的基础上添加了 HTTP 协议相关的内容。Netty 作为一个高性能高吞吐量的网络开发库，为 Ratpack 的简单高效打下了良好的基础。与传统的 Web 开发框架相比，Ratpack 所提供的功能相对较少。但是 Ratpack 所缺失的这些功能，在微服务开发中通常也是不需要的。从这个角度来说，Ratpack 是作为一个库，而不是框架来使用的。Ratpack 的作用类似于 NodeJS 社区中的 Connect。

在 Ratpack 编程中有两个重要的概念，分别是上下文对象（Context）和处理器（Handler）。处理器可以看成是在上下文对象上进行处理的函数。处理器在进行处理时，需要从上下文对象中获取所需的信息，如 HTTP 的请求和响应对象，或者通过上下文对象进行错误处理和重定向等。在创建 Ratpack 服务器时需要提供一个处理器的实现。这个处理器实际上就是 Ratpack 应用本身。一个处理器实现可以代理给其他处理器来完成其工作。多个代理器可以嵌套和链接起来，形成复杂的处理器链结构。使用 Ratpack 开发 HTTP 微服务，实际上就是开发处理微服务请求的处理器链。

### 上下文对象

Ratpack 的 Context 接口表示处理器当前的处理上下文。在处理器执行时，当前的 Context 对象会被作为唯一的参数传递给处理器。Context 接口中所提供的功能非常繁多，主要有下面几个类别。

首先是可以获取当前 HTTP 连接的请求和响应对象。通过 Context 接口的 getRequest 方法可以获取到表示当前 HTTP 请求的 Request 接口对象。Request 接口提供了一系列方法用来获取与 HTTP 请求相关的信息，如 getBody 方法来读取请求的内容，getHeaders 方法来获取 HTTP 头，getQuery 方法来获取查询字符串等。Context 接口的 getResponse 方法可以获取到作为 HTTP 请求响应的 Response 接口对象。Response 接口的方法用来对响应进行操作，如 status 方法设置 HTTP 响应状态码，getHeaders 方法返回可以进行修改的 HTTP 响应头，cookie 方法用来设置响应的 Cookie 信息。当 HTTP 响应配置完成之后，通过 send 方法来发送响应。Response 可以发送不同类型的响应，如 byte[]、String 类型或文件等。

Context 接口的另外一个重要功能是作为上下文对象的注册表。Context 的父接口 Registry 提供了与对象注册和获取相关的方法。在 Ratpack 应用开发中，处理器的实现可能会需要用到相关的上下文对象。这些对象都保存在 Context 中。比如一个获取用户订单的微服务，可能由一个处理器负责根据订单号进行数据库查询，而另外一个处理器则负责处理订单对象并生成 JSON 格式的响应。前一个处理器可以把获取到的订单对象保存到 Context 中，而后一个处理器则直接从 Context 中获取订单对象。通过这样的职责划分，前一个处理器所执行的逻辑可以在不同的地方被复用。

Context 接口还提供了 insert 和 next 方法来构建复杂的处理器链。insert 方法负责插入一个或多个处理器到当前处理链中，并代理给第一个处理器；next 方法则直接代理给处理链中的下一个处理器。

在 [清单 1\. Context 中的处理器链](#清单-1-context-中的处理器链) 中，首先在 Context 中添加了一个 ArrayList 类型的对象。addOutput 方法的实现中向该 ArrayList 对象中添加一个字符串，然后通过 Context 的 next 方法代理给下一个处理器来执行。通过 Context 的 insert 方法添加了 3 个不同的处理器，前两个处理器用来修改 Context 中的对象，最后一个处理器用来把结果以 JSON 格式输出。

##### 清单 1\. Context 中的处理器链

```
public class HandlerChain {
public static void main(String[] args) throws Exception {
new HandlerChain().start();
}

public void start() throws Exception {
RatpackServer.start(server ->
server.registry(Registry.single(new ArrayList<String>()))
.handlers(chain ->
chain.get(ctx -> ctx.insert(
addOutput("Hello"),
addOutput("World"),
ctx3 -> ctx3.render(json(ctx3.get(List.class)))
))
)
);
}

private Handler addOutput(final String text) {
return ctx -> {
ctx.get(List.class).add(text);
ctx.next();
};
}
}

```

Show moreShow more icon

Context 接口中还有一些辅助方法，如 clientError 来返回一个与客户端错误相关的状态码，error 来返回服务器错误相关的状态码，notFound 来返回 404 错误，redirect 来返回重定向的响应。Context 还提供了 parse 方法来解析 HTTP 请求内容。对于 HTTP 响应，Context 接口提供了 render 方法来把不同类型的对象作为响应的内容。

### 处理器

处理器 Handler 的接口定义很简单，只有一个 handle 方法。handle 方法只有唯一的参数是 Context 接口的对象。每个 Handler 负责对 Context 进行操作。每个处理器所执行的操作可以简单或复杂。单个处理器可以获取 HTTP 请求的内容，调用服务层的代码来执行业务逻辑，并生成相应的响应。也可以只完成某个简单的任务，再通过 Context 接口的 insert 和 next 方法来代理给其他处理器继续处理。从可复用和可测试的角度出发，每个处理器的实现应该尽可能的简单。这样可以充分利用 Ratpack 提供的处理器代理和链接功能。

## 基本的 HTTP 处理

在开发 HTTP 微服务时，某些基本的操作是不可或缺的。Ratpack 也通过 Request 和 Response 接口对这些操作提供了支持。

### 表单和文件上传

表单处理是 HTTP 应用中的一个重要内容。一般的 HTML 表单提交时的内容类型是 application/x-www-form-urlencoded。当有文件上传时，使用的是 multipart 的内容格式。代码清单演示了一般的表单处理和文件上传。基本的做法是使用 Context 接口的 parse 方法把请求内容解析成 Form 接口的对象。对于一般的表单元素，使用 Form 接口的 get 方法可以获取到给定名称的元素的值；对于上传文件，通过 file 方法来获取到给定名称的 UploadedFile 接口对象。通过 UploadedFile 接口可以读取文件的内容。 [清单 2\. 表单和文件上传](#清单-2-表单和文件上传) 中直接把上传文件的文本内容返回给客户端。

##### 清单 2\. 表单和文件上传

```
RatpackServer.start(server -> server
.handlers(chain -> chain.post("form",
ctx -> ctx.parse(Form.class)
.then(form -> ctx.render(form.get("name")))
).post("file", ctx ->
ctx.parse(Form.class).then(form -> ctx.render(form.file("file").getText()))
)
)
);

```

Show moreShow more icon

### JSON 请求和响应

HTTP 微服务通常使用 JSON 作为数据的传输格式，因此需要对 JSON 格式的请求内容进行解析，并且生成 JSON 格式的响应。解析 JSON 请求内容同样需要使用 Context 接口的 parse 方法；而生成 JSON 响应，也是使用 Context 接口的 render 方法。Ratpack 的 JSON 支持使用的是 Jackson 库，支持 POJO 对象和 JSON 格式之间的转换。

在 [清单 3\. JSON 请求和响应](#清单-3-json-请求和响应) 中，通过 Jackson 类的 fromJson 方法可以把 HTTP 请求内容转换成 User 类的对象。当需要生成 JSON 响应时，只需要用 Jackson 类的 json 方法包装 User 类的对象即可。

##### 清单 3\. JSON 请求和响应

```
RatpackServer.start(server ->
server.handlers(chain ->
chain.post(ctx -> {
ctx.parse(fromJson(User.class))
.map(user -> new User(user.getId(), user.getName().toUpperCase()))
.then(user -> ctx.render(json(user)));
})
)
);

```

Show moreShow more icon

### 静态文件

在 HTTP 微服务开发中，可能会需要直接发送磁盘上的静态文件，如 HTML、JavaScript、CSS 和图片文件等。Ratpack 提供了对静态文件的直接支持，可以很方便的把某个目录暴露为静态资源。在 [清单 4\. 静态文件处理](#清单-4-静态文件处理) 中，使用 Chain 接口的 files 方法来声明静态资源的特征，即目录是 public，首页文件是 index.html。在启动 Ratpack 服务器时，需要配置服务器的根目录。这是通过 baseDir 方法来设置的。files 方法中的静态资源声明中的路径都是相对于该根目录的。

##### 清单 4\. 静态文件处理

```
public class StaticFiles {
public static void main(String[] args) throws Exception {
Path baseDir = Paths.get("app").toAbsolutePath();
new StaticFiles().start(baseDir);
}

public void start(final Path baseDir) throws Exception {
RatpackServer.start(server -> server
.serverConfig(c -> c.baseDir(baseDir))
.handlers(chain -> chain
.files(f -> f.dir("public").indexFiles("index.html"))
)
);
}
}

```

Show moreShow more icon

Ratpack 的静态文件支持提供了基于文件修改时间的 ETag 头，以及 GZIP 的支持。但是相比于 Apache 和 Nginx 来说，所提供的功能还是很弱。因此 Ratpack 的静态文件支持一般用在开发和测试环境中，并不适合直接在生产环境中使用。

除了直接暴露静态文件目录之外，还可以在处理器中通过 Response 接口的 sendFile 来发送文件作为响应。

### HTTP 头

Request 和 Response 接口都提供了 getHeaders 方法来获取到 HTTP 头，不过 Request 接口的 getHeaders 方法返回的是不可变的 Headers 接口对象，而 Response 接口的 getHeaders 方法返回的是可变的 MutableHeaders 接口对象。Headers 接口中的 get 和 getAll 方法可以根据名称获取 HTTP 头的值。MutableHeaders 接口提供了 add、set 和 remove 方法来对 HTTP 头进行修改。
[清单 5\. HTTP 头处理](#清单-5-http-头处理) 中获取了 HTTP 请求中 X-UID 头的内容，并作为 HTTP 响应的内容发送回来。

##### 清单 5\. HTTP 头处理

```
RatpackServer.start(server ->
server.handlers(chain ->
chain.get(ctx -> {
String userId = ctx.getRequest().getHeaders().get("X-UID");
ctx.getResponse().getHeaders().set("X-Poweredby", "Ratpack");
ctx.render(userId);
})
)
);

```

Show moreShow more icon

### Cookie

Cookie 也是 HTTP 处理中的重要部分。Request 接口中的 getCookies 方法可以获取到全部 Cookie 的集合，oneCookie 方法可以根据名称获取某个 Cookie 的值。Response 接口的 cookie 方法用来创建一个 Cookie。

[清单 6\. Cookie 处理](#清单-6-cookie-处理) 中使用 Request 接口的 oneCookie 方法来读取 Cookie 的值，使用 Response 接口的 cookie 方法来设置 Cookie 的值。

##### 清单 6\. Cookie 处理

```
RatpackServer.start(server ->
server.handlers(chain ->
chain.get(ctx -> {
String userId = ctx.getRequest().oneCookie("userId");
ctx.getResponse().cookie("userId", "value");
ctx.render(userId);
})
)
);

```

Show moreShow more icon

## Spring 集成

Ratpack 提供了与 Spring 框架的集成。集成的方式有两种。第一种是在 Ratpack 应用中使用 Spring 框架提供的依赖注入功能。之前介绍过 Ratpack 中的 Registry 接口提供了对象的注册和查找功能。这个功能可以用 Spring 来替代，也就是由 Spring 来负责 Registry 中对象的创建和管理。这样就可以充分利用 Spring 强大的依赖注入功能来简化 Ratpack 应用的开发。

在 [清单 7\. 在 Ratpack 中使用 Spring 的依赖注入功能](#清单-7-在-ratpack-中使用-spring-的依赖注入功能) 中，Greeting 是一个作为示例的接口。AppConfiguration 是 Spring 框架使用的配置类，在其中定义了一个 Greeting 接口的实现对象。在启动 Ratpack 服务器时，通过 registry 方法来指定一个 Registry 接口的实现，而 ratpack.spring.Spring.spring 方法则把一个 Spring 的配置类转换成 Registry 对象。这样在处理器的实现中就可以从 Context 中查找到 Greeting 接口的实现对象并使用。

##### 清单 7\. 在 Ratpack 中使用 Spring 的依赖注入功能

```
RatpackServer.start(server -> server
.registry(spring(AppConfiguration.class))
.handlers(chain -> chain
.get(ctx -> ctx
.render(ctx.get(Greeting.class).say("Alex")))
)
);

public interface Greeting {
String say(final String name);
}

@Configuration
class AppConfiguration {
@Bean
public Greeting greeting() {
return name -> String.format("Hello, %s", name);
}
}

```

Show moreShow more icon

另外一种与 Spring 集成的方式是在 Spring Boot 中把 Ratpack 作为 Servlet 容器的替代。在 [清单 8\. 在 Spring Boot 中使用 Ratpack](#清单-8-在-spring-boot-中使用-ratpack) 中，在 Spring Boot 启动类中添加注解 @EnableRatpack 来启用 Ratpack。Ratpack 中处理器的逻辑定义在 Action类型的对象中。这与代码清单中启动 Ratpack 服务器时使用的 handlers 方法的作用是相同的。

##### 清单 8\. 在 Spring Boot 中使用 Ratpack

```
@SpringBootApplication
@EnableRatpack
public class SpringBootApp {
@Bean
public Action<Chain> index() {
return chain -> chain
.get(ctx -> ctx
.render(greeting().say("Alex"))
);
}

public Greeting greeting() {
return name -> String.format("Hello, %s", name);
}

public static void main(String... args) throws Exception {
SpringApplication.run(SpringBootApp.class, args);
}
}

```

Show moreShow more icon

## 异步和非阻塞

Ratpack 从其底层设计和实现都是异步的，极大的减少线程阻塞在 I/O 上的等待时间，从而提升应用的性能。Ratpack 的底层 Netty 提供了事件驱动和非阻塞的 HTTP I/O，类似 NodeJS。Ratpack 库提供的请求处理也是异步完成的。相比于传统的 Servlet 框架的同步方式，Ratpack 的异步方式的性能更优，所耗费的系统资源较少，不过对习惯了同步方式的开发人员来说，需要一定时间的学习和适应。

Promise 是 Ratpack 中代码执行的重要接口，表示的是一个会在将来变得可用的值。Ratpack 中的 Promise 的概念，与 JavaScript 框架常用的 Promise 的概念是相同的。当 Promise 中所表示的值变得可用时，该 Promise 上关联的回调函数会被调用。代码清单中已经展示了 Promise 的用法。Context 的 parse 返回的就是一个 Promise 对象。通过 Promise 的 then 方法可以添加一个以 Action 接口来表示的回调函数。当 parse 完成解析时，所得到的结果会被传递给回调函数。在使用 Ratpack 库时，经常会看到类似这样的代码结构。

在 Ratpack 应用中可能需要调用第三方服务，而这些服务本身是同步执行的，如调用一个第三方的 Web 服务。在这种情况下，Ratpack 提供了一个单独的线程池，用来执行这些同步操作，并提供了相关的 API 来把同步操作转换成异步操作。Blocking 对象的 get 方法可以把一个同步操作转换成 Promise 对象，从而可以与 Ratpack 应用中的其他部分进行整合。

在 [清单 9\. 在 Ratpack 中使用同步方法](#清单-9-在-ratpack-中使用同步方法) 中，Blocking 的 get 方法从一个 Factory 函数接口中获取到所需的结果。这里通过线程睡眠来模拟耗时较长的同步操作。由于 Blocking 的 get 方法返回的是 Promise 对象，可以直接通过 then 方法来对返回的结果进行处理。

##### 清单 9\. 在 Ratpack 中使用同步方法

```
public class BlockingGet {
public static void main(String[] args) throws Exception {
RatpackServer.start(server ->
server.handlers(chain ->
chain.get(ctx -> Blocking.get(() -> {
Thread.sleep(5000);
return "Hello World";
}).then(str -> ctx.render(str)))
)
);
}
}

```

Show moreShow more icon

需要注意的是，Ratpack 的 Promise 可能会造成所谓的回调函数地狱问题（callback hell），即过多层次的回调函数嵌套。对于这样的情况，可以考虑使用 RxJava 来解决。

## 响应式流

Ratpack 提供了对流处理的功能。Ratpack 的流处理 API 基于标准的 Reactive Streams API。Reactive Streams 是 JVM 上进行非阻塞带背压（back pressure) 的异步流处理的规范。Reactive Streams 的一个重要特征是支持带背压的流控制，其核心思想是由流的消费者来通知流的生产者其所能处理的数据量。这样可以避免处理速度较慢的消费者占用生产者的过多资源。

Reactive Streams API 中定义了 3 个最基本的接口，Publisher、Subscriber 和 Subscription。Publisher 是数据的发送者，其中的 subscribe 方法允许数据的消费者 Subscriber 进行注册并开始消费数据。Subscriber 是数据的消费者，其中定义了不同的事件回调方法。当在 Publisher 上注册成功之后，Subscriber 接口的 onSubscribe 方法会被调用，表示当前注册的 Subscription 接口的对象作为参数传入。当有数据可以消费时，Subscriber 接口的 onNext 方法会被调用，产生的数据作为 onNext 的参数传入；当不再有任何数据可用时，onComplete 方法会被调用；当产生数据出现错误时，onError 方法会被调用，与错误相关的 Throwable 对象作为参数传入。前面提到的支持带背压的流控制体现在 Subscription 接口中。Subscription 接口的 request 方法用来通知 Publisher 发送指定数量的数据，这样 Subscriber 就可以只要求它所能处理的数据量。cancel 方法用来通知 Publisher 接口停止发送数据。

Ratpack 对响应式流的支持体现在可以直接使用响应式流作为 HTTP 的响应，服务器推送事件（Server Sent Events）和 WebSocket 的数据源。以作为 HTTP 的响应为例，Ratpack 可以把响应式流的内容以 HTTP 分块传输编码（Chunk Encoding）的格式进行发送。

[清单 10\. 使用响应式流读取文件并发送](#清单-10-使用响应式流读取文件并发送) 中展示了一个从文件系统中读取文件内容并作为 HTTP 请求的响应的示例。为了代码可以工作，需要额外添加 ratpack-rx 和 rxjava-reactive-streams 两个依赖包，其中 ratpack-rx 是 Ratpack 提供的与 RxJava 集成的库，而 rxjava-reactive-streams 是 RxJava 与 Reactive Streams API 之间的桥接库。首先使用 Files 类的 readAllLines 方法从文件中读取所有行来得到一个 List对象，再从该对象中创建一个 RxJava 的 Observable，再把该 Observable 对象转换成 Reactive Streams API 的 Publisher 对象。Ratpack 的 ratpack.http.ResponseChunks.stringChunks 方法可以把一个 Publisher 对象以 HTTP 分块编码的格式直接发送。

##### 清单 10\. 使用响应式流读取文件并发送

```
public class FileWriter {
public static void main(String[] args) throws Exception {
RxRatpack.initialize();
new FileWriter().start();
}

public void start() throws Exception {
RatpackServer.start(server ->
server
.serverConfig(c -> c.baseDir(Paths.get("app").toAbsolutePath()))
.handlers(chain ->
chain.get(ctx -> {
Publisher<String> publisher = RxReactiveStreams.toPublisher(
Observable.from(
Files.readAllLines(
ctx.getFileSystemBinding().file("largefile.txt"))));
ctx.render(stringChunks(publisher));
})
)
);
}
}

```

Show moreShow more icon

## 测试

HTTP 微服务的测试是非常重要的一环。Ratpack 对 HTTP 应用的单元测试和集成测试都提供了良好的支持。由于 Ratpack 的应用的逻辑都在处理器链中，只需要对处理器进行单元测试，就可以覆盖 Ratpack 应用测试的绝大部分内容。Ratpack 提供了单元测试所需的支持工具 RequestFixture。通过 RequestFixture 可以向处理器的实现发送指定内容的 HTTP 请求，再对返回的响应进行验证。

在 [清单 11\. 处理器单元测试](#清单-11-处理器单元测试) 中，SetHeaderHandler 是一个简单的处理器实现，根据请求中的 HTTP 头信息来生成新的 HTTP 头和响应。在实际的测试中，使用 RequestFixture 的 handle 方法来包装一个处理器的实现对象，同时对请求进行设置，接着对返回的结果 HandlingResult 进行验证。

##### 清单 11\. 处理器单元测试

```
public class HandlerTests {

public static class SetHeaderHandler implements Handler {
@Override
public void handle(Context context) throws Exception {
String userId = context.getRequest().getHeaders().get("X-UID");
String greeting = String.format("Hello, %s", userId);
context.getResponse().getHeaders().set("X-GREETING", greeting);
context.render(userId);
}
}

@Test
public void testSetHeader() throws Exception {
HandlingResult result = RequestFixture.handle(new SetHeaderHandler(), fixture -> {
fixture.header("X-UID", "Alex");
});

assertEquals("Alex", result.rendered(String.class));
assertEquals("Hello, Alex", result.getHeaders().get("X-GREETING"));
}
}

```

Show moreShow more icon

在集成测试方面，Ratpack 提供了 ApplicationUnderTest 来运行一个 Ratpack 应用，并提供对应的 TestHttpClient 来对该应用进行测试。在 [清单 12\. Ratpack 应用集成测试](#清单-12-ratpack-应用集成测试) 中，simpleServer 方法创建了一个简单的 RatpackServer 对象。在进行测试时，使用 ApplicationUnderTest 的 of 方法从该 RatpackServer 中创建一个 ApplicationUnderTest 对象，然后使用 getHttpClient 方法得到进行测试的 TestHttpClient 对象，就可以使用该对象来发送请求并验证结果。

##### 清单 12\. Ratpack 应用集成测试

```
public class IntegrationTests {
public RatpackServer simpleServer() throws Exception {
return RatpackServer.of(server -> {
server.handlers(chain -> {
chain.get(ctx -> {
ctx.render("Hello");
});
});
});
}

@Test
public void testSimpleServer() throws Exception {
ApplicationUnderTest application = ApplicationUnderTest.of(simpleServer());
TestHttpClient httpClient = application.getHttpClient();
assertEquals("Hello", httpClient.get().getBody().getText());
}
}

```

Show moreShow more icon

## 结束语

随着微服务架构的流行，会有越来越多的大型系统被重新组织成多个独立的微服务。这些微服务之间可以使用 HTTP 协议来通讯和集成。Ratpack 作为一个轻量级的 HTTP 应用开发库，可以开发出高效并且占用资源少的 HTTP 微服务。本文对 Ratpack 库的重要部分进行了详细的说明，包括上下文对象和处理器的概念，基本的 HTTP 请求处理，异步和非阻塞的实现方式，响应式流处理和对 Ratpack 应用进行测试等。

## 下载例代码

[source\_code.zip](http://www.ibm.com/developerWorks/cn/java/j-lo-ratack-http-microservice/source_code.zip): 示例代码

## 相关主题

- [Ratpack 官方网站](https://ratpack.io/)
- [微服务架构](https://baike.baidu.com/item/微服务架构)
- [Reactive Streams](https://www.reactive-streams.org/)
- [Reactive Extensions](https://github.com/ReactiveX)
- [RxJava](https://github.com/ReactiveX/RxJava)