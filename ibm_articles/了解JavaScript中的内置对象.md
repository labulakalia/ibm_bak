# 了解 JavaScript 中的内置对象
最常用的对象以及它们有哪些功能

**标签:** JavaScript,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/wa-objectsinjs-v1b/)

Kris Hadlock

发布: 2012-07-23

* * *

所有编程语言都具有内部（或内置的）对象来创建 语言的基本功能。内部对象是 您编写自定义代码所用语言的基础， 该代码基于您的想象实现自定义功能。JavaScript 有许多 将其定义为语言的内部对象。本文介绍了一些 最常用的对象，并简要介绍了它们 有哪些功能以及如何使用这些功能。

## Number

JavaScript `Number` 对象是 一个数值包装器。您可以将其与 `new` 关键词结合使用，并将其设置为一个稍后要在 JavaScript 代码中使用的变量：

```
var myNumber = new Number(numeric value);

```

Show moreShow more icon

或者，您可以通过将一个变量设置为一个数值来创建一个 `Number` 对象。然后，该变量将 能够访问该对象可用的属性和方法。

除了存储数值， `Number` 对象包含各种属性和 方法，用于操作或检索关于数字的信息。 `Number` 对象可用的所有属性 都是 _只读常量_ ，这意味着它们的值始终保持 不变，不能更改。有 4 个属性包含在 `Number` 对象里：

- `MAX_VALUE`
- `MIN_VALUE`
- `NEGATIVE_INFINITY`
- `POSITIVE_INFINITY`

`MAX_VALUE` 属性返回 `1.7976931348623157e+308` 值，它是 JavaScript 能够处理的最大数字：

```
document.write(Number.MAX_VALUE);
// Result is: 1.7976931348623157e+308

```

Show moreShow more icon

另外，使用 `MIN_VALUE` 返回 `5e-324` 值，这是 JavaScript 中最小的数字：

```
document.write(Number.MIN_VALUE);
// Result is: 5e-324

```

Show moreShow more icon

`NEGATIVE_INFINITY` 是 JavaScript 能够处理的最大负数，表示为 `-Infinity` ：

```
document.write(Number.NEGATIVE_INFINITY);
// Result is: -Infinity

```

Show moreShow more icon

`POSITIVE_INFINITY` 属性是大于 `MAX_VALUE` 的任意数，表示为 `Infinity` ：

```
document.write(Number.POSITIVE_INFINITY);
// Result is: Infinity

```

Show moreShow more icon

`Number` 对象还有一些方法，您可以 用这些方法对数值进行格式化或进行转换。这些方法包括：

- `toExponential`
- `toFixed`
- `toPrecision`
- `toString`
- `valueOf`

每种方法基本上执行如其名称所暗示的操作。例如， `toExponential` 方法以指数形式返回 数字的字符串表示。每种 方法的独特之处在于它接受的参数。 `toExponential` 方法有一个可选参数， 可用于设置要使用多少有效数字， `toFixed` 方法基于所传递的参数确定小数 精度， `toPrecision` 方法基于所传递的参数确定 要显示的有效数字。

JavaScript 中的每个对象都包含一个 `toString` 和 `valueOf` 方法，因此这些方法 在前面的章节中不介绍。 `toString` 方法返回 数字的字符串表示（在本例中），但是在其他对象中，它返回 相应对象类型的字符串表示。 `valueOf` 方法返回调用它的对象类型的原始值，在本例中为 `Number` 对象。

仅 `Number` 对象似乎并不十分 强大，但它是任何编程语言的一个重要组成部分， JavaScript 也不例外。JavaScript `Number` 对象为任何 数学程序提供基础，这基本上是所有 编程语言的基础。

## Boolean

`Boolean` 在尝试 用 JavaScript 创建任何逻辑时是必要的。 _Boolean_ 是一个 代表 true 或 false 值的对象。 `Boolean` 对象有多个值，这些值 相当于 false 值（ `0` 、 `-0` 、 `null` 或 `""` [一个空字串]），未定义的 (`NaN`)，当然还有 false。所有其他布尔 值相当于 true 值。该对象可以 通过 `new` 关键词进行实例化，但通常是 一个被设为 true 或 false 值的变量：

