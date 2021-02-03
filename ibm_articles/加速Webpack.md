# 加速 Webpack
学习 Webpack 构建的技巧，提高构建速度

**标签:** JavaScript,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/wa-lo-expedite-webpack/)

吴浩麟

发布: 2018-01-09

* * *

Web 应用日益复杂，相关开发技术也百花齐放，这对前端构建工具提出了更高的要求。 Webpack 从众多构建工具中脱颖而出成为目前最流行的构建工具，几乎成为目前前端开发里的必备工具之一。 大多数人在使用 Webpack 的过程中都会遇到构建速度慢的问题，在项目大时显得尤为突出，这极大的影响了我们的开发体验，降低了我们的开发效率。

本文将传授你一些加速 Webpack 构建的技巧，下面来一一介绍。

## 通过多进程并行处理

由于有大量文件需要解析和处理，构建是文件读写和计算密集型的操作，特别是当文件数量变多后，Webpack 构建慢的问题会显得严重。 运行在 Node.js 之上的 Webpack 是单线程模型的，也就是说 Webpack 需要处理的任务需要一件件挨着做，不能多个事情一起做。

文件读写和计算操作是无法避免的，那能不能让 Webpack 同一时刻处理多个任务，发挥多核 CPU 电脑的威力，以提升构建速度呢？

### 使用 HappyPack

HappyPack 就能让 Webpack 做到上面抛出的问题，它把任务分解给多个子进程去并发的执行，子进程处理完后再把结果发送给主进程。

接入 HappyPack 的相关代码如下：

```
    const path = require('path');
    const  ExtractTextPlugin =  require('extract-text-webpack-plugin');
    const  HappyPack = require('happypack');
    module.exports = {
        module: {
            rules: [
                {    test: /\.js$/,
                    // 把对 .js 文件的处理转交给 id 为 babel 的 HappyPack 实例
                    use:['happypack/loader?id=babel'],
                    // 排除 node_modules 目录下的文件，node_modules目录下的文件都是采用的 ES5 语法，没必要再通过 Babel 去转换
                    exclude: path.resolve(__dirname, 'node_modules'),
                 },
                {
                    // 把对 .css 文件的处理转交给 id 为 css 的 HappyPack 实例
                     test: /\.css$/,
                     use:ExtractTextPlugin.extract({
                        use: ['happypack/loader?id=css'],
             }),
        },
] },
    plugins: [
        new HappyPack({
            // 用唯一的标识符 id 来代表当前的HappyPack 是用来处理一类特定的文件
        id: 'babel',
            // 如何处理 .js 文件，用法和 Loader配置中一样
        loaders: ['babel-loader?cacheDirectory'],
     }),
        new HappyPack({
                id: 'css',
                    // 如何处理 .css 文件，用法和Loader 配置中一样
                loaders: ['css-loader'], }),
                new ExtractTextPlugin({
                    filename: `[name].css`,
            }),
        ],
    };

```

Show moreShow more icon

以上代码有两点重要的修改：

- 在 Loader 配置中，所有文件的处理都交给了 happypack/loader 去处理，使用紧跟其后的 querystring ?id=babel 去告诉 happypack/loader 去选择哪个 HappyPack 实例去处理文件。
- 在 Plugin 配置中，新增了两个 HappyPack 实例分别用于告诉 happypack/loader 去如何处理 .js 和 .css 文件。选项中的 id 属性的值和上面 querystring 中的 ?id=babel 相对应，选项中的 loaders 属性和 Loader 配置中一样。

接入 HappyPack 后，你需要给项目安装新的依赖：

`npm i -D happypack`

安装成功后重新执行构建，你就会看到以下由 HappyPack 输出的日志：

```
Happy[babel]: Version: 4.0.0-beta.5. Threads: 3

Happy[babel]: All set; signaling webpack to proceed.Happy[css]: Version: 4.0.0-beta.5. Threads: 3Happy[css]: All set; signaling webpack to proceed.

```

