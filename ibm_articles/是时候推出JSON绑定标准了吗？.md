# 是时候推出 JSON 绑定标准了吗？
通过比较 Gson、Jackson 和 JSON-B 来突出基本特性和行为中的不一致性

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-javaee8-json-binding-4/)

Alex Theedom

发布: 2017-01-09

* * *

**关于本系列**：Java EE 支持 XML 已久，但显然遗漏了对 JSON 数据的内置支持。Java EE 8 改变了这一状况，给核心 Java 企业平台带来了强大的 JSON 绑定特性。开始使用 JSON-B，了解如何结合使用它与 JSON Processing API 及其他技术来处理 Java 企业应用程序中的 JSON 文档。

JSON 成为首选的数据交换格式已有近 10 年，但直到现在，我们仍没有一个稳健到足以用作标准的规范。在 JSON Binding API 发布后，供应商有了对 JSON 解析器中的默认和自定义绑定操作执行标准化的机会。作为开发人员，我们应坚持要求他们这么做。

比较 3 种 JSON 解析器的序列化和反序列化行为，揭示这些工具在处理常规绑定操作方面的差异和特有风格。同样明显的是，许多关键特性受到的支持并不一致，一些工具提供了这些特性，而另一些工具则没有提供。

尽管这反映出解析器之间的差异化，但不应以牺牲功能为代价。允许各种实现间的默认行为存在如此大的差异，这对工具选择的影响可能比对性能差异的影响更大，而性能差异应该是主要决定因素。

## JSON 绑定标准的现状

构建有效的标准需要时间，而且在基于既定的用例和最佳实践时最有效。一个好标准是基于多年实验、更正和创新的领域知识的结晶。对于 JSON 解析器，我相信建立标准的时机成熟了。

**您觉得呢？** 我相信，采用 JSON 绑定标准是 JSON 进化的下一步，而 JSR 367 正是实现这一进化的合适规范。

对 JSON 绑定操作和特性执行标准化，使开发人员能基于基础引擎的效率和性能来比较工具。实现人员应该更多地关注性能而不是特性。工具的整体质量将有所改进

标准化不会扼杀创新；相反，它会反映已发生的创新。它为未来的创新和进步提供了良好基础，而且最佳的增补功能将成为合并到标准中的候选功能。

我相信实现标准化的理由很充分，而且有了 JSON Binding API 规范，没有理由不去执行标准化。在以下各节中，我将比较 3 种 JSON 绑定工具的核心特性。尽管这些比较很重要，但我的主要目标是展示工具在处理 Java™ 中的绑定操作和类型格式上的差异。

## 比较 JSON 解析器

