# 从命名风格等方面解读阿里巴巴 Java 代码规范
阿里巴巴的 Java 代码规范深度解读的上篇

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/deconding-code-specification-part-1/)

周明耀

发布: 2018-01-10

* * *

## 前言

2017 年阿里云栖大会，阿里发布了针对 Java 程序员的《阿里巴巴 Java 开发手册（终极版）》，这篇文档作为阿里数千位 Java 程序员的经验积累呈现给公众，并随之发布了适用于 Eclipse 和 Intellim 的代码检查插件。为了能够深入了解 Java 程序员编码规范，也为了深入理解为什么阿里这样规定，是否规定有误，本文以阿里发布的这篇文档作为分析起源，扩大范围至业界其他公司的规范，例如谷歌、FaceBook、微软、百度、华为，并搜索网络上技术大牛发表的技术文章，深入理解每一条规范的设计背景和目标。

由于解读文章仅有两篇，所以按照阿里的篇幅权重分为上篇仅针对 Java 语言本身的编码规约，下篇包含日志管理、异常处理、单元测试、MySQL 规范、工程规范等方面内容进行解读。本文是上篇，主要针对编码规约部分进行解读，由于篇幅限制，仅挑选一小部分进行解读，如果需要全篇，请联系本文作者。

## 编码规约

### 命名风格

#### 下划线或美元符号

阿里强制规定代码中的命名均不能以下划线或美元符号开始，也不能以下划线或美元符号结束。

例如以下为错误，如清单 1 所示：

##### 清单 1 错误示例

```
_name/__name/$Object/name_/name$/Object$。

```

Show moreShow more icon

##### 我的理解 1

Oracle 官网建议不要使用$或者 _开始变量命名，并且建议在命名中完全不要使用”$”字符，原文是”The convention,however,is to always begin your variable names with a letter,not ‘$’ or ‘_‘”。对于这一条，腾讯的看法是一样的，百度认为虽然类名可以支持使用”$”符号，但只在系统生成中使用（如匿名类、代理类），编码不能使用。

这类问题在 StackOverFlow 上有很多人提出，主流意见为人不需要过多关注，只需要关注原先的代码是否存在” _“，如果存在就继续保留，如果不存在则尽量避免使用。也有一位提出尽量不适用”_“的原因是低分辨率的显示器，肉眼很难区分” _“（一个下划线）和”\__“（两个下划线）。

我个人觉得可能是由于受 C 语言的编码规范所影响。因为在 C 语言里面，系统头文件里将宏名、变量名、内部函数名用 _开头，因为当你#include 系统头文件时，这些文件里的名字都有了定义，如果与你用的名字冲突，就可能引起各种奇怪的现象。综合各种信息，建议不要使用”_“、”$”、空格作为命名开始，以免不利于阅读或者产生奇怪的问题。

#### 类命名

阿里强制规定类名使用 UpperCamelCase 风格，必须遵从驼峰形式，但以下情形例外：DO/BO/DTO/VO/AO。

#### 清单 2 类命名例子

```
正例：MarcoPolo/UserDO/XmlService/TcpUdpDeal/TarPromotion
反例：macroPolo/UserDo/XMLService/TCPUDPD/TAPromotion

```

Show moreShow more icon

##### 我的理解 2

百度除了支持阿里的规范以外，规定虽然类型支持”$”符号，但只在系统生成中使用（如匿名类、代理类），编码中不能使用。

对于类名，俄罗斯 Java 专家 Yegor Bugayenko 给出的建议是尽量采用现实生活中实体的抽象，如果类的名字以”-er”结尾，这是不建议的命名方式。他指出针对这一条有一个例外，那就是工具类，例如 StringUtils、FileUtils、IOUtils。对于接口名称，不要使用 IRecord、IfaceEmployee、RedcordInterface，而是使用现实世界的实体命名。如清单 3 所示。

##### 清单 3 示例

```
Class SimpleUser implements User{};
Class DefaultRecord implements Record{};
Class Suffixed implements Name{};
Class Validated implements Content{};

```

Show moreShow more icon

#### 抽象类的命名

阿里强制规定抽象类命名使用 Abstratc 或 Base 开头。

