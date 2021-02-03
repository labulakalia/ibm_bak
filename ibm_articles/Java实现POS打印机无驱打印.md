# Java 实现 POS 打印机无驱打印
探讨使用 Java 实现对 POS 打印机的直接控制

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-lo-pos/)

于丙超

发布: 2009-06-29

* * *

## 行业需求

我们是一家专业做酒店餐饮软件的公司，餐饮软件一个重要的功能就是后厨打印问题，前台点菜完毕，后厨立刻打印出单子，这样就减少人工递单的麻烦，节省时间，提高翻台率。这种信息化解决方案对打印技术要求很高，理论上最好 100% 不丢单，也就是每次点菜后厨都会相应出单子，但是实际上行不通，为什么呢？因为网线、打印机、网卡等都有可能有问题，别说打印机等硬件因为厨房油烟问题损坏，我们甚至碰到过网线被老鼠咬断的情况，总之硬件网络故障防不胜防，所以只能退而求其次，就是有问题不可怕，程序能够判断是否出了问题，并能给出提示，便于服务员处理，及时补单。

如果我们用安装 Windows 驱动的方法来实现后厨打印，那么肯定是不行的，因为我们只能单向向驱动程序抛包，不能从驱动程序获得任何返回值，没有办法了解是否打印成功。而且更为严重的是，有时候因为后厨打印机过多，Windows 驱动甚至会因为网络堵塞自作主张将包丢弃，没有任何提示。

这在行业应用中是不行的，会给用户带来损失，所以想到了绕过 Windows 驱动，直接写端口的方法。

## 无驱打印的可行性

所谓直接写端口的方法，就是不用安装打印机驱动，不使用 PrinterJob 获得打印机的名字的方法进行打印。

众所周知，之所以安装打印机驱动，一个重要的原因就是打印机厂商千差万别，不同的打印机往往都有各自的驱动，很难实现万能驱动。但是，在 POS 打印机行业却有一条捷径，就是现在市面上的 POS 打印机基本上都支持爱普生指令，也就是说，只要将程序和打印机联通，直接向端口里面写爱普生指令就可以控制打印机。

打印机接受到爱普生指令以后，自行进行解析，然后打印出相应的内容。

## 爱普生指令

日本的 EPSON 公司在目前的 POS 打印机市场，尤其是针式打印机市场占有很大一部分份额。它所推行的 ESC 打印控制命令 (EPSON StandardCode for Pr5nter) 已经成为了针式打印机控制语言事实上的工业标准，ESC/POS 打印命令集是 ESC 打印控制命令的简化版本，现在大多数 POS 打印都采用 ESC/POS 指令集。绝大多数打印机都有 EPSON ESC 的软件命令仿真功能，而且其它打印控制命令的格式和功能也都与 ESC 代码集类似。

由于早期的操作系统 DOS 与现在 Windows 的结构不同，在打印机内部软件和应用软件之间没有由硬件厂商提供的打印驱动程序，必须由应用软件直接通过硬件接口来控制打印机，所以从 ESC 指令出现开始，它就是公开的，否则没有应用软件可以使用它，而除了标准的 ESC 指令外，每种型号的打印机其指令又不太一样，所以在 DOS 软件中，你可以看到每个应用软件都只是支持为数不多的几种常用打印机。

ESC 指令在形式上分为两种格式，一种是文本方式控制码，一种是 Escape 转义序列码。文本方式控制码由一字节字符码表示，实现的是与打印机硬件操作有关的指令，Escape 序列码由转义字符和参数字符或打印数据组成。

## 建立打印连接

通过上面的介绍，了解了实现无驱打印原来只是一层窗户纸，具体的方法就是首先建立打印机连接，然后写入爱普生指令即可。那么如何建立打印机连接？以网口 POS 打印机举例。

第一步，首先要给网口打印机赋一个 IP 地址，例如叫做 192.168.0.18 。

第二步，编写连接代码。

```
Socket client=new java.net.Socket();
PrintWriter socketWriter;
client.connect(new InetSocketAddress("192.168.0.18" , 9100),1000); // 创建一个 socket
socketWriter = new PrintWriter(client.getOutputStream());// 创建输入输出数据流

```

Show moreShow more icon

看起来跟一般的 socket 连接没有很大的区别，就是赋一个 IP 地址和一个端口号，并设置一下超时时间即可，只需要说明的是，一般 POS 打印机的端口都是 9100 。

## 写入打印内容

连接建立完毕，写入内容就非常容易，只要使用 write 或者 println 方法写入即可，其中 write 方法是写入数字或字符，println 写入一行字符串。

例如：写入数字 socketWriter.write(0);

写入一行字符串 socketWriter.println( “巧富餐饮软件后厨单据” );

再入一行字符串 socketWriter.println( “桌位 14 桌，人数 3 ” );

再入一行字符串 socketWriter.println( “跺脚鱼头 1 份” );

您或许有疑问？内容已经成功写入，好像我们还没有用到爱普生指令。是的，如果只是普通的写入内容，不需要用到爱普生指令，爱普生指令主要帮助实现放大字体，自动走纸，打印条形码等功能。

## 放大字体

放大字体需要用到爱普生的 0x1c 指令，使用爱普生指令的方法很简单，只要向端口写入指令即可，例如：

```
socketWriter.write(0x1c);

```

