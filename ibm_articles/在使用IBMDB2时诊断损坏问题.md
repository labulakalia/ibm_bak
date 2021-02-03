# 在使用 IBM DB2 时诊断损坏问题
了解和排除故障

**标签:** 数据库

[原文链接](https://developer.ibm.com/zh/articles/dm-1208corruptiondb2/)

Amitkumar Bamane

发布: 2012-10-15

* * *

免费下载： IBM® DB2® Express-C 10.1 免费版 或者 DB2® 10.1 for Linux®, UNIX®, and Windows® 试用版下载更多的 [IBM 软件试用版](http://www.ibm.com/developerworks/cn/downloads/) ，并加入 [IBM 软件下载与技术交流群组](https://www.ibm.com/developerworks/mydeveloperworks/groups/service/html/communityview?communityUuid=38997896-bb16-451a-aa97-189a27a3cd5a/?lang=zh) ，参与在线交流。

## 简介

被视为是最麻烦的业务问题之一，数据库损坏常常在不知不觉中逐渐形成，给企业带来不利影响。简言之，可以将 _损坏_ 定义为数据库中的任何意外项。损坏问题可能会对系统造成严重的性能冲击。在某些情况下，它可能会导致频繁的系统崩溃，引发关键业务系统宕机。数据库损坏可发生在任何层面，从 DB2 到操作系统以及硬件层。因此，了解和排除故障很重要，即分析所有可能受影响的层，并收集可能尽快需要的任何可用的诊断数据。

在本文中，您将了解为何数据库会在遇到损坏问题时离线。您还将学习分析损坏症状，区分易于修复的故障和灾难性故障。本文将阐明使用 IBM DB2 时的损坏问题，并帮助 DB2 用户理解和选择处理这种关键的高影响问题的最佳方法。

本文首先讨论可能的损坏来源，然后解释以下任务：

1. 识别和排除损坏故障，在使用 DB2 时识别数据库中的损坏问题并对其进行分类，辅以 db2diag.log 中出现的样例症状消息。损坏问题可以大体分为五个类别：数据页面损坏（或表损坏）、索引损坏、CBIT 损坏、日志损坏和压缩描述符损坏。
2. 使用 `db2dart` 和 `INSPECT` 识别损坏问题，洞悉有用的 DB2 命令， `db2dart` 和 `INSPECT` ，来检查数据库损坏。
3. 从损坏中恢复的方法，一旦识别到一个损坏问题，如何着手处理这些情况、要收集什么数据、如何从该状况中恢复过来，这些至关重要。学习可能的恢复方法以及如何选择可用方案。
4. 避免可能的损坏的预防性战略，讨论最佳实践。

## 来源

数据库损坏可能在写入、读取、存储、传输或处理过程中发生，这会向原始数据引入非计划中的更改。损坏问题的一些常见原因：

1. 损坏的文件系统是数据库中出现损坏的最常见原因之一。突然的系统关闭、电涌、文件系统双机挂载、迁移磁盘、文件系统级活动，比如数据库上线运行时检查和修复文件系统（使用的实用程序包括 Linux® 上的 fsck），在文件打开时使用 **Ctrl+Alt+Delete** 以及病毒，都可能在数据库中引入意外的变更。
2. 硬件故障。
3. 内存损坏。
4. DB2 缺陷。
5. I/O 和网络问题（如光纤适配器和交换机中的问题）。
6. 不正确的应用程序编码。
7. 缓冲池 (sqldPage) 和文件系统中存储的页面的值不一致。
8. 重写磁盘数据会导致损坏问题。
9. 用户对数据库的重要配置文件、日志文件、日志控制文件等的干扰都会使数据库处于不一致的状态。

虽说损坏问题由各种原因而致，确切地查明是什么导致了数据损坏是极具挑战的。在大部分情况下，该问题是由文件系统问题和硬件问题引起的。

## 识别和排除故障

对于一个 DBMS， _页面_ 是由操作系统为一个程序执行的内存分配的数据的最小单元，在主内存与任何其他辅助存储（比如硬盘驱动器）之间传输。因此所谓数据库损坏也就是说数据库中的某些页面被损坏了。

如果 DB2 有无法得体处理的错误情况，panic 是它会用来招致崩溃的一种方法。当 DB2 检测到一个页面损坏时，它通过一个受控崩溃 (panic) 停止所有处理，因为它无法确定数据库完整性。这也是为了阻止进一步的数据损害或丢失。

当 DB2 遇到数据库损坏时，db2diag.log 中转储很多错误消息。当出现意外中断且启用了自动的首次出现数据捕获 (FODC) 时，会基于症状来收集数据。当满足以下条件之一时，DB2 9.5 上会自动收集 FODC 数据：

1. FOCD\_Trap，当发生一个实例范围内的陷阱时。
2. FODC\_Panic，当一个 DB2 引擎检测到不连贯且决定不继续时。
3. FODC\_BadPage，当检测到坏页面时。
4. FODC\_DBMarkedBad，当数据库因一个错误而被标记为 “坏” 时。

要搜集必要的信息，比如 OS 诊断（例如，AIX® 上的 `errpt –a` 、 `snap` 和 `fileplace` 输出）以及任何硬件诊断（状态保存和错误日志等），关键是要包含 OS 和硬件支持。重要的是要确保关键的文件系统有足够的磁盘空间，比如转储空间和日志目录，从而确保完全捕获关键事件。

您可以查看 db2diag.log，确认 panic 是因为损坏还是另外的原因引起的。下面您会看到如何识别 DB2 中的损坏并对其进行分类。以下是识别损坏的最常见的一些 db2diag.log 错误消息。

### 数据页面损坏

_数据页面损坏_ 表示一个表中真实数据的损坏。当检测到数据损坏时，大量消息被记录到 db2diag.log。这些消息是识别受影响对象所必需的。清单 1 显示 db2diag.log 中的一个错误消息样例。

##### 清单 1 数据页面损坏的错误消息样例

```
2012-02-04-03.13.05.261082-360 I3442A358          LEVEL: Error
PID     : 393470               TID  : 1           PROC : db2pfchr 0
INSTANCE: inst1                NODE : 000
FUNCTION: DB2 UDB, buffer pool services, sqlbReadAndReleaseBuffers,
probe:13
RETCODE : ZRC=0x86020001=-2046689279=SQLB_BADP "page is bad"
DIA8400C A bad page was encountered.
.
.
.
2012-02-04-03.13.05.264301-360 I3801A437          LEVEL: Error
PID     : 393470               TID  : 1           PROC : db2pfchr 0
INSTANCE: inst1                NODE : 000
FUNCTION: DB2 UDB, buffer pool services, sqlbReadAndReleaseBuffers,
probe:13
DATA #1 : String, 158 bytes
Obj={pool:9;obj:20;type:0} State=x27 Parent={9;20}, EM=1120, PP0=1152
Page=51235 Cont=17 Offset=2848 BlkSize=15
sqlbReadAndReleaseBuffers error: num-pages=16

```

Show moreShow more icon

如上述错误消息所示，DB2 遇到了一个表空间 ID 为 9、表 ID 为 20 的坏页面。存档的对象类型被标记为 0，表明有数据页面损坏。

您可以查询目录表，确定哪个表有损坏页面：

```
db2 "select char(tabname,20), char(tabschema,20) from
                    syscat.tables where tableid=20 and tbspaceid=9"

```

Show moreShow more icon

注意，DB2 只为它尝试访问的那些页面转储坏页面错误。这并不意味着只有那些页面是受损的。您需要借助 `db2dart` 或 `INSPECT` 命令显式地检查数据库中的所有页面，找到损坏范围。

对于临时表中的损坏，类似的错误消息被转储在 db2diag.log 中。如果 db2diag.log 错误消息中的数据类型字段存放一个大于 128 的值，这表明临时表中有损坏。如果对象类型是 3，表明表中存在 LOB 数据损坏。

### 索引损坏

_索引_ 是一个数据库对象，包含一个有序的指针集，指向基表中的行。有多种索引相关的损坏问题，包括：

1. 具有不同 RID 或相同 RID 的惟一索引复制
2. 指向同一 RID 的多个索引条目
3. 摆放不当的索引键（错误的索引键顺序）
4. 行存在，但索引键不存在于任何或一些索引中
5. 指向一个空数据槽或未使用数据槽或 RID 的索引条目是无效的
6. 不正确的上一个或下一个索引页面指针，不正确的高键或索引页面上的其他损坏

以下是索引损坏的错误消息样例。

##### 清单 2 索引损坏的错误消息样例

```
2012-01-30-01.35.50.952434+000 I29308542A2532     LEVEL: Severe
PID     : 1175792              TID  : 33926       PROC : db2sysc 0
INSTANCE: inst1                NODE : 000         DB   : SAMPLE
APPHDL  : 0-7                  APPID: *LOCAL.inst1.120130013528
AUTHID  : TP0ADM
EDUID   : 33926                EDUNAME: db2redow (TP0) 0
FUNCTION: DB2 UDB, buffer pool services, sqlb_verify_page, probe:3
MESSAGE : ZRC=0x86020001=-2046689279=SQLB_BADP "page is bad"
DIA8400C A bad page was encountered.
DATA #1 : String, 64 bytes
Error encountered trying to read a page - information follows :
DATA #2 : String, 23 bytes
Page verification error
DATA #3 : Page ID, PD_TYPE_SQLB_PAGE_ID, 4 bytes
23046981
DATA #4 : Object descriptor, PD_TYPE_SQLB_OBJECT_DESC, 72 bytes
Obj: {pool:9;obj:11076;type:1} Parent={8;11076}

```

Show moreShow more icon

如您所见，错误消息有对象类型：1，这表示索引页面损坏。索引位于 ID 为 9 的表空间中，且索引 ID 为 11076。该索引所在的表的表空间 ID 为 8，表 ID 为 11076。您可以通过查询目录表来检索基表名称和索引名称。

db2diag.log 中的上述代码段表示有一个坏的索引页面。db2diag.log 中表明索引损坏的其他常见错误是 SQLI\_NOKEY，以及无法从索引函数找到行。

### CBIT 损坏

CBIT 是 DB2 使用的一种方法，用于确认从磁盘读入缓冲池的一个页面不是一个部分页面，或者没有从某种形式的损坏改变而来。

CBIT 背后的基本思想是，在写入页面之前将页面上每个扇区（512 字节）的一个位数设置为相同的值。在 DB2 将一个页面刷新到磁盘上之前，计算校验和并记录在页面上。在缓冲池中读回一个页面时，重新计算这一校验和，并根据存储的值对其进行检查。如果有些位数不同，表明有一个部分页面写操作或磁盘损坏。

##### 清单 3 CBIT 损坏的错误消息样例

```
2012-03-12-04.45.17.559235-240 I1104A2616 LEVEL: Severe
PID : 2551866                TID : 1       PROC : db2pfchr
INSTANCE: inst1              NODE : 000
FUNCTION: DB2 UDB, buffer pool services, sqlbVerifyCBITS, probe:1110
MESSAGE : ZRC=0x86020019=-2046689255=SQLB_CSUM "Bad Page, Checksum Error"
DIA8426C A invalid page checksum was found for page "".
DATA #1 : String, 64 bytes
Error encountered trying to read a page - information follows :
DATA #2 : String, 95 bytes
CBIT verification error
bitExpected is 0, userByte is 33, sector 7 (from head of page, 0 based)
DATA #3 : Page ID, PD_TYPE_SQLB_PAGE_ID, 4 bytes

```

Show moreShow more icon

注意，如果损坏不触及存放 CBIT 的字节，在 DB2 以外对页面做出的更改（文件系统故障、磁盘故障等）不会被注意到。CBIT 错误（校验和错误）多半是因初始写或读期间的硬件或 OS 错误引起的。

### 日志损坏

DB2 中的事务日志只是对数据库中发生的所有变更的一个记录。要跟踪事务做出的变更，需要一个方法来对数据变更和日志记录加时间戳。在 DB2 中，这一时间戳机制是使用一个日志序列号 (LSN) 执行的。如果您遇到一个损坏的日志，可能会在 db2diag.log 中有一条类似于清单 4 的错误消息。

##### 清单 4 日志损坏的错误消息样例

```
2010-06-07-12.37.21.143998+120 I8673583A553       LEVEL: Severe
PID     : 2498668              TID  : 27358       PROC : db2sysc 56
INSTANCE: inst1              NODE : 056         DB   : SAMPLE
APPHDL  : 998-22947            APPID: *N998.inst1.100607192315
AUTHID  : LOADDSF
EDUID   : 27358                EDUNAME: db2agntp (SAMPLE) 56
FUNCTION: DB2 UDB, data protection services, sqlpgrlg, probe:291
DATA #1 : < reformatted >
Error -2028994519 when reading LSN 00000B1C261B4FD3 from log file
S0119292.LOG tellMe 1 dpsAcbFlags 1000 setSkipOutputBuf 0

2010-06-07-12.37.21.144202+120 I8674137A487       LEVEL: Severe
PID     : 2498668              TID  : 27358       PROC : db2sysc 56
INSTANCE: inst1              NODE : 056         DB   : SAMPLE
APPHDL  : 998-22947            APPID: *N998.inst1.100607192315
AUTHID  : LOADDSF
EDUID   : 27358                EDUNAME: db2agntp (SAMPLE) 56
FUNCTION: DB2 UDB, data protection services, sqlpgrlg, probe:291
DATA #1 : < preformatted >
HeadLsn 00000B1B8B996EB3, copyLookForLsn 00000B1C261B4FD3

2010-06-07-12.37.21.158065+120 I8675153A549       LEVEL: Error
PID     : 2498668              TID  : 27358       PROC : db2sysc 56
INSTANCE: inst1              NODE : 056         DB   : SAMPLE
APPHDL  : 998-22947            APPID: *N998.inst1.100607192315
AUTHID  : LOADDSF
EDUID   : 27358                EDUNAME: db2agntp (SAMPLE) 56
FUNCTION: DB2 UDB, data protection services, sqlptudo, probe:1010
RETCODE : ZRC=0x87100029=-2028994519=SQLP_BADLSN "Invalid LSN value."
DIA8538C An invalid log sequence number (LSN), the value was "".

```

Show moreShow more icon

如果日志被损坏了，在数据库前滚、崩溃恢复和 HADR 日志回复等需要回应日志的情况下可能会引发严重的问题。数据库前滚能够让您维护数据库中的一致性。它通过应用记录在数据库日志文件中的事务来恢复数据库。前滚是在还原数据库或表空间备份映像之后调用的一个过程。

### 压缩描述符损坏

压缩的描述符是系统目录表中的一个列，DB2 使用它来识别数据库对象的细节。如果它出于某些原因被损坏了，您会在 db2diag.log 中看到如下所示的错误。

##### 清单 5 压缩描述符损坏的错误消息样例

```
2011-08-22-20.50.00.922275-300 I154161182E497      LEVEL: Severe
PID     : 14152                TID  : 184633256288 PROC : db2sysc 0
INSTANCE: inst1             NODE : 000          DB   : SAMPLE
APPHDL  : 0-64465              APPID: 161.166.44.48.27625.11082301341
AUTHID  : DTADM1
EDUID   : 606                  EDUNAME: db2agent (SAMPLE) 0
FUNCTION: DB2 UDB, catcache support, sqlrlc_systables_fetch_from_disk,
probe:200
MESSAGE : Corrupt PD->length in table:DTWF    .JBPM_TASKINSTANCE

2011-08-22-20.50.00.922476-300 I154161680E1588     LEVEL: Severe
PID     : 14152                TID  : 184633256288 PROC : db2sysc 0
INSTANCE: inst1             NODE : 000          DB   : SAMPLE
APPHDL  : 0-64465              APPID: 161.166.44.48.27625.11082301341
AUTHID  : DTADM1
EDUID   : 606                  EDUNAME: db2agent (SAMPLE) 0
FUNCTION: DB2 UDB, catcache support, sqlrlc_systables_fetch_from_disk,
probe:201
MESSAGE : Length in PD=3792, LFD length=10104, DMS length=10104
DATA #1 : String, 11 bytes
Corrupt PD:
DATA #2 : Dumped object of size 3792 bytes at offset 1386248, 53 bytes
/volumes/work/db2dump/xnawmtp5/14152.606.000.dump.bin
DATA #3 : LOB Descriptor, PD_TYPE_LOB_DESCRIPTOR, 60 bytes
SQLDX_LD: Size:60
    x0000        lfd_check                        0x49
    x0001        lfd_version                      6
    x0002        lfd_numsegs                      1
    x0003        lfd_flags                        0x00
    x0004        lfd_size                         10104
    x000C        lfd_life_lsn                     0000009860A6FF2D
    x0014        lfd_mini_numsegs                 0
    x0015        lfd_first                        4
    x0016        lfd_descsize                     60
    x0018        lfd_last_pages                   16
    x001C        lfd_last_bytes                   10104
    x0038        lfd_dir                          Regular Directory
Offsets
        lfd_dir[0]: 33888 (16K)
    Hexdump of LOB descriptor follows:
        4906 0100 0000 0000 7827 0000 0000 0098
        60A6 FF2D 0004 3C00 1000 0000 7827 0000
        0000 0000 0000 0000 0000 0000 0000 0000
        0000 0000 0000 0000 6084 0000

```

Show moreShow more icon

如果您有损坏的 PD，表就变得无法访问了。

## 使用 `db2dart` 和 `INSPECT` 识别损坏问题

`db2dart` 是用于验证数据库及其中对象架构正确与否的一个命令。它还可用于从表（可能会因损坏而无法访问）中提取数据。

绝不应该对仍然拥有有效连接的数据库运行 `db2dart` 。 `db2dart` 通过直接从磁盘中读取来访问数据库中的数据和元数据。因此，如果有连接， `db2dart` 就不会意识到缓冲池中的页面、内存中的控制结构等，因而可能会报告虚假的错误。类似地，如果您对需要崩溃恢复或没有完成前滚恢复的数据库运行 `db2dart` ，那么可能会得到不一致的结果。

要显示所有可能的选项，在没有任何参数的情况下执行 `db2dart` 实用程序。 `db2dart` 的一些选项需要用户输入，比如表空间 ID 和表 ID，如果不在代码行明确指定的话，系统会提示输入这些信息。

### 使用 `db2dart` 检查数据库、表空间和表

`db2dart` 的默认行为是要检查整个数据库。在这种情况下只有数据库名称必须提供。默认情况下， `db2dart` 会创建一个名为 databaseName.RPT 的报告文件。对于单分区数据库环境，文件会被创建到了当前目录中。对于多分区数据库环境，则文件会被创建到了诊断目录的一个子目录中。子目录名为 DARTnnnn，其中 nnnn 是分区号。

如果数据库很大，且您只对一个表空间感兴趣，您可以使用 `/TS` 选项。使用该选项时，您必须在命令行上提供表空间 ID（通过指定 `/TSI` 参数）或者您可以让 `db2dart` 提示您输入该 ID。如果您不知道表空间 ID，可以通过 `DB2 LIST TABLESPACES` 获取。

类似的，可以使用 `/T` 选项检查一个表及其相关对象（LOB 和索引等）。使用该选项时，您必须提供表名称或对象 ID 以及表所在的表空间的 ID。要确定一个表的对象 ID 和表空间 ID，您可以查询 SYSIBM.SYSTABLES 目录表。

### 使用 `db2dart` 转储格式化的表数据

如果一个表空间或表出于任何原因（例如，因坏的磁盘或磁盘控制器）被损坏，通过 SQL 访问表的尝试可能行不通。SQL 语句可能会失败并伴有一个错误，或者数据库可能被标记为 “坏” 且所有连接都会被丢弃。

如果出现这种情况，可能就需要尽量提取所有数据，这样就可以重建表空间和表。在这样的情况下， `db2dart` 的 `/DDEL` 选项可用于提取表数据并将其放到一个分隔的 ASCII 中。下面几节将详细讨论这些选项。

### `INSPECT` 命令

`INSPECT` 命令类似于 `db2dart` 命令。它允许您检查数据库、表空间和表的架构完整性，方法就是检查数据库的页面，看页面是否一致。两个命令之间的显著区别在于，在运行 `db2dart` 之前需要停用数据库，而 `INSPECT` 需要一个数据库连接，且可在同时有多个活跃的数据库连接时运行。

不同于 `db2dart` ， `INSPECT` 命令不能用于格式化和转储数据页面、格式化和转储索引页面、格式化数据行为分隔的 ASCII，并标记索引为无效。因此，它只能检查数据库或其在线对象是否有任何损坏。

## 从损坏中恢复的方法

数据库损坏问题有时候不是很简单明了，需要 IBM 的支持专业知识来选择摆脱这种局面的最佳方式。这里，您会了解如何处理最常见类型的损坏问题并找到最好的恢复方案。当检测到损坏问题时，确定问题来源通常是次要的。您会希望首先纠正当前的问题。

### 从数据页面损坏中恢复

前面已经讨论过，您需要使用 `db2dart` 等工具识别损坏页面。通过查看生成的 `db2dart` 报告文件，您可以估计损坏范围。根据数据损坏的数量和涉及到的复杂度，您需要决定最好的恢复计划。下面是一些恢复选项：

1. 如果有数据库备份，还原数据库并前滚到日志的末尾。如果可行的话这是最简洁的方法，如果数据库很小时首选该方法。
2. 您还可以还原表空间并前滚到日志的末尾。如果损坏是小范围的，这可能是最好的方法。
3. 如果您有为损坏的表重新创建数据的其他方式或者有表数据的副本，删除并重新创建表。您会需要表的 DDL，而且如果您有表数据，您应当能够删除表，使用 DDL 重新创建表，通过可用的任何方法重新创建数据。
4. 如果您没有有效的备份映像以及重新创建表的任何方式，您可以使用 `db2dart` 以及 `/ddel` 来从受损的表上挽救数据。在此之前，您需要使用 `db2look` 查找表的 DDL。您可以使用以下代码获取一个受损表的 DDL：





    ```
    db2look -d <dbname> -e
                            -z <schema_name> -t <table_name> -o <output_file_name>

    ```





    Show moreShow more icon

    下面是使用 `db2dart` 以及 `/ddel` 从受损的表上挽救数据的一个示例： `db2dart <dbname> /ddel` 。该命令需要四个输入值：表对象 ID 或表名称，表空间 ID，起始的页码，以及页数。页数可以是一个具体的值或一个足够大的值（例如 999999999），才能提取表中的所有页面。另外，如果表中的某个页面受到太大损坏，可能需要在页面范围内执行多次 `db2dart /DDEL` 并忽略受损的页面。

    转储的分隔 ASCII 文件在数据代码页进行编码。 `db2dart` 命令不执行代码页面转换。一旦从受损的表上挽救了所有数据之后，您可以检查输出文件 \*.DEL，确保所有数据都存在。完成这个之后，您可以删除受损的表，稍后通过使用 `db2dart` 提取的数据重新创建表。注意， `db2dart /DDEL` 不妨碍 LOB 数据。

5. 如果您有重新创建受损的表空间的方法，您可以使用 `restart database` 将受损的表空间标记为处于删除暂挂状态。稍后您可以重新创建受损的表空间。
6. 如果以上选项都不可行或在提取数据的过程中报出某个错误，在重新创建表的同时删除表等，向 IBM 支持团队寻求帮助。IBM 支持团队可以帮助您删除受损的表，将受损页面初始化为 NULL，等等，视具体情况而定。

### 从索引损坏中恢复

如果在 db2diag.log 和/或 `db2dart` 报告中有索引损坏的迹象，您可以使用 `db2dart` 将索引标记为无效，并去除坏的索引。稍后您可以重新创建索引。

`db2dart` 有一个可将索引标记为无效且使其处于删除暂挂状态的选项。您可以使用 `db2dart /MI` 将受损的索引标记为无效。例如： `db2dart <dbname> /MI /TSI 9 /OI 11076` 。

通过将 `INDEXREC` 参数设置为 restart 和 access 等值，您可以决定何时重新创建一个索引。要在一个应用程序试图访问索引时重新创建索引，您可以将 `INDEXREC` 更新为 access： `db2 update db cfg for <dbname> INDEXREC ACCESS` 。

在设定坏的索引无效并更新 `INDEXREC` 之后，您可以连接到数据库。

如果问题是一个索引 NOKEY 错误，那么在 db2diag.log 文件中会打印大量 `db2dart` 命令。必要时，您可以运行这些 `db2dart <dbname> /di` 命令来转储格式化的索引数据，用于索引损坏的根源分析。您需要使用 UNIX® 上的 `grep` 或 Windows® 上的 **Find** 保存这些命令，将它们保存到一个文件中。编辑文件并用数据库名称替换 DBNAME。如果多次遇到这个问题，可能是有重复的条目。您仅需要保留最近的一组 `db2dart` 命令。要使用 `db2dart` 确认坏的索引，您可以发出 `db2dart <dbname> /t /tsi <tablespace_id> /oi <table_id>` ，其中 `tablespace_id` 和 `table_id` 是定义索引所在的基表的表空间 ID 和对象 ID。

### 从 CBIT 错误中恢复

要修复 CBIT 损坏错误，通过至少在受损表上运行 `db2dart` （最好对整个数据库运行它）来检查问题的严重程度。您可以根据 CBIT 错误是在数据页面上还是索引上来决定上述讨论的方法。从 CBIT 错误中恢复的最可行方法是还原数据库或表空间（如果错误是小范围的话）。

### 从日志损坏中恢复

在日志回复期间日志损坏是个令人担忧的问题。您需要在数据库或表空间前滚、HADR 备用数据库上的日志重播以及崩溃恢复过程中重播日志。如果发生日志损坏，数据库可能没事，只有日志会被损坏。

如果在前滚过程中由于坏的日志而发生错误，可以做的第一件事就是检查 DB2 在为哪个日志文件报告错误。 `db2flsn` 可用于返回包含日志记录的文件的名称，该日志记录根据指定日志序列号 (LSN) 识别。因此如果您在 db2diag.log 中有 _‘bad\_lsn’_ 消息，您可以使用 `db2flsn` 查找相应的日志文件。

如果日志损坏是因为一个缺失的日志文件或一个不正确的日志链上的一个日志文件，则您可以查找正确的日志文件。如果前滚操作由于一个损坏的日志而失败，则您可以选择时间点前滚。前滚操作的指定时间点必须等于或晚于最小恢复时间。最小恢复时间 (MRT) 是当一个数据库一致时前滚过程中最早的时间点。如果您不能将日志至少前滚到 MRT，就需要联系 IBM 支持团队获取协助。另一种方案是从一个离线的数据库备份中还原，而不前滚日志。在这种情况下日志中的事务将不会被应用于数据库。

如果您在崩溃恢复中因一个坏的日志文件而出现问题，您需要还原最新的备份或联系 IBM 支持团队获取协助。

在 HADR 术语中，处理事务的服务器称为 _主服务器_ ，而接收日志并重播它们的合伙服务器称为 _备用服务器_ 。在 HADR 备用数据库上，在日志回复期间，备用服务器可能会因为一个坏的日志而崩溃。您可以检查备用服务器上的 db2diag.log，找到坏的日志文件并尝试从主服务器上传送该日志文件的一个良好备份。一旦有了一个好的日志文件，您就可以启动 HADR：

1. 在备用节点上对 HADR 运行 start： `db2 start hadr on db <dbname> as standby` 。
2. 在主节点上对 HADR 运行 start： `db2 start hadr on db <dbname> as primary` 。

如果上述尝试失败了，您可能需要使用来自主服务器上的全新备份重新配置 HADR，并在备用服务器上还原它。如果您将坏日志文件的一个副本传递给 IBM DB2 支持团队，他们可以检查其中的内容，看看里面有什么内容可能就错误根源提供一些指示。

您可以使用 `db2cklog` 检查归档日志文件的有效性，确定在数据库或表空间的前滚恢复过程中是否可以使用日志文件。可以检查一个归档日志文件或一系列归档日志文件。在使用 `db2cklog` 进行验证的过程中返回错误的日志文件会导致恢复操作失败。

### 从压缩描述符损坏中恢复

您可以使用 db2cat 工具修复损坏的压缩描述符。在修改压缩描述符之前，您需要咨询 IBM 支持团队。可以执行以下步骤来修改压缩的描述符：

1. `export DB2SVCPW=<service_password from IBM support>`
2. `db2cat -d <dbname> -s <schema> -n <tablename> -f <raw pd output file> -o <message file>`
3. `db2cat -d <dbname> -s <schema> -n <tablename> -g <generated pd output file> -o <message file>`
4. `db2cat -d <dbname> -s <schema> -n <tablename> -r <generated pd output file> -o <message file>`

    （使用生成的文件，以 -g 作为替换选项 “-r” 的输入）

5. 从表中导出数据（如有需要）
6. 删除表
7. 重新创建表（如有需要）
8. 导入用户数据（如有需要）

现在您可以在数据库上再次运行 `db2cat` 的验证选项： `db2cat -d <dbname> -s % -n % -v` 。

如果您不需要表，可以删除它们。否则，您应当替换压缩描述符，提取数据，然后删除并重新创建表。

## 避免可能的损坏的预防战略

数据库损坏可能很微妙，难以检测到。因此永远不太可能开发出一个能够检测出可能发生的每一个损坏的工具。我们应当始终避免损坏问题和可能出现损坏风险的过程。

很有必要了解用于避免宕机和数据库崩溃的各种预防措施。这里是帮助在出现系统崩溃之前识别数据库中的损坏的一些最佳实践：

1. 跟踪所有变更。
2. 安装最新的补丁包，如果可能的话，安装最新版本的 DB2 和操作系统（适用的话）。
3. 定期检查文件系统健康状况、网络问题。
4. 尽可能得体地关闭 DB2。
5. 在离线时对数据库运行 `db2dart` ，以检查是否有损坏。如果您禁不起在生产系统上运行 `db2dart` 所需的停机，还原最新的生产备份到测试机器并运行 `db2dart` 。另外，您还可以在数据库在线时使用 `INSPECT` 。它可以及早检测损坏问题或执行主动的健康状况检查。
6. 拥有一个好的备份策略。备份不会检测到一个页面中的损坏问题，因此建议您有一个强大的备份策略和足够的备份生成。
7. 像 RAID 这样的磁盘配置使用冗余磁盘来备份数据，从而有助于最大限度地降低数据损坏。
8. 好的备份电源能够克服由于电涌而引起的损坏问题。
9. 跟踪 IBM DB2 和操作系统中的最新缺陷。

## 结束语

本文讨论了使用 IBM DB2 时最常见的损坏问题，展示了 db2diag.log 中损坏问题的各种症状消息，以及如何基于这些消息识别损坏类型，如何解决这些问题。另外还描述了有助于处理损坏问题的 `db2dart` 和 `INSPECT` 命令。

本文翻译自： [Diagnosing corruption when using IBM DB2](https://developer.ibm.com/articles/dm-1208corruptiondb2/)（2012-10-15）