# 五种常见的 PHP 设计模式
设计模式对于每个人都非常有用

**标签:** PHP,Web 开发,数据库,软件开发

[原文链接](https://developer.ibm.com/zh/articles/os-php-designptrns/)

Jack Herrington

发布: 2006-10-19

* * *

_设计模式_ 一书将设计模式引入软件社区，该书的作者是 Erich Gamma、Richard Helm、Ralph Johnson 和 John Vlissides Design（俗称 “四人帮”）。所介绍的设计模式背后的核心概念非常简单。经过多年的软件开发实践，Gamma 等人发现了某些具有固定设计的模式，就像建筑师设计房子和建筑物一样，可以为浴室的位置或厨房的构造方式开发模板。使用这些模板或者说 _设计模式_ 意味着可以更快地设计更好的建筑物。同样的概念也适用于软件。

设计模式不仅代表着更快开发健壮软件的有用方法，而且还提供了以友好的术语封装大型理念的方法。例如，您可以说您正在编写一个提供松散耦合的消息传递系统，也可以说你正在编写名称为 _观察者_ 的模式。

用较小的示例展示模式的价值是非常困难的。这往往有些大材小用的意味，因为模式实际上是在大型代码库中发挥作用的。本文不展示大型应用程序，所以您需要思索的是在您自己的大型应用程序中应用示例原理的方法 —— 而不是本文演示的代码本身。这不是说您不应该在小应用程序中使用模式。很多良好的应用程序都以小应用程序为起点，逐渐发展到大型应用程序，所以没有理由不以此类扎实的编码实践为基础。

既然您已经了解了设计模式以及它们的有用之处，现在我们来看看 PHP V5 的五种常用模式。

## 工厂模式

最初在 _设计模式_ 一书中，许多设计模式都鼓励使用 _松散耦合_ 。要理解这个概念，让我们最好谈一下许多开发人员从事大型系统的艰苦历程。在更改一个代码片段时，就会发生问题，系统其他部分 —— 您曾认为完全不相关的部分中也有可能出现级联破坏。

该问题在于 _紧密耦合_ 。系统某个部分中的函数和类严重依赖于系统的其他部分中函数和类的行为和结构。您需要一组模式，使这些类能够相互通信，但不希望将它们紧密绑定在一起，以避免出现联锁。

在大型系统中，许多代码依赖于少数几个关键类。需要更改这些类时，可能会出现困难。例如，假设您有一个从文件读取的 `User` 类。您希望将其更改为从数据库读取的其他类，但是，所有的代码都引用从文件读取的原始类。这时候，使用工厂模式会很方便。

_工厂模式_ 是一种类，它具有为您创建对象的某些方法。您可以使用工厂类创建对象，而不直接使用 `new` 。这样，如果您想要更改所创建的对象类型，只需更改该工厂即可。使用该工厂的所有代码会自动更改。

清单 1 显示工厂类的一个示列。等式的服务器端包括两个部分：数据库和一组 PHP 页面，这些页面允许您添加反馈、请求反馈列表并获取与特定反馈相关的文章。

##### 清单 1\. Factory1.php

```
<?php
interface IUser
{
function getName();
}

class User implements IUser
{
public function __construct( $id ) { }

public function getName()
{
    return "Jack";
}
}

class UserFactory
{
public static function Create( $id )
{
    return new User( $id );
}
}

$uo = UserFactory::Create( 1 );
echo( $uo->getName()."\n" );
?>

```

Show moreShow more icon

`IUser` 接口定义用户对象应执行什么操作。 `IUser` 的实现称为 `User` ， `UserFactory` 工厂类则创建 `IUser` 对象。此关系可以用图 1 中的 UML 表示。

##### 图 1\. 工厂类及其相关 IUser 接口和用户类

![工厂类及其相关 IUser 接口和用户类](../ibm_articles_img/os-php-designptrns_images_factory1.gif)

如果您使用 `php` 解释器在命令行上运行此代码，将得到如下结果：

```
% php factory1.php
Jack
%

```

Show moreShow more icon

测试代码会向工厂请求 `User` 对象，并输出 `getName` 方法的结果。

有一种工厂模式的变体使用工厂方法。类中的这些公共静态方法构造该类型的对象。如果创建此类型的对象非常重要，此方法非常有用。例如，假设您需要先创建对象，然后设置许多属性。此版本的工厂模式会将该进程封装在单个位置中，这样，不用复制复杂的初始化代码，也不必将复制好的代码在在代码库中到处粘贴。

清单 2 显示使用工厂方法的一个示例。

##### 清单 2\. Factory2.php

```
<?php
interface IUser
{
function getName();
}

class User implements IUser
{
public static function Load( $id )
{
        return new User( $id );
}

public static function Create( )
{
        return new User( null );
}

public function __construct( $id ) { }

public function getName()
{
    return "Jack";
}
}

$uo = User::Load( 1 );
echo( $uo->getName()."\n" );
?>

```

Show moreShow more icon

这段代码要简单得多。它仅有一个接口 `IUser` 和一个实现此接口的 `User` 类。 `User` 类有两个创建对象的静态方法。此关系可用图 2 中的 UML 表示。

##### 图 2\. IUser 接口和带有工厂方法的 user 类

![IUser 接口和带有工厂方法的用户类](../ibm_articles_img/os-php-designptrns_images_factory2.gif)

在命令行中运行脚本产生的结果与清单 1 的结果相同，如下所示：

```
% php factory2.php
Jack
%

```

Show moreShow more icon

如上所述，有时此类模式在规模较小的环境中似乎有些大材小用。不过，最好还是学习这种扎实的编码形式，以便应用于任意规模的项目中。

## 单元素模式

某些应用程序资源是 _独占的_ ，因为有且只有一个此类型的资源。例如，通过数据库句柄到数据库的连接是独占的。您希望在应用程序中共享数据库句柄，因为在保持连接打开或关闭时，它是一种开销，在获取单个页面的过程中更是如此。

单元素模式可以满足此要求。如果应用程序每次包含且仅包含一个对象，那么这个对象就是一个 _单元素_ （Singleton）。清单 3 中的代码显示了 PHP V5 中的一个数据库连接单元素。

##### 清单 3\. Singleton.php

```
<?php
require_once("DB.php");

class DatabaseConnection
{
public static function get()
{
    static $db = null;
    if ( $db == null )
      $db = new DatabaseConnection();
    return $db;
}

private $_handle = null;

private function __construct()
{
    $dsn = 'mysql://root:password@localhost/photos';
    $this->_handle =& DB::Connect( $dsn, array() );
}

public function handle()
{
    return $this->_handle;
}
}

print( "Handle = ".DatabaseConnection::get()->handle()."\n" );
print( "Handle = ".DatabaseConnection::get()->handle()."\n" );
?>

```

Show moreShow more icon

此代码显示名为 `DatabaseConnection` 的单个类。您不能创建自已的 `DatabaseConnection` ，因为构造函数是专用的。但使用静态 `get` 方法，您可以获得且仅获得一个 `DatabaseConnection` 对象。此代码的 UML 如图 3 所示。

##### 图 3\. 数据库连接单元素

![数据库连接单元素](../ibm_articles_img/os-php-designptrns_images_singleton.gif)

在两次调用间， `handle` 方法返回的数据库句柄是相同的，这就是最好的证明。您可以在命令行中运行代码来观察这一点。

```
% php singleton.php
Handle = Object id #3
Handle = Object id #3
%

```

Show moreShow more icon

返回的两个句柄是同一对象。如果您在整个应用程序中使用数据库连接单元素，那么就可以在任何地方重用同一句柄。

您可以使用全局变量存储数据库句柄，但是，该方法仅适用于较小的应用程序。在较大的应用程序中，应避免使用全局变量，并使用对象和方法访问资源。

## 观察者模式

观察者模式为您提供了避免组件之间紧密耦合的另一种方法。该模式非常简单：一个对象通过添加一个方法（该方法允许另一个对象，即 _观察者_ 注册自己）使本身变得可观察。当可观察的对象更改时，它会将消息发送到已注册的观察者。这些观察者使用该信息执行的操作与可观察的对象无关。结果是对象可以相互对话，而不必了解原因。

一个简单示例是系统中的用户列表。清单 4 中的代码显示一个用户列表，添加用户时，它将发送出一条消息。添加用户时，通过发送消息的日志观察者可以观察此列表。

##### 清单 4\. Observer.php

```
<?php
interface IObserver
{
function onChanged( $sender, $args );
}

interface IObservable
{
function addObserver( $observer );
}

class UserList implements IObservable
{
private $_observers = array();

public function addCustomer( $name )
{
    foreach( $this->_observers as $obs )
      $obs->onChanged( $this, $name );
}

public function addObserver( $observer )
{
    $this->_observers []= $observer;
}
}

class UserListLogger implements IObserver
{
public function onChanged( $sender, $args )
{
    echo( "'$args' added to user list\n" );
}
}

$ul = new UserList();
$ul->addObserver( new UserListLogger() );
$ul->addCustomer( "Jack" );
?>

```

Show moreShow more icon

此代码定义四个元素：两个接口和两个类。 `IObservable` 接口定义可以被观察的对象， `UserList` 实现该接口，以便将本身注册为可观察。 `IObserver` 列表定义要通过怎样的方法才能成为观察者， `UserListLogger` 实现 `IObserver` 接口。图 4 的 UML 中展示了这些元素。

##### 图 4\. 可观察的用户列表和用户列表事件日志程序

![可观察的用户列表和用户列表事件日志程序](../ibm_articles_img/os-php-designptrns_images_observer.gif)

如果在命令行中运行它，您将看到以下输出：

```
% php observer.php
'Jack' added to user list
%

```

Show moreShow more icon

测试代码创建 `UserList` ，并将 `UserListLogger` 观察者添加到其中。然后添加一个消费者，并将这一更改通知 `UserListLogger` 。

认识到 `UserList` 不知道日志程序将执行什么操作很关键。可能存在一个或多个执行其他操作的侦听程序。例如，您可能有一个向新用户发送消息的观察者，欢迎新用户使用该系统。这种方法的价值在于 `UserList` 忽略所有依赖它的对象，它主要关注在列表更改时维护用户列表并发送消息这一工作。

此模式不限于内存中的对象。它是在较大的应用程序中使用的数据库驱动的消息查询系统的基础。

## 命令链模式

_命令链_ 模式以松散耦合主题为基础，发送消息、命令和请求，或通过一组处理程序发送任意内容。每个处理程序都会自行判断自己能否处理请求。如果可以，该请求被处理，进程停止。您可以为系统添加或移除处理程序，而不影响其他处理程序。清单 5 显示了此模式的一个示例。

##### 清单 5\. Chain.php

```
<?php
interface ICommand
{
function onCommand( $name, $args );
}

class CommandChain
{
private $_commands = array();

public function addCommand( $cmd )
{
    $this->_commands []= $cmd;
}

public function runCommand( $name, $args )
{
    foreach( $this->_commands as $cmd )
    {
      if ( $cmd->onCommand( $name, $args ) )
        return;
    }
}
}

class UserCommand implements ICommand
{
public function onCommand( $name, $args )
{
    if ( $name != 'addUser' ) return false;
    echo( "UserCommand handling 'addUser'\n" );
    return true;
}
}

class MailCommand implements ICommand
{
public function onCommand( $name, $args )
{
    if ( $name != 'mail' ) return false;
    echo( "MailCommand handling 'mail'\n" );
    return true;
}
}

$cc = new CommandChain();
$cc->addCommand( new UserCommand() );
$cc->addCommand( new MailCommand() );
$cc->runCommand( 'addUser', null );
$cc->runCommand( 'mail', null );
?>

```

Show moreShow more icon

此代码定义维护 `ICommand` 对象列表的 `CommandChain` 类。两个类都可以实现 `ICommand` 接口 —— 一个对邮件的请求作出响应，另一个对添加用户作出响应。 图 5 给出了 UML。

##### 图 5\. 命令链及其相关命令

![命令链及其相关命令](../ibm_articles_img/os-php-designptrns_images_chain.gif)

如果您运行包含某些测试代码的脚本，则会得到以下输出：

```
% php chain.php
UserCommand handling 'addUser'
MailCommand handling 'mail'
%

```

Show moreShow more icon

代码首先创建 `CommandChain` 对象，并为它添加两个命令对象的实例。然后运行两个命令以查看谁对这些命令作出了响应。如果命令的名称匹配 `UserCommand` 或 `MailCommand` ，则代码失败，不发生任何操作。

为处理请求而创建可扩展的架构时，命令链模式很有价值，使用它可以解决许多问题。

## 策略模式

我们讲述的最后一个设计模式是 _策略_ 模式。在此模式中，算法是从复杂类提取的，因而可以方便地替换。例如，如果要更改搜索引擎中排列页的方法，则策略模式是一个不错的选择。思考一下搜索引擎的几个部分 —— 一部分遍历页面，一部分对每页排列，另一部分基于排列的结果排序。在复杂的示例中，这些部分都在同一个类中。通过使用策略模式，您可将排列部分放入另一个类中，以便更改页排列的方式，而不影响搜索引擎的其余代码。

作为一个较简单的示例，清单 6 显示了一个用户列表类，它提供了一个根据一组即插即用的策略查找一组用户的方法。

##### 清单 6\. Strategy.php

```
<?php
interface IStrategy
{
function filter( $record );
}

class FindAfterStrategy implements IStrategy
{
private $_name;

public function __construct( $name )
{
    $this->_name = $name;
}

public function filter( $record )
{
    return strcmp( $this->_name, $record ) <= 0;
}
}

class RandomStrategy implements IStrategy
{
public function filter( $record )
{
    return rand( 0, 1 ) >= 0.5;
}
}

class UserList
{
private $_list = array();

public function __construct( $names )
{
    if ( $names != null )
    {
      foreach( $names as $name )
      {
        $this->_list []= $name;
      }
    }
}

public function add( $name )
{
    $this->_list []= $name;
}

public function find( $filter )
{
    $recs = array();
    foreach( $this->_list as $user )
    {
      if ( $filter->filter( $user ) )
        $recs []= $user;
    }
    return $recs;
}
}

$ul = new UserList( array( "Andy", "Jack", "Lori", "Megan" ) );
$f1 = $ul->find( new FindAfterStrategy( "J" ) );
print_r( $f1 );

$f2 = $ul->find( new RandomStrategy() );
print_r( $f2 );
?>

```

Show moreShow more icon

此代码的 UML 如图 6 所示。

##### 图 6\. 用户列表和用于选择用户的策略

![用户列表和用于选择用户的策略](../ibm_articles_img/os-php-designptrns_images_strategy.gif)

`UserList` 类是打包名称数组的一个包装器。它实现 `find` 方法，该方法利用几个策略之一来选择这些名称的子集。这些策略由 `IStrategy` 接口定义，该接口有两个实现：一个随机选择用户，另一个根据指定名称选择其后的所有名称。运行测试代码时，将得到以下输出：

```
% php strategy.php
Array
(
    [0] => Jack
    [1] => Lori
    [2] => Megan
)
Array
(
    [0] => Andy
    [1] => Megan
)
%

```

Show moreShow more icon

测试代码为两个策略运行同一用户列表，并显示结果。在第一种情况中，策略查找排列在 `J` 后的任何名称，所以您将得到 Jack、Lori 和 Megan。第二个策略随机选取名称，每次会产生不同的结果。在这种情况下，结果为 Andy 和 Megan。

策略模式非常适合复杂数据管理系统或数据处理系统，二者在数据筛选、搜索或处理的方式方面需要较高的灵活性。

## 结束语

本文介绍的仅仅是 PHP 应用程序中使用的几种最常见的设计模式。在 _设计模式_ 一书中演示了更多的设计模式。不要因架构的神秘性而放弃。模式是一种绝妙的理念，适用于任何编程语言、任何技能水平。