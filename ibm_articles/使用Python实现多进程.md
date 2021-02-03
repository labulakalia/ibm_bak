# 使用 Python 实现多进程
学习使用 Python 2.6 管理一组进程

**标签:** Python

[原文链接](https://developer.ibm.com/zh/articles/au-multiprocessing/)

Noah Gift

发布: 2009-05-21

* * *

## 简介

在 IBM® Developer® 的 [早期文章](http://www.ibm.com/developerworks/cn/aix/library/au-threadingpython/index.html) 中，我演示了使用 Python 实现线程式编程的一种简单且有效的模式。但是，这种方法的一个缺陷就是它并不总是能够提高应用程序的速度，因为全局解释器锁（Global Interpreter Lock，GIL）将线程有效地限制到一个核中。如果需要使用计算机中的所有核，那么通常都需通过 对 经常使用 fork 操作来实现，从而提高速度。处理进程组是件困难的事情，因为为了在进程之间进行通信，需要对所有调用进行协调，这通常会使事情变得更复杂。

幸运的是，自 2.6 版本起，Python 包括了一个名为 “多进程（multiprocessing）” 的模块来帮助处理进程。该进程模块的 API 与线程 API 的工作方式有些相似点，但是也存在一些需要特别注意的不同之处。主要区别之一就是进程拥有的一些微妙的底层行为，这是高级 API 永远无法完全抽象出来的。

### fork 简介

进程和线程在并发性的工作原理方面存在一些明显的差异。在进程执行 fork 时，操作系统将创建具有新进程 ID 的新的子进程，复制父进程的状态（内存、环境变量等）。首先，在我们实际使用进程模块之前，先看一下 Python 中的一个非常基本的 fork 操作。

##### Listing 1. fork.py

```
#!/usr/bin/env python

"""A basic fork in action"""

import os

def myfork():
    childpid = os.fork()
    if childpid == 0:
        print "Child Process: PID#%s" % os.getpid()
    else:
        print "Parent Process: PID#%s" % os.getpid()

if _name == "__main":
    my_fork()
```

Show moreShow more icon

现在来看一下以上代码的输出：

```
mac% python fork.py
Parent Process: PID#5285
Child Process: PID#5286
```

Show moreShow more icon

在下一个示例中，增强初始 fork 的代码，并设置一个环境变量。该环境变量随后将被复制到子进程中。下面给出了相应的代码：

##### Listing 2. 示例 1. Python 中的 fork 操作

```
#!/usr/bin/env python

"""A fork that demonstrates a copied environment"""

import os
from os import environ

def myfork():
    environ['FOO']="baz"
    print "FOO environmental variable set to: %s" % environ'FOO'    environ['FOO']="bar"
    print "FOO environmental variable changed to: %s" % environ'FOO'    childpid = os.fork()
    if childpid == 0:
        print "Child Process: PID#%s" % os.getpid()
        print "Child FOO environmental variable == %s" % environ'FOO'    else:
        print "Parent Process: PID#%s" % os.getpid()
        print "Parent FOO environmental variable == %s" % environ'FOO'
if _name == "__main":
    my_fork()
```

Show moreShow more icon

下面给出了 fork 的输出：

```
mac% python env_fork.py
FOO environmental variable set to: baz
FOO environmental variable changed to: bar
Parent Process: PID#5333
Parent FOO environmental variable == bar
Child Process: PID#5334
Child FOO environmental variable == bar
```

Show moreShow more icon

在输出中，可以看到 “修改后的” 环境变量 FOO 留在了子进程和父进程中。您可以通过在父进程中再次修改环境变量来进一步测试这个示例，您将看到子进程现在是完全独立的，它有了自己的生命。注意，子进程模块也可用于 fork 进程，但是实现方式没有多进程模块那么复杂。

### 多进程简介

现在您已经了解 Python fork 操作的基本知识，让我们通过一个简单例子了解它在更高级的多进程库中的使用。在这个示例中，仍然会出现 fork，但是已经为我们处理了大部分标准工作。

##### Listing 3. 示例 2. 简单的多进程

```
#!/usr/bin/env python
from multiprocessing import Process
import os
import time

def sleeper(name, seconds):
print 'starting child process with id: ', os.getpid()
print 'parent process:', os.getppid()
print 'sleeping for %s ' % seconds
time.sleep(seconds)
print "Done sleeping"

if name == 'main':
print "in parent process (id %s)" % os.getpid()
p = Process(target=sleeper, args=('bob', 5))
p.start()
print "in parent process after child process start"
print "parent process about to join child process"
p.join()
print "in parent process after child process join"
print "parent process exiting with id ", os.getpid()
print "The parent's parent process:", os.getppid()
```

Show moreShow more icon

如果查看输出，将会看到下面的内容：

```
mac% python simple.py
in parent process (id 5245)
in parent process after child process start
parent process about to join child process
starting child process with id:  5246
parent process: 5245
sleeping for 5
Done sleeping
in parent process after child process join
parent process exiting with id  5245
The parent's parent process: 5231
```

Show moreShow more icon

可以看到从主进程分出了一个子进程，该子进程随后休眠了 5 秒种。子进程分配是在调用 `p.start()` 时发生的。在下一节中，您将看到这个基础的程序将扩展为更大的程序。

### 构建异步 Net-SNMP 引擎

到目前为止，您尚未构建任何特别有用的内容。下一个示例将解决一个实际问题，为 Net-SNMP 异步生成 Python 绑定。默认情况下，Net-SNMP 将阻塞每一个 Python 调用。使用多进程库可以非常简单地将 Net-SNMP 库转换为完全异步的操作。

在开始之前，需要检查是否安装了一些必备的内容，以便使用 Python 2.6 多进程库和 Net-SNMP 绑定：

1. 下载 Python 2.6 并针对所使用的操作系统进行编译： [Python 2.6 下载](http://www.python.org/download/releases/2.6/)
2. 调整 shell 路径，这样在输入 `python` 时就会启动 Python 2.6。例如，如果将 Python 编译到 /usr/local/bin/，您就需要预先处理 $PATH 变量，从而确保它位于一个较旧的 Python 版本之前。
3. 下载并安装设置工具： [设置工具](http://peak.telecommunity.com/DevCenter/setuptools)
4. 下载 Net-SNMP，除了使用其他操作系统所需的标记（参见相应的 README 文件）外，另外使用一个 “–with-python-modules” 标记进行配置。`./configure --with-python-modules`

按如下所示编译 Net-SNMP：

```
‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑
            Net‑SNMP configuration summary:
‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑‑

SNMP Versions Supported:    1 2c 3
Net‑SNMP Version:           5.4.2.1
Building for:               darwin9
Network transport support:  Callback Unix TCP UDP
SNMPv3 Security Modules:     usm
Agent MIB code:             default_modules =>  snmpv3mibs mibII ucd_snmp notification
notification‑log‑mib target agent_mibs agentx disman/event disman/schedule utilities
Embedded Perl support:      enabled
SNMP Perl modules:          building ‑‑ embeddable
SNMP Python modules:        building for /usr/local/bin//python
Authentication support:     MD5 SHA1
Encryption support:         DES AES
]]
```

Show moreShow more icon

查看以下模块的代码，您将随后运行它。

##### Listing 4. 示例 3. Net-SNMP 的多进程包装器

```
#!/usr/bin/env python2.6
"""
This is a multiprocessing wrapper for Net‑SNMP.
This makes a synchronous API asynchronous by combining
it with Python2.6
"""

import netsnmp
from multiprocessing import Process, Queue, currentprocess

class HostRecord():
    """This creates a host record"""
    def init(self,
                 hostname = None,
                 query = None):
        self.hostname = hostname
        self.query = query

class SnmpSession():
    """A SNMP Session"""
    def init(self,
                oid = "sysDescr",
                Version = 2,
                DestHost = "localhost",
                Community = "public",
                Verbose = True,
                ):
        self.oid = oid
        self.Version = Version
        self.DestHost = DestHost
        self.Community = Community
        self.Verbose = Verbose
        self.var = netsnmp.Varbind(oid, 0)
        self.hostrec = HostRecord()
        self.hostrec.hostname = self.DestHost

    def query(self):
        """Creates SNMP query

        Fills out a Host Object and returns result
        """
        try:
            result = netsnmp.snmpget(self.var,
                                Version = self.Version,
                                DestHost = self.DestHost,
                                Community = self.Community)
            self.hostrec.query = result
        except Exception, err:
            if self.Verbose:
                print err
            self.hostrec.query = None
        finally:
            return self.hostrec

def makequery(host):
    """This does the actual snmp query

    This is a bit fancy as it accepts both instances
    of SnmpSession and host/ip addresses.  This
    allows a user to customize mass queries with
    subsets of different hostnames and community strings
    """
    if isinstance(host,SnmpSession):
        return host.query()
    else:
        s = SnmpSession(DestHost=host)
        return s.query()

#Function run by worker processes
def worker(input, output):
    for func in iter(input.get, 'STOP'):
        result = makequery(func)
        output.put(result)

def main():
    """Runs everything"""

    #clients
    hosts = "localhost", "localhost"    NUMBEROF_PROCESSES = len(hosts)

    #Create queues
    task_queue = Queue()
    done_queue = Queue()

    #submit tasks
    for host in hosts:
        task_queue.put(host)

    #Start worker processes
    for i in range(NUMBER_OF_PROCESSES):
        Process(target=worker, args=(task_queue, done_queue)).start()

     #Get and print results
    print 'Unordered results:'
    for i in range(len(hosts)):
        print '\t', done_queue.get().query

    #Tell child processes to stop
    for i in range(NUMBER_OF_PROCESSES):
        task_queue.put('STOP')
        print "Stopping Process #%s" % i

if __name == "__main":
    main()
```

Show moreShow more icon

这里有两个类，一个 HostRecord 类和一个 SnmpSession 类。SnmpSession 类包含一个使用 Net-SNMP 的 SNMP 库实际执行查询的方法。由于调用一般都会进行阻塞，因此需要导入多进程库并使用 Process 运行它。此外，传入一个 task\_queue 和一个 done\_queue，这样可以同步并保护进出进程池的数据。如果对线程比较熟悉的话，将会注意到这种方式非常类似于线程 API 使用的方法。

需要特别关注一下主函数中 `#clients` 部分的主机列表。注意，可以对 50 或 100 台或更多主机运行异步 SNMP 查询，具体取决于当前使用的硬件。NUMBER\_OF\_PROCESSES 变量的设置非常简单，只是被设置为主机列表中的主机数。最终，最后两个部分在处理过程中从队列获取结果，然后将一个 “STOP” 消息放到队列中，表示可以终止进程。

如果在对 Net-SNMP 进行监听的 OS X 机器上运行代码，那么将会得到如下所示的非阻塞输出：

```
mac% time python multisnmp.py
Unordered results:
     ('Darwin mac.local 9.6.0 Darwin Kernel Version 9.6.0: Mon Nov 24 17:37:00 PST 2008;
root:xnu‑1228.9.59~1/RELEASE_I386 i386',)
     ('Darwin mac.local 9.6.0 Darwin Kernel Version 9.6.0: Mon Nov 24 17:37:00 PST 2008;
root:xnu‑1228.9.59~1/RELEASE_I386 i386',)
Stopping Process #0
Stopping Process #1
python multisnmp.py  0.18s user 0.08s system 107% cpu 0.236 total
```

Show moreShow more icon

### 配置 OS X 的 SNMPD

如果希望配置 OS X 的 SNMP Daemon 以针对本文进行测试，那么需要执行下面的操作。首先，在 shell 中使用三个命令重写配置文件：

```
            $ sudo cp /etc/snmp/snmpd.conf /etc/snmp/snmpd.conf.bak.testing
            $ sudo echo "rocommunity public" > /etc/snmp/snmpd.conf
            $ sudo snmpd
```

Show moreShow more icon

这将有效地备份您的配置，生成一个新配置，然后重新启动 snmpd。步骤与许多 UNIX 平台类似，但步骤 3 是除外，该步骤需要重新启动 snmpd，或发送一个 HUP。如果希望 OS X 在启动后永久运行 snmpd，那么可以按如下所示编辑这个 plist 文件：
`/System/Library/LaunchDaemons/org.net-snmp.snmpd.plist`

```
        <?xml version="1.0" encoding="UTF‑8"?> <!DOCTYPE plist PUBLIC "‑//Apple//DTD
PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList‑1.0.dtd">
plist version="1.0">
<dict>
    <key>Disabled</key>
         <false/>
    <key>KeepAlive</key>
    <true/>
    <key>Label</key>
    <string>org.net‑snmp.snmpd</string>
    <key>OnDemand</key>
    <false/>
    <key>Program</key>
    <string>/usr/sbin/snmpd</string>
    <key>ProgramArguments</key>
    <array>
        <string>snmpd</string>
            <string>‑f</string>
    </array>
     <key>RunAtLoad</key>
    <true/>
    <key>ServiceIPC</key>
    <false/>
</dict>
</plist>
```

Show moreShow more icon

如果希望对多台机器进行测试，那么使用下面的内容替换主机行就可以轻松执行修改：

```
hosts = ["192.168.1.100", SnmpSession(DestHost="example.com", Community="mysecret"),
"example.net", "example.org"]
```

Show moreShow more icon

运行作业的 worker 函数将获得两个字符串形式的主机名和完整的 SnmpSession 对象。

## 结束语

官方文档与多进程库一样有用，您应当特别关注其中提到的以下这些事项：避免共享状态；最好显式地连接所创建的进程；尽量避免终止具有共享状态的进程；最后确保在连接前删除队列中的所有队列项，否则将出现死锁。

除了以上的注意事项之外，多进程也是 Python 编程语言的一大增强。尽管 GIL 对线程的限制曾经被认为是一个弱点，但是通过包含强大灵活的多进程库，Python 不仅弥补了这个弱点，而且还得到了增强。非常感谢 David Goodger 担任本文的技术审校！

本文翻译自： [Multiprocessing with Python](https://developer.ibm.com/articles/au-multiprocessing/)（2009-03-24）