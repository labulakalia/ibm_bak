# Db2 LUW federation server 的安全机制
深入浅出了解 Db2 LUW federation server 的安全机制

**标签:** 分析

[原文链接](https://developer.ibm.com/zh/articles/ba-lo-security-mechanism-for-db2-luw-federation-server/)

刘 平

发布: 2019-09-11

* * *

## Db2 的安全认证

Db2 的用户管理与其它数据库比如 Oracle，SQLServer 不同，Db2 需要依赖外部的认证系统做认证，Db2 本身并没有认证组件。Db2 支持多种认证机制，本文重点介绍普通用户名和密码的认证方式。当我们想要使用一个用户名登录 Db2 时，这个用户名必须已经在操作系统中存在，因为对于普通用户名和密码的认证方式，Db2 是通过操作系统的安全机制进行认证的。如果在登录 Db2 时没有提供用户名和密码，那么 Db2 会自动使用当前登录系统并发出登录 Db2 请求的用户名和密码。如果一个用户成功的通过认证，那么他会自动登录并连接到 Db2，Db2 会将该用户名和密码保存在一个特殊的数据结构里。而 Federation Server 在进行用户认证时就是在这个数据结构里获取 Db2 的认证用户信息的。Db2 也支持其它认证方式，请参见 [这里](https://www.ibm.com/support/knowledgecenter/en/SSEPGG_11.1.0/com.ibm.db2.luw.admin.sec.doc/doc/c0021804.html) ，本文主要针对普通的用户名和密码方式进行探讨。

## Db2 安装账户(db2inst1, db2fenc1)

首先我们需要了解 Db2 运行时所需要的必要账户，在安装完 Db2 后，通常需要做的第一件事情就是创建 Db2 实例，在创建实例之前，需要首先创建至少两个用户。 一个是 Db2 实例用户，另一个是 Db2 fenced 进程用户。Db2 实例的所有者默认拥有该实例的所有权限。Db2 fenced 进程用户主要用于 Db2 执行 UDF 或存储过程等位于数据库执行地址空间以外的程序，通常在这些程序遇到错误或者异常退出时不会对 Db2 内核造成影响。Fenced 用户对 Db2 非常重要，因为 Db2 的许多功能都是基于 UDF 或存储过程执行完成的。Federation Server 就属于这些功能中的一个，当用户创建了 Fenced wrapper 时，那么所有 Federation Server Wrapper 及其包含的组件都将运行在 Fenced 模式下，而这些进程的所有者就是这个 Fenced 用户。通常 Federation server 遇到的权限有关的问题都是因为不正确的 Fenced 用户引起的。我们可以通过以下方式确认当前实例的 Fenced 用户。

### 第一种方式使用 db2pd -fmp 命令

启动 db2 instance，然后执行 db2pd -fmp，以下内容为该命令的输出结果。

```
        Database Member 0 -- Active -- Up 6 days 05:38:07 -- Date 2019-06-20-00.58.24.695785

FMP:
Pool Size:       1
Max Pool Size:   200 ( Automatic )
Keep FMP:        YES
Initialized:     YES
Trusted Path:    /home/db2inst1/sqllib/function/unfenced
Fenced User:     db2fenc1
Shared Memory:   0x0000000201C90440
IPC Pool:        0x0000000201C904A0

FMP Process:
Address            FmpPid     Bit   Flags      ActiveThrd PooledThrd ForcedThrd Active IPCList
0x0000000203F4FB00 18187      64    0x00000000 0          0          0          No     0x0000000203B3EA80

Active Threads:
Address            FmpPid     EduPid     ThreadId
No active threads.

Pooled Threads:
Address            FmpPid     ThreadId
No pooled threads.

Forced Threads:
Address            FmpPid     ThreadId
No forced threads.

FMP Process:
Address            FmpPid     Bit   Flags      ActiveThrd PooledThrd ForcedThrd Active IPCList
0x0000000203C0E780 2432       64    0x00000002 0          2          0          Yes    0x0000000203C0F860

Active Threads:
Address            FmpPid     EduPid     ThreadId
No active threads.

Pooled Threads:
Address            FmpPid     ThreadId
0x0000000203B3FD00 2432       1325307648
0x0000000203C0F480 2432       1287427840

Forced Threads:
Address            FmpPid     ThreadId
No forced threads.

```

Show moreShow more icon

在该命令的输出结果中 Fenced User 表征的即是当前的 Fenced User。

### 第二种方法查看.fenced 文件

.fenced 文件位于实例的安装目录下，具体为：/sqllib/adm。可以执行以下命令查看该文件的属性。

```
ls -la .fenced

```

Show moreShow more icon

通过检查这个文件的所有者，可以得到当前该 instance 的 fenced 用户。

```
$ ls -la .fenced
-r--r--r-- 1 db2fenc1 db2fenc1 0 Jun 13 19:06 .fenced

```

Show moreShow more icon

## Federation server 安全认证机制

在了解了 Db2 的实例账户和 fenced 账户后，可以继续深入讨论 Federation Server 的安全机制。Federation Server 安全认证机制包括两部分，第一部分为 Db2 的安全认证，用户只有通过 Db2 的安全认证后才能使用 Federation Server 的功能。第二部分是 Federation Server 的安全认证的核心部分，这一部分负责管理和维护本地用户和远程数据源用户之间的映射关系。比如本地用户 UserA 与远端数据源用户 RUserA 建立了映射关系，那么当 UserA 登录到 db2 时，他可以使用远程数据源用户 RUserA 查询远程数据源 DatasourceA 的数据，而本地用户 UserB 因为没有和远端数据源用户 RUserA 建立映射关系，那么他就不可以使用远程数据源用户 RUserA 查询远程数据源 DatasourceA 的数据。Federation Server 通过提供 Create user mapping 这样的 DDL 完成这样的认证机制。

CREATE USER MAPPING 的语法如下：

```
CREATE USER MAPPING
FOR local_userID
SERVER server_definition_name
OPTIONS (REMOTE_AUTHID 'remote_userID',
REMOTE_PASSWORD 'remote_password')

```

Show moreShow more icon

Local\_userID 有三种不同的使用方式：

- 第一种 create user mapping FOR user SERVER server _definition\_name options (REMOTE\_AUTHID_‘remote _userID’_, REMOTE _PASSWORD_‘remote _password’_), 这个 SQL 语句将为当前登录到 Db2 的用户和远端数据源用户 remote\_user 建立映射关系。比如当前登录 Db2 的用户为 localuser1，那么执行这个 DDL 后将建立 localuser1 和 remote\_userID 的映射关系。
- 第二种 create user mapping FOR local\_username SERVER server\_definition\_name options (REMOTE\_AUTHID ‘remote\_userID’, REMOTE\_PASSWORD ‘remote\_password’), 这个 SQL 语句将为用户名为 local\_username 的用户和远程数据源用户 remote\_userID 建立映射关系，比如当前登录 Db2 的用户为 local\_user1，使用该方法可以为 local\_user2 建立一个映射关系。
- 第三种 create user mapping FOR public SERVER server\_definition\_name options (REMOTE\_AUTHID ‘remote\_userID’, REMOTE\_PASSWORD ‘remote\_password’)，这个 SQL 语句为任何一个可以登录到 Db2 database 的用户建立一个和远程数据源用户 remote\_userID 的映射关系，当一个用户登录到本地数据库时，如果该用户执行 Nickname 查询，那么 Federation Server 在与远程数据源建立连接时会首先查询当前用户是否存在与远程数据源用户的映射关系，如果没有则会继续查询是否有 public 到远程用户的映射关系，如果有那么就会使用该映射关系中对应的远程用户 remote\_userID 做为登录用户连接远程数据源。如果对一个远端数据源已经建立一个非 public 映射的映射关系，那么此时不能建立 public 映射，错误信息为：

```
db2 => create user mapping for public server datastore2 options(remote_authid
'db2inst1', remote_password 'Liuping123')
DB21034E The command was processed as an SQL statement because it was not a
valid Command Line Processor command. During SQL processing it returned:
SQL1515N The user mapping cannot be created for server "DATASTORE2" because
of a conflict with an existing user mapping or federated server option. Reason
code "2". SQLSTATE=428HE

```

Show moreShow more icon

当 CREATE USER MAPPING 语句执行成功后，它建立的映射关系会被存储到本地数据库的 CATALOG 中。

Db2 用户可以通过以下 DDL 查询 CATALOG 中 user mapping 的信息。

```
SELECT * FROM
                SYSCAT.USEROPTIONS

```

Show moreShow more icon

这三种创建用户映射关系的方法都有自己的应用场景，下面我们对每一种方法给出应用示例。

### 使用场景一

为远程数据源创建一个 public 映射关系。

如果远程数据源所存储的信息为非机密信息，要求的安全性不高，或者允许所有用户访问，那么我们可以在 Federation Server 这端建立一个 public 的映射关系。

```
Connect to db user instance_user using ******
      create server rcac_server1 type db2/udb version 11 authorization remote_user password
        "******" options (host 'remote_host_name', port '50000',dbname 'dbname', pushdown 'Y',
        db2_maximal_pushdown 'Y');
      create user mapping for public server rcac_server1 options (remote_authid 'remote_user',
        remote_password '******');

```

Show moreShow more icon

### 使用场景二

为当前本地用户建立一个映射关系：本地用户可以是任意已经通过 Db2 认证的用户。

```
Connect to db user dasusr1 using ******
      create server rcac_server1 type db2/udb version 11 authorization remote_user password
        "******" options (host 'remote_host_name', port '50000',dbname 'dbname', pushdown 'Y',
        db2_maximal_pushdown 'Y');
      create user mapping for user server rcac_server1 options (remote_authid 'remote_user',
        remote_password '******');

```

Show moreShow more icon

### 使用场景三

为其他用户建立映射关系：由于 Db2 有多种授权管理，不同的用户可以有不同的权限，Federation Server 的权限管理也可以基于 Db2 的授权管理实现。比如 Db2 的一些用户拥有 DDL 权限，但是没有 DML 权限，而一些用户拥有 DML 权限却没有 DDL 权限，为了能够让那些拥有 DML 权限的用户可以使用 Federation Server 功能，我们可以执行以下操作：

```
Connect to db user user_with_ddl_permission using ******
      Create user mapping for user_with_dml_permission for server server_name options
        (remote_authid 'remote_user'， remote_password'******')
     Connect reset
     Connect to db user user_with_dml_permission using *****
     Select * from nickname

```

Show moreShow more icon

当我们使用 user\_with\_ddl\_permission 用户执行 create user mapping 时，有可能遇到授权的错误消息，比如：

```
db2 => create user mapping for dasusr1 server datastore2 options(remote_authid
'db2inst1', remote_password 'Liuping123')
DB21034E The command was processed as an SQL statement because it was not a
valid Command Line Processor command. During SQL processing it returned:
SQL0551N The statement failed because the authorization ID does not have the
required authorization or privilege to perform the operation. Authorization
ID: "SESSIONUSER1". Operation: "CREATE USER MAPPING". Object: "DASUSR1".
SQLSTATE=42501

```

Show moreShow more icon

此时需要对该用户进行授权：

```
Connect to db user instance_user using ******

```

Show moreShow more icon

Grant dbadm on db to user user\_with\_ddl\_permission，在拥有了 dbadm 权限后可以继续执行 create user mapping。在成功的执行了 create user mapping 后就可以以另一个账户登录并查询远程数据源的数据。

以上我们介绍了 Federation server 的安全认证机制，即通过创建本地 Db2 用户与远程数据源用户的映射关系，并将这个映射关系存储到 CATALOG 中。当执行远端数据查询时，Federation Server 首先会以登录到 Db2 系统的用户名为查询关键字在 CATALOG 中查询远程数据源的用户名，如果有与之匹配的记录，那么就会使用该用户名所对应的远程数据源用户名与远程数据源建立连接。接下来我们继续深入分析一下这个认证过程。

## Federation Server 安全认证的隐藏功能

当用户需要通过 Federation Server 访问远程数据源时，通常需要做以下几件事情：

```
        create wrapper WRAPPER1 LIBRARY 'libdb2drda.so' OPTIONS(DB2_FENCED 'Y');
create server DATASTORE2 type db2/udb version 11 wrapper WRAPPER1 authorization "DB2INST1" password "Liuping123" options (host 'remote_host_name, port '50000', dbname 'dbname', password 'Y',  pushdown 'Y');
create user mapping for user server DATASTORE2 options (REMOTE_AUTHID 'remote_user', REMOTE_PASSWORD '******');
create nickname nktest1 for DATASTORE2."remote_user ".TEST1;
select * from nktest1;

```

Show moreShow more icon

在 create user mapping 时我们指定了远程数据源的用户名和密码，如果远程数据源的用户名密码和登录本地数据库的用户名密码相同，那么我们还需要创建这个映射关系吗，这种情况下的映射关系左右两边的值都是相同的，这样的一个映射关系是不是多余呢？

我们可以做以下测试：

```
        connect to db
create wrapper "drdawrapper" library 'libdb2drda.so' options(DB2_FENCED 'Y');
create server DATASTORE2 type db2/udb version 11 wrapper "drdawrapper" authorization db2inst1 password "******" options (host 'remote_host_name', port '50000',dbname 'testdb', pushdown 'Y', db2_maximal_pushdown 'Y');
create nickname N1 for DATASTORE2.db2inst1.T1;
DB21034E  The command was processed as an SQL statement because it was not a
valid Command Line Processor command.  During SQL processing it returned:
SQL1101N  Remote database "testdb" on node "" could not be accessed with the
specified authorization id and password.  SQLSTATE=08004

```

Show moreShow more icon

我们可以看到在没有 create user mapping 时会失败，尽管我们本地 db2 的用户名和密码与远程一致，为了进一步分析这个问题，我们把 dialog 里的信息贴出来。

```
        2019-06-13-23.30.20.608324-420 I231443E420           LEVEL: Error
PID     : 1788                 TID : 139674621990656 PROC : db2fmp (
INSTANCE: db2inst1             NODE : 000
HOSTNAME: db2source1.fyre.ibm.com
FUNCTION: DB2 UDB, Query Gateway, sqlqg_fedstp_hook, probe:40
MESSAGE : Unexpected error returned from outer RC=
DATA #1 : Hexdump, 4 bytes
0x00007F08883A4660 : 1200 2681                                  ..& .

2019-06-13-23.30.20.609354-420 I231864E678           LEVEL: Error
PID     : 1371                 TID : 140627123627776 PROC : db2sysc 0
INSTANCE: db2inst1             NODE : 000            DB   : TESTDBU
APPHDL  : 0-400                APPID: *LOCAL.db2inst1.190614063025
AUTHID  : DB2INST1             HOSTNAME: db2source1.fyre.ibm.com
EDUID   : 436                  EDUNAME: db2agent (TESTDBU) 0
FUNCTION: DB2 UDB, Query Gateway, sqlqgPassthruPrepare, probe:30
MESSAGE : ZRC=0x80260160=-2144992928=SQLQG_ERROR "Error constant for gateway."
DATA #1 : String, 71 bytes
FedStart Failed, check the remote data source or the user mapping first
DATA #2 : String, 10 bytes
DATASTORE2

2019-06-14-01.35.20.730819-420 I232543E420           LEVEL: Error
PID     : 1788                 TID : 139674621990656 PROC : db2fmp (
INSTANCE: db2inst1             NODE : 000
HOSTNAME: db2source1.fyre.ibm.com
FUNCTION: DB2 UDB, Query Gateway, sqlqg_fedstp_hook, probe:40
MESSAGE : Unexpected error returned from outer RC=
DATA #1 : Hexdump, 4 bytes
0x00007F08883A4660 : 1200 2681                                  ..& .

2019-06-14-01.35.20.731915-420 I232964E678           LEVEL: Error
PID     : 1371                 TID : 140627236873984 PROC : db2sysc 0
INSTANCE: db2inst1             NODE : 000            DB   : TESTDBU
APPHDL  : 0-675                APPID: *LOCAL.db2inst1.190614083525
AUTHID  : DB2INST1             HOSTNAME: db2source1.fyre.ibm.com
EDUID   : 435                  EDUNAME: db2agent (TESTDBU) 0
FUNCTION: DB2 UDB, Query Gateway, sqlqgPassthruPrepare, probe:30
MESSAGE : ZRC=0x80260160=-2144992928=SQLQG_ERROR "Error constant for gateway."
DATA #1 : String, 71 bytes
FedStart Failed, check the remote data source or the user mapping first
DATA #2 : String, 10 bytes
DATASTORE2

2019-06-14-03.40.19.538373-420 I233643E420           LEVEL: Error
PID     : 1788                 TID : 139674621990656 PROC : db2fmp (
INSTANCE: db2inst1             NODE : 000
HOSTNAME: db2source1.fyre.ibm.com
FUNCTION: DB2 UDB, Query Gateway, sqlqg_fedstp_hook, probe:40
MESSAGE : Unexpected error returned from outer RC=
DATA #1 : Hexdump, 4 bytes
0x00007F08883A4660 : 1200 2681                                  ..& .

2019-06-14-03.40.19.539376-420 I234064E678           LEVEL: Error
PID     : 1371                 TID : 140627123627776 PROC : db2sysc 0
INSTANCE: db2inst1             NODE : 000            DB   : TESTDBU
APPHDL  : 0-829                APPID: *LOCAL.db2inst1.190614104024
AUTHID  : DB2INST1             HOSTNAME: db2source1.fyre.ibm.com
EDUID   : 436                  EDUNAME: db2agent (TESTDBU) 0
FUNCTION: DB2 UDB, Query Gateway, sqlqgPassthruPrepare, probe:30
MESSAGE : ZRC=0x80260160=-2144992928=SQLQG_ERROR "Error constant for gateway."
DATA #1 : String, 71 bytes
FedStart Failed, check the remote data source or the user mapping first
DATA #2 : String, 10 bytes
DATASTORE2

```

Show moreShow more icon

通过分析错误消息的内容，我们可以猜测到这个错误应该和 user mapping 有关系。通过进一步对该语句进行 trace 分析，我们得到下面的信息。

```
        49568   data DB2 UDB drda wrapper DRDA_Connection::build_connection_string fnc (***********)
        pid 1788 tid 139674621990656 cpid 9061 node 0 probe 41
        bytes 72

        Data1   (PD_TYPE_STRING,21) String:
        Remote_host_name
        Data2   (PD_TYPE_STRING,5) String:
        50000
        Data3   (PD_TYPE_STRING,6) String:
        testdb
        Data4   (PD_TYPE_STRING,8) String:
        DB2INST1

```

Show moreShow more icon

通过对比源代码，我们可以确认在与远程数据源建立连接时使用的用户名是正确的，那么唯一可能错误的地方就是密码不正确。

进一步分析 trace 文件可以看到以下信息：

```
49515   data DB2 UDB drda wrapper DRDA_Connection::fold_identifier fnc (*********)
        pid 1788 tid 139674621990656 cpid 9061 node 0 probe 10
        bytes 27

        Data1   (PD_TYPE_DEFAULT,1) Hexdump:
        00                                         .

        Data2   (PD_TYPE_DEFAULT,10) Hexdump:
        4C69 7570 696E 6731 3233

```

Show moreShow more icon

在 trace 文件中我们并没有得到密码信息，这说明 Db2 在与远程数据源建立连接时没有取得正确的密码。

回顾我们登录本地数据库时，我们使用的是默认的用户名和密码（没有指定用户名那么会使用当前登录系统的默认用户名）。

如果我们显示的指定用户名和密码会不会可以成功访问远程数据源呢？

```
Connect to db user db2inst1 using ******
create wrapper "drdawrapper" library 'libdb2drda.so' options(DB2_FENCED 'Y');
create server DATASTORE2 type db2/udb version 11 wrapper "drdawrapper" authorization db2inst1 password "******" options (host 'remote_host_name', port '50000',dbname 'testdb', pushdown 'Y', db2_maximal_pushdown 'Y');
create nickname N1 for DATASTORE2.db2inst1.T1;
DB20000I  The SQL command completed successfully

```

Show moreShow more icon

我们发现此时可以成功的创建 nickname。

这说明如果远程数据源的用户名和本地数据库认证通过的用户名密码相同时我们可以省却掉 create user mapping，此时 Federation Server 会使用本地登录数据库的用户名和密码登录远程数据源并建立连接。

通过上面的分析我们知道 Federation Server 实现安全机制的原理是：

1. 通过有 DDL 权限的用户创建本地用户和远程数据源用户之间的映射关系；
2. 当一个用户成功通过认证并登录到 Db2 后，Federation Server 会在 Catalog 中以该用户为键值查找相应的映射关系，如果有那么就会使用该映射关系所对应的远程数据源的用户名和密码登录远程数据源。如果没有那么在连接远程数据源时就会报相应的错误。比如我们在数据库中为用户 db2inst1 创建一个映射关系，并为之创建 Nickname。如果我们用 db2inst1 登录并查询 Nickname 时我们可以得到以下信息：

```
        db2 => select * from nktest1

C1          C2
----------- --------------------------------
          1 First Record
          2 Second Record
          3 Third Record

3 record(s) selected.

```

Show moreShow more icon

如果我们用另一个账户登录 Db2，并试图查询 Nickname 时，因为没有这个用户的映射关系，所以执行时就会报权限相关的错误。

错误信息：

```
        Connect to db user db2inst1 using ******
GRANT SELECT ON TABLE NKTEST1 TO USER DASUSR1;
CONNECT RESET
Connect to db user dasusr1 using ******

db2 => select * from db2inst1.nktest1

COL
---
SQL1101N  Remote database "testdb" on node "< unknown>" could not be accessed
with the specified authorization id and password.  SQLSTATE=08004

```

Show moreShow more icon

## SESSION 用户

Db2 LUW 有 SESSION USER 的概念，具体定义可以参考 [这里](https://www.ibm.com/support/knowledgecenter/en/SSEPGG_11.1.0/com.ibm.db2.luw.sql.ref.doc/doc/r0011134.html) ，如何使用 SESSION USER 可以参考 [这里](https://www.ibm.com/support/knowledgecenter/en/SSEPGG_11.1.0/com.ibm.db2.luw.sql.ref.doc/doc/r0011139.html) 。接下来我们继续深入分析 Db2 session user 与 Federation server 的关系。

Session user 和 Db2 for z/OS 的 set current sqlid 非常相似。一个用户可以被授权成为另一个用户，这个过程就是通过 set session authorization 完成的。当一个用户通过被授权成为另一个用户后，那么这个新用户将成为所有在此之后创建的数据库对象的所有者。这个过程完成的功能和以下相同：

首先我们使用 db2inst1 登录到 Db2：

```
db2 => connect to testdbu user db2inst1 using ******

```

Show moreShow more icon

此时如果我们想更换当前的用户，我们可以执行以下操作：

```
db2 => connect reset
      DB20000I The SQL command completed successfully.
      db2 => connect to testdbu user dasusr1 using ******
      db2 =>

```

Show moreShow more icon

通过以上操作，我们把登录账户从 db2inst1 切换到了 dasusr1。Set session authorization 实现的功能和上述操作相同，set session authorization 的执行需要特殊的权限。下面我们通过实例来解释如何使用 set session authorization。

```
      db2 => connect to testdbu user db2inst1 using ******
db2 => grant dbadm on database to user dasusr1
DB20000I  The SQL command completed successfully.
db2 => grant setsessionuser on user sessionusr1 to user dasusr1
DB20000I  The SQL command completed successfully.
db2 => connect reset
DB20000I  The SQL command completed successfully.
db2 => connect to testdbu user dasusr1 using ******
db2 => select * from dasusr1.dasusr_table1

C1
-----------

0 record(s) selected.

db2 => set session authorization sessionuser1
DB20000I  The SQL command completed successfully.
db2 => select * from dasusr1.dasusr_table1
SQL0551N  The statement failed because the authorization ID does not have the
required authorization or privilege to perform the operation.  Authorization
ID: "SESSIONUSER1".  Operation: "SELECT". Object: "DASUSR1.DASUSR_TABLE1".
SQLSTATE=42501

```

Show moreShow more icon

这里我们可以发现在执行了 set session authorization 之后在此查询 dasusr1.dasusr\_table1 这张表的时候我们已经没有了权限。说明此时的认证用户已经发生了变化。

值得注意的是 Session user 可以为非操作系统的账户，事实上在实际的应用中 Session user 都不是操作系统的账户。同时 Db2 instance user 不可以执行 set session authorization。

以下操作演示了如何用一个非操作系统用户作为 session user。

```
        db2 => connect to testdbu user db2inst1 using ******
db2 => grant setsessionuser on user nonsystemuser to user dasusr1
DB20000I  The SQL command completed successfully.
db2 => connect reset
DB20000I  The SQL command completed successfully.
db2 => connect to testdbu user dasusr1 using ******
db2 => set session authorization nonsystemuser
DB20000I  The SQL command completed successfully.
db2 => select * from dasusr1.dasusr1_table
SQL0551N  The statement failed because the authorization ID does not have the
required authorization or privilege to perform the operation.  Authorization
ID: "NONSYSTEMUSER".  Operation: "SELECT". Object: "DASUSR1.DASUSR1_TABLE".
SQLSTATE=42501
db2 =>

```

Show moreShow more icon

除了可以应用 set session authorization 设置 session user 外，我门还可以通过 set session\_user=user 完成相同的功能。

```
        db2 => connect to testdbu user dasusr1 using ******
db2 => set session_user=sessionuser1
DB20000I  The SQL command completed successfully.
db2 =>

```

Show moreShow more icon

通过上面的介绍我们发现 Db2 可以通过设置 session user 从而更换当前的 authorization ID，而当我们执行 Federation Nickname 查询时 Federation Server 首先要获取 Db2 当前的 authorization ID 并用此 ID 查询与远程数据源用户的映射关系。那么当 Db2 通过 set session authorization/set session\_user=username 更换当前 authorization ID 后，Federation 是否能感知到这个变化呢，比如当以用户 db2inst1 登录到 Db2 后，并执行以下操作：

```
        db2 => connect to testdbu user db2inst1 using ******
db2 => create wrapper WRAPPER1 LIBRARY 'libdb2drda.so' OPTIONS(DB2_FENCED 'Y')
DB20000I  The SQL command completed successfully.
db2 => create server DATASTORE2 type db2/udb version 11 wrapper WRAPPER1 authorization "DB2INST1" password "******" options (host 'remote_host_name', port '50000', dbname 'testdb', password 'Y',  pushdown 'Y')
DB20000I  The SQL command completed successfully.
db2 => create user mapping for user server DATASTORE2 options (REMOTE_AUTHID 'DB2INST1', REMOTE_PASSWORD '******')
DB20000I  The SQL command completed successfully.
db2 => create nickname nktest1 for DATASTORE2."DB2INST1".TEST1
DB20000I  The SQL command completed successfully.
db2 => select * from nktest1

C1          C2
----------- --------------------------------
          1 First Record
          2 Second Record
          3 Third Record

3 record(s) selected.

db2 =>

```

Show moreShow more icon

值得注意的是这些操作中，我们为当前用户 db2inst1 创建了一个与远程数据源用户的映射关系。也就是说当下次以 db2inst1 用户登录 Db2 并执行对 Nickname nktest1 的查询操作时我们可以成功得到以下结果。

```
          db2 => connect to testdbu user db2inst1 using ******
db2 => select * from nktest1

C1          C2
----------- --------------------------------
          1 First Record
          2 Second Record
          3 Third Record

3 record(s) selected.

db2 =>

```

Show moreShow more icon

如果为 sessionuser1 创建一个与远程的映射关系，那么 Federation Server 是否可以在 Db2 切换 session user 到 session user1 后成功连接远程数据源呢？

```
        db2 => create user mapping for sessionuser1 server DATASTORE2 options (REMOTE_AUTHID 'DB2INST1', REMOTE_PASSWORD '******')
DB20000I  The SQL command completed successfully.
db2 => connect reset
DB20000I  The SQL command completed successfully.
db2 => connect to testdbu user dasusr1 using ******
db2 => grant select on table db2inst1.nktest1 to user sessionuser1
DB20000I  The SQL command completed successfully.
db2 => set session_user=sessionuser1
DB20000I  The SQL command completed successfully.
db2 => select * from db2inst1.nktest1

C1          C2
----------- --------------------------------
SQL1101N  Remote database "testdb" on node "< unknown>" could not be accessed
with the specified authorization id and password.  SQLSTATE=08004
db2 =>

```

Show moreShow more icon

可以发现即使为 sessionuser1 建立了与远程数据源用户的映射关系，在切换到 sessionuser1 后依旧不能连到到远程数据源，错误消息提示为错误的用户名或密码。

如果直接以用户 sessionuser1 登录到 Db2，是否可以成功连接远程数据源呢？

```
db2 => connect reset
        DB20000I  The SQL command completed successfully.
        db2 => connect to testdbu user sessionuser1 using ******
        db2 => select * from db2inst1.nktest1

        C1          C2
        ----------- --------------------------------
        1 First Record
        2 Second Record
        3 Third Record

        3 record(s) selected.

        db2 =>

```

Show moreShow more icon

可以发现这种情况下是可以直接登录到远程数据源并获取数据的，这说明 Federation Server 并不能感知 Db2 set session authorization 带来的变化。通过分析 trace 文件也能证明我们的结论。

```
          23008   data DB2 UDB drda wrapper DRDA_Connection::fold_identifier fnc (3.3.44.108.0.10)
        pid 1788 tid 139674621990656 cpid 12519 node 0 probe 10
        bytes 24

        Data1   (PD_TYPE_DEFAULT,1) Hexdump:
        00                                         .

        Data2   (PD_TYPE_DEFAULT,7) Hexdump:
        4441 5355 5352 31                          DASUSR1

23012   data DB2 UDB drda wrapper DRDA_Connection::fold_identifier fnc (3.3.44.108.0.10)
        pid 1788 tid 139674621990656 cpid 12519 node 0 probe 10
        bytes 27

        Data1   (PD_TYPE_DEFAULT,1) Hexdump:
        00                                         .

        Data2   (PD_TYPE_DEFAULT,10) Hexdump:
        4C61 6E6C 616E 2131 3233                   ******

```

Show moreShow more icon

发送给远端数据源进行认证的用户名和密码为本地登录 Db2 的用户名和密码，根据之前得出的结论说明在 Catalog 中以 Sessionuser1 为键值查找与远程数据源用户的映射关系时没有找到相应的映射关系。而这与事实不符，因为我们已经为 sessionuser1 创建了映射关系。

```
        db2 => Select * from SYSCAT.USEROPTIONS where authid='SESSIONUSER1'

---------------------    ------------------  --------------------  ------------------------      ----------------

2 record(s) selected.

```

Show moreShow more icon

这说明在 Federation Server 执行安全认证时并没有以 sessionuser1 为键值查找映射关系，而是以当前的登录用户 dasusr1 为键值查找的。

在分析了源代码后发现当前的 Federation Server 确实是这样实现的。它并不能感知到 Db2 session user 的变化，事实上 session user 存储在不同的结构体中，而 Federation Server 只会读取第一次登录到本地数据库时的 authentication ID。

这为 Federation Server 客户带来了使用上的不便，限制了 Federation Server 的安全认证灵活性，因为 Db2 session user 可以让用户更灵活的管理自己的权限。

Federation Server 已经着手开发这个新功能，并计划在 Db2 的新版本中发布。最新的计划是 Db2 V11.5 的第一个 fixpack。

## 总结

本文系统的讲述了 Db2 Federation server 的安全认证机制，包括 Db2 的安装用户，fenced 进程用户，fenced 进程用户与普通用户的关系和区别。阐述了 Federation server 如何管理本地的 Db2 用户和远程的数据源用户之间的映射关系，如何根据不同应用场景创建不同的映射关系。本文还深入分析了一个 Federation server 的隐藏功能即在本地用户名和密码远程数据源相同时，用户可以省略掉创建用户映射关系，Federation server 会默认用本地的用户名和密码与远程数据源进行认证。最后本文介绍和分析了 SESSION 用户，发现 Federation Server 不能感知到 Db2 SESSION 用户的变化，在后续版本中 Federation server 部门将会实现这个功能。

## 资源推荐

- 参考技术文档 [Db2 安全模型](https://www.ibm.com/support/knowledgecenter/en/SSEPGG_11.1.0/com.ibm.db2.luw.admin.sec.doc/doc/c0021804.html) ，这里读者可以了解 Db2 的基本安全认证过程和方法。
- 参考技术文档 Db2 SESSION 用户，这里读者可以了解 Db2 SESSION 用户的定义及其使用方法的链接。
- 参考技术文档 [Public user mapping](https://www.ibm.com/support/knowledgecenter/en/SSEPGG_11.1.0/com.ibm.data.fluidquery.doc/topics/iiynfpubusrmap.html) ，这里读者可以了解 Db2 Federation server 的 public user mapping。