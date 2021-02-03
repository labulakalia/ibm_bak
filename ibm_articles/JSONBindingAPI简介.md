# JSON Binding API 简介
了解 JSON-B 的默认特性和自定义注解、运行时配置等

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-javaee8-json-binding-1/)

Alex Theedom

发布: 2017-12-07

* * *

**关于本系列**：Java EE 支持 XML 已久，但显然遗漏了对 JSON 数据的内置支持。Java EE 8 改变了这一情况，向核心 Java 企业平台带来了强大的 JSON 绑定特性。了解 JSON-B，以及如何结合它与 JSON Processing API 及其他技术来处理 Java 企业应用程序中的 JSON 文档。

作为 Java™ EE 8 中的一个重大特性，Java API for JSON Binding (JSON-B) 1.0 ( [JSR 367](https://www.jcp.org/en/jsr/detail?id=367)) 加强了 Java 对 JavaScript Object Notation (JSON) 的支持。

JSON 已成为在联网与未联网的系统间传递信息的标准。它在云和 Web 服务应用程序中大受欢迎，已成为基于微服务的系统中的 RESTful Web 服务首选的数据交换格式。

JSON-B 完善了 Java EE 的 API 套件和管理 JSON 的特性，提供了一种在 Java 对象与 JSON 文档间来回转换的简单、统一标准。它很容易使用，这一定程度上得益于最少的程序，以及合理且直观的默认配置。

JSON Binding 还可以与 JSON Processing 紧密结合，后者是 Java EE 的另一个管理企业系统中的 JSON 的 API。JSON-P 是一个较低级的 API，为 JSON 处理和转换提供了两种模型（流和对象）。这两个 API 相结合，为在 Java 企业系统中使用 JSON 提供了一个牢固的框架。

## JSON-B 的默认和自定义绑定

JSON-B 以开箱即用的方式提供了适合日常编程场景的默认序列化和反序列化映射。对于较简单的自定义，该 API 提供了适配器及一些自定义序列化器和反序列化器。对于较高级的场景，可以使用类中的注解或运行时配置构建器来覆盖 JSON-B 的默认设置。

总而言之，JSON-B 整理了已在企业开发人员中广泛应用的行业实践和方法。它使用注解通过映射语义来标记类和字段，提供了处理复杂数据结构时常常需要的可扩展性。

最重要的是，这个 API 支持以一种直观且容易的方式在 Java 类与 JSON 文档之间进行绑定。甚至没有 JSON 经验的开发人员，应该也能够轻松上手使用 Java API for JSON Binding。对 JSON 反/序列化库（比如 [GSON](https://github.com/google/gson) 、 [Boon](http://richardhightower.github.io/site/Boon/Welcome.html) 和 [Jackson](https://github.com/FasterXML/jackson) ）经验丰富的开发人员将发现 JSON-B 很熟悉且用起来很舒适。

## 两个接口：Jsonb 和 JsonBuilder

从功能上讲，JSON-B API 包含两个入口点接口： [javax.json.bind.Jsonb](https://static.javadoc.io/javax.json.bind/javax.json.bind-api/1.0/javax/json/bind/Jsonb.html) 和 [javax.json.bind.JsonBuilder](https://static.javadoc.io/javax.json.bind/javax.json.bind-api/1.0/javax/json/bind/JsonbBuilder.html) 。 `Jsonb` 接口分别通过方法 `toJson()` 和 `fromJson()` 实现序列化和反序列化。 `JsonbBuilder` 接口基于一组可选的操作而构造 `Jsonb` 实例，并为客户端提供对实例的访问权。

### 往返示例

要了解一个新 API，没有比亲自试用更好的方法了。在这个示例中，我们将使用一个简单类来构建一个基本的示例应用程序。在清单 1 中，我们首先将 `Book` 类序列化为一个 JSON 文档，然后将它反序列化为 `Book` 对象。

##### 清单 1\. 一个最简单的 Book 类

```
public class Book {
private String id;
private String title;
private String author;
// Plumbing code removed for brevity
}

```

Show moreShow more icon

要序列化一个 `Book` 对象，首先需要创建一个 `Jsonb` 实例，我们使用 `JsonbBuilder` API 的静态 `create()` 工厂方法完成此任务。拥有 `Jsonb` 实例后，我们调用一个重载的 `toJson()` 方法并将要序列化的对象传递给它，如清单 2 所示。

##### 清单 2\. 序列化一个 Book 实例

```
Book book = new Book("ABCD-1234", "Fun with Java", "Alex Theedom");
Jsonb jsonb = JsonbBuilder.create();
String json = jsonb.toJson(book);

```

Show moreShow more icon

`toJson()` 方法返回序列化的 `Book` 对象，如清单 3 所示。

##### 清单 3\. 返回序列化的 Book 实例

```
{
"author": "Alex Theedom",
"id": "ABCD-1234",
"title": "Fun with Java"
}

```

Show moreShow more icon

请注意清单 3 中的属性采用字母顺序。 这是默认的属性顺序，而且可以自定义。本文后面将介绍自定义 JSON-B 的一些细节。

使用 JSON-B 执行反序列化同样很简单。给定一个 `Jsonb` 实例，我们调用一个重载的 `fromJson()` 方法并将要序列化的 JSON 文档和对象类型传递给它，如清单 4 所示。

##### 清单 4\. 反序列化一个 JSON 字符串

```
Book book = jsonb.fromJson(json, Book.class);

```

Show moreShow more icon

**并发性与 JSON-B**：尽可能重用 `JsonbBuilder` 和 `Jsonb` 实例是一种最佳实践，如清单 4 所示。根据 [Java 文档](https://static.javadoc.io/javax.json.bind/javax.json.bind-api/1.0/index.html?java.json.bind-summary.html) 中针对 JSON-B API 的建议，在典型的用例中，每个应用程序只需要一个 `Jsonb` 实例。JSON-B 的所有方法都可由多个并发线程安全地使用。

尽管 JSON-B 规范没有强制要求保留往返相等性，但它推荐这么做。尽管无法保证，但是在大部分情况下，将序列化的 JSON 输出提供给反序列化的 `fromJson()` 方法，将得到一个与原始对象等效的对象。

## JSON-B 的默认设置

为了方便， `Jsonb` 实例预先配置了反映公认标准的默认设置。下面简要介绍 JSON-B 的默认策略和映射。本系列后面的文章将更深入地探索其中一些方面。

- **属性命名策略** 确定类的字段将如何在序列化 JSON 中表示。默认情况下，JSON 属性名采用与它序列化的类中的字段完全相同的名称。
- **属性可视性策略** 确定如何访问字段，以及为字段和 bean 方法上的访问修饰符赋予的重要性。在序列化和反序列化期间，JSON-B 的默认行为都是仅接受公共字段、公共访问器和更改器方法。如果不存在公共方法或字段，则会忽略该属性。
- **属性顺序策略** 指定属性在输出 JSON 中出现的顺序。JSON-B 的默认顺序是字母顺序。
- **二进制数据策略** 编码二进制数据（默认使用一个字节数组）。
- **Null 值策略** 指定忽略包含 null 值的对象字段和保留集合中的 null 值。JSON 属性中的 null 值将相应对象字段值设置为 _null_ 。

了解 3 个额外的默认设置很重要：

- **JSON-B 遵守 [I-JSON 概要文件](https://tools.ietf.org/html/rfc7493)** ，但有 3 处例外：JSON-B 不会将顶级 JSON 的序列化限制到非对象或数组文本；不会序列化采用 `base64url` 编码的二进制数据，也不会对基于时间的属性执行额外的限制。也可以通过一个自定义选项开启与 I-JSON 的严格一致性。
- **JSON-B 的数据格式和区域设置** 将取决于您的机器设置。
- **JSON 数据** 的表示不使用换行符和缩进，而且使用 UTF-8 进行编码。

## 使用 JSON-B 处理类型

Java API for JSON Binding 在序列化和反序列化上 [兼容 RFC 7159](https://tools.ietf.org/html/rfc7159) ，而且支持所有 Java 类型和原语。本系列将探讨和使用以下许多原则和技术：

- **基本类型** ：这些类型包括所有原语，以及它们关联的包装器类、 `Number` 和 `String` ：

    - 在序列化期间，会在实例字段上调用 `toString()` 方法来生成一个 `String` 值。
    - `Number` 类型的字段是个例外，因为 `Number` 使用 `doubleValue()` 方法生成原语。
    - 原语值会保持不变。
    - 在反序列化时，会调用合适的 `parseXXX()` 方法来转换为需要的基本类型。
- **Java SE 类型** ：这些类型指的是一系列 `Optional` 类（比如 `OptionalInt` ）、 `BigInteger` 和 `BigDecimal` 类，以及 `URI` 和 `URL` 类。

    - 要将可选值转换为 `String` ，可以获取它们的内部类型并按照对象的约定转换为 `String` 或原语。（一个示例是在 `Integer` 实例上调用 `intValue()` 。）
    - 空可选值视为 null 值，依据处理 null 的自定义设置而处理。
    - 反序列化过程利用该类型的接受 `String` 参数的构造方法。一个示例是 `new BigDecimal("1234")` 。
- **日期类型** ：JSON-B 支持所有标准 Java 日期类型，包括 `Calendar` 、 `TimeZone` 和 `java.time.*` 包。序列化期间使用的格式包括 `ISO_DATE` 、 `ISO_DATE_TIME` 和 `ISO_INSTANT/LOCAL_*/ZONED_*/OFFSET_*` 。API 规范（ [3.5 节](https://jcp.org/aboutJava/communityprocess/final/jsr374/index.html) ）对用于每种日期类型的格式进行了详细分类。
- **数组和集合类型** ：可以序列化和反序列化所有支持的 Java 类型和原语的单维和多维数组。Null 值元素在结果 JSON 或 Java 数组中表示为 null 值。
- **JSON 处理类型** ：所有 JSON 处理类型都支持序列化和反序列化。
- **Java 类** ：要实现可序列化，Java 类需要有一个无参数的 `public` 或 `protected` 构造方法，或者显式指定要用于自定义对象构造的方法。要注意的一些规则：

    - 序列化不需要约束：JSON-B 能够反序列化 `public` 、 `protected` 和 `protected static` 类。
    - 不支持匿名类：序列化匿名类会生成 JSON 对象。
    - 除了已经介绍的与 Java 基本/Java SE 中 `Date/Time` 和 `Collection/Map` 类型相关的类型，反序列化期间不支持接口。运行时类型用于序列化。

## 默认映射

接下来，我们将看一个示例，其中包含目前介绍的许多默认设置和类型。清单 5 的 `Magazine` 类提供了一些 Java 类型、原语和集合。请注意， `internalAuditCode` 字段只有一个 setter 方法。

##### 清单 5\. Magazine 类

```
public class Magazine {
private String id;
private String title;
private Author author;
private Float price;
private int pages;
private boolean inPrint;
private Binding binding;
private List<String> languages;
private URL website;
private String internalAuditCode; // Only has setter method
private LocalDate published;
private String alternativeTitle;
// Plumbing code removed for brevity
}

```

Show moreShow more icon

`Magazine` 对象将在清单 6 中构造，在 `language` 数组中包含 null 字段，而且 `alternativeTitle` 字段也设置为 null。

##### 清单 6\. 构造 Magazine 对象

```
Magazine magazine = new Magazine();
magazine.setId("ABCD-1234");
magazine.setTitle("Fun with Java");
magazine.setAuthor(new Author("Alex", "Theedom"));
magazine.setPrice(45.00f);
magazine.setPages(300);
magazine.setInPrint(true);
magazine.setBinding(Binding.SOFT_BACK);
magazine.setLanguages(
    Arrays.asList("French", "English", "Spanish", null));
magazine.setWebsite(new URL("https://www.readlearncode.com"));
magazine.setInternalAuditCode("IN-675X-NF09");
magazine.setPublished(LocalDate.parse("01/01/2018",
    DateTimeFormatter.ofPattern("MM/dd/yyyy")));
magazine.setAlternativeTitle("IN-675X-NF09");
String json = JsonbBuilder.create().toJson(magazine);

```

Show moreShow more icon

清单 7 给出了序列化这个对象的代码。序列化时，数组中的 null 值会保留，而 `alternativeTitle` 字段中的 null 值会导致该字段被排除。

另请注意， `internalAuditCode` 值未包含在 JSON 输出中。不能包含它，是因为没有在获取该值时要使用的 getter 方法或公共字段。

##### 清单 7\. Magazine 的 JSON 表示

```
{
"author": {
    "firstName": "Alex",
    "lastName": "Theedom"
},
"binding": "SOFT_BACK",
"id": "ABCD-1234",
"inPrint": true,
"languages": [
    "French",
    "English",
    "Spanish",
    null
],
"pages": 300,
"price": 45.0,
"published": "2018-01-01",
"title": "Fun with Java",
"website": "https://www.readlearncode.com"
}

```

Show moreShow more icon

## 自定义映射

对于大部分简单场景，JSON-B 的默认映射行为应该就够了。该 API 还提供了各种不同途径来针对更复杂的场景和需求而自定义映射行为。

有两个自定义模型，您可以结合使用或单独使用它们：注解模型和运行时模型。注解模型使用内置的注解来标记您希望自定义其行为的字段。运行时模型构建一个在 `Jsonb` 实例上设置的配置策略。注解和运行时配置可同时应用于序列化和反序列化。

### 通过注解自定义 JSON-B

通过使用注解，可以标记您希望自定义其关联的映射行为的字段、JavaBean 属性、类型或包。以下是一个示例：

##### 清单 8\. 注解一个字段

```
@JsonbProperty("writer")
private Author author;

```

Show moreShow more icon

`author` 字段上的注解将 JSON 输出中的属性名称从 `author` 更改为 `writer` 。

### 通过运行时配置自定义 JSON-B

运行时模型构建一个 `JsonbConfig` 实例，它将该实例传递给 `JsonBuilder` 的 `create()` 方法，如清单 9 所示。

##### 清单 9\. 构建一个 Jsonb 配置实例

```
JsonbConfig jsonbConfig = new JsonbConfig()
       .withPropertyOrderStrategy(PropertyOrderStrategy.REVERSE)
       .withNullValues(true);

Jsonb jsonb = JsonbBuilder.create(jsonbConfig);

```

Show moreShow more icon

在这个示例中，该配置将属性顺序设置为 `reverse` 并保留所有 null 值。

## 基本自定义

尽管可以通过开箱即用的配置满足大部分简单用例，但默认映射只能为您做这么多。对于需要更多地控制序列化和反序列化流程的情况，JSON-B 提供了一系列自定义选项。

### 自定义属性名称和顺序

JSON-B 提供了多种与属性相关的自定义，您应该了解它们。

如果字段使用 `@JsonbTransient` 进行了注解，则会排除一个属性，属性的名称可使用 `@JsonProperty` 注解更改。要设置命名策略，可以将一个 `PropertyNamingStrategy` 常量传递到 `JsonbConfig` 实例的 `withPropertyNamingStrategy()` 方法。

属性在 JSON 输出中出现的顺序在类级别上设置，可使用 `@JsonbPropertyOrder()` 注解或 `withPropertyOrderStrategy()` 方法完成。传递给该注解和方法的 `PropertyOrderStrategy` 常量决定了策略。

### 配置忽略属性和可视性

可以使用 `@JsonbVisibility()` 注解或 `withPropertyVisibilityStrategy()` 方法来指定给定属性的可视性。

### Null 处理

可通过 3 种方式更改处理 null 的默认行为：

1. 使用 `@JsonbNillable` 注解来注解包或类。
2. 使用 `@JsonbProperty` 来注解相关的字段或 JavaBean 属性，将 `nillable` 参数设置为 _true_ 。
3. 将 _true_ 或 _false_ 传递给 `withNullValues()` 方法。

清单 10 演示了第二种和第三种方法：使用 `@JsonbNillable` 来注解类和将 `nillable` 设置为 _true_ 。可以使用类级注解 `@JsonbNillable` 来配置所有字段的 null 处理，或者使用字段级注解来实现字段级的更精细控制：

##### 清单 10\. 配置 null 处理

```
public class Magazine {
    @JsonbProperty(value = "writer", nillable = true)
    private Author author;
}

@JsonbNillable
    public class Magazine {
    private Author author;
}

```

Show moreShow more icon

### 自定义创建器

回想一下，JSON-B 要求所有类拥有一个无参数的公共构造方法，该方法用于构造类实例。如果这还不够，可以使用一个自定义构造方法或静态工厂方法。对于每种选项，可以使用注解 `@JsonCreator` 。

### 自定义日期和数字格式

JSON-B 指定使用 `@JsonbDateFormat()` 注解来设置日期格式并传递要使用的地区和日期模式。该注解可用在从包到字段的任何地方。另一个选择是使用 `withDateFormat()` 方法，如清单 11 所示。

##### 清单 11\. 日期格式

```
new JsonbConfig()
     .withDateFormat("MM/dd/yyyy", Locale.ENGLISH);

```

Show moreShow more icon

### 处理二进制数据

JSON-B 提供了 3 种处理二进制数据的策略：

- BYTE
- BYTE\_64
- BYTE\_64\_URL

这 3 个静态常量之一会从 `BinaryDataStrategy` 传递给 `withBinaryDataStrategy()` 方法。清单 12 给出了一个示例。

##### 清单 12\. 二进制数据策略

```
new JsonbConfig()
     .withBinaryDataStrategy(BinaryDataStrategy.BASE_64_URL);

```

Show moreShow more icon

### 配置 I-JSON

在大部分情况下，默认设置都与 I-JSON 概要文件一致，除了前面的”默认设置”节中提到的例外。对于需要更严格的一致性的场景，可以使用 `withStrictIJSON(true)` 方法。

## 高级自定义

在一些情况下，注解和运行时配置都爱莫能助。例如，常常无法注解第三方类。没有默认构造方法的类也是一大麻烦。JSON-B 为这些场景提供了两个解决方案：适配器和序列化器/反序列化器。

### 通过适配器自定义 JSON-B

适配器允许创建自定义 Java 对象和序列化 JSON 代码。适配器类必须实现 `JsonbAdapter` 接口并提供它的两个方法的代码： `adaptToJson()` 和 `adaptFromJson()` 。

在清单 13 中，实现 `adaptToJson()` 方法的代码将清单 14 中所示的 `Booklet` 对象转换为 JSON 字符串：

##### 清单 13\. 实现 adaptToJson() 方法

```
public JsonObject adaptToJson(Booklet booklet) throws Exception {
return Json.createObjectBuilder()
        .add("title", booklet.getTitle())
        .add("firstName", booklet.getAuthor().getFirstName())
        .add("lastName", booklet.getAuthor().getLastName())
        .build();
}

```

Show moreShow more icon

请注意，在清单 13 中，我们使用了 JSON-P `JsonObjectBuilder` 。这是 `Booklet` 对象：

##### 清单 14\. Booklet

```
public class Booklet {
       private String title;
       private Author author;
    // Plumbing code removed for brevity
}

```

Show moreShow more icon

可以看到， `adaptToJson()` 方法将 `Author` 对象”扁平化”为两个属性： `firstName` 和 `secondName` 。这只是许多使用适配器自定义序列化的示例之一。

### 注册一个 JSON-B 适配器

可以使用 `JsonbConfig` 实例来注册 JSON-B 适配器，如下所示：

```
new JsonbConfig().withAdapters(new bookletAdapter())

```

Show moreShow more icon

或者通过在给定类上使用以下注解：

```
@JsonbTypeAdapter(BookletAdapter.class)

```

Show moreShow more icon

清单 15 显示了此转换生成的 JSON 文档。

##### 清单 15\. Booklet 类的序列化

```
{
"title": "Fun with Java",
"firstName": "Alex",
"lastName": "Theedom"
}

```

Show moreShow more icon

`adaptFromJson()` 方法未给出；但是 GitHub 存储库中有一个与本文相关的代码示例。

### 通过序列化器和反序列化器自定义 JSON-B

JSON-B 序列化器和反序列化器是可用的最低级自定义，让您能够访问 JSON Processing API 中提供的解析器和生成器。

要使用此选项，必须实现 `JsonbDeserializer` 和 `JsonbSerializer` ，而且覆盖合适的 `serialize()` 或 `deserialize()` 方法。然后可以编写自定义代码来执行必要的工作，将该代码放在该方法中。

序列化器和反序列化器通过一个包含合适方法的 `JsonbConfig` 实例来注册： `withDeserializers()` 或 `withSerializers()` 。请参阅本文的 GitHub 存储库来查看创建和使用 JSON-B 序列化器和反序列化器的示例。

## 结束语

为让 JSON-B 准备好在 Java EE 8 中发布，专家组付出了艰苦的努力，我们应该祝贺他们。这个 API 完善了 Java EE 的 JSON API 套件。在本文中，我们粗略介绍了 JSON Binding API 的广泛功能。接下来的 3 篇文章将更深入地介绍 Java EE 的最受欢迎的补充组件的重要方面。

## 测试您的知识

1. 在结合使用时，以下哪些注解可以自定义 JSON 输出来保留 null 值，将属性按字母逆序排序，以及将一个属性重命名为”cost”？

    1. `@JsonbNillable`
    2. `@JsonbNullable`
    3. `@JsonbPropertyOrder(PropertyOrderStrategy.REVERSE)`
    4. `@JsonbPropertyOrder(PropertyOrderStrategy.RESERVED)`
    5. `@JsonbPropertyName("cost")`
    6. `@JsonbProperty("cost")`
2. 您应该调用以下 `JsonbConfig` 构建器方法的哪个组合，才能构造一个满足以下条件的自定义配置：使用严格的 I-JSON 概要文件；将日期格式指定为月/日/年和英语地区；以及以一种美化的格式输出 JSON？

    1. `withNonStrictIJSON(false)`
    2. `withStrictIJSON(true)`
    3. `withDateFormat(Locale.ENGLISH, "MM/dd/yyyy")`
    4. `withDateFormat("MM/dd/yyyy", Locale.ENGLISH)`
    5. `withFormatting(true)`
    6. `withPrettyFormatting(true)`
3. 假设适配器的名称为 `OrderAdapter.class` ，可使用以下哪种方式来指定要使用的适配器？

    1. `new JsonbConfig().withAdapters(new orderAdapter())`
    2. `new JsonbConfig().withTypeAdapters(new orderAdapter())`
    3. `@JsonbTypeAdapter(OrderAdapter.class)`
    4. `@JsonbAdapter(OrderAdapter.class)`
    5. 上述选项都不是
4. 序列化和反序列化的默认可视性配置是什么？

    1. 仅公共访问器和更改器方法可见
    2. 仅公共访问器和更改器方法及字段可见
    3. 仅公共字段可见
    4. 公共和私有字段可见
    5. 公共和受保护的访问器和更改器方法及字段可见
5. 可使用以下哪个选项来自定义从 JSON 文档创建对象的过程？

    1. 一个注解的自定义构造方法： `@JsonCreator`
    2. 一个注解的静态工厂方法： `@JsonCreator`
    3. 一个 JSON B 适配器
    4. 一个 JSON B 序列化器
    5. 无法实现

## 核对您的答案

1. 在结合使用时，以下哪些注解可以自定义 JSON 输出来保留 null 值，将属性按字母逆序排序，以及将一个属性重命名为”cost”？

    1. **`@JsonbNillable`**
    2. `@JsonbNullable`
    3. **`@JsonbPropertyOrder(PropertyOrderStrategy.REVERSE)`**
    4. `@JsonbPropertyOrder(PropertyOrderStrategy.RESERVED)`
    5. **`@JsonbPropertyName("cost")`**
    6. `@JsonbProperty("cost")`
2. 您应该调用以下 `JsonbConfig` 构建器方法的哪个组合，才能构造一个满足以下条件的自定义配置：使用严格的 I-JSON 概要文件；将日期格式指定为月/日/年和英语地区；以及以一种美化的格式输出 JSON？

    1. `withNonStrictIJSON(false)`
    2. **`withStrictIJSON(true)`**
    3. `withDateFormat(Locale.ENGLISH, "MM/dd/yyyy")`
    4. **`withDateFormat("MM/dd/yyyy", Locale.ENGLISH)`**
    5. **`withFormatting(true)`**
    6. `withPrettyFormatting(true)`
3. 假设适配器的名称为 `OrderAdapter.class` ，可使用以下哪种方式来指定要使用的适配器？

    1. **`new JsonbConfig().withAdapters(new orderAdapter())`**
    2. `new JsonbConfig().withTypeAdapters(new orderAdapter())`
    3. **`@JsonbTypeAdapter(OrderAdapter.class)`**
    4. `@JsonbAdapter(OrderAdapter.class)`
    5. 上述选项都不是
4. 序列化和反序列化的默认可视性配置是什么？

    1. 仅公共访问器和更改器方法可见
    2. **仅公共访问器和更改器方法及字段可见**
    3. 仅公共字段可见
    4. 公共和私有字段可见
    5. 公共和受保护的访问器和更改器方法及字段可见
5. 可使用以下哪个选项来自定义从 JSON 文档创建对象的过程？

    1. **一个注解的自定义构造方法： `@JsonCreator`**
    2. **一个注解的静态工厂方法： `@JsonCreator`**
    3. **一个 JSON B 适配器**
    4. **一个 JSON B 序列化器**
    5. 无法实现

本文翻译自： [The JSON Binding API in a nutshell](https://developer.ibm.com/articles/j-javaee8-json-binding-1/)（2017-11-10）