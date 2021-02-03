# 当 Hibernate 遇上 Spring
Hibernate 事务天生适合 Spring AOP

**标签:** Java,Spring

[原文链接](https://developer.ibm.com/zh/articles/wa-spring2/)

Naveen Balani

发布: 2005-09-26

* * *

在这个系列的 [前一期中](https://www.ibm.com/developerworks/cn/java/wa-spring1/) ，我介绍了 Spring 框架的 7 个模块，包括 Spring AOP 和控制反转（IOC）容器。然后我用一个简单的示例演示了 IOC 模式（由 Spring IOC 容器实现）如何用松散耦合的方式集成分散的系统。

现在，我从我上次结束的地方开始，采用与上次类似的示例，演示 Spring AOP 和 Spring Hibernate 持久性支持的声明性事务处理，所以我首先从对这两项技术的深入研究开始。

[下载这篇文章的源代码](http://www.ibm.com/developerworks/apps/download/index.jsp?contentid=94734&filename=wa-spring2-SpringProjectPart2.zip&method=http&locale=zh_CN)。请访问 Spring 框架和 Apache Ant，运行这篇文章的示例应用程序需要它们。

## Spring AOP

软件系统通常由多个组件构成，每个组件负责一个特定的功能领域。但是，这些组件也经常承担它们的核心功能之外的额外责任。系统服务（例如日志、事务管理和安全性）经常发现自己跑到了别的组件的领域里，而这些组件的核心职责是其他事情。结果就是所谓的”代码纠缠”，或者更简单点儿说”一团糟”。面向方面编程是一种试图解决这个问题的编程技术，它把关注点的隔离提升为核心的编程概念。

使用 AOP 时，仍然是在一个地方定义系统的公共功能，但是可以声明性地定义 _如何_ 和 _在哪里_ 应用这个功能。如果对横切关注点（例如日志和事务管理）进行了模块化，那么不用修改每个单独的类，就可以向代码中添加新特性。这类模块化的关注点称作 _方面_ 。

##### 您知道么？

可以在任何 Java™ 企业版（J2EE） 服务器中使用 Spring 框架的功能。而且，还可以调整它的大多数功能，使其适合不受控环境。Spring 的中心焦点就是支持不被束缚在特定 J2EEE 服务上的可重用业务和数据访问对象。可以跨 J2EE 环境（Web 或企业 JavaBean（EJB））、独立应用程序、测试环境等等重用这类对象，而不会有任何麻烦。

以一个企业应用程序为例。这类应用程序通常要求类似于安全性和事务支持的服务。显然，可以把这些服务的支持直接编写到要求服务的每个类当中，但是更希望能够不必为大量事务性上下文编写同样的事务处理代码。如果使用 Spring AOP 进行事务处理，那么可以声明性地安排适当的方法调用，而不必逐个安排。

Spring AOP 提供了几个方面，可以为 JavaBean 声明事务。例如，`TransactionProxyFactoryBean` 是个方便的代理类，能够拦截对现有类的方法调用，并把事务上下文应用到事务 bean。在下面的示例中会看到这个类的实际应用。

## Hibernate

Spring 框架提供了对 Hibernate、JDO 和 iBATIS SQL Maps 的集成支持。Spring 对 Hibernate 的支持是第一级的，整合了许多 IOC 的方便特性，解决了许多典型的 Hibernate 集成问题。框架对 Hibernate 的支持符合 Spring 通用的事务和数据访问对象（DAO）异常层次结构。

Spring 为使用选择的 OR 映射层来创建数据访问应用程序提供了支持。因为所有东西都设计成一组可重用 JavaBean，所以不管选择什么技术，都能以库的格式访问大多数 Spring 的 OR 映射支持。`ApplicationContext` 或 `BeanFactory` 内部的 OR 映射的好处是简化了配置和部署。

Hibernate 是 Java 平台上一个功能全面的、开源的 OR 映射框架。Hibernate 支持开发符合常规 Java 理念的持久性类 —— 包括关联、继承、多态、复合以及 Java 集合框架。Hibernate 查询语言（HQL）被设计成 SQL 的一个微型面向对象扩展，它是对象和关系世界之间的桥梁。Hibernate 也支持用原始 SQL 或基于 Java 的标准和示例查询表达查询。Hibernate 使用 XML（ _\*.hbm.xml_ ） 文件把 Java 类映射到表，把 JavaBean 属性映射到数据库表。

通过 JDBC 技术，支持所有的 SQL 数据库管理系统。Hibernate 与所有流行的 J2EE 应用程序服务器和 Web 容器都很好地集成。

## 实际示例

一个银行应用程序示例可以让您自己看到 Spring AOP 和 Hibernate 一起工作有多么好。银行帐户用例允许用户 （ `Customer` ） 在一个事务中打开一个或多个银行帐户。用户可以申请多个银行帐户，可以选择是支票帐户类型或者是储蓄帐户类型。

应用程序数据库（Cloudscape™）容纳所有客户和帐户信息。在这个例子中，假设在 `Customer` 和 `Account` 类之间存在 _1:N_ 的关联。在实际生活场景中，关联可能需要按 _m:n_ 建模，才能支持联合帐户。

由于用户必须可以在一个事务中申请多个帐户，所以首先要为数据库交互实现一个 DOA 模式。然后要设置 Spring AOP 的 `TransactionProxyFactoryBean` ，让它拦截方法调用并声明性地把事务上下文应用到 DOA。

## Hibernate 实践

在 Spring 框架中，像 JDBC `DataSource` 或 Hibernate `SessionFactory` 这样的资源，在应用程序上下文中可以用 bean 实现。需要访问资源的应用程序对象只需通过 bean 引用得到这类预先定义好的实例的引用即可（这方面的更多内容在 [分析这个!](#分析这个) ）。在清单 1 中，可以看到示例银行应用程序的一段摘录：XML 应用程序上下文定义显示了如何设置 JDBC `DataSource` ，并在上面放一个 Hibernate `SessionFactory` 。

##### 清单 1\. JDBC DataSource 和 HibernateSessionFactory 连接

```
<!-- DataSource Property -->
<bean id="exampleDataSource"
    class="org.apache.commons.dbcp.BasicDataSource">
<property name="driverClassName">
<value>org.apache.derby.jdbc.EmbeddedDriver</value>
</property>
<property name="url">
<value>jdbc:derby:springexample;create=true</value>
</property>
</bean>
<!-- Database Property -->
<bean id="exampleHibernateProperties"
class="org.springframework.beans.factory.config.PropertiesFactoryBean">
<property name="properties">
<props>
<prop key="hibernate.hbm2ddl.auto">update</prop>
<prop
    key="hibernate.dialect">net.sf.hibernate.dialect.DerbyDialect</prop>
<prop
    key="hibernate.query.substitutions">true 'T', false 'F'</prop>
<prop key="hibernate.show_sql">false</prop>
<prop key="hibernate.c3p0.minPoolSize">5</prop>
<prop key="hibernate.c3p0.maxPoolSize">20</prop>
<prop key="hibernate.c3p0.timeout">600</prop>
<prop key="hibernate.c3p0.max_statement">50</prop>
<prop
     key="hibernate.c3p0.testConnectionOnCheckout">false</prop>
</props>
</property>
</bean>
<!-- Hibernate SessionFactory -->
<bean id="exampleSessionFactory"
class="org.springframework.orm.hibernate.LocalSessionFactoryBean">
<property name="dataSource">
     <ref local="exampleDataSource"/>
</property>
<property name="hibernateProperties">
     <ref bean="exampleHibernateProperties" />
</property>
<!--  OR mapping files. -->
<property name="mappingResources">
    <list>
       <value>Customer.hbm.xml</value>
       <value>Account.hbm.xml</value>
    </list>
</property>
</bean>

```

Show moreShow more icon

清单 1 显示了如何为示例应用程序数据库（是 Cloudscape）配置数据源 bean （ `exampleDataSource` ）。 `exampleDatasource` 被连接到 Spring Hibernate 的 `SessionFactory` 。请注意 _\*.hbm.xml_ 指定了示例应用程序的 OR 映射文件。

数据源和会话工厂设置好之后，下一步就是在 DAO 中连接，在 `CustomerDAOImpl` 示例中，要使用 `SessionFactory` 。接下来，插入 Spring 的 `TransactionProxyFactoryBean` ，它会拦截对应用程序的 `CustomerDAOImpl` 对象的方法调用，并声明性地在它上面应用事务。

##### 清单 2\. 将应用程序 DAO 和 TransactionManager 编写在一起

```
<!-- Pass the session factory to our CustomerDAO -->
<bean id="customerDAOTarget"
     class="springexample.hibernate.CustomerDAOImpl">
<property name="sessionFactory"><ref local="exampleSessionFactory"/>
</property>
</bean>
<bean id="transactionManager"
class="org.springframework.orm.hibernate.HibernateTransactionManager">
<property name="sessionFactory">
     <ref bean="exampleSessionFactory"/>
</property>
</bean>
<bean id="userDAO"
class="org.springframework.transaction.interceptor.TransactionProxyFactoryBean">
<property name="transactionManager"><ref local="transactionManager"/>
</property>
<property name="target"><ref local="customerDAOTarget"/>
</property>
<property name="transactionAttributes">
     <props>
       <prop key="addCustomer">PROPAGATION_REQUIRED</prop>
     </props>
</property>
</bean>

```

Show moreShow more icon

在清单 2 的这个示例中， `CustomerDAOImpl` 类的 `addCustomer` 方法是作为事务的一部分执行的，有一个事务属性 `PROPAGATION_REQUIRED` 。这个属性等价于 EJB 容器的 `TX_REQUIRED` 。如果想让这个方法一直在事务中运行，可以使用 `PROPAGATION_REQUIRED` 。如果事务已经在运行，那么 bean 方法会加入事务，否则 Spring 的轻量级事务管理器会启动一个事务。如果想在调用组件服务时总是启动新事务，可以使用 `PROPAGATION_REQUIRES_NEW` 属性。

应用程序的连接完成之后，现在来进一步查看源代码。

## 分析这个!

如果以前没这么做过，那么请 [下载这篇文章的源代码](http://www.ibm.com/developerworks/apps/download/index.jsp?contentid=94734&filename=wa-spring2-SpringProjectPart2.zip&method=http&locale=zh_CN)。把源 Zip 文件释放到计算机中的任何位置上，例如 _c:\_ 。会创建一个叫作 _SpringProjectPart2_ 的文件夹。 _src\\spring_ 文件夹包含示例应用程序的 Hibernate 映射文件和 Spring 配置文件。 _src\\springexample\\hibernate_ 文件包含应用程序的源代码。

在这里会发现两个类，即 `Customer` 和 `Account`，它们用 Hibernate 映射文件映射到两个表。`Customer` 类代表客户信息，`Account` 代表客户的帐户信息。正如前面提到的，我把这两个类按照 _1: N_ 关系进行建模，即一个 `Customer` 可以拥有多个 `Account` 。清单 3 显示了 `Customer` 对象的 Hibernate 映射文件。

##### 清单 3\. Customer 对象的 Hibernate 映射文件

```
<?xml version="1.0"?>
<!DOCTYPE hibernate-mapping PUBLIC
         "-//Hibernate/Hibernate Mapping DTD 2.0//EN"
         "http://hibernate.sourceforge.net/hibernate-mapping-2.0.dtd">
<hibernate-mapping>
<class
      name="springexample.hibernate.Customer"
      table="TBL_CUSTOMER"
      dynamic-update="false"
      dynamic-insert="false">
      <id
         name="id"
         column="CUSTOMER_ID"
         type="java.lang.Long"
         unsaved-value="-1"
      >
         <generator class="native">
         </generator>
      </id>
      <set name ="accounts"
         inverse = "true"
         cascade="all-delete-orphan">
         <key column ="CUSTOMER_ID"/>
         <one-to-many class="springexample.hibernate.Account"/>
      </set>
      <property
         name="email"
         type="string"
         update="false"
         insert="true"
         column="CUSTOMER_EMAIL"
         length="82"
         not-null="true"
      />
      <property
         name="password"
         type="string"
         update="false"
         insert="true"
         column="CUSTOMER_PASSWORD"
         length="10"
         not-null="true"
      />
      <property
         name="userId"
         type="string"
         update="false"
         insert="true"
         column="CUSTOMER_USERID"
         length="12"
         not-null="true"
         unique="true"
      />
      <property
         name="firstName"
         type="string"
         update="false"
         insert="true"
         column="CUSTOMER_FIRSTNAME"
         length="25"
         not-null="true"
      />
      <property
         name="lastName"
         type="string"
         update="false"
         insert="true"
         column="CUSTOMER_LASTTNAME"
         length="25"
         not-null="true"

      />
</class>
</hibernate-mapping>

```

Show moreShow more icon

`set name="accounts"` 和一对多类标签指定了 `Customer` 和 `Account` 之间的关系。我还在 _Account.hbm.xml_ 文件中定义了 `Account` 对象的映射。

`CustomerDAOImpl.java` 代表应用程序的 DAO，它在应用程序数据库中插入客户和帐户信息。 `CustomerDAOImpl` 扩展了 Spring 的 `HibernateDaoSupport` ，它用 Spring HibernateTemplate 简化了会话管理。这样，可以通过 `getHibernateTemplate()` 方法保存或检索数据。下面显示的 `getCustomerAccountInfo()` 对 `Customer` 进行 _查找_ ，通过 `getHibernateTemplate().find` 方法用 HQL 得到客户的帐户信息，如清单 4 所示。

##### 清单 4\. DAO 实现

```
public class CustomerDAOImpl extends HibernateDaoSupport
           implements CustomerDAO{
public void addCustomer(Customer customer) {
    getHibernateTemplate().save(customer);
    // TODO Auto-generated method stub
}
public Customer getCustomerAccountInfo(Customer customer) {
    Customer cust = null;
    List list = getHibernateTemplate().find("from Customer customer " +
          "where customer.userId = ?" ,
          customer.getUserId(),Hibernate.STRING);
    if(list.size() > 0){
          cust = (Customer)  list.get(0);
    }
    return cust;
}

```

Show moreShow more icon

所有这些都应当很容易掌握。现在来看代码的实际应用!

## 运行应用程序

要运行示例应用程序，必须首先 [下载 Spring 框架](https://sourceforge.net/projects/springframework/files/springframework/1.2%20RC2/spring-framework-1.2-rc2-with-dependencies.zip/download?use_mirror=master&download=) 和它的全部依赖文件。接下来，释放框架到某一位置（比如 _c:\_ ），这会创建文件夹 _C:\\spring-framework-1.2-rc2_ （针对当前发行版）。在继续之前还必须下载和释放 [Apache Ant](https://ant.apache.org/) 和 [Cloudscape](http://db.apache.org/derby/)（已被 Apache Derby 取代） 。下载 Cloudscape 之后，把它释放到 _c:\_ ，这会创建文件夹 _C:\\Cloudscape\_10.0_ 。

接下来，释放源代码到 _c:\_ ，这会创建 _SpringProject2_ 文件夹。接下来修改 _build.xml_ 文件的入口，用实际安装 Spring 的位置代替 _C:\\spring-framework-1.2-rc2_ ，用实际安装 Cloudscape 的位置代替 _C:\\Program Files\\IBM\\Cloudscape\_10.0_ 。

打开命令行提示符，进入 _SpringProject_ 目录，在命令行提示符下输入以下命令：`build`.

这会构建并运行 `CreateBankCustomerClient` 类，它会创建 `Customer` 类对象，用数据填充它，创建 `Account` 对象，填充它，并把它添加到 `Customer` 对象。

然后 `CreateBankCustomerClient` 会调用 `CustomerDAOImpl.addCustomer` 类，添加客户和帐户信息。一旦插入完成， `CreateBankCustomerClient` 会调用 `CustomerDAOImpl.getCustomerAccountInfo` 方法，根据 `userid` 得到客户和帐户信息。如果 `CreateBankCustomerClient` 执行成功，会在控制台上看到打印出 `userid` 。也可以查询 Cloudscape 数据库检索客户和帐户信息。

## 结束语

在三部分的 _Spring 系列_ 的第 2 部分中，我介绍了如何集成 Spring Hibernate 和 Spring AOP。结果是一个强健的持久性框架，支持声明性的实现事务。

在这个系列的下一篇，我将介绍 Spring 的 MVC 模块，介绍如何用它来简化基于 Web 的应用程序的创建。