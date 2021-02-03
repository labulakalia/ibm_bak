# PowerMock 简介
使用 PowerMock 以及 Mockito 实现单元测试

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-lo-powermock/)

张羽, 吴长侠

发布: 2012-04-16

* * *

EasyMock 以及 Mockito 都因为可以极大地简化单元测试的书写过程而被许多人应用在自己的工作中，但是这 2 种 Mock 工具都不可以实现对静态函数、构造函数、私有函数、Final 函数以及系统函数的模拟，但是这些方法往往是我们在大型系统中需要的功能。PowerMock 是在 EasyMock 以及 Mockito 基础上的扩展，通过定制类加载器等技术，PowerMock 实现了之前提到的所有模拟功能，使其成为大型系统上单元测试中的必备工具。

## 单元测试模拟框架的功能及其实现简介

单元测试在软件开发过程中的重要性不言而喻，特别是在测试驱动开发的开发模式越来越流行的前提下，单元测试更成为了软件开发过程中不可或缺的部分。于是相应的，各种单元测试技术也应运而生。本文要介绍的 PowerMock 以及 Mockito 都是简化单元测试书写过程的工具。

Mockito 是一个针对 Java 的单元测试模拟框架，它与 EasyMock 和 jMock 很相似，都是为了简化单元测试过程中测试上下文 ( 或者称之为测试驱动函数以及桩函数 ) 的搭建而开发的工具。在有这些模拟框架之前，为了编写某一个函数的单元测试，程序员必须进行十分繁琐的初始化工作，以保证被测试函数中使用到的环境变量以及其他模块的接口能返回预期的值，有些时候为了单元测试的可行性，甚至需要牺牲被测代码本身的结构。单元测试模拟框架则极大的简化了单元测试的编写过程：在被测试代码需要调用某些接口的时候，直接模拟一个假的接口，并任意指定该接口的行为。这样就可以大大的提高单元测试的效率以及单元测试代码的可读性。

相对于 EasyMock 和 jMock，Mockito 的优点是通过在执行后校验哪些函数已经被调用，消除了对期望行为（expectations）的需要。其它的 mocking 库需要在执行前记录期望行为（expectations），而这导致了丑陋的初始化代码。

但是，Mockito 也并不是完美的，它不提供对静态方法、构造方法、私有方法以及 Final 方法的模拟支持。而程序员时常都会发现自己有对以上这些方法的模拟需求，特别是当一个已有的软件系统摆在面前时。幸好 , 还有 PowerMock。

PowerMock 也是一个单元测试模拟框架，它是在其它单元测试模拟框架的基础上做出的扩展。通过提供定制的类加载器以及一些字节码篡改技巧的应用，PowerMock 现了对静态方法、构造方法、私有方法以及 Final 方法的模拟支持，对静态初始化过程的移除等强大的功能。因为 PowerMock 在扩展功能时完全采用和被扩展的框架相同的 API, 熟悉 PowerMock 所支持的模拟框架的开发者会发现 PowerMock 非常容易上手。PowerMock 的目的就是在当前已经被大家所熟悉的接口上通过添加极少的方法和注释来实现额外的功能，目前，PowerMock 仅支持 EasyMock 和 Mockito。

本文的目的就是和大家一起学习在 Mockito 框架上扩展的 PowerMock 的强大功能。

## 环境配置方法

