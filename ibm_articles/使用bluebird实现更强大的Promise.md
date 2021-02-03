# 使用 bluebird 实现更强大的 Promise
JavaScript 应用的首选实现

**标签:** JavaScript,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/wa-lo-use-bluebird-implements-power-promise/)

成富

发布: 2018-01-22

* * *

Promise 是 JavaScript 开发中的一个重要概念。按照 Promises/A+ 规范的定义，Promise 表示的是一个异步操作的最终结果。与 Promise 进行交互的方式是通过其 then(onFulfilled, onRejected) 方法，用来注册处理 Promise 最终结果的回调方法。Promise 所表示的异步操作可能成功或失败。当成功时，回调方法 onFulfilled 会接受到一个值作为最终结果；当失败时，回调方法 onRejected 会接受到操作失败的原因。一个 Promise 可能处于三种状态之中：进行中（pending）表示 Promise 所对应的异步操作还在进行中；已满足（fulfilled）表示异步操作已经成功完成；已拒绝（rejected）表示异步操作无法完成。Promise 对象的 then 方法接受两个参数，分别是已满足状态和已拒绝状态的回调方法。已满足状态的回调方法的第一个参数是 Promise 的最终结果值；已拒绝状态的回调方法的第一个参数是 Promise 被拒绝的原因。

Promise 是 ECMAScript 6 规范的一部分，已经有很多的第三方库提供了与 Promise 相关的实现。Promises/A+ 规范只是定义了 then 方法，其目的是为了不同实现库之间的互操作性，侧重的是对于 Promise 的一般使用场景。ECMAScript 6 规范中还额外定义了其他 Promise 的方法，包括 Promise.all、Promise.race、Promise.reject 和 Promise.resolve 等。在 NodeJS 和浏览器环境中可以使用不同的 Promise 实现。本文中介绍的 Bluebird 是一个功能丰富而且性能优异的 Promise 实现。

根据不同的运行环境，可以有不同的安装方式。在 NodeJS 应用中，直接使用 npm install bluebird 就可以安装。在浏览器应用中，可以直接下载发布版本的 JavaScript 文件或使用 Bower 来安装。本文的实例基于 NodeJS 环境，对应版本是 LTS 8.9.4 版本，代码使用 ECMAScript 6 的语法。在 NodeJS 环境中，通过 const Promise = require(‘bluebird’) 就可以开始使用 Bluebird 提供的 Promise 对象。

## Promise 对象创建

在使用 Bluebird 作为 Promise 实现时的第一个问题是如何创建 Promise 对象。实际上，在大部分时候都不需要显式的创建 Promise。很多第三方库本身的方法返回的就已经是 Promise 对象或是包含了 then 方法的 thenable 对象，可以直接使用。Bluebird 的一个实用功能是把不使用 Promise 的已有API包装成返回 Promise 的新 API。大部分 NodeJS 的标准库 API 和不少第三方库的 API 都使用了回调方法的模式，也就是在执行异步操作时，需要传入一个回调方法来接受操作的执行结果和可能出现的错误。对于这样的方法，Bluebird 可以很容易的将它们转换成使用 Promise 的形式。

比如 NodeJS 中读取文件的 fs.readFile 方法，其常见的使用模式如 代码清单 1 所示。使用回调方法的问题在于可能导致过多的代码嵌套层次，造成所谓的回调地狱问题（callback hell）。

#### 清单 1\. NodeJS 中的 fs.readFile 方法的基本使用方式

```
const fs = require('fs'),
path = require('path');

fs.readFile(path.join(__dirname, 'sample.txt'), 'utf-8', (err, data) => {
if (err) {
console.error(err);
} else {
console.log(data);
}
});

```

Show moreShow more icon

