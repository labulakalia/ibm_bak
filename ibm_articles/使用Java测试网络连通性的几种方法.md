# 使用 Java 测试网络连通性的几种方法
Java 中三种不同的网络可达的判断方法以及针对 IPv4 和 IPv6 混合网络的编程方法

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-lo-connectiontest/)

吴校军, 刘冠群

发布: 2012-11-22

* * *

## 概述

在网络编程中，有时我们需要判断两台机器之间的连通性，或者说是一台机器到另一台机器的网络可达性。在系统层面的测试中，我们常常用 Ping 命令来做验证。尽管 Java 提供了比较丰富的网络编程类库（包括在应用层的基于 URL 的网络资源读取，基于 TCP/IP 层的 Socket 编程，以及一些辅助的类库），但是没有直接提供类似 Ping 命令来测试网络连通性的方法。本文将介绍如何通过 Java 已有的 API，编程实现各种场景下两台机器之间的网络可达性判断。在下面的章节中，我们会使用 Java 网络编程的一些类库 java.net.InetAddress 和 java.net.Socket，通过例子解释如何模拟 Ping 命令。

## 简单判断两台机器的可达性

一般情况下，我们仅仅需要判断从一台机器是否可以访问（Ping）到另一台机器，此时，可以简单的使用 Java 类库中 java.net.InetAddress 类来实现，这个类提供了两个方法探测远程机器是否可达

```
boolean isReachable(int timeout) // 测试地址是否可达 boolean
                isReachable(NetworkInterface netif, int ttl, int timeout) // 测试地址是否可达.

```

Show moreShow more icon

简单说来，上述方法就是通过远端机器的 IP 地址构造 InetAddress 对象，然后调用其 isReachable 方法，测试调用机器和远端机器的网络可达性。注意到远端机器可能有多个 IP 地址，因而可能要迭代的测试所有的情况。

##### 清单 1：简单判断两台机器的可达性

```
void isAddressAvailable(String ip){ try{ InetAddress
                address = InetAddress.getByName(ip);//ping this IP if(address instanceof
                java.net.Inet4Address){ System.out.println(ip + " is ipv4 address"); }else
                if(address instanceof java.net.Inet6Address){ System.out.println(ip + " is ipv6
                address"); }else{ System.out.println(ip + " is unrecongized"); }
                if(address.isReachable(5000)){ System.out.println("SUCCESS - ping " + IP + " with no
                interface specified"); }else{ System.out.println("FAILURE - ping " + IP + " with no
                interface specified"); } System.out.println("\n-------Trying different
                interfaces--------\n"); Enumeration<NetworkInterface> netInterfaces =
                NetworkInterface.getNetworkInterfaces(); while(netInterfaces.hasMoreElements()) {
                NetworkInterface ni = netInterfaces.nextElement(); System.out.println( "Checking
                interface, DisplayName:" + ni.getDisplayName() + ", Name:" + ni.getName());
                if(address.isReachable(ni, 0, 5000)){ System.out.println("SUCCESS - ping " + ip);
                }else{ System.out.println("FAILURE - ping " + ip); }
                Enumeration<InetAddress> ips = ni.getInetAddresses();
                while(ips.hasMoreElements()) { System.out.println("IP: " +
                ips.nextElement().getHostAddress()); }
                System.out.println("-------------------------------------------"); }
                }catch(Exception e){ System.out.println("error occurs."); e.printStackTrace(); } }

```

Show moreShow more icon

程序输出

```
--------------START-------------- 10.13.20.70 is ipv4 address
                SUCCESS - ping 10.13.20.70 with no interface specified -------Trying different
                interfaces-------- Checking interface, DisplayName:MS TCP Loopback interface,
                Name:lo FAILURE - ping 10.13.20.70 IP: 127.0.0.1
                ------------------------------------------- Checking interface, DisplayName:Intel(R)
                Centrino(R) Advanced-N 6200 AGN - Teefer2 Miniport, Name:eth0 FAILURE - ping
                10.13.20.70 IP: 9.123.231.40 ------------------------------------------- Checking
                interface, DisplayName:Intel(R) 82577LM Gigabit Network Connection - Teefer2
                Miniport, Name:eth1 SUCCESS - ping 10.13.20.70
                ------------------------------------------- Checking interface, DisplayName:WAN
                (PPP/SLIP) Interface, Name:ppp0 SUCCESS - ping 10.13.20.70 IP: 10.0.50.189
                ------------------------------------------- --------------END--------------

```

Show moreShow more icon

从上可以看出 isReachable 的用法，可以不指定任何接口来判断远端网络的可达性，但这不能区分出数据包是从那个网络接口发出去的 ( 如果本地有多个网络接口的话 )；而高级版本的 isReachable 则可以指定从本地的哪个网络接口测试，这样可以准确的知道远端网络可以连通本地的哪个网络接口。

但是，Java 本身没有提供任何方法来判断本地的哪个 IP 地址可以连通远端网络，Java 网络编程接口也没有提供方法来访问 ICMP 协议数据包，因而通过 ICMP 的网络不可达数据包实现这一点也是不可能的 ( 当然可以用 JNI 来实现，但就和系统平台相关了 ), 此时可以考虑本文下一节提出的方法。

## 指定本地和远程网络地址，判断两台机器之间的可达性

在某些情况下，我们可能要确定本地的哪个网络地址可以连通远程网络，以便远程网络可以回连到本地使用某些服务或发出某些通知。一个典型的应用场景是，本地启动了文件传输服务 ( 如 FTP)，需要将本地的某个 IP 地址发送到远端机器，以便远端机器可以通过该地址下载文件；或者远端机器提供某些服务，在某些事件发生时通知注册了获取这些事件的机器 ( 常见于系统管理领域 )，因而在注册时需要提供本地的某个可达 ( 从远端 ) 地址。

