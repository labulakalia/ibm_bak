# 使用 Quarkus 开发 Kubernetes 原生 Java 应用
面向云平台的 Java 应用开发

**标签:** Docker,Java,Kubernetes,云计算,容器,微服务

[原文链接](https://developer.ibm.com/zh/articles/j-use-quarkus-to-develop-kubernetes-native-java-application/)

成 富

发布: 2019-12-09

* * *

随着 Docker 和 Kubernetes 的流行，容器化成为很多应用的部署选择。Kubernetes 也成为流行的应用部署平台。其实，容器化的思想和微服务架构可以很好的结合在一起。从实现上来说，微服务架构把应用垂直切分成多个相互协同的单元。从应用部署的角度来说，把应用的每个微服务用容器的方式部署在 Kubernetes 平台，可以充分利用 Kubernetes 平台提供的功能。通过 Kubernetes 平台提供的服务发现、自动伸缩、自动化容器部署和监控等功能，可以实现易维护和可伸缩的微服务架构。对 Java 应用来说，进行容器化并不是一件复杂的事情，只需要在 Docker 镜像中添加 JDK 来运行 Java 应用即可。每个容器都需要自己的 JDK 运行时支持。这样的部署方式对一般的 Java 应用没有问题。但是，当微服务数量和容器数量增加时，JDK 所带来的成本代价变得越来越高，甚至超过应用本身，这就造成容器化 Java 应用启动速度慢和占用内存资源过多等问题。虽然 Java 9 引入的模块系统允许对 JDK 自身的模块进行定制，只保留应用需要的 JDK 模块，但这只在一定程度上缓解了这个问题。Quarkus 的出现改变了这一现状，它是一个面向容器的 Java 应用开发框架，能够解决容器化 Java 应用的启动速度和内存占用问题。

## Quarkus 简介

Quarkus 是一个 Java 应用开发框架。与传统开发框架的不同，Quarkus 的目标是创建在容器中可以快速启动和占用更少资源的 Java 应用，其设计时的基本理念是容器优先。Quarkus 提供了对 GraalVM 及其 Substrate 虚拟机的良好支持。通过 GraalVM 的原生镜像功能，可以把 Quarkus 应用打包成体积小且启动速度快的原生镜像。一个使用 Quarkus 创建的简单 REST API 应用，它所占用的内存只有 13MB，启动时间只需要 14 毫秒。相对于传统的 Java 应用开发栈来说，这是一个极大的性能提升。

此外，Quarkus 会在应用构建时进行大量的分析，使得构建的应用只包含运行时所需的 Java 类。这可以进一步降低 Quarkus 应用运行时的内存占用。同时，Quarkus 减少了对 Java 反射 API 的使用。另外，在构建原生镜像时，Quarkus 框架会预先启动自身，从而降低启动时间。

Quarkus 使用了大量开源框架和库，包括 Hibernate、Netty、RESTEasy、Eclipse MicroProfile、Eclipse Vert.x 和 Apache Camel 等。Quarkus 自身的开发也得到了 Red Hat 的支持，因此它的版本更新和维护有很好的保障。本文使用的是 Quarkus 0.27.0 版本。

## Quarkus 开发快速入门

下面先介绍如何快速创建 Quarkus 应用、构建原生镜像并部署到 Kubernetes 平台上。

### 创建 Quarkus 应用

Quarkus 提供了两种不同的方式来创建 Quarkus 应用。

1. 使用 [code.quarkus.io](https://code.quarkus.io/) 网站生成 Quarkus 应用。在这个网站上可以选择 Quarkus 应用所需要使用的扩展（extension）。
2. 使用 Quarkus 提供的 Maven 插件以命令行的方式创建新应用。代码清单 1 给出了使用 Maven 插件生成 Quarkus 应用的示例命令。

##### 清单 1\. 使用 Maven 插件生成 Quarkus 应用

```
mvn io.quarkus:quarkus-maven-plugin:0.27.0:create \
    -DprojectGroupId=io.vividcode \
    -DprojectArtifactId=quarkus-starter \
    -DclassName="io.vividcode.quarkus.ExampleResource" \
    -Dpath="/example"

```

Show moreShow more icon

当上述命令执行完成之后，`quarkus-starter`目录中包含的是新创建的 Quarkus 应用。在该目录下执行 `mvn compile quarkus:dev` 可以启动在端口 8080 运行的开发服务器。在浏览器中访问 `http://localhost:8080` 可以看到 Quarkus 的默认页面。生成的代码中包含了一个名为 `ExampleResource` 的 JAX-RS 资源，如代码清单 2 所示。当用浏览器访问路径 `/example` 时，会看到字符串 `hello`。

##### 清单 2\. Quarkus 应用默认生成的 JAX-RS 资源

```
@Path("/example")
public class ExampleResource {

@GET
@Produces(MediaType.TEXT_PLAIN)
public String hello() {
    return "hello";
}
}

```

Show moreShow more icon

当使用 Maven 命令 `quarkus:dev` 启动 Quarkus 开发服务器时，Quarkus 支持应用的热重载（hot reload）。在 IDE 中修改代码之后，只需要刷新浏览器，Quarkus 会自动重新加载应用代码，可以即时查看修改结果。Java 远程调试服务在端口 5005 启动，可以使用 IDE 的远程调试功能进行调试。

### 构建原生镜像

构建 Quarkus 应用的原生镜像需要 GraalVM 的支持。GraalVM 的详细介绍，可以参考“ [使用 GraalVM 开发多语言应用](https://www.ibm.com/developerworks/cn/java/j-use-graalvm-to-run-polyglot-apps/index.html)”一文。根据 GraalVM 官方网站上的文档说明，下载安装 GraalVM 的社区版或企业版。安装完成之后，需要添加环境变量 `GRAALVM_HOME` 指向 GraalVM 的安装目录，如下面的代码所示。

```
export GRAALVM_HOME=<somedir>/graalvm-ce-19.2.1/Contents/Home

```

Show moreShow more icon

在 Quarkus 应用的当前目录下，使用 `./mvnw package -Pnative` 命令生成原生镜像。当该命令执行结束后，会在 `target` 目录产生可执行文件 `quarkus-starter-1.0-SNAPSHOT-runner`。直接运行该文件，可以启动相应的 REST 服务。

注意，使用上述命令生成的可执行文件与当前开发环境相关，只能在当前环境上运行。当需要在 Kubernetes 上运行时，需要生成适合于 Linux 环境的可执行文件。这是通过命令 `./mvnw package -Pnative -Dnative-image.docker-build=true` 来完成的。生成的 Quarkus 应用的 `src/main/docker` 目录中已经包含了使用原生镜像的 Docker 文件 `Dockerfile.native`。只需要通过 `docker build` 命令来创建 Docker 镜像即可，如下面的代码所示：

```
docker build -f src/main/docker/Dockerfile.native -t vividcode/quarkus-starter .

```

Show moreShow more icon

当 Docker 镜像构建完成之后，可以使用 Docker 命令来运行，如下面的代码所示：

```
docker run -i --rm -p 8080:8080 vividcode/quarkus-starter

```

Show moreShow more icon

#### 部署到 Kubernetes

在创建了 Quarkus 应用的 Docker 镜像之后，就可以部署到 Kubernetes 或其他云平台。具体的部署方式与其他应用并没有什么不同，包括把 Docker 镜像发布到 Docker 仓库中，创建 Kubernetes 中的部署和服务等。具体的步骤可以参考 Kubernetes 相关文档。

##### 免费试用 IBM Cloud

利用 [IBM Cloud Lite](https://cloud.ibm.com/registration?cm_sp=ibmdev-_-developer-articles-_-cloudreg) 快速轻松地构建您的下一个应用程序。您的免费帐户从不过期，而且您会获得 256 MB 的 Cloud Foundry 运行时内存和包含 Kubernetes 集群的 2 GB 存储空间。 [了解所有细节](https://www.ibm.com/cloud/free/?cm_sp=ibmdev-_-developer-articles-_-cloudreg) 并确定如何开始。

下面介绍 Quarkus 的核心功能及重要扩展。

## 依赖注入

在 Quarkus 应用开发中同样可以使用类似 Spring 框架所提供的依赖注入功能。Quarkus 的依赖注入实现基于 JSR 365 定义的 CDI (Contexts and Dependency Injection for Java 2.0) 规范，并且只实现了 CDI 的部分功能。这些功能对于开发应用已经足够。

代码清单 3 中的 `UserService` 类是管理 User 对象的服务层实现。UserService 类上的注解 `@ApplicationScoped` 声明了该对象是一个应用作用域中的 bean。Quarkus 支持的其他作用域相关的注解包括`@Dependent`、`@Singleton`、`@RequestScoped` 和 `@SessionScoped`。

##### 清单 3\. 示例服务 UserService

```
@ApplicationScoped
public class UserService {
private Map<String, User> users = new HashMap<>();

public UserService() {
    addUser(new User("test1", "test1@example.com"));
    addUser(new User("test2", "test2@example.com"));
    addUser(new User("test3", "test3@example.com"));
}

public void addUser(User user) {
    if (user != null) {
      users.put(user.getId(), user);
    }
}

public User deleteUser(String userId) {
    return users.remove(userId);
}

public List<User> list() {
    return new ArrayList<>(users.values());
}
}

```

Show moreShow more icon

代码清单 4 中的 `UserResource` 类通过 `@Inject` 注解来声明它所使用的 `UserService` 对象通过依赖注入的方式提供。使用依赖注入的字段的可访问性一般设置为仅包可见（package private），这样可以避免使用反射 API 设置字段值时产生问题。

##### 清单 4\. 通过依赖注入使用 UserService

```
public class UserResource {
@Inject
UserService userService;
}

```

Show moreShow more icon

## 创建 REST 服务

通过 [代码清单 1](#清单-1-使用-maven-插件生-quarkus-应用) 中的 Maven 命令生成的 Quarkus 应用已经提供了作为示例的 REST 服务，也就是 [代码清单 2](#清单-2-quarkus-应用默认生成的-jax-rs-资源) 中的示例 JAX-RS 资源。不过该服务使用纯文本作为内容格式。在实际的 REST 服务中，JSON 是最常用的内容格式。在 Quarkus 应用中创建使用 JSON 的 REST 服务，需要添加 `resteasy-jsonb` 或 `resteasy-jackson` 扩展。这两个扩展的区别在于，`resteasy-jsonb` 扩展使用 JSON-B，而 `resteasy-jackson` 使用 Jackson。对于一个已有的 Quarkus 应用，可以使用 Quarkus 的 Maven 插件中的 `add-extension` 命令来添加扩展，如下面的代码所示：

```
./mvnw quarkus:add-extension -Dextensions="resteasy-jsonb, undertow"

```

Show moreShow more icon

如代码清单 5 中所示，UserResource 类的 `list()` 方法的返回值为 `List<User>` 对象。由于通过 `@Produces(MediaType.APPLICATION_JSON)` 注解声明了 REST 服务产生的内容类型为 JSON，`List<User>` 对象会被自动序列化为 JSON 格式。当用浏览器访问路径 `/user` 时，可以看到 JSON 格式的内容。

##### 清单 5\. 使用 JSON 格式的 REST 服务

```
@Path("/user")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
public class UserResource {
@Inject
UserService userService;

@GET
public List<User> list() {
    return userService.list();
}
}

```

Show moreShow more icon

JSON 序列化库使用 Java 反射 API 来查找对象中的属性。当 Quarkus 应用以原生镜像的方式运行在 GraalVM 上时，由于 Substrate 虚拟机的限制，通过反射 API 访问的类都需要事先声明。`UserResource` 类中对 User 类的使用出现在 `list()` 方法的返回值中，Quarkus 可以自动识别出 User 类并进行声明。如果方法的返回值类型是 Response 对象，则 Quarkus 无法识别出其中包含的实体 Java 类。这种情况下需要使用 `@RegisterForReflection` 注解声明 `Response` 对象中包含的实体类。

另外一个与 REST 服务相关的功能是对 OpenAPI 的支持。通过 Quarkus 的 `openapi` 扩展，可以生成 OpenAPI 的规范文档。在添加了 `openapi` 扩展之后，访问路径 `/openapi` 可以得到 REST 服务的基于 OpenAPI v3 规范的文档。该扩展也自带了 Swagger 界面，可以通过路径 `/swagger-ui` 来访问。

## 访问关系数据库

Quarkus 应用可以使用 Hibernate 访问关系数据库。首先需要添加 Hibernate 对应的 `hibernate-orm` 扩展以及相关的数据库驱动。以 MySQL 为例，需要添加对应的 `jdbc-mysql` 扩展。代码清单 6 展示了与 Hibernate 和 MySQL 相关的 `application.properties` 文件中的配置内容。

##### 清单 6\. Hibernate 和 MySQL 相关的配置内容

```
quarkus.datasource.url = jdbc:mysql://localhost:3306/quarkus_starter
quarkus.datasource.driver = com.mysql.cj.jdbc.Driver
quarkus.datasource.username = quarkus
quarkus.datasource.password = quarkus
quarkus.hibernate-orm.dialect = org.hibernate.dialect.MySQL8Dialect
quarkus.hibernate-orm.dialect.storage-engine = InnoDB
quarkus.hibernate-orm.database.generation = drop-and-create

```

Show moreShow more icon

### 使用 JPA 和 Hibernate

完成对 Hibernate 和数据库的配置之后，接着需要对实例类添加 JPA 相关的注解。在代码清单 7 中，`@Entity` 注解被添加到 `User` 类中。`User` 类中包含 3 个字段，其中 `id` 是数据库表的主键，使用 `@Id` 注解进行声明。

##### 清单 7\. JPA 中的实体类

```
@Entity
public class User {
private String id;
private String name;
private String email;

public User() {
    this.id = UUID.randomUUID().toString();
}

public User(String name, String email) {
    this();
    this.name = name;
    this.email = email;
}

@Id
public String getId() {
    return id;
}

public void setId(String id) {
    this.id = id;
}

public String getName() {
    return name;
}

public void setName(String name) {
    this.name = name;
}

public String getEmail() {
    return email;
}

public void setEmail(String email) {
    this.email = email;
}
}

```

Show moreShow more icon

代码清单 8 中的 UserService 类使用 JPA 中的 `EntityManager` 对象进行数据库操作。`EntityManager` 对象通过依赖注入的方式来获取。相关的数据库操作与一般的 JPA 应用没有区别。

##### 清单 8\. 使用 JPA 的 UserService 类

```
@ApplicationScoped
public class UserService {

@Inject
EntityManager entityManager;

@Transactional
public User addUser(User user) {
    if (user != null) {
      entityManager.persist(user);
    }
    return user;
}

public List<User> list() {
    CriteriaQuery<User> query = entityManager.getCriteriaBuilder().createQuery(User.class);
    TypedQuery<User> allQuery = entityManager.createQuery(query.select(query.from(User.class)));
    return allQuery.getResultList();
}

@Transactional
public void delete(String id) {
    CriteriaBuilder criteriaBuilder = entityManager.getCriteriaBuilder();
    CriteriaDelete<User> deleteQuery = criteriaBuilder.createCriteriaDelete(User.class);
    Root<User> root = deleteQuery.from(User.class);
    deleteQuery.where(criteriaBuilder.equal(root.get("id"), id));
    entityManager.createQuery(deleteQuery).executeUpdate();
}
}

```

Show moreShow more icon

### 使用 Panache

使用 Hibernate 和 JPA 进行数据库访问的代码不够直观和简洁，我们可以使用 Panache 来简化对 Hibernate 的使用。使用 Panache 之前需要添加 `hibernate-orm-panache` 扩展。代码清单 9 中的 `Product` 类是另外一个实体类。`Product` 与 `User` 的不同之处在于，`Product` 类继承自 `PanacheEntity` 类。`PanacheEntity` 类提供了很多实用方法来简化 JPA 相关的操作。Product 类的静态方法 `findByName()` 使用 `PanacheEntity` 类的父类 `PanacheEntityBase` 中的 `find()` 方法来根据 `name` 字段查询并返回第一个结果；`priceGte()` 方法使用 `find()` 方法来查询价格大于或等于给定值的 `Product` 对象。相对于使用 JPA 中的 `EntityManager` 和 `CriteriaBuilder`，`PanacheEntity` 类提供的实用方法要简单很多。

##### 清单 9\. 使用 Panache 的实体类

```
@Entity
public class Product extends PanacheEntity {

public String name;
public BigDecimal price;

public static Product findByName(String name) {
    return find("name", name).firstResult();
}

public static List<Product> priceGte(BigDecimal price) {
    return find("price >= ?1", price).list();
}
}

```

Show moreShow more icon

代码清单 10 中的 `ProductService` 类可以直接使用 `Product` 类中的静态方法完成相关查询。

##### 清单 10\. Product 相关的 ProductService 类

```
@ApplicationScoped
public class ProductService {
@Transactional
public Product addProduct(Product product) {
    product.persist();
    return product;
}

public List<Product> list() {
    return Product.listAll();
}

@Transactional
public void delete(Long id) {
    Product.delete("id", id);
}

public List<Product> findPriceGte(BigDecimal price) {
    return Product.priceGte(price);
}
}

```

Show moreShow more icon

## 异步消息传递

Quarkus 中的异步消息传递有两种方式：第一种是使用 AMQP，第二种是使用 Eclipse Vert.x 中的事件总线（Event Bus）。

### 使用 AMQP

Quarkus 的 AMQP 支持使用 SmallRye 的反应式消息库。使用 AMQP 时需要首先添加 `amqp` 扩展。本节中的示例使用 Apache Artemis 作为 AMQP 协商器（broker）。

SmallRye 的编程模型是抽象的异步数据流。代码清单 11 中的 `NumberGenerator` 类的 `generate()` 方法是数据流中消息的生产者。该方法每隔 10 秒钟会产生一个随机的 Long 类型的整数。方法的返回值是 RxJava 2 中的 `Flowable` 类型。`@Outgoing("generated-number")` 注解表示 `generate()` 方法产生的消息会被发送到名为 generated-number 的流中。

##### 清单 11\. 数据流的生产者

```
@ApplicationScoped
public class NumberGenerator {

@Outgoing("generated-number")
public Flowable<Long> generate() {
    return Flowable.interval(10, TimeUnit.SECONDS)
        .map(tick -> ThreadLocalRandom.current().nextLong());
}

}

```

Show moreShow more icon

代码清单 12 中的 `PowCalculator` 类的 `calculate()` 方法对于输入的数字，返回其 10 次方的 `BigInteger` 对象。该方法的 `@Incoming("numbers")` 注解表示接收来自数据流 numbers 的值作为输入参数，`@Outgoing("pow10")` 注解表示该方法的返回值被发送到数据流 `pow10` 中。`@Broadcast`注解表示产生的消息会被发送到所有匹配的接收者。

##### 清单 12\. 数据流的消费者

```
@ApplicationScoped
public class PowCalculator {
@Incoming("numbers")
@Outgoing("pow10")
@Broadcast
public BigInteger calculate(Long number) {
    return BigInteger.valueOf(number).pow(10);
}
}

```

Show moreShow more icon

代码清单 13 中的 `NumberResource` 类中的 `Publisher<BigInteger>` 类型的字段 `numbers` 表示的是来自通道 `pow10` 中的值。这些值以服务器推送事件的形式发送到浏览器。

##### 清单 13\. 消费流的 JAX-RS 资源

```
@Path("/numbers")
public class NumberResource {
@Inject
@Channel("pow10")
Publisher<BigInteger> numbers;

@GET
@Produces(MediaType.SERVER_SENT_EVENTS)
public Publisher<BigInteger> stream() {
    return numbers;
}
}

```

Show moreShow more icon

Quarkus 中对数据流的生成和消费都是抽象的。通过`@Outgoing` 和 `@Incoming` 注解可以把多个流串联起来，形成消息的处理链条。在抽象的数据流之下，SmallRye 依靠不同的连接器进行实际的消息传递。代码清单 14 中是 AMQP 相关的配置项。其中以 `mp.messaging.outgoing` 开头的配置项表示输出流的配置，而 `mp.messaging.incoming` 表示输入流的配置。可以对每个流进行配置，比如前缀 `mp.messaging.outgoing.generated-number` 表示的是输出流 `generated-number` 相关的配置项。在这些配置项中，`connector` 属性的值 `smallrye-amqp` 表示使用 SmallRye 的 AMQP 作为连接器实现，`address` 表示 `generated-number` 中的消息被发布到 `numbers` 流中。

##### 清单 14\. AMQP 相关的配置

```
amqp-host=localhost
amqp-port=5672
amqp-username=quarkus
amqp-password=quarkus
mp.messaging.outgoing.generated-number.connector=smallrye-amqp
mp.messaging.outgoing.generated-number.address=numbers
mp.messaging.outgoing.generated-number.durable=true
mp.messaging.incoming.numbers.connector=smallrye-amqp
mp.messaging.incoming.numbers.durable=true

```

Show moreShow more icon

### 使用事件总线

使用事件总线需要添加扩展 vertx。代码清单 15 中的 `EncodingService` 类的 `encode()` 方法对输入的数据进行 Base64 编码。注解 `@ConsumeEvent("encoding")` 表示该方法的输入参数来自名为 encoding 的事件。

##### 清单 15\. 事件的消费者

```
@ApplicationScoped
public class EncodingService {

@ConsumeEvent("encoding")
public String encode(String data) {
return Base64.getEncoder()
.encodeToString(data.getBytes(StandardCharsets.UTF_8));
}
}

```

Show moreShow more icon

代码清单 16 中的 `EncodingResource` 类是对应的 REST 服务资源。当接收到 HTTP 请求时，查询参数 `data` 的值通过 EventBus 的 `request()` 方法发送到事件总线上。事件的处理结果以 `CompletionStage<String>` 对象的形式返回。这说明对事件的处理是异步进行的。

##### 清单 16\. 使用事件总线的 JAX-RS 资源

```
@Path("/encoding")
public class EncodingResource {
@Inject
EventBus eventBus;

@GET
@Produces(MediaType.TEXT_PLAIN)
public CompletionStage<String> encode(@QueryParam("data") String data) {
    return eventBus.<String>request("encoding", data).thenApply(Message::body);
}
}

```

Show moreShow more icon

## 配置管理

Quarkus 应用使用 `application.properties` 文件来配置。该文件中可以包含 Quarkus 及其扩展所需的配置，也可以包含应用自身的配置。在代码中，可以通过 `@ConfigProperty` 注解来访问配置项的值。

代码清单 17 中的 `PowCalculator` 类使用 `@ConfigProperty` 注解来访问名为 `numbers.exponent` 的配置项，默认值为 `10`。

##### 清单 17\. 使用配置项

```
@ApplicationScoped
public class PowCalculator {
@ConfigProperty(name = "numbers.exponent", defaultValue = "10")
Integer exponent;

@Incoming("numbers")
@Outgoing("pow10")
@Broadcast
public BigInteger calculate(Long number) {
    return BigInteger.valueOf(number).pow(exponent);
}
}

```

Show moreShow more icon

如果有多个配置项属于同一个分组，可以 `@ConfigProperties` 注解对 Java 接口进行标注。在代码清单 18 中，`@ConfigProperties` 注解的 prefix 属性指定了分组中的配置项的前缀 `numbers`。`NumbersConfig` 接口的 `exponent()` 方法上的 `@ConfigProperty` 注解则只需要使用分组下的名称 `exponent` 即可。在使用时，只需要通过依赖注入的方式声明 `NumbersConfig` 对象即可。

##### 清单 18\. 使用 Java 接口

```
@ConfigProperties(prefix = "numbers")
public interface NumbersConfig {
@ConfigProperty(name = "exponent", defaultValue = "10")
Integer exponent();
}

```

Show moreShow more icon

Quarkus 的配置管理支持不同的概要文件（profile）。Quarkus 默认提供了 `dev`、`test` 和 `prod` 等三个概要文件，会在不同的运行场景下自动启用。可以通过 `%{profile}` 前缀来添加对于特定概要文件的配置项，如 `%dev.numbers.exponent=100` 为概要文件 `dev` 提供了不同的值。自定义的概要文件名称需要通过系统属性 `quarkus.profile` 或环境变量 `QUARKUS_PROFILE` 进行设置。

## 测试

在 Quarkus 应用中，可以使用 [`Rest Assured`](http://rest-assured.io/) 来测试 REST API。代码清单 19 中的 UserResourceTest 类测试 [代码清单 8](#清单-8-使用-jpa-的-userservice-类) 中的 `UserService` 类。`@QuarkusTest` 注解表明这是一个 Quarkus 测试用例。在 `testCreateUser()` 方法中，首先发送 POST 请求到 `/user` 来创建 `User` 对象，接着发送 GET 请求到 `/user` 来验证新创建的 `User` 对象出现在返回的列表中。

##### 清单 19\. 测试 REST 服务

```
@QuarkusTest
public class UserResourceTest {

@Test
public void testCreateUser() {
    given()
        .contentType(ContentType.JSON)
        .body(new User("test", "test@example.com"))
        .post("/user")
        .then()
        .statusCode(200)
        .body("name", equalTo("test"), "id", notNullValue());

    get("/user")
        .then()
        .statusCode(200)
        .body("name", hasItem("test"));
}

}

```

Show moreShow more icon

## 下载示例源码

本文示例代码可以在我的 GitHub 代码仓库里找到。

[获得代码](https://github.com/VividcodeIO/quarkus-starter)

## 结束语

Quarkus 为开发面向容器和 Kubernetes 平台的 Java 应用提供了一种新的选项。依靠 GraalVM 的原生镜像功能，Quarkus 应用有更快的启动时间和更少的内存占用。通过各种不同的扩展，Quarkus 应用所能实现的功能是非常丰富的。Quarkus 可以作为开发下一个微服务应用的良好选择。