# 开发 Spring Redis 应用程序
使用 Redis 作为数据存储来构建基于 Spring 的应用程序

**标签:** Java,Spring,数据库

[原文链接](https://developer.ibm.com/zh/articles/os-springredis/)

Shekar Gulati

发布: 2015-02-12

* * *

开源的 Spring 框架是企业应用程序开发的一根中流砥柱，它的用户群中包含数百万的 Java 开发人员。Spring Data 是保护性开源项目，用于简化受 Spring 支持的、使用了数据访问技术的应用程序的构建，这些数据访问技术包括非关系数据库、MapReduce 框架和基于云的数据服务等现代技术。其中一项技术是 Redis（ **远程 字典 服务器** ），它是一个开源的、高级的、NoSQL 键-值数据存储，是使用 ANSI C 编写的。本文将介绍 Redis、它的数据模型和数据类型，以及它的优点。然后将展示如何使用 Spring Data Redis 构建一个样例应用程序。

## Redis 简介

Redis 是一种内存型数据存储，也可以将它写入磁盘中来实现耐久性。Redis 可通过两种方式来持久存储数据：RDB 和 AOF。RDB 持久性按照指定的间隔对您的数据集执行时间点快照。它不是非常耐久，而且您可能会丢失一些数据，但它非常快。AOF 的持久性要长得多，而且记录了服务器收到的每个写入操作。这些写入操作在服务器启动时重新执行，重新构建原始数据集。在查询 Redis 时，将从内存中获取数据，绝不会从磁盘获取数据，Redis 对内存中存储的键和值执行所有操作。

Redis 采用了一种客户端/服务器模型，借用该模型来监听 TCP 端口并接受命令。Redis 中的所有命令都是原子性的，所以您可以从不同的客户端处理同一个键，没有任何争用条件。如果您使用的是 memcached（一个内存型对象缓存系统），您会发现自己对它很熟悉，但 Redis（可以说）是 memcached++。（请参阅 Andrew Glover 的 “ [Java 开发 2.0：现实世界中的 Redis](http://www.ibm.com/developerworks/cn/java/j-javadev2-22/) ”，了解 Redis 与 memcached 之间的对比。）Redis 也支持数据复制。

### 数据模型

Redis 数据模型不仅与关系数据库管理系统 (RDBMS) 不同，也不同于任何简单的 NoSQL 键-值数据存储。Redis 数据类型类似于编程语言的基础数据类型，所以开发人员感觉很自然。每个数据类型都支持适用于其类型的操作。受支持的数据类型包括：

- 字符串
- 列表
- 集合
- 有序集
- 哈希值

### 关键优势

Redis 的优势包括它的速度、它对富数据类型的支持、它的操作的原子性，以及它的通用性：

- Redis 非常快。它每秒可执行约 100,000 个 `SET` 以及约 100,000 个 `GET` 操作。您可以使用 `redis-benchmark` 实用程序在自己的机器上对它的性能进行基准测试。（ `redis-benchmark` 模拟在它发送总共 _M_ 个查询的同时， _N_ 个客户端完成的 `SET` / `GET` 操作。）
- Redis 对大多数开发人员已知道的大多数数据类型提供了原生支持，这使得各种问题得以轻松解决。经验会告诉您哪个问题最好由何种数据类型来处理。
- 因为所有 Redis 操作都是原子性的，所以多个客户端会并发地访问一个 Redis 服务器，获取相同的更新值。
- Redis 是一个多效用工具，对许多用例很有用，这些用例包括缓存、消息队列（Redis 原生支持发布/订阅）、短期应用程序数据（比如 Web 会话、Web 页面命中计数）等。

## 开始使用 Redis

要开始在 Linux® 或 UNIX® 上使用 Redis，可以下载压缩的 .tar 文件（参见 参考资料 ），解压它，然后运行 `make` 命令：

```
wget http://redis.googlecode.com/files/redis-2.6.7.tar.gz
tar xzf redis-2.6.7.tar.gz
cd redis-2.6.7
make

```

Show moreShow more icon

编译的二进制文件现在包含在 src 目录中。使用以下命令运行 Redis：

```
src/redis-server

```

Show moreShow more icon

Redis 使用以下输出进行响应：

````
[988] 05 Jan 14:41:00.230 # Warning: no config file specified, using the default config.
In order to specify a config file use src/redis-server /path/to/redis.conf

[988] 05 Jan 14:41:00.231 * Max number of open files set to 10032
                _._
           _.-``__ ''-._
      _.-``    `.  `_.  ''-._           Redis 2.6.7 (00000000/0) 64 bit
.-``.-```.  ```\/    _.,_ ''-._
(    '      ,       .-`  | `,    )     Running in stand alone mode
|`-._`-...-` __...-.``-._|'` _.-'|     Port: 6379
|    `-._   `._    /     _.-'    |     PID: 988
`-._    `-._ `-./  _.-'    _.-'
|`-._`-._    `-.__.-'    _.-'_.-'|
|    `-._`-._        _.-'_.-'    |           http://redis.io
`-._    `-._`-.__.-'_.-'    _.-'
|`-._`-._    `-.__.-'    _.-'_.-'|
|    `-._`-._        _.-'_.-'    |
`-._    `-._`-.__.-'_.-'    _.-'
      `-._    `-.__.-'    _.-'
          `-._        _.-'
              `-.__.-'
[988] 05 Jan 14:41:00.238 # Server started, Redis version 2.6.7
[988] 05 Jan 14:41:00.239 * DB loaded from disk: 0.000 seconds
[988] 05 Jan 14:41:00.239 * The server is now ready to accept connections on port 6379

````

Show moreShow more icon

要使用内置的客户端与 Redis 交互，可从命令行启动该客户端：

```
src/redis-cli

```

Show moreShow more icon

客户端会话显示了 Redis 对 `ping` 和 `INFO` 命令的响应：

```
redis 127.0.0.1:6379> ping
PONG
redis 127.0.0.1:6379>
redis 127.0.0.1:6379>
redis 127.0.0.1:6379> INFO
# Server
redis_version:2.6.7
redis_git_sha1:00000000
redis_git_dirty:0
redis_mode:standalone
os:Darwin 12.0.0 x86_64
arch_bits:64
multiplexing_api:kqueue
gcc_version:4.2.1
process_id:3449
run_id:270454ebad19fbc851194548569efca6ac63e00a
tcp_port:6379
uptime_in_seconds:95
uptime_in_days:0
lru_clock:1407736
....

```

Show moreShow more icon

## Redis 数据类型示例

现在我将简要介绍一下受 Redis 支持的数据类型，展示一些使用内置客户端的简单示例。

### 字符串

字符串是 Redis 支持的一种基本数据类型。您可以使用 `SET` 命令设置一个键的字符串值，也可以使用 `GET` 命令获取一个键的字符串值：

```
redis> SET firstname shekhar
OK
redis> SET lastname gulati
OK
redis> GET firstname
"shekhar"
redis> GET lastname
"gulati"

```

Show moreShow more icon

如果您的键的值是整数，那么可使用 `DECR` 或 `DECRBY` 递减这些值，使用 `INCR` 或 `INCRBY` 递增它们。这些操作在您希望维护一些对象的数量（比如网页的命中次数）的情形中很有用：

```
redis> INCR votes
(integer) 1
redis> INCR votes
(integer) 2
redis> INCR votes
(integer) 3
redis> DECR votes
(integer) 2

```

Show moreShow more icon

其他一些操作（包括 `APPEND` 、 `GETRANGE` 、 `MSET` 和 `STRLENGTH` 也可用于字符串。请参见 参考资料 ，获取所有 Redis 数据类型的完整操作列表的链接。

### 列表

Redis 中的列表是一个有序的字符串集合，您可以向其中添加任意数量的（惟一或非惟一）元素。除了向列表添加元素和从中获取元素的操作之外，Redis 还支持对列表使用取出、推送、范围和其他一些操作。

##### 从 0 开始的列表

在 Redis 中，列表从 0 开始。

作为一个例子，假设您希望维护一个单词列表（无论是否惟一），并获取您最近添加到系统中的三个单词。

要将单词添加到列表中，可以使用 `LPUSH` 命令，它将一个或多个值附加到列表前部：

```
redis> LPUSH words austerity
(integer) 1
redis> LPUSH words socialism moratorium socialism socialism
(integer) 5

```

Show moreShow more icon

使用 `LPUSH` ，最近添加的单词位于列表顶部，以前添加的单词会在后面。可使用 `LRANGE` 命令查看列表中顶部的三个单词：

```
redis> LRANGE words 0 2
1) "socialism"
2) "socialism"
3) "moratorium"

