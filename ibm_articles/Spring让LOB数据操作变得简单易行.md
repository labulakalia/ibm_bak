# Spring 让 LOB 数据操作变得简单易行
在 Spring 中处理 LOB 数据的原理和方法

**标签:** Java,Spring,数据库

[原文链接](https://developer.ibm.com/zh/articles/j-lo-spring-lob/)

陈 雄华

发布: 2007-07-30

* * *

## 概述

LOB 代表大对象数据，包括 BLOB 和 CLOB 两种类型，前者用于存储大块的二进制数据，如图片数据，视频数据等，而后者用于存储长文本数据，如论坛的帖子内容，产品的详细描述等。值得注意的是：在不同的数据库中，大对象对应的字段类型是不尽相同的，如 DB2 对应 BLOB/CLOB，MySql 对应 BLOB/LONGTEXT，SqlServer 对应 IMAGE/TEXT。需要指出的是，有些数据库的大对象类型可以象简单类型一样访问，如 MySql 的 LONGTEXT 的操作方式和 VARCHAR 类型一样。在一般情况下， LOB 类型数据的访问方式不同于其它简单类型的数据，我们经常会以流的方式操作 LOB 类型的数据。此外，LOB 类型数据的访问不是线程安全的，需要为其单独分配相应的数据库资源，并在操作完成后释放资源。最后，Oracle 9i 非常有个性地采用非 JDBC 标准的 API 操作 LOB 数据。所有这些情况给编写操作 LOB 类型数据的程序带来挑战，Spring 在 `org.springframework.jdbc.support.lob` 包中为我们提供了相应的帮助类，以便我们轻松应对这头拦路虎。

Spring 大大降低了我们处理 LOB 数据的难度。首先，Spring 提供了 `NativeJdbcExtractor` 接口，您可以在不同环境里选择相应的实现类从数据源中获取本地 JDBC 对象；其次，Spring 通过 `LobCreator` 接口取消了不同数据厂商操作 LOB 数据的差别，并提供了创建 LobCreator 的 `LobHandler` 接口，您只要根据底层数据库类型选择合适的 LobHandler 进行配置即可。

本文将详细地讲述通过 Spring JDBC 插入和访问 LOB 数据的具体过程。不管是以块的方式还是以流的方式，您都可以通过 LobCreator 和 LobHandler 方便地访问 LOB 数据。对于 ORM 框架来说，JPA 拥有自身处理 LOB 数据的配置类型，Spring 为 Hibernate 和 iBatis 分别提供了 LOB 数据类型的配置类，您仅需要使用这些类进行简单的配置就可以像普通类型一样操作 LOB 类型数据。

## 本地 JDBC 对象

当您在 Web 应用服务器或 Spring 中配置数据源时，从数据源中返回的数据连接对象是本地 JDBC 对象（如 DB2Connection、OracleConnection）的代理类，这是因为数据源需要改变数据连接一些原有的行为以便对其进行控制：如调用 `Connection#close()` 方法时，将数据连接返回到连接池中而非将其真的关闭。

在访问 LOB 数据时，根据数据库厂商的不同，可能需要使用被代理前的本地 JDBC 对象（如 DB2Connection 或 DB2ResultSet）特有的 API。为了从数据源中获取本地 JDBC 对象， Spring 定义了 `org.springframework.jdbc.support.nativejdbc.NativeJdbcExtractor` 接口并提供了相应的实现类。 `NativeJdbcExtractor` 定义了从数据源中抽取本地 JDBC 对象的若干方法：

方法说明`Connection getNativeConnection(Connection con)`获取本地 Connection 对象`Connection getNativeConnectionFromStatement(Statement stmt)`获取本地 Statement 对象`PreparedStatement getNativePreparedStatement(PreparedStatement ps)`获取本地 PreparedStatement 对象`ResultSet getNativeResultSet(ResultSet rs)`获取本地 ResultSet 对象`CallableStatement getNativeCallableStatement(CallableStatement cs)`获取本地 CallableStatement 对象



有些简单的数据源仅对 `Connection` 对象进行代理，这时可以直接使用 `SimpleNativeJdbcExtractor` 实现类。但有些数据源（如 Jakarta Commons DBCP）会对所有的 JDBC 对象进行代理，这时，就需要根据具体的情况选择适合的抽取器实现类了。下表列出了不同数据源本地 JDBC 对象抽取器的实现类：

数据源类型说明WebSphere 4 及以上版本的数据源`org.springframework.jdbc.support.nativejdbc.WebSphereNativeJdbcExtractor`WebLogic 6.1+ 及以上版本的数据源`org.springframework.jdbc.support.nativejdbc.WebLogicNativeJdbcExtractor`JBoss 3.2.4 及以上版本的数据源`org.springframework.jdbc.support.nativejdbc.JBossNativeJdbcExtractor`C3P0 数据源`org.springframework.jdbc.support.nativejdbc.C3P0NativeJdbcExtractor`DBCP 数据源`org.springframework.jdbc.support.nativejdbc.CommonsDbcpNativeJdbcExtractor`ObjectWeb 的 XAPool 数据源`org.springframework.jdbc.support.nativejdbc.XAPoolNativeJdbcExtractor`

下面的代码演示了从 DBCP 数据源中获取 DB2 的本地数据库连接 DB2Connection 的方法：

##### 清单 1\. 获取本地数据库连接

```
package com.baobaotao.dao.jdbc;

import java.sql.Connection;
import COM.ibm.db2.jdbc.net.DB2Connection;
import org.springframework.jdbc.core.support.JdbcDaoSupport;
import org.springframework.jdbc.datasource.DataSourceUtils;

public class PostJdbcDao extends JdbcDaoSupport implements PostDao {
public void getNativeConn(){
    try {
      Connection conn = DataSourceUtils.getConnection(getJdbcTemplate()
          .getDataSource()); ① 使用 DataSourceUtils 从模板类中获取连接
      ② 使用模板类的本地 JDBC 抽取器获取本地的 Connection
      conn = getJdbcTemplate().getNativeJdbcExtractor().getNativeConnection(conn);
      DB2Connection db2conn = (DB2Connection) conn; ③ 这时可以强制进行类型转换了
     ...
    } catch (Exception e) {
      e.printStackTrace();
    }
}
}

```

Show moreShow more icon

在 ① 处我们通过 `DataSourceUtils` 获取当前线程绑定的数据连接，为了使用线程上下文相关的事务，通过 `DataSourceUtils` 从数据源中获取连接是正确的做法，如果直接通过 `dateSource` 获取连接，则将得到一个和当前线程上下文无关的数据连接实例。

JdbcTemplate 可以在配置时注入一个本地 JDBC 对象抽取器，要使代码 清单 1 正确运行，我们必须进行如下配置：

##### 清单 2\. 为 JdbcTemplate 装配本地 JDBC 对象抽取器

```
...
<bean id="dataSource" class="org.apache.commons.dbcp.BasicDataSource"
    destroy-method="close">
    <property name="driverClassName"
      value="${jdbc.driverClassName}" />
    <property name="url" value="${jdbc.url}" />
    <property name="username" value="${jdbc.username}" />
    <property name="password" value="${jdbc.password}" />
</bean>
① 定义 DBCP 数据源的 JDBC 本地对象抽取器
<bean id="nativeJdbcExtractor"
class="org.springframework.jdbc.support.nativejdbc.CommonsDbcpNativeJdbcExtractor"
lazy-init="true" />
<bean id="jdbcTemplate" class="org.springframework.jdbc.core.JdbcTemplate">
    <property name="dataSource" ref="dataSource" />
    ② 设置抽取器
    <property name="nativeJdbcExtractor" ref="nativeJdbcExtractor"/>
</bean>
<bean id="postDao" class="com.baobaotao.dao.jdbc.PostJdbcDao">
    <property name="jdbcTemplate" ref="jdbcTemplate" />
</bean>

```

Show moreShow more icon

在获取 DB2 的本地 Connection 实例后，我们就可以使用该对象的一些特有功能了，如使用 DB2Connection 的特殊 API 对 LOB 对象进行操作。

### LobCreator

虽然 JDBC 定义了两个操作 LOB 类型的接口： `java.sql.Blob` 和 `java.sql.Clob` ，但有些厂商的 JDBC 驱动程序并不支持这两个接口。为此，Spring 定义了一个独立于 `java.sql.Blob/Clob` 的 `LobCreator` 接口，以统一的方式操作各种数据库的 LOB 类型数据。因为 LobCreator 本身持有 LOB 所对应的数据库资源，所以它不是线程安全的，一个 LobCreator 只能操作一个 LOB 数据。

为了方便在 PreparedStatement 中使用 LobCreator，您可以直接使用 `JdbcTemplate#execute(String sql,AbstractLobCreatingPreparedStatementCallback lcpsc)` 方法。下面对 LobCreator 接口中的方法进行简要说明：

方法说明`void close()`关闭会话，并释放 LOB 资源`void setBlobAsBinaryStream(PreparedStatement ps, int paramIndex, InputStream contentStream, int contentLength)`通过流填充 BLOB 数据`void setBlobAsBytes(PreparedStatement ps, int paramIndex, byte[] content)`通过二进制数据填充 BLOB 数据`void setClobAsAsciiStream(PreparedStatement ps, int paramIndex, InputStream asciiStream, int contentLength)`通过 Ascii 字符流填充 CLOB 数据`void setClobAsCharacterStream(PreparedStatement ps, int paramIndex, Reader characterStream, int contentLength)`通过 Unicode 字符流填充 CLOB 数据`void setClobAsString(PreparedStatement ps, int paramIndex, String content)`通过字符串填充 CLOB 数据

### LobHandler

`LobHandler` 接口为操作 BLOB/CLOB 提供了统一访问接口，而不管底层数据库究竟是以大对象的方式还是以一般数据类型的方式进行操作。此外，LobHandler 还充当了 LobCreator 的工厂类。

大部分数据库厂商的 JDBC 驱动程序（如 DB2）都以 JDBC 标准的 API 操作 LOB 数据，但 Oracle 9i 及以前的 JDBC 驱动程序采用了自己的 API 操作 LOB 数据，Oracle 9i 直接使用自己的 API 操作 LOB 数据，且不允许通过 PreparedStatement 的 `setAsciiStream()` 、 `setBinaryStream()` 、 `setCharacterStream()` 等方法填充流数据。Spring 提供 `LobHandler` 接口主要是为了迁就 Oracle 特立独行的作风。所以 Oracle 必须使用 `OracleLobHandler` 实现类，而其它的数据库统一使用 `DefaultLobHandler` 就可以了。Oracle 10g 改正了 Oracle 9i 这个异化的风格，终于天下归一了，所以 Oracle 10g 也可以使用 `DefaultLobHandler` 。 下面，我们来看一下 `LobHandler` 接口的几个重要方法：

方法说明`InputStream getBlobAsBinaryStream(ResultSet rs, int columnIndex)`从结果集中返回 InputStream，通过 InputStream 读取 BLOB 数据`byte[] getBlobAsBytes(ResultSet rs, int columnIndex)`以二进制数据的方式获取结果集中的 BLOB 数据；`InputStream getClobAsAsciiStream(ResultSet rs, int columnIndex)`从结果集中返回 InputStream，通过 InputStreamn 以 Ascii 字符流方式读取 BLOB 数据`Reader getClobAsCharacterStream(ResultSet rs, int columnIndex)`从结果集中获取 Unicode 字符流 Reader，并通过 Reader以Unicode 字符流方式读取 CLOB 数据`String getClobAsString(ResultSet rs, int columnIndex)`从结果集中以字符串的方式获取 CLOB 数据`LobCreator getLobCreator()`生成一个会话相关的 LobCreator 对象

## 在 Spring JDBC 中操作 LOB 数据

### 插入 LOB 数据

假设我们有一个用于保存论坛帖子的 `t_post` 表，拥有两个 LOB 字段，其中 `post_text` 是 CLOB 类型，而 `post_attach` 是 BLOB 类型。下面，我们来编写插入一个帖子记录的代码：

##### 清单 3\. 添加 LOB 字段数据

```
package com.baobaotao.dao.jdbc;
...
import java.sql.PreparedStatement;
import java.sql.SQLException;
import org.springframework.jdbc.core.support.AbstractLobCreatingPreparedStatementCallback;
import org.springframework.jdbc.support.lob.LobCreator;
import org.springframework.jdbc.support.lob.LobHandler;

public class PostJdbcDao extends JdbcDaoSupport implements PostDao {
private LobHandler lobHandler; ① 定义 LobHandler 属性
public LobHandler getLobHandler() {
    return lobHandler;
}
public void setLobHandler(LobHandler lobHandler) {
    this.lobHandler = lobHandler;
}
public void addPost(final Post post) {
    String sql = " INSERT INTO t_post(post_id,user_id,post_text,post_attach)"
        + " VALUES(?,?,?,?)";
    getJdbcTemplate().execute(sql,
      new AbstractLobCreatingPreparedStatementCallback(this.lobHandler) { ②
          protected void setValues(PreparedStatement ps,LobCreator lobCreator)
                      throws SQLException {
            ps.setInt(1, 1);
            ps.setInt(2, post.getUserId());
            ③ 设置 CLOB 字段
            lobCreator.setClobAsString(ps, 3, post.getPostText());
            ④ 设置 BLOB 字段
            lobCreator.setBlobAsBytes(ps, 4, post.getPostAttach());
          }
        });
}
...
}

```

Show moreShow more icon

首先，我们在 `PostJdbcDao` 中引入了一个 `LobHandler` 属性，如 ① 所示，并通过 `JdbcTemplate#execute(String sql,AbstractLobCreatingPreparedStatementCallback lcpsc)` 方法完成插入 LOB 数据的操作。我们通过匿名内部类的方式定义 `LobCreatingPreparedStatementCallback` 抽象类的子类，其构造函数需要一个 `LobHandler` 入参，如 ② 所示。在匿名类中实现了父类的抽象方法 `setValues(PreparedStatement ps,LobCreator lobCreator)` ，在该方法中通过 `lobCreator` 操作 LOB 对象，如 ③、④ 所示，我们分别通过字符串和二进制数组填充 BLOB 和 CLOB 的数据。您同样可以使用流的方式填充 LOB 数据，仅需要调用 lobCreator 相应的流填充方法即可。

我们需要调整 Spring 的配置文件以配合我们刚刚定义的 PostJdbcDao。假设底层数据库是 Oracle，可以采用以下的配置方式：

##### 清单 4\. Oracle 数据库的 LobHandler 配置

```
...
<bean id="nativeJdbcExtractor"
class="org.springframework.jdbc.support.nativejdbc.CommonsDbcpNativeJdbcExtractor"
lazy-init="true"/>
<bean id="oracleLobHandler"
class="org.springframework.jdbc.support.lob.OracleLobHandler"
lazy-init="true">
<property name="nativeJdbcExtractor" ref="nativeJdbcExtractor"/> ① 设置本地 Jdbc 对象抽取器
</bean>
<bean id="postDao" class="com.baobaotao.dao.jdbc.PostJdbcDao">
<property name="lobHandler" ref="oracleLobHandler"/> ② 设置 LOB 处理器
</bean>

```

Show moreShow more icon

大家可能已经注意到 `nativeJdbcExtractor` 和 `oracleLobHandler` Bean 都设置为 `lazy-init="true"` ，这是因为 `nativeJdbcExtractor` 需要通过运行期的反射机制获取底层的 JDBC 对象，所以需要避免在 Spring 容器启动时就实例化这两个 Bean。

LobHandler 需要访问本地 JDBC 对象，这一任务委托给 `NativeJdbcExtractor` Bean 来完成，因此我们在 ① 处为 LobHandler 注入了一个 `nativeJdbcExtractor` 。最后，我们把 `lobHandler` Bean 注入到需要进行 LOB 数据访问操作的 PostJdbcDao 中，如 ② 所示。

如果底层数据库是 DB2、SQL Server、MySQL 等非 Oracle 的其它数据库，则只要简单配置一个 `DefaultLobHandler` 就可以了，如下所示：

##### 清单 5\. 一般数据库 LobHandler 的配置

```
<bean id="defaultLobHandler"
class="org.springframework.jdbc.support.lob.DefaultLobHandler"
lazy-init="true"/>
<bean id="postDao" class="com.baobaotao.dao.jdbc.PostJdbcDao">
<property name="lobHandler" ref=" defaultLobHandler"/>
<property name="jdbcTemplate" ref="jdbcTemplate" />
</bean>

```

Show moreShow more icon

DefaultLobHandler 只是简单地代理标准 JDBC 的 PreparedStatement 和 ResultSet 对象，由于并不需要访问数据库驱动本地的 JDBC 对象，所以它不需要 NativeJdbcExtractor 的帮助。您可以通过以下的代码测试 PostJdbcDao 的 `addPost()` 方法：

##### 清单 6\. 测试 PostJdbcDao 的 addPost() 方法

```
package com.baobaotao.dao.jdbc;

import org.springframework.core.io.ClassPathResource;
import org.springframework.test.AbstractDependencyInjectionSpringContextTests;
import org.springframework.util.FileCopyUtils;
import com.baobaotao.dao.PostDao;
import com.baobaotao.domain.Post;
public class TestPostJdbcDaoextends AbstractDependencyInjectionSpringContextTests {
private PostDao postDao;
public void setPostDao(PostDao postDao) {
    this.postDao = postDao;
}
protected String[] getConfigLocations() {
    return new String[]{"classpath:applicationContext.xml"};
}
public void testAddPost() throws Throwable{
    Post post = new Post();
    post.setPostId(1);
    post.setUserId(2);
    ClassPathResource res = new ClassPathResource("temp.jpg"); ① 获取图片资源
    byte[] mockImg = FileCopyUtils.copyToByteArray(res.getFile());  ② 读取图片文件的数据
    post.setPostAttach(mockImg);
    post.setPostText("测试帖子的内容");
    postDao.addPost(post);
}
}

```

Show moreShow more icon

这里，有几个知识点需要稍微解释一下： `AbstractDependencyInjectionSpringContextTests` 是 Spring 专门为测试提供的类，它能够直接从 IoC 容器中装载 Bean。此外，我们使用了 `ClassPathResource` 加载图片资源，并通过 `FileCopyUtils` 读取文件的数据。 `ClassPathResource` 和 `FileCopyUtils` 都是 Spring 提供的非常实用的工具类。

### 以块数据方式读取 LOB 数据

您可以直接用数据块的方式读取 LOB 数据：用 `String` 读取 CLOB 字段的数据，用 `byte[]` 读取 BLOB 字段的数据。在 PostJdbcDao 中添加一个 `getAttachs()` 方法，以便获取某一用户的所有带附件的帖子：

##### 清单 7\. 以块数据访问 LOB 数据

```
public List getAttachs(final int userId){
String sql = "SELECT post_id,post_attach FROM t_post "+
"where user_id =? and post_attach is not null ";
return getJdbcTemplate().query(
        sql,new Object[] {userId},
        new RowMapper() {
          public Object mapRow(ResultSet rs, int rowNum) throws SQLException {
            int postId = rs.getInt(1);
            ① 以二进制数组方式获取 BLOB 数据。
            byte[] attach = lobHandler.getBlobAsBytes(rs, 2);
            Post post = new Post();
            post.setPostId(postId);
            post.setPostAttach(attach);
            return post;
          }
        });
}

```

Show moreShow more icon

通过 JdbcTemplate 的 `List query(String sql, Object[] args, RowMapper rowMapper)` 接口处理行数据的映射。在 RowMapper 回调的 `mapRow()` 接口方法中，通过 LobHandler 以 `byte[]` 获取 BLOB 字段的数据。

### 以流数据方式读取 LOB 数据

由于 LOB 数据可能很大（如 100M），如果直接以块的方式操作 LOB 数据，需要消耗大量的内存资源，对应用程序整体性能产生巨大的冲击。对于体积很大的 LOB 数据，我们可以使用流的方式进行访问，减少内存的占用。JdbcTemplate 为此提供了一个 `Object query(String sql, Object[] args, ResultSetExtractor rse)` 方法， `ResultSetExtractor` 接口拥有一个处理流数据的抽象类 `org.springframework.jdbc.core.support.AbstractLobStreamingResultSetExtractor` ，可以通过扩展此类用流的方式操作 LOB 字段的数据。下面我们为 PostJdbcDao 添加一个以流的方式获取某个帖子附件的方法：

##### 清单 8\. 以流方式访问 LOB 数据

```
...
public void getAttach(final int postId,final OutputStream os){ ① 用于接收 LOB 数据的输出流
String sql = "SELECT post_attach FROM t_post WHERE post_id=? ";
getJdbcTemplate().query(
      sql, new Object[] {postId},
      new AbstractLobStreamingResultSetExtractor() { ② 匿名内部类
③ 处理未找到数据行的情况
protected void handleNoRowFound() throws LobRetrievalFailureException {
          System.out.println("Not Found result!");
        }
        ④ 以流的方式处理 LOB 字段
        public void streamData(ResultSet rs) throws SQLException, IOException {
          InputStream is = lobHandler.getBlobAsBinaryStream(rs, 1);
          if (is != null) {
            FileCopyUtils.copy(is, os);
          }
        }
      }
);
}

```

Show moreShow more icon

通过扩展 `AbstractLobStreamingResultSetExtractor` 抽象类，在 `streamData(ResultSet rs)` 方法中以流的方式读取 LOB 字段数据，如 ④ 所示。这里我们又利用到了 Spring 的工具类 `FileCopyUtils` 将输入流的数据拷贝到输出流中。在 `getAttach()` 方法中通过入参 `OutputStream os` 接收 LOB 的数据，如 ① 所示。您可以同时覆盖抽象类中的 `handleNoRowFound()` 方法，定义未找到数据行时的处理逻辑。

## 在 JPA 中操作 LOB 数据

在 JPA 中 LOB 类型的持久化更加简单，仅需要通过特殊的 LOB 注释（Annotation）就可以达到目的。我们对 Post 中的 LOB 属性类型进行注释：

##### 清单 9\. 注释 LOB 类型属性

```
package com.baobaotao.domain;
...
import javax.persistence.Basic;
import javax.persistence.Lob;
import javax.persistence. Column;
@Entity(name = "T_POST")
public class Post implements Serializable {
...
@Lob ①-1 表示该属性是 LOB 类型的字段
@Basic(fetch = FetchType.EAGER) ①-2 不采用延迟加载机制
@Column(name = "POST_TEXT", columnDefinition = "LONGTEXT NOT NULL") ①-3 对应字段类型
private String postText;

@Lob ②-1 表示该属性是 LOB 类型的字段
@Basic(fetch = FetchType. LAZY) ②-2 采用延迟加载机制
@Column(name = "POST_ATTACH", columnDefinition = "BLOB") ②-3 对应字段类型
private byte[] postAttach;
...
｝

```

Show moreShow more icon

`postText` 属性对应 `T_POST` 表的 `POST_TEXT` 字段，该字段的类型是 `LONTTEXT` ，并且非空。JPA 通过 `@Lob` 将属性标注为 LOB 类型，如 ①-1 和 ②-1 所示。通过 `@Basic` 指定 LOB 类型数据的获取策略， `FetchType.EAGER` 表示非延迟加载，而 `FetchType.LAZY` 表示延迟加载，如 ①-2 和 ②-2 所示。通过 `@Column` 的 `columnDefinition` 属性指定数据表对应的 LOB 字段类型，如 ①-3 和 ②-3 所示。

关于 JPA 注释的更多信息，请阅读 参考资源 中的相关技术文章。

## 在 Hibernate 中操作 LOB 数据

##### 提示

使用 Spring JDBC 时，我们除了可以按 byte[]、String 类型处理 LOB 数据外，还可以使用流的方式操作 LOB 数据，当 LOB 数据体积较大时，流操作是唯一可行的方式。可惜，Spring 并未提供以流方式操作 LOB 数据的 UserType（记得 Spring 开发组成员认为在实现上存在难度）。不过，`www.atlassian.com` 替 Spring 完成了这件难事，读者可以通过 [这里](https://docs.atlassian.com/atlassian-confluence/) 了解到这个满足要求的 `BlobInputStream` 类型。

Hibernate 为处理特殊数据类型字段定义了一个接口： `org.hibernate.usertype.UserType` 。Spring 在 `org.springframework.orm.hibernate3.support` 包中为 BLOB 和 CLOB 类型提供了几个 `UserType` 的实现类。因此，我们可以在 Hibernate 的映射文件中直接使用这两个实现类轻松处理 LOB 类型的数据。

- `BlobByteArrayType` ：将 BLOB 数据映射为 `byte[]` 类型的属性；
- `BlobStringType` ：将 BLOB 数据映射为 `String` 类型的属性；
- `BlobSerializableType` ：将 BLOB 数据映射为 `Serializable` 类型的属性；
- `ClobStringType` ：将 CLOB 数据映射为 `String` 类型的属性；

下面我们使用 Spring 的 `UserType` 为 Post 配置 Hibernate 的映射文件，如 清单 10 所示：

##### 清单 10 . LOB 数据映射配置

```
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE hibernate-mapping PUBLIC "-//Hibernate/Hibernate Mapping DTD//EN"
    "http://hibernate.sourceforge.net/hibernate-mapping-3.0.dtd">
<hibernate-mapping auto-import="true" default-lazy="false">
<class name="com.baobaotao.domain.Post" table="t_post">
    <id name="postId" column="post_id">
      <generator class="identity" />
    </id>
    <property name="userId" column="user_id"/>
    <property name="postText" column="post_text"
    type="org.springframework.orm.hibernate3.support.ClobStringType"/>①对应 CLOB 字段
    <property name="postAttach" column="post_attach"
     type="org.springframework.orm.hibernate3.support.BlobByteArrayType"/>② BLOB 字段
    <property name="postTime" column="post_time" type="date" />
    <many-to-one name="topic" column="topic_id" class="com.baobaotao.domain.Topic" />
</class>
</hibernate-mapping>

```

Show moreShow more icon

`postText` 为 `String` 类型的属性，对应数据库的 CLOB 类型，而 `postAttach` 为 `byte[]` 类型的属性，对应数据库的 BLOB 类型。分别使用 Spring 所提供的相应 `UserType` 实现类进行配置，如 ① 和 ② 处所示。

在配置好映射文件后，还需要在 Spring 配置文件中定义 LOB 数据处理器，让 SessionFactory 拥有处理 LOB 数据的能力：

##### 清单 11 . 将 LobHandler 注入到 SessionFactory 中

```
...
<bean id="lobHandler" class="org.springframework.jdbc.support.lob.DefaultLobHandler"
lazy-init="true" />
<bean id="sessionFactory"
class="org.springframework.orm.hibernate3.LocalSessionFactoryBean">
<property name="dataSource" ref="dataSource" />
<property name="lobHandler" ref="lobHandler" /> ① 设置 LOB 处理器
...
</bean>

```

Show moreShow more icon

在一般的数据库（如 DB2）中，仅需要简单地使用 `HibernateTemplate#save(Object entity)` 等方法就可以正确的保存 LOB 数据了。如果是 Oracle 9i 数据库，还需要配置一个本地 JDBC 抽取器，并使用特定的 LobHandler 实现类，如 清单 4 所示。

使用 LobHandler 操作 LOB 数据时，需要在事务环境下才能工作，所以必须事先配置事务管理器，否则会抛出异常。

## 在 iBatis 中操作 LOB 数据

iBatis 为处理不同类型的数据定义了一个统一的接口： `com.ibatis.sqlmap.engine.type.TypeHandler` 。这个接口类似于 Hibernate 的 UserType。iBatis 本身拥有该接口的众多实现类，如 LongTypeHandler、DateTypeHandler 等，但没有为 LOB 类型提供对应的实现类。Spring 在 `org.springframework.orm.ibatis.support` 包中为我们提供了几个处理 LOB 类型的 TypeHandler 实现类：

- `BlobByteArrayTypeHandler` ：将 BLOB 数据映射为 `byte[]` 类型；
- `BlobSerializableTypeHandler` ：将 BLOB 数据映射为 `Serializable` 类型的对象；
- `ClobStringTypeHandler` ：将 CLOB 数据映射为 `String` 类型；

当结果集中包括 LOB 数据时，需要在结果集映射配置项中指定对应的 Handler 类，下面我们采用 Spring 所提供的实现类对 Post 结果集的映射进行配置。

##### 清单 12 . 对 LOB 数据进行映射

```
<?xml version="1.0" encoding="UTF-8" ?>
<!DOCTYPE sqlMap PUBLIC "-//ibatis.apache.org//DTD SQL Map 2.0//EN"
    "http://ibatis.apache.org/dtd/sql-map-2.dtd">
<sqlMap namespace="Post">
<typeAlias alias="post" type="com.baobaotao.domain.Post"/>
<resultMap id="result" class="post">
        <result property="postId" column="post_id"/>
        <result property="userId" column="user_id"/>
        <result property="postText" column="post_text" ① 读取 CLOB 类型数据
        typeHandler="org.springframework.orm.ibatis.support.ClobStringTypeHandler"/>
        <result property="postAttach" column="post_attach" ② 读取 BLOB 类型数据
        typeHandler="org.springframework.orm.ibatis.support.BlobByteArrayTypeHandler"/>
    </resultMap>
<select id="getPost" resultMap="result">
    SELECT post_id,user_id,post_text,post_attach,post_time
    FROM t_post  WHERE post_id =#postId#
</select>
<insert id="addPost">
    INSERT INTO t_post(user_id,post_text,post_attach,post_time)
    VALUES(#userId#,
#postText,handler=org.springframework.orm.ibatis.support.ClobStringTypeHandler#, ③
#postAttach,handler=org.springframework.orm.ibatis.support.BlobByteArrayTypeHandler#, ④
#postTime#)
</insert>
</sqlMap>

```

Show moreShow more icon

##### 提示

为每一个 LOB 类型字段分别指定处理器并不是一个好主意，iBatis 允许在 sql-map-config.xml 配置文件中通过  标签统一定义特殊类型数据的处理器，如：

` `

当 iBatis 引擎从结果集中读取或更改 LOB 类型数据时，都需要指定处理器。我们在 ① 和 ② 处为读取 LOB 类型的数据指定处理器，相似的，在 ③ 和 ④ 处为插入 LOB 类型的数据也指定处理器。

此外，我们还必须为 SqlClientMap 提供一个 LobHandler：

##### 清单 13\. 将 LobHandler 注入到 SqlClientMap 中

```
<bean id="lobHandler" class="org.springframework.jdbc.support.lob.DefaultLobHandler"
    lazy-init="true" />
<bean id="sqlMapClient" class="org.springframework.orm.ibatis.SqlMapClientFactoryBean">
    <property name="dataSource" ref="dataSource" />
    <property name="lobHandler" ref="lobHandler" /> ①设置LobHandler
    <property name="configLocation"
        value="classpath:com/baobaotao/dao/ibatis/sql-map-config.xml" />
</bean>

```

Show moreShow more icon

处理 LOB 数据时，Spring 要求在事务环境下工作，所以还必须配置一个事务管理器。iBatis 的事务管理器和 Spring JDBC 事务管理器相同，此处不再赘述。

## 小结

本文就 Spring 中如何操作 LOB 数据进行较为全面的讲解，您仅需简单地配置 LobHandler 就可以直接在程序中象一般数据一样操作 LOB 数据了。对于 ORM 框架来说，Spring 为它们分别提供了支持类，您仅要使用相应的支持类进行配置就可以了。因此您会发现在传统 JDBC 程序操作 LOB 头疼的问题将变得轻松了许多。