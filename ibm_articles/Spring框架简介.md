# Spring 框架简介
Spring AOP 和 IOC 容器入门

**标签:** Java,Spring,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/wa-spring1/)

Naveen Balani

发布: 2005-08-18

* * *

Spring 是一个开源框架，是为了解决企业应用程序开发复杂性而创建的。框架的主要优势之一就是其分层架构，分层架构允许您选择使用哪一个组件，同时为 J2EE 应用程序开发提供集成的框架。

在这篇由四部分组成的 _Spring 系列_ 的第 1 部分中，我将介绍 Spring 框架。我先从框架底层模型的角度描述该框架的功能，然后将讨论两个最有趣的模块：Spring 面向方面编程（AOP）和控制反转 （IOC） 容器。接着将使用几个示例演示 IOC 容器在典型应用程序用例场景中的应用情况。这些示例还将成为本系列后面部分进行的展开式讨论的基础，在本文的后面部分，将介绍 Spring 框架通过 Spring AOP 实现 AOP 构造的方式。

[下载 Spring 框架和 Apache Ant](http://www.ibm.com/developerworks/apps/download/index.jsp?contentid=93033&filename=wa-spring1-SpringProject.zip&method=http&locale=zh_CN)，运行本系列的示例应用程序需要它们。

## Spring 框架

Spring 框架是一个分层架构，由 7 个定义良好的模块组成。Spring 模块构建在核心容器之上，核心容器定义了创建、配置和管理 bean 的方式，如图 1 所示。

##### 图 1\. Spring 框架的 7 个模块

![Spring 框架图示](../ibm_articles_img/wa-spring1_images_spring_framework.gif)

组成 Spring 框架的每个模块（或组件）都可以单独存在，或者与其他一个或多个模块联合实现。每个模块的功能如下：

- **核心容器** ：核心容器提供 Spring 框架的基本功能。核心容器的主要组件是 `BeanFactory`，它是工厂模式的实现。 `BeanFactory` 使用 _控制反转_ （IOC） 模式将应用程序的配置和依赖性规范与实际的应用程序代码分开。
- **Spring 上下文** ：Spring 上下文是一个配置文件，向 Spring 框架提供上下文信息。Spring 上下文包括企业服务，例如 JNDI、EJB、电子邮件、国际化、校验和调度功能。
- **Spring AOP** ：通过配置管理特性，Spring AOP 模块直接将面向方面的编程功能集成到了 Spring 框架中。所以，可以很容易地使 Spring 框架管理的任何对象支持 AOP。Spring AOP 模块为基于 Spring 的应用程序中的对象提供了事务管理服务。通过使用 Spring AOP，不用依赖 EJB 组件，就可以将声明性事务管理集成到应用程序中。
- **Spring DAO** ：JDBC DAO 抽象层提供了有意义的异常层次结构，可用该结构来管理异常处理和不同数据库供应商抛出的错误消息。异常层次结构简化了错误处理，并且极大地降低了需要编写的异常代码数量（例如打开和关闭连接）。Spring DAO 的面向 JDBC 的异常遵从通用的 DAO 异常层次结构。
- **Spring ORM** ：Spring 框架插入了若干个 ORM 框架，从而提供了 ORM 的对象关系工具，其中包括 JDO、Hibernate 和 iBatis SQL Map。所有这些都遵从 Spring 的通用事务和 DAO 异常层次结构。
- **Spring Web 模块** ：Web 上下文模块建立在应用程序上下文模块之上，为基于 Web 的应用程序提供了上下文。所以，Spring 框架支持与 Jakarta Struts 的集成。Web 模块还简化了处理多部分请求以及将请求参数绑定到域对象的工作。
- **Spring MVC 框架** ：MVC 框架是一个全功能的构建 Web 应用程序的 MVC 实现。通过策略接口，MVC 框架变成为高度可配置的，MVC 容纳了大量视图技术，其中包括 JSP、Velocity、Tiles、iText 和 POI。

Spring 框架的功能可以用在任何 J2EE 服务器中，大多数功能也适用于不受管理的环境。Spring 的核心要点是：支持不绑定到特定 J2EE 服务的可重用业务和数据访问对象。毫无疑问，这样的对象可以在不同 J2EE 环境 （Web 或 EJB）、独立应用程序、测试环境之间重用。

## IOC 和 AOP

控制反转模式（也称作依赖性介入）的基本概念是：不创建对象，但是描述创建它们的方式。在代码中不直接与对象和服务连接，但在配置文件中描述哪一个组件需要哪一项服务。容器 （在 Spring 框架中是 IOC 容器） 负责将这些联系在一起。

在典型的 IOC 场景中，容器创建了所有对象，并设置必要的属性将它们连接在一起，决定什么时间调用方法。下表列出了 IOC 的一个实现模式。

类型实现模式类型 1服务需要实现专门的接口，通过接口，由对象提供这些服务，可以从对象查询依赖性（例如，需要的附加服务）类型 2通过 JavaBean 的属性（例如 setter 方法）分配依赖性类型 3依赖性以构造函数的形式提供，不以 JavaBean 属性的形式公开

Spring 框架的 IOC 容器采用类型 2 和类型3 实现。

### 面向方面的编程

_面向方面的编程_ ，即 AOP，是一种编程技术，它允许程序员对横切关注点或横切典型的职责分界线的行为（例如日志和事务管理）进行模块化。AOP 的核心构造是 _方面_ ，它将那些影响多个类的行为封装到可重用的模块中。

AOP 和 IOC 是补充性的技术，它们都运用模块化方式解决企业应用程序开发中的复杂问题。在典型的面向对象开发方式中，可能要将日志记录语句放在所有方法和 Java 类中才能实现日志功能。在 AOP 方式中，可以反过来将日志服务 _模块化_ ，并以声明的方式将它们应用到需要日志的组件上。当然，优势就是 Java 类不需要知道日志服务的存在，也不需要考虑相关的代码。所以，用 Spring AOP 编写的应用程序代码是松散耦合的。

AOP 的功能完全集成到了 Spring 事务管理、日志和其他各种特性的上下文中。

## IOC 容器

Spring 设计的核心是 `org.springframework.beans` 包，它的设计目标是与 JavaBean 组件一起使用。这个包通常不是由用户直接使用，而是由服务器将其用作其他多数功能的底层中介。下一个最高级抽象是 `BeanFactory` 接口，它是工厂设计模式的实现，允许通过名称创建和检索对象。 `BeanFactory` 也可以管理对象之间的关系。

`BeanFactory` 支持两个对象模型。

- **单态** 模型提供了具有特定名称的对象的共享实例，可以在查询时对其进行检索。Singleton 是默认的也是最常用的对象模型。对于无状态服务对象很理想。
- **原型** 模型确保每次检索都会创建单独的对象。在每个用户都需要自己的对象时，原型模型最适合。

bean 工厂的概念是 Spring 作为 IOC 容器的基础。IOC 将处理事情的责任从应用程序代码转移到框架。正如我将在下一个示例中演示的那样，Spring 框架使用 JavaBean 属性和配置数据来指出必须设置的依赖关系。

### BeanFactory 接口

因为 `org.springframework.beans.factory.BeanFactory` 是一个简单接口，所以可以针对各种底层存储方法实现。最常用的 `BeanFactory` 定义是 `XmlBeanFactory`，它根据 XML 文件中的定义装入 bean，如清单 1 所示。

##### 清单 1\. XmlBeanFactory

```
BeanFactory factory = new XMLBeanFactory(new FileInputSteam("mybean.xml"));

```

Show moreShow more icon

在 XML 文件中定义的 Bean 是被消极加载的，这意味在需要 bean 之前，bean 本身不会被初始化。要从 `BeanFactory` 检索 bean，只需调用 `getBean()` 方法，传入将要检索的 bean 的名称即可，如清单 2 所示。

##### 清单 2\. getBean()

```
MyBean mybean = (MyBean) factory.getBean("mybean");

```

Show moreShow more icon

每个 bean 的定义都可以是 POJO （用类名和 JavaBean 初始化属性定义） 或 `FactoryBean` 。 `FactoryBean` 接口为使用 Spring 框架构建的应用程序添加了一个间接的级别。

## IOC 示例

理解控制反转最简单的方式就是看它的实际应用。在对由三部分组成的 _Spring 系列_ 的第 1 部分进行总结时，我使用了一个示例，演示了如何通过 Spring IOC 容器注入应用程序的依赖关系（而不是将它们构建进来）。

我用开启在线信用帐户的用例作为起点。对于该实现，开启信用帐户要求用户与以下服务进行交互：

- 信用级别评定服务，查询用户的信用历史信息。
- 远程信息链接服务，插入客户信息，将客户信息与信用卡和银行信息连接起来，以进行自动借记（如果需要的话）。
- 电子邮件服务，向用户发送有关信用卡状态的电子邮件。

## 三个接口

对于这个示例，我假设服务已经存在，理想的情况是用松散耦合的方式把它们集成在一起。以下清单显示了三个服务的应用程序接口。

##### 清单 3\. CreditRatingInterface

```
public interface CreditRatingInterface {
public boolean getUserCreditHistoryInformation(ICustomer iCustomer);
}

```

Show moreShow more icon

清单 3 所示的信用级别评定接口提供了信用历史信息。它需要一个包含客户信息的 `Customer` 对象。该接口的实现是由 `CreditRating` 类提供的。

##### 清单 4\. CreditLinkingInterface

```
public interface CreditLinkingInterface {
public String getUrl();
        public void setUrl(String url);
        public void linkCreditBankAccount() throws Exception ;
}

```

Show moreShow more icon

信用链接接口将信用历史信息与银行信息（如果需要的话）连接在一起，并插入用户的信用卡信息。信用链接接口是一个远程服务，它的查询是通过 `getUrl()` 方法进行的。URL 由 Spring 框架的 bean 配置机制设置，我稍后会讨论它。该接口的实现是由 `CreditLinking` 类提供的。

##### 清单 5\. EmailInterface

```
public interface EmailInterface {
      public void sendEmail(ICustomer iCustomer);
      public String getFromEmail();
      public void setFromEmail(String fromEmail) ;
      public String getPassword();
      public void setPassword(String password) ;
      public String getSmtpHost() ;
      public void setSmtpHost(String smtpHost);
      public String getUserId() ;
      public void setUserId(String userId);
}

```

Show moreShow more icon

`EmailInterface` 负责向客户发送关于客户信用卡状态的电子邮件。邮件配置参数（例如 SMPT 主机、用户名、口令）由前面提到的 bean 配置机制设置。 `Email` 类提供了该接口的实现。

## Spring 使其保持松散

这些接口就位之后，接下来要考虑的就是如何用松散耦合方式将它们集成在一起。在清单 6 中可以看到信用卡帐户用例的实现。

##### 清单 6\. CreateCreditCardAccount

```
public class CreateCreditCardAccount implements
CreateCreditCardAccountInterface {

public CreditLinkingInterface getCreditLinkingInterface() {
return creditLinkingInterface;
}

public void setCreditLinkingInterface(
CreditLinkingInterface creditLinkingInterface) {
this.creditLinkingInterface = creditLinkingInterface;
}

public CreditRatingInterface getCreditRatingInterface() {
return creditRatingInterface;
}

public void setCreditRatingInterface(CreditRatingInterface creditRatingInterface)
{
this.creditRatingInterface = creditRatingInterface;
}

public EmailInterface getEmailInterface() {
return emailInterface;
}

public void setEmailInterface(EmailInterface emailInterface) {
this.emailInterface = emailInterface;
}

//Client will call this method
public void createCreditCardAccount(ICustomer icustomer) throws Exception{
         boolean crediRating =
         getCreditRatingInterface().getUserCreditHistoryInformation(icustomer);
         icustomer.setCreditRating(crediRating);
         //Good Rating
         if(crediRating){
         getCreditLinkingInterface().linkCreditBankAccount(icustomer);
         }

         getEmailInterface().sendEmail(icustomer);

}

```

Show moreShow more icon

注意，所有的 setter 方法都是由 Spring 的配置 bean 实现的。所有的依赖关系 （也就是三个接口）都可以由 Spring 框架用这些 bean 注入。 `createCreditCardAccount()` 方法会用服务去执行其余实现。在清单 7 中可以看到 Spring 的配置文件。我用箭头突出了这些定义。

##### 清单 7\. 配置文件

```
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE beans PUBLIC "-//SPRING//DTD BEAN//EN"
         "http://www.springframework.org/dtd/spring-beans.dtd">
<beans>

<bean id="createCreditCard" -->define createCreditCard definition
class="springexample.creditcardaccount.CreateCreditCardAccount">
<property name="creditRatingInterface">-->inject creditRatingInterface dependency  via creditRating reference bean
<ref bean="creditRating" />
</property>
<property name="creditLinkingInterface">">-->inject creditLinkingInterface dependency  via creditLinking reference bean
<ref bean="creditLinking" />
</property>
<property name="emailInterface">">">-->inject emailInterface dependency  via email reference bean
<ref bean="email" />
/property>
</bean>

<bean id="creditLinking" class="springexample.creditlinking.CreditLinking">
<property name="url">
<value>http://localhost/creditLinkService</value>-->set url property value
</property>
</bean>

<bean id="creditRating" class="springexample.creditrating.CreditRating">
</bean>

<bean id="email" class="springexample.email.Email">
<property name="smtpHost">
<value>localhost</value>>-->set smpHtpHost property value
</property>
<property name="fromEmail">
<value>mycompanyadmin@mycompanyadmin.com</value>
</property>
<property name="userId">
<value>myuserid</value>
</property>
<property name="password">
<value>mypassword</value>
</property>
</bean>
</beans>

```

Show moreShow more icon

## 运行应用程序

要运行示例应用程序，首先必须 [下载 Spring 框架](https://sourceforge.net/projects/springframework/files/springframework/1.2%20RC2/spring-framework-1.2-rc2-with-dependencies.zip/download?use_mirror=master&download=) 及其所有依赖文件。接下来，将框架释放到（比如说）磁盘 _c:\_ ，这会创建 _C:\\spring-framework-1.2-rc2_ （适用于当前发行版本） 这样的文件夹。在继续后面的操作之前，还必须下载和释放 [Apache Ant](http://ant.apache.org/) 。

接下来，将源代码释放到文件夹，例如 _c:\_ 盘，然后创建 _SpringProject_ 。将 Spring 库（即 _C:\\spring-framework-1.2-rc2\\dist_ 下的 _spring.jar_ 和 _C:\\spring-framework-1.2-rc2\\lib\\jakarta-commons_ 下的 _commons-logging.jar_ ）复制到 _SpringProject\\lib_ 文件夹中。完成这些工作之后，就有了必需的构建依赖关系集。

打开命令提示符，将当前目录切换到 _SpringProject_ ，在命令提示符中输入以下命令： `build` 。

这会构建并运行 `CreateCreditAccountClient` 类，类的运行将创建 `Customer` 类对象并填充它，还会调用 `CreateCreditCardAccount` 类创建并链接信用卡帐户。 `CreateCreditAccountClient` 还会通过 `ClassPathXmlApplicationContext` 装入 Spring 配置文件。装入 bean 之后，就可以通过 `getBean()` 方法访问它们了，如清单 8 所示。

##### 清单 8\. 装入 Spring 配置文件

```
ClassPathXmlApplicationContext appContext =
                    new ClassPathXmlApplicationContext(new String[] {
     "springexample-creditaccount.xml"
    });
CreateCreditCardAccountInterface creditCardAccount =
                    (CreateCreditCardAccountInterface)
    appContext.getBean("createCreditCard");

```

Show moreShow more icon

## 结束语

在这篇由三部分组成的 _Spring 系列_ 的第一篇文章中，我介绍了 Spring 框架的基础。我从讨论组成 Spring 分层架构的 7 个模块开始，然后深入介绍了其中两个模块：Spring AOP 和 IOC 容器。

由于学习的最佳方法是实践，所以我用一个工作示例介绍了 IOC 模式 （像 Spring 的 IOC 容器实现的那样）如何用松散耦合的方式将分散的系统集成在一起。在这个示例中可以看到，将依赖关系或服务注入工作中的信用卡帐户应用程序，要比从头开始构建它们容易得多。

请继续关注这一系列的下一篇文章，我将在这里学习的知识基础上，介绍 Spring AOP 模块如何在企业应用程序中提供持久支持，并让您开始了解 Spring MVC 模块和相关插件。