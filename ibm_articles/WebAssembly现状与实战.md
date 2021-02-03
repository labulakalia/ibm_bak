# WebAssembly 现状与实战
WebAssembly 入门

**标签:** JavaScript,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/wa-lo-webassembly-status-and-reality/)

吴 浩麟

发布: 2018-06-13

* * *

## 为什么需要 WebAssembly

自从 JavaScript 诞生起到现在已经变成最流行的编程语言，这背后正是 Web 的发展所推动的。Web 应用变得更多更复杂，但这也渐渐暴露出了 JavaScript 的问题：

- 语法太灵活导致开发大型 Web 项目困难；
- 性能不能满足一些场景的需要。

针对以上两点缺陷，近年来出现了一些 JS 的代替语言，例如：

- 微软的 [TypeScript](http://www.typescriptlang.org/) 通过为 JS 加入静态类型检查来改进 JS 松散的语法，提升代码健壮性；
- 谷歌的 [Dart](https://www.dartlang.org/) 则是为浏览器引入新的虚拟机去直接运行 Dart 程序以提升性能；
- 火狐的 [asm.js](http://asmjs.org/) 则是取 JS 的子集，JS 引擎针对 asm.js 做性能优化。

以上尝试各有优缺点，其中：

- TypeScript 只是解决了 JS 语法松散的问题，最后还是需要编译成 JS 去运行，对性能没有提升；
- Dart 只能在 Chrome 预览版中运行，无主流浏览器支持，用 Dart 开发的人不多；
- asm.js 语法太简单、有很大限制，开发效率低。

三大浏览器巨头分别提出了自己的解决方案，互不兼容，这违背了 Web 的宗旨； 是技术的规范统一让 Web 走到了今天，因此形成一套新的规范去解决 JS 所面临的问题迫在眉睫。

于是 WebAssembly 诞生了，WebAssembly 是一种新的字节码格式，主流浏览器都已经支持 WebAssembly。 和 JS 需要解释执行不同的是，WebAssembly 字节码和底层机器码很相似可快速装载运行，因此性能相对于 JS 解释执行大大提升。 也就是说 WebAssembly 并不是一门编程语言，而是一份字节码标准，需要用高级编程语言编译出字节码放到 WebAssembly 虚拟机中才能运行，浏览器厂商需要做的就是根据 WebAssembly 规范实现虚拟机。

## WebAssembly 原理

要搞懂 WebAssembly 的原理，需要先搞懂计算机的运行原理。 电子计算机都是由电子元件组成，为了方便处理电子元件只存在开闭两种状态，对应着 0 和 1，也就是说计算机只认识 0 和 1，数据和逻辑都需要由 0 和 1 表示，也就是可以直接装载到计算机中运行的机器码。 机器码可读性极差，因此人们通过高级语言 C、C++、Rust、Go 等编写再编译成机器码。

由于不同的计算机 CPU 架构不同，机器码标准也有所差别，常见的 CPU 架构包括 x86、AMD64、ARM， 因此在由高级编程语言编译成可自行代码时需要指定目标架构。

WebAssembly 字节码是一种抹平了不同 CPU 架构的机器码，WebAssembly 字节码不能直接在任何一种 CPU 架构上运行， 但由于非常接近机器码，可以非常快的被翻译为对应架构的机器码，因此 WebAssembly 运行速度和机器码接近，这听上去非常像 Java 字节码。

相对于 JS，WebAssembly 有如下优点：

- 体积小：由于浏览器运行时只加载编译成的字节码，一样的逻辑比用字符串描述的 JS 文件体积要小很多；
- 加载快：由于文件体积小，再加上无需解释执行，WebAssembly 能更快的加载并实例化，减少运行前的等待时间；
- 兼容性问题少：WebAssembly 是非常底层的字节码规范，制订好后很少变动，就算以后发生变化,也只需在从高级语言编译成字节码过程中做兼容。可能出现兼容性问题的地方在于 JS 和 WebAssembly 桥接的 JS 接口。

每个高级语言都去实现源码到不同平台的机器码的转换工作是重复的，高级语言只需要生成底层虚拟机(LLVM)认识的中间语言(LLVM IR)， [LLVM](https://llvm.org/) 能实现：

- LLVM IR 到不同 CPU 架构机器码的生成；
- 机器码编译时性能和大小优化。

除此之外 LLVM 还实现了 LLVM IR 到 WebAssembly 字节码的编译功能，也就是说只要高级语言能转换成 LLVM IR，就能被编译成 WebAssembly 字节码，目前能编译成 WebAssembly 字节码的高级语言有：

- [AssemblyScript](https://github.com/AssemblyScript/assemblyscript):语法和 TypeScript 一致，对前端来说学习成本低，为前端编写 WebAssembly 最佳选择；
- c\\c++:官方推荐的方式，详细使用见 [文档](http://webassembly.org.cn/getting-started/developers-guide/);
- [Rust](https://www.rust-lang.org/):语法复杂、学习成本高，对前端来说可能会不适应。详细使用见 [文档](https://github.com/rust-lang-nursery/rust-wasm);
- [Kotlin](http://kotlinlang.org/):语法和 Java、JS 相似，语言学习成本低，详细使用见 [文档](https://kotlinlang.org/docs/reference/native-overview.html);
- [Golang](https://golang.org/):语法简单学习成本低。但对 WebAssembly 的支持还处于未正式发布阶段，详细使用见 [文档](https://blog.gopheracademy.com/advent-2017/go-wasm/) 。

通常负责把高级语言翻译到 LLVM IR 的部分叫做编译器前端，把 LLVM IR 编译成各架构 CPU 对应机器码的部分叫做编译器后端； 现在越来越多的高级编程语言选择 LLVM 作为后端，高级语言只需专注于如何提供开发效率更高的语法同时保持翻译到 LLVM IR 的程序执行性能。

## 编写 WebAssembly

### AssemblyScript 初体验

接下来详细介绍如何使用 AssemblyScript 来编写 WebAssembly，实现斐波那契序列的计算。 用 TypeScript 实现斐波那契序列计算的模块 f.ts 如下：

```
export function f(x: i32): i32 {
    if (x === 1 || x === 2) {
        return 1;
    }
    return f(x - 1) + f(x - 2)
}

```

Show moreShow more icon

在按照 [AssemblyScript 提供的安装教程](https://github.com/AssemblyScript/assemblyscript#installation) 成功安装后， 再通过

```
asc f.ts -o f.wasm

```

Show moreShow more icon

就能把以上代码编译成可运行的 WebAssembly 模块。

为了加载并执行编译出的 f.wasm 模块，需要通过 JS 去加载并调用模块上的 f 函数，为此需要以下 JS 代码：

```
fetch('f.wasm') // 网络加载 f.wasm 文件
    .then(res => res.arrayBuffer()) // 转成 ArrayBuffer
    .then(WebAssembly.instantiate) // 编译为当前 CPU 架构的机器码 + 实例化
    .then(mod => { // 调用模块实例上的 f 函数计算
    console.log(mod.instance.f(50));
    });

```

Show moreShow more icon

以上代码中出现了一个新的内置类型 i32，这是 AssemblyScript 在 TypeScript 的基础上内置的类型。 AssemblyScript 和 TypeScript 有细微区别，AssemblyScript 是 TypeScript 的子集，为了方便编译成 WebAssembly 在 TypeScript 的基础上加了更严格的 [类型限制](https://github.com/AssemblyScript/assemblyscript/wiki/Limitations) ，区别如下：

- 比 TypeScript 多了很多更细致的内置类型，以优化性能和内存占用，详情 [文档](https://github.com/AssemblyScript/assemblyscript/wiki/Types);
- 不能使用 any 和 undefined 类型，以及枚举类型；
- 可空类型的变量必须是引用类型，而不能是基本数据类型如 string、number、boolean；
- 函数中的可选参数必须提供默认值，函数必须有返回类型，无返回值的函数返回类型需要是 void；
- 不能使用 JS 环境中的内置函数，只能使用 [AssemblyScript 提供的内置函数](https://github.com/AssemblyScript/assemblyscript/wiki/Built-in-functions) 。

总体来说 AssemblyScript 比 TypeScript 又多了很多限制，编写起来会觉得局限性很大； 用 AssemblyScript 来写 WebAssembly 经常会出现 tsc 编译通过但运行 WebAssembly 时出错的情况，这很可能就是你没有遵守以上限制导致的；但 AssemblyScript 通过修改 TypeScript 编译器默认配置能在编译阶段找出大多错误。

AssemblyScript 的实现原理其实也借助了 LLVM，它通过 TypeScript 编译器把 TS 源码解析成 AST，再把 AST 翻译成 IR，再通过 LLVM 编译成 WebAssembly 字节码实现； 上面提到的各种限制都是为了方便把 AST 转换成 LLVM IR。

### 为什么选 AssemblyScript 作为 WebAssembly 开发语言

AssemblyScript 相对于 C、Rust 等其它语言去写 WebAssembly 而言，好处除了对前端来说无额外新语言学习成本外，还有对于不支持 WebAssembly 的浏览器，可以通过 TypeScript 编译器编译成可正常执行的 JS 代码，从而实现从 JS 到 WebAssembly 的平滑迁移。

### 接入 Webpack 构建

任何新的 Web 开发技术都少不了构建流程，为了提供一套流畅的 WebAssembly 开发流程，接下来介绍接入 Webpack 具体步骤。

1. 安装以下依赖，以便让 TS 源码被 AssemblyScript 编译成 WebAssembly。

```
{
"devDependencies": {
    "assemblyscript": "github:AssemblyScript/assemblyscript",
    "assemblyscript-typescript-loader": "^1.3.2",
    "typescript": "^2.8.1",
    "webpack": "^3.10.0",
    "webpack-dev-server": "^2.10.1"
}
}

```

Show moreShow more icon

1. 修改 webpack.config.js，加入 loader：

```
module.exports = {
    module: {
        rules: [
            {
                test: /\.ts$/,
                loader: 'assemblyscript-typescript-loader',
                options: {
                    sourceMap: true,
                }
            }
        ]
    },
};

```

Show moreShow more icon

1. 修改 TypeScript 编译器配置 tsconfig.json，以便让 TypeScript 编译器能支持 AssemblyScript 中引入的内置类型和函数。

```
{
"extends": "../../node_modules/assemblyscript/std/portable.json",
"include": [
    "./**/*.ts"
]
}

```

Show moreShow more icon

1. 配置直接继承自 assemblyscript 内置的配置文件。

## WebAssembly 相关文件格式

前面提到了 WebAssembly 的二进制文件格式 wasm，这种格式的文件人眼无法阅读，为了阅读 WebAssembly 文件的逻辑，还有一种文本格式叫 wast； 以前面讲到的计算斐波那契序列的模块为例，对应的 wast 文件如下：

```
func $src/asm/module/f (param f64) (result f64)
(local i32)
get_local 0
f64.const 1
f64.eq
tee_local 1
if i32
    get_local 1
else
    get_local 0
    f64.const 2
    f64.eq
end
i32.const 1
i32.and
if
    f64.const 1
    return
end
get_local 0
f64.const 1
f64.sub
call 0
get_local 0
f64.const 2
f64.sub
call 0
f64.add
end

```

Show moreShow more icon

这和汇编语言非常像，里面的 f64 是数据类型，f64.eq f64.sub f64.add 则是 CPU 指令。

为了把二进制文件格式 wasm 转换成人眼可见的 wast 文本，需要安装 WebAssembly 二进制工具箱 [WABT](https://github.com/WebAssembly/wabt) ， 在 Mac 系统下可通过 brew install WABT 安装，安装成功后可以通过命令 wasm2wast f.wasm 获得 wast；除此之外还可以通过 wast2wasm f.wast -o f.wasm 逆向转换回去。

## WebAssembly 相关工具

除了前面提到的 WebAssembly 二进制工具箱，WebAssembly 社区还有以下常用工具：

- [Emscripten](http://kripken.github.io/emscripten-site/): 能把 C、C++代码转换成 wasm、asm.js；
- [Binaryen](https://github.com/WebAssembly/binaryen): 提供更简洁的 IR，把 IR 转换成 wasm，并且提供 wasm 的编译时优化、wasm 虚拟机，wasm 压缩等功能，前面提到的 AssemblyScript 就是基于它。

## WebAssembly JS API

目前 WebAssembly 只能通过 JS 去加载和执行，但未来在浏览器中可以通过像加载 JS 那样 `<script src='f.wasm'></script>` 去加载和执行 WebAssembly，下面来详细介绍如何用 JS 调 WebAssembly。

JS 调 WebAssembly 分为 3 大步： **加载字节码 > 编译字节码 \> 实例化** ，获取到 WebAssembly 实例后就可以通过 JS 去调用了，以上 3 步具体的操作是：

1. 对于浏览器可以通过网络请求去加载字节码，对于 Nodejs 可以通过 fs 模块读取字节码文件；
2. 在获取到字节码后都需要转换成 ArrayBuffer 后才能被编译，通过 WebAssembly 通过的 JS API [WebAssembly.compile](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/WebAssembly/compile) 编译后会通过 Promise resolve 一个 [WebAssembly.Module](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/WebAssembly/Module) ，这个 module 是不能直接被调用的需要；
3. 在获取到 module 后需要通过 [WebAssembly.Instance](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/WebAssembly/Instance) API 去实例化 module，获取到 Instance 后就可以像使用 JS 模块一个调用了。

其中的第 2、3 步可以合并一步完成，前面提到的 [WebAssembly.instantiate](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/WebAssembly/instantiate) 就做了这两个事情。

```
WebAssembly.instantiate(bytes).then(mod=>{
mod.instance.f(50);
})

```

Show moreShow more icon

### WebAssembly 调 JS

之前的例子都是用 JS 去调用 WebAssembly 模块，但是在有些场景下可能需要在 WebAssembly 模块中调用浏览器 API，接下来介绍如何在 WebAssembly 中调用 JS。

WebAssembly.instantiate 函数支持第二个参数 WebAssembly.instantiate(bytes,importObject)，这个 importObject 参数的作用就是 JS 向 WebAssembly 传入 WebAssembly 中需要调用 JS 的 JS 模块。举个具体的例子，改造前面的计算斐波那契序列在 WebAssembly 中调用 Web 中的 window.alert 函数把计算结果弹出来，为此需要改造加载 WebAssembly 模块的 JS 代码：

```
WebAssembly.instantiate(bytes,{
window:{
    alert:window.alert
}
}).then(mod=>{
mod.instance.f(50);
})

```

Show moreShow more icon

对应的还需要修改 AssemblyScript 编写的源码：

```
// 声明从外部导入的模块类型
declare namespace window {
    export function alert(v: number): void;
}

function _f(x: number): number {
    if (x == 1 || x == 2) {
        return 1;
    }
    return _f(x - 1) + _f(x - 2)
}

export function f(x: number): void {
    // 直接调用 JS 模块
    window.alert(_f(x));
}

```

Show moreShow more icon

修改以上 AssemblyScript 源码后重新用 asc 通过命令 asc f.ts 编译后输出的 wast 文件比之前多了几行：

```
(import "window" "alert" (func $src/asm/module/window.alert (type 0)))

(func $src/asm/module/f (type 0) (param f64)
    get_local 0
    call $src/asm/module/_f
    call $src/asm/module/window.alert)

```

Show moreShow more icon

多出的这部分 wast 代码就是在 AssemblyScript 中调用 JS 中传入的模块的逻辑。

除了以上常用的 API 外，WebAssembly 还提供一些 API，你可以通过这个 [d.ts 文件](https://github.com/01alchemist/webassembly-types/blob/master/webassembly.d.ts) 去查看所有 WebAssembly JS API 的细节。

## 不止于浏览器

WebAssembly 作为一种底层字节码，除了能在浏览器中运行外，还能在其它环境运行。

### 直接执行 wasm 二进制文件

前面提到的 Binaryen 提供了在命令行中直接执行 wasm 二进制文件的工具，在 Mac 系统下通过 brew install binaryen 安装成功后，通过 wasm-shell f.wasm 文件即可直接运行。

### 在 Node.js 中运行

目前 V8 JS 引擎已经添加了对 WebAssembly 的支持，Chrome 和 Node.js 都采用了 V8 作为引擎，因此 WebAssembly 也可以运行在 Node.js 环境中；

V8 JS 引擎在运行 WebAssembly 时，WebAssembly 和 JS 是在同一个虚拟机中执行，而不是 WebAssembly 在一个单独的虚拟机中运行，这样方便实现 JS 和 WebAssembly 之间的相互调用。

要让上面的例子在 Node.js 中运行，可以使用以下代码：

```
const fs = require('fs');

function toUint8Array(buf) {
    var u = new Uint8Array(buf.length);
    for (var i = 0; i < buf.length; ++i) {
        u[i] = buf[i];
    }
    return u;
}

function loadWebAssembly(filename, imports) {
    // 读取 wasm 文件，并转换成 byte 数组
    const buffer = toUint8Array(fs.readFileSync(filename));
    // 编译 wasm 字节码到机器码
    return WebAssembly.compile(buffer)
        .then(module => {
            // 实例化模块
            return new WebAssembly.Instance(module, imports)
        })
}

loadWebAssembly('../temp/assembly/module.wasm')
    .then(instance => {
        // 调用 f 函数计算
        console.log(instance.exports.f(10))
    });

```

Show moreShow more icon

在 Nodejs 环境中运行 WebAssembly 的意义其实不大，原因在于 Nodejs 支持运行原生模块，而原生模块的性能比 WebAssembly 要好。 如果你是通过 C、Rust 去编写 WebAssembly，你可以直接编译成 Nodejs 可以调用的原生模块。

## WebAssembly 展望

从上面的内容可见 WebAssembly 主要是为了解决 JS 的性能瓶颈，也就是说 WebAssembly 适合用于需要大量计算的场景，例如：

- 在浏览器中处理音视频， [flv.js](https://github.com/Bilibili/flv.js/) 用 WebAssembly 重写后性能会有很大提升；
- React 的 dom diff 中涉及到大量计算，用 WebAssembly 重写 React 核心模块能提升性能。Safari 浏览器使用的 JS 引擎 JavaScriptCore 也已经支持 WebAssembly，RN 应用性能也能提升；
- 突破大型 3D 网页游戏性能瓶颈， [白鹭引擎已经开始探索用 WebAssembly](https://feday.fequan.com/2017/WebAssembly在白鹭引擎的实践.pdf) 。

## 结束语

WebAssembly 标准虽然已经定稿并且得到主流浏览器的实现，但目前还存在以下问题：

- 浏览器兼容性不好，只有最新版本的浏览器支持，并且不同的浏览器对 JS WebAssembly 互调的 API 支持不一致；
- 生态工具不完善不成熟，目前还不能找到一门体验流畅的编写 WebAssembly 的语言，都还处于起步阶段；
- 学习资料太少，还需要更多的人去探索去踩坑。

总之现在的 WebAssembly 还不算成熟，如果你的团队没有不可容忍的性能问题，那现在使用 WebAssembly 到产品中还不是时候， 因为这可能会影响到团队的开发效率，或者遇到无法轻易解决的坑而阻塞开发。