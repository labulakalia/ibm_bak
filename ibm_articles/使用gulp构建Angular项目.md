# 使用 gulp 构建 Angular 项目
gulp 实战

**标签:** DevOps,Node.js,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/wa-cn-gulpangular/)

徐玉珠

发布: 2016-08-25

* * *

## 引言

随着 Web “前后端分离” 架构的日渐普及，Angular.js 因为其 MVVM 模式使其受到广大前端开发者的青睐。然而，在前台项目开发过程中，依赖包的管理、JavaScript / SASS / Less 的编译和压缩、图片的压缩、版本更新，这些零碎繁多的任务阻碍着项目的构建和部署。所以，前端项目的自动化构建也成为了开发中的必备环节。

## 什么是 gulp?

gulp 是一个构建工具，可以通过它自动执行网站开发过程中的公共任务，比如编译 SASS/Less，编译压缩混淆 JavaScript,，合并编译模板和版本控制等。因为 gulp 是基于 Node.js 构建的，所以 gulp 源文件和开发者自己定义的 gulpfile 都被写进 JavaScript 里，前端开发者可以用自己熟悉的语言来编写 gulp 任务。

gulp 本身并不能完成这么多种任务，不过它可以借助 npm 丰富的插件库。开发者可以在 npm 中搜索 gulpplugin 找到想要的插件。例如本文中将要提到的 gulp-cssmin, gulp-jshint, gulp-concat、gulp-inject 等等。

## 为什么选择 gulp？

其实现有的基于 Node.js 的构建工具有很多，比如 Bower，Yeoman，grunt 等。而且自 2013 年 grunt v0.4.0 发布以后，grunt 已经改变了前端的开发方式。那么为什么我们要选 gulp？

gulp 最大的特点是所有的任务都是以 Node.js Stream 的形式处理，构建流程可以由 Stream 之间的 pipe 来定义，省去了把中间文件写到磁盘再读取的过程，而且任务都是默认并行，速度比 grunt 快很多，配置也感觉更省心。

- 易于使用：采用代码优于配置策略，gulp 让简单的事情继续简单，复杂的任务变得可管理。
- 高效：gulp 基于 Node.js 流 Unix 管道连接的方式，不需要往磁盘写中间文件，可以更快地完成构建。
- 高质量：gulp 每个 task 只完成一个任务，提高 task 的重用度。
- 易于学习：gulp 核心 API 约 5 个，开发者能在很短的时间内学会，之后就可以通过管道流来组合自己需要的 task。

## 安装 gulp

首先，gulp 需要全局安装。注意一点，gulp 的安装依赖于 npm，所以要在环境中装好 Node.js。安装好 Node.js 以后运行以下命令：

##### 清单 1\. 全局安装 gulp

```
$ npm install –g gulp

```

Show moreShow more icon

其次，在项目目录中安装 gulp 为本地模块

##### 清单 2\. 本地安装 gulp

```
$ npm install –-save-dev gulp

```

Show moreShow more icon

最后，安装项目中依赖的 gulp plugin 模块，如 gulp-cssmin, gulp-jshint

##### 清单 3\. 安装 gulp 插件

```
$ npm install –-save-dev gulp-cssmin gulp-jshint

```

Show moreShow more icon

## 使用 gulp

首先，在项目根目录中创建一个 gulpfile.js 文件。gulpfile.js 是所有 gulp 的任务配置和启动文件，以下是一个最简易的 gulpfile.js。

##### 清单 4\. gulpfile.js

```
var gulp = require('gulp');
var cssmin = require('gulp-cssmin');

gulp.task('css-task', function() {
return gulp.src('style/main.css')
.pipe(cssmin())
.pipe(gulp.dest('./build'));
});

gulp.task('default', function() {
gulp.run(css-task);
});

```

Show moreShow more icon

其次，开发者通过调用 gulp 命令执行默认的名为 default 的任务，或者执行 gulp  来执行指定任务。

在上面这段代码中，gulpfile.js 加载了 gulp 和 gulp-cssmin 两个模块之后，执行了开发者指定的 task。Task 的定义有两个参数，第一个为 task 名称 (css-task)，第二个是任务的调用函数。在任务函数中，使用 gulp 模块的 src 方法，指定所要处理的文件，然后使用 pipe 方法，将上一步的输出转为当前的输入，进行链式处理。最后使用 gulp 模块的 dest 方法，将上一步的输出写入本地文件。

在这个例子中，可以清楚的看到 gulp 使用了 Node.js 的 stream 概念：src 方法读入文件产生数据流，dest 方法将数据流写入文件，中间是一些中间步骤，每一步都对数据流进行一些处理。

## 使用 gulp 自动化构建 Angular 项目

### 示例项目介绍