```
var myBoolean = true;

```

Show moreShow more icon

`Boolean` 对象包括 `toString` 和 `valueOf` 方法，尽管您不太可能需要使用这些方法。 `Boolean` 最常用于在 条件语句中 true 或 false 值的简单判断。 布尔值和条件语句的组合提供了一种使用 JavaScript 创建逻辑的方式。此类条件语句的示例包括 `if` 、 `if...else` 、 `if...else...if` 以及 `switch` 语句。当与 条件语句结合使用时，您可以基于 您编写的条件使用布尔值确定结果。 [清单 1\. 与布尔值相结合的条件语句](#清单-1-与布尔值相结合的条件语句) 显示了 条件语句与布尔值相结合的一个简单示例。

##### 清单 1\. 与布尔值相结合的条件语句

```
var myBoolean = true;
if(myBoolean == true) {
// If the condition evaluates to true
}
else {
// If the condition evaluates to false
}

```

Show moreShow more icon

不言而喻， `Boolean` 对象 是 JavaScript 一个极其重要的组成部分。如果没有 Boolean 对象， 在条件语句内便无法进行判断。

## String

JavaScript `String` 对象是 文本值的包装器。除了存储文本， `String` 对象包含一个属性和各种 方法来操作或收集有关文本的信息。与 `Boolean` 对象类似， `String` 对象不需要进行实例化 便能够使用。例如，您可以将一个变量设置为一个字符串， 然后 `String` 对象的所有属性或 方法都可用于该变量：

```
var myString = "My string";

```

Show moreShow more icon

`String` 对象只有一个 属性，即 `length` ，它是 只读的。 `length` 属性可用于只返回 字符串的长度：您不能在外部修改它。随后的代码 提供了使用 `length` 属性确定一个字符串中的字符数的示例：

```
var myString = "My string";
document.write(myString.length);
// Results in a numeric value of 9

```

Show moreShow more icon

该代码的结果是 `9` ，因为 两个词之间的空格也作为一个字符计算。

在 `String` 对象中有相当多的方法可用于操作和收集有关文本的信息。 以下是可用的方法列表：

- `charAt`
- `charCodeAt`
- `concat`
- `fromCharCode`
- `indexOf`
- `lastIndexOf`
- `match`
- `replace`
- `search`
- `slice`
- `split`
- `substr`
- `substring`
- `toLowerCase`
- `toUpperCase`

`chartAt` 方法可用于基于您作为参数传递的索引检索 特定字符。 下面的代码说明了如何返回 字符串的第一个字符：

```
var myString = "My string";
document.write(myString.chartAt(0);
// Results in M

```

Show moreShow more icon

如果您需要相反的结果，有几个方法 可返回字符串中的指定字符或字符集，而不 使用索引返回字符。这些方法包括 `indexOf` 和 `lastIndexOf` ，这两个方法都包含两个 参数： `searchString` 和 `start` 。 `searchString` 参数是起始索引， `start` 参数告诉方法 从哪里开始搜索。这两个方法之间的区别在于， `indexOf` 返回第一个索引， `lastIndexOf` 返回最后一个索引。

`charCodeAt` 方法类似于 `charAt` ：惟一的区别在于它返回 Unicode 字符。另一种与 Unicode 相关的方法（包括在 `String` 对象中）是 `fromCharCode` ，它将 Unicode 转换为 字符。

如果您想要组合字符串，可以使用加号 (`+`) 将这些字符串加起来，或者您可以 更适当地使用 `concat` 方法。该 方法接受无限数量的字符串参数，连接它们，并 将综合结果作为新字符串返回。 [清单 2](#清单-2-使用-concat-方法合并多个字符串) 展示了如何使用 `concat` 实例将多个字符串合并成一个。

##### 清单 2\. 使用 concat 方法合并多个字符串

```
var myString1 = "My";
var myString2 = " ";
var myString3 = "string";
document.write(myString.concat(myString1, myString2, myString3);
// Results in "My String"

```

Show moreShow more icon

还有一组 `String` 方法 接受正则表达式作为一个参数，以查找或修改一个字符串。 这些包括 `match` 、 `replace` 和 `search` 方法。 `match` 方法使用正则 表达式搜索特定字符串并返回所有的匹配的字符串。 `replace` 方法实际上接受子字符串或 正则表达式和替换字符串作为其第二个参数， 用替换字符串更换所有匹配项，并返回更新的 字符串。这些方法的最后一个是 `search` 方法，它搜索正则表达式的匹配结果并返回其 位置。

如果需要修改字符串，有多个方法派得上用场。 第一个方法是 `slice` 方法，它基于索引或 索引的开始和结尾的组合提取 并返回一个字符串的一部分。另一个方法是 `split` 方法。 `split` 方法每当找到分隔符参数时就将一个字符串分割成一系列 子字符串。例如，如果将逗号 (`,`) 作为一个参数传递，那么字符串 将在每个逗号处分割成一个新的子字符串。能够修改字符串的方法还包括 `substr` 方法，它 基于指定为参数的起始位置和长度，从字符串提取字符， 还有 `substring` 方法，该方法基于指定为参数的两个索引从一个字符串提取字符。能够改变字符串的最后的方法分别是 `toLowerCase` 和 `toUpperCase` ，它们将字符串中的字符分别转换为 小写和大写字母。这些方法在 比较字符串值时非常有用，因为字符串有时可能 大小写不一致。这些方法确保您在比较 值，而不是大小写。

## Date

JavaScript `Date` 对象提供了一种方式 来处理日期和时间。您可以用许多不同的 方式对其进行实例化，具体取决于想要的结果。例如，您可以在没有参数的情况下对其进行实例化：

```
var myDate = new Date();

```

Show moreShow more icon

或传递 `milliseconds` 作为一个参数：

```
var myDate = new Date(milliseconds);

```

Show moreShow more icon

您可以将一个日期字符串作为一个参数传递：

```
var myDate = new Date(dateString);

```

Show moreShow more icon

或者您可以传递多个参数来创建一个完整的日期：

```
var myDate = new Date(year, month, day, hours, minutes, seconds, milliseconds);

```

Show moreShow more icon

此外，有几种方法可用于 `Date` 对象，一旦该对象 得到实例化，您便可以使用这些方法。大多数可用的方法围绕 获取当前时间的特定部分。以下方法是 可用于 `Date` 对象的 getter 方法：

- `getDate`
- `getDay`
- `getFullYear`
- `getHours`
- `getMilliseconds`
- `getMinutes`
- `getMonth`
- `getSeconds`
- `getTime`
- `getTimezoneOffset`

如您所见，每个方法所 返回的值都相当简单。区别在于所返回的值范围。例如， `getDate` 方法返回 一个月份的天数，范围从 1 到 31； `getDay` 方法返回每周的天数，范围从 0 到 6； `getHours` 方法返回小时数值， 范围从 0 到 23； `getMilliseconds` 函数返回毫秒数值，范围从 0 到 999。 `getMinutes` 和 `getSeconds` 方法返回一个范围从 0 到 59 的值， `getMonth` 方法返回一个 从 0 到 11 之间的月份数值。本列表中惟一独特的方法 是 `getTime` 和 `getTimezoneOffset` 。 `getTime` 方法返回 自 1/1/1970 中午 12 点的毫秒数，而 `getTimezoneOffset` 方法返回 格林尼治标准时间和本地时间之间的时间差，以分钟为单位。

对于大多数 getter 方法，还有一个 setter 方法，接受 相应的值范围内的数值参数。setter 方法 如下所示：

- `setDate`
- `setFullYear`
- `setHours`
- `setMilliseconds`
- `setMinutes`
- `setMonth`
- `setSeconds`
- `setTime`

对于上述所有 getter 方法，有一些匹配的方法 返回相同的值范围，只是这些值以 国际标准时间设置。这些方法包括：

- `getUTCDate`
- `getUTCDay`
- `getUTCFullYear`
- `getUTCHours`
- `getUTCMilliseconds`
- `getUTCMinutes`
- `getUTCMonth`
- `getUTCSeconds`

当然，由于对于所有原始 getter 方法都有 setter 方法， 对于国际标准时间也一样。这些方法包括：

- `setUTCDate`
- `setUTCFullYear`
- `setUTCHours`
- `setUTCMilliseconds`
- `setUTCMinutes`
- `setUTCMonth`
- `setUTCSeconds`

正如在本文开头提到的，我不提供许多 关于 `toString` 方法的信息，但是 在 `Date` 对象中有一些方法可将日期转换为一个字符串，值得一提。在某些 情况下，需要将日期或日期的一部分转换为一个 字符串。例如，如果您将其追加到一个字符串或在 比较语句中使用它。有几个方法可用于 `Date` 对象，提供略微不同的 方法将其转换成字符串，包括：

- `toDateString`
- `toLocaleDateString`
- `toLocaleTimeString`
- `toLocaleString`
- `toTimeString`
- `toUTCString`

`toDateString` 方法将日期转换为 字符串：

```
var myDate = new Date();
document.write(myDate.toDateString());

```

Show moreShow more icon

`toDateString` 返回当前日期， 格式为 _Tue Jul 19 2011。_

`toTimeString` 方法将时间从 `Date` 对象转换为字符串：

```
var myDate = new Date();
document.write(myDate.toTimeString());

```

Show moreShow more icon

`toTimeString` 将时间作为字符串返回， 格式为 _23:00:00 GMT-0700 (MST)。_

最后一种将日期转换为字符串的方法是 `toUTCString` ，它将日期转换为 国际标准时间的字符串。

有几种方法使用区域设置将日期转换成字符串，但是在撰写本文之时 Google Chrome 还不支持这几种方法。不支持的方法 包括 `toLocaleDateString` 、 `toLocaleTimeString` 和 `toLocaleString` 。

JavaScript `Date` 对象乍看起来似乎很简单， 但是它不仅仅是一种显示 当前日期的有用方式。它取决于您要创建的功能。 例如， `Date` 对象是 创建倒计时钟表或其他与时间相关的功能的基础。

## Array

JavaScript `Array` 对象是一个存储变量的变量：您可以用它一次在一个变量中存储多个值， 它有许多方法允许您操作或收集 有关它所存储的值的信息。尽管 `Array` 对象不差别对待值类型，但是 在一个单一数组中使用同类值是很好的做法。因此， 在同一数组中使用数字和字符串不是好的做法。所有 可用于 `Array` 对象的属性 都是只读的，这意味着它们的值不能从外部予以更改。

可用于 `Array` 对象的惟一属性 是 `length` 。该属性返回 一个数组中的元素数目，通常在使用 循环迭代数组中的值时用到：

```
var myArray = new Array(1, 2, 3);
for(var i=0; i<myArray.length; i++) {
document.write(myArray[i]);
}

```

Show moreShow more icon

有多种方法可用于 `Array` 对象，您可以使用各种方法来向数组添加元素，或从数组删除元素。 这些方法包括 `pop` 、 `push` 、 `shift` 和 `unshift` 。 `pop` 和 `shift` 方法都从 数组中删除元素。 `pop` 方法删除并返回 一个数组中的最后一个元素，而 `shift` 方法删除并返回一个数组中的第一个元素。相反的 功能可以通过 `push` 和 `unshift` 方法实现，它们将元素添加到 数组中。 `push` 方法将元素作为新元素添加到 数组的结尾，并返回新长度，而 `unshift` 方法将元素添加到 数组的前面，并返回新长度。

在 JavaScript 中对数组进行排序可以通过两个方法实现，其中之一 实际上称为 `sort` 。另一个方法是 `reverse` 。 `sort` 方法的复杂之处在于，它基于可选的 `sort` 函数排列数组。 `sort` 函数可以是 您编写的任何自定义函数。 `reverse` 方法不像 `sort` 那样复杂，尽管它的确通过颠倒元素更改 数组中元素的顺序。

在处理数组时，索引非常重要，因为它们定义 数组中每个元素的位置。有两个方法可基于索引更改 字符串： `slice` 和 `splice` 。 `slice` 方法接受索引或 索引开始和结尾的组合作为参数，然后提取数组的一部分并基于参数将其作为 新数组返回。 `splice` 方法包括 `index` 、 `length` 和 `unlimited element` 参数。该方法基于指定的索引将 元素添加到数组，并基于指定的索引将元素从 数组中删除，或基于指定的索引将元素添加到数组或从 数组删除元素。还有一种方法 可以基于匹配值返回一个索引： `indexOf` 。然后您可以使用该索引截取 或拼接数组。

用任何编程语言编写好代码的关键是编写 有条理的代码。正如其各种方法所示， JavaScript `Array` 对象是一种 组织数据并创建复杂功能的强大方式。

## Math

JavaScript `Math` 对象用于执行 数学函数。它不能加以实例化：您只能依据 `Math` 对象的原样使用它，在没有任何实例的情况下从该对象调用属性和 方法：

```
var pi = Math.PI;

```

Show moreShow more icon

`Math` 对象有许多属性和方法 向 JavaScript 提供数学功能。所有的 `Math` 属性都是只读常量， 包括以下各项：

- `E`
- `LN2`
- `LN10`
- `LOG2E`
- `LOG10E`
- `PI`
- `SQRT1_2`
- `SQRT2`

`E` 属性返回 自然对数的底数的值，或 _欧拉指数。_ 该值是惟一的 实数，以 Leonhard Euler 命名。调用 `E` 属性会产生数字 2.718281828459045。其他两个属性也用于返回自然 对数： `LN2` 和 `LN10` 。 `LN2` 属性返回值为 2 的自然对数，而 `LN10` 属性返回值为 10 的自然 对数。 `LOG2E` 和 `LOG10E` 属性可用于返回 `E` 以 2 或 10 为底的对数。 `LOG2E` 的结果是 1.4426950408889633，而 `LOG10E` 的结果是 0.4342944819032518。通常您不需要 大部分这些属性，除非您正在构建 计算器或其他数学密集型项目。然而， `PI` 和平方根比较常见。 `PI` 方法返回圆周与直径的比率。两个属性返回平方根值： `SQRT1_2` 和 `SQRT2` 。 第一个属性返回 0.5 的平方根，而 `SQRT2` 返回 2 的平方根。

除了这些属性，还有几种方法可用来 返回一个数的不同值。其中每种方法都接受 数值，并根据方法名称返回一个值。 遗憾的是，方法名称不总是显而易见的：

- **`abs` 。** 一个数的 绝对值
- **`acos` 。** 反余弦
- **`asin` 。** 反正弦
- **`atan` 。** 反正切
- **`atan2` 。** 多个数的 反正切
- **`cos` 。** 余弦
- **`exp` 。** 幂
- **`log` 。** 一个数的自然 对数
- **`pow` 。** x 的 y 次方值
- **`sin` 。** 正弦
- **`sqrt` 。** 平方根
- **`tan` 。** 一个角的 正切

有三种方法可用于在 JavaScript 中取整数： `ceil` 、 `floor` 和 `round` 。 `ceil` 方法返回一个数的向上舍入值。该方法在 您需要将数字向上舍入到最接近的整数时非常有用。 `floor` 方法提供 与 `ceil` 相反的功能：它返回 一个数字的向下舍入值。该方法在需要 将数字向下舍入到最近的整数时非常有用。 `round` 方法提供了普通的四舍五入 功能，基于现有的 小数将数字向上或向下舍入。

`Math` 对象中包括的最后三个方法分别是 `max` 、 `min` 和 `random` 。 `max` 方法接受多个数字参数并返回最高值， 而 `min` 方法接受多个数字 参数并返回最低值。这些方法在 比较拥有数值的变量时非常有用，特别是当您事先不 知道是什么数值时。您使用 `random` 方法返回 0 与 1 之间的一个随机数。您可以将该方法用作多种目的，比如在 网站主页上显示一个随机图像，或返回一个随机数， 该随机数可用作包含图像的文件路径的数组的一个索引。 从该数组选择的随机图像文件路径然后可 用于将该图像写到 HTML `<img>` 标记。

## 结束语

JavaScript 提供的属性和方法仅仅是可以实现的 功能的开始：是您的想象力创建了 自定义功能。由于您的想象力没有界限，因此 您编写的代码也没有界限。JavaScript 是一种灵活的语言，这有时使它 名声较差，但是往好的一面看，它也向您提供了 快速、创造性地编写代码的能力。如需进一步了解 JavaScript 对象，以及如何使用 JavaScript 语言创建您自己的自定义对象，请务必查看参考资源部分。