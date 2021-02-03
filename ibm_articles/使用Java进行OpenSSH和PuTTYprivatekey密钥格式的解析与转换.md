# 使用 Java 进行 OpenSSH 和 PuTTY private key 密钥格式的解析与转换
常用的两种 SSH 密钥的格式以及转换的原理

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-lo-opensshppk/)

罗 文刚

发布: 2015-01-29

* * *

## 概述

Secure Shell(SSH) 是建立在应用层和传输层基础上的安全协议，由 IETF 的网络工作小组（Network Working Group）所制定。SSH 是目前较可靠，专为远程登录会话和其他网络服务提供安全性的协议。利用 SSH 协议可以有效防止远程管理过程中的信息泄露问题。SSH 客户端与服务器端的通信的安全验证除了用户名密码的口令验证方式以外，还可以使用密钥的验证方式。

在程序开发中，我们经常需要用程序与服务器或存储设备建立 SSH 连接来管理设备（例如重启设备，查看硬件信息等等），目前 Java 的开源 SSH 库一般只支持 OpenSSH 格式的密钥连接，但日常使用中大量客户会使用 PuTTY 来创建 SSH 密钥，这样程序就需要对密钥进行解析和转换。本文描述了在 Java 程序中如何使用 Orion SSH2 通过密钥 SSH 连接到服务器端并执行操作，同时分析了常用的两种 SSH 密钥的格式以及解析的原理。

## 使用 Java 建立 SSH 连接并执行 cli

在本文中我们将使用 Orion SSH2 作为 SSH 客户端，然后使用密钥去访问 linux 服务器并执行 linux 命令。

首先，我们需要在 linux 服务器上生成一对密钥，将公钥注册到服务器上，然后将私钥下载到本地，接着使用 PuTTY 来载入私钥并连接 Linux 服务器。很多设备如服务器的管理工具，存储设备和交换机也提供 SSH 连接功能，有时候我们需要根据用户需求开发程序来对这些设备进行管理，这时可以使用 SSH client 库来进行开发。

Orion SSH2 是一个纯 Java 实现的 SSH-2 协议包，可让 Java 程序透过 SSH 协议连接到服务器上执行远程命令和文件传输功能。通过 Orion SSH2，我们只需要传入 SSH 服务端的 ip、port 还有密钥信息，就可以建立 SSH 连接，执行 cli 命令的具体 code 如下所示：

```
        String hostname = "192.168.1.2";
        String username = "root";
         //输入密钥所在路径
        File keyfile = new File("C:\\temp\\private");
         //输入密钥的加密密码，没有可以设为 null
        String keyfilePass = "joespass";
        try
        {
            /* 创建一个 SSH 连接 */
            Connection conn = new Connection(hostname);
            /* 尝试连接 */
            conn.connect();
            /* 传入 SSH key 进行验证 */
            boolean isAuthenticated = conn.authenticateWithPublicKey(username,
                                    keyfile,keyfilePass);
            if (isAuthenticated == false)
            throw new IOException("Authentication failed.");
            /* 验证通过，开始 SSH 会话 */
            Session sess = conn.openSession();
            //执行 linux 命令
             sess.execCommand("uname -a && date && uptime && who");
             //获取命令行输出
            InputStream stdout = new StreamGobbler(sess.getStdout());
            BufferedReader br = new BufferedReader(new InputStreamReader(stdout));
            System.out.println("Here is some information about the remote host:");
            while (true)
            {
                String line = br.readLine();
                if (line == null)
                    break;
                System.out.println(line);
            }
            /* 关闭 SSH 会话 */
            sess.close();
            /* 关闭 SSH 连接 */
            conn.close();
        }
        catch (IOException e)
        {
            e.printStackTrace(System.err);
            System.exit(2);
        }

```

Show moreShow more icon

