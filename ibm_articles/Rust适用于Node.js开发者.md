# Rust 适用于 Node.js 开发者
阐述 JavaScript 和 Rust 之间的主要区别

**标签:** Node.js,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/rust-for-nodejs-developers/)

[Josh Hannaford](https://developer.ibm.com/zh/profiles/joshua.hannaford)

发布: 2020-09-15

* * *

您是一位 Node.js 或 JavaScript 开发者，对 Rust 编程语言很感兴趣，但不知道从何开始？您很好奇如何将 Node.js 和 JavaScript 知识转换为 “Rust 语言”。在本文中，我将介绍 JavaScript 与 Rust 之间最明显的一些区别和联系，以便您在开始使用 Rust 语言时心里有底。

## 类型系统

我们首先来了解最大的区别：类型系统。我知道，对于熟悉 JavaScript 但不熟悉 TypeScript 的开发者来说，一听到类型系统就会很焦虑。不必担心，Rust 系统是我使用过的最简单、表达性最强的系统之一。

Rust 是一种静态的强类型语言，而 JavaScript 是一种动态的弱类型语言，让我们来了解下这两种语言类型的含义以及对您日常代码的影响。

### 静态类型与动态类型

静态类型是指必须在编译时而不是运行时已知所有类型的变量。这意味着静态类型语言的值不依赖于用户输入，不能在运行时更改，也不包含编译时未知的任何其他内容。

Rust 使用静态类型，而 JavaScript 使用动态类型。让我们看看这两种语言如何转换为一段代码：

**JavaScript（动态）**

```lang-js
let a = 3 // a is a `number` type
a = "Hello, world!"// This is perfectly valid, and a is now a `string`

```

Show moreShow more icon

**Rust（静态）**

```lang-rust
let mut a = 3; // a is now an `i32` value
a = "Hello, world!"; // Compiler Error: The type cannot change at runtime.

```

Show moreShow more icon

关于上述 Rust 代码值得注意的一点是，与其他静态类型的语言（例如 C++ 或 Java）相比，Rust 实际上具有更高级的类型推断功能。这意味着，除了在结体和函数定义中，通常不必指定变量的类型。我认为这是一种均衡有效的方法：对于研究函数签名的任何人来说，这不难理解；而实现该函数的人员也不必担心要在其中键入所有内容。

另一点是，Rust 实际上拥有一种处理变量的模式，您可能希望在函数执行过程中更改变量的类型。例如，如果您的函数中有一个数字，但您希望将该数字作为另一个操作的字符串，则可以通过以下代码来完成。

```lang-rust
let num = 10 * get_random_number(); // num is an `i32`
let num = num.to_string(); // num has been shadowed and is now a `String`

```

Show moreShow more icon

在 Rust 中，这种模式被称为 _遮蔽_（shadowing）。您将在许多 Rust 代码库中看到这种很常见的模式。它的工作原理是：每次使用 `let` 关键字时，Rust 都会在内部创建一个新变量，因此它并不关心赋值与赋值前变量的类型是否不同。在 Rust 看来，它只是另一个新变量，和任何其他变量没什么不同。

现在，我们已经探究了 Rust 的静态类型与动态类型的 JavaScript 之间的区别，下面我们来看看强类型和弱类型语言之间的差异。

### 强类型与弱类型

虽然静态类型和动态类型之间的比较显而易见（类型在编译时已知还是未知），但强类型和弱类型的区别却有点模糊，因为这些类型更像是语言可归属的一个范围。这是因为，这些术语衡量的是一种语言允许在一种类型和另一种类型之间进行隐式转换的程度。

如果一种语言属于强类型，那么与具有许多此类转换实例的语言相比，它允许的隐式转换更少。根据该定义，Rust 在 _强_ 程度上的评级为 **非常高**，因为只有一个隐式转换实例，并且只能在非常特殊的情况下访问此转换。

另一方面，JavaScript 在 _弱_ 程度上的评级很高，因为到处都有隐式转换。实际上，无论知道与否，您分配的 **每个** 变量都有一个隐式转换的特定元素。因为您声明的每个变量都是一个由 JavaScript 运行时赋予了额外属性的对象。

这对您的代码有什么影响？

我们来看看 JavaScript 和 Rust 中两个等效的代码段。

**JavaScript（弱）**

```lang-js
// Implicitly converts num to a `string`
const num = 10 + "0"

// Prints "Num is a string with a value of 100"
console.log(`Num is a ${typeof num} with a value of ${num}`)

```

Show moreShow more icon

**Rust（强）**

```lang-rust
#![feature(type_name_of_val)]
use std::any::type_name_of_val;

// No implicit conversion.You have to manually make 10 a `String` first
let num = 10.to_string() + "0";

// Prints "Num is a std::string::String with a value of 100"
println!("Num is a {} with a value of {}", type_name_of_val(num), num);

```

Show moreShow more icon

正如您在上面的代码中所见，JavaScript 可以将 `number` 隐式转换为 `string`。使用 Rust 时，您必须在添加包含 _“0”_ 的字符串切片之前，明确指出要将 _10_ 转换为字符串。

虽然每个人喜好不同，但我个人更喜欢在我的代码中具有这种明确性，尤其是在处理较大的代码库时（您并未自行编写所有代码）。这种明确性让您可以更容易推断出原始作者的意图并避免错误，例如想知道在添加 _0_ 后您的数字为什么是 _100_ 而不是 _10_。

现在，我们已经了解了类型系统中的一些区别，接下来，我们将探究另一个非常重要的主题，即可变性。

## 可变性

众所周知，在现代 JavaScript 中，可以通过三种方式来声明变量：

- `let`
- `const`
- `var`

虽然每种方法的工作方式各有不同，但出于讨论可变性的目的，我将仅提及 `let` 和 `const`，因为`let` 和 `var` 在可变性方面功能相同。

在 Rust 中，可以通过三种方式来声明变量：

- `let`
- `let mut`
- `const`

就可变性而言，它们的功能各不相同。看看 Rust 中声明变量的不同选项，就很容易理解 `let mut` 可以声明一个可变变量，但如果是这样，`let` 和 `const` 有什么区别？另外，为什么必须使用两个词来表示变量应该是可变的？

要回答第一个问题，您需要记住 Rust 与 JavaScript 不同，Rust 是一种编译语言。这一事实就是 `let` 和 `const` 之间产生主要区别的地方。即使二者都不可变，区别将在于变量初始化的时间 _（注意：我在这里专门使用了“初始化”一词，而不是“声明”）_。

`const` 变量必须在编译时进行初始化，并且该变量的每次出现都在引用它的站点实现内联。 **PI** 值就是一个很典型的例子。您可能希望在整个应用程序中使用该值并保持不变。而另一方面，`let` 变量允许您将其初始化为仅在运行时已知的变量，例如用户输入。虽然上述解释回答了第一个问题，但您可能仍想知道为什么要花更多的精力来创建可变变量，而不是不可变变量？

这是因为 Rust 很大程度上借鉴了函数式编程范例，在纯函数式编程中，每个值都是不可变的。虽然这会大大增加 CPU 负载和内存消耗，但也有很多好处，尽管我不能在这里一一介绍，但最大的好处就是在尝试推理并发/并行代码之时。Rust 并不采用纯函数式编程模式，因此出于性能方面的考虑，您仍然可以声明一个可变变量，但这会让您思考是否真的需要将该变量设为可变。

现在，我们来看看 Rust 和 JavaScript 之间的最后一个主要区别：可空性。

## 可空性

`null` 的发明者 Tony Hoare 爵士有一句 [名言](https://www.infoq.com/presentations/Null-References-The-Billion-Dollar-Mistake-Tony-Hoare/) 称这是一个“价值十亿美元的错误”。他之所以称其为错误，是因为根据他的估计，过去 40 年来，`null references` 已导致了“无数错误、漏洞和系统崩溃……”，损失高达 10 亿美元。

我敢肯定，任何使用允许空引用的语言的开发者在其职业生涯的某些时候一定遇到过 `null` 异常情况。我当然也遇到过！但在 JavaScript 中，当您发现不仅可以传递为 `null` 的变量，还可以传递为 `undefined` 的变量，这种区别会变得更加令人困惑。如果必须检查这两个实例，会让人非常困惑，并且只能在运行时强制执行，因此会出现许多执行失败的情况。庆幸的是，Rust 并不会像这样处理 _空_ 的概念。

### Rust 如何处理可空性

在 Rust 中，有两种主要方法来处理有可能为 _空_ 的值：

- 如果不一定是不存在任何值的错误，则使用 `Option<T>` 类型
- 如果确实是该值不存在的错误，则使用 `Result<T, E>` 类型

在以后的教程中，我将详细介绍如何使用这些模式，但为了快速地展示不同之处，我们来看看使用 API 时这些模式如何工作。如果从数据库返回的对象的特定字段可能为空，那么您可以使用 `Option` 类型来表示这种情况；但如果您要调用 API 并期望得到响应，那么可以使用 `Result` 类型来模拟这种行为。

## 结束语

尽管还可以进行更多的比较，但我希望您现在已经更好地了解了 Rust 的基本语言功能与 JavaScript 的关系！我也希望您可以勇敢尝试 Rust。虽然确实需要一些时间才能熟练使用，但这样做的回报是值得的，您可能会发现您的 JavaScript 代码因此而变得更出色。

此外，敬请继续关注我的下一篇内容，将是一篇教程，该教程将通过在 Rust 中构建 API 来探讨这两种语言之间更底层的一些比较，还将介绍元组、数组、哈希图和字符串操作等概念。

本文翻译自： [Rust for Node.js developers](https://developer.ibm.com/articles/rust-for-nodejs-developers/)（2020-08-05）