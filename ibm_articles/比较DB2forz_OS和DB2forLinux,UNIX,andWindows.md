# 比较 DB2 for z/OS 和 DB2 for Linux, UNIX, and Windows
面向需要跨平台技能的 DB2 数据库管理员的概述

**标签:** 数据库

[原文链接](https://developer.ibm.com/zh/articles/dm-1108compdb2luwzos/)

Balachandran Chandrasekaran

发布: 2011-11-14

* * *

免费下载： IBM® DB2® Express-C 9.7.2 免费版 或者 DB2® 9.7 for Linux®, UNIX®, and Windows® 试用版下载更多的 [IBM 软件试用版](http://www.ibm.com/developerworks/cn/downloads/) ，并加入 [IBM 软件下载与技术交流群组](https://www.ibm.com/developerworks/mydeveloperworks/groups/service/html/communityview?communityUuid=38997896-bb16-451a-aa97-189a27a3cd5a/?lang=zh) ，参与在线交流。

## 简介

针对 z/OS 和 Linux, UNIX, Windows (LUW) 的 DB2 产品版本如此不同，所以 DB2 数据库管理员通常必须具备不同的技能才能在不同的平台上使用 DB2。本文将帮助您开始使用在其他平台上可用的 DB2。本文假定您拥有各个硬件和操作系统所需的知识。本文旨在介绍跨平台处理数据库管理职责的不同技术，不应用作比较两个产品的特性或功能的比较性文章。

DB2 LUW 也有各种版本可用，包括 DB2 Express-C（免费）、DB2 Express Edition、DB2 Personal Edition、DB2 Workgroup Server Edition、DB2 Enterprise Server Edition 和 DB2 Advanced Enterprise Server Edition。对于您的企业来说，理想的功能组合将取决于基础产品的功能、企业的规模、用户数量和其他许可条款。在这篇文章中，DB2 z/OS 是指 DB2 10 for z/OS；而 DB2 LUW 是指 DB2 9.7 for Linux, UNIX and Windows。还应指出的是，System Z 服务器上的 Linux 会在 DB2 LUW 而非 DB2 z/OS 上运行。

## 关于子系统、实例和数据库

您会发现在这两个平台上提及和访问 DB2 的方式有所不同。在 z/OS 中，DB2 服务器称为 _子系统_ ，在 LUW 中则称为 _实例_ 。DB2 在 z/OS 系统中的一组地址空间下运行其任务（在 [DB2 内存结构](#db2-内存结构) 中会对此进行讨论），在 Linux 和 UNIX 操作系统中，会将任务作为进程运行，而在 Windows 中，则是将任务作为服务运行（除非以 /D 选项开头）。每个 DB2 z/OS 服务器都用一个 4 个字符的子系统标识符 (SSID) 来标识，而每个 DB2 LUW 服务器都用一个实例名来标识，该实例名通过 DB2INSTANCE 环境变量生效。

_数据库_ 一词在每个平台中的含义是什么？DB2 LUW 中的数据库有其自己的内存区域、进程和恢复日志，而 DB2 z/OS 中的数据库是若干表空间和索引空间的逻辑集合，对于其下创建的所有对象，很少将参数定义为默认参数。DB2 LUW 下的每个数据库都拥有其独有的目录表，即 _每个 DB2 子系统_ 保存目录表。

### 系统对象呢？

- DB2 z/OS 子系统中的 DSNDB06 数据库以及 DB2 LUW 中的 SYSCATSPACE 表空间会保存所有的 DB2 目录表。
- 在非数据共享的 DB2 z/OS 中，DSNDB07 是工作文件数据库。您在 DSNDB07 创建的表空间将用作处理 SQL 联接和排序以及临时表的存储。在 DB2 LUW 中，，一个名为 TEMPSPACE1 的默认系统临时表空间可用来满足 DB2 DBM 需求（排序或联接等）。您应当创建一个或多个 USER 临时表空间来存储临时表的数据。
- 如果您在 DB2 z/OS 中创建表空间或表时没有指定数据库名称，那么 DSNDB04 会成为所创建对象的默认用户数据库。如果您创建了一个表并且没有指定表空间，那么 DB2 会为您创建一个隐式的表空间。在 DB2 LUW 中，假定您还没有创建任何用户表空间，如果没有显式指定表空间，则会用 USERSPACE1 存放您创建的表。
- DB2 z/OS 使用了在特殊数据库 DSNDB01 中正常运行所需的一些信息。Boot Strap Dataset (BSDS) 是一个 z/OS VSAM 数据集，它包含一些有关 DB2 恢复日志数据集名称、DB2 检查点和 DB2 位置细节的关键信息。另一方面，DB2 LUW 数据库还存储跨本地数据库目录下的多个操作系统文件正常运行所需的若干信息。

### 您将如何建立您的第一个远程连接？

要访问远程 DB2 LUW 服务器，则应当通过指定 DB2 实例为节点编写目录，然后在已编写目录的节点上为该实例下的目标数据库编写目录。 [清单 1](#以下命令会为您提供远程-db2-luw-服务器的主机名和端口-服务) 显示了如何在 AIX 主机上查找有关 DB2 服务器的细节。另一方面，您必须使用其位置位一个远程 DB2 z/OS 编写目录。而且，您可以通过 `CATALOG DCS DATABASE` 命令在 Database Connection Services (DCS) 目录中存储与主机 DB2 相关的信息，如 [清单 2](#在-db2-客户端上运行以下命令为上述服务器编写目录) 中所示。DB2 z/OS 的位置会在安装 DB2 子系统时予以定义，您可以通过在 DB2 z/OS 服务器上执行 `-DIS DDF` 命令来找到该位置，如 [清单 3](#以下是-db2-z-os-上的一个-dis-ddf-命令输出样例) 中所示。

DB2 z/OS 子系统在名为 Communication Database (CDB) 的一组目录表中包含与远程子系统建立连接的信息。另外需要一个 IBM DB2 Connect 许可，以便远程连接到位于主机系统（包括 z/OS）上的 DB2 数据库。您可以参阅 参考资料 部分的链接，了解有关 IBM DB2 Connect 的更多信息。

##### 以下命令会为您提供远程 DB2 LUW 服务器的主机名和端口/服务

```
$ db2set -all
[i] DB2COMM=TCPIP
[g] DB2SYSTEM=host.ibm.com
$ db2 get DB2INSTANCE

The current database manager instance is:  db2inst1

$ db2 "get dbm cfg " | grep -i '(svcename)'
TCP/IP Service name                          (SVCENAME) = DB2_db2inst1
$ cat /etc/services | grep -i DB2_db2inst1
DB2_db2inst1    60000/tcp

```

Show moreShow more icon

##### 在 DB2 客户端上运行以下命令为上述服务器编写目录

```
$ db2 "CATALOG TCPIP NODE aixnode REMOTE host.ibm.com SERVER 60000
> REMOTE_INSTANCE db2inst1"
$ db2 "CATALOG DATABASE sampledb AS sampledb AT NODE aixnode"
$ db2 "CONNECT TO sampledb USER userid USING passwd"

$ db2 "CATALOG TCPIP NODE zosnode REMOTE zserver.ibm.com SERVER 446"
$ db2 "CATALOG DATABASE nyc AS zosdb2 AT NODE zosnode"
$ db2 "CATALOG DCS DATABASE nyc AS newyork"
$ db2 "CONNECT TO zosdb2 USER zosuser USING passwd"

```

Show moreShow more icon

##### 以下是 DB2 z/OS 上的一个 -DIS DDF 命令输出样例

```
-DIS DDF
DSNL080I  @ DSNLTDDF DISPLAY DDF REPORT FOLLOWS:
DSNL081I STATUS=STARTD
DSNL082I LOCATION           LUNAME            GENERICLU
DSNL083I NEWYORK            A.B               -NONE
DSNL084I TCPPORT=446       SECPORT=0     RESPORT=4463  IPNAME=-NONE
DSNL085I IPADDR=::ip-addr
DSNL086I SQL    DOMAIN=zserver.ibm.com
DSNL099I DSNLTDDF DISPLAY DDF REPORT COMPLETE

```

Show moreShow more icon

在 DB2 LUW 中，您可以发出 `LIST DB DIRECTORY` 和 `LIST NODE DIRECTORY` 命令来分别查询已编写目录的数据库和节点。

## DB2 内存结构

在 DB2 LUW 中，DB2 服务器所消耗的内存大致分为实例级的数据库管理器共享内存和数据库级的数据库共享内存。以下列表总结了 DB2 z/OS 地址空间及其内存组件。

- System Services Address Space (ssidMSTR) 控制整个系统服务，包括日志管理和命令处理。
- Database Manager Address Space (DBM1) 托管所有缓冲池、Environment Descriptor Manager (EDM) 池、Sort 池和 Record Identifier (RID) 池。拥有最大 DB2 地址空间的 DBM1 可与 DB2 LUW 的数据库共享内存相匹敌。EDM 池包括用于目录访问、动态 SQL、计划/包和游标的缓存块。
- Distributed Data Facility 地址空间 (ssidDDF) 负责从远程客户端进行分布式数据访问，以及从 DB2 子系统访问远程数据。
- Internal Resource Lock Manager (IRLM) 负责锁和其他序列化相关活动。在 IRLM 中，您可以定义可与 DB2 LUW 中的 locklist 参数相匹敌的锁存储。
- z/OS Workload Managed (WLM) DB2 地址空间运行存储过程和例程。您还会发现 DB2 Stored Procedure Address Space (SPAS)，它运行在 DB2 v8 出现之前创建的例程。
- Agent Allied 地址空间针对连接到 DB2 的每个地址空间运行。这些类似于 DB2 LUW 下的私有代理内存。
- DB2 实用程序运行于其批量地址空间之下，根据需要使用来自其他 DB2 地址空间的内存。

DB2 z/OS 不允许您创建具有任何名称的缓冲池，而 DB2 LUW 允许您通过 `CREATE BUFFERPOOL` SQL 语句随意命名您的缓冲池。在 DB2 z/OS 中，有 50 个页面大小不同的预定义缓冲池。您必须根据您想要的页面大小选择一个缓冲池名称，并通过 `-ALTER BUFFERPOOL` 命令激活它。而且，在 DB2 z/OS 中创建表空间时，您将通过具有所需页面大小的缓冲池 _间接_ 定义其页面大小。清单 4 和清单 5 展示和对比如何分别在 DB2 z/OS 和 DB2 LUW 中查询有关缓冲池的详细信息。

##### z/OS 中的缓冲池

```
-DIS BPOOL(BP0)
DSNB401I  DB1S BUFFERPOOL NAME BP0, BUFFERPOOL ID 0, USE COUNT 357
DSNB402I  DB1S BUFFER POOL SIZE = 12500 BUFFERS  AUTOSIZE = NO
              ALLOCATED       =    12500   TO BE DELETED   =        0
              IN-USE/UPDATED  =      284   BUFFERS ACTIVE  =    12500
DSNB406I  DB1S PGFIX ATTRIBUTE -
              CURRENT = NO
              PENDING = NO
            PAGE STEALING METHOD = LRU
DSNB404I  DB1S THRESHOLDS -
             VP SEQUENTIAL    = 80
             DEFERRED WRITE   = 30   VERTICAL DEFERRED WRT  =  5,  0
             PARALLEL SEQUENTIAL =50   ASSISTING PARALLEL SEQT=  0
DSN9022I  DB1S DSNB1CMD '-DIS BPOOL' NORMAL COMPLETION

```

Show moreShow more icon

##### LUW 中的缓冲池

```
$ db2  "SELECT SUBSTR(BPNAME,1,15) BPOOL, BUFFERPOOLID, NPAGES, PAGESIZE,
>  NUMBLOCKPAGES, BLOCKSIZE FROM SYSIBM.SYSBUFFERPOOLS"

BPOOL           BUFFERPOOLID NPAGES      PAGESIZE    NUMBLOCKPAGES BLOCKSIZE
--------------- ------------ ----------- ----------- ------------- -----------
IBMDEFAULTBP               1        2000        4096             0           0

1 record(s) selected.

```

Show moreShow more icon

作为一个新平台上的一名 DBA，您需要知道当前的 DB2 服务器配置和按需设置或更新它们的方法。通过 IBM 提供的 MVS Job DSNTIJUZ，DB2 z/OS 子系统参数被汇编并连接到一个可执行模块。这些子系统参数被称为 DSNZPARM（或简称为 ZPARM）。作为一名 DB2 系统管理员，要更新任何这些 DB2 子系统参数，您应当执行 DSNTIJUZ，后接 `-SET SYSPARM RELOAD` 命令，通过这些来动态刷新 DSNZPARM 模块。一个类似的任务将需要在 LUW 平台上执行 `UPDATE DBM CFG` 或 `UPDATE DB CFG` 命令，具体取决于您希望更新的是数据库管理器 (DBM) 配置还是数据库 (DB) 配置。DB2 LUW 环境还可以通过不同类型的注册表加以控制。您可以使用 `db2set` 系统命令查看和设置这些注册表下的变量。

DB2 LUW 为您提供一种管理各种内存结构的简单方法，即自调优内存管理 (STMM)。通过 DB CFG 参数 _self\_tuning\_mem_ 启用 STMM 之后，您可以让 DB2 根据工作负载的波动调优内存结构。

另外，DB2 LUW 提供了一个名为 `AUTOCONFIGURE` 的非常有用的命令，您可以使用它让 DB2 根据某种类型的工作负载（OLTP 或 DSS 类型）设置 DBM 和 DB 参数。DB2 z/OS 提供一个选项，让 DB2 根据工作负载自动调整缓冲池的大小。您可以使用带有 AUTORESIZE(YES) 操作数的 `ALTER BUFFERPOOL` 命令来实现这一点。DB2 会在 z/OS WLM 的帮助下调优缓冲池。

## 表空间和存储布局

表空间将物理容器分组，以容纳各个表的数据及其索引。因为他们将数据库数据映射到底层物理存储（原始磁盘、文件系统、卷），在一个新的平台中理解它们对于 DBA 来说很重要。DB2 z/OS 将表数据存储在 Linear Virtual Storage Access Method（VSAM，类型 LDS）数据集上。DB2 LUW 还会根据表空间的类型将表数据存储到文件、目录和原始设备上。

在 DB2 z/OS 中保存索引的表空间称为 _索引空间_ ，当您创建一个索引时，DB2 自动为新的索引创建索引空间。与 DB2 LUW 中不同，您不能在一个索引空间上放置多个索引。在 DB2 z/OS 中创建一个范围分区的表时，底层的表空间也被分区，并且这样一个分区的表空间无法容纳多个表。在 DB2 for LUW 中，您可以将这种表的各个分区放在一个或多个表空间上。

在 DB2 LUW 中，通常几个表纳入一个单独的表空间，而在 DB2 z/OS 中大部分时候您最终会为每个表空间创建一个表。这是因为在 DB2 z/OS 中，大部分数据管理实用程序在表空间级别（如果表被分区，则在分区级别上）运行，而且您可以更加自由地选择存储属性或表维护窗口。然而在表的尺寸较小时，您在每个表空间可以创建多个表。您还应当知道，在表空间下创建多个表在运行多个 DB2 实用程序时会影响整个表空间。

### 如何定义 DB2 存储？

在 DB2 LUW 中，您有以下两种定义 DB2 存储的方法。清单 6 显示了如何使用这些方法的示例。

- 使用 Database Managed Space (DMS) 方法时，需要显式地为每个表指定一组文件或原始设备，并让 DB2 数据库管理它们。现在鼓励将 DB2 数据放在原始设备上，因为在默认情况下，DB2 会在没有文件系统缓存时创建 DMS 表空间。使用该方法时需要您指定最初空间大小，并且应当可以使用该打小成功创建表空间。对于 DMS 表空间，您可以按需添加或删除容器，还可以指定是否应当允许扩展容器，对于这样的容器，您可以定义应当扩展的尺寸和各个容器的最大尺寸。您可以通过对整个数据库启用 DB2 Automatic Storage 来简化 DB2 存储管理。您可以通过 `CREATE DATABASE` 或 `ALTER DATABASE` 命令在数据库级启用 Automatic Storage（默认设置适用于新数据库）。该功能需要您在数据库级指定一个或多个存储路径，您希望 DB2 在该路径中为每个表空间创建容器。
- 使用 System Managed Space (SMS) 方法，您要指定位置和目录，让 OS 文件管理器处理空间。空间会按需予以分配。只有通过添加空间到底层文件系统，您才能够将空间添加到您的表空间。此外，您不能直接将一个 SMS 表空间转化为 DMS 表空间。

##### DB2 LUW 中的不同存储选项

```
-- Create a database with automatic storage disabled; Each table space
-- will use a default type depending upon its type.
db2 "CREATE DATABASE testdb AUTOMATIC STORAGE NO"

-- Enable Automatic Storage on automatic storage disabled database
db2 "ALTER DATABASE testdb ADD STORAGE ON '/home/dbstor1'"

-- Create a Tablespace managed by automatic storage
db2 "CREATE TABLESPACE autotbsp MANAGED BY AUTOMATIC STORAGE AUTORESIZE YES
> INITIALSIZE 1G INCREASESIZE 10M MAXSIZE 2G"

-- Create a DMS Tablespace with 100 4KB pages. DB2 will fix increment size
-- and max size.
db2 "CREATE TABLESPACE dmstbsp PAGESIZE 4096 MANAGED BY DATABASE USING
> (FILE '/home/dbstor1/dbfile01.dbf' 100)"

-- Add additional container to DMS table space
db2 "ALTER TABLESPACE dmstbsp ADD (FILE '/home/dbstor1/dbfile02.dbf' 100)"

-- Alter the DMS table space to use Automatic Storage
db2 "ALTER TABLESPACE dmstbsp MANAGED BY AUTOMATIC STORAGE"

-- Spread data to newly added storage paths
db2 "ALTER TABLESPACE dmstbsp REBALANCE"

-- Create SMS table space
db2 "CREATE TABLESPACE smstbsp MANAGED BY SYSTEM USING ('/home/db2inst1/smsdir')"
db2 "CREATE TABLE T1(C1 CHAR(1)) IN smstbsp INDEX IN smstbsp"
db2 "CREATE INDEX I1 ON T1(C1)"
-- At least one file would be created for every object created in smstbsp
-- Operating system files for the SMS table space just created:
$ ls -l /home/db2inst1/smsdir
total 40
-rw-------    1 db2inst1 db2adm1        4096 Jul 27 17:12 SQL00002.DAT
-rw-------    1 db2inst1 db2adm1       12288 Jul 27 17:12 SQL00002.INX
-rw-------    1 db2inst1 db2adm1         512 Jul 27 17:06 SQLTAG.NAM

```

Show moreShow more icon

##### 在 DB2 z/OS 中定义存储

您可以指定管理底层存储的方式。

- 由用户管理
- 由 DB2 管理
- 由 z/OS SMS 管理
- 在 CREATE DDL 中支持 z/OS SMS 构造

在 DB2 z/OS 中创建表空间和索引时，您可以指定是由用户 (DBA) 还是由 DB2 来管理底层存储。在创建用户管理的表空间时，您必须已经使用 IDCAMS 程序创建了底层 z/OS LDS VSAM 数据集。该方法允许您规划 DB2 存储并小心地将 DB2 数据集放在选定的直接访问卷上。对于 DB2 管理的表空间，DB2 z/OS 提供了一个名为存储组 (STOGROUP) 的数据库对象。您可以为单个存储组分配一个或多个直接访问卷，并在创建对象时使用存储组名称。DB2 会管理所列卷内的空间，并将数据放在所列卷上的可用空间中。另一方面，为了管理存储，DB2 z/OS 可以与另一个名为 Storage Management System (SMS) 的 z/OS 子系统协同工作。在该技术中，整个 DB2 存储由 z/OS Storage 管理员编写的规则管理。该方法很常见，因为它让存储管理员完全控制不同的消费者所消耗的存储。DB2 还能够让 DB2 管理员在创建 DB2 存储组时指定各个 SMS 构造。清单 7 向您展示在 DB2 z/OS 中定义存储的示例。

##### DB2 z/OS 中的典型存储定义

```
-- Create an user managed table space in database TESTDB
CREATE TABLESPACE USERTS USING VCAT DB2X  IN TESTDB;
-- Create Non-SMS managed table space
CREATE STOGROUP dbasg  VOLUMES(VOLSER1,VOLSER2,VOLSER3) VCAT DB2T;
CREATE TABLESPACE NONSMSTS  USING STOGROUP dbasg
PRIQTY 1440 SECQTY 720;
-- Create SMS managed table space through VOLUMES ('*') operand
CREATE STOGROUP smssg  VOLUMES('*')  VCAT DB2T;
CREATE TABLESPACE SMSTS1 USING STOGROUP smssg;
-- Create storage group  with SMS clauses
CREATE STOGROUP newsmssg VCAT DB2T
MGMTCLAS mgmt1 STORCLAS storclas1 DATACLAS dc1 ;

```

Show moreShow more icon

对于 DB2 z/OS，不管您创建表空间使用的是什么方法，底层的 LDS VSAM 数据集都会有一个标准命名约定。在创建一个用户管理的对象时，DB2 会强制您遵循相同的命名约定。这种 LDS VSAM 数据集的最大尺寸将间接取决于您创建的表空间的类型和您在 CREATE 语句上指定的参数。

## 数据移动技术

如果您是一名 DBA，您可能已经执行了一些涉及到 z/OS 服务器的跨平台表数据移动。对于跨异构产品和平台的简单移动，DB2 z/OS 支持分隔格式的数据的卸载和加载。在 DB2 z/OS 中，基本 DB2 工具附带一些核心工具，而几个实用程序则分装在一个名为 Utilities Suite 的不同函数下。

另一方面，DB2 LUW 实用程序提供基本的 DB2 工具。通常通过在批处理模式下执行 DSNUTILB 程序来运行 DB2 z/OS 实用程序。z/OS 中的批处理作业使用 Job Control Language (JCL) 进行编码，如清单 8 所示。因此，在着手使用 DB2 z/OS 实用程序之前，您可能希望学习一点或更多关于 JCL 的知识。除了 DB2 z/OS 中也提供的 LOAD 实用程序之外，DB2 LUW 还提供了 EXPORT 和 IMPORT。IMPORT 实用程序可以执行一系列表中的插入操作，从而触发在表上定义的任何触发器，这在您使用 LOAD 实用程序将数据移动到表中时，不会出现这种情况。在 z/OS 中，您还可以执行大批数据移动，例如，在 DB2 单机实用程序的帮助下使用 DFSMSdss 实用程序移动整个 DB2 子系统或一组大型表空间。对于 LUW 中这么大规模的数据移动，您可以使用 BACKUP 和 RESTORE 命令，以及重定向的还原方法或 _db2move_ 工具。

##### 用于执行 DB2 实用程序的典型 z/OS JCL

```
//*Job Statement line-1
//*Job Statement line-2
//* It will run against DB2 sub-system DB2T
//STEP01   EXEC PGM=DSNUTILB,PARM='DB2T'
//STEPLIB  DD DISP=SHR,DSN=DB2.V9R1.DB2T.SDSNEXIT
//         DD DISP=SHR,DSN=DB2.V9R1.SDSNLOAD
//SYSPRINT DD SYSOUT=*
//SYSOUT   DD SYSOUT=*
//SYSREC   DD DISP=SHR,DSN=DDS1923.TABLE.DATA
//SYSIN    DD *
    LOAD DATA RESUME NO LOG NO
    INTO TABLE DEPT
/*

```

Show moreShow more icon

### 在运行 LOAD 之前了解您的环境！

##### Load 考虑事项

- 在 DB2 z/OS 中，表空间中有多少表？
- 在 DB2 LUW 中，规划 INTEGRITY PENDING
- 了解进度，使用 DB2 LUW 中的 `LOAD QUERY` 和 DB2 z/OS 中的 `-DIS UTILITY(util-id)`

尽管 DB2 z/OS 和 DB2 LUW 上都有 LOAD 实用程序，但您会发现它们之间存在显著差异。在运行 REPLACE 函数时，DB2 z/OS LOAD 实用程序会清空整个表空间。因此在运行这样的函数之前，您可能需要检查一下底层表空间中是否包含多个表。DB2 LUW 不支持加载一个或多个分区同时维护其他分区的可用性。另一方面，在 DB2 z/OS 中，您可以运行 LOAD 实用程序来影响一组分区，且如果记录落入分区范围内，DB2 会将数据加载到这些分区中。DB2 LUW LOAD 实用程序不对受影响的数据强制实施关联完整性 (RI)。相反地，它会使受影响的表处于 INTEGRITY PENDING 状态。这个受限的状态相当于您在 DB2 z/OS 中遇到的 CHECK PENDING 状态。

DB2 z/OS 默认情况下会强制实施 RI，也可以让您绕过 RI 检查，并且 DB2 只有在这种情况下将您的表空间置于检查待定状态。如果您的表有一个或多个检查约束条件，DB2 z/OS LOAD 实用程序始终会检查这样的约束并拒绝违反它们的行。DB2 LUW LOAD 实用程序不会检查约束。相反地，它将所有行加载到表中，并使表处于 INTEGRITY PENDING 状态。

在使用 DB2 z/OS 时，您可以使用 UNLOAD 实用程序。IBM 还提供一个名为 DSNTIAUL 的卸载示例程序。此示例程序可以为您提供 UNLOAD 实用程序不支持的其他 SQL 功能（比如通过联接表获得的数据）。但是，UNLOAD 实用程序可以让您从之前做的备份副本上获取数据，这是一个将有助于您获得历史数据的功能。在执行 UNLOAD 时，您还可以生成相应的 LOAD 控制语句，并使用它加载数据到任何 DB2 子系统中的相同或类似的表中。

### 如何检查数据完整性？

清单 9 显示如何查询一个表空间的状态，从而了解 CHECK PENDING 的状态。要检查表数据的完整性，您可以在表的数据行设置索引键，并将表的数据行关联到相关的 LOB。在使用 DB2 z/OS 时，您可以运行 CHECK 实用程序；而在使用 DB2 LUW 时，您可以运行 SET INTEGRITY 语句和/或 INSPECT 语句，如清单 10 和清单 11 所示。

##### DB2 z/OS 中的一个 Display 表空间命令上报表的 CHECK pending 状态

```
-DIS DB(TESTDB) SPACENAM(TESTTS)
DSNT360I  DB1S ***********************************
DSNT361I  DB1S *  DISPLAY DATABASE SUMMARY
               *    GLOBAL
DSNT360I  DB1S ***********************************
DSNT362I  DB1S     DATABASE = TESTDB  STATUS = RW
                  DBD LENGTH = 16142
DSNT397I  DB1S
NAME     TYPE PART  STATUS            PHYERRLO PHYERRHI CATALOG  PIECE
-------- ---- ----- ----------------- -------- -------- -------- -----
TESTTS   TS    0001 RW,CHKP
******* DISPLAY OF DATABASE TESTDB  ENDED      **********************

```

Show moreShow more icon

##### DB2 z/OS 的 CHECK 实用程序示例

```
-- Check violations of check constraints and RI
CHECK DATA dbname.tbspname SCOPE REFONLY

-- Check violations and if any found, move them to exception tables
-- and delete from original table
CHECK DATA dbname.tbspname
FOR EXCEPTION in tabowner.tabname USE tabowner.excep_tabname
DELETE YES

-- Check all dependant table spaces (containing all kinds of dependant tables)
-- for violations of RI, Check constraints, LOB values and XML values
CHECK DATA dbname.tbspname SCOPE ALL

-- Check invalid LOB values and LOB references
CHECK LOB TABLESPACE dbname.lobtbsp

-- Check index to table data integrity with concurrent read/write access
CHECK INDEX tabowner.indexname SHRLEVEL CHANGE

-- Check all the indexes on all the tables under table space 'tbspname'
-- of the database 'dbname'
CHECK INDEX (ALL) TABLESPACE dbname.tbspname

-- Finally, reset the Check-pending w/o running the integrity check
REPAIR SET TABLESPACE dbname.tbspname NOCHECKPEND

```

Show moreShow more icon

##### 在 DB2 LUW 中检查完整性

```
-- A Sample query to query integrity status and access modes of individual tables
-- 'U' and 'Y' denotes integrity checked by user or system respectively
-- 'N' denotes particular constraint not checked yet
$ db2 "select substr(creator,1,12) creator, substr(name,1,30) name,
> case when access_mode='R' then 'Read Only'
> when access_mode='F' then 'Full access'
> when access_mode='N' then 'No Access'
> when access_mode='D' then 'No Data Movement'
> end as "Access mode",
> substr(const_checked,1,1) FK, substr(const_checked,2,1) CC
> from sysibm.systables order by 4,5"

CREATOR      NAME                           Access mode      FK CC
------------ ------------------------------ ---------------- -- --
DB2INST1     EMPPROJACT                     No Access        N  Y
DB2INST1     ADEFUSR                        No Access        U  U
DB2INST1     EMPLOYEE                       No Access        N  N
DB2INST1     DEPARTMENT                     No Access        N  Y
DB2INST1     CL_SCHED                       Full access      Y  Y
DB2INST1     DEPT                           Full access      Y  Y

-- SET INTEGRITY will not run if the table is not in INTEGRITY PENDING,
-- So, you can set INTEGRITY PENDING on a particular table alone by:
$ db2 "SET INTEGRITY FOR db2inst1.employee OFF CASCADE DEFERRED"

-- Put table in INTEGRITY PENDING with read access to the data
$ db2 "SET INTEGRITY FOR db2inst1.customer OFF READ ACCESS CASCADE DEFERRED"

-- Check violations of check constraints and referential integrity
$ db2 "SET INTEGRITY FOR db2inst1.employee IMMEDIATE CHECKED"

-- Check violations and move violating rows off to exception tables
$ db2 "SET INTEGRITY FOR db2inst1.employee IMMEDIATE CHECKED
> FOR EXCEPTION in db2inst1.employee USE db2inst1.emp_exception"

-- Reset Integrity pending for RI violations without actually checking
$ db2 "SET INTEGRITY FOR db2inst1.employee FOREIGN KEY IMMEDIATE UNCHECKED"

```

Show moreShow more icon

## 备份和恢复方法

RECOVER Utility/ 命令在 DB2 LUW 和 DB2 z/OS 中均可用。RECOVER 会恢复最近的备份，然后执行一个日志应用进程，将表空间或数据库带到一个指定的数据点。您应当查询 DB2 目录表 SYSIBM.SYSCOPY 了解备份细节（位置、RBA、备份时间），这类似于您在 DB2 LUW 中使用 `LIST HISTORY` 命令读取存储在 DB2RHIST.ASC 文件中的与恢复相关的信息。

在 DB2 LUW 中，您还可以使用 `ROLLFORWARD DATABASE` 命令来应用恢复日志中的数据变更，这通常发生在还原一个联机副本之后（使用 `RESTORE DATABASE` 命令）。DB2 z/OS 中的 RECOVER _LOGONLY_ 实用程序执行类一些似于 ROLLFORWARD 的命令。

在使用 DB2 z/OS 时，您可以利用以下两个实用程序在不同级别备份数据。

- BACKUP SYSTEM 实用程序：

    - 使用 z/OS SMS 的 DFSMShsm 组件的整个 DB2 子系统（有/没有日志）
- COPY 实用程序：

    - 完整的表空间（非分区的或分区表的所有分区）
    - 分区表的一组分区
    - 为 COPY 启用的索引空间
    - 分区索引的一组分区
    - 非分区表空间的单个数据集
    - 非分区索引的单个索引

对于 DB2 LUW，您可以备份的内容取决于数据库的日志模式。如果数据库处于循环日志记录模式，只允许脱机备份数据库。在数据库处于循环日志记录模式时，单个表空间的备份不受支持。您可以在数据库上启用归档日志记录模式，并开始联机备份数据库以及备份数据库中的一组表空间。清单 12 展示了 DB2 LUW 中的联机数据库备份和恢复方法。

##### 在 DB2 LUW 中使用恢复相关的命令

```
-- Have archiving enabled and perform an online backup.
$ db2 "backup db sample online"

Backup successful. The timestamp for this backup image is : 20110728153122

-- Perform a recovery; DB2 will automatically pickup
-- most recent backup it should restore to.
$ db2 "recover db sample to 2011-07-28-15.33.00.000000"

                                 Rollforward Status

Input database alias                   = sample
Number of nodes have returned status   = 1

Node number                            = 0
Rollforward status                     = not pending
Next log file to be read               =
Log files processed                    = S0000024.LOG - S0000025.LOG
Last committed transaction             = 2011-07-28-15.31.32.000000 Local

DB20000I  The RECOVER DATABASE command completed successfully.

-- Also you can do a RESTORE from a backup copy followed by ROLLFORWARD
$ db2 "restore db sample taken at 20110728153122 without prompting"
SQL2540W  Restore is successful, however a warning "2539" was encountered
during Database Restore while processing in No Interrupt mode.

$ db2 "restore db sample logs taken at 20110728153122 logtarget /home/db2inst1/logs'
$ ls -l /home/db2inst1/logs
total 104
-rw-------    1 db2inst1 db2adm1       53248 Jul 28 15:53 S0000024.LOG

$ db2 "rollforward db sample query status using local time"

                                 Rollforward Status

Input database alias                   = sample
Number of nodes have returned status   = 1

Node number                            = 0
Rollforward status                     = DB  pending
Next log file to be read               = S0000024.LOG
Log files processed                    =  -
Last committed transaction             = 2011-07-28-15.31.32.000000 Local

$ db2 "rollforward db sample to 2011-07-28-15.33.00.000000 using local time
> and complete overflow log path (/home/db2inst1/logs)"
                                 Rollforward Status

Input database alias                   = sample
Number of nodes have returned status   = 1

Node number                            = 0
Rollforward status                     = not pending
Next log file to be read               =
Log files processed                    = S0000024.LOG - S0000025.LOG
Last committed transaction             = 2011-07-28-15.31.32.000000 Local

DB20000I  The ROLLFORWARD command completed successfully.

```

Show moreShow more icon

在使用 DB2 z/OS 时，您还可以运行 `REPORT RECOVERY` 实用程序来为恢复情况提前做好准备。在 DB2 z/OS 中恢复一个表时，不会自动构建依赖于受影响表的索引。您应当使用 `REBUILD INDEXES` 实用程序计划和构建这样的索引。

### 清理备份历史记录

作为一名 DBA，业务规则通常会决定您应该保留备份副本的时间长度。在 DB2 z/OS 中，您应当执行 `MODIFY RECOVERY` 来从 SYSIBM.SYSCOPY 目录表总删除表空间的过时条目。该实用程序绝不会删除备份数据集的实体。

- `MODIFY RECOVERY dbname.tsname DELETE DATE(111201)` ：清理 2011 年 12 月 01 日之前的条目。
- `MODIFY RECOVERY dbname.tsname DELETE AGE(60)` ：清理 60 天以上的条目。
- `MODIFY RECOVERY dbname.tsname RETAIN LAST(30)` ：清理条目，但保留最近的 30 个条目。

例如在 DB2 LUW 中， `PRUNE HISTORY 201105` 会删除 2011 年 5 月之前的条目。如果不指定一个完整的时间戳，DB2 会替代月、日值和时间戳的其他字段。通过将 `AND DELETE` 附加到前面的命令，您还可以删除备份副本的实体。

### Quiesce 简介

尽管 Quiesce 在 DB2 z/OS 和 DB2 LUW 中均可使用该命令/实用程序，但对它的使用会有所不同。在 DB2 z/OS 中，您会发现需要建立一致的日志点来执行 Point In Time (PIT) 恢复。在运行 `QUIESCE TABLESPACE dbname.tbspname` 时，DB2 会将所有页面从缓冲区写入磁盘，而且会在 SYSIBM.SYSCOPY 表中记录日志序列号（Relative Byte Address/RBA 或 Log Range Sequence Number/LRSN）。在不需要进行频繁备份的情况下可以这么做，或者在 DB2 之外执行复制活动之前强制发生变更的页面离开缓冲池。您可以在 `RECOVER TABLESPACE` 命令中使用记录的 RBA 来执行 Point In Time (PIT) 恢复。

另一方面，在 DB2 LUW 中，Quiesce 将处于持续锁定模式，可在执行维护任务时将其应用于实例、数据库或表空间来限制普通用户。在维护活动结束之后，您应当发出 `UNQUIESCE INSTANCE db2inst1` 、 `UNQUIESCE DB sampledb` 或 `QUIESCE TABLESPACES FOR TABLE myschema.tabname RESET` 命令来允许普通用户访问数据。

## 数据维护

要在两个平台上保持数据完整无损，则需要使用 REORG 和 RUNSTATS。Real Time Statistics (RTS) 也可同时用于 DB2 z/OS 和 DB2 LUW。当涉及维护索引时，DB2 z/OS 有以下两种不同的实用程序：REORG INDEX 和 REBUILD INDEX。在 DB2 LUW 中，REORG INDEX 实用程序可以同时执行两个函数。由于 DB2 z/OS 实用程序在分区表空间的表空间级或分区级运行，所以它们会影响表空间中存在的所有表。因此在 DB2 z/OS 中工作时，您可以重组多个表，或一次重建多个表的索引。表 1 列出并比较了这些平台中可用的数据维护方法。

### 运行 REORG 的不同要求和方式

目标在 DB2 z/OS 中在 DB2 LUW 中重组一个非分区索引REORG INDEX creator.ixnameREORG INDEX schema.npi\_index CLEANUP ONLY重建一个非分区索引REBUILD INDEX (creator.ixname)REORG INDEX schema.npi\_index重建单个表的表空间中的所有索引REBUILD INDEX (ALL) TABLESPACE dbname.tbspaceREORG INDEXES ALL FOR TABLE myschema.tabname重组分区索引的一个分区REORG INDEX creator.ixname PART 10不受支持脱机重组单个表的表空间REORG TABLESPACE creator.ixname SHRLEVEL NONEREORG TABLE ALLOW NO ACCESS联机重组单个表的表空间指定 SHRLEVEL CHANGE 并使用 MAPPINGTABLE指定 INPLACE ALLOW WRITE ACCESS重组一个表的 2 个或多个分区REORG TABLESPACE dbname.tsname PART m:n不受支持联机重组具有至少一个非分区索引的分区表受支持不受支持检查 Reorg（的需要）运行 RUNSTATS 实用程序和 Query Catalog 表或使用 OFFPOSLIMIT m INDREFLIMIT n 和 REPORTONLY 关键字运行 REORG TABLESPACE；使用 LEAFDIST n 和 REPORTONLY 关键字运行 REORG INDEXREORGCHK ON TABLE myschema.tabname

## 监控 DB2

在使用 DB2 LUW 时，如果您希望在任何对象级监控 DB2，就必须在进行任何活动之前启用 Database Monitor (DBM) 开关。您可以运行 `GET DBM M0NITOR SWITCHES` 命令来了解您服务器上的监控级别。

类似地，在 DB2 z/OS 中，会在子系统级运行不同类型的 _跟踪_ ，以收集监控和性能调优通常所需的信息。在设置监控工作时，您可以启动和停止跟踪，也可以在子系统 DSNZPARM 中对其进行定义，让 DB2 自动运行它们。

##### 了解 DB2 z/OS 中的活动踪迹及其目的地

```
-DIS TRACE(*)
DSNW127I  DB1S CURRENT TRACE ACTIVITY IS -
TNO TYPE   CLASS        DEST QUAL IFCID
01  STAT   01,03,04,05, SMF  NO
01         06,08
02  AUDIT  01           SMF  NO
03  ACCTG  01,02,03,07, SMF  NO
03         08
04  MON    01           OP1  NO
05  AUDIT  03           OP1  NO
06  AUDIT  01,02,04,05, OP1  NO   090,091

```

Show moreShow more icon

System Management Facility (SMF) 在先前的清单中被视为一个目的地 (DEST)，它是各种 z/OS 子系统用于编写不同类型的跟踪记录的一个 z/OS 辅助工具。要读取 SMF 记录，您可以使用自定义编写的例程（自主开发或者从供应商处购买）来制作有关特定事件或 DB2 对象的若干批量报表。另外，您可以使用 IBM Tivoli OMEGAMON XE 这样的监控工具（需另行购买）来提供 DB2 子系统或其应用程序/线程的一个快照视图。

##### DB2 LUW 中的其他监控和故障排除辅助工具

- db2pd：用于报告来自内存集的信息，类似于快照监视器产生的那些信息。
- db2mtrk：用于跟踪或报告数据库、实例和应用程序的内存。
- db2top：是可以联机或在批处理模式下运行的一个监控工具（可用于 AIX、Linux 和 Solaris）。

在 DB2 LUW 中，您可以通过一个 _事件监视器_ 和 _快照监视器_ 执行监控。您可以为特定类型的事件定义 Event Monitor，而且根据事件的类型，您可以告诉 DB2 将捕获的信息写入 DB2 表、文件、管道或未格式化的时间表（DB2 9.7 中的新增特性）。

在分别使用表或文件类型的事件监视器时，您还可以使用一些辅助工具，比如 `db2evtbl` 命令或 `db2evmon` 命令。而在 DB2 LUW 中，您还可以在监控数据库的 DB2 实例时使用 SYSIBMADM 模式下的各种管理视图。这些动态视图会为您提供 _快照监视器_ 生成的系统 _快照_ 。您还可以通过 `GET SNAPSHOT` 命令或通过各种快照 _函数_ 收集快照信息。清单 14 显示在 DB2 LUW 中创建和管理事件监视器的一个示例。

##### 在 DB2 LUW 中运行简单的事件监视器

```
-- Create a simple file event monitor for every statement
$ db2 "CREATE EVENT MONITOR file_event_stmt FOR STATEMENTS
> WRITE TO FILE 'stmtevent' MANUALSTART"
-- Event is not started yet because of keyword MANUALSTART
$ db2 "SELECT SUBSTR(EVMONNAME,1,18) EVMONNAME, CASE WHEN
> EVENT_MON_STATE(EVMONNAME)=0 THEN 'STOPPED'
> WHEN EVENT_MON_STATE(EVMONNAME)=1 THEN 'STARTED' END AS STATUS
> FROM SYSCAT.EVENTMONITORS"

EVMONNAME          STATUS
------------------ -------
DB2DETAILDEADLOCK  STARTED
FILE_EVENT_STMT    STOPPED

2 record(s) selected.
-- Make sure target directory is created and start the event
$ db2 "SET EVENT MONITOR FILE_EVENT_STMT STATE=1"
-- Format the DB2 event output
$ db2evmon -path /home/db2inst1/db2inst1/NODE0000/SQL00002/db2event/stmtevent > event_out

```

Show moreShow more icon

## 可用性和可伸缩性

##### DB2 数据共享

- 在 z/OS Parallel sysplex 上运行。
- 所有 DB2 成员位于同一个 sysplex 上。
- DB2 成员共享目录和表。
- 获益包括提高的可用性、工作负载平衡、对渐进式增长的支持和维护/升级的灵活性。
- 支持组访问，具体成员以及单个成员的访问。

IBM 在每次升级 DB2 都会添加一些功能，以提高数据和数据服务器的可用性。DB2 Data Sharing for z/OS 就是这样一个基于 z/OS 的软硬件组件的高可用性解决方案。

数据共享通过形成两个或多个 DB2 子系统来实现，在 z/OS Parallel sysplex 中，这些子系统通常在不同的物理机上运行。单个 DB2 子系统称为 _成员_ ，会向数据请求者呈现单一视图，并共享磁盘上的数据。在连接到这样一个数据共享组时，DB2 会在可用的、有工作负载容量的成员上运行您的工作。一个或多个 z/OS 机器可以充当 _Coupling Facility_ ，您可以在此基础上创建 DB2 组件，比如组缓冲池、一个锁结构或一个共享区，从而有效地管理数据访问。该解决方案还提供一个执行单个 DB2 子系统的迁移和维护的简单方法。扩展 DB2 服务器涉及将一个或多个成员添加到现有数据共享组中。

用于 Linux 和 UNIX 的 DB2 pureScale 特性可与 z/OS 的 DB2 Data Sharing 特性相匹敌。参见 参考资料 部分，了解有关这两个特性的更多比较。注意，该特性在 Enterprise Server Edition 和 Advanced Enterprise Server Edition 上仅可作为单独定价的选项。

扩展数据库服务器的另一种方式称为 Database Partitioning Feature (DPF)，可用它来跨多个分区（节点）分散数据。DB2 会基于一个分布图将整个表数据分布到位于不同物理机上的数据库分区中。与 DB2 LUW 9.7 一样，DPF 不再只是 DB2 版本的一部分，所有 IBM InfoSphere Warehouse 版本中都有提供。

DB2 LUW 还提供以下高可用性功能：

- 自动客户端重新路由。您可以使用 `UPDATE ALTERNATE SERVER FOR DATABASE` 命令在主服务器注册一个备用的 DB2 LUW 服务器。通过 TCP/IP 连接到主 DB2 服务器的客户端如果失去与主服务器的通信，则可以连接到已定义的备用服务器。
- 高可用性灾难恢复 (High Availability Disaster Recovery, HADR)。该解决方案将一个额外的 DB2 服务器用作备用服务器，它在主服务器出现故障时接管其工作。在通过各种数据库配置参数配置了主备数据库并在备用位置恢复了主数据库之后，您就可以启动 HADR 函数了。DB2 服务器会自动将日志变更从主服务器传递给备用服务器。您还可以混合使用自动客户端重新路由技术与 HADR 数据库。

## 结束语

在本文中，我们讨论了在 DB2 for z/OS 和 DB2 for LUW 上执行最常见的数据库管理任务通常使用的方法、工具和命令。由于 DB2 是跨这两个平台使用的，因此具有较高的 SQL 可移植性、相似的目录表，以及类似的功能和数据管理命令。本文可帮助您通过利用现有技能完成更多任务。

## 致谢

作者要感谢他的同事们 Ken Taylor、Camalla Haley Baruwa 和 Ramesh Chejarla 为本文内容提出的宝贵建议。

本文翻译自： [Compare DB2 for z/OS and DB2 for Linux, UNIX, and Windows](https://developer.ibm.com/articles/dm-1108compdb2luwzos/)（2011-11-14）