以上代码展示了连接 Linux 系统并执行命令 uname -a && date && uptime && who 的过程， 这个命令用来打印操作系统信息，当前系统时间，系统运行时间以及当前登录用户的信息。同样我们可以利用这个代码去连接服务器的 IMM（Integrated Management Module）来执行命令查看服务器的网络配置和硬件信息或对服务器进行重启操作，也可以连接到 IBM Storwize V7000 的控制台对 V7000 进行管理。

## OpenSSH 密钥格式简介以及 Java 解析

### OpenSSH 密钥格式

OpenSSH 在文档 RFC4716 中定义了公钥（Public Key）和私钥 (Private Key) 的格式，简单来说一个密钥由开始标识（Begin Maker）、文件头信息（Header）、文件体（Body）和结束标志（End Maker）组成。下面是一个没有加密的私钥（本文将使用未加密的 RSA 私钥来进行解析）。

```
-----BEGIN RSA PRIVATE KEY-----
MIICWgIBAAKBgGAbzjs7fo9HeoxvZqnjHu9WTOFbNX+RMGtRtrA0PqG0ekOri7Mn
SkD3K8Zgs3tGjVEwW48YdUfVqNDwGpdwZtlhBgPydbqk1Ki87arAZrExAEly7TQ6
0+CPj8vGkVEC3K1ZH0/TkRUHlP/6C8V3SOqI56d1/ZjqWerZ8FVFqMy9AgElAoGA
IcSUkVoXtc0Bi0m8SYcmi3FZSEKkGBBq9UY5RNQWAXbDLIhhhCKPtfX6n6VvflcP
DrAgK1ue1AzMnHAJV82L6xM4S1M06jmq2UvVmpNjpeCN978U9Qpjke+iQDYS1gS6
cK3xQYFKM1zGrPCkbqFl4FsWH62pWuH5amxLzu71R60CQQC5EbbwYgTn7BnOQGTV
/gDZDwQWiMbfRpBQFGZiWlwGVLb8kCJ7TrkDUzG6yt8ZCSzdkgsJXkGkqS51f+uk
/b0BAkEAhPGeAZDMU9rpZZwX0MM6sXDmsrKgPfvElXKHQHeAKiaORhuaSr3xN0dr
4R4hppPwWVG4aqmjYBI+ug7N5N1DvQJBAKAPUhwBvw3FRsA3sSfG6/n/JiFTsuqd
5JhI/ppAT5bFzrDrXBeeB8uGONjmzsmLZRKn0jGdoI5ozjwbm15DO60CQCPuRmFJ
uq7hOClMyCqVoSkJwc9ujCxtjxOiab5lfJXFOzWK624lf3a5W21GafWrcWQ/maBJ
hhn3F99CRXwgICUCQQCcDZ8y99aUOPppd1t4Z1Wemc44Nyl/M+IxGxiEvEAUWv0U
dBBTyuQX+9Gh2WImVYfl8ugGTqCI22n6Xt/u+w8n
-----END RSA PRIVATE KEY-----

```

Show moreShow more icon

OpenSSH 支持 RSA 和 DSA 类型的密钥，本文将以 RSA 密钥为例来进行解析。在文档 RFC3447 中定义了 RSA 密钥的语法结构，私钥的语法结构如下所示：

```
RSAPrivateKey ::= SEQUENCE {
version Version,
modulus INTEGER, -- n
publicExponent INTEGER, -- e
privateExponent INTEGER, -- d
prime1 INTEGER, -- p
prime2 INTEGER, -- q
exponent1 INTEGER, -- d mod (p-1)
exponent2 INTEGER, -- d mod (q-1)
coefficient INTEGER, -- (inverse of q) mod p
otherPrimeInfos OtherPrimeInfos OPTIONAL
}

```

Show moreShow more icon

version 是 RSA 的版本号，本文中使用的密钥版本号是 0，如果密钥是使用多素数（多于 2 个素数）版本号则为 1。

modulus 是 RSA 的合数模 n。

publicExponent 是 RSA 的公开幂 e。

privateExponent 是 RSA 的私有幂 d。

prime1 是 n 的素数因子 p。

