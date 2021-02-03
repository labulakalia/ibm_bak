# 运用加密技术保护 Java 源代码
如何在不修改原有程序的情况下通过加密技术保护源代码

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/l-secureclass/)

俞良松

发布: 2001-10-15

* * *

## 为什么要加密？

对于传统的C或C++之类的语言来说，要在Web上保护源代码是很容易的，只要不发布它就可以。遗憾的是，Java程序的源代码很容易被别人偷看。只要有一个反编译器，任何人都可以分析别人的代码。Java的灵活性使得源代码很容易被窃取，但与此同时，它也使通过加密保护代码变得相对容易，我们唯一需要了解的就是Java的ClassLoader对象。当然，在加密过程中，有关Java Cryptography Extension（JCE）的知识也是必不可少的。

有几种技术可以”模糊”Java类文件，使得反编译器处理类文件的效果大打折扣。然而，修改反编译器使之能够处理这些经过模糊处理的类文件并不是什么难事，所以不能简单地依赖模糊技术来保证源代码的安全。

我们可以用流行的加密工具加密应用，比如PGP（Pretty Good Privacy）或GPG（GNU Privacy Guard）。这时，最终用户在运行应用之前必须先进行解密。但解密之后，最终用户就有了一份不加密的类文件，这和事先不进行加密没有什么差别。

Java运行时装入字节码的机制隐含地意味着可以对字节码进行修改。JVM每次装入类文件时都需要一个称为ClassLoader的对象，这个对象负责把新的类装入正在运行的JVM。JVM给ClassLoader一个包含了待装入类（比如java.lang.Object）名字的字符串，然后由ClassLoader负责找到类文件，装入原始数据，并把它转换成一个Class对象。

我们可以通过定制ClassLoader，在类文件执行之前修改它。这种技术的应用非常广泛――在这里，它的用途是在类文件装入之时进行解密，因此可以看成是一种即时解密器。由于解密后的字节码文件永远不会保存到文件系统，所以窃密者很难得到解密后的代码。

由于把原始字节码转换成Class对象的过程完全由系统负责，所以创建定制ClassLoader对象其实并不困难，只需先获得原始数据，接着就可以进行包含解密在内的任何转换。

Java 2 在一定程度上简化了定制ClassLoader的构建。在Java 2中，loadClass的缺省实现仍旧负责处理所有必需的步骤，但为了顾及各种定制的类装入过程，它还调用一个新的findClass方法。

这为我们编写定制的ClassLoader提供了一条捷径，减少了麻烦：只需覆盖findClass，而不是覆盖loadClass。这种方法避免了重复所有装入器必需执行的公共步骤，因为这一切由loadClass负责。

不过，本文的定制ClassLoader并不使用这种方法。原因很简单。如果由默认的ClassLoader先寻找经过加密的类文件，它可以找到；但由于类文件已经加密，所以它不会认可这个类文件，装入过程将失败。因此，我们必须自己实现loadClass，稍微增加了一些工作量。

## 定制类装入器

每一个运行着的JVM已经拥有一个ClassLoader。这个默认的ClassLoader根据CLASSPATH环境变量的值，在本地文件系统中寻找合适的字节码文件。

应用定制ClassLoader要求对这个过程有较为深入的认识。我们首先必须创建一个定制ClassLoader类的实例，然后显式地要求它装入另外一个类。这就强制JVM把该类以及所有它所需要的类关联到定制的ClassLoader。清单 1.显示了如何用定制C lassLoader装入类文件。

清单 1\. 利用定制的 ClassLoader 装入类文件

```
// 首先创建一个ClassLoader对象
ClassLoader myClassLoader = new myClassLoader();
// 利用定制ClassLoader对象装入类文件
// 并把它转换成Class对象
Class myClass = myClassLoader.loadClass( "mypackage.MyClass" );
// 最后，创建该类的一个实例
Object newInstance = myClass.newInstance();
// 注意，MyClass所需要的所有其他类，都将通过
// 定制的ClassLoader自动装入

```

Show moreShow more icon

如前所述，定制ClassLoader只需先获取类文件的数据，然后把字节码传递给运行时系统，由后者完成余下的任务。

ClassLoader有几个重要的方法。创建定制的ClassLoader时，我们只需覆盖其中的一个，即loadClass，提供获取原始类文件数据的代码。这个方法有两个参数：类的名字，以及一个表示JVM是否要求解析类名字的标记（即是否同时装入有依赖关系的类）。如果这个标记是true，我们只需在返回JVM之前调用resolveClass。