Show moreShow more icon

说明你的 HappyPack 配置生效了，并且可以得知 HappyPack 分别启动了3个子进程去并行的处理任务。

在整个 Webpack 构建流程中，最耗时的流程可能就是 Loader 对文件的转换操作了，因为要转换的文件数据巨多，而且这些转换操作都只能一个个挨着处理。 HappyPack 的核心原理就是把这部分任务分解到多个进程去并行处理，从而减少了总的构建时间。

从前面的使用中可以看出所有需要通过 Loader 处理的文件都先交给了 happypack/loader 去处理，收集到了这些文件的处理权后 HappyPack 就好统一分配了。

每通过 new HappyPack() 实例化一个 HappyPack 其实就是告诉 HappyPack 核心调度器如何通过一系列 Loader 去转换一类文件，并且可以指定如何给这类转换操作分配子进程。

核心调度器的逻辑代码在主进程中，也就是运行着 Webpack 的进程中，核心调度器会把一个个任务分配给当前空闲的子进程，子进程处理完毕后把结果发送给核心调度器，它们之间的数据交换是通过进程间通信 API 实现的。

核心调度器收到来自子进程处理完毕的结果后会通知 Webpack 该文件处理完毕。

### 使用 ParallelUglifyPlugin

在使用 Webpack 构建出用于发布到线上的代码时，都会有压缩代码这一流程。 最常见的 JavaScript 代码压缩工具是 UglifyJS，并且 Webpack 也内置了它。

用过 UglifyJS 的你一定会发现在构建用于开发环境的代码时很快就能完成，但在构建用于线上的代码时构建一直卡在一个时间点迟迟没有反应，其实卡住的这个时候就是在进行代码压缩。

由于压缩 JavaScript 代码需要先把代码解析成用 Object 抽象表示的 AST 语法树，再去应用各种规则分析和处理 AST，导致这个过程计算量巨大，耗时非常多。

为什么不把多进程并行处理的思想也引入到代码压缩中呢？

ParallelUglifyPlugin 就做了这个事情。 当 Webpack 有多个 JavaScript 文件需要输出和压缩时，原本会使用 UglifyJS 去一个个挨着压缩再输出， 但是 ParallelUglifyPlugin 则会开启多个子进程，把对多个文件的压缩工作分配给多个子进程去完成，每个子进程其实还是通过 UglifyJS 去压缩代码，但是变成了并行执行。 所以 ParallelUglifyPlugin 能更快的完成对多个文件的压缩工作。

使用 ParallelUglifyPlugin 也非常简单，把原来 Webpack 配置文件中内置的 UglifyJsPlugin 去掉后，再替换成 ParallelUglifyPlugin，相关代码如下：

```
const path = require('path');
const  ParallelUglifyPlugin =  require('webpack-parallel-uglify-plugin');
module.exports = {
        plugins: [
            // 使用 ParallelUglifyPlugin 并行压缩输出的 JS 代码
            new ParallelUglifyPlugin({
                // 传递给 UglifyJS 的参数
                uglifyJS: {
                 },
            }),
        ],
    };

```

Show moreShow more icon

接入 ParallelUglifyPlugin 后，项目需要安装新的依赖：

`npm i -D webpack-parallel-uglify-plugin`

安装成功后，重新执行构建你会发现速度变快了许多。如果设置 cacheDir 开启了缓存，在之后的构建中会变的更快。

### 缩小文件搜索范围

Webpack 启动后会从配置的 Entry 出发，解析出文件中的导入语句，再递归的解析。 在遇到导入语句时 Webpack 会做两件事情：

1. 根据导入语句去寻找对应的要导入的文件。例如 require(‘react’) 导入语句对应的文件是 ./node\_modules/react/react.js，而require(‘./util’)导入语句 对应的文件是 ./util.js。
2. 根据找到的要导入文件的后缀，使用配置中的 Loader 去处理文件。例如使用 ES6 开发的 JavaScript 文件需要使用 babel-loader 去处理。

