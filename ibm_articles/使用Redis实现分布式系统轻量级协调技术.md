# 使用 Redis 实现分布式系统轻量级协调技术
一种基于共享内存模式的分布式系统轻量级协调技术方案

**标签:** Python,数据库,消息传递

[原文链接](https://developer.ibm.com/zh/articles/os-cn-redis-coordinate/)

朱 培旭

发布: 2015-06-16

* * *

在分布式系统中，各个进程（本文使用进程来描述分布式系统中的运行主体，它们可以在同一个物理节点上也可以在不同的物理节点上）相互之间通常是需要协调进行运作的，有时是不同进程所处理的数据有依赖关系，必须按照一定的次序进行处理，有时是在一些特定的时间需要某个进程处理某些事务等等，人们通常会使用分布式锁、选举算法等技术来协调各个进程之间的行为。因为分布式系统本身的复杂特性，以及对于容错性的要求，这些技术通常是重量级的，比如 Paxos 算法，欺负选举算法，ZooKeeper 等，侧重于消息的通信而不是共享内存，通常也是出了名的复杂和难以理解，当在具体的实现和实施中遇到问题时都是一个挑战。

Redis 经常被人们认为是一种 NoSQL 软件，但其本质上是一种分布式的数据结构服务器软件，提供了一个分布式的基于内存的数据结构存储服务。在实现上，仅使用一个线程来处理具体的内存数据结构，保证它的数据操作命令的原子特性；它同时还支持基于 Lua 的脚本，每个 Redis 实例使用同一个 Lua 解释器来解释运行 Lua 脚本，从而 Lua 脚本也具备了原子特性，这种原子操作的特性使得基于共享内存模式的分布式系统的协调方式成了可能，而且具备了很大的吸引力，和复杂的基于消息的机制不同，基于共享内存的模式对于很多技术人员来说明显容易理解的多，特别是那些已经了解多线程或多进程技术的人。在具体实践中，也并不是所有的分布式系统都像分布式数据库系统那样需要严格的模型的，而所使用的技术也不一定全部需要有坚实的理论基础和数学证明，这就使得基于 Redis 来实现分布式系统的协调技术具备了一定的实用价值，实际上，人们也已经进行了不少尝试。本文就其中的一些协调技术进行介绍。

## signal/wait 操作

在分布式系统中，有些进程需要等待其它进程的状态的改变，或者通知其它进程自己的状态的改变，比如，进程之间有操作上的依赖次序时，就有进程需要等待，有进程需要发射信号通知等待的进程进行后续的操作，这些工作可以通过 Redis 的 Pub/Sub 系列命令来完成，比如：

```
import redis, time
rc = redis.Redis()

def wait( wait_for ):
ps = rc.pubsub()
ps.subscribe( wait_for )
ps.get_message()
wait_msg = None
while True:
msg = ps.get_message()
if msg and msg['type'] == 'message':
wait_msg = msg
break
time.sleep(0.001)
ps.close()
return wait_msg

def signal_broadcast( wait_in, data ):
wait_count = rc.publish(wait_in, data)
return wait_count

```

Show moreShow more icon

用这个方法很容易进行扩展实现其它的等待策略，比如 try wait，wait 超时，wait 多个信号时是要等待全部信号还是任意一个信号到达即可返回等等。因为 Redis 本身支持基于模式匹配的消息订阅（使用 psubscribe 命令），设置 wait 信号时也可以通过模式匹配的方式进行。

和其它的数据操作不同，订阅消息是即时易逝的，不在内存中保存，不进行持久化保存，如果客户端到服务端的连接断开的话也是不会重发的，但是在配置了 master/slave 节点的情况下，会把 publish 命令同步到 slave 节点上，这样我们就可以同时在 master 以及 slave 节点的连接上订阅某个频道，从而可以同时接收到发布者发布的消息，即使 master 在使用过程中出故障，或者到 master 的连接出了故障，我们仍然能够从 slave 节点获得订阅的消息，从而获得更好的鲁棒性。另外，因为数据不用写入磁盘，这种方法在性能上也是有优势的。

上面的方法中信号是广播的，所有在 wait 的进程都会收到信号，如果要将信号设置成单播，只允许其中一个收到信号，则可以通过约定频道名称模式的方式来实现，比如：

频道名称 = 频道名前缀 (channel) + 订阅者全局唯一 ID(myid)

其中唯一 ID 可以是 UUID，也可以是一个随机数字符串，确保全局唯一即可。在发送 signal 之前先使用 “pubsub channels channel\* ”命令获得所有的订阅者订阅的频道，然后发送信号给其中一个随机指定的频道；等待的时候需要传递自己的唯一 ID，将频道名前缀和唯一 ID 合并为一个频道名称，然后同前面例子一样进行 wait。示例如下：

```
import random

single_cast_script="""
local channels = redis.call('pubsub', 'channels', ARGV[1]..'*');
if #channels == 0 then
return 0;
end;
local index= math.mod(math.floor(tonumber(ARGV[2])), #channels) + 1;
return redis.call( 'publish', channels[index], ARGV[3]);
"""

def wait_single( channel, myid):
return wait( channel + myid )

def signal_single( channel, data):
rand_num = int(random.random() * 65535)
return rc.eval( single_cast_script, 0, channel, str(rand_num), str(data) )

```

Show moreShow more icon

## 分布式锁 Distributed Locks

分布式锁的实现是人们探索的比较多的一个方向，在 Redis 的官方网站上专门有一篇文档介绍 [基于 Redis 的分布式锁](http://redis.io/topics/distlock) ，其中提出了 Redlock 算法，并列出了多种语言的实现案例，这里作一简要介绍。

Redlock 算法着眼于满足分布式锁的三个要素：

- 安全性：保证互斥，任何时间至多只有一个客户端可以持有锁。
- 免死锁：即使当前持有锁的客户端崩溃或者从集群中被分开了，其它客户端最终总是能够获得锁。
- 容错性：只要大部分的 Redis 节点在线，那么客户端就能够获取和释放锁。

锁的一个简单直接的实现方法就是用 SET NX 命令设置一个设定了存活周期 TTL 的 Key 来获取锁，通过删除 Key 来释放锁，通过存活周期来保证避免死锁。不过这个方法存在单点故障风险，如果部署了 master/slave 节点，则在特定条件下可能会导致安全性方面的冲突，比如：

1. 客户端 A 从 master 节点获得锁。
2. master 节点在将 key 复制到 slave 节点之前崩溃了。
3. slave 节点提升为新的 master 节点。
4. 客户端 B 从新的 master 节点获得了锁，而这个锁实际上已经由客户端 A 所持有，导致了系统中有两个客户端在同一时间段内持有同一个互斥锁，破坏了互斥锁的安全性。

在 Redlock 算法中，通过类似于下面这样的命令进行加锁：

```
SET resource_name my_random_value NX PX 30000

```

Show moreShow more icon

这里的 my\_random\_value 为全局不同的随机数，每个客户端需要自己产生这个随机数并且记住它，后面解锁的时候需要用到它。

解锁则需要通过一个 Lua 脚本来执行，不能简单地直接删除 Key，否则可能会把别人持有的锁给释放了：

```
if redis.call("get",KEYS[1]) == ARGV[1] then
return redis.call("del",KEYS[1])
else
return 0
end

```

Show moreShow more icon

这个 ARGV[1] 的值就是前面加锁的时候的 my\_random\_value 的值。

如果需要更好的容错性，可以建立一个有 N（N 为奇数）个相互独立完备的 Redis 冗余节点的集群，这种情况下，一个客户端获得锁和释放锁的算法如下：

1. 先获取当前时间戳 timestamp\_1，以毫秒为单位。
2. 以相同的 Key 和随机数值，依次从 N 个节点获取锁，每次获取锁都设置一个超时，超时时限要保证小于所有节点上该锁的自动释放时间，以免在某个节点上耗时过长，通常都设的比较短。
3. 客户端将当前时间戳减去第一步中的时间戳 timestamp\_1，计算获取锁总消耗时间。只有当客户端获得了半数以上节点的锁，而且总耗时少于锁存活时间，该客户端才被认为已经成功获得了锁。
4. 如果获得了锁，则其存活时间为开始预设锁存活时间减去获取锁总耗时间。
5. 如果客户端不能获得锁，则应该马上在所有节点上解锁。
6. 如果要重试，则在随机延时之后重新去获取锁。
7. 获得了锁的客户端要释放锁，简单地在所有节点上解锁即可。

Redlock 算法不需要保证 Redis 节点之间的时钟是同步的（不论是物理时钟还是逻辑时钟），这点和传统的一些基于同步时钟的分布式锁算法有所不同。Redlock 算法的具体的细节可以参阅 Redis 的官方文档，以及文档中列出的多种语言版本的实现。

## 选举算法

在分布式系统中，经常会有些事务是需要在某个时间段内由一个进程来完成，或者由一个进程作为 leader 来协调其它的进程，这个时候就需要用到选举算法，传统的选举算法有欺负选举算法（霸道选举算法）、环选举算法、Paxos 算法、Zab 算法 (ZooKeeper) 等，这些算法有些依赖于消息的可靠传递以及时钟同步，有些过于复杂，难以实现和验证。新的 Raft 算法相比较其它算法来说已经容易了很多，不过它仍然需要依赖心跳广播和逻辑时钟，leader 需要不断地向 follower 广播消息来维持从属关系，节点扩展时也需要其它算法配合。

选举算法和分布式锁有点类似，任意时刻最多只能有一个 leader 资源。当然，我们也可以用前面描述的分布式锁来实现，设置一个 leader 资源，获得这个资源锁的为 leader，锁的生命周期过了之后，再重新竞争这个资源锁。这是一种竞争性的算法，这个方法会导致有比较多的空档期内没有 leader 的情况，也不好实现 leader 的连任，而 leader 的连任是有比较大的好处的，比如 leader 执行任务可以比较准时一些，查看日志以及排查问题的时候也方便很多，如果我们需要一个算法实现 leader 可以连任，那么可以采用这样的方法：

```
import redis

rc = redis.Redis()
local_selector = 0

def master():
global local_selector
master_selector = rc.incr('master_selector')
if master_selector == 1:
# initial / restarted
local_selector = master_selector
else:
if local_selector > 0:
# I'm the master before
if local_selector > master_selector:
# lost, maybe the db is fail-overed.
local_selector = 0
else:
# continue to be the master
local_selector = master_selector
if local_selector > 0:
# I'm the current master
rc.expire('master_selector', 20)
return local_selector > 0

```

Show moreShow more icon

这个算法鼓励连任，只有当前的 leader 发生故障或者执行某个任务所耗时间超过了任期、或者 Redis 节点发生故障恢复之后才需要重新选举出新的 leader。在 master/slave 模式下，如果 master 节点发生故障，某个 slave 节点提升为新的 master 节点，即使当时 master\_selector 值尚未能同步成功，也不会导致出现两个 leader 的情况。如果某个 leader 一直连任，则 master\_selector 的值会一直递增下去，考虑到 master\_selector 是一个 64 位的整型类型，在可预见的时间内是不可能溢出的，加上每次进行 leader 更换的时候 master\_selector 会重置为从 1 开始，这种递增的方式是可以接受的，但是碰到 Redis 客户端（比如 Node.js）不支持 64 位整型类型的时候就需要针对这种情况作处理。如果当前 leader 进程处理时间超过了任期，则其它进程可以重新生成新的 leader 进程，老的 leader 进程处理完毕事务后，如果新的 leader 的进程经历的任期次数超过或等于老的 leader 进程的任期次数，则可能会出现两个 leader 进程，为了避免这种情况，每个 leader 进程在处理完任期事务之后都应该检查一下自己的处理时间是否超过了任期，如果超过了任期，则应当先设置 local\_selector 为 0 之后再调用 master 检查自己是否是 leader 进程。

## 消息队列

消息队列是分布式系统之间的通信基本设施，通过消息可以构造复杂的进程间的协调操作和互操作。Redis 也提供了构造消息队列的原语，比如 Pub/Sub 系列命令，就提供了基于订阅/发布模式的消息收发方法，但是 Pub/Sub 消息并不在 Redis 内保持，从而也就没有进行持久化，适用于所传输的消息即使丢失了也没有关系的场景。

如果要考虑到持久化，则可以考虑 list 系列操作命令，用 PUSH 系列命令（LPUSH, RPUSH 等）推送消息到某个 list，用 POP 系列命令（LPOP, RPOP,BLPOP,BRPOP 等）获取某个 list 上的消息，通过不同的组合方式可以得到 FIFO，FILO，比如：

```
import redis
rc = redis.Redis()

def fifo_push(q, data):
rc.lpush(q, data)

def fifo_pop(q):
return rc.rpop(q)

def filo_push(q, data):
rc.lpush(q, data)

def filo_pop(q):
return rc.lpop(q)

```

Show moreShow more icon

如果用 BLPOP,BRPOP 命令替代 LPOP, RPOP，则在 list 为空的时候还支持阻塞等待。不过，即使按照这种方式实现了持久化，如果在 POP 消息返回的时候网络故障，则依然会发生消息丢失，针对这种需求 Redis 提供了 RPOPLPUSH 和 BRPOPLPUSH 命令来先将提取的消息保存在另外一个 list 中，客户端可以先从这个 list 查看和处理消息数据，处理完毕之后再从这个 list 中删除消息数据，从而确保了消息不会丢失，示例如下：

```
def safe_fifo_push(q, data):
rc.lpush(q, data)

def safe_fifo_pop(q, cache):
msg = rc.rpoplpush(q, cache)
# check and do something on msg
rc.lrem(cache, 1) # remove the msg in cache list.
return msg

```

Show moreShow more icon

如果使用 BRPOPLPUSH 命令替代 RPOPLPUSH 命令，则可以在 q 为空的时候阻塞等待。

## 结束语

使用 Redis 作为分布式系统的共享内存，以共享内存模式为基础来实现分布式系统协调技术，虽然不像传统的基于消息传递的技术那样有着坚实的理论证明的基础，但是它在一些要求不苛刻的情况下不失为一种简单实用的轻量级解决方案，毕竟不是每个系统都需要严格的容错性等要求，也不是每个系统都会频繁地发生进程异常，而且 Redis 本身已经经受了工业界的多年实践和考验。另外，用 Redis 技术还有一些额外的好处，比如在开发过程中和生产环境中都可以直接观察到锁、队列的内容，实施的时候也不需要额外的特别配置过程等，它足够简单，在调试问题的时候逻辑清晰，进行排查和临时干预也比较方便。在可扩展性方面也比较好，可以动态扩展分布式系统的进程数目，而不需要事先预定好进程数目。

Redis 支持基于 Key 值 hash 的集群，在集群中应用本文所述技术时建议另外部署专用 Redis 节点（或者冗余 Redis 节点集群）来使用，因为在基于 Key 值 hash 的集群中，不同的 Key 值会根据 hash 值被分布到不同的集群节点上，而且对于 Lua 脚本的支持也受到限制，难以保证一些操作的原子性，这一点是需要考虑到的。使用专用节点还有一个好处是专用节点的数据量会少很多，当应用了 master/slave 部署或者 AOF 模式的时候，因为数据量少，master 和 slave 之间的同步会少很多，AOF 模式实时写入磁盘的数据也少很多，这样子也可以大大提高可用性。

本文示例所列 Python 代码在 Python3.4 下运行，Redis 客户端采用 [redis 2.10.3](https://pypi.python.org/pypi/redis/2.10.3) ，Redis 服务端版本为 3.0.1 版。