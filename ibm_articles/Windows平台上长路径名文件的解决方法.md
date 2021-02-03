# Windows 平台上长路径名文件的解决方法
两种支持长路径名文件的 C/C++ 编程方法

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/java-j-lo-longpath/)

强晟, 韩兆兵

发布: 2008-01-31

* * *

## Windows 对长路径名文件的限制

众所周知，微软的文件系统经历了 fat->fat32->NTFS 的技术变革。且不论安全和文件组织方式上的革新，单就文件名而言，已经从古老的 DOS 8.3 文件格式（仅支持最长 8 个字符的文件名和 3 个字符的后缀名）转变为可以支持长达 255 个字符的文件名。而对于路径长度，NTFS 也已经支持长达 32768 个字符的路径名。

然而，Windows 操作系统并没有完全放开路径名长度的限制，在 windef.h 中，可以找到如下的宏：

```
#define MAX_PATH 260

```

Show moreShow more icon

事实上，所有的 Windows API 都遵循这个限制。因此，每当我们试图更改某一文件的文件名时，当输入的文件名长度 (全路径) 到达一定限度时，虽然文件名本身还未达到 255 个字符的限制，但是任何输入将不再被接受，这其实正是由于操作系统不允许 260 个字符（byte）的文件全路径。

实际应用中，这种 260 个字符的全路径的限制给应用开发带来了很大的不便。试想如下应用：我们希望给应用服务器增加一个本地 cache 的功能，该功能可以把远程服务器上的文件留下一个本地的副本。一个合理的实现可以把 url 映射为文件名，当 url 很长时，cache 文件的长度也会很长。当文件名长度超过 255，我们可以把映射文件名的前 255 个字符作为目录名称。但是，我们仍然无法解决 260 个字符的全路径限制。另外，如果一个应用软件的目录结构过深，很容易出现某些文件名长度（含路径）超过 260 个字符，并因此造成安装或删除的失败。总而言之，该限制给我们的开发测试工作带来了诸多不便。

对于一些网络服务器，往往需要将 Java 代码用于上层逻辑控制/事务处理的开发，同时将 C/C++ 用于底层核心功能的实现。为此，我们研究了这两种程序语言对长路径名文件的支持情况。其中，对于 Java，比较了两个常用版本 1.4 和 5.0 对长路径支持的差异性；对于 C/C++ 语言的局限性，提出了我们的解决方法。

实验环境 :

操作系统： Windows xp

文件系统： NTFS 文件系统

Java 编译环境： IBM JDK 1.4.2 以及 IBM JDK 5.0

C++ 编译环境： VC.net

## 在 Java 中使用长路径名文件

Java 语言并不需要对长路径名文件进行特殊的处理，就可以支持长路径名文件的创建、读写和删除操作等基本操作。但是，JDK 1.4.2 和 JDK 5.0 在长路径的支持上是不同的，JDK 1.4.2 并不是完全支持所有的长路径名文件操作，比如访问文件属性的操作是不支持的。我们设计了如下代码来验证 JDK 1.4.2 和 JDK 5.0 对长路径名文件支持的区别。

##### 清单 1\. 对长路径名文件操作的 Java 实验代码：

```
try {
    String fileName = "E:\\VerylongpathVerylongpathVerylongpath
        VerylongpathVerylongpathVerylongpathVerylongpath
        VerylongpathVerylongpathVerylongpathVerylongpath\\
    VerylongpathVerylongpathVerylongpathVery
        longpathVerylongpathVerylongpathVerylongpath
    VerylongpathVerylongpathVerylongpathVerylongpa
        th.txt";
    System.out.println("Filename: " + fileName);
    System.out.println("File path length: " + fileName.length());
    String renameFileName = "E:\\VerylongpathVerylongpathVerylongpath
        VerylongpathVerylongpathVerylongpathVerylongpath
        VerylongpathVerylongpathVerylongpathVerylongpath\\Short.txt";

    //Create the file.
    File file = new File(fileName);
    if (!file.exists())
        file.createNewFile();
    if (file.exists())
        System.out.println("The file exists!");
    if (file.canRead())
        System.out.println("The file can be read!");
    if (file.canWrite())
        System.out.println("The file can be written!");
    if (file.isFile())
        System.out.println("It's a file!");

    //Write to the created file.
    FileOutputStream out = new FileOutputStream(file);
    PrintStream p = new PrintStream(out);
    p.println("This is only a test!");
    p.close();

    //Read the information from that file.
    BufferedReader br = new BufferedReader(new FileReader(file));
    StringBuffer sb = new StringBuffer();
    while (true) {
        String sl = br.readLine();
        if (sl == null) {
            break;
        } else {
            sb.append(sl + "\n");
        }
    }
    br.close();
    System.out.println("The content in the file:");
    System.out.print("\t" + sb.toString());

    //File rename
    File newfile = new File(renameFileName);
    if (newfile.exists())
        System.out.println(renameFileName + "exsited");
    else {
        if (file.renameTo(newfile)){
            System.out.println("Rename sucessful!");
        } else {
            System.out.println("Rename failed!");
        }
    }

    //delete file
    if (file.delete())
        System.out.println("The old file deleted!");
    if (newfile.delete())
        System.out.println("The renamed file deleted!");
    }  catch (IOException e) {
        //Error happened
        e.printStackTrace();
        System.out.println("Error occurs in writing to the file.");
    }
}

```