##### 我的理解 3

Oracle 的抽象类和方法规范并没有要求必须采用 Abstract 或者 Base 开头命名，事实上官网上的示例没有这种命名规范要求，如清单 4 所示。

#### 清单 4 示例

```
public abstract class GraphicObject{
//declare fields
//declare nonabstract methods
    abstract void draw();
}

```

Show moreShow more icon

我也查了一下 JDK，确实源码里很多类都是以这样的方式命名的，例如抽象类 java.util.AbstractList。

Stackoverflow 上对于这个问题的解释是，由于这些类不会被使用，一定会由其他的类继承并实现内部细节，所以需要明白地告诉读者这是一个抽象类，那以 Abstract 开头比较合适。

JoshuaBloch的理解是支持以 Abstract 开头。我的理解是不要以 Base 开头命名，因为实际的基类也以 Base 开头居多，这样意义有多样性，不够直观。

### 常量定义

#### 避免魔法值的使用

阿里强制规定不允许任何魔法值（未经定义的常量）直接出现在代码中，反例如清单 5 所示。

##### 清单 5 反例

```
String key = "Id#taobao_" + tradeId；
cache.put(key,value);

```

Show moreShow more icon

##### 我的理解 4

魔法值确实让你很疑惑，比如你看下面这个例子：

int priceTable[] = new int[16];//这样定义错误；这个 16 究竟代表什么？

正确的定义方式是这样的：

static final int PRICE\_TABLE\_MAX = 16; //这样定义正确，通过使用完整英语单词的常量名明确定义

int price Table[] = new int[PRICE\_TABLE\_MAX];

魔法值会让代码的可读性大大降低，而且如果同样的数值多次出现时，容易出现不清楚这些数值是否代表同样的含义。另一方面，如果本来应该使用相同的数值，一旦用错，也难以发现。因此可以采用以下两点，极力避免使用魔法数值。

1. 不适用魔法数值，使用带名字的 Static final 或者 enum 值；

2. 原则上 0 不用于魔法值，这是因为 0 经常被用作数组的最小下标或者变量初始化的缺省值。


#### 变量值范围

阿里推荐如果变量值仅在一个范围内变化，且带有名称之外的延伸属性，定义为枚举类。下面这个正例中的数字就是延伸信息，表示星期几。正例如清单 6 所示。

##### 清单 6 正例

```
public Enum {MONDAY(1),TUESDAY(2),WEDNESDAY(3),THURSDAY(4),FRIDAY(5),SATURDAY(6),SUNDAY(7);}

```

Show moreShow more icon

##### 我的理解 5

对于固定并且编译时对象，如 Status、Type 等，应该采用 enum 而非自定义常量实现，enum 的好处是类型更清楚，不会再编译时混淆。这是一个建议性的试用推荐，枚举可以让开发者在 IDE 下使用更方便，也更安全。另外就是枚举类型是一种具有特殊约束的类类型，这些约束的存在使得枚举类本身更加简洁、安全、便捷。

### 代码格式

#### 大括号的使用约定

阿里强制规定如果是大括号为空，则简洁地写成{}即可，不需要换行；如果是非空代码块则：

1. 左大括号前不换行

2. 左大括号后换行

3. 右大括号前换行

4. 右大括号后还有 else 等代码则不换行表示终止的右大括号后必须换行


##### 我的理解 6

阿里的这条规定应该是参照了 SUN 公司 1997 年发布的代码规范（SUN 公司是 JAVA 的创始者），Google 也有类似的规定，大家都是遵循 K&R 风格（Kernighan 和 Ritchie），Kernighan 和 Ritchie 在《The C Programming Language》一书中推荐这种风格，JAVA 语言的大括号风格就是受到了 C 语言的编码风格影响。

注意，SUN 公司认为方法名和大括号之间不应该有空格。

#### 单行字符数限制

阿里强制规定单行字符数限制不超过 120 个，超出需要换行，换行时遵循如下原则：

1. 第二行相对第一行缩进 4 个空格，从第三行开始，不再继续缩进，参考示例。

2. 运算符与下文一起换行。

3. 方法调用的点符号与下文一起换行。

