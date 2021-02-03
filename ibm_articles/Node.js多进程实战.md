# Node.js 多进程实战
Node.js 进阶

**标签:** Node.js

[原文链接](https://developer.ibm.com/zh/articles/os-cn-nodejs-practice/)

郭晋

发布: 2016-08-15

* * *

## Node.js 多进程实战

Node.js 以其天生的处理高并发 I/O 的强大能力闻名于世，我们选用 Node.js 也大多是看上了其这一特性。然而无论是随着开发的逐步深入还是产品在生产环境上的部署，往往会沮丧的发现，高并发 I/O 并不总是生活的全部，CPU 密集型计算依然或多或少的不可避免。

怎么办？无疑只有两条路可以走：

1. 忽视它，也许 Node.js 本身就能帮我们处理好呢
2. 正视它，Java 里有多线程，Node 里面应该也能找到类似办法处理

## 当单线程遇到 CPU

首先看看第一条路能否走通，如果可以则皆大欢喜，毕竟这样可以不改变原有的代码逻辑和结构。来看下面的这段代码：

##### 清单 1\. 单线程示例

```
var Promise = require('bluebird');
var fiboTasks = [44, 42, 42, 43];
function fibo(n) {
return n == 0 ? 0 : n > 1 ? fibo(n - 1) + fibo(n - 2) : 1;
}
function excuteFibo(seq, taskID) {
return new Promise((resolve) => {
setTimeout(() => {
var st = Date.now();
var result = fibo(seq);
console.log(`Task ${taskID} was complete and using ${Date.now() -
     st} ms`);
resolve(result);
}, Math.random()*10);
});
}
var st = Date.now();
Promise.map(fiboTasks, function (item, index) { return excuteFibo(item,
        index) }).then(function (result) {
console.log(`All tasks were complete and using ${Date.now() - st} ms`);
})

```

Show moreShow more icon

首先，为了让 CPU 忙起来，这里用了大家喜闻乐见的斐波那契数列函数，用来来模拟一个消耗 CPU 的应用。

其次，使用 setTimeout() 方法使得 Fibo 函数不立即执行，而且交由 Node.js 系统来调度。

最后，使用了 Promise.map() 方法使得多次调用 Fibo 函数的行为得以并发执行。

考虑到一些系统开销，我们期望所有 task 的总耗时应该略长于耗时最长的 task，而结果如何呢？

##### 清单 2\. 单线程示例结果

```
c:\work\NodeClusterTest>node concurrentTest.js
Task 0 was complete 701408733 and using 11971 ms
Task 3 was complete 433494437 and using 7272 ms
Task 1 was complete 267914296 and using 4477 ms
Task 2 was complete 267914296 and using 4573 ms
All tasks were complete and using 28308 ms

```

Show moreShow more icon

从日志信息中可以看出，所有 task 的总耗时（28308 ms）基本上等于每个 task 执行时间的总和（28293 ms），这明显和我们的期望（~ 11971 ms）不符。果然，并没有什么奇迹发生，单线程只能依次的执行 CPU 计算。

## The Way Out – 多进程初探

在 Node.js 中处理 CPU 密集型计算一般有三种方案：

1. 写独立的 c 代码
2. 使用 Node.js 自带的 cluster 模块
3. 使用其他开源项目，如 threads-a-gogo

首先我们一般不会选 1，除了 coding 的难度外，还增加了编译和部署的复杂度，与我们选用 Node.js 的初衷相违背。

方法 2 中的 cluster 模块提供的是多进程的解决方案，而方法 3 则采用多线程的方式，我们在这里不去比较两种方案孰优孰劣，仅从使用的便捷性（不依赖第三方库）以及维护成本（有 Node.js 团队维护和持续开发）考虑，本文选择了方法 2，即使用 cluster 模块。

cluster 模块的原理，Node.js 官方文档里有详细的介绍，这里就不赘述了。这里先将之前的代码改成一个最简单多进程的方式实现：

##### 清单 3\. 多进程示例

```
console.log('####====START====###');
const cluster = require('cluster');
const numCPUs = require('os').cpus().length;
function fibo(n) {...}

if (cluster.isMaster) {
var collection = [44, 42, 42, 43];
var st = Date.now();
for (var i = 0; i < Math.min(numCPUs, collection.length); i++) {
var wk = cluster.fork();
wk.send(collection[i]);
}
cluster.on('fork', function (worker) {
console.log(`[master] : fork worker ${worker.id}`);
});
cluster.on('exit', function (worker, code, signal) {
console.log(`[master] : worker ${worker.id} died`);
});
var numOfCompelete = 0
Object.keys(cluster.workers).forEach(function (id) {
cluster.workers[id].on('message', function (msg) {
console.log(`[master] receive message from [worker ${id}]: ${msg}`);
numOfCompelete++;
if (numOfCompelete === collection.length) {
console.log(`[master] finish all work and using ${Date.now() -
st} ms`);
cluster.disconnect();
}
});
});

} else {
process.on('message', function (msg) {
var st = Date.now();
var result = fibo(msg);
console.log(`[worker ${cluster.worker.id}] finish work and using
${Date.now() - st} ms`);
process.send(result);
});
}

```

Show moreShow more icon

可以看到，整段代码分为两块，master 段和 worker 段。在 master 段里主要做了如下几件事：

1. 建立子进程
2. master 注册了”fork”, “exit”,”message” 事件以便在事件发生时做出相应的操作。这里最重要的是”message”事件，当 worker 也就是子进程返回计算结果时会被触发，而 master 进程需要将各个子进程的结果汇总以便后续处理。

在 worker 段比较简单，主要有如下操作：

1. 注册了”message”用以接收 master 进程传来的内容
2. 执行具体计算并将结果返回 master 进程

下面来看具体执行的结果：

##### 清单 4\. 多进程示例结果

```
C:\work\NodeClusterTest> node test.js
####====START====###
[master] : fork worker 1
[master] : fork worker 2
[master] : fork worker 3
[master] : fork worker 4
####====START====###
####====START====###
####====START====###
####====START====###
[worker 2] finish work and using 9284 ms
[master] receive message from [worker 2]: 267914296
[worker 3] finish work and using 9456 ms
[master] receive message from [worker 3]: 267914296
[worker 4] finish work and using 13030 ms
[master] receive message from [worker 4]: 433494437
[worker 1] finish work and using 17632 ms
[master] receive message from [worker 1]: 701408733
[master] finish all work and using 18099 ms
[master] : worker 4 died
[master] : worker 3 died
[master] : worker 2 died
[master] : worker 1 died

```

Show moreShow more icon

从日志信息中可以看出，使用了多进程之后，所有 task 的总耗时为 18099 ms，远小于之前的 28308 ms，且略长于耗时最长的 task 17632 ms，这和我们最初的预期一致，我们的目标似乎已经达到了？

如果细心的话，我们会发现，日志中显示程序的开始位置被打印了 4 遍。也就是说，master 进程 fork 的每一个 worker 子进程都会将整段代码执行一遍。在真正生产环境中是绝不会允许这种情况出现的。因此，此代码仅仅能说是实现了功能，还需要继续的优化才能进入真正的生产环境。

## 从 Demo 程序到生产环境有

首先，我们需要将 worker 段的代码抽离出来，以避免 worker 子进程将 master 进程的内容重新执行一遍，例如，我们可以生成一个 worker.js 的文件：

##### 清单 5.worker.js 代码片段

```
var cluster = require('cluster');
function fibo(n) {...}
console.log(`[worker ${cluster.worker.id}] start ...` );
process.on('message', function(msg) {
var st = Date.now();
console.log(`[worker ${cluster.worker.id}] start to work`);
var result = worker.fibo(msg);
console.log(`[worker ${cluster.worker.id}] work finish work and using ${Date.now() - st} ms`);
process.send(result);
});

```

Show moreShow more icon

可以看出，worker.js 中除了多包含了计算逻辑 fibo 函数外，和 清单 3 中的 worker 段的代码基本相同。

接下来，来修改 master 段的代码，同样的，可以新建一个 master.js 的文件：

##### 清单 6.master.js 代码片段

```
console.log('####====START====###');
var cluster = require('cluster');

cluster.setupMaster({
exec: 'worker.js',
slient: true
});

// 剩余代码和清单 3 中 isMaster 段相同
......

```

Show moreShow more icon

在 master.js 中，绝大多数代码和清单 3 中 cluster.isMaster 段中的相同。唯一的区别是使用了 cluster.setupMaster() 方法，并指定了 worker 进程的执行路径为 worer.js。有关 setupMaster() 的详细说明，请参考 Node.js 官方文档。

修改过后的代码执行结果如下：

##### 清单 7.master.js 代码执行结果

```
c:\work\NodeClusterTest>node master1.js
####====START====###
[master] : fork worker 1
[master] : fork worker 2
[master] : fork worker 3
[master] : fork worker 4
[worker 1] start ...
......

```

Show moreShow more icon

很显然，经过修改后的代码已经不会再出现 worker 进程重复执行 master 段代码的情况了，这点很重要。接下来需要做的是将这段多进程的代码和主程序做集成，因为我们最终是需对将多线程计算的结果进行后续操作。我们在这里选择将整个 master 段代码做为一个 Promise 返回给主程序，这样做可以最大限度的保证主程序的逻辑结构不被破坏。

下面来继续修改 master.js 文件：

##### 清单 8.master.js 代码片段

```
module.exports = exuteFibo;
function exuteFibo () {
return (new Promise(
function (reslove, reject) {
var cluster = require('cluster');
var result = [];
cluster.setupMaster({
exec: 'worker.js',
slient: true
});
// 以下清单 3 中 isMaster 段相同
......
Object.keys(cluster.workers).forEach(function (id) {
cluster.workers[id].on('message', function (msg) {
console.log(`[master] receive message from [worker${id}]:
${msg}`);
result.push(msg);
numOfCompelete++;
if (numOfCompelete === collection.length) {
console.log(`[master] finish all work and using
${Date.now() - st} ms`);
cluster.disconnect();
reslove(result);
}
});
});
}
))}

```

Show moreShow more icon

如清单 5 中所示，我们将 master 进程所做的操作封装成了一个 Promise 并返回，这样主程序就可以像调用其他所有返回类型为 Promise 的函数一样来处理。接下来我们来写一个简单的主程序并命名为 main.js

##### 清单 9.main.js 代码片段

```
var exuteFibo = require('./master');
console.log('=====Start========');
var st = Date.now();
exuteFibo().then(function(result){
console.log(`Finish all work and using ${Date.now() - st}`);
console.log('####Get result from multiple-processes: '+ result);
});

```

Show moreShow more icon

同样的，还是来看一下具体的执行结果：

##### 清单 10.main.js 代码执行结果

```
PS C:\work\NodeClusterTest> node .\main.js
=====Start========
[master] : fork worker 1
......
[master] finish all work and using 17357 ms
Finish all work and using 17365
####Get result from multiple-processes: 433494437,433494437,701408733,1134903170
[master] : worker 4 died
[master] : worker 3 died
[master] : worker 1 died
[master] : worker 2 died

```

Show moreShow more icon

正如我们期待的，主程序通过调用 master.js 得到了由多线程计算的结果以便继续后续处理，一切看上去都已经很美好了，这篇文章也似乎可以结束了。等等，并发调用呢？做为 Nodejs 的程序怎么能忽视这个问题呢。

## 当多进程遇到并发

继续修改 main.js，让主程序并发调用 master.js 试试，看看会出现什么情况：

##### 清单 11.main.js 代码片段

```
var exuteFibo = require('./master');
console.log('=====Start========');
var st = Date.now();
exuteFibo().then(function(result){
console.log(`Finish all work and using ${Date.now() - st}`);
console.log('####Get result from multiple-processes: '+ result);
});

st = Date.now();
exuteFibo().then(function(result){
console.log(`Finish all work and using ${Date.now() - st}`);
console.log('####Get result1 from mutliple-processes: '+ result);
});

```

Show moreShow more icon

而运行的结果十分的令人沮丧，从日志的输出结果上来看是一团糟！

##### 清单 12.main.js 代码执行结果

```
PS C:\work\NodeClusterTest> node .\main.js
=====Start========
[master] : fork worker 1  ①
[master] : fork worker 1
......
[master] : fork worker 8
[master] : fork worker 8  ②

......
[worker 6] start to work
[worker 3] work finish work and using 17823 ms
[master] receive message from [worker3]: 433494437
......
[master] receive message from [worker6]: 433494437
[master] finish all work and using 19618 ms
Finish all work and using 19634
[master] : worker 3 died
[master] : worker 3 died
[master] : worker 7 died
[master] : worker 7 died ③
[master] : worker 6 died
[master] : worker 6 died
[master] : worker 2 died
[master] : worker 2 died
[worker 8] work finish work and using 24262 ms
[master] receive message from [worker8]: 701408733
[master] : worker 8 died
[master] : worker 8 died
[worker 4] work finish work and using 23825 ms
[master] receive message from [worker4]: 701408733
[master] receive message from [worker4]: 701408733
[master] : worker 4 died
[master] : worker 4 died
[worker 1] work finish work and using 28957 ms
[master] receive message from [worker1]: 1134903170
[master] finish all work and using 29816 ms
events.js:141
throw er; // Unhandled 'error' event ④
^
Error: write EPIPE
at exports._errnoException (util.js:874:11)

```

Show moreShow more icon

基本上，从日志中可以发现有四处很明显的问题：

1. 每一个 worker 进程似乎被 fork 了两次
2. Worker 的 id 最大为什么会是 8？我们的 master 进程可只 fork 了 4 个 work 进程
3. 为什么一个 master 进程结束后会杀死所有的 worker 进程
4. 什么导致了这个异常？

怎么解决呢？现在让我们逐一来看一看。

首先，我们修改一下代码，将 master 进程的 pid 加入日志信息，这样可以更加清楚 master 进程和 worker 进程直接的关系。修改过后的代码执行结果如下：

##### 清单 13.main.js 代码执行结果

```
PS C:\work\NodeClusterTest> node .\main.js
=====Start========
[master 8364] : fork worker 1
[master 8364] : fork worker 1
[master 8364] : fork worker 2
[master 8364] : fork worker 2
[master 8364] : fork worker 3
[master 8364] : fork worker 3
[master 8364] : fork worker 4
[master 8364] : fork worker 4
[master 8364] : fork worker 5
[master 8364] : fork worker 5
[master 8364] : fork worker 6
[master 8364] : fork worker 6
[master 8364] : fork worker 7
[master 8364] : fork worker 7
......

```

Show moreShow more icon

通过日志我们清楚的看到，Node.js 内核并没有因为我们并发调用 cluster 模块而创建多个 master 进程，相反的，它复用了已存在的 master 进程，类似于 singleton 模式。

而 worker 为什么会被 fork 两次呢？仔细阅读代码，会发现相关日志是由如下代码产生的

##### 清单 14\. 进程 fork 代码片段

```
cluster.on('fork', function (worker) {
console.log(`[master] : fork worker ${worker.id}`);
});

```

Show moreShow more icon

这是一个典型的 Nodejs 的事件注册方法，cluster 进程监听了”fork”事件，每当一个 worker 被成功 fork 后就会触发一个”fork”事件，而 cluster 进程收到”fork”事件后就会执行预定义好的操作。

在本例中，同一个 cluster 进程并且都监听了两次”fork”事件，因而当一个 worker 被 fork 后，cluster 进程都会触发两次相应的操作。

因此，worker 进程不是被 fork 了两遍而只是被记录了两遍。而由于只有一个 master 进程存在，也就不难解释为什么最大的 worker id 会是 8 而不是 4，进一步查看 Node.js 官方文档关于 worker.id 的描述也印证了这一点：

_“Each new worker is given its own unique id, this id is stored in the id. While a worker is alive, this is the key that indexes it in cluster.workers”_

这就是说每创建一个新的 worker 进程就会生成一个唯一的 id，并且这个 id 在整个 worker 进程的生命周期内都会被保存在 cluster.workers 中。而我们的 master 进程在结束 worker 进程使使用的命令如下：

##### 清单 15\. 进程关闭代码片段

```
cluster.disconnect();

```

Show moreShow more icon

有着如下说明：

_“Calls.disconnect() on each worker in cluster.workers”_

我们刚刚知道每个 worker id 在其整个生命周期内都保存在 cluster.workers 中，因而当我们调用 cluster.disconnect() 的时候，master 进程会尝试关闭当前 cluster.workers 中所有存在的 worker 进程。

至于那个抛出的异常，也就很容易解释了，这是由于一个 master 进程在关闭一个已经被另一个 master 进程关闭的 worker 进程的时候出现的错误。至此，我们已经清楚了诸多问题的原因了，该是到解决问题的时候了。经过之前的分析，能够知道问题的核心在于 master 进程错误的操作了未由它创建的子进程，解决此问题的大体思路就是给 master 进程创建一个数组用来保存自己创建的子进程的 id，master 进程只对数组中存在的 worker id 进行操作。关于这一点，其实 Node.js 官方文档里也有说明：

_“Node.js does not automatically manage the number of workers for you, however. It is your responsibility to manage the worker pool for your application’s needs “_

修改过的 master.js 为：

##### 清单 16.main.js 代码片段

```
module.exports = exuteFibo;
function exuteFibo() {
return (new Promise(
function (reslove, reject) {
var cluster = require('cluster');
......
var workerID = [];
// Fork workers.
for (var i = 0; i < Math.min(numCPUs, collection.length); i++) {
var wk = cluster.fork();
workerID.push(wk.id);  ①
wk.send(collection[i]);
}
cluster.on('fork', function (worker) {
if (workerID.indexOf(worker.id) !== -1) { ②
console.log(`[master ${process.pid}] : fork worker
${worker.id}`);
}
});
......
var numOfCompelete = 0
workerID.forEach(function (id) {
cluster.workers[id].on('message', function (msg) {
......
numOfCompelete++;
if (numOfCompelete === collection.length) {
console.log(`[master ${process.pid}] finish all work
and using ${Date.now() - st} ms`);
workerID.forEach(function (id) {
if (!cluster.workers[id].suicide) { ③
cluster.workers[id].disconnect();
}
});
reslove(result);
}
});
});
}
))
}

```

Show moreShow more icon

重要的改动有以下三点：

1. 增加了 workID 数组用来保存 worker ID
2. Master 进程只相应保存在 workID 数组中对应的 worker 进程
3. 在 master 进程关闭 worker 进程之前判断此 worker 进程的状态，以避免”EPIEP”异常的出现

经过此番修改，我们的多线程程序可以正确的响应并发操作了：

##### 清单 17.main.js 代码执行结果

```
PS C:\work\NodeClusterTest> node .\main.js
=====Start========
[master 11796] : fork worker 1
[master 11796] : fork worker 2
[master 11796] : fork worker 3
[master 11796] : fork worker 4
.......
[master 11796] finish all work and using 26499 ms
Finish all work and using 26499
####Get result1 from mutliple-processes: 165580141,433494437,701408733,1134903170
[master 11796] : worker 6 died
[master 11796] : worker 5 died
[master 11796] : worker 7 died
[master 11796] : worker 8 died
[worker 1] work finish work and using 25670 ms
[master 11796] receive message from [worker1]: 1134903170
[master 11796] finish all work and using 26631 ms
Finish all work and using 26561
####Get result from mutliple-processes: 165580141,433494437,701408733,1134903170
[master 11796] : worker 1 died
[master 11796] : worker 2 died
[master 11796] : worker 3 died
[master 11796] : worker 4 died

```

Show moreShow more icon

## 结术语

至此，我们的 Node.js 多线程实战可以告一段落，我们之前得到的数据表明，Node.js 所提供的 cluster 模块确实能显著的提高 Node.js 程序在处理 CPU 密集型应用时的效率。但是，细心的朋友大概已经发现了 master 模块每个进程的执行效率是不如直接使用单进程模式的，我们都只执行一个同样的任务，用多进程和单进程的结果如下：

##### 清单 18\. 单进程与多进程比较

```
PS C:\work\NodeClusterTest> node .\main.js
[master 14752] : fork worker 1
    ......
Finish all work and using 11497 ms

PS C:\work\NodeClusterTest> node .\single.js
Task were complete with result 1134903170 and using 11079 ms

```

Show moreShow more icon

可见两者有着 400+ ms 的差距，这对于 Node.js 系统来说也并不是一个可以容易被忽略的时间。因此，是否要使用多进程模式或者是否一直需要使用多进程模式，这个问题就要根据大家的具体问题来考虑了。