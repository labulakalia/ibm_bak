# 使用 REST-Assured 测试 REST API 的进阶技巧和最佳实践
使用 Rest-Assured 和 JSON Schema 测试 REST API 的方法及技巧

**标签:** API 管理,Java,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/j-lo-rest-assured2/)

吴斌

发布: 2015-10-20

* * *

## REST API 的测试要点

随着 Web 时代的发展，无论是互联网网站还是企业应用，都开始或者已经公布了自己的 REST API，API 的应用的与集成也越来越广泛，因此 API 的测试也越来越受到重视。各种针对 REST API 的测试工具也应运而生，《 [使用 Rest-Assured 测试 REST API](http://www.ibm.com/developerworks/cn/java/j-lo-rest-assured/index.html) 》已进行了初步的介绍。

REST API 的测试有其自己的特点，虽然测试执行很快，很适合自动化测试，但是通常参数特别多，请求体和返回体有时也很复杂。从本质上说，REST API 的测试主要是测试 HTTP 的 GET/POST/DELETE/PUT 这几个方法。其中，最复杂的主要是 GET 和 POST/PUT 两种情况。GET 方法主要测试返回的 xml 或 JSON 返回体。返回体的属性及内容越多，测试就越复杂。

POST/PUT 方法主要测试发送过去的请求体是否能通过验证，并且是否根据请求体创建或修改相应的内容。这里的难点是请求体的复杂性，以及需要测试请求体中 property 的各种参数组合。这时候，使用 JSON scheme 来验证返回体能大大简化测试代码。

## REST-assured 的测试实践

REST-assured 是一套测试框架，本质上就是一组 Jar 包，测试人员可以使用其中的各种 API 来实现自己的测试目的。它的安装和简单的使用本文就不再赘述，请参考《 [使用 Rest-Assured 测试 REST API](http://www.ibm.com/developerworks/cn/java/j-lo-rest-assured/index.html) 》。

我们首先看前面提到的第一个复杂点–验证返回体。JSON 返回体因为其结构简单，非常常用。在返回体中可能有着十几或者几十个 property，每个 property 的类型不同，取值范围也不同。

Rest-Assured 可以直接在 GET 的时候，同时进行验证。如下例子：

```
Get(url).then().body("server.name”,equalTo("apache”));

```

Show moreShow more icon

如果有很多个属性都需要验证，则可以使用 from(body) 方法来从返回体中获取到具体某个属性，之后进行格式和内容的验证。如：

```
assertEquals(from(body).getInt("errorCode"),400);

```

Show moreShow more icon

from 使用相当灵活，既可以做验证，也可以用来获取 body 中的某一些值做为中间值来计算，或者用来做后续的验证。

如果返回体是一个数组，还可以用 from 来获取数组中的每一个对象来分别做验证。

例如下面这段代码：

```
List<HashMap> aList = from(body).getList("", HashMap.class);

```

Show moreShow more icon

将返回体中的数组转型成 hashmap 组成的一个列表。每一个 JSON 对象都成为了一个 hashmap 对象，我们就可以方便地在循环中获得其中具体的值做验证。

## JSON Schema

然而，如果返回体非常庞大，属性非常多，这样的话，一个个参数的去验证，测试用例会非常多，代码也会很冗长。

这时候如果我们使用 JSON schema 去验证的话，就会大大减少用例和代码数量。

JSON schema 描述了 JSON 的数据格式，是一种元数据，它非常简单易读，我们先来看一个例子：

```
{
"type": "object",
"required”: true,
"properties": {
"name": {
"type": "string"
"required”: true,
},
"badgeNumber": {
"type": "string"
"pattern": "^[0-9].*$",
"required”: true,
},
"isActive”: {
"type”:"string”,
"enum": [ "false","true" ],
"required”: true
}
"age": {
"description": "Age in years",
"type": "integer",
"minimum": 1
}
},
}

```

Show moreShow more icon

从这个例子我们可以看出，JSON schema 本身就是 JSON 格式的。

它所描述的这个 JSON 对象，有 4 个属性，name, badgeNumber,isActive 和 age。另外 type 还描述了每一个属性的类型，除了 age 为整数型，其余均为字符串型。required 表示该属性是否是必须的。这个例子中，除了 age 外，其他属性是必须的。

对于整数型，我们还可以限制其取值范围，例如在上面这个例子中，我们使用 minimum=1，将 age 的最小值限制为为 1。

对于字符串类型，我们更可以用正则表达式来做更具体的描述。

例如上例中的 badgeNumber，我们限定了这个字符串必须以数字开头。

在 isActive 属性中我们用枚举的方式，限定了取值只能为 false 或者 true。

## 生成 JSON Schema

对于简单的 JSON 返回体，我们可以根据需求来自己创建 JSNO Schema，但是对于复杂的返回体，这个过程也挺累人的。为了方便起见，我可以用 JSON Schema 生成工具。例如： [http://JSONschema.net/#/home](http://JSONschema.net/#/home) 。

通常我们可以先用任何方式（如测试代码或者 REST Client 等插件）得到一个需要测试的返回体，然后用自动生成工具生成一个 schema 模板。

一般来说，生成的 schema 模板会列出所有的属性及其类型。

然后在这个 schema 基础上我们来分析每个属性，根据不同的类型加上必要的限制条件。每种限制条件都相当于测试用例中的一个验证点。像上面这个例子中的 badgeNumber，如果在返回体中这个属性的值如果是整数型，就能使测试失败，如果这个属性的值以字母开头，同样会使测试失败。

## 使用 REST-Assured 验证 JSON Schema

首先我们需要安装 JSON-schema-validator，在 [https://github.com/fge/JSON-schema-validator](https://github.com/fge/json-schema-validator) 上下载 JSON-schema-validator 的 lib 包，将其添加到我们的 classpath 中。

其次我们还需要将之前生成的 JSON schema 文件添加到我们的 classpath 中。

然后，我们就能在测试代码中仅用一句代码验证返回体是否符合指定的 JSON Schema

例如：

```
expect().statusCode(200).given().auth().preemptive().basic(user, user).
headers("Accept","application/JSON").when().get("http://xyz.com/abc/").
then().assertThat().body(matchesJSONSchemaInClasspath("abc.JSON"));

```

Show moreShow more icon

其中 abc.JSON 文件就是我们准备 JSON schema 文件。

**REST API 其他测试技巧**

1. REST API 测试中经常需要对于返回体中的部分元素进行验证。我们可以利用 org.hamcrest.MatcherAssert 提供的各种断言来帮助简化测试代码。REST-assured 可以与 org.hamcrest.MatcherAssert 一起使用，进行很多方便而有意思的验证。

     如下例子：





    ```
    //验证 index 的值不大于 20
    given().auth().preemptive().basic(user, user).headers("Accept","application/json").
    get("http://xyz.com/abc/").then().assertThat().body("index",lessThanOrEqualTo(20));
    //验证 data.items 数组的 size 为 3
    given().auth().preemptive().basic(user, user).headers("Accept","application/json").
    get("http://xyz.com/abc/").then().assertThat().body("data.items",hasSize(3));
    //验证 data.items 中每个元素的 id 属性值都大于 5
    given().auth().preemptive().basic(user, user).headers("Accept","application/json").
    get("http://xyz.com/abc/").then().assertThat().body("data.items.id", everyItem(greaterThan(5)));

    ```





    Show moreShow more icon

2. REST-assured 也支持类似于 Ruby block 的方来进行搜索验证某些属性。

     例如：





    ```
    getInt("data.transferQualities.find {it.status=='SUCCESS'}.count");

    ```





    Show moreShow more icon

     在这段代码中，data.transferQualities 是一个 array，这段代码可以在其中找寻某元素，它的 status 的属性为 success，然后得到它的 count 属性。

3. 在测试 POST/PUT 方法时，最麻烦的地方是请求体中 property 非常多，而且各自有不同的限制条件。为了测试非法的输入能正确被系统识别出来，要测试很多参数组合。我们可以使用 Combinatorial Testing（又称 All-pairs testing）的方法来得到参数组合，然后使用 Rest-Assured 进行测试。如果需要了解关于 Combinatorial Testing 的更详细信息，可以从 [维基百科上的资料](https://en.wikipedia.org/wiki/All-pairs_testing) 开始了解。


## 结束语

本文介绍了如何使用 Rest-Assured 和 JSON Schema 测试 REST API 的方法及其他技巧。由于笔者水平有限，如果文章中存在错误，欢迎读者联系并进行指正，也欢迎读者一起分享经验与想法。

## 相关主题：

- [Rest Assured 官网](http://rest-assured.io/)
- [Rest-Assured GitHub 地址](https://github.com/jayway/rest-assured)
- [使用 Rest-Assured 测试 REST API](https://www.ibm.com/developerworks/cn/java/j-lo-rest-assured/index.html)
- [JSON-Schema-validator 官方网站](https://github.com/fge/JSON-schema-validator)