对于需要的开发包，PowerMock 网站提供了”一站式”下载 : 从 [此页面](http://code.google.com/p/powermock/downloads/list) 中选择以类似 PowerMock 1.4.10 with Mockito and JUnit including dependencies 为注释的链接，该包中包含了最新的 JUnit 库，Mockito 库，PowerMock 库以及相关的依赖。

如果是使用 Eclipse 开发，只需要在 Eclipse 工程中包含这些库文件即可。

如果是使用 Maven 开发，则需要根据版本添加以下清单内容到 POM 文件中：

JUnit 版本 4.4 以上请参考清单 1，

##### 清单 1

```
<properties>
    <powermock.version>1.4.10</powermock.version>
</properties>
<dependencies>
<dependency>
      <groupId>org.powermock</groupId>
      <artifactId>powermock-module-junit4</artifactId>
      <version>${powermock.version}</version>
      <scope>test</scope>
</dependency>
<dependency>
      <groupId>org.powermock</groupId>
      <artifactId>powermock-api-mockito</artifactId>
      <version>${powermock.version}</version>
      <scope>test</scope>
</dependency>
</dependencies>

```

Show moreShow more icon

JUnit 版本 4.0-4.3 请参考清单 2，

##### 清单 2

```
<properties>
    <powermock.version>1.4.10</powermock.version>
</properties>
<dependencies>
<dependency>
      <groupId>org.powermock</groupId>
      <artifactId>powermock-module-junit4-legacy</artifactId>
      <version>${powermock.version}</version>
      <scope>test</scope>
</dependency>
<dependency>
      <groupId>org.powermock</groupId>
      <artifactId>powermock-api-mockito</artifactId>
      <version>${powermock.version}</version>
      <scope>test</scope>
</dependency>
</dependencies>

```

Show moreShow more icon

JUnit 版本 3 请参考清单 3，

##### 清单 3

```
<properties>
    <powermock.version>1.4.10</powermock.version>
</properties>
<dependencies>
<dependency>
      <groupId>org.powermock</groupId>
      <artifactId>powermock-module-junit3</artifactId>
      <version>${powermock.version}</version>
      <scope>test</scope>
</dependency>
<dependency>
      <groupId>org.powermock</groupId>
      <artifactId>powermock-api-mockito</artifactId>
      <version>${powermock.version}</version>
      <scope>test</scope>
</dependency>
</dependencies>

```

Show moreShow more icon

## PowerMock 在单元测试中的应用

### 模拟 Static 方法

在任何需要用到 PowerMock 的类开始之前，首先我们要做如下声明：

@RunWith(PowerMockRunner.class)

然后，还需要用注释的形式将需要测试的静态方法提供给 PowerMock：

@PrepareForTest( { YourClassWithEgStaticMethod.class })

然后就可以开始写测试代码：

首先，需要有一个含有 static 方法的代码 , 如清单 4：

##### 清单 4

```
public class IdGenerator {

    ...

    public static long generateNewId() {
        ...
    }

    ...
}

```

Show moreShow more icon

然后，在被测代码中，引用了以上方法 , 如清单 5 所示：

##### 清单 5

```
public class ClassUnderTest {
    ...
    public void methodToTest() {
        ..
        final long id = IdGenerator.generateNewId();
        ..
     }
    ...
}

```

Show moreShow more icon

为了达到单元测试的目的，需要让静态方法 `generateNewId()` 返回各种值来达到对被测试方法 `methodToTest()` 的覆盖测试，实现方式如清单 6 所示：

##### 清单 6

```
@RunWith(PowerMockRunner.class)
//We prepare the IdGenerator for test because the static method is normally not mockable
@PrepareForTest(IdGenerator.class)
public class MyTestClass {
    @Test
    public void demoStaticMethodMocking() throws Exception {
        mockStatic(IdGenerator.class);
        /*
         * Setup the expectation using the standard Mockito syntax,
         * generateNewId() will now return 2 everytime it's invoked
         * in this test.
         */
        when(IdGenerator.generateNewId()).thenReturn(2L);

        new ClassUnderTest().methodToTest();

        // Optionally verify that the static method was actually called
        verifyStatic();
        IdGenerator.generateNewId();
    }
}

```

Show moreShow more icon

如清单 6 中所展示，在测试代码中，可以使用 `When().thenReturn(`) 语句来指定被引用的静态方法返回任意需要的值，达到覆盖测试的效果。

### 模拟构造函数

有时候，能模拟构造函数，从而使被测代码中 `new` 操作返回的对象可以被随意定制，会很大程度的提高单元测试的效率，考虑如清单 7 的代码：

##### 清单 7

```
public class DirectoryStructure {
    public boolean create(String directoryPath) {
        File directory = new File(directoryPath);

        if (directory.exists()) {
            throw new IllegalArgumentException(
            "\"" + directoryPath + "\" already exists.");
        }

        return directory.mkdirs();
    }
}

```

Show moreShow more icon

为了充分测试 `create()` 函数，我们需要被 `new` 出来的 File 对象返回文件存在和不存在两种结果。在 PowerMock 出现之前，实现这个单元测试的方式通常都会需要在实际的文件系统中去创建对应的路径以及文件。然而，在 PowerMock 的帮助下，本函数的测试可以和实际的文件系统彻底独立开来：使用 PowerMock 来模拟 File 类的构造函数，使其返回指定的模拟 File 对象而不是实际的 File 对象，然后只需要通过修改指定的模拟 File 对象的实现，即可实现对被测试代码的覆盖测试，参考如清单 8 的代码：

