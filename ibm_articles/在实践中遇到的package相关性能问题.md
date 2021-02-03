# 在实践中遇到的 package 相关性能问题
将 Db2 package 应用到你的工作中

**标签:** 数据管理

[原文链接](https://developer.ibm.com/zh/articles/ba-lo-improve-oltp-performance2/)

郝 庆运

发布: 2018-11-13

* * *

## 引言

本系列的第一部分用了较长的篇幅来较为深入的探索了 Db2 package 相关概念和内部机制，现在来介绍从 Db2 package 层面来优化性能。其实，有了前文的理论基础再来看问题，很多问题都可以迎刃而解甚至做到知其然并且知其所以然。比如说一个简单的事情，有些时候会利用网络抓包工具然后解析 DRDA 协议来抓取应用发起的 SQL 语句并分析 SQL 请求响应时间，以此来监控和分析数据库系统的性能，而且市场上也有类似功能的产品，请问，对于嵌入式静态 SQL，能用网络抓包方式来抓到具体 SQL 语句以及响应时间吗？

接下来介绍几个与 Db2 package 相关的性能问题场景，都是作者工作中曾经遇到过的。

## 高并发场景中迟早会遇到的 xhshlatch 问题

这是前几年发生过的事情，一个重要的渠道类银行业务系统，Java 程序使用 JDBC 方式连接 Db2 服务器，典型的高并发、短交易系统，当时还是 Db2 9.7.10 版本，突然某一天交易响应率下降，同时在数据库上看到大量的 latch。如清单 1 所示。

#### 清单 1 问题发生时的 latch 信息

```
> db2pd -latches：
Address Holder Waiter Filename LOC LatchType HoldCount

0x0A00020001DA22A8 511780     69051      /view/db2_v97fp10_aix64_s141015/vbs/engn/sqp/inc/sqlpLockInternal.h 513        SQLO_LT_SQLP_LHSH__xhshlatch 1
0x0A00020001DA22A8 511780 261703
/view/db2_v97fp10_aix64_s141015/vbs/engn/sqp/inc/sqlpLockInternal.h 513
SQLO_LT_SQLP_LHSH__xhshlatch 1
......

```

Show moreShow more icon

第一次看到这个 SQLO\_LT\_SQLP\_LHSH\_\_xhshlatch（在本文中简称 xhshlatch），完全不知所云，我们知道分析 latch 问题的最有效的方法就是收集 stack，通过 stack 来分析 latch 堵塞的场景。当问题发生时收的 stack 如清单 2 所示。

#### 清单 2 问题发生时的 latch 信息

```
            -------Frame------ ------Function + Offset------
0x090000002104B258 thread_waitlock@glue881 + 0x8C
0x090000002104AFE0 sqloXlatchConflict + 0x1E8
0x090000002104AD3C sqloXlatchConflict@glue1AB + 0x78
0x0900000020F3D3C4 sqlplrq__FP9sqeBsuEduP14SQLP_LOCK_INFO + 0x24
0x0900000020F85FBC sqlrr_get_lock__FP8sqlrr_cbP14SQLP_LOCK_INFOUiPUi +
    0x100
0x0900000020F85DD0
    sqlra_pkg_lock__FP8sqlrr_cbPUcsT2T3T2iUiT7P14SQLP_LOCK_INFO + 0xB8
0x0900000020F852C4
    sqlra_find_pkg__FP8sqlrr_cbPUcsT2T3T2T3iT8P14SQLP_LOCK_INFOPP20sqlra_cached_package
    + 0x190
0x0900000020F84ADC sqlra_revalidate_pkg__FP8sqlrr_cb + 0xA0
0x0900000020F84E24 .sqlra_load_pkg.fdpr.clone.2194__FP8sqlrr_cbPUcsT2T3T2b
    + 0x190
0x0900000020F7DB3C
    sqlrr_sql_request_pre__FP14db2UCinterfaceUiiP16db2UCprepareInfoP15db2UCCursorInfo
    + 0x1678
0x0900000021093F54
    sqlrr_sql_request_pre__FP14db2UCinterfaceUiiP16db2UCprepareInfoP15db2UCCursorInfo
    + 0x9C

```

Show moreShow more icon

有了本系列第一部分文章的基础，再来看这个 stack 就很好理解，这里看到了熟悉的 sqlra\_load\_pkg，sqlra\_find\_pkg，sqlra\_pkg\_lock 等函数，可以知道这是 package 上相关的竞争。问题解决的具体过程不细说，最后是紧急联系 IBM Db2 国内的支持团队，告知可以设置一个注册变量来解决，db2set DB2\_APM\_PERFORMANCE=ON，设置后需要重启数据库实例，性能问题马上得以解决（但是有功能性的副作用）。

本着严谨的态度，后来对这种方案进行了深入的分析。根据相关的 technote 的说明，该注册变量打开之后，Db2 服务器在处理 package 时能够避免 package lock，于是利用 db2trc 工具进行了验证，对比了该注册变量打开和关闭情况下的 trace 情况，结果如清单 3 所示。

#### 清单 3 DB2\_APM\_PERFORMANCE 开关后不同的 trace

```
# 默认 APM_PERFORMANCE 为 OFF 的 trace
| | | | | sqlra_get_section entry
| | | | | | sqlra_load_pkg entry
| | | | | | | sqlra_revalidate_pkg entry
| | | | | | | | sqlra_find_pkg entry
| | | | | | | | | sqlra_pkg_lock entry
| | | | | | | | | | sqlra_fmt_pkg_lock entry
| | | | | | | | | | sqlrr_get_lock entry
| | | | | | | | | | | sqlplrq entry
......
| | | | | | sqlra_load_pkg exit
| | | | | | sqlra_set_stmt_authid_dynrules entry
| | | | | | sqlra_set_stmt_authid_dynrules exit
| | | | | sqlra_get_section exit

# APM_PERFORMANCE为ON的trace
| | | | | sqlra_get_section entry
| | | | | | sqlra_set_stmt_authid_dynrules entry
| | | | | | sqlra_set_stmt_authid_dynrules exit
| | | | | sqlra_get_section exit

```

Show moreShow more icon

清单 3 中，再次看到了感觉眼熟的函数们，例如 sqlra\_get\_section 和 sqlra\_pkg\_lock 等，在本系列第一部分的清单 7 中即嵌入式 SQC 的 trace 中出现过，再次说明 JDBC 方式的 package 工作机制和其他类型的 package 工作机制是基本一致的。

经过对比，明显的看到 DB2\_APM\_PERFORMANCE 为 OFF 时，执行 SECTION 即执行计划时需要有锁定 package（这里是 CLI package）的过程，当应用程序高并发的情况下，这个锁定会导致堵塞。而该注册变量为 ON 时，则完全没有锁定 package 的过程，从而避免了堵塞，解决之前的问题。至于为什么问题发生时体现在 xhshlatch 堵塞，这应该是因为 package 加锁解锁的过程中，Db2 内部 hash 对象上的互斥访问用到了 xhshlatch。

关于这个问题的信息，可以参考 technote “ [Slow performance with latch contention on SQLO\_LT\_SQLP\_LHSH\_\_xhshlatch](http://www-01.ibm.com/support/docview.wss?uid=swg21990243)“。

总结一下 DB2\_APM\_PERFORMANCE 的用法，在 Db2 9.7 版本时，该注册变量还只有 ON 和 OFF 两个值，当设置为 ON 时，所有类型的 package lock 都被避免，包括嵌入式 SQL 的 package，但没有了锁的保护，所有 package 绑定或者重新绑定都无法执行，新的嵌入式 SQL 程序也无法预编译。在 Db2 10.5 版本（或者某个 fix pack）之后，可以设置该注册变量为 16，避免了 CLI package 上的锁定，同样的也就无法对 CLI package 进行重新绑定等操作，该设置对于嵌入式 SQL 没有影响。顺便提一句，设置为 16 的这个功能，只在 technote 中有过介绍，至今在 Db2 11 版本的知识中心的介绍中，还只有 ON 和 OFF。

现在作者所在的银行已全面升级到 Db2 10.5 版本，把 DB2\_APM\_PERFORMANCE=16 作为标准，在所有上线的 Db2 生产系统进行了配置。

另外，xhshlatch 会在多种场景中出现，并不限于清单 2 中的场景。

## 必选的 BLOCKING ALL 选项

在命令 db2 list packages for all 的结果中有一项是”Blocking”，利用Db2实例程序创建的 TBREAD package 该项为”U”，而 CLI package 的该项全部是”B”。这是 package 一个重要的属性，对应的是绑定 package 命令（PREP 或者 bind）的 BLOCK 相关选项，在 Db2 知识中心 [Specifying row blocking to reduce overhead](https://www.ibm.com/support/knowledgecenter/en/SSEPGG_11.1.0/com.ibm.db2.luw.admin.perf.doc/doc/t0005281.html) 解释了这些选项的作用。简单的说，如果绑定 package 时指定了 BLOCKING ALL 选项，应用程序用游标方式获取数据时，对于某些类型的 SQL 语句，Db2 服务器会把需要的数据以数据块的形式发送到应用程序，从而减少交互成本。

为了理解 BLOCKING ALL 选项的作用，可以做一个简单的对比测试，还是用 tbread.sqc 并把其中一个函数 TbSelectUsingFetchIntoHostVariables() 改为无限循环调用，同时准备一个监控脚本 sqlmonitor.sh（如清单 4 所示）来显示 Db2 处理 SQL 的速度。

#### 清单 4 监控脚本sqlmonitor.sh

```
            > vi sqlmonitor.sh
#!/usr/bin/ksh
if [ $# -eq 0 ] ; then
echo "Usage:sh sqlmonitor.sh <dbname>"
exit
fi
dbname=$1
interval=10
STATIC_SQL_CUR=0
while :
do
STATIC_SQL_PRE=$STATIC_SQL_CUR
STATIC_SQL_CUR=`db2 get snapshot for database on $dbname | grep -i
    "Static" | awk '{print $5}'`
let STATIC_SQL=STATIC_SQL_CUR-STATIC_SQL_PRE
let STATIC_SQL_SPEED=STATIC_SQL/$interval
echo `date +%Y%m%d%H%M%S`" : STATIC_SQL_SPEED:"$STATIC_SQL_SPEED
sleep $interval
done

```

Show moreShow more icon

设置 BLOCKING ALL 选项是在同目录（~/sqllib/samples/c）的 embprep 文件中，默认是没有 BLOCKING ALL 选项的，根据文档的介绍，默认是 BLOCKING UNAMBIG。相关代码修改和对比结果如清单 5 所示，可以看到加了 BLOCKING ALL 后，同样是单线程的应用程序，SQL 处理速度大幅度提高（个人电脑的虚拟机上做的测试）。因此，可以认为 BLOCKING ALL 是必选项。

#### 清单 5 对比 BLOCKING ALL 的性能影响

```
# 在客户端
    > vi embprep
    Bind the program to the database.

db2 bind $1.bnd blocking all （也可以与 blocking no 对比）
    > vi tbread.sqc
    while(1)
     rc = TbSelectUsingFetchIntoHostVariables();

> rm tbread; make tbread
> nohup ./tbread sample db2inst1 passw0rd  >/dev/null &

# 在服务器端
>  ./sqlmonitor.sh sample
    // 未使用 BLOCKING ALL 的 SQL 速度（SQL/s）
    20181107082154 : STATIC_SQL_SPEED:1697
    20181107082205 : STATIC_SQL_SPEED:1711
    20181107082216 : STATIC_SQL_SPEED:1705
    20181107082227 : STATIC_SQL_SPEED:1529
    20181107082238 : STATIC_SQL_SPEED:1428

    // 使用 BLOCKING ALL 的 SQL 速度（SQL/s）
    20181107082025 : STATIC_SQL_SPEED:5981
    20181107082036 : STATIC_SQL_SPEED:6608
    20181107082047 : STATIC_SQL_SPEED:6569
    20181107082059 : STATIC_SQL_SPEED:5693
    20181107082110 : STATIC_SQL_SPEED:5900

```

Show moreShow more icon

设置 BLOCKING ALL 选项还需要与另外两个参数相互配合，aslheapsz 和 rqrioblk，分别影响本地应用和远程应用。相关的用法在文档中都详细的介绍，本文不再重复。

## CLI0129E 错误和绑定的 CLIPKG 参数

前段时间某一天，开发的同事发过来一个开发测试过程中遇到的报错，错误信息是”[IBM][CLI Driver] CLI0129E An attempt to allocate a handle failed because there are no more handles to allocate. SQLSTATE=HY014″。看到这个信息，首先知道是个 CLI 的程序，然后去 Db2 知识中心搜索”CLI0129E”马上找到答案，简单的说是应用程序中申请的某种句柄（handle）达到了数量上限，要么是应用程序 bug 没有及时释放相关句柄，要么是确实需要那么多的句柄。解决方法是及时的释放句柄，或者通过 CLIPKG 来增加可用句柄数量。

接着通过 CLIPKG 关键字进行搜索，找到一篇 technote 讲的是 Java 程序遇到的错误信息”Out of Package Error Occurred”，然后引出了一个 Java 程序最多并发执行多少个 statement 的问题。结合这篇 technote 和本系列第一部分文章中认识 CLI package 章节的内容，我们可以知道，Java 程序执行时也是在 Db2 服务器端执行 CLI package 中的动态 SECTION，每个动态语句需要一个 SECTION。默认情况下 Db2 服务器绑定了 3 个大 CLI package 和 3 个小 CLI package，这也解答了之前提的问题，为什么除了 SYSSTAT 之外的每一种类型（相同的 size、HOLD 属性、隔离级别）的 CLI package 都是三个，例如 SYSLH100/101/102， 因为这是 CLIPKG 参数决定的，该参数的默认值为 3，也就是最多支持(3 _63) + (3_ 383) = 1338 个并行 SQL（每个 package 有两个预留的 SECTION）。当应用程序需要更大的并发来执行 SQL 时，可以通过 CLIPKG 设置更大的值来扩大句柄的数量，例如使用命令 db2 “bind @db2cli.lst CLIPKG 10 grant public blocking all”来重新绑定 package。当然，如果设置了上文介绍的 DB2\_APM\_PERFORMANCE=16 是无法重新绑定 CLI package 的。

由于这篇很有信息量的 technote 是从 Java 的问题场景展开的，这里来模拟重现 CLI 方式的问题。还是用 Db2 自带的示例程序 tbread.c，另外对该文件的代码稍作改动，人为的引入一个 bug，即 statement handle 用完之后不释放，代码修改部分以及运行结果如清单 6 所示。

#### 清单 6 重现 CLI 程序的报错

```
# 在客户端
    > cd sqllib/samples/cli

> export PATH=.:$PATH   //将当前目录加入 PATH 环境变量
    > vi tbread.c

    int i = 0;
    while(i < 3000)  /*main 函数中*/
        {
         rc = TbBasicSelectUsingFetchAndBindCol(hdbc);
         printf("\n------------COUNTER = %d",i++);
        }

    /* cliRC = SQLFreeHandle(SQL_HANDLE_STMT, hstmt);*/  /*第 207 行附近*/
        > rm tbread; make tbread

> ./tbread sample  db2inst1 passw0rd |egrep 'CLI0129E|COUNTER'
    ------------COUNTER = 1325
    [IBM][CLI Driver] CLI0129E An attempt to allocate a handle failed because
        there are no more handles to allocate. SQLSTATE=HY014

> db2 "bind @db2cli.lst CLIPKG 6 grant public blocking all"
> ./tbread sample  db2inst1 passw0rd |egrep 'CLI0129E|COUNTER'
    ------------COUNTER = 2468
    [IBM][CLI Driver] CLI0129E An attempt to allocate a handle failed because
        there are no more handles to allocate. SQLSTATE=HY014

```

Show moreShow more icon

在清单 6 中可以看到，CLI 示例程序调用函数 1326 遍之后报错 CLI0129E，这个数量与 technote 中的 1338 个略有偏差，然后指定 CLIPKG 6 重新绑定 CLI package 之后，CLI 示例程序调用函数 2469 遍之后报错，这个数字也有偏差，作者还不知道具体什么原因，无论如何，我们已经重现了 CLI 程序也有最大句柄限制的问题，并且验证了指定较大的 CLIPKG 可以扩大这个限制。

另外，根据相关文档上的介绍，较大的 CLIPKG 会带来性能的负面影响，大致的原因应该是，数据库内 package 及 SECTION 数量增加了，在加载 SECTION、搜索 package（回忆一下本系列第一部分文章清单 7 的那些函数）等过程会有更多的资源消耗或者资源竞争。因此，对于 CLI0129E 报错还是需要从应用程序的角度进行分析，合理控制句柄的申请数量并及时释放。

### Db2 package“分身术”

最后简要的介绍一种解决性能问题的思路，可以称之为“分身术”，同样也可以一定程度的解决 xhshlatch 问题，例如某些无法设置 DB2\_APM\_PERFORMANCE 注册变量的系统。

对于使用 CLI/JCC 程序，可以多次绑定 CLI package 并指定不同的 COLLECTION，例如 db2 bind @db2cli.lst collection USERx blocking all grant public，同时在应用程序建立连接时配置属性来指定 COLLECTION，例如 props.setProperty(“currentPackageSet”,”USER2″)或者 props.setProperty(“jdbcCollection”,”USER1″)。

对于 sqc 程序，默认情况下预编译绑定之后在数据库服务器只有一个 package，并且其 package name 与 sqc 文件同名（取前 8 位），为了避免单一 package 上的压力和竞争，也可以编译发布多份可执行文件，在每次预编译时指定不同的 package 名字，例如 db2 prep abc.sqc package using abc1。根据作者很崇拜的一位 Db2 大师的观点，通常把 package 分成两个会有很大的性能提升，再分多了则效果不明显。

## 总结

行至文末，笔者感觉仍然有非常多的内容可供探索，有些内容限于本人的知识水平也没讲的很透彻。 Db2 package 是一个非常重要的对象，本系列的第一部分笔者详细介绍了嵌入式 SQL 的 package 和 CLI package，包括基本概念和内部工作机制，相信大家也基本理解了本系列开头所说的 package 无处不在的含义。第二部分文章结合笔者的经历，介绍了几个与 package 相关的性能问题，并给出了相应的方案，希望对大家有所帮助。

最后，本文仅代表笔者观点，不代表 IBM Db2 官方观点。

## 参考资料

**学习**

- [DB2 packages: Concepts, examples, and common problems](https://www.ibm.com/developerworks/data/library/techarticle/dm-0606chun/index.html) ，这是 2006 年 developerWorks 英文网站的文章，虽然很老且少量内容不准确，但非常有帮助，其介绍了基本概念和常见使用上的问题，而本文关注点是内部机制和性能优化。
- [Slow performance with latch contention on SQLO\_LT\_SQLP\_LHSH\_\_xhshlatch](http://www-01.ibm.com/support/docview.wss?uid=swg21990243) ，xhshlatch 问题的终结者。
- [Performance degradation when running CLI / JCC applications with many short transactions](http://www-01.ibm.com/support/docview.wss?uid=swg21656754) ，这是 package”分身术”之一。
- [Specifying row blocking to reduce overhead](https://www.ibm.com/support/knowledgecenter/en/SSEPGG_11.1.0/com.ibm.db2.luw.admin.perf.doc/doc/t0005281.html) ，介绍了为什么要用 BLOCKING ALL 选项。
- [CLI0129E](https://www.ibm.com/support/knowledgecenter/en/SSEPGG_11.1.0/com.ibm.db2.luw.messages.cli.doc/com.ibm.db2.luw.messages.cli.doc-gentopic1.html#cli0129e) ，开发团队成员问这个报错是 DBA 们就知道该怎么解决了。
- [SQL0805N How many concurrently running statements allowed for a DB2 Java application and how to increase it?](http://www-01.ibm.com/support/docview.wss?uid=swg21670200) ，从一个 Java 报错引出了对 CLI package 中动态 SECTION 数量限制的分析。
- 通过 [DB2 知识中心](https://www.ibm.com/support/knowledgecenter/en/SSEPGG_11.1.0/com.ibm.db2.luw.apdv.embed.doc/doc/c0005532.html) ，了解 DB2 的详细产品信息和相关技术等全面的内容。
- 随时关注 developerWorks [技术活动](http://www.ibm.com/developerworks/cn/offers/techbriefings/) 。

**获得产品和技术**

- 使用可直接从 developerWorks 下载的 [IBM 产品评估试用软件](http://www.ibm.com/developerworks/cn/downloads/) 构建您的下一个开发项目。