Bluebird 的 Promise.promisifyAll 方法可以为一个对象的属性中的所有方法创建一个对应的使用 Promise 的版本。这些新创建方法的名称在已有方法的名称后加上”Async”后缀。代码清单1中的实现可以改成代码清单2中的使用 Promise 的形式。新的方法 readFileAsync 对应的是已有的 readFile 方法，但是返回值是 Promise 对象。除了 readFile，fs 中的其他方法也都有了对应的 Async 版本，如 writeFileAsync 和 fstatAsync 等。

#### 清单 2\. 使用 Promise.promisifyAll 来转换方法

```
const Promise = require('bluebird'),
fs = require('fs'),
path = require('path');
Promise.promisifyAll(fs);

fs.readFileAsync(path.join(__dirname, 'sample.txt'), 'utf-8')
.then(data => console.log(data))
.catch(err => console.error(err));

```

Show moreShow more icon

如果不希望把一个对象的所有方法都自动转换成使用 Promise 的形式，可以使用 Promise.promisify 来转换单个方法，如 Promise.promisify(require(“fs”).readFile)。对于一个 NodeJS 格式的回调方法，可以使用 Promise.fromCallback 将其转换成一个 Promise。回调方法的结果决定了 Promise 的状态。

对于一个已有值，可以使用 Promise.resolve 方法将其转换成一个状态为已满足的 Promise 对象。同样的，使用 Promise.reject 方法可以从给定的拒绝原因中创建一个状态为已拒绝的 Promise 对象。这两个方法都可以快速地创建 Promise 对象。

如果上述这些方法都不能满足创建 Promise 的需求，可以使用 Promise 构造方法。构造方法的参数是一个接受两个方法作为参数的方法，这两个方法分别用来把 Promise 标记为已满足或已拒绝。在代码清单3中，创建的 Promise 的状态取决于生成的随机数是否大于 0.5。

#### 清单 3\. 创建 Promise

```
const Promise = require('bluebird');
const promise = new Promise((resolve, reject) => {
const value = Math.random();
if (value > 0.5) {
    resolve(value);
} else {
    reject(`Invalid value ${value}`);
}
});
promise.then(console.log).catch(console.error);

```

Show moreShow more icon

## 使用 Promise

在之前的代码示例中，我们已经看到了 Promise 对象的基本用法。其中 then 方法用来添加处理 Promise 结果的回调方法。可以同时添加已满足状态和已拒绝状态的回调方法。spread 方法的作用与 then 相似。当 Promise 的结果值为数组时，可以使用 spread 方法把数组的元素打平（flatten）为回调方法的不同参数。在代码清单4中，Promise 的值是一个包含3个值的数组，分别对应处理方法的3个参数。

#### 清单 4\. spread 使用示例

```
Promise.resolve([1, 2, 3])
.spread((v1, v2, v3) => console.log(v1 + v2 + v3));

```

Show moreShow more icon

catch 方法用来添加已拒绝状态时的回调方法。catch 方法有两种使用方式。第一种方式只添加回调方法，会捕获所有错误情况；第二种方式是使用错误对象类型或断言（predicate）来对错误进行过滤。error 方法的作用与 catch 类似，但是 error 只会处理真正的 Error 对象。JavaScript 中的 throw 语句是可以抛出任何值的，不一定是 Error 对象。throw 抛出的任何值都会被 catch 处理，但是只有 Error 对象才会被 error 处理。在代码清单5中，error 添加的处理方法并不会被调用。

#### 清单 5\. error 使用示例

```
Promise.reject('not an Error')
.error(err => console.log('should not appear'));

```

Show moreShow more icon

finally 方法添加的回调方法不管 Promise 的最终状态如何，都会被调用。 [代码清单6](_清单6catch和finally使用示例) 中的 catch 只会捕获 TypeError，而 finally 中的逻辑始终会被调用。

#### 清单 6\. catch 和 finally 使用示例

```
Promise.reject(new TypeError('some error'))
.catch(TypeError, console.error)
.finally(() => console.log('done'));

```

Show moreShow more icon

对于一个 Promise 对象，可以使用其提供的方法来查看其内部状态。

