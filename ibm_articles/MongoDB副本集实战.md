# MongoDB 副本集实战
MongoDB 高可用的那些事儿

**标签:** 数据库

[原文链接](https://developer.ibm.com/zh/articles/os-mongodb-replication/)

王 拓

发布: 2019-08-19

* * *

## 背景

对于互联网应用来说，我们要尽量保证服务的不可间断性，一旦出现故障需要尽快的恢复。MongoDB 的副本集模式主要解决了在其主从模式下故障无法自动转移的情况，自动实现高可用。举个例子，如果使用主从模式，一旦主库出现故障，首先需要人为去修改程序连接 MongoDB 的一些相关配置信息，还需登录到从库注释掉和从库相关配置才能使得 MongoDB 对程序提供正常读写服务。如果使用了副本集，那么主库挂掉了，自动会有别的节点来承担主库的职责，此过程无需人为干预。这一点给了我们充分的理由去使用 MongoDB 副本集。另外，MongoDB 从 4.0 版本开始，已经删除了对主从模式的支持。3.0 之后的版本中，每个 replica set 架构模式可以支持最多 50 个节点，先前版本只能支持 12 个。本文将介绍如何搭建、配置和使用 MongoDB 的副本集。

**[Databases for MongoDB](https://cloud.ibm.com/catalog/services/databases-for-mongodb?cm_sp=ibmdev-_-developer-articles-_-cloudreg)**

试用 [IBM Cloud](https://cloud.ibm.com/login?cm_sp=ibmdev-_-developer-articles-_-cloudreg) 上提供的 MongoDB 数据库服务。

**[Hyper Protect DBaaS for MongoDB](https://cloud.ibm.com/catalog/services/hyper-protect-dbaas-for-mongodb?cm_sp=ibmdev-_-developer-articles-_-cloudreg)**

试用 [IBM Cloud](https://cloud.ibm.com/login?cm_sp=ibmdev-_-developer-articles-_-cloudreg) 上提供的更加安全的 MongoDB 企业服务，您可以通过标准化的界面管理 MongoDB。

## 副本集搭建

首先我们来搭建一个副本集看看它长什么样子。搭建环境为：

- OS: CentOS 7 64 位。
- Software: 采用的 `mongodb-linux-x86_64-4.0.10` 二进制包。
- 准备三台虚拟机：mg1(primary)->10.0.4.23，mg2(secondary)->10.0.4.64，mg3(secondary)->10.0.4.33。

### 快速搭建 3 个 MongoDB 实例

以下命令在三个节点 (mg1, mg2, mg3) 分别执行。

1. 禁用 THP（数据库应用对内存的访问一般都是稀疏访问模式而很少是连续访问模式,如果启用大页面会导致更多的磁盘 IO）。





    ```
    echo never > /sys/kernel/mm/transparent_hugepage/enabled
    echo never > /sys/kernel/mm/transparent_hugepage/defrag

    ```





    Show moreShow more icon

2. 下载二进制安装包。





    ```
    wget https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-4.0.10.tgz
    && tar xvf mongodb-linux-x86_64-4.0.10.tgz -C /usr/local/ && cd
    /usr/local/ && ln -s mongodb-linux-x86_64-4.0.10 mongodb

    ```





    Show moreShow more icon

3. 创建配置文件。





    ```
    cd /usr/local/mongodb && mkdir config && cd config
    cat<<"EOF">mongodb-27017.conf
    bind_ip = 0.0.0.0
    #bind_ip_all = true # add in 3.6 和 bind_ip 互斥
    port = 27017
    dbpath = /var/lib/mongodb/27017 # 数据目录
    logpath = /var/log/mongodb_27017.log # 日志目录
    logappend = true
    fork = true
    auth = false
    #keyFile = # 副本集启用认证的配置
    directoryperdb = true
    storageEngine = wiredTiger
    #wiredTigerCacheSizeGB = 1
    #profile = 2
    #slowms = 1000
    #oplogSize = 100
    #replSet = test # 副本集名称
    EOF

    ```





    Show moreShow more icon

4. 创建数据目录。





    ```
    mkdir -pv /var/lib/mongodb/27017

    ```





    Show moreShow more icon

5. 启动 mongodb。





    ```
    numactl --interleave=all /usr/local/mongodb/bin/mongod -f /usr/local/mongodb/config/mongodb-27017.conf

    ```





    Show moreShow more icon

6. 确认启动，通过 `ps` 命令可以看到以下进程，说明成功启动，如果启动失败可以通过 `/var/log/mongodb_27017.log` 日志文件观察具体失败原因。

    以下为命令和输出示例：





    ```
    ps -ef | grep mongodb
    root 5137 1 4 09:28 ? 00:00:00 /usr/local/mongodb/bin/mongod -f
    /usr/local/mongodb/config/mongodb-27017.conf

    ```





    Show moreShow more icon

7. 登录。





    ```
    /usr/local/mongodb/bin/mongo --host 127.0.0.1 --port 27017

    ```





    Show moreShow more icon


通过以上步骤我们就成功的启动了 3 个 MongoDB 实例。

### 配置副本集

副本集配置比较简单，主要关注两个参数 `replSet`（每个节点必须配置一致）和 `keyFile` （认证相关）。

1. 任意一台服务器生成 `keyfile` 文件，用于副本集认证，然后传送到每台服务器。





    ```
    cd /usr/local/mongodb/config
    openssl rand -base64 741 > mongodb-keyfile

    ```





    Show moreShow more icon

2. 注意修改权限，否则会报错。





    ```
    chmod 600 mongodb-keyfile
    权限未修改为 600 的报错如下输出：
    ACCESS   permissions on xxx are too open

    ```





    Show moreShow more icon

3. 编辑每个实例的配置文件，修改 `replSet` 和 `keyfile`。





    ```
    vi /usr/local/mongodb/config/mongodb-27017.conf
    replSet = myReplSet
    keyFile = /usr/local/mongodb/config/mongodb-keyfile

    ```





    Show moreShow more icon

4. 重启 `mongodb`。





    ```
    /usr/local/mongodb/bin/mongod -f /usr/local/mongodb/config/mongodb-27017.conf --shutdown
    numactl --interleave=all /usr/local/mongodb/bin/mongod -f /usr/local/mongodb/config/mongodb-27017.conf

    ```





    Show moreShow more icon

5. `primary` 执行如下命令:





    ```
    config = {_id: "myReplSet", members: []} # _id 为副本集的名称
    config.members.push({_id: 0, host: "10.0.4.23:27017", "priority": 1})
    config.members.push({_id: 1, host: "10.0.4.64:27017", "priority": 1})
    config.members.push({_id: 2, host: "10.0.4.33:27017", "priority": 1})
    rs.initiate(config)

    ```





    Show moreShow more icon

    关键参数解释:

    - `"priority"`: 优先级，默认为 `1` ，优先级 `0` 为被动节点，不能成为活跃节点。优先级不为 `0` 则按照由大到小选出活跃节点。
    - `"arbiterOnly"`: 仲裁节点，只参与投票，不接收数据，也不能成为活跃节点。如果节点变为主库，则 `mongo shell` 会变成 `myReplSet:PRIMARY>` 。
6. 创建相关账号 `3.x` 或 `4.x` 需要提前创建账号才能使用 `rs.add()` 。





    ```
    use admin
    db.createUser({user:"root",pwd:"root",roles:[{role:"root",db:"admin"}]});

    ```





    Show moreShow more icon

7. 增加节点->primary 执行(先验证账号密码)。





    ```
    use admin
    db.auth('root','root')

    ```





    Show moreShow more icon

8. 增加第一个 `secondary->primary` 执行。

    `rs.add("10.0.4.64:27017")`

    以下为命令输出:





    ```
    {
        "ok" : 1,
        "operationTime" : Timestamp(1562292397, 1),
        "$clusterTime" : {
            "clusterTime" : Timestamp(1562292397, 1),
            "signature" : {
                "hash" : BinData(0,"Jkf3msZaFQu+XDOxbayQvL/3dyM="),
                "keyId" : NumberLong("6709993403284717569")
            }
        }
    }

    ```





    Show moreShow more icon

9. 增加第二个 `secondary->primary` 执行。

    `rs.add("10.0.4.33:27017")`

    以下为命令输出:





    ```
    {
        "ok" : 1,
        "operationTime" : Timestamp(1562292418, 1),
        "$clusterTime" : {
            "clusterTime" : Timestamp(1562292418, 1),
            "signature" : {
                "hash" : BinData(0,"y18FzdXn3/pcHJ4gyGlq4wYfR1I="),
                "keyId" : NumberLong("6709993403284717569")
            }
        }
    }

    ```





    Show moreShow more icon

10. 确认副本集状态。

    `rs.status()`

    以下为部分关键命令输出:





    ```
    {
        "set" : "myReplSet",
        "members" : [
            {
                "_id" : 0,
                "name" : "10.0.4.23:27017",
                "health" : 1,
                "state" : 1,
                "stateStr" : "PRIMARY",
            },
            {
                "_id" : 1,
                "name" : "10.0.4.64:27017",
                "health" : 1,
                "state" : 2,
                "stateStr" : "SECONDARY",
            },
            {
                "_id" : 2,
                "name" : "10.0.4.33:27017",
                "health" : 1,
                "state" : 2,
                "stateStr" : "SECONDARY",
            }
        ]
    }

    ```





    Show moreShow more icon


至此为止一个三节点的副本集已经搭建完成。

## 副本集如何选举主库

主库（Primary）主要对程序提供读写服务，从库（Secondary）在主库挂掉的时候接替其成为主库，以保证程序的读写服务不间断。

MongoDB r3.2.0 版本之前选举协议是基于 Bully 算法，从 r3.2.0 版本开始默认使用基于 [Raft 算法](https://blog.csdn.net/xu__cg/article/details/73555161) 的选举策略。

### 影响选举的因素和条件

以下因素都会决定一个从库能否成为主库。

- **心跳** ：副本集成员每 2 秒发送一次心跳（pings），如果心跳在 10 秒内没有返回，则其它成员将其标记为不可访问。
- **优先级** ：其它成员更愿投票给 priority 值高的。

    - Priority 为 0 的不能成为 primary 也不会寻求选举
    - 只要当前主节点具有最高优先级值，或者没有具有较高优先级的辅助节点在该集合中的最新 oplog（操作日志）条目的 10 秒内，副本集就不会保持选举。
    - 如果优先级较高的成员在当前主节点的最新 oplog 条目的 10 秒内赶上，则该集合保持选举，以便为优先级较高的节点提供成为主节点的机会。
- **Optime** ：时间戳比较老的不能成为 primary。
- **Connections** ：如果大多数副本集不可访问或不可用，则所有剩余成员变为只读。何为大多数？例如 3 个副本集中 2 个为大多数，5 个副本集中有 3 个为大多数。选举需要时间，在选举的过程中，副本集由于没有 primary，不能接受写入，所有剩余成员都变为只读。

### 选举过程

整个选举过程其实非常快，对用户来说几乎是透明的。以下步骤为具体选举步骤。

1. 得到每个服务器节点的最后操作时间戳。每个 mongodb 都有 oplog 机制会记录本机的操作，方便和主服务器进行对比数据是否同步还可以用于错误恢复。
2. 如果集群中大部分服务器 down 机了，保留活着的节点都为 secondary 状态并停止，不选举了。
3. 如果集群中选举出来的主节点或者所有从节点最后一次同步时间看起来很旧了，停止选举等待人来操作。
4. 如果上面都没有问题就选择最后操作时间戳最新（保证数据是最新的）的服务器节点作为主节点。

### 选举触发条件

正常情况下不会触发选举，如有以下条件之一则会出触发选举。

1. 初始化一个副本集时。
2. 副本集和主节点断开连接，可能是网络问题。
3. 主节点挂掉。

### 如何降级为 Standalone 模式

有些场景例如核心业务从 MongoDB 下架等，则可以从副本集模式降级为 Standalone 模式。

1. 注释掉 `replicaSet` 参数。
2. 删除数据库目录下的 `local.*` 文件（注意删除操作必须在 MongoDB 停止服务的情况下进行）。

_**注意** ：当切换回 Standalone 模式下时所做的任何事情都是没有 oplog 的，这意味着当再次回到 Replica Set 模式中时会丢失部分数据（直接使用 SECONDARY 未经过删除 `local.*` 文件的原配置启动），此时MongoDB 会报错，此处是个坑。所以这个时候得重新搭建副本集。_

### 将 Standalone 模式转化为副本集模式

项目前期对于 MongoDB 的要求没有那么高，随着业务的增长，发现项目对 MongoDB 的可靠性要求越来越高，此时可以将 MongoDB 切换为副本集模式。

1. 关闭 mongodb 实例。
2. 配置 replSet 参数以及生成配置 keyFile 。
3. 连接到 mongodb。
4. 初始化，增加成员。

具体可以参考： [Convert a Standalone to a Replica Set](https://docs.mongodb.com/v3.6/tutorial/convert-standalone-to-replica-set/) 。

## 模拟副本集中单实例故障

现在我们有了一个三节点成员的副本集，现在我们来模拟主库宕机，看看是否会自动切换。

1. 在 `PRIMARY` 节点将 `mongodb` 关闭。





    ```
    /usr/local/mongodb/bin/mongod -f /usr/local/mongodb/config/mongodb-27017.conf --shutdown

    ```





    Show moreShow more icon

2. 在任意一台 `SECONDARY` 节点登录查看状态，可以看到主库已经由 mg1 `（10.0.4.23）` 切换为 mg2 `（10.0.4.64）` 。

    以下为命令和输出示例：





    ```
    myReplSet:PRIMARY> use admin
    switched to db admin
    myReplSet:PRIMARY> db.auth('root','root')
    1
    myReplSet:PRIMARY> rs.status()
    {
        "set" : "myReplSet",
        "members" : [
            {
                "_id" : 0,
                "name" : "10.0.4.23:27017",
                "health" : 0,
                "state" : 8,
                "stateStr" : "(not reachable/healthy)",
            },
            {
                "_id" : 1,
                "name" : "10.0.4.64:27017",
                "health" : 1,
                "state" : 1,
                "stateStr" : "PRIMARY",
            },
            {
                "_id" : 2,
                "name" : "10.0.4.33:27017",
                "health" : 1,
                "state" : 2,
                "stateStr" : "SECONDARY",
            }
        ],
        "ok" : 1,
        "operationTime" : Timestamp(1562293557, 1),
        "$clusterTime" : {
            "clusterTime" : Timestamp(1562293557, 1),
            "signature" : {
                "hash" : BinData(0,"KXVbe3bUu27xEwgNmv6Y628m9mY="),
                "keyId" : NumberLong("6709993403284717569")
            }
        }
    }

    ```





    Show moreShow more icon


上述只是个简单的模拟主库宕机，可以看到在主库宕机的这段时间到重新选举主库花费的时间几乎感知不到。

## 副本集数据同步原理

Oplog（操作日志）是一种特殊的 Capped collections（固定集合），特殊之处在于它是系统级 Collection，记录了数据库的所有操作，集群之间依靠 oplog 进行数据同步。由于 local 数据不允许创建用户，如果要访问 oplog 需要借助其它数据库的用户，并且赋予该用户访问 local 数据库的权限。

Oplog 记录的操作记录是幂等的 (idempotent)，这意味着你可以多次执行这些操作而不会导致数据丢失或不一致。

Oplog 位于 local 库下：

- `master/slave->db.oplog.$main.find()`
- `replicaset->db.oplog.rs.find()`
- `sharding->` 需要切换到每个分片去查看

Oplog 大小：

- 对于 64 位 Linux，Solaris，FreeBSD 和 Windows 系统，MongoDB 分配 5％的可用磁盘空间，但总是至少分配 1G，不超过 50G。
- 对于 64 位 OS X 系统，MongoDB 为 oplog 分配 183MB 的空间。
- 对于 32 位系统，MongoDB 为 oplog 分配大约 48MB 的空间。

Oplog 的大小设置是值得考虑的一个问题，如果 oplog 过大，会浪费存储空间；如果 Oplog 过小，老的 oplog 记录很快就会被覆盖，那么宕机的节点就很容易出现无法同步数据的现象。

下面我们通过生产环境的一个示例来详细了解 oplog 如何存储在 MongoDB 中。登录 Primary 执行以下命令：

```
use local;
db.oplog.rs.find().pretty();

```

Show moreShow more icon

以下为命令输出示例及部分关键字段解释：

- `"ts" : Timestamp(154806363, 19)`

    8 字节的时间戳，由 4 字节 `unix timestamp + 4` 字节自增计数表示。这个值很重要，在选举(如 `master` 宕机时)新 `primary` 时，会选择 `ts` 最大的那个 `secondary` 作为新 `primary` 。

- `"h" : NumberLong("-2056314559875839750")`

    操作的唯一 `ID` 。

- `"op" : "i" 1`

    字节的操作类型。其中 `"i"` 代表 `insert` 操作， `"u"` 代表 `update` 操作， `"d"` 代表 `delete` 操作。

- `"ns" : "test.test"`

    该操作所在的集合。

- ```
    "o" : {
    "_id" : ObjectId("5c459393a4a6afa8f49832ad"),
    "user" : "user3",
    "sku" : 7,
    "price" : 7.84
    }

    ```





    Show moreShow more icon

    对文档字段进行操作的值。

- `o2:`

    该字段只在进行 `update` 时才会体现。


## 常用命令

以下命令能帮助您更快的上手副本集。

- 查看 `oplog` 大小





    ```
    rs.printReplicationInfo() # 2.6+
    db.getReplicationInfo()

    ```





    Show moreShow more icon

- 查看延迟





    ```
    db.printSlaveReplicationInfo()

    ```





    Show moreShow more icon

- 查看副本集的状态





    ```
    rs.status()

    ```





    Show moreShow more icon

- 设置 `SECONDARY` 为可读





    ```
    rs.slaveOk()
    db.isMaster()
    rs.config()

    ```





    Show moreShow more icon

- 增加副本集成员





    ```
    rs.add()

    ```





    Show moreShow more icon

- 删除副本集成员





    ```
    rs.remove()

    ```





    Show moreShow more icon

- `rs` 帮助命令





    ```
    rs.help()

    ```





    Show moreShow more icon

- 使 `Primary` 节点退化为 `Secondary` ，并在一段时间内不能参与选举。





    ```
    rs.stepDown(time)

    ```





    Show moreShow more icon

- 阻止选举，始终处于备份节点状态。比如主节点需要做一些维护，不希望其他成员选举为主节点，可以在每个备份节点上执行。强制他们出于备份节点状态。





    ```
    rs.freeze(time)

    ```





    Show moreShow more icon


## 常见问题

在使用副本集的过程中，对于一些有较高频率碰到的问题，在此做了一些总结，用以帮助想要学习 `MongoDB` 副本集的小伙伴们少走弯路。

1. 初始化节点的时候报错怎么办？

    - 如果某个节点有数据，则必须到该节点初始化，否则会报错 “errmsg” : “couldn’t initiate : member 127.0.0.1:27017 has data already, cannot initiate set. All members except initiator must be empty.”
    - 对于一机多实例，如果单独用 rs.initiate()则会默认用主机名:端口初始化”me” : “mytest201.com:27018″导致 rs.add (“127.0.0.1:27018″)”errmsg” : “exception: can’t use localhost in repl set member names except when using it for all members”。
2. 为什么使用奇数个数成员？

    MongoDB 本身设计的就是一个可以跨 IDC 的分布式数据库，假设一种场景，偶数 8 节点，IDC 网络故障发生分裂，这两个分裂后的 IDC 各持有 4 个节点，因为投票数未超过半数，所以无法选举出新 Primary。

3. 仲裁节点作用？

    仲裁节点并不需要太多系统资源，也并不持有数据本身，而是参与投票。 投票选举机制，根据数据最后操作、更新时间戳等判定，若有两方都为最新，且票数相等，此环节需要等待若干分钟。仲裁节点打破这个僵局。当节点数目为奇数时，可以不需要。 当节点数目为偶数个时，需要部署一个仲裁节点，否则偶数个节点，当主节点挂了后，其他节点会变为只读。

4. 将一台节点移除副本集，并且注释掉相关配置文件后，重启启动该节点登录后会有 `REMOVED` 字样，如何去掉？

    登录后在 `local` 数据库下查找 `db.system.replset.find()` ，如果还有之前的副本集相关记录，删掉即可， `db.system.replset.remove({})` 。

5. 什么情况下会出现 OTHER 状态？

    当使用 `rs.reconfig` 重新对副本集进行初始化的时候， `SECONDARY` 和 `ARBITER` 节点会出现 `OTHER` 状态。

6. 对于连接 MongoDB 的客户端都是外网的情况，如何配置副本级？

    初始化副本级的时候 host 必须为外网 IP。


## 结束语

通过对本文的学习，希望您能对 MongoDB 副本集技术有一定的了解，并且可以将其应用到工作当中，以帮助您构建一个健壮的 MongoDB 应用程序。