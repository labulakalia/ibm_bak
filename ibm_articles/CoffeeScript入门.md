# CoffeeScript 入门
初步了解 CoffeeScript 以及它为开发人员带来的特有优势

**标签:** JavaScript,Node.js,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/wa-coffee1/)

Michael Galpin

发布: 2012-08-27

* * *

## 简介

CoffeeScript 编程语言是构建于 JavaScript 之上，它可编译成高效 JavaScript，您可以在 Web 浏览器上运行它或者将其与诸如 Node.js 一类的技术相结合用于构建服务器端应用程序。编译过程通常都很简单，生成的 JavaScript 均与许多最佳实践保持一致。在本文中，我们将了解有关于 CoffeeScript 编程语言的特性。在安装并运行 CoffeeScript 编译器后，将向您展示一个在 Web 页面中使用 CoffeeScript 的简单示例。

[下载](http://public.dhe.ibm.com/software/dw/web/wa-coffee1/cs1.zip) 本文所使用的源代码。

## CoffeeScript 的吸引力

##### 常用的缩写

- NPM：Node Package Manager
- REPL：Read-Evaluate-Print-Loop

现在人们动辄就会说，JavaScript 是最重要的编程语言。这是一门浏览器语言，正越来越频繁地出现在桌面和移动应用程序中。随着 Node.js 的日益普及，JavaScript 已经成为了服务器和系统应用程序一种切实可行的选择。一些开发人员强烈抵制 JavaScript，大部分原因是因为其不一致的语法和古怪的实现。不过，随着 JavaScript 虚拟机变得越来越标准化，古怪的实现已经逐渐减少。不一致的语法则可能会随着 JavaScript 的下一次革命：ECMAScript.next 的到来而得到一定的解决。ECMAScript.next 是一个新兴标准，深受 CoffeeScript 的影响。不过，在新标准被接受并且被流行虚拟机实现之前，JavaScript 语法还需要进一步改进。

如果您正在等待 JavaScript 运行时，CoffeeScript 会是一个有吸引力的选择。从语法的角度来看，JavaScript 是不择不扣的大杂烩，它具有许多函数式编程语言的功能，特别是深受 Scheme 的影响。然而，Scheme 是构建于 S-表达式之上的一种非常简单的语法。JavaScript 沿袭了 Scheme 中的许多概念，但是却没有采用它的语法。相反，JavaScript 拥有类似 C 的语法。结果是，一种有着函数式概念的语言，但却有着冗长的语法，而没有用于表达这些概念的自然构造，例如，JavaScript 允许高阶函数，比如，那些其输入参数包含其他函数的函数。这既有用且功能强大，是许多语言都缺少的功能。不过，JavaScript 的语法并不那么优雅，如清单 1 所示。

##### 清单 1\. 丑陋的 JavaScript

```
pmb.requestPaymentInfo('type', function(info){
    $('result').innerHTML = info.name;
});

```

Show moreShow more icon

示例中含有许多样板代码，如括号、逗号、大括号、分号，以及一些并不是真正需要的语言关键字。

JavaScript 的主要用途是作为一种客户端 Web 应用程序语言。诸如 Cocoa、Windows® 、Forms 和 Android 这类桌面和移动应用程序框架均是面向对象的。面向对象范型并非完美的，但非常适用于使用图形用户界面的应用程序。JavaScript 也是一种拥有继承性的面向对象语言，但它是原型继承，并非像被大多数应用程序框架所采用的基于类的语言。因此，使用 JavaScript 进行应用程序编程可能非常繁琐。

CoffeeScript 解决了 JavaScript 这些痛点。CoffeeScript：

- 提供一种比较简单的语法，减少了样板代码，诸如括号和逗号
- 使用空格作为一种组织代码块的方法
- 提供拥有表达函数的简单语法
- 提供基于类的继承（可选项，但是在进行应用程序开发时非常有用）

您可能会想，相较于 JavaScript，CoffeeScript 肯定有一些劣势，因为它的语法比较抽象。例如，CoffeeScript 会不会比 JavaScript 更慢？会不会需要一个较大的运行时库？实际上，CoffeeScript 将编译成简洁、高效的 JavaScript，您总能清楚地看到正在编译的内容，因此，您可以很自信没有引入过多的东西。此外，因为 CoffeeScript 将完全编译成函数式 JavaScript，所以不需要任何类型的运行时库。CoffeeScript 所提供的语法允许您充分利用 JavaScript 的强大功能，而只需要较小的运行时开销。

## 先决条件

如上所述，您可以使用 CoffeeScript 来编写运行在 Node.js 上的服务器和系统应用程序，不过 CoffeeScript 和 Node.js 之间就要建立起更深一层的关系。若要安装 CoffeeScript，需要先安装 Node.js，因为：

- CoffeeScript 使用节点的包管理程序 NPM，作为 Node.js 的一个包进行分布。
- CoffeScript 必须被编译，其编译器实际上就是使用 CoffeeScript 编写的，所以，需要一个 JavaScript 运行时来完成其编译。作为 Node.js 核心的 V8 JavaScript 虚拟机很适合这一任务。

若要按照本文中的示例进行实践，需要先安装 Node.js。

## 安装

您是否曾希望可以从命令行运行 JavaScript？反正我是没有，不过 CoffeeScript 可能会对此有所改变。有了 Node.js，您就可以从命令行运行 JavaScript，或者或将其作为可执行脚本的一部分。Node.js 的这一主要功能允许在命令行执行 CoffeeScript 代码，提供了 CoffeeScript 编译器（使用 CoffeeScript 编写的）所需要的运行时。

第一步就是安装 Node.js。有几种安装选择；您可以编译源代码或者运行可用于各种系统的安装程序。从命令行运行 node -v 来确认 Node.js 是否已安装，以及是否位于您的路径中。

Node.js 为您带来了额外收益：节点包管理器 (NPM)。从命令行运行 node -v 来确认 Node.js 已安装且已在您的路径中之后，就可以按照如下步骤使用 NPM 安装 CoffeeScript 了。

1. 从命令行运行 `npm install --global coffee-script` 。

    `--global` 标记使 CoffeeScript 可广泛用于系统，而不仅仅用于特定项目。

2. `npm` 命令输出类似 /usr/bin/coffee -> /usr/lib/node\_modules/coffee-script/bin/coffee 的信息。

    NPM 在 /usr/bin 中创建了一个快捷键，因此，现在可执行 Coffee 位于正确的路径中，这就是 CoffeeScript 编译器和解释器。

3. 若要验证可执行 coffee 是否位于路径中，在命令行运行 `coffee -v` 。

还剩最后一步，确认 CoffeeScript 环境是否正确设置。为了使 CoffeeScript 可用于任何已启动的 Node.js 进程，您需要将其添加到 Node.js 的 NODE\_PATH 中。在遇到不能识别的函数时，Node.js 将搜索 NODE\_PATH 获取模块（库）。

本文示例，您主要可以将 Node.js 作为可执行 CoffeeScript 运行时来使用。最简单的方法是仅将所有 NPM 模块添加到 NODE\_PATH。要找出 NPM 模块的位置，输入 `npm ls -g` 。您需要添加一个环境变量将 NODE\_PATH 指向其位置。例如，如果 npm ls -g 输出 /usr/lib，那么模块则位于 /usr/lib/node\_modules。要设置一个 NODE\_PATH 环境变量，运行 export NODE\_PATH=/usr/lib/node\_modules。

您可以通过将上述命令放入启动脚本（比如说 ~/.bash\_profile）中来流水线化这些操作。若要验证更改，通过执行 Node 启动一个 Node.js shell，然后输入 `require('coffee-script')` 。Node.js shell 会加载 CoffeeScript 库。如果可以正常运行，那么 CoffeeScript 环境已经配置好了。现在您可以从编译器开始探索 CoffeeScript 之旅了。

## 编译器

运行 CoffeeScript 编译器很简单，输入 `coffee -c` 即可。这启动了 CoffeeScript 的读取-计算-输出-循环 (read-evaluate-print-loop, [REPL](#repl))。若要执行该编译器，需要将您准备编译 CoffeScript 文件给传递它。创建一个名为 cup0.coffee 的文件，将清单 2 中的内容粘贴到其中。

##### 清单 2\. Cup 0

```
for i in [0..5]
    console.log "Hello #{i}"

```

Show moreShow more icon

您可能会想，清单 2 中的这两行代码能干什么，清单 3 展示了运行 `coffee cup0.coffee` 的输出。

##### 清单 3\. 运行您的第一个 CoffeeScript

```
$ coffee cup0.coffee
Hello 0
Hello 1
Hello 2
Hello 3
Hello 4
Hello 5

```

Show moreShow more icon

为了更好地了解到底发生了什么，试着运行一下编译器。输入 `coffee -c cup0.coffee` ，这将创建了一个名为 cup0.js 的文件，清单 4 展示了 cup0.js 内容。

##### 清单 4\. Cup 0 JavaScript

```
(function() {
    var i;
    for (i = 0; i <= 5; i++) {
        console.log("Hello " + i);
    }
}).call(this);

```

Show moreShow more icon

CoffeeScript 一个优势是，即使它提供的是一种比 JavaScript 更简洁的语法，但是可以编译成非常简单、有逻辑的 JavaScript。您可能会好奇为什么所有代码都包装在一个函数中，这是因为 JavaScript 只支持函数级范围。通过把所有内容包装在一个函数中，确保了变量的作用域只在该函数范围内，不会变成全局的（或者修改现有全局变量）。

打开一个命名为 cup1.coffee 的新文件，输入更复杂化的代码，如清单 5 所示。

##### 清单 5\. Cup 1

```
stdin = process.openStdin()
stdin.setEncoding 'utf8'

stdin.on 'data', (input) ->
    name = input.trim()
    process.exit() if name == 'exit'
    console.log "Hello #{name}"
    console.log "Enter another name or 'exit' to exit"
console.log 'Enter your name'

```

Show moreShow more icon

清单 5 中的程序提示用户输入他们的名称，然后对他们表示欢迎。JavaScript 并没有任何内置库，用于从标准输入读取，但 Node.js 有，这是另一个利用 CoffeeScript/Node.js 的共生优势的示例。在诸如 C 一类的语言中，从标准输入读取是一种阻塞式调用。在从标准输入读取完成之前，不能执行任何代码。如果您熟悉 Node.js 的话，那么就会知道不能那样处理，因为 Node.js 不支持阻塞式 I/O。相反，您必须使用 `stdin.on` 注册一个回调。

运行 `coffee -c cup1.coffee` 来看一下 CoffeeScript 编译器产生的 JavaScript，如清单 6 所示。

##### 清单 6\. Cup 1 JavaScript

```
(function() {
    var stdin;
    stdin = process.openStdin();
    stdin.setEncoding('utf8');
    stdin.on('data', function(input) {
        var name;
        name = input.trim();
        if (name === 'exit') {
            process.exit();
        }
        console.log("Hello " + name);
        return console.log("Enter another name or 'exit' to exit");
    });
console.log('Enter your name');
}).call(this);

```

Show moreShow more icon

`stdin.on` 函数使用一种典型的事件绑定格式。您指定想要监听的事件 (‘data’) 的种类，然后给它一个回调函数，在事件被触发时执行。在编译的 JavaScript 中，看到的是冗长的 JavaScript，其创建一个内联函数，并把它传递给另一个函数。把这段代码和相当的 CoffeeScript 比较一下，是不是在后者中，这些括号、大括号、分号和关键字都不见了呢？

现在您已经知道如何编译 CoffeeScript 程序了，下一节内容我们将了解在学习 CoffeeScript 时最常用到的功能之一：REPL。

## REPL

REPL 是一个在许多编程语言中，特别是在各种函数式语言中都能找到的标准工具。REPL 就相当于 Ruby 的 IRB，只需输入 `coffee` 就会启动 CoffeeScript 的 REPL。例如，使用这一 CoffeeScript 功能来进行实验以及解决一些简单的问题，如清单 7 所示。

##### 清单 7\. 使用 REPL

```
coffee> nums = [1..10]
[ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
coffee> isOdd = (n) -> n % 2 == 1
[Function]
coffee> odds = nums.filter isOdd
[ 1, 3, 5, 7, 9 ]
coffee> odds.reduce (a,b) -> a + b
25

```

Show moreShow more icon

每次只要您在 REPL 中输入一个表达式，它就会计算表达式、输出结果，然后等待下一表达式。该示例定义了一个名为 `nums` 的变量，其取值区间为从 1 到 10。REPL 输出刚刚定义的变量值，该项功能可以立刻发挥作用。也许您并未记得定义的区间是否是闭区间（包含了最后一个数值，在本用例中为 10）或是开区间的。REPL 显示 10 被包含在内，所以是一个闭区间。如果您想要一个开区间的话，使用 `nums = [1...10]` 即可。

接下来定义的是 `isOdd` 函数。CoffeeScript 有着非常简洁的函数声明语法，鉴于 JavaScript 函数式特性，这是一个很好的功能。在本示例中，REPL 仅显示 `[Function]` ，让您知道 `isOdd` 变量相当于一个函数。接下来声明一个新变量 `odds` 。通过调用 `nums` 的过滤器函数并传入 `isOdd` 来获得 `odds` 的值。这会生成一个新数组，其中的元素在函数传递至 `odds` 时必须为真值。在 `odds` 上调用 `reduce` 函数。传递进一个函数，把数组中的每个值与之前的合数相加，求出了数组中的值总和，然后 REPL 显示和为 25。

下一节所讨论的主题对于任何 JavaScript 开发人员来说都会有亲近感：浏览器端脚本。

## 简单 Web 示例

您已经了解了如何编写 CoffeeScript 文件，以及如何把它们编译成 JavaScript，这些 JavaScript 当然可用于 Web 应用程序中。对于开发目的，则有一种更简单的做法，CoffeeScript 编译器可以在浏览器中运行，所以允许直接在网页中使用 CoffeeScript。不过，如果您想构建一个高性能 Web 应用程序，不建议以这种方式使用 CoffeeScript。CoffeeScript 编译器是一个大文件，要动态编译 CoffeeScript 肯定会影响某些方面的速度。CoffeeScript 的确提供了一种简单的方法来开发有着明显的 “生产路径” 的应用程序。

CoffeeScript 并非一种 JavaScript 工具包或框架，它是一种编程语言，因此，它并未包含许多 DOM 相关的便利函数。不过，您可以将 CoffeeScript 和您所喜欢的工具包结合使用，最常见的组合是将其与 jQuery 一起使用。清单 8 展示了一个使用 jQuery 和 CoffeeScript 的简单网页。

##### 清单 8\. Web 页面中的 CoffeeScript

```
<html>
<head>
    <script src="http://ajax.googleapis.com/ajax/libs/jquery/1.4.2/jquery.min.js">
</script>
    <script
src="http://jashkenas.github.com/coffee-script/extras/coffee-script.js">
</script>
    <script type="text/coffeescript">
        gcd = (a,b) -> if (b==0) then a else gcd(b, a % b)
        $("#button").click ->
            a = $("#a").val()
            b = $("#b").val()
            $("#c").html gcd(a,b)
    </script>
</head>
<body>
    A: <input type="text" id="a"/><br/>
    B: <input type="text" id="b"/><br/>
    <input type="button" value="Calculate GCD" id="button"/> <br/>
    C: <span id="c"/>
</body>
</html>

```

Show moreShow more icon

该网页使用了 Google 的 Ajax 库来加载 jQuery，从 CoffeeScript 的创建者 Jeremy Ashkenas 的 Github 资源库（参见 参考资料 ）载入 CoffeeScript 编译器库。该代码包含了一个脚本块，不过它的类型并非 `text/javascript` ，而是 `text/coffeescript` ，这就是 CoffeeScript 编译器如何知道要编译的脚本内容。该脚本接着创建了一个名为 `gcd` 的函数，该函数计算两个整数的最大公约数。jQuery 被用来为页面上的按钮创建单击处理程序，您可以从处理程序中获取两个文本输入域的值，并把它们传递给函数 `gcd` ，计算所得的结果被写回到网页中。 `$()` 、 `val()` 和 `html()` 等均是 jQuery 函数，但可以轻易地和 CoffeeScript 结合使用，并利用 CoffeeScript 简洁语法的优势。

## 结束语

在本文中，我们快速地了解一下 CoffeeScript。在开发环境设置完成并正常运行后，您现在可以使用 REPL 来探索 CoffeeScript 了。您已经学会了使用编译器来查看其 JavaScript 类型，而且也知道了如何直接在网页中编写和运行 CoffeeScript。尽管没有太多的解释，不过这些示例还是让我们尝到了 CoffeeSCript 语法的一些味道。

## 下载本文源代码

[cs1.zip](http://public.dhe.ibm.com/software/dw/web/wa-coffee1/cs1.zip)