- isPending 检查 Promise 的状态是否为进行中。
- isFulfilled 检查 Promise 的状态是否为已满足。
- isRejected 检查 Promise 的状态是否为已拒绝。
- isCancelled 检查 Promise 是否被取消。
- value 获取 Promise 满足后的结果。
- reason 获取 Promise 被拒绝的原因。

可以使用 cancel 方法来取消一个 Promise，表示不再对 Promise 的结果感兴趣。当一个 Promise 被取消之后，其回调方法都不会被调用。取消功能默认是禁用的，需要使用 Promise.config 来启用该功能。需要注意的是，取消 Promise 只是表示 Promise 的回调方法不被调用，并不会自动取消正在进行的异步操作，比如不会自动终止进行中的 XHR 请求。如果需要添加自定义的取消逻辑，代码清单3中的 Promise 构造方法的唯一参数所对应的方法可以添加第三个参数 onCancel，用来注册 Promise 取消时的回调方法。

#### 清单 7\. 使用 Promise 的取消功能

```
Promise.config({
cancellation: true,
});

Promise.delay(1000, 'hello').then(console.log).cancel();

```

Show moreShow more icon

## 处理集合

之前的代码示例都针对单个 Promise。在实际中，经常会需要与多个 Promise 进行交互。比如，读取单个文件的操作返回的是 Promise 对象，而我们需要等待若干个文件的读取操作都成功之后，再执行其他操作。Promise.all 可以接受一个 Iterable 参数，并返回一个新的 Promise。结果的 Promise 只有在 Iterable 中包含的全部 Promise 对象都满足时，才会满足；任何 Promise 出现错误会导致结果 Promise 被拒绝。新的 Promise 的最终结果是一个包含 Iterable 中 Promise 的结果的数组。该 Iterable 中可以包含任何值。如果值不是 Promise，则其值会直接出现在结果中，并不需要等待 Promise 完成。在代码清单8中，我们使用 Promise.all 来等待3个读取文件操作完成。

#### 清单 8\. Promise.all 使用示例

```
Promise.all([
fs.readFileAsync(path.join(__dirname, '1.txt'), 'utf-8'),
fs.readFileAsync(path.join(__dirname, '2.txt'), 'utf-8'),
fs.readFileAsync(path.join(__dirname, '3.txt'), 'utf-8')
]).then(results => console.log(results.join(', '))).catch(console.error);

```

Show moreShow more icon

在很多情况下的做法都需要把 Iterable 中的对象转换成 Promise 之后，再等待这些 Promise 完成。对于这样的场景，可以使用 Promise.map。map 方法的第一个参数是 Iterable 对象，第二个参数是转换为 Promise 的方法，第三个参数是可选的配置对象，可以使用属性 concurrency 来控制同时运行的 Promise 的数量。在代码清单9代码清单9中，我们使用 Promise.map 把包含文件名称的数组转换成读取文件内容的 Promise 对象，再等待读取操作完成。的功能与代码清单8一样，但是实现更简单。

#### 清单 9\. Promise.map 的使用示例

```
Promise.map(['1.txt', '2.txt', '3.txt'],
name => fs.readFileAsync(path.join(__dirname, name), 'utf-8')
).then(results => console.log(results.join(', '))).catch(console.error);

```

Show moreShow more icon

Promise.mapSeries 的作用与 Promise.map 相同，只不过 mapSeries 按照 Iterable 中的顺序依次遍历每个元素。

在有多个 Promise 时，如果只需要等待其中部分 Promise 完成，可以使用 Promise.some 并指定完成的 Promise 数量。在代码清单10中，只需要等待2个文件读取操作完成即可。

#### 清单 10\. Promise.some 的使用示例

```
Promise.some([
fs.readFileAsync(path.join(__dirname, '1.txt'), 'utf-8'),
fs.readFileAsync(path.join(__dirname, '2.txt'), 'utf-8'),
fs.readFileAsync(path.join(__dirname, '3.txt'), 'utf-8')
], 2).then(results => console.log(results.join(', '))).catch(console.error);

```

