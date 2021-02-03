# 用 jsp:include 控制动态内容
用于构建动态网站的简单 JSP 标记

**标签:** Java,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/j-jsp04293/)

Brett McLaughlin

发布: 2003-07-06

* * *

在新的 _JSP 最佳实践_ 系列的前一篇文章中，您了解了如何使用 JSP `include` 伪指令将诸如页眉、页脚和导航组件之类的静态内容包含到 Web 页面中。和服务器端包含一样，JSP `include` 伪指令允许某个页面从另一个页面提取内容或数据。清单 1 重温了 `include` 伪指令。

##### 清单 1\. JSP include 伪指令

```
<![CDATA[
<%@ page language="java" contentType="text/html" %>
<html>
<head>
      <title>newInstance.com</title>
      <meta http-equiv="Content-Type"
        content="text/html; charset=iso-8859-1" />
      <link href="/styles/default.css"
        rel="stylesheet" type="text/css" />
</head>
<body>
<%@ include file="header.jsp" %>
<%@ include file="navigation.jsp" %>
<%@ include file="bookshelf.jsp" %>
<%@ include file="/mt-blogs/index.jsp" %>
<%@ include file="footer.jsp" %>
</body>
</html>
]]>

```

Show moreShow more icon

##### 您需要什么

本系列中的所有最佳实践都基于 JavaServer Pages 技术。要运行其中的任何最佳实践，都需要在您的本地机器或测试服务器上安装符合 JSP 技术的 Web 容器。您还需要使用文本编辑器或 IDE 来对 JSP 页面编码。

虽然 `include` 非常适于将静态内容并入 Web 页面，但对于动态内容却不尽如人意。我们在前一篇文章中在试图重新装入高速缓存文件时发现了这一问题。与大多数页眉文件及页脚文件不同，动态内容变化频繁，必须时刻更新。我们将首先扼要地重述一下 `include` 伪指令的局限性，然后我将向您演示如何用 `jsp:include` 标记来扩展 JSP 的包含能力。

## 高速缓存问题

JSP `include` 伪指令的不足之处有一个是：它会导致 Web 浏览器高速缓存所有页面。在处理诸如页脚、版权声明或一组静态链接之类的静态组件时，这是有意义的。这些文件不会改变，因此没有理由让 JSP 解释器不断地重新轮询其中的数据。凡是可能的地方，都应该实现高速缓存，因为它改善了应用程序的性能。

##### JSP 测试和开发

在构建 Web 应用程序或网站时，可能需要大量更新页眉、页脚和导航链接。仅仅为了看到对所包含文件所做的更改，而被迫不断地关闭浏览器或清除其高速缓存，这可能是件痛苦的事情。另一方面，为了结束开发周期，而不得不彻底检查一遍并修改数百个使用了 `include` 伪指令的页面，这也是一件痛苦的事情。我的建议是，在测试期间禁用浏览器高速缓存。在大多数情形下，这样做能够彻底解决问题。也有极少数情形，这样做并不奏效，这时可以在浏览器或服务器上不断地重新启动 Web 容器来确保不进行高速缓存。

但是，有时侯，进行高速缓存会得不偿失。如果提入的内容来自使用动态数据（如 Weblog 或数据库驱动的 JSP 文件）的程序，甚至如果所包含的内容是经常变化的 HTML（如时间戳记），那么每当装入 Web 页面时，都需要显示这些文件或程序的最新版本。遗憾的是，JSP `include` 伪指令并不具备这一功能。在测试和开发周期（请参阅侧栏”JSP 测试和开发”）中，在浏览器中禁用高速缓存通常能够解决这一问题。但是，对于实际使用的应用程序而言，性能是任何设计决策过程中的一项重要因素，禁用高速缓存并不是一种可行的长远之计。更好的解决方案是使用 `jsp:include` 标记。

## jsp:include 标记

`jsp:include` 只不过是一个不同于 `include` 的伪指令而已。 `jsp:include` 的优点在于：它 _总是_ 会检查所含文件中的变化。过一会儿我们将研究这一新标记的工作方式。但首先看一下两种 `include` 各自的代码，以便能够看到二者之间的异同。

清单 2 显示了一个简单页面，它使用了原始的 JSP `include` 伪指令。

##### 清单 2\. JSP include 伪指令

```
<![CDATA[
<%@ page language="java" contentType="text/html" %>
<html>
     <head>
      <title>JSP include element test</title>
     </head>
     <body>
      This content is statically in the main JSP file.<br />
      <%@ include file="included.html" %>
     </body>
</html>
]]>

```

Show moreShow more icon

清单 3 是同一个页面，只不过这里转成使用 `jsp:include` 标记。

##### 清单 3\. 转成使用 jsp:include

```
<![CDATA[
<%@ page language="java" contentType="text/html" %>
<html>
     <head>
      <title>JSP include element test</title>
     </head>
     <body>
      This content is statically in the main JSP file.<br />
      <jsp:include page="included.html" flush="true" />
     </body>
</html>
]]>

```

Show moreShow more icon

