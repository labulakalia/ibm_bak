# 文件资源操作和 Web 相关工具类
了解 Spring 中有哪些好用的工具类并在程序编写时适当使用

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-lo-spring-utils1/)

陈雄华

发布: 2007-08-17

* * *

## 文件资源操作

文件资源的操作是应用程序中常见的功能，如当上传一个文件后将其保存在特定目录下，从指定地址加载一个配置文件等等。我们一般使用 JDK 的 I/O 处理类完成这些操作，但对于一般的应用程序来说，JDK 的这些操作类所提供的方法过于底层，直接使用它们进行文件操作不但程序编写复杂而且容易产生错误。相比于 JDK 的 File，Spring 的 Resource 接口（资源概念的描述接口）抽象层面更高且涵盖面更广，Spring 提供了许多方便易用的资源操作工具类，它们大大降低资源操作的复杂度，同时具有更强的普适性。这些工具类不依赖于 Spring 容器，这意味着您可以在程序中象一般普通类一样使用它们。

### 加载文件资源

Spring 定义了一个 org.springframework.core.io.Resource 接口，Resource 接口是为了统一各种类型不同的资源而定义的，Spring 提供了若干 Resource 接口的实现类，这些实现类可以轻松地加载不同类型的底层资源，并提供了获取文件名、URL 地址以及资源内容的操作方法。

**访问文件资源**

假设有一个文件地位于 Web 应用的类路径下，您可以通过以下方式对这个文件资源进行访问：

- 通过 FileSystemResource 以文件系统绝对路径的方式进行访问；
- 通过 ClassPathResource 以类路径的方式进行访问；
- 通过 ServletContextResource 以相对于 Web 应用根目录的方式进行访问。

相比于通过 JDK 的 File 类访问文件资源的方式，Spring 的 Resource 实现类无疑提供了更加灵活的操作方式，您可以根据情况选择适合的 Resource 实现类访问资源。下面，我们分别通过 FileSystemResource 和 ClassPathResource 访问同一个文件资源：

##### 清单 1\. FileSourceExample

```
package com.baobaotao.io;
import java.io.IOException;
import java.io.InputStream;
import org.springframework.core.io.ClassPathResource;
import org.springframework.core.io.FileSystemResource;
import org.springframework.core.io.Resource;
public class FileSourceExample {
    public static void main(String[] args) {
        try {
            String filePath =
            "D:/masterSpring/chapter23/webapp/WEB-INF/classes/conf/file1.txt";
            // ① 使用系统文件路径方式加载文件
            Resource res1 = new FileSystemResource(filePath);
            // ② 使用类路径方式加载文件
            Resource res2 = new ClassPathResource("conf/file1.txt");
            InputStream ins1 = res1.getInputStream();
            InputStream ins2 = res2.getInputStream();
            System.out.println("res1:"+res1.getFilename());
            System.out.println("res2:"+res2.getFilename());
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}

```

Show moreShow more icon

在获取资源后，您就可以通过 Resource 接口定义的多个方法访问文件的数据和其它的信息：如您可以通过 getFileName() 获取文件名，通过 getFile() 获取资源对应的 File 对象，通过 getInputStream() 直接获取文件的输入流。此外，您还可以通过 createRelative(String relativePath) 在资源相对地址上创建新的资源。

在 Web 应用中，您还可以通过 ServletContextResource 以相对于 Web 应用根目录的方式访问文件资源，如下所示：

```
<%@ page language="java" contentType="text/html; charset=utf-8" pageEncoding="utf-8"%>
<jsp:directive.page import="
    org.springframework.web.context.support.ServletContextResource"/>
<jsp:directive.page import="org.springframework.core.io.Resource"/>
<%
    // ① 注意文件资源地址以相对于 Web 应用根路径的方式表示
    Resource res3 = new ServletContextResource(application,
        "/WEB-INF/classes/conf/file1.txt");
    out.print(res3.getFilename());
%>

```

Show moreShow more icon

对于位于远程服务器（Web 服务器或 FTP 服务器）的文件资源，您则可以方便地通过 UrlResource 进行访问。

