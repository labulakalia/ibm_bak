# MicroProfile 1.3
了解 Eclipse MicroProfile 1.3 中的 5 个新 API

**标签:** Eclipse MicroProfile,Jakarta EE,Java

[原文链接](https://developer.ibm.com/zh/articles/j-5things18/)

Alex Theedom

发布: 2018-05-21

* * *

**关于本系列**

您觉得自己了解 Java 编程？事实是，大多数开发人员只是了解 Java 平台的皮毛，所学知识也仅够应付工作。在这个 [连载系列](/zh/series/5-things-you-didnt-know-about/) 中，Java 技术侦探们将深度挖掘 Java 平台的核心功能，揭示一些可帮助您解决最棘手的编程挑战的技巧和诀窍。

Eclipse MicroProfile 社区成立于 2016 年，致力于为 Java™ 中的微服务架构快速开发新工具和规范。每个新版本都会扩充一些最适合开发 Java 云原生微服务的技术集合。MicroProfile 现在为第 3 版，即 V1.3。

本文将引导您了解 MicroProfile 1.3 的 5 个新 API 的亮点：

- REST Client 1.0
- Metrics 1.1
- OpenAPI 1.0
- OpenTracing 1.0
- Config 1.2

所有 MicroProfile API 都是基于以下 Java EE 7 API 而构建的：

- 所有 API 都利用了 CDI 1.2，后者是该平台的核心。
- JAX-RS 2.1 为 RESTful API 的创建提供了客户端和服务器特性。
- JSON-P 1.0 提供了 JSON 处理功能。
- Commons Annotation 1.2 为常用语义概念提供了注解。

请参阅 [MicroProfile 项目页](https://microprofile.io) ，了解 V1.2 中发布的 Java MicroProfile API：Fault Tolerance 1.0、Health Check 1.0 和 JWT Authentication 1.0。

## 安装 MicroProfile 1.3

MicroProfile 1.3 至少需要 Java SE 8 才能运行。清单 1 给出了要添加到 _pom.xml_ 中的 MicroProfile 1.3 的依赖项。

##### 清单 1\. MicroProfile 1.3 的 Maven 坐标

```
<dependency>
<groupId>org.eclipse.microprofile</groupId>
<artifactId>microprofile</artifactId>
<version>1.3</version>
<type>pom</type>
<scope>provided</scope>
</dependency>

```

Show moreShow more icon

也可以通过 [IBM Cloud Developer](https://www.ibm.com/cloud/cli) 工具或 [IBM Cloud UI](https://cloud.ibm.com/developer/appservice/starter-kits?cm_sp=ibmdev-_-developer-articles-_-cloudreg) 为 IBM Liberty 生成一个 MicroProfile 项目。

[获得代码](https://github.com/eclipse/microprofile-bom)

## MicroProfile 项目的历史

MicroProfile 的诞生是为了应对 2016 年 Oracle 在 Java EE 开发上的止步不前。在那时，Java EE 缓慢的版本发布节奏不利于集成对微服务和云原生解决方案的支持。一些供应商和用户组提出了用 Java MicroProfile 来解决这一需求，他们的提案很快获得了大力支持。

自那以后，其他许多供应商加入了 MicroProfile 社区，MicroProfile 也迁移到了 Eclipse Foundation，在 Eclipse Foundation 中作为一个 [Eclipse 项目](https://projects.eclipse.org/projects/technology.microprofile) 进行维护。

Java MicroProfile 规范是一项社区成果：欢迎新贡献者参与其中。您可以先参加 [MicroProfile Community Google Group](https://groups.google.com/forum/#!forum/microprofile) 中的讨论。也欢迎开发人员向 [托管在 GitHub 上的 MicroProfile 项目](https://github.com/eclipse/microprofile) 贡献代码。

接下来的几节将揭示 MicroProfile 1.3 的 5 个方面：即此版本中包含的 5 个新 API。

## 1.Rest Client 1.0

Rest Client API 基于 JAR-RS 2.0，提供了一种类型安全方式来通过 HTTP 调用 RESTful 端点。清单 2 中的 `BookClientService` 接口表示远程接口，它能够创建图书并基于 ISBN 检索图书。

##### 清单 2\. 创建 BookClientService 接口

```
@Path("/books")
public interface BookClientService {
    @GET
    @Path("/{isbn}")
    Book getBook(@PathParam("isbn") String isbn);

    @POST
    String addBook(Book book);
}

```

Show moreShow more icon

在清单 3 中， `BookClientService` 类被传递给 `RestClientBuilder` 的一个构建器方法。然后调用接口方法来创建和检索一本书。

##### 清单 3\. 使用 BookClientService 接口

```
URI bookSrvUri = new URI("http://localhost:8081/bookservice");
BookClientService bookSrv = RestClientBuilder.newBuilder()
            .baseUri(bookSrvUri)
            .build(BookClientService.class);

Book book = new Book ("Fun with Microprofile","Alex Theedom");
String isbn = bookSrv.addBook(book);
Book newBook = bookSrv.getBook(isbn);

```

Show moreShow more icon

## 2.Metrics 1.1

Metrics API 为开发人员提供了以便利方式检测微服务的一种方法。该 API 的设计类似于 Dropwizard Metrics API，许多开发人员对后者已经非常熟悉。使用 Metrics API 会消除对第三方库的依赖。

为了更易于使用，熟悉流行的 [Prometheus](https://prometheus.io/) 格式的开发人员可以简单设置 Prometheus 来获取这些指标并监控他们的服务。

Metrics API 与 Health Check API 1.2（在 MicroProfile 1.2 中发布）的不同之处在于，它并不只是提供了服务健康状况的简单的”是/否”二进制表示。它还提供了关于服务性能的详细元数据。

### 在云和微服务应用程序中启用指标

要在 MicroProfile 应用程序中启用指标，您需要创建指标并向应用程序注册表注册它们。可通过以下方式之一完成此任务：

- 使用 MetricRegistry。此方法要求显式创建和注册您的指标。
- 注入指标并使用指标注解。指标由 CDI 容器自动实例化，并向 MetricRegistry 应用程序进行注册。

### 查看指标值

指标通过调用服务中的 URI /metrics 端点来呈现，如清单 4 所示。该指标端点位于 `localhost` 上，在端口 9091 上运行。

##### 清单 4\. 指标端点

```
https://localhost:9091/metrics

```

Show moreShow more icon

### 通过 MetricRegistry 启用指标

MetricRegistry 应用程序的用途是存储所有指标和它们的元数据。指标和元数据可使用上面列出的任何方式来注册和检索。让我们更仔细地查看每一种方法。

#### 手动创建和注册指标

将一个 `MetricRegistry` 实例注入您的类中，如清单 5 所示，并调用各种方法来注册和检索这些指标。

##### 清单 5\. 注入 MetricRegistry

```
@Inject private MetricRegistry registry;

```

Show moreShow more icon

`MetricRegistry` 类中包含检索和注册指标的方法，如清单 6 所示。这个代码片段展示了如何创建和注册一个计数类型的指标。5 种指标类型包括计数器、计量器、仪表、直方图和计时器。

##### 清单 6\. 创建一个新的元数据计数器指标

```
Metadata counter = new Metadata(
    "hitsCounter",                            // name
    "Hits Count",                             // display name
    "Number of hits",                         // description
    MetricType.COUNTER,                       // type
    MetricUnits.NONE);                        // units

```

Show moreShow more icon

通过 MetricRegistry 以编程方式注册指标，如清单 7 所示。

##### 清单 7\. 注册计数器指标

```
Counter hitsCounter = registry.counter(counter);

```

Show moreShow more icon

注册指标后，可以使用它来记录指标数据，如清单 8 所示。

##### 清单 8\. 递增 hitsCounter 指标

```
@GET
@Path("/messages")
public String getMessage() {
    hitsCounter.inc();
    return "Hello world";
}

```

Show moreShow more icon

指标在记录后存储在 MetricRegistry 中，并在调用 `/metric` URI 时检索。要检索上面创建的点击计数器指标，可以调用以下 URI： `https://localhost:9091/metrics/application/` 。

清单 9 中的输出是从 REST 端点生成的结果：

##### 清单 9\. REST 端点生成的指标统计数据

```
https://localhost:9091/metrics/application/hitsCounter
# TYPE application:hits_count counter
# HELP application:hits_count Number of hits
application:stats_hits 499

```

Show moreShow more icon

这些指标表示为 Prometheus 格式。如果您想要 JSON 格式的统计数据，可以调用同一个端点，但将 `HTTP Accept` 标头配置为 `application/json` ，如清单 10 所示。

##### 清单 10\. 获取 JSON 格式的统计数据

```
https://localhost:9091/metrics/application/statsHits
{"hitsCount":499}

```

Show moreShow more icon

#### 注入指标并使用注解

CDI 注入是在应用程序中启用指标的更简单、更自然的方式。您仍会创建指标并向 MetricRegistry 注册指标，但不会创建元数据实例，而是将指标配置数据作为 `@Metric` 注解的参数进行传递，如清单 11 所示。

##### 清单 11\. 通过 @Metric 配置指标

```
https://localhost:9091/metrics/application/statsHits
{"hitsCount":499}

```

Show moreShow more icon

`@Metric` 注解促使服务器在 MetricRegistry 中创建和注册一个计数器，并将它提供给应用程序使用。

## 3.OpenAPI 1.0

作为 MicroProfile 1.3 中的新特性，OpenAPI 1.0 为 RESTful 服务提供了一个契约，类似于 SOAP Web 服务的 WSDL 契约。OpenAPI 1.0 基于 [OpenAPI 规范](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.0.md) ，提供了一组 Java 接口和编程模型，以便允许开发人员使用 JAX-RS 应用程序生成 OpenAPI v3 文档。

OpenAPI 使用 JAX-RS 2.0 应用程序生成一个有效的 OpenAPI 文档。它处理所有 JAX-RS 注解（包括 `@Path` 和 `@Consumes/@Produces` 注解），以及用作 JAX-RS 操作的输入或输出的 POJO。

### 启动并配置 OpenAPI

要开始使用 OpenAPI 1.0，只需将您的现有 JAX-RS 应用程序部署到诸如 Open Liberty 之类的 MicroProfile 服务器中，然后从 OpenAPI URI `/openapi` 签出输出。

看到调用 `/openapi` 的输出后，您可能想要配置您的 RESTful 资源来提供更好的文档。这是通过对 RESTful API 的各种元素进行注解来完成的。

表 1 列出了 OpenAPI 1.0 规范的关键注解。所有 OpenAPI 注解都已在 [org.eclipse.microprofile.openapi.annotations](https://github.com/eclipse/microprofile-open-api/tree/master/api/src/main/java/org/eclipse/microprofile/openapi/annotations) 包中列出。

##### 表 1\. 关键 OpenAPI 注解

注解描述@Operation描述针对某个特定路径的操作，或者通常是一个 HTTP 方法。@APIResponse描述来自一个 API 操作的一次响应。@RequestBody描述一个请求主体。@Content提供一种特定媒体类型的模式和示例。@Schema允许定义输入和输出数据类型。@Server一个用于多个服务器定义的容器。@ServerVariable表示一个用于服务器 URL 模板置换的服务器变量。@OpenAPIDefinition一个 OpenAPI 定义的综合元数据

### 使用 OpenAPI 注解

OpenAPI 注解将对 JAX-RS 资源定义中的特征进行标记，并提供该资源的信息。在清单 12 中， `findBookByISBN` 方法使用 `@Operation` 进行了注解，并传递了关于其功能的信息。

##### 清单 12\. 使用 OpenAPI 注解标记 REST 资源

```
@GET @Path("/findByISBN") @Operation( summary = "Finds a book by its ISBN", description = "Only one ISBN should be provided") public Response findBookByISBN(...){ ...}

```

Show moreShow more icon

来自此注解的输出如清单 13 所示。

##### 清单 13\. REST 文档的 YAML 输出

```
/books/findBookByISBN:
get:
    summary: Finds a book by its ISBN
    description: Only one ISBN should be provided
    operationId: findBookByISBN

```

Show moreShow more icon

调用 OpenAPI URI `/openapi` 的默认输出为 YAML。如果调用此端点时将 `HTTP Accept` 标头设置为 `application/json` ，那么响应将为 JSON 格式。

## 4.OpenTracing 1.0

OpenTracing API 是 [OpenTracing Project](http://opentracing.io/) 的一种实现。它允许您跟踪一个跨服务边界的请求流。这在微服务环境中特别重要，其中的请求通常会经历许多服务。

为了完成分布式跟踪，每个服务都将一个关联 ID 记录到一个指定用于存储分布式跟踪记录的服务中。然后，各个服务提供与某个请求有关联的跟踪信息的视图。

OpenTracing 在运行时可以包含或不包含对应用程序代码的显式检测。我们将查看这两种选项。

### 不含代码检测的 OpenTracing

此操作模式允许 JAX-RS 应用程序参与跟踪记录，而不需要向应用程序添加任何分布式跟踪代码。在此用法中，开发人员甚至不需要知道任何关于应用程序将部署到的跟踪环境的信息。

### 包含代码检测的 OpenTracing

在某些情况下，您可能希望主动与跟踪信息交互。在这种情况下，您将使用 `@Traced` 注解来指定要跟踪的类或方法。可以结合使用显式代码检测和无代码的检测。在某些场景中，集成两种方法会创建一种适合您的用例的解决方案。

### 使用 @Traced 注解

将 `@Traced` 注解应用于某个类时，它会应用于该类中的所有方法。如果您想仅跟踪部分方法，可以仅将该注解添加到这些方法中。该注解在方法的开头处启用了一个 `Span` ，并在方法的末尾处结束了该 `Span` 。在清单 14 中，在类级注解中将跟踪的操作命名为 `Book Operation` 。从跟踪中排除方法 `getBookByIsbn()` ，因为它将 false 指定为包含值。

##### 清单 14\. @Traced 用法示例

```
@Path("/books")
@Traced(operationName = "Book Operation")
public class BookResource {

    @GET
    public Response getBooks() {
      // implementation code removed for brevity
    }

    @Traced(value = false)
    @GET
    @Path("/{isbm}")
    public Response getBookByIsbn(String @Param("isbn") isbn {
      // implementation code removed for brevity
    }

}

```

Show moreShow more icon

## 5.Config 1.2

为了避免在操作环境改变时重新打包应用程序，人们开发了许多方法。这对微服务特别重要，因为它们是专为在多个环境中可移植和可运行而设计的。

Config 1.2.1 API 聚合了来自各种不同 `ConfigSource` 的配置，并向用户呈现单一视图。

配置被存储为 `String/String` 键值格式。配置键可以使用一种点分隔表示法，类似于 Java 包中的名称空间。清单 15 给出了一个兼容 Config API 的配置文件的示例。

##### 清单 15\. 示例配置文件

```
com.readlearncode.bookservice.url = http://readlearncode.com/api/books
com.readlearncode.bookservice.port = 9091

```

Show moreShow more icon

### 使用配置对象

可以通过注入，使用 `@Inject Config` 来自动获取配置对象的实例，如清单 16 所示：

##### 清单 16\. 对配置使用注入

```
public class AppConfiguration {
    public void getConfiguration() {
        Config config = ConfigProvider.getConfig();
        String bookServiceUrl = config.getValue("com.readlearncode.bookservice.url", String.class);
    }
}

```

Show moreShow more icon

也可以通过 `ConfigProvider` 以编程方式获取它，如清单 17 所示：

##### 清单 17\. 对配置使用 CDI

```
@ApplicationScoped
public class InjectedConfigUsageSample {

    @Inject
    private Config config;

    @Inject
    @ConfigProperty(name="com.readlearncode.bookservice.url")
    private String bookServiceUrl ;

}

```

Show moreShow more icon

## 结束语

MicroProfile 为 Java 开发人员提供了大量的特性，这些特性变得对开发基于微服务的系统不可或缺。较短的版本发布周期可以确保定期添加新特性，我们已经看到，该规范自创建以来在不断快速成长。每个 MicroProfile API 都遵循业界认可的最佳实践，以帮助确保基于 Java 的微服务既灵活又稳健。

本文翻译自： [MicroProfile 1.3](https://developer.ibm.com/articles/j-5things18/)（2018-04-30）