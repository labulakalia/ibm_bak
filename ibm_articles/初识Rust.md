# 初识 Rust
一种替代 C 和 C++ 的编译语言

**标签:** 软件开发

[原文链接](https://developer.ibm.com/zh/articles/os-know-rust/)

Dylan Hicks

发布: 2018-06-14

* * *

如果您编写过软件，您肯定问过自己“我应该使用哪种语言来编写这个程序？”这是一个合理的问题。您的代码是否需要尽可能快？它是否会在网络上运行？该代码将位于后端还是前端？所有语言都有自己的定位，Rust 也不例外。

Rust 是一种静态类型的编译语言，满足了大多数用户使用 `C` 或 `C++` 能够实现的目标。但是不同于 `C` 和 `C++` ，Rust 还侵占了 `C#` 和 Java™ 语言在本世纪统治了很长时间的领域：Rust 语言是内存安全且与操作系统无关的，这意味着它可以在任何计算机上运行。实质上，您会获得系统语言的所有速度和低阶优势，而没有我提及的后几种语言中麻烦的垃圾收集。激动不已？是的，我也很激动。欢迎使用 Rust！

**注意：** 本文中演示的所有代码和命令都是我在运行 Linux® 的计算机上写的，但它们很容易移植到 macOS® 和 Windows®。

## 安装并运行 Rust

大多数初学者指南都会直接分析代码，但我认为在开始开发之前，最好首先了解如何编译和运行该软件。无论您采用何种操作系统， [Rust 安装说明](https://www.rust-lang.org/en-US/install.html) 都介绍了必要的安装步骤。

接下来，我将设置一个项目，将我的示例放入其中。我可以自行创建项目的目录结构，但我更愿意使用 Rust 附带的强大的包管理器 Cargo。参阅 [Cargo 文档](https://doc.rust-lang.org/cargo/getting-started/first-steps.html) ，使用 Cargo 创建一个新的二进制可执行程序很容易。只需在终端发出此命令：

```
$ cd ~/Documents
$ cargo new rusty_start –bin

```

Show moreShow more icon

`cargo new rusty_start` 命令告诉 Cargo 创建一个名为 `rusty_start` 的新程序目录结构。 `--bin` 选项告诉 Cargo，此程序将是一个二进制可执行程序。反之，如果保持此选项关闭，Cargo 将创建一个库。 实质上，将文件 main.rs 替换为 lib.rs 很简单。查看 `tree` 程序中的新 `rusty_start` 目录，我看到了以下代码：

```
$ cd rusty_start
$ tree .
.
??? Cargo.toml
??? src
    ??? main.rs

```

Show moreShow more icon

完美！不过，Cargo.toml 文件是什么？它是 Cargo 用来构建程序的文件。现在，它描述了程序的作者和程序的依赖项。对于本示例，我是唯一的作者，而且不需要任何外部包。

在 main.rs 中可以注意到，Cargo 还生成了一些代码：

```
fn main() {
    println!("Hello, world!");
}

```

Show moreShow more icon

这是 Rust 定义其 `main` 函数的方式。就现在而言，只需知道此函数会输出到终端即可。现在的问题是，我如何运行此代码？ 这实际上很简单 ( [Rust 中的 Hello World!](#rust-中的-hello-world))。

##### Rust 中的 Hello World!

```
$ cargo build
    Compiling...
     Finished...
$ cargo run
     Finished...
      Running...
Hello, world!

```

Show moreShow more icon

与 `; java` 或设置 CMake 相比，这无疑是非常简单的。为了更容易操作，我可以执行 `cargo run` ，这个命令会同时编译和执行。

安装 Rust 并运行它之后，我就可以介绍一些重要功能了。

## Rust 的核心功能

本节包括两个主题：绑定变量和函数，它们使得理解后续小节中的示例变得容易得多。

### 了解绑定

不同于前面构建的 Hello, World! 程序，最终我需要在我的 Rust 代码中使用变量。为了将值绑定到名称，我在 main.rs 文件中使用了 `let` 关键字：

```
fn main() {
    let x = 5;
}

```

Show moreShow more icon

在此代码段中，我将值 `5` 绑定到了变量 `x` 。 但是，类型该如何处理？我之前是否说过 Rust 是一种静态类型语言？Rust 的另一个有趣事实是，它可以在编译时通过值来推断变量的类型 – 一个无需在运行时付出代价即可获得的出色功能。如果我想显式声明类型（比如一个 32 位有符号整数），我可以轻松地完成此工作：

```
let x: i32 = 5;

```

Show moreShow more icon

您也可以从多种原语类型中进行选择，下一节将介绍此主题。

但是，在介绍函数之前，我想使用模式来演示 Rust 的 `let` 绑定的一个有趣且有用的方面：

```
let (x, y) = (5, '6'); // types (i32, char)

```

Show moreShow more icon

在这里，我将整数 `5` 绑定到变量 `x` ，将字符 `6` 绑定到变量 `y` 。此技巧非常有用，而且模式在 Rust 社区中随处可见（该社区的成员被称为 _Rustacean_ ）。使用模式可以做的事情要多得多，但为了简洁起见，在此不做累述。

默认情况下，Rust 中的一切都是不可变的，这意味着您无法在以后执行更改。例如，如果我编译了这个命令：

```
let x: i32 = 5;
x = 6;

```

Show moreShow more icon

编译器告诉我 `x` 不可变。所以，要更改此变量，必须使用 `mut` 关键字：

```
let mut x: i32 = 5;
x = 6;

```

Show moreShow more icon

现在，一切恢复正常！

我要介绍的变量的最后一个方面是作用域 (scoping)。 _作用域_ 体现的思想是，变量仅在实例化它们的一组给定的花括号 (`{}`) 内有效。从作用域外访问变量会导致出现编译错误。考虑 [作用域变量](#作用域变量) 中的代码。

##### 作用域变量

```
fn main() {
    let x: i32 = 5; // x lives in this scope

    {
        let y: i32 = 4; // y only live in this scope.
    }

    println!("x: {}, y: {}, x, y); // y is not in this scope
}

```

Show moreShow more icon

在涉及到 Rust 的所有权和借用属性时，理解作用域至关重要。现在，我们将介绍一些函数！

### 您的函数是什么？

变量固然很好，但如果没有随行的伙伴 – 函数，它们终将无用。每个程序至少有一个函数，在 Rust 中，该函数就是 `main` 函数。但是，我可以使用 `fn` 关键字定义新函数：

```
fn foo() {
}

```

Show moreShow more icon

您会问，如何采用参数？很简单：

```
fn print_value(value: i32) {
    println!("The value given was: {}”, value);
}

```

Show moreShow more icon

此函数采用一个整数参数并输出它的值。 [输出一个参数的值](#输出一个参数的值) 展示了如何在 main.rs 中使用参数。

##### 输出一个参数的值

```
fn main() {
    print_value(17);
}

fn print_value(value: i32) {
    println!("The value given was: {}”, value);
}

```

Show moreShow more icon

就像使用 `let` 分配变量一样，参数遵循的模式是一个名称后跟一个冒号 (`:`)，然后是一个类型。对于多个参数，只需使用逗号 (`,`) 分隔变量即可：

```
fn print_values(value_1: i32, value_2: i32) {
    println!("Values given were: {} and {}”, value_1, value_2);
}

```

Show moreShow more icon

既然可以创建包含和不含参数的函数，我想尝试创建返回一个值的函数。如何创建一个将整数加 1 的函数？

```
fn increase_by_one(value: i32) -> i32 {
    value += 1
}

```

Show moreShow more icon

在 Rust 中，函数仅返回一个值，而且类型是在箭头 (`->`) 后声明的。但是，我在何处返回 i32？有趣的事实：在 Rust 中，常见的做法是用分号结束来返回值。（该语言有一个 `return` 关键字，但它不经常使用。）如果我想使用 `return` 关键字，可以编写一个类似这样的函数：

```
fn divmod(x: i32, y: i32) -> (i32, i32) {
    return (x / y, x % y);
}

```

Show moreShow more icon

就这么简单！这就是绑定变量和编写函数的基础知识。正如我所提到的，Rust 拥有大量要处理的类型。我会在下一节中简要介绍。

## 我是您的类型吗？

像其他编程语言一样，Rust 有一个类型系统。以下各节会简要介绍这些类型。

### 布尔值

对于布尔类型，没有太多要讲的。Rust 有一种内置的布尔类型（ `bool` ），该类型有两个值： `true` 和 `false` 。这是一个示例：

```
let x = true;
let y: bool = !x; // y == false

```

Show moreShow more icon

### 字符

像 `bool` 一样，Rust 有一种表示单一 Unicode 标量值的类型： `char` 。但与大多数语言不同， `char` 不是一种填入数字值的类型。字符使用单引号来初始化，如下所示：

```
let x = 'x';
let y: char = 'y';

```

Show moreShow more icon

### 字符串

Rust 是一种有着两种字符串类型的特殊语言： `str` 和 `String` 。关于 `str` 没有特别重要的地方，因为它”未确定大小”，但是如果在它前面添加一个逻辑与符号 (`&`)， `&str` 类型就会变得非常有用。简言之，它使用了 `str` 的引用 – 我将在下一节中介绍此主题。 `String` 是 `&str` 的由堆分配的对象版本。

首先，让我们详细分析一下 `&str` 。此类型称为 _字符串切片_ 。它是固定大小的、不可变的 UTF-8 字节序列。我可以使用双引号将此类型绑定到变量：

```
let x = "Hello, World!”;
let y: &str = "Isn't it a wonderful life?!”;

```

Show moreShow more icon

`String` 对象是可变的，我可以使用 `String::from` 构造函数来初始化该类型：

```
let x = String::from("Hello, World”);
let y: String = String::from("Isn't it a wonderful life?!”);

```

Show moreShow more icon

### 数字

Rust 中的数字类型涵盖所有符号和大小。以下是目前类型的列表：

- `i8`
- `i16`
- `i32`
- `i64`
- `u8`
- `u16`
- `u32`
- `u64`
- `isize`
- `usize`
- `f32`
- `f64`

以 `i` 为前缀的类型表示一个有符号整数，这意味着该数字可以为负值。以 `u` 为前缀的类型是无符号整数。以 `f` 为前缀的类型是浮点数。后缀的数字表示该类型的字节大小，后缀的单词 size 表示系统的指针大小，用于访问序列数组的元素。

### 数组

谈到数组，我想讲讲 Rustacean 是如何表示数据序列的。基本策略是使用一个数组，比如：

```
let x = [1, 2, 3];
let y: [i32; 3] = [4, 5, 6];

```

Show moreShow more icon

显式声明数组的类型时，公式为 `[T; N]` ，其中 _T_ 是元素的类型， _N_ 是元素的数量。Rust 中的数组与其他语言中的数组类似，通过方括号 (`[]`) 来存取元素，但它们也有内置的函数，比如 `len` （数组中的元素数量）：

```
let x = [1, 2, 3];
let s: usize = x.len(); // s == 3
let e = x[1]; // e is an i32, and e == 2;

```

Show moreShow more icon

现在，尽管数组看似很有意思，但它们不经常使用。Rust 社区最常使用的是矢量。不同于数组，矢量在堆上分配它们的数据，这非常类似于 `&str` 与 `String` 之间的区别。 您可以使用宏 `vec!` 创建这些矢量，它们具有类型 `Vec<T>` ，其中 _T_ 是已包含的元素的类型：

```
let x = vec![1, 2, 3];
let y: Vec<i32> = [4, 5, 6];

```

Show moreShow more icon

### 元组

最后一个 Rust 基础类型是 _元组_ ，它是一种有序、不可变的对象列表。它们的优势在于类型不必是统一的：

```
let x = (5, '6');
let y: (i32, char) = (7, '8');

```

Show moreShow more icon

不同于数组，这些类型必须以稍微不同的方式建立索引：

```
let x = (5, '6');
let y: i32 = x.0; // y == 5
let z: char = x.1; // z == '6'

```

Show moreShow more icon

如果我想创建包含单个值的元组，可以在第一个元素后添加一个逗号。否则，会忽略圆括号：

```
let x: (bool) = (true,); x == (true)
let y: bool = (false); y == false

```

Show moreShow more icon

## 流控制

如果没有一组方便的控制关键字，很难在 Rust 中编写和编程。以下是两个重要的关键字。

### If

Rust 类似于大多数其他编程语言，所以您应该熟悉 [Rust 中的 if 语句](#rust-中的-if-语句) 中的 Rust `if` 语句示例。

##### Rust 中的 if 语句

```
let x = '5';

if x == '5' {
    println!("X is the char '5'”);
}

This can also be extended with the other keywords else if and else:

let x = '5';

if x == '5' {
    println!("X is the char '5'!”);
} else if x == '6' {
    println!("X is the char '6'!”);
} else {
    println!("I don't know what X is.”);
}

```

Show moreShow more icon

也可以像三元不表达式那样使用 `if` 代码块：

```
let x = '5';
let y = if x == '5' { 5 } else if x == '6' { 6 } else { -1 };
// y == 5

```

Show moreShow more icon

此代码运行正常，因为等号 (`=`) 左侧的信息是返回一个值的表达式，它非常类似于一个函数。 请注意， `if` 代码块中的整数上缺少分号。

### 循环

同样地，如果没有循环，Rust 就不太像一种语言。该语言提供了 3 种循环方法： `loop` 、 `for` 和 `while` 。首先看看 `loop` ，它会无限循环：

```
loop {
    println!("Looping for eternity!”);
}

```

Show moreShow more icon

如果您像我一样，或许并不总是喜欢无限循环。幸运的是，Rust 提供了一种通用的 `while` 循环，它在给定条件 ( [Rust 中的 while 循环](#rust-中的-while-循环)) 上中断。

##### Rust 中的 while 循环

```
let mut x = 0;

while x != 5 {
    println!("x: {}”, x);
    x += 1; // this is equal to x = x + 1;
}

```

Show moreShow more icon

[在 Rust 中创建一个带 break 关键字的循环](#在-rust-中创建一个带-break-关键字的循环) 展示了如何结合使用 `loop` 和 `break` 关键字，这会让循环在那一时刻停止。

##### 在 Rust 中创建一个带 break 关键字的循环

```
let mut x = 0;

loop {
    if x == 5 {
        break;
    }

    println!("x: {}”, x);
    x += 1;
}

```

Show moreShow more icon

最后，Rust 还拥有 `for` 循环。不同于其他循环，此循环不像其他大多数 `C` 风格的语言，这些语言使用的格式是：

```
for(int x = 0; x < 5; x++) {
    printf("%d\n”, x);
}

```

Show moreShow more icon

而 Rust 采用了 [一种典型的 Rust for 循环（选项 1）](#一种典型的-rust-for-循环（选项-1）) 中的格式。

##### 一种典型的 Rust for 循环（选项 1）

```
for x in 0..5 {
    println!("x: {}”, x);
}

```

Show moreShow more icon

或者更通用的格式 ( [一种典型的 Rust for 循环（选项 2）](#一种典型的-rust-for-循环（选项-2）))…

##### 一种典型的 Rust for 循环（选项 2）

```
for var in expression {
...
}

```

Show moreShow more icon

… 其中 `expression` 可以转换为一个迭代器。在清单 7 中的示例中， `0..5` 转换为一个从第一个数字开始到第二个数字结束（不含第二个数字）的迭代器。 `0..5` 称为一个 _值域_ 。Rust 未遵循标准的 `C` 型格式，所以它可以避免处理循环的每个元素时的常见错误。此外，它通过类似矢量的对象来支持 Python 型迭代，比如在 [Rust 中的迭代](#rust-中的迭代) 中。

##### Rust 中的迭代

```
let x = vec![0, 1, 2, 3, 4];

for element in &x {
    println!("element: {}”, element);
}

// prints out the elements in x: 0, 1, 2, 3, 4

```

Show moreShow more icon

`&x` 表示 `for` 循环”借用了”矢量，该概念将在下一节中介绍。

## 从所有者借用

传递、借用是本文以及整个 Rust 中最重要部分。巧合的是，借用也是最让新用户感到沮丧的概念。此系统使得 Rust 能够像 `C` 一样快但更安全地运行。

### 所有权

在 Rust 中，每个变量都有它绑定到的值的所有权。在变量超出作用域时，Rust 就会清理该资源。例如：

```
fn foo() {
    let x = String::from("Hello, World!”);
}

```

Show moreShow more icon

在函数 `foo` 的作用域内创建 `x` 时，字符串中分配的数据会在堆中分配。 现在，当函数运行结束时， `x` 将超出作用域，Rust 会清理相关的字符串数据，甚至是堆中的字节。我不必再生活在 `C/C++` 的痛苦中，也不需要记住包含 `free` 或 `delete` 。

另一种实际查看所有权的好方法是使用移动语义。除了每个变量都拥有其绑定的值的所有权之外，Rust 还可以确保每个资源都被绑定到单个值。换言之，如果一个资源从一个变量重新分配给另一个变量，则第一个变量不再有效 ( [Rust 中的移动语义](#rust-中的移动语义))。

##### Rust 中的移动语义

```
fn main() {
    let x = String::from("Hello, World!”);

    println!("{}”, x); // x is valid here

    let y = x; // the resource assigned to x is move to y

    println!("{}”, x);
}

```

Show moreShow more icon

如果运行此代码，程序将会失败，而且编译器会告诉我，我尝试使用的变量 `x` 已转移到变量 `y` 。此行为也适用于函数 ( [Rust 函数中的移动语义](#rust-函数中的移动语义))。

##### Rust 函数中的移动语义

```
fn main() {
    let x = vec![1, 2, 3];

    take_ownership(x);

    println!("{:?}”, x); // ignore the ':?' for now.
                         // it just prints the full vector.
}

fn take_ownership(v: Vec<i32>) {
println!("I took this data: {:?}”, v);
}

```

Show moreShow more icon

此程序同样会失败，而且编译器会告诉我，我尝试使用的变量 `x` 已转移到函数 `take_ownership` 中的变量 `v` 。我可以使用 [修复编译器错误](#修复编译器错误) 中的构造来修复此问题。

##### 修复编译器错误

```
fn main() {
    let x = vec![1, 2, 3];

    let x = print_vector(x);

    println!("{:?}”, x);
}

fn print_vector(v: Vec<i32>) -> Vec<i32> {
println!("I took this data: {:?}, and returned it”, v);
v
}

```

Show moreShow more icon

但是，这看起来很傻，不是吗？如果我想对矢量中的元素求和并返回该值，该怎么办？我真的需要返回一个包含原始矢量及其和的元组吗？不需要：这里可以采用借用概念。但是在介绍该主题之前，请记住 3 条所有权规则：

- Rust 中的每个值都有一个称为其 _所有者_ 的变量。
- 一次只能有一个所有者。
- 当所有者超出作用域时，该值将被丢弃。

### 引用和借用

有了所有权的明确定义后，我将深入研究另一个同样重要的概念：借用。Rust 在借用资源方面有 3 条规则：

- 任何借用的有效期都不能超过原始所有者的作用域。
- 在任何给定时刻，您都能不可变地借用（或 _引用_ ）一个资源许多次。
- 在任何给定时刻，您都能可变地借用（或 _引用_ ）一个资源一次。

请注意，规则 2 和规则 3 是互斥的。

在 [所有权](#所有权) 小节末尾，我留下了一个关于对象所有权和函数的令人讨厌的难题，如 [对象所有权难题](#对象所有权难题) 中所示。

##### 对象所有权难题

```
fn main() {
    let x = vec![1, 2, 3];

    let x = print_vector(x);

    println!("{:?}”, x);
}

fn print_vector(v: &Vec<i32>) -> Vec<i32> {
println!("I took this data: {:?}, and returned it”, v);
v
}

```

Show moreShow more icon

我必须经常将原始资源返回给原始所有者。此功能不符合 Rust 的语言习惯。我可以使用借用系统和引用来实现同样的目的，但没有可怕的返回结果。 （ [通过引用实现借用](#通过引用实现借用) ）。

##### 通过引用实现借用

```
fn main() {
    let x = vec![1, 2, 3];

    print_vector(&x);

    println!("{:?}”, x);
}

fn print_vector(v: &Vec<i32>) {
println!("I borrowed this data: {:?}”, v);
}

```

Show moreShow more icon

清单 14 中的程序能正常编译。通过在将 `x` 传递给 `print_vector` 和 `print_vector` 的参数时向其添加一个逻辑与符号作为前缀，我解决了一个大问题。在用语习惯上，我将此称为 `print_vector` 借用 `x` 。现在，如果我想写一个函数来将值 `5` 推送给 `x` ，该怎么做？首先，我必须让 `m` 可变。 （ [让 m 可变](#让-m-可变) ）。

##### 让 m 可变

```
fn main() {
    let mut x = vec![1, 2, 3];

    push_five(&x);

    println!("{:?}”, x);
}

fn push_five(v: &Vec<i32>) {
    v.push(5);
}

```

Show moreShow more icon

哪里出错了？构建失败了，而且编译器抛出了错误！这是因为 `print_five` 的参数是不可变的。它只是一个对该资源的不可变引用。让我修改一下（ [修复可变构建](#修复可变构建) ）。

##### 修复可变构建

```
fn main() {
    let mut x = vec![1, 2, 3];

    push_five(&x);

    println!("{:?}”, x);
}

fn push_five(v: &mut Vec<i32>) {
    v.push(5);
}

```

Show moreShow more icon

仍然有问题？！不错：我仅将 `x` 的一个不可变引用传递给了 `print_five` 。让我修改一下这一行：

```
push_five(&mut x);

```

Show moreShow more icon

大功告成。编译器最终不再报错。现在看看这 3 条规则，让我看看能否打破这些规则。在这里，我打破了规则 1（ [打破 Rust 规则 1](#打破-rust-规则-1) ）。

##### 打破 Rust 规则 1

```
fn main() {
    let x = 4;
    let mut x_ref: &i32 = &x;

    // this works, and prints "x_ref: 4”
    println!("x_ref: {}”, x_ref);

    {
        let y = 5;
        x_ref = &y;
    }

    // this, however, will fail, according to rule 1!
    println!("x_ref: {}”, x_ref);
}

```

Show moreShow more icon

编译器告诉我，借用的值 `y` 的寿命不足以让 `x_ref` 保持有效，此规则很不错，因为我在 `C` 中犯了很多这样的错误，而我甚至不知道！现在，我打破了规则 2 和规则 3 的互斥特性 ( [打破 Rust 规则 2](#打破-rust-规则-2))。

##### 打破 Rust 规则 2

```
fn main() {
    let mut x = 4; // needs to be mut to borrow mutable.
    let x_ref_1 = &x;
    let x_ref_2: &i32 = &x; // this is fine because we can have
                            // multiple immutable references.

    let x_mut_ref = &mut x; // this is not fine because it
                            // breaks the rule 2 and 3 mutual
                            // exclusivity.
}

```

Show moreShow more icon

大功告成！编译器警告我，我不能既可变地借用 `x` ，又不可变地借用它。最后，我想打破规则 3（ [打破 Rust 规则 3](#打破-rust-规则-3) ）。

##### 打破 Rust 规则 3

```
fn main() {
    let mut x = 4;
    let x_mut_ref_1 = &x;
    let x_mut_ref_2: &mut i32 = &x;
}

```

Show moreShow more icon

好哇！我再次打破了规则。显然这 3 条规则得到了严格执行，原因只有一个。让多个变量能够通过引用来更改资源，可能引发大量错误！Rust 确实很有趣，也很安全。

## 那么对象呢？

在本节中，我将介绍如何在 Rust 中创建对象。在这个面向对象的世界里，对象似乎必不可少。

### 结构

不同于大多数语言，Rust 没有类。相反，它拥有仅包含数据的结构。以下是 Rust 中的一种典型结构：

```
struct Circle {
    center: (i32, i32),
    radius: u32,
}

```

Show moreShow more icon

我们分解一下该结构，我使用关键字 `struct` 来声明一个自定义结构类型。添加成员变量时，我使用了键-值语法和 `variable_name: type` 格式。我使用逗号分隔多个成员变量。最后，不需要使用分号来终止声明，这与函数非常相像。现在，我已拥有此结构， [Rust 中的结构语法](#rust-中的结构语法) 展示了我如何使用它。

##### Rust 中的结构语法

```
fn main() {
    let c_1 = Circle {
        center: (0, 0),
        radius: 1,
    };

    let c_2: Circle = Circle {
        center: (-1, 1),
        radius: 2,
    };

    println!("c_1's radius is {}”, c_1.radius);
}

struct Circle {
    center: (i32, i32),
    radius: u32,
}

```

Show moreShow more icon

首先， `impl` 关键字告诉 Rust，它必须将以下函数与给定类型（在本例中为 `Circle` ）相关联。接下来，Rust 中的语言习惯是将构造函数命名为”new”。这个 `new` 函数被称为 _关联函数_ 。换言之，类型被用于调用该函数，就像清单 20 中的 `Circle::new` 一样。

创建结构后，如何创建一些方法呢？请看 [Rust 中的示例方法](#rust-中的示例方法) 。

##### Rust 中的示例方法

```
fn main() {
    let mut c = Circle::new((0, 0), 1);

    println!("The circle's area is {}”, c.area());
    println!("The circles location is {:?}”, c.center);

    c.move_to((-1, 1));

    println!("The circles location is {:?}”, c.center);}

struct Circle {
    center: (i32, i32),
    radius: u32,
}

impl Circle {
    fn new(center: (i32, i32), radius: u32) -> Circle {
        Circle {
            center: center,
            radius: radius,
        }
    }

    fn area(&self) -> f64 {
        // converting self.radius, a u32, to an f64.
        let f_radius = self.radius as f64;

        f_radius * f_radius * 3.14159
    }

    fn move_to(&mut self, new_center: (i32, i32)) {
        self.center = new_center;
    }
}

```

Show moreShow more icon

添加方法似乎很直接，但这些新的 `self` 参数是什么？此参数通过点语法传递给该方法（如清单 21 中的 `main` 函数所示），而且它的值是调用该方法的对象。

Rust 中的方法有 3 个唯一变量： `self` 、 `&self` 、 `&mut self` 。回想一下关于所有权和借用的所有知识，这些变量分别将对象所有权转移给该方法，不可变地借用一个对象，再可变地借用一个对象。

现在我已拥有一个包含方法和构造函数的结构，我可以像使用其他编程语言中提供的基本对象一样使用它。还剩最后一种类型没有介绍： `enum` 。

### 枚举

`enum` 是一种表示变化状态的类型。许多 `C/C++` 用户可能认为此类型非常不错，因为它被大量用于表示程序的某种特定状态。下面展示了如何使用 `enum` 关键字创建此类型：

```
enum Color {
    Red,
    Green,
    Blue,
}

```

Show moreShow more icon

这看起来很像一个 `struct` 声明，但没有给变量分配类型。借助这个新的 `enum` ，我可以分配一个该类型的变量，如 [分配变量](#分配变量) 所示。

##### 分配变量

```
fn main() {
    let red = Color::Red;
    let blue: Color = Color::Blue;
}

enum Color {
    Red,
    Green,
    Blue,
}

```

Show moreShow more icon

现在，Rust 在 `enum` 类型上完成了一些独特的工作，比如向变量添加关联数据，但这不属于本文的介绍范畴。

## 其他功能

目前为止，我介绍了开展一个 Rust 项目需要掌握的所有必要知识。本节将介绍一些有助于完善这些知识的细节。

### Match

`match` 是一个关键字，查看用 Rust 编写的程序时会经常看到它。一些人可能将它视为对 `switch/case` 的简单替换，或者认为 `if/else` 代码块是万能的，但他们错了。随着条件变复杂， `if/else` 代码块会严重失控，而且 `C` 和 `C++` 中的 `switch/case` 代码块饱受错误困扰（如果忘记了 default case）。 [match 的简单示例](#match-的简单示例) 提供了 `match` 的一个简单示例。

##### match 的简单示例

```
fn main() {
    let x = 5;

    match x {
        1 => println!("Matched to 1!”),
        2 => println!("Matched to 2!”),
        3 => println!("Matched to 3!”),
        4 => println!("Matched to 4!”),
        5 => println!("Matched to 5!”),
        6 => println!("Matched to 6!”),
        _ => println!("Matched to some other number!”),
    };
}

```

Show moreShow more icon

它的工作原理如下： `match` 接受一个表达式（比如 `x` ）并查找与该值匹配的分支。 `match` 语句的每条”手臂”都采用 `value => expression` 格式，以冒号分隔。找到合适的”手臂”时，就会执行该分支的表达式。 `_ value` 是 Rust 表示”其他任何值”或”默认值”的语法。

`match` 的一个特殊方面在于，它强制要求您全面涵盖要匹配的给定表达式的每种可能结果。换言之，如果我要删除 `_ branch` ，构建将会失败，而且编译器会告知我不是 `i32` 的所有可能值（即从 -2,147,483,648 到 2,147,483,647 的值）都得到了处理。

`match` 的另一种有趣用法是与 `enum` 结合使用，如 [结合 enum 使用 match](#结合-enum-使用-match) 所示。

##### 结合 enum 使用 match

```
fn main() {
    let color = Color::Red;

    match color {
        Color::Red => println!("Got a red color!”),
        Color::Green => println!("Got a green color!”),
        Color::Blue => println!("Got a blue color!”),
    };
}

enum Color {
    Red,
    Green,
    Blue,
}

```

Show moreShow more icon

最后，由于 `match` 是一个表达式，所以我可以使用它分配变量，这很像我之前用作三元表达式的 `if` 代码块 ( [使用 match 分配变量](#使用-match-分配变量))。

##### 使用 match 分配变量

```
let x = 5;

let s = match x {
    1 => "Matched to 1!”,
    2 => "Matched to 2!”,
    3 => "Matched to 3!”,
    4 => "Matched to 4!”,
    5 => "Matched to 5!”,
    6 => "Matched to 6!”,
    _ => "Matched to some other number!”,
};

println!("{}”, s);

```

Show moreShow more icon

### 模块

您可能很震惊我还没有提到名称空间。这是因为 Rust 中没有名称空间。作为替代，Rust 采用了通过 `mod` 关键字创建的模块。您可以按照 [定义模块](#定义模块) 中所示的方式定义模块。

##### 定义模块

```
fn main() {
    let c = Circle {
        center: (0, 0),
        radius: 1,
    };
}

mod shapes {
    struct Circle {
        center: (i32, i32),
        radius: u32,
    }
}

```

Show moreShow more icon

请注意，我将之前定义的 `Circle` 结构放在了 `mod shapes` 代码块中。这意味着 `Circle` 现在属于该名称空间，这正是我编译该程序时构建失败的原因。编译器警告我，类型 `Circle` 未在 `main` 的作用域内定义。我通过使用后跟双冒号的模块名称来修复该问题 ( [定义 Circle](#定义-circle))。

##### 定义 Circle

```
fn main() {
    let c = shapes::Circle {
        center: (0, 0),
        radius: 1,
    };
}

mod shapes {
    struct Circle {
        center: (i32, i32),
        radius: u32,
    }
}

```

Show moreShow more icon

构建仍然失败，但这一次编译器告诉我， `Circle` 结构是私有的。这是 Rust 中的标准，因为模块中的所有类型、变量和函数默认都是私有的。为了关闭此行为，我不得不使用 `pub` 关键字。 [修复 Circle 声明](#修复-circle-声明) 给出了对 `Circle` 声明的修复。

##### 修复 Circle 声明

```
mod shapes {
    pub struct Circle {
        center: (i32, i32),
        radius: u32,
    }
}

```

Show moreShow more icon

该代码修复了 `Circle` 隐私问题，但是现在我被警告成员变量 `center` 和 `radius` 是私有的。现在我将修复此问题：

```
mod shapes {
    pub struct Circle {
        pub center: (i32, i32),
        pub radius: u32,
    }
}

```

Show moreShow more icon

该程序终于开始编译了。但是，这回避了核心问题：我真的想将该模块放在 main.rs 文件中吗？答案是不想。我想将它放在一个单独的文件中。我将在包含 main.rs 的同一个 src 文件夹中创建一个新的 shapes.rs 文件。 在此文件中，我放入这个新的 `Circle` 声明。 [新的 shapes.rs](#新的-shapes-rs) 给出了这些新文件及其内容。

##### 新的 shapes.rs

```
main.rs:
fn main() {
    let c = shapes::Circle {
        center: (0, 0),
        radius: 1,
    };
}

shapes.rs:
pub struct Circle {
    pub center: (i32, i32),
    pub radius: u32,
}

```

Show moreShow more icon

现在我已将 `shapes` 内容转移到 shapes.rs，我如何访问 `Circle` 呢？很容易，可以使用 `mod` 关键字找到具有给定名称的文件，将它们用作一个名称空间。 [main.rs 的内容](#main-rs-的内容) 给出了我的新 main.rs 文件的内容。

##### main.rs 的内容

```
mod shapes;

fn main() {
    let c = shapes::Circle {
        center: (0, 0),
        radius: 1,
    };
}

```

Show moreShow more icon

基本上讲， `mod` 语句告诉 Rust 查找一个 shapes.rs 文件或一个包含 mod.rs 文件的 shapes 文件夹。如果我不想经常使用 `shapes::Circle` ，可以利用 Rust 的一个小诀窍： `use` 关键字。 [use 关键字](#use-关键字) 展示了它的使用方法。

##### use 关键字

```
mod shapes;

use shapes::Circle;

fn main() {
    let c = Circle {
        center: (0, 0),
        radius: 1,
    };
}

```

Show moreShow more icon

现在，我可以使用 `Circle` 了，就像我在此作用域中声明了该结构一样。您可以在 [使用一个哈希映射](#使用一个哈希映射) 中看到更高级的用法，我在其中想要使用一个哈希映射（位于 Rust 的 `std::collections` 模块中）。

##### 使用一个哈希映射

```
use std::collections::HashMap;

fn main() {
    let mut counter = HashMap::new(); // HashMap<char, i32>

    counter.insert('a', 1);
}

```

Show moreShow more icon

呼！基本上完成了。要介绍的最后一点是，如何使用第三方包，也就是 Rust 所称的 _Crate_ 。

### Crate

现在，我已知道如何创建我自行设置的数字变量，但是，如果我想要使用一个随机变量，该怎么做？读遍 Rust 文档，我也没能找到随机变量的标准实现。所以我将使用一个第三方库。在 Google 上快速搜索可以找到 [Crate rand](https://docs.rs/rand/0.6.5/rand/) 。

首先，我必须将 Cargo.toml 文件更改为 [针对 rand 进行了修改的 Cargo.toml](#针对-rand-进行了修改的-cargo-toml) 中的代码。

##### 针对 rand 进行了修改的 Cargo.toml

```
[package]
name = "rusty_start"
version = "0.1.0"
authors = ["Dylan Hicks <dirtgrub.dylanhicks@gmail.com>"]

[dependencies]
rand = "0.4"

```

Show moreShow more icon

通过将 `rand = "0.4"` 添加到我的依赖项，我告诉 Cargo 按此名称和版本来搜索该 crate，并将它编译到我的二进制程序中。 在我的 main.rs 文件中，我使用了 `extern crate` 关键字：

```
extern crate rand;

fn main() {
    let x = 5;
}

```

Show moreShow more icon

我现在可以使用 rand crate 借助 `use` 关键字来生成一个随机数字，如 [使用 rand crate 生成一个随机数字](#使用-rand-crate-生成一个随机数字) 所示。

##### 使用 rand crate 生成一个随机数字

```
extern crate rand;

use rand::Rng;

fn main() {
    let mut rng = rand::thread_rng(); // random number generator
    let x = rng.get::<i32>(); // ignore the ::<i32>
                              // this just tells the rng to make
                              // an i32.

    println!("x: {}”, x);
}

```

Show moreShow more icon

## 最后的思考

难以置信，我成功了！但事实是，我仅介绍了 Rust 功能的一小部分。如果启动一个需要某种系统语言的新项目，或者在 Java 或 `C#` 之间摇摆不定，请考虑使用 Rust 作为替代选择。Rust 快速、有趣、容易理解，比已有的其他任何语言都更安全，而且消除了尽可能多的人为错误。要了解更多信息，请访问 [https://doc.rust-lang.org/book](https://doc.rust-lang.org/book/)，阅读 Mozilla 编写的有关它的著名语言的两本书。

本文翻译自： [Get to know Rust](https://developer.ibm.com/articles/os-know-rust/)（2018-06-14）