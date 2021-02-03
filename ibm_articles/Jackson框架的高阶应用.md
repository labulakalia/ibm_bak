# Jackson 框架的高阶应用
充分发挥 Jackson 的使用价值

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/jackson-advanced-application/)

刘万振

发布: 2017-11-15

* * *

Jackson 是当前用的比较广泛的，用来序列化和反序列化 json 的 Java 的开源框架。Jackson 社 区相对比较活跃，更新速度也比较快， 从 Github 中的统计来看，Jackson 是最流行的 json 解析器之一 。 Spring MVC 的默认 json 解析器便是 Jackson。 Jackson 优点很多。 Jackson 所依赖的 jar 包较少 ，简单易用。与其他 Java 的 json 的框架 Gson 等相比， Jackson 解析大的 json 文件速度比较快；Jackson 运行时占用内存比较低，性能比较好；Jackson 有灵活的 API，可以很容易进行扩展和定制。

Jackson 的 1.x 版本的包名是 org.codehaus.jackson ，当升级到 2.x 版本时，包名变为 com.fasterxml.jackson，本文讨论的内容是基于最新的 Jackson 的 2.9.1 版本。

Jackson 的核心模块由三部分组成。

- jackson-core，核心包，提供基于”流模式”解析的相关 API，它包括 JsonPaser 和 JsonGenerator。 Jackson 内部实现正是通过高性能的流模式 API 的 JsonGenerator 和 JsonParser 来生成和解析 json。
- jackson-annotations，注解包，提供标准注解功能；
- jackson-databind ，数据绑定包， 提供基于”对象绑定” 解析的相关 API （ ObjectMapper ） 和”树模型” 解析的相关 API （JsonNode）；基于”对象绑定” 解析的 API 和”树模型”解析的 API 依赖基于”流模式”解析的 API。

在了解 Jackson 的概要情况之后，下面介绍 Jackson 的基本用法。

## Jackson 的 基本用法

若想在 Java 代码中使用 Jackson 的核心模块的 jar 包 ，需要在 pom.xml 中添加如下信息。

##### 清单 1.在 pom.xml 的 Jackson 的配置信息

```
<dependency>
<groupId>com.fasterxml.jackson.core</groupId>
<artifactId>jackson-databind</artifactId>
<version>2.9.1</version>
</dependency>

```

Show moreShow more icon

jackson-databind 依赖 jackson-core 和 jackson-annotations，当添加 jackson-databind 之后， jackson-core 和 jackson-annotations 也随之添加到 Java 项目工程中。在添加相关依赖包之后，就可以使用 Jackson。

### ObjectMapper 的 使用

Jackson 最常用的 API 就是基于”对象绑定” 的 ObjectMapper。下面是一个 ObjectMapper 的使用的简单示例。

##### 清单 2 . ObjectMapper 使用示例

```
ObjectMapper mapper = new ObjectMapper();
Person person = new Person();
person.setName("Tom");
person.setAge(40);
String jsonString = mapper.writerWithDefaultPrettyPrinter()
.writeValueAsString(person);
Person deserializedPerson = mapper.readValue(jsonString, Person.class);

```

Show moreShow more icon

ObjectMapper 通过 writeValue 系列方法 将 java 对 象序列化 为 json，并 将 json 存 储成不同的格式，String（writeValueAsString），Byte Array（writeValueAsString），Writer， File，OutStream 和 DataOutput。

ObjectMapper 通过 readValue 系列方法从不同的数据源像 String ， Byte Array， Reader，File，URL， InputStream 将 json 反序列化为 java 对象。

### 信息配置

在调用 writeValue 或调用 readValue 方法之前，往往需要设置 ObjectMapper 的相关配置信息。这些配置信息应用 java 对象的所有属性上。示例如下：

##### 清单 3 . 配置信息使用示例

```

//在反序列化时忽略在 json 中存在但 Java 对象不存在的属性
mapper.configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES,
    false);
//在序列化时日期格式默认为 yyyy-MM-dd'T'HH:mm:ss.SSSZ
mapper.configure(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS,false)
//在序列化时忽略值为 null 的属性
mapper.setSerializationInclusion(Include.NON_NULL);
//忽略值为默认值的属性
mapper.setDefaultPropertyInclusion(Include.NON_DEFAULT);

```

Show moreShow more icon