您应该注意这两种代码类型之间的两大区别。首先， `jsp:include` 元素不使用属于 `include` 伪指令的 `%@` 语法。实际上， `jsp` 前缀让 JSP 编译器知道：它应该寻找标准 JSP 标记集中的元素。其次，指定要包含的文件的属性从 `file` 变成了 `page` 。如果愿意，可以自己测试一下新标记的结果。只需更改上一篇文章（请参阅 参考资料 ）中 `included.html` 文件的内容，然后重新装入浏览器页面，就会立即看到新内容。

##### flush 属性

您可能已注意到 `jsp:include` 代码示例中的 `flush` 属性。顾名思义， `flush` 指示在读入包含内容之前是否清空任何现有的缓冲区。JSP 1.1 中需要 `flush` 属性，因此，如果代码中不用它，会得到一个错误。但是，在 JSP 1.2 中， `flush` 属性缺省为 false。由于清空大多数时候不是一个重要的问题，因此，我的建议是：对于 JSP 1.1，将 `flush` 设置为 true；而对于 JSP 1.2 及更高版本，将其设置为关闭。

## jsp:include 是如何工作的

如果您有点爱刨根问底，那么可能十分想知道 `jsp:include` 标记的行为为什么与 `include` 伪指令不同。道理其实十分简单： `jsp:include` 包含的是所包含 URI 的 _响应_ ，而不是 URI 本身。这意味着：对所指出的 URI 进行 _解释_ ，因而包含的是 _生成的响应_ 。如果页面是 HTML，那么将得到一点也没有变化的 HTML。但是，如果是 Perl 脚本、Java servlet 或者 CGI 程序，那么得到的将是从该程序解释而得的结果。虽然页面通常就是 HTML，但实际程序恰好是达到目的的手段。而且，由于每次请求页面的时候都会进行解释，因此从来不会象使用 `include` 伪指令时那样高速缓存结果。虽然这只是很小的变动，但它却导致了您所见到的行为中的全部差异。

## 一种混合搭配的解决方案

`include` 伪指令在某些网站上有其用武之地。例如，如果站点包含一些（如果有变化，也很少）几乎没有变化的页眉、页脚和导航文件，那么基本的 `include` 伪指令是这些组件的最佳选项。由于 `include` 伪指令采用了高速缓存，因此只需放入包含文件一次，其内容就会被高速缓存，其结果会是极大地提高了站点的性能。

然而，对于现在许多 Web 应用程序或站点而言，地毯式的高速缓存并不能解决问题。虽然页眉和页脚可能是静态的，但是不可能整个站点都是静态的。例如，从数据库提取导航链接是很常见的，并且许多基于 JSP 技术的站点还从其它站点或应用程序上的动态 JSP 页面提取内容。如果正在处理动态内容，那么需要采用 `jsp:include` 来处理该内容。

当然，最好的解决方案是经常把这两种方法混合搭配使用，将每种构造用到最恰当的地方。清单 4 是混合搭配包含解决方案的一个示例。

##### 清单 4\. 混合搭配解决方案

```
<![CDATA[
<%@ page language="java" contentType="text/html" %>
<html>
<head>
<title>newInstance.com</title>
<meta http-equiv="Content-Type" content="text/html; charset=iso-8859-1" />
<link href="/styles/default.css" rel="stylesheet" type="text/css" />
</head>
<body>
<jsp:include page="header.jsp" flush="true">
<jsp:param name="pageTitle" value="newInstance.com"/>
<jsp:param name="pageSlogan" value=" " />
</jsp:include>
<%@ include file="/navigation.jsp" %>
<jsp:include page="bookshelf.jsp" flush="true" />
<jsp:include page="/mt-blogs/index.jsp" flush="true" />
<%@ include file="/footer.jsp" %>
</body>
</html>
]]>

```

Show moreShow more icon

上面的代码显示了前一篇文章中的示例索引页面。导航链接和页脚是静态内容，一年最多更改一次。对于这些文件，我使用了 `include` 伪指令。内容窗格包含 Weblog 和”bookshelf”组件，它们是动态生成的。这两个组件需要一直更新，因此对它们，我使用了 `jsp:include` 标记。 `header.jsp` 文件有点奇怪。这个组件是从另一个本质上是静态的 JSP 页面提取的。但是，正如您将注意到的那样，它从包含页提取页”标语”，然后将它显示出来。要处理这一共享信息，我们必须向页眉文件传入参数。而要处理那些参数，就必须使用 `jsp:include` 元素。

如果您想了解有关那些参数的内容，请放心，后续文章不久就将能满足您的需要。在下一篇文章中，我将解释 JSP 参数以及它们是如何与 JavaBeans 组件交互的。届时，我们网上见。

## 相关主题

- [Apache Tomcat](http://jakarta.apache.org/tomcat)
- [jEdit](http://www.jedit.org)
- [NetBeans](https://netbeans.org)
- [Eclipse](http://www.eclipse.org)
- [JSP 规范](http://jcp.org/aboutJava/communityprocess/first/jsr152/index2.html)
- [JavaServer Pages](http://www.oreilly.com/catalog/jserverpages2/)