以上两件事情虽然对于处理一个文件非常快，但是当项目大了以后文件量会变的非常多，这时候构建速度慢的问题就会暴露出来。 虽然以上两件事情无法避免，但需要尽量减少以上两件事情的发生，以提高速度。

接下来一一介绍可以优化它们的途径。

#### 缩小 resolve.modules 的范围

Webpack的resolve.modules 用于配置 Webpack 去哪些目录下寻找第三方模块。

resolve.modules 的默认值是 [‘node\_modules’]，含义是先去当前目录下的 ./node\_modules 目录下去找想找的模块，如果没找到就去上一级目录 ../node\_modules 中找，再没有就去 ../../node\_modules 中找，以此类推，这和 Node.js 的模块寻找机制很相似。

当安装的第三方模块都放在项目根目录下的 ./node\_modules 目录下时，没有必要按照默认的方式去一层层的寻找，可以指明存放第三方模块的绝对路径，以减少寻找，配置如下：

```
module.exports = {
            resolve: {
                // 使用绝对路径指明第三方模块存放的位置，以减少搜索步骤
                // 其中 __dirname 表示当前工作目录，也就是项目根目录
                modules: [path.resolve(__dirname, 'node_modules')]
         },
};

```

Show moreShow more icon

#### 缩小 Loader 的命中范围

除此之外在使用 Loader 时可以通过 test 、 include 、 exclude 三个配置项来命中 Loader 要应用规则的文件。 为了尽可能少的让文件被 Loader 处理，可以通过 include 去命中只有哪些文件需要被处理。

以采用 ES6 的项目为例，在配置 babel-loader 时，可以这样：

```
module.exports = {
module: {
    rules: [
      {
        // 如果项目源码中只有 js 文件就不要写成 /\.jsx?$/，提升正则表达式性能
        test: /\.js$/,
        // babel-loader 支持缓存转换出的结果，通过 cacheDirectory 选项开启
        use: ['babel-loader?cacheDirectory'],
        // 只对项目根目录下的 src 目录中的文件采用 babel-loader
        include: path.resolve(__dirname, 'src'),
      },
    ]
},
};

```

Show moreShow more icon

你可以适当的调整项目的目录结构，以方便在配置 Loader 时通过 include 去缩小命中范围。

#### 缩小 resolve.extensions 的数量

在导入语句没带文件后缀时，Webpack 会自动带上后缀后去尝试询问文件是否存在。 Webpack 配置中的 resolve.extensions 用于配置在尝试过程中用到的后缀列表，默认是：

`extensions: ['.js', '.json']`

也就是说当遇到 require(‘./data’) 这样的导入语句时，Webpack 会先去寻找 ./data.js 文件，如果该文件不存在就去寻找 ./data.json 文件，如果还是找不到就报错。

如果这个列表越长，或者正确的后缀在越后面，就会造成尝试的次数越多，所以 resolve.extensions 的配置也会影响到构建的性能。 在配置 resolve.extensions 时你需要遵守以下几点，以做到尽可能的优化构建性能：

- 后缀尝试列表要尽可能的小，不要把项目中不可能存在的情况写到后缀尝试列表中。
- 频率出现最高的文件后缀要优先放在最前面，以做到尽快的退出寻找过程。
- 在源码中写导入语句时，要尽可能的带上后缀，从而可以避免寻找过程。例如在你确定的情况下把 require(‘./data’) 写成 require(‘./data.json’)。

相关 Webpack 配置如下：

```
module.exports = {
resolve: {
    // 尽可能的减少后缀尝试的可能性
    extensions: ['js'],
},
};

```

Show moreShow more icon

#### 缩小 resolve.mainFields 的数量

Webpack 配置中的 resolve.mainFields 用于配置第三方模块使用哪个入口文件。

安装的第三方模块中都会有一个 package.json 文件用于描述这个模块的属性，其中有些字段用于描述入口文件在哪里，resolve.mainFields 用于配置采用哪个字段作为入口文件的描述。

