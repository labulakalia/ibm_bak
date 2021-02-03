# 使用 Rest-Assured 测试 REST API
简单介绍 Rest-Assured 并通过一些示例程序展示了它的用法

**标签:** API 管理,Java,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/j-lo-rest-assured/)

王 群峰

发布: 2015-06-03

* * *

现在，越来越多的 Web 应用转向了 RESTful 的架构，很多产品和应用暴露给用户的往往就是一组 REST API，这样有一个好处，用户可以根据需要，调用不同的 API，整合出自己的应用出来。从这个角度来讲，Web 开发的成本会越来越低，人们不必再维护自己的信息孤岛，而是使用 REST API 互联互通。

那么，作为 REST API 的提供者，如何确保 API 的稳定性与正确性呢？全面系统的测试是必不可少的。Java 程序员常常借助于 JUnit 来测试自己的 REST API，不，应该这样说，Java 程序员常常借助于 JUnit 来测试 REST API 的 **实现** ！从某种角度来说，这是一种“白盒测试”，Java 程序员清楚地知道正在测试的是哪个类、哪个方法，而不是从用户的角度出发，测试的是哪个 REST API。Rest-Assured 是一套由 Java 实现的 REST API 测试框架，它是一个轻量级的 REST API 客户端，可以直接编写代码向服务器端发起 HTTP 请求，并验证返回结果；它的语法非常简洁，是一种专为测试 REST API 而设计的 DSL。使用 Rest-Assured 测试 REST API，就和真正的用户使用 REST API 一样，只不过 Rest-Assured 让这一切变得自动化了。

## 安装和使用

