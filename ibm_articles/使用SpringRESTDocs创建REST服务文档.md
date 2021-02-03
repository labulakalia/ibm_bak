# 使用 Spring REST Docs 创建 REST 服务文档
基本的配置和对 HTTP 请求和响应的不同部分添加文档

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-use-spring-rest-docs-to-creat-rest/)

成富

发布: 2018-04-17

* * *

当一个应用需要对第三方提供服务接口时，REST API 是目前主流的选择。一个良好的 REST API 必须由同样高质量的文档的支持，否则第三方无法正确有效的使用 REST API。高质量的 REST API 文档的创建并不是一件容易的事情。创建文档并不难，难的是如何维护文档，让文档与代码的变化保持同步。错误的文档还不如没有文档。如果第三方按照文档的要求使用a API，却得到了错误的结果，那会极大降低他们对产品的信任度。Spring REST Doc 给出了创建和维护 REST API 文档的另外一种思路，也就是手写文档和自动生成相结合。手写文档用来提供必要的背景知识和相关介绍，而基于单元测试的自动生成文档用来产生真实有效的 HTTP 请求和响应相关的内容。由于自动生成的部分是基于单元测试的，可以保证文档的准确性，否则单元测试就会失败。

## Spring REST Docs

我们首先对 Spring REST Docs 进行基本的介绍。前面我们提到了 Spring REST Docs 采用的是手写文档和自动生成相结合的形式。手写文档时可以使用 Asciidoc 和 Markdown 两种格式。推荐使用支持度更好的 Asciidoc。Spring 项目的文档也是使用 Asciidoc 来编写的。相对于更流行的 Markdown 来说，Asciidoc 的语法较为复杂，相应的功能也更强大。本文的示例基于 Asciidoc，也会介绍 Markdown 相关的内容。

在自动生成部分，Spring REST Docs 依赖单元测试来产生所需的 HTTP 请求和响应的文档。在单元测试框架部分，Spring REST Docs 支持主流的单元测试框架，包括 JUnit 4、JUnit 5 和 TestNG。测试 REST API 通常需要专门的工具来支持 HTTP 请求的发送和响应结果的验证。Spring REST Docs 提供了对 Spring MVC 的 MockMvc、WebFlux 的 WebTestClient 和 REST Assured 的支持。本文的示例使用的是 JUnit 4 和 MockMvc。

