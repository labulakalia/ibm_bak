# 用 Gearman 分发 PHP 应用程序的工作负载
将一个计算密集型的或专门化的功能放在一个或多个独立的专用服务器上运行

**标签:** PHP,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/os-php-gearman/)

Martin Streicher

发布: 2010-02-01

* * *

尽管一个 Web 应用程序的大部分内容都与表示有关，但它的价值与竞争优势却可能体现在若干专有服务或算法方面。如果这类处理过于复杂或拖沓，最好是进行异步执行，以免 Web 服务器对传入的请求没有响应。实际上，将一个计算密集型的或专门化的功能放在一个或多个独立的专用服务器上运行，效果会更好。

##### 常用的缩略词

- **API**：应用程序编程接口
- **HTTP**：超文本传输协议
- **LAMP**：Linux、Apache、MySQL 与 PHP

PHP 的 Gearman 库能把工作分发给一组机器。Gearman 会对作业进行排队并少量分派作业，而将那些复杂的任务分发给为此任务预留的机器。这个库对 Perl、Ruby、 `C` 、Python 及 PHP 开发人员均可用，并且还可以运行于任何类似 UNIX® 的平台上，包括 Mac OS X、 Linux® 和 Sun Solaris。

向一个 PHP 应用程序添加 Gearman 非常简单。假设您将 PHP 应用程序托管在一个典型的 LAMP 配置上，那么 Gearman 将需要一个额外的守护程序以及一个 PHP 扩展。截止到 2009 年 11 月，Gearman 守护程序的最新版本是 0.10，并且有两个 PHP 扩展可以用 — 一个用 PHP 包裹了 Gearman `C` 库，另一个用纯 PHP 编写。我们要用的是前者。它的最新版本是 0.6.0，可以从 PECL 或 Github（参见 参考资料 ）获取它的源代码。

**注意：** 对于本文而言， _producer_ 指的是生成工作请求的机器； _consumer_ 是执行工作的机器；而 _agent_ 则是连接 producer 与适当 consumer 的中介。

## 安装 Gearman

向一个机器添加 Gearman 需要两步：第一步构建并启动这个守护程序，第二步构建与 PHP 版本相匹配的 PHP 扩展。这个守护程序包包括构建此扩展所需的所有库。

首先，下载 Gearman 守护程序 `gearmand` 的最新源代码，解压缩这个 tarball，构建并安装此代码（安装需要有超级用户的权限，即根用户权限）。

```
$ wget http://launchpad.net/gearmand/trunk/\
0.10/+download/gearmand-0.10.tar.gz
$ tar xvzf gearmand-0.10.tar.gz
$ cd gearmand-0.10
$ ./configure
$ make
$ sudo make install

```

Show moreShow more icon

安装 `gearmand` 后，构建 PHP 扩展。您可以从 PECL 获取这个 tarball，也可以从 Github 复制该存储库。

```
$ wget http://pecl.php.net/get/gearman-0.6.0.tgz
$ cd pecl-gearman
#
# or
#
$ git clone git://github.com/php/pecl-gearman.git
$ cd pecl-gearman

```

Show moreShow more icon

有了这些代码后，就可以开始构建扩展了：

```
$ phpize
$ ./configure
$ make
$ sudo make install

```

Show moreShow more icon

这个 Gearman 守护程序通常被安装在 /usr/sbin。可以从命令行直接启动此守护程序，也可以将这个守护程序添加到启动配置中，以便在机器每次重启时就可以启动这个守护程序。

接下来，需要安装 Gearman 扩展。打开 php.ini 文件（可以通过 `php --ini` 命令快速找到这个文件），然后添加代码行 `extension = gearman.so` ：

```
$ php --ini
Loaded Configuration File:         /etc/php/php.ini
$ vi /etc/php/php.ini
...
extension = gearman.so

```

Show moreShow more icon

保存此文件。要想验证扩展是否启用，请运行 `php --info` ，然后查找 Gearman：

```
$ php --info | grep "gearman support"
gearman
gearman support => enabled
libgearman version => 0.10

```

Show moreShow more icon

此外，还可以用一个 PHP 代码片段来验证构建和安装是否得当。将这个小应用程序保存到 verify\_gearman.php：

```
<?php
print gearman_version() . "\n";
?>

```

Show moreShow more icon

接下来，从命令行运行此程序：

```
$ php verify_gearman.php
0.10

```

Show moreShow more icon

如果这个版本号与之前构建和安装的 Gearman 库的版本号相匹配，那么系统就已准备好了。

## 运行 Gearman

我们前面提到过，一个 Gearman 配置有三个角色：

- 一个或多个 producer 生成工作请求。每个工作请求命名它所想要的函数，例如 `email_all` 或 `analyze` 。
- 一个或多个 consumer 完成请求。每个 consumer 命名它所提供的一个或多个函数并向 agent 注册这些功能。一个 consumer 也可以被称为是一个 _worker_ 。
- 代理对与之建立连接的那些 consumer 提供的所有服务进行集中编制。它将 producer 与恰当的 consumer 联系起来。