首先，读者需要访问 Rest-Assured 的 [官方网站](https://github.com/jayway/rest-assured) 下载最新版本的 Rest-Assured。写作此文时，作者使用的是 2.4.1，下载解压成功后，会发现软件包还包含了其他第三方依赖，这些都是运行期需要的类库，确保将它们都加入你的 classpath。

现在我们以 [豆瓣 API](https://github.com/zce/douban-api-docs) 为例，简单介绍一下 Rest-Assured 的使用方法。为此，我们需要安装一个 REST 客户端插件方便调试。如果使用 FireFox 浏览器，可安装 RESTClient 插件，安装成功后，在地址一栏输入：`http://api.douban.com/v2/book/1220562`，方法选择 GET，请求发送成功后，会返回如下 JSON 字符串：

##### 清单 1\. 豆瓣读书 API 返回结果

```
{
"id":"1220562",
"alt":"http:\/\/book.douban.com\/book\/1220562",
"rating":{"max":10, "average":"7.0", "numRaters":282, "min":0},
"author":[{"name":"片山恭一"}, {"name":"豫人"}],
"alt_title":"",
"image":"http:\/\/img1.douban.com\/spic\/s1747553.jpg",
"title":"满月之夜白鲸现",
"mobile_link":"http:\/\/m.douban.com\/book\/subject\/1220562\/",
"summary":"那一年，是听莫扎特、钓鲈鱼和家庭破裂的一年。",
"attrs":{
"publisher":["青岛出版社"],
"pubdate":["2005-01-01"],
"author":["片山恭一", "豫人"],
"price":["18.00 元"],
"title":["满月之夜白鲸现"],
"binding":["平装 (无盘)"],
"translator":["豫人"],
"pages":["180"]
},
"tags":[
{"count":106, "name":"片山恭一"},
{"count":50, "name":"日本"},
{"count":42, "name":"日本文学"},
{"count":30, "name":"满月之夜白鲸现"},
{"count":28, "name":"小说"},
{"count":10, "name":"爱情"},
{"count":7, "name":"純愛"},
{"count":6, "name":"外国文学"}
]
}

```

Show moreShow more icon

现在，我们使用 Rest-Assured 来编写一个简单的测试程序调用相同的 API：

##### 清单 2\. DouBanTest.java

```
//import 要使用的类和方法
import com.jayway.restassured.RestAssured;
import static com.jayway.restassured.RestAssured.*;
import static com.jayway.restassured.matcher.RestAssuredMatchers.*;
import static org.hamcrest.Matchers.*;

public class DouBanTest
{
@Before
public void setUP(){
//指定 URL 和端口号
RestAssured.baseURI = "http://api.douban.com/v2/book";
RestAssured.port = 80;
}

@Test
public void testGETBook()
{
get("/1220562").then().body("title", equalTo("满月之夜白鲸现"));
}

}

```

Show moreShow more icon

上面这段测试调用了同样的 API，返回了同样的结果，然后检查其 `"title"` 属性是否为 `"满月之夜白鲸现"` 。对比上面的 JSON 字符串，我们知道该测试可以通过。Rest-Assured 默认访问的 URL 为 localhost，端口号为 8080，这是为了方便开发人员做测试。如果要访问互联网上的 API，需要指定基础 URL 和端口号，如上面的代码所示。

发送请求时，还可以指定参数，比如我们使用 Java8 作为关键字，调用豆瓣读书的查询 API，可得到有关 Java8 的图书：

##### 清单 3\. 带参数的请求

```
@Test
public void testSearchBook(){
given().param("q", "java8").when().get("/search").then().body("count", equalTo(2));
}

```

Show moreShow more icon

同理，用户还可以在发送请求时指定 header、cookie、content type 等选项：

##### 清单 4\. 指定 header、cookie、content type

```
given().cookie("name", "xx").when().get("/xxx").then().body();
given().header("name", "xx").when().get("/xxx").then().body();
given().contentType("application/json").when().get("/xxx").then().body();

```

Show moreShow more icon

## 解析 JSON

使用 Rest-Assured 测试 REST API，当请求的 JSON 字符串返回时，需要选取相应的字段进行验证，Rest-Assured 提供了各种灵活的方式来选择字段，本节将通过示例，对一些常用方式进行说明。

##### 清单 5\. 常用的 JSON 字段选取方式

```
get("/1220562").then()
//取顶级属性"title”
.body("title", equalTo("满月之夜白鲸现"))
//取下一层属性
.body("rating.max", equalTo(10))
//调用数组方法
.body("tags.size()", is(8))
//取数组第一个对象的"name”属性
.body("tags[0].name", equalTo("片山恭一"))
//判断数组元素
.body("author", hasItems("[日] 片山恭一"));

```

Show moreShow more icon

字段的选取方式符合 JSON 的习惯，上面的例程分别展示了如何选取多层属性、数组元素以及调用数组的方法。

Rest-Assured 不仅支持 JSON，也支持 XML，读者可参考文末的链接自行学习。

## 身份验证

系统的某些 API 只对注册用户开放，Rest-Assured 提供了相应的验证方法：

##### 清单 6\. 验证

```
given().auth().basic(username, password).when().get("/secured").then().statusCode(200);
given().auth().oauth(..);
given().auth().oauth2(..);

```

Show moreShow more icon

## 其他 HTTP 请求

Rest-Assured 支持所有的标准 HTTP 请求：GET、POST、PUT、DELETE、PATCH、HEAD、OPTION。编写测试代码时，只需要替换为相应的方法即可。

## 集成

大家可以看到，在编写 Rest-Assured 测试代码时，我们使用了 JUnit，这样可以方便的将其集成进构建脚本中。每当有新的代码提交，系统自动触发构建，成功部署启动服务器后，就可以自动运行使用 Rest-Assured 编写的单元测试，对 REST API 进行测试，从而实现了系统测试的自动化。

## 结束语

本文对 Rest-Assured 做了简单的介绍，并通过一些示例程序展示了它的用法。需要注意的是，Rest-Assured 的功能强大，本文只是选取了它的一个子集进行介绍，但这些已经足够编写丰富的测试代码，来测试 REST API 了。读者若想进一步学习 Rest-Assured，请参考文末的链接。