4. 方法调用时，多个参数，需要换行时，在逗号后进行。

5. 在括号前不要换行，见反例。


如清单 7 所示。

#### 清单 7 示例

```
StringBuffer sb = new StringBuffer();
//超过 120 个字符的情况下，换行缩进 4 个空格，点号和方法名称一起换行
sb.append("zi").append("xin")...
.append("huang")...
.append("huang")...
.append("huang")...
反例：
StringBuffer sb = new StringBuffer();
//超过 120 个字符的情况下，不要在括号前换行
sb.append("zi").append("xin").append
("huang");
//参数很多的方法调用可能超过 120 个字符，不要在逗号前换行
method(args1,args2,args3,....,argsX);

```

Show moreShow more icon

##### 我的理解 7

的规范中指出单行不要超过 80 个字符，对于文档里面的代码行，规定不要超过 70 个字符单行。当表达式不能在一行内显示的时候，genuine 以下原则进行切分：

1. 在逗号后换行；

2. 在操作符号前换行；

3. 倾向于高级别的分割；

4. 尽量以描述完整作为换行标准；

5. 如果以下标准造成代码阅读困难，直接采用 8 个空格方式对第二行代码留出空白。


示例代码如清单 8 所示。

#### 清单 8 示例

```
function(longExpression1, longExpression2, longExpression3,
longExpression4, longExpression5);
var = function(longExpression1,
function2(longExpression2,
longExpression3));
longName1 = longName2 * (longName3 + longName4 – longName5)
+ 4 * longName6;//做法正确
longName1 = longName2 * (longName3 + longName4
– longName5) + 4 * longName6;//做法错误
if ((condition1 && condition2)
|| (condition3 && condition4)
|| !(condition5 && condition6) {
doSomethingAboutIt();
}//这种做法错误
if ((condition1 && condition2)
|| (condition3 && condition4)
        || !(condition5 && condition6) {
    doSomethingAboutIt();
    }//这种做法正确
if ((condition1 && condition2) || (condition3 && condition4)
|| !(condition5 && condition6) {
    doSomethingAboutIt();
}//这种做法正确

```

Show moreShow more icon

### OOP 规约

#### 静态变量及方法调用

阿里强制规定代码中避免通过一个类的对象引用访问此类的静态变量或静态方法，暂时无谓增加编译器解析成本，直接用类名来访问即可。

##### 我的理解 8

谷歌公司在代码规范中指出必须直接使用类名对静态成员进行引用，并同时举例说明，如清单 9 所示。

##### 清单 9 示例

```
Foo aFoo =...;
Foo.aStaticMethod();//good
aFoo.aStaticMethod();//bad
somethingThatYieldsAFoo().aStaticMethod();//very bad

```

Show moreShow more icon

SUN 公司 1997 年发布的代码规范也做了类似的要求。

为什么需要这样做呢？因为被 static 修饰过的变量或者方法都是随着类的初始化产生的，在堆内存中有一块专门的区域用来存放，后续直接用类名访问即可，避免编译成本的增加和实例对象存放空间的浪费。

StackOverflow 上也有人提出了相同的疑问，网友较为精辟的回复是”这是由于生命周期决定的，静态方法或者静态变量不是以实例为基准的，而是以类为基准，所以直接用类访问，否则违背了设计初衷”。那为什么还保留了实例的访问方式呢？可能是因为允许应用方无污染修改吧。

#### 可变参数编程

阿里强制规定相同参数类型、相同业务类型，才可以使用 Java 的可变参数，避免使用 Object，并且要求可变参数必须放置在参数列表的最后（提倡同学们尽量不用可变参数编程）。

##### 我的理解 9

我们先来了解可变参数的使用方式：

1. 在方法中定义可变参数后，我们可以像操作数组一样操作该参数。

2. 如果该方法除了可变参数还有其他的参数，可变参数必须放到最后。

3. 拥有可变参数的方法可以被重载，在被调用时，如果能匹配到参数定长的方法则优先调用参数定长的方法。

4. 可变参数可以兼容数组参数，但数组参数暂时无法兼容可变参数。


