# Spring 事务管理高级应用难点剖析，第 3 部分
走出迷茫：数据访问技术数据连接泄漏的应对方案

**标签:** Spring

[原文链接](https://developer.ibm.com/zh/articles/j-lo-spring-ts3/)

陈雄华

发布: 2010-03-29

* * *

## 概述

对于应用开发者来说，数据连接泄漏无疑是一个可怕的梦魇。如果存在数据连接泄漏问题，应用程序将因数据连接资源的耗尽而崩溃，甚至还可能引起数据库的崩溃。数据连接泄漏像黑洞一样让开发者避之唯恐不及。

Spring DAO 对所有支持的数据访问技术框架都使用模板化技术进行了薄层的封装。只要您的程序都使用 Spring DAO 模板（如 JdbcTemplate、HibernateTemplate 等）进行数据访问，一定不会存在数据连接泄漏的问题 ―― 这是 Spring 给予我们郑重的承诺！因此，我们无需关注数据连接（Connection）及其衍生品（Hibernate 的 Session 等）的获取和释放的操作，模板类已经通过其内部流程替我们完成了，且对开发者是透明的。

但是由于集成第三方产品，整合遗产代码等原因，可能需要直接访问数据源或直接获取数据连接及其衍生品。这时，如果使用不当，就可能在无意中创造出一个魔鬼般的连接泄漏问题。

我们知道：当 Spring 事务方法运行时，就产生一个事务上下文，该上下文在本事务执行线程中针对同一个数据源绑定了一个唯一的数据连接（或其衍生品），所有被该事务上下文传播的方法都共享这个数据连接。这个数据连接从数据源获取及返回给数据源都在 Spring 掌控之中，不会发生问题。如果在需要数据连接时，能够获取这个被 Spring 管控的数据连接，则使用者可以放心使用，无需关注连接释放的问题。

那么，如何获取这些被 Spring 管控的数据连接呢？ Spring 提供了两种方法：其一是使用数据资源获取工具类，其二是对数据源（或其衍生品如 Hibernate SessionFactory）进行代理。在具体介绍这些方法之前，让我们先来看一下各种引发数据连接泄漏的场景。

## Spring JDBC 数据连接泄漏

如果直接从数据源获取连接，且在使用完成后不主动归还给数据源（调用 Connection#close()），则将造成数据连接泄漏的问题。

### 一个具体的实例

下面，来看一个具体的实例：

##### 清单 1.JdbcUserService.java：主体代码

```
package user.connleak;
import org.apache.commons.dbcp.BasicDataSource;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.ApplicationContext;
import org.springframework.context.support.ClassPathXmlApplicationContext;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Service;
import java.sql.Connection;

@Service("jdbcUserService")
public class JdbcUserService {
    @Autowired
    private JdbcTemplate jdbcTemplate;

    public void logon(String userName) {
        try {
            // ①直接从数据源获取连接，后续程序没有显式释放该连接
            Connection conn = jdbcTemplate.getDataSource().getConnection();
            String sql = "UPDATE t_user SET last_logon_time=? WHERE user_name =?";
            jdbcTemplate.update(sql, System.currentTimeMillis(), userName);
            Thread.sleep(1000);// ②模拟程序代码的执行时间
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}

```

Show moreShow more icon

JdbcUserService 通过 Spring AOP 事务增强的配置，让所有 public 方法都工作在事务环境中。即让 logon() 和 updateLastLogonTime() 方法拥有事务功能。在 logon() 方法内部，我们在①处通过调用 `jdbcTemplate.getDataSource().getConnection()` 显式获取一个连接，这个连接不是 logon() 方法事务上下文线程绑定的连接，所以如果开发者如果没有手工释放这连接（显式调用 Connection#close() 方法），则这个连接将永久被占用（处于 active 状态），造成连接泄漏！下面，我们编写模拟运行的代码，查看方法执行对数据连接的实际占用情况：

##### 清单 2.JdbcUserService.java：模拟运行代码

```
...
@Service("jdbcUserService")
public class JdbcUserService {
...
    //①以异步线程的方式执行JdbcUserService#logon()方法，以模拟多线程的环境
    public static void asynchrLogon(JdbcUserService userService, String userName) {
        UserServiceRunner runner = new UserServiceRunner(userService, userName);
        runner.start();
    }
    private static class UserServiceRunner extends Thread {
        private JdbcUserService userService;
        private String userName;
        public UserServiceRunner(JdbcUserService userService, String userName) {
            this.userService = userService;
            this.userName = userName;
        }
        public void run() {
            userService.logon(userName);
        }
    }

    //② 让主执行线程睡眠一段指定的时间
    public static void sleep(long time) {
        try {
            Thread.sleep(time);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }

    //③ 汇报数据源的连接占用情况
    public static void reportConn(BasicDataSource basicDataSource) {
        System.out.println("连接数[active:idle]-[" +
            basicDataSource.getNumActive()+":"+basicDataSource.getNumIdle()+"]");
    }

    public static void main(String[] args) {
        ApplicationContext ctx =
            new ClassPathXmlApplicationContext("user/connleak/applicatonContext.xml");
        JdbcUserService userService = (JdbcUserService) ctx.getBean("jdbcUserService");

        BasicDataSource basicDataSource = (BasicDataSource) ctx.getBean("dataSource");

        //④汇报数据源初始连接占用情况
        JdbcUserService.reportConn(basicDataSource);

        JdbcUserService.asynchrLogon(userService, "tom");
        JdbcUserService.sleep(500);

        //⑤此时线程A正在执行JdbcUserService#logon()方法
        JdbcUserService.reportConn(basicDataSource);

        JdbcUserService.sleep(2000);
        //⑥此时线程A所执行的JdbcUserService#logon()方法已经执行完毕
        JdbcUserService.reportConn(basicDataSource);

        JdbcUserService.asynchrLogon(userService, "john");
        JdbcUserService.sleep(500);

        //⑦此时线程B正在执行JdbcUserService#logon()方法
        JdbcUserService.reportConn(basicDataSource);

        JdbcUserService.sleep(2000);

        //⑧此时线程A和B都已完成JdbcUserService#logon()方法的执行
        JdbcUserService.reportConn(basicDataSource);
    }

```

Show moreShow more icon

在 JdbcUserService 中添加一个可异步执行 logon() 方法的 asynchrLogon() 方法，我们通过异步执行 logon() 以及让主线程睡眠的方式模拟多线程环境下的执行场景。在不同的执行点，通过 reportConn() 方法汇报数据源连接的占用情况。

使用如下的 Spring 配置文件对 JdbcUserServie 的方法进行事务增强：

##### 清单 3.applicationContext.xml

```
<?xml version="1.0" encoding="UTF-8" ?>
<beans xmlns="http://www.springframework.org/schema/beans"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xmlns:context="http://www.springframework.org/schema/context"
    xmlns:p="http://www.springframework.org/schema/p"
    xmlns:aop="http://www.springframework.org/schema/aop"
    xmlns:tx="http://www.springframework.org/schema/tx"
    xsi:schemaLocation="http://www.springframework.org/schema/beans
        http://www.springframework.org/schema/beans/spring-beans-3.0.xsd
        http://www.springframework.org/schema/context
        http://www.springframework.org/schema/context/spring-context-3.0.xsd
        http://www.springframework.org/schema/aop
        http://www.springframework.org/schema/aop/spring-aop-3.0.xsd
        http://www.springframework.org/schema/tx
        http://www.springframework.org/schema/tx/spring-tx-3.0.xsd">
    <context:component-scan base-package="user.connleak"/>
    <bean id="dataSource"
        class="org.apache.commons.dbcp.BasicDataSource"
            destroy-method="close"
            p:driverClassName="oracle.jdbc.driver.OracleDriver"
            p:url="jdbc:oracle:thin:@localhost:1521:orcl"
            p:username="test"
            p:password="test"
            p:defaultAutoCommit="false"/>

    <bean id="jdbcTemplate"
        class="org.springframework.jdbc.core.JdbcTemplate"
        p:dataSource-ref="dataSource"/>

    <bean id="jdbcManager"
        class="org.springframework.jdbc.datasource.DataSourceTransactionManager"
        p:dataSource-ref="dataSource"/>

    <!-- 对JdbcUserService的所有方法实施事务增强 -->
    <aop:config proxy-target-class="true">
        <aop:pointcut id="serviceJdbcMethod"
            expression="within(user.connleak.JdbcUserService+)"/>
        <aop:advisor pointcut-ref="serviceJdbcMethod"
            advice-ref="jdbcAdvice" order="0"/>
    </aop:config>
    <tx:advice id="jdbcAdvice" transaction-manager="jdbcManager">
        <tx:attributes>
            <tx:method name="*"/>
        </tx:attributes>
    </tx:advice>
</beans>

```

Show moreShow more icon

保证 BasicDataSource 数据源的配置默认连接为 0，运行以上程序代码，在控制台中将输出以下的信息：

##### 清单 4\. 输出日志

```
连接数 [active:idle]-[0:0]
连接数 [active:idle]-[2:0]
连接数 [active:idle]-[1:1]
连接数 [active:idle]-[3:0]
连接数 [active:idle]-[2:1]

```

Show moreShow more icon

我们通过下表对数据源连接的占用和泄漏情况进行描述：

##### 表 1\. 执行过程数据源连接占用情况

**时间****执行线程 1****执行线程 2****数据源连接 (active idle leak)**T0未启动未启动0 0 0T1正在执行方法未启动2 0 0T2执行完毕未启动1 1 1T3执行完毕正式执行方法3 0 1T4执行完毕执行完毕2 1 2

可见在执行线程 1 执行完毕后，只释放了一个数据连接，还有一个数据连处于 active 状态，说明泄漏了一个连接。相似的，执行线程 2 执行完毕后，也泄漏了一个连接：原因是直接通过数据源获取连接（jdbcTemplate.getDataSource().getConnection()）而没有显式释放造成的。

### 通过 DataSourceUtils 获取数据连接

Spring 提供了一个能从当前事务上下文中获取绑定的数据连接的工具类，那就是 DataSourceUtils。Spring 强调必须使用 DataSourceUtils 工具类获取数据连接，Spring 的 JdbcTemplate 内部也是通过 DataSourceUtils 来获取连接的。DataSourceUtils 提供了若干获取和释放数据连接的静态方法，说明如下：

- `static Connection doGetConnection(DataSource dataSource)` ：首先尝试从事务上下文中获取连接，失败后再从数据源获取连接；
- `static Connection getConnection(DataSource dataSource)` ：和 doGetConnection 方法的功能一样，实际上，它内部就是调用 doGetConnection 方法获取连接的；
- `static void doReleaseConnection(Connection con, DataSource dataSource)` ：释放连接，放回到连接池中；
- `static void releaseConnection(Connection con, DataSource dataSource)` ：和 doReleaseConnection 方法的功能一样，实际上，它内部就是调用 doReleaseConnection 方法获取连接的；

来看一下 DataSourceUtils 从数据源获取连接的关键代码：

##### 清单 5\. DataSourceUtils.java 获取连接的工具类

```
public abstract class DataSourceUtils {
...
    public static Connection doGetConnection(DataSource dataSource) throws SQLException {

        Assert.notNull(dataSource, "No DataSource specified");

        //①首先尝试从事务同步管理器中获取数据连接
        ConnectionHolder conHolder =
            (ConnectionHolder) TransactionSynchronizationManager.getResource(dataSource);
        if (conHolder != null && (conHolder.hasConnection() ||
            conHolder.isSynchronizedWithTransaction())) {
            conHolder.requested();
            if (!conHolder.hasConnection()) {
                logger.debug(
                    "Fetching resumed JDBC Connection from DataSource");
                conHolder.setConnection(dataSource.getConnection());
            }
            return conHolder.getConnection();
        }

        //②如果获取不到，则直接从数据源中获取连接
        Connection con = dataSource.getConnection();

        //③如果拥有事务上下文，则将连接绑定到事务上下文中
        if (TransactionSynchronizationManager.isSynchronizationActive()) {
            ConnectionHolder holderToUse = conHolder;
            if (holderToUse == null) {
                holderToUse = new ConnectionHolder(con);
            }
            else {holderToUse.setConnection(con);}
            holderToUse.requested();
            TransactionSynchronizationManager.registerSynchronization(
                new ConnectionSynchronization(holderToUse, dataSource));
            holderToUse.setSynchronizedWithTransaction(true);
            if (holderToUse != conHolder) {
                TransactionSynchronizationManager.bindResource(
                dataSource, holderToUse);
            }
        }
        return con;
    }
...
}

```

Show moreShow more icon

它首先查看当前是否存在事务管理上下文，并尝试从事务管理上下文获取连接，如果获取失败，直接从数据源中获取连接。在获取连接后，如果当前拥有事务上下文，则将连接绑定到事务上下文中。

我们在清单 1 的 JdbcUserService 中，使用 DataSourceUtils.getConnection() 替换直接从数据源中获取连接的代码：

##### 清单 6\. JdbcUserService.java：使用 DataSourceUtils 获取数据连接

```
public void logon(String userName) {
    try {
        //Connection conn = jdbcTemplate.getDataSource().getConnection();
        //①使用DataSourceUtils获取数据连接
        Connection conn = DataSourceUtils.getConnection(jdbcTemplate.getDataSource());
        String sql = "UPDATE t_user SET last_logon_time=? WHERE user_name =?";
        jdbcTemplate.update(sql, System.currentTimeMillis(), userName);
        Thread.sleep(1000);
    } catch (Exception e) {
        e.printStackTrace();
    }
}

```

Show moreShow more icon

重新运行代码，得到如下的执行结果：

##### 清单 7\. 输出日志

```
连接数 [active:idle]-[0:0]
连接数 [active:idle]-[1:0]
连接数 [active:idle]-[0:1]
连接数 [active:idle]-[1:0]
连接数 [active:idle]-[0:1]

```

Show moreShow more icon

对照清单 4 的输出日志，我们可以看到已经没有连接泄漏的现象了。一个执行线程在运行 JdbcUserService#logon() 方法时，只占用一个连接，而且方法执行完毕后，该连接马上释放。这说明通过 DataSourceUtils.getConnection() 方法确实获取了方法所在事务上下文绑定的那个连接，而不是像原来那样从数据源中获取一个新的连接。

### 使用 DataSourceUtils 获取数据连接也可能造成泄漏！

是否使用 DataSourceUtils 获取数据连接就可以高枕无忧了呢？理想很美好，但现实很残酷：如果 DataSourceUtils 在没有事务上下文的方法中使用 getConnection() 获取连接，依然会造成数据连接泄漏！

保持代码清单 6 的代码不变，调整 Spring 配置文件，将清单 3 中 Spring AOP 事务增强配置的代码注释掉，重新运行清单 6 的代码，将得到如下的输出日志：

##### 清单 8\. 输出日志

```
连接数 [active:idle]-[0:0]
连接数 [active:idle]-[1:1]
连接数 [active:idle]-[1:1]
连接数 [active:idle]-[2:1]
连接数 [active:idle]-[2:1]

```

Show moreShow more icon

我们通过下表对数据源连接的占用和泄漏情况进行描述：

##### 表 2\. 执行过程数据源连接占用情况

**时间****执行线程 1****执行线程 2****数据源连接 (active idle leak)**T0未启动未启动0 0 0T1正在执行方法未启动1 1 0T2执行完毕未启动1 1 1T3执行完毕正式执行方法2 1 1T4执行完毕执行完毕2 1 2

仔细对照表 1 的执行过程，我们发现在 T1 时，有事务上下文时的 active 为 2，idle 为 0，而此时由于没有事务管理，则 active 为 1 而 idle 也为 1。这说明有事务上下文时，需要等到整个事务方法（即 logon()）返回后，事务上下文绑定的连接才释放。但在没有事务上下文时，logon() 调用 JdbcTemplate 执行完数据操作后，马上就释放连接。

在 T2 执行线程完成 logon() 方法的执行后，有一个连接没有被释放（active），所以发生了连接泄漏。到 T4 时，两个执行线程都完成了 logon() 方法的调用，但是出现了两个未释放的连接。

要堵上这个连接泄漏的漏洞，需要对 logon() 方法进行如下的改造：

##### 清单 9.JdbcUserService.java：手工释放获取的连接

```
public void logon(String userName) {
    Connection conn = null;
    try {
        conn = DataSourceUtils.getConnection(jdbcTemplate.getDataSource());
        String sql = "UPDATE t_user SET last_logon_time=? WHERE user_name =?";
        jdbcTemplate.update(sql, System.currentTimeMillis(), userName);
        Thread.sleep(1000);
        // ①
    } catch (Exception e) {
        e.printStackTrace();
    }finally {
        // ②显式使用DataSourceUtils释放连接
        DataSourceUtils.releaseConnection(conn,jdbcTemplate.getDataSource());
    }
}

```

Show moreShow more icon

在 ② 处显式调用 `DataSourceUtils.releaseConnection()` 方法释放获取的连接。特别需要指出的是：一定不能在 ① 处释放连接！因为如果 logon() 在获取连接后，① 处代码前这段代码执行时发生异常，则①处释放连接的动作将得不到执行。这将是一个非常具有隐蔽性的连接泄漏的隐患点。

### JdbcTemplate 如何做到对连接泄漏的免疫

分析 JdbcTemplate 的代码，我们可以清楚地看到它开放的每个数据操作方法，首先都使用 DataSourceUtils 获取连接，在方法返回之前使用 DataSourceUtils 释放连接。

来看一下 JdbcTemplate 最核心的一个数据操作方法 execute()：

##### 清单 10.JdbcTemplate\#execute()

```
public <T> T execute(StatementCallback<T> action) throws DataAccessException {
    //① 首先根据DataSourceUtils获取数据连接
    Connection con = DataSourceUtils.getConnection(getDataSource());
    Statement stmt = null;
    try {
        Connection conToUse = con;
       ...
        handleWarnings(stmt);
        return result;
    }
    catch (SQLException ex) {
        JdbcUtils.closeStatement(stmt);
        stmt = null;
        DataSourceUtils.releaseConnection(con, getDataSource());
        con = null;
        throw getExceptionTranslator().translate(
            "StatementCallback", getSql(action), ex);
    }
    finally {
        JdbcUtils.closeStatement(stmt);
        //② 最后根据DataSourceUtils释放数据连接
        DataSourceUtils.releaseConnection(con, getDataSource());
    }
}

```

Show moreShow more icon

在 ① 处通过 DataSourceUtils.getConnection() 获取连接，在 ② 处通过 DataSourceUtils.releaseConnection() 释放连接。所有 JdbcTemplate 开放的数据访问方法最终都是通过 `execute(StatementCallback<T> action)` 执行数据访问操作的，因此这个方法代表了 JdbcTemplate 数据操作的最终实现方式。

正是因为 JdbcTemplate 严谨的获取连接，释放连接的模式化流程保证了 JdbcTemplate 对数据连接泄漏问题的免疫性。所以，如有可能尽量使用 JdbcTemplate，HibernateTemplate 等这些模板进行数据访问操作，避免直接获取数据连接的操作。

### 使用 TransactionAwareDataSourceProxy

如果不得已要显式获取数据连接，除了使用 DataSourceUtils 获取事务上下文绑定的连接外，还可以通过 TransactionAwareDataSourceProxy 对数据源进行代理。数据源对象被代理后就具有了事务上下文感知的能力，通过代理数据源的 getConnection() 方法获取的连接和使用 DataSourceUtils.getConnection() 获取连接的效果是一样的。

下面是使用 TransactionAwareDataSourceProxy 对数据源进行代理的配置：

##### 清单 11.applicationContext.xml：对数据源进行代理

```
<bean id="dataSource"
    class="org.apache.commons.dbcp.BasicDataSource"
    destroy-method="close"
    p:driverClassName="oracle.jdbc.driver.OracleDriver"
    p:url="jdbc:oracle:thin:@localhost:1521:orcl"
    p:username="test"
    p:password="test"
    p:defaultAutoCommit="false"/>

<!-- ①对数据源进行代理-->
<bean id="dataSourceProxy"
    class="org.springframework.jdbc.datasource.TransactionAwareDataSourceProxy"
    p:targetDataSource-ref="dataSource"/>

<!-- ②直接使用数据源的代理对象-->
<bean id="jdbcTemplate"
    class="org.springframework.jdbc.core.JdbcTemplate"
    p:dataSource-ref="dataSourceProxy"/>

<!-- ③直接使用数据源的代理对象-->
<bean id="jdbcManager"
    class="org.springframework.jdbc.datasource.DataSourceTransactionManager"
    p:dataSource-ref="dataSourceProxy"/>

```

Show moreShow more icon

对数据源进行代理后，我们就可以通过数据源代理对象的 getConnection() 获取事务上下文中绑定的数据连接了。

因此，如果数据源已经进行了 TransactionAwareDataSourceProxy 的代理，而且方法存在事务上下文，那么清单 1 的代码也不会生产连接泄漏的问题。

## 其它数据访问技术的等价类

理解了 Spring JDBC 的数据连接泄漏问题，其中的道理可以平滑地推广到其它框架中去。Spring 为每个数据访问技术框架都提供了一个获取事务上下文绑定的数据连接（或其衍生品）的工具类和数据源（或其衍生品）的代理类。

### DataSourceUtils 的等价类

下表列出了不同数据访问技术对应 DataSourceUtils 的等价类：

##### 表 3\. 不同数据访问框架 DataSourceUtils 的等价类

**数据访问技术框架****连接 ( 或衍生品 ) 获取工具类**Spring JDBCorg.springframework.jdbc.datasource.DataSourceUtilsHibernateorg.springframework.orm.hibernate3.SessionFactoryUtilsiBatisorg.springframework.jdbc.datasource.DataSourceUtilsJPAorg.springframework.orm.jpa.EntityManagerFactoryUtilsJDOorg.springframework.orm.jdo.PersistenceManagerFactoryUtils

### TransactionAwareDataSourceProxy 的等价类

下表列出了不同数据访问技术框架下 TransactionAwareDataSourceProxy 的等价类：

##### 表 4\. 不同数据访问框架 TransactionAwareDataSourceProxy 的等价类

**数据访问技术框架****连接 ( 或衍生品 ) 获取工具类**Spring JDBCorg.springframework.jdbc.datasource.TransactionAwareDataSourceProxyHibernateorg.springframework.orm.hibernate3.LocalSessionFactoryBeaniBatisorg.springframework.jdbc.datasource.TransactionAwareDataSourceProxyJPA无JDOorg.springframework.orm.jdo. TransactionAwarePersistenceManagerFactoryProxy

## 结束语

在本文中，我们通过剖析了解到以下的真相：

- 使用 Spring JDBC 时如果直接获取 Connection，可能会造成连接泄漏。为降低连接泄漏的可能，尽量使用 DataSourceUtils 获取数据连接。也可以对数据源进行代理，以便将其拥有事务上下文的感知能力；
- 可以将 Spring JDBC 防止连接泄漏的解决方案平滑应用到其它的数据访问技术框架中。