借助如下的命令行，可以立即体验 Gearman：

1. 启动这个 agent，即 Gearman 守护程序：





    ```
    $ sudo /usr/sbin/gearmand --daemon

    ```





    Show moreShow more icon

2. 用命令行实用工具 `gearman` 运行一个 worker。这个 worker 需要一个名字并能运行任何命令行实用工具。例如，可以创建一个 worker 来列出某个目录的内容。 `-f` 参数命名了该 worker 所提供的函数：





    ```
    $ gearman -w -f ls -- ls -lh

    ```





    Show moreShow more icon

3. 最后的一个代码块是一个 producer，或用来生成查找请求的一个作业。也可以用 `gearman` 生成一个请求。同样，用 `-f` 选项来指定想要从中获得帮助的那个服务：





    ```
    $ gearman -f ls < /dev/null
    drwxr-xr-x@ 43 supergiantrobot  staff   1.4K Nov 15 15:07 gearman-0.6.0
    -rw-r--r--@  1 supergiantrobot  staff    29K Oct  1 04:44 gearman-0.6.0.tgz
    -rw-r--r--@  1 supergiantrobot  staff   5.8K Nov 15 15:32 gearman.html
    drwxr-xr-x@ 32 supergiantrobot  staff   1.1K Nov 15 14:04 gearmand-0.10
    -rw-r--r--@  1 supergiantrobot  staff   5.3K Jan  1  1970 package.xml
    drwxr-xr-x  47 supergiantrobot  staff   1.6K Nov 15 14:45 pecl-gearman

    ```





    Show moreShow more icon


## 从 PHP 使用 Gearman

从 PHP 使用 Gearman 类似于之前的示例，惟一的区别在于这里是在 PHP 内创建 producer 和 consumer。每个 consumer 的工作均封装在一个或多个 PHP 函数内。

[清单 1\. Worker.php](#清单-1-worker-php) 给出了用 PHP 编写的一个 Gearman worker。将这些代码保存在一个名为 _worker.php_ 的文件中。

##### 清单 1\. Worker.php

```
<?php
$worker= new GearmanWorker();
$worker->addServer();
$worker->addFunction("title", "title_function");
while ($worker->work());

function title_function($job)
{
    return ucwords(strtolower($job->workload()));
}
?>

```

Show moreShow more icon

[清单 2\. Client.php](#清单-2-client-php) 给出了用 PHP 编写的一个 producer，或 _client_ 。将此代码保存在一个名为 _client.php_ 的文件内。

##### 清单 2\. Client.php

```
<?php
$client= new GearmanClient();
$client->addServer();
print $client->do("title", "AlL THE World's a sTagE");
print "\n";
?>

```

Show moreShow more icon

现在，可以用如下的命令行连接客户机与 worker 了：

```
$ php worker.php &
$ php client.php
All The World's A Stage
$ jobs
[3]+  Running                 php worker.php &

```

Show moreShow more icon

这个 worker 应用程序继续运行，准备好服务另一个客户机。

## Gearman 的高级特性

在一个 Web 应用程序内可能有许多地方都会用到 Gearman。可以导入大量数据、发送许多电子邮件、编码视频文件、挖据数据并构建一个中央日志设施 — 所有这些均不会影响站点的体验和响应性。可以并行地处理数据。而且，由于 Gearman 协议是独立于语言和平台的，所以您可以在解决方案中混合编程语言。比如，可以用 PHP 编写一个 producer，用 `C` 、Ruby 或其他任何支持 Gearman 库的语言编写 worker。

一个连接客户机和 worker 的 Gearman 网络实际上可以使用任何您能想象得到的结构。很多配置能够运行多个代理并将 worker 分配到许多机器上。负载均衡是隐式的：每个可操作的可用 worker（可能是每个 worker 主机具有多个 worker）从队列中拉出作业。一个作业能够同步或异步运行并具有优先级。

Gearman 的最新版本已经将系统特性扩展到了包含持久的作业队列和用一个新协议来通过 HTTP 提交工作请求。对于前者，Gearman 工作队列保存在内存并在一个关系型数据库内存有备份。这样一来，如果 Gearman 守护程序故障，它就可以在重启后重新创建这个工作队列。另一个最新的改良通过一个 memcached 集群增加队列持久性。memcached 存储也依赖于内存，但被分散于几个机器以避免单点故障。

Gearman 是一个刚刚起步却很有实力的工作分发系统。据 Gearman 的作者 Eric Day 介绍，Yahoo! 在 60 或更多的服务器上使用 Gearman 每天处理 600 万个作业。新闻聚合器 Digg 也已构建了一个相同规模的 Gearman 网络，每天可处理 400,000 个作业。Gearman 的一个出色例子可以在 Narada 这个开源搜索引擎（参见参考资源）中找到。

Gearman 的未来版本将收集并报告统计数据、提供高级监视和缓存作业结果等。为了跟踪这个 Gearman 项目，可以订阅它的 Google 组，或访问 Freenode 上它的 IRC 频道 `#gearman` 。