至于为什么可变参数需要被放在最后一个，这是因为参数个数不定，所以当其后还有相同类型参数时，编译器无法区分传入的参数属于前一个可变参数还是后边的参数，所以只能让可变参数位于最后一项。

可变参数编程有一些好处，例如反射、过程建设、格式化等。对于阿里同学提出的尽量不使用可变参数编程，我猜测的原因是不太可控，比如 Java8 推出 Lambda 表达式之后，可变参数编程遇到了实际的实现困难。

我们来看一个例子。假设我们想要实现以下功能，如清单 10 所示。

#### 清单 10 实现功能

```
test((arg0,arg1) -> me.call(arg0,arg1));
>test((arg0,arg1,arg2)->me.call(arg0,arg1,arg2));
...

```

Show moreShow more icon

对应的实现定义接口的继承关系，并且使用默认方法避免失败，如清单 11 所示。

#### 清单 11 实现方式代码段 1

```
interface VarArgsRunnable{
    default void run(Object...arguments){
         throw new UnsupportedOperationException("not possible");
    }
    default int getNumberOfArguments(){
         throw new UnsupportedOperationException("unknown");
    }
}
@FunctionalInterface
Interface VarArgsRunnable4 extends VarArgsRnnable {
    @Override
    default void run(Object...arguments){
        assert(arguments.length == 4);
        run(arguments[0], arguments[1], arguments[2], arguments[3]);
}
    void run(Object arg0, Object arg1, Object arg2, Object arg3, Object arg4);
    @Override
    default int getNumberOfArguments(){
    return 4;
}
}

```

Show moreShow more icon

这样我们就可以定义 11 个接口，从 VarArgsRnnable0 到 VarArgsRnnable10，并且覆盖方法，调用方式如清单 12 所示。

##### 清单 12 实现方式代码段 2

```
public void myMethod(VarArgsRnnable runnable,Object...arguments){
     runnable.run(arguments);
}

```

Show moreShow more icon

针对上述需求，我们也可以编写代码如清单 13 所示。

#### 清单 13 实现方式 2 代码段

```
public class Java8VariableArgumentsDemo{
    interface Invoker{
         void invoke(Object...args);
}
    public static void invokeInvoker(Invoker invoker,Object...args){
        invoker.invoke(args);
    }
    public static void applyWithStillAndPrinting(Invoker invoker){
    invoker.invoke("Still","Printing");
    }
    Public static void main(String[] args){
        Invoker printer = new Invoker(){
    Public void invoke(Object...args){
        for(Object arg:args){
        System.out.println(arg);
        }
    }
    };
         printer.invoke("I","am","printing");
         invokeInvoker(printer, "Also","printing");
         applyWithStillAndPrinting(printer);
         applyWithStillAndPrinting((Object...args)->System.out.println("Not done"));
         applyWithStillAndPrinting(printer::invoke);
    }
}

```

Show moreShow more icon

运行后输出如清单 14 所示。

#### 清单 14 实现方式 2 代码段运行结果

```
I
am
printing
Also
printing
Still
Printing
Not done
Still
Printing

```

Show moreShow more icon

### 并发处理

#### 单例模式需要保证线程安全

阿里强制要求获取单例对象需要保证线程安全，其中的方法也要保证线程安全，并进一步说明资源驱动类、工具类、单例工厂类都需要注意。

##### 我的理解 10

对于这一条规范是通识化规定，我这里进一步讲讲如何做好针对单例对象的线程安全，主要有以下几种方式：

1. 方法中申明 synchronized 关键字

     出现非线程安全问题，是由于多个线程可以同时进入 getInstance()方法，那么只需要对该方法进行 synchronized 锁同步即可，如清单 15 所示。


    #### 清单 15 synchronized 关键字方式






    ```
    public class MySingleton{
         private static MySingleton instance = null;
         private MySingleton(){}
             public synchronized static MySingleton getInstance(){
             try{
                 if(instance != null){//懒汉式
             }else{
                 //创建实例之前可能会有一些准备性的耗时工作
                 Thread.sleep(500);
                 Instance = new MySingleton();
              }
             }catch(InterruptedException e){
                 e.printStackTrace();
             }
                 return instance;
         }
    }

    ```





    Show moreShow more icon

     执行结果如清单 16 所示。


    #### 清单 16 synchronized 关键字方式运行结果






    ```
    174342932
    174342932
    174342932
    174342932
    174342932
    174342932

    ```





    Show moreShow more icon

     从执行结果上来看，多线程访问的问题已经解决了，返回的是一个实例。但是这种实现方式的运行效率很低。我们接下来采用同步方法块实现。

