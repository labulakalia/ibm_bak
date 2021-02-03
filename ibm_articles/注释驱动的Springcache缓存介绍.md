# 注释驱动的 Spring cache 缓存介绍
spring 3.1 激动人心的新特性

**标签:** Java,Spring

[原文链接](https://developer.ibm.com/zh/articles/os-cn-spring-cache/)

赵才文

发布: 2013-03-15

* * *

## 概述

Spring 3.1 引入了激动人心的基于注释（annotation）的缓存（cache）技术，它本质上不是一个具体的缓存实现方案（例如 EHCache 或者 OSCache），而是一个对缓存使用的抽象，通过在既有代码中添加少量它定义的各种 annotation，即能够达到缓存方法的返回对象的效果。

Spring 的缓存技术还具备相当的灵活性，不仅能够使用 SpEL（Spring Expression Language）来定义缓存的 key 和各种 condition，还提供开箱即用的缓存临时存储方案，也支持和主流的专业缓存例如 EHCache 集成。

其特点总结如下：

- 通过少量的配置 annotation 注释即可使得既有代码支持缓存
- 支持开箱即用 Out-Of-The-Box，即不用安装和部署额外第三方组件即可使用缓存
- 支持 Spring Express Language，能使用对象的任何属性或者方法来定义缓存的 key 和 condition
- 支持 AspectJ，并通过其实现任何方法的缓存支持
- 支持自定义 key 和自定义缓存管理者，具有相当的灵活性和扩展性

本文将针对上述特点对 Spring cache 进行详细的介绍，主要通过一个简单的例子和原理介绍展开，然后我们将一起看一个比较实际的缓存例子，最后会介绍 spring cache 的使用限制和注意事项。OK，Let ‘ s begin!

## 原来我们是怎么做的

这里先展示一个完全自定义的缓存实现，即不用任何第三方的组件来实现某种对象的内存缓存。

场景是：对一个账号查询方法做缓存，以账号名称为 key，账号对象为 value，当以相同的账号名称查询账号的时候，直接从缓存中返回结果，否则更新缓存。账号查询服务还支持 reload 缓存（即清空缓存）。

首先定义一个实体类：账号类，具备基本的 id 和 name 属性，且具备 getter 和 setter 方法

##### 清单 1\. Account.java

```
package cacheOfAnno;

public class Account {
private int id;
private String name;

public Account(String name) {
     this.name = name;
}
public int getId() {
     return id;
}
public void setId(int id) {
     this.id = id;
}
public String getName() {
     return name;
}
public void setName(String name) {
     this.name = name;
}
}

```

Show moreShow more icon

然后定义一个缓存管理器，这个管理器负责实现缓存逻辑，支持对象的增加、修改和删除，支持值对象的泛型。如下：

##### 清单 2\. MyCacheManager.java

```
package oldcache;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

public class MyCacheManager<T> {
private Map<String,T> cache =
       new ConcurrentHashMap<String,T>();

public T getValue(Object key) {
     return cache.get(key);
}

public void addOrUpdateCache(String key,T value) {
     cache.put(key, value);
}

public void evictCache(String key) {// 根据 key 来删除缓存中的一条记录
     if(cache.containsKey(key)) {
       cache.remove(key);
     }
}

public void evictCache() {// 清空缓存中的所有记录
     cache.clear();
}
}

```

Show moreShow more icon

好，现在我们有了实体类和一个缓存管理器，还需要一个提供账号查询的服务类，此服务类使用缓存管理器来支持账号查询缓存，如下：

##### 清单 3\. MyAccountService.java

```
package oldcache;

import cacheOfAnno.Account;

public class MyAccountService {
private MyCacheManager<Account> cacheManager;

public MyAccountService() {
     cacheManager = new MyCacheManager<Account>();// 构造一个缓存管理器
}

public Account getAccountByName(String acctName) {
     Account result = cacheManager.getValue(acctName);// 首先查询缓存
     if(result!=null) {
       System.out.println("get from cache..."+acctName);
       return result;// 如果在缓存中，则直接返回缓存的结果
     }
     result = getFromDB(acctName);// 否则到数据库中查询
     if(result!=null) {// 将数据库查询的结果更新到缓存中
       cacheManager.addOrUpdateCache(acctName, result);
     }
     return result;
}

public void reload() {
     cacheManager.evictCache();
}

private Account getFromDB(String acctName) {
     System.out.println("real querying db..."+acctName);
     return new Account(acctName);
}
}

```

Show moreShow more icon

现在我们开始写一个测试类，用于测试刚才的缓存是否有效

##### 清单 4\. Main.java

```
package oldcache;

public class Main {

public static void main(String[] args) {
     MyAccountService s = new MyAccountService();
     // 开始查询账号
     s.getAccountByName("somebody");// 第一次查询，应该是数据库查询
     s.getAccountByName("somebody");// 第二次查询，应该直接从缓存返回

     s.reload();// 重置缓存
     System.out.println("after reload...");

     s.getAccountByName("somebody");// 应该是数据库查询
     s.getAccountByName("somebody");// 第二次查询，应该直接从缓存返回

}

}

```

Show moreShow more icon

按照分析，执行结果应该是：首先从数据库查询，然后直接返回缓存中的结果，重置缓存后，应该先从数据库查询，然后返回缓存中的结果，实际的执行结果如下：

##### 清单 5\. 运行结果

```
real querying db...somebody// 第一次从数据库加载
get from cache...somebody// 第二次从缓存加载
after reload...// 清空缓存
real querying db...somebody// 又从数据库加载
get from cache...somebody// 从缓存加载

```

Show moreShow more icon

可以看出我们的缓存起效了，但是这种自定义的缓存方案有如下劣势：

- 缓存代码和业务代码耦合度太高，如上面的例子，AccountService 中的 getAccountByName（）方法中有了太多缓存的逻辑，不便于维护和变更
- 不灵活，这种缓存方案不支持按照某种条件的缓存，比如只有某种类型的账号才需要缓存，这种需求会导致代码的变更
- 缓存的存储这块写的比较死，不能灵活的切换为使用第三方的缓存模块

如果你的代码中有上述代码的影子，那么你可以考虑按照下面的介绍来优化一下你的代码结构了，也可以说是简化，你会发现，你的代码会变得优雅的多！

## Hello World，注释驱动的 Spring Cache

### Hello World 的实现目标

本 Hello World 类似于其他任何的 Hello World 程序，从最简单实用的角度展现 spring cache 的魅力，它基于刚才自定义缓存方案的实体类 Account.java，重新定义了 AccountService.java 和测试类 Main.java（注意这个例子不用自己定义缓存管理器，因为 spring 已经提供了缺省实现）

### 需要的 jar 包

为了实用 spring cache 缓存方案，在工程的 classpath 必须具备下列 jar 包。

##### 图 1\. 工程依赖的 jar 包图

![图 1. 工程依赖的 jar 包图](../ibm_articles_img/os-cn-spring-cache_images_image001.jpg)

注意这里我引入的是最新的 spring 3.2.0.M1 版本 jar 包，其实只要是 spring 3.1 以上，都支持 spring cache。其中 spring-context-\*.jar 包含了 cache 需要的类。

### 定义实体类、服务类和相关配置文件

实体类就是上面自定义缓存方案定义的 Account.java，这里重新定义了服务类，如下：

##### 清单 6\. AccountService.java

```
package cacheOfAnno;

import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;

public class AccountService {
@Cacheable(value="accountCache")// 使用了一个缓存名叫 accountCache
public Account getAccountByName(String userName) {
     // 方法内部实现不考虑缓存逻辑，直接实现业务
     System.out.println("real query account."+userName);
     return getFromDB(userName);
}

private Account getFromDB(String acctName) {
     System.out.println("real querying db..."+acctName);
     return new Account(acctName);
}
}

```

Show moreShow more icon

注意，此类的 getAccountByName 方法上有一个注释 annotation，即 @Cacheable(value=”accountCache”)，这个注释的意思是，当调用这个方法的时候，会从一个名叫 accountCache 的缓存中查询，如果没有，则执行实际的方法（即查询数据库），并将执行的结果存入缓存中，否则返回缓存中的对象。这里的缓存中的 key 就是参数 userName，value 就是 Account 对象。”accountCache”缓存是在 spring\*.xml 中定义的名称。

好，因为加入了 spring，所以我们还需要一个 spring 的配置文件来支持基于注释的缓存

##### 清单 7\. Spring-cache-anno.xml

```
<beans xmlns="http://www.springframework.org/schema/beans"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns:cache="http://www.springframework.org/schema/cache"
    xmlns:p="http://www.springframework.org/schema/p"
xsi:schemaLocation="http://www.springframework.org/schema/beans
http://www.springframework.org/schema/beans/spring-beans.xsd
     http://www.springframework.org/schema/cache
     http://www.springframework.org/schema/cache/spring-cache.xsd">

<cache:annotation-driven />

<bean id="accountServiceBean" class="cacheOfAnno.AccountService"/>

    <!-- generic cache manager -->
<bean id="cacheManager"
class="org.springframework.cache.support.SimpleCacheManager">
     <property name="caches">
       <set>
         <bean
           class="org.springframework.cache.concurrent.ConcurrentMapCacheFactoryBean"
           p:name="default" />

         <bean
           class="org.springframework.cache.concurrent.ConcurrentMapCacheFactoryBean"
           p:name="accountCache" />
       </set>
     </property>
</bean>
</beans>

```

Show moreShow more icon

注意这个 spring 配置文件有一个关键的支持缓存的配置项： `<cache:annotation-driven />，` 这个配置项缺省使用了一个名字叫 cacheManager 的缓存管理器，这个缓存管理器有一个 spring 的缺省实现，即 org.springframework.cache.support.SimpleCacheManager，这个缓存管理器实现了我们刚刚自定义的缓存管理器的逻辑，它需要配置一个属性 caches，即此缓存管理器管理的缓存集合，除了缺省的名字叫 default 的缓存，我们还自定义了一个名字叫 accountCache 的缓存，使用了缺省的内存存储方案 ConcurrentMapCacheFactoryBean，它是基于 java.util.concurrent.ConcurrentHashMap 的一个内存缓存实现方案。

OK，现在我们具备了测试条件，测试代码如下：

##### 清单 8\. Main.java

```
package cacheOfAnno;

import org.springframework.context.ApplicationContext;
import org.springframework.context.support.ClassPathXmlApplicationContext;

public class Main {
public static void main(String[] args) {
     ApplicationContext context = new ClassPathXmlApplicationContext(
        "spring-cache-anno.xml");// 加载 spring 配置文件

     AccountService s = (AccountService) context.getBean("accountServiceBean");
     // 第一次查询，应该走数据库
     System.out.print("first query...");
     s.getAccountByName("somebody");
     // 第二次查询，应该不查数据库，直接返回缓存的值
     System.out.print("second query...");
     s.getAccountByName("somebody");
     System.out.println();
}
}

```

Show moreShow more icon

上面的测试代码主要进行了两次查询，第一次应该会查询数据库，第二次应该返回缓存，不再查数据库，我们执行一下，看看结果

##### 清单 9\. 执行结果

```
first query...real query account.somebody// 第一次查询
real querying db...somebody// 对数据库进行了查询
second query...// 第二次查询，没有打印数据库查询日志，直接返回了缓存中的结果

```

Show moreShow more icon

可以看出我们设置的基于注释的缓存起作用了，而在 AccountService.java 的代码中，我们没有看到任何的缓存逻辑代码，只有一行注释：@Cacheable(value=”accountCache”)，就实现了基本的缓存方案，是不是很强大？

### 如何清空缓存

好，到目前为止，我们的 spring cache 缓存程序已经运行成功了，但是还不完美，因为还缺少一个重要的缓存管理逻辑：清空缓存，当账号数据发生变更，那么必须要清空某个缓存，另外还需要定期的清空所有缓存，以保证缓存数据的可靠性。

为了加入清空缓存的逻辑，我们只要对 AccountService.java 进行修改，从业务逻辑的角度上看，它有两个需要清空缓存的地方

- 当外部调用更新了账号，则我们需要更新此账号对应的缓存
- 当外部调用说明重新加载，则我们需要清空所有缓存

##### 清单 10\. AccountService.java

```
package cacheOfAnno;

import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;

public class AccountService {
@Cacheable(value="accountCache")// 使用了一个缓存名叫 accountCache
public Account getAccountByName(String userName) {
     // 方法内部实现不考虑缓存逻辑，直接实现业务
     return getFromDB(userName);
}
@CacheEvict(value="accountCache",key="#account.getName()")// 清空 accountCache 缓存  public void updateAccount(Account account) {
     updateDB(account);
}

@CacheEvict(value="accountCache",allEntries=true)// 清空 accountCache 缓存
public void reload() {
}

private Account getFromDB(String acctName) {
     System.out.println("real querying db..."+acctName);
     return new Account(acctName);
}

private void updateDB(Account account) {
     System.out.println("real update db..."+account.getName());
}

}

```

Show moreShow more icon

##### 清单 11\. Main.java

```
package cacheOfAnno;

import org.springframework.context.ApplicationContext;
import org.springframework.context.support.ClassPathXmlApplicationContext;

public class Main {

public static void main(String[] args) {
     ApplicationContext context = new ClassPathXmlApplicationContext(
        "spring-cache-anno.xml");// 加载 spring 配置文件

     AccountService s = (AccountService) context.getBean("accountServiceBean");
     // 第一次查询，应该走数据库
     System.out.print("first query...");
     s.getAccountByName("somebody");
     // 第二次查询，应该不查数据库，直接返回缓存的值
     System.out.print("second query...");
     s.getAccountByName("somebody");
     System.out.println();

     System.out.println("start testing clear cache...");    // 更新某个记录的缓存，首先构造两个账号记录，然后记录到缓存中
     Account account1 = s.getAccountByName("somebody1");
     Account account2 = s.getAccountByName("somebody2");
     // 开始更新其中一个    account1.setId(1212);
     s.updateAccount(account1);
     s.getAccountByName("somebody1");// 因为被更新了，所以会查询数据库    s.getAccountByName("somebody2");// 没有更新过，应该走缓存    s.getAccountByName("somebody1");// 再次查询，应该走缓存    // 更新所有缓存
     s.reload();
     s.getAccountByName("somebody1");// 应该会查询数据库    s.getAccountByName("somebody2");// 应该会查询数据库    s.getAccountByName("somebody1");// 应该走缓存    s.getAccountByName("somebody2");// 应该走缓存
}
}

```

Show moreShow more icon

##### 清单 12\. 运行结果

```
first query...real querying db...somebody
second query...
start testing clear cache...
real querying db...somebody1
real querying db...somebody2
real update db...somebody1
real querying db...somebody1
real querying db...somebody1
real querying db...somebody2

```

Show moreShow more icon

结果和我们期望的一致，所以，我们可以看出，spring cache 清空缓存的方法很简单，就是通过 @CacheEvict 注释来标记要清空缓存的方法，当这个方法被调用后，即会清空缓存。注意其中一个 @CacheEvict(value=”accountCache”,key=”#account.getName()”)，其中的 Key 是用来指定缓存的 key 的，这里因为我们保存的时候用的是 account 对象的 name 字段，所以这里还需要从参数 account 对象中获取 name 的值来作为 key，前面的 # 号代表这是一个 SpEL 表达式，此表达式可以遍历方法的参数对象，具体语法可以参考 Spring 的相关文档手册。

### 如何按照条件操作缓存

前面介绍的缓存方法，没有任何条件，即所有对 accountService 对象的 getAccountByName 方法的调用都会起动缓存效果，不管参数是什么值，如果有一个需求，就是只有账号名称的长度小于等于 4 的情况下，才做缓存，大于 4 的不使用缓存，那怎么实现呢？

Spring cache 提供了一个很好的方法，那就是基于 SpEL 表达式的 condition 定义，这个 condition 是 @Cacheable 注释的一个属性，下面我来演示一下

##### 清单 13\. AccountService.java（getAccountByName 方法修订，支持条件）

```
@Cacheable(value="accountCache",condition="#userName.length() <= 4")// 缓存名叫 accountCache
public Account getAccountByName(String userName) {
// 方法内部实现不考虑缓存逻辑，直接实现业务
return getFromDB(userName);
}

```

Show moreShow more icon

注意其中的 condition=”#userName.length() <=4”，这里使用了 SpEL 表达式访问了参数 userName 对象的 length() 方法，条件表达式返回一个布尔值，true/false，当条件为 true，则进行缓存操作，否则直接调用方法执行的返回结果。

##### 清单 14\. 测试方法

```
s.getAccountByName("somebody");// 长度大于 4，不会被缓存
s.getAccountByName("sbd");// 长度小于 4，会被缓存
s.getAccountByName("somebody");// 还是查询数据库
s.getAccountByName("sbd");// 会从缓存返回

```

Show moreShow more icon

##### 清单 15\. 运行结果

```
real querying db...somebody
real querying db...sbd
real querying db...somebody

```

Show moreShow more icon

可见对长度大于 4 的账号名 (somebody) 没有缓存，每次都查询数据库。

### 如果有多个参数，如何进行 key 的组合

假设 AccountService 现在有一个需求，要求根据账号名、密码和是否发送日志查询账号信息，很明显，这里我们需要根据账号名、密码对账号对象进行缓存，而第三个参数”是否发送日志”对缓存没有任何影响。所以，我们可以利用 SpEL 表达式对缓存 key 进行设计

##### 清单 16\. Account.java（增加 password 属性）

```
private String password;
public String getPassword() {
return password;
}
public void setPassword(String password) {
this.password = password;
}

```

Show moreShow more icon

##### 清单 17\. AccountService.java（增加 getAccount 方法，支持组合 key）

```
@Cacheable(value="accountCache",key="#userName.concat(#password)")
public Account getAccount(String userName,String password,boolean sendLog) {
// 方法内部实现不考虑缓存逻辑，直接实现业务
return getFromDB(userName,password);

}

```

Show moreShow more icon

注意上面的 key 属性，其中引用了方法的两个参数 userName 和 password，而 sendLog 属性没有考虑，因为其对缓存没有影响。

##### 清单 18\. Main.java

```
public static void main(String[] args) {
ApplicationContext context = new ClassPathXmlApplicationContext(
      "spring-cache-anno.xml");// 加载 spring 配置文件

AccountService s = (AccountService) context.getBean("accountServiceBean");
s.getAccount("somebody", "123456", true);// 应该查询数据库
s.getAccount("somebody", "123456", true);// 应该走缓存
s.getAccount("somebody", "123456", false);// 应该走缓存
s.getAccount("somebody", "654321", true);// 应该查询数据库
s.getAccount("somebody", "654321", true);// 应该走缓存
}

```

Show moreShow more icon

上述测试，是采用了相同的账号，不同的密码组合进行查询，那么一共有两种组合情况，所以针对数据库的查询应该只有两次。

##### 清单 19\. 运行结果

```
real querying db...userName=somebody password=123456
real querying db...userName=somebody password=654321

```

Show moreShow more icon

和我们预期的一致。

### 如何做到：既要保证方法被调用，又希望结果被缓存

根据前面的例子，我们知道，如果使用了 @Cacheable 注释，则当重复使用相同参数调用方法的时候，方法本身不会被调用执行，即方法本身被略过了，取而代之的是方法的结果直接从缓存中找到并返回了。

现实中并不总是如此，有些情况下我们希望方法一定会被调用，因为其除了返回一个结果，还做了其他事情，例如记录日志，调用接口等，这个时候，我们可以用 @CachePut 注释，这个注释可以确保方法被执行，同时方法的返回值也被记录到缓存中。

##### 清单 20\. AccountService.java

```
@Cacheable(value="accountCache")// 使用了一个缓存名叫 accountCache
public Account getAccountByName(String userName) {
// 方法内部实现不考虑缓存逻辑，直接实现业务
return getFromDB(userName);
}
@CachePut(value="accountCache",key="#account.getName()")// 更新 accountCache 缓存
public Account updateAccount(Account account) {
return updateDB(account);
}
private Account updateDB(Account account) {
System.out.println("real updating db..."+account.getName());
return account;
}

```

Show moreShow more icon

##### 清单 21\. Main.java

```
public static void main(String[] args) {
ApplicationContext context = new ClassPathXmlApplicationContext(
      "spring-cache-anno.xml");// 加载 spring 配置文件

AccountService s = (AccountService) context.getBean("accountServiceBean");

Account account = s.getAccountByName("someone");
account.setPassword("123");
s.updateAccount(account);
account.setPassword("321");
s.updateAccount(account);
account = s.getAccountByName("someone");
System.out.println(account.getPassword());
}

```

Show moreShow more icon

如上面的代码所示，我们首先用 getAccountByName 方法查询一个人 someone 的账号，这个时候会查询数据库一次，但是也记录到缓存中了。然后我们修改了密码，调用了 updateAccount 方法，这个时候会执行数据库的更新操作且记录到缓存，我们再次修改密码并调用 updateAccount 方法，然后通过 getAccountByName 方法查询，这个时候，由于缓存中已经有数据，所以不会查询数据库，而是直接返回最新的数据，所以打印的密码应该是”321”

##### 清单 22\. 运行结果

```
real querying db...someone
real updating db...someone
real updating db...someone
321

```

Show moreShow more icon

和分析的一样，只查询了一次数据库，更新了两次数据库，最终的结果是最新的密码。说明 @CachePut 确实可以保证方法被执行，且结果一定会被缓存。

### @Cacheable、@CachePut、@CacheEvict 注释介绍

通过上面的例子，我们可以看到 spring cache 主要使用两个注释标签，即 @Cacheable、@CachePut 和 @CacheEvict，我们总结一下其作用和配置方法。

表 1\. @Cacheable 作用和配置方法 {: #表-1-cacheable-作用和配置方法}

**@Cacheable 的作用** 主要针对方法配置，能够根据方法的请求参数对其结果进行缓存 **@Cacheable 主要的参数** value缓存的名称，在 spring 配置文件中定义，必须指定至少一个例如：

@Cacheable(value=”mycache”) 或者

@Cacheable(value={”cache1”,”cache2”}key缓存的 key，可以为空，如果指定要按照 SpEL 表达式编写，如果不指定，则缺省按照方法的所有参数进行组合例如：

@Cacheable(value=”testcache”,key=”#userName”)condition缓存的条件，可以为空，使用 SpEL 编写，返回 true 或者 false，只有为 true 才进行缓存例如：

@Cacheable(value=”testcache”,condition=”#userName.length()>2”)

表 2\. @CachePut 作用和配置方法 {: #表-2-cacheput-作用和配置方法}

**@CachePut 的作用** 主要针对方法配置，能够根据方法的请求参数对其结果进行缓存，和 @Cacheable 不同的是，它每次都会触发真实方法的调用 **@CachePut 主要的参数** value缓存的名称，在 spring 配置文件中定义，必须指定至少一个例如：

@Cacheable(value=”mycache”) 或者

@Cacheable(value={”cache1”,”cache2”}key缓存的 key，可以为空，如果指定要按照 SpEL 表达式编写，如果不指定，则缺省按照方法的所有参数进行组合例如：

@Cacheable(value=”testcache”,key=”#userName”)condition缓存的条件，可以为空，使用 SpEL 编写，返回 true 或者 false，只有为 true 才进行缓存例如：

@Cacheable(value=”testcache”,condition=”#userName.length()>2”)

表 3\. @CacheEvict 作用和配置方法 {: #表-3-cacheevict-作用和配置方法}

**@CachEvict 的作用** 主要针对方法配置，能够根据一定的条件对缓存进行清空 **@CacheEvict 主要的参数** value缓存的名称，在 spring 配置文件中定义，必须指定至少一个例如：

@CachEvict(value=”mycache”) 或者

@CachEvict(value={”cache1”,”cache2”}key缓存的 key，可以为空，如果指定要按照 SpEL 表达式编写，如果不指定，则缺省按照方法的所有参数进行组合例如：

@CachEvict(value=”testcache”,key=”#userName”)condition缓存的条件，可以为空，使用 SpEL 编写，返回 true 或者 false，只有为 true 才清空缓存例如：

@CachEvict(value=”testcache”,

condition=”#userName.length()>2”)allEntries是否清空所有缓存内容，缺省为 false，如果指定为 true，则方法调用后将立即清空所有缓存例如：

@CachEvict(value=”testcache”,allEntries=true)beforeInvocation是否在方法执行前就清空，缺省为 false，如果指定为 true，则在方法还没有执行的时候就清空缓存，缺省情况下，如果方法执行抛出异常，则不会清空缓存例如：

@CachEvict(value=”testcache”，beforeInvocation=true)

## 基本原理

和 spring 的事务管理类似，spring cache 的关键原理就是 spring AOP，通过 spring AOP，其实现了在方法调用前、调用后获取方法的入参和返回值，进而实现了缓存的逻辑。我们来看一下下面这个图：

##### 图 2\. 原始方法调用图

![图 2. 原始方法调用图](../ibm_articles_img/os-cn-spring-cache_images_image002.jpg)

上图显示，当客户端”Calling code”调用一个普通类 Plain Object 的 foo() 方法的时候，是直接作用在 pojo 类自身对象上的，客户端拥有的是被调用者的直接的引用。

而 Spring cache 利用了 Spring AOP 的动态代理技术，即当客户端尝试调用 pojo 的 foo（）方法的时候，给他的不是 pojo 自身的引用，而是一个动态生成的代理类

##### 图 3\. 动态代理调用图

![图 3. 动态代理调用图](../ibm_articles_img/os-cn-spring-cache_images_image003.jpg)

如上图所示，这个时候，实际客户端拥有的是一个代理的引用，那么在调用 foo() 方法的时候，会首先调用 proxy 的 foo() 方法，这个时候 proxy 可以整体控制实际的 pojo.foo() 方法的入参和返回值，比如缓存结果，比如直接略过执行实际的 foo() 方法等，都是可以轻松做到的。

## 扩展性

直到现在，我们已经学会了如何使用开箱即用的 spring cache，这基本能够满足一般应用对缓存的需求，但现实总是很复杂，当你的用户量上去或者性能跟不上，总需要进行扩展，这个时候你或许对其提供的内存缓存不满意了，因为其不支持高可用性，也不具备持久化数据能力，这个时候，你就需要自定义你的缓存方案了，还好，spring 也想到了这一点。

我们先不考虑如何持久化缓存，毕竟这种第三方的实现方案很多，我们要考虑的是，怎么利用 spring 提供的扩展点实现我们自己的缓存，且在不改原来已有代码的情况下进行扩展。

首先，我们需要提供一个 CacheManager 接口的实现，这个接口告诉 spring 有哪些 cache 实例，spring 会根据 cache 的名字查找 cache 的实例。另外还需要自己实现 Cache 接口，Cache 接口负责实际的缓存逻辑，例如增加键值对、存储、查询和清空等。利用 Cache 接口，我们可以对接任何第三方的缓存系统，例如 EHCache、OSCache，甚至一些内存数据库例如 memcache 或者 h2db 等。下面我举一个简单的例子说明如何做。

##### 清单 23\. MyCacheManager

```
package cacheOfAnno;

import java.util.Collection;

import org.springframework.cache.support.AbstractCacheManager;

public class MyCacheManager extends AbstractCacheManager {
private Collection<? extends MyCache> caches;

/**
* Specify the collection of Cache instances to use for this CacheManager.
*/
public void setCaches(Collection<? extends MyCache> caches) {
     this.caches = caches;
}

@Override
protected Collection<? extends MyCache> loadCaches() {
     return this.caches;
}

}

```

Show moreShow more icon

上面的自定义的 CacheManager 实际继承了 spring 内置的 AbstractCacheManager，实际上仅仅管理 MyCache 类的实例。

##### 清单 24\. MyCache

```
package cacheOfAnno;

import java.util.HashMap;
import java.util.Map;

import org.springframework.cache.Cache;
import org.springframework.cache.support.SimpleValueWrapper;

public class MyCache implements Cache {
private String name;
private Map<String,Account> store = new HashMap<String,Account>();;

public MyCache() {
}

public MyCache(String name) {
     this.name = name;
}

@Override
public String getName() {
     return name;
}

public void setName(String name) {
     this.name = name;
}

@Override
public Object getNativeCache() {
     return store;
}

@Override
public ValueWrapper get(Object key) {
     ValueWrapper result = null;
     Account thevalue = store.get(key);
     if(thevalue!=null) {
       thevalue.setPassword("from mycache:"+name);
       result = new SimpleValueWrapper(thevalue);
     }
     return result;
}

@Override
public void put(Object key, Object value) {
     Account thevalue = (Account)value;
     store.put((String)key, thevalue);
}

@Override
public void evict(Object key) {
}

@Override
public void clear() {
}
}

```

Show moreShow more icon

上面的自定义缓存只实现了很简单的逻辑，但这是我们自己做的，也很令人激动是不是，主要看 get 和 put 方法，其中的 get 方法留了一个后门，即所有的从缓存查询返回的对象都将其 password 字段设置为一个特殊的值，这样我们等下就能演示”我们的缓存确实在起作用！”了。

这还不够，spring 还不知道我们写了这些东西，需要通过 spring\*.xml 配置文件告诉它

##### 清单 25\. Spring-cache-anno.xml

```
<beans xmlns="http://www.springframework.org/schema/beans"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns:cache="http://www.springframework.org/schema/cache"
xmlns:p="http://www.springframework.org/schema/p"
xsi:schemaLocation="http://www.springframework.org/schema/beans
http://www.springframework.org/schema/beans/spring-beans.xsd
     http://www.springframework.org/schema/cache
     http://www.springframework.org/schema/cache/spring-cache.xsd">

<cache:annotation-driven />

<bean id="accountServiceBean" class="cacheOfAnno.AccountService"/>

    <!-- generic cache manager -->
<bean id="cacheManager" class="cacheOfAnno.MyCacheManager">
     <property name="caches">
       <set>
         <bean
           class="cacheOfAnno.MyCache"
           p:name="accountCache" />
       </set>
     </property>
</bean>

</beans>

```

Show moreShow more icon

注意上面配置文件的黑体字，这些配置说明了我们的 cacheManager 和我们自己的 cache 实例。

好，什么都不说，测试！

##### 清单 26\. Main.java

```
public static void main(String[] args) {
ApplicationContext context = new ClassPathXmlApplicationContext(
      "spring-cache-anno.xml");// 加载 spring 配置文件

AccountService s = (AccountService) context.getBean("accountServiceBean");

Account account = s.getAccountByName("someone");
System.out.println("passwd="+account.getPassword());
account = s.getAccountByName("someone");
System.out.println("passwd="+account.getPassword());
}

```

Show moreShow more icon

上面的测试代码主要是先调用 getAccountByName 进行一次查询，这会调用数据库查询，然后缓存到 mycache 中，然后我打印密码，应该是空的；下面我再次查询 someone 的账号，这个时候会从 mycache 中返回缓存的实例，记得上面的后门么？我们修改了密码，所以这个时候打印的密码应该是一个特殊的值

##### 清单 27\. 运行结果

```
real querying db...someone
passwd=null
passwd=from mycache:accountCache

```

Show moreShow more icon

结果符合预期，即第一次查询数据库，且密码为空，第二次打印了一个特殊的密码。说明我们的 myCache 起作用了。

## 注意和限制

### 基于 proxy 的 spring aop 带来的内部调用问题

上面介绍过 spring cache 的原理，即它是基于动态生成的 proxy 代理机制来对方法的调用进行切面，这里关键点是对象的引用问题，如果对象的方法是内部调用（即 this 引用）而不是外部引用，则会导致 proxy 失效，那么我们的切面就失效，也就是说上面定义的各种注释包括 @Cacheable、@CachePut 和 @CacheEvict 都会失效，我们来演示一下。

##### 清单 28\. AccountService.java

```
public Account getAccountByName2(String userName) {
return this.getAccountByName(userName);
}

@Cacheable(value="accountCache")// 使用了一个缓存名叫 accountCache
public Account getAccountByName(String userName) {
// 方法内部实现不考虑缓存逻辑，直接实现业务
return getFromDB(userName);
}

```

Show moreShow more icon

上面我们定义了一个新的方法 getAccountByName2，其自身调用了 getAccountByName 方法，这个时候，发生的是内部调用（this），所以没有走 proxy，导致 spring cache 失效

##### 清单 29\. Main.java

```
public static void main(String[] args) {
ApplicationContext context = new ClassPathXmlApplicationContext(
      "spring-cache-anno.xml");// 加载 spring 配置文件

AccountService s = (AccountService) context.getBean("accountServiceBean");

s.getAccountByName2("someone");
s.getAccountByName2("someone");
s.getAccountByName2("someone");
}

```

Show moreShow more icon

##### 清单 30\. 运行结果

```
real querying db...someone
real querying db...someone
real querying db...someone

```

Show moreShow more icon

可见，结果是每次都查询数据库，缓存没起作用。要避免这个问题，就是要避免对缓存方法的内部调用，或者避免使用基于 proxy 的 AOP 模式，可以使用基于 aspectJ 的 AOP 模式来解决这个问题。

### @CacheEvict 的可靠性问题

我们看到，@CacheEvict 注释有一个属性 beforeInvocation，缺省为 false，即缺省情况下，都是在实际的方法执行完成后，才对缓存进行清空操作。期间如果执行方法出现异常，则会导致缓存清空不被执行。我们演示一下

##### 清单 31\. AccountService.java

```
@CacheEvict(value="accountCache",allEntries=true)// 清空 accountCache 缓存
public void reload() {
throw new RuntimeException();
}

```

Show moreShow more icon

注意上面的代码，我们在 reload 的时候抛出了运行期异常，这会导致清空缓存失败。

##### 清单 32\. Main.java

```
public static void main(String[] args) {
ApplicationContext context = new ClassPathXmlApplicationContext(
      "spring-cache-anno.xml");// 加载 spring 配置文件

AccountService s = (AccountService) context.getBean("accountServiceBean");

s.getAccountByName("someone");
s.getAccountByName("someone");
try {
     s.reload();
} catch (Exception e) {
}
s.getAccountByName("someone");
}

```

Show moreShow more icon

上面的测试代码先查询了两次，然后 reload，然后再查询一次，结果应该是只有第一次查询走了数据库，其他两次查询都从缓存，第三次也走缓存因为 reload 失败了。

##### 清单 33\. 运行结果

```
real querying db...someone

```

Show moreShow more icon

和预期一样。那么我们如何避免这个问题呢？我们可以用 @CacheEvict 注释提供的 beforeInvocation 属性，将其设置为 true，这样，在方法执行前我们的缓存就被清空了。可以确保缓存被清空。

##### 清单 34\. AccountService.java

```
@CacheEvict(value="accountCache",allEntries=true,beforeInvocation=true)
// 清空 accountCache 缓存
public void reload() {
throw new RuntimeException();
}

```

Show moreShow more icon

注意上面的代码，我们在 @CacheEvict 注释中加了 beforeInvocation 属性，确保缓存被清空。

执行相同的测试代码

##### 清单 35\. 运行结果

```
real querying db...someone
real querying db...someone

```

Show moreShow more icon

这样，第一次和第三次都从数据库取数据了，缓存清空有效。

### 非 public 方法问题

和内部调用问题类似，非 public 方法如果想实现基于注释的缓存，必须采用基于 AspectJ 的 AOP 机制，这里限于篇幅不再细述。

## 其他技巧

### Dummy CacheManager 的配置和作用

有的时候，我们在代码迁移、调试或者部署的时候，恰好没有 cache 容器，比如 memcache 还不具备条件，h2db 还没有装好等，如果这个时候你想调试代码，岂不是要疯掉？这里有一个办法，在不具备缓存条件的时候，在不改代码的情况下，禁用缓存。

方法就是修改 spring\*.xml 配置文件，设置一个找不到缓存就不做任何操作的标志位，如下

##### 清单 36\. Spring-cache-anno.xml

```
<beans xmlns="http://www.springframework.org/schema/beans"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns:cache="http://www.springframework.org/schema/cache"
xmlns:p="http://www.springframework.org/schema/p"
xsi:schemaLocation="http://www.springframework.org/schema/beans
    http://www.springframework.org/schema/beans/spring-beans.xsd
     http://www.springframework.org/schema/cache
     http://www.springframework.org/schema/cache/spring-cache.xsd">

<cache:annotation-driven />

<bean id="accountServiceBean" class="cacheOfAnno.AccountService"/>

    <!-- generic cache manager -->
<bean id="simpleCacheManager"
class="org.springframework.cache.support.SimpleCacheManager">
     <property name="caches">
       <set>
         <bean
           class="org.springframework.cache.concurrent.ConcurrentMapCacheFactoryBean"
           p:name="default" />
       </set>
     </property>
</bean>

<!-- dummy cacheManager  -->
<bean id="cacheManager"
class="org.springframework.cache.support.CompositeCacheManager">
     <property name="cacheManagers">
       <list>
         <ref bean="simpleCacheManager" />
       </list>
     </property>
     <property name="fallbackToNoOpCache" value="true" />
</bean>

</beans>

```

Show moreShow more icon

注意以前的 cacheManager 变为了 simpleCacheManager，且没有配置 accountCache 实例，后面的 cacheManager 的实例是一个 CompositeCacheManager，他利用了前面的 simpleCacheManager 进行查询，如果查询不到，则根据标志位 fallbackToNoOpCache 来判断是否不做任何缓存操作。

##### 清单 37\. 运行结果

```
real querying db...someone
real querying db...someone
real querying db...someone

```

Show moreShow more icon

可以看出，缓存失效。每次都查询数据库。因为我们没有配置它需要的 accountCache 实例。

如果将上面 xml 配置文件的 fallbackToNoOpCache 设置为 false，再次运行，则会得到

##### 清单 38\. 运行结果

```
Exception in thread "main" java.lang.IllegalArgumentException:
Cannot find cache named [accountCache] for CacheableOperation
     [public cacheOfAnno.Account
     cacheOfAnno.AccountService.getAccountByName(java.lang.String)]
     caches=[accountCache] | condition='' | key=''

```

Show moreShow more icon

可见，在找不到 accountCache，且没有将 fallbackToNoOpCache 设置为 true 的情况下，系统会抛出异常。

## 结束语

总之，注释驱动的 spring cache 能够极大的减少我们编写常见缓存的代码量，通过少量的注释标签和配置文件，即可达到使代码具备缓存的能力。且具备很好的灵活性和扩展性。但是我们也应该看到，spring cache 由于急于 spring AOP 技术，尤其是动态的 proxy 技术，导致其不能很好的支持方法的内部调用或者非 public 方法的缓存设置，当然这都是可以解决的问题，通过学习这个技术，我们能够认识到，AOP 技术的应用还是很广泛的，如果有兴趣，我相信你也能基于 AOP 实现自己的缓存方案。