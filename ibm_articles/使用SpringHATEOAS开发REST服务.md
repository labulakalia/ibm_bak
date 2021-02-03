# 使用 Spring HATEOAS 开发 REST 服务
了解 REST 和 HATEOAS 的相关概念及其框架的使用

**标签:** API 管理,Java,Spring

[原文链接](https://developer.ibm.com/zh/articles/j-lo-springhateoas/)

成富

发布: 2015-01-08

* * *

绝大多数开发人员对于 REST 这个词都并不陌生。自从 2000 年 Roy Fielding 在其博士论文中创造出来这个词之后，REST 架构风格就很快地流行起来，已经成为了构建 Web 服务时应该遵循的事实标准。很多 Web 服务和 API 都宣称满足了 REST 架构风格的要求，即所谓的”RESTful”服务。不过就如同其他很多流行的概念一样，不少人对于 REST 的含义还是存在或多或少的种种误解。REST 在某些时候被当成了一种营销的手段。不少所谓的”RESTful” Web 服务或 API 实际上并不满足 REST 架构风格的要求。这其中的部分原因在于 REST 的含义比较复杂，包含很多不同方面的内容。本文首先对 REST 架构做一个简单的说明以澄清某些误解。

## REST 架构

REST 是 Representational state transfer 的缩写，翻译过来的意思是表达性状态转换。REST 是一种架构风格，它包含了一个分布式超文本系统中对于组件、连接器和数据的约束。REST 是作为互联网自身架构的抽象而出现的，其关键在于所定义的架构上的各种约束。只有满足这些约束，才能称之为符合 REST 架构风格。REST 的约束包括：

- 客户端-服务器结构。通过一个统一的接口来分开客户端和服务器，使得两者可以独立开发和演化。客户端的实现可以简化，而服务器可以更容易的满足可伸缩性的要求。
- 无状态。在不同的客户端请求之间，服务器并不保存客户端相关的上下文状态信息。任何客户端发出的每个请求都包含了服务器处理该请求所需的全部信息。
- 可缓存。客户端可以缓存服务器返回的响应结果。服务器可以定义响应结果的缓存设置。
- 分层的系统。在分层的系统中，可能有中间服务器来处理安全策略和缓存等相关问题，以提高系统的可伸缩性。客户端并不需要了解中间的这些层次的细节。
- 按需代码（可选）。服务器可以通过传输可执行代码的方式来扩展或自定义客户端的行为。这是一个可选的约束。
- 统一接口。该约束是 REST 服务的基础，是客户端和服务器之间的桥梁。该约束又包含下面 4 个子约束。

    - 资源标识符。每个资源都有各自的标识符。客户端在请求时需要指定该标识符。在 REST 服务中，该标识符通常是 URI。客户端所获取的是资源的表达（representation），通常使用 XML 或 JSON 格式。
    - 通过资源的表达来操纵资源。客户端根据所得到的资源的表达中包含的信息来了解如何操纵资源，比如对资源进行修改或删除。
    - 自描述的消息。每条消息都包含足够的信息来描述如何处理该消息。
    - 超媒体作为应用状态的引擎（HATEOAS）。客户端通过服务器提供的超媒体内容中动态提供的动作来进行状态转换。这也是本文所要介绍的内容。

在了解 REST 的这些约束之后，就可以对”表达性状态转换”的含义有更加清晰的了解。”表达性”的含义是指对于资源的操纵都是通过服务器提供的资源的表达来进行的。客户端在根据资源的标识符获取到资源的表达之后，从资源的表达中可以发现其可以使用的动作。使用这些动作会发出新的请求，从而触发状态转换。

### HATEOAS 约束

HATEOAS（Hypermedia as the engine of application state）是 REST 架构风格中最复杂的约束，也是构建成熟 REST 服务的核心。它的重要性在于打破了客户端和服务器之间严格的契约，使得客户端可以更加智能和自适应，而 REST 服务本身的演化和更新也变得更加容易。

在介绍 HATEOAS 之前，先介绍一下 Richardson 提出的 REST 成熟度模型。该模型把 REST 服务按照成熟度划分成 4 个层次：

- 第一个层次（Level 0）的 Web 服务只是使用 HTTP 作为传输方式，实际上只是远程方法调用（RPC）的一种具体形式。SOAP 和 XML-RPC 都属于此类。
- 第二个层次（Level 1）的 Web 服务引入了资源的概念。每个资源有对应的标识符和表达。
- 第三个层次（Level 2）的 Web 服务使用不同的 HTTP 方法来进行不同的操作，并且使用 HTTP 状态码来表示不同的结果。如 HTTP GET 方法来获取资源，HTTP DELETE 方法来删除资源。
- 第四个层次（Level 3）的 Web 服务使用 HATEOAS。在资源的表达中包含了链接信息。客户端可以根据链接来发现可以执行的动作。

从上述 REST 成熟度模型中可以看到，使用 HATEOAS 的 REST 服务是成熟度最高的，也是推荐的做法。对于不使用 HATEOAS 的 REST 服务，客户端和服务器的实现之间是紧密耦合的。客户端需要根据服务器提供的相关文档来了解所暴露的资源和对应的操作。当服务器发生了变化时，如修改了资源的 URI，客户端也需要进行相应的修改。而使用 HATEOAS 的 REST 服务中，客户端可以通过服务器提供的资源的表达来智能地发现可以执行的操作。当服务器发生了变化时，客户端并不需要做出修改，因为资源的 URI 和其他信息都是动态发现的。

## 示例

本文将通过一个 [完整的示例](http://www.ibm.com/developerworks/cn/java/j-lo-SpringHATEOAS/sample_code.zip) 来说明 HATEOAS。该示例是一个常见的待办事项的服务，用户可以创建新的待办事项、进行编辑或标记为已完成。该示例中包含的资源如下：

- 用户：应用中的用户。
- 列表：待办事项的列表，属于某个用户。
- 事项：具体的待办事项，属于某个列表。

应用提供相关的 REST 服务来完成对于列表和事项两个资源的 CRUD 操作。

## Spring HATEOAS

如果 Web 应用基于 Spring 框架开发，那么可以直接使用 Spring 框架的子项目 HATEOAS 来开发满足 HATEOAS 约束的 Web 服务。本文的示例应用基于 Java 8 和使用 Spring Boot 1.1.9 来创建，Spring HATEOAS 的版本是 0.16.0.RELEASE。

### 基本配置

满足 HATEOAS 约束的 REST 服务最大的特点在于服务器提供给客户端的表达中包含了动态的链接信息，客户端通过这些链接来发现可以触发状态转换的动作。Spring HATEOAS 的主要功能在于提供了简单的机制来创建这些链接，并与 Spring MVC 框架有很好的集成。对于已有的 Spring MVC 应用，只需要一些简单的改动就可以满足 HATEOAS 约束。对于一个 Maven 项目来说，只需要添加清单 1 中的依赖即可。

##### 清单 1\. Spring HATEOAS 的 Maven 依赖声明

```
<dependency>
<groupId>org.springframework.hateoas</groupId>
<artifactId>spring-hateoas</artifactId>
<version>0.16.0.RELEASE</version>
</dependency>

```

Show moreShow more icon

### 资源

REST 架构中的核心概念之一是资源。服务器提供的是资源的表达，通常使用 JSON 或 XML 格式。在一般的 Web 应用中，服务器端代码会对所使用的资源建模，提供相应的模型层 Java 类，这些模型层 Java 类通常包含 JPA 相关的注解来完成持久化。在客户端请求时，服务器端代码通过 Jackson 或 JAXB 把模型对象转换成 JSON 或 XML 格式。清单 2 给出了示例应用中表示列表的模型类 List 的声明。

##### 清单 2\. 表示列表的模型类 List 的声明

```
@Entity
public class List extends AbstractEntity {
private String name;

@ManyToOne
@JsonIgnore
private User user;

@OneToMany(mappedBy = "list", fetch = FetchType.LAZY)
@JsonIgnore
private Set<Item> items = new HashSet<>();

protected List() {
}

public List(String name, User user) {
this.name = name;
this.user = user;
}

public String getName() {
return name;
}

public User getUser() {
return user;
}

public Set<Item> getItems() {
return items;
}
}

```

Show moreShow more icon

当客户端请求某个具体的 List 类的对象时，服务器端返回如清单 3 所示的 JSON 格式的表达。

##### 清单 3\. List 类的对象的 JSON 格式的表达

```
{
"id": 1,
"name": "Default"
}

```

Show moreShow more icon

在清单 3 中，服务器端返回的只是模型类对象本身的内容，并没有提供相关的链接信息。为了把模型对象类转换成满足 HATEOAS 要求的资源，需要添加链接信息。Spring HATEOAS 使用 org.springframework.hateoas.Link 类来表示链接。Link 类遵循 Atom 规范中对于链接的定义，包含 rel 和 href 两个属性。属性 rel 表示的是链接所表示的关系（relationship），href 表示的是链接指向的资源标识符，一般是 URI。资源通常都包含一个属性 rel 值为 self 的链接，用来指向该资源本身。

在创建资源类时，可以继承自 Spring HATEOAS 提供的 org.springframework.hateoas.Resource 类，Resource 类提供了简单的方式来创建链接。清单 4给出了与模型类 List 对应的资源类 ListResource 的声明。

##### 清单 4\. 模型类 List 对应的资源类 ListResource 的声明

```
public class ListResource extends Resource {
private final List list;

public ListResource(List list) {
super(list);
this.list = list;
add(new Link("http://localhost:8080/lists/1"));
add(new Link("http://localhost:8080/lists/1/items", "items"));
}

public List getList() {
return list;
}
}

```

Show moreShow more icon

如清单 4 所示，ListResource 类继承自 Resource 类并对 List 类的对象进行了封装，添加了两个链接。在使用 ListResource 类之后，服务器端返回的表达格式如清单 5 所示。

##### 清单 5\. 使用 ListResource 类之后的 JSON 格式的表达

```
{
"list": {
"id": 1,
"name": "Default"
},
"links": [
{
"rel": "self",
"href": "http://localhost:8080/lists/1"
},
{
"rel": "items",
"href": "http://localhost:8080/lists/1/items"
}
]
}

```

Show moreShow more icon

清单 5 的 JSON 内容中添加了额外的 links 属性，并包含了两个链接。不过模型类对象的内容被封装在属性 list 中。这是因为 ListResource 类直接封装了整个的 List 类的对象，而不是把 List 类的属性提取到 ListResource 类中。如果需要改变输出的 JSON 表达的格式，可以使用另外一种封装方式的 ListResource 类，如清单 6 所示。

##### 清单 6\. 不同封装格式的 ListResource 类

```
public class ListResource extends Resource {
private final Long id;
private final String name;

public ListResource(List list) {
super(list);
this.id = list.getId();
this.name = list.getName();
add(new Link("http://localhost:8080/lists/1"));
add(new Link("http://localhost:8080/lists/1/items", "items"));
}

public Long getId() {
return id;
}
public String getName() {
return name;
}
}

```

Show moreShow more icon

对应的资源的表达如清单 7 所示。

##### 清单 7\. 使用不同封装方式的 JSON 格式的表达

```
{
"id": 1,
"name": "Default",
"links": [
{
"rel": "self",
"href": "http://localhost:8080/lists/1"
},
{
"rel": "items",
"href": "http://localhost:8080/lists/1/items"
}
]
}

```

Show moreShow more icon

这两种不同的封装方式各有优缺点。第一种方式的优点是实现起来很简单，只需要把模型层的对象直接包装即可；第二种方式虽然实现起来相对比较复杂，但是可以对资源的表达格式进行定制，使得资源的表达格式更直接。

在代码实现中经常会需要把模型类对象转换成对应的资源对象，如把 List 类的对象转换成 ListResource 类的对象。一般的做法是通过”new ListResource(list)”这样的方式来进行转换。可以使用 Spring HATEOAS 提供的资源组装器把转换的逻辑封装起来。资源组装器还可以自动创建 rel 属性为 self 的链接。 清单 8 中给出了组装资源类 ListResource 的 ListResourceAssembler 类的实现。

##### 清单 8\. 组装资源类 ListResource 的 ListResourceAssembler 类的实现

```
public class ListResourceAssembler extends ResourceAssemblerSupport<List, ListResource> {

public ListResourceAssembler() {
super(ListRestController.class, ListResource.class);
}

@Override
public ListResource toResource(List list) {
ListResource resource = createResourceWithId(list.getId(), list);
return resource;
}

@Override
protected ListResource instantiateResource(List entity) {
return new ListResource(entity);
}
}

```

Show moreShow more icon

在创建 ListResourceAssembler 类的对象时需要指定使用资源的 Spring MVC 控制器 Java 类和资源 Java 类。对于 ListResourceAssembler 类来说分别是 ListRestController 和 ListResource。ListRestController 类在下一节中会具体介绍，其作用是用来创建 rel 属性为 self 的链接。ListResourceAssembler 类的 instantiateResource 方法用来根据一个模型类 List 的对象创建出 ListResource 对象。ResourceAssemblerSupport 类的默认实现是通过反射来创建资源对象的。toResource 方法用来完成实际的转换。此处使用了 ResourceAssemblerSupport 类的 createResourceWithId 方法来创建一个包含 self 链接的资源对象。

在代码中需要创建 ListResource 的地方，都可以换成使用 ListResourceAssembler，如清单 9 所示。

##### 清单 9\. 使用 ListResourceAssembler 的示例

```
//组装单个资源对象
new ListResourceAssembler().toResource(list);

//组装资源对象的集合
new ListResourceAssembler().toResources(lists);

```

Show moreShow more icon

清单 9 中的 toResources 方法是 ResourceAssemblerSupport 类提供的。当需要转换一个集合的资源对象时，这个方法非常实用。

### 链接

HATEOAS 的核心是链接。链接的存在使得客户端可以动态发现其所能执行的动作。在上一节中介绍过链接由 rel 和 href 两个属性组成。其中属性 rel 表明了该链接所代表的关系含义。应用可以根据需要为链接选择最适合的 rel 属性值。由于每个应用的情况并不相同，对于应用相关的 rel 属性值并没有统一的规范。不过对于很多常见的链接关系，IANA 定义了规范的 rel 属性值。在开发中可能使用的常见 rel 属性值如表 1 所示。

##### 表 1\. 常用的 rel 属性

rel 属性值描述self指向当前资源本身的链接的 rel 属性。每个资源的表达中都应该包含此关系的链接。edit指向一个可以编辑当前资源的链接。item如果当前资源表示的是一个集合，则用来指向该集合中的单个资源。collection如果当前资源包含在某个集合中，则用来指向包含该资源的集合。related指向一个与当前资源相关的资源。search指向一个可以搜索当前资源及其相关资源的链接。first、last、previous、next这几个 rel 属性值都有集合中的遍历相关，分别用来指向集合中的第一个、最后一个、上一个和下一个资源。

如果在应用中使用自定义 rel 属性值，一般的做法是属性值全部为小写，中间使用”-”分隔。

链接中另外一个重要属性 href 表示的是资源的标识符。对于 Web 应用来说，通常是一个 URL。URL 必须指向的是一个绝对的地址。在应用中创建链接时，在 URL 中使用硬编码的主机名和端口号显然不是好的选择。Spring MVC 提供了相关的工具类可以获取 Web 应用启动时的主机名和端口号，不过创建动态的链接 URL 还需要可以获取资源的访问路径。对于一个典型的 Spring MVC 控制器来说，其声明如清单 10 所示。

##### 清单 10\. Spring MVC 控制器 ListRestController 类的实现

```
@RestController
@RequestMapping("/lists")
public class ListRestController {

@Autowired
private ListService listService;

@RequestMapping(method = RequestMethod.GET)
public Resources<ListResource> readLists(Principal principal) {
String username = principal.getName();
return new Resources<ListResource>(new
ListResourceAssembler().toResources(listService.findByUserUsername(username)));
@RequestMapping(value = "/{listId}", method = RequestMethod.GET)
public ListResource readList(@PathVariable Long listId) {
return new ListResourceAssembler().toResource(listService.findOne(listId));
}
}

```

Show moreShow more icon

从清单 10 中可以看到，Spring MVC 控制器 ListRestController 类通过”@RequestMapping”注解声明了其访问路径是”/lists”，而访问单个资源的路径是类似”/lists/1”这样的形式。在创建资源的链接时，指向单个资源的链接的 href 属性值是类似”`http://localhost:8080/lists/1`”这样的格式。而其中的”/lists”不应该是硬编码的，否则当修改了 ListRestController 类的”@RequestMapping”时，所有相关的生成链接的代码都需要进行修改。Spring HATEOAS 提供了 org.springframework.hateoas.mvc.ControllerLinkBuilder 来解决这个问题，用来根据 Spring MVC 控制器动态生成链接。清单 11 给出了创建单个资源的链接的方式。

##### 清单 11\. 使用 ControllerLinkBuilder 类创建链接

```
import static org.springframework.hateoas.mvc.ControllerLinkBuilder.*;

Link link = linkTo(ListRestController.class).slash(listId).withSelfRel();

```

Show moreShow more icon

通过 ControllerLinkBuilder 类的 linkTo 方法，先指定 Spring MVC 控制器的 Java 类，再通过 slash 方法来找到下一级的路径，最后生成属性值为 self 的链接。在使用 ControllerLinkBuilder 生成链接时，除了可以使用控制器的 Java 类之外，还可以使用控制器 Java 类中包含的方法。如清单 12 所示。

##### 清单 12\. 通过控制器 Java 类中的方法生成链接

```
Link link = linkTo(methodOn(ItemRestController.class).readItems(listId)).withRel("items");

```

Show moreShow more icon

清单 12 中的链接使用的是 ItemRestController 类中的 readItems 方法。参数 listId 是组成 URI 的一部分，在调用 readItems 方法时需要提供。

上面介绍的是通过 Spring MVC 控制器来创建链接，另外一种做法是从模型类中创建。这是因为控制器通常用来暴露某个模型类。如 ListRestController 类直接暴露模型类 List，并提供了访问 List 资源集合和单个 List 资源的接口。对于这样的情况，并不需要通过控制器来创建相关的链接，而可以使用 EntityLinks。

首先需要在控制器类中通过”@ExposesResourceFor”注解声明其所暴露的模型类，如清单 13 中的 ListRestController 类的声明。

##### 清单 13\. “@ExposesResourceFor”注解的使用

```
@RestController
@ExposesResourceFor(List.class)
@RequestMapping("/lists")
public class ListRestController {

}

```

Show moreShow more icon

另外在 Spring 应用的配置类中需要通过”@EnableEntityLinks”注解来启用 EntityLinks 功能。此外还需要添加清单 14 中给出的 Maven 依赖。

##### 清单 14\. EntityLinks 功能所需的 Maven 依赖

```
<dependency>
<groupId>org.springframework.plugin</groupId>
<artifactId>spring-plugin-core</artifactId>
<version>1.1.0.RELEASE</version>
</dependency>

```

Show moreShow more icon

在需要创建链接的代码中，只需要通过依赖注入的方式添加对 EntityLinks 的引用，就可以使用 linkForSingleResource 方法来创建指向单个资源的链接，如 清单 15 所示。

##### 清单 15\. 使用 EntityLinks 创建链接

```
@Autowired
private EntityLinks entityLinks;

entityLinks.linkForSingleResource(List.class, 1)

```

Show moreShow more icon

需要注意的是，为了 linkForSingleResource 方法可以正常工作，控制器类中需要包含访问单个资源的方法，而且其”@RequestMapping”是类似”/{id}”这样的形式。

## 超媒体控制与 HAL

在添加了链接之后，服务器端提供的表达可以帮助客户端更好的发现服务器端所支持的动作。在具体的表达中，应用虽然可以根据需要选择最适合的格式，但是在表达的基本结构上应该遵循一定的规范，这样可以保证最大程度的适用性。这个基本结构主要是整体的组织方式和链接的格式。HAL（Hypertxt Application Language）是一个被广泛采用的超文本表达的规范。应用可以考虑遵循该规范，Spring HATEOAS 提供了对 HAL 的支持。

### HAL 规范

HAL 规范本身是很简单的，清单 16 给出了示例的 JSON 格式的表达。

##### 清单 16\. HAL 规范的示例 JSON 格式的表达

```
{
"_links": {
"self": {
"href": "http://localhost:8080/lists"
}
},
"_embedded": {
"lists": [
{
"id": 1,
"name": "Default",
"_links": {
"todo:items": {
"href": "http://localhost:8080/lists/1/items"
},
"self": {
"href": "http://localhost:8080/lists/1"
},
"curies": [
{
"href": "http://www.midgetontoes.com/todolist/rels/{rel}",
"name": "todo",
"templated": true
}
]
}
}
]
}
}

```

Show moreShow more icon

HAL 规范围绕资源和链接这两个简单的概念展开。资源的表达中包含链接、嵌套的资源和状态。资源的状态是该资源本身所包含的数据。链接则包含其指向的目标（URI）、所表示的关系和其他可选的相关属性。对应到 JSON 格式中，资源的链接包含在\_links 属性对应的哈希对象中。该\_links 哈希对象中的键（key）是链接的关系，而值（value）则是另外一个包含了 href 等其他链接属性的对象或对象数组。当前资源中所包含的嵌套资源由\_embeded 属性来表示，其值是一个包含了其他资源的哈希对象。

链接的关系不仅是区分不同链接的标识符，同样也是指向相关文档的 URL。文档用来告诉客户端如何对该链接所指向的资源进行操作。当开发人员获取到了资源的表达之后，可以通过查看链接指向的文档来了解如何操作该资源。

使用 URL 作为链接的关系带来的问题是 URL 作为属性名称来说显得过长，而且不同关系的 URL 的大部分内容是重复的。为了解决这个问题，可以使用 Curie。简单来说，Curie 可以作为链接关系 URL 的模板。链接的关系声明时使用 Curie 的名称作为前缀，不用提供完整的 URL。应用中声明的 Curie 出现在\_links 属性中。代码中定义了 URI 模板为 `http://www.midgetontoes.com/todolist/rels/{rel}` 的名为 todo 的 Curie。在使用了 Curie 之后，名为 items 的链接关系变成了包含前缀的”todo:items”的形式。这就表示该链接的关系实际上是”`http://www.midgetontoes.com/todolist/rels/items`”。

### Spring HATEOAS 的 HAL 支持

目前 Spring HATEOAS 仅支持 HAL 一种超媒体表达格式，只需要在应用的配置类上添加”@EnableHypermediaSupport(type= {HypermediaType.HAL})”注解就可以启用该超媒体支持。在启用了超媒体支持之后，服务器端输出的表达格式会遵循 HAL 规范。另外，启用超媒体支持会默认启用”@EnableEntityLinks”。在启用超媒体支持之后，应用需要进行相关的定制使得生成的 HAL 表达更加友好。

首先是内嵌资源在\_embedded 对应的哈希对象中的属性值，该属性值是由 org.springframework.hateoas.RelProvider 接口的实现来提供的。对于应用来说，只需要在内嵌资源对应的模型类中添加 org.springframework.hateoas.core.Relation 注解即可，如清单 17 所示。

##### 清单 17\. 在模型类中添加 @Relation 注解

```
@Relation(value = "list", collectionRelation = "lists")
public class List extends AbstractEntity {
}

```

Show moreShow more icon

清单 17 中声明了当模型类 List 的对象作为内嵌资源时，单个资源使用 list 作为属性值，多个资源使用 lists 作为属性值。

如果需要添加 Curie，则提供 org.springframework.hateoas.hal.CurieProvider 接口的实现，如清单 18 所示。利用已有的 org.springframework.hateoas.hal.DefaultCurieProvider 类并提供 Curie 的前缀和 URI 模板即可。

##### 清单 18\. 添加 CurieProvider 接口的实现

```
@Bean
public CurieProvider curieProvider() {
return new DefaultCurieProvider("todo",
new UriTemplate("http://www.midgetontoes.com/todolist/rels/{rel}"));
}

```

Show moreShow more icon

## 结束语

在开发一个新的 Web 服务或 API 时，REST 架构风格已经成为事实上的标准。在开发时需要明白 REST 架构风格中所包含的约束的含义。HATEOAS 作为 REST 服务约束中最复杂的一个，目前还没有得到广泛的使用。但是采用 HATEOAS 所带来的好处是很大的，可以帮助客户端和服务器更好的解耦，可以减少很多潜在的问题。Spring HATEOAS 在 Spring MVC 框架的基础上，允许开发人员通过简单的配置来添加 HATEOAS 约束。如果应用本身已经使用了 Spring MVC，则同时启用 HATEOAS 是一个很好的选择。本文对 REST 和 HATEOAS 的相关概念以及 Spring HATEOAS 框架的使用做了详细的介绍。

## 下载示例代码

[sample\_code.zip](http://www.ibm.com/developerworks/cn/java/j-lo-SpringHATEOAS/sample_code.zip): 示例代码下载