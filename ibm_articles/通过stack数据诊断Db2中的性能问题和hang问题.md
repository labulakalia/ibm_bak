# 通过 stack 数据诊断 Db2 中的性能问题和 hang 问题
深入了解 stack

**标签:** 分析

[原文链接](https://developer.ibm.com/zh/articles/ba-cn-stack-db2-performance-hang/)

苗庆松

发布: 2017-12-05

* * *

在 Db2 数据库日常运维中，很多的性能问题和数据库 hang 住的问题都可以通过分析 stack 数据来找到 root cause，学会正确地收集和分析 stack 数据是资深 Db2 从业者必备的技能。本文首先介绍了什么是 stack 以及如何收集 Db2 的 stack，之后通过分享三个生产系统中遇到的具体的例子，让读者对分析 stack 有一个直观的认识。

## 什么是 stack？

如果学习过程序设计和数据结构，就知道 stack 是一种数据结构，它遵循着先进后出的原则。函数调用时使用的就是 stack 数据结构，因此函数调用的过程也被叫做调用栈(call stack)，以下面这个简单的 C 语言程序 test.c 为例：

##### 清单 1\. test.c

```
1 #include <stdio.h>
2
3 void fa(void);
4 void fb(void);
5 void fc(void);
6
7 int main(void) {
8 fa();
9 return 0;
10 }
11
12 void fa(void) {
13 fb();
14 }
15
16 void fb(void) {
17 fc();
18 }
19
20 void fc(void) {
21 printf("Last function\n");
22 }

```

Show moreShow more icon

它有四个函数 main()、fa()、fb()、fc()，首先被调用的是 main()函数，main()函数被压入栈中；由于 main()函数会调用 fa()函数，fa()函数也会被压入栈中；之后 fb()、fc()也会被依次压入栈中。当一个函数执行完成后会从栈弹出，所以出栈的顺序正好和入栈顺序相反，依次为 fc()、fb()、fa()、main()。如果想知道某个程序正在运行的是哪个函数，则需要获取程序的 call stack，下面的例子就使用了 gdb 的 backtrace 选项查看清单 1 程序运行 fc()函数时的 call stack：

##### 清单 2\. 使用 gdb 查看 call stack

```
$ gcc -g -o myPro test.c
$ gdb myPro
...
Reading symbols from myPro...done.
(gdb) b 21 <---在第 21 行，即 fc()函数内设置断点
Breakpoint 1 at 0x400557: file test.c, line 21.
(gdb) run <---运行程序
Starting program: /home/qingsong/myPro
Breakpoint 1, fc () at test.c:21
21 printf("Last function\n");
(gdb) backtrace <---查看 call stack
#0 fc () at test.c:21
#1 0x0000000000400551 in fb () at test.c:17
#2 0x0000000000400546 in fa () at test.c:13
#3 0x0000000000400536 in main () at test.c:8

```

Show moreShow more icon

可以看到，完整的 call stack 为 fc() <- fb() <- fa() <- main() 。更一般地，在 linux 上有 pstack 工具、在 AIX 上有 procstack 工具可以查看进程/线程的 stack。

## 收集 Db2 stack 的意义

Db2 日常运维中，可能会遇到应用性能问题、甚至 hang 住。由于 Db2 是线程模型，每个应用都有对应的 EDU 线程，这时候如果要想知道线程正在运行的是哪个函数，就需要收集线程执行函数的 stack。不管是您自行分析，还是求助于 IBM 技术支持团队，stack 都是一项非常重要的数据。Db2 不需要通过操作系统层面的命令收集 stack，db2pd -stack 命令就可以收集到所需 stack。

## 收集 Db2 stack 的方法

### 收集单个应用或者 EDU 的 stack

如果要收集单个应用(对应 EDU)的 stack，可以使用 apphdl 选项，后面接上 application handle，例如下面的命令收集 application handle 为 7 的应用的 stack

##### 清单 3\. 根据 apphdl 收集单个应用的 stack

```
$ db2pd -stack apphdl=7
Attempting to produce dump file and stack trace for coordinator agent, with EDUID 20. See
        current DIAGPATH for dumpfile and stack trace file.

```

Show moreShow more icon

也可以直接跟 EDUID，例如，application handle 7 对应的 EDUID 为 20，那么下面的命令和上一条命令效果一样：

##### 清单 4\. 根据 EDUID 收集单个应用的 stack

```
$ db2pd -stack 20
Attempting to produce stack trace for EDUID 20.
See current DIAGPATH for stack trace file.

```

Show moreShow more icon

收集完成后，在 db2dump 目录下会生成一个 2525.20.000.stack.txt 文本文件，其中 2525 是 db2sysc 进程的 PID，20 是应用的 EDUID。Linx 和 Unix 环境下，stack 文件的命名规则为 pid.eduid.node.stack.txt （注：信息中心的说法 pid.tid.node.stack.txt 并不准确），在 windows 环境下生成的是 pid.tid.stack.bin 文件，需要格式化。stack 文件中和之间部分的就是该 EDU 的 Stack。

### 收集所有 EDU 的 stack

有时候应用之间会出现资源等待的现象，所以只收集某一个应用或者几个应用的 stack 数据是不够的，这时候可以使用 “db2pd -stack all” 收集所有应用的 stack 数据，每个 EDU 的 stack 数据都对应一个 _.stack.txt 文件。如果是多分区环境，要收集所有节点上的 stack，可以使用 db2\_all “; db2pd -stack all” 命令。清单 5 是发出”db2pd -stack all”命令之后，db2dump 目录下生成的文件，可以看到，有很多_.stack.txt 文件，每个 EDUID 对应一个：

##### 清单 5\. 收集所有应用的 stack

```
$ db2pd -stack all
Attempting to produce all stack traces for database partition.
See current DIAGPATH for stack trace file.
$ ls $HOME/sqllib/db2dump
2523.2.000.stack.txt 2525.38.000.stack.txt
2525.000.processObj.txt 2525.39.000.stack.txt
2525.1.000.stack.txt 2525.40.000.stack.txt
2525.11.000.stack.txt 2525.41.000.stack.txt
2525.12.000.stack.txt 2525.42.000.stack.txt
2525.14.000.stack.txt 2525.43.000.stack.txt
...
2525.26.000.stack.txt 2525.55.000.stack.txt
2525.27.000.stack.txt 2531.2.000.stack.txt
2525.28.000.stack.txt 2532.2.000.stack.txt
2525.29.000.stack.txt 2533.2.000.stack.txt
2525.30.000.stack.txt 2535.1936533376.000.stack.txt
2525.31.000.stack.txt 2543.139732824094592.000.nonEDU.app_stack.txt
2525.32.000.stack.txt clientrecords
2525.33.000.stack.txt db2diag.log
2525.34.000.stack.txt db2eventlog.000
2525.35.000.stack.txt events
2525.36.000.stack.txt inst105.nfy
2525.37.000.stack.txt stmmlog

```

Show moreShow more icon

### 重复收集 stack

重复收集 stack 至少有两个原因：首先，一次的 stack 数据只能显示收集数据的这个时间点上函数调用 stack，如果多收集几次，并且每次结果都是相同的，那么就有更有充分的理由证明应用一直是卡在某个函数上；其次，某些性能问题出现时间是随机的，如果只收集一次 stack 数据，可能抓取不到问题出现期间的 stack 数据。可以加上 repeat 选项来重复收集，例如命令 “db2pd -stack all -repeat 5 3” 表示每隔 5 秒收集一次数据，一共收集 3 次，后两次收集的数据会追加到第一次生成的文件中。

##### 清单 6\. 重复收集 stack

```
$ db2pd -stack all -repeat 5 3
Attempting to produce all stack traces for database partition.
See current DIAGPATH for stack trace file.
Attempting to produce all stack traces for database partition.
See current DIAGPATH for stack trace file.
Attempting to produce all stack traces for database partition.
See current DIAGPATH for stack trace file.

```

Show moreShow more icon

## 分析 Db2 stack 的方法

可以看到，如果有单一的应用出现性能问题，分析 Stack 很容易，只要找到其对应的 stack 文件就可以分析了。如果是很多应用都有性能问题，或者涉及到多个应用之间的资源冲突，从哪个 stack 文件入手呢？

这里介绍一些思路：

- 首先，IBM 已经提供了一个脚本工具，在 sqllib/pd 目录下，文件名为 analyzestack，不加任何参数直接运行会显示帮助。它能汇总、统计所有的 stack 文件并生成一个报告，有兴趣的读者可以自行研究一下如何使用。
- 第二种思路是根据 latch wait，很多性能/hang 住问题，都有 latch wait 现象，如果应用在等某个 latch，那么要分析 latch holder 对应的 stack。如果 latch 没有 holder，即 holder 为 0，则只需要分析 later waiter 的 stack
- 第三种思路就是根据关键字去搜索，比如应用 connect to db 时 hang 住，那么可以尝试从包含关键字”connect”的 stack 文件入手，下面的命令可以查找包含”connect”关键字的文件

##### 清单 7\. 查找包含关键字的文件

```
$ grep -i "connect" -l *.stack.txt

```

Show moreShow more icon

## 一个 stack 示例

下面的清单 8 显示了一个 Uow-executing 状态的应用的完整的 stack。它有三个部分：Frame, Function 和 Offset，我们只需要关注 Function 部分。stack 中最上面的 Function 表示是最后被调用的函数，这里面有操作系统的调用，也有 Db2 自己的函数，虽然不一定每个函数都能看懂什么含义，但根据里面一些函数很容易看出来是做的 Insert 操作： pwrite <- sqloseekwrite64 <- sqloWriteBlocks <- sqlbSMSDirectWrite <- … sqlbGetPageFromDisk <- .. <- sqldInsertRow <- sqldRowInsert。

##### 清单 8\. insert 操作的 stack 示例

```
<StackTrace>
-------Frame------ ------Function + Offset------
0x00007FE7DD596D23 pwrite + 0x0033
                (/lib/x86_64-linux-gnu/libpthread.so.0)
0x00007FE7D90D528C sqloseekwrite64 + 0x05c8
                (/home/inst97/sqllib/lib64/libdb2e.so.1)
0x00007FE7D90D4A16 sqloWriteBlocks + 0x00cc
                (/home/inst97/sqllib/lib64/libdb2e.so.1)
0x00007FE7D8BC3C90 _Z15sqlbWriteBlocksP16SqlbOpenFileInfoPvlmjPmP12SQLB_GLOBALS + 0x0060
                (/home/inst97/sqllib/lib64/libdb2e.so.1)
0x00007FE7D9880181 _Z18sqlbSMSDirectWriteP20SQLB_DIRECT_WRITE_CB + 0x031b
                (/home/inst97/sqllib/lib64/libdb2e.so.1)
0x00007FE7D98D7DDD _Z15sqlbDirectWriteP20SQLB_DIRECT_WRITE_CB + 0x013f
                (/home/inst97/sqllib/lib64/libdb2e.so.1)
0x00007FE7D8B27E96 _Z15sqlbDirectWriteP16SQLB_OBJECT_DESCt14SQLB_AS_SCHEMEljPctjP12SQLB_GL
OBALS + 0x0082
                (/home/inst97/sqllib/lib64/libdb2e.so.1)
0x00007FE7D987EA32 sqlbAppendPagesForObject + 0x01b4
                (/home/inst97/sqllib/lib64/libdb2e.so.1)
0x00007FE7D8BBCB57 address: 0x00007FE7D8BBCB57 ; dladdress: 0x00007FE7D7DCC000 ; offset in
lib: 0x0000000000DF0B57 ;
                (/home/inst97/sqllib/lib64/libdb2e.so.1)
0x00007FE7DB40D059 _Z19sqlbGetPageFromDiskP11SQLB_FIX_CBi + 0x14f5
                (/home/inst97/sqllib/lib64/libdb2e.so.1)
0x00007FE7DB3FDC0C _Z7sqlbfixP11SQLB_FIX_CB + 0x10ae
                (/home/inst97/sqllib/lib64/libdb2e.so.1)
0x00007FE7D8BFDFFE _Z18sqlbFixNextNewPageP11SQLB_FIX_CBbjj + 0x0076
                (/home/inst97/sqllib/lib64/libdb2e.so.1)
0x00007FE7D99444DF _Z17sqldCreateNewPageP13SQLD_DFM_WORKP11SQLB_FIX_CB + 0x005b
                (/home/inst97/sqllib/lib64/libdb2e.so.1)
0x00007FE7D9943FB2 _Z15sqldExtendTableP13SQLD_DFM_WORKi + 0x01aa
                (/home/inst97/sqllib/lib64/libdb2e.so.1)
0x00007FE7D9942BB0 address: 0x00007FE7D9942BB0 ; dladdress: 0x00007FE7D7DCC000 ; offset in
lib: 0x0000000001B76BB0 ;
                (/home/inst97/sqllib/lib64/libdb2e.so.1)
0x00007FE7D8BFB188 address: 0x00007FE7D8BFB188 ; dladdress: 0x00007FE7D7DCC000 ; offset in
lib: 0x0000000000E2F188 ;
                (/home/inst97/sqllib/lib64/libdb2e.so.1)
0x00007FE7D8BF9CC8 _Z13sqldInsertRowP13SQLD_DFM_WORK + 0x00e8
                (/home/inst97/sqllib/lib64/libdb2e.so.1)
0x00007FE7D8C13998 _Z13sqldRowInsertP8sqeAgenttthmiPP10SQLD_VALUEP13SQLD_TDATARECP8SQLZ_RI
D + 0x0768
                (/home/inst97/sqllib/lib64/libdb2e.so.1)
...
</StackTrace>

```

Show moreShow more icon

有一点需要注意的是在 Db2 收集的 stack 中，有的时候前几行函数并不是来源于应用线程的逻辑，而是线程收到 signal 后去执行 signal handler 压到堆栈中的，比如用 db2pd -stack 生成的堆栈经常看到前几个函数是 ossDumpStackTrace, OSSTrapFile, sqlo\_trce，这几个函数就是收到 db2pd 命令发过去的 signal 之后调用的，在实际分析问题时，要注意忽略。另外由于在分析 stack 时 offset 并非关注点，在文章后面的几个例子中，我们将给出的 stack 清单做了精简，省略了对分析问题没有帮助的部分，只保留了 Frame 和 Function 部分。

## 通过分析 stack 解决性能问题

这一部分举了两个生产系统中遇到的性能问题，通过分析 stack 找到了问题产生的根源，具有一定的代表性。

### load 性能问题

首先来看一个 load 性能问题：DBA 发现最近一段时间所有的 load 命令明显变慢，即使只 load 一行数据到一张空表，也需要超过两分钟的时间，于是在这期间收集了三次 load 应用的 stack 数据，三次 StackTrace 完全一样，如清单 9 所示：

##### 清单 9\. load 应用的 stack

```
<StackTrace>
-------Frame------ ------Function + Offset------
0x0900000000029F60 lseek64
0x090000000ACD3A0C sqloseek
0x090000000D0039C0 sqluhReadEntry
0x090000000D0037BC sqluhUpdate
0x090000000CC39ABC sqluCommitLoadEntryInHistoryFile
0x090000000CC316B0 sqlu_process_pending_operation
...
0x090000000E5F341C sqlrr_commit
0x090000000CC8FAA8 sqluv_commit
0x090000000D093440 sqlu_register_table_load
...
</StackTrace>

```

Show moreShow more icon

可以看到 load 操作卡在 lseek64 上，这个函数作用是重新定位读/写文件偏移量。往前看，有一个函数 sqluhUpdate，表示要更新文件，具体更新的是哪一个文件呢？再往前看，发现 sqluCommitLoadEntryInHistoryFile 函数，看到这里问题就明确了：由于历史文件过大，导致 load 完成之后，需要花很多时间更新历史文件。解决的方法也很简单：使用 prune history 命令清理了历史文件之后，load 恢复正常。

### 取 sequence 慢问题

再来看另外一个性能问题：Db2 数据库上有很多应用同时运行，每隔两个小时，这些应用都会有一次明显的性能下降，正常情况下只需要几秒钟就能完成，但出问题的时候却需要 1-3 分钟左右，之后就恢复正常。问题期间，有大量的 latch wait 现象，于是在问题期间收集了几次 latch 数据、所有应用的 stack 和 application snapshot。因为有 latch wait 现象，所以首先查看了 latch 数据， db2pd -latch 的输出可以在清单 10 中看到，其中显示大量的 latch wait on SQLO\_LT\_SQLD\_CHAIN **fastChainLatch 和 SQLO\_LT\_SQLD\_SEQ** seqLatch，这俩 latch 对应文件名都是 sqldsequence.C，看起来和 sequence 有关系：

##### 清单 10\. 等 latch 应用的 stack

```
Latches:
Address Holder Waiter Filename LOC LatchType
No latch holders.
0x00002AAFC2C3FD00 0 38 sqldsequence.C 375 SQLO_LT_SQLD_CHAIN__fastChainLatch
0x00002AAFC2C3FD00 0 437 sqldsequence.C 375 SQLO_LT_SQLD_CHAIN__fastChainLatch
0x00002AAFC2C3FD00 0 442 sqldsequence.C 375 SQLO_LT_SQLD_CHAIN__fastChainLatch
0x00002AAFC2C3FD00 0 454 sqldsequence.C 375 SQLO_LT_SQLD_CHAIN__fastChainLatch
0x00002AAFC2C3FD00 0 458 sqldsequence.C 375 SQLO_LT_SQLD_CHAIN__fastChainLatch
0x00002AAFC2F072E8 0 460 sqldsequence.C 388 SQLO_LT_SQLD_SEQ__seqLatch
...<skip>..

```

Show moreShow more icon

由于这些 latch 没有 holder，所以只需要分析它本身的内容，找到 EDU 对应的 stack，发现这些应用对应 EDU 的 stack 内容都一样的，如清单 11 所示：

##### 清单 11\. 等 latch 的应用的 stack

```
<StackTrace>
...
0x0000003532ABB5A7 __sched_yield
0x00002AAAABF85F55 sqloSpinLockConflict
0x00002AAAADE27F19 address: 0x00002AAAADE27F19 ;
0x00002AAAAC576441 address: 0x00002AAAAC576441 ;
0x00002AAAAC578FD6 sqldSeqGenerate
0x00002AAAAD62B79B sqlriSeqNextval
...
</StackTrace>

```

Show moreShow more icon

可以从 stack 里看到，有两个函数：sqldSeqGenerate 、SeqNextval，字面意思是在产生 Sequence 的 Next Value。实际上应用是不是在做这件事呢？于是查看应用的快照，发现这些应用正在执行的 SQL 语句果然都是在取 sequence 的 NEXTVAL： VALUES NEXTVAL FOR SEQ\_NAME。解决方法是将该 sequence 的 cache 从 1000 扩大到 5000，扩大之后问题得到解决。

## 通过分析 stack 诊断 hang 住问题

再来看一个连库 hang 住的问题：问题现象为某 Db2 系统会随机出现应用连不上数据库的情况（连库一直 hang 住），过一段时间自己就好了，在 Db2 server 上连库也会 hang 住，初步排除网络原因。问题期间观察到有很多 latch wait 现象，于是收集了 db2pd -latch 和 db2pd -stack all 的数据。数据分析过程如下：

第一步，因为是连库 hang 住，先在所有的 stack.txt 里搜索关键字 agentConnect，发现有几十个应用的 stack 都如清单 12 所示：

##### 清单 12\. 连库应用的 stack

```
<StackTrace>
-------Frame------ ------Function + Offset------
0x09000000078484FC sqloXlatchConflict
0x0900000007A56F48 StartUsingLocalDatabase
0x0900000007A53510 AppStartUsing
0x0900000007A8B460 AppLocalStart
0x0900000007994828 sqlelostWrp
0x09000000079948D8 sqleUCengnInit
0x0900000007A94D80 sqleUCagentConnect
0x0900000007A96294 sqljsConnectAttach
0x0900000007A958E4 sqljs_ddm_accsec
0x0900000007A9550C sqljsParseConnect
...
</StackTrace>
<LatchInformation>
Waiting on latch type: (SQLO_LT_sqeLocalDatabase__dblatch) - Address: (7800000002eb208),
        Line: 1013, File: sqle_database_services.C
</LatchInformation>

```

Show moreShow more icon

有几个关键的函数 sqloXlatchConflict <- StartUsingLocalDatabase <- sqleUCagentConnect，看起来意思是尝试连库的时候遇到 latch conflict，具体等的 latch 名为 SQLO\_LT\_sqeLocalDatabase\_\_dblatch，下一步需要看一下 hold 住这个 latch 的应用在做什么事。

第二步，通过 db2pd -latch 的输出，查找到 hold 住 SQLO\_LT\_sqeLocalDatabase\_\_dblatch 的应用对应的 EDUID 为 200876，于是查看 xxxxxxxx.200876.000.stack.txt 文件，stacktrace 如下：

##### 清单 13\. latch holder 的 stack

```
xxxxxxxx.200876.000.stack.txt
<StackTrace>
-------Frame------ ------Function + Offset------
0x0900000000007758 stat64
0x0900000007983800 sqlofindn2
0x0900000007B8C20C next
0x0900000009283660 sqlbMonSMSGetContainerPages
0x0900000009564370 sqlbSnapshotTablespace
0x09000000092ABE98 write_sqlm_ts
0x09000000092AB838 snap_tspaces_for_db
0x0900000009512DB0 sqlmonssagnt
...
</StackTrace>

```

Show moreShow more icon

可以通过 stack 看到，有两个函数 MonSMSGetContainerPages <- sqlbSnapshotTablespace，意思是正在获取表空间的快照，再具体一点，是获取 SMS 表空间容器(container)的数据页(pages)信息。由于只有系统临时表空间是 SMS 的，所以，根据以上信息可以推断出连不上库的原因是临时表空间页数过多，导致 snapshot 获取它的容器信息时间过长，没有释放连库所需要的 latch。

第三步，为了验证上面的说法，再看一下其他应用的 stack 信息，使用 grep 按照 Temp 关键字搜索所有的\*.stack.txt，发现有很多应用的 stack 如下：

##### 清单 14\. 包含关键字”temp”的应用的 stack

```
<StackTrace>
-------Frame------ ------Function + Offset------
0x090000000002A678 pwrite64
0x09000000078AFC7C sqloseekwrite64
0x09000000078AFABC sqloWriteBlocks
0x09000000078AE83C sqlbWritePageToDisk
0x09000000078AE60C sqlbWritePage
0x09000000078D0728 sqlbCheckVictimSlot
0x09000000078CEDF4 sqlbFreeUpSlot
0x09000000078CF1A0 sqlbGetVictimSlot
0x09000000078D3970 sqlbGetPageFromDisk
0x09000000078D637C sqlbfix
0x09000000078D6A68 sqlbFixPage
0x09000000078B5664 sqlbFixNextNewPage
0x0900000007972D58 sqldAppendTempRow
0x0900000007972820 sqldInsertTemp
0x09000000078868C4 sqldRowInsert
0x090000000A3016AC sqlri_hsjnFlushSinglePartition
0x090000000A300188 sqlri_hsjnFlushBlocks
0x09000000099F5DB8 sqlri_hsjnPartitionFull
0x090000000788F824 sqlri_hsjnNewTuple
0x090000000788FF64 sqlrihsjnpd
0x09000000078FD2EC sqldEvalDataPred
0x09000000078FA970 sqldReadNorm
0x09000000078A0248 sqldFetchNext
0x09000000078A0938 sqldFetchNext
0x09000000078A0B64 sqldfrd
0x09000000078EB400 sqldRowFetch
...
</StackTrace>

```

Show moreShow more icon

几个关键的函数为：pwrite64 <- sqloseekwrite64 <- sqloWriteBlocks <- sqlbWritePageToDisk <- .. sqlbGetPageFromDisk <- .. sqldAppendTempRow <- sqldInsertTemp <- sqldRowInsert <- sqlrihsjnpd。咦，是不是看起来比较面熟？没错，几乎和本文第一个 insert 例子中的 stack 完全一样，不过这次 insert 的是 temp 数据，而且看起来像是 hash join 造成的。在后续的监控中，果然发现出问题时临时表空间使用量很高，并且在快速增长。解决方案为找到应用对应的 SQL，并对其进行调优以避免使用临时表空间，问题消失。

综上：通过对多个 stack 数据一层层地分析，找到了问题的根本原因：未经优化的 SQL 语句使用了大量临时表空间 -> 运行监控函数时间长，一直持有 latch -> 新的应用不能连接。

## 结束语

在使用 stack 分析问题的时候，我们并不需要知道每个函数的作用，往往根据一些关键的函数就能推断出导致问题的原因，这里就需要一定的经验和积累。

另外，这篇文章的重点不是如何去分析性能/hang 问题，而是展示 stack 数据在分析性能/hang 问题时所扮演的重要角色，因为 stack 虽然强大，但并非所有的性能问题都能通过 stack 找到原因，性能问题往往要综合系统监控函数/视图、快照信息、执行计划甚至 Db2 trace 才能查明原因，而 hang 问题很多是 Db2 产品的 BUG，这时候就需要将收集的 stack 以及其他数据提交至 IBM 技术支持中心来分析。

## 参考资料

文章 [What data should I collect for a PERFORMANCE issue with DB2 For Linux, Unix and Windows?](http://www-01.ibm.com/support/docview.wss?uid=swg21700589) 讲述了在出现性能问题时，应该如何收集数据。

信息中心链接 [Troubleshooting DB2 application hangs](https://www.ibm.com/support/knowledgecenter/en/SSEPGG_9.7.0/com.ibm.db2.luw.admin.trb.doc/doc/t0061155.html) 讲述了如何 troubleshooting Db2 应用程序 hang 住问题，其中包括了 db2\_hang\_analyze 工具的使用。

文章 [DB2 系统临时表空间过大引发的性能问题](https://www.ibm.com/developerworks/cn/data/library/ba-cn-db2-tempspace-oversize-performance/index.html) 中有一个类似的问题，也是由于临时表空间过大引起的性能问题，和本文中最后一个例子略有不同，不过在关键的数据分析上是类似的。

文章 [Where does db2 allocate sequence cache? Does db2 allocate more memory if the cache size increases?](http://www-01.ibm.com/support/docview.wss?uid=swg21997112) 解释了 sequence cache 在内存中的分配机制。