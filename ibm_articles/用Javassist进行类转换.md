# 用 Javassist 进行类转换
用 Javassist 转换字节码中的方法

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-dyn0916/)

Dennis M. Sosnoski

发布: 2003-10-25

* * *

讲过了 Java 类格式和利用反射进行的运行时访问后，本系列到了进入更高级主题的时候了。本月我将开始本系列的第二部分，在这里 Java 类信息只不过是由应用程序操纵的另一种形式的数据结构而已。我将这个主题的整个内容称为 _classworking_ 。

我将以 Javassist 字节码操作库作为对 classworking 的讨论的开始。Javassist 不仅是一个处理字节码的库，而且更因为它的另一项功能使得它成为试验 classworking 的很好的起点。这一项功能就是：可以用 Javassist 改变 Java 类的字节码，而无需真正了解关于字节码或者 Java 虚拟机(Java virtual machine JVM)结构的任何内容。从某方面将这一功能有好处也有坏处 — 我一般不提倡随便使用不了解的技术 — 但是比起在单条指令水平上工作的框架，它确实使字节码操作更可具有可行性了。

## Javassist 基础

Javassist 使您可以检查、编辑以及创建 Java 二进制类。检查方面基本上与通过 Reflection API 直接在 Java 中进行的一样，但是当想要修改类而不只是执行它们时，则另一种访问这些信息的方法就很有用了。这是因为 JVM 设计上并没有提供在类装载到 JVM 中后访问原始类数据的任何方法，这项工作需要在 JVM 之外完成。

##### 不要错过本系列的其余部分