清单 2\. ClassLoader.loadClass() 的一个简单实现

```
      public Class loadClass( String name, boolean resolve )
      throws ClassNotFoundException {
    try {
      // 我们要创建的Class对象
       Class clasz = null;
      // 必需的步骤1：如果类已经在系统缓冲之中，
      // 我们不必再次装入它
      clasz = findLoadedClass( name );
      if (clasz != null)
        return clasz;
      // 下面是定制部分
      byte classData[] = /* 通过某种方法获取字节码数据 */;
      if (classData != null) {
        // 成功读取字节码数据，现在把它转换成一个Class对象
        clasz = defineClass( name, classData, 0, classData.length );
      }
      // 必需的步骤2：如果上面没有成功，
      // 我们尝试用默认的ClassLoader装入它
      if (clasz == null)
        clasz = findSystemClass( name );
      // 必需的步骤3：如有必要，则装入相关的类
      if (resolve && clasz != null)
        resolveClass( clasz );
      // 把类返回给调用者
      return clasz;
    } catch( IOException ie ) {
      throw new ClassNotFoundException( ie.toString() );
    } catch( GeneralSecurityException gse ) {
      throw new ClassNotFoundException( gse.toString() );
    }
}

```

Show moreShow more icon

清单 2 显示了一个简单的loadClass实现。代码中的大部分对所有ClassLoader对象来说都一样，但有一小部分（已通过注释标记）是特有的。在处理过程中，ClassLoader对象要用到其他几个辅助方法：

- findLoadedClass：用来进行检查，以便确认被请求的类当前还不存在。loadClass方法应该首先调用它。
- defineClass：获得原始类文件字节码数据之后，调用defineClass把它转换成一个Class对象。任何loadClass实现都必须调用这个方法。
- findSystemClass：提供默认ClassLoader的支持。如果用来寻找类的定制方法不能找到指定的类（或者有意地不用定制方法），则可以调用该方法尝试默认的装入方式。这是很有用的，特别是从普通的JAR文件装入标准Java类时。
- resolveClass：当JVM想要装入的不仅包括指定的类，而且还包括该类引用的所有其他类时，它会把loadClass的resolve参数设置成true。这时，我们必须在返回刚刚装入的Class对象给调用者之前调用resolveClass。

## 加密、解密

Java加密扩展即Java Cryptography Extension，简称JCE。它是Sun的加密服务软件，包含了加密和密匙生成功能。JCE是JCA（Java Cryptography Architecture）的一种扩展。

JCE没有规定具体的加密算法，但提供了一个框架，加密算法的具体实现可以作为服务提供者加入。除了JCE框架之外，JCE软件包还包含了SunJCE服务提供者，其中包括许多有用的加密算法，比如DES（Data Encryption Standard）和Blowfish。

为简单计，在本文中我们将用DES算法加密和解密字节码。下面是用JCE加密和解密数据必须遵循的基本步骤：

1. 生成一个安全密匙。在加密或解密任何数据之前需要有一个密匙。密匙是随同被加密的应用一起发布的一小段数据，清单 3 显示了如何生成一个密匙。

    清单 3\. 生成一个密匙





    ```
    // DES算法要求有一个可信任的随机数源
    SecureRandom sr = new SecureRandom();
    // 为我们选择的DES算法生成一个KeyGenerator对象
    KeyGenerator kg = KeyGenerator.getInstance( "DES" );
    kg.init( sr );
    // 生成密匙
    SecretKey key = kg.generateKey();
    // 获取密匙数据
    byte rawKeyData[] = key.getEncoded();
    /* 接下来就可以用密匙进行加密或解密，或者把它保存
        为文件供以后使用 */
    doSomething( rawKeyData );

    ```





    Show moreShow more icon