为了方便访问不同类型的资源，您必须使用相应的 Resource 实现类，是否可以在不显式使用 Resource 实现类的情况下，仅根据带特殊前缀的资源地址直接加载文件资源呢？ Spring 提供了一个 ResourceUtils 工具类，它支持”classpath:”和”file:”的地址前缀，它能够从指定的地址加载文件资源，请看下面的例子：

##### 清单 2\. ResourceUtilsExample

```
package com.baobaotao.io;
import java.io.File;
import org.springframework.util.ResourceUtils;
public class ResourceUtilsExample {
    public static void main(String[] args) throws Throwable{
        File clsFile = ResourceUtils.getFile("classpath:conf/file1.txt");
        System.out.println(clsFile.isFile());

        String httpFilePath = "file:D:/masterSpring/chapter23/src/conf/file1.txt";
        File httpFile = ResourceUtils.getFile(httpFilePath);
        System.out.println(httpFile.isFile());
    }
}

```

Show moreShow more icon

ResourceUtils 的 getFile(String resourceLocation) 方法支持带特殊前缀的资源地址，这样，我们就可以在不和 Resource 实现类打交道的情况下使用 Spring 文件资源加载的功能了。

**本地化文件资源**

本地化文件资源是一组通过本地化标识名进行特殊命名的文件，Spring 提供的 LocalizedResourceHelper 允许通过文件资源基名和本地化实体获取匹配的本地化文件资源并以 Resource 对象返回。假设在类路径的 i18n 目录下，拥有一组基名为 message 的本地化文件资源，我们通过以下实例演示获取对应中国大陆和美国的本地化文件资源：

##### 清单 3\. LocaleResourceTest

```
package com.baobaotao.io;
import java.util.Locale;
import org.springframework.core.io.Resource;
import org.springframework.core.io.support.LocalizedResourceHelper;
public class LocaleResourceTest {
    public static void main(String[] args) {
        LocalizedResourceHelper lrHalper = new LocalizedResourceHelper();
        // ① 获取对应美国的本地化文件资源
        Resource msg_us = lrHalper.findLocalizedResource("i18n/message", ".properties",
        Locale.US);
        // ② 获取对应中国大陆的本地化文件资源
        Resource msg_cn = lrHalper.findLocalizedResource("i18n/message", ".properties",
        Locale.CHINA);
        System.out.println("fileName(us):"+msg_us.getFilename());
        System.out.println("fileName(cn):"+msg_cn.getFilename());
    }
}

```

Show moreShow more icon

虽然 JDK 的 java.util.ResourceBundle 类也可以通过相似的方式获取本地化文件资源，但是其返回的是 ResourceBundle 类型的对象。如果您决定统一使用 Spring 的 Resource 接表征文件资源，那么 LocalizedResourceHelper 就是获取文件资源的非常适合的帮助类了。

### 文件操作

在使用各种 Resource 接口的实现类加载文件资源后，经常需要对文件资源进行读取、拷贝、转存等不同类型的操作。您可以通过 Resource 接口所提供了方法完成这些功能，不过在大多数情况下，通过 Spring 为 Resource 所配备的工具类完成文件资源的操作将更加方便。

**文件内容拷贝**

第一个我们要认识的是 FileCopyUtils，它提供了许多一步式的静态操作方法，能够将文件内容拷贝到一个目标 byte[]、String 甚至一个输出流或输出文件中。下面的实例展示了 FileCopyUtils 具体使用方法：

##### 清单 4\. FileCopyUtilsExample

```
package com.baobaotao.io;
import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.FileReader;
import java.io.OutputStream;
import org.springframework.core.io.ClassPathResource;
import org.springframework.core.io.Resource;
import org.springframework.util.FileCopyUtils;
public class FileCopyUtilsExample {
    public static void main(String[] args) throws Throwable {
        Resource res = new ClassPathResource("conf/file1.txt");
        // ① 将文件内容拷贝到一个 byte[] 中
        byte[] fileData = FileCopyUtils.copyToByteArray(res.getFile());
        // ② 将文件内容拷贝到一个 String 中
        String fileStr = FileCopyUtils.copyToString(new FileReader(res.getFile()));
        // ③ 将文件内容拷贝到另一个目标文件
        FileCopyUtils.copy(res.getFile(),
        new File(res.getFile().getParent()+ "/file2.txt"));

        // ④ 将文件内容拷贝到一个输出流中
        OutputStream os = new ByteArrayOutputStream();
        FileCopyUtils.copy(res.getInputStream(), os);
    }
}

```