第一部分：” [类和类装入](https://www.ibm.com/developerworks/cn/java/j-dyn0429/index.html)”(2003年4月)

第二部分：” [引入反射](https://www.ibm.com/developerworks/cn/java/j-dyn0603/index.html)” (2003年6月)

第三部分：” [应用返射](https://www.ibm.com/developerworks/cn/java/j-dyn0715/)” (2003年7月)

Javassist 使用 `javassist.ClassPool` 类跟踪和控制所操作的类。这个类的工作方式与 JVM 类装载器非常相似，但是有一个重要的区别是它不是将装载的、要执行的类作为应用程序的一部分链接，类池使所装载的类可以通过 Javassist API 作为数据使用。可以使用默认的类池，它是从 JVM 搜索路径中装载的，也可以定义一个搜索您自己的路径列表的类池。甚至可以直接从字节数组或者流中装载二进制类，以及从头开始创建新类。

装载到类池中的类由 `javassist.CtClass` 实例表示。与标准的 Java `java.lang.Class` 类一样， `CtClass` 提供了检查类数据（如字段和方法）的方法。不过，这只是 `CtClass` 的部分内容，它还定义了在类中添加新字段、方法和构造函数、以及改变类、父类和接口的方法。奇怪的是，Javassist 没有提供删除一个类中字段、方法或者构造函数的任何方法。

字段、方法和构造函数分别由 `javassist.CtField、 javassist.CtMethod` 和 `javassist.CtConstructor` 的实例表示。这些类定义了修改由它们所表示的对象的所有方法的方法，包括方法或者构造函数中的实际字节码内容。

### 所有字节码的源代码

Javassist 让您可以完全替换一个方法或者构造函数的字节码正文，或者在现有正文的开始或者结束位置选择性地添加字节码(以及在构造函数中添加其他一些变量)。不管是哪种情况，新的字节码都作为类 Java 的源代码声明或者 `String` 中的块传递。Javassist 方法将您提供的源代码高效地编译为 Java 字节码，然后将它们插入到目标方法或者构造函数的正文中。

Javassist 接受的源代码与 Java 语言的并不完全一致，不过主要的区别只是增加了一些特殊的标识符，用于表示方法或者构造函数参数、方法返回值和其他在插入的代码中可能用到的内容。这些特殊标识符以符号 `$` 开头，所以它们不会干扰代码中的其他内容。

对于在传递给 Javassist 的源代码中可以做的事情有一些限制。第一项限制是使用的格式，它必须是单条语句或者块。在大多数情况下这算不上是限制，因为可以将所需要的任何语句序列放到块中。下面是一个使用特殊 Javassist 标识符表示方法中前两个参数的例子，这个例子用来展示其使用方法：

```
{
System.out.println("Argument 1: " + $1);
System.out.println("Argument 2: " + $2);
}

```

Show moreShow more icon

对于源代码的一项更实质性的限制是不能引用在所添加的声明或者块外声明的局部变量。这意味着如果在方法开始和结尾处都添加了代码，那么一般不能将在开始处添加的代码中的信息传递给在结尾处添加的代码。有可能绕过这项限制，但是绕过是很复杂的 — 通常需要设法将分别插入的代码合并为一个块。

## 用 Javassist 进行 Classworking

作为使用 Javassist 的一个例子，我将使用一个通常直接在源代码中处理的任务：测量执行一个方法所花费的时间。这在源代码中可以容易地完成，只要在方法开始时记录当前时间、之后在方法结束时再次检查当前时间并计算两个值的差。如果没有源代码，那么得到这种计时信息就要困难得多。这就是 classworking 方便的地方 — 它让您对任何方法都可以作这种改变，并且不需要有源代码。

清单 1 显示了一个(不好的)示例方法，我用它作为我的计时试验的实验品： `StringBuilder` 类的 `buildString` 方法。这个方法使用一种所有 Java 性能优化的高手都会叫您 _不_ 要使用的方法构造一个具有任意长度的 `String`— 它通过反复向字符串的结尾附加单个字符来产生更长的字符串。因为字符串是不可变的，所以这种方法意味着每次新的字符串都要通过一个循环来构造：使用从老的字符串中拷贝的数据并在结尾添加新的字符。最终的效果是用这个方法产生更长的字符串时，它的开销越来越大。

##### 清单 1\. 需要计时的方法

```
public class StringBuilder
{
    private String buildString(int length) {
        String result = "";
        for (int i = 0; i < length; i++) {
            result += (char)(i%26 + 'a');
        }
        return result;
    }

    public static void main(String[] argv) {
        StringBuilder inst = new StringBuilder();
        for (int i = 0; i < argv.length; i++) {
            String result = inst.buildString(Integer.parseInt(argv[i]));
            System.out.println("Constructed string of length " +
                result.length());
        }
    }
}

```

Show moreShow more icon

### 添加方法计时

因为有这个方法的源代码，所以我将为您展示如何直接添加计时信息。它也作为使用 Javassist 时的一个模型。清单 2 只展示了 `buildString()` 方法，其中添加了计时功能。这里没有多少变化。添加的代码只是将开始时间保存为局部变量，然后在方法结束时计算持续时间并打印到控制台。

##### 清单 2\. 带有计时的方法

```
private String buildString(int length) {
        long start = System.currentTimeMillis();
        String result = "";
        for (int i = 0; i < length; i++) {
            result += (char)(i%26 + 'a');
        }
        System.out.println("Call to buildString took " +
            (System.currentTimeMillis()-start) + " ms.");
        return result;
    }

```

Show moreShow more icon

### 用 Javassist 来做

来做 使用 Javassist 操作类字节码以得到同样的效果看起来应该不难。Javassist 提供了在方法的开始和结束位置添加代码的方法，别忘了，我在为该方法中加入计时信息就是这么做的。

不过，还是有障碍。在描述 Javassist 是如何让您添加代码时，我提到添加的代码不能引用在方法中其他地方定义的局部变量。这种限制使我不能在 Javassist 中使用在源代码中使用的同样方法实现计时代码，在这种情况下，我在开始时添加的代码中定义了一个新的局部变量，并在结束处添加的代码中引用这个变量。

那么还有其他方法可以得到同样的效果吗？是的，我 _可以_ 在类中添加一个新的成员字段，并使用这个字段而不是局部变量。不过，这是一种糟糕的解决方案，在一般性的使用中有一些限制。例如，考虑在一个递归方法中会发生的事情。每次方法调用自身时，上次保存的开始时间值就会被覆盖并且丢失。

幸运的是有一种更简洁的解决方案。我可以保持原来方法的代码不变，只改变方法名，然后用原来的方法名增加一个新方法。这个 _拦截器(interceptor_)方法可以使用与原来方法同样的签名，包括返回同样的值。清单 3 展示了通过这种方法改编后源代码看上去的样子：

##### 清单 3\. 在源代码中添加一个拦截器方法

```
private String buildString$impl(int length) {
        String result = "";
        for (int i = 0; i < length; i++) {
            result += (char)(i%26 + 'a');
        }
        return result;
    }
    private String buildString(int length) {
        long start = System.currentTimeMillis();
        String result = buildString$impl(length);
        System.out.println("Call to buildString took " +
            (System.currentTimeMillis()-start) + " ms.");
        return result;
    }

```

Show moreShow more icon

通过 Javassist 可以很好地利用这种使用拦截器方法的方法。因为整个方法是一个块，所以我可以毫无问题地在正文中定义并且使用局部变量。为拦截器方法生成源代码也很容易 — 对于任何可能的方法，只需要几个替换。

### 运行拦截

实现添加方法计时的代码要用到在 [Javassist 基础](#javassist-基础) 中描述的一些 Javassist API。清单 4 展示了该代码，它是一个带有两个命令行参数的应用程序，这两个参数分别给出类名和要计时的方法名。 `main()` 方法的正文只给出类信息，然后将它传递给 `addTiming()` 方法以处理实际的修改。 `addTiming()` 方法首先通过在名字后面附加”`$impl”` 重命名现有的方法，接着用原来的方法名创建该方法的一个拷贝。然后它用含有对经过重命名的原方法的调用的计时代码替换拷贝方法的正文。

##### 清单 4\. 用 Javassist 添加拦截器方法

```
public class JassistTiming
{
    public static void main(String[] argv) {
        if (argv.length == 2) {
            try {

                // start by getting the class file and method
                CtClass clas = ClassPool.getDefault().get(argv[0]);
                if (clas == null) {
                    System.err.println("Class " + argv[0] + " not found");
                } else {

                    // add timing interceptor to the class
                    addTiming(clas, argv[1]);
                    clas.writeFile();
                    System.out.println("Added timing to method " +
                        argv[0] + "." + argv[1]);

                }

            } catch (CannotCompileException ex) {
                ex.printStackTrace();
            } catch (NotFoundException ex) {
                ex.printStackTrace();
            } catch (IOException ex) {
                ex.printStackTrace();
            }

        } else {
            System.out.println("Usage: JassistTiming class method-name");
        }
    }

    private static void addTiming(CtClass clas, String mname)
        throws NotFoundException, CannotCompileException {

        //  get the method information (throws exception if method with
        //  given name is not declared directly by this class, returns
        //  arbitrary choice if more than one with the given name)
        CtMethod mold = clas.getDeclaredMethod(mname);

        //  rename old method to synthetic name, then duplicate the
        //  method with original name for use as interceptor
        String nname = mname+"$impl";
        mold.setName(nname);
        CtMethod mnew = CtNewMethod.copy(mold, mname, clas, null);

        //  start the body text generation by saving the start time
        //  to a local variable, then call the timed method; the
        //  actual code generated needs to depend on whether the
        //  timed method returns a value
        String type = mold.getReturnType().getName();
        StringBuffer body = new StringBuffer();
        body.append("{\nlong start = System.currentTimeMillis();\n");
        if (!"void".equals(type)) {
            body.append(type + " result = ");
        }
        body.append(nname + "($$);\n");

        //  finish body text generation with call to print the timing
        //  information, and return saved value (if not void)
        body.append("System.out.println(\"Call to method " + mname +
            " took \" +\n (System.currentTimeMillis()-start) + " +
            "\" ms.\");\n");
        if (!"void".equals(type)) {
            body.append("return result;\n");
        }
        body.append("}");

        //  replace the body of the interceptor method with generated
        //  code block and add it to class
        mnew.setBody(body.toString());
        clas.addMethod(mnew);

        //  print the generated code block just to show what was done
        System.out.println("Interceptor method body:");
        System.out.println(body.toString());
    }
}

```

Show moreShow more icon

构造拦截器方法的正文时使用一个 `java.lang.StringBuffer` 来累积正文文本(这显示了处理 `String` 的构造的正确方法，与在 `StringBuilder` 的构造中使用的方法是相对的)。这种变化取决于原来的方法是否有返回值。如果它 _有_ 返回值，那么构造的代码就将这个值保存在局部变量中，这样在拦截器方法结束时就可以返回它。如果原来的方法类型为 `void` ，那么就什么也不需要保存，也不用在拦截器方法中返回任何内容。

除了对(重命名的)原来方法的调用，实际的正文内容看起来就像标准的 Java 代码。它是代码中的 `body.append(nname + "($$);\n")` 这一行，其中 `nname` 是原来方法修改后的名字。在调用中使用的 `$$` 标识符是 Javassist 表示正在构造的方法的一系列参数的方式。通过在对原来方法的调用中使用这个标识符，在调用拦截器方法时提供的参数就可以传递给原来的方法。

清单 5 展示了首先运行未修改过的 `StringBuilder` 程序、然后运行 `JassistTiming` 程序以添加计时信息、最后运行修改后的 `StringBuilder` 程序的结果。可以看到修改后的 `StringBuilder` 运行时会报告执行的时间，还可以看到因为字符串构造代码效率低下而导致的时间增加远远快于因为构造的字符串长度的增加而导致的时间增加。

##### 清单 5\. 运行这个程序

```
[dennis]$ java StringBuilder 1000 2000 4000 8000 16000
Constructed string of length 1000
Constructed string of length 2000
Constructed string of length 4000
Constructed string of length 8000
Constructed string of length 16000
[dennis]$ java -cp javassist.jar:. JassistTiming StringBuilder buildString
Interceptor method body:
{
long start = System.currentTimeMillis();
java.lang.String result = buildString$impl($$);
System.out.println("Call to method buildString took " +
(System.currentTimeMillis()-start) + " ms.");
return result;
}
Added timing to method StringBuilder.buildString
[dennis]$ java StringBuilder 1000 2000 4000 8000 16000
Call to method buildString took 37 ms.
Constructed string of length 1000
Call to method buildString took 59 ms.
Constructed string of length 2000
Call to method buildString took 181 ms.
Constructed string of length 4000
Call to method buildString took 863 ms.
Constructed string of length 8000
Call to method buildString took 4154 ms.
Constructed string of length 16000

```

Show moreShow more icon

## 可以信任源代码吗?

Javassist 通过让您处理源代码而不是实际的字节码指令清单而使 classworking 变得容易。但是这种方便性也有一个缺点。正如我在 [所有字节码的源代码](#所有字节码的源代码) 中提到的，Javassist 所使用的源代码与 Java 语言并不完全一样。除了在代码中识别特殊的标识符外，Javassist 还实现了比 Java 语言规范所要求的更宽松的编译时代码检查。因此，如果不小心，就会从源代码中生成可能会产生令人感到意外的结果的字节码。

作为一个例子，清单 6 展示了在将方法开始时的拦截器代码所使用的局部变量的类型从 `long` 变为 `int` 时的情况。Javassist 会接受这个源代码并将它转换为有效的字节码，但是得到的时间是毫无意义的。如果试着直接在 Java 程序中编译这个赋值，您就会得到一个编译错误，因为它违反了 Java 语言的一个规则：一个窄化的赋值需要一个类型覆盖。

##### 清单 6\. 将一个 `long` 储存到一个 `int` 中

```
[dennis]$ java -cp javassist.jar:. JassistTiming StringBuilder buildString
Interceptor method body:
{
int start = System.currentTimeMillis();
java.lang.String result = buildString$impl($$);
System.out.println("Call to method buildString took " +
(System.currentTimeMillis()-start) + " ms.");
return result;
}
Added timing to method StringBuilder.buildString
[dennis]$ java StringBuilder 1000 2000 4000 8000 16000
Call to method buildString took 1060856922184 ms.
Constructed string of length 1000
Call to method buildString took 1060856922172 ms.
Constructed string of length 2000
Call to method buildString took 1060856922382 ms.
Constructed string of length 4000
Call to method buildString took 1060856922809 ms.
Constructed string of length 8000
Call to method buildString took 1060856926253 ms.
Constructed string of length 16000

```

Show moreShow more icon

取决于源代码中的内容，甚至可以让 Javassist 生成无效的字节码。清单7展示了这样的一个例子，其中我将 `JassistTiming` 代码修改为总是认为计时的方法返回一个 `int` 值。Javassist 同样会毫无问题地接受这个源代码，但是在我试图执行所生成的字节码时，它不能通过验证。

##### 清单 7\. 将一个 `String` 储存到一个 `int` 中

```
[dennis]$ java -cp javassist.jar:. JassistTiming StringBuilder buildString
Interceptor method body:
{
long start = System.currentTimeMillis();
int result = buildString$impl($$);
System.out.println("Call to method buildString took " +
(System.currentTimeMillis()-start) + " ms.");
return result;
}
Added timing to method StringBuilder.buildString
[dennis]$ java StringBuilder 1000 2000 4000 8000 16000
Exception in thread "main" java.lang.VerifyError:
(class: StringBuilder, method: buildString signature:
(I)Ljava/lang/String;) Expecting to find integer on stack

```

Show moreShow more icon

只要对提供给 Javassist 的源代码加以小心，这就不算是个问题。不过，重要的是要认识到 Javassist 没有捕获代码中的所有错误，所以有可能会出现没有预见到的错误结果。

## 后续内容

Javassist 比我们在本文中所讨论的内容要丰富得多。下一个月，我们将进行更进一步的分析，看一看 Javassist 为批量修改类以及为在运行时装载类时对类进行动态修改而提供的一些特殊的功能。这些功能使 Javassist 成为应用程序中实现方面的一个很棒的工具，所以一定要继续跟随我们了解这个强大工具的全部内容。

本文翻译自： [Class transformation with Javassist](https://www.ibm.com/developerworks/library/j-dyn0916/index.html)（2003-09-16）