```

Show moreShow more icon

要获得列表的长度，可使用 `LLEN` 命令：

```
redis > LLEN words
(integer) 5

```

Show moreShow more icon

要从列表中删除元素，可使用 `LREM` 命令。例如，使用此命令删除所有 `socialism` ：

```
redis> LREM words 0 socialism
(integer) 2

```

Show moreShow more icon

要删除列表，必须删除 `words` 键：

```
redis 127.0.0.1:6379> DEL words
(integer) 1

```

Show moreShow more icon

### 集合

集合（set）是惟一元素的无序集合。要将一个元素添加到一个集合中，可使用 `SADD` 命令；要获取一个集合的所有成员，可使用 `SMEMBERS` 命令。作为一个例子，假设我希望维护所有添加到系统中的惟一单词的集合（称为 `uniquewords` ）。可以看到，因为我使用了 `SET` ，所以无法将同一个单词 (`socialism`) 添加两次：

```
redis> SADD uniquewords austerity
(integer) 1
redis> SADD uniquewords pragmatic
(integer) 1
redis> SADD uniquewords moratorium
(integer) 1
redis> SADD uniquewords socialism
(integer) 1
redis> SADD uniquewords socialism
(integer) 0

```

Show moreShow more icon

要查看 `uniquewords` 集合中的所有单词，可使用 `SMEMBERS` 命令：

```
redis 127.0.0.1:6379> SMEMBERS uniquewords
1) "moratorium"
2) "austerity"
3) "socialism"
4) "pragmatic"