Show moreShow more icon

往往我们都通过直接操作 InputStream 读取文件的内容，但是流操作的代码是比较底层的，代码的面向对象性并不强。通过 FileCopyUtils 读取和拷贝文件内容易于操作且相当直观。如在 ① 处，我们通过 FileCopyUtils 的 copyToByteArray(File in) 方法就可以直接将文件内容读到一个 byte[] 中；另一个可用的方法是 copyToByteArray(InputStream in)，它将输入流读取到一个 byte[] 中。

如果是文本文件，您可能希望将文件内容读取到 String 中，此时您可以使用 copyToString(Reader in) 方法，如 ② 所示。使用 FileReader 对 File 进行封装，或使用 InputStreamReader 对 InputStream 进行封装就可以了。

FileCopyUtils 还提供了多个将文件内容拷贝到各种目标对象中的方法，这些方法包括：

方法说明`static void copy(byte[] in, File out)`将 byte[] 拷贝到一个文件中`static void copy(byte[] in, OutputStream out)`将 byte[] 拷贝到一个输出流中`static int copy(File in, File out)`将文件拷贝到另一个文件中`static int copy(InputStream in, OutputStream out)`将输入流拷贝到输出流中`static int copy(Reader in, Writer out)`将 Reader 读取的内容拷贝到 Writer 指向目标输出中`static void copy(String in, Writer out)`将字符串拷贝到一个 Writer 指向的目标中

在实例中，我们虽然使用 Resource 加载文件资源，但 FileCopyUtils 本身和 Resource 没有任何关系，您完全可以在基于 JDK I/O API 的程序中使用这个工具类。

**属性文件操作**

我们知道可以通过 java.util.Properties 的 load(InputStream inStream) 方法从一个输入流中加载属性资源。Spring 提供的 PropertiesLoaderUtils 允许您直接通过基于类路径的文件地址加载属性资源，请看下面的例子：

```
package com.baobaotao.io;
import java.util.Properties;
import org.springframework.core.io.support.PropertiesLoaderUtils;
public class PropertiesLoaderUtilsExample {
    public static void main(String[] args) throws Throwable {
        // ① jdbc.properties 是位于类路径下的文件
        Properties props = PropertiesLoaderUtils.loadAllProperties("jdbc.properties");
        System.out.println(props.getProperty("jdbc.driverClassName"));
    }
}

```

Show moreShow more icon

一般情况下，应用程序的属性文件都放置在类路径下，所以 PropertiesLoaderUtils 比之于 Properties#load(InputStream inStream) 方法显然具有更强的实用性。此外，PropertiesLoaderUtils 还可以直接从 Resource 对象中加载属性资源：

方法说明`static Properties loadProperties(Resource resource)`从 Resource 中加载属性`static void fillProperties(Properties props, Resource resource)`将 Resource 中的属性数据添加到一个已经存在的 Properties 对象中

**特殊编码的资源**

当您使用 Resource 实现类加载文件资源时，它默认采用操作系统的编码格式。如果文件资源采用了特殊的编码格式（如 UTF-8），则在读取资源内容时必须事先通过 EncodedResource 指定编码格式，否则将会产生中文乱码的问题。

##### 清单 5\. EncodedResourceExample

```
package com.baobaotao.io;
import org.springframework.core.io.ClassPathResource;
import org.springframework.core.io.Resource;
import org.springframework.core.io.support.EncodedResource;
import org.springframework.util.FileCopyUtils;
public class EncodedResourceExample {
        public static void main(String[] args) throws Throwable  {
            Resource res = new ClassPathResource("conf/file1.txt");
            // ① 指定文件资源对应的编码格式（UTF-8）
            EncodedResource encRes = new EncodedResource(res,"UTF-8");
            // ② 这样才能正确读取文件的内容，而不会出现乱码
            String content  = FileCopyUtils.copyToString(encRes.getReader());
            System.out.println(content);
    }
}

```

Show moreShow more icon

