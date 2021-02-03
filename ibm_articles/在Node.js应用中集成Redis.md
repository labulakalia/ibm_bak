# 在 Node.js 应用中集成 Redis
了解 Redis

**标签:** Node.js

[原文链接](https://developer.ibm.com/zh/articles/os-cn-nodejs-redis/)

朱盛浩

发布: 2015-12-21

* * *

## Redis 简介

Redis 是一个基于内存的键（key）值（value）类型的数据结构存储容器，它既可以完全工作在内存中，也可以持久化存储。当 Redis 工作于持久化模式时，可以将它当作一个非关系型数据库使用。而工作于内存中时，则可以用作数据的高速缓存。不过和普通的键值结构缓存不同是：Redis 的值可以拥有种类繁多并且灵活的内建数据结构。这些数据结构具有及其高效的增删改查时间复杂度，在能够满足更多业务场景的数据存储需求同时还提供极为快速的处理速度。

### Redis 值数据结构类型

Redis 内建的值数据结构有如下类型：

- **字符串（Strings）**

    字符串是最基本也是最简单的值元素类型，普通的键值存储都可以归为此类。字符串类型的值如果是数字，还可以做一些数值运算，例如 incr 操作和 decr 操作对数值的原子增减。字符串类型还允许存储二进制数据，比如图片。

- - **哈希（Hashes）**



    存储比较复杂数据结构时使用。Redis 的哈希结构其实就是一个哈希表（hashmap），它允许在一个哈希中存储一个或多个键值映射，同时提供获取，修改，删除哈希表其中一个字段或者全部字段的方法。通常，我们可以将 JSON 结构的数据类型存储于哈希中。

- **线性列表（Lists）**

    常用的线性队列，实现方式是一个双向链表，它提供了列表头尾添加/删除，索引访问/添加/删除，获取队列长度等操作。可以用来实现消息队列，任务列表以及排行等功能。访问队列两端的元素是非常快的，但如果您试着访问一个非常大的列表的中间元素仍然是十分慢的，列表操作的时间复杂度为 O(N)。

- **集合（Sets）**

    集合对外提供的功能与线性列表类似是一个列表的功能，特殊之处在于集合是可以自动排重的，当您需要存储一个列表数据，又不希望出现重复数据时，集合是一个很好的选择，并且集合提供了判断某个成员是否在一个集合内的重要接口，这个也是列表所不能提供的。同时，Redis 还为集合提供了求交集、并集、差集等操作，非常方便、实用。集合是无序存储的，所有关于集合的增删改查操作都是常量 O(1) 时间复杂度。

- **有序集合（Sorted Sets）**

    有序集合与普通集合非常相似，是一个没有重复元素的字符串集合。不同之处是有序集合的没有成员都关联了一个评分，这个评分被用来按照从最低分到最高分的方式排序集合中的成员。集合的成员是唯一的，但是评分可以是重复的。使用有序集合您可以以非常快的速度（O(log(N))）添加，删除和更新元素。 您也可以很快的根据评分（score）或者次序（position）来获取一个范围的元素。访问有序集合的中间元素也是非常快的, 因此您能够使用有序集合作为一个没有重复成员的智能列表。在有序集合中，您可以很快捷的访问一切您需要的东西：有序的元素，快速的存在性测试，快速访问集合的中间元素。


### Redis 数据过期策略

在实际的生产应用中，服务器的内存都是有限的，而数据缓存的需求空间则可能是无限的。因此高速缓存都会有一定的机制将频繁访问的数据内容优先保存在缓存中，而将长时间不访问的数据从缓存中删除。Redis 提供了以下几种数据过期方式：

直接用 expire 命令设置键过期，这种方式要求开发者自己掌握键值对过期时刻，通常在有修改时使用。

LRU（近期最少使用）模式，当在启动 Redis 的配置中加入 maxmemory 选项时，Redis 就工作在 LRU 模式了。这种模式下，还有几种子类型可以配置：

allkeys-lru：所有的键值对按照最近最少使用原则被删除。

volatile-lru: 只有设置了 expire 的键值对会按最近最少原则删除。

allkeys-random: 随机删除某一个键值对。

volatile-random: 随机删除某一个设置了 expire 的键值对。

这些配置都是在 redis.conf 文件中配置的，也可以通过代码在创建 redis client 时设置。例如，我们设置 Redis 的最大内存池为 100MB，并且使用 LRU 模式中的 volatile-lru 子模式。

##### 清单 1\. 配置 redis.conf 使 Redis 工作在 LRU 模式

```
maxmemory 100mb maxmemory-policy
volatile-lru

```

Show moreShow more icon

## Redis 的 Node.js 插件

我们使用 Node.js 来进行开发，Redis 推荐使用的 Node.js 驱动为 Node\_Redis 插件。这是一个在 Github 上的开源项目，获取它的方式非常方便，可以直接使用 Node.js 的包管理 npm 获得。Redis 针对不同数据结构的操作有很多的命令，为了能够和 Redis 的命令保持一致，Node\_Redis 也采用了命令名作为调用方法名。例如我们要从 Redis 中获取一个键名为 key 的值，用 Redis 的命令如下：

```
GET 'key'

```

Show moreShow more icon

而利用 Node\_Redis，我们执行如下函数：

```
redis.get('key',function(res) {
// 业务逻辑
});

```

Show moreShow more icon

Node\_Redis 在 Github 上一直有专门维护人员且一直处于更新状态。更多的 Node\_Redis 介绍，可以参考其在 Github 的主页。

## 实例：Node.js 应用集成 Redis

### 一个简易的博客应用

接下来我们将手动搭建一个实例，我们将搭建一个简易的个人博客。假设我们有一个访问量较大的博客，为了能够减轻对数据库读的压力，我们将热门文章缓存在 Redis 中。同时，我们有博客的评论系统，读者在阅读完文章之后允许评论允许回复。这些评论和回复需要更新进 Redis，这样，这些评论和回复可以在较少的延迟内被别的读者所看到。

### 使用 Node\_Redis 插件

通过 Node.js 的 npm 命令我们可以很容易的将 Node\_Redis 插件下载到博客应用的 node\_modules 目录下。

##### 清单 2\. 通过命令行添加 Node\_Redis 插件

```
npm install redis --save

```

Show moreShow more icon

也可以通过在项目包依赖文件 package.json 中配置 Node\_Redis 的依赖项来导入：

##### 清单 3\. 通过管理文件添加 Node\_Redis 插件

```
"dependencies": { "redis" : "latest" }

```

Show moreShow more icon

安装完 Node\_Redis 插件之后，在 Node.js 程序中我们便可以直接通过 require 导入。接着我们需要配置 Node\_Redis 的连接参数来连接到 Redis-Server。假设我们的 Redis-Server 部署在 127.0.0.1，端口为 6379 ，那么通过下列配置我们便创建了一个 Redis 客户端。

##### 清单 4\. 创建一个 Redis 客户端并连接到 Redis-Server

```
var rediz = require('redis'); var redis =
rediz.createClient({ "host": "127.0.0.1", "port": "6379" });

```

Show moreShow more icon

通常我们会注册一个连接失败消息，当连接到 Redis-Server 出现错误的时候，错误信息会以回调方式返回，err 包含了许多有用信息来帮助定位和解决错误。

##### 清单 5\. 注册 Redis 错误事件监听

```
redis.on('error', function (err) { console.log('error
event - ' + redis.host + ':' + redis.port + ' - ' + err); });

```

Show moreShow more icon

### 操作异步化

Node\_Redis 提供的 API 都是同步调用，然后通过异步回调函数返回查询结果，例如以下的代码我们向 Redis 中存入一个值为字符串类型的键值对，并且读取存入的值以验证字符串已被正确存入：

##### 清单 6\. 同步模式调用 Node\_Redis

```
function redis_get_string() { redis.set('key',
['string'], function(err, res) { redis.get('key', function(err, res) {
console.log(res); //打印'string' }); }); };

```

Show moreShow more icon

因为 Node\_Redis 返回结果是异步的，也就是后续操作必须要在前一操作结果的回调函数中进行才能保证程序逻辑不出错。我们可以看到其中的弊端：业务逻辑必须写在回调中，业务逻辑一复杂那么回调层级也会很多，代码变得紧耦合，可读性变差，这就是所谓的”回调陷阱”。为了解决以上弊端，Node\_Redis 插件也提供了全异步的调用方式来操作 Redis，方法是引入 JavaScript 异步模块来配合，其中比较流行的异步模块为 bluebird。

##### 清单 7\. 引入 JavaScript 异步模块 bluebird

```
"dependencies" : { "redis" : "latest",
"bluebird" : "latest" }

```

Show moreShow more icon

同样的，当下载完 bluebird 后，通过 require 引入 bluebird 模块，同时对创建 Redis 客户端的代码稍加改造：

##### 清单 8\. 引入 bluebird 使 Node\_Redis API 异步化

```
var Q = require('bluebird'); var rediz =
require('redis'); Q.promisifyAll(rediz.RedisClient.prototype);
Q.promisifyAll(rediz.Multi.prototype);

```

Show moreShow more icon

通过这样的配置之后，原先的同步操作 API 依旧保留，在同步 API 方法名后追加 Async 即为新的异步 API。这样，Node\_Redis 便支持全异步模式来操作 Redis 了，异步 API 的返回值一律为一个 Promise 对象，用来延迟获取执行结果。同样的，使用上面的例子写入一个字符串并读取，异步代码可以改造成如下基于 Promise 的链式调用：

##### 清单 9\. 异步 Promise 链式调用

```
function redis_setget_string_async() { var promise =
redis.setAsync('key', ['string']).then(function(res) { return redis.getAsync('key');
//返回 Promise }).then(function(res) { console.log(res); //打印'string' Q.resolve(res);
//返回 Promise }); return promise; }

```

Show moreShow more icon

### ArticleHelper 类

每篇文章都应该有其唯一的 id，我们博客应用总体的设计思路是将文章 id 列表和文章主内容，文章评论组分离缓存。这是因为，文章主体内容和评论组可能十分庞大，会占用更多的内存空间，应当有一定的过期策略来使一定时间内未被访问的文章失效。而 id 列表则占用非常小的空间，可以长时间存储不失效。因此，我们的 ArticleHelper 类应该类似如下：

##### 清单 10\. ArticleHelper 类构造方法

```
var ArticleHelper = function () { this.ArticleIDSet
= "AIDSet"; // 存放文章 ID 的集合 this.ArticleIDPrefix = "Article"; // 文章 KEY 前缀
this.CommentIDPrefix = "Comment"; // 评论 KEY 前缀 };

```

Show moreShow more icon

Article 类构造函数中我们设置了 3 个成员变量：他们分别是文章 id 集合的键，文章主内容和评论组的键前缀。例如：一篇 id 为 100000 的文章在 Redis 中的键为 Article100000，评论组键为 Comment100000。

### 读取文章 ID 列表

设计数据库时，我们设置文章的 id 作为主键并自动递增。为了方便文章按新旧顺序显示，我们要在 Reids 中储存一份有序的文章 id 列表，而文章 id 大小可以直接做为文章的排序标准。在我们的博客服务启动时，读取数据库将所有的文章 id 缓存至 Redis，id 值越大，则文章越新。同时，我们应当使用一个 Redis 提供的有序集合来保存这些 id。

##### 清单 11\. ArticleHelper 有序存储文章 id 方法

```
ArticleHelper.prototype.getArticleIDs() { var
that = this; var promise = ArticleMySQLHelper.getIDs().then(function (artis) { if(
artis.length == 0 ) { return null; } var _ = []; artis.forEach(function(arti) {
_.push(redis.zaddAsync(that.ArticleIDSet, [1, arti.ID])); }); return Q.all(_); });
return promise; }

```

Show moreShow more icon

以上代码中我们用 zadd 方法添加 id 时，统一 score=1，这样有序集合就只会根据 id 大小来排序。

### 读取文章

建立了文章的 id 列表之后接下来我们应该考虑的事情就是读取一篇文章的全部内容并且缓存到 Redis 了。我们的博客应用应该提供 REST 服务，通过文章的 id 来访问文章。

##### 清单 12\. ArticleHelper 读取文章方法

```
ArticleHelper.prototype.getArticle = function(id)
{ var that = this; var promise =
redis.hgetallAsync(that.ArticleIDPrefix+id).then(function(res) { if (res == null) {
return ArticleMySQLHelper.getOne(id).then(function(article) { return
redis.hmsetAsync(that.ArticleIDPrefix+id, article). then(function(res) { return
Q.resolve(article); }); }); } else { return Q.resolve(res); } }); return promise; }

```

Show moreShow more icon

我们用 ArticleIDPrefix+id 的值作为该文章在 Redis 中的键。一篇文章应该有很多字段，例如作者，日期，内容等等。为了能在一个文章键中存储，我们选用了哈希表。Node\_Redis 支持将 JSON 结构的数据通过 hmset 命令直接存入 Redis，我们不需了解文章数据字段组成，更不需要遵循 hmset 多个字段多个值命令结构。

### 显示文章列表

我们的文章数量很多，在一个页面上肯定无法全部展示，因此读者在浏览网页时必须分页浏览我们的全部文章。我们的每一页文章数量都是固定的并且按文章新旧顺序已经排序，因此，读者在切换上下页时，我们将只需要读取文章 id 集合中某一段区间的文章 id，然后再根据这些 id 我们去读取数据库，最后将读取到的文章缓存在 Redis 中。

##### 清单 13\. ArticleHelper 读取文章列表方法

```
ArticleHelper.prototype.getArticleList =
function(offset, size) { var that = this; var promise =
redis.zrevrangeAsync([that.ArticleIDSet, offset, offset+size-1]).then(function(res){
var _ = []; res.forEach(function(id) { _.push(that.getArticle(id)); }); return
Q.all(_); }); return promise; };

```

Show moreShow more icon

Node\_Redis 提供 zrevrange() 方法从 Redis 中倒序读取有序集合中的一段区间 id。获取到了区间 id 组之后可以直接通过调用上面定义的 getArticle() 方法缓存文章。

### 文章评论

文章发布出去之后，和读者最大的互动方式就是评论了。首先读者应该能看到之前的老评论，这部分应该缓存在 Redis 中。当有新读者对文章进行评论时，我们需要将评论内容保存起来，比较好的方案是采用先写数据库，然后让 Redis 的老缓存过期再重读数据库。这样做至少有如下好处：符合读写分离，符合 Redis 作为缓存只读的定位。

通常情况下一篇文章可能有不止一条评论，并且可能有回复，一条评论和若干回复一起构成一个评论组。我们需要考虑下文章评论组在 Redis 中的存储结构，每个评论组应该有对应的文章，然后这些评论组都应该有日期先后顺序，较新的评论组应该靠前显示。对评论的回复也是类似，它将按照读者发言的先后顺序显示。根据以上的分析，一个有序线集合结构更适合存储这些评论组。

##### 清单 14\. 获取评论方法

```
ArticleHelper.prototype.getComment = function(id) { var that = this; var promise =
redis.zrangeAsync(that.CommentIDPrefix+id, [0, -1]).then(function(res) { if (res ==
null) { var p = ArticleMySQLHelper.getComments(id).then(function(cmmts) { var _ =
[]; cmmts.forEach(function(cmmt) { _.push(redis.zaddAsync(that.CommentIDPrefix+id,
[cmmt.index, cmmt.body])); }); return Q.all(_); }); return p; } else { return
Q.resolve(res); } }) return promise; }

```

Show moreShow more icon

一个评论组将共享一个有序集合键，我们将评论组的第一条也就是评论在有序集合中设置 score=1，而接下来的回复 score 依次递增，这样，评论组在有序集合中就可以顺序读取出来正常显示了。当有一条新的评论或者回复时，我们需要将其写入关系型数据库，同时设置该评论 Redis 缓存过期。

##### 清单 15\. 用户添加评论方法

```
ArticleHelper.prototype.addComment = function(id, comment) { var that = this; var
promise = ArticleMySQLHelper.saveComment(id, comment).then(function(res){ return
redis.expireAsync(that.ArticleIDPrefix+id, [0]); }); return promise; };

```

Show moreShow more icon

执行 expire 命令之后，Redis 就会将这对评论键值对删除。更新后评论组将重新通过 getComment() 方法获取。

### 文章/评论过期

随着存储的文章数量不断上升并且接近 Redis 设定的内存上限时，我们就需要考虑下一定的策略将部分文章清理出缓存了。在我们的缓存中，存在文章 id 集合，文章集合，评论组集合，其中文章 id 集合不应该被过期销毁，而文章集合和评论组集合都可以被过期删除。因此，我们设置 Redis 的工作模式为 LRU，子工作模式为 volatile-lru，这种模式需要将所有文章和评论组键值对都设置过期时间，而过期时间的颗粒度可以自由选择，但是不建议太短，因为设置太短的过期时间加重了数据库的读压力，这里我们选择 2 小时过期时间。因此，我们需要将 getArticle 方法做如下修改：

##### 清单 16\. 带过期设置的读取文章方法

```
ArticleHelper.prototype.getArticleExpire = function(id) {
var that = this; var promise =
redis.hgetallAsync(that.ArticleIDPrefix+id).then(function(res) { if (res == null) {
return ArticleMySQLHelper.getOne(id).then(function(article) { return
redis.hmsetAsync(that.ArticleIDPrefix+id, article).then(function(res) { return
redis.expireAsync(that.ArticleIDPrefix+id, 3600*2); }).then(function(){ return
Q.resolve(article); }); }); } else { return Q.resolve(res); } }); return promise; }

```

Show moreShow more icon

设置了过期时间之后，无论是否到达过期时间，只要出现内存不足 volatile-lru 机制都会将带有 expire set 的缓存给删除。

## 结束语

到这里一个比较简单的博客系统做好了。我们可以看到 Redis 在其中扮演的角色，在数据库和应用服务之间假设了一道缓存层，减少应用直接与数据库交互，加大吞吐量。随着我们应用规模进一步扩大，我们可以预期单机部署的 Redis 也将会达到性能极限。为了解决这些问题，多机集群，主从备份，容灾等更多高级的主题将会被引入，在这些高级主题上，Node.js 依旧能够和 Redis 良好结合。有兴趣的读者，可以进一步研究。