##### 清单 8

```
@RunWith(PowerMockRunner.class)
@PrepareForTest(DirectoryStructure.class)
public class DirectoryStructureTest {
    @Test
    public void createDirectoryStructureWhenPathDoesntExist() throws Exception {
        final String directoryPath = "mocked path";

        File directoryMock = mock(File.class);

        // This is how you tell PowerMockito to mock construction of a new File.
        whenNew(File.class).withArguments(directoryPath).thenReturn(directoryMock);

        // Standard expectations
        when(directoryMock.exists()).thenReturn(false);
        when(directoryMock.mkdirs()).thenReturn(true);

        assertTrue(new NewFileExample().createDirectoryStructure(directoryPath));

        // Optionally verify that a new File was "created".
        verifyNew(File.class).withArguments(directoryPath);
    }
}

```

Show moreShow more icon

使用 `whenNew().withArguments().thenReturn()` 语句即可实现对具体类的构造函数的模拟操作。然后对于之前创建的模拟对象 `directoryMock` 使用 `When().thenReturn()` 语句，即可实现需要的所有功能，从而实现对被测对象的覆盖测试。在本测试中，因为实际的模拟操作是在类 `DirectoryStructureTest` 中实现，所以需要指定的 `@PrepareForTest` 对象是 `DirectoryStructureTest.class` 。

### 模拟私有以及 Final 方法

为了实现对类的私有方法或者是 Final 方法的模拟操作，需要 PowerMock 提供的另外一项技术：局部模拟。

在之前的介绍的模拟操作中，我们总是去模拟一整个类或者对象，然后使用 `When().thenReturn()` 语句去指定其中值得关心的部分函数的返回值，从而达到搭建各种测试环境的目标。对于没有使用 `When().thenReturn()` 方法指定的函数，系统会返回各种类型的默认值（具体值可参考官方文档）。

局部模拟则提供了另外一种方式，在使用局部模拟时，被创建出来的模拟对象依然是原系统对象，虽然可以使用方法 `When().thenReturn()` 来指定某些具体方法的返回值，但是没有被用此函数修改过的函数依然按照系统原始类的方式来执行。

这种局部模拟的方式的强大之处在于，除开一般方法可以使用之外，Final 方法和私有方法一样可以使用。

参考如清单 9 所示的被测代码：

##### 清单 9

```
public final class PrivatePartialMockingExample {
    public String methodToTest() {
        return methodToMock("input");
    }

    private String methodToMock(String input) {
        return "REAL VALUE = " + input;
    }
}

```

Show moreShow more icon

为了保持单元测试的纯洁性，在测试方法 `methodToTest()` 时，我们不希望受到私有函数 `methodToMock()` 实现的干扰，为了达到这个目的，我们使用刚提到的局部模拟方法来实现 , 实现方式如清单 10：

##### 清单 10

```
@RunWith(PowerMockRunner.class)
@PrepareForTest(PrivatePartialMockingExample.class)
public class PrivatePartialMockingExampleTest {
    @Test
    public void demoPrivateMethodMocking() throws Exception {
        final String expected = "TEST VALUE";
        final String nameOfMethodToMock = "methodToMock";
        final String input = "input";

        PrivatePartialMockingExample underTest = spy(new PrivatePartialMockingExample());

        /*
         * Setup the expectation to the private method using the method name
         */
        when(underTest, nameOfMethodToMock, input).thenReturn(expected);

        assertEquals(expected, underTest.methodToTest());

        // Optionally verify that the private method was actually called
        verifyPrivate(underTest).invoke(nameOfMethodToMock, input);
    }
}

```

Show moreShow more icon

可以发现，为了实现局部模拟操作，用来创建模拟对象的函数从 `mock()` 变成了 `spy()` ，操作对象也从类本身变成了一个具体的对象。同时， `When()` 函数也使用了不同的版本：在模拟私有方法或者是 Final 方法时， `When()` 函数需要依次指定模拟对象、被指定的函数名字以及针对该函数的输入参数列表。

## 结束语

以上列举了扩展于 Mockito 版本的 PowerMock 的一部分强大的功能，特别是针对已有的软件系统，利用以上功能可以轻易的完成清晰独立的单元测试代码，帮助我们提高代码质量。

注：本文中的部分测试代码引用自 Johan Haleby 的 [Untestable code with Mockito and PowerMock](http://blog.jayway.com/2009/10/28/untestable-code-with-mockito-and-powermock/) 。