更多配置信息可以查看 Jackson 的 DeserializationFeature，SerializationFeature 和 I nclude。

### Jackson 的 注解的使用

Jackson 根据它的默认方式序列化和反序列化 java 对象，若根据实际需要，灵活的调整它的默认方式，可以使用 Jackson 的注解。常用的注解及用法如下。

##### 表 1\. Jackson 的 常用注解

注解用法@JsonProperty用于属性，把属性的名称序列化时转换为另外一个名称。示例： @JsonProperty(“birth\_ d ate”) private Date birthDate;@JsonFormat用于属性或者方法，把属性的格式序列化时转换成指定的格式。示例： @JsonFormat(timezone = “GMT+8”, pattern = “yyyy-MM-dd HH:mm”) public Date getBirthDate()@JsonPropertyOrder用于类， 指定属性在序列化时 json 中的顺序 ， 示例： @JsonPropertyOrder({ “birth\_Date”, “name” }) public class Person@JsonCreator用于构造方法，和 @JsonProperty 配合使用，适用有参数的构造方法。 示例： @JsonCreator public Person(@JsonProperty(“name”)String name) {…}@JsonAnySetter用于属性或者方法，设置未反序列化的属性名和值作为键值存储到 map 中 @JsonAnySetter public void set(String key, Object value) { map.put(key, value); }@JsonAnyGetter用于方法 ，获取所有未序列化的属性 public Map<string, object> any() { return map; }

在了解 Jackson 的基本用法后，下面详细地介绍它的一些高阶应用。

## Jackson 的 高阶应用

### 格式处理（含日期格式）

不同类型的日期类型，Jackson 的处理方式也不同。

- 对于日期类型为 java.util.Calendar,java.util.GregorianCalendar,java.sql.Date,java.util.Date,java.sql.Timestamp，若不指定格式， 在 json 文件中将序列化 为 long 类型的数据。显然这种默认格式，可读性差，转换格式是必要的。Jackson 有 很多方式转换日期格式。
- 注解方式，请参照” Jackson 的注解的使用”的@ JsonFormat 的示例。
- ObjectMapper 方式，调用 ObjectMapper 的方法 setDateFormat，将序列化为指定格式的 string 类型的数据。
- 对于日期类型为 java.time.LocalDate，还需要添加代码 mapper.registerModule(new JavaTimeModule())，同时添加相应的依赖 jar 包

##### 清单 4 . JSR31 0 的配置信息

```
<dependency>
<groupId>com.fasterxml.jackson.datatype</groupId>
<artifactId>jackson-datatype-jsr310</artifactId>
<version>2.9.1</version>
</dependency>

```

Show moreShow more icon

对于 Jackson 2.5 以下版本，需要添加代码 objectMapper.registerModule(new JSR310Module ())

- 对于日期类型为 org.joda.time.DateTime，还需要添加代码 mapper.registerModule(new JodaModule())，同时添加相应的依赖 jar 包

##### 清单 5\. joda 的 配置信息

```

<dependency>
<groupId>com.fasterxml.jackson.datatype</groupId>
<artifactId>jackson-datatype-joda</artifactId>
<version>2.9.1</version>
</dependency>

```

Show moreShow more icon

### 泛型反序列化

Jackson 对泛型反序列化也提供很好的支持。

- 对于 List 类型 ，可以调用 constructCollectionType 方法来序列化，也可以构造 TypeReference 来序列化。

##### 清单 6 . List 泛 型使用示例

```

CollectionType javaType = mapper.getTypeFactory()
.constructCollectionType(List.class, Person.class);
List<Person> personList = mapper.readValue(jsonInString, javaType);
List<Person> personList = mapper.readValue(jsonInString, new
    TypeReference<List<Person>>(){});

```

Show moreShow more icon

- 对于 map 类型， 与 List 的实现方式相似。

##### 清单 7 . Map 泛型使用示例

```

//第二参数是 map 的 key 的类型，第三参数是 map 的 value 的类型
MapType javaType =
    mapper.getTypeFactory().constructMapType(HashMap.class,String.class,
    Person.class);
Map<String, Person> personMap = mapper.readValue(jsonInString,
    javaType);
Map<String, Person> personMap = mapper.readValue(jsonInString, new
    TypeReference<Map<String, Person>>() {});

```

Show moreShow more icon

Array 和 Collection 的处理与 List，Map 相似，这里不再详述。