可以存在多个字段描述入口文件的原因是因为有些模块可以同时用在多个环境中，针对不同的运行环境需要使用不同的代码。 以 isomorphic-fetchfetch API 为例，它是 的一个实现，但可同时用于浏览器和 Node.js 环境。

为了减少搜索步骤，在你明确第三方模块的入口文件描述字段时，你可以把它设置的尽量少。 由于大多数第三方模块都采用 main 字段去描述入口文件的位置，可以这样配置 Webpack：

```
module.exports = {
resolve: {
    // 只采用 main 字段作为入口文件描述字段，以减少搜索步骤
    mainFields: ['main'],
},
};

```

Show moreShow more icon

使用本方法优化时，你需要考虑到所有运行时依赖的第三方模块的入口文件描述字段，就算有一个模块搞错了都可能会造成构建出的代码无法正常运行。

## 善用现存的文件

### 通过 module.noParse 忽略文件

Webpack 配置中的 module.noParse 配置项可以让 Webpack 忽略对部分没采用模块化的文件的递归解析处理，这样做的好处是能提高构建性能。 原因是一些库，例如 jQuery 、ChartJS， 它们庞大又没有采用模块化标准，让 Webpack 去解析这些文件耗时又没有意义。

在上面的 _优化 resolve.alias 配置_ 中讲到单独完整的 react.min.js 文件就没有采用模块化，让我们来通过配置 module.noParse 忽略对 react.min.js 文件的递归解析处理， 相关 Webpack 配置如下：

```
module.exports = {
module: {
    // 独完整的 `react.min.js` 文件就没有采用模块化，忽略对 `react.min.js` 文件的递归解析处理
    noParse: [/react\.min\.js$/],
},
};

```

Show moreShow more icon

注意被忽略掉的文件里不应该包含 import、require、 define 等模块化语句，不然会导致构建出的代码中包含无法在浏览器环境下执行的模块化语句。

### 通过 resolve.alias 映射文件

Webpack 配置中的 resolve.alias 配置项通过别名来把原导入路径映射成一个新的导入路径。

在实战项目中经常会依赖一些庞大的第三方模块，以 React 库为例，库中包含两套代码：

- 一套是采用 CommonJS 规范的模块化代码，这些文件都放在 lib 目录下，以 package.json 中指定的入口文件 react.js 为模块的入口。
- 一套是把 React 所有相关的代码打包好的完整代码放到一个单独的文件中，这些代码没有采用模块化可以直接执行。其中 dist/react.js 是用于开发环境，里面包含检查和警告的代码。dist/react.min.js 是用于线上环境，被最小化了。

默认情况下 Webpack 会从入口文件 ./node\_modules/react/react.js 开始递归的解析和处理依赖的几十个文件，这会时一个耗时的操作。 通过配置 resolve.alias 可以让 Webpack 在处理 React 库时，直接使用单独完整的 react.min.js 文件，从而跳过耗时的递归解析操作。

相关 Webpack 配置如下：

```
module.exports = {
resolve: {
    // 使用 alias 把导入 react 的语句换成直接使用单独完整的 react.min.js 文件，
    // 减少耗时的递归解析操作
    alias: {
      'react': path.resolve(__dirname, './node_modules/react/dist/react.min.js'),
    }
},
};

```

Show moreShow more icon

除了 React 库外，大多数库发布到 Npm 仓库中时都会包含打包好的完整文件，对于这些库你也可以对它们配置 alias。