prime2 i 是 n 的素数因子 q。

exponent1 等于 d mod (p − 1)。

exponent2 等于 d mod (q − 1)。

coefficient 是 CRT 系数 q–1 mod p。

otherPrimeInfos 按顺序包含了其它素数 r3，……，ru 的信息。如果 version 是 0 ，它应该被忽略；而如果 version 是 1，它应该至少包含 OtherPrimeInfo 的一个实例。

RSA 公钥的语法结构如下所示，可见公钥所需的因子信息都包含在私钥中。

```
RSAPublicKey ::= SEQUENCE {
modulus INTEGER, -- n
publicExponent INTEGER -- e
}

```

Show moreShow more icon

### OpenSSH 密钥解析

OpenSSH 的 RSA 密钥的文件体使用 DER 编码，JDK 提供了 DerInputStream 来解析 DER 编码的字符串。所以 OpenSSH 密钥的解析很简单，首先读取密钥，过滤掉开始和结束标志，文件头（如果是加密的密钥则需要根据文件头信息来确定解密方式，因本文使用未加密的密钥故去掉文件头），然后使用 DER inputstream 来解析密钥，代码如下：

```
String keyinfo = "";
String line = null;
//去掉文件头尾的注释信息
while ((line = br.readLine()) != null) {
if (line.indexOf("---") == -1) {
keyinfo += line;
}
}
//密钥信息用 BASE64 编码加密过，需要先解密
byte[] decodeKeyinfo = (new BASE64Decoder()).decodeBuffer(keyinfo);
//使用 DerInputStream 读取密钥信息
DerInputStream dis = new DerInputStream(decodeKeyinfo);
//密钥不含 otherPrimeInfos 信息，故只有 9 段
DerValue[] ders = dis.getSequence(9);
//依次读取 RSA 因子信息
int version = ders[0].getBigInteger().intValue();
BigInteger modulus = ders[1].getBigInteger();
BigInteger publicExponent = ders[2].getBigInteger();
BigInteger privateExponent = ders[3].getBigInteger();
BigInteger primeP = ders[4].getBigInteger();
BigInteger primeQ = ders[5].getBigInteger();
BigInteger primeExponentP = ders[6].getBigInteger();
BigInteger primeExponentQ = ders[7].getBigInteger();
BigInteger crtCoefficient = ders[8].getBigInteger();
//generate public key and private key
KeyFactory keyFactory = KeyFactory.getInstance("RSA");
RSAPublicKeySpec rsaPublicKeySpec =
new RSAPublicKeySpec(modulus, publicExponent);
PublicKey publicKey = keyFactory.generatePublic(rsaPublicKeySpec);
RSAPrivateCrtKeySpec rsaPrivateKeySpec =
new RSAPrivateCrtKeySpec(modulus,publicExponent,privateExponent,
primeP,primeQ,primeExponentP,primeExponentQ,crtCoefficient);
PrivateKey privateKey = keyFactory.generatePrivate(rsaPrivateKeySpec);

```

Show moreShow more icon

在 Orion SSH2 的源代码中亦可见 SSH key 的解析代码，详见 PemDecoder 类, 以下为解析 RSA private key 的片段：

```
if (ps.pemType == PEM_RSA_PRIVATE_KEY){
SimpleDERReader dr = new SimpleDERReader(ps.data);
byte[] seq = dr.readSequenceAsByteArray();
if (dr.available() != 0)
throw new IOException("Padding in RSA PRIVATE KEY DER stream.");
dr.resetInput(seq);
BigInteger version = dr.readInt();
if ((version.compareTo(BigInteger.ZERO) != 0)
            && (version.compareTo(BigInteger.ONE) != 0))
throw new IOException("Wrong version (" + version + ")
        in RSA PRIVATE KEY DER stream.");
BigInteger n = dr.readInt();
System.out.println("n-----"+n);
BigInteger e = dr.readInt();
System.out.println("e-----"+e);
BigInteger d = dr.readInt();
System.out.println("d-----"+d);
return new RSAPrivateKey(d, e, n);
    }

```