### 属性可视化

是 java 对象的所有的属性都被序列化和反序列化，换言之，不是所有属性都可视化，默认的属性可视化的规则如下：

- 若该属性修饰符是 public，该属性可序列化和反序列化。
- 若属性的修饰符不是 public，但是它的 getter 方法和 setter 方法是 public，该属性可序列化和反序列化。因为 getter 方法用于序列化， 而 setter 方法用于反序列化。
- 若属性只有 public 的 setter 方法，而无 public 的 getter 方 法，该属性只能用于反序列化。

若想更改默认的属性可视化的规则，需要调用 ObjectMapper 的方法 setVisibility。

下面的示例使修饰符为 protected 的属性 name 也可以序列化和反序列化。

##### 清单 8 . 属性可视化示例

```

mapper.setVisibility(PropertyAccessor.FIELD, Visibility.ANY);
public class Person {
public int age;
protected String name;
}
PropertyAccessor 支持的类型有 ALL,CREATOR,FIELD,GETTER,IS_GETTER,NONE,SETTER
Visibility 支持的类型有 A
    NY,DEFAULT,NON_PRIVATE,NONE,PROTECTED_AND_PUBLIC,PUBLIC_ONLY

```

Show moreShow more icon

### 属性过滤

在将 Java 对象序列化为 json 时 ，有些属性需要过滤掉，不显示在 json 中 ， Jackson 有多种实现方法。

- 注解方式， 可以用 @JsonIgnore 过滤单个属性或用 @JsonIgnoreProperties 过滤多个属性，示例如下：

##### 清单 9 . 属性过滤示例一

```

@JsonIgnore
public int getAge()
@JsonIgnoreProperties(value = { "age","birth_date" })
public class Person

```

Show moreShow more icon

- addMixIn 方法加注解方式@JsonIgnoreProperties。

addMixIn 方法签名如下：

public ObjectMapper addMixIn(Class<?> target, Class<?> mixinSource);

addMixIn 方法的作用是用 mixinSource 接口或类的注解会重写 target 或 target 的子类型的注解。 用ixIn 设置

Person peixIn 的 @JsonIgnoreProperties(“name”)所重写，最终忽略的属性为 name，最终生成的 json 如下：

{“birthDate”:”2017/09/13″,”age”:40}

- SimpleBeanPropertyFilter 方式。这种方式比前两种方式更加灵活，也更复杂一些。

首先需要设置@JsonFilter 类或接口，其次设置 addMixIn，将@JsonFilter 作用于 java 对象上，最后调用 SimpleBeanPropertyFilter 的 serializeAllExcept 方法或重写 S impleBeanPropertyFilter 的 serializeAsField 方法来过滤相关属性。示例如下：

##### 清单 11 . 属性过滤示例三

```

//设置 Filter 类或接口
@JsonFilter("myFilter")
public interface MyFilter {}
//设置 addMixIn
mapper.addMixIn(Person.class, MyFilter.class);
//调用 SimpleBeanPropertyFilter 的 serializeAllExcept 方法
SimpleBeanPropertyFilter newFilter =
SimpleBeanPropertyFilter.serializeAllExcept("age");
//或重写 SimpleBeanPropertyFilter 的 serializeAsField 方法
SimpleBeanPropertyFilter newFilter = new SimpleBeanPropertyFilter() {
@Override
public void serializeAsField(Object pojo, JsonGenerator jgen,
SerializerProvider provider, PropertyWriter writer)
throws Exception {
if (!writer.getName().equals("age")) {
writer.serializeAsField(pojo, jgen, provider);
}
}
};
//设置 FilterProvider
FilterProvider filterProvider = new SimpleFilterProvider()
.addFilter("myFilter", newFilter);
mapper.setFilterProvider(filterProvider).writeValueAsString(person);

```

Show moreShow more icon

### 自定义序列化和反序列化

当 Jackson 默认序列化和反序列化的类不能满足实际需要，可以自定义新的序列化和反序列化的类。

- 自定义序列化类。自定义的序列化类需要直接或间接继承 StdSerializer 或 JsonSerializer，同时需要利用 JsonGenerator 生成 json，重写方法 serialize，示例如下：

##### 清单 12 . 自定义序列化