Show moreShow more icon

Promise.any 相当于使用 Promise.some 并把数量设为1，不过 Promise.any 的结果不是一个长度为1的数组，还是具体的单个值。Promise.race 的作用类似 any，不过 race 的结果有可能是被拒绝的 Promise。因此推荐的做法是使用 any。

Promise.filter 可以等待多个 Promise 的完成，并对结果进行过滤。它实际上的效果相当于在 Promise.map 之后使用 Array 的 filter 方法来进行过滤。在代码清单11中，Promise.filter 用来过滤掉内容长度不大于1的文件。

#### 清单 11\. Promise.filter 的使用示例

```
Promise.filter([
fs.readFileAsync(path.join(__dirname, '1.txt'), 'utf-8'),
fs.readFileAsync(path.join(__dirname, '2.txt'), 'utf-8'),
fs.readFileAsync(path.join(__dirname, '3.txt'), 'utf-8')
], value => value.length > 1).then(results => console.log(results.join(', '))).catch(console.error);

```

Show moreShow more icon

Promise.each 对 Iterable 中的元素进行依次处理。处理方法也可以返回 Promise 对象来表示异步处理。只有对上一个元素的处理完成之后，才会处理下一个元素。

Promise.reduce 把多个 Promise 的结果缩减成单个值。其作用类似 Array 的 reduce 方法，但是可以处理 Promise 对象。在代码清单12中，第一个参数是待处理的文件名称数组；第二个参数是进行累积操作的方法，total 表示当前累积值，用来与读取到的当前文件的长度相加；第三个参数是累积的初始值。从代码中可以看到，累积操作也可以返回 Promise 对象。Promise.reduce 会等待 Promise 完成。

#### 清单 12\. Promise.reduce 的使用示例

```
Promise.reduce(['1.txt', '2.txt', '3.txt'],
(total, name) => {
return fs.readFileAsync(path.join(__dirname, name), 'utf-8').then(data => total + data.length);
}
, 0).then(result => console.log(`Total size: ${result}`)).catch(console.error);

```

Show moreShow more icon

Promise.all 用来处理动态数量的 Promise，Promise.join 用来处理固定数量的不相关的 Promise。

Promise.method 用来封装一个方法，使其返回 Promise。Promise.try 用来封装可能产生问题的方法调用，并根据调用结果来决定 Promise 的状态。

## 资源使用

如果在 Promise 的代码中使用了需要释放的资源，如数据库连接，确保这些资源被释放并不是一件容易的事情。一般的做法是在 finally 中添加资源释放逻辑。但是当 Promise 互相嵌套时，很容易产生资源没有被释放的情况。Bluebird 提供了更加健壮的资源管理方式，即使用资源释放器（disposer）和 Promise.using 方法。

资源释放器以 Disposer 对象表示，通过 Promise 的 disposer 方法来创建。创建时的参数是一个用来释放资源的方法。该方法的第一个参数是资源对象，第二个参数是 Promise.using 产生的 Promise 对象。Disposer 对象可以传递给 Promise.using 来保证其资源释放逻辑被执行。

在代码清单13中，Connection 表示数据库连接，DB 表示数据库。DB 的 connect 方法创建一个数据库连接，返回的是 Promise 对象。Connection 的 query 方法表示执行数据库查询，返回结果也是一个 Promise 对象。Connection 的 close 方法关闭数据库连接。在使用时，在 connect 方法返回的 Promise 对象之上使用 disposer 来创建一个 Disposer 对象，所对应的资源释放逻辑是调用 close 方法。接着就可以通过 Promise.using 来使用 Connection 对象。

#### 清单 13\. Promise 的资源管理

```
const Promise = require('bluebird');

class Connection {
query() {
return Promise.resolve('query');
}
close() {
console.log('close');
}
}

class DB {
connect() {
return Promise.resolve(new Connection());
}
}

const disposer = new DB().connect().disposer(connection => connection.close());
Promise.using(disposer, connection => connection.query())
.then(console.log).catch(console.error);

```