Show moreShow more icon

在解析好密钥以后就可以使用这些密钥因子信息与 SSH 服务器进行验证，具体验证过程可参看 Orion SSH2 里的 AuthenticationManager 类。

## PuTTY Private Key（PPK）密钥格式简介以及 Java 解析

### PPK 密钥格式

PuTTY 没有提供文档说明其密钥的语法，但 PuTTY 是一个开源项目，我们可以从其官网下载 PuTTY 的源代码从而获知其构建密钥的方式。

从 PuTTY 源码可得知 PuTTY 的密钥结构如下：

```
PuTTY-User-Key-File-2:<key algorithm>
Encryption:<encrypt type>
Comment: <comments>
Public-Lines:<line number>
<Publick Key Body>
Private-Lines:<line number>
<Private Key Body>
Private-MAC:<hmac>

```

Show moreShow more icon

在首行的 key algorithm 标记了 key 的算法，“ssh-rsa” 表示采用 RSA 算法。下面一行的 encryption 表示 key 的加密方式，目前只支持 “aes256-cbc”和 “none”。Public lines 和 private lines 后面的数字分别代表 public key 和 private key 的行数。Private MAC 为这个密钥信息的 HMAC-SHA1 值。

一个实际的 ppk 密钥如下：

```
PuTTY-User-Key-File-2: ssh-rsa
Encryption: none
Comment: imported-openssh-key
Public-Lines: 4
AAAAB3NzaC1yc2EAAAABJQAAAIBgG847O36PR3qMb2ap4x7vVkzhWzV/kTBrUbaw
ND6htHpDq4uzJ0pA9yvGYLN7Ro1RMFuPGHVH1ajQ8BqXcGbZYQYD8nW6pNSovO2q
wGaxMQBJcu00OtPgj4/LxpFRAtytWR9P05EVB5T/+gvFd0jqiOendf2Y6lnq2fBV
RajMvQ==
Private-Lines: 8
AAAAgCHElJFaF7XNAYtJvEmHJotxWUhCpBgQavVGOUTUFgF2wyyIYYQij7X1+p+l
b35XDw6wICtbntQMzJxwCVfNi+sTOEtTNOo5qtlL1ZqTY6Xgjfe/FPUKY5HvokA2
EtYEunCt8UGBSjNcxqzwpG6hZeBbFh+tqVrh+WpsS87u9UetAAAAQQC5EbbwYgTn
7BnOQGTV/gDZDwQWiMbfRpBQFGZiWlwGVLb8kCJ7TrkDUzG6yt8ZCSzdkgsJXkGk
qS51f+uk/b0BAAAAQQCE8Z4BkMxT2ullnBfQwzqxcOaysqA9+8SVcodAd4AqJo5G
G5pKvfE3R2vhHiGmk/BZUbhqqaNgEj66Ds3k3UO9AAAAQQCcDZ8y99aUOPppd1t4
Z1Wemc44Nyl/M+IxGxiEvEAUWv0UdBBTyuQX+9Gh2WImVYfl8ugGTqCI22n6Xt/u
+w8n
Private-MAC: 4c039beb12ae005acb763225572cc4eeec767542

```

Show moreShow more icon

### PPK 密钥解析

由 PuTTY 的源码可知 PPK 中公钥包含以下信息：

```
string "ssh-rsa"
mpint exponent
mpint modulus

```

Show moreShow more icon

私钥包含以下信息：

```
mpint private_exponent
mpint p (the larger of the two primes)
mpint q (the smaller prime)
mpint iqmp (the inverse of q modulo p)
data padding (to reach a multiple of the cipher block size)

```

Show moreShow more icon

这些因子信息可以参看上文 OpenSSH 私钥格式的解释。

