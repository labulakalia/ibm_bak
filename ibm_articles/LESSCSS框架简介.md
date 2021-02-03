# LESS CSS 框架简介
使用 LESS 简化层叠样式表（CSS）的编写

**标签:** Web 开发

[原文链接](https://developer.ibm.com/zh/articles/1207-zhaoch-lesscss/)

赵春红, 王群峰, 张建峰

发布: 2012-07-11

* * *

## 简介

CSS（层叠样式表）是一门历史悠久的标记性语言，同 HTML 一道，被广泛应用于万维网（World Wide Web）中。HTML 主要负责文档结构的定义，CSS 负责文档表现形式或样式的定义。

作为一门标记性语言，CSS 的语法相对简单，对使用者的要求较低，但同时也带来一些问题：CSS 需要书写大量看似没有逻辑的代码，不方便维护及扩展，不利于复用，尤其对于非前端开发工程师来讲，往往会因为缺少 CSS 编写经验而很难写出组织良好且易于维护的 CSS 代码，造成这些困难的很大原因源于 CSS 是一门非程序式语言，没有变量、函数、SCOPE（作用域）等概念。LESS 为 Web 开发者带来了福音，它在 CSS 的语法基础之上，引入了变量，Mixin（混入），运算以及函数等功能，大大简化了 CSS 的编写，并且降低了 CSS 的维护成本，就像它的名称所说的那样，LESS 可以让我们用更少的代码做更多的事情。

## LESS 原理及使用方式

本质上，LESS 包含一套自定义的语法及一个解析器，用户根据这些语法定义自己的样式规则，这些规则最终会通过解析器，编译生成对应的 CSS 文件。LESS 并没有裁剪 CSS 原有的特性，更不是用来取代 CSS 的，而是在现有 CSS 语法的基础上，为 CSS 加入程序式语言的特性。下面是一个简单的例子：

##### 清单 1\. LESS 文件

```
@color: #4D926F;

#header {
color: @color;
}
h2 {
color: @color;
}

```

Show moreShow more icon

经过编译生成的 CSS 文件如下：

##### 清单 2\. CSS 文件

```
#header {
color: #4D926F;
}
h2 {
color: #4D926F;
}

```

Show moreShow more icon

从上面的例子可以看出，学习 LESS 非常容易，只要你了解 CSS 基础就可以很容易上手。

LESS 可以直接在客户端使用，也可以在服务器端使用。在实际项目开发中，我们更推荐使用第三种方式，将 LESS 文件编译生成静态 CSS 文件，并在 HTML 文档中应用。

### 客户端

我们可以直接在客户端使用 .less（LESS 源文件），只需要从 [http://lesscss.org](http://lesscss.org/) 下载 less.js 文件，然后在我们需要引入 LESS 源文件的 HTML 中加入如下代码：

```
<link rel="stylesheet/less" type="text/css" href="styles.less">

```

Show moreShow more icon

LESS 源文件的引入方式与标准 CSS 文件引入方式一样：

```
<link rel="stylesheet/less" type="text/css" href="styles.less">

```

Show moreShow more icon

需要注意的是：在引入 .less 文件时，rel 属性要设置为”stylesheet/less”。还有更重要的一点需要注意的是：LESS 源文件一定要在 less.js 引入之前引入，这样才能保证 LESS 源文件正确编译解析。

### 服务器端

LESS 在服务器端的使用主要是借助于 LESS 的编译器，将 LESS 源文件编译生成最终的 CSS 文件，目前常用的方式是利用 node 的包管理器 (npm) 安装 LESS，安装成功后就可以在 node 环境中对 LESS 源文件进行编译。

在项目开发初期，我们无论采用客户端还是服务器端的用法，我们都需要想办法将我们要用到的 CSS 或 LESS 文件引入到我们的 HTML 页面或是桥接文件中，LESS 提供了一个我们很熟悉的功能— Importing。我们可以通过这个关键字引入我们需要的 .less 或 .css 文件。 如：

@import “variables.less”;

.less 文件也可以省略后缀名，像这样：

@import “variables”;

引入 CSS 同 LESS 文件一样，只是 .css 后缀名不能省略。

### 使用编译生成的静态 CSS 文件

我们可以通过 LESS 的编译器，将 LESS 文件编译成为 CSS 文件，在 HTML 文章中引入使用。这里要强调的一点，LESS 是完全兼容 CSS 语法的，也就是说，我们可以将标准的 CSS 文件直接改成 .less 格式，LESS 编译器可以完全识别。

## 语法

### 变量

LESS 允许开发者自定义变量，变量可以在全局样式中使用，变量使得样式修改起来更加简单。

我们可以从下面的代码了解变量的使用及作用：

##### 清单 3\. LESS 文件

```
@border-color : #b5bcc7;

.mythemes tableBorder{
border : 1px solid @border-color;
}

```

Show moreShow more icon

经过编译生成的 CSS 文件如下：

##### 清单 4\. CSS 文件

```
.mythemes tableBorder {
border: 1px solid #b5bcc7;
}

```

Show moreShow more icon

从上面的代码中我们可以看出，变量是 VALUE（值）级别的复用，可以将相同的值定义成变量统一管理起来。

该特性适用于定义主题，我们可以将背景颜色、字体颜色、边框属性等常规样式进行统一定义，这样不同的主题只需要定义不同的变量文件就可以了。当然该特性也同样适用于 CSS RESET（重置样式表），在 Web 开发中，我们往往需要屏蔽浏览器默认的样式行为而需要重新定义样式表来覆盖浏览器的默认行为，这里可以使用 LESS 的变量特性，这样就可以在不同的项目间重用样式表，我们仅需要在不同的项目样式表中，根据需求重新给变量赋值即可。

LESS 中的变量和其他编程语言一样，可以实现值的复用，同样它也有生命周期，也就是 Scope（变量范围，开发人员惯称之为作用域），简单的讲就是局部变量还是全局变量的概念，查找变量的顺序是先在局部定义中找，如果找不到，则查找上级定义，直至全局。下面我们通过一个简单的例子来解释 Scope。

##### 清单 5\. LESS 文件

```
@width : 20px;
#homeDiv {
@width : 30px;
#centerDiv{
       width : @width;// 此处应该取最近定义的变量 width 的值 30px
              }
}
#leftDiv {
     width : @width; // 此处应该取最上面定义的变量 width 的值 20px

}

```

Show moreShow more icon

经过编译生成的 CSS 文件如下：

##### 清单 6\. CSS 文件

```
#homeDiv #centerDiv {
width: 30px;
}
#leftDiv {
width: 20px;
}

```

Show moreShow more icon

### Mixins（混入）

Mixins（混入）功能对用开发者来说并不陌生，很多动态语言都支持 Mixins（混入）特性，它是多重继承的一种实现，在 LESS 中，混入是指在一个 CLASS 中引入另外一个已经定义的 CLASS，就像在当前 CLASS 中增加一个属性一样。

我们先简单看一下 Mixins 在 LESS 中的使用：

##### 清单 7\. LESS 文件

```
// 定义一个样式选择器
.roundedCorners(@radius:5px) {
 -moz-border-radius: @radius;
 -webkit-border-radius: @radius;
border-radius: @radius;
}
// 在另外的样式选择器中使用
#header {
.roundedCorners;
}
#footer {
.roundedCorners(10px);
}

```

Show moreShow more icon

经过编译生成的 CSS 文件如下：

##### 清单 8\. CSS 文件

```
#header {
 -moz-border-radius:5px;
 -webkit-border-radius:5px;
border-radius:5px;
}
#footer {
 -moz-border-radius:10px;
 -webkit-border-radius:10px;
border-radius:10px;
}

```

Show moreShow more icon

从上面的代码我们可以看出：Mixins 其实是一种嵌套，它允许将一个类嵌入到另外一个类中使用，被嵌入的类也可以称作变量，简单的讲，Mixins 其实是规则级别的复用。

Mixins 还有一种形式叫做 Parametric Mixins（混入参数），LESS 也支持这一特性：

##### 清单 9\. LESS 文件

```
// 定义一个样式选择器
.borderRadius(@radius){
 -moz-border-radius: @radius;
 -webkit-border-radius: @radius;
border-radius: @radius;
}
// 使用已定义的样式选择器
#header {
.borderRadius(10px); // 把 10px 作为参数传递给样式选择器
}
.btn {
.borderRadius(3px);// // 把 3px 作为参数传递给样式选择器

}

```

Show moreShow more icon

经过编译生成的 CSS 文件如下：

##### 清单 10\. CSS 文件

```
#header {
 -moz-border-radius: 10px;
 -webkit-border-radius: 10px;
border-radius: 10px;
}
.btn {
 -moz-border-radius: 3px;
 -webkit-border-radius: 3px;
border-radius: 3px;
}

```

Show moreShow more icon

我们还可以给 Mixins 的参数定义一人默认值，如

##### 清单 11\. LESS 文件

```
.borderRadius(@radius:5px){
 -moz-border-radius: @radius;
 -webkit-border-radius: @radius;
border-radius: @radius;
}
.btn {
.borderRadius;
}

```

Show moreShow more icon

经过编译生成的 CSS 文件如下：

##### 清单 12\. CSS 文件

```
.btn {
 -moz-border-radius: 5px;
 -webkit-border-radius: 5px;
border-radius: 5px;
}

```

Show moreShow more icon

像 JavaScript 中 **arguments** 一样，Mixins 也有这样一个变量： **@arguments** 。@arguments 在 Mixins 中具是一个很特别的参数，当 Mixins 引用这个参数时，该参数表示所有的变量，很多情况下，这个参数可以省去你很多代码。

##### 清单 13\. LESS 文件

```
.boxShadow(@x:0,@y:0,@blur:1px,@color:#000){
 -moz-box-shadow: @arguments;
 -webkit-box-shadow: @arguments;
box-shadow: @arguments;
}
#header {
.boxShadow(2px,2px,3px,#f36);
}

```

Show moreShow more icon

经过编译生成的 CSS 文件如下：

##### 清单 14\. CSS 文件

```
#header {
 -moz-box-shadow: 2px 2px 3px #FF36;
 -webkit-box-shadow: 2px 2px 3px #FF36;
box-shadow: 2px 2px 3px #FF36;
}

```

Show moreShow more icon

Mixins 是 LESS 中很重要的特性之一，我们这里也写了很多例子，看到这些例子你是否会有这样的疑问：当我们拥有了大量选择器的时候，特别是团队协同开发时，如何保证选择器之间重名问题？如果你是 java 程序员或 C++ 程序员，我猜你肯定会想到命名空间 Namespaces，LESS 也采用了命名空间的方法来避免重名问题，于是乎 LESS 在 mixins 的基础上扩展了一下，看下面这样一段代码：

##### 清单 15\. LESS 文件

```
#mynamespace {
.home {...}
.user {...}
}

```

Show moreShow more icon

这样我们就定义了一个名为 mynamespace 的命名空间，如果我们要复用 user 这个选择器的时候，我们只需要在需要混入这个选择器的地方这样使用就可以了。#mynamespace > .user。

### 嵌套的规则

在我们书写标准 CSS 的时候，遇到多层的元素嵌套这种情况时，我们要么采用从外到内的选择器嵌套定义，要么采用给特定元素加 CLASS 或 ID 的方式。在 LESS 中我们可以这样写：

##### 清单 16\. HTML 片段

```
<div id="home">
<div id="top">top</div>
<div id="center">
<div id="left">left</div>
<div id="right">right</div>
</div>
</div>

```

Show moreShow more icon

##### 清单 17\. LESS 文件

```
#home{
color : blue;
width : 600px;
height : 500px;
border:outset;
#top{
        border:outset;
        width : 90%;
}
#center{
        border:outset;
        height : 300px;
        width : 90%;
        #left{
          border:outset;
          float : left;
width : 40%;
        }
        #right{
          border:outset;
          float : left;
width : 40%;
        }
    }
}

```

Show moreShow more icon

经过编译生成的 CSS 文件如下：

##### 清单 18\. CSS 文件

```
#home {
color: blue;
width: 600px;
height: 500px;
border: outset;
}
#home #top {
border: outset;
width: 90%;
}
#home #center {
border: outset;
height: 300px;
width: 90%;
}
#home #center #left {
border: outset;
float: left;
width: 40%;
}
#home #center #right {
border: outset;
float: left;
width: 40%;
}

```

Show moreShow more icon

从上面的代码中我们可以看出，LESS 的嵌套规则的写法是 HTML 中的 DOM 结构相对应的，这样使我们的样式表书写更加简洁和更好的可读性。同时，嵌套规则使得对伪元素的操作更为方便。

##### 清单 19\. LESS 文件

```
a {
color: red;
text-decoration: none;
&:hover {// 有 & 时解析的是同一个元素或此元素的伪类，没有 & 解析是后代元素
color: black;
text-decoration: underline;
}
}

```

Show moreShow more icon

经过编译生成的 CSS 文件如下：

##### 清单 20\. CSS 文件

```
a {
color: red;
text-decoration: none;
}
a:hover {
color: black;
text-decoration: underline;
}

```

Show moreShow more icon

### 运算及函数

在我们的 CSS 中充斥着大量的数值型的 value，比如 color、padding、margin 等，这些数值之间在某些情况下是有着一定关系的，那么我们怎样利用 LESS 来组织我们这些数值之间的关系呢？我们来看这段代码：

##### 清单 21\. LESS 文件

```
@init: #111111;
@transition: @init*2;
.switchColor {
color: @transition;
}

```

Show moreShow more icon

经过编译生成的 CSS 文件如下：

##### 清单 22\. CSS 文件

```
.switchColor {
color: #222222;
}

```

Show moreShow more icon

上面的例子中使用 LESS 的 operation 是 特性，其实简单的讲，就是对数值型的 value（数字、颜色、变量等）进行加减乘除四则运算。同时 LESS 还有一个专门针对 color 的操作提供一组函数。下面是 LESS 提供的针对颜色操作的函数列表：

```
lighten(@color, 10%); // return a color which is 10% *lighter* than @color
darken(@color, 10%); // return a color which is 10% *darker* than @color
saturate(@color, 10%); // return a color 10% *more* saturated than @color
desaturate(@color, 10%);// return a color 10% *less* saturated than @color
fadein(@color, 10%); // return a color 10% *less* transparent than @color
fadeout(@color, 10%); // return a color 10% *more* transparent than @color
spin(@color, 10); // return a color with a 10 degree larger in hue than @color
spin(@color, -10); // return a color with a 10 degree smaller hue than @color

```

Show moreShow more icon

PS: 上述代码引自 LESS CSS 官方网站，详情请见 [http://lesscss.org/#-color-functions](http://lesscss.org/#-color-functions)

使用这些函数和 JavaScript 中使用函数一样。

##### 清单 23 LESS 文件

```
init: #f04615;
#body {
background-color: fadein(@init, 10%);
}

```

Show moreShow more icon

经过编译生成的 CSS 文件如下：

##### 清单 24\. CSS 文件

```
#body {
background-color: #f04615;
}

```

Show moreShow more icon

从上面的例子我们可以发现，这组函数像极了 JavaScript 中的函数，它可以被调用和传递参数。这些函数的主要作用是提供颜色变换的功能，先把颜色转换成 HSL 色，然后在此基础上进行操作，LESS 还提供了获取颜色值的方法，在这里就不举例说明了。

LESS 提供的运算及函数特性适用于实现页面组件特性，比如组件切换时的渐入渐出。

### Comments（注释）

适当的注释是保证代码可读性的必要手段，LESS 对注释也提供了支持，主要有两种方式：单行注释和多行注释，这与 JavaScript 中的注释方法一样，我们这里不做详细的说明，只强调一点：LESS 中单行注释 (// 单行注释 ) 是不能显示在编译后的 CSS 中，所以如果你的注释是针对样式说明的请使用多行注释。

### LESS VS SASS

同类框架还有 SASS : [http://sass-lang.com/](http://sass-lang.com/), 与 LESS 相比，两者都属于 CSS 预处理器，功能上大同小异，都是使用类似程序式语言的方式书写 CSS, 都具有变量、混入、嵌套、继承等特性，最终目的都是方便 CSS 的书写及维护。

LESS 和 SASS 互相促进互相影响，相比之下 LESS 更接近 CSS 语法，更多两者之间的比较，请参考这篇帖子： [https://gist.github.com/674726](https://gist.github.com/674726) 。

## 结束语

本文提到的只是 LESS 的基本功能，更高级的功能如：字符串插值，服务器端使用配置，JavaScript 表达式，避免编译等可以参看 LESS 的官方网站。

LESS 以 CSS 语法为基础，又借用了很多我们熟知编程式语言的特性，这对于我们开发人员来讲学习成本几乎可以忽略，它在保留 CSS 语法的基础上扩展了更多实用的功能，LESS 为我们提供了一种新的编写样式表的方法，我们可以根据我们的项目特性选择使用 LESS 的部分特性，我们只需用很少的成本就可以换了很大的回报，一句话，Less is more，借助 LESS 可以更便捷的进行 Web 开发。