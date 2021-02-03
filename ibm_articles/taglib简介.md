# taglib 简介
将 scriptlet 转换成自定义标记以获得更好的 JSP 代码

**标签:** Java,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/j-jsp07233/)

Brett McLaughlin

发布: 2003-09-28

* * *

在 _JSP 最佳实践_ 的 [上一期](https://www.ibm.com/developerworks/cn/java/j-jsp07013/)，您学习了一种基于 scriptlet 的技术，这种技术被用来将上次修改的时间戳添加到 JavaServer Page（JSP）文件中。不幸的是，比起它所提供的短期利益，scriptlet 会将更多的长期复杂性引入到您的页面中来。这些 scriptlet 会用 Java 代码将各种类型的 HTML 混杂在一起，从而使得 程序的调试和设计极其错综复杂。scriptlet 不能重用，这常常导致开发者不得不在 JSP 页面之间进行复制-粘贴操作，进而导致同一段代码出现多个版本。而且，scriptlet 还加大了错误处理的难度，因为 JSP 没有提供干净利落的方式来报告脚本错误。

因此，这次我们将设计一种新的解决方案。在本期的 _JSP 最佳实践_ 中，您将学习一些基础知识，主要是关于如何将 scriptlet 转换成自定义标记，并对其进行设置以便在您的 JSP 开发项目中使用。

## 为什么使用 taglib？

所谓 _标记库（tag library）_ ，是指由在 JSP 页面中使用的标记所组成的库。JSP 容器推出时带有一个小型的、默认的标记库。而 _自定义标记库_ 是人们为了某种特定的用途或者目的，将一些标记放到一起而形成的一种库。在一个团队中协同工作的开发者们可能会为各自的项目创建一些非常特定化的自定义标记库，同时也会创建一个通用自定义标记库，以供当前使用。

JSP 标记替代了 scriptlet，并缓解了由 scriptlet 所招致的所有令人头痛的事情。例如，您可以看到这样的标记：

```
<store:shoppingCart id="1097629"/>

```

Show moreShow more icon

或者这样的标记：

```
<tools:usageGraph />

```

Show moreShow more icon

每个标记都包含了指向一个 Java 类的引用，但是类中的代码仍然在它该在的地方：在标签之外，一个编译好的类文件之中。

## 从 scriptlet 到标记

创建一个自定义标记的第一步就是决定您想怎样使用它，如何称呼它，以及它允许使用或者需要什么属性（如果有的话）。对于时间戳标记，我们所需要的很简单：只要一个能够输出一个页面的最后修改数据的简单标记。

因为不需要属性，这个标记看上去就是这个样子：

```
<site-utils:lastModified />

```

Show moreShow more icon

这个标记的名称和前缀是一样的：都是 `site-utils` 。元素的内容是空的，这意味着该元素中不允许有子元素存在。定义了这个标记之后，接下来的一步就是实现它的行为。

## 实现行为

实现标记行为的第一步是将 scriptlet 代码从原先所在的地方移到一个 Java 类（ `LastModifiedTag` ）中，如清单 1 所示：

##### 清单 1\. 创建一个时间戳标记

```
package com.newInstance.site.tags;
import java.io.File;
import java.io.IOException;
import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.Date;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.jsp.tagext.TagSupport;
public class LastModifiedTag extends TagSupport {
    public int doEndTag() {
      try {
        HttpServletRequest request =
(HttpServletRequest)pageContext.getRequest();
        String path = pageContext.getServletContext().getRealPath(
          request.getServletPath());
        File file = new File(path);
        DateFormat formatter = DateFormat.getDateInstance(
          DateFormat.LONG);
        pageContext.getOut().println(
          formatter.format(new Date(file.lastModified())));
      } catch (IOException ignored) { }
      return EVAL_PAGE;
    }
}

```

Show moreShow more icon

这个方法中的代码看上去比较熟悉；实质上，它正是我们在先前用到的相同的时间戳代码。由于不需要用户输入，而且该标记也没有属性或者嵌入的内容，我们惟一需要关心的一个新方法就是 `doEndTag()` ()，在这个方法中该标记可以输出内容（在这个例子中是最后修改数据）到 JSP 页面。

清单 1 中其他的更改更多地与作为一个 JSP 标记的代码有关，而与在一个页面内运行的 scriptlet 没有多大关系。例如，所有的 JSP 标记都应该扩展 JSP 类 `javax.servlet.jsp.tagext.TagSupport` ，这个类为 JSP 标记提供了基本框架。可能您还注意到 ，该标记返回的 `EVAL_PAGE`. `EVAL_PAGE` 是一个预定义的整型常量，它指示容器处理页面的剩下部分。另一种选项就是使用 `SKIP_PAGE` ，它将中止对页面剩下部分的处理。如果您要将控制转移到另一个页面，例如您要前进（forward）或者重定向（redirect）用户，那么只需要使用 `SKIP_PAGE` 。剩下来的细节都是与时间戳自身有关的事情。

接下来，编译这个类，并将 LastModifiedTag.class 文件放到一个 WEB-INF/classes 目录下，注意要放到正确的路径层次结构中。这个路径应该匹配该标记的包名，包名中的圆点（.）用斜杠（/）代替。在本例中，目录的路径是基路径（WEB-INF/classes）再加上层次结构 com/newInstance/site/tags。如果有一个名为 `foo.bar.tag.MyTag` 的标记，那么它将被放在 WEB-INF/classes/foo/bar/tag 中。这种路径层次结构确保了 Web 容器在任何需要装载该标记的时候都能够找到这个类。

## 创建 TLD

接下来的一步是创建一个 _标记库描述符_ （tag library descriptor ，TLD）文件。TLD 向容器和任何要使用该标记库的 JSP 页面描述您的标记库。清单 2 显示了一个非常标准的 TLD，其中只包含了一个标记。当您将更多的标记添加到库中时，TLD 文件的长度和复杂性将随之增长。

##### 清单 2\. 一个标记库描述符

```
<?xml version="1.0" encoding="ISO-8859-1"?>
<!DOCTYPE taglib
    PUBLIC "-//Sun Microsystems, Inc.//DTD JSP Tag Library 1.2/EN"
           "http://java.sun.com/dtd/web-jsptaglibrary_1_2.dtd">
<taglib>
    <tlib-version>1.0</tlib-version>
    <jsp-version>1.2</jsp-version>
    <short-name>site-utils</short-name>
    <uri>http://www.newInstance.com/taglibs/site-utils</uri>
    <tag>
      <name>lastModified</name>
      <tag-class>com.newInstance.site.tags.LastModifiedTag</tag-class>
      <body-content>empty</body-content>
    </tag>
</taglib>

```

Show moreShow more icon

TLD 文件顶部的信息应用于整个标记库。在本例中，我提供了一个版本（这对于跟踪某个标记库的 JSP 创建者拥有哪个版本很有用，尤其是在您需要经常修改标记库的情况下更是如此）；该标记库所依赖的 JSP 版本；一个为该标记库推荐的前缀；以及用于引用这个标记库的 URI。注意，我使用了前缀 `short-name` 作为 URI 的一部分，这样 比较容易将前缀和标记库的 URI 看成一个整体。

剩下的信息用于一个特定的标记，这些信息用 `tag` 元素表示。我指定了该标记的名称、用于该标记的类（这个类应该被编译好并放在适当的地方，以便容器能够装载），最后还指定了该标记是否有嵌入的内容。在本例中，标记没有嵌入的内容，因此使用”empty”。

保存这个文件，并将其放到 WEB-INF/tlds 目录下（您可能需要在您的容器中创建这个目录）。我将这个文件保存为 site-utils.tld，并在该标记库的 URI（推荐的前缀）和 TLD 文件本身之间再次创建一个干净的链接。对于这个特定的标记库，要使其可以使用，最后一步要做的是让您的 Web 应用知道如何连接一个 JSP 页面中的 URI，如何请求使用一个标记库。这可以通过应用的 web.xml 文件来做。清单 3 显示了一个非常简单的 web.xml 片段，正是 它为我们的标记库做这样的事情。

##### 清单 3\. 将一个 URI 与一个标记库链接起来

```
<taglib>
      <taglib-uri>http://www.newInstance.com/taglibs/site-utils</taglib-uri>
      <taglib-location>/WEB-INF/tlds/site-utils.tld</taglib-location>
    </taglib>

```

Show moreShow more icon

## 包装起来

如果您已经按照上述步骤执行了，那么现在您应该能够在 JSP 页面中引用新标记了。清单 4 向我们展示了新改进的 footer.jsp，这个文件中现在完全没有 scriptlet，也没有指向具有 scriptlet 的 JSP 页面的引用。

##### 清单 4\. 使用新的标记库

```
<%@ taglib prefix="site-utils"
             uri="http://www.newInstance.com/taglibs/site-utils" %>
          </td>
          <td width="16" align="left" valign="top"> </td>
    </tr>
    <!-- End main content -->
<!-- Begin footer section -->
    <tr>
      <td width="91" align="left" valign="top" bgcolor="#330066"> </td>
      <td align="left" valign="top"> </td>
      <td class="footer" align="left" valign="top"><div align="center"><br>
        © 2003
        <a href="mailto:webmaster@newInstance.com">Brett McLaughlin</a><br>
        Last Updated: <site-utils:lastModified />
        </div></td>
      <td align="left" valign="top"> </td>
      <td width="141" align="right" valign="top" bgcolor="#330066"> </td>
    </tr>
</table>
<!-- End footer section -->

```

Show moreShow more icon

在前面几期中看了 JSTL 如何工作（参见” [利用 JSTL 更新您的 JSP 页面](https://www.ibm.com/developerworks/cn/java/j-jsp05273/) ”和” [导入内容到您的 Web 站点](https://www.ibm.com/developerworks/cn/java/j-jsp06173/) ”）之后，接下来该做什么您应该很清楚了：我们通过使用 web.xml 文件中的 URI 来引用这个标记库，为之分配一个前缀（来自 TLD 的 `short-name` 始终是最好的选择），然后就像使用任何其他 JSP 标记一样使用这个标记。最终得到的是一个简洁的、更好的 JSP 页面，这个 JSP 页面运行起来不比有 scriptlet 的时候差。

当前我们的自定义标记还很简单，在下一期，我们将扩展其功能。到时候网上见。