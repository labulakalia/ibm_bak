# 使用 Spring Security 保护 Web 应用的安全
了解如何使用 Spring Security 来实现不同的用户认证和授权机制

**标签:** Java,Spring,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/j-lo-springsecurity/)

成富

发布: 2010-12-02

* * *

在 Web 应用开发中，安全一直是非常重要的一个方面。安全虽然属于应用的非功能性需求，但是应该在应用开发的初期就考虑进来。如果在应用开发的后期才考虑安的问题，就可能陷入一个两难的境地：一方面，应用存在严重的安全漏洞，无法满足用户的要求，并可能造成用户的隐私数据被攻击者窃取；另一方面，应用的基本架构已经确定，要修复安全漏洞，可能需要对系统的架构做出比较重大的调整，因而需要更多的开发时间，影响应用的发布进程。因此，从应用开发的第一天就应该把安全相关的因素考虑进来，并在整个应用的开发过程中。

本文详细介绍了如何使用 Spring Security 来保护 Web 应用的安全。Spring Security 本身以及 Spring 框架带来的灵活性，能够满足一般 Web 应用开发的典型需求，并允许开发人员进行定制。下面首先简单介绍 Spring Security。

## Spring Security 简介

Spring 是一个非常流行和成功的 Java 应用开发框架。Spring Security 基于 Spring 框架，提供了一套 Web 应用安全性的完整解决方案。一般来说，Web 应用的安全性包括用户认证（Authentication）和用户授权（Authorization）两个部分。用户认证指的是验证某个用户是否为系统中的合法主体，也就是说用户能否访问该系统。用户认证一般要求用户提供用户名和密码。系统通过校验用户名和密码来完成认证过程。用户授权指的是验证某个用户是否有权限执行某个操作。在一个系统中，不同用户所具有的权限是不同的。比如对一个文件来说，有的用户只能进行读取，而有的用户可以进行修改。一般来说，系统会为不同的用户分配不同的角色，而每个角色则对应一系列的权限。

对于上面提到的两种应用情景，Spring Security 框架都有很好的支持。在用户认证方面，Spring Security 框架支持主流的认证方式，包括 HTTP 基本认证、HTTP 表单验证、HTTP 摘要认证、OpenID 和 LDAP 等。在用户授权方面，Spring Security 提供了基于角色的访问控制和访问控制列表（Access Control List，ACL），可以对应用中的领域对象进行细粒度的控制。

