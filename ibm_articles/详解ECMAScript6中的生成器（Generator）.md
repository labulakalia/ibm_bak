# 详解 ECMAScript 6 中的生成器（Generator）
生成器的基本用法、高级用法和实际应用

**标签:** JavaScript,Node.js,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/wa-es6-generator/)

成 富

发布: 2018-09-04

* * *

生成器（Generator）是 ECMAScript 6 中引入的新概念。生成器本身是一个很强大的工具，应该出现在每个 JavaScript 开发人员的工具箱之中。不过，生成器目前主要被使用在 JavaScript 框架和第三方库中，而在日常开发中的使用较少。这主要是因为生成器的用法并不容易掌握。本文的目标是让更多的开发人员了解生成器，并使用这一工具在日常开发中简洁高效地完成一些任务。本文中的示例代码在 Node.js 8.10 上测试运行。生成器在浏览器的支持请见参考资源。

## 生成器的基本概念

在讨论生成器的具体用法之前，先从生成器的基本概念开始进行介绍。与生成器相关的概念有两个：

- **生成器函数** ：用来创建生成器对象的一类特殊函数。
- **生成器对象** ：生成器函数的具体实例。

生成器对象的最大特点在于它们的执行可以被暂停和继续。在 JavaScript 中，我们并不能控制普通函数的执行。对于一个函数，当它开始执行之后，就会一直执行到结束，并把返回值传递给调用者。考虑到 JavaScript 引擎执行时的单线程特性（不考虑到使用 WebWorker 的情况），在一个函数的执行过程中，是不能终止该函数的。如果你的函数中不小心引入了无限循环，那么整个应用的执行都会被阻塞。

我们先从最简单的生成器函数开始，如代码清单 1 所示。生成器函数与普通函数的差别在于 function 和函数名称之间的星号（`*`）。生成器函数中使用 `yield` 表达式来产生值。

##### 清单 1\. 基本的生成器函数

```
function *sample() {
yield 1;
yield 2;
yield 3;
}

```

Show moreShow more icon

在调用代码清单 1 中的生成器函数之后，得到的是一个生成器对象。生成器对象在其 `next` 方法被调用时，可以依次返回多个不同的值。这些返回值通过 `yield` 来声明。代码清单 1 中使用 `yield` 来产生了 `1`，`2` 和 `3` 共 3 个值。当 `next` 方法调用时，这些值会被依次返回。在代码清单 2 中，首先从生成器函数 `sample` 中创建了一个新的生成器对象 `func`。该对象 `func` 的执行一开始是被暂停的。当 `next` 方法调用时，`func` 对象开始执行，并执行到第一个 `yield` 表达式，并把结果 `1` 返回给调用者。返回值是一个包含了两个属性 `value` 和 `done` 的对象。属性 `value` 包含的是 `yield` 表达式所产生的值，而 `done` 用来表示是否还有更多值可以被获取。再次调用 `next`，可以继续 `func` 的执行，并执行到第二个 `yield` 表达式。如此循环，直到第四个 `next` 方法调用，`done` 的值才变为 `true`，表明已经没有更多值可以被获取了。

##### 清单 2\. 使用生成器对象

```
let func = sample();
func.next();
// -> {value: 1, done: false}
func.next();
// -> {value: 2, done: false}
func.next();
// -> {value: 3, done: false}
func.next();
// -> {value: undefined, done: true}
func.next();
// -> {value: undefined, done: true}

```

Show moreShow more icon

生成器的强大之处正是来源于可以暂停和继续生成器对象执行的能力。每个生成器对象都可以被看成是一个状态机。同一个生成器函数所创建的每个对象都在内部维护自己的状态，彼此并不会互相影响。调用 `next` 方法会继续生成器的执行，触发内部的状态转换，运行到下一个 `yield` 表达式所在的位置。接着执行会被暂停，等待下一次 `next` 方法的调用。代码清单 3 中创建了生成器函数 `sample` 的 2 个不同的实例。调用其中一个对象的 `next` 方法并不会影响到另一个对象的内部状态。