2. 同步方法块实现


    #### 清单 17 同步方法块方式






    ```
    public class MySingleton {
         private static MySingleton instance = null;
         private MySingleton(){}
         //public synchronized static MySingleton getInstance() {
         public static MySingleton getInstance() {
             try {
                 synchronized (MySingleton.class) {
                         if(instance != null){//懒汉式
                 }else{
                         //创建实例之前可能会有一些准备性的耗时工作
                     Thread.sleep(300);
                     instance = new MySingleton();
                 }
                 }
                 } catch (InterruptedException e) {
                     e.printStackTrace();
             }
          return instance;
         }
    }

    ```





    Show moreShow more icon

     这里的实现能够保证多线程并发下的线程安全性，但是这样的实现将全部的代码都被锁上了，同样的效率很低下。

3. 针对某些重要的代码来进行单独的同步

     针对某些重要的代码进行单独的同步，而不是全部进行同步，可以极大的提高执行效率，代码如清单 18 所示。


    ##### 清单 18 单独同步方式






    ```
    public class MySingleton {
         private static MySingleton instance = null;
         private MySingleton(){}
         public static MySingleton getInstance() {
             try {
                 if(instance != null){//懒汉式
                 }else{
                     //创建实例之前可能会有一些准备性的耗时工作
                 Thread.sleep(300);
                 synchronized (MySingleton.class) {
                 instance = new MySingleton();
                 }
             }
             }catch (InterruptedException e) {
                 e.printStackTrace();
             }
         return instance;
         }
    }

    ```





    Show moreShow more icon

     从运行结果来看，这样的方法进行代码块同步，代码的运行效率是能够得到提升，但是却没能保住线程的安全性。看来还得进一步考虑如何解决此问题。

4. 双检查锁机制（Double Check Locking）

     为了达到线程安全，又能提高代码执行效率，我们这里可以采用 DCL 的双检查锁机制来完成，代码实现如清单 19 所示。


    #### 清单 19 双检查锁机制






    ```
    public class MySingleton {
         //使用 volatile 关键字保其可见性
         volatile private static MySingleton instance = null;
         private MySingleton(){}
         public static MySingleton getInstance() {
             try {
                 if(instance != null){//懒汉式
             }else{
                 //创建实例之前可能会有一些准备性的耗时工作
                 Thread.sleep(300);
                 synchronized (MySingleton.class) {
                 if(instance == null){//二次检查
                 instance = new MySingleton();
             }
             }
             }
             } catch (InterruptedException e) {
                 e.printStackTrace();
             }
         return instance;
         }
    }

    ```





    Show moreShow more icon

     这里在声明变量时使用了 volatile 关键字来保证其线程间的可见性；在同步代码块中使用二次检查，以保证其不被重复实例化。集合其二者，这种实现方式既保证了其高效性，也保证了其线程安全性。

5. 静态内置类方式

     DCL 解决了多线程并发下的线程安全问题，其实使用其他方式也可以达到同样的效果，代码实现如清单 20 所示。


    #### 清单 20 静态内置类方式






    ```
    public class MySingleton {
         //内部类
         private static class MySingletonHandler{
             private static MySingleton instance = new MySingleton();
         }
         private MySingleton(){}
         public static MySingleton getInstance() {
             return MySingletonHandler.instance;
         }
    }

    ```





    Show moreShow more icon

6. 序列化与反序列化方式

     静态内部类虽然保证了单例在多线程并发下的线程安全性，但是在遇到序列化对象时，默认的方式运行得到的结果就是多例的。


    #### 清单 21 序列化与反序列化






    ```
    import java.io.Serializable;
    public class MySingleton implements Serializable {
         private static final long serialVersionUID = 1L;
         //内部类
         private static class MySingletonHandler{
             private static MySingleton instance = new MySingleton();
         }
         private MySingleton(){}
         public static MySingleton getInstance() {
             return MySingletonHandler.instance;
         }
    }

    ```





    Show moreShow more icon

