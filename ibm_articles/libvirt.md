# libvirt
使用 Python 编写 KVM 脚本，第一部分

**标签:** Python

[原文链接](https://developer.ibm.com/zh/articles/os-python-kvm-scripting1/)

Paul Ferrill

发布: 2012-02-20

* * *

虚拟化是目前市场上大多数服务器操作系统的标准设备。在 Linux® 的世界里，服务器虚拟化有两个主要选择：基于 Kernel 的虚拟机 (KVM) 和 Xen。KVM 是 Red Hat 和其他公司采用的主要技术。虽然 Citrix 拥有 Xen，但是大多数核心功能是向外公开的。

虚拟机管理器（VMM 或 `virt-manager` ）项目为管理 KVM 和 Xen 虚拟机 (VM) 实例的创建和运行提供了一个工具。VMM 是用 Python 编写的，它使用 GTK+ 库来完成图形用户界面的构造。真实的构造是通过 `libvirt` 库完成的，本系列文章中会用到这个库。虽然 `libvirt` 是由 Red Hat 赞助研发的产品，但它是为一个开源项目，可通过 GNU Lesser General Public License 获得。

`libvirt` 由几个不同的部分组成，其中包括应用程序编程接口 (API) 库、一个守护进程 (`libvirtd`)，以及一个默认命令行实用工具 (`virsh`)。出于本文的目的，所有测试均使用 Ubuntu Server version 11.04 完成。 [安装与设置](#安装与设置) 部分介绍了我通过配置服务器来开发此处展示的脚本的所有步骤。第 1 部分介绍了 `libvirt` 和 KVM 虚拟化的基础知识，以及引发您的兴趣的一些命令行脚本。第 2 部分将深入探讨如何使用 `libvirt` 、Python 和 wxPython 构建您自己的虚拟化管理工具。

## 准备开始

在深入研究实际代码示例之前，先让我们复习一些关于虚拟化和 KVM 的术语和概念。当您在 Ubuntu Server 11.04 等服务器上安装 KVM 时，同时也在创建一个 _虚拟化主机_ 或 _虚拟机管理程序_ 。这意味着，您的服务器可以管理运行在 KVM 主机上的多个来宾操作系统 (guest operating system)。每个惟一的来宾均被称为一个 _域_ 和函数，它们以与您期望的大致相同的方式出现在个人机器上的单个服务器实例中。您可以通过安全外壳 (SSH) 或虚拟网络计算连接至服务器，就像连接至物理机器一样。

虽然 KVM 的功能与虚拟机管理程序或来宾管理器相似，但是 QEMU 提供了实际的机器仿真，这意味着 QEMU 可以执行目标机器的本机指令集。对于 x86 来宾操作系统，该执行会转换成为底层硬件上的直接执行的本机指令。对于其他架构，如 ARM，必须执行转换进程。KVM 和 QEMU 的组合提供了虚拟化所有目前可用以及某些不再可用的操作系统所需的所有支持函数。

_来宾域_ 包含若干个文件，其中包括一个或多个磁盘映像文件和一个基于 XML 的配置文件。该设置使得管理多台 VM 变得极其简单，方法是先创建一个基准系统映像，然后修改配置文件，使之符合您的需求。配置 KVM/QEMU 并与之通信的一个方法是使用 `libvirt` 工具包。一些供应商已经基于 `libvirt` 标准化其管理产品。

下面来看一个典型的域配置文件的内容。 [清单 1](#清单-1-设备-xml-定义) 显示了来自 `libvirt` 示例的 testdev.xml 文件。

##### 清单 1\. 设备 XML 定义

```
<device>
<name>File_test_device</name>
<capability type='system'>
         <hardware>
               <vendor>Libvirt</vendor>
               <version>Test driver</version>
               <serial>123456</serial>
               <uuid>11111111-2222-3333-4444-555555555555</uuid>
         </hardware>
         <firmware>
               <vendor>Libvirt</vendor>
               <version>Test Driver</version>
               <release_date>01/22/2007</release_date>
         </firmware>
</capability>
</device>

```

Show moreShow more icon

从 [清单 2](#清单-2-domfv0-xml-设备定义文件) 显示的测试 domfv0.xml 文件中，您可以看到关于配置虚拟设备的更多详细信息。

##### 清单 2\. domfv0.xml 设备定义文件

```
<devices>
<emulator>/usr/lib/xen/bin/qemu-dm</emulator>
         <interface type='bridge'>
               <source bridge='xenbr0'/>
               <mac address='00:16:3e:5d:c7:9e'/>
               <script path='vif-bridge'/>
         </interface>
         <disk type='file'>
               <source file='/root/fv0'/>
               <target dev='hda'/>
         </disk>
         <disk type='file' device='cdrom'>
               <source file='/root/fc5-x86_64-boot.iso'/>
               <target dev='hdc'/>
               <readonly/>
         </disk>
         <disk type='file' device='floppy'>
               <source file='/root/fd.img'/>
               <target dev='fda'/>
         </disk>
         <graphics type='vnc' port='5904'/>
</devices>

```

Show moreShow more icon

这里的重点是，您可以相对轻松地读取这些文件并继，然后创建自己的文件。虽然您可以手动构建任意数量的配置文件，但也可以使用像 Python 这样的脚本语言来自动化构建过程。

## 安装与设置

由于本文谈论的是为 KVM 编写脚本，所以有一个基础假设：您拥有一个安装了 KVM 的服务器。在使用 Ubuntu Server 11.04 的情况下，通过选择 **Software selection** 屏幕上的 **Virtual Machine Host** 选项，您在设置过程中有一个安装虚拟化的选项。您可能还想要选择 OpenSSH 服务器，因为您想远程连接至该机器。

第一项任务是安装最新版的 `libvirt` 。为此，您必须执行一些命令行命令。当安装 Ubuntu Server 11.04 之后，获得的 `libvirt` 版本为 0.8.8。而 `libvirt` 网站上提供的最新、最强大的版本是 0.9.5。要安装较新的版本，则需要将一个 Personal Package Archive (PPA) 储存库添加到包含新版 `libvirt` 的系统。在 launchpad.net 网站上快速搜索 `libvirt` 会显示一些可能的候选对象。在尝试执行更新之前查看储存库详细信息页面非常重要，因为一些存储库中可能包含损坏的软件包。Ubuntu Virtualization Team 维护了一个包含一些软件包（包括 `libvirt` ）的 PPA 存储库。在撰写本文时，libvirt 的最新版本为 0.9.2-4。

请执行以下步骤来安装此版本：

1. 安装 python-software-properties 软件包如下所示：





    ```
    sudo apt-get install python-software-properties

    ```





    Show moreShow more icon

    该命令使得 `add-apt-repository` 命令变得可用，您需要引用第三方来源。

2. 输入以下命令：





    ```
    sudo add-apt-repository ppa:ubuntu-virt/ppa
    sudo apt-get update
    sudo apt-get install libvirt-bin

    ```





    Show moreShow more icon

3. 由于您要使用 Python 完成本文中的所有脚本编写，所以请安装 IDLE Shell，它会使脚本的编写和测试变得更容易。

    此步骤假设：您已在 Ubuntu 服务器上安装了桌面环境。安装桌面的最快捷方法是使用下列的命令：





    ```
    sudo apt-get install ubuntu-desktop

    ```





    Show moreShow more icon


完成上述操作之后，您就可以访问任意数量的图形应用程序和 Ubuntu 软件安装程序。您可以使用 Ubuntu Software Center 来安装 Python IDLE 工具。

## 样例脚本

此时，在深入研究代码之前，让我们先来看一下使用 `libvirt` 的基础知识。应用程序与 `libvirt` 库之间的通信使用了一个简单的远程过程调用机制，该机制使得构建应用程序、从而通过 TCP/IP 连接与远程虚拟机管理程序进行通信成为可能。统一资源标识符（URI，由 Internet Engineering Task Force [IETF] Request for Comments [RFC] 2396 定义）用于标识您想用来建立连接的特定虚拟机管理程序。

虽然有些远程连接需要进行身份验证，但本地连接通常不需要进行身份验证。libvirt.conf 文件控制着安全配置。通过与某个域进行通信来进行最为广泛的控制是通过网络过滤来实现的。下面列出了如何使用过滤器控制网络通讯量的一个示例：

```
<devices>
    <interface type='bridge'>
      <mac address='00:16:3e:5d:c7:9e'/>
      <filterref filter='clean-traffic'/>
    </interface>
</devices>

```

Show moreShow more icon

这个代码片段定义了一个名为 `clean-traffic` 的过滤器，可用它来计算通过指定的媒体访问控制 (MAC) 地址的所有网络通讯量。如果您查看一下 `clean-traffic` XML，就会发现它包含以下内容：

```
<filter name='clean-traffic' chain='root'>
     <uuid>6f145c54-e3de-4c33-544a-70b69c16d9da</uuid>
     <filterref filter='no-mac-spoofing'/>
     <filterref filter='no-ip-spoofing'/>
     <filterref filter='allow-incoming-ipv4'/>
     <filterref filter='no-arp-spoofing'/>
     <filterref filter='no-other-l2-traffic'/>
     <filterref filter='qemu-announce-self'/>
</filter>

```

Show moreShow more icon

过滤器功能有大量完整的相关文档记录。如果您想要 `libvirt` 文档和样例文件的本地副本，只需使用一条命令即可。在这里，您只需执行以下命令即可：

```
sudo apt-get install libvirt-doc

```

Show moreShow more icon

操作完成后，所有的文档都会出现在 /usr/share/doc/libvirt-doc 目录中。稍后您会看到一些 Python 示例。如果您已经更新至较新的 `libvirt` 版本，那么您可能需要明确安装 Python 绑定。完成此项操作只需执行一条命令：

```
sudo apt-get install python-libvirt

```

Show moreShow more icon

可以使用 Python 的 IDLE 控制台查看 Python 代码，从而使用本地 QEMU 实例建立连接，然后查看已定义的域。 [清单 3](#清单-3-在-idlev-控制台中查看-python-代码) 显示了您使用此方法会看到的内容。

##### 清单 3\. 在 IDLEV 控制台中查看 Python 代码

```
Python 2.7.1+ (r271:86832, Apr 11 2011, 18:13:53)
[GCC 4.5.2] on linux2
Type "copyright", "credits" or "license()" for more information.
==== No Subprocess ====
>>> import libvirt
>>> conn=libvirt.open("qemu:///system")
>>> names = conn.listDefinedDomains()
>>> print names
['Test1', 'SBSLite', 'UbuntuServer1104', 'Win7_64-bit']]
>>>

```

Show moreShow more icon

此代码显示了如何获取所有已定义域的列表。 `listDefinedDomains()` 函数的返回结果中显示了一个包含 4 个已命名域的列表。建立到虚拟机管理程序的连接之后，就可以访问可用函数的细目清单。以下是一个较短的脚本，显示了如何获取 `conn` 对象上所有可用函数的列表：

```
clist = dir(conn)
for item in clist:
            print item

```

Show moreShow more icon

要查看已定义过滤器的列表，您可以使用相似的方法：

```
filts = conn.listNWFilters()
for item in filts:
            print item

```

Show moreShow more icon

IDLE 工具是一个很不错的方法，可用它来调查各种 API 调用，并在执行调用之后快速查看返回的结果。有些函数只能在运行域上执行。Python `dir()` 函数会返回指定对象的有效属性列表。它是一个便利的命令行工具，可快速查看特定对象提供的属性。您可以如上所示，在与虚拟机管理程序创建连接之后，使用该工具来获得可用函数列表。

要进行演示，您可以在 IDLE 控制台中使用几行 Python 代码，通过这些了解您可以在特定域上执行哪些类型的操作。 [清单 4](#清单-4-域对象的-python-输出) 提供了您能执行的操作的一个示例。

##### 清单 4\. 域对象的 Python 输出

```
>>> import libvirt
>>> import pprint
>>> conn=libvirt.open("qemu:///system")
>>> p = conn.lookupByName('ubuntu100403')
>>> pprint.pprint(dir(p))
['ID',
'OSType',
'UUID',
'UUIDString',
'XMLDesc',
'__del__',
'__doc__',
'__init__',
'__module__',
'_conn',
'_o',
'abortJob',
'attachDevice',
'attachDeviceFlags',
'autostart',
'blkioParameters',
'blockInfo',
'blockPeek',
'blockStats',
'connect',
'coreDump',
'create',
'createWithFlags',
'destroy',
'detachDevice',
'detachDeviceFlags',
'hasCurrentSnapshot',
'hasManagedSaveImage',
'info',
'injectNMI',
'interfaceStats',
'isActive',
'isPersistent',

```

Show moreShow more icon

您可以使用这种基本方法来构建一个简单的脚本，以列出有关所有运行域的信息。可以使用 `listDomainsID()` 和 `lookupByID()` 函数调用来完成大多数的工作，如 [清单 5](#清单-5-用-python-编写的列出域的脚本) 所示。

##### 清单 5\. 用 Python 编写的列出域的脚本

```
import libvirt
conn=libvirt.open("qemu:///system")

for id in conn.listDomainsID():
dom = conn.lookupByID(id)
infos = dom.info()
print 'ID = %d' % id
print 'Name =  %s' % dom.name()
print 'State = %d' % infos[0]
print 'Max Memory = %d' % infos[1]
print 'Number of virt CPUs = %d' % infos[3]
print 'CPU Time (in ns) = %d' % infos[2]
print ' '

```

Show moreShow more icon

在此脚本的输出中，有一个域处于激活状态，另一个为挂起状态，如下所示：

```
ID = 3
Name =  ubuntu100403
State = 3
Max Memory = 1048576
Number of virt CPUs = 1
CPU Time (in ns) = 1048576

ID = 4
Name =  Win7_64-bit
State = 1
Max Memory = 2097152
Number of virt CPUs = 2
CPU Time (in ns) = 2097152

```

Show moreShow more icon

`libvirt` 还为所有的类和方法实现了 Python 文档字符串。您可以访问其信息，方法是输入 `help (libvirt)` 获得顶级帮助信息，或通过 `help(libvirt.class)` 获取特定的类。您必须在输入 `help()` 命令前先导入 `libvirt` 模块。我所测试的版本实现下列 11 个类：

- `libvirtError`
- `virConnect`
- `virDomain`
- `virDomainShapshot`
- `virInterface`
- `virNWFilter`
- `virNetwork`
- `virSecret`
- `virStoragePool`
- `virStorageVol`
- `virStream`

此列表能帮助您从 Python 解码访问 `libvirt` 函数的语法。该列表还为您提供所有已命名常数的列表，比如 `VIR_DOMAIN_RUNNING` ，该常数等同 1。上面使用过的一些函数（如 `dom.info()` ）会返回一个整数值，您需要对照常数表对其进行解码。

## 利用实用工具脚本实现自动化

您可以使用 `libvirt` 和 Python 编写一些脚本来管理 KVM 安装。虽然对于少部分域而言这种方法是无效的，但在域的数量达到两位数时，这样做可以快速节省时间。一个简单的任务就可能让所有域映像的静态 IP 地址做出大量修改。您可以通过迭代所有 .conf 文件并进行相应修改来完成此任务。Python 拥有许多内置功能，可帮助完成此项任务。

[清单 6](#清单-6-网络配置-xml-文件) 显示了一个 XML 网络定义示例。

##### 清单 6\. 网络配置 XML 文件

```
<network>
<name>testnetwork</name>
<bridge name="virbr1" />
<forward/>
<ip address="192.168.100.1" netmask="255.255.255.0">
    <dhcp>
      <range start="192.168.100.2" end="192.168.100.254" />
      <host mac='de:af:de:af:00:02' name='vm-1' ip='192.168.100.2' />
      <host mac='de:af:de:af:00:03' name='vm-2' ip='192.168.100.3' />
      <host mac='de:af:de:af:00:04' name='vm-3' ip='192.168.100.4' />
      <host mac='de:af:de:af:00:05' name='vm-4' ip='192.168.100.5' />
      <host mac='de:af:de:af:00:06' name='vm-5' ip='192.168.100.6' />
      <host mac='de:af:de:af:00:07' name='vm-6' ip='192.168.100.7' />
      <host mac='de:af:de:af:00:08' name='vm-7' ip='192.168.100.8' />
      <host mac='de:af:de:af:00:09' name='vm-8' ip='192.168.100.9' />
      <host mac='de:af:de:af:00:10' name='vm-9' ip='192.168.100.10' />
    </dhcp
</ip>
</network>

```

Show moreShow more icon

如果想要将主要的子网络从 192.168.100 改为 192.168.200，可以在编辑器中打开配置文件并进行全局搜索和替换。当您想要做一些较复杂的操作时，比如向所有以 2 开始的 IP 地址和 MAC 地址添加 10，这个技巧非常有用。 [清单 7](#清单-7-更改-mac-和-ip-地址的-python-脚本) 显示了如何只使用 20 多行 Python 代码就能完成此项任务。

##### 清单 7\. 更改 MAC 和 IP 地址的 Python 脚本

```
#!/usr/bin/env python

from xml.dom.minidom import parseString
import sys

def main():
    target = sys.argv[1]
    number = int(sys.argv[2])

    xml = open(target, 'r').read()
    doc = parseString(xml)
    for host in doc.getElementsByTagName('host'):
        ip = host.getAttribute('ip')
        parts = ip.split('.')
        parts[-1] = str(int(parts[-1]) + number)
        host.setAttribute('ip', '.'.join(parts))

        mac = host.getAttribute('mac')
        parts = mac.split(':')
        parts[-1] = str(int(parts[-1]) + number)
        host.setAttribute('mac', ':'.join(parts))

    f = open(target, 'w')
    f.write(doc.toxml())
    f.close()

if __name__ == '__main__':
    main()

```

Show moreShow more icon

此脚本演示了您使用 Python Standard Library 时所展现的 Python 的强大功能。在这里，使用了 `xml.dom.minidom` 的 `parseString` 来完成解析 XML 文件的重任。拥有特定的 XML 属性后，只需使用 Python `string.split` 函数将其分成几块即可。接着，要进行匹配，并将字符串再连在一起。您可以扩展这个方法，将它用于任何 XML 文件的批量更改，这些文件中包括 `libvirt` 的 .conf 文件。

另一个有用的脚本将获取所有正在运行的域的快照。该脚本首先需要获取所有正在运行的域的列表，然后分别暂停每个域并获取每个域的快照。虽然此操作对生产环境不太实用，但是您可以将其设置为在午夜作为一个 `CRON` 作业运行。该脚本可以轻松实现目前为止突出显示的命令和 `snapshotCreateXML()` 调用。

## 结束语

本文只是粗浅地介绍了 `libvirt` 所包含的功能。请参阅 参考资料 部分，从那里获得更深入地研究 `libvirt` 和虚拟化的更多文章的链接。当您开始尝试通过实现代码来监视和管理您的环境时，掌握 KVM 基础知识只是一个开始。本系列文章的下一期文章将以此为基础，构建一些实际的虚拟管理工具。

本文翻译自： [libvirt](https://developer.ibm.com/articles/os-python-kvm-scripting1/)（2011-12-06）