EncodedResource 拥有一个 getResource() 方法获取 Resource，但该方法返回的是通过构造函数传入的原 Resource 对象，所以必须通过 EncodedResource#getReader() 获取应用编码后的 Reader 对象，然后再通过该 Reader 读取文件的内容。

## Web 相关工具类

您几乎总是使用 Spring 框架开发 Web 的应用，Spring 为 Web 应用提供了很多有用的工具类，这些工具类可以给您的程序开发带来很多便利。在这节里，我们将逐一介绍这些工具类的使用方法。

### 操作 Servlet API 的工具类

当您在控制器、JSP 页面中想直接访问 Spring 容器时，您必须事先获取 WebApplicationContext 对象。Spring 容器在启动时将 WebApplicationContext 保存在 ServletContext 的属性列表中，通过 WebApplicationContextUtils 工具类可以方便地获取 WebApplicationContext 对象。

**WebApplicationContextUtils**

当 Web 应用集成 Spring 容器后，代表 Spring 容器的 WebApplicationContext 对象将以 WebApplicationContext.ROOT\_WEB\_APPLICATION\_CONTEXT\_ATTRIBUTE 为键存放在 ServletContext 属性列表中。您当然可以直接通过以下语句获取 WebApplicationContext：

```
WebApplicationContext wac = (WebApplicationContext)servletContext.
getAttribute(WebApplicationContext.ROOT_WEB_APPLICATION_CONTEXT_ATTRIBUTE);

```

Show moreShow more icon

但通过位于 org.springframework.web.context.support 包中的 WebApplicationContextUtils 工具类获取 WebApplicationContext 更方便：

```
WebApplicationContext wac =WebApplicationContextUtils.
getWebApplicationContext(servletContext);

```

Show moreShow more icon

当 ServletContext 属性列表中不存在 WebApplicationContext 时，getWebApplicationContext() 方法不会抛出异常，它简单地返回 null。如果后续代码直接访问返回的结果将引发一个 NullPointerException 异常，而 WebApplicationContextUtils 另一个 getRequiredWebApplicationContext(ServletContext sc) 方法要求 ServletContext 属性列表中一定要包含一个有效的 WebApplicationContext 对象，否则马上抛出一个 IllegalStateException 异常。我们推荐使用后者，因为它能提前发现错误的时间，强制开发者搭建好必备的基础设施。

**WebUtils**

位于 org.springframework.web.util 包中的 WebUtils 是一个非常好用的工具类，它对很多 Servlet API 提供了易用的代理方法，降低了访问 Servlet API 的复杂度，可以将其看成是常用 Servlet API 方法的门面类。

下面这些方法为访问 HttpServletRequest 和 HttpSession 中的对象和属性带来了方便：

方法说明`Cookie getCookie(HttpServletRequest request, String name)`获取 HttpServletRequest 中特定名字的 Cookie 对象。如果您需要创建 Cookie， Spring 也提供了一个方便的 CookieGenerator 工具类；`Object getSessionAttribute(HttpServletRequest request, String name)`获取 HttpSession 特定属性名的对象，否则您必须通过 request.getHttpSession.getAttribute(name) 完成相同的操作；`Object getRequiredSessionAttribute(HttpServletRequest request, String name)`和上一个方法类似，只不过强制要求 HttpSession 中拥有指定的属性，否则抛出异常；`String getSessionId(HttpServletRequest request)`获取 Session ID 的值；`void exposeRequestAttributes(ServletRequest request, Map attributes)`将 Map 元素添加到 ServletRequest 的属性列表中，当请求被导向（forward）到下一个处理程序时，这些请求属性就可以被访问到了；

此外，WebUtils 还提供了一些和 ServletContext 相关的方便方法：

方法说明`String getRealPath(ServletContext servletContext, String path)`获取相对路径对应文件系统的真实文件路径；`File getTempDir(ServletContext servletContext)`获取 ServletContex 对应的临时文件地址，它以 File 对象的形式返回。

下面的片断演示了使用 WebUtils 从 HttpSession 中获取属性对象的操作：

```
protected Object formBackingObject(HttpServletRequest request) throws Exception {
    UserSession userSession = (UserSession) WebUtils.getSessionAttribute(request,
        "userSession");
    if (userSession != null) {
        return new AccountForm(this.petStore.getAccount(
        userSession.getAccount().getUsername()));
    } else {
        return new AccountForm();
    }
}

```