Show moreShow more icon

##### 清单 2\. 使用 ibm-java2-sdk-142 的结果

```
Filename: E:\VerylongpathVerylongpathVerylongpath
VerylongpathVerylongpathVerylongpathVerylongpathVer
ylongpathVerylongpathVerylongpathVerylongpath\
VerylongpathVerylongpathVerylongpathVerylong
pathVerylongpathVerylongpathVerylongpath
VerylongpathVerylongpathVerylongpathVerylongpath.t
xt

File path length: 272

The content in the file:

This is only a test!

Rename failed!

The old file deleted!

```

Show moreShow more icon

从实验结果来看，JDK 1.4.2 得到了该长路径名文件的内容，因此，对于该长路径名文件的创建以及读写操作都是支持的。但是对比下文使用 JDK 5.0 的结果，可以看到，所有对于文件属性的判断都是错误的，同时，重命名的操作也无法实现。更为重要的是，JDK 1.4.2 存在着一个很致命的问题，即方法 `File.exists()` 是失效的。通常，在删除文件前，需要调用该方法判断文件是否存在，对于 JDK 1.4.2，如果直接去删除一个不知道是否存在的文件，就会存在比较大的风险。因此，JDK 1.4.2 在 Windows 平台对长路径名文件的操作只是有限的支持，使用的时候，一定要注意。

##### 清单 3\. 使用 ibm-java2-sdk-50 的结果

```
Filename: E:\VerylongpathVerylongpathVerylongpath
VerylongpathVerylongpathVerylongpathVerylongpathVer
ylongpathVerylongpathVerylongpathVerylongpath\
VerylongpathVerylongpathVerylongpathVerylong
pathVerylongpathVerylongpathVerylongpath
VerylongpathVerylongpathVerylongpathVerylongpath.t
xt
File path length: 272
The file exists!
The file can be read!
The file can be written!
It's a file!
The content in the file:
    This is only a test!
Rename sucessful!
The renamed file deleted!

```

Show moreShow more icon

从实验中可以清楚的看到，在版本 JDK 5.0 中，所有的文件操作（新建、读写、属性操作、重命名、删除等）都能够得到正确的处理。使用 JDK 5.0 就可以完全不用担心长路径名文件的使用问题。

## 在 C/C++ 中使用长路径名文件

相对于 JDK 5.0 不需要任何改动就可以支持长路径名文件，在 C/C++ 中使用超过 260 个字符的路径长度的文件，会复杂得多。下面介绍两种支持长路径名文件的方法。

### 方法一：使用 Unicode 版本的 API