本文将比较两个流行框架（FastXML 的 [Jackson Project](https://github.com/FasterXML/jackson) 和 Google 的 [Gson](https://github.com/google/gson) ）与 JSON Binding API。

我们首先比较各个框架对结构化 JSON 文档中的关键数据类型的默认处理方式。然后，我们将了解每个框架如何处理这些数据类型的自定义，从基本 Java 类型到 Java 8 中出现的新的日期和时间格式。

您将亲眼见到每个 JSON 框架在 Java 对象与 JSON 之间执行序列化和反序列化时的行为。我的目标不是执行科学评测或基准测试比较，而是概述框架之间的差异和不一致性。最后，我最关注的是演示在 JSON Binding 规范上执行标准化的好处。

对于本文中的示例，我使用了带 JSON Processing V1.1 的 JSON Binding API V1.0 作为提供程序，还使用了 Gson V2.8.2 和 Jackson V2.9.2。我还创建了一个 [GitHub 存储库](https://github.com/readlearncode/JSON-framework-comparison) ，其中包含代码示例和单元测试，以便您可以自行复制这些示例。

[获得代码](https://github.com/readlearncode/JSON-framework-comparison)

## 绑定模型

在 JSON 解析器中， _绑定模型_ 是序列化和反序列化的机制。所有 3 个框架都为常见绑定场景提供了默认绑定行为，每个框架都提供了自定义绑定解决方案。对于自定义绑定，开发人员可以选择使用编译时注解或运行时配置。在本节，我将展示 JSON-B、Jackson 和 Gson 如何通过运行时配置实现自定义绑定。各个模型完全不同，所以这是一个执行比较的不错起点。后面的小节会详细探讨自定义特性。

### 通过 JSON-B 执行运行时绑定

开发人员通过 `JsonConfig` 类访问 JSON Binding 的运行时配置。您可以通过调用 `with(...)` 方法并传入配置属性来自定义运行时配置。例如， `new JsonConfig.withPropertyOrderStrategy(PropertyOrderStrategy.REVERSE)` 。

然后在 `JsonBuilder` 实例上设置 `JsonConfig` 实例， `toJson()` 或 `fromJson()` 方法会在 JsonBuilder 实例上调用。清单 1 给出了一个示例。

##### 清单 1\. JSON Binding API 的运行时配置

```
JsonbConfig jsonbConfig = new JsonbConfig()
    .with...(...);

String json = JsonbBuilder
    .create(jsonbConfig)
    .toJson(...);

AnObject anObj = JsonbBuilder
    .create(jsonbConfig)
    .fromJson("{JSON}", AnObject.class);

```

Show moreShow more icon

### 通过 Jackson 执行运行时绑定

Jackson 的运行时模型不同于 JSON-B，您可以直接配置对象映射器引擎。然后使用这个引擎在目标对象上执行序列化和反序列化。在清单 2 中可以看到， `ObjectMapper` 提供了非常多的配置方法。

##### 清单 2\. Jackson 的运行时模型

```
ObjectMapper objectMapper = new ObjectMapper()
    .set...(...)
    .configure(...)
    .addHandler(...)
    .disable(...)
    .enable(...)
    .registerModule(...)
    ...;

String json = objectMapper
    .writeValueAsString(...);

AnObject anObj = objectMapper
    .readValue("{JSON}", AnObject.class);

```

Show moreShow more icon

### 通过 Gson 执行运行时绑定

Gson 配置了 `GsonBuilder` 的一个实例， `Gson` 实例就是通过该实例创建的。一个 Gson 对象提供了执行序列化和反序列化操作的方法。可以在清单 3 中看到 Gson 运行时模型。

##### 清单 3\. Gson 运行时模型

```
GsonBuilder gsonBuilder = new GsonBuilder()
    .set...(...)
    .add....(...)
    .register...(...)
    .serialize...(...)
    .enable...(...)
    .disable...(...)
    ...;

String json = gsonBuilder
    .create()
    .toJson(...);

AnObject anObj = gsonBuilder
    .create()
    .fromJson("{JSON}", AnObject.class);

```

Show moreShow more icon

接下来，我们将看看这 3 个解析器如何处理类型，首先从基本 Java 类型开始。

## 基本 Java 类型

在 Java 中， **基本类型** 指的是所有原语类型和它们各自的包装器类，以及 `String` 类型。 我们将比较一个简单示例，该示例将所有这些类型都序列化为 JSON。

首先，清单 4 中的类包含 Java 语言的所有基本类型。

##### 清单 4\. 基本 Java 类型

```
public class AllBasicTypes {

    // Primitive types
    private byte bytePrimitive;
    private short shortPrimitive;
    private char charPrimitive;
    private int intPrimitive;
    private long longPrimitive;
    private float floatPrimitive;
    private double doublePrimitive;
    private boolean aBoolean;

    // Wrapper types
    private Byte byteWrapper = 0;
    private Short shortWrapper = 0;
    private Character charWrapper = 0;
    private Integer intWrapper = 0;
    private Long longWrapper = 0L;
    private Float floatWrapper = 0F;
    private Double doubleWrapper = 0D;
    private Boolean booleanWrapper = false;

    private String string = "Hello World";

    // Getters and setter omitted for brevity
}

```

Show moreShow more icon

通过对所有 3 个框架执行序列化，得到了清单 5 中所示的 JSON 文档。您会看到，这种序列化行为是目前唯一在这些框架中保持一致的行为。请注意，字段具有自描述性，这使匹配 JSON 属性与源类变得更容易。

##### 清单 5\. JSON 文档

```
{
    "bytePrimitive": 0,
    "shortPrimitive": 0,
    "charPrimitive": "\u0000",
    "intPrimitive": 0,
    "longPrimitive": 0,
    "floatPrimitive": 0,
    "doublePrimitive": 0,
    "aBoolean": false,
    "byteWrapper": 0,
    "shortWrapper": 0,
    "charWrapper": "\u0000",
    "intWrapper": 0,
    "longWrapper": 0,
    "floatWrapper": 0,
    "doubleWrapper": 0,
    "booleanWrapper": false,
    "string": "Hello World"
}

```

Show moreShow more icon

所有 3 个 JSON 解析器都对基本类型上的操作使用了默认配置，以及少量配置选项。（请参阅 [JSON Binding API 入门，第 2 部分](https://www.ibm.com/developerworks/cn/java/j-javaee8-json-binding-2/index.html?ca=drs-) ，查看使用 JSON-B 执行数字格式化的演示。）

## 特殊 Java 类型

所谓的 **特殊类型** 包括更复杂的 JDK 类型，比如扩展 `java.lang.Number` 的类型。其中包括 `BigDecimal` 、 `AtomicInteger` 和 `LongAdder` ；包含 `OptionalInt` 和 `Optional<type>` 在内的 `Optional` 类型；以及 `URL` 和 `URI` 实例。

首先比较 `Number` 类型。所有 JSON 框架都通过仅提取内部数字作为属性值来执行序列化。处理 `BigDecimal` 、 `BigInteger` 、 `AtomicInteger` 和 `LongInteger` 的方式在 3 个框架中保持一致。

处理 `URI` 和 `URL` 类的方式在所有 3 个框架中也相同：执行序列化后，在获得的 JSON 文档中，内部值表示为一个 `String` 。

但是，处理 `Optional` 类型的方式是不一致的。表 1 给出了对清单 6 中的两个实例执行序列化的结果。

##### 清单 6\. Optional 测试类

```
Optional<String> stringOptional = Optional.of("Hello World");
OptionalInt optionalInt = OptionalInt.of(10);

```

Show moreShow more icon

在表 1 中，我将这两个字段分成两列，以便更容易比较。

##### 表 1\. Optional 类型

OptionalOptionalIntJSON-B`"stringOptional": "Hello World"``"optionalInt": 10`Jackson`"stringOptional": { "present": true }``"optionalInt": { "asInt": 10,"present": true }`Gson`"stringOptional": { "value": "Hello World" }``"optionalInt": { "isPresent": true, "value": 10 }`

如您所见，JSON-B 通过内部探查并检索内部值来实现 `Optional` 值的全面支持。Jackson 和 Gson 都无法直接支持 `Optional` ，而且 `String Optional` 的序列化最不一致。

Gson 将 `Optional` 的结构序列化为一个包含 `String` 的 JSON 对象，Jackson 甚至不包含基础实例的值，只是表明该值存在。

Gson 和 Jackson 无法很好地支持 `OptionalInt` 和 `Optional` 数字系列的其他类型，这些类型仅仅被序列化为合理的 JSON 对象，其中包含基础值和表示该值存在的布尔值。

## 日期、时间和日历类型

日期、时间和日历类型都是数字，既包括旧的 `java.util.Date` 和 `java.util.Calendar` 类型，也包含 `java.time` 中的新 Java 8 日期和时间类。

考虑到要比较的类型有很多，我仅选择了 3 种类型：

- 一个使用构建器设置日期的 `Calendar` 实例。
- 一个利用 `String` 设置日期并应用 `SimpleDateFormat` 的 `Date` 实例。
- 一个使用 `parse()` 方法设置日期的 Java 8 `LocalDate` 实例。

### Calendar 类型

`Calendar` 实例配置了日期 2017 年 12 月 25 日，如清单 7 所示。

##### 清单 7\. Calendar 类型

```
new Calendar.Builder().setDate(2017, 1, 25).build()

```

Show moreShow more icon

表 2 给出了分别使用 3 个框架对此实例执行序列化的结果。

##### 表 2\. 序列化后的 Calendar 实例

Calendar 实例JSON-B`2017-12-25T00:00:00Z[Europe/London]`Jackson`1514160000000`Gson`{ "year": 2017, "month": 11, "dayOfMonth": 25, "hourOfDay": 0, "minute": 0, "second": 0 }`

在表 2 中可以清楚地看到，每个框架都实现了不同方法来序列化 `Calendar` 实例。JSON-B 输出了包含时间零点和时区的日期。Jackson 生成自新时代（1970 年 1 月 1 日）以来的毫秒数，Gson 生成一个包含该日期的数字部分的 JSON 对象。

默认行为的差别很大。尽管可以将这些格式配置为更容易理解的格式，但输出格式其实不应存在如此大的差别。

### Date 和 LocalDate 类型

现在看看如何处理 `Date` 和 `LocalDate` 。清单 8 显示了创建这两个实例的代码，表 3 显示了它们的序列化。

##### 清单 8\. LocalDate 和 Date 实例

```
new SimpleDateFormat("dd/MM/yyyy").parse("25/12/2017");
LocalDate.parse("2017-12-25");

```

Show moreShow more icon

##### 表 3\. 序列化后的 Date 和 LocalDate 实例

DateLocalDateJSON-B`2017-12-25T00:00:00Z[UTC]``2017-12-25`Jackson`1514160000000``{ "year": 2017, "month": "DECEMBER", "chronology": { "id": "ISO", "calendarType": "iso8601" }, "dayOfMonth": 25, "dayOfWeek": "MONDAY", "dayOfYear": 359, "leapYear": false, "monthValue": 12, "era": "CE" }`Gson`Dec 25, 2017 12:00:00 AM``{ "year": 2017, "month": 12, "day": 25 }`

### Date 类型

JSON Binding 忽略原始日期格式，应用默认格式样式，并添加时间和时区信息。Jackson 返回自新时代以来的毫秒数，Gson 应用了 `MEDIUM` 日期/时间样式。

### LocalDate 类型

`LocalDate` 实例的处理方式与 `Date` 不同，而且从 JSON-B 的角度讲，该处理方式更合理。JSON-B 仅以默认格式输出日期，没有时间和时区信息。Jackson 通过调用实例中的 accessor 方法来生成一个 JSON 对象，实际上，它将 `LocalDate` 作为 POJO 进行处理。Gson 使用年、月和日 3 个属性构造了一个对象。

`Date` 、 `LocalDate` 和 `Calendar` 实例在 3 个 JSON 框架中的处理方式明显不一致。

### 配置日期和时间格式

可以对日期和时间格式执行编译时和运行时自定义。只有 JSON-B 和 Jackson 提供了注解，以便为给定属性、方法或类指定日期和时间格式，如清单 9 和清单 10 所示。

##### 清单 9\. 在编译时通过 JSON Binding 配置日期格式

```
@JsonbDateFormat(value = "MM/dd/yyyy", locale = "Locale.ENGLISH")

```

Show moreShow more icon

##### 清单 10\. 在编译时通过 Jackson 配置日期格式

```
@JsonFormat(pattern = "MM/dd/yyyy", locale = "Locale.ENGLISH")

```

Show moreShow more icon

所有 3 个解析器都有全局设置该格式的运行时配置，如清单 11、12 和 13 所示。

##### 清单 11\. 在运行时通过 JSON Binding 配置日期格式

```
new JsonbConfig().withDateFormat("MM/dd/yyyy", Locale.ENGLISH)

```

Show moreShow more icon

##### 清单 12\. 在运行时通过 Jackson 配置日期格式

```
new ObjectMapper().setDateFormat(new SimpleDateFormat("MM/dd/yyyy"))

```

Show moreShow more icon

##### 清单 13\. 在运行时通过 Gson 配置日期格式

```
new GsonBuilder().setDateFormat("MM/dd/yyyy")

```

Show moreShow more icon

## 数组、集合和映射

在处理数组、集合和映射的方式上，这些解析器惊人地一致。原因在于，每个结构都分别直接映射到等价的 JSON 类型。所以无论采用哪个 JSON 框架，清单 14 中的代码都会序列化为清单 15 中的 JSON 文档。

##### 清单 14\. 示例数组、集合和映射

```
private int[] intArray = new int[]{ 1, 2, 3, 4 };
private String[] stringArray = new String[]{ "one", "two" };

Collection<Object> objectCollection = new ArrayList<Object>() {{
add("one");
add("two");
}};

private Map<String, Integer> stringIntegerMap = new HashMap<String, Integer>() {{
put("one", 1);
put("two", 2);
}};

```

Show moreShow more icon

这是 JSON 文档。

##### 清单 15\. 序列化后的数组、集合和映射

```
{
    "intArray": [1, 2, 3, 4],
    "objectCollection": ["one", "two"],
    "stringArray": ["one", "two"],
    "stringIntegerMap": { "one": 1, "two": 2 }
}

```

Show moreShow more icon

## null、集合中的 null 和 Optional.empty()

JSON Binding 或 Gson 不会序列化 null，但 Jackson 保留了它们。在执行反序列化时，JSON 文档中缺少一个值不会导致在目标对象中调用相应的 setter 方法。但是，如果该值为 null，则会将它设置为正常值。

对于所有 3 个框架，数组、映射和集合中默认情况下都会保留 null。

### Optional.empty()

`Optional.empty()` 值表示一个不存在的值，因此 JSON Binding 采用与 null 类似的方式处理它：在执行序列化时，JSON 文档不会包含该属性。但是，Jackson 和 Gson 都会尝试将此值序列化为 JSON 对象：Jackson 将 `Optional.empty()` 作为 POJO 进行处理，而 Gson 会生成一个空 JSON 对象。

表 4 总结了 `Optional<Object> emptyOptional = Optional.empty()` 的序列化。

##### 表 4\. Optional.Empty() 实例的序列化

JSON-BJacksonGson“Property not included in JSON document.”`{ "emptyOptional": { "present": false } }``{ "emptyOptional": {} }`

### null 的配置选项

JSON Binding API 同时为 null 提供了运行时和编译时配置选项。要通过编译时配置在特定字段中包含 null，可以使用 `@JsonbProperty` 注解并将 `nillable` 标志设置为 `true` 。 也可以在类或包级别上使用 `@JsonbNillable` 注解来全局设置包含 null，如清单 16 所示。

##### 清单 16\. JSON Binding 编译时 null 配置

```
@JsonbNillable
public class CompileTimeSampler {
@JsonbProperty(nillable = true)
private String nillable;
// field level configuration overrides class and package level configuration settings.
}

```

Show moreShow more icon

对于运行时配置，可以使用方法 `.withNullValues()` 设置为 true，如清单 17 所示。

##### 清单 17\. JSON Binding 的运行时 null 配置

```
new JsonbConfig().withNullValues(true);

```

Show moreShow more icon

Jackson 还为排除 null 值同时提供了运行时和编译时选项。清单 18 展示了如何使用 `@JsonInclude` 注解排除 null 值。

##### 清单 18\. 一个排除 null 的 Jackson 编译时注解

```
@JsonInclude(JsonInclude.Include.NON_NULL)

```

Show moreShow more icon

Jackson 提供了两个运行时配置选项。一个选项允许全局排除 null，而另一个支持从映射中排除 null（但已弃用），如清单 19 所示。

##### 清单 19\. 一种排除 null 的 Jackson 运行时配置

```
new ObjectMapper()
    .setSerializationInclusion(JsonInclude.Include.NON_NULL)
    .configure(SerializationFeature.WRITE_NULL_MAP_VALUES, false);

```

Show moreShow more icon

Gson 的运行时 null 配置会全局开启序列化并保留 null，如清单 20 所示。

##### 清单 20\. 一种排除 null 的 Gson 运行时配置

```
new GsonBuilder().serializeNulls();

```

Show moreShow more icon

我们再次看到，3 个框架处理类型序列化的方式是不一致的。不仅对 null 的处理不一致，甚至连配置选项也不同。

## 字段可视性

字段的可视性决定了是否会将字段序列化为 JSON 文档和反序列化为目标实例的字段。 此外，每个框架对可视性的解释都不同。

清单 21 中的代码使用所有 3 个框架执行了序列化。这些字段具有自描述性。一些字段有相应的 setter/getter 方法（从字段名称可以看出），而剩余字段没有。此外，一个 virtual 字段没有相应的字段，仅通过一个 getter 方法来表示。

##### 清单 21\. 包含字段可视性选项的类

```
public class FieldsVisibility {

    // Without getter and setters
    private String privateString = "";
    String defaultString = "";
    protected String protectedString = "";
    public String publicString = "";
    final private String finalPrivateString = "";
    final public String finalPublicString = "";
    static public String STATIC_PUBLIC_STRING = "";

    // With getter and setters.Omitted for brevity
    private String privateStringWithSetterGetter = "";
    String defaultStringWithSetterGetter = "";
    protected String protectedStringWithSetterGetter = "";
    public String publicStringWithSetterGetter = "";
    final private String finalPrivateStringWithSetterGetter = "";
    final public String finalPublicStringWithSetterGetter = "";
    static public String STATIC_PUBLIC_STRING_WITH_SETTER_GETTER = "";

    public String getVirtualField() {
        return "";
    }
}

```

Show moreShow more icon

因此，我们可以比较每个框架如何处理字段可视性，表 5 给出了序列化清单 21 中的代码并编译结果后的效果。

##### 表 5\. 每个框架的字段可视性默认处理方式的总结

访问方式字段修饰符JSON-BJacksonGson字段`private`[否](images/no.png)[否](images/no.png)[是](images/yes.png)`default`[否](images/no.png)[否](images/no.png)[是](images/yes.png)`protected`[是](images/yes.png)[是](images/yes.png)[是](images/yes.png)`public`[是](images/yes.png)[是](images/yes.png)[是](images/yes.png)`final public`[否](images/no.png)[否](images/no.png)[是](images/yes.png)`final private`[否](images/no.png)[否](images/no.png)[否](images/no.png)`static public`[是](images/yes.png)[是](images/yes.png)[是](images/yes.png)公共 Getter`private`[是](images/yes.png)[是](images/yes.png)[是](images/yes.png)`default`[是](images/yes.png)[是](images/yes.png)[是](images/yes.png)`protected`[是](images/yes.png)[是](images/yes.png)[是](images/yes.png)`public`[是](images/yes.png)[是](images/yes.png)[是](images/yes.png)`final private`[是](images/yes.png)[是](images/yes.png)[是](images/yes.png)`final public`[是](images/yes.png)[是](images/yes.png)[是](images/yes.png)`static public`[否](images/no.png)[否](images/no.png)[否](images/no.png)`virtual field`[是](images/yes.png)[是](images/yes.png)[否](images/no.png)

唯一一致的地方是，所有 3 个框架都很重视公共非静态访问修饰符。此外，所有框架都不会序列化 public static 字段，即使该字段有一个公共 getter 方法。

显然，Gson 完全不重视字段的访问修饰符，而且在序列化中包含所有字段（除了 public static 字段），无论指定的修饰符是什么。Gson 也不包含 virtual 字段，而 JSON-B 和 Jackson 都包含该字段。

JSON-B 和 Jackson 处理字段可视性的方法相同。Gson 的不同之处在于，它包含 final private 字段，而且不包含 virtual 字段。

## 可视性配置

可视性配置选项非常丰富，而且可能很复杂。让我们看看每个框架如何配置可视性。

### JSON Binding

JSON-B 提供了一个简单的 `@JsonbTransient` 注解，以便利用序列化和反序列化来排除任何带注解的字段。

在 JSON-B 中，有一种控制字段可视性的更复杂方式。可通过以下方式创建自定义可视性策略：实现 `PropertyVisibilityStrategy` 接口，并将它设置为 `JsonbConfig` 实例上的一个运行时属性。清单 22 展示了这一过程。（要获得 `PropertyVisibilityStrategy` 类的实现示例，请参阅 [第 2 部分](https://www.ibm.com/developerworks/cn/java/j-javaee8-json-binding-2/index.html?ca=drs-) 。）

##### 清单 22\. 在 JSON-B 中设置自定义可视性策略

```
new JsonbConfig()
    .withPropertyVisibilityStrategy(new CustomPropertyVisibilityStrategy());

```

Show moreShow more icon

### Jackson

Jackson 提供了编译时注解来忽略任何带 `@JsonIgnore` 注解的字段，以及在传递给 `@JsonIgnoreProperties({"aField"})` 类级注解的列表中明确提及的任何字段。方法和字段可视性可使用 `@JsonAutoDetect` 进行配置，如清单 23 所示。

##### 清单 23\. 通过 Jackson 注解设置字段可视性

```
@JsonAutoDetect(fieldVisibility = JsonAutoDetect.Visibility.PROTECTED_AND_PUBLIC)
public class CompileTimeSampler {
    // fields and methods omitted for brevity
}

```

Show moreShow more icon

Jackson 还提供了一种运行时配置，允许基于访问修饰符来显式包含字段。清单 24 中的代码示例展示了如何在序列化期间包含 public 和 protected 字段。

##### 清单 24\. 配置以包含 public 和 protected 字段

```
new ObjectMapper()
.setDefaultVisibility(
    JsonAutoDetect.Value.construct(PropertyAccessor.FIELD,
      JsonAutoDetect.Visibility.PROTECTED_AND_PUBLIC));

```

Show moreShow more icon

### Gson

Gson 提供了一个编译时注解来指定要在序列化和反序列化操作中包含的字段。在这种情况下，可以使用 `@Expose` 注解来标记每个字段，并将 `serialize` 和/或 `deserialize` 标志设置为 `true` 或 `false` ，如清单 25 所示。

##### 清单 25\. 使用 @Expose 注解

```
@Expose(serialize = false, deserialize = false)
private String aField;

```

Show moreShow more icon

要使用此注解，必须在 GsonBuilder 实例上调用 `excludeFieldsWithoutExposeAnnotation()` 方法来启用它，如清单 26 所示。

##### 清单 26\. 启用 @Expose 注解

```
new GsonBuilder().excludeFieldsWithoutExposeAnnotation();

```

Show moreShow more icon

也可以基于字段访问修饰符来显式排除字段，如清单 27 所示。

##### 清单 27\. 基于修饰符来排除字段

```
new GsonBuilder().excludeFieldsWithModifiers(Modifier.PROTECTED);

```

Show moreShow more icon

## 配置属性顺序

[IETF RFC 7159](https://tools.ietf.org/html/rfc7159#section-1) JSON 交换格式规范将 JSON 对象定义为”零个或多个名称/值对的无序集合”。但是，在某些情况下，需要按给定顺序显示属性。默认情况下，JSON-B 按词典顺序对属性进行排序，而 Jackson 和 Gson 采用字段在类中出现的顺序。

在所有 3 个框架中，都能通过显式列出字段名称或指定顺序策略来指定属性顺序。

### JSON Binding

JSON-B 同时提供了运行时和编译时机制来指定属性顺序。在运行时，可以通过将一个字段列表传递给 `@JsonbPropertyOrder` 注解来指定顺序，如清单 28 所示。

##### 清单 28\. 在 JSON-B 中显式指定字段顺序

```
@JsonbPropertyOrder({"firstName", "title", "author"})

```

Show moreShow more icon

对于运行时配置，可以通过在一个 `JsonbConfig` 实例上指定想要的顺序策略，从而设置一个 **全局顺序策略** ，如清单 29 所示。

##### 清单 29\. 在 JSON-B 中指定按词典逆序进行排序

```
new JsonbConfig()
    .withPropertyOrderStrategy(PropertyOrderStrategy.REVERSE);

```

Show moreShow more icon

### Jackson

Jackson 也提供了一种显式规定字段顺序的方式，并指定按字母顺序进行排序，如清单 30 所示。

##### 清单 30\. 在 Jackson 中指定字段顺序的两种方式

```
@JsonPropertyOrder(value = {"firstName", "title", "author"}, alphabetic = true)

@JsonPropertyOrder(alphabetic = true)

```

Show moreShow more icon

也可以使用运行时配置来指定全局属性顺序，如清单 31 所示。

##### 清单 31\. 在 Jackson 中指定全局属性顺序

```
new ObjectMapper().configure(MapperFeature.SORT_PROPERTIES_ALPHABETICALLY, true);

```

Show moreShow more icon

### Gson

Gson 没有提供能轻松对属性进行排序的配置设置。 您必须创建一个自定义序列化器和反序列化器来封装字段排序逻辑。

## 继承

JSON-B 处理继承的方式是先放置父类，后放置子类。清单 32 给出了一个 `Parent` 类的代码，清单 33 给出了扩展该 `Parent` 类的 `Child` 类的代码。

##### 清单 32.Parent 类

```
public class Parent {
private String parentName = "Parent";
// Getters/setter omitted
}

```

Show moreShow more icon

##### 清单 33.扩展 Parent 的 Child 类

```
public class Child extends Parent {
private String child = "Child";
// Getters/setter omitted
}

```

Show moreShow more icon

Jackson 也遵守继承顺序，它将父类放置在子类前，但 Gson 不遵守这种顺序。表 6 展示了所有 3 个框架如何处理 `Child` 类的序列化。

##### 表 6\. 每个框架的字段可视性默认处理方式的总结

JSON BindingJacksonGson`{ "parentName": "Parent", "child": "Child" }``{ "parentName": "Parent", "child": "Child" }``{ "child": "Child", "parentName": "Parent" }`

## 对 JSON Processing 类型的支持

目前，只有 JSON Binding API 提供了对 JSON Processing 类型的支持 ( [JSR 374](https://www.jcp.org/en/jsr/detail?id=374))，但 Jackson 提供了一个能处理 JSON Processing 类型的 [扩展模块](https://github.com/FasterXML/jackson-datatype-jsr353) 。回想第 1 部分，其中介绍了 JSON Processing 类型如何为 JSON 结构建模。例如， `JsonObject` 类为 JSON 对象建模， `JsonArray` 为 JSON 数组建模。JSON Binding API 也提供了两个用于生成 JSON 结构的模型。Jackson 和 Gson 都尝试利用基础映射实现来直接序列化 JSON Processing 实例，如表 7 所示。

##### 表 7\. JSONObject 的序列化

JSON BindingJacksonGson`"jsonObject": { "firstName": "Alex", "lastName": "Theedom" }``"jsonObject": { "firstName": { "chars": "Alex", "string": "Alex", "valueType": "STRING" }, "lastName": { "chars": "Theedom", "string": "Theedom", "valueType": "STRING" } }``"jsonObject": { "firstName": { "value": "Alex" }, "lastName": { "value": "Theedom" } }`

## 结束语

Gson 和 Jackson 是已投入使用的、非常受欢迎的 JSON 框架，它们提供了丰富的特性，尤其是自定义方面的特性。但是，正如我所演示的，各个框架在处理简单场景的方式上很少一致。甚至在最简单的 null 处理上也有所不同。

在类型处理上保持了一致性，这是因为没有其他选项可供选择。基本类型只能实际序列化为它们的实际值，而映射和集合与等价的 JSON 对象可以直接来回转换。最明显的不一致表现在日期和时间的处理上：各个框架序列化这些类型的方式各不相同。尽管有控制格式的配置选项，但诸多的不一致表明需要一个标准。

开发人员不应被迫牺牲性能来换取特性，而 JSON 绑定标准可以解决这一问题。我相信，采用 JSON 绑定标准是 JSON 进化的下一步，而 JSR 367 正是实现这一进化的合适规范。

本文翻译自： [Is it time for a JSON binding standard?](https://developer.ibm.com/articles/j-javaee8-json-binding-4/)（2017-12-15）