公钥和私钥的信息均遵从 RFC 4251 中定义的数据格式，简单来说即每一个元素以字节流的形式进行存储，前 4 个字节为整数，数值为该元素的字节长度。由于私钥仅储存了 iqpm（即 OpenSSH 所说的 coefficient），所以我们还需要计算出 primeP（即上文 prime1）和 primeQ（即上文 prime2）。具体解析代码如下：

```
public static void testPPK() {
        String file = "private1.ppk";
        try {
            Map<String, String> keyMap = parseKV(file);
             //获取公钥信息
            String publicKeyInfo = keyMap.get("Public-Lines");
             //密钥信息被 BASE64 加密过，需要先解密
            byte[] decodedPubKey = (new BASE64Decoder())
                    .decodeBuffer(publicKeyInfo);
            System.out.println("publick key : "
                    + bytesToHexString(decodedPubKey));
            DataInputStream dis = new DataInputStream(new ByteArrayInputStream(
                    decodedPubKey));
             //读取前 4 个字节，获得该元素数据长度
            int leng = dis.readInt();
             //根据长度读入字节信息
             //公钥第一段为固定字符串信息，所以单独解析
            byte[] tmpBytes = new byte[leng];
            dis.readFully(tmpBytes);
            String keyAlgo = new String(tmpBytes);
            BigInteger publicExponent = readInt(dis);
            BigInteger modulus = readInt(dis);
            dis.close();
             //用 BASE64 编码解密私钥
            byte[] decodedPriKey =
            (new BASE64Decoder()).decodeBuffer(keyMap.get("Private-Lines"));
            dis = new DataInputStream(new ByteArrayInputStream(decodedPriKey));
            BigInteger privateExponent = readInt(dis);
            BigInteger primeP = readInt(dis);
            BigInteger primeQ = readInt(dis);
            BigInteger iqmp = readInt(dis);
            BigInteger primeExponentP = privateExponent.mod(primeP.subtract(BigInteger.ONE));
            BigInteger primeExponentQ = privateExponent.mod(primeQ.subtract(BigInteger.ONE));
            KeyFactory keyFactory = KeyFactory.getInstance("RSA");
            RSAPublicKeySpec rsaPublicKeySpec = new RSAPublicKeySpec(modulus,
                    publicExponent);
            PublicKey publicKey = keyFactory.generatePublic(rsaPublicKeySpec);
            RSAPrivateCrtKeySpec rsaPrivateKeySpec = new RSAPrivateCrtKeySpec(
            modulus, publicExponent, privateExponent, primeP, primeQ,
            primeExponentP, primeExponentQ, iqmp);
            PrivateKey privateKey = keyFactory
                    .generatePrivate(rsaPrivateKeySpec);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
     //读入前 4 个字节获得元素长度，然后读取该元素字节信息并转换为 BigInteger
    public static BigInteger readInt(DataInputStream dis) throws IOException {
        int leng = dis.readInt();
        byte[] tmpBytes = new byte[leng];
        dis.readFully(tmpBytes);
        return new BigInteger(1, tmpBytes);
    }
     //将密钥信息解析为键值对
    private static Map<String, String> parseKV(String file) throws IOException {
        HashMap<String, String> kv = new HashMap<String, String>();
        BufferedReader r = null;
        try {
            r = new BufferedReader(new FileReader(file));
            String k = null;
            String line;
            while ((line = r.readLine()) != null) {
                int idx = line.indexOf(": ");
            if (idx > 0) {
                    k = line.substring(0, idx);
                    if ((!"Public-Lines".equals(k))
                            && (!"Private-Lines".equals(k))) {
                        kv.put(k, line.substring(idx + 2));
                    }
                } else {
                    String s = (String) kv.get(k);
                    if (s == null) {
                        s = "";
                    }
                    s = s + line;
                    kv.put(k, s);
                }
            }
        } finally {
            r.close();
        }
        return kv;
    }

```

Show moreShow more icon

## 两种格式的 key 的转换