虽然我们可以用 InetAddress.isReachabl 方法判断出本地的哪个网络接口可连通远程玩过，但是由于单个网络接口是可以配置多个 IP 地址的，因而在此并不合适。我们可以使用 Socket 建立可能的 TCP 连接，进而判断某个本地 IP 地址是否可达远程网络。我们使用 java.net.Socket 类中的 connect 方法

```
void connect(SocketAddress endpoint, int timeout)
                //使用Socket连接服务器，指定超时的时间

```

Show moreShow more icon

这种方法需要远程的某个端口，该端口可以是任何基于 TCP 协议的开放服务的端口（如一般都会开放的 ECHO 服务端口 7, Linux 的 SSH 服务端口 22 等）。实际上，建立的 TCP 连接被协议栈放置在连接队列，进而分发到真正处理数据的各个应用服务，由于 UDP 没有连接的过程，因而基于 UDP 的服务（如 SNMP）无法在此方法中应用。

具体过程是，枚举本地的每个网络地址，建立本地 Socket，在某个端口上尝试连接远程地址，如果可以连接上，则说明该本地地址可达远程网络。

##### 程序清单 2：指定本地地址和远程地址，判断两台机器之间的可达性

```
void printReachableIP(InetAddress
                remoteAddr, int port){ String retIP = null;
                Enumeration<NetworkInterface> netInterfaces; try{ netInterfaces =
                NetworkInterface.getNetworkInterfaces(); while(netInterfaces.hasMoreElements()) {
                NetworkInterface ni = netInterfaces.nextElement();
                Enumeration<InetAddress> localAddrs = ni.getInetAddresses();
                while(localAddrs.hasMoreElements()){ InetAddress localAddr =
                localAddrs.nextElement(); if(isReachable(localAddr, remoteAddr, port, 5000)){ retIP
                = localAddr.getHostAddress(); break; } } } } catch(SocketException e) {
                System.out.println( "Error occurred while listing all the local network
                addresses."); } if(retIP == null){ System.out.println("NULL reachable local IP is
                found!"); }else{ System.out.println("Reachable local IP is found, it is " + retIP);
                } } boolean isReachable(InetAddress localInetAddr, InetAddress remoteInetAddr, int
                port, int timeout) { booleanisReachable = false; Socket socket = null; try{ socket =
                newSocket(); // 端口号设置为 0 表示在本地挑选一个可用端口进行连接 SocketAddress localSocketAddr = new
                InetSocketAddress(localInetAddr, 0); socket.bind(localSocketAddr); InetSocketAddress
                endpointSocketAddr = new InetSocketAddress(remoteInetAddr, port);
                socket.connect(endpointSocketAddr, timeout); System.out.println("SUCCESS -
                connection established! Local: " + localInetAddr.getHostAddress() + " remote: " +
                remoteInetAddr.getHostAddress() + " port" + port); isReachable = true; }
                catch(IOException e) { System.out.println("FAILRE - CAN not connect! Local: " +
                localInetAddr.getHostAddress() + " remote: " + remoteInetAddr.getHostAddress() + "
                port" + port); } finally{ if(socket != null) { try{ socket.close(); }
                catch(IOException e) { System.out.println("Error occurred while closing socket..");
                } } } return isReachable; }

```

Show moreShow more icon

运行结果

```
--------------START-------------- FAILRE - CAN not connect! Local:
                127.0.0.1 remote: 10.8.1.50 port22 FAILRE - CAN not connect! Local: 9.123.231.40
                remote: 10.8.1.50 port22 SUCCESS - connection established! Local: 10.0.50.189
                remote: 10.8.1.50 port22 Reachable local IP is found, it is 10.0.50.189
                --------------END--------------

```

Show moreShow more icon

## IPv4 和 IPv6 混合网络下编程

当网络环境中存在 IPv4 和 IPv6，即机器既有 IPv4 地址，又有 IPv6 地址的时候，我们可以对程序进行一些优化，比如

- 由于 IPv4 和 IPv6 地址之间是无法互相访问的，因此仅需要判断 IPv4 地址之间和 IPv6 地址之间的可达性。
- 对于 IPv4 的换回地址可以不做判断，对于 IPv6 的 Linklocal 地址也可以跳过测试
- 根据实际的需要，我们可以优先考虑选择使用 IPv4 或者 IPv6，提高判断的效率

##### 程序清单 3: 判断本地地址和远程地址是否同为 IPv4 或者 IPv6

```
// 判断是 IPv4 还是 IPv6 if(!((localInetAddr
                instanceofInet4Address) && (remoteInetAddr instanceofInet4Address)
                || (localInetAddr instanceofInet6Address) && (remoteInetAddr
                instanceofInet6Address))){ // 本地和远程不是同时是 IPv4 或者 IPv6，跳过这种情况，不作检测 break; }

```

Show moreShow more icon

##### 程序清单 4：跳过本地地址和 LinkLocal 地址

```
if( localAddr.isLoopbackAddress() ||
                localAddr.isAnyLocalAddress() || localAddr.isLinkLocalAddress() ){ // 地址为本地环回地址，跳过
                break; }

```

Show moreShow more icon

## 结束语

本文列举集中典型的场景，介绍了通过 Java 网络编程接口判断机器之间可达性的几种方式。在实际应用中，可以根据不同的需要选择相应的方法稍加修改即可。对于更加特殊的需求，还可以考虑通过 JNI 的方法直接调用系统 API 来实现，能提供更加强大和灵活的功能，这里就不再赘述了。