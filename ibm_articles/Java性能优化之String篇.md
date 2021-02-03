# Java 性能优化之 String 篇
通过分析 String 在 JVM 中的存储结构，以及常见 String 操作对内存的影响阐述问题产生的原因及解决

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-lo-optmizestring/)

杨博文, 应乐年, 杨雯雯

发布: 2012-05-14

* * *

## String 在 JVM 的存储结构

一般而言，Java 对象在虚拟机的结构如下：

- 对象头（object header）：8 个字节
- Java 原始类型数据：如 int, float, char 等类型的数据，各类型数据占内存如表 1。
- 引用（reference）：4 个字节
- 填充符（padding）

表 1\. Java 各数据类型所占内存

数据类型占用内存（字节数）boolean 和 byte1char 和 short2int 和 float4long 和 double8

然而，一个 Java 对象实际还会占用些额外的空间，如：对象的 class 信息、ID、在虚拟机中的状态。在 Oracle JDK 的 Hotspot 虚拟机中，一个普通的对象需要额外 8 个字节。

如果对于 String（JDK 6）的成员变量声明如下：

```
private final char value[];
private final int offset;
private final int count;
private int hash;

```

Show moreShow more icon

那么因该如何计算该 String 所占的空间？

首先计算一个空的 char 数组所占空间，在 Java 里数组也是对象，因而数组也有对象头，故一个数组所占的空间为对象头所占的空间加上数组长度，即 8 + 4 = 12 字节 , 经过填充后为 16 字节。

那么一个空 String 所占空间为：

对象头（8 字节）+ char 数组（16 字节）+ 3 个 int（3 × 4 = 12 字节）+1 个 char 数组的引用 (4 字节 ) = 40 字节。

因此一个实际的 String 所占空间的计算公式如下：

```
8*( ( 8+2*n+4+12)+7 ) / 8 = 8*(int) ( ( ( (n) *2 )+43) /8 )

```

Show moreShow more icon

其中，n 为字符串长度。

## 案例分析

在我们的大规模文本分析的案例中，程序需要统计一个 300MB 的 csv 文件所有单词的出现次数，分析发现共有 20,000 左右的唯一单词，假设每个单词平均包含 15 个字母，这样根据上述公式，一个单词平均占用 75 bytes. 那么这样 75 \* 20,000 = 1500000，即约为 1.5M 左右。但实际发现有上百兆的空间被占用。 实际使用的内存之所以与预估的产生如此大的差异是因为程序大量使用 `String.split()` 或 `String.substring()` 来获取单词。在 JDK 1.6 中 `String.substring(int, int)` 的源码为：

```
public String substring(int beginIndex, int endIndex) {
      if (beginIndex < 0) {
           throw new StringIndexOutOfBoundsException(beginIndex);
      }
      if (endIndex > count) {
           throw new StringIndexOutOfBoundsException(endIndex);
      }
      if (beginIndex > endIndex) {
           throw new StringIndexOutOfBoundsException(endIndex - beginIndex);
      }
      return ((beginIndex == 0) && (endIndex == count)) ? this :
           new String(offset + beginIndex, endIndex - beginIndex, value);
}

```

Show moreShow more icon

调用的 String 构造函数源码为：

```
String(int offset, int count, char value[]) {
this.value = value;
this.offset = offset;
this.count = count;
}

```

Show moreShow more icon

仔细观察粗体这行代码我们发现 `String.substring()` 所返回的 String 仍然会保存原始 String, 这就是 20,000 个平均长度的单词竟然占用了上百兆的内存的原因。 一个 csv 文件中每一行都是一份很长的数据，包含了上千的单词，最后被 `String.split()` 或 `String.substring()` 截取出的每一个单词仍旧包含了其原先所在的上下文中，因而导致了出乎意料的大量的内存消耗。

当然，JDK String 的源码设计当然有着其合理之处，对于通过 `String.split()` 或 `String.substring()` 截取出大量 String 的操作，这种设计在很多时候可以很大程度的节省内存，因为这些 String 都复用了原始 String，只是通过 int 类型的 start, end 等值来标识每一个 String。 而对于我们的案例，从一个巨大的 String 截取少数 String 为以后所用，这样的设计则造成大量冗余数据。 因此有关通过 `String.split()` 或 `String.substring()` 截取 String 的操作的结论如下：

- 对于从大文本中截取少量字符串的应用， `String.substring()` 将会导致内存的过度浪费。
- 对于从一般文本中截取一定数量的字符串，截取的字符串长度总和与原始文本长度相差不大，现有的 `String.substring()` 设计恰好可以共享原始文本从而达到节省内存的目的。

既然导致大量内存占用的根源是 `String.substring()` 返回结果中包含大量原始 String，那么一个显而易见的减少内存浪费的的途径就是去除这些原始 String。办法有很多种，在此我们采取比较直观的一种，即再次调用 `newString` 构造一个的仅包含截取出的字符串的 String，我们可调用 `String. toCharArray ()` 方法：