##### 清单 3\. 同一个生成器函数创建的不同生成器对象

```
let func1 = sample();
let func2 = sample();
func1.next();
// -> {value: 1, done: false}
func2.next();
// -> {value: 1, done: false}
func1.next();
// -> {value: 2, done: false}

```

Show moreShow more icon

## 生成器的基本用法

在上一节中介绍了生成器的基本概念，下面介绍生成器的一些基本用法。

### 调用 next 方法时的参数

首先来看一下代码清单 4 中的生成器函数 `doMath`。如果只是简单的看代码，可能会认为在调用生成器对象的 `next` 方法之后，`x` 的值应该是 `1`，`y` 的值应该是 `11`，`z` 的值应该是 `110`。

##### 清单 4\. 生成器函数 doMath

```
function *doMath() {
let x = yield 1;
let y = yield x + 10;
let z = yield y * 10;
}

```

Show moreShow more icon

然而实际的执行结果并不符合我们的预期。如代码清单 5 所示，实际的值分别是 `1`，`NaN` 和 `NaN`。了解产生这样结果的关键点在于，在调用 `next` 方法时所传递的参数值，被作为上一个 `yield` 表达式的实际值。由于我们在调用 `next` 时没有传入任何参数，每个 `yield` 表达式的实际值都是 `undefined`。在第一个 `next` 调用中，由于没有上一个 `yield` 表达式，因此该值被忽略；在第二个 `next` 调用中，上一个 `yield` 表达式，也就是 `yield 1` 的值是 `next` 调用时的参数值，也就是 `undefined`，因此 `x` 的值是 `undefined`，所以 `yield x + 10` 的值是 `NaN`。由此类推，第三个 `next` 调用中，`y` 的值同样也是 `undefined`，因此产生的值也是 `NaN`。

##### 清单 5\. 使用生成器函数 doMath

```
let func = doMath();
func.next();
// -> {value: 1, done: false}
func.next();
// -> {value: NaN, done: false}
func.next();
// -> {value: NaN, done: false}
func.next();
// -> {value: undefined, done: true}

```

Show moreShow more icon

现在可以尝试在调用 `next` 方法时传递值。在代码清单 6 中，第二个 `next` 调用时传入了值 1，因此 `1` 被作为上一个 `yield` 表达式 `yield 1` 的值，把 `x` 设为 `1`，所以该 `next` 方法调用的值为 `11`；第三个 `next` 调用时传入了值 `2`，因此 `2` 被作为第二个 `yield` 表达式 `yield x + 10` 的值，把 `y` 设为 `2`，所以该 `next` 方法调用的值为 `20`。

##### 清单 6\. 在调用 next 方法时传递值

```
let func = doMath();
func.next();
// -> {value: 1, done: false}
func.next(1);
// -> {value: 11, done: false}
func.next(2);
// -> {value: 20, done: false}
func.next(3);
// -> {value: undefined, done: true}

```

Show moreShow more icon

### 使用 return 语句

在生成器函数中，同样可以使用 `return` 语句。通过 `retur`n 返回的值也会被传递给 `next` 方法的调用者，同时会结束掉生成器对象的执行，也就是把属性 `done` 的值设为 `tru`e。在代码清单 7 中，第二个 `next` 的调用会执行到 `return` 语句并结束执行。

##### 清单 7\. 使用 return 语句

```
function *withReturn() {
let x = yield 1;
return x + 2;
}

let func = withReturn();
func.next();
// -> {value: 1, done: false}
func.next(1);
// -> {value: 3, done: true}
func.next();
// -> {value: undefined, done: true}

```

Show moreShow more icon

### 生成器与迭代器