Orion SSH2 不支持 PPK 格式的 SSH key，但通过以上方法我们可以自行从 PPK 中读取到 RSA 的各个元素，这样在 Orion SSH2 的 AuthenticationManager 的 authenticatePublicKey 中加入对 PPK 格式密钥的判断和处理，就可以让 Orion SSH2 使用 PPK 去连接 SSH 服务器。

大部分的 Java SSH client 库都只提供 OpenSSH key 的支持，但实际应用中大部分用户都是使用 PuTTY 来生成管理 SSH 密钥的，使用上文的方法就可以让我们编写的程序直接支持 PPK 格式的密钥而不用要求用户在使用程序前先对密钥进行格式的转换。

同样我们可以在程序中提供密钥格式转换的功能方便用户使用。

将这两种密钥进行转换，首先都需要将其解析成对应密钥算法的因子，解析方法在上两节已有描述，故不再赘述，下面仅说明如何根据 RSA 的因子生成 OpenSSH 格式或 ppk 格式密钥。

### 转换成 OpenSSH 格式

OpenSSH 密钥采用 DER 编码方式记录信息，Java 提供了 DerOutputStream 来写入 Der 数据，唯一需要注意的是构建完 key 的字节信息之后还需要写入 DER 的 Sequence tag 信息来标明该段信息是按顺序排列的。转换的代码如下：

```
DerOutputStream deros = new DerOutputStream();
deros.putInteger(version);
deros.putInteger(modulus);
deros.putInteger(publicExponent);
deros.putInteger(privateExponent);
deros.putInteger(primeP);
deros.putInteger(primeQ);
deros.putInteger(primeExponentP);
deros.putInteger(primeExponentQ);
deros.putInteger(crtCoefficient);
deros.flush();
byte[] rsaInfo = deros.toByteArray();
deros.close();
deros = new DerOutputStream();
//写入信息时需要写一个 Sequence 的 tag
deros.write(DerValue.tag_Sequence, rsaInfo);
deros.flush();
//写入时需要用 BASE64 编码进行加密
String keyInfo = (new BASE64Encoder()).encode(deros.toByteArray());

```

Show moreShow more icon

### 转换成 PPK 格式

在 PPK 密钥格式解析中已经分析了 ppk 的信息格式，在获取到 RSA 的计算因子后写入信息时对每个因子先计算其字节长度并写入长度信息，然后再写入该因子的字节信息即可，以 public key 信息为例：

```
ByteArrayOutputStream bos = new ByteArrayOutputStream();
//写入该元素长度信息
bos.write(convertIntToByteArray(keyAlgo.length()));
//写入该元素的字节信息
bos.write(keyAlgo.getBytes());
bos.write(convertIntToByteArray(publicExponent.toByteArray().length));
bos.write(publicExponent.toByteArray());
bos.write(convertIntToByteArray(modulus.toByteArray().length));
bos.write(modulus.toByteArray());
bos.flush();
String keyInfo = (new BASE64Encoder()).encode(bos.toByteArray());
bos.close();
public static byte[] convertIntToByteArray(int value){
byte[] bytevalue = new byte[4];
bytevalue[0] = (byte)((value >> 24)&0xFF);
bytevalue[1] = (byte)((value >> 16)&0xFF);
bytevalue[2] = (byte)((value >> 8)&0xFF);
bytevalue[3] = (byte)(value&0xFF);
return bytevalue;
}

```

Show moreShow more icon

## 结束语

SSH 是一个使用非常广泛的协议，很多服务器的管理工具（如 IMM，CMM）、存储设备（如 IBM Storwize V7000）和交换机都提供了 SSH 服务端供管理员使用 SSH 客户端（如 PuTTY）登录并使用命令行来查看设备信息、对设备进行远程操控，因此很多管理程序也使用了 SSH 客户端来对这些设备进行管理。

本文以开源的 Java SSH client（Orion SSH2）为例，解释了 SSH 密钥的结构以及 SSH client 对密钥的解析过程，希望能对大家理解 SSH 密钥、使用 SSH client 进行开发有所帮助。