Show moreShow more icon

## 定时器

Promise.delay 返回的 Promise 会在指定的时间之后转换到已满足状态，并使用指定的值作为结果。Promise.timeout 用来为一个已有的 Promise 添加超时功能。如果已有的 Promise 没有在指定的超时时间内完成，结果 Promise 会产生 TimeoutError 错误。

在代码清单14中，第一个 delay 方法会在1秒后输出 hello。第二个方法会抛出 TimeoutError 错误，因为 delay 的延迟时间1秒超过了超时时间500毫秒。

#### 清单 14\. Promise.delay 和 Promise.timeout 的使用示例

```
Promise.delay(1000, 'hello').then(console.log);

Promise.delay(1000, 'hello')
.timeout(500, 'timed out')
.then(console.log)
.catch(console.error);

```

Show moreShow more icon

## 实用方法

Promise 中还包含了一些实用方法。tap 和 tapCatch 分别用来查看 Promise 中的结果和出现的错误。这两个方法中的处理方法不会影响 Promise 的结果，适合用来执行日志记录。call 用来调用 Promise 结果对象中的方法。get 用来获取 Promise 结果对象中的属性值。return 用来改变 Promise 的结果。throw 用来抛出错误。catchReturn 用来在捕获错误之后，改变 Promise 的值。catchThrow 用来在捕获错误之后，抛出新的错误。

代码清单15给出了这些实用方法的使用示例。第一个语句使用 tap 来输出 Promise 的结果。第二个语句使用 call 来调用 sayHi 方法，并把该方法的返回值作为 Promise 的结果。第三个语句使用 get 来获取属性b的值，并作为 Promise 的结果。第四个语句使用 return 来改变了 Promise 的结果。第五个语句使用 throw 来抛出错误，再用 catch 来处理错误。第六个语句使用 catchReturn 来处理错误并改变 Promise 的值。

#### 清单 15\. 实用方法的使用示例

```
Promise.resolve(1)
.tap(console.log)
.then(v => v + 1)
.then(console.log);

Promise.resolve({
sayHi: () => 'hi',
}).call('sayHi').then(console.log);

Promise.resolve({
a: 'hello',
b: 1,
}).get('b').then(console.log);

Promise.resolve(1)
.return(2)
.then(console.log);

Promise.resolve(1)
.throw(new TypeError('type error'))
.catch(console.error);

Promise.reject(new TypeError('type error'))
.catchReturn('default value')
.then(console.log);

```

Show moreShow more icon

## 错误处理

虽然可以使用 Promise 的 catch 来捕获错误并进行处理，在很多情况下，由于程序本身或第三方库的问题，仍然可能出现错误没有被捕获的情况。未捕获的错误可能导致 NodeJS 进程意外退出。Bluebird 提供了全局和本地两种错误处理机制。Bluebird 会触发与 Promise 被拒绝相关的全局事件，分别是 unhandledRejection 和 rejectionHandled。可以添加全局的事件处理器来处理这两种事件。对于每个 Promise，可以使用 onPossiblyUnhandledRejection 方法来添加对于未处理的拒绝错误的默认处理方法。如果不关注未处理的错误，可以使用 suppressUnhandledRejections 来忽略错误。

在代码清单16中，Promise 的拒绝错误未被处理，由全局的处理器来进行处理。

#### 清单 16\. 错误处理示例

```
Promise.reject(new TypeError('error'))
.then(console.log);

process.on('unhandledRejection', (reason, promise) => console.error(`unhandled ${reason}`));

```

Show moreShow more icon

## 结束语

Promise 是 JavaScript 应用开发中重要的异步操作相关的抽象。Bluebird 作为流行的 Promise 库，提供了很多与 Promise 相关的实用功能。本文对 Bluebird 的重要特性做了详细的介绍，包括如何创建 Promise 对象、使用 Promise、处理集合、Promise 实用方法和错误处理等。