```

Show moreShow more icon

也可以在多个集合上执行 intersection 和 union 命令。我将创建另一个集合，称为 `newwords` ，并向其中添加一些元素：

```
redis 127.0.0.1:6379> SADD newwords austerity good describe strange
(integer) 4

```

Show moreShow more icon

为了找出两个集合中的所有共同元素（ `uniquewords` 和 `newwords` ），我执行了 `SINTER` 命令：

```
redis 127.0.0.1:6379> SINTER uniquewords newwords
1) "austerity"

```

Show moreShow more icon

为了合并多个集合，我使用了 `SUNION` 命令。可以看到，单词 `austerity` 仅被添加了一次。

```
redis 127.0.0.1:6379> SUNION uniquewords newwords
1) "austerity"
2) "strange"
3) "describe"
4) "socialism"
5) "pragmatic"
6) "good"
7) "moratorium"

```

Show moreShow more icon

### 有序集

有序集是可基于一个分数进行排序并且仅包含惟一元素的集合。每次您向集合中添加一个值，您都会提供一个用于排序的分数。

例如，假设您希望基于单词的长度对一组单词进行排序。使用 `ZADD` 命令将一个元素添加到一个有序集中，使用语法 `ZADD _键分数 值` 。使用 `ZRANGE` 命令按分数查看一个有序集的元素。

```
redis> ZADD wordswithlength 9 austerity
(integer) 1
redis> ZADD wordswithlength 7 furtive
(integer) 1
redis> ZADD wordswithlength 5 bigot
(integer) 1
redis> ZRANGE wordswithlength 0 -1
1) "bigot"
2) "furtive"
3) "austerity"

```

Show moreShow more icon

要获得有序集的大小，可使用 `ZCARD` 命令：

```
redis 127.0.0.1:6379> ZCARD wordswithlength
(integer) 3

```

Show moreShow more icon

### 哈希值

哈希值允许您针对一个哈希值存储一个键-值对。此选项可能对您希望保存具有多个属性的对象的情形很有用，就像以下示例一样：

```
redis> HSET user:1 name shekhar
(integer) 1
redis> HSET user:1 lastname gulati
(integer) 1
redis> HGET user:1
redis> HGET user:1 name
"shekhar"
redis> HGETALL user:1
1) "name"
2) "shekhar"
3) "lastname"
4) "gulati"

