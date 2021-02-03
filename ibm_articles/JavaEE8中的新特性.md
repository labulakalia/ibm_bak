# Java EE 8 中的新特性
抢先了解针对 Java 安全性、JSON 绑定和处理、HTTP/2 等方面的新 API 和特性

**标签:** Jakarta EE,Java

[原文链接](https://developer.ibm.com/zh/articles/j-whats-new-in-javaee-8/)

Alex Theedom

发布: 2017-10-24

* * *

万众期待的 Java™ EE 8 即将发布。这是自 2013 年 6 月以来 Java 企业平台的首次发布，是分成两部分的发布中的前半部分（后一部分是 Java EE 9）。Oracle 对 Java EE 进行了战略性重新定位，重点关注一些支持云计算、微服务和反应式编程的技术。反应式编程现在正融入到许多 Java EE API 的架构中，而 JSON 交换格式为核心平台提供了支撑。

我们将大体了解 Java EE 8 中提供的主要特性。重点包括 API 更新和引入、对 HTTP/2 的新支持、反应式编程，以及 JSON。我们首先会介绍 Java EE 规范和升级，它们无疑将决定企业 Java 编程未来几年的发展方向。

## 全新的 API 和更新的 API

Java EE 8 向核心 API 引入了一些主要和次要更新，比如 Servlet 4.0 和 Context and Dependency Injection 2.0。还引入了两个新的 API — Java API for JSON Binding ( [JSR 367](https://www.jcp.org/en/jsr/detail?id=367)) 和 Java EE Security API ( [JSR 375](https://www.jcp.org/en/jsr/detail?id=375))。我们首先介绍新的 API，然后探索对存在已久的 Java EE 规范的改动。

## JSON Binding API

新的 JSON Binding API (JSON-B) 支持在 Java 对象与 [兼容 RFC 7159 的](https://tools.ietf.org/html/rfc7159) JSON 之间执行序列化和反序列化，同时维护与 JAXB ( [Java API for XML Binding 2.0](https://jcp.org/en/jsr/detail?id=222)) 的一致性。它提供了从 Java 类和实例到符合公认约定的 JSON 文档的默认映射。

JSON-B xu 还允许开发人员自定义序列化和反序列化。您可以使用注释自定义各个类的这些流程，或者使用运行时配置构建器来开发自定义策略。后一种方法包括使用适配器来支持用户定义的自定义。JSON-B 与 Java API for JSON Processing (JSON-P) 1.1（本文后面将讨论）紧密集成。

### JSON Binding API 入口

两个接口为新的 JSON Binding API 提供了入口： `JsonbBinding` 和 `Jsonb` 。

- `JsonbBinding` 提供了 JSON Binding API 的客户端访问点。为此，它通过根据所设置的配置和参数来构建 `Jsonb` 实例。
- `Jsonb` 通过方法 `toJson()` 和 `fromJson()` 提供序列化和反序列化操作。

JSON-B 还可以向插入的外部 JSON Binding 提供者描述功能，所以您不受 API 附带的绑定逻辑的限制。

### 使用 JsonB 实现序列化和反序列化

清单 1 首先序列化 `Book` 类 `book` 的实例，然后将其反序列化。

##### 清单 1\. 序列化和反序列化的最简单示例

```
String bookJson = JsonbBuilder.create().toJson(book);
Book book = JsonbBuilder.create().fromJson(bookJson, Book.class);

```

Show moreShow more icon

静态 `create()` 工厂方法返回一个 `Jsonb` 实例。您可以在该实例上调用许多重载的 `toJson()` 和 `fromJson()` 方法。另请注意，该规范没有强制要求来回等效转换：在上面的示例中，将 `bookJson` 字符串传入 `fromJson()` 方法中可能不会反序列化为一个等效对象。

JSON-B 还支持采用与对象大致相同的方式来绑定集合类与原语和/或实例数组 — 包括多维数组。

### 自定义 Jsonb

可通过为字段、JavaBeans 方法和类添加注释，对 `Jsonb` 方法的默认行为进行自定义。

例如，可以使用 `@JsonbNillable` 和 `@JsonbPropertyOrder` 注释来自定义 null 处理和属性顺序，该顺序应在类级别上进行自定义：

##### 清单 2\. 自定义 Jsonb

```
@JsonbNillable
@JsonbPropertyOrder(PropertyOrderStrategy.REVERSE)
public class Booklet {

private String title;

@JsonbProperty("cost")
@JsonbNumberFormat("#0.00")
private Float price;

private Author author;

@JsonbTransient
public String getTitle() {
       return title;
}

@JsonbTransient
public void setTitle(String title) {
       this.title = title;
}

// price and author getters/setter removed for brevity
}

```

Show moreShow more icon

调用方法 `toJson()` 会生成清单 3 中所示的 JSON 结构。

##### 清单 3\. 自定义的 JSON 结构

```
{
"cost": "10.00",
"author": {
    "firstName": "Alex",
    "lastName": "Theedom"
}
}

```

Show moreShow more icon

也可以选择通过运行时配置构建器 `JsonbConfig` 来执行自定义：

##### 清单 4\. Jsonb 的运行时配置

```
JsonbConfig jsonbConfig = new JsonbConfig()
    .withPropertyNamingStrategy(
        PropertyNamingStrategy.LOWER_CASE_WITH_DASHES)
    .withNullValues(true)
    .withFormatting(true);

Jsonb jsonb = JsonbBuilder.create(jsonbConfig);

```

Show moreShow more icon

清单 4 将 JSON-B 配置为使用 `LOWER_CASE_WITH_DASHES` 约定，以便将 null 保留在它们所在的位置，并输出经过优化的 JSON。

### 开源绑定

前面已经提到过，您不需要使用现成的 JSON-B 选项。清单 5 展示了如何配置开源绑定实现：

##### 清单 5\. 开源绑定配置

```
JsonbBuilder builder = JsonbBuilder.newBuilder("aProvider");

```

Show moreShow more icon

## Java EE Security API

新 Java EE Security API 的引入，是为了更正各种 servlet 容器在解决安全问题的方式上的不一致。这个问题在 Java Web 配置文件中尤为突出，主要是因为 Java EE 仅规定了完整的 Java EE 配置文件必须如何实现标准 API。新规范还引入了一些现有 API 没有利用的现代功能，比如 CDI。

这个 API 的美妙之处在于，它提供了一种配置身份存储和身份验证机制的备选方法，但没有取代现有安全机制。开发人员应该非常希望拥有在 Java EE Web 应用程序中启用安全性的机会，无论是否使用特定于供应商的解决方案或专用的解决方案。

### 规范内容

Java EE Security API 规范解决了 3 个关键问题：

- `HttpAuthenticationMechanism` 支持对 servlet 容器执行身份验证。
- `IdentityStore` 标准化了 JAAS `LoginModule` 。
- `SecurityContext` 提供了一个实现编程安全的访问点。

下面将介绍上述每个组件。

### HttpAuthenticationMechanism

Java EE 已指定了两种验证 Web 应用程序用户的机制：Java Servlet 规范 3.1 (JSR-340) 指定了一种声明性的应用程序配置机制，而 JASPIC (Java Authentication Service Provider Interface for Containers) 定义了一个名为 `ServerAuthModule` 的 SPI，该 SPI 支持开发身份验证模块来处理任何凭证类型。

这两种机制既有意义又很有效，但从 Web 应用程序开发人员的角度讲，每种机制都存在局限性。servlet 容器机制仅支持小范围的凭证类型。JASPIC 非常强大和灵活，但使用起来也非常复杂。

Java EE Security API 希望通过一个新接口解决这些问题： `HttpAuthenticationMechanism` 。这个新接口实际上是 JASPIC `ServerAuthModule` 接口的一个简化的 servlet 容器变体，它在减少现有机制的局限性的同时利用现有机制。

`HttpAuthenticationMechanism` 类型的实例是一个 CDI bean，可将它用于实现注入操作的容器，而且该实例被指定仅用于 servlet 容器。该规范明确排除了 EJB 和 JMS 等其他容器。

`HttpAuthenticationMechanism` 接口定义了 3 个方法： `validateRequest()` 、 `secureResponse()` 和 `cleanSubject()` 。这些方法非常类似于 JASPIC `ServerAuth` 接口上声明的方法，所以开发人员应该对它们感到熟悉。唯一需要重写的方法是 `validateRequest()` ；其他所有方法都有默认实现。

### IdentityStore

_身份存储_ 是一个数据库，用于存储用户身份数据，比如用户名、组成员关系，以及用于验证凭证的信息。在新的 Java EE Security API 中，使用了一个名为 `IdentityStore` 的身份存储抽象来与身份存储进行交互。它的用途是验证用户并检索组成员信息。

正如规范中所写， `IdentityStore` 的意图是供 `HttpAuthenticationMechanism` 实现使用，不过这不是必须的。结合使用 `IdentityStore` 和 `HttpAuthenticationMechanism` ，使得应用程序能以一种便携的、标准的方式控制其身份存储。

### SecurityContext

`IdentityStore` 和 `HttpAuthenticationMechanism` 相结合，成为了一个强大的新的用户身份验证工具。然而，声明性模型可能无法满足系统级安全需求。 `SecurityContext` 可在这里派上用场：编程安全使得 Web 应用程序能执行所需的测试，从而允许或拒绝访问应用程序资源。

## 主要更新：Servlet 4.0、Bean Validation 2.0、CDI 2.0

Java EE 8 中主要提供了 3 个企业标准 API：Servlet 4.0 ( [JSR 369](https://jcp.org/en/jsr/detail?id=369))、Bean Validation 2.0 ( [JSR 380](https://jcp.org/en/jsr/detail?id=380)) 和 Contexts and Dependency Injection for Java 2.0 ( [JSR 365](https://jcp.org/en/jsr/detail?id=365))。

下面将介绍每个 API 的重要特性。

### Servlet 4.0

Java Servlet API 是 Java 企业开发人员最早接触、最熟悉的 API 之一。它于 1999 年在 J2EE 1.2 中首次面世，现在在 Java Server Pages (JSP)、JavaServer Faces (JSF)、JAX-RS 和 MVC ( [JSR 371](https://jcp.org/en/jsr/detail?id=371)) 中发挥着重要作用。

Java EE 8 中对 Servlet 进行了重大修订，主要是为了适应 [HTTP/2](https://www.ibm.com/developerworks/library/wa-http2-under-the-hood/index.html) 的性能增强特性。服务器推送目前是这一领域的首要特性。

**服务器推送是什么？**

_服务器推送_ 通过将客户端资源推送到浏览器的缓存中来预先满足对这些资源的需求。客户端发送请求并收到服务器响应时，所需的资源已在缓存中。

#### PushBuilder

在 Servlet 4.0 中，服务器推送是通过一个 `PushBuilder` 实例公开的。清单 6 展示了一个从 `HttpServletResponse` 实例获取的 `PushBuilder` 实例，该实例被传递到一个请求处理方法。

##### 清单 6\. servlet 中的 PushBuilder

```
@Override
protected void doGet(HttpServletRequest request,
           HttpServletResponse response)
           throws ServletException, IOException {

PushBuilder pushBuilder = request.newPushBuilder();
pushBuilder.path("images/header.png").push();
pushBuilder.path("css/menu.css").push();
pushBuilder.path("js/ajax.js").push();

// Do some processing and return JSP that
// requires these resources
}

```

Show moreShow more icon

在清单 6 中， `header.png` 的路径是通过 `path()` 方法在 `PushBuilder` 实例上设置的，并通过调用 `push()` 推送到客户端。该方法返回时，会清除路径和条件标头，以供构建器重用。

#### servlet 映射的运行时发现

Servlet 4.0 提供了一个新的 API，用它来实现 URL 映射的运行时发现。 `HttpServletMapping` 接口的用途是让确定导致 servlet 激活的映射变得更容易。在该 API 内，会从一个 `HttpServletRequest` 实例获得 servlet 映射，该实例包含 4 个方法：

- `getMappingMatch()` 返回匹配的类型。
- `getPattern()` 返回激活 servlet 请求的 URL 模式。
- `getMatchValue()` 返回匹配的 `String`
- `getServletName()` 返回通过该请求激活的 servlet 类的完全限定名称。

##### 清单 7\. HttpServletMapping 接口上的所有 4 个方法

```
HttpServletMapping mapping = request.getHttpServletMapping();
String mapping = mapping.getMappingMatch().name();
String value = mapping.getMatchValue();
String pattern = mapping.getPattern();
String servletName = mapping.getServletName();

```

Show moreShow more icon

除了这些更新之外，Servlet 4.0 还包含更细微的管理性更改和对 [HTTP Trailer](https://tools.ietf.org/html/rfc7230) 的支持。新的 `GenericFilter` 和 `HttpFilter` 类简化了过滤器的编写，实现了对 Java SE 8 的一般性改进。

### Bean Validation 2.0

Bean Validation 2.0 通过一系列新特性得到了增强，其中许多特性是 Java 开发人员社区所请求的。Bean 验证是一个横切关注点，所以 2.0 版规范希望确保数据在从客户端传输到数据库的过程中是完整的，并通过对字段值、返回值和方法参数应用约束来实现此目的。

这些增强中包含一些约束，它们可以验证电子邮箱地址，确保数字是正的或负的，测试日期是过去还是现在，并测试字段不是空的或 null。这些约束包括： `@Email` 、 `@Positive` 、 `@PositiveOrZero` 、 `@Negative` 、 `@NegativeOrZero` 、 `@PastOrPresent` 、 `@FutureOrPresent` 、 `@NotEmpty` 和 `@NotBlank` 。这些约束现在也可应用于更广泛的地方。例如，它们可以继续处理参数化类型的参数。添加了对按类型参数来验证容器元素的支持，如清单 8 所示。

##### 清单 8\. 按类型验证容器元素

```
private List<@Size(min = 30) String> chapterTitles;

```

Show moreShow more icon

更新后的 Bean Validation API 改进了 Java SE 8 的 `Date` 和 `Time` 类型，提供了对 `java.util.Optional` 的支持，如清单 9 所示。

##### 清单 9\. Bean Validation 2.0 支持 Date 和 Time 类型，以及 Optional

```
Optional<@Size(min = 10) String> title;

private @PastOrPresent Year released;
private @FutureOrPresent LocalDate nextVersionRelease;
private @Past LocalDate publishedDate;

```

Show moreShow more icon

容器的级联验证是一个很方便的新特性。使用 `@Valid` 对容器的任何类型参数进行注释，都会导致在验证父对象时验证每个元素。在清单 10 中，将会验证每个 `String` 和 `Book` 元素，如下所示：

##### 清单 10\. 容器类型的级联验证

```
Map<@Valid String, @Valid Book> otherBooksByAuthor;

```

Show moreShow more icon

Bean Validation 还通过插入值提取器来添加对自定义容器类型的支持。内置的约束被标记为 repeatable，参数名是通过反射来检索的， `ConstraintValidator#initialize()` 是默认方法。JavaFX 还获得了对其类型的支持。

### Contexts and Dependency Injection for Java 2.0

Context and Dependency Injection API (CDI) 是一种自第 6 版开始就存在于 Java EE 中的主干技术。从那时起，它就成为了一个简化开发的关键特性。

在最新版本中，该 API 已扩展用于 Java SE。为了适应这一更改，CDI 规范被分解为 3 部分：第 1 部分处理 Java EE 和 Java SE 中的通用概念；第 2 部分处理仅针对 CDI with Java SE 的规则；第 3 部分处理仅针对 Java EE 的规则。

CDI 2.0 还对观察者和事件的行为及交互方式进行了重大更改。

#### CDI 2.0 中的观察者和事件

在 CDI 1.1 中，在触发事件时会同步调用观察者，没有任何机制来定义执行它们的顺序。此行为的问题在于，如果某个观察者抛出异常，则不会再调用所有后续观察者，观察者链将会中断。在 CDI 2.0 中，通过引入 `@Priority` 注释，从某种程度上减轻了这一问题，该注释指定了调用观察者应该采用的顺序，编号越小的观察者越先调用。

清单 11 展示了一次事件触发和两个具有不同优先级的观察者。观察者 `AuditEventReciever1` （优先级为 10）在 `AuditEventReciever2` （优先级为 100）前调用。

##### 清单 11\. 观察者优先级演示

```
@Inject
private Event<AuditEvent> event;

public void send(AuditEvent auditEvent) {
event.fire(auditEvent);
}

// AuditEventReciever1.class
public void receive(@Observes @Priority(10) AuditEvent auditEvent) {
    // react to event
}

// AuditEventReciever2.class
public void receive(@Observes @Priority(100) AuditEvent auditEvent) {
    // react to event
}

```

Show moreShow more icon

值得注意的是，具有相同优先级的观察者会按无法预测的顺序进行调用，默认顺序为 `javax.interceptor.Interceptor.Priority.APPLICATION + 500` 。

#### 异步事件

对观察者特性添加的另一个有用功能是，异步触发事件。添加了一个新触发方法 (`Async()`) 和相应的观察者注释 (`@ObservesAsyncfire`) 来支持此特性。清单 12 展示了一个异步调用的 `AuditEvent` ，以及收到该事件的通知的观察者方法。

##### 清单 12\. 异步触发事件

```
@Inject
private Event<AuditEvent> event;

public CompletionStage<AuditEvent> sendAsync(AuditEvent auditEvent) {
return event.fireAsync(auditEvent);
}

// AuditEventReciever1.class
public void receiveAsync(@ObservesAsync AuditEvent auditEvent) {}

// AuditEventReciever2.class
public void receiveAsync(@ObservesAsync AuditEvent auditEvent) {}

```

Show moreShow more icon

如果任何观察者抛出异常， `CompletionStage` 就会包含 `CompletionException` 。这个实例包含对在观察者调用期间抛出的所有受抑制异常的引用。清单 13 展示了如何管理此场景。

##### 清单 13\. 管理异步观察者中的异常

```
public CompletionStage<AuditEvent> sendAsync(AuditEvent auditEvent) {
System.out.println("Sending async");
CompletionStage<AuditEvent> stage = event.fireAsync(auditEvent)
           .handle((event, ex) -> {
               if (event != null) {
                   return event;
               } else {
                   for (Throwable t : ex.getSuppressed()) {}
                   return auditEvent;
               }
           });

return stage;
}

```

Show moreShow more icon

前面已经提到过，CDI 2.0 还对一些 Java SE 8 特性进行了一般性改进，比如流、lambda 表达式和可重复的修饰符。其他值得注意的增添内容包括：

- 一个新的 `Configurators` 接口。
- 配置或否决观察者方法的能力。
- 内置的注释文字。
- 在 `producer` 上应用 `interceptor` 的能力。

[这里](https://issues.jboss.org/secure/ReleaseNote.jspa?version%3D12325406%26styleName%3D%26projectId%3D12311062%26_sscc%3Dt) 提供了所有更改的完整列表。请参阅规范的 [建议最终草案](https://docs.jboss.org/cdi/spec/2.0-PFD/cdi-spec-with-assertions.html) 了解完整细节。

## 次要更新：JAX-RS 2.1、JSF 2.3、JSON-P 1.1

尽管是相对次要的更新，但对 Java API for RESTful Web Services ( [JSR 370](https://jcp.org/en/jsr/detail?id=370))、JavaServer Faces 2.3 ( [JSR 372](https://jcp.org/en/jsr/detail?id=372)) 和 Java API for JSON Processing 1.1 ( [JSR 374](https://jcp.org/en/jsr/detail?id=374)) 的更改仍值得注意，尤其应该注意反应式和函数式编程的元素。

### Java API for RESTful Web Services 2.1

JAX-RS 2.1 API 版本专注于两个主要特性：一个新的反应式客户端 API，以及对服务器发送的事件的支持。

#### 反应式客户端 API

RESTful Web Services 自 1.1 版开始就包含一个客户端 API，提供了一种访问 Web 资源的高级方式。JAX-RS 2.1 向此 API 添加了反应式编程支持。最明显的区别是 `Invocation.Builder` ，它被用于构造客户端实例。如清单 14 所示，新的 `rx()` 方法具有一个返回类型 `CompletionStage` 和参数化类型 `Response` ：

##### 清单 14\. 包含新 rx() 方法的调用构建器

```
CompletionStage<Response> cs1 = ClientBuilder.newClient()
       .target("http://localhost:8080/jax-rs-2-1/books")
       .request()
       .rx()
       .get();

```

Show moreShow more icon

`CompletionStage` 接口是在 Java 8 中引入的，它提供了一些有用的可能性。在清单 15 中，对不同端点执行了两次调用并组合了结果：

##### 清单 15\. 调用不同端点后的组合结果

```
CompletionStage<Response> cs1 = // from Listing 14
CompletionStage<Response> cs2 = ClientBuilder.newClient()
       .target("http://localhost:8080/jax-rs-2-1/magazines")
       .request()
       .rx()
       .get();

cs1.thenCombine(cs2, (r1, r2) ->
    r1.readEntity(String.class) + r2.readEntity(String.class))
        .thenAccept(System.out::println);

```

Show moreShow more icon

#### 服务器发送的事件

Server Sent Events API (SSE) 由 W3C 在 HTML 5 中引入，并由 [WHATWG 社区](https://html.spec.whatwg.org/multipage/server-sent-events.html) 维护，该 API 允许客户端订阅服务器生成的事件。在 SSE 架构中，创建了一条从服务器到客户端的单向通道，服务器可通过该通道发送多个事件。该连接长期存在并一直处于打开状态，直到一端关闭它。

JAX-RS API 包含一个用于 SSE 的客户端和服务器 API。来自客户端的入口点是 `SseEventSource` 接口，可以通过一个配置好的 `WebTarget` 来修改它，如清单 16 所示。在这段代码中，客户端注册了一个使用者。使用者输出到控制台，然后打开连接。 `onComplete` 和 `onError` 生命周期事件的处理函数也受到支持。

##### 清单 16\. 针对服务器发送的事件的 JAX-RS 客户端 API

```
WebTarget target = ClientBuilder.newClient()
    .target("http://localhost:8080/jax-rs-2-1/sse/");

try (SseEventSource source = SseEventSource
    .target(target).build()) {
       source.register(System.out::println);
       source.open();
}

```

Show moreShow more icon

在另一端，SSE 服务器 API 接受来自客户端的连接，并将事件发送到所有已连接的客户端：

##### 清单 17\. 用于服务器发送的事件的服务器 API

```
@POST
@Path("progress/{report_id}")
@Produces(MediaType.SERVER_SENT_EVENTS)
public void eventStream(@PathParam("report_id")String id,
                        @Context SseEventSink es,
                        @Context Sse sse) {
    executorService.execute(() -> {
    try {
        eventSink.send(
            sse.newEventBuilder().name("report-progress")
                .data(String.class,
                "Commencing process for report " + id)
                .build());
            es.send(sse.newEvent("Progress", "25%"));
            Thread.sleep(500);
            es.send(sse.newEvent("Progress", "50%"));
            Thread.sleep(500);
            es.send(sse.newEvent("Progress", "75%"));
    } catch (InterruptedException e) {
        e.printStackTrace();
    }
    });
}

```

Show moreShow more icon

在清单 17 中，将 `SseEventSink` 和 `Sse` 资源注入到了 resource 方法中。 `Sse` 实例获取一个新的出站事件构建器，并将进度事件发送到客户端。请注意，介质类型为一个单独用于事件流的新 `text/event-stream` 类型。

#### 广播服务器发送的事件

事件可同时广播到多个客户端。此操作是通过在 `SseBroadcaster` 上注册多个 `SseEventSink` 实例来实现的：

##### 清单 18\. 广播到所有已注册的订阅者

```
@Path("/")
@Singleton
public class SseResource {

    @Context
    private Sse sse;

    private SseBroadcaster broadcaster;

    @PostConstruct
    public void initialise() {
       this.broadcaster = sse.newBroadcaster();
    }

    @GET
    @Path("subscribe")
    @Produces(MediaType.SERVER_SENT_EVENTS)
    public void subscribe(@Context SseEventSink eventSink) {
       eventSink.send(sse.newEvent("You are subscribed"));
       broadcaster.register(eventSink);
    }

    @POST
    @Path("broadcast")
    @Consumes(MediaType.APPLICATION_FORM_URLENCODED)
    public void broadcast(@FormParam("message") String message) {
       broadcaster.broadcast(sse.newEvent(message));
    }
}

```

Show moreShow more icon

清单 18 中的示例接收 URI `/subscribe` 上的一个订阅请求。该订阅请求会导致为该订阅者创建一个 `SseEventSink` 实例，然后使用 `SseBroadcaster` 向资源注册该实例。从提交到 URL `/broadcast` 的网页表单中检索广播到所有订阅者的消息，然后通过 `broadcast()` 方法处理它。

#### 对 JAX-RS 2.1 的其他更新

JAX-RS 2.1 还引入了其他一些值得提及的细微更改：

- `Resource` 方法现在支持使用 `CompletionStage<T>` 作为一种返回类型。
- 添加了对 Java API for JSON Processing (JSON-P) 和 Java API for JSON Binding (JSON-B) 的全面支持。
- 支持对所有提供者使用 `@Priority` ，包括实体提供者 — `MessageBodyReader` 和 `MessageBodyWriter` 。
- 对于所有不支持 Java Concurrency Utilities API 的环境，删除了默认设置。

### JavaServer Faces 2.3

JavaServer Faces (JSF) 2.3 版的目的是包含社区请求的特性，以及更紧密地集成 CDI 和 websocket。该版本考虑了大量问题（最终数量达数百个），尝试纠正了许多细微但较为突出的问题。JSF 是一项成熟的技术，所以现在解决的问题都是以前很难解决的或次要的问题。在这些问题中，JSF 2.3 通过对 Java `Date/Time` 类型和类级 bean 验证的支持，对 Java SE 8 进行了改进。要获得完整列表，建议下载 [最终版本规范](http://download.oracle.com/otndocs/jcp/jsf-2_3-final-eval-spec/index.html) 。

#### JSF 2.3 中的服务器推送

JSF 2.3 集成了 Servlet 4.0 中的服务器推送支持。JSF 非常适合在服务器推送实现中识别所需的资源，现在已将它设置为在 `RenderResponsePhase` 生命周期阶段执行该任务。

### JSON Processing 1.1

JSON-P 获得了一个单点版本，这使它能够跟上最新的 [IEFT](https://www.ietf.org/) 标准。这些标准包括 [JSON Pointer](http://tools.ietf.org/html/rfc6901) 、 [JSON Patch](http://tools.ietf.org/html/rfc6902) 和 [JSON Merge Patch](http://tools.ietf.org/html/rfc7396) 。

通过在 `javax.json.streams` 中引入的新 `JsonCollectors` 类，查询 JSON 也得到了简化。

#### JSON-Pointer

_JSON Pointer_ 定义一个字符串表达式来识别 JSON 文档中的特定值。类似于用于识别 XML 中的片段的 [XPointer](https://www.w3.org/TR/WD-xptr) ，JSON Pointer 引用了 JSON 文档中的值。例如，考虑清单 19 中的 JSON 文档， `topics` 数组中的第 2 个元素将通过 JSON 指针表达式 `/topics/1` 进行引用。

##### 清单 19\. 一个包含数组的 JSON 对象

```
{
"topics": [
"Cognitive",
"Cloud",
"Data",
"IoT",
"Java" ]
}

```

Show moreShow more icon

入口 API 是 `JsonPointer` 接口。通过调用 `Json` 类上的静态工厂方法 `createPointer()` ，创建了一个实例。清单 20 中的代码段创建了一个 `JsonPointer` ，并引用了 topics 数组中的第 2 个元素：

##### 清单 20\. JsonPointer 引用一个数组元素

```
Json.createPointer("/topics/1").getValue(jsonTopicData);

```

Show moreShow more icon

`JsonPointer` API 也可以通过 `adding` 、 `replacing` 和 `removing` 属性来修改 JSON 文档。清单 21 向 topics 列表添加了值”Big Data”：

##### 清单 21\. 使用 JsonPointer 向数组添加值

```
Json.createPointer("/topics/0")
    .add(jsonTopicData, Json.createValue("Big Data"));

```

Show moreShow more icon

#### JSON Patch

JSON Patch 表达了应用于 JSON Pointer 符号中的目标 JSON 文档的一系列操作。它可以执行以下操作： `add` 、 `copy` 、 `move` 、 `remove` 、 `replace` 和 `test` 。

`JsonPatchBuilder` 接口是进入此 API 的网关，是利用 `Json` 类上的静态方法 `createPatchBuilder()` 创建的。JSON Pointer 表达式被传递给某个操作方法并应用于 JSON 文档。清单 22 展示了如何将 topics 数组中的第一个元素替换为值”Spring 5”：

##### 清单 22\. 替换了数组中的值的 JsonPatchBuilder

```
Json.createPatchBuilder()
    .replace("/topics/0", "Spring 5")
    .build()
    .apply(jsonTopicData);

```

Show moreShow more icon

多个操作可链接起来并按顺序应用于前一次的修订结果。

#### JSON Merge Patch

JSON Merge Patch 是一个 JSON 文档，描述了将对目标 JSON 文档执行的一组更改。表 1 给出了 3 种可用的操作。

##### 表 1\. 选择合并修订操作

操作目标修订结果Replace{“color”:”blue”}{“color”:”red”}{“color”:”red”}Add{“color”:”blue”}{“color”:”red”}{“color”: null}Remove{“color”:”red”}{“color”:”blue”,”color”:”red”}{}

`Json` 类上的静态方法 `createMergePatch()` 提供了一个 `JsonMergePatch` 类型的实例，您可以向该实例传递修订操作。然后，向得到的 `JsonMergePatch` 实例的 `apply()` 方法传递目标 JSON 并应用该修订。清单 23 展示了如何执行表 1 中的 `replace` 操作。

##### 清单 23\. 使用 JSON Merge Patch 来替换值

```
Json.createMergePatch(Json.createValue("{\"colour\":\"red\"}"))
    .apply(Json.createValue("{\"colour\":\"red\"}"));

```

Show moreShow more icon

#### JsonCollectors

由于 Java 8 向 `javax.json.streams` 包引入 `JsonCollectors` 类，查询 JSON 变得简单得多。清单 24 展示了如何按字母 _c_ 过滤 `topics` 数组并将结果收集到一个 `JsonArray` 中。

##### 清单 24\. 使用 JsonCollectors 过滤数组并将结果收集到 JsonArray 中

```
JsonArray topics = jsonObject.getJsonArray("topics")
    .stream()
    .filter(jv -> ((JsonString) jv).getString().startsWith("C"))
    .collect(JsonCollectors.toJsonArray());

```

Show moreShow more icon

## 结束语

Java EE 针对云进行了重新定位，计划的版本包含两部分，前半部分包含促进实现该目标的技术。Java EE 8 一面世，Java EE 9 的研发工作就会开始。目前的目标是在一年内发布 Java EE 9。

Java EE 9 开发路线图已包含对 Java EE Security API 的增强，这些增强将包含不适用于 Java EE Security 1.0 中的一些特性。开发人员也有望看到对云中企业开发的深化支持，看到更多 Java EE 技术将使用 [反应式宣言](http://www.reactivemanifesto.org/) 作为蓝图来提高灾备能力和可伸缩性。

我对 Java EE 9 的愿望清单包括针对微服务友好的技术的改进。两个将支持微服务的 API 是 Configuration API 1.0 ( [JSR 382](https://www.jcp.org/en/jsr/detail?id=382)) 和 Health Check，二者都已添加到 Java EE 8 中，但正在考虑进行下一次更新。

最新 [公告](https://blogs.oracle.com/theaquarium/opening-up-java-ee) 表明 Eclipse 正在采用 Java EE，我也很高兴听到此消息。预计开源基础将促进来自供应商和开发人员的更广泛的社区参与。我们应当预料到 Java EE 的发布节奏将会加快，这个强大的企业规范将会得到持续增强。

本文翻译自： [What’s new in Java EE 8](https://developer.ibm.com/articles/j-whats-new-in-javaee-8/)（2017-09-26）