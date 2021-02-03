# MySQL 事务隔离级别和锁
事务特性和 InnoDB 锁

**标签:** 数据库

[原文链接](https://developer.ibm.com/zh/articles/os-mysql-transaction-isolation-levels-and-locks/)

王 拓

发布: 2019-08-26

* * *

## 事务及其特性

数据库事务（简称:事务）是数据库管理系统执行过程中的一个逻辑单位，由一个有限的数据库操作序列构成。事务的使用是数据库管理系统区别文件系统的重要特征之一。

事务拥有四个重要的特性：原子性（Atomicity）、一致性（Consistency）、隔离性（Isolation）、持久性（Durability），人们习惯称之为 ACID 特性。下面我逐一对其进行解释。

**[IBM Compose for MySQL](https://cloud.ibm.com/catalog/services/compose-for-mysql?cm_sp=ibmdev-_-developer-articles-_-cloudreg)**

[IBM Cloud](https://cocl.us/IBM_CLOUD_GCG) 上提供的 Compose for MySQL 数据库服务可以帮助您更好地管理 MySQL，其特性包含自扩展部署、高可用、以及自动化无停止备份。

- **原子性** （Atomicity）

    事务开始后所有操作，要么全部做完，要么全部不做，不可能停滞在中间环节。事务执行过程中出错，会回滚到事务开始前的状态，所有的操作就像没有发生一样。例如，如果一个事务需要新增 100 条记录，但是在新增了 10 条记录之后就失败了，那么数据库将回滚对这 10 条新增的记录。也就是说事务是一个不可分割的整体，就像化学中学过的原子，是物质构成的基本单位。

- **一致性** （Consistency）

    指事务将数据库从一种状态转变为另一种一致的的状态。事务开始前和结束后，数据库的完整性约束没有被破坏。例如工号带有唯一属性，如果经过一个修改工号的事务后，工号变的非唯一了，则表明一致性遭到了破坏。

- **隔离性** （Isolation）

    要求每个读写事务的对象对其他事务的操作对象能互相分离，即该事务提交前对其他事务不可见。 也可以理解为多个事务并发访问时，事务之间是隔离的，一个事务不应该影响其它事务运行效果。这指的是在并发环境中，当不同的事务同时操纵相同的数据时，每个事务都有各自的完整数据空间。由并发事务所做的修改必须与任何其他并发事务所做的修改隔离。例如一个用户在更新自己的个人信息的同时，是不能看到系统管理员也在更新该用户的个人信息（此时更新事务还未提交）。

    注：MySQL 通过锁机制来保证事务的隔离性。

- **持久性** （Durability）

    事务一旦提交，则其结果就是永久性的。即使发生宕机的故障，数据库也能将数据恢复，也就是说事务完成后，事务对数据库的所有更新将被保存到数据库，不能回滚。这只是从事务本身的角度来保证，排除 RDBMS（关系型数据库管理系统，例如 Oracle、MySQL 等）本身发生的故障。

    注：MySQL 使用 `redo log` 来保证事务的持久性。


## 事务的隔离级别

SQL 标准定义的四种隔离级别被 ANSI（美国国家标准学会）和 ISO/IEC（国际标准）采用，每种级别对事务的处理能力会有不同程度的影响。

我们分别对四种隔离级别从并发程度由高到低进行描述，并用代码进行演示，数据库环境为 MySQL 5.7。

### READ UNCOMMITTED（读未提交）

该隔离级别的事务会读到其它未提交事务的数据，此现象也称之为 **脏读** 。

1. 准备两个终端，在此命名为 mysql 终端 1 和 mysql 终端 2，再准备一张测试表 `test` ，写入一条测试数据并调整隔离级别为 `READ UNCOMMITTED` ，任意一个终端执行即可。





    ```
    SET @@session.transaction_isolation = 'READ-UNCOMMITTED';
    create database test;
    use test;
    create table test(id int primary key);
    insert into test(id) values(1);

    ```





    Show moreShow more icon

2. 登录 mysql 终端 1，开启一个事务，将 ID 为 `1` 的记录更新为 `2` 。





    ```
    begin;
    update test set id = 2 where id = 1;
    select * from test; -- 此时看到一条ID为2的记录

    ```





    Show moreShow more icon

3. 登录 mysql 终端 2，开启一个事务后查看表中的数据。





    ```
    use test;
    begin;
    select * from test; -- 此时看到一条 ID 为 2 的记录

    ```





    Show moreShow more icon


最后一步读取到了 mysql 终端 1 中未提交的事务（没有 commit 提交动作），即产生了 **脏读** ，大部分业务场景都不允许脏读出现，但是此隔离级别下数据库的并发是最好的。

### READ COMMITTED（读提交）

一个事务可以读取另一个已提交的事务，多次读取会造成不一样的结果，此现象称为不可重复读问题，Oracle 和 SQL Server 的默认隔离级别。

1. 准备两个终端，在此命名为 mysql 终端 1 和 mysql 终端 2，再准备一张测试表 `test` ，写入一条测试数据并调整隔离级别为 `READ COMMITTED` ，任意一个终端执行即可。





    ```
    SET @@session.transaction_isolation = 'READ-COMMITTED';
    create database test;
    use test;
    create table test(id int primary key);
    insert into test(id) values(1);

    ```





    Show moreShow more icon

2. 登录 mysql 终端 1，开启一个事务，将 ID 为 `1` 的记录更新为 `2` ，并确认记录数变更过来。





    ```
    begin;
    update test set id = 2 where id = 1;
    select * from test; -- 此时看到一条记录为 2

    ```





    Show moreShow more icon

3. 登录 mysql 终端 2，开启一个事务后，查看表中的数据。





    ```
    use test;
    begin;
    select * from test; -- 此时看一条 ID 为 1 的记录

    ```





    Show moreShow more icon

4. 登录 mysql 终端 1，提交事务。





    ```
    commit;

    ```





    Show moreShow more icon

5. 切换到 mysql 终端 2。





    ```
    select * from test; -- 此时看到一条 ID 为 2 的记录

    ```





    Show moreShow more icon


mysql 终端 2 在开启了一个事务之后，在第一次读取 `test` 表（此时 mysql 终端 1 的事务还未提交）时 ID 为 `1` ，在第二次读取 `test` 表（此时 mysql 终端 1 的事务已经提交）时 ID 已经变为 `2` ，说明在此隔离级别下已经读取到已提交的事务。

### REPEATABLE READ（可重复读）

该隔离级别是 MySQL 默认的隔离级别，在同一个事务里， `select` 的结果是事务开始时时间点的状态，因此，同样的 `select` 操作读到的结果会是一致的，但是，会有 **幻读** 现象。MySQL 的 InnoDB 引擎可以通过 `next-key locks` 机制（参考下文 [行锁的算法](#行锁的算法) 一节）来避免幻读。

1. 准备两个终端，在此命名为 mysql 终端 1 和 mysql 终端 2，准备一张测试表 test 并调整隔离级别为 `REPEATABLE READ` ，任意一个终端执行即可。





    ```
    SET @@session.transaction_isolation = 'REPEATABLE-READ';
    create database test;
    use test;
    create table test(id int primary key,name varchar(20));

    ```





    Show moreShow more icon

2. 登录 mysql 终端 1，开启一个事务。





    ```
    begin;
    select * from test; -- 无记录

    ```





    Show moreShow more icon

3. 登录 mysql 终端 2，开启一个事务。





    ```
    begin;
    select * from test; -- 无记录

    ```





    Show moreShow more icon

4. 切换到 mysql 终端 1，增加一条记录并提交。





    ```
    insert into test(id,name) values(1,'a');
    commit;

    ```





    Show moreShow more icon

5. 切换到 msyql 终端 2。





    ```
    select * from test; --此时查询还是无记录

    ```





    Show moreShow more icon

    通过这一步可以证明，在该隔离级别下已经读取不到别的已提交的事务，如果想看到 mysql 终端 1 提交的事务，在 mysql 终端 2 将当前事务提交后再次查询就可以读取到 mysql 终端 1 提交的事务。我们接着实验，看看在该隔离级别下是否会存在别的问题。

6. 此时接着在 mysql 终端 2 插入一条数据。





    ```
    insert into test(id,name) values(1,'b'); -- 此时报主键冲突的错误

    ```





    Show moreShow more icon


也许到这里您心里可能会有疑问，明明在第 5 步没有数据，为什么在这里会报错呢？其实这就是该隔离级别下可能产生的问题，MySQL 称之为 **幻读** 。注意我在这里强调的是 MySQL 数据库，Oracle 数据库对于幻读的定义可能有所不同。

### SERIALIZABLE（序列化）

在该隔离级别下事务都是串行顺序执行的，MySQL 数据库的 InnoDB 引擎会给读操作隐式加一把读共享锁，从而避免了脏读、不可重读复读和幻读问题。

1. 准备两个终端，在此命名为 mysql 终端 1 和 mysql 终端 2，分别登入 mysql，准备一张测试表 test 并调整隔离级别为 `SERIALIZABLE` ，任意一个终端执行即可。





    ```
    SET @@session.transaction_isolation = 'SERIALIZABLE';
    create database test;
    use test;
    create table test(id int primary key);

    ```





    Show moreShow more icon

2. 登录 mysql 终端 1，开启一个事务，并写入一条数据。





    ```
    begin;
    insert into test(id) values(1);

    ```





    Show moreShow more icon

3. 登录 mysql 终端 2，开启一个事务。





    ```
    begin;
    select * from test; -- 此时会一直卡住

    ```





    Show moreShow more icon

4. 立马切换到 mysql 终端 1,提交事务。





    ```
    commit;

    ```





    Show moreShow more icon


一旦事务提交，msyql 终端 2 会立马返回 ID 为 1 的记录，否则会一直卡住，直到超时，其中超时参数是由 `innodb_lock_wait_timeout` 控制。由于每条 `select` 语句都会加锁，所以该隔离级别的数据库并发能力最弱，但是有些资料表明该结论也不一定对，如果感兴趣，您可以自行做个压力测试。

表 1 总结了各个隔离级别下产生的一些问题。

##### 表 1\. 各个隔离级别下产生的一些问题

**隔离级别****脏读****不可重复读****幻读**读未提交可以出现可以出现可以出现读提交不允许出现可以出现可以出现可重复读不允许出现不允许出现可以出现序列化不允许出现不允许出现不允许出现

## MySQL 中的锁

锁也是数据库管理系统区别文件系统的重要特征之一。锁机制使得在对数据库进行并发访问时，可以保障数据的完整性和一致性。对于锁的实现，各个数据库厂商的实现方法都会有所不同。本文讨论 MySQL 中的 InnoDB 引擎的锁。

### 锁的类型

InnoDB 实现了两种类型的行级锁：

- **共享锁** （也称为 S 锁）：允许事务读取一行数据。

    可以使用 SQL 语句 `select * from tableName where... lock in share mode;` 手动加 S 锁。

- **独占锁** （也称为 X 锁）：允许事务删除或更新一行数据。

    可以使用 SQL 语句 `select * from tableName where... for update`; 手动加 X 锁。


S 锁和 S 锁是 **兼容** 的，X 锁和其它锁都 **不兼容** ，举个例子，事务 T1 获取了一个行 r1 的 S 锁，另外事务 T2 可以立即获得行 r1 的 S 锁，此时 T1 和 T2 共同获得行 r1 的 S 锁，此种情况称为 **锁兼容** ，但是另外一个事务 T2 此时如果想获得行 r1 的 X 锁，则必须等待 T1 对行 r 锁的释放，此种情况也成为 **锁冲突** 。

为了实现多粒度的锁机制，InnoDB 还有两种内部使用的 **意向锁** ，由 InnoDB 自动添加，且都是表级别的锁。

- **意向共享锁** （IS）：事务即将给表中的各个行设置共享锁，事务给数据行加 S 锁前必须获得该表的 IS 锁。
- **意向排他锁** （IX）：事务即将给表中的各个行设置排他锁，事务给数据行加 X 锁前必须获得该表 IX 锁。

意向锁的主要目的是为了使得 **行锁** 和 **表锁** 共存。表 2 列出了行级锁和表级意向锁的兼容性。

##### 表 2\. 行级锁和表级意向锁的兼容性

**锁类型****X****IX****S****IS****X**冲突冲突冲突冲突**IX**冲突兼容冲突兼容**S**冲突冲突兼容兼容**IS**冲突兼容兼容兼容

### 行锁的算法

InnoDB 存储引擎使用三种行锁的算法用来满足相关事务隔离级别的要求。

- **Record Locks**

    该锁为索引记录上的锁，如果表中没有定义索引，InnoDB 会默认为该表创建一个隐藏的聚簇索引，并使用该索引锁定记录。

- **Gap Locks**

    该锁会锁定一个范围，但是不括记录本身。可以通过修改隔离级别为 `READ COMMITTED` 或者配置 `innodb_locks_unsafe_for_binlog` 参数为 `ON` 。

- **Next-key Locks**

    该锁就是 Record Locks 和 Gap Locks 的组合，即锁定一个范围并且锁定该记录本身。InnoDB 使用 Next-key Locks 解决幻读问题。需要注意的是，如果索引有唯一属性，则 InnnoDB 会自动将 Next-key Locks 降级为 Record Locks。举个例子，如果一个索引有 1, 3, 5 三个值，则该索引锁定的区间为 `(-∞,1], (1,3], (3,5], (5,+ ∞)` 。


### 死锁

**死锁** 是指两个或两个以上的进程在执行过程中，由于竞争资源或者由于彼此通信而造成的一种阻塞的现象，若无外力作用，它们都将无法推进下去。此时称系统处于死锁状态或系统产生了死锁，这些永远在互相等待的进程称为死锁进程。

InnoDB 引擎采取的是 `wait-for graph` 等待图的方法来自动检测死锁，如果发现死锁会自动回滚一个事务。

下面我们通过一个示例来了解死锁。

1. 准备两个终端，在此命名为 mysql 终端 1 和 mysql 终端 2，分别登入 mysql，再准备一张测试表 `test` 写入两条测试数据，并调整隔离级别为 `SERIALIZABLE` ，任意一个终端执行即可。





    ```
    SET @@session.transaction_isolation = 'REPEATABLE-READ';
    create database test;
    use test;
    create table test(id int primary key);
    insert into test(id) values(1),(2);

    ```





    Show moreShow more icon

2. 登录 mysql 终端 1，开启一个事务，手动给 ID 为 `1` 的记录加 X 锁。





    ```
    begin;
    select * from test where id = 1 for update;

    ```





    Show moreShow more icon

3. 登录 mysql 终端 2，开启一个事务，手动给 ID 为 `2` 的记录加 X 锁。





    ```
    begin;
    select * from test where id = 2 for update;

    ```





    Show moreShow more icon

4. 切换到 mysql 终端 1，手动给 ID 为 `2` 的记录加 X 锁，此时会一直卡住，因为此时在等待第 3 步中 X 锁的释放，直到超时，超时时间由 `innodb_lock_wait_timeout` 控制。





    ```
    select * from test where id = 2 for update;

    ```





    Show moreShow more icon

5. 在锁超时前立刻切换到 mysql 终端 2，手动给 ID 为 `1` 的记录加 X 锁，此时又会等待第 2 步中 X 所的释放，两个终端都在等待资源的释放，所以 InnoDB 引擎会立马检测到死锁产生，自动回滚一个事务，以防止死锁一直占用资源。





    ```
    select * from test where id = 1 for update;
    ERROR 1213 (40001): Deadlock found when trying to get lock; try restarting transaction

    ```





    Show moreShow more icon

    此时，通过 `show engine innodb status\G` 命令可以看到 `LATEST DETECTED DEADLOCK` 相关信息，即表明有死锁发生；或者通过配置 `innodb_print_all_deadlocks` （MySQL 5.6.2 版本开始提供）参数为 `ON` 将死锁相关信息打印到 MySQL 的错误日志。


### 锁的优化建议

锁如果利用不好，会给业务造成大量的卡顿现象，在了解了锁相关的一些知识点后，我们可以有意识的去避免锁带来的一些问题。

1. 合理设计索引，让 InnoDB 在索引键上面加锁的时候尽可能准确，尽可能的缩小锁定范围，避免造成不必要的锁定而影响其他 Query 的执行。
2. 尽可能减少基于范围的数据检索过滤条件，避免因为间隙锁带来的负面影响而锁定了不该锁定的记录。
3. 尽量控制事务的大小，减少锁定的资源量和锁定时间长度。
4. 在业务环境允许的情况下，尽量使用较低级别的事务隔离，以减少 MySQL 因为实现事务隔离级别所带来的附加成本。

## 结束语

通过阅读本文，可以让您对数据库的事务还有事务的隔离级别有个基本的了解，同时也介绍了 MySQL 中 InnoDB 引擎中一些锁相关的知识，从而可以让您利用关系型数据库系统设计一个更为健壮的业务模型。