2. 加密数据。得到密匙之后，接下来就可以用它加密数据。除了解密的ClassLoader之外，一般还要有一个加密待发布应用的独立程序（见清单 4）。

    清单 4\. 用密匙加密原始数据





    ```
    // DES算法要求有一个可信任的随机数源
    SecureRandom sr = new SecureRandom();
    byte rawKeyData[] = /* 用某种方法获得密匙数据 */;
    // 从原始密匙数据创建DESKeySpec对象
    DESKeySpec dks = new DESKeySpec( rawKeyData );
    // 创建一个密匙工厂，然后用它把DESKeySpec转换成
    // 一个SecretKey对象
    SecretKeyFactory keyFactory = SecretKeyFactory.getInstance( "DES" );
    SecretKey key = keyFactory.generateSecret( dks );
    // Cipher对象实际完成加密操作
    Cipher cipher = Cipher.getInstance( "DES" );
    // 用密匙初始化Cipher对象
    cipher.init( Cipher.ENCRYPT_MODE, key, sr );
    // 现在，获取数据并加密
    byte data[] = /* 用某种方法获取数据 */
    // 正式执行加密操作
    byte encryptedData[] = cipher.doFinal( data );
    // 进一步处理加密后的数据
    doSomething( encryptedData );

    ```





    Show moreShow more icon

3. 解密数据。运行经过加密的应用时，ClassLoader分析并解密类文件。操作步骤如清单 5所示。

    清单 5\. 用密匙解密数据





    ```
    // DES算法要求有一个可信任的随机数源
    SecureRandom sr = new SecureRandom();
    byte rawKeyData[] = /* 用某种方法获取原始密匙数据 */;
    // 从原始密匙数据创建一个DESKeySpec对象
    DESKeySpec dks = new DESKeySpec( rawKeyData );
    // 创建一个密匙工厂，然后用它把DESKeySpec对象转换成
    // 一个SecretKey对象
    SecretKeyFactory keyFactory = SecretKeyFactory.getInstance( "DES" );
    SecretKey key = keyFactory.generateSecret( dks );
    // Cipher对象实际完成解密操作
    Cipher cipher = Cipher.getInstance( "DES" );
    // 用密匙初始化Cipher对象
    cipher.init( Cipher.DECRYPT_MODE, key, sr );
    // 现在，获取数据并解密
    byte encryptedData[] = /* 获得经过加密的数据 */
    // 正式执行解密操作
    byte decryptedData[] = cipher.doFinal( encryptedData );
    // 进一步处理解密后的数据
    doSomething( decryptedData );

    ```





    Show moreShow more icon


## 应用实例

前面介绍了如何加密和解密数据。要部署一个经过加密的应用，步骤如下：

1. 创建应用。我们的例子包含一个App主类，两个辅助类（分别称为Foo和Bar）。这个应用没有什么实际功用，但只要我们能够加密这个应用，加密其他应用也就不在话下。
2. 生成一个安全密匙。在命令行，利用GenerateKey工具（参见GenerateKey.java）把密匙写入一个文件：





    ```
    % java GenerateKey key.data

    ```





    Show moreShow more icon

3. 加密应用。在命令行，利用EncryptClasses工具（参见EncryptClasses.java）加密应用的类：





    ```
    % java EncryptClasses key.data App.class Foo.class Bar.class

    ```





    Show moreShow more icon

    该命令把每一个.class文件替换成它们各自的加密版本。

4. 运行经过加密的应用。用户通过一个DecryptStart程序运行经过加密的应用。DecryptStart程序如 清单 6 所示。

清单 6\. DecryptStart.java，启动被加密应用的程序