Show moreShow more icon

### Spring 所提供的过滤器和监听器

Spring 为 Web 应用提供了几个过滤器和监听器，在适合的时间使用它们，可以解决一些常见的 Web 应用问题。

**延迟加载过滤器**

Hibernate 允许对关联对象、属性进行延迟加载，但是必须保证延迟加载的操作限于同一个 Hibernate Session 范围之内进行。如果 Service 层返回一个启用了延迟加载功能的领域对象给 Web 层，当 Web 层访问到那些需要延迟加载的数据时，由于加载领域对象的 Hibernate Session 已经关闭，这些导致延迟加载数据的访问异常。

Spring 为此专门提供了一个 OpenSessionInViewFilter 过滤器，它的主要功能是使每个请求过程绑定一个 Hibernate Session，即使最初的事务已经完成了，也可以在 Web 层进行延迟加载的操作。

OpenSessionInViewFilter 过滤器将 Hibernate Session 绑定到请求线程中，它将自动被 Spring 的事务管理器探测到。所以 OpenSessionInViewFilter 适用于 Service 层使用 HibernateTransactionManager 或 JtaTransactionManager 进行事务管理的环境，也可以用于非事务只读的数据操作中。

要启用这个过滤器，必须在 web.xml 中对此进行配置：

```
...
<filter>
    <filter-name>hibernateFilter</filter-name>
    <filter-class>
    org.springframework.orm.hibernate3.support.OpenSessionInViewFilter
    </filter-class>
</filter>
<filter-mapping>
    <filter-name>hibernateFilter</filter-name>
    <url-pattern>*.html</url-pattern>
</filter-mapping>
...

```

Show moreShow more icon

上面的配置，我们假设使用 .html 的后缀作为 Web 框架的 URL 匹配模式，如果您使用 Struts 等 Web 框架，可以将其改为对应的”\*.do”模型。

**中文乱码过滤器**

在您通过表单向服务器提交数据时，一个经典的问题就是中文乱码问题。虽然我们所有的 JSP 文件和页面编码格式都采用 UTF-8，但这个问题还是会出现。解决的办法很简单，我们只需要在 web.xml 中配置一个 Spring 的编码转换过滤器就可以了：

```
<web-app>
<!---listener 的配置 -->
<filter>
    <filter-name>encodingFilter</filter-name>
    <filter-class>
        org.springframework.web.filter.CharacterEncodingFilter ① Spring 编辑过滤器
    </filter-class>
    <init-param> ② 编码方式
        <param-name>encoding</param-name>
        <param-value>UTF-8</param-value>
    </init-param>
    <init-param> ③ 强制进行编码转换
        <param-name>forceEncoding</param-name>
        <param-value>true</param-value>
    </init-param>
    </filter>
    <filter-mapping> ② 过滤器的匹配 URL
        <filter-name>encodingFilter</filter-name>
        <url-pattern>*.html</url-pattern>
    </filter-mapping>

<!---servlet 的配置 -->
</web-app>

```

Show moreShow more icon

这样所有以 .html 为后缀的 URL 请求的数据都会被转码为 UTF-8 编码格式，表单中文乱码的问题就可以解决了。

**请求跟踪日志过滤器**

除了以上两个常用的过滤器外，还有两个在程序调试时可能会用到的请求日志跟踪过滤器，它们会将请求的一些重要信息记录到日志中，方便程序的调试。这两个日志过滤器只有在日志级别为 DEBUG 时才会起作用：

方法说明`org.springframework.web.filter.ServletContextRequestLoggingFilter`该过滤器将请求的 URI 记录到 Common 日志中（如通过 Log4J 指定的日志文件）；`org.springframework.web.filter.ServletContextRequestLoggingFilter`该过滤器将请求的 URI 记录到 ServletContext 日志中。

以下是日志过滤器记录的请求跟踪日志的片断：

```
(JspServlet.java:224) - JspEngine --> /htmlTest.jsp
(JspServlet.java:225) - ServletPath: /htmlTest.jsp
(JspServlet.java:226) - PathInfo: null
(JspServlet.java:227) - RealPath: D:\masterSpring\chapter23\webapp\htmlTest.jsp
(JspServlet.java:228) - RequestURI: /baobaotao/htmlTest.jsp
...

```