7. 使用枚举数据类型方式

     枚举 enum 和静态代码块的特性相似，在使用枚举时，构造方法会被自动调用，利用这一特性也可以实现单例。


    #### 清单 22 枚举数据方式 1






    ```
    public enum EnumFactory{
         singletonFactory;
         private MySingleton instance;
         private EnumFactory(){//枚举类的构造方法在类加载是被实例化
             instance = new MySingleton();
         }
         public MySingleton getInstance(){
             return instance;
         }
    }
         class MySingleton{//需要获实现单例的类，比如数据库连接 Connection
         public MySingleton(){}
    }

    ```





    Show moreShow more icon

     这样写枚举类被完全暴露了，据说违反了”职责单一原则”，我们可以按照下面的代码改造。


    #### 清单 23 枚举数据方式 2






    ```
    public class ClassFactory{
         private enum MyEnumSingleton{
             singletonFactory;
             private MySingleton instance;
             private MyEnumSingleton(){//枚举类的构造方法在类加载是被实例化
                 instance = new MySingleton();
             }
             public MySingleton getInstance(){
                 return instance;
             }
         }
         public static MySingleton getInstance(){
             return MyEnumSingleton.singletonFactory.getInstance();
         }
    }
    class MySingleton{//需要获实现单例的类，比如数据库连接 Connection
         public MySingleton(){}
    }

    ```





    Show moreShow more icon


### 控制语句

#### Switch 语句的使用

阿里强制规定在一个 switch 块内，每个 case 要么通过 break/return 等来终止，要么注释说明程序将继续执行到哪一个 case 为止；在一个 switch 块内，都必须包含一个 default 语句并且放在最后，即使它什么代码也没有。

##### 我的理解 11

首先理解前半部分，”每个 case 要么通过 break/return 等来终止，要么注释说明程序将继续执行到哪一个 case 为止”。因为这样可以比较清楚地表达程序员的意图，有效防止无故遗漏的 break 语句。我们来看一个示例，如清单 24 所示。

#### 清单 24 synchronized 关键字方式运行结果

```
switch(condition){
case ABC:
statements;
/*程序继续执行直到 DEF 分支*/
case DEF:
statements;
break;
case XYZ:
statements;
break;
default:
statements;
break;
}

```

Show moreShow more icon

上述示例中，每当一个 case 顺着往下执行时（因为没有 break 语句），通常应在 break 语句的位置添加注释。上面的示例代码中就包含了注释”/ _程序继续执行直到 DEF 分支_/”（这一条也是 SUN 公司 1997 年代码规范的要求）。

语法上来说，default 语句中的 break 是多余的，但是如果后续添加额外的 case，可以避免找不到匹配 case 项的错误。

### 集合处理

#### 集合转数组处理

阿里强制规定使用集合转数组的方法，必须使用集合的 toArray(T[] arrays)，传入的是类型完全一样的数组，大小就是 list.size()。使用 toArray 带参方法，入参分配的数组空间不够大时，toArray 方法内部将重新分配内存空间，并返回新数组地址；如果数组元素大于实际所需，下标为[list.size()]的数组元素将被置为 null，其它数组元素保持原值，因此最好将方法入参数组大小定义与集合元素个数一致。正例如清单 25 所示。

#### 清单 25 正例

```
List<String> list = new ArrayList<String>(2);
list.add("guan");
list.add("bao");
String[] array = new String[list.size()];
array = list.toArray(array);

```

Show moreShow more icon

反例：直接使用 toArray 暂时无参方法存在问题，此方法返回值只能是 Object[]类，若强转其他类型数组将出现 ClassCastException 错误。

##### 我的理解 12

ArrayList 类的 toArray()源码如清单所示，toArray()方法暂时无需传入参数，可以直接将集合转成 Object 数组进行返回，而且也只能返回 Object 类型。

##### 清单 26 toArray()源码