Spring REST Docs 通常是作为已有项目的一部分，与需要编写文档的 REST API 的模块属于同一个项目。以 Gradle 为例，在已有的项目中添加 Spring REST Docs 的支持时，需要首先添加对 org.springframework.restdocs:spring-restdocs-mockmvc:2.0.0-RELEASE 的依赖。接着对 build.gradle 进行修改，如 [清单 1\. 使用 Spring REST Docs 的 Gradle 配置](#清单-1-使用-spring-rest-docs-的-gradle-配置) 所示。除了 Spring Boot 相关的配置之外，添加的内容包括启用插件 org.asciidoctor.convert 来把 Asciidoc 文件转换成 HTML。这里使用的是 1.5.3 版本。较新的 1.5.6 版本在使用 Gradle 构建时有问题，不建议使用。参数 snippetsDir 表示的是单元测试对应生成的 Asciidoc 片段所在的目录。配置 Asciidoctor 插件以该目录作为片段的输入，这样就可以在 Asciidoc 文档中包含这些代码片段。

#### 清单 1\. 使用 Spring REST Docs 的 Gradle 配置

```
buildscript {
ext {
springBootVersion = '2.0.0.RELEASE'
}
repositories {
mavenCentral()
}
dependencies {
classpath("org.springframework.boot:spring-boot-gradle-plugin:${springBootVersion}")
}
}
plugins {
id "org.asciidoctor.convert" version "1.5.3"
}
apply plugin: 'java'
apply plugin: 'org.springframework.boot'
apply plugin: 'io.spring.dependency-management'
group = 'io.vividcode'
version = '0.0.1-SNAPSHOT'
sourceCompatibility = 1.8
repositories {
mavenCentral()
}
ext {
snippetsDir = file('build/generated-snippets')
}
dependencies {
asciidoctor "org.springframework.restdocs:spring-restdocs-asciidoctor:2.0.0.RELEASE"
compile('org.springframework.boot:spring-boot-starter-data-rest')
compile('org.springframework.boot:spring-boot-starter-data-jpa')
compile('org.springframework.boot:spring-boot-starter-web')
runtime 'com.h2database:h2'
runtime 'org.atteo:evo-inflector:1.2.1'
testCompile 'com.jayway.jsonpath:json-path'
testCompile 'org.springframework.boot:spring-boot-starter-test'
testCompile 'org.springframework.restdocs:spring-restdocs-mockmvc'
}
test {
outputs.dir snippetsDir
}
asciidoctor {
inputs.dir snippetsDir
dependsOn test
}

```

Show moreShow more icon

本文的示例应用是一个基于 Spring Data 的 REST API，用来管理用户的书签。应用中有两类领域实体，分别是用户（User）和书签（Bookmark）。API 提供了对这两类实体的增删改查操作。应用的实现很简单，这里不进行介绍，可以直接查看附带的源代码。

为了生成文档中需要的 HTTP 请求和响应相关的代码片段，需要编写对应的单元测试用例。 [清单 2\. 使用 Spring REST Docs 的单元测试用例](#清单-2-使用-spring-rest-docs-的单元测试用例) 给出了单元测试用例类的示例。JUnitRestDocumentation 的作用是根据 MockMvc 的测试过程生成对应的代码片段。相应的配置在 setUp 方法中完成。在测试用例 indexExample 中，通过 andDo 来添加相应的创建文档的动作。document(“index-example”)的含义是产生名称为 index-example 的代码片段。

#### 清单 2\. 使用 Spring REST Docs 的单元测试用例

```
@RunWith(SpringRunner.class)
@SpringBootTest
public class ApiDocumentation {
@Rule
public final JUnitRestDocumentation restDocumentation = new JUnitRestDocumentation();
@Autowired
private UserRepository userRepository;
@Autowired
private BookmarkRepository bookmarkRepository;
@Autowired
private WebApplicationContext context;
private MockMvc mockMvc;
@Before
public void setUp() {
this.mockMvc = MockMvcBuilders.webAppContextSetup(this.context)
.apply(documentationConfiguration(this.restDocumentation)).build();
}
@Test
public void indexExample() throws Exception {
this.mockMvc.perform(get("/"))
.andExpect(status().isOk())
.andDo(document("index-example"));
}
}

```

Show moreShow more icon

通过 Gradle 命令运行测试之后，可以在 build/generated-snippets 目录下看到生成的 Asciidoc 文档片段。其中的每个目录代表一个文档中可以使用的代码片段，通常与一个单元测试用例相对应。目录的名称与调用 document 方法时给定的名称相同。在目录下面默认有 6 个 Asciidoc 文件，分别表示与 HTTP 请求和响应相关的内容，具体的说明如下所示。

- curl-request.adoc：使用 curl 发送请求的格式
- httpie-request.adoc：使用 HTTPie 发送请求的格式
- http-request.adoc：完整的 HTTP 请求的内容
- http-response.adoc：完整的 HTTP 响应的内容
- request-body.adoc：HTTP 请求体的内容
- response-body.adoc：HTTP 响应体的内容

最终要生成的 Asciidoc 文档应该放在 src/docs/asciidoc 目录下。在 Asciidoc 文档中，除了手动创建的内容之外，还可以添加之前生成的代码片段。 [清单 3\. API 文档 api.adoc](#清单-3-api-文档-api-adoc) 是 API 文件的 Asciidoc 文件 api.adoc 的内容，其中 operation::index-example[]用来包含 index-example 的代码片段。

#### 清单 3\. API 文档 api.adoc

```
= REST API 文档
:doctype: book
:icons: font
:source-highlighter: highlightjs
:toc: left
:toclevels: 4
:sectlinks:
== 首页
使用 `GET` 请求来访问首页。
operation::index-example[]

```

Show moreShow more icon

在运行 Gradle 构建之后，可以在 build/asciidoc/html5 下面找到生成的 api.html 文件，这是最终的 API 文档，可以直接在浏览器中查看。

## 为 API 添加文档

对一个 REST API 来说，最重要的文档是说明请求和响应的内容格式，告诉 API 使用者以什么样的格式发送请求，会得到什么样格式的响应。下面对请求和响应的不同组成部分进行分别的介绍。

### 超媒体链接

超媒体链接是 REST API 的一个重要组成部分。通过这些链接可以访问与当前资源相关的其他资源。Spring REST Docs 也提供了对超媒体链接的文档支持。默认支持的包括 Atom 和 HAL 两种格式。Atom 使用 links 字段来表示链接，HAL 使用\_links 来表示链接。本文的实例应用使用的是 HAL 规范。

[清单 4\. 为超媒体链接添加文档](#清单-4-为超媒体链接添加文档) 对上一节的测试用例 indexExample 进行修改，添加了链接相关的文档。方法 links 用来对链接添加文档。linkWithRel 描述了每个链接的关系及其说明。

#### 清单 4\. 为超媒体链接添加文档

```
@Test
public void indexExample() throws Exception {
this.mockMvc.perform(get("/"))
.andExpect(status().isOk())
.andDo(document("index-example",
links(
linkWithRel("users").description("指向用户资源的链接"),
linkWithRel("bookmarks").description("指向书签资源的链接"),
linkWithRel("profile").description("指向 ALPS 服务概要文件的链接")),
responseFields(
subsectionWithPath("_links").description("到其他资源的链接"))));
}

```

Show moreShow more icon

在描述了链接之后，会发现 build/generated-snippets 目录下对应测试的子目录 index-example 下会增加一个新的 links.adoc 文件，里面包含了链接的信息。这个文件可以在 API 文档中引用。

### 请求和响应的字段说明

Spring REST Docs 已经默认生成了 HTTP 请求和响应的完整内容，就是之前提到的 http-request.adoc 和 http-response.adoc 两个文件。请求和响应一般使用 JSON 或 XML 格式。对于其中包含的每个字段，都应该添加相应的文档说明。

[清单 5\. 为响应的字段添加文档](#清单-5-为响应的字段添加文档) 给出了对响应中包含的字段添加文档的示例。这个用例测试的是获取用户信息的 API。首先使用 POST 请求创建一个新用户，然后再通过 GET 请求获取该用户的信息，并验证响应中包含的内容。responseFields 方法用来声明响应中的字段，fieldWithPath 方法描述单个字段的路径和含义，subsectionWithPath 方法描述一个路径下的全部内容。

#### 清单 5\. 为响应的字段添加文档

```
@Test
public void userGetExample() throws Exception {
Map<String, Object> user = new HashMap<>();
user.put("username", "alex");
user.put("email", "alex@example.com");
String userLocation = this.mockMvc
.perform(post("/users").contentType(MediaTypes.HAL_JSON).content(this.objectMapper.wr
iteValueAsString(user)))
.andExpect(status().isCreated()).andReturn().getResponse().getHeader("Location");
this.mockMvc.perform(get(userLocation))
.andExpect(status().isOk())
.andExpect(jsonPath("username", is(user.get("username"))))
.andExpect(jsonPath("email", is(user.get("email"))))
.andExpect(jsonPath("_links.self.href", is(userLocation)))
.andExpect(jsonPath("_links.bookmarks", is(notNullValue())))
.andDo(document("user-get-example",
links(
linkWithRel("self").description("当前资源的链接"),
linkWithRel("user").description("当前用户实体的链接"),
linkWithRel("bookmarks").description("当前用户书签的链接")),
responseFields(
fieldWithPath("username").description("用户名"),
fieldWithPath("email").description("Email 地址"),
subsectionWithPath("_links").description("到其他资源的链接")
)));
}

```

Show moreShow more icon

使用 responseFields 方法会产生新的代码片段文件 response-fields.adoc，可以在 api.adoc 中引用。在引用 user-get-example 时，通过 snippets 属性限制了出现在文档中的代码片段的数量，如 [清单 6\. 在 API 文档中添加新的代码片段文件](#清单-6-在-api-文档中添加新的代码片段文件) 所示。

#### 清单 6\. 在 API 文档中添加新的代码片段文件

```
== 获取用户信息
使用 `GET` 请求来获取用户信息。
operation::user-get-example[snippets='response-fields,curl-request,http-response']

```

Show moreShow more icon

### 请求参数

HTTP 请求的参数也应该添加文档。 [清单 7\. 为请求参数添加文档](#清单-7-为请求参数添加文档) 中，requestParameters 方法用来添加请求参数的文档，parameterWithName 方法表示单个参数的名称和含义。所对应的代码片段的文件名是 request-parameters.adoc。

#### 清单 7\. 为请求参数添加文档

```
@Test
public void usersListExample() throws Exception {
this.mockMvc.perform(get("/users?offset=10&limit=10"))
.andExpect(status().isOk())
.andDo(document("users-list-example", requestParameters(
parameterWithName("offset").description("跳过的记录数量"),
parameterWithName("limit").description("返回的记录的最大数量")
)));
}

```

Show moreShow more icon

### 路径参数

REST API 中也可以使用路径参数，同样需要添加相应的文档。在 [清单 8\. 为路径参数添加文档](#清单-8-为路径参数添加文档) 中，pathParameters 方法用来添加路径参数的文档，同样使用 parameterWithName 方法来描述参数。所对应的代码片段的文件名是 path-parameters.adoc。

#### 清单 8\. 为路径参数添加文档

```
@Test
public void userGetPathExample() throws Exception {
User user = new User();
user.setUsername("alex");
user.setEmail("alex@example.com");
user = this.userRepository.save(user);
this.mockMvc.perform(get("/users/{user_id}", user.getId()))
.andExpect(status().isOk())
.andDo(document("user-get-path-example",
pathParameters(
parameterWithName("user_id").description("用户 ID")
)));
}

```

Show moreShow more icon

### HTTP 头

HTTP 请求和响应中的 HTTP 头也可以添加文档。在 [清单 9\. 为 HTTP 头添加文档](#清单-9-为-http-头添加文档) 中，requestHeaders 方法用来添加 HTTP 请求中头的文档，headerWithName 方法描述一个 HTTP 头的名称和含义。所对应的代码片段的文件名是 request-headers.adoc。对于 HTTP 响应中的头，可以使用 responseHeaders 方法来描述，同样使用 headerWithName 方法来描述单个 HTTP 头，对应的代码片段文件名为 response-headers.adoc。

#### 清单 9\. 为 HTTP 头添加文档

```
@Test
public void usersListHeaderExample() throws Exception {
this.mockMvc.perform(get("/users").header("X-Sample", "Hello World"))
.andExpect(status().isOk())
.andDo(document("users-list-header-example",
requestHeaders(
headerWithName("X-Sample").description("测试 HTTP 头")
)));
}

```

Show moreShow more icon

### 输入约束

在使用 REST API 时，领域模型的约束也是需要在文档中说明的。比如有的字段有长度限制，有的字段只允许 Email 地址或 URL。API 的使用者只有清楚这些约束，才能发出正确的 HTTP 请求。Spring REST Docs 提供了对使用 Bean Validation 2.0 规范描述的领域约束的支持。如果领域模型中使用了 Bean Validation 2.0 规范中的注解来描述约束，Spring REST Docs 可以读取这些约束并生成描述文本，并在文档中使用。生成的描述文本可以直接在字段的描述中使用，也可以作为字段的额外属性。

在示例应用中，领域模型 User 类的 email 字段的值必须为 Emai 地址。在 [清单 10\. 读取领域模型的约束并添加文档](#清单-10-读取领域模型的约束并添加文档) 中，创建一个查找 User 类中约束的 ConstraintDescriptions 对象，并通过 descriptionsForProperty 方法来得到每个字段的约束的描述信息。约束的描述信息被添加为字段描述的一部分。

#### 清单 10\. 读取领域模型的约束并添加文档

```
@Test
public void userCreateExample() throws Exception {
Map<String, Object> user = new HashMap<>();
user.put("username", "alex");
user.put("email", "alex@example.com");
ConstraintDescriptions userConstraints = new ConstraintDescriptions(User.class);
this.mockMvc
.perform(post("/users").contentType(MediaTypes.HAL_JSON).content(this.objectMapper.wr
iteValueAsString(user)))
.andExpect(status().isCreated())
.andDo(document("user-create-example", requestFields(
fieldWithPath("username")
.description(String.format("用户名（%s）",
userConstraints.descriptionsForProperty("username"))),
fieldWithPath("email")
.description(String.format("Email 地址（%s）",
userConstraints.descriptionsForProperty("email")))
)));
}

```

Show moreShow more icon

### 定制输出格式

前面提到了 Spring REST Docs 会自动生成代码片段。这些代码片段基于可以定制的 Mustache 模板。可以通过提供自定义模板文件的方式来覆盖默认的模板，从而改变代码片段文档的格式。以上一节中介绍的字段约束来说，更好的显示方式是在字段表格中增加专门的一列来显示约束。要实现这样的定制，需要做两个改变。

首先是把字段的约束添加为字段的额外属性。在 [清单 11\. 把字段的约束添加为额外的属性](#清单-11-把字段的约束添加为额外的属性) 中，使用 attributes 方法来添加额外的属性 constraints。属性值是字段约束的描述。

#### 清单 11\. 把字段的约束添加为额外的属性

```
this.mockMvc.perform(post("/users").contentType(MediaTypes.HAL_JSON).content(this.objectMa
pper.writeValueAsString(user)))
.andExpect(status().isCreated())
.andDo(document("user-create-example", requestFields(
fieldWithPath("username")
.description("用户名").attributes(
key("constraints").value(userConstraints.descriptionsForProperty("username"))
),
fieldWithPath("email")
.description("Email 地址")
.attributes(key("constraints").value(userConstraints.descriptionsForProperty("emai
l")))
)));

```

Show moreShow more icon

接着需要提供自定义的模板。模板文件需要保存在 src/test/resources/org/springframework/restdocs/templates/asciidoctor 目录中，模板文件的名称与对应生成的 Asciidoc 文件相同，只不过扩展名为 snippet。因为要定制的是请求字段的格式，因此模板文件的名称为 request-fields.snippet。模板文件的内容如 [清单 12\. 请求字段文件的自定义模板](#清单-12-请求字段文件的自定义模板) 所示，在表格中添加了一个新的列，其值为属性值 constraints。

#### 清单 12\. 请求字段文件的自定义模板

```
|===
|路径|类型|描述|约束
{{#fields}}
|{{#tableCellContent}}`{{path}}`{{/tableCellContent}}
|{{#tableCellContent}}`{{type}}`{{/tableCellContent}}
|{{#tableCellContent}}{{description}}{{/tableCellContent}}
|{{#tableCellContent}}{{constraints}}{{/tableCellContent}}
{{/fields}}
|===

```

Show moreShow more icon

## 定制 HTTP 请求和响应

在默认情况下，Spring REST Docs 根据实际发生的 HTTP 请求和响应来生成文档。但是在某些情况下，可能需要在生成文档之前对请求和响应内容进行预处理。Spring REST Docs 提供了相应的预处理功能，同时也提供了一些内置的预处理器，包括对内容进行重新格式化、删除 HTTP 头、替换内容、修改请求参数和修改 URI 等。

在 [清单 13\. 定制 HTTP 请求和响应](#清单-13-定制-http-请求和响应) 中，使用 preprocessRequest 方法添加了一个删除 HTTP 头 X-Sample 的请求预处理器，使用 preprocessResponse 方法添加了对内容进行重新格式化的响应预处理器。

#### 清单 13\. 定制 HTTP 请求和响应

```
@Test
public void indexExample() throws Exception {
this.mockMvc.perform(get("/"))
.andExpect(status().isOk())
.andDo(document("index-example",
preprocessRequest(
removeHeaders("X-Sample")
),
preprocessResponse(
prettyPrint()
),
links(
linkWithRel("users").description("指向用户资源的链接"),
linkWithRel("bookmarks").description("指向书签资源的链接"),
linkWithRel("profile").description("指向 ALPS 服务概要文件的链接")),
responseFields(
subsectionWithPath("_links").description("到其他资源的链接"))));
}

```

Show moreShow more icon

## 配置与自定义

如果 Spring REST Docs 的默认设置不能满足要求，可以进行不同的定制。在使用 MockMvc 时，可以定制文档中出现的 REST API 的 HTTP Scheme、主机名和端口。在前面介绍过，Spring REST Docs 默认创建 6 个代码片段，可以通过配置去掉一些不常用的代码片段。

如 [清单 14\. 配置 Spring REST Docs](#清单-14-配置-spring-rest-docs) 所示，配置工作是在创建 MockMvc 对象时完成的。通过方法 uris 来配置 REST API 的 URI，通过 snippets 来配置默认生成的代码片段。

#### 清单 14\. 配置 Spring REST Docs

```
@Before
public void setUp() {
this.mockMvc = MockMvcBuilders.webAppContextSetup(this.context)
.apply(documentationConfiguration(this.restDocumentation)
.uris().withScheme("https").withHost("api.mycompany.com").withPort(443)
.and()
.snippets().withDefaults(curlRequest(), requestBody(), httpResponse()))
.build();
}

```

Show moreShow more icon

## 使用 Asciidoctor 和 Markdown

Spring REST Docs 默认使用 Asciidoctor 作为文档的编写格式。Asciidoctor 的功能比较强大，可以使用宏。前面我们已经看到了宏 operation，可以用来包含某个测试用例所生成的代码片段文件，其中属性 snippets 可以选择要包含的代码片段的名称。使用宏 include 可以直接包含单个代码片段文件。

Markdown 的设计初衷是网页发表，对编写 API 文档这样的技术文档的支持比较弱，比如 Markdown 并没有原生的表格支持，也不支持包含其他文件。可以通过 [清单 15\. 使用 Markdown 作为文档格式](#清单-15-使用-markdown-作为文档格式) 中的方式把默认的文档格式改为 Markdown。

#### 清单 15\. 使用 Markdown 作为文档格式

```
this.mockMvc = MockMvcBuilders.webAppContextSetup(this.context)
.apply(documentationConfiguration(this.restDocumentation)
.snippets().withTemplateFormat(TemplateFormats.markdown()))
.build();

```

Show moreShow more icon

## 打包和发布

Spring REST Docs 生成的 API 文档是 HTML 格式的，可以直接发布到网上。也可以把生成的文档打包到 JAR 包中，作为发布内容的一部分。 [清单 16\. 把生成的文档打包到 JAR 文件中](#清单-16-把生成的文档打包到-jar-文件中) 给出了 build.gradle 中添加的内容，用来把生成的文档打包到 JAR 文件中。

#### 清单 16\. 把生成的文档打包到 JAR 文件中

```
jar {
dependsOn asciidoctor
from("${asciidoctor.outputDir}/html5") {
into 'static/docs'
}
}

```

Show moreShow more icon

## 小结

对于 REST API 的开发者来说，不管 API 的用户是内部团队，还是第三方，高质量的文档都是不可或缺的。长久以来，API 文档的正确性一直困扰着开发人员。Spring REST Docs 采用的手动编写内容与从单元测试自动生成代码片段相结合的方式，为解决 API 文档与代码的同步问题提供了一个很好的思路。本文对 Spring REST Docs 进行了详细的介绍，包括基本的配置和对 HTTP 请求和响应的不同部分添加文档。

## 参考资源

- 参考 [Spring REST Docs 官方指南](https://docs.spring.io/spring-restdocs/docs/2.0.0.RELEASE/reference/html5/) ，了解更多内容。
- 了解如何使用 [MockMvc](https://spring.io/guides/gs/testing-web/) 来测试 Spring MVC 项目。
- 了解 [Asciidoc](https://asciidoctor.org/docs/asciidoc-writers-guide/) 文档格式。