Show moreShow more icon

通过这个请求跟踪日志，程度调试者可以详细地查看到有哪些请求被调用，请求的参数是什么，请求是否正确返回等信息。虽然这两个请求跟踪日志过滤器一般在程序调试时使用，但是即使程序部署不将其从 web.xml 中移除也不会有大碍，因为只要将日志级别设置为 DEBUG 以上级别，它们就不会输出请求跟踪日志信息了。

**转存 Web 应用根目录监听器和 Log4J 监听器**

Spring 在 org.springframework.web.util 包中提供了几个特殊用途的 Servlet 监听器，正确地使用它们可以完成一些特定需求的功能。比如某些第三方工具支持通过 ${key} 的方式引用系统参数（即可以通过 System.getProperty() 获取的属性），WebAppRootListener 可以将 Web 应用根目录添加到系统参数中，对应的属性名可以通过名为”webAppRootKey”的 Servlet 上下文参数指定，默认为”webapp.root”。下面是该监听器的具体的配置：

##### 清单 6\. WebAppRootListener 监听器配置

```
...
<context-param>
    <param-name>webAppRootKey</param-name>
    <param-value>baobaotao.root</param-value> ① Web 应用根目录以该属性名添加到系统参数中
</context-param>
...
② 负责将 Web 应用根目录以 webAppRootKey 上下文参数指定的属性名添加到系统参数中
<listener>
    <listener-class>
    org.springframework.web.util.WebAppRootListener
    </listener-class>
</listener>
...

```

Show moreShow more icon

这样，您就可以在程序中通过 System.getProperty(“baobaotao.root”) 获取 Web 应用的根目录了。不过更常见的使用场景是在第三方工具的配置文件中通过 ${baobaotao.root} 引用 Web 应用的根目录。比如以下的 log4j.properties 配置文件就通过 ${baobaotao.root} 设置了日志文件的地址：

```
log4j.rootLogger=INFO,R
log4j.appender.R=org.apache.log4j.RollingFileAppender
log4j.appender.R.File=${baobaotao.root}/WEB-INF/logs/log4j.log ① 指定日志文件的地址
log4j.appender.R.MaxFileSize=100KB
log4j.appender.R.MaxBackupIndex=1
log4j.appender.R.layout.ConversionPattern=%d %5p [%t] (%F:%L) - %m%n

```

Show moreShow more icon

另一个专门用于 Log4J 的监听器是 Log4jConfigListener。一般情况下，您必须将 Log4J 日志配置文件以 log4j.properties 为文件名并保存在类路径下。Log4jConfigListener 允许您通过 log4jConfigLocation Servlet 上下文参数显式指定 Log4J 配置文件的地址，如下所示：

```
① 指定 Log4J 配置文件的地址
<context-param>
    <param-name>log4jConfigLocation</param-name>
    <param-value>/WEB-INF/log4j.properties</param-value>
</context-param>
...
② 使用该监听器初始化 Log4J 日志引擎
<listener>
    <listener-class>
    org.springframework.web.util.Log4jConfigListener
    </listener-class>
</listener>
...

```

Show moreShow more icon

##### 提示

一些 Web 应用服务器（如 Tomcat）不会为不同的 Web 应用使用独立的系统参数，也就是说，应用服务器上所有的 Web 应用都共享同一个系统参数对象。这时，您必须通过 webAppRootKey 上下文参数为不同 Web 应用指定不同的属性名：如第一个 Web 应用使用 webapp1.root 而第二个 Web 应用使用 webapp2.root 等，这样才不会发生后者覆盖前者的问题。此外，WebAppRootListener 和 Log4jConfigListener 都只能应用在 Web 应用部署后 WAR 文件会解包的 Web 应用服务器上。一些 Web 应用服务器不会将 Web 应用的 WAR 文件解包，整个 Web 应用以一个 WAR 包的方式存在（如 Weblogic），此时因为无法指定对应文件系统的 Web 应用根目录，使用这两个监听器将会发生问题。