```

Show moreShow more icon

现在，您已经了解了 Redis 数据类型，可以尝试使用 Spring Data Redis 框架构建一个应用程序，使用 Redis 作为后端数据存储。

## 开发一个 Spring Redis 应用程序

##### 前提条件

要构建本文的示例应用程序，您需要安装：

- Redis
- JDK 6 或更高版本
- Apache Maven 3 或更高版本
- Spring Tool Suite
- Git

使用 Spring Data Redis，Java 开发人员可以编程方式访问 Redis 并执行相应操作。Spring 框架总是推荐一种基于 POJO（plain old Java object，简单 Java 对象）的编程模型，高度重视生产力、一致性和可移植性。这些值会传输到 Spring Data Redis 项目。

Spring Data Redis 在现有 Redis 客户端库（比如 Jedis、JRedis、redis-protocol 和 RJC，参见 参考资料 ）上提供了一种抽象。通过消除与 Redis 交互所需的样板代码，它使得使用 Redis 键-值数据存储变得很容易，无需了解低级 Redis API。它还提供了一个名为 `RedisTemplate` 的泛化的模板类（类似于 `JDBCTemplate` 或 `HibernateTemplate` ）来与 Redis 进行交互。 `RedisTemplate` 是与 Redis 执行面向对象的交互的主要类。它处理对象序列化和类型转换，使得作为开发人员的您在处理对象时无需担忧序列化和数据转换。

您将构建的应用程序是一个简单的词典应用程序。（可从 GitHub 下载该应用程序的完整代码；参见 参考资料 。）词典是一个单词集合，其中每个单词可能拥有多种含义。这个词典应用程序可轻松地建模为 Redis 列表数据类型，其中每个单词是列表键，单词的含义是它的值。（或者，除了列表，如果您希望含义是惟一的，那么还可以使用集合。如果希望排序含义，那么可以使用有序集。）例如，单词 `astonishing` 可以是一个键， `astounding` 和 `staggering` 是它的值。

首先创建一个简单的单词-含义列表。使用 `redis-server` 命令启动 Redis，使用 `redis-cli` 命令启动客户端。然后，执行这个客户端会话中所示的命令：

```
redis> RPUSH astonishing astounding
(integer) 1
redis> RPUSH astonishing staggering
(integer) 2
redis> LRANGE astonishing 0 -1
1) "astounding"
2) "staggering"

```

Show moreShow more icon

您创建了一个名为 `astonishing` 的列表，并将含义 `astounding` 和 `staggering` 推送到 `astonishing` 列表的末尾。您还使用 `LRANGE` 命令获取了 `astonishing` 的所有含义。

### 使用 Spring Tool Suite 创建一个模板项目

现在必须创建一个 Spring 模板项目，以便可以将它用于应用程序。打开 Spring Tool Suite 并转到 File -> New -> Spring Template Project -> Simple Spring Utility Project ->。在系统提示您时单击 **Yes** 。输入 `dictionary` 作为项目名称，提供顶级包名称（我的为 `com.shekhar.dictionary` ）。现在您在 Spring Tool Suite 工作区中已经有了一个名为 `dictionary` 的样例项目。

### 使用依赖关系更新 pom.xml

`dictionary` 项目没有与 Spring Data Redis 项目相关的依赖关系。要添加它们，可以将您创建的模板项目的 pom.xml 替换为清单 1 中的 pom.xml。（这个文件使用了 Spring Data Redis 项目 1.0.2 版，这是编写本文时的最新版本。）

##### 清单 1\. 放置添加 Spring Data Redis 依赖关系的 pom.xml

```
<project xmlns="http://maven.apache.org/POM/4.0.0"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xsi:schemaLocation="http://maven.apache.org/POM/4.0.0
http://maven.apache.org/xsd/maven-4.0.0.xsd">
<modelVersion>4.0.0</modelVersion>

<groupId>com.shekhar</groupId>
<artifactId>dictionary</artifactId>
<version>0.0.1-SNAPSHOT</version>
<packaging>jar</packaging>

<name>dictionary</name>
<url>http://maven.apache.org</url>

<properties>
      <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
</properties>

