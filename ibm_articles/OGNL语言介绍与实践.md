# OGNL 语言介绍与实践
手把手教你 OGNL

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/os-cn-ognl/)

金 剑

发布: 2009-10-31

* * *

## OGNL 的历史

OGNL 最初是为了能够使用对象的属性名来建立 UI 组件 (component) 和 控制器 (controllers) 之间的联系，简单来说就是：视图 与 控制器 之间数据的联系。后来为了应付更加复杂的数据关系，Drew Davidson 发明了一个被他称为 KVCL(Key-Value Coding Language) 的语言。 Luke 参与进来后，用 ANTLR 来实现了该语言，并给它取了这个新名字，他后来又使用 JavaCC 重新实现了该语言。目前 OGNL 由 Drew 来负责维护。目前很多项目中都用到了 OGNL，其中不乏为大家所熟知的，例如几个流行的 web 应用框架：WebWork，Tapestry 等。

## 什么是 OGNL？

OGNL 是 Object-Graph Navigation Language 的缩写，从语言角度来说：它是一个功能强大的表达式语言，用来获取和设置 java 对象的属性 , 它旨在提供一个更高抽象度语法来对 java 对象图进行导航，OGNL 在许多的地方都有应用，例如：

1. 作为 GUI 元素（textfield,combobox, 等）到模型对象的绑定语言。
2. 数据库表到 Swing 的 TableModel 的数据源语言。
3. web 组件和后台 Model 对象的绑定语言 (WebOGNL,Tapestry,WebWork,WebObjects) 。
4. 作为 Jakarata Commons BeanUtils 或者 JSTL 的表达式语言的一个更具表达力的替代语言。

另外，java 中很多可以做的事情，也可以使用 OGNL 来完成，例如：列表映射和选择。 对于开发者来说，使用 OGNL，可以用简洁的语法来完成对 java 对象的导航。通常来说： 通过一个”路径”来完成对象信息的导航，这个”路径”可以是到 java bean 的某个属性，或者集合中的某个索引的对象，等等，而不是直接使用 get 或者 set 方法来完成。

## 为什么需要表达式语言 (EL)

表达式语言（EL）本质上被设计为：帮助你使用简单的表达式来完成一些”常用”的工作。通常情况下，ELs 可以在一些框架中找到，它被是用来简化我们的工作。例如：大家熟知的 Hibernate，使用 HQL(Hibernate Query Language) 来完成数据库的操作，HQL 成了开发人员与复查的 SQL 表达式之间的一个桥梁。 在 web 框架下，表达式语言起到了相似的目的。它的存在消除了重复代码的书写。例如：当没有 EL 的时候，为了从 session 中得到购物车并且将 ID 在网页上呈现出来，当直接在 jsp 中使用 java 代码来完成的时候，一般是：

```
<%
    ShoppingCart cart = (ShoppingCart) session.get("cart");
    int id = cart.getId();
    %>
    <%= id%>

```

Show moreShow more icon

你也可以将这些 code 压缩成一句，如下，但是现在代码就很不直观，且不可读。另外，虽然变成了一句，但是与上面的原始的例子一样，也包含了同样的表达式。例如：类型转换：转换成 ShoppingCart 。这里只不过是将原来的三个表达式变成了一句，其复杂度是没有得到简化的。

```
<%= ((ShoppingCart) session.get("cart")).getId() %>

```

Show moreShow more icon

当在 web 框架中使用表达式语言的时候，则可以有效的处理这种代码的复杂性。而不需要你，调用 servelet API，类型转换，然后再调用 getter 方法，多数的 Els 都可将这个过程简化为类似于：#session.cart.id 这中更可读的表达式。 表达式：#session.cart.id 与 java 代码不一样的是：没有 java 代码的 get 方法调用和类型转换。因为这些操作是非常”常用”的，这时候使用 EL 就顺理成章了，使用 EL 可以”消除”这些代码。

## OGNL 的基本语法

OGNL 表达式一般都很简单。虽然 OGNL 语言本身已经变得更加丰富了也更强大了，但是一般来说那些比较复杂的语言特性并未影响到 OGNL 的简洁：简单的部分还是依然那么简单。比如要获取一个对象的 name 属性，OGNL 表达式就是 name, 要获取一个对象的 headline 属性的 text 属性，OGNL 表达式就是 headline.text 。 OGNL 表达式的基本单位是”导航链”，往往简称为”链”。最简单的链包含如下部分：

表达式组成部分示例属性名称如上述示例中的 name 和 headline.text方法调用hashCode() 返回当前对象的哈希码。数组元素listeners[0] 返回当前对象的监听器列表中的第一个元素。

所有的 OGNL 表达式都基于当前对象的上下文来完成求值运算，链的前面部分的结果将作为后面求值的上下文。你的链可以写得很长，例如：

```
name.toCharArray()[0].numericValue.toString()

```

