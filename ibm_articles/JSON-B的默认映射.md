# JSON-B 的默认映射
通过 JSON Binding API 执行日常的序列化和反序列化

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-javaee8-json-binding-3/)

Alex Theedom

更新: 2018-01-24 \| 发布: 2017-11-18

* * *

[本系列的第一篇文章](https://www.ibm.com/developerworks/cn/java/j-javaee8-json-binding-1/index.html?ca=drs-) 概述了 Java API for JSON Binding 1.0 (JSON-B)，包括自定义和默认绑定选项。本文将介绍 JSON-B 的默认映射，以及对包括特殊数据类型在内的大部分 Java 类型执行序列化和反序列化的自动行为。

**关于本系列**：Java EE 长期以来一直支持 XML，但它显然遗漏了对 JSON 数据的内置支持。Java EE 8 改变了这一状况，给核心 Java 企业平台带来了强大的 JSON 绑定特性。开始使用 JSON-B，了解如何结合使用它与 JSON Processing API 及其他技术来处理 Java 企业应用程序中的 JSON 文档。

## 默认映射和内部原理

JSON-B API 采用了一种熟悉的方法来映射字段和属性，如果您之前使用过 JSON Binder，应该很容易上手。JSON-B 的映射应用了一组规则，而不需要注解或自定义配置。这些映射是开箱即用的，开发人员无需了解 JSON 文档的结构（甚至是 JSON 数据交换格式本身）就能上手使用它们。

JSON-B API 提供了两个入口点接口： `javax.json.bind.Jsonb` 和 `javax.json.bind.JsonBuilder` 。 `Jsonb` 接口通过重写方法 `toJson()` 和 `fromJson()` 来提供序列化和反序列化功能。 `JsonbBuilder` 接口为 `Jsonb` 实例提供客户端访问点，使用一组可选的配置来构建这些实例。

同样值得重点注意的是，JSON-B 高度依赖于 JSON Processing API 所提供的功能。借助 JSON-P 的流模型， `javax.json.stream.JsonGenerator` 提供了将 JSON 数据写到输出源（通过 `toJson()` 方法）的功能，而 `javax.json.stream.JsonParser` 提供了对 JSON 数据的正向只读访问能力（通过 `fromJson()` 方法）。

## 安装 Yasson

要上手使用 JSON-B，需要获得它的参考实现 [Eclipse Yasson](https://projects.eclipse.org/projects/ee4j.yasson) ，该实现可从 [Maven 中央存储库](https://repo.eclipse.org/content/repositories/yasson-releases/org/eclipse/yasson/) 获得。如下所示，您需要将 JSON Processing API 添加到 JSON-B 的 Maven POM 依赖项列表中。

##### 清单 1\. Maven 坐标

```
<dependency>
<groupId>javax.json</groupId>
<artifactId>javax.json-api</artifactId>
<version>1.1</version>
</dependency>

<!-- JSON-P 1.1 RI -->
<dependency>
<groupId>org.glassfish</groupId>
<artifactId>javax.json</artifactId>
<version>${javax.json.version}</version>
</dependency>

<!-- JSON-B 1.0 API -->
<dependency>
     <groupId>javax.json.bind</groupId>
     <artifactId>javax.json.bind-api</artifactId>
     <version>1.0</version>
</dependency>

<!-- JSON-B 1.0 RI -->
<dependency>
<groupId>org.eclipse</groupId>
<artifactId>yasson</artifactId>
<version>1.0</version>
</dependency>

```

Show moreShow more icon

## 最简单的示例

我们先来看一个简单示例，该示例对清单 2 中所示的 `Book` 类实例进行来回转换。

##### 清单 2\. Book 对象

```
public class Book {
public String title;
}

```

Show moreShow more icon

要开始序列化或反序列化，您需要一个 `Jsonb` 实例。可通过调用 `JsonBuilder` 接口的静态工厂方法 `create()` 来创建此实例。借助这个实例，可通过选择适当重载的 `toJson()` 或 `fromJson()` 方法来执行序列化和反序列化操作。

在清单 3 中，我调用了最简单的 `toJson()` 方法，并将一个 `Book` 对象传递给它。

##### 清单 3\. 对 Book 实例执行序列化

```
Book book = new Book("Fun with Java"); // Sets the title field
String bookJson = JsonbBuilder.create().toJson(book);

```

Show moreShow more icon

此方法的值是一个 `String` ，其中包含传递给 `toJson()` 的对象的 JSON 数据表示，如清单 4 所示。

##### 清单 4\. Book 实例的 JSON 表示

```
{
"title": "Fun with Java"
}

```

Show moreShow more icon

现在让我们将注意力转移到反序列化操作上。反序列化像序列化一样简单，而且也需要一个 `Jsonb` 实例。在清单 5 中，我将调用最简单的 `fromJson()` ，并将要反序列化的 JSON 数据及其目标类型传递给它。

##### 清单 5\. 对 Book 的 JSON 表示执行反序列化

```
String json = "{\"title\":\"Fun with Java\"}";
Book book = JsonbBuilder.create().fromJson(json, Book.class);

```

Show moreShow more icon

在这些示例中，我使用了 `toJson()` 和 `fromJson()` ，它们是 `Jsonb` 上提供的两个最简单的重载方法。接下来您将了解来自此接口的其他一些更复杂的方法。

**往返等效性**

请注意，JSON-B 规范没有保证往返等效性，但在大部分情况下，您会发现这一特性得到了保留。在整个系列中，我将指出这条一般规则的一些例外情况，以及如何解决它们。

## 重载方法

`Jsonb` 接口提供了重载的 `toJson()` 和 `fromJson()` 方法来执行 JSON-B 的序列化和反序列化功能。表 1 概述了这 12 个方法及其签名。

##### 表 1\. Jsonb 接口方法

修饰符和类型方法 TfromJson(InputStream stream, Class type) TfromJson(InputStream stream, Type runtimeType) TfromJson(Reader reader, Class type) TfromJson(Reader reader, Type runtimeType) TfromJson(String string, Class type) TfromJson(String string, Type runtimeType)StringtoJson(Object object)voidtoJson(Object object, OutputStream stream)voidtoJson(Object object, Writer writer)StringtoJson(Object object, Type runtimeType)voidtoJson(Object object, Type runtimeType, OutputStream stream)voidtoJson(Object object, Type runtimeType, Writer writer)

如您所见，方法的范围非常广泛，足以支持大部分使用场景。一个最有用的解决方案是连接到 `InputStream` 或 `OutputStream` 实例的能力。根据执行的是序列化还是反序列化操作，要么接受 `Stream` 实例并输出一个 Java™ 对象，要么写入到 `Stream` 。

`InputStream` 和 `OutputStream` 相结合，为 JSON 数据的流处理带来了许多富有创意的可能性。举例而言，考虑一个包含 JSON 数据并执行往返操作的平面文件。在清单 6 中，调用了 `toJson(Object object, OutputStream stream)` 方法对一个 `Book` 类实例执行序列化并输出到文本文件： `book.json` 。

##### 清单 6\. 将 JSON 序列化结果发送到文本文件 OutputStream

```
Jsonb jsonb = JsonbBuilder.create();
jsonb.toJson(book, new FileOutputStream("book.json"));

```

Show moreShow more icon

在清单 7 中，反序列化操作是通过调用 `fromJson(InputStream stream, Class<T> type)` 方法来执行的。

##### 清单 7\. 对存储在文本文件中的 JSON 数据执行反序列化

```
Book book = jsonb.fromJson(new FileInputStream("book.json"), Book.class);

```

Show moreShow more icon

**性能和重用**

`Jsonb` 接口的 [Java 文档](https://static.javadoc.io/javax.json.bind/javax.json.bind-api/1.0/index.html?overview-summary.html) 包含这些方法的操作的完整细节，以及其他一些示例。此外，该文档还建议，要达到最佳的使用效果，应该重用 `JsonbBuilder` 和 `Jsonb` 的实例，而且对于典型的用例，每个应用程序只需一个 `Jsonb` 实例。这是可能做到的，因为它的方法是线程安全的。

## 基本 Java 类型的默认处理方式

JSON-B 对基本 Java 类型的序列化和反序列化应用了一些简单的规则。Java 中的基本类型包括 `String` 、 `Boolean` ，以及 `Character` 和 `AtomicInteger` 等所有 `Number` 类型及其相应的原语类型。

对于原语及其等效包装器，JSON-B 使用 `toString()` 方法生成一个被转换为 JSON `Number` 的 `String` 。对于一些 `Number` 类型（比如 `java.util.concurrent.atomic.LongAdder` ），调用 `toString()` 方法不会生成合适的等效 `String` 。在这些情况下，可以调用 `doubleValue()` 方法生成一个双精度原语。

### Number 类型的异常处理

**编码**

默认情况下，UTF-8 是 `toJson()` 的序列化操作的编码。对于使用 `fromJson()` 的反序列化，会从 JSON 数据中自动检测编码，而且该编码应该是 RFC 7159 中定义的有效字符编码。也可以根据需要对编码进行自定义。

JSON-B 规范要求为除包装器类型外的所有 `Number` 类型实例化一个 `BigDecimal` 。在反序列化期间，反序列化操作创建了一个 `BigDecimal` 实例，方法是将 JSON 数字作为 `String` 传递给该实例的构造函数，如下所示： `new BigDecimal("10")` 。如果目标类型不是 `BigDecimal` ，而是 `AtomicInteger` 等其他许多 `Number` 类型之一，这会成为一个问题。

对于大部分基本类型，反序列化操作会自动调用合适的 `parse$Type` 方法 — 就像在 `Integer.parseInt("10")` 或 `Boolean.parseBoolean("true")` 中一样。但是，包装器类型之外的 Number 类型可能没有 `parse$Type` 方法，所以在反序列化器尝试调用 `parse$Type` 方法时，将会抛出一个异常。

要了解这一过程的工作原理，可以先来看看清单 8 中简化的 `Book` 类。

##### 清单 8\. 包含 AtomicInteger 的 SimpleBook 类

```
public class SimpleBook {
private AtomicInteger bookVersion;
// Getters and setters omitted for brevity
}

```

Show moreShow more icon

JSON 数据仅包含一个 `Integer` 类型的数字，如清单 9 所示。

##### 清单 9\. SimpleBook 的 JSON 表示

```
{
    "bookVersion":10
}

```

Show moreShow more icon

当 JSON-B 尝试将 `SimpleBook` 的 JSON 表示转换回 `SimpleBook` 类时，反序列化操作会根据它的名称 `bookVersion` 来识别目标字段。它确定 `bookVersion` 是一种不属于包装器类型分组的 `Number` 类型，并创建一个 `BigDecimal` 实例。

`BigDecimal` 类型与 `AtomicInteger` 类型不兼容。因此，反序列化操作会抛出一个 `IllegalArgumentException` 和消息”`argument type mismatch` ”。清单 10 中的代码会引起这种异常：

##### 清单 10\. 抛出 IllegalArgumentException

```
Jsonb jsonb = JsonbBuilder.create();
jsonb.fromJson(
       jsonb.toJson(new SimpleBook(new AtomicInteger(10))),
       SimpleBook.class);

```

Show moreShow more icon

总之，为原语和它们的等效包装器保留了往返等效性，为 `BigDecimal` 和 `BigInteger` 也保留了往返等效性。其他 `Number` 类型没有直接支持往返等效性。可以使用我将在 [第 2 部分](https://www.ibm.com/developerworks/cn/java/j-javaee8-json-binding-2/index.html?ca=drs-) 中介绍的 JSON-B 适配器应对这种情况。

### JSON-B 中的 null 值

请注意，使用一个 null 值来序列化一个字段的结果是，该属性将从输出 JSON 文档中排除。在执行反序列化时，缺少的属性不会导致目标属性被设置为 null，而且不会调用 setter 方法（或一个公共字段集）。该属性的值会保持不变，并允许在目标类中设置默认值。

## 标准 Java 类型的默认处理方式

JSON-B 支持标准 Java SE 类型，包括 `BigDecimal` 、 `BigInteger` 、 `URL` 、 `URI` 、 `Optional` 、 `OptionalInt` 、 `OptionalLong` 和 `OptionalDouble` 。

我前面已经提到过， `Number` 类型通过调用 `toString()` 方法来实现序列化，并使用合适的构造函数或工厂方法来实现反序列化。 `URL` 和 `URI` 类型的序列化行为相同，都通过调用 `toString()` 来完成。反序列化是通过使用合适的构造函数或工厂静态方法来执行的。

可选值的序列化方式是，检索它们的内部实例，并以适合该类型的方式将实例转换为 JSON（通常通过调用 `toString()` ）。清单 11 表明，一个 `Integer` 的 `Optional` 将通过调用该 `Integer` 的 `toString()` 值来实现序列化。

##### 清单 11\. 序列化一个 Optional 实例

```
public class OptionalExample {
public Optional<Integer> value;
// Constructors, getters and setters omitted for brevity
}

JsonBuilder.create().toJson(new OptionalExample(10))

```

Show moreShow more icon

输出 JSON 数字如清单 12 所示。

##### 清单 12\. OptionalExample 的 JSON 表示

```
{
"value": 10
}

```

Show moreShow more icon

根据目标字段的类型，JSON 值会反序列化为合适的 `Optional<T>` 或 `OptionalInt/Long/Double` 值。

清单 13 给出了一个将 JSON 数字反序列化为 `OptionalInt` 实例的示例。

##### 清单 13\. 将 JSON 数字反序列化为 OptionalInt

```
JsonBuilder.create().fromJson("{\"value\":10}", OptionalIntExample.class)

```

Show moreShow more icon

JSON 数据被反序列化为 `OptionalIntExample` 类的一个实例（清单 14），该实例有一个名为 _value_ 的 `OptionalInt` 字段。反序列化后，value 字段拥有值 `OptionalInt.of(10)` 。

##### 清单 14\. OptionalIntExample 类

```
public class OptionalIntExample {
public OptionalInt value;
// Constructors, getters and setters omitted for brevity
}

```

Show moreShow more icon

如您所见， `Optional` 类完全支持序列化和反序列化。

空 `Optional` 的序列化的处理方式与在序列化基本类型时的 null 处理方式相同。这意味着默认情况下不会保留它们。在执行反序列化时，null 值表示为 `Optional.empty()` （或 `Optional/Int/Long/Double.empty()` ）实例。相反地，存储在数组和映射数据结构中的空 `Optional` 被序列化为输出 JSON 数组中的 null。

## 特殊数据类型的默认处理方式

现在看看特殊数据类型的默认处理方式。（在 [第 2 部分](https://www.ibm.com/developerworks/cn/java/j-javaee8-json-binding-2/index.html?ca=drs-) 中，您将学习如何使用 `@JsonbNillable` 映射注解来自定义 null 的处理。）

### 枚举

枚举值的序列化方式是，调用 `name()` 方法来返回枚举常量的 `String` 表示。这特别重要，因为反序列化操作调用了 `valueOf()` 方法并将属性值传递给它。枚举类型应具有往返等效性。

### 日期类型

JSON Binding API 既支持来自（JDK 1.1 中引入的） `java.util` 包中的旧 `Date` 类的日期和时间实例，也支持 `java.time` 包中发现的新 Java 8 日期和时间类的日期和时间实例。

采用的默认时区是 GMT 标准时区，偏移指定为 UTC Greenwich。日期-时间格式是没有偏移的 ISO 8601。这些默认处理方式可使用第 3 部分中介绍的自定义选项进行覆盖。

除了已弃用的 3 字母时区 ID 外，时区实例受到全面支持。清单 15 和清单 16 给出了如何序列化日期和时间的两个示例。

##### 清单 15\. 对一个传统的 Date 实例执行序列化

```
SimpleDateFormat sdf = new SimpleDateFormat("dd/MM/yyyy");
Date date = sdf.parse("25/12/2018");
String json = jsonb.toJson(date);

Serialize to JSON: "2018-12-25T00:00:00Z[UTC]"

```

Show moreShow more icon

##### 清单 16\. 对一个 Java 8 Duration 实例执行序列化

```
Duration.ofHours(4).plusMinutes(3).plusSeconds(2)

Serialize to JSON: "PT4H3M2S"

```

Show moreShow more icon

请参阅 JSON-B 规范文档（第 3.5 小节），了解为每种类型采用何种格式的内幕。

### 数组和集合类型

所有类型的集合和数组都受支持并序列化为 JSON 数组，而 `Map` 类型序列化为 JSON 映射。清单 17 展示了一个 `Map` 的序列化，清单 18 展示了它的 JSON 表示。

##### 清单 17\. 一个 String 和 Integer 的映射

```
new HashMap<String, Integer>() {{
put("John", 10);
put("Jane", 10);
}};

```

Show moreShow more icon

##### 清单 18\. 一个 Map 的 JSON 表示

```
{
"John": 10,
"Jane": 10
}

```

Show moreShow more icon

当一个 `Collection` 或 `Map` 是某个类的成员时，它的字段名称被用作结果 JSON 对象的属性名称，如清单 19 和清单 20 所示。

##### 清单 19\. 一个作为类成员的 Map

```
public class ScoreBoard {
public Map<String, Integer> scores =
new HashMap<String, Integer>() {{
       put("John", 12);
       put("Jane", 34);
}};

}

```

Show moreShow more icon

##### 清单 20\. JSON 表示

```
{
"scores": {
"John": 12,
"Jane": 34
}
}

```

Show moreShow more icon

JSON-B 还支持多维原语和 Java 类型数组。

## 集合和映射中的 null 值

JSON-B 对数组和映射结构中的 null 值的处理不符合约定，它在序列化和反序列化时都保留了 null。这意味着空元素被序列化为 null，而 null 被反序列化为空元素。空 `Optional` 值被序列化为输出 JSON 数组中的 null，被反序列化为 `Optional.empty()` 实例。

### 无类型映射

日期、时间和日历类型都是数字，既包括旧的 `java.util.Date` 和 `java.util.Calendar` 类型，也包含 `java.time` 中发现的新 Java 8 日期和时间类。

我们已经看到，JSON 属性数据被反序列化为目标类中相应字段的 Java 类型；但是，如果输出类型未指定或指定为 `Object` ，JSON 值将被反序列化为相应的默认 Java 类型，如表 2 所示。

##### 表 2\. 默认反序列化类型

JSON 值Java 类型objectjava.util.Map<string, object>arrayjava.util.Liststringjava.lang.Stringnumberjava.math.BigDecimaltrue, falsejava.lang.Booleannullnull

清单 21 展示了这在实际中的意义，其中的一个 JSON 对象被反序列化为一个 Java `Map` 实例。

##### 清单 21\. 被反序列化为 Java Map 的一个 JSON 对象

```
String json = "{
\"title\":\"Fun with Java\",
\"price\":24.99,
\"issue\":null}";
Map<String, Object>() map =
JsonbBuilder.create().fromJson(json, Map.class);

```

Show moreShow more icon

## Java 类序列化和反序列化

_反序列化_ 操作要求目标类有一个无参数的 public 或 protected 构造函数，或者显式指定一种用于自定义对象构造的方法。如果未能给构造对象提供一个方法，则会导致抛出 `javax.json.bind.JsonbException` 。序列化没有这种需求。

JSON-B 可以反序列化 public、protected 和 protected-static 类，但不支持匿名类。匿名类的序列化会生成一个 JSON 对象。在反序列化期间，不支持已提及的接口以外的接口，而且在序列化时，会使用运行时类型。

序列化和反序列化都支持具有公开和受保护访问权的类的嵌套和静态嵌套类。

## 对 JSON-P 类型的支持

因为 JSON-B 是基于 JSON Processing API 而构建的，所以它支持所有 JSON-P 类型。这些类型来自 `javax.json` 包，包含 `JsonValue` 的所有子类型。清单 22 展示了一个被序列化为 JSON 数组的 `JsonArray` 实例（清单 23）。

##### 清单 22\. JSON Value 的序列化

```
JsonArray value = Json.createArrayBuilder()
       .add(Json.createValue("John"))
       .add(JsonValue.NULL)
       .build();

```

Show moreShow more icon

##### 清单 23\. 一个 JSON Value 实例的 JSON 表示

```
[
"John",
null
]

```

Show moreShow more icon

## 属性顺序

默认情况下，属性是按词典顺序进行序列化的。清单 24 中的类的实例将被序列化为清单 25 所示的 JSON 文档。

##### 清单 24\. 包含未按词典顺序排序的字段的类

```
public class LexicographicalOrder {
public String dog = "Labradoodle";
public String animal = "Cat";
public String bread = "Chiapata";
public String car = "Ford";
}

```

Show moreShow more icon

##### 清单 25\. 按词典顺序显示属性的序列化

```
{
"animal": "Cat",
"bread": "Chiapata",
"car": "Ford",
"dog": "Labradoodle"
}

```

Show moreShow more icon

JSON-B 的属性排序策略可通过一个自定义选项进行配置。 [第 2 部分](https://www.ibm.com/developerworks/cn/java/j-javaee8-json-binding-2/index.html?ca=drs-) 会更详细讨论此主题。

## 默认访问策略

正如我在 [第 1 部分](https://www.ibm.com/developerworks/cn/java/j-javaee8-json-binding-1/index.html?ca=drs-) 中简单讨论的那样，JSON-B 的默认字段访问策略要求方法或字段允许通过指定 public 访问修饰符来实现公共访问。在序列化期间，会调用字段的公共 getter 方法来检索它的值。如果这个 getter 不存在（或者如果它没有公共访问权），则会直接访问该字段 — 但仅在访问修饰符为 public 时才会出现这种情况。

类似地，反序列化依靠 setter 方法的公共访问性来设置属性值。如果该方法不存在或不是公共的，则会直接设置该 public 字段。

可通过调节相关的自定义选项，将 JSON-B 的字段访问策略配置为更加受限。 [第 2 部分](https://www.ibm.com/developerworks/cn/java/j-javaee8-json-binding-2/index.html?ca=drs-) 将更详细地介绍此主题和其他自定义方法。

## 结束语

JSON Binding API 可直接满足大部分简单用例场景的需求，而且提供了合理的默认处理方式。对于熟悉其他任何 JSON 绑定技术的开发人员，使用 JSON-B 将会很简单，而且它很容易使用，即使不熟悉 JSON 的开发人员也应能轻松上手。

您已初步认识 JSON-B，并了解了它的主要默认特性和功能。要发挥出 JSON Binding 的真正威力，有时需要自定义这些默认操作。在第 3 部分中，我将介绍 JSON-B 的自定义模型和特性。您将了解如何使用编译时注解和运行时配置来自定义该 API 的几乎所有方面，包括低级别序列化和反序列化操作。

本文翻译自： [Default mapping with JSON-B](https://developer.ibm.com/articles/j-javaee8-json-binding-2/)（2017-11-18）