文中使用的例子是一个基于 Angular.js 实现的网页版 Todo App，在 Github 中下载 [angular-quickstart](https://github.com/timelyxyz/angular-quickstart) 。项目代码结构如下

##### 清单 5\. 项目目录结构

```
|--bower_components
|--build
|--node_modules
|--static_pages
|--js
|--controllers
|--services
|--app.js             // app 启动配置
|--style
|--main.css
|--view
|--note.html
|--gulpfile.js         // gulp 配置
|--bower.json             // bower 配置
|--package.json         // node module 配置
|--index.html             // app 启动文件

```

Show moreShow more icon

架构上使用如下 libs

- Angular：JavaScript 架构
- Bootstrap：样式控件库
- bower：管理项目依赖包
- npm：管理项目运行的依赖工具
- gulp：自动化构建工具

### 设计 build process

因为项目中使用到的 source code 有 JavaScript,，CSS，HTML 以及依赖 bower 导入的依赖包，所以我们将构建分成以下几个 task：

1. jshint：编译 JavaScript
2. clean：清空 build 目录
3. template：合并并打包所有 HTML 模板文件生成 template.js
4. js：压缩、混淆 JavaScript，并生成 sourcemap
5. deployCSS：合并压缩 CSS
6. devIndex：组织 develop 时的 index.html
7. deployIndex：组织 deploy 时的 index.html，区别于第 6 个 task，这个 index.html 中 link 的将是压缩过的 CSS 和 JavaScript

接下来，我们分条介绍上面的 task。

### 设定全局变量 paths，记录 source code 的路径

##### 清单 6\. 项目 paths

```
var paths = {
js: ['./js/app.js', 'js/**/*.js'],
css: ['./style/**/*.css'],
templates: './js/templates.js',
buildjs: ['./build/**/*.js'],
buildcss: ['./build/**/*.css']
};

```

Show moreShow more icon

#### task 1: jshint

这个 task 主要用来编译 JavaScript 代码。依赖插件 glup-jshint，需要在配置之前将这个插件安装好，之后如下配置：

##### 清单 7\. jshint

```
var jshint = require('gulp-jshint');
gulp.task('jshint', function() {
gulp.src(paths.js)
.pipe(jshint())
.pipe(jshint.reporter('default'));
});

```

Show moreShow more icon

#### task 2: clean

我们可以看到项目目录中有个 build 文件夹，这个目录是用来存放所有 build 之后的 CSS/JS/HTML 文件。所以 clean 的任务是在每次 build 之前清理上一次 build 的 outputs，也就是清空 build 文件夹。依赖 del 模块：

##### 清单 8\. clean

```
var del = require('del');
gulp.task('clean', function() {
// You can use multiple globbing patterns as you would with `gulp.src`
return del(['./build', paths.templates]);
});

```

Show moreShow more icon

#### task 3: template

开发过程中一般按照功能模块将页面分成不同的 HTML 片段模板，每个片段对应一个相应 controller，客户端浏览器访问网页时会分别请求下载 HTML 文件。当业务逻辑较为复杂时，项目中将会存在很多的 HTML 小文件，此时下载的请求数目也会很大，影响 performance。我们的想法是使用插件 gulp-angular-templatecache，将这些 HTML 模版片段合并成一个文件。这个插件会将所有模板文件合并成一个 template.js, 并将其作为项目中的一个 module 而存在。

##### 清单 9\. template

```
var templateCache = require('gulp-angular-templatecache');
gulp.task('template', function () {
return gulp.src('./view/**/*.html')
.pipe(templateCache({module: 'myApp'}))
.pipe(gulp.dest('./js'))
});

```

Show moreShow more icon

#### task 4: js

JavaScript 的 build 分成两种情况：

1. DEV 模式：在开发环境中，为了方便 debug，JavaScript 通常不做压缩和混淆；
2. PROD 模式：上生产环境时，JavaScript 必须并合压缩、混淆并成一个文件，同时需要去除所有的 log 输出语句。

JavaScript 混淆使用了插件 gulp-uglify；合并压缩使用了插件 gulp-concat，以及 gulp-sourcemaps；

去除 log 语句使用了插件 gulp-strip-debug。

##### 清单 10\. gulp 任务 js

```
gulp.task('js', function() {
if (deployEnvironment == Environment.DEV) { // DEV
return gulp.src(paths.js)
.pipe(concat('all.js'))
.pipe(gulp.dest('./build'));
} else { // PROD
return gulp.src(paths.js)
.pipe(sourcemaps.init())
.pipe(stripDebug())
.pipe(uglify())
.pipe(concat('all.min.js'))
.pipe(sourcemaps.write())
.pipe(gulp.dest('./build'));
}
});

```

Show moreShow more icon

#### task 5: css

文中示例的样式表使用的是原生的 CSS 文件，所以使用 cssmin 完成了所有 CSS 文件的合并及压缩。

##### 清单 11\. gulp 任务 css

```
var cssmin = require('gulp-cssmin');
gulp.task('deployCSS', function() {
return gulp.src(paths.css)
.pipe(cssmin())
.pipe(concat('all.css'))
.pipe(gulp.dest('./build'));
});

```

Show moreShow more icon

#### task 6: devIndex

这个 task 主要是管理开发使用的 index.html 文件中对静态资源的引用。因为我们使用了 bower 来管理项目的依赖包，比如 jQuery 和 Angular，我们使用下面这种方式来导入项目：

##### 清单 12\. 引用静态资源

```
<script src="bower_components/jquery/dist/jquery.js"></script>
<script src="bower_components/angular/angular.js"></script>

```

Show moreShow more icon

可以看出，当我们需要引用的静态资源变得很多时，这个 list 也将会变得很长很不好管理。这时可以使用插件 gulp-inject。

第 1 步，使用 bower 命令安装项目依赖包时添加–save-dev 参数，使配置写入 bower.json 文件。比如安装 angular 使用如下命令：

##### 清单 13\. 安装本地 angular 模块

```
$ bower install –-save-dev angular

```

Show moreShow more icon

安装成功以后 bower.json 文件的 dependencies 中会自动生成 angular 的依赖关系。

##### 清单 14.bower.json

```
"dependencies": {
"angular": "^1.5.7"
},

```

Show moreShow more icon

第 2 步，在 index.html 中使用 inject 标签管理需要插入依赖包的位置。整理 index.html 成如下格式：

##### 清单 15\. 为 gulp-inject 准备的 index.html 模板

```
<!DOCTYPE html>
<html lang="en" data-ng-app="myApp">
<head>
<meta charset="utf-8">
<title>Angular UI Router</title>
<!-- bower:css -->
<!-- endinject -->

<!-- inject:css -->
<!-- endinject -->
</head>
<body>
<a ui-sref="note">go note</a>
<div ui-view></div>
<!-- bower:js -->
<!-- endinject -->

<!-- inject:js -->
<!-- endinject -->
</body>
</html>

```

Show moreShow more icon

第 3 步，配置 inject 的 task

##### 清单 16.gulp 任务 inject

```
gulp.task('devIndex', ['clean', 'jshint'], function () {
// It's not necessary to read the files (will speed up things), we're only after their paths:
return gulp.src('./index.html')
.pipe(inject(gulp.src(paths.js, {read: false}), {relative: true}))
.pipe(inject(gulp.src(paths.css, {read: false}), {relative: true}))
.pipe(inject(gulp.src(bowerFiles(), {read: false}), {name: 'bower', relative: true}))
.pipe(gulp.dest('./'));
});

```

Show moreShow more icon

第 4 步，执行命令 gulp devIndex，我们可以得到注入以后的 index.html 如下：

##### 清单 17.Inject 后生成的 index.html

```
<!DOCTYPE html>
<html lang="en" data-ng-app="myApp">
<head>
<meta charset="utf-8">
<title>Angular UI Router</title>
<!-- bower:css -->
<link rel="stylesheet" href="bower_components/normalize-css/normalize.css">
<link rel="stylesheet" href="bower_components/bootstrap/dist/css/bootstrap.min.css">
<link rel="stylesheet" href="bower_components/bootstrap/dist/css/bootstrap-theme.css">
<!-- endinject -->

<!-- inject:css -->
<link rel="stylesheet" href="style/main.css">
<!-- endinject -->
</head>
<body>
<a ui-sref="note">go note</a>
<div ui-view></div>

<!-- bower:js -->
<script src="bower_components/angular/angular.js"></script>
<!-- endinject -->

<!-- inject:js -->
<script src="js/app.js"></script>
<script src="js/controllers/note.js"></script>
<script src="js/services/note.js"></script>
<!-- endinject -->
</body>
</html>

```

Show moreShow more icon

#### task 7: deployIndex

deployIndex 主要是用来组织部署到测试或者生产环境时需要的 index.html。和 devIndex 的区别在于页面中引用的是合并压缩混淆后的静态资源，也就是项目目录 build 中的文件。此处同样也是使用的 gulp-inject，我们直接看一下 task。

##### 清单 18\. gulp 任务 deployIndex

```
gulp.task('deployIndex', ['clean', 'jshint', 'template', 'deployJS', 'deployCSS'], function () {
// It's not necessary to read the files (will speed up things), we're only after their paths:
return gulp.src('./index.html')
.pipe(inject(gulp.src(paths.buildjs, {read: false}), {relative: true}))
.pipe(inject(gulp.src(paths.buildcss, {read: false}), {relative: true}))
.pipe(inject(gulp.src(bowerFiles(), {read: false}), {name: 'bower', relative: true}))
.pipe(gulp.dest('./'));
});

```

Show moreShow more icon

上面这段代码在配置 deployIndex 时使用了三个参数，其中第二个参数表示当前 task 所依赖的 task list。

调用 gulp deployIndex，生成的 index.html 和上一个 task 类似，不过 CSS 和 JS 的引用会修改如下：

##### 清单 19\. Inject 之后生成的 index.html

```
<!-- inject:css -->
<link rel="stylesheet" href="build/all.css">
<!-- endinject -->
<!-- inject:js -->
<script src="build/all.min.js"></script>
<!-- endinject -->

```

Show moreShow more icon

## 结束语

随着”前后端分离”架构的日渐普及，前台的自动化构建将越来越被重视。本文着重介绍了 gulp 的使用方式，以及 Web 前台开发过程中涉及的一些自动化构建切入点，希望读者可以通过本文了解 Web 前台的自动化构建相关知识。