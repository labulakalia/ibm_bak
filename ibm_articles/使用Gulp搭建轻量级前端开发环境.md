# 使用 Gulp 搭建轻量级前端开发环境
加快中小型项目开发的速度, 适用于简单前端页面开发

**标签:** Node.js,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/wa-using-gulp-to-build-lightweight-frontend-environment/)

党 王虎

发布: 2019-01-21

* * *

## 背景

在 Web 开发大行其道的时代，迭代开发越来越注重开发效率，使用 React、Angular 这些大型前端框架来开发一些复杂交互的页面，第一步搭框架可能就会费很多精力。在中小型项目中，没有复杂交互下，有可能您的网页只是嵌入到大型前端 APP 的一部分，业务逻辑不复杂，更多的只是展示页面，而简单的不用任何的模板框架又不能复用页面的公共部分，一点不大的修改可能涉及到每个页面的变更，严重影响开发效率。如果您正面临此问题，那么您有必要阅读下去，这个使用 Gulp 搭建的框架核心功能是模板复用。

Gulp 是基于 Node.js 的一个前端自动化构建工具，您可以使用它构建自动化工作流程（前端集成开发环境），简化工作量，从而把重点放在功能的开发上，提高您的开发效率和工作质量。关于 Gulp 的详细介绍与使用可以参考： [https://www.ibm.com/developerworks/cn/web/wa-cn-gulpangular/index.html](https://www.ibm.com/developerworks/cn/web/wa-cn-gulpangular/index.html)

本文将手把手教您运用最简单的方法实现一个轻量级的前端开发框架，您可以在日常开发一些交互不复杂的应用中使用它来提升您的开发效率。

## 先决条件

在我们开始构建项目之前，需要先配置开发环境：

- 安装 Node 环境

    [下载 Node 包](https://nodejs.org/en/download/), 您可以自行搜索如何安装，本文不做延伸介绍。

- 全局安装 Gulp 包

    运行命令： `npm run gulp -g`


## 开启本地 HTTP 服务器

在日常的开发过程中，常遇到开发一些简单的 HTML demo 页面，但是以浏览器打开本地文件的的方式，页面中如果有请求远程数据服务器的情况，浏览器会就会报跨域错误，无法获取服务器的数据，当然您可以把页面放到 Tomcat 或者 Apache HTTP Server 这样的应用服务器里面，但是在项目开发过程中这是件繁琐的事情，这时候就需要启一个简单的本地 HTTP 服务器，最好是基于 Node 平台，那么 gulp-webserver 是个不错的选择，它的使用如下：

**清单 1\. 配置 gulp-webserver 服务器**

```
const WebServer = require('gulp-webserver');
WebServer({
    port: 8080, // 端口
    host: 'localhost' // 地址
});

```

Show moreShow more icon

## 构建项目

### 示例项目结构

首先，我们创建一个项目，命名为 example，项目代码结构如清单 2 所示：

**清单 2\. 示例项目代码结构**

```
|--build
|--example
      |--css
      |--images
      |--js
      |--index.html
      |--overview.html
|--node_modules
|--src
|--css
      |--overview.css
      |--base.css
|--images
|--include
      |--base.html
      |--header.html
      |--footer.html
|--js
      |--overview.js
      |--base.js
|--view
      |--overview.html
      |--index.html
|--gulpfile.js
|--package.json

```

Show moreShow more icon

构建这个项目的目的是把我们开发目录下（ `src` ）的相关 `css、image、js、html` 文件都按照设定好的目录结构打包以静态文件形式放到 `build` 目录下， `build` 目录下的文件可以直接放在 Apache、Tomcat 等任何一种应用服务器中运行，不依赖额外的 lib。

### 如何构建项目

按照以下步骤构建项目。以下 `Dist` 变量为 `build` 的路径：’ `build/example'` 。

**注意：** _由于 本文 属 轻量级框架介绍 ，所以 暂时 不 涉及`js/css`文件合并。_

1. 打包 `html` 文件，如清单 3 所示：

    **清单 3\. 打包 HTML**





    ```
    Gulp.task('copy-html', () => {
        return Gulp.src('src/view/*.html')
            .pipe(FileInclude({ // HTML模板替换，具体用法见下文
                prefix: '##',
                basepath: '@file'
            })).on('error', function(err) {
                console.error('Task:copy-html,', err.message);
                this.end();
            }).pipe(Gulp.dest(Dist)); // 拷贝
    });

    ```





    Show moreShow more icon

2. 打包 `js` 文件，如清单 4 所示：

    **清单 4\. 打包 JS**





    ```
    const Uglify = require('gulp-uglify');
    Gulp.task('copy-js', () => {
    Gulp.src('src/js/**')
    .pipe(Uglify()) // 压缩js
    .pipe(Gulp.dest(Dist + '/js')) // 拷贝
    });

    ```





    Show moreShow more icon

3. 打包 css 文件，如清单 5 所示：

    **清单 5\. 打包 CSS**





    ```
    const Minifycss = require('gulp-minify-css');
    Gulp.task('copy-css', () => {
    Gulp.src('src/css/*.css')
    .pipe(Minifycss()) // 压缩css
    .pipe(Gulp.dest(Dist + '/css')) // 拷贝
    });

    ```





    Show moreShow more icon

4. 打包 `image` 文件，如清单 6 所示：

    **清单 6\. 打包 Image**





    ```
    Gulp.task('copy-images', () => {
    Gulp.src('src/images/*').pipe(Gulp.dest(Dist + '/images'));
    });

    ```





    Show moreShow more icon


### 自动构建并刷新浏览器

日常开发中，我们经常会频繁的修改代码来验证某个功能，从而需要频繁的执行构建操作，构建完成后还需手动刷新浏览器来验证代码，往往浪费了很多精力和时间，下面将介绍如何通过 Gulp 监听文件变化来自动构建以及自动刷新浏览器，这样可以从繁琐的操作中解放出来而专注于业务开发。

1. 自动构建监听文件变化。通过 `gulp watch` 命令监听所有的 `js、css、html` 以及 `image` 的更新，自动拷贝对应的最新资源到 `build` 目录下，具体代码如下：

    **清单 7\. 监听文件变化**





    ```
    Gulp.task('watch', () => {
        Gulp.watch('src/view/*', ['copy-html']); // 监听HTML变化
        Gulp.watch('src/js/**', ['copy-js']); // 监听js变化
        Gulp.watch('src/css/*', ['copy-css']); // 监听css变化
        Gulp.watch('src/images/*', ['copy-images']); // 监听image变化
    });

    ```





    Show moreShow more icon

2. 监听 `build` 目录变化实现浏览器自动刷新。通过配置 `webserver` 监听 `build` 目录变化来实现浏览器自动刷新，如清单 8 所示：

    **清单 8\. 监听 build 目录变化**





    ```
    const WebServer = require('gulp-webserver');
    Gulp.src(Dist) // 监听build目录
        .pipe(WebServer({
        livereload: true, // 配置热刷新
        open: 'http://localhost:8080/overview.html' // 配置启动成功自动打开指定页面
    }));

    ```





    Show moreShow more icon


## 配置 HTML 模板

### 配置替换 HTML 模板

配置 HTML 模板的原理是通过 `gulp-file-include` 组件在构建过程中自动去替换文件名为该文件内容，配置方法如清单 9 所示：

**清单 9\. HTML 模板配置方法**

```
const FileInclude = require('gulp-file-include');
Gulp.src('src/view/*.html') // 需要替换的文件路径
        .pipe(FileInclude({
            prefix: '##', // 自定义字符作为引用前缀
            basepath: '@file'
        }))

```

Show moreShow more icon

### 配置公共页面

清单 10 是一个公共页面示例：

**清单 10\. 公共页面示例**

```
base.html:
< meta charset="UTF-8">
    < script type="text/javascript" src="./js/plugin/jquery-3.2.1.min.js">< /script>
    < link rel="stylesheet" type="text/css" href="./css/base.css">< /link>

```

Show moreShow more icon

### 其他页面引用方法

在其他页面中引用方法如清单 11 所示：

**清单 11\. 其他页面中引用方法**

```
<head>
    <title>Test Page</title>
##include('../include/base.html') // ##include
<link rel="stylesheet" type="text/css" href="./css/overview.css"></link>
</head>

```

Show moreShow more icon

## 完整的 gulpfile 示例

清单 12 展示了一个完整的 gulpfile 示例：

**清单 12\. 完整 gulpfile 示例**

```
const Gulp = require('gulp');
const Minifycss = require('gulp-minify-css');
const Uglify = require('gulp-uglify');
const FileInclude = require('gulp-file-include');
const Watch = require('gulp-watch');
const WebServer = require('gulp-webserver');
const RunSequence = require('gulp-run-sequence');
const Clean = require('gulp-clean');

const Dist = 'build/example';

Gulp.task('copy-html', () => {
    return Gulp.src('src/view/*.html')
        .pipe(FileInclude({
            prefix: '##',
            basepath: '@file'
        })).on('error', function(err) {
            console.error('Task:copy-html,', err.message);
            this.end();
        }).pipe(Gulp.dest(Dist));
});

Gulp.task('copy-js', () => {
    Gulp.src('src/js/**').pipe(Uglify()).pipe(Gulp.dest(Dist + '/js'));
});

Gulp.task('copy-css', () => {
    Gulp.src('src/css/*.css').pipe(Minifycss()).pipe(Gulp.dest(Dist + '/css'));
});

Gulp.task('copy-images', () => {
    Gulp.src('src/images/*').pipe(Gulp.dest(Dist + '/images'));
});

Gulp.task('watch', () => {
    Gulp.watch('src/view/*', ['copy-html']);
    Gulp.watch('src/js/**', ['copy-js']);
    Gulp.watch('src/css/*', ['copy-css']);
    Gulp.watch('src/images/*', ['copy-images']);
});

Gulp.task('copy-sources', ['copy-css','copy-js','copy-html', 'copy-images']);

Gulp.task('web-server', () => {
    Gulp.src(Dist).pipe(WebServer({
        port: 8080,
        host: 'localhost',
        livereload: true,
        open: 'http://localhost:8080/overview.html'
    }));
});

Gulp.task('clean', () => {
    return Gulp.src(Dist).pipe(Clean());
});

Gulp.task('start', () => {
    RunSequence('clean', ['copy-sources', 'watch'], 'web-server'); // RunSequence是用来设置任务串行执行，因为有些任务是有先后顺序依赖，[]内的并行执行，()内的串行执行
});

Gulp.task('default', ['start']);

```

Show moreShow more icon

## 完整的项目依赖包

在 `package.json` 文件中添加如下依赖，因为是开发时依赖，所以在 `devDependencies` 模块中声明。

**清单 13\. 在 `package.json` 文件中添加依赖**

```
"devDependencies": {
    "gulp": "^3.9.0",
    "gulp-clean": "^0.4.0",
    "gulp-file-include": "^2.0.1",
    "gulp-minify-css": "^1.2.4",
    "gulp-run-sequence": "^0.3.2",
    "gulp-uglify": "^3.0.0",
    "gulp-watch": "^5.0.0",
    "gulp-webserver": "^0.9.1"
}

```

Show moreShow more icon

## 结束语

本文详细介绍了如何使用 Gulp 来搭建一个小型简单的框架，希望本文能够帮助您在日常开发中快速建立一个小型的前端项目。