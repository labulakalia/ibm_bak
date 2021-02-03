# 使用 JAX-RS 简化 REST 应用开发
通过一个完整的示例应用程序展示 JAX-RS 关键的设计细节以及与 JPA 的结合使用

**标签:** Java,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/j-lo-jaxrs/)

殷 钢, 张 昊

发布: 2009-09-07

* * *

## REST 简介

REST 是英文 Representational State Transfer 的缩写，有中文翻译为”具象状态传输”。REST 这个术语是由 Roy Fielding 在他的博士论文 _《 Architectural Styles and the Design of Network-based Software Architectures 》_ 中提出的。REST 并非标准，而是一种开发 Web 应用的架构风格，可以将其理解为一种设计模式。REST 基于 HTTP，URI，以及 XML 这些现有的广泛流行的协议和标准，伴随着 REST，HTTP 协议得到了更加正确的使用。

相较于基于 SOAP 和 WSDL 的 Web 服务，REST 模式提供了更为简洁的实现方案。目前，越来越多的 Web 服务开始采用 REST 风格设计和实现，真实世界中比较著名的 REST 服务包括：Google AJAX 搜索 API 、 [Amazon Simple Storage Service (Amazon S3)](http://aws.amazon.com/s3/) 等。

基于 REST 的 Web 服务遵循一些基本的设计原则：

- 系统中的每一个对象或是资源都可以通过一个唯一的 URI 来进行寻址，URI 的结构应该简单、可预测且易于理解，比如定义目录结构式的 URI。
- 以遵循 RFC-2616 所定义的协议的方式显式地使用 HTTP 方法，建立创建、检索、更新和删除（CRUD：Create, Retrieve, Update and Delete）操作与 HTTP 方法之间的一对一映射：

    - 若要在服务器上创建资源，应该使用 POST 方法；
    - 若要检索某个资源，应该使用 GET 方法；
    - 若要更改资源状态或对其进行更新，应该使用 PUT 方法；
    - 若要删除某个资源，应该使用 DELETE 方法。
- URI 所访问的每个资源都可以使用不同的形式加以表示（比如 XML 或者 JSON），具体的表现形式取决于访问资源的客户端，客户端与服务提供者使用一种内容协商的机制（请求头与 MIME 类型）来选择合适的数据格式，最小化彼此之间的数据耦合。

## JAX-RS — Java API for RESTful Web Services

Java EE 6 引入了对 JSR-311 的支持。JSR-311（JAX-RS：Java API for RESTful Web Services）旨在定义一个统一的规范，使得 Java 程序员可以使用一套固定的接口来开发 REST 应用，避免了依赖于第三方框架。同时，JAX-RS 使用 POJO 编程模型和基于标注的配置，并集成了 JAXB，从而可以有效缩短 REST 应用的开发周期。

JAX-RS 定义的 API 位于 javax.ws.rs 包中，其中一些主要的接口、标注和抽象类如 [图 1\. javax.ws.rs 包概况](#图-1-javax-ws-rs-包概况) 所示。

##### 图 1\. javax.ws.rs 包概况

![图 1. javax.ws.rs 包概况](../ibm_articles_img/j-lo-jaxrs_images_image001.jpg)

JAX-RS 的具体实现由第三方提供，例如 Sun 的参考实现 Jersey、Apache 的 CXF 以及 JBoss 的 RESTEasy。

在接下来的文章中，将结合一个记账簿应用向读者介绍 JAX-RS 一些关键的细节。

## 示例简介

记账簿示例应用程序中包含了 3 种资源：账目、用户以及账目种类，用户与账目、账目种类与账目之间都是一对多的关系。记账簿实现的主要功能包括：

1. 记录某用户在什么时间花费了多少金额在哪个种类上
2. 按照用户、账目种类、时间或者金额查询记录
3. 对用户以及账目种类的管理

## Resource 类和 Resource 方法

Web 资源作为一个 Resource 类来实现，对资源的请求由 Resource 方法来处理。Resource 类或 Resource 方法被打上了 Path 标注，Path 标注的值是一个相对的 URI 路径，用于对资源进行定位，路径中可以包含任意的正则表达式以匹配资源。和大多数 JAX-RS 标注一样，Path 标注是可继承的，子类或实现类可以继承超类或接口中的 Path 标注。

Resource 类是 POJO，使用 JAX-RS 标注来实现相应的 Web 资源。Resource 类分为根 Resource 类和子 Resource 类，区别在于子 Resource 类没有打在类上的 Path 标注。Resource 类的实例方法打上了 Path 标注，则为 Resource 方法或子 Resource 定位器，区别在于子 Resource 定位器上没有任何 @GET、@POST、@PUT、@DELETE 或者自定义的 @HttpMethod。 [清单 1\. 根 Resource 类](#清单-1-根-resource-类) 展示了示例应用中使用的根 Resource 类及其 Resource 方法。

##### 清单 1\. 根 Resource 类

```
@Path("/")
public class BookkeepingService {
    ......
    @Path("/person/")
    @POST
    @Consumes("application/json")
    public Response createPerson(Person person) {
        ......
    }

    @Path("/person/")
    @PUT
    @Consumes("application/json")
    public Response updatePerson(Person person) {
        ......
    }

    @Path("/person/{id:\\d+}/")
    @DELETE
    public Response deletePerson(@PathParam("id")
    int id) {
        ......
    }

    @Path("/person/{id:\\d+}/")
    @GET
    @Produces("application/json")
    public Person readPerson(@PathParam("id")
    int id) {
        ......
    }

    @Path("/persons/")
    @GET
    @Produces("application/json")
    public Person[] readAllPersons() {
        ......
    }

    @Path("/person/{name}/")
    @GET
    @Produces("application/json")
    public Person readPersonByName(@PathParam("name")
    String name) {
        ......
}
......

```

Show moreShow more icon

### 参数标注

JAX-RS 中涉及 Resource 方法参数的标注包括：@PathParam、@MatrixParam、@QueryParam、@FormParam、@HeaderParam、@CookieParam、@DefaultValue 和 @Encoded。这其中最常用的是 @PathParam，它用于将 @Path 中的模板变量映射到方法参数，模板变量支持使用正则表达式，变量名与正则表达式之间用分号分隔。例如对 [清单 1\. 根 Resource 类](#清单-1-根-resource-类) 中所示的 BookkeepingService 类，如果使用 Get 方法请求资源”/person/jeffyin”，则 readPersonByName 方法将被调用，方法参数 name 被赋值为”jeffyin”；而如果使用 Get 方法请求资源”/person/123”，则 readPerson 方法将被调用，方法参数 id 被赋值为 123。要了解如何使用其它的参数标注 , 请参考 JAX-RS API 。

JAX-RS 规定 Resource 方法中只允许有一个参数没有打上任何的参数标注，该参数称为实体参数，用于映射请求体。例如 清单 1 中所示的 BookkeepingService 类的 createPerson 方法和 updatePerson 方法的参数 person。

### 参数与返回值类型

Resource 方法合法的参数类型包括：

1. 原生类型
2. 构造函数接收单个字符串参数或者包含接收单个字符串参数的静态方法 valueOf 的任意类型
3. List，Set，SortedSet（T 为以上的 2 种类型）
4. 用于映射请求体的实体参数

Resource 方法合法的返回值类型包括：

1. void：状态码 204 和空响应体
2. Response：Response 的 status 属性指定了状态码，entity 属性映射为响应体
3. GenericEntity：GenericEntity 的 entity 属性映射为响应体，entity 属性为空则状态码为 204，非空则状态码为 200
4. 其它类型：返回的对象实例映射为响应体，实例为空则状态码为 204，非空则状态码为 200

对于错误处理，Resource 方法可以抛出非受控异常 WebApplicationException 或者返回包含了适当的错误码集合的 Response 对象。

### Context 标注

通过 Context 标注，根 Resource 类的实例字段可以被注入如下类型的上下文资源：

1. Request、UriInfo、HttpHeaders、Providers、SecurityContext
2. HttpServletRequest、HttpServletResponse、ServletContext、ServletConfig

要了解如何使用第 1 种类型的上下文资源 , 请参考 JAX-RS API 。

## CRUD 操作

JAX-RS 定义了 @POST、@GET、@PUT 和 @DELETE，分别对应 4 种 HTTP 方法，用于对资源进行创建、检索、更新和删除的操作。

### POST 标注

POST 标注用于在服务器上创建资源，如 [清单 2\. POST 标注](#清单-2-post-标注) 所示。

##### 清单 2\. POST 标注

```
@Path("/")
public class BookkeepingService {
    ......
    @Path("/account/")
    @POST
    @Consumes("application/json")
    public Response createAccount(Account account) {
        ......
    }
......

```

Show moreShow more icon

如果使用 POST 方法请求资源”/account”，则 createAccount 方法将被调用，JSON 格式的请求体被自动映射为实体参数 account。

### GET 标注

GET 标注用于在服务器上检索资源，如 [清单 3\. GET 标注](#清单-3-get-标注) 所示。

##### 清单 3\. GET 标注

```
@Path("/")
public class BookkeepingService {
    ......
    @Path("/person/{id}/accounts/")
    @GET
    @Produces("application/json")
    public Account[] readAccountsByPerson(@PathParam("id")
    int id) {
        ......
    }
    ......
    @Path("/accounts/{beginDate:\\d{4}-\\d{2}-\\d{2}},{endDate:\\d{4}-\\d{2}-\\d{2}}/")
    @GET
    @Produces("application/json")
    public Account[] readAccountsByDateBetween(@PathParam("beginDate")
    String beginDate, @PathParam("endDate")
    String endDate) throws ParseException {
        ......
    }
......

```

Show moreShow more icon

如果使用 GET 方法请求资源”/person/123/accounts”，则 readAccountsByPerson 方法将被调用，方法参数 id 被赋值为 123，Account 数组类型的返回值被自动映射为 JSON 格式的响应体；而如果使用 GET 方法请求资源”/accounts/2008-01-01,2009-01-01”，则 readAccountsByDateBetween 方法将被调用，方法参数 beginDate 被赋值为”2008-01-01”，endDate 被赋值为”2009-01-01”，Account 数组类型的返回值被自动映射为 JSON 格式的响应体。

### PUT 标注

PUT 标注用于更新服务器上的资源，如 [清单 4\. PUT 标注](#清单-4-put-标注) 所示。

##### 清单 4\. PUT 标注

```
@Path("/")
public class BookkeepingService {
    ......
    @Path("/account/")
    @PUT
    @Consumes("application/json")
    public Response updateAccount(Account account) {
        ......
    }
......

```

Show moreShow more icon

如果使用 PUT 方法请求资源”/account”，则 updateAccount 方法将被调用，JSON 格式的请求体被自动映射为实体参数 account。

### DELETE 标注

DELETE 标注用于删除服务器上的资源，如 [清单 5\. DELETE 标注](#清单-5-delete-标注) 所示。

##### 清单 5\. DELETE 标注

```
@Path("/")
public class BookkeepingService {
    ......
    @Path("/account/{id:\\d+}/")
    @DELETE
    public Response deleteAccount(@PathParam("id")
    int id) {
        ......
    }
......

```

Show moreShow more icon

如果使用 DELETE 方法请求资源”/account/323”，则 deleteAccount 方法将被调用，方法参数 id 被赋值为 323。

## 内容协商与数据绑定

Web 资源可以有不同的表现形式，服务端与客户端之间需要一种称为内容协商（Content Negotiation）的机制：作为服务端，Resource 方法的 Produces 标注用于指定响应体的数据格式（MIME 类型），Consumes 标注用于指定请求体的数据格式；作为客户端，Accept 请求头用于选择响应体的数据格式，Content-Type 请求头用于标识请求体的数据格式。

JAX-RS 依赖于 MessageBodyReader 和 MessageBodyWriter 的实现来自动完成返回值到响应体的序列化以及请求体到实体参数的反序列化工作，其中，XML 格式的请求／响应数据与 Java 对象的自动绑定依赖于 JAXB 的实现。

用户可以使用 Provider 标注来注册使用自定义的 MessageBodyProvider，如 [清单 6\. GsonProvider](#清单-6-gsonprovider) 所示，GsonProvider 类使用了 Google Gson 作为 JSON 格式的 MessageBodyProvider 的实现。

##### 清单 6\. GsonProvider

```
@Provider
@Produces("application/json")
@Consumes("application/json")
public class GsonProvider implements MessageBodyWriter<Object>,
    MessageBodyReader<Object> {

    private final Gson gson;

    public GsonProvider() {
        gson = new GsonBuilder().excludeFieldsWithoutExposeAnnotation().setDateFormat(
                "yyyy-MM-dd").create();
    }

    public boolean isReadable(Class<?> type, Type genericType, Annotation[] annotations,
            MediaType mediaType) {
        return true;
    }

    public Object readFrom(Class<Object> type, Type genericType,
            Annotation[] annotations, MediaType mediaType,
            MultivaluedMap<String, String> httpHeaders, InputStream entityStream)
            throws IOException, WebApplicationException {
        return gson.fromJson(new InputStreamReader(entityStream, "UTF-8"), type);
    }

    public boolean isWriteable(Class<?> type, Type genericType, Annotation[] annotations,
            MediaType mediaType) {
        return true;
    }

    public long getSize(Object obj, Class<?> type, Type genericType,
            Annotation[] annotations, MediaType mediaType) {
        return -1;
    }

    public void writeTo(Object obj, Class<?> type, Type genericType,
            Annotation[] annotations, MediaType mediaType,
            MultivaluedMap<String, Object> httpHeaders, OutputStream entityStream)
            throws IOException, WebApplicationException {
        entityStream.write(gson.toJson(obj, type).getBytes("UTF-8"));
    }

}

```

Show moreShow more icon

## JAX-RS 与 JPA 的结合使用

由于 JAX-RS 和 JPA 同样都使用了基于 POJO 和标注的编程模型，因而很易于结合在一起使用。示例应用中的 Web 资源 ( 如账目 ) 同时也是持久化到数据库中的实体，同一个 POJO 类上既有 JAXB 的标注，也有 JPA 的标注 ( 或者还有 Gson 的标注 ) ，这使得应用中类的个数得以减少。如 [清单 7\. Account](#清单-7-account) 所示，Account 类可以在 JAX-RS 与 JPA 之间得到复用，它不但可以被 JAX-RS 绑定为请求体 / 响应体的 XML/JSON 数据，也可以被 JPA 持久化到关系型数据库中。

##### 清单 7\. Account

```
@Entity
@Table(name = "TABLE_ACCOUNT")
@XmlRootElement
public class Account {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "COL_ID")
    @Expose
    private int id;

    @ManyToOne
    @JoinColumn(name = "COL_PERSON")
    @Expose
    private Person person;

    @Column(name = "COL_AMOUNT")
    @Expose
    private BigDecimal amount;

    @Column(name = "COL_DATE")
    @Expose
    private Date date;

    @ManyToOne
    @JoinColumn(name = "COL_CATEGORY")
    @Expose
    private Category category;

    @Column(name = "COL_COMMENT")
    @Expose
    private String comment;
......

```

Show moreShow more icon

## 结束语

REST 作为一种轻量级的 Web 服务架构被越来越多的开发者所采用，JAX-RS 的发布则规范了 REST 应用开发的接口。本文首先阐述了 REST 架构的基本设计原则，然后通过一个示例应用展示了 JAX-RS 是如何通过各种标注来实现以上的设计原则的，最后还介绍了 JAX-RS 与 JPA、Gson 的结合使用。本文的示例应用使用了 Jersey 和 OpenJPA，部署在 Tomcat 容器上，替换成其它的实现只需要修改 web.xml 和 persistence.xml 配置文件。

## 下载

[Bookkeeping.zip](https://www.ibm.com/developerworks/cn/java/j-lo-jaxrs/Bookkeeping.zip): 本文示例代码