但是对于有些库使用本优化方法后会影响到后面要讲的使用 Tree-Shaking 去除无效代码的优化，因为打包好的完整文件中有部分代码你的项目可能永远用不上。 一般对整体性比较强的库采用本方法优化，因为完整文件中的代码是一个整体，每一行都是不可或缺的。 但是对于一些工具类的库，例如 [lodash](https://lodash.com/) ，你的项目可能只用到了其中几个工具函数，你就不能使用本方法去优化，因为这会导致你的输出代码中包含很多永远不会执行的代码。

## 使用 DllPlugin

在介绍 DllPlugin 前先给大家介绍下 DLL。 用过 Windows 系统的人应该会经常看到以 .dll 为后缀的文件，这些文件称为 **动态链接库** ，在一个动态链接库中可以包含给其他模块调用的函数和数据。

要给 Web 项目构建接入动态链接库的思想，需要完成以下事情：

- 把网页依赖的基础模块抽离出来，打包到一个个单独的动态链接库中去。一个动态链接库中可以包含多个模块。
- 当需要导入的模块存在于某个动态链接库中时，这个模块不能再次被打包，而是去动态链接库中获取。
- 页面依赖的所有动态链接库需要被加载。

为什么给 Web 项目构建接入动态链接库的思想后，会大大提升构建速度呢？ 原因在于包含大量复用模块的动态链接库只需要编译一次，在之后的构建过程中被动态链接库包含的模块将不会在重新编译，而是直接使用动态链接库中的代码。 由于动态链接库中大多数包含的是常用的第三方模块，例如 react、react-dom，只要不升级这些模块的版本，动态链接库就不用重新编译。

### 接入 Webpack

Webpack 已经内置了对动态链接库的支持，需要通过2个内置的插件接入，它们分别是：

- DllPlugin 插件：用于打包出一个个单独的动态链接库文件。
- DllReferencePlugin 插件：用于在主要配置文件中去引入 DllPlugin 插件打包好的动态链接库文件。

下面以基本的 React 项目为例，为其接入 DllPlugin，在开始前先来看下最终构建出的目录结构：

```
├── main.js

├── polyfill.dll.js

├── polyfill.manifest.json

├── react.dll.js

└── react.manifest.json

```

Show moreShow more icon

其中包含两个动态链接库文件，分别是：

- polyfill.dll.js 里面包含项目所有依赖的 polyfill，例如 Promise、fetch 等 API。
- react.dll.js 里面包含 React 的基础运行环境，也就是 react 和 react-dom 模块。

以 react.dll.js 文件为例，其文件内容大致如下：

```
var _dll_react = (function(modules) {
// ... 此处省略 webpackBootstrap 函数代码
}([
function(module, exports, __webpack_require__) {
    // 模块 ID 为 0 的模块对应的代码
}
// ... 此处省略剩下的模块对应的代码
]));

```

Show moreShow more icon

可见一个动态链接库文件中包含了大量模块的代码，这些模块存放在一个数组里，用数组的索引号作为 ID。 并且还通过 \_dll\_react 变量把自己暴露在了全局中，也就是可以通过 window.\_dll\_react 可以访问到它里面包含的模块。

其中 polyfill.manifest.json 和 react.manifest.json 文件也是由 DllPlugin 生成，用于描述动态链接库文件中包含哪些模块， 以 react.manifest.json 文件为例，其文件内容大致如下：

```
{
// 描述该动态链接库文件暴露在全局的变量名称
"name": "_dll_react",
"content": {
    "./node_modules/process/browser.js": {
      "id": 0,
      "meta": {}
    },
    // ... 此处省略部分模块
}
}

```

Show moreShow more icon

可见 manifest.json 文件清楚地描述了与其对应的 dll.js 文件中包含了哪些模块，以及每个模块的路径和 ID。

main.js 文件是编译出来的执行入口文件，当遇到其依赖的模块在 dll.js 文件中时，会直接通过 dll.js 文件暴露出的全局变量去获取打包在 dll.js 文件的模块。 所以在 index.html 文件中需要把依赖的两个 dll.js 文件给加载进去，index.html 内容如下：

```
<!--导入依赖的动态链接库文件-->
<script src="./dist/polyfill.dll.js"></script>
<script src="./dist/react.dll.js"></script>
<!--导入执行入口文件-->
<script src="./dist/main.js"></script>

```

Show moreShow more icon

以上就是所有接入 DllPlugin 后最终编译出来的代码，接下来教你如何实现。

#### 构建出动态链接库文件

构建输出的以下这四个文件

```
├── polyfill.dll.js

├── polyfill.manifest.json

├── react.dll.js

└── react.manifest.json

```

Show moreShow more icon

和以下这一个文件

```
├── main.js

```

Show moreShow more icon

是由两份不同的构建分别输出的。

与动态链接库相关的文件需要由一个独立的构建输出，用于给主构建使用。新建一个 Webpack 配置文件 webpack\_dll.config.js 专门用于构建它们，文件内容如下：

```
const path = require('path');
const DllPlugin = require('webpack/lib/DllPlugin');

module.exports = {
// JS 执行入口文件
entry: {
    // 把 React 相关模块的放到一个单独的动态链接库
    react: ['react', 'react-dom'],
    // 把项目需要所有的 polyfill 放到一个单独的动态链接库
    polyfill: ['core-js/fn/object/assign', 'core-js/fn/promise', 'whatwg-fetch'],
},
output: {
    // 输出的动态链接库的文件名称，[name] 代表当前动态链接库的名称，
    // 也就是 entry 中配置的 react 和 polyfill
    filename: '[name].dll.js',
    // 输出的文件都放到 dist 目录下
    path: path.resolve(__dirname, 'dist'),
    // 存放动态链接库的全局变量名称，例如对应 react 来说就是 _dll_react
    // 之所以在前面加上 _dll_ 是为了防止全局变量冲突
    library: '_dll_[name]',
},
plugins: [
    // 接入 DllPlugin
    new DllPlugin({
      // 动态链接库的全局变量名称，需要和 output.library 中保持一致
      // 该字段的值也就是输出的 manifest.json 文件 中 name 字段的值
      // 例如 react.manifest.json 中就有 "name": "_dll_react"
      name: '_dll_[name]',
      // 描述动态链接库的 manifest.json 文件输出时的文件名称
      path: path.join(__dirname, 'dist', '[name].manifest.json'),
    }),
],
};

```

Show moreShow more icon

#### 使用动态链接库文件

构建出的动态链接库文件用于在其它地方使用，在这里也就是给执行入口使用。

用于输出 main.js 的主 Webpack 配置文件内容如下：

```
const DllReferencePlugin = require('webpack/lib/DllReferencePlugin');

module.exports = {
plugins: [
    // 告诉 Webpack 使用了哪些动态链接库
    new DllReferencePlugin({
      // 描述 react 动态链接库的文件内容
      manifest: require('./dist/react.manifest.json'),
    }),
    new DllReferencePlugin({
      // 描述 polyfill 动态链接库的文件内容
      manifest: require('./dist/polyfill.manifest.json'),
    }),
],
devtool: 'source-map'
};

```

Show moreShow more icon

注意：在 webpack\_dll.config.js 文件中，DllPlugin 中的 name 参数必须和 output.library 中保持一致。 原因在于 DllPlugin 中的 name 参数会影响输出的 manifest.json 文件中 name 字段的值， 而在 webpack.config.js 文件中 DllReferencePlugin 会去 manifest.json 文件读取 name 字段的值， 把值的内容作为在从全局变量中获取动态链接库中内容时的全局变量名。

#### 执行构建

在修改好以上两个 Webpack 配置文件后，需要重新执行构建。 重新执行构建时要注意的是需要先把动态链接库相关的文件编译出来，因为主 Webpack 配置文件中定义的 DllReferencePlugin 依赖这些文件。

执行构建时流程如下：

1. 如果动态链接库相关的文件还没有编译出来，就需要先把它们编译出来。方法是执行 webpack –config webpack\_dll.config.js 命令。
2. 在确保动态链接库存在的前提下，才能正常的编译出入口执行文件。方法是执行 webpack 命令。这时你会发现构建速度有了非常大的提升。

相信给你的项目加上以上优化方法后，构建速度会大大提高，赶快去试试把！