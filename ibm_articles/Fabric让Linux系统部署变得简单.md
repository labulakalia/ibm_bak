# Fabric 让 Linux 系统部署变得简单
解决工作中的自动化运维需求

**标签:** Linux,Python

[原文链接](https://developer.ibm.com/zh/articles/simplyfy-linux-deployment-with-fabric/)

曹 元其

发布: 2018-07-06

* * *

## 前言

Fabric 是一个是基于 Python 库的管理工具。系统管理员在企业运维中有时要管理许多 Linux 计算机，许多工作经常要重复操作，自动化和批量化工作是系统管理员面临的问题，其实只要目标计算机可以使用 SSH 登录（如在主机 A 上对主机 B 远程运行 shell 命令），就可以借助 Fabric 来进行远程自动化和批量化的任务操作，比如：查看本地和远程计算机信息，管理密码，在服务器之间进行数据传输，本文将以这些操作为例进行讲解。需要说明的是 Fabric 工具比较适合于 PHP、Python 等解释型语言的工作，而对于编译型的 Java、C# 等就不太有效了。

## Fabric 简介

正如 Fabric 自述文档中所说：“Fabric 是 Python 语言(2.5-2.7) 的库和命令行工具，可以用来提高基于 SSH 的应用部署和系统管理效率。” Fabric 的特点包括：

- 代码可读性强，基于 Python 语言库。
- 可以进行了本地、远程操作。
- 内置了一些 API，可以根据需求灵活的定义参数。
- 可以根据角色定义，方便批量管理。
- 它可以运行本地或远程 shell 命令，包括：执行命令、文件上传和下载、日志输出等。

## 安装软件包

### 软件包依赖说明

使用 Fabric 需要预先安装 Python（ [https://www.python.org/](https://www.python.org/) ，2.5 -2.7 版本），以及 OpenSSH（ [http://www.openssh.com/](http://www.openssh.com/) ）软件包。

### 安装测试 Fabric

对于基于 Red Hat 的发行版本的系统，可以使用 epel 源来安装 Fabric。 Debian 及其派生的系统可以使用 apt-get 来安装，也可以使用 pip 工具来安装最新版本 Fabric。安装完成后使用如下命令测试安装是否成功，如果显示出版本说明安装成功，如清单 1 。

##### 清单 1\. 测试安装

```
#python -c "from fabric.api import * ; print env.version"
1.6.2

```

Show moreShow more icon

说明： `1.6.2` 就是当前 Fabric 版本号。

## fab 命令简介

Fabric 最常用的命令行工具是 fab，fab 命令格式如下：

`fab [options] <command>[:arg1,arg2=val2,host=foo,hosts='h1;h2',...] ...`

fab 命令常用参数：

- `-l` ：显示可用的 task（任务）。
- `-f` ：指定 fab 命令执行文件，默认文件是： `fabfile/fabfile.py` 。
- `-g` ：指定网关（中转）设备，填写 IP 即可。
- `-H` ：指定目标计算机，多台计算机用’,’号分隔。
- `-p` ：远程账号的密码，fab 执行时默认使用 root 账户。
- `-P` ：以异步并行方式运行多计算机任务，默认为串行运行。
- `-R` ：指定角色，以角色名区分不同业务组设备。
- `-t` ：设置设备连接超时时间（秒）。
- `-T` ：设置远程计算机命令执行超时时间（秒）。
- `-w` ：当命令运行失败，发出警告信息。

## 常用的 Fabric API

Fabric 的核心 API 主要有 7 种：颜色输出（color output）、上下文管理（context managers）、装饰器（decorators）、网络（network）、操作（operations）、任务（tasks）、工具（utils）。这些 API 提供了一系列可以与远程服务器交互的方法。下面列举几个常用的 Fabric API：

- run：执行远程命令，可以在在一台或者多台远程计算机上面执行 shell 命令。
- local：执行本地的命令或者脚本。
- lcd：切换本地目录。
- cd：切换远程目录。
- sudo：sudo 方式执行远程命令。
- put：上传和分发文件到远程计算机。
- get：是从远程计算机复制文件到本地,功能跟 scp 类似。一般用来从远程计算机下载备份文件，或者日志文件等。
- prompt：获得用户输入信息。
- confirm：确认获得的提示信息。
- reboot：重启远程计算机。
- @task：函数修饰符，标识的函数为 fab 可调用的，非标记对 fab 命令不可见。
- @runs\_once：函数修饰符，标识的函数只会执行一次，不受多台计算机影响。

Fabric API 使用了一个名为 env 的关联数组（Python 中的词典）作为配置目录，来储存 Fabric 要控制的机器的相关信息。 `env` 对象作用是定义 `fabfile` 的全局设定，下面对各属性进行说明：

- `env.hosts` ：定义目标计算机，可以用 IP 或计算机名表示。
- `env.exclude_hosts` ：排除指定计算机。
- `env.user` ：定义用户名。
- `env.port` ：定义端口，默认为 22。
- `env.password` ：定义密码。
- `env.passwords` ：定义多个密码，不同计算机对应不同密码。
- `env.gateway` ：定义网关 IP。
- `env.roledefs` ：定义角色分组。
- `env.deploy_release_dir` ：自定义全局变量。
- `env.roledefs` ：定义角色。

## fabfile 文件的编写

fabfile 文件是用来指定 Fabric 执行的命令，它通常被命名为 `fabfile.py` 并用 fab 命令运行。用户可以使用编辑器（vi 或者 gedit）创建一个名为 `fabfile.py` 的 Python 脚本。也可以使用其他名字来命名这个 Python 脚本，但是执行命令的时候需要指定这个脚本的路径： `fabric --fabfile /tmp/file.py` 。

## 输出本地和远程计算机信息

清单 2 代码的作用是在终端输出本地和远程计算机信息。这里涉及几个命令： `uname` ， `ps` 和 `uptime` 。通过一定的参数设置我们可以获取本地主机和远程主机的系统信息。

##### 清单 2\. 输出本地和远程计算机信息

```
#-*- coding: utf-8 -*-
from fabric.api import *
env.hosts=[ "10.0.2.15", "47.104.174.146"]
env.user = "root"
def local_info():
local('uname -a')
def remote_info():
with cd('/home'):
run('ps auxw|head -1;ps auxw|sort -rn -k3|head -3 && uptime ')

```

Show moreShow more icon

**清单 2 的操作流程:**

1. 首先计算机 Server 1 在执行本地计算机的一个命令操作：查看系统版本信息（使用的命令是 `uname -a` 。
2. 然后登录远程计算机 Server 2，在 Server 2 执行三个命令：首先切换到 home 目录下，然后连续执行两个命令： `ps auxw|head -1` 和 `ps auxw|sort -rn -k3|head -3 && uptime` 。这两个命令是显示 CPU 占用最多的前 3 个进程，然后查看系统负载。
3. 使用命令执行完成后退出连接。

##### 清单 3\. 命令执行界面

```
# fab local_info remote_info
[localhost] local: uname -a
Linux localhost 3.10.0-862.2.3.el7.x86_64 ：1 SMP Wed May 9 18:05:47 UTC 2018 x86_64 x86_64x86_64 GNU/Linux
[47.104.174.146] run: ps auxw|head -1;ps auxw|sort -rn -k3|head -3 && uptime
[47.104.174.146] out: USER PID %CPU %MEM VSZ RSS TTY STAT START TIME COMMAND
[47.104.174.146] out: root 1332 0.1 0.2 132040 17280 ? Ssl May20 15:07/usr/local/aegis/aegis_client/aegis_10_41/AliYunDun
[47.104.174.146] out: www 13969 0.0 0.0 168832 6464 ? S 09:32 0:00/alidata/server/httpd/bin/httpd -k start
[47.104.174.146] out: www 13906 0.0 0.0 168832 6464 ? S 08:56 0:00/alidata/server/httpd/bin/httpd -k start
[47.104.174.146] out: 09:35:21 up 6 days, 40 min, 1 user, load average: 0.00, 0.00,0.00
[47.104.174.146] out:
Done.
Disconnecting from 47.104.174.146... done.

```

Show moreShow more icon

## 密码管理

Fabric 远程执行命令的所有操作都是基于 SSH 执行的，必要时它会提示输入口令。清单 2 的例子中两台计算机的用户名一样是 root，但密码不一样。在执行远程计算机操作的时候要输入 root 密码。如果要在多台机器上执行任务，密码输入的过程也可以自动化完成，实现方式有以下两种。

### 通过 `env.passwords` 配置目标机器

一种是通过 `env.passwords` 配置目标机器的登录信息，方法是在清单 2 中第 4 行开始添加 3 行：

```
env.passwords={
      'root@10.0.2.15:22 ：'1111',
      'root@47.104.174.146:22" : '1111',}

```

Show moreShow more icon

由于通过 `env.passwords` 配置了目标机器的登录用户名/密码，所以，当我们在终端运行 `fab` 命令时，就不用手动输入密码了，大大方便了脚本远程自动执行。但是， **这种明文指定登录名/密码的方式存在安全性问题 。**

### 用证书配置无密码的 SSH 连接

另一种更安全、方便的办法是用证书配置无密码的 SSH 连接。步骤如下：

1. 首先在 Server 1 生成 SSH Key：

    `$ ssh-keygen -t rsa -b 4096`

    键入以上命令后，会出现一连串的提示，忽略它，一直按回车键即可。执行完成后，在的 ~/.ssh/ 目录生成以下两个文件：

    `~/.ssh/id_rsa 私钥`

    `~/.ssh/id_rsa.pub 公钥`

2. 把生成的公钥文件 `~/.ssh/id_rsa.pub` 里面的数据添加进远程服务器 Server 2 的 `authorized_keys file` 文件中。
3. 操作完成后可以把以下三行：





    ```
    env.passwords={
          'root@10.0.2.15:22 ：'1111',
          'root@47.104.174.146:22" : '1111',}

    ```





    Show moreShow more icon

    修改为一行：

    `env.key_filename = '~/.ssh/id_rsa'`

4. 重新执行命令：

    `#fab local_info remote_info`


整个过程和清单 3 命令执行界面相同，不过用户无需手工输入密码，这样就更加简单也更加安全。

## 在计算机之间传输数据

下面这个实例是用来在两台服务器之间进行数据传输的，清单 4 代码涉及 Fabric 几个的常用 API：local，put 和 run，通过和 Linux 命令的组合我们可以实现网络中计算机之间进行数据传输的功能。fabfile 文件内容见清单 4, 清单 4 的网络拓扑和清单 2 相同。

##### 清单 4\. fabfile 文件的内容

```
# -*- coding: utf-8 -*-
from fabric.api import *
env.hosts=[ "10.0.2.15", "47.104.174.146"]
env.user = "root"
env.key_filename = '~/.ssh/id_rsa'
def send_file ():
       local(' tar -czvf vm.tar.gz ./vm')
       put('~/wm.tar.gz', '~')
       run('tar -xzvf ~/wm.tar.gz && mv /vm /var/www chown -R apache.apache /var /www/vm  && /etc/init.d/httpd restart')
if __name__ == '__main__':
send_file ()

```

Show moreShow more icon

清单 4 的执行命令是：

`# fab send_file`

清单 4 的操作流程：

1. 首先是在 Server 1 把目录 `vm` 中的所有文件进行打包。
2. 然后把打包文件推送到 Server 2 的 `/var/www` 目录下，然后解压缩，设置权限。
3. 重新启动 Server 2 的 Apache 服务器。

## 结束语

Fabric 是一个的轻量级运维工具，熟练掌握其用法能够解决工作中的很多自动化运维需求，最大特点是不用登录远程服务器，在本地运行远程命令，几行 Python 脚本就可以完成任务。你可以定义一系列的任务函数，然后灵活的指定在哪些主机上执行哪些任务。Fabric 适合管理大量主机的场景，比系统如运维，私有云管理，自动化部署等。Fabric 还一些功能，比如角色的定义，远程交互及异常处理，并发执行等，希望本文能够引起你对 Fabric 的兴趣，帮助你在实际应用中解决问题。