Show moreShow more icon

上面的表达式的求值步骤：

- 提取根 (root) 对象的 name 属性。
- 调用上一步返回的结果字符串的 toCharArray() 方法。
- 提取返回的结果数组的第一个字符。
- 获取字符的 numericValue 属性，该字符是一个 Character 对象，Character 类有一个 getNumericValue() 方法。
- 调用结果 Integer 对象的 toString() 方法。

上面的例子只是用来得到一个对象的值，OGNL 也可以用来去设置对象的值。当把上面的表达式传入 Ognl.setValue() 方法将导致 InappropriateExpressionException，因为链的最后的部分（ `toString()` ）既不是一个属性的名字也不是数组的某个元素。 了解了上面的语法基本上可以完成绝大部分工作了。

## 如何使用 OGNL

OGNL 不仅能够去获取或者设置对象的属性，而也可以用来：完成实例方法的调用，静态方法的调用，表达式求值，Lambda 表达式等，下面我们看看如何使用 OGNL 来完成这些任务。

### `ognl.Ognl` 类

最简单的使用是直接使用 ognl.Ognl 类来评估一个 OGNL 表达式。 Ognl 类提供一些静态方法用来解析和解释 OGNL 表达式，最简单的示例是不使用上下文从一个对象中获取某个表达式的值，示例如下：

```
import ognl.Ognl; import ognl.OgnlException;
try {
result = Ognl.getValue(expression, root);
}
catch (OgnlException ex)
{   // Report error or recover   }

```

Show moreShow more icon

上述代码将基于 root 对象评估 expression，返回结果，如果表达式有错，比如没有找到指定的属性，将抛出 OgnlException 。 更复杂一点的应用是使用预解析的表达式。这种方式允许在表达式求值之前就能捕获表达式的解析错误，应用开发人员可以缓存表达式解析出来的结果（AST），从而能在重复使用的时候提高性能。 Ognl 的 parseExpression 方法就是用来执行预解析操作的。 Ognl 2.7 版本后由于添加了” Expression Compilation ”性能得到了质的提高，在后面章节会有介绍。 Ognl 类的获取和设置方法也可以接受一个 context map 参数，他允许你放一些自己的变量，并能在 Ognl 表达式中使用。缺省的上下文里只包含 #root 和 #context 两个键。下面的示例展示如何从 root 对象中解析出 documentName 属性，然后将当前用户名称添加到返回的结果后面：

```
private Map context = new HashMap();
public void setUserName(String value)
{
      context.put("userName", value);
}
try {
     // get value using our own custom context map
     result = Ognl.getValue("userName"", context, root);
} catch (OgnlException ex) {
      // Report error or recover
}

```

Show moreShow more icon

上面提到的 #root 变量指向的就是当前的 root 变量（表达式求值的初始对象（initial object））, 而 #context 就是指向的 Map 对象，下面的例子可以更直观的说明

```
User root = new User();
     root.setId(19612);
     root.setName("sakura");
     Map context = new HashMap();
     context.put("who", "Who am i?");
     try {
       String who1 = (String)Ognl.getValue("#who", context, root);
       String who2 = (String)Ognl.getValue("#context.who", context, root);
         Object whoExp = Ognl.parseExpression("#who");
       String who3 = (String)Ognl.getValue(whoExp, context, root);
       //who1 who2 who3 返回同样的值， whoExp 重复使用可以提高效率
       String name1 = (String)Ognl.getValue("name", root);
       String name2 = (String)Ognl.getValue("#root.name", root);
       //name1 name2 返回同样的值
     } catch (OgnlException e) {
       //error handling
     }

```

Show moreShow more icon

### OGNL 表达式

1. 常量： 字符串：” ello” 字符：’h’ 数字：除了像 java 的内置类型 int,long,float 和 double,Ognl 还有如例：10.01B，相当于 java.math.BigDecimal，使用 ‘b’ 或者 ‘B’ 后缀。 100000H，相当于 java.math.BigInteger，使用 ‘h’ 或 ‘H’ 后缀。
2. 属性的引用 例如：user.name
3. 变量的引用 例如：#name
4. 静态变量的访问 使用 @class@field
5. 静态方法的调用 使用 @class@method(args), 如果没有指定 class 那么默认就使用 java.lang.Math.
6. 构造函数的调用 例如：new java.util.ArrayList();

其它的 Ognl 的表达式可以参考 Ognl 的语言手册。

### OGNL 的基本用法

OGNL 的 API，前面我们已经介绍过了，OGNL 的一个主要功能就是对象图的导航，我们看一下 OGNL 的最基本的用去取值和设置的 API。

##### getValue