Log4jConfigListener 监听器包括了 WebAppRootListener 的功能，也就是说，Log4jConfigListener 会自动完成将 Web 应用根目录以 webAppRootKey 上下文参数指定的属性名添加到系统参数中，所以当您使用 Log4jConfigListener 后，就没有必须再使用 WebAppRootListener 了。

**Introspector 缓存清除监听器**

Spring 还提供了一个名为 org.springframework.web.util.IntrospectorCleanupListener 的监听器。它主要负责处理由 JavaBean Introspector 功能而引起的缓存泄露。IntrospectorCleanupListener 监听器在 Web 应用关闭的时会负责清除 JavaBean Introspector 的缓存，在 web.xml 中注册这个监听器可以保证在 Web 应用关闭的时候释放与其相关的 ClassLoader 的缓存和类引用。如果您使用了 JavaBean Introspector 分析应用中的类，Introspector 缓存会保留这些类的引用，结果在应用关闭的时候，这些类以及 Web 应用相关的 ClassLoader 不能被垃圾回收。不幸的是，清除 Introspector 的唯一方式是刷新整个缓存，这是因为没法准确判断哪些是属于本 Web 应用的引用对象，哪些是属于其它 Web 应用的引用对象。所以删除被缓存的 Introspection 会导致将整个 JVM 所有应用的 Introspection 都删掉。需要注意的是，Spring 托管的 Bean 不需要使用这个监听器，因为 Spring 的 Introspection 所使用的缓存在分析完一个类之后会马上从 javaBean Introspector 缓存中清除掉，并将缓存保存在应用程序特定的 ClassLoader 中，所以它们一般不会导致内存资源泄露。但是一些类库和框架往往会产生这个问题。例如 Struts 和 Quartz 的 Introspector 的内存泄漏会导致整个的 Web 应用的 ClassLoader 不能进行垃圾回收。在 Web 应用关闭之后，您还会看到此应用的所有静态类引用，这个错误当然不是由这个类自身引起的。解决这个问题的方法很简单，您仅需在 web.xml 中配置 IntrospectorCleanupListener 监听器就可以了：

```
<listener>
    <listener-class>
    org.springframework.web.util.IntrospectorCleanupListener
    </listener-class>
</listener>

```

Show moreShow more icon

## 结束语

本文介绍了一些常用的 Spring 工具类，其中大部分 Spring 工具类不但可以在基于 Spring 的应用中使用，还可以在其它的应用中使用。使用 JDK 的文件操作类在访问类路径相关、Web 上下文相关的文件资源时，往往显得拖泥带水、拐弯抹角，Spring 的 Resource 实现类使这些工作变得轻松了许多。

在 Web 应用中，有时你希望直接访问 Spring 容器，获取容器中的 Bean，这时使用 WebApplicationContextUtils 工具类从 ServletContext 中获取 WebApplicationContext 是非常方便的。WebUtils 为访问 Servlet API 提供了一套便捷的代理方法，您可以通过 WebUtils 更好的访问 HttpSession 或 ServletContext 的信息。

Spring 提供了几个 Servlet 过滤器和监听器，其中 ServletContextRequestLoggingFilter 和 ServletContextRequestLoggingFilter 可以记录请求访问的跟踪日志，你可以在程序调试时使用它们获取请求调用的详细信息。WebAppRootListener 可以将 Web 应用的根目录以特定属性名添加到系统参数中，以便第三方工具类通过 ${key} 的方式进行访问。Log4jConfigListener 允许你指定 Log4J 日志配置文件的地址，且可以在配置文件中通过 ${key} 的方式引用 Web 应用根目录，如果你需要在 Web 应用相关的目录创建日志文件，使用 Log4jConfigListener 可以很容易地达到这一目标。

Web 应用的内存泄漏是最让开发者头疼的问题，虽然不正确的程序编写可能是这一问题的根源，也有可能是一些第三方框架的 JavaBean Introspector 缓存得不到清除而导致的，Spring 专门为解决这一问题配备了 IntrospectorCleanupListener 监听器，它只要简单在 web.xml 中声明该监听器就可以了。

在 [下一篇](https://www.ibm.com/developerworks/cn/java/j-lo-spring-utils2/) 文章中，我们将继续介绍 Spring 中用于特殊字符转义和方法入参检测相关的工具类。