从之前使用生成器对象的代码中可以发现，`next` 方法的返回值并不是特别直观，需要通过属性 `done` 来判断是否还有值。实际上，这是因为生成器对象本身也是迭代器（iterator）对象，而迭代器对象用 `next` 方法来获取其中的值。同时生成器对象也是可被迭代的（iterable）。因此，我们可以用 ECMAScript 6 中的其他新特性来遍历其中的值，包括 `for-of` 循环，`spread` 操作符和新的集合类型。在代码清单 8 中，首先使用了 `for-of` 循环来遍历一个生成器对象，接着使用 `spread` 操作符把生成器对象的值添加到数组中，最后使用生成器中的值创建了一个 `Set` 对象。

##### 清单 8\. 以迭代器的方式来使用生成器

```
for (let value of sample()) {
console.log(value);
}
// -> 输出 1，2 和 3
['a', ...sample(), 'b']
// -> [ 'a', 1, 2, 3, 'b' ]

let set = new Set(sample())
set.size
// -> 3

```

Show moreShow more icon

### 生成器函数的参数

与普通函数一样，生成器函数也可以接受输入参数。这些参数可以在 `yield` 表达式中使用。在代码清单 9 中，生成器函数 `seq` 有两个参数 `start` 和 `numbe`r，分别表示所产生值的初始值和值的数量。两个参数都有默认值。函数 `debug` 的作用是输出生成器中的全部值，在后面的代码中也会用到。

##### 清单 9\. 生成器函数的参数

```
function debug(values) {
for (let value of values) {
    console.log(value);
}
}

function *seq(start = 0, number = 10) {
while (number-- > 0) {
    yield start++;
}
}

debug(seq());
// -> 输出从 0 到 9 的值

debug(seq(3));
// -> 输出从 3 到 12 的值

debug(seq(3, 5));
// -> 输出从 3 到 7 的值

```

Show moreShow more icon

## 生成器的高级用法

在介绍完生成器的基本用法之后，下面介绍生成器的一些高级用法。

### 生成器对象的 return 方法

生成器对象的 `return` 方法可以用来返回给定值并结束它的执行。其使用效果类似于在生成器函数中使用 `return` 语句。在代码清单 10 中，调用 `func.return('d')` 会返回传入的值 `d`，并结束生成器，也就是 `done` 的值变为 `true`，即使生成器中仍然还有值 `b` 和 `c` 未被生成。方法 `return` 可以被多次调用，每次调用都返回传入的值。

##### 清单 10\. 生成器对象的 return 方法

```
function *values() {
yield 'a';
yield 'b';
yield 'c';
}

let func = values();
func.next();
// -> {value: "a", done: false}
func.return('d');
// -> {value: "d", done: true}
func.next();
// -> {value: undefined, done: true}

```

Show moreShow more icon

### 生成器对象的 throw 方法

生成器对象的 `throw` 方法可以用来传入一个值，并使其抛出异常。`throw` 和之前提到的 `next` 都可以传入值到生成器对象中来改变其行为。通过 `next` 传入的值会作为上一个 `yield` 表达式的值，而通过 `throw` 传入的值则相当于把上一个 `yield` 语句替换到一个 `throw` 语句。在代码清单 11 中，当 `func.throw('hello')` 被调用时，上一个 `yield` 表达式 `yield x + 1` 被替换成 `throw 'hello'`。由于抛出的对象没有被处理，会被直接传递到 JavaScript 引擎，导致生成器的执行终止。

##### 清单 11\. 生成器对象的 throw 方法

```
function *sample() {
let x = yield 1;
let y = yield x + 1;
yield y * 10;
}

let func = sample();
func.next();
// -> {value: 1, done: false}
func.next(1);
// -> {value: 2, done: false}
func.throw('hello');
// -> Uncaught hello
func.next();
// -> {value: undefined, done: true}

```

Show moreShow more icon

我们可以在生成器函数中使用 try-catch 来捕获异常并处理。代码清单 12 中，在调用 `func.throw(new Error('boom!'))` 时，上一个 `yield` 表达式 `yield 2` 被替换成了 `throw new Error('boom!')`。抛出的对象由 `try-catch` 进行了处理，因此生成器的执行可以被继续。