本文将通过三个具体的示例来介绍 Spring Security 的使用。第一个示例是一个简单的企业员工管理系统。该系统中存在三类用户，分别是普通员工、经理和总裁。不同类别的用户所能访问的资源不同。对这些资源所能执行的操作也不相同。Spring Security 能帮助开发人员以简单的方式满足这些安全性相关的需求。第二个示例展示了如何与 LDAP 服务器进行集成。第三个示例展示了如何与 OAuth 进行集成。完整的示例代码 [可在此获取](http://www.ibm.com/developerworks/cn/java/j-lo-springsecurity/Company.zip)。下面首先介绍基本的用户认证和授权的实现。

## 基本用户认证和授权

本节从最基本的用户认证和授权开始对 Spring Security 进行介绍。一般来说，Web 应用都需要保存自己系统中的用户信息。这些信息一般保存在数据库中。用户可以注册自己的账号，或是由系统管理员统一进行分配。这些用户一般都有自己的角色，如普通用户和管理员之类的。某些页面只有特定角色的用户可以访问，比如只有管理员才可以访问 `/admin` 这样的网址。下面介绍如何使用 Spring Security 来满足这样基本的认证和授权的需求。

首先需要把 Spring Security 引入到 Web 应用中来，这是通过在 `web.xml` 添加一个新的过滤器来实现的，如代码清单 1 所示。

##### 清单 1\. 在 web.xml 中添加 Spring Security 的过滤器

```
<filter>
    <filter-name>springSecurityFilterChain</filter-name>
    <filter-class>org.springframework.web.filter.DelegatingFilterProxy</filter-class>
</filter>

<filter-mapping>
    <filter-name>springSecurityFilterChain</filter-name>
    <url-pattern>/*</url-pattern>
</filter-mapping>

```

Show moreShow more icon

Spring Security 使用的是 Servlet 规范中标准的过滤器机制。对于特定的请求，Spring Security 的过滤器会检查该请求是否通过认证，以及当前用户是否有足够的权限来访问此资源。对于非法的请求，过滤器会跳转到指定页面让用户进行认证，或是返回出错信息。需要注意的是，代码清单 1 中虽然只定义了一个过滤器，Spring Security 实际上是使用多个过滤器形成的链条来工作的。

下一步是配置 Spring Security 来声明系统中的合法用户及其对应的权限。用户相关的信息是通过 `org.springframework.security.core.userdetails.UserDetailsService` 接口来加载的。该接口的唯一方法是 `loadUserByUsername(String username)` ，用来根据用户名加载相关的信息。这个方法的返回值是 `org.springframework.security.core.userdetails.UserDetails` 接口，其中包含了用户的信息，包括用户名、密码、权限、是否启用、是否被锁定、是否过期等。其中最重要的是用户权限，由 `org.springframework.security.core.GrantedAuthority` 接口来表示。虽然 Spring Security 内部的设计和实现比较复杂，但是一般情况下，开发人员只需要使用它默认提供的实现就可以满足绝大多数情况下的需求，而且只需要简单的配置声明即可。

在第一个示例应用中，使用的是数据库的方式来存储用户的信息。Spring Security 提供了 `org.springframework.security.core.userdetails.jdbc.JdbcDaoImpl` 类来支持从数据库中加载用户信息。开发人员只需要使用与该类兼容的数据库表结构，就可以不需要任何改动，而直接使用该类。代码清单 2 中给出了相关的配置。

##### 清单 2\. 声明使用数据库来保存用户信息

```
<bean id="dataSource"
    class="org.springframework.jdbc.datasource.DriverManagerDataSource">
    <property name="driverClassName" value="org.apache.derby.jdbc.ClientDriver" />
    <property name="url" value="jdbc:derby://localhost:1527/mycompany" />
    <property name="username" value="app" />
    <property name="password" value="admin" />
</bean>

<bean id="userDetailsService"
    class="org.springframework.security.core.userdetails.jdbc.JdbcDaoImpl">
    <property name="dataSource" ref="dataSource" />
</bean>

<sec:authentication-manager>
    <sec:authentication-provider user-service-ref="userDetailsService" />
</sec:authentication-manager>

```

Show moreShow more icon

如代码清单 2 所示，首先定义了一个使用 Apache Derby 数据库的数据源，Spring Security 的 `org.springframework.security.core.userdetails.jdbc.JdbcDaoImpl` 类使用该数据源来加载用户信息。最后需要配置认证管理器使用该 `UserDetailsService` 。

接着就可以配置用户对不同资源的访问权限了。这里的资源指的是 URL 地址。配置的内容如代码清单 3 所示。 `sec` 是 Spring Security 的配置元素所在的名称空间的前缀。

##### 清单 3\. 配置对不同 URL 模式的访问权限

```
<sec:http>
    <sec:intercept-url pattern="/president_portal.do**" access="ROLE_PRESIDENT" />
    <sec:intercept-url pattern="/manager_portal.do**" access="ROLE_MANAGER" />
    <sec:intercept-url pattern="/**" access="ROLE_USER" />
    <sec:form-login />
    <sec:logout />
</sec:http>

```

Show moreShow more icon

第一个示例应用中一共定义了三种角色：普通用户、经理和总裁，分别用 `ROLE_USER` 、 `ROLE_MANAGER` 和 `ROLE_PRESIDENT` 来表示。代码清单 3 中定义了访问不同的 URL 模式的用户所需要的角色。这是通过 `<sec:intercept-url>` 元素来实现的，其属性 `pattern` 声明了请求 URL 的模式，而属性 `access` 则声明了访问此 URL 时所需要的权限。需要按照 URL 模式从精确到模糊的顺序来进行声明。因为 Spring Security 是按照声明的顺序逐个进行比对的，只要用户当前访问的 URL 符合某个 URL 模式声明的权限要求，该请求就会被允许。如果把代码清单 3 中本来在最后的 URL 模式 `/**` 声明放在最前面，那么当普通用户访问 `/manager_portal.do` 的时候，该请求也会被允许。这显然是不对的。通过 `<sec:form-login>` 元素声明了使用 HTTP 表单验证。也就是说，当未认证的用户试图访问某个受限 URL 的时候，浏览器会跳转到一个登录页面，要求用户输入用户名和密码。 `<sec:logout>` 元素声明了提供用户注销登录的功能。默认的注销登录的 URL 是 `/j_spring_security_logout` ，可以通过属性 `logout-url` 来修改。

当完成这些配置并运行应用之后，会发现 Spring Security 已经默认提供了一个登录页面的实现，可以直接使用。开发人员也可以对登录页面进行定制。通过 `<form-login>` 的属性 `login-page` 、 `login-processing-url` 和 `authentication-failure-url` 就可以定制登录页面的 URL、登录请求的处理 URL 和登录出现错误时的 URL 等。从这里可以看出，一方面 Spring Security 对开发中经常会用到的功能提供了很好的默认实现，另外一方面也提供了非常灵活的定制能力，允许开发人员提供自己的实现。

在介绍如何用 Spring Security 实现基本的用户认证和授权之后，下面介绍其中的核心对象。

## SecurityContext 和 Authentication 对象

下面开始讨论几个 Spring Security 里面的核心对象。 `org.springframework.security.core.context.SecurityContext` 接口表示的是当前应用的安全上下文。通过此接口可以获取和设置当前的认证对象。 `org.springframework.security.core.Authentication` 接口用来表示此认证对象。通过认证对象的方法可以判断当前用户是否已经通过认证，以及获取当前认证用户的相关信息，包括用户名、密码和权限等。要使用此认证对象，首先需要获取到 `SecurityContext` 对象。通过 `org.springframework.security.core.context.SecurityContextHolder` 类提供的静态方法 `getContext()` 就可以获取。再通过 `SecurityContext` 对象的 `getAuthentication()` 就可以得到认证对象。通过认证对象的 `getPrincipal()` 方法就可以获得当前的认证主体，通常是 `UserDetails` 接口的实现。联系到上一节介绍的 `UserDetailsService` ，典型的认证过程就是当用户输入了用户名和密码之后， `UserDetailsService` 通过用户名找到对应的 `UserDetails` 对象，接着比较密码是否匹配。如果不匹配，则返回出错信息；如果匹配的话，说明用户认证成功，就创建一个实现了 `Authentication` 接口的对象，如 `org.springframework.security.authentication.UsernamePasswordAuthenticationToken` 类的对象。再通过 `SecurityContext` 的 `setAuthentication()` 方法来设置此认证对象。

代码清单 4 给出了使用 `SecurityContext` 和 `Authentication` 的一个示例，用来获取当前认证用户的用户名。

##### 清单 4\. 获取当前认证用户的用户名

```
public static String getAuthenticatedUsername() {
    String username = null;
    Object principal = SecurityContextHolder.getContext()
        .getAuthentication().getPrincipal();
    if (principal instanceof UserDetails) {
        username = ((UserDetails) principal).getUsername();
    } else {
        username = principal.toString();
    }
    return username;
}

```

Show moreShow more icon

默认情况下， `SecurityContextHolder` 使用 `ThreadLocal` 来保存 `SecurityContext` 对象。因此， `SecurityContext` 对象对于当前线程上所有方法都是可见的。这种实现对于 Web 应用来说是合适的。不过在有些情况下，如桌面应用，这种实现方式就不适用了。Spring Security 允许开发人员对此进行定制。开发人员只需要实现接口 `org.springframework.security.core.context.SecurityContextHolderStrategy` 并通过 `SecurityContextHolder` 的 `setStrategyName(String)` 方法让 Spring Security 使用此实现即可。另外一种设置方式是使用系统属性。除此之外，Spring Security 默认提供了另外两种实现方式： `MODE_GLOBAL` 表示当前应用共享唯一的 `SecurityContextHolder` ； `MODE_INHERITABLETHREADLOCAL` 表示子线程继承父线程的 `SecurityContextHolder` 。 代码清单 5 给出了使用全局唯一的`SecurityContextHolder` 的示例。

##### 清单 5\. 使用全局唯一的 SecurityContextHolder

```
public void useGlobalSecurityContextHolder() {
    SecurityContextHolder.setStrategyName(SecurityContextHolder.MODE_GLOBAL);
}

```

Show moreShow more icon

在介绍完 Spring Security 中的 `SecurityContext` 和 `Authentication` 之后，下面介绍如何保护服务层的方法。

## 服务层方法保护

之前章节中介绍的是在 URL 这个粒度上的安全保护。这种粒度的保护在很多情况下是不够的。比如相同的 URL 对应的页面上，不同角色的用户所能看到的内容和执行的操作是有可能不同的。在第一个示例应用中，系统中记录了每个员工的工资收入。所有员工都可以查看自己的工资，但是只有员工的直接经理才可以修改员工的工资。这就涉及到对应用中服务层的方法进行相应的权限控制，从而避免安全漏洞。

保护服务层方法涉及到对应用中的方法调用进行拦截。通过 Spring 框架提供的良好面向方面编程（AOP）的支持，可以很容易的对方法调用进行拦截。Spring Security 利用了 AOP 的能力，允许以声明的方式来定义调用方式时所需的权限。代码清单 6 中给出了对方法调用进行保护的配置文件示例。

##### 清单 6\. 对方法调用进行保护

```
<bean id="userSalarySecurity"
    class="org.springframework.security.access.intercept.aspectj.
        AspectJMethodSecurityInterceptor">
    <property name="authenticationManager" ref="authenticationManager" />
    <property name="accessDecisionManager" ref="accessDecisionManager" />
    <property name="securityMetadataSource">
        <value>
            mycompany.service.UserService.raiseSalary=ROLE_MANAGER
        </value>
    </property>
</bean>

```

Show moreShow more icon

如代码清单 6 所示，通过 `mycompany.service.UserService.raiseSalary=ROLE_MANAGER` 声明了 `mycompany.service.UserService` 类的 `raiseSalary` 方法只有具有角色 `ROLE_MANAGER` 的用户才能执行。这就使得只具有角色 `ROLE_USER` 的用户无法调用此方法。

不过仅对方法名称进行权限控制并不能解决另外的一些问题。比如在第一个示例应用中的增加工资的实现是通过发送 HTTP POST 请求到 `salary.do` 这个 URL 来完成的。 `salary.do` 对应的控制器 `mycompany.controller.SalaryController` 会调用 `mycompany.service.UserService` 类的 `raiseSalary` 方法来完成增加工资的操作。存在的一种安全漏洞是具有 `ROLE_MANAGER` 角色的用户可以通过其它工具（如 cURL 或 Firefox 扩展 Poster 等）来创建 HTTP POST 请求来更改其它员工的工资。为了解决这个问题，需要对 `raiseSalary` 的调用进行更加细粒度的控制。通过 Spring Security 提供的 AspectJ 支持就可以编写相关的控制逻辑，如代码清单 7 所示。

##### 清单 7\. 使用 AspectJ 进行细粒度的控制

```
public aspect SalaryManagementAspect {
    private AspectJMethodSecurityInterceptor securityInterceptor;

    private UserDao userDao;

    pointcut salaryChange(): target(UserService)
        && execution(public void raiseSalary(..)) &&!within(SalaryManagementAspect);

    Object around(): salaryChange() {
        if (this.securityInterceptor == null) {
            return proceed();
        }
        AspectJCallback callback = new AspectJCallback() {
            public Object proceedWithObject() {
                return proceed();
            }
        };
        Object[] args = thisJoinPoint.getArgs();
        String employee = (String) args[0]; // 要修改的员工的用户名
        User user = userDao.getByUsername(employee);
        String currentUser = UsernameHolder.getAuthenticatedUsername(); // 当前登录用户
        if (!currentUser.equals(user.getManagerId())) {
            throw new AccessDeniedException
                ("Only the direct manager can change the salary.");
        }

        return this.securityInterceptor.invoke(thisJoinPoint, callback);
    }
}

```

Show moreShow more icon

如代码清单 7 所示，定义了一个切入点（pointcut）`salaryChange` 和对应的环绕增强。当方法 `raiseSalary` 被调用的时候，会比较要修改的员工的经理的用户名和当前登录用户的用户名是否一致。当不一致的时候就会抛出 `AccessDeniedException` 异常。

在介绍了如何保护方法调用之后，下面介绍如何通过访问控制列表来保护领域对象。

## 访问控制列表

之前提到的安全保护和权限控制都是只针对 URL 或是方法调用，只对一类对象起作用。而在有些情况下，不同领域对象实体所要求的权限控制是不同的。以第一类示例应用来说，系统中有报表这一类实体。由于报表的特殊性，只有具有角色 `ROLE_PRESIDENT` 的用户才可以创建报表。对于每份报表，创建者可以设定其对于不同用户的权限。比如有的报表只允许特定的几个用户可以查看。对于这样的需求，就需要对每个领域对象的实例设置对应的访问控制权限。Spring Security 提供了对访问控制列表（Access Control List，ACL）的支持，可以很方便的对不同的领域对象设置针对不同用户的权限。

Spring Security 中的访问控制列表的实现中有 3 个重要的概念，对应于 4 张数据库表。

- 授权的主体：一般是系统中的用户。由 `ACL_SID` 表来表示。
- 领域对象：表示系统中需要进行访问控制的实体。由 `ACL_CLASS` 和 `ACL_OBJECT_IDENTITY` 表来表示，前者保存的是实体所对应的 Java 类的名称，而后者保存的是实体本身。
- 访问权限：表示一个用户对一个领域对象所具有的权限。由表 `ACL_ENTRY` 来表示。

Spring Security 已经提供了参考的数据库表模式和相应的基于 JDBC 的实现。在大多数情况下，使用参考实现就可以满足需求了。类 `org.springframework.security.acls.jdbc.JdbcMutableAclService` 可以对访问控制列表进行查询、添加、更新和删除的操作，是开发人员最常直接使用的类。该类的构造方法需要 3 个参数，分别是 `javax.sql.DataSource` 表示的数据源、 `org.springframework.security.acls.jdbc.LookupStrategy` 表示的数据库的查询策略和 `org.springframework.security.acls.model.AclCache` 表示的访问控制列表缓存。数据源可以使用第一个示例应用中已有的数据源。查询策略可以使用默认的实现 `org.springframework.security.acls.jdbc.BasicLookupStrategy`。缓存可以使用基于 EhCache 的缓存实现 `org.springframework.security.acls.domain.EhCacheBasedAclCache`。代码清单 8 中给出了相关代码。

##### 清单 8\. 使用 JDBC 的访问控制列表服务基本配置

```
<bean id="aclService"
    class="org.springframework.security.acls.jdbc.JdbcMutableAclService">
    <constructor-arg ref="dataSource" />
    <constructor-arg ref="lookupStrategy" />
    <constructor-arg ref="aclCache" />
    <property name="classIdentityQuery" value="values IDENTITY_VAL_LOCAL()"/>
    <property name="sidIdentityQuery" value="values IDENTITY_VAL_LOCAL()"/>
</bean>

```

Show moreShow more icon

如代码清单 8 所示，需要注意的是 `org.springframework.security.acls.jdbc.JdbcMutableAclService` 的属性 `classIdentityQuery` 和 `sidIdentityQuery` 。Spring Security 的默认数据库模式使用了自动增长的列作为主键。而在实现中，需要能够获取到新插入的列的 ID。因此需要与数据库实现相关的 SQL 查询语言来获取到这个 ID。Spring Security 默认使用的 HSQLDB，因此这两个属性的默认值是 HSQLDB 支持的 `call identity()` 。如果使用的数据库不是 HSQLDB 的话，则需要根据数据库实现来设置这两个属性的值。第一个示例应用使用的是 Apache Derby 数据库，因此这两个属性的值是 `values IDENTITY_VAL_LOCAL()` 。对于 MySQL 来说，这个值是 `select @@identity` 。 代码清单 9 给出了使用 `org.springframework.security.acls.jdbc.JdbcMutableAclService` 来管理访问控制列表的 Java 代码。

##### 清单 9\. 使用访问控制列表服务

```
public void createNewReport(String title, String content) throws ServiceException {
    final Report report = new Report();
    report.setTitle(title);
    report.setContent(content);

    transactionTemplate.execute(new TransactionCallback<Object>() {
        public Object doInTransaction(TransactionStatus status) {
            reportDao.create(report);
            addPermission(report.getId(), new PrincipalSid(getUsername()),
                BasePermission.ADMINISTRATION);
            return null;
        }
    });
}

public void grantRead(final String username, final Long reportId) {
    transactionTemplate.execute(new TransactionCallback<Object>() {
        public Object doInTransaction(TransactionStatus status) {
            addPermission(reportId, new PrincipalSid(username), BasePermission.READ);
            return null;
        }
    });
}

private void addPermission(Long reportId, Sid recipient, Permission permission) {
    MutableAcl acl;
    ObjectIdentity oid = new ObjectIdentityImpl(Report.class, reportId);

    try {
        acl = (MutableAcl) mutableAclService.readAclById(oid);
    } catch (NotFoundException nfe) {
        acl = mutableAclService.createAcl(oid);
    }

    acl.insertAce(acl.getEntries().size(), permission, recipient, true);
    mutableAclService.updateAcl(acl);
}

```

Show moreShow more icon

代码清单 9 中的 `addPermission(Long reportId, Sid recipient, Permission permission)` 方法用来为某个报表添加访问控制权限，参数 `reportId` 表示的是报表的 ID，用来标识一个报表； `recipient` 表示的是需要授权的用户； `permission` 表示的是授予的权限。 `createNewReport()` 方法用来创建一个报表，同时给创建报表的用户授予管理权限（ `BasePermission.ADMINISTRATION` ）。 `grantRead()` 方法用来给某个用户对某个报表授予读权限（ `BasePermission.READ` ）。这里需要注意的是，对访问控制列表的操作都需要在一个事务中进行处理。利用 Spring 提供的事务模板（ `org.springframework.transaction.support.TransactionTemplate` ）就可以很好的处理事务。对于权限，Spring Security 提供了 4 种基本的权限：读、写、删除和管理。开发人员可以在这基础上定义自己的权限。

在介绍完访问控制列表之后，下面介绍 Spring Security 提供的 JSP 标签库。

## JSP 标签库

之前的章节中介绍了在 Java 代码中如何使用 Spring Security 提供的能力。很多情况下，用户可能有权限访问某个页面，但是页面上的某些功能对他来说是不可用的。比如对于同样的员工列表，普通用户只能查看数据，而具有经理角色的用户则可以看到对列表进行修改的链接或是按钮等。Spring Security 提供了一个 JSP 标签库用来方便在 JSP 页面中根据用户的权限来控制页面某些部分的显示和隐藏。使用这个 JSP 标签库很简单，只需要在 JSP 页面上添加声明即可： `<%@ taglib prefix="sec" uri="http://www.springframework.org/security/tags" %>` 。这个标签库包含如下 3 个标签：

- `authorize` 标签：该标签用来判断其中包含的内容是否应该被显示出来。判断的条件可以是某个表达式的求值结果，或是是否能访问某个 URL，分别通过属性 `access` 和 `url` 来指定。如 `<sec:authorize access="hasRole('ROLE_MANAGER')">` 限定内容只有具有经理角色的用户才可见。 `<sec:authorize url="/manager_portal.do">` 限定内容只有能访问 URL `/manager_portal.do` 的用户才可见。
- `authentication` 标签：该标签用来获取当前认证对象（ `Authentication` ）中的内容。如 `<sec:authentication property="principal.username" />` 可以用来获取当前认证用户的用户名。
- `accesscontrollist` 标签：该标签的作用与 `authorize` 标签类似，也是判断其中包含的内容是否应该被显示出来。所不同的是它是基于访问控制列表来做判断的。该标签的属性 `domainObject` 表示的是领域对象，而属性 `hasPermission` 表示的是要检查的权限。如 `<sec:accesscontrollist hasPermission="READ" domainObject="myReport">` 限定了其中包含的内容只在对领域对象 `myReport` 有读权限的时候才可见。

值得注意的是，在使用 `authorize` 标签的时候，需要通过 `<sec:http use-expressions="true">` 来启用表达式的支持。查看 [权限控制表达式](#权限控制表达式) 一节了解关于表达式的更多内容。

在介绍完 JSP 标签库之后，下面介绍如何与 LDAP 进行集成。

## 使用 LDAP

很多公司都使用 LDAP 服务器来保存员工的相关信息。内部的 IT 系统都需要与 LDAP 服务器做集成来进行用户认证与访问授权。Spring Security 提供了对 LDAP 协议的支持，只需要简单的配置就可以让 Web 应用使用 LDAP 来进行认证。第二个示例应用使用 OpenDS LDAP 服务器并添加了一些测试用户。代码清单 10 中给出了配置文件的示例，完整的代码见 参考资料 。

##### 清单 10\. 集成 LDAP 服务器的配置文件

```
<bean id="contextSource"
    class="org.springframework.security.ldap.DefaultSpringSecurityContextSource">
    <constructor-arg value="ldap://localhost:389" />
</bean>

<bean id="ldapAuthProvider"
    class="org.springframework.security.ldap.authentication.LdapAuthenticationProvider">
    <constructor-arg>
        <bean class="org.springframework.security.ldap.authentication.BindAuthenticator">
            <constructor-arg ref="contextSource" />
            <property name="userSearch">
                <bean id="userSearch"
        class="org.springframework.security.ldap.search.FilterBasedLdapUserSearch">
                    <constructor-arg index="0" value="ou=People,dc=mycompany,dc=com" />
                    <constructor-arg index="1"
                value="(&(uid={0})(objectclass=person))" />
                    <constructor-arg index="2" ref="contextSource" />
                </bean>
            </property>
        </bean>
    </constructor-arg>
    <constructor-arg>
        <bean class="mycompany.CompanyAuthoritiesPopulator"></bean>
    </constructor-arg>
</bean>

<sec:authentication-manager>
    <sec:authentication-provider ref="ldapAuthProvider" />
</sec:authentication-manager>

```

Show moreShow more icon

如代码清单 10 所示，配置中的核心部分是类 `org.springframework.security.ldap.authentication.LdapAuthenticationProvider`，它用来与 LDAP 服务器进行认证以及获取用户的权限信息。一般来说，与 LDAP 服务器进行认证的方式有两种。一种是使用用户提供的用户名和密码直接绑定到 LDAP 服务器；另外一种是比较用户提供的密码与 LDAP 服务器上保存的密码是否一致。前者通过类 `org.springframework.security.ldap.authentication.BindAuthenticator` 来实现，而后者通过类 `org.springframework.security. ldap.authentication.PasswordComparisonAuthenticator` 来实现。第二个示例应用中使用的是绑定的方式来进行认证。在进行绑定的时候，需要在 LDAP 服务器上搜索当前的用户。搜索的时候需要指定基本的识别名（Distinguished Name）和过滤条件。在该应用中，用户登录时使用的是其唯一识别符（ `uid` ），如 `user.0` ，而在 LDAP 服务器上对应的识别名是 `uid=user.0,ou=People,dc=mycompany,dc=com` 。通过使用过滤条件 `(&(uid={0})(objectclass=person))` 就可以根据 `uid` 来搜索到用户并进行绑定。当认证成功之后，就需要获取到该用户对应的权限。一般是通过该用户在 LDAP 服务器上所在的分组来确定的。不过在示例应用中展示了如何提供自己的实现来为用户分配权限。类 `mycompany.CompanyAuthoritiesPopulator` 实现了 `org.springframework.security.ldap.userdetails.LdapAuthoritiesPopulator` 接口，并为所有的用户分配了单一的角色 `ROLE_USER` 。

在介绍完与 LDAP 进行集成之后，下面介绍如何与 OAuth 进行集成。

## OAuth 集成

现在的很多 Web 服务都提供 API 接口，允许第三方应用使用其数据。当第三方应用需要访问用户私有数据的时候，需要进行认证。OAuth 是目前流行的一种认证方式，被很多 Web 服务采用，包括 Twitter、LinkedIn、Google Buzz 和新浪微博等。OAuth 的特点是第三方应用不能直接获取到用户的密码，而只是使用一个经过用户授权之后的令牌（token）来进行访问。用户可以对能够访问其数据的第三方应用进行管理，通过回收令牌的方式来终止第三方应用对其数据的访问。OAuth 的工作方式涉及到服务提供者、第三方应用和用户等 3 个主体。其基本的工作流程是：第三方应用向服务提供者发出访问用户数据的请求。服务提供者会询问用户是否同意此请求。如果用户同意的话，服务提供者会返回给第三方应用一个令牌。第三方应用只需要在请求数据的时候带上此令牌就可以成功获取。

第三方应用在使用 OAuth 认证方式的时候，其中所涉及的交互比较复杂。Spring Security 本身并没有提供 OAuth 的支持，通过另外一个开源库 OAuth for Spring Security 可以实现。OAuth for Spring Security 与 Spring Security 有着很好的集成，可以很容易在已有的使用 Spring Security 的应用中添加 OAuth 的支持。不过目前 OAuth for Spring Security 只对 Spring Security 2.0.x 版本提供比较好的支持。对 OAuth 的支持包括服务提供者和服务消费者两个部分：服务提供者是数据的提供者，服务消费者是使用这些数据的第三方应用。一般的应用都是服务消费者。OAuth for Spring Security 对服务提供者和消费者都提供了支持。下面通过获取 LinkedIn 上的状态更新的示例来说明其用法。

作为 OAuth 的服务消费者，需要向服务提供者申请表示其应用的密钥。服务提供者会提供 3 个 URL 来与服务消费者进行交互。代码清单 11 中给出了使用 OAuth for Spring Security 的配置文件。

##### 清单 11\. 使用 OAuth for Spring Security 的配置文件

```
<oauth:consumer resource-details-service-ref="linkedInResourceDetails"
    oauth-failure-page="/oauth_error.jsp">
    <oauth:url pattern="/linkedin.do**" resources="linkedIn" />
</oauth:consumer>

<bean id="oauthConsumerSupport"
    class="org.springframework.security.oauth.consumer.CoreOAuthConsumerSupport">
    <property name="protectedResourceDetailsService" ref="linkedInResourceDetails" />
</bean>

<oauth:resource-details-service id="linkedInResourceDetails">
    <oauth:resource id="linkedIn"
        key="***" secret="***"
        request-token-url="https://api.linkedin.com/uas/oauth/requestToken"
        user-authorization-url="https://www.linkedin.com/uas/oauth/authorize"
        access-token-url="https://api.linkedin.com/uas/oauth/accessToken" />
</oauth:resource-details-service>

```

Show moreShow more icon

如代码清单 11 所示，只需要通过对 `<oauth:resource>` 元素进行简单的配置，就可以声明使用 LinkedIn 的服务。每个 `<oauth:resource>` 元素对应一个 OAuth 服务资源。该元素的属性包含了与该服务资源相关的信息。OAuth for Spring Security 在 Spring Security 提供的过滤器的基础上，额外增加了处理 OAuth 认证的过滤器实现。通过 `<oauth:consumer>` 的子元素 `<oauth:url>` 可以定义过滤器起作用的 URL 模式和对应的 OAuth 服务资源。当用户访问指定的 URL 的时候，应用会转到服务提供者的页面，要求用户进行授权。当用户授权之后，应用就可以访问其数据。访问数据的时候，需要在 HTTP 请求中添加额外的 `Authorization` 头。代码清单 12 给出了访问数据时使用的代码。

##### 清单 12\. 获取访问令牌和构建 HTTP 请求

```
public OAuthConsumerToken getAccessTokenFromRequest(HttpServletRequest request) {
    OAuthConsumerToken token = null;

    List<OAuthConsumerToken> tokens = (List<OAuthConsumerToken>) request
        .getAttribute(OAuthConsumerProcessingFilter.ACCESS_TOKENS_DEFAULT_ATTRIBUTE);
    if (tokens != null) {
        for (OAuthConsumerToken consumerToken : tokens) {
            if (consumerToken.getResourceId().equals(resourceId)) {
                token = consumerToken;
                break;
            }
        }
    }
    return token;
}

public GetMethod getGetMethod(OAuthConsumerToken accessToken, URL url) {
    GetMethod method = new GetMethod(url.toString());
    method.setRequestHeader("Authorization",
                 getHeader(accessToken, url, "GET"));
    return method;
}

public String getHeader(OAuthConsumerToken accessToken, URL url,
             String method) {
    ProtectedResourceDetails details = support
        .getProtectedResourceDetailsService()
        .loadProtectedResourceDetailsById(accessToken.getResourceId());
    return support.getAuthorizationHeader(details, accessToken, url, method, null);
}

```

Show moreShow more icon

如代码清单 12 所示，OAuth for Spring Security 的过滤器会把 OAuth 认证成功之后的令牌保存在当前的请求中。通过 `getAccessTokenFromRequest()` 方法就可以从请求中获取到此令牌。有了这个令牌之后，就可以通过 `getHeader()` 方法构建出 HTTP 请求所需的 `Authorization` 头。只需要在请求中添加此 HTTP 头，就可以正常访问到所需的数据。默认情况下，应用的 OAuth 令牌是保存在 HTTP 会话中的，开发人员可以提供其它的令牌保存方式，如保存在数据库中。只需要提供 `org.springframework.security.oauth.consumer.token.OAuthConsumerTokenServices` 接口的实现就可以了。

在介绍完与 OAuth 的集成方式之后，下面介绍一些高级话题。

## 高级话题

这些与 Spring Security 相关的高级话题包括权限控制表达式、会话管理和记住用户等。

### 权限控制表达式

有些情况下，对于某种资源的访问条件可能比较复杂，并不只是简单的要求当前用户具有某一个角色即可，而是由多种条件进行组合。权限控制表达式允许使用一种简单的语法来描述比较复杂的授权条件。Spring Security 内置了一些常用的表达式，包括 `hasRole()` 用来判断当前用户是否具有某个角色， `hasAnyRole()` 用来判断当前用户是否具备列表中的某个角色，以及 `hasPermission()` 用来判断当前用户是否具备对某个领域对象的某些权限等。这些基本表达式可以通过 `and` 和 `or` 等组合起来，表示复杂的语义。当通过 `<sec:http use-expressions="true">` 启用了表达式支持之后，就可以在 `<sec:intercept-url>` 元素的 `access` 属性上使用表达式。

表达式还可以用来对方法调用进行权限控制，主要是用在方法注解中。要启用 Spring Security 提供的方法注解，需要添加元素 `<global-method-security pre-post-annotations="enabled"/>` 。这几个方法注解分别是：

- `@PreAuthorize` ：该注解用来确定一个方法是否应该被执行。该注解后面跟着的是一个表达式，如果表达式的值为真，则该方法会被执行。如 `@PreAuthorize("hasRole('ROLE_USER')")` 就说明只有当前用户具有角色 `ROLE_USER` 的时候才会执行。
- `@PostAuthorize` ：该注解用来在方法执行完之后进行访问控制检查。
- `@PostFilter` ：该注解用来对方法的返回结果进行过滤。从返回的集合中过滤掉表达式值为假的元素。如 `@PostFilter("hasPermission(filterObject, 'read')")` 说明返回的结果中只保留当前用户有读权限的元素。
- `@PreFilter` ：该注解用来对方法调用时的参数进行过滤。

### 会话管理

Spring Security 提供了对 HTTP 会话的管理功能。这些功能包括对会话超时的管理、防范会话设置攻击（Session fixation attack）和并发会话管理等。

如果当前用户的会话因为超时而失效之后，如果用户继续使用此会话来访问，Spring Security 可以检测到这种情况，并跳转到适当的页面。只需要在 `<sec:http>` 元素下添加 `<sec:session-management invalid-session-url="/sessionTimeout.jsp" />` 元素即可，属性 `invalid-session-url` 指明了会话超时之后跳转到的 URL 地址。

有些 Web 应用会把用户的会话标识符直接通过 URL 的参数来传递，并且在服务器端不进行验证，如用户访问的 URL 可能是 `/myurl;jsessionid=xxx` 。攻击者可以用一个已知的会话标识符来构建一个 URL，并把此 URL 发给要攻击的对象。如果被攻击者访问这个 URL 并用自己的用户名登录成功之后，攻击者就可以利用这个已经通过认证的会话来访问被攻击者的数据。防范这种攻击的办法就是要求用户在做任何重要操作之前都重新认证。Spring Security 允许开发人员定制用户登录时对已有会话的处理，从而可以有效的防范这种攻击。通过 `<sec:session-management>` 元素的属性 `session-fixation-protection` 可以修改此行为。该属性的可选值有 `migrateSession` 、`newSession` 和 `none`。 `migrateSession` 是默认值。在这种情况下，每次用户登录都会创建一个新的会话，同时把之前会话的数据复制到新会话中。 `newSession` 表示的是只创建新的会话，而不复制数据。 `none` 表示的是保持之前的会话。

在有些情况下，应用需要限定使用同一个用户名同时进行登录所产生的会话数目。比如有些应用可能要求每个用户在同一时间最多只能有一个会话。可以通过 `<sec:session-management>` 元素的子元素 `<sec:concurrency-control>` 来限制每个用户的并发会话个数。如 `<sec:concurrency-control max-sessions="2" />` 就限定了每个用户在同一时间最多只能有两个会话。如果当前用户的会话数目已经达到上限，而用户又再次登录的话，默认的实现是使之前的会话失效。如果希望阻止后面的这次登录的话，可以设置属性 `error-if-maximum-exceeded` 的值为 `true`。这样的话，后面的这次登录就会出错。只有当之前的会话失效之后，用户才能再次登录。

### 记住用户

有些 Web 应用会在登录界面提供一个复选框，询问用户是否希望在当前计算机上记住自己的密码。如果用户勾选此选项的话，在一段时间内用户访问此应用时，不需要输入用户名和密码进行登录。Spring Security 提供了对这种记住用户的需求的支持。只需要在 `<sec:http>` 中添加 `<sec:remember-me>` 元素即可。

一般来说，有两种方式可以实现记住用户的能力。一种做法是利用浏览器端的 cookie。当用户成功登录之后，特定内容的字符串被保存到 cookie 中。下次用户再次访问的时候，保存在 cookie 中的内容被用来认证用户。默认情况下使用的是这种方式。使用 cookie 的做法存在安全隐患，比如攻击者可能窃取用户的 cookie，并用此 cookie 来登录系统。另外一种更安全的做法是浏览器端的 cookie 只保存一些随机的数字，而且这些数字只能使用一次，在每次用户登录之后都会重新生成。这些数字保存在服务器端的数据库中。如果希望使用这种方式，需要创建一个数据库表，并通过 `data-source-ref` 属性来指定包含此表的数据源。

## 结束语

对于使用 Spring 开发的 Web 应用来说，Spring Security 是增加安全性时的最好选择。本文详细介绍了 Spring Security 的各个方面，包括实现基本的用户认证和授权、保护服务层方法、使用访问控制列表保护具体的领域对象、JSP 标签库和与 LDAP 和 OAuth 的集成等。通过本文，开发人员可以了解如何使用 Spring Security 来实现不同的用户认证和授权机制。

## 下载示例代码

[Company.zip](http://www.ibm.com/developerworks/cn/java/j-lo-springsecurity/Company.zip): 简单的企业员工管理系统