```
Public Object[] toArray(){
    Object aobj[] = new Object[size];
    System.arraycopy(((Object)(elementData)),0,((Object)(aobj)),0,size);
    return aobj;
}
public <T> T[] toArray(T[] a){
    if(a.length < size)
       //Makeanewarrayofa'sruntimetype,butmycontents:
    return(T[])Arrays.copyOf(elementData,size,a.getClass());
    System.arraycopy(elementData,0,a,0,size);
    if(a.length> a[size]=null;
    returna;
}

```

Show moreShow more icon

由源码可知，不带参数的 toArray()构造一个 Object 数组，然后进行数据拷贝，此时进行转型就会产生 ClassCastException。原因是不能将 Object[]转化为 Strng[]。Java 中的强制类型转换只是针对单个对象，想要将一种类型数组转化为另一种类型数组是不可行的。

针对传入参数的数组大小，测试大于 list、等于 list 和小于 list 三种情况，测试代码如清单 27 所示。

#### 清单 27 toArray()测试

```
public static void main(String[] args){
    List<String> list = new ArrayList<String>();
    for(int i=0;i<20;i++){
        list.add("test");
    }
    long start = System.currentTimeMills();
    for(int i=0;i<10000000;i++){
        String[] array = new String[list.size()];
        Array = list.toArray(array);
    }
    System.out.println("数组长度等于 list 耗时："+(System.currentTimeMills()-start)+"ms");
    start = System.currentTimeMills();
    for(int i=0;i<10000000;i++){
        String[] array = new String[list.size()*2];
        Array = list.toArray(array);
    }
     System.out.println("数组长度等于 list 耗时："+(System.currentTimeMills()-start)+"ms");
     start = System.currentTimeMills();
    for(int i=0;i<10000000;i++){
        String[] array = new String[0];
        Array = list.toArray(array);
    }
        System.out.println("数组长度等于 list 耗时："+(System.currentTimeMills()-start)+"ms");
}

```

Show moreShow more icon

清单运行后输出结果如清单 28 所示。

##### 清单 28 清单运行输出

```
数组长度等于 list 耗时：431ms
数组长度等于 list 耗时：509ms
数组长度等于 list 耗时：1943ms

```

Show moreShow more icon

通过测试可知无论数据大小如何，数组转换都可以成功，只是耗时不同，数组长度等于 list 时性能最优，因此强制方法入参数组大小与集合元素个数一致。

### 注释规约

#### 方法注释要求

阿里强制要求方法内部单行注释，在被注释语句上方另起一行，使用//注释。方法内部多行注释使用/\*\*/注释，注意与代码对照。

##### 我的理解 13

百度规定方法注释采用标准的 Javadoc 注释规范，注释中必须提供方法说明、参数说明及返回值和异常说明。腾讯规定采用 JavaDoc 文档注释，在方法定义之前应该对其进行注释，包括方法的描述、输入、输出以及返回值说明、抛出异常说明、参考链接等。

### 其他

#### 数据结构初始化大小

阿里推荐任何数据结构的构造或初始化，都应指定大小，避免数据结构暂时无限增长吃光内存。

##### 我的理解 14

首先明确一点，阿里这里指的大小具体是指数据结构的最大长度。大部分 Java 集合类在构造时指定的大小都是初始尺寸（initial Capacity），而不是尺寸上限（Capacity），只有几种队列除外，例如 ArrayBlockingQueue、LinkedBlockingQueue，它们在构造时可以指定队列的最大长度。阿里推荐的目的是为了合理规划内存，避免出现 OOM（Out of Memory）异常。

## 结束语

本文主要介绍了阿里巴巴针对命名风格、常量定义、代码格式、OOP 规约、并发处理、控制语句、集合处理、注释规约、其他这些关于编码规约的要求。本文仅覆盖了阿里代码规范的少数内容，更多内容请咨询本文作者。

## 参考资源

参考文档《阿里巴巴 Java 开发手册（又名阿里巴巴 Java 代码规约）》。

参考 IBM Developer 上的 Java 文章，了解更多 Java 知识。

参考书籍 《Effective Java Second Edition》Joshua Bloch。