```
import java.io.*;
import java.security.*;
import java.lang.reflect.*;
import javax.crypto.*;
import javax.crypto.spec.*;
public class DecryptStart extends ClassLoader
{
     // 这些对象在构造函数中设置，
     // 以后loadClass()方法将利用它们解密类
     private SecretKey key;
     private Cipher cipher;
     // 构造函数：设置解密所需要的对象
     public DecryptStart( SecretKey key ) throws GeneralSecurityException,
         IOException {
       this.key = key;
       String algorithm = "DES";
       SecureRandom sr = new SecureRandom();
       System.err.println( "[DecryptStart: creating cipher]" );
       cipher = Cipher.getInstance( algorithm );
       cipher.init( Cipher.DECRYPT_MODE, key, sr );
     }
     // main过程：我们要在这里读入密匙，创建DecryptStart的
     // 实例，它就是我们的定制ClassLoader。
     // 设置好ClassLoader以后，我们用它装入应用实例，
     // 最后，我们通过Java Reflection API调用应用实例的main方法
     static public void main( String args[] ) throws Exception {
       String keyFilename = args[0];
       String appName = args[1];
        // 这些是传递给应用本身的参数
       String realArgs[] = new String[args.length-2];
       System.arraycopy( args, 2, realArgs, 0, args.length-2 );
       // 读取密匙
       System.err.println( "[DecryptStart: reading key]" );
       byte rawKey[] = Util.readFile( keyFilename );
       DESKeySpec dks = new DESKeySpec( rawKey );
       SecretKeyFactory keyFactory = SecretKeyFactory.getInstance( "DES" );
       SecretKey key = keyFactory.generateSecret( dks );
       // 创建解密的ClassLoader
       DecryptStart dr = new DecryptStart( key );
       // 创建应用主类的一个实例
       // 通过ClassLoader装入它
       System.err.println( "[DecryptStart: loading "+appName+"]" );
       Class clasz = dr.loadClass( appName );
       // 最后，通过Reflection API调用应用实例
       // 的main()方法
       // 获取一个对main()的引用
       String proto[] = new String[1];
       Class mainArgs[] = { (new String[1]).getClass() };
       Method main = clasz.getMethod( "main", mainArgs );
       // 创建一个包含main()方法参数的数组
       Object argsArray[] = { realArgs };
       System.err.println( "[DecryptStart: running "+appName+".main()]" );
       // 调用main()
       main.invoke( null, argsArray );
     }
     public Class loadClass( String name, boolean resolve )
         throws ClassNotFoundException {
       try {
         // 我们要创建的Class对象
         Class clasz = null;
         // 必需的步骤1：如果类已经在系统缓冲之中
         // 我们不必再次装入它
         clasz = findLoadedClass( name );
         if (clasz != null)
           return clasz;
         // 下面是定制部分
         try {
           // 读取经过加密的类文件
           byte classData[] = Util.readFile( name+".class" );
           if (classData != null) {
             // 解密...
             byte decryptedClassData[] = cipher.doFinal( classData );
             // ... 再把它转换成一个类
             clasz = defineClass( name, decryptedClassData,
               0, decryptedClassData.length );
             System.err.println( "[DecryptStart: decrypting class "+name+"]" );
           }
         } catch( FileNotFoundException fnfe ) {
         }
         // 必需的步骤2：如果上面没有成功
         // 我们尝试用默认的ClassLoader装入它
         if (clasz == null)
           clasz = findSystemClass( name );
         // 必需的步骤3：如有必要，则装入相关的类
         if (resolve && clasz != null)
           resolveClass( clasz );
         // 把类返回给调用者
         return clasz;
       } catch( IOException ie ) {
         throw new ClassNotFoundException( ie.toString()
);
       } catch( GeneralSecurityException gse ) {
         throw new ClassNotFoundException( gse.toString()
);
       }
     }
}

```

Show moreShow more icon

对于未经加密的应用，正常执行方式如下：

```
% java App arg0 arg1 arg2

```

Show moreShow more icon

对于经过加密的应用，则相应的运行方式为：

```
% java DecryptStart key.data App arg0 arg1 arg2

```

Show moreShow more icon

DecryptStart有两个目的。一个DecryptStart的实例就是一个实施即时解密操作的定制ClassLoader；同时，DecryptStart还包含一个main过程，它创建解密器实例并用它装入和运行应用。示例应用App的代码包含在App.java、Foo.java和Bar.java内。Util.java是一个文件I/O工具，本文示例多处用到了它。完整的代码请从本文最后下载。

## 注意事项

我们看到，要在不修改源代码的情况下加密一个Java应用是很容易的。不过，世上没有完全安全的系统。本文的加密方式提供了一定程度的源代码保护，但对某些攻击来说它是脆弱的。

虽然应用本身经过了加密，但启动程序DecryptStart没有加密。攻击者可以反编译启动程序并修改它，把解密后的类文件保存到磁盘。降低这种风险的办法之一是对启动程序进行高质量的模糊处理。或者，启动程序也可以采用直接编译成机器语言的代码，使得启动程序具有传统执行文件格式的安全性。

另外还要记住的是，大多数JVM本身并不安全。狡猾的黑客可能会修改JVM，从ClassLoader之外获取解密后的代码并保存到磁盘，从而绕过本文的加密技术。Java没有为此提供真正有效的补救措施。

不过应该指出的是，所有这些可能的攻击都有一个前提，这就是攻击者可以得到密匙。如果没有密匙，应用的安全性就完全取决于加密算法的安全性。虽然这种保护代码的方法称不上十全十美，但它仍不失为一种保护知识产权和敏感用户数据的有效方案。