<repositories>
      <repository>
         <id>spring-release</id>
         <name>Spring Maven RELEASE Repository</name>
         <url>http://maven.springframework.org/release</url>
      </repository>
</repositories>
<dependencies>
      <dependency>
         <groupId>javax.inject</groupId>
         <artifactId>javax.inject</artifactId>
         <version>1</version>
      </dependency>
      <dependency>
         <groupId>cglib</groupId>
         <artifactId>cglib</artifactId>
         <version>2.2.2</version>
      </dependency>
      <dependency>
         <groupId>org.springframework.data</groupId>
         <artifactId>spring-data-redis</artifactId>
         <version>1.0.2.RELEASE</version>
      </dependency>
      <dependency>
         <groupId>org.springframework</groupId>
         <artifactId>spring-test</artifactId>
         <version>3.1.2.RELEASE</version>
         <scope>test</scope>
      </dependency>
      <dependency>
         <groupId>junit</groupId>
         <artifactId>junit</artifactId>
         <version>4.8.1</version>
         <scope>test</scope>
      </dependency>
</dependencies>

<build>
      <plugins>
         <plugin>
            <groupId>org.apache.maven.plugins</groupId>
            <artifactId>maven-compiler-plugin</artifactId>
            <configuration>
               <source>1.6</source>
               <target>1.6</target>
            </configuration>
         </plugin>
      </plugins>
</build>
</project>

```

Show moreShow more icon

### 配置 `RedisConnectionFactory` 和 `RedisTemplate`

`RedisConnectionFactory` 是一个用来与 Redis 建立连接的线程安全的连接工厂， `RedisConnection` 是连接到 Redis 的一个短期、非线程安全的连接。 `RedisConnection` 提供了与 Redis 命令的一对一映射，而 `RedisConnectionFactory` 提供了有助于消除样板代码的便捷方法。 `RedisConnectionFactory` 使不同 Redis 客户端 API 之间的切换就像定义一个 bean 那么简单。我们将对样例应用程序使用 `JedisConnectionFactory` ，但也可使用其他任何 `ConnectionFactory` 变体。

创建一个 `LocalRedisConfig` 类，如清单 2 所示：

##### 清单 2\. 创建 `LocalRedisConfig` 类

```
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.redis.connection.RedisConnectionFactory;
import org.springframework.data.redis.connection.jedis.JedisConnectionFactory;
import org.springframework.data.redis.core.StringRedisTemplate;

import redis.clients.jedis.JedisPoolConfig;

@Configuration
@ComponentScan(basePackages="com.shekhar.dictionary.dao")
public class LocalRedisConfig {

@Bean
public RedisConnectionFactory jedisConnectionFactory(){
JedisPoolConfig poolConfig = new JedisPoolConfig();
poolConfig.maxActive = 10;
poolConfig.maxIdle = 5;
poolConfig.minIdle = 1;
poolConfig.testOnBorrow = true;
poolConfig.testOnReturn = true;
poolConfig.testWhileIdle = true;
JedisConnectionFactory jedisConnectionFactory = new JedisConnectionFactory(poolConfig);
return jedisConnectionFactory;
}

@Bean
public StringRedisTemplate redisTemplate(){
StringRedisTemplate redisTemplate = new StringRedisTemplate(jedisConnectionFactory());
return redisTemplate;
}
}

```

Show moreShow more icon

`LocalRedisConfig` 定义了 `JedisConnectionFactory` 和 `SpringRedisTemplate` bean。 `SpringRedisTemplate` 是 `RedisTemplate` 的一个处理字符串的特殊版本。

可使用一个 JUnit 测试来测试配置，如清单 3 所示：

##### 清单 3\. 测试配置

```
import javax.inject.Inject;

import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.data.redis.connection.jedis.JedisConnectionFactory;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.test.context.ContextConfiguration;
import org.springframework.test.context.junit4.SpringJUnit4ClassRunner;

@ContextConfiguration(classes = LocalRedisConfig.class)
@RunWith(SpringJUnit4ClassRunner.class)
public class LocalRedisConfigTest {

@Inject
private JedisConnectionFactory jedisConnectionFactory;

@Inject
private StringRedisTemplate redisTemplate;

@Test
public void testJedisConnectionFactory() {
      assertNotNull(jedisConnectionFactory);
}

@Test
public void testRedisTemplate() {
      assertNotNull(redisTemplate);
}

}