```
String newString = new String(smallString.toCharArray());

```

Show moreShow more icon

举一个极端例子，假设要从一个字符串中获取所有连续的非空子串，字符串长度为 n，如果用 JDK 本身提供的 `String.substring() 方` 法，则总共的连续非空子串个数为：

```
n+(n-1)+(n-2)+... +1 = n*(n+1)/2 =O(n2)

```

Show moreShow more icon

由于每个子串所占的空间为常数，故空间复杂度也为 O(n2)。

如果用本文建议的方法，即构造一个内容相同的新的字符串，则所需空间正比于子串的长度，则所需空间复杂度为：

```
1*n+2*(n-1)+3*(n-2)+... +n*1 = (n3+3*n2+2*n)/6 = O(n3)

```

Show moreShow more icon

所以，从以上定量的分析看来，当需要截取的字符串长度总和大于等于原始文本长度，本文所建议的方法带来的空间复杂度反而高了，而现有的 `String.substring()` 设计恰好可以共享原始文本从而达到节省内存的目的。反之，当所需要截取的字符串长度总和远小于原始文本长度时，用本文所推荐的方法将在很大程度上节省内存，在大文本数据处理中其优势显而易见。

## 其他 String 使用的优化建议

以上我们描述了在我们的大量文本分析案例中调用 String 的 `subString 方法` 导致内存消耗的问题，下面再列举一些其他将导致内存浪费的 String 的 API 的使用：

### String 拼接的方法选择

在拼接静态字符串时，尽量用 +，因为通常编译器会对此做优化，如：

```
String test = "this " + "is " + "a " + "test " + "string"

```

Show moreShow more icon

编译器会把它视为：

```
String test = "this is a test string"

```

Show moreShow more icon

在拼接动态字符串时，尽量用 `StringBuffer` 或 `StringBuilder` 的 `append` ，这样可以减少构造过多的临时 String 对象。

### String 构造的方法选择

常见的创建一个 String 可以用赋值操作符”=” 或用 new 和相应的构造函数。初学者一定会想这两种有何区别，举例如下：

```
String a1 = "Hello”;
String a2 = new String("Hello”);

```

Show moreShow more icon

第一种方法创建字符串时 JVM 会查看内部的缓存池是否已有相同的字符串存在：如果有，则不再使用构造函数构造一个新的字符串，直接返回已有的字符串实例；若不存在，则分配新的内存给新创建的字符串。

第二种方法直接调用构造函数来创建字符串，如果所创建的字符串在字符串缓存池中不存在则调用构造函数创建全新的字符串，如果所创建的字符串在字符串缓存池中已有则再拷贝一份到 Java 堆中。

尽管这是一个简单明显的例子，然而在实际项目中编程者却不那么容易洞察因为这两种方式的选择而带来的性能问题。

### 使用构造函数 string() 带来的内存性能隐患和缓解

仍然以之前的从 csv 文件中截取 String 为例，先前我们通过用 new String() 去除返回的 String 中附带的原始 String 的方法优化了 `subString` 导致的内存消耗问题。然而，当我们下意识地使用 `newString` 去构造一个全新的字符串而不是用赋值符来创建（重用）一个字符串时，就导致了另一个潜在的性能问题，即：重复创建大量相同的字符串。说到这里，您也许会想到使用缓存池的技术来解决这一问题，大概有如下两种方法：

方法一，使用 String 的 `intern()` 方法返回 JVM 对字符串缓存池里相应已存在的字符串引用，从而解决内存性能问题，但这个方法并不推荐！原因在于：首先， `intern()` 所使用的池会是 JVM 中一个全局的池，很多情况下我们的程序并不需要如此大作用域的缓存；其次，intern() 所使用的是 JVM heap 中 PermGen 相应的区域，在 JVM 中 PermGen 是用来存放装载类和创建类实例时用到的元数据。程序运行时所使用的内存绝大部分存放在 JVM heap 的其他区域，过多得使用 `intern()` 将导致 PermGen 过度增长而最后返回 `OutOfMemoryError` ，因为垃圾收集器不会对被缓存的 String 做垃圾回收。所以我们建议使用第二种方式。

方法二，用户自己构建缓存，这种方式的优点是更加灵活。创建 HashMap，将需缓存的 String 作为 key 和 value 存放入 HashMap。假设我们准备创建的字符串为 key，将 Map cacheMap 作为缓冲池，那么返回 key 的代码如下：

```
private String getCacheWord(String key) {
     String tmp = cacheMap.get(key);
     if(tmp != null) {
            return tmp;
     } else {
             cacheMap.put(key, key);
             return key;
     }
}

```

Show moreShow more icon

## 结束语

本文通过一个实际项目中遇到的因使用 String 而导致的性能问题讲述了 String 在 JVM 中的存储结构，String 的 API 使用可能造成的性能问题以及解决方法。相信这些建议能对处理大文本分析的朋友有所帮助，同时希望文中提到的某些优化方法能被举一反三的应用在其他有关 String 的性能优化的场合。