从微软官方网站 [Path Field Limits](http://msdn2.microsoft.com/en-us/library/930f87yf(vs.80).aspx) ，可以查到，使用 Unicode 版本的 API，对于使用 NTFS 文件系统的 Windows NT 4.0, Windows 2000, Windows XP Home Edition, Windows XP Professional 和 Windows Server 2003 操作系统，可以支持 32768 字节的文件路径长度。同时，路径名必须使用 \\?\ 的前缀。依照这个思路，我们设计了实验。

##### 清单 4\. 对长路径名文件操作的 C 的示例代码（Unicode API）

```
{
FILE *from, *to;
char filename[1024];
strcpy(filename,"\\\\?\\E:\\VerylongpathVerylongpathVerylongpathVerylongpathVerylongpathV
erylongpathVerylongpathVerylongpathVerylongpathVerylongpathVerylongpath\\VerylongpathVeryl
ongpathVerylongpathVerylongpathVerylongpathVerylongpathVerylongpathVerylongpathVerylongpat
hVerylongpathVerylongpath.txt");
int iL1=MultiByteToWideChar(CP_ACP, 0, filename, strlen(filename), NULL, 0);
WCHAR* wfilename=new WCHAR[iL1+1];
wfilename[iL1] = '\0';
int iL2=MultiByteToWideChar(CP_ACP, 0, filename, strlen(filename), wfilename, iL1);
from = _wfopen( wfilename ,L"rb");
to = fopen(".\\longpath.txt", "wb");
if((from ==NULL)||(to==NULL))
    return -1;
char buffer[1024];
int count = 0;
while ( (count = fread(buffer, sizeof(char), 1024, from)) != 0)
    fwrite( buffer, sizeof(char), count, to);
delete []wfilename;
fclose (from); fclose(to);
}

```

Show moreShow more icon

使用如上的方法，我们可以拷贝某长路径名的文件到当前文件夹中。从试验结果看，该方法是有效的。但是，由于该方法要求系统使用 Unicode 的 API，同时需要更改路径名称以及编码方式。因此，对于一个已经存在的系统，由于需要改变所有文件操作相关的 API，因此改动将会很大。

### 方法二：创建 8.3 格式的短路径名

对于每一个长路径名，都有一个 8.3 格式（8 个字符的文件名和 3 个字符的后缀名）的短路径名与其相对应，任意的文件夹或者文件名都可以映射成一个 8 字符的文件名（A~B），其中 A 是文件名前缀，B 是表示字母序的顺序。操作系统可以保证这样的映射是一对一的，只要使用 `GetShortPathName()` 将长路径名转成相应的短路径名，就可以进行对该文件进行普通的文件操作。同时，在任何时候都可以用函数 `GetLongPathName()` 把 8.3 格式的短路径名恢复成初始的长路径名。

如 [GetShortPathName Function](http://msdn2.microsoft.com/en-us/library/aa364989.aspx) 叙述，我们需要一个 Unicode 版本的 API，同时在路径名前加上 \\?\ 的前缀，才能实现长短路径名间的切换。但从实验来看，即使不使用 Unicode 的 API，依然可以实现上述功能。

##### 清单 4\. 对长路径名文件操作的 c 的示例代码（ShortPath）

```
{
char pathName [1024];
strcpy(pathName,"\\\\?\\E:\\VerylongpathVerylongpathVerylongpathVerylongpathVerylongpathV
erylongpathVerylongpathVerylongpathVerylongpathVerylongpathVerylongpath\\VerylongpathVeryl
ongpathVerylongpathVerylongpathVerylongpathVerylongpathVerylongpathVerylongpathVerylongpat
hVerylongpathVerylongpath.txt");

const int MaxPathLength = 2048;
char shortPath[MaxPathLength];

if (strlen(pathName) >= MAX_PATH)
{
    char prePath[] = "\\\\?\\";
    if (strlen(pathName) >= MaxPathLength - strlen(pathName))
        return false;

    sprintf(shortPath, "%s%s", prePath, pathName);

    for (int iPathIndex = 0; iPathIndex < strlen(shortPath); iPathIndex++)
        if (shortPath[iPathIndex] == '/')
            shortPath[iPathIndex] = '\\';

    int dwlen = GetShortPathName(shortPath, shortPath, MaxPathLength);
    if (dwlen <= 0)
        return false;
}
}

```

Show moreShow more icon

经过上述的代码，超过 `MAX_PATH` 限制的路径名都可以转变成一个 8.3 格式的短路径名，可以把这个文件名 (`shortPath)` 作为后续文件操作函数的参数。这种情况下，对于该文件的所有操作都可以被支持了。我们用这种缩短路径名长度的方式解决了长路径名文件的操作问题。

## 结束语

本文首先列出了不同的 JDK 版本在 Windows 操作系统上对于长路径名文件处理的区别，同时指出了 JDK 5.0 开始才完全支持长路径名；在第二部分中给出了两种支持长路径名文件的 C/C++ 编程方法。使用上文中的任一方法，我们都可以实现对长路径名文件的操作，这将在很大程度上方便我们的开发工作，解决在 Windows 平台上标准 API 函数对长路径名文件支持的局限性问题。

## 声明

以上实验代码仅在 Windows XP 操作系统和 VC.NET 编译环境中测试通过，作者不对其提供任何种类的保证。如果有任何问题，欢迎来信与作者讨论。