```

Show moreShow more icon

### 编写 `DictionaryDao`

现在创建一个 `DictionaryDao` 类，如清单 4 所示：

##### 清单 4\. 创建 `DictionaryDao` 类

```
import org.springframework.data.redis.core.StringRedisTemplate;

@Repository
public class DictionaryDao {

private static final String ALL_UNIQUE_WORDS = "all-unique-words";

private StringRedisTemplate redisTemplate;

@Inject
public DictionaryDao(StringRedisTemplate redisTemplate){
      this.redisTemplate = redisTemplate;
}

public Long addWordWithItsMeaningToDictionary(String word, String meaning) {
      Long index = redisTemplate.opsForList().rightPush(word, meaning);
      return index;
}
}

```

Show moreShow more icon

您可以使用 `DictionaryDao` 执行 Redis 中的操作。从 [清单 4\. 创建 `DictionaryDao` 类](#清单-4-创建-code-dictionarydao-code-类) 中可以看到，您将 `RedisTemplate` （Spring Data Redis 项目中的核心类）注入到了 `DictionaryDao` 中。 [清单 4\. 创建 `DictionaryDao` 类](#清单-4-创建-code-dictionarydao-code-类) 中的 `StringRedisTemplate` 是 `RedisTemplate` 的一个处理字符串数据类型的子类。

`RedisTemplate` 提供了键类型操作，比如 `ValueOperations` 、 `ListOperations` 、 `SetOperations` 、 `HashOperations` 和 `ZSetOperations` 。 [清单 4\. 创建 `DictionaryDao` 类](#清单-4-创建-code-dictionarydao-code-类) 使用 `ListOperations` 将一个新单词存储在 Redis 数据存储中。 `rightPush()` 操作将该单词和它的含义添加到列表末尾处。 `dictionaryDao.addWordWithItsMeaningToDictionary()` 方法返回添加到列表中的元素的索引。

清单 5 显示了一个针对 `dictionaryDao.addWordWithItsMeaningToDictionary()` 的 JUnit 测试案例：

##### 清单 5\. 测试 `dictionaryDao.addWordWithItsMeaningToDictionary()` 方法

```
import static org.hamcrest.CoreMatchers.equalTo;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.CoreMatchers.notNullValue;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertThat;
import static org.junit.matchers.JUnitMatchers.hasItems;

import java.util.List;
import java.util.Set;

import javax.inject.Inject;

import org.junit.After;
import org.junit.Ignore;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.test.context.ContextConfiguration;
import org.springframework.test.context.junit4.SpringJUnit4ClassRunner;
import com.shekhar.dictionary.config.LocalRedisConfig;

