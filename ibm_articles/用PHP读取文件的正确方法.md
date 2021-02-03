# 用 PHP 读取文件的正确方法
了解使用 fopen、fclose、feof、fgets、fgetss 和 fscanf 的正确时机

**标签:** PHP,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/os-php-readfiles/)

Roger McCoy

发布: 2007-03-06

* * *

## 让我们算一算有多少种方法

处理诸如 PHP 之类的现代编程语言的乐趣之一就是有大量的选项可用。PHP 可以轻松地赢得 Perl 的座右铭”There’s more than one way to do it”（并非只有一种方法可做这件事），尤其是在文件处理上。但是在这么多可用的选项中，哪一种是完成作业的最佳工具？当然，实际答案取决于解析文件的目标，因此值得花时间探究所有选项。

## 传统的 fopen 方法

`fopen` 方法可能是以前的 C 和 C++ 程序员最熟悉的，因为如果您使用过这些语言，那么它们或多或少都是您已掌握多年的工具。对于这些方法中的任何一种，通过使用 `fopen` （用于读取数据的函数）的标准方法打开文件，然后使用 `fclose` 关闭文件，如清单 1 所示。

##### 清单 1\. 用 fgets 打开并读取文件

```
$file_handle = fopen("myfile", "r");
while (!feof($file_handle)) {
$line = fgets($file_handle);
echo $line;
}
fclose($file_handle);

```

Show moreShow more icon

虽然大多数具有多年编程经验的程序员都熟悉这些函数，但是让我对这些函数进行分解。有效地执行以下步骤：

1. 打开文件。 `$file_handle` 存储了一个对文件本身的引用。
2. 检查您是否已到达文件的末尾。
3. 继续读取文件，直至到达文件末尾，边读取边打印每行。
4. 关闭文件。

记住这些步骤，我将回顾在这里使用的每个文件函数。

### fopen

`fopen` 函数将创建与文件的连接。我之所以说”创建连接”，是因为除了打开文件之外， `fopen` 还可以打开一个 URL：

```
$fh = fopen("http://127.0.0.1/", "r");

```

Show moreShow more icon

这行代码将创建一个与以上页面的连接，并允许您开始像读取一个本地文件一样读取它。

**注：** `fopen` 中使用的 `"r"` 将指示文件以只读方式打开。由于将数据写入文件不在本文的讨论范围内，因此我将不列出所有其他选项。但是，如果是从二进制文件读取以获得跨平台兼容性，则应当将 `"r"` 更改为 `"rb"` 。稍后您将看到这样的示例。

### feof