##### 清单 12\. 使用 try-catch 捕获异常

```
function *sample() {
yield 1;
try {
    yield 2;
} catch (e) {
    console.error(e);
}
yield 3;
yield 4;
}

let func = sample();
func.next();
// -> {value: 1, done: false}
func.next();
// -> {value: 2, done: false}
func.throw(new Error('boom!'));
// -> Error: boom!
// -> {value: 3, done: false}
func.next();
// -> {value: 4, done: false}

```

Show moreShow more icon

### 使用 yield \* 表达式

目前我们看到的生成器对象每次只通过 `yield` 表达式来产生一个值。实际上，我们可以使用 `yield *` 表达式来生成一个值的序列。当使用 `yield *` 时，当前生成器对象的序列生成被代理给另外一个生成器对象或可迭代对象。代码清单 13 中的生成器函数 `oneToThree` 通过 `yield* [1, 2, 3]` 来生成 3 个值，与 [清单 1](#清单-1-基本的生成器函数) 中的生成器函数 `sample` 的结果是相同的，不过使用 `yield *` 的方式更加简洁易懂。

##### 清单 13\. 使用 yield \* 表达式

```
function *oneToThree() {
yield* [1, 2, 3];
}

debug(oneToThree());
// -> 输出 1, 2, 3

```

Show moreShow more icon

在一个生成器函数中可以使用多个 `yield *` 表达式。在这种情况下，来自每个 `yield *` 表达式的值会被依次生成。 代码清单 14 中的生成器 multipleYieldStars 使用了 2 个 `yield *` 语句和一个 `yield` 语句。这里需要注意的是字符串 `hello` 会被当成一个可迭代的对象，也就是会输出其中包含的每个字符。

##### 清单 14\. 使用多个 yield \* 表达式

```
function *multipleYieldStars() {
yield* [1, 2, 3];
yield 'x';
yield* 'hello';
}

debug(multipleYieldStars());
// -> 输出 1, 2, 3, 'x', 'h', 'e', 'l', 'l', 'o'

```

Show moreShow more icon

由于 `yield *` 也是表达式，它是有值的。它的值取决于在 `yield *` 之后的表达式。`yield *` 表达式的值是其后面的生成器对象或可迭代对象所产生的最后一个值，也就是属性 done 为 true 时的那个值。如果 `yield *` 后面是可迭代对象，那么 `yield *` 表达式的值总是 `undefined`，这是因为最后一个生成的值总是 `{value: undefined, done: true}`。如果 `yield *` 后面是生成器对象，我们可以通过在生成器函数中使用 `return` 来控制最后一个产生的值。在代码清单 15 中，通过 `return` 来改变了生成器函数 abc 的返回值，因此 `yield *abc()` 的值为 `d`。

##### 清单 15\. yield \* 表达式的值

```
var result;

function loop(iterable) {
for (let value of iterable) {
    //ignore
}
}

function *abc() {
yield* 'abc';
return 'd';
}

function *generator() {
result = yield* abc();
}

loop(generator());
console.log(result);
// -> "d"

```

Show moreShow more icon

表达式 `yield` 和 `yield *` 都可以进行嵌套。在代码清单 16 中，最内层的 `yield` 表达式生成值 `1`，然后中间的 `yield` 表达生成 `yield 1` 的值，也就是 `undefined`。这是因为在遍历调用 `next` 时并没有传入参数。最外层的 `yield` 的值也是 `undefined`。

##### 清单 16\. 嵌套的 yield 表达式

```
function *manyYields() {
yield yield yield 1;
}

debug(manyYields());
// 输出 1, undefined, undefined

```

Show moreShow more icon

在代码清单 17 中，内层的 `yield *` 首先产生由生成器 `oneToThree` 所生成的值 `1`，`2` 和 `3`，然后外层的 `yield` 再产生 `yield *` 的值，也就是 `undefined`。

##### 清单 17\. 嵌套的 yield \* 和 yield 表达式

```
function *oneToThree() {
yield* [1, 2, 3];
}

function *values() {
yield yield* oneToThree();
}

debug(values());
// -> 输出 1, 2, 3, undefined

```

Show moreShow more icon

## 生成器的使用案例

在介绍完生成器的高级用法之后，下面介绍生成器在框架和实际开发中的具体案例。

### 生成器在框架中的使用

下面通过一个例子来说明生成器在框架和第三方库中的使用。Babel 是一个流行的 JavaScript 编译工具，可以把使用了最新规范的 JavaScript 代码转换成可以在今天的平台上运行的代码。async/await 是一个实用的新特性。为了在老的平台上支持这一特性，Babel 提供了一个插件来把使用了 async/await 的代码转换成使用生成器。代码清单 18 中给出了一个使用 async/await 的简单函数。

##### 清单 18\. 使用 async/await 的代码

```
async function foo() {
await bar();
}

```

Show moreShow more icon

上述代码经过 Babel 转换之后的结果如代码清单 19 所示。首先，async 函数被转换成生成器函数，而 `await` 被转换成 `yield` 表达式。接着通过函数`_asyncToGenerator` 把生成器函数转换成一个返回 `Promise` 的普通函数。

##### 清单 19\. 转换之后的代码

```
let foo = (() => {
var _ref = _asyncToGenerator(function* () {
    yield bar();
});

return function foo() {
    return _ref.apply(this, arguments);
};
})();

function _asyncToGenerator(fn) { return function () { var gen = fn.apply(this, arguments);
return new Promise(function (resolve, reject) { function step(key, arg) { try { var info =
gen[key](arg); var value = info.value; } catch (error) { reject(error); return; } if
(info.done) { resolve(value); } else { return Promise.resolve(value).then(function (value)
{ step("next", value); }, function (err) { step("throw", err); }); } } return
step("next"); }); }; }

```

Show moreShow more icon

### 生成器在实际开发中的使用

在实际开发中，生成器的一个常见使用场景是动态生成序列。代码清单 20 中的生成器函数 `numbers` 使用了复杂的逻辑来定义所产生的值。要理解为什么会产生这些值，关键在于理解到每次 `next` 方法调用会使得生成器对象运行到下一个 `yield` 表达式：

- 当 i 的值是 0 到 4 时，满足条件 i<5，所以会由 yield i 来生成；
- 当 i 变为 5 之后，没有任何条件满足，不产生任何值；
- 当 i 为 6 时，满足第二个条件，因此产生值 12；
- 当 i 为 7 时，满足最后一个条件，产生值 49；
- 当 i 为 8 时，满足第二个条件，产生值 16；
- 当 i 为 9 时，满足第三个条件，产生值 27；
- 直到 i 为 12，才由第三个条件产生值 36；
- 最后在当 i 为 14 时，产生值 98。

##### 清单 20\. 使用生成器来产生序列

```
function *numbers() {
for (let i = 0; i < 20; i++) {
    if (i < 5) {
      yield i;
    } else if (i < 10 && i % 2 === 0) {
      yield i * 2;
    } else if (i < 15 && i % 3 === 0) {
      yield i * 3;
    } else if (i % 7 === 0) {
      yield i * 7;
    }
}
}

debug(numbers());
// -> 输出数字: 0, 1, 2, 3, 4, 12, 49, 16, 27, 36, 98

```

Show moreShow more icon

## 结束语

ECMAScript 6 中的生成器是一个强大的工具。目前在很多 JavaScript 框架和库中都得到了应用。本文对生成器的基本概念、基本用法、`return` 和 `throw` 方法、以及 `yield *` 表达式的用法都做了详细的介绍。通过本文的介绍，读者可以对生成器有一个更深入的了解，并尝试在实际开发中使用生成器来解决一些具体的问题。