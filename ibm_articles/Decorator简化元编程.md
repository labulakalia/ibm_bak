# Decorator 简化元编程
用于元编程的最新 Python 工具简介

**标签:** Linux,Python,系统

[原文链接](https://developer.ibm.com/zh/articles/l-cpdecor/)

David Mertz

发布: 2007-01-23

* * *

## 少劳多得

Decorator 与 Python 之前引入的元编程抽象有着某些共同之处：即使没有这些技术，您也一样可以实现它们所提供的功能。即使在 Python 1.5 中，也可以实现 Python 类的创建，而不需要使用 “元类” 挂钩。

Decorator 根本上的平庸与之非常类似。Decorator 所实现的功能就是修改紧接 Decorator 之后定义的函数和方法。这总是可能的，但这种功能主要是由 Python 2.2 中引入的 `classmethod()` 和 `staticmethod()` 内置函数驱动的。在旧式风格中，您可以调用 `classmethod()` ，如下所示：

##### 清单 1\. 典型的 “旧式” classmethod

```
class C:
    def foo(cls, y):
        print "classmethod", cls, y
    foo = classmethod(foo)

```

Show moreShow more icon

虽然 `classmethod()` 是内置函数，但并无独特之处；您也可以使用自己的方法转换函数。例如：

##### 清单 2\. 典型的 “旧式” 方法的转换

```
def enhanced(meth):
    def new(self, y):
        print "I am enhanced"
        return meth(self, y)
    return new
class C:
    def bar(self, x):
        print "some method says:", x
    bar = enhanced(bar)

```

Show moreShow more icon

decorator 所做的一切就是使您避免重复使用方法名，并且将 decorator 放在方法定义中第一处提及其名称的地方。例如：

##### 清单 3\. 典型的 “旧式” classmethod

```
class C:
    @classmethod
    def foo(cls, y):
        print "classmethod", cls, y
    @enhanced
    def bar(self, x):
        print "some method says:", x

```

Show moreShow more icon

decorator 也可以用于正则函数，采用的是与类中的方法相同的方式。令人惊奇的是，这一切是如此简单（严格来说，甚至有些不必要），只需要对语法进行简单修改，所有东西就可以工作得更好，并且使得程序的论证更加轻松。通过在方法定义的函数之前列出多个 decorator，即可将 decorator 链接在一起；良好的判断可以有助于防止将 _过多_ decorator 链接在一起，不过有时候将几个 decorator 链接在一起是有意义的：

##### 清单 4\. 链接 decorator

```
@synchronized
@logging
def myfunc(arg1, arg2, ...):
    # ...do something
# decorators are equivalent to ending with:
#    myfunc = synchronized(logging(myfunc))
# Nested in that declaration order

```

Show moreShow more icon

Decorator 只是一个语法糖，如果您过于急切，那么它就会使您搬起石头砸了自己的脚。decorator 其实就是一个至少具有一个参数的函数 —— 程序员要负责确保 decorator 的返回内容仍然是一个有意义的函数或方法，并且实现了原函数为使连接有用而做的事情。例如，下面就是 decorator 两个不正确的用法：

##### 清单 5\. 没有返回函数的错误 decorator

```
>>> def spamdef(fn):
...     print "spam, spam, spam"
...
>>> @spamdef
... def useful(a, b):
...     print a**2 + b**2
...
spam, spam, spam
>>> useful(3, 4)
Traceback (most recent call last):
File "<stdin>", line 1, in ?
TypeError: 'NoneType' object is not callable

```

Show moreShow more icon

decorator 可能会返回一个函数，但这个函数与未修饰的函数之间不存在有意义的关联：

##### 清单 6\. 忽略传入函数的 decorator

```
>>> def spamrun(fn):
...     def sayspam(*args):
...         print "spam, spam, spam"
...     return sayspam
...
>>> @spamrun
... def useful(a, b):
...     print a**2 + b**2
...
>>> useful(3,4)
spam, spam, spam

```

Show moreShow more icon

最后，一个表现更良好的 decorator 可以在某些方面增强或修改未修饰函数的操作：

##### 清单 7\. 修改未修饰函数行为的 decorator

```
>>> def addspam(fn):
...     def new(*args):
...         print "spam, spam, spam"
...         return fn(*args)
...     return new
...
>>> @addspam
... def useful(a, b):
...     print a**2 + b**2
...
>>> useful(3,4)
spam, spam, spam
25

```

Show moreShow more icon

您可能会质疑， `useful()` 到底有多么有用？ `addspam()` 真的是那样出色的 _增强_ 吗？但这种机制至少符合您通常能在有用的 decorator 中看到的那种模式。

## 高级抽象简介

根据我的经验，元类应用最多的场合就是在类实例化之后对类中的方法进行修改。decorator 目前并不允许您修改类实例化 _本身_ ，但是它们可以修改依附于类的方法。这并不能让您在实例化过程中动态添加或删除方法或类属性，但是它让这些方法可以在运行时根据环境的条件来变更其行为。现在从技术上来说，decorator 是在运行 `class` 语句时应用的，对于顶级类来说，它更接近于 “编译时” 而非 “运行时”。但是安排 decorator 的运行时决策与创建类工厂一样简单。例如：

##### 清单 8\. 健壮但却深度嵌套的 decorator

```
def arg_sayer(what):
    def what_sayer(meth):
        def new(self, *args, **kws):
            print what
            return meth(self, *args, **kws)
        return new
    return what_sayer

def FooMaker(word):
    class Foo(object):
        @arg_sayer(word)
        def say(self): pass
    return Foo()

foo1 = FooMaker('this')
foo2 = FooMaker('that')
print type(foo1),; foo1.say()  # prints: <class '__main__.Foo'> this
print type(foo2),; foo2.say()  # prints: <class '__main__.Foo'> that

```

Show moreShow more icon

`@arg_sayer()` 绕了很多弯路，但只获得非常有限的结果，不过对于它所阐明的几方面来说，这是值得的：

- `Foo.say()` 方法对于不同的实例有不同的行为。在这个例子中，不同之处只是一个数据值，可以轻松地通过其他方式改变这个值；不过原则上来说，decorator 可以根据运行时的决策来彻底重写这个方法。
- 本例中未修饰的 `Foo.say()` 方法是一个简单的占位符，其整个行为都是由 decorator 决定的。然而，在其他情况下，decorator 可能会将未修饰的方法与一些新功能相 _结合_ 。
- 正如我们已经看到的一样， `Foo.say()` 的修改是通过 `FooMaker()` 类工厂在运行时严格确定的。可能更加典型的情况是在顶级定义类中使用 decorator，这些类只依赖于编译时可用的条件（这通常就足够了）。
- decorator 都是参数化的。或者更确切地说， `arg_sayer()` 本身根本就不是一个真正的 decorator； `arg_sayer()`所返回的\_ 函数 —— `what_sayer()` 就是一个使用了闭包来封装其数据的 decorator 函数。参数化的 decorator 较为常见，但是它们将所需的函数嵌套为三层。

## 迈进元类领域

正如上一节中介绍的一样，decorator 并不能完全取代元类挂钩，因为它们只修改了方法，而未添加或删除方法。实际上，这样说并不完全正确。作为一个 Python 函数，decorator 完全可以实现其他 Python 代码所实现的任何功能。通过修饰一个类的 `.__new__()` 方法（甚至是其占位符版本），您实际上可以更改附加到该类的方法。尽管尚未在现实中看到这种模式，不过我认为它有着某种必然性，甚至可以作为 `_metaclass_` 指派的一项改进：

##### 清单 9\. 添加和删除方法的 decorator

```
def flaz(self): return 'flaz'     # Silly utility method
def flam(self): return 'flam'     # Another silly method

def change_methods(new):
    "Warning: Only decorate the __new__() method with this decorator"
    if new.__name__!= '__new__':
        return new  # Return an unchanged method
    def __new__(cls, *args, **kws):
        cls.flaz = flaz
        cls.flam = flam
        if hasattr(cls, 'say'): del cls.say
        return super(cls.__class__, cls).__new__(cls, *args, **kws)
    return __new__

class Foo(object):
    @change_methods
    def __new__(): pass
    def say(self): print "Hi me:", self

foo = Foo()
print foo.flaz()  # prints: flaz
foo.say()         # AttributeError: 'Foo' object has no attribute 'say'

```

Show moreShow more icon

在 `change_methods()` decorator 示例中，我们添加并删除了几个固定的方法，不过这是毫无意义的。在更现实的情况中，应使用上一节中提到的几个模式。例如，参数化的 decorator 可以接受一个能表示要添加或删除的方法的数据结构；或者由数据库查询之类的某些环境特性做出这一决策。这种对附加方法的操作也可以像之前一样打包到一个函数工厂中，这将使最终决策延迟到运行时。这些新兴技术也许比 `_metaclass_` 指派更加万能。例如，您可以调用一个增强了的 `change_methods()` ，如下所示：

##### 清单 10\. 增强的 change\_methods()

```
class Foo(object):
    @change_methods(add=(foo, bar, baz), remove=(fliz, flam))
    def __new__(): pass

```

Show moreShow more icon

## 修改调用模型

您将看到，有关 decorator 的最典型的例子可能是使一个函数或方法来实现 “其他功能”，同时完成其基本工作。例如，在诸如 Python Cookbook Web 站点之类的地方，您可以看到 decorator 添加了诸如跟踪、日志记录、存储/缓存、线程锁定以及输出重定向之类的功能。与这些修改相关（但实质略有区别）的是修饰 “之前” 和 “之后”。对于修饰之前/之后来说，一种有趣的可能性就是检查传递给函数的参数和函数返回值的类型。如果这些类型并非如我们预期的一样，那么这种 `type_check()` decorator 就可能会触发一个异常，或者采取一些纠正操作。

与这种 decorator 前/后类似的情况，我想到了 R 编程语言和 `NumPy` 特有的函数的 “elementwise” 应用。在这些语言中，数学函数通常应用于元素序列中的 _每个元素_ ，但也会应用于单个数字。

当然， `map()` 函数、列表内涵（list-comprehension）和最近的生成器内涵（generator-comprehension 都可以让您实现 elementwise 应用。但是这需要较小的工作区来获得类似于 R 语言的行为： `map()` 所返回的序列类型通常是一个列表；如果您传递的是单个元素而不是一个序列，那么调用将失败。例如：

##### 清单 11\. map() 调用失败

```
>>> from math import sqrt
>>> map(sqrt, (4, 16, 25))
[2.0, 4.0, 5.0]
>>> map(sqrt, 144)
TypeError: argument 2 to map() must support iteration

```

Show moreShow more icon

创建一个可以 “增强” 普通数值函数的 decorator 并不困难：

##### 清单 12\. 将函数转换成 elementwise 函数

```
def elementwise(fn):
    def newfn(arg):
        if hasattr(arg,'__getitem__'):  # is a Sequence
            return type(arg)(map(fn, arg))
        else:
            return fn(arg)
    return newfn

@elementwise
def compute(x):
    return x**3 - 1

print compute(5)        # prints: 124
print compute([1,2,3])  # prints: [0, 7, 26]
print compute((1,2,3))  # prints: (0, 7, 26)

```

Show moreShow more icon

当然，简单地编写一个具有不同返回类型的 `compute()` 函数并不困难；毕竟 decorator 只需占据几行。但是作为对面向方面编程的一种认可，这个例子让我们可以 _分离_ 那些在不同层次上运作的关注事项。我们可以编写各种数值计算函数，希望它们都可转换成 elementwise 调用模型，而不用考虑参数类型测试和返回值类型强制转换的细节。

对于那些对单个事物或事物序列（此时要保留序列类型）进行操作的函数来说， `elementwise()` decorator 均可同样出色地发挥作用。作为一个练习，您可尝试去解决如何允许相同的修饰后调用来接受和返回迭代器（提示：如果您只是想迭代一次完整的 elementwise 计算，那么当且仅当传入的是一个迭代对象时，才能这样简化一些。）

您将碰到的大多数优秀的 decorator 都在很大程度上采用了这种组合正交关注的范例。传统的面向对象编程，尤其是在诸如 Python 之类允许多重继承的语言中，都会试图使用一个继承层次结构来模块化关注事项。然而，这仅会从一个祖先那里获取一些方法，而从其他祖先那里获取其他方法，因此需要采用一种概念，使关注事项比在面向方面的思想中更加分散。要充分利用生成器，就要考虑一些与混搭方法不同的问题：可以处于方法本身的 “核心” 之外的关注事项为依据，使 _各_ 方法以不同方式工作。

## 修饰 decorator

在结束本文之前，我想为您介绍一种确实非常出色的 Python 模块，名为 `decorator` ，它是由与我合著过一些图书的 Michele Simionato 编写的。该模块使 decorator 的开发变得更加美妙。 `decorator` 模块的主要组件具有某种自反式的优雅，它是一个称为 `decorator()` 的 decorator。与未修饰的函数相比，使用 `@decorator` 修饰过的函数可以通过一种更简单的方式编写。

Michele 已经为自己的模块编写了很好的文档，因此这里不再赘述；不过我非常乐意介绍一下它所解决的基本问题。 `decorator` 模块有两大主要优势。一方面，它使您可以编写出嵌套层次更少的 decorator，如果没有这个模块，您就只能使用更多层次（”平面优于嵌套”）；但更加有趣的是这样一个事实：它使得修饰过的函数可以真正地与其在元数据中未修饰的版本相匹配，这是我的例子中没有做到的。例如，回想一下我们上面使用过的简单 “跟踪” decorator `addspam()` ：

##### 清单 13\. 一个简单的 decorator 是如何造成元数据崩溃的

```
>>> def useful(a, b): return a**2 + b**2
>>> useful.__name__
'useful'
>>> from inspect import getargspec
>>> getargspec(useful)
(['a', 'b'], None, None, None)
>>> @addspam
... def useful(a, b): return a**2 + b**2
>>> useful.__name__
'new'
>>> getargspec(useful)
([], 'args', None, None)

```

Show moreShow more icon

尽管这个修饰过的函数 _的确完成_ 了自己增强过的工作，但若进一步了解，就会发现这并不是完全正确的，尤其是对于那些关心这种细节的代码分析工具或 IDE 来说更是如此。使用 `decorator` ，我们就可以改进这些问题：

##### 清单 14\. decorator 更聪明的用法

```
>>> from decorator import decorator
>>> @decorator
... def addspam(f, *args, **kws):
...     print "spam, spam, spam"
...     return f(*args, **kws)
>>> @addspam
... def useful(a, b): return a**2 + b**2
>>> useful.__name__
'useful'
>>> getargspec(useful)
(['a', 'b'], None, None, None)

```

Show moreShow more icon

这对于编写 decorator 更加有利，同时，其保留行为的元数据的也更出色了。当然，阅读 Michele 开发这个模块所使用的全部资料会使您回到大脑混沌的世界，我们将这留给 Simionato 博士一样的宇宙学家好了。

本文翻译自： [Decorators make magic easy](https://developer.ibm.com/articles/l-cpdecor/)（2007-01-23）