```
public static java.lang.Object getValue(java.lang.String expression,
                                            java.util.Map context,
                                            java.lang.Object root,
                                            java.lang.Class resultType)
                                     throws OgnlException
     Evaluates the given OGNL expression to extract a value
      from the given root object in a given context
     Parameters:
     expression - the OGNL expression to be parsed
     context - the naming context for the evaluation
     root - the root object for the OGNL expression
     resultType - the converted type of the resultant object,
                 using the context's type converter
     Returns:
     the result of evaluating the expression

```

Show moreShow more icon

##### setValue

```
public static void setValue(java.lang.String expression,
                                 java.util.Map context,
                                 java.lang.Object root,
                                 java.lang.Object value)
                          throws OgnlException
     Evaluates the given OGNL expression to insert a value into the
     object graph rooted at the given root object given the context.
     Parameters:
     expression - the OGNL expression to be parsed
     root - the root object for the OGNL expression
     context - the naming context for the evaluation
     value - the value to insert into the object graph

```

Show moreShow more icon

OGNL 的 API 设计得是很简单的，context 提供上下文，为变量和表达式的求值过程来提供命名空间，存储变量 等，通过 root 来指定对象图遍历的初始变量，使用 expression 来告诉 Ognl 如何完成运算。看看下面两个简单的代码片段：

```
User user1 = new User();
user1.setId(1);
user1.setName("firer");
User user2 = new User();
user2.setId(2);
user2.setName("firer2");
List users = new ArrayList();
users.add(user1);
users.add(user2);
Department dep = new Department();
dep.setUsers(users);
dep.setName("dep");
dep.setId(11);
Object o = Ognl.getValue("users[1].name", dep);

```

Show moreShow more icon

这里我们可以看到前面介绍的使用表达式语言的有点，使用 “users[1].name”，就能完成对 name 的取值，而不用去进行类型转换等工作。下面是一个简单的设值的例子。

```
User user = new User();
user.setId(1);
user.setName("ffirer");
Ognl.setValue("department.name", user, "dep1");

```

Show moreShow more icon

就想前面介绍的，Ognl 也可以完成一些其它的工作，一个例子就是在我们的日常工作中，我们经常需要从列表中去”搜索”符合我们要求的对象，使用 java 的时候我们需要对列表进行遍历、类型转换、取值然后比较来得到我们想要的值，而使用 Ognl 将使这个过程变得简单，优雅。代码片段：

```
User user1 = new User();
user1.setId(1);
user1.setName("firer");
// 如上例创建一些 User
List users = new ArrayList();
users.add(user1);
// 将创建的 User 添加到 List 中
Department dep = new Department();
dep.setUsers(users);
List names = (List)Ognl.getValue("users.{name}", dep);
List ids = (List)Ognl.getValue("users.{? #this.id > 1}", dep);

```

Show moreShow more icon

这里表达式 “users.{name}” 将取得列表中所有 Users 的 name 属性，并以另外一个列表返回。 “users.{? #this.id > 1}” 将返回所有 id 大于 1 的 User，也以一个列表返回，包含了所有的 User 对象。这里使用的是 Ognl 的”列表投影”操作，叫这个名字是因为比较像数据库中返回某些列的操作。 Ognl 还有很多其它功能，在 Ognl 的 SVN 的 testcase 中找到如何使用它们。

## OGNL 的性能

OGNL，或者说表达式语言的性能主要又两方面来决定，一个就是对表达式的解析 (Parser)，另一个是表达式的执行，OGNL 采用 javaCC 来完成 parser 的实现，在 OGNL 2.7 中又对 OGNL 的执行部分进行了加强，使用 javasisit 来 JIT(Just-In-Time) 的生成 byte code 来完成表达式的执行。 Ognl 给这个功能的名字是：OGNL Expression Compilation 。 基本的使用方法是：

```
SimpleObject root = new SimpleObject();
OgnlContext context =  (OgnlContext) Ognl.createDefaultContext(null);

Node node =  (Node) Ognl.compileExpression(context, root, "user.name");
String userName = （String）node.getAccessor().get(context, root);

```

Show moreShow more icon

那么使用 Expression Compilation 能给 Ognl 的性能带来什么了？ 一个性能方面的简单测试：对比了 4 中情况：分别是：直接 Java 调用 (1), OGNL 缓存使用 expression compilation (2), OGNL 缓存使用 OGNL expression parsed（3） 以及不使用任何缓存的结果（4）。

##### 图 1\. 性能对比

![性能对比](../ibm_articles_img/os-cn-ognl_images_performanceCompare.jpg)

可以看到 expression compilation 非常接近使用 java 调用的时间，所以可以看到当表达式要被多次使用，使用 expression compilation 并且已经做好了缓存的情况下，OGNL 非常接近 java 直接调用的时间。

## 结束语

本文介绍了 OGNL 的概念、表达式语法以及如何使用 OGNL， 并提供了一些简单的示例代码，OGNL 在实际应用中还可以对其进行扩展，在本文中并为涉及， 感兴趣的读者可以进一步进行相关的学习和研究