```

public class CustomSerializer extends StdSerializer<Person> {
@Override
public void serialize(Person person, JsonGenerator jgen,
SerializerProvider provider) throws IOException {
jgen.writeStartObject();
jgen.writeNumberField("age", person.getAge());
jgen.writeStringField("name", person.getName());
jgen.writeEndObject();
}
}

```

Show moreShow more icon

JsonGenerator 有多种 write 方法以支持生成复杂的类型的 json，比如 writeArray，writeTree 等 。若想单独创建 JsonGenerator，可以通过 JsonFactory() 的 createGenerator。

- 自定义反序列化类。自定义的反序列化类需要直接或间接继承 StdDeserializer 或 StdDeserializer，同时需要利用 JsonParser 读取 json，重写方法 deserialize，示例如下：

##### 清单 13 . 自定义序列化

```

public class CustomDeserializer extends StdDeserializer<Person> {
@Override
public Person deserialize(JsonParser jp, DeserializationContext ctxt)
throws IOException, JsonProcessingException {
JsonNode node = jp.getCodec().readTree(jp);
Person person = new Person();
int age = (Integer) ((IntNode) node.get("age")).numberValue();
String name = node.get("name").asText();
person.setAge(age);
person.setName(name);
return person;
}
}

```

Show moreShow more icon

JsonParser 提供很多方法来读取 json 信息， 如 isClosed(), nextToken(), getValueAsString()等。若想单独创建 JsonParser，可以通过 JsonFactory() 的 createParser。

- 定义好自定义序列化类和自定义反序列化类，若想在程序中调用它们，还需要注册到 ObjectMapper 的 Module，示例如下：

##### 清单 14 . 注 册 M odule 示例

```

SimpleModule module = new SimpleModule("myModule");
module.addSerializer(new CustomSerializer(Person.class));
module.addDeserializer(Person.class, new CustomDeserializer());
mapper.registerModule(module);
也可通过注解方式加在 java 对象的属性，方法或类上面来调用它们，
@JsonSerialize(using = CustomSerializer.class)
@JsonDeserialize(using = CustomDeserializer.class)
public class Person

```

Show moreShow more icon

### 树模型处理

Jackson 也提供了树模型(tree model)来生成和解析 json。若想修改或访问 json 部分属性，树模型是不错的选择。树模型由 JsonNode 节点组成。程序中常常使用 ObjectNode，ObjectNode 继承于 JsonNode，示例如下：

##### 清单 15 . ObjectNode 生成和解析 json 示例

```

ObjectMapper mapper = new ObjectMapper();
//构建 ObjectNode
ObjectNode personNode = mapper.createObjectNode();
//添加/更改属性
personNode.put("name","Tom");
personNode.put("age",40);
ObjectNode addressNode = mapper.createObjectNode();
addressNode.put("zip","000000");
addressNode.put("street","Road NanJing");
//设置子节点
personNode.set("address",addressNode);
//通过 path 查找节点
JsonNode searchNode = personNode.path("street ")；
//删除属性
((ObjectNode) personNode).remove("address");
//读取 json
JsonNode rootNode = mapper.readTree(personNode.toString());
//JsonNode 转换成 java 对象
Person person = mapper.treeToValue(personNode, Person.class);
//java 对象转换成 JsonNode
JsonNode node = mapper.valueToTree(person);

```

Show moreShow more icon

## 结束语

本文首先通过与其他 Java 的 json 的框架比较，介绍了 Jackson 的优点，并且描述了 Jackson 的 核心模块的组成，以及每个部分的作用。然后， 本文通过示例，讲解 Jackson 的基本用法，介绍了 ObjectMapper 的 write 和 read 方法，ObjectMapper 的配置信息设定，以及 jackson-annotations 包下注释的运用。最后，本文详细的介绍了 Jackson 的高阶用法，这也是本文的重点。这些高阶用法包括不同类型的日期格式处理（普通日期的类型，jdk 8 的日期类型，joda 的日期类型），List 和 Map 等泛型的反序列化，属性的可视化管理，Jackson 的 三种属性过滤方式，自定义序列化和反序列化的实现以及树模型的使用。通过本文的系统地讲解，相信读者对 Jackson 会有更深刻而全面的掌握。

## 参考资源

Jenkov.com 的 [Jackson – ObjectMapper 部分](http://tutorials.jenkov.com/java-json/jackson-objectmapper.html) ，介绍 Jackson 的 ObjectMapper 的各种 write 和 read 方法 。