@ContextConfiguration(classes = {LocalRedisConfig.class })
@RunWith(SpringJUnit4ClassRunner.class)
public class DictionaryDaoTest {

@Inject
private DictionaryDao dictionaryDao;

@Inject
private StringRedisTemplate redisTemplate;

@Test
public void testAddWordWithItsMeaningToDictionary() {
      String meaning = "To move forward with a bounding, drooping motion.";
      Long index = dictionaryDao.addWordWithItsMeaningToDictionary("lollop",
            meaning);
      assertThat(index, is(notNullValue()));
      assertThat(index, is(equalTo(1L)));
}

```

Show moreShow more icon

第一次运行此测试时，它通过了，这个单词被存储到了 Redis 中。但是，如果随后再次运行该测试，测试会失败，因为 Redis 再次添加了相同的含义，并返回 `2` 作为它的索引。所以，在每次运行之后必须清理 Redis 数据存储，这可以使用 `flushAll()` 或 `flushDb` 服务器命令来完成。 `flushAll()` 命令从数据库中删除所有键，而 `flushDb()` 仅删除当前数据库中的键。清单 6 给出了已修改的测试：

##### 清单 6\. 已修改的 `dictionaryDao.addWordWithItsMeaningToDictionary()` 测试

```
@ContextConfiguration(classes = {LocalRedisConfig.class })
@RunWith(SpringJUnit4ClassRunner.class)
public class DictionaryDaoTest {

@Inject
private DictionaryDao dictionaryDao;

@Inject
private StringRedisTemplate redisTemplate;

@After
public void tearDown() {
      redisTemplate.getConnectionFactory().getConnection().flushDb();
}

@Test
public void testAddWordWithItsMeaningToDictionary() {
      String meaning = "To move forward with a bounding, drooping motion.";
      Long index = dictionaryDao.addWordWithItsMeaningToDictionary("lollop",
            meaning);
      assertThat(index, is(notNullValue()));
      assertThat(index, is(equalTo(1L)));
}

@Test
public void shouldAddMeaningToAWordIfItExists() {
      Long index = dictionaryDao.addWordWithItsMeaningToDictionary("lollop",
            "To move forward with a bounding, drooping motion.");
      assertThat(index, is(notNullValue()));
      assertThat(index, is(equalTo(1L)));
      index = dictionaryDao.addWordWithItsMeaningToDictionary("lollop",
            "To hang loosely; droop; dangle.");
      assertThat(index, is(equalTo(2L)));
}
}

```

Show moreShow more icon

现在您已经拥有将一个单词存储到 Redis 数据存储中的功能，下一步是实现获取一个单词的所有含义的功能。这一步可使用 `List` 的 `range` 操作轻松完成。 `range()` 方法接受三个参数：键的名称、范围的起点和范围的终点。要获取一个单词的所有含义，可使用 0 表示起点，使用 -1 表示终点：

```
public List<String> getAllTheMeaningsForAWord(String word) {
      List<String> meanings = redisTemplate.opsForList().range(word, 0, -1);
      return meanings;
}

```

Show moreShow more icon

您想在该应用程序中实现的另一个基本功能是删除单词。可以使用 `RedisTemplate` 类的 `delete` 实现此目标。 `delete` 操作接受一个键集合：

```
public void removeWord(String word) {
      redisTemplate.delete(Arrays.asList(word));
}

public void removeWords(String... words) {
      redisTemplate.delete(Arrays.asList(words));
}

```

Show moreShow more icon

针对新添加的 `read` 和 `delete` 操作的 JUnit 测试案例如清单 7 所示：

##### 清单 7\. 针对 `read` 和 `delete` 操作的 JUnit 测试

```
@Test
public void shouldGetAllTheMeaningForAWord() {
      setupOneWord();
      List<String> allMeanings = dictionaryDao
            .getAllTheMeaningsForAWord("lollop");
      assertThat(allMeanings.size(), is(equalTo(2)));
      assertThat(
            allMeanings,
            hasItems("To move forward with a bounding, drooping motion.",
                  "To hang loosely; droop; dangle."));
}

@Test
public void shouldDeleteAWordFromDictionary() throws Exception {
      setupOneWord();
      dictionaryDao.removeWord("lollop");
      List<String> allMeanings = dictionaryDao
            .getAllTheMeaningsForAWord("lollop");
      assertThat(allMeanings.size(), is(equalTo(0)));
}

@Test
public void shouldDeleteMultipleWordsFromDictionary() {
      setupTwoWords();
      dictionaryDao.removeWords("fain", "lollop");
      List<String> allMeaningsForLollop = dictionaryDao
            .getAllTheMeaningsForAWord("lollop");
      List<String> allMeaningsForFain = dictionaryDao
            .getAllTheMeaningsForAWord("fain");
      assertThat(allMeaningsForLollop.size(), is(equalTo(0)));
      assertThat(allMeaningsForFain.size(), is(equalTo(0)));
}
}

```

Show moreShow more icon

## 结束语

本文介绍了 Redis 和如何使用 Spring Data Redis 项目构建 Spring Redis 应用程序。我还介绍了 Redis 的基本知识、它的数据模型、数据类型，以及如何开始使用 Spring Redis。Redis 包含其他的丰富功能，包括 Redis PubSub 和 MultI-EXEC 等。请参见本文的 参考资料 ，了解这些主题的信息。

本文翻译自： [Archive: Develop Spring Redis applications](https://www.ibm.com/developerworks/library/os-springredis/index.html)（2013-08-21）