Show moreShow more icon

注意 0x1c，是 16 进制的数字，当然也可以转换成 10 进制来写。需要说明的是，使用爱普生指令放大字体不能随意放大，因为它不是图形化打印，而是文本化打印，所以纵向或者横向只能按照倍数放大，不能矢量放大。例如在 POS58 打印机上将”巧富餐饮软件”几个字放大打印，可以有如下放大方法。

```
/* 横向放大一倍 */
socketWriter.write(0x1c);
socketWriter.write(0x21);
socketWriter.write(4);
/* 纵向放大一倍 */
socketWriter.write(0x1c);
socketWriter.write(0x21);
socketWriter.write(8);
/* 横向纵向都放大一倍 */
socketWriter.write(0x1c);
socketWriter.write(0x21);
socketWriter.write(12);

```

Show moreShow more icon

一般情况下，我们倾向采用纵向放大一倍的方法，放大后的字体看起来有点像仿宋体，视觉效果还不错。

## 兼容多种类型打印机

现在知道了使用爱普生指令的方法，所以只要有一本爱普生指令手册在手里，就可以用 Java 控制打印机进行无驱打印。但是现在问题是，同样是爱普生指令，不同的 pos 打印机可能不一样，就拿放大字体来说，pos58 打印机和 pos80 打印机指令就不尽相同。这时候怎么办呢？如何兼容多种类型打印机？

比如说，有的打印机并不是使用 0x1c 作为放大指令，而是使用 0x1b 作为放大指令，怎么办？容易。

```
/* 横向放大一倍 */
socketWriter.write(0x1c);
socketWriter.write(0x21);
socketWriter.write(4);
socketWriter.write(0x1b);
socketWriter.write(0x21);
socketWriter.write(4);
/* 纵向放大一倍 */
socketWriter.write(0x1c);
socketWriter.write(0x21);
socketWriter.write(8);
socketWriter.write(0x1b);
socketWriter.write(0x21);
socketWriter.write(8);
/* 横向纵向都放大一倍 */
socketWriter.write(0x1c);
socketWriter.write(0x21);
socketWriter.write(12);
socketWriter.write(0x1b);
socketWriter.write(0x21);
socketWriter.write(12);

```

Show moreShow more icon

看明白了吗？就是写两遍就行，因为如果 0x1b 指令若不存在，打印机自动将其抛弃。

## 实现自动走纸

POS 打印机因为出纸口有一些深度，打印完毕为了避免撕裂文字内容，一般需要适当走纸才行，当然可以使用爱普生指令来走纸，但是这样并不稳妥，为什么呢 ? 因为要考虑 POS 机的兼容性，所以一般采用打印空行的方式实现走纸。

```
for(int i=0;i<10;i++){
    socketWriter.println(" ");// 打印完毕自动走纸
}

```

Show moreShow more icon

显然，打印空行的方式有更好地兼容性。

## 打印条形码

条形码在各个行业中现在有广泛的应用，所以让打印机打印条形码是非常重要的功能，不过你不需要费好多精力去研究条形码知识，因为爱普生指令中有一个打印条形码指令，例如我们要打印条形码” 091955826335 ”，只要使用如下命令即可。

```
socketWriter.write(0x1d);
socketWriter.write(0x68);
socketWriter.write(120);
socketWriter.write(0x1d);
socketWriter.write(0x48);
socketWriter.write(0x01);
socketWriter.write(0x1d);
socketWriter.write(0x6B);
socketWriter.write(0x02);
socketWriter.println "091955826335" );
socketWriter.write(0x00);

```

Show moreShow more icon

## 完整的代码

好了，下面举一个完整的例子，我们来建立一个叫做 print 的方法，向某个打印机打印一个字符和条形码，并实现自动走纸，代码如下：

```
private boolean print(String ip, int port, String str,String code,int skip)
throws Exception{
    Socket client=new java.net.Socket();
    PrintWriter socketWriter;
    client.connect(new InetSocketAddress(ip,port),1000); // 创建一个 socket
    socketWriter = new PrintWriter(client.getOutputStream());// 创建输入输出数据流
    /* 纵向放大一倍 */
    socketWriter.write(0x1c);
    socketWriter.write(0x21);
    socketWriter.write(8);
    socketWriter.write(0x1b);
    socketWriter.write(0x21);
    socketWriter.write(8);
    socketWriter.println(str);
    // 打印条形码
    socketWriter.write(0x1d);
    socketWriter.write(0x68);
    socketWriter.write(120);
    socketWriter.write(0x1d);
    socketWriter.write(0x48);
    socketWriter.write(0x01);
    socketWriter.write(0x1d);
    socketWriter.write(0x6B);
    socketWriter.write(0x02);
    socketWriter.println(code);
    socketWriter.write(0x00);
    for(int i=0;i<skip;i++){
        socketWriter.println(" ");// 打印完毕自动走纸
    }
}

```

Show moreShow more icon

## 结束语

本文虽然只是讲述了网口打印机的直接写端口方式，似乎对并口打印机无效，其实不是这样，并口打印机只要接一个打印服务器就可以用了，缺点就是非一体机，然后还要安装打印服务器驱动。

这种无驱打印在非常广泛的范围内可以得到应用，包括餐饮、超市、医药等等其他需要用到 POS 打印机的行业。