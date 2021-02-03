# 使用 PostCSS 进行 CSS 处理
PostCSS 及其常用插件

**标签:** Web 开发

[原文链接](https://developer.ibm.com/zh/articles/1604-postcss-css/)

成富

发布: 2016-04-26

* * *

在 Web 应用开发中，CSS 代码的编写是重要的一部分。CSS 规范从最初的 CSS1 到现在的 CSS3，再到 CSS 规范的下一步版本，规范本身一直在不断的发展演化之中。这给开发人员带来了效率上的提高。不过与其他 Web 领域的规范相似的处境是，CSS 规范在浏览器兼容性方面一直存在各种各样的问题。不同浏览器在 CSS 规范的实现方面的进度也存在很大差异。另外，CSS 规范本身的发展速度与社区的期待还有一定的差距。这也是为什么 SASS 和 LESS 等 CSS 预处理语言可以流行的重要原因。SASS 和 LESS 等提供了很多更实用的功能，也体现了开发人员对 CSS 语言的需求。本文中要介绍的 PostCSS 是目前流行的一个对 CSS 进行处理的工具。PostCSS 依托其强大的插件体系为 CSS 处理增加了无穷的可能性。

## PostCSS 介绍

PostCSS 本身是一个功能比较单一的工具。它提供了一种方式用 JavaScript 代码来处理 CSS。它负责把 CSS 代码解析成抽象语法树结构（Abstract Syntax Tree，AST），再交由插件来进行处理。插件基于 CSS 代码的 AST 所能进行的操作是多种多样的，比如可以支持变量和混入（mixin），增加浏览器相关的声明前缀，或是把使用将来的 CSS 规范的样式规则转译（transpile）成当前的 CSS 规范支持的格式。从这个角度来说，PostCSS 的强大之处在于其不断发展的插件体系。目前 PostCSS 已经有 200 多个功能各异的插件。开发人员也可以根据项目的需要，开发出自己的 PostCSS 插件。

PostCSS 从其诞生之时就带来了社区对其类别划分的争议。这主要是由于其名称中的 post，很容易让人联想到 PostCSS 是用来做 CSS 后处理（post-processor）的，从而与已有的 CSS 预处理（pre-processor）语言，如 SASS 和 LESS 等进行类比。实际上，PostCSS 的主要功能只有两个：第一个就是前面提到的把 CSS 解析成 JavaScript 可以操作的 AST，第二个就是调用插件来处理 AST 并得到结果。因此，不能简单的把 PostCSS 归类成 CSS 预处理或后处理工具。PostCSS 所能执行的任务非常多，同时涵盖了传统意义上的预处理和后处理。PostCSS 是一个全新的工具，给前端开发人员带来了不一样的处理 CSS 的方式。

## 使用 PostCSS

PostCSS 一般不单独使用，而是与已有的构建工具进行集成。PostCSS 与主流的构建工具，如 Webpack、Grunt 和 Gulp 都可以进行集成。完成集成之后，选择满足功能需求的 PostCSS 插件并进行配置。下面将具体介绍如何在 Webpack、Grunt 和 Gulp 中使用 PostCSS 的 Autoprefixer 插件。

### Webpack

Webpack 中使用 postcss-loader 来执行插件处理。在清单 1 中，postcss-loader 用来对.css 文件进行处理，并添加在 style-loader 和 css-loader 之后。通过一个额外的 postcss 方法来返回所需要使用的 PostCSS 插件。require(‘autoprefixer’) 的作用是加载 Autoprefixer 插件。

##### 清单 1\. 在 Webpack 中使用 PostCSS 插件

```
var path = require('path');

module.exports = {
context: path.join(__dirname, 'app'),
entry: './app',
output: {
path: path.join(__dirname, 'dist'),
filename: 'bundle.js'
},
module: {
loaders: [
{
test: /\.css$/,
loader: "style-loader!css-loader!postcss-loader"
}
]
},
postcss: function () {
return [require('autoprefixer')];
}
}

Gulp

```

Show moreShow more icon

Gulp 中使用 gulp-postcss 来集成 PostCSS。在清单 2 中，CSS 文件由 gulp-postcss 处理之后，产生的结果写入到 dist 目录。

##### 清单 2\. 在 Gulp 中使用 PostCSS 插件

```
var gulp = require('gulp');

gulp.task('css', function() {
var postcss = require('gulp-postcss');
return gulp.src('app/**/*.css')
.pipe(postcss([require('autoprefixer')]))
.pipe(gulp.dest('dist/'));
});

Grunt

```

Show moreShow more icon

Grunt 中使用 grunt-postcss 来集成 PostCSS。Grunt 中需要使用 grunt.loadNpmTasks 方法来加载插件，如清单 3 所示。

##### 清单 3\. 在 Grunt 中使用 PostCSS 插件

```
module.exports = function(grunt) {
grunt.initConfig({
postcss: {
options: {
processors: [
require('autoprefixer')()
]
},
dist: {
src: 'app/**/*.css',
expand: true,
dest: 'dist'
}
}
});
grunt.loadNpmTasks('grunt-postcss');
}

```

Show moreShow more icon

下面介绍常用的 PostCSS 插件。

## 常用插件

### Autoprefixer

Autoprefixer 是一个流行的 PostCSS 插件，其作用是为 CSS 中的属性添加浏览器特定的前缀。由于 CSS 规范的制定和完善一般需要花费比较长的时间，浏览器厂商在实现某个 CSS 新功能时，会使用特定的浏览器前缀来作为正式规范版本之前的实验性实现。比如 Webkit 核心的浏览器使用-webkit-，微软的 IE 使用-ms-。为了兼容不同浏览器的不同版本，在编写 CSS 样式规则声明时通常需要添加额外的带前缀的属性。这是一项繁琐而无趣的工作。Autoprefixer 可以自动的完成这项工作。Autoprefixer 可以根据需要指定支持的浏览器类型和版本，自动添加所需的带前缀的属性声明。开发人员在编写 CSS 时只需要使用 CSS 规范中的标准属性名即可。

清单 4 中给出了使用 CSS 弹性盒模型的 display 属性声明。

##### 清单 4\. 标准的 CSS 弹性盒模型的 display 属性声明

```
#content {
display: flex;
}

```

Show moreShow more icon

在经过 Autoprefixer 处理之后得到的 CSS 如清单 5 所示。

##### 清单 5\. 使用 Autoprefixer 处理之后的 CSS

```
#content {
display: -webkit-box;
display: -webkit-flex;
display: -ms-flexbox;
display: flex;
}

```

Show moreShow more icon

Autoprefixer 使用 Can I Use 网站提供的数据来确定所要添加的不同浏览器的前缀。随着浏览器版本的升级，浏览器在新版本中可能已经提供了对标准属性的支持，从而不再需要添加额外的前缀。Autoprefixer 可以配置需要支持的浏览器。如”last 2 versions”表示主流浏览器的最近两个版本，”ie 6-8”表示 IE 6 到 8，”> 1%”表示全球使用率大于 1%的浏览器版本。清单 6 中给出了配置 Autoprefixer 插件的示例。

##### 清单 6\. 配置 Autoprefixer 插件

```
require('autoprefixer')({
browsers: ['last 2 versions']
})

```

Show moreShow more icon

Autoprefixer 除了添加所需要的属性名称前缀之外，还可以移除 CSS 代码中冗余的属性名称前缀。遗留 CSS 代码中可能包含由开发人员手动添加的旧版本的浏览器所支持的带前缀的属性名称。Autoprefixer 默认情况下会移除这些冗余的前缀。可以通过配置对象中的 remove 属性来配置该行为。

### cssnext

cssnext 插件允许开发人员在当前的项目中使用 CSS 将来版本中可能会加入的新特性。cssnext 负责把这些新特性转译成当前浏览器中可以使用的语法。从实现角度来说，cssnext 是一系列与 CSS 将来版本相关的 PostCSS 插件的组合。比如，cssnext 中已经包含了对 Autoprefixer 的使用，因此使用了 cssnext 就不再需要使用 Autoprefixer。

**自定义属性和变量**

CSS 的层叠变量的自定义属性规范（CSS Custom Properties for Cascading Variables）允许在 CSS 中定义属性并在样式规则中作为变量来使用它们。自定义属性的名称以”–”开头。当声明了自定义属性之后，可以在样式规则中使用”var()”函数来引用，如清单 7 所示。

##### 清单 7\. CSS 自定义属性和变量

```
:root {
--text-color: black;
}

body {
color: var(--text-color);
}

```

Show moreShow more icon

在经过 cssnext 转换之后的 CSS 代码如清单 8 所示。

##### 清单 8\. 转换之后的 CSS 代码

```
body {
color: black;
}

自定义选择器

```

Show moreShow more icon

CSS 扩展规范（CSS Extensions）中允许创建自定义选择器，比如可以对所有的标题元素（h1 到 h6）创建一个自定义选择器并应用样式。通过”@custom-selector”来定义自定义选择器。在清单 9 中，”–heading”是自定义选择器的名称，其等同于选择器声明”h1, h2, h3, h4, h5, h6”。

##### 清单 9\. 自定义选择器

```
@custom-selector :--heading h1, h2, h3, h4, h5, h6;

:--heading {
font-weight: bold;
}

```

Show moreShow more icon

经过 cssnext 处理之后的 CSS 如清单 10 所示。

##### 清单 10\. 转换之后的 CSS 代码

```
h1,
h2,
h3,
h4,
h5,
h6 {
font-weight: bold;
}

样式规则嵌套

```

Show moreShow more icon

样式规则嵌套是一个非常实用的功能，可以减少重复的选择器声明。这也是 SASS 和 LESS 等 CSS 预处理器流行的一个重要原因。CSS 嵌套模块规范（CSS Nesting Module）中定义了标准的样式规则嵌套方式。可以通过 cssnext 把标准的样式嵌套格式转换成当前的格式。CSS 嵌套规范中定义了两种嵌套方式：第一种方式要求嵌套的样式声明使用”&”作为前缀，”&”只能作为声明的起始位置；第二种方式的样式声明使用”@nest”作为前缀，并且”&”可以出现在任意位置。清单 11 中给出了两种不同声明方式的示例。

##### 清单 11\. 样式规则嵌套

```
.message {
font-weight: normal;
& .header {
font-weight: bold;
}
@nest .body & {
color: black;
}
}

```

Show moreShow more icon

经过 cssnext 转换之后的 CSS 代码如清单 12 所示。

##### 清单 12\. 转换之后的 CSS 代码

```
.message {
font-weight: normal
}

.message .header {
font-weight: bold;
}

.body .message {
color: black;
}

CSS 模块化

```

Show moreShow more icon

在编写 CSS 代码时会遇到的一个很重要的问题是 CSS 代码的组织方式。当项目中包含的 CSS 样式非常多时，该问题尤其突出。这主要是由于不同 CSS 文件中的样式声明可能产生名称冲突。现在的 Web 开发大多采用组件化的组织方式，即把应用划分成多个不同的组件。每个组件都可以有自己的 CSS 样式声明。比如，两个组件的 CSS 中可能都定义了对于 CSS 类 title 的样式规则。当这两个组件被放在同一个页面上时，冲突的 CSS 类名称可能造成样式的错乱。对于此类 CSS 类名冲突问题，一般的做法是避免全局名称空间中的样式声明，而是每个组件有自己的名称空间。BEM 通过特殊的命名 CSS 类名的方式来解决名称冲突的问题。两个不同组件中的标题的 CSS 类名可能为 component1 **title 和 component2** title。

CSS 模块（CSS modules）并不要求使用 BEM 那样复杂的命名规范。每个组件可以自由选择最合适的简单 CSS 类名。组件的 CSS 类名在使用时会被转换成带唯一标识符的形式。这样就避免了名称冲突。在组件开发中可以继续使用简单的 CSS 类名，而不用担心名称冲突问题。清单 13 中给出了使用 CSS 模块规范的 CSS 代码。样式规则之前的”:global”表示这是一个全局样式声明。其他的样式声明是局部的。

##### 清单 13\. 使用 CSS 模块规范的 CSS 代码

```
:global .title {
font-size: 20px;
}

.content {
font-weight: bold;
}

```

Show moreShow more icon

经过转换之后的 CSS 样式声明如清单 14 所示。全局的 CSS 类名 title 保存不变，局部的 CSS 类名 content 被转换成\_content\_6xmce\_5。这样就确保了不会与其他组件中名称为 content 的类名冲突。

##### 清单 14\. 转换之后的 CSS 代码

```
.title {
font-size: 20px;
}

._content_6xmce_5 {
font-weight: bold;
}

```

Show moreShow more icon

由于在组件的 HTML 代码中引用的 CSS 类名和最终生成的类名并不相同，因此需要一个中间的过程来进行类名的转换。对于 React 来说，可以使用 react-css-modules 插件；在其他情况下，可以使用 PostHTML 对 HTML 进行处理。postcss-modules 插件把 CSS 模块中的 CSS 类名的对应关系保存在一个 JavaScript 对象中，可以被 PostHTML 中的 posthtml-css-modules 插件来使用。

在清单 15 中，在使用 postcss-modules 插件时提供了一个方法 getJSON。当 CSS 模块的转换完成时，该方法会被调用。该方法的参数 json 参数表示的是转换结果的 JavaScript 对象。该对象被以 JavaScript 文件的形式保存到 cssModules 目录下，并添加了模块导出的逻辑。

##### 清单 15\. 保存 CSS 模块的输出文件

```
require('postcss-modules')({
getJSON: function(cssFileName, json) {
var cssName = path.basename(cssFileName, '.css');
var jsonFileName = path.resolve(dist, 'cssModules', cssName + '.js');
mkdirp.sync(path.dirname(jsonFileName));
fs.writeFileSync(jsonFileName, "module.exports = " + JSON.stringify(json) + ";");
}
})

```

Show moreShow more icon

清单 16 中给出了使用 Gulp 的 gulp-posthtml 来处理 HTML 文件的任务。posthtml-css-modules 可以处理一个目录下的多个 CSS 模块输出文件。

##### 清单 16\. 使用 PostHTML 处理 HTML 里支持 CSS 模块

```
gulp.task('posthtml', function() {
var posthtml = require('gulp-posthtml');
return gulp.src('app/**/*.html')
.pipe(posthtml([ require('posthtml-css-modules')(path.join(dist, 'cssModules')) ]))
.pipe(gulp.dest('dist/'));
});

```

Show moreShow more icon

在 HTML 文件中使用”css-module”属性来指定对应的 CSS 类名。在 清单 17 中，名称”header.content”的 header 表示的是 CSS 文件名，而 content 是该文件中定义的 CSS 类名。

##### 清单 17\. 使用 CSS 模块的 HTML 文件

```
<div css-module="header.content">Hello world</div>

```

Show moreShow more icon

在经过处理之后，得到的 HTML 内容如清单 18 所示。

##### 清单 18\. 转换之后的 HTML 文件

```
<div class="_content_6xmce_5">Hello world</div>

资源文件处理

```

Show moreShow more icon

在 CSS 中经常会需要引用外部资源，如图片和字体等。在 CSS 代码中处理这些资源时会遇到一些常见的问题，比如图片的路径问题，内联图片内容，生成图片 Sprites 等。对于这些问题，都有专门的 PostCSS 插件来提供所需的功能。

postcss-assets 插件用来处理图片和 SVG。在 CSS 声明中引用图片时，可以使用 resolve 加上图片路径的形式，如”resolve(‘logo.png’)”。在插件处理时，会按照一定的顺序在指定的目录中查找该文件，如果找到，会用图片的真实路径来替换。可以通过选项 loadPaths 来指定查找的路径，basePath 来指定项目的根目录。在 CSS 声明中，可以使用 width、height 和 size 方法来获取到图片的宽度、高度和尺寸。当需要内联一个图片时，可以使用 inline 方法。inline 会把图片转换成 Base64 编码的 data url 的格式，这样可以减少对图片的 HTTP 请求。清单 19 给出了使用示例。

##### 清单 19\. Postcss-assets 插件使用示例

```
require('postcss-assets')({
loadPaths: ['assets/images']
})

```

Show moreShow more icon

清单 20 中给出了使用 resolve 的 CSS 样式声明。

##### 清单 20\. 使用 resolve 的 CSS 样式声明

```
.logo {
background-image: resolve('logo.png');
}

```

Show moreShow more icon

清单 21 中给出了 cssnext 处理之后的 CSS 代码。

##### 清单 21\. 转换之后的 CSS 代码

```
.logo {
background-image: url('/assets/images/logo.png');
}

其他插件

```

Show moreShow more icon

还有其他实用的 PostCSS 插件可以在开发中使用。如 postcss-stylelint 用来检查 CSS 中可能存在的格式问题。cssnano 用来压缩 CSS 代码。postcss-font-magician 用来生成 CSS 中的 @font-face 声明。precss 允许在 CSS 中使用类似 SASS 的语法。

## 开发 PostCSS 插件

虽然 PostCSS 已经有 200 多个插件，但在开发中仍然可能存在已有插件不能满足需求的情况。这个时候可以开发自己的 PostCSS 插件。开发插件是一件很容易的事情。每个插件本质只是一个 JavaScript 方法，用来对由 PostCSS 解析的 CSS AST 进行处理。

每个 PostCSS 插件都是一个 NodeJS 的模块。使用 postcss 的 plugin 方法来定义一个新的插件。插件需要一个名称，一般以”postcss-”作为前缀。插件还需要一个进行初始化的方法。该方法的参数是插件所支持的配置选项，而返回值则是另外一个方法，用来进行实际的处理。该处理方法会接受两个参数，css 代表的是表示 CSS AST 的对象，而 result 代表的是处理结果。清单 22 中给出了一个简单的 PostCSS 插件。该插件使用 css 对象的 walkDecls 方法来遍历所有的”color”属性声明，并对”color”属性值进行检查。如果属性值为 black，就使用 result 对象的 warn 方法添加一个警告消息。

##### 清单 22\. PostCSS 插件示例

```
var postcss = require('postcss');

module.exports = postcss.plugin('postcss-checkcolor', function(options) {
return function(css, result) {
css.walkDecls('color', function(decl) {
if (decl.value == 'black') {
result.warn('No black color.', {decl: decl});
}
});
};
})

```

Show moreShow more icon

清单 22 中的插件的功能比较简单。PostCSS 插件一般通过不同的方法来对当前的 CSS 样式规则进行修改。如通过 insertBefore 和 insertAfter 方法来插入新的规则。

## 结束语

对于 CSS 的处理一直都是 Web 开发中的一个复杂问题，其中一部分来源于 CSS 规范，一部分来源于不同浏览器实现带来的兼容性问题。PostCSS 为处理 CSS 提供了一种新的思路。通过 PostCSS 强大的插件体系，可以对 CSS 进行各种不同的转换和处理，从而尽可能的把繁琐复杂的工作交由程序去处理，而把开发人员解放出来。本文对 PostCSS 及其常用插件进行了详细介绍，希望可以帮助开发人员提高开发效率。

## 下载示例代码

[sourcecode.zip](http://www.ibm.com/developerWorks/cn/web/1604-postcss-css/sourcecode.zip)