# Sails 中的路由和控制器
使用控制器自定义和管理复杂操作

**标签:** Node.js,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/wa-sailsjs4/)

Ted Neward

发布: 2016-11-01

* * *

**关于本系列**

像同类产品 Ruby on Rails 一样，Sails.js 是一个旨在帮助开发人员构建 Web 应用程序的框架。Rails 默认情况下用于构建在服务器上生成 HTML 并将其发回给浏览器的应用程序，与它不同的是，Sails 为构建 HTTP API 提供了很好的开箱即用支持。这种设计有助于在客户端和服务器之间实现更好的关注点分离，并促进客户端和服务器之间的更多的互操作性。在本系列中，Web 开发人员和教师 Ted Neward 将向您介绍 Sails.js。HTTP API 示例应用程序可以是 AngularJS、React 或几乎其他任何客户端应用程序的后端。

欢迎！Sails 粉丝们！在我们的 [上一期教程](http://www.ibm.com/developerworks/cn/web/wa-sailsjs3/index.html) 中，您学习了如何在 Sails.js 中将模型与其他模型关联。使用显式关联将通用关系非常好地集成到代码中，但不是所有操作都适合放入一个模型实例中。本期教程将介绍 Sails 控制器，您将使用它们来管理 Sails 应用程序中的更复杂操作。

## 自定义模型操作

目前为止，我们的进展非常顺利，我们使用了 Sails 的默认路由来访问或修改模型实例。这些默认设置（包含在 Sails Blueprint API 中）负责我们期望从 Web 或移动应用程序获得的基本的创建（create）、读取（read）、更新（update）和删除（delete）功能。但是所有开发生产 HTTP API 的开发人员都会告诉您，简单的 CRUD 只有这点能耐。即使您从这里着手，也需要能够定制基本路由与控制器的映射。

Blueprint API 对您很有帮助，但最终您需要某个更强大、灵活和可自定义（或同时具备这三种特性）的工具来构建您的客户和用户想要的应用程序。对于大部分开发周期，都可以使用蓝图路由来建立原型，然后将此框架替换为自定义的路由和关联的控制器。

除了 CRUD 路由，一个常规 Sails.js 安装中已经预先定义了其他一些控制器-路由组合。但是，从很大程度上讲，您往往想创建自己的映射来获得所需的行为。

## 映射复杂的查询

在 [上一篇教程](http://www.ibm.com/developerworks/cn/web/wa-sailsjs3/index.html) 中，您将我们最开始拥有的博客 API 扩展为了一个更庞大的内容管理系统（CMS）后端。尽管您目前构想的是让这个应用程序为一个网络博客提供支持，但它还可以用于其他用途。该 RESTful API 可供几乎任何想要获取并显式博客文章或 RSS 提要的前端应用程序访问，而且它允许搜索查询。扩展该 API 后得到了多种新模型类型，分别是 `Author` 、 `Entry` 和 `Comment` 。每个 `Entry` 拥有一个 `Author` 和 `Comment` ，这些类型也链接回它们引用的 `Entry` 。

切换成 CMS 后，越来越明确地表明您希望能够设置 Sails 的默认路由所不支持的操作。假设您想能够发出复杂的查询（比如针对一个由特定作者在特定日期后编写的，按给定条件进行组织的文章），或者获取与一个特定标签对应的所有文章。您可能需要提供的不是采用原生 JSON 的格式，而是采用兼容 RSS 或 Atom 的 XML 格式的文章。在输入端，您可能希望添加一个导入功能，使 CMS 能够通过 RSS 或 Atom 提要获取并存储完整的现有博客。

实际上，您可以做很多事情，而且许多事情都是 Sails 默认不支持的。

## 控制器的作用

在传统 MVC 模式中， _控制器_ （controller）定义模型（model）与其视图（view）之间的交互。当某个模型的数据发生更改时，控制器会确保附加到该模型的每个视图都会相应地更新。用户在视图内执行某项操作时，控制器也会获得通知。如果该操作必须更改某个模型，控制器会向所有受影响的视图发出通知。从架构角度讲，模型、控制器和视图都是在服务器上定义的。视图通常是某种形式的 HTML 模板（理想情况下包含极少的代码）；模型是域类或对象；控制器是路由背后的代码块。

在 HTTP API 中，三个应用程序组件之间的关系类似但不等同于 MVC。不同于 MVC 架构，HTTP API 的模型、视图和控制器通常并不都包含在同一个服务器上。具体地讲，视图通常位于服务器外，会是一个单页 Web 应用程序或移动应用程序的形式。

在架构上，HTTP API 或多或少与 MVC 应用程序有些类似：

**文档**：目前，对于错误场景的含义，HTTP API 设计中还没有标准。错误编码 500 究竟是表明请求未被处理，还是表明它已被处理，但由于执行请求的逻辑中的内部错误而失败了？所以应该始终显式地文档化您的 HTTP API 中的错误编码的含义。

- _视图_ 显示了通过 HTTP API 收到的数据。
- _模型_ 是通过线路交换的工件。（HTTP API 使用了一种专为简化传输而设计的模型，该模型有时称为 _ViewModel_ 。）
- _控制器_ 是 API 的动词形式；它们的存在是为”执行”某个操作而不是”成为”某个东西。

在 HTTP API 中，当一个视图将 HTTP 请求传到控制器端点时，该请求直接依靠控制器来执行。控制器获取通过 URL 或请求正文传入的数据，对一个或多个现有模型执行某个操作（或创建新模型），并生成响应返回给视图来进行更新。控制器的响应通常是 API 开发人员定义的 HTTP 状态代码和 JSON 主体的组合。

理论介绍得已经足够多了，现在让我们开始实现控制吧。

## 创建一个控制器

开始使用 Sails 控制器的最简单方法是，创建一个返回静态数据的控制器。在本例中，您将创建一个简单的控制器，它返回 CMS API 的一个用户友好的版本。这将是一个运行迅速且容易使用的 API，业务客户可使用它测试服务器是否在正常运行。客户还能够大体了解任何最近的更改，比如对 API 的升级，并相应地调整其 Web 或移动前端。

可通过两种方式在 Sails 中创建新控制器：可以使用 `sails` 命令生成框架，或者可以让这些文件为您生成控制器。在后一种情况下，生成器使用 Sails 约定规则所定义的命名系统来识别和放置文件，这些文件通常是空的。对于第一个控制器，我们将会使用生成器：

```
~$ sails generate controller System
info: Created a new controller ("System") at api/controllers/SystemController.js!

```

Show moreShow more icon

控制器位于您的 Sails 项目的 `api/controllers` 目录中。根据约定，它们拥有 `controller` 后缀。如果您在控制器名称后附加额外的描述符，Sails 会假设这些是要在控制器上执行的操作（方法）。通过提前指定操作，您可以节省一些步骤；否则默认生成器会生成一个空的控制器，如这里所示：

```
/**
* SystemController
*
* @description :: Server-side logic for managing the System
* @help        :: See http://sailsjs.org/#!/documentation/concepts/Controllers
*/

module.exports = {

};

```

Show moreShow more icon

运行一个稍微不同的命令—— `sails generate controller System version` ——会留下更多工作让您处理：

```
/**
* SystemController
*
* @description :: Server-side logic for managing Systems
* @help        :: See http://sailsjs.org/#!/documentation/concepts/Controllers
*/

module.exports = {

    /**
     * `SystemController.version()`
     */
    version: function (req, res) {
      return res.json({
        todo: 'version() is not implemented yet!'
      });
    }
};

```

Show moreShow more icon

添加命令会对搭建的方法端点进行布局，但您可以看到，这些命令目前未执行太多工作，仅返回有帮助的错误消息。使用该生成器创建已建立框架的控制器，最初可能很有帮助。随着时间的推移，许多开发人员最终仅在自己最喜欢用的文本编辑器中执行 `File|New` 。如果您决定这么做，则需要正确地命名该文件和端点（例如 `foocontroller.js` 中的 `FooController` ），并手动编写导出的函数。这些方式都没有对错，所以请使用最适合您的方式。

实现第一个简单的控制器非常容易：

```
module.exports = {

/**
* `SystemController.version()`
*/
version: function (req, res) {
    return res.json({
      version: '0.1'
    });
}
};

```

Show moreShow more icon

### 绑定和调用控制器

接下来，您希望绑定您的控制器，然后调用它。 _绑定_ （Binding）将控制器映射到一个路由； _调用_ （invoking）向该路由发出合适的 HTTP 请求。基于此控制器的文件名和所调用方法的名称，此控制器的默认路由将是 `/System/version` 。

就个人而言，我不喜欢将”`system` ”放在 URL 中，而是更喜欢使用”`/version` ”。Sails 允许将控制器的调用绑定到选择的任何路由，所以我们将它设置为绑定到 `/version` 。

您可以在 Sails 路由表中创建一个条目来设置控制器路由，该路由表存储在 `config/routes.js` 中。只要您创建一个 Sails 应用程序，就会生成这个默认文件。除去所有注释后，此文件几乎是空的：

```
module.exports.routes = {

    /***************************************************************************
    *                                                                          *
    * Make the view located at `views/homepage.ejs` (or `views/homepage.jade`, *
    * etc. depending on your default view engine) your home page.              *
    *                                                                          *
    * (Alternatively, remove this and add an `index.html` file in your         *
    * `assets` directory)                                                      *
    *                                                                          *
    ***************************************************************************/

    '/': {
      view: 'homepage'
    }

};

```

Show moreShow more icon

大体上讲， `config/routes.js` 包含一组 _路由_ （routes）和 _目标_ （targets）。路由是相对 URL，目标是您希望 Sails 调用的对象。默认路由是 `/` URL 模式，它调出 Sails 的默认主页。如果您愿意的话，还可以将此默认路由替换为一个常规 HTML 页面。这样，该 HTML 将位于您的 `assets` 目录中，就在我们将一直使用的 `api` 目录旁边。

升级 CMS 的 HTML 对于像 [React.js](https://facebook.github.io/react/) 这样的单页应用程序框架可能是一种有趣的练习，但这里的目的是添加一个 `/version` 路由，并让它指向 `SystemController.version` 方法。为此，您需要向 `config.js` 所导入的 JSON 对象添加额外的一行代码：

```
module.exports.routes = {

'get /version': 'SystemController.version'

};

```

Show moreShow more icon

将此路由保存到 `config.js` 文件中后，向 `/version` 发出 HTTP GET 请求会发回您之前设置的相同 JSON 结果。

### Blueprint API 中的控制器

创建模型时，Sails 的工具会自动为该模型创建一个控制器。所以，即使您没有打算创建它们，您的每个方法也已经有一个控制器： `AuthorController` 、 `EntryController` 等。

为了进一步熟悉 Sails 控制器和路由，我们假设您希望 `AuthorController` 持有一个方法，该方法返回您 CMS 中所有作者的简历的聚合列表。该实现非常简单——只需要获取所有 `Authors` ，提取它们的简历，并传回该列表：

```
module.exports = {

bios: function(req, res) {
    Author.find({})
      .then(function (authors) {
        console.log("authors = ",authors);
        var bs = [];
        authors.forEach(function (author) {
          bs.push({
            name: author.fullName,
            bio: author.bio
          });
        });
        res.json(bs);
      })
      .catch(function (err) {
        console.log(err);
        res.status(500)
          .json({ error: err });
      });
}

};

```

Show moreShow more icon

如果您使用了默认的 `/author/bios` 路由，您的 `routes.js` 中甚至不需要特殊条目。在这种情况下，默认条目就够用了（如果您不这么认为，您知道如何更改它），所以我们暂时保留默认路由。

### 测试状态

您可能发现 _种子控制器_ （seed controller）对一些测试很有用。这种控制器将数据库初始化为一种已知状态，就象这样：

```
module.exports = {
    run: function(req, res) {
        Author.create({
            fullName: "Fred Flintstone",
            bio: "Lives in Bedrock, blogs in cyberspace",
            username: "fredf",
            email: "fred@flintstone.com"
        }).exec(function (err, author) {
            Entry.create({
                title: "Hello",
                body: "Yabba dabba doo!",
                author: author
            }).exec(function (err, created) {
                Entry.create({
                    title: "Quit",
                    body: "Mr Slate is a jerk",
                    author: author.id
                }).exec(function (err, created) {
                    return res.send("Database seeded");
                });
            });
        });
    }
};

```

Show moreShow more icon

在这种情况下，该已知状态被绑定到控制器的默认路由 `/seed/run` 。您还可以为不同的测试和/或开发场景设置不同的种子方法（seed method）。对于所关注的路由，可使用 `curl` 命令将数据库设置为特定状态。但需要确保您在生产代码中禁用或删除了这些路由。

## 管理控制器输入

现在，您已拥有一个没有输入的非常简单的控制器。但是，控制器通常需要从调用方获取输入。有三种类型的控制器输入：

1. _请求正文中发送的表单参数_ 。这是通过 Web 接受输入的传统机制。
2. _通过请求正文中的 JSON 对象发送的输入数据_ 。此概念与表单参数相同，但发送的内容类型为 `application/json` ，而不是 `form/multipart-form-data` 。客户端通常更容易生成输入数据，而且服务器也更容易使用。
3. _通过参数指定并通过 URL 路由中的占位符发送的输入_ 。请求某个特定作者的文章时，您通常希望在 URL 自身中传递作者标识符。一个示例是 `/author/1/entries` ，其中的”1”是作者的唯一标识符。这样，您就保留了博客文章包含在作者资源中的外观，即使这些文章在物理上未与该作者存储在一起。（示例应用程序就是如此， `Entry` 对象存储在一个与 `Author` 对象不同的集合或表中。）

表单参数属于传统的 Express 样式 `request.getParam()` 调用的领域，已在其他地方具有明确规定。而且 HTTP API 也不经常使用表单参数，所以我们暂时放弃该方法。第二种方法非常适合 CMS 应用程序。

### 获取输入

捕获通过 JSON 对象发送的值通常很简单，只需使用一个 `request.body.field` 。如果输入是一个位于 JSON 最高层级的数组，您则可以使用 `request.body[idx].field` 。

首先，您将创建一个端点，它将返回 CMS 数据库中所有 `Entry` 对象的 RSS XML 提要。（在实际的系统中，需要限制此数字来支持分页模式，比如返回最近的 20 篇文章，但我们暂时保持简单即可。）命名您的控制器（我在下面选择了 `FeedController` ），并在它之上放置一个 RSS 方法，使默认路由（ `/feed/rss` ）变得有意义：

```
var generateRSS = function(entries) {
var rss =
    '<rss version="2.0">' +
    '<channel>' +
    '<title>SailsBlog</title>';

// Items
entries.forEach(function (entry) {
    rss += '<item>' +
      '<title>' + entry.title + '</title>' +
      '<description>' + entry.body + '</description>' +
      '</item>';
});

// Closing
rss += '</channel>' +
    '</rss>';
return rss;
}

module.exports = {

rss: function (req, res) {
    Entry.find({})
      .then(function (entries) {
        var rss = generateRSS(entries);
        res.type("rss");
        res.status(200);
        res.send(rss);
      })
      .catch(function (err) {
        console.log(err);
        res.status(500)
          .json({ error: err });
      });

    return res;
}
};

```

Show moreShow more icon

现在，这是一个很小的提要，但是，如果您想将提要限制到一个或多个特定作者，该怎么办？在这种情况下，您需要设置一个新路由（ `/author/{id}/rss` ）。新路由将获取 URL 中传递的标识符，然后使用它限制查询，仅查找给定作者编写的文章。该 RSS 方法的剩余部分基本相同。

我们看看该方法在代码中的效果。首先， `FeedController` 获取一个针对每个作者的 RSS 方法，就像之前一样：

```
var generateRss = // as before

module.exports = {

rss: // as before
authorRss: function(req, res) {
    Entry.find({ 'author' : req.param("authorID") })
      .then(function (entries) {
        var rss = generateRSS(entries);
        res.type("rss");
        res.status(200);
        res.send(rss);
      })
      .catch(function (err) {
        console.log(err);
        res.status(500)
          .json({ error: err });
      });

    return res;
}
};

```

Show moreShow more icon

不同的是上面是一个查询，它现在被限制为仅查找所有具有某个作者 ID 的 `Entry` 对象。请注意 HTTP 请求中包含的参数 `authorID` 的指令。 `authorID` （在本例中为传入的 URL 模式的第三部分）的映射在 `routes.js` 文件中指定，如下所示：

```
var generateRss = // as before

module.exports = {

rss: // as before
authorRss: function(req, res) {
    Entry.find({ 'author' : req.param("authorID") })
      .then(function (entries) {
        var rss = generateRSS(entries);
        res.type("rss");
        res.status(200);
        res.send(rss);
      })
      .catch(function (err) {
        console.log(err);
        res.status(500)
          .json({ error: err });
      });

    return res;
}
};

```

Show moreShow more icon

routes 文件中指定的参数是区分大小写的，所以请确保您保持了一致的命名约定。大小写差异是应用程序代码中一种常见但不易察觉的错误来源。

## 结束语

模型仅处理与数据相关的定义，而控制器可以处理应用程序逻辑中所有非特定于数据的部分。除非您需要特定的逻辑或约束，否则没有理由替换 Blueprint API 中使用的标准 REST 控制器。即使您需要特定的逻辑或约束，也可以经常使用现有的 Blueprint 路由作为起点。

如果您是从头开始学习本系列的，那么您现在已经很好地掌握了 Sails.js。您可以在 Bluemix 上开始使用它，定义一些模型，访问和更新您的模型，插入/更新/删除和查询它们，并创建非特定于所传递模型的端点。简而言之，您获得了在一个非平凡的 HTTP API 项目中开始使用 Sails.js 所需的所有工具。

显然，还有更多功能可供探索。在以后某个时刻，您可能希望知道如何插入您自己的 ORM/ODM（而不是使用 [第 2 部分](http://www.ibm.com/developerworks/cn/web/wa-build-deploy-web-app-sailsjs-2-bluemix/index.html) 中介绍的默认的 Waterline），或者钻研 Sails 引擎深入的配置选项。您可能还想配合使用 Sails 和 Ember.js 来实现一个强大的可行 SANE 堆栈。但就目前而言，您已经入门，也就是说，是时候最后一次说 _路途愉快_ 了！

本文翻译自： [Routes and controllers in Sails](https://developer.ibm.com/articles/wa-sailsjs4/)（2016-06-15）