`feof` 命令将检测您是否已经读到文件的末尾并返回 True 或 False。 [清单 1](#清单-1-用-fgets-打开并读取文件) 中的循环将继续执行，直至您达到文件”myfile”的末尾。注：如果读取的是 URL 并且套接字由于不再有任何数据可以读取而超时，则 `feof` 也将返回 False。

### fclose

向前跳至清单 1 的末尾， `fclose` 将实现与 `fopen` 相反的功能：它将关闭指向文件或 URL 的连接。执行此函数后，您将不再能够从文件或套接字中读取任何信息。

### fgets

在清单 1 中回跳几行，您就到达了文件处理的核心：实际读取文件。 `fgets` 函数是处理第一个示例的首选武器。它将从文件中提取一行数据并将其作为字符串返回。在那之后，您可以打印或者以别的方式处理数据。清单 1 中的示例将精细地打印整个文件。

如果决定限制处理数据块的大小，您可以将一个参数添加到 `fgets` 中限制最大行长度。例如，使用以下代码将行长度限制为 80 个字符：

```
$string = fgets($file_handle, 81);

```

Show moreShow more icon

回想 C 中的”\\0”字符串末尾终止符，将长度设为比实际所需值大一的数字。因而，如果需要 80 个字符，则以上示例使用 81。应养成以下习惯：只要对此函数使用行限制，就添加该额外字符。

### fread

`fgets` 函数是多个文件读取函数中惟一一个可用的。它是一个更常用的函数，因为逐行解析通常会有意义。事实上，几个其他函数也可以提供类似功能。但是，您并非总是需要逐行解析。

这时就需要使用 `fread` 。 `fread` 函数与 `fgets` 的处理目标略有不同：它趋于从二进制文件（即，并非主要包含人类可阅读的文本的文件）中读取信息。由于”行”的概念与二进制文件无关（逻辑数据结构通常都不是由新行终止），因此您必须指定需要读入的字节数。

```
$fh = fopen("myfile", "rb");
$data = fread($file_handle, 4096);

```

Show moreShow more icon

使用二进制数据注意：此函数的示例已经使用了略微不同于 `fopen` 的参数。当处理二进制数据时，始终要记得将 `b` 选项包含在 `fopen` 中。如果跳过这一点，Microsoft Windows 系统可能无法正确处理文件，因为它们将以不同的方式处理新行。如果处理的是 Linux 系统（或其他某个 UNIX 变种），则这可能看似没什么关系。但即使不是针对 Windows 开发的，这样做也将获得良好的跨平台可维护性，并且也是应当遵循的一个好习惯。

以上代码将读取 4,096 字节 (4 KB) 的数据。注：不管指定多少字节， `fread` 都不会读取超过 8,192 个字节 (8 KB)。

假定文件大小不超过 8 KB，则以下代码应当能将整个文件读入一个字符串。

```
$fh = fopen("myfile", "rb");
$data = fread($fh, filesize("myfile"));
fclose($fh);

```

Show moreShow more icon

如果文件长度大于此值，则只能使用循环将其余内容读入。

### fscanf

回到字符串处理， `fscanf` 同样遵循传统的 C 文件库函数。如果您不熟悉它，则 `fscanf` 将把字段数据从文件读入变量中。

```
list ($field1, $field2, $field3) = fscanf($fh, "%s %s %s");

```

Show moreShow more icon

此函数使用的格式字符串在很多地方都有描述（如 PHP.net 中），故在此不再赘述。可以这样说，字符串格式化极为灵活。值得注意的是所有字段都放在函数的返回值中。（在 C 中，它们都被作为参数传递。）

### fgetss

`fgetss` 函数不同于传统文件函数并使您能更好地了解 PHP 的力量。该函数的功能类似于 `fgets` 函数，但将去掉发现的任何 HTML 或 PHP 标记，只留下纯文本。查看如下所示的 HTML 文件。

##### 清单 2\. 样例 HTML 文件

```
<html>
    <head><title>My title</title></head>
    <body>
        <p>If you understand what "Cause there ain't no one for to give you no pain"
            means then you listen to too much of the band America</p>
    </body>
</html>

```

Show moreShow more icon

然后通过 `fgetss` 函数过滤它。

##### 清单 3\. 使用 fgetss

```
$file_handle = fopen("myfile", "r");
while (!feof($file_handle)) {
echo = fgetss($file_handle);
}
fclose($file_handle);

```

Show moreShow more icon

以下是输出：

```
My title

        If you understand what "Cause there ain't no one for to give you no pain"
            means then you listen to too much of the band America

```

Show moreShow more icon

### fpassthru 函数

无论怎样读取文件，您都可以使用 `fpassthru` 将其余数据转储到标准输出通道。

```
fpassthru($fh);

```

Show moreShow more icon

此外，此函数将打印数据，因此无需使用变量获取数据。

## 非线性文件处理：跳跃访问

当然，以上函数只允许顺序读取文件。更复杂的文件可能要求您来回跳转到文件的不同部分。这时就用得着 `fseek` 了。

```
fseek($fh, 0);

```

Show moreShow more icon

以上示例将跳转回文件的开头。如果不需要完全返回 —— 我们可设定返回千字节 —— 然后就可以这样写：

```
fseek($fh, 1024);

```

Show moreShow more icon

从 PHP V4.0 开始，您有一些其他选项。例如，如果需要从当前位置向前跳转 100 个字节，则可以尝试使用：

```
fseek($fh, 100, SEEK_CUR);

```

Show moreShow more icon

类似地，可以使用以下代码向后跳转 100 个字节：

```
fseek($fh, -100, SEEK_CUR);

```

Show moreShow more icon

如果需要向后跳转至文件末尾前 100 个字节处，则应使用 `SEEK_END` 。

```
fseek($fh, -100, SEEK_END);

```

Show moreShow more icon

在到达新位置后，可以使用 `fgets` 、 `fscanf` 或任何其他方法读取数据。

**注：** 不能将 `fseek` 用于引用 URL 的文件处理。

## 提取整个文件

现在，我们将接触到一些 PHP 的更独特的文件处理功能：用一两行处理大块数据。例如，如何提取文件并在 Web 页面上显示其全部内容？好的，您看到了 `fgets` 使用循环的示例。但是如何能够使此过程变得更简单？用 `fgetcontents` 会使过程超级简单，该方法将把整个文件放入一个字符串中。

```
$my_file = file_get_contents("myfilename");
echo $my_file;

```

Show moreShow more icon

虽然它不是最好的做法，但是可以将此命令更简明地写为：

```
echo file_get_contents("myfilename");

```

Show moreShow more icon

本文主要介绍的是如何处理本地文件，但是值得注意的是您还可以用这些函数提取、回显和解析其他 Web 页面。

```
echo file_get_contents("http://127.0.0.1/");

```

Show moreShow more icon

此命令等效于：

```
$fh = fopen("http://127.0.0.1/", "r");
fpassthru($fh);

```

Show moreShow more icon

您一定会查看此命令并认为：”那还是太费力”。PHP 开发人员同意您的看法。因此可以将以上命令缩短为：

```
readfile("http://127.0.0.1/");

```

Show moreShow more icon

`readfile` 函数将把文件或 Web 页面的全部内容转储到默认的输出缓冲区。默认情况下，如果失败，此命令将打印错误消息。要避免此行为（如果需要），请尝试：

```
@readfile("http://127.0.0.1/");

```

Show moreShow more icon

当然，如果确实需要解析文件，则 `file_get_contents` 返回的单个字符串可能有些让人吃不消。您的第一反应可能是用 `split()` 函数将它分解一下。

```
$array = split("\n", file_get_contents("myfile"));

```

Show moreShow more icon

但是既然已经有一个很好的函数为您执行此操作为什么还要这样大费周章？PHP 的 `file()` 函数一步即可完成此操作：它将返回分为若干行的字符串数组。

```
$array = file("myfile");

```

Show moreShow more icon

应当注意的是，以上两个示例有一点细微差别。虽然 `split` 命令将删除新行，但是当使用 `file` 命令（与 `fgets` 命令一样）时，新行仍将被附加到数组中的字符串上。

但是，PHP 的力量还远不止于此。您可以在一条命令中使用 `parse_ini_file` 解析整个 PHP 样式的 .ini 文件。 `parse_ini_file` 命令接受类似清单 4 所示的文件。

##### 清单 4\. 样例 .ini 文件

```
; Comment
[personal information]
name = "King Arthur"
quest = To seek the holy grail
favorite color = Blue

[more stuff]
Samuel Clemens = Mark Twain
Caryn Johnson = Whoopi Goldberg

```

Show moreShow more icon

以下命令将把此文件转储为数组，然后打印该数组：

```
$file_array = parse_ini_file("holy_grail.ini");
print_r $file_array;

```

Show moreShow more icon

以下输出的是结果：

##### 清单 5\. 输出

```
Array
(
    [name] => King Arthur
    [quest] => To seek the Holy Grail
    [favorite color] => Blue
    [Samuel Clemens] => Mark Twain
    [Caryn Johnson] => Whoopi Goldberg
)

```

Show moreShow more icon

当然，您可能注意到此命令合并了各个部分。这是默认行为，但是您可以通过将第二个参数传递给 `parse_ini_file` 轻松地修正它： `process_sections` ，这是一个布尔型变量。将 `process_sections` 设为 True。

```
$file_array = parse_ini_file("holy_grail.ini", true);
print_r $file_array;

```

Show moreShow more icon

并且您将获得以下输出：

##### 清单 6\. 输出

```
Array
(
    [personal information] => Array
        (
            [name] => King Arthur
            [quest] => To seek the Holy Grail
            [favorite color] => Blue
        )

    [more stuff] => Array
        (
            [Samuel Clemens] => Mark Twain
            [Caryn Johnson] => Whoopi Goldberg
        )

)

```

Show moreShow more icon

PHP 将把数据放入可以轻松解析的多维数组中。

对于 PHP 文件处理来说，这只是冰山一角。诸如 `tidy_parse_file` 和 `xml_parse` 之类的更复杂的函数可以分别帮助您处理 HTML 和 XML 文档。有关这些特殊函数的使用细节，请参阅 参考资料 。如果您要处理那些类型的文件，则那些参考资料值得一看，但不必过度考虑本文中谈到的每种可能遇到的文件类型，下面是一些用于处理到目前为止介绍的函数的很好的通用规则。

## 最佳实践

绝不要假定程序中的一切都将按计划运行。例如，如果您要查找的文件已被移动该当如何？如果权限已被改变而无法读取其内容又当如何？您可以通过使用 `file_exists` 和 `is_readable` 预先检查这些问题。

##### 清单 7\. 使用 file\_exists 和 is\_readable

```
$filename = "myfile";
if (file_exists($filename) && is_readable ($filename)) {
    $fh = fopen($filename, "r");
    # Processing
    fclose($fh);
}

```

Show moreShow more icon

但是，在实践中，用这样的代码可能太繁琐了。处理 `fopen` 的返回值更简单并且更准确。

```
if ($fh = fopen($filename, "r")) {
    # Processing
    fclose($fh);
}

```

Show moreShow more icon

由于失败时 `fopen` 将返回 False，这将确保仅当文件成功打开后才执行文件处理。当然，如果文件不存在或者不可读，您可以期望一个负返回值。这将使这个检查可以检查所有可能遇到的问题。此外，如果打开失败，可以退出程序或让程序显示错误消息。

如 `fopen` 函数一样， `file_get_contents` 、 `file` 和 `readfile` 函数都在打开失败或处理文件失败时返回 False。 `fgets` 、 `fgetss` 、 `fread` 、 `fscanf` 和 `fclose` 函数在出错时也返回 False。当然，除 `fclose` 以外，您可能已经对这些函数的返回值都进行了处理。使用 `fclose` 时，即使文件处理未正常关闭，也不会执行什么操作，因此通常不必检查 `fclose` 的返回值。

## 由您来选择

PHP 不缺读取和解析文件的有效方法。诸如 `fread` 之类的典型函数可能在大多数时候都是最佳的选择，或者当 `readfile` 刚好能满足任务需要时，您可能会发现自己更为 `readfile` 的简单所吸引。它实际上取决于所要完成的操作。

如果要处理大量数据， `fscanf` 将能证明自己的价值并比使用 `file` 附带 `split` 和 `sprintf` 命令更有效率。相反，如果要回显只做了少许修改的大量文本，则使用 `file` 、 `file_get_contents` 或 `readfile` 可能更合适。使用 PHP 进行缓存或者创建权宜的代理服务器时可能就属于这种情况。

PHP 给您提供了大量处理文件的工具。深入了解这些工具并了解哪些工具最适合于要处理的项目。您已拥有很多的选择，因此好好地利用它们享受使用 PHP 处理文件的乐趣。

本文翻译自： [The right way to read files with PHP](https://developer.ibm.com/articles/os-php-readfiles/)（2007-02-13）