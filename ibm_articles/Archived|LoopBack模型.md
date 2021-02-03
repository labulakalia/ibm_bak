# Archived | LoopBack 模型
定义和使用模型和数据源

**标签:** 

[原文链接](https://developer.ibm.com/zh/articles/wa-get-started-with-loopback-neward-2/)

Ted Neward

发布: 2017-05-31

* * *

**本文已归档**

**归档日期：:** 2019-08-30

此内容不再被更新或维护。 内容是按“原样”提供。鉴于技术的快速发展，某些内容，步骤或插图可能已经改变。

欢迎回来！在本系列的 [第 1 部分](https://developer.ibm.com/zh/articles/wa-get-started-with-loopback-neward-1/) 中，您快速了解了 [LoopBack](https://loopback.io) 服务器端 JavaScript 框架，安装了该框架，创建了一个基本应用程序，并了解了 LoopBack 的一些 API 工具。在第二期中，将了解 LoopBack 如何处理 _模型_ — 表示如何存储和检索数据的对象。LoopBack 通过多种机制来定义模型对象，每种机制都有不同的持久性结构。

在了解它们之前，我们需要确定将要建模和存储的内容。

## LoopBack 模型

为了演示如何应用模型，我们将开发一个简单 API — 但是首先，我们需要确定模型将使用哪些数据，以及如何使用它们。每个这样的示例都需要解决一个相对简单的问题 — 这个问题需要具有一定的复杂性，使我们能够了解如何建模更复杂的定义和关系，但又不能过于复杂，导致我们需要花大部分时间来解释（虚构的）应用程序模型。为此，我想对于目前的千禧代的人而言，他们可能想构建一个系统（公开为基于 HTTP 的 Web API）来跟踪超级英雄电影中的各种人物。毕竟，如果不知道有哪些超级英雄人物，他们有何能力，以及他们的弱点是什么，您如何能真正理解现代超级英雄电影的复杂细节？（请跟我一起完成这项任务 — 请记住这只是一个演示。）

我们的 Superhero-Information-as-a-Service (SIaaS) 应用程序首先会跟踪一些基本的人物信息：人物的真实姓名（姓和名）；他们的 “代号”（他们的超级英雄名称），如果他们有；以及他们的原居地 — 地球还是其他星球。最后，我们还将跟踪他们的 “因缘” 指数 — 这个抽象数字使用了简单的赞同票/反对票数据来表示他们的相对 “好” 或 “坏”，其中的赞同票表示 “好人行为”，而反对票将表示他们的 “坏人行为” 。这样，我们就能跟踪某些无法单纯地归为简单的好或坏类别的人物。

在适当的时候，我们还会记录加入超级英雄团队的人物，考虑到许多人是多个联盟的成员，这可能有点棘手 — 我们将此任务推迟到本系列后面的文章中进行解决。现在，我们将重点关注存储 “简单数据”（姓名、原居地等字符串），并建立一种机制来允许使用赞同票（英雄行为）和反对票（罪恶行为）来表示人物的宇宙因缘，但不允许用户直接修改该因缘（也就是说，我们希望封装一些数据，以便只能通过明确定义的入口点进行操作）。

> 如果您想到了一种方法，只要它不是特定于模型的，LoopBack 或许已经提供了该方法。

## 定义模型

我们首先定义一个简单的模型类，使用 LoopBack CLI 来搭建文件和核心的框架。如果您尚未这样做，请使用 npm 安装 loopback-cli 包，以便安装 lb 命令行工具。在第 1 部分中使用 hello-world 示例的地方，我们将运行 lb 来生成一个 _empty-server_ 应用程序，它将提供一个白板。

```
    lb

         _-----_
        |       |    ????????????????????????????
        |--(o)--|    ?  Let's create a LoopBack ?
       `---------   ?       application!       ?
        ( _U`_ )    ????????????????????????????
        /___A___\   /
         |  ~  |
       __'.___.'__
        `  |  Y `

    ? What's the name of your application? heroes
    ? Which version of LoopBack would you like to use? 3.x (current)
    ? What kind of application do you have in mind? empty-server (An empty LoopBack
    API, without any configured models or datasources)
    Generating .yo-rc.json

    I'm all done. Running npm install for you to install the required dependencies. If this fails, try running the command yourself.

       create .editorconfig
       create .eslintignore
       create .eslintrc
       create server/boot/root.js
       create server/middleware.development.json
       create server/middleware.json
       create server/server.js
       create .gitignore
       create client/README.md

```

Show moreShow more icon

此命令使用我们需要的核心基础信息创建了一个空应用程序。请注意，client 目录仅包含一个 README 文件；它提醒 LoopBack 仅生成一个全栈的应用程序的服务器端/API 端，没有对 API 可以做什么或应该如何使用该 API 做出任何假设。LoopBack 提供了模型和控制器；视图需要由您构建。

要验证应用程序已生成，可在应用程序代码目录中运行 `node` ，浏览到 `http://localhost:3000` 来查看 JSON 响应，其中包含服务器当前的正常运行时间（参见第 1 部分）。现在我们需要构建一个模型。我们的初始模型类型将是 **Hero** 类型，表示一个来自超级英雄宇宙的人物。但是首先，我们需要一个地方来存储这些模型和检索数据。在 LoopBack 术语中，这是一个 _数据源_ 。empty-server 应用程序默认情况下没有定义数据源。因为这是一个新建的应用程序，预先不存在任何信息，所以现在使用 `lb datasource` 命令定义一个内存型数据源是最容易的：

```
Teds-MBP15-2:code ted$ lb datasource
? Enter the data-source name: memdb
? Select the connector for memdb: In-memory db (supported by StrongLoop)
Connector-specific configuration:
? window.localStorage key to use for persistence (browser only): memdb
? Full path to file for persistence (server only): ./mem.db

```

Show moreShow more icon

需要为您的数据源定义一个唯一名称，因为 LoopBack 可能在单个应用程序中处理多个数据源，每个数据源必须有不同的名称。命名数据源后，回答一系列配置问题；对于内存型数据库，只需要一个配置元素 — 存储数据的文件的名称。

在幕后，lb 工具会创建一些文件。具体来讲，它会在文件 `server/datasources.json` 中配置数据源，这是一个简单的 JSON 文件：

```
{
"memdb": {
    "name": "memdb",
    "localStorage": "memdb",
    "file": "./mem.db",
    "connector": "memory"
}
}

```

Show moreShow more icon

**LoopBack 连接器：** LoopBack 支持各种各样的数据源；一些 [连接器](https://loopback.io/doc/en/lb2/Connecting-models-to-data-sources.html) 来自 LoopBack 工程师本身，其他连接器来自用户社区。截至撰写本文时，LoopBack 既支持大量传统的关系数据库，比如 MySQL、Oracle、PostgreSQL、SQLServer、SQLite 3 和 DB2，也支持一些非关系数据库，比如 Cloudant、MongoDB 和 Redis。社区连接器更加多样化，包括用于推送连接、REST 服务、SOAP 服务，甚至是电子邮件的连接器。

每个 JSON 元素中的值将因数据源类型不同而有所不同，但它们通常与 lb 工具要求用于命令行上的配置参数一一对应。在这种情况下，如果我突然认为 memb.db 是一个不好的数据库存储文件名称，而是更喜欢 database.json，我可以轻松地执行更改：只需更改前面的 JSON 清单中的 “file” 条目即可，如果数据已经存在，可以重命名该文件，使之与新名称相匹配。砰！数据存储被重构了。从 localStorage 切换到一个不同类型的数据源（比如 MongoDB 实例或关系数据库）需要做更多工作，我们将暂时延后该任务。

配置数据源后，就可以添加模型了。在新建的应用程序中，比如这个应用程序，lb 工具很容易直接从命令行使用 `lb model` 创建一个简单模型；但是在这么做之后，您将面临一个有趣的选择：

```
$ lb model
? Enter the model name: Hero
? Select the data-source to attach Hero to: memdb (memory)
? Select model's base class (Use arrow keys)
Model
? PersistedModel
ACL
AccessToken
Application
Change
Checkpoint
(Move up and down to reveal more choices)

```

Show moreShow more icon

前两个查询不言自明：应公开的模型名称和模型被附加到的数据源。但是模型的基类值得说明一下。

`Model` （也就是扩展 LoopBack `Model` 类型的模型）不同于 `PersistedModel` ，因为 `PersistedModel` 引入了我们期望与数据库交互的基本方法： `create` 、 `update` 、 `delete` 、 `find` 、 `findById` 等。但是，不是系统中的所有对象都需要或应该存储到数据库中，所以 LoopBack 让我们有机会扩展基础 `Model` 类型，以便在模型类型不需要持久化时保持它们的简单性。 `Model` 类型提供了有用的事件方法（比如 `changed` 和 `deleted` ），这些方法本身就很有趣，包括可用于提供对特定对象的任意访问权的 `checkAccess` 方法。

`PersistedModel` 提供了以下方法：

- `create` ：创建一个实例并将它保存到数据库
- `count` ：返回满足传入的断言条件（如果有）的对象数量
- `destroyById` ：从数据库中删除一个实例
- `destroyAll` ：删除这些模型对象的完整集合
- `find` ：查找所有满足传入的条件（如果有）的模型实例
- `findById` ：按唯一标识符查找给定模型实例
- `findOne` ：查找满足传入的条件（如果有）的第一个模型实例
- `findOrCreate` ：查找与传入的过滤器对象相匹配的一条记录；如果数据库中没有该对象，则创建该对象并返回它
- `upsert` ：更新或插入（借鉴自 MongoDB 术语），根据数据存储中是否已存在该对象来选择一个或另一个对象
- `updateAll` ：使用新数据来修改所有满足传入条件的实例（类似于经典的关系 `UPDATE` 语句）

除了上述静态方法之外， `PersistedModel` 可确保所有实例都将拥有用于单对象的类似便捷方法：

- `destroy` ：删除此对象实例
- `getId/setId` ：删除或修改此对象的唯一标识符
- `isNewRecord` ：返回对此实例是否是新实例的判断
- `reload` ：从数据库重新加载此实例，丢弃在此期间对该对象执行的所有更改
- `save` ：将对象存储到数据库

这算不上是完整的列表；请参阅 LoopBack API 文档查看完整的方法列表和每个方法的细节。如果您想到了一种方法，只要它不是特定于模型的，LoopBack 或许已经提供了该方法。由于其 I/O 性质和 NodeJS 中的相应约定，所有这些方法都会在操作完成时调用回调函数。

我们的 `Hero` 类将扩展 `PersistedModel` 类型。选择 `PersistedModel` 会引出更多问题：

```
    $ lb model
    ? Enter the model name: Hero
    ? Select the data-source to attach Hero to: memdb (memory)
    ? Select model's base class PersistedModel
    ? Expose Hero via the REST API? Yes
    ? Custom plural form (used to build REST URL): Heroes
    ? Common model or server only? server
    Let's add some Hero properties now.

    Enter an empty property name when done.
    ? Property name: Codename
       invoke   loopback:property
    ? Property type: string
    ? Required? Yes
    ? Default value[leave blank for none]:

    Let's add another Hero property.
    Enter an empty property name when done.
    ? Property name:

```

Show moreShow more icon

许多 Web API 开发人员希望通过类似 REST 的 API 端点公开所定义的模型；LoopBack 也考虑到了这一点，所以为了方便起见，对 “Expose via REST API” 问题选择 **Yes** 会自动定义一些端点和 URL 参数。因为机器并不总是能推断出模型名称的复数形式，所以 LoopBack 要求我们采用复数形式，以便端点具有正确的语法。

然后，因为模型定义通常由前端和后端共享，所以 LoopBack 会询问是否应通过在客户端和服务器之间的一个共同目录共享来导出模型的定义。因为该应用程序只有服务器端，所以选择 **server** 会保持代码的独立性。

最后，LoopBack 会要求您定义模型类型的属性。第一个属性是英雄的代号，这是一个字符串，该属性必须存在且没有默认值。LoopBack 将继续询问要定义的属性，直到您为一个属性输入空白名称。

## 模型定义和文件结构

完成该操作后，将会在 server/models 目录中定义两个文件：hero.js 和 hero.json。JSON 文件是 `Hero` 类型的定义，反映出了在命令行做出的选择：

```
{
"name": "Hero",
"plural": "Heroes",
"base": "PersistedModel",
"idInjection": true,
"options": {
    "validateUpsert": true
},
"properties": {
    "Codename": {
      "type": "string",
      "required": true
    },
    "FirstName": {
      "type": "string"
    },
    "LastName": {
      "type": "string"
    },
    "Origin": {
      "type": "string"
    },
    "Karma": {
      "type": "number",
      "required": true,
      "default": 0
    }
},
"validations": [],
"relations": {},
"acls": [],
"methods": {}
}

```

Show moreShow more icon

## 定义其他方法

hero.js 文件提供了向 `Hero` 类型添加其他方法的机会，可以获取传入的参数（ `Hero` 原型对象）并添加想要的方法，比如添加 `upvote` 和 `downvote` 方法来适当调整英雄的因缘。（我们稍后会查看这些方法；现在我们只填充数据字段/属性。）

此模型定义文件中提供的全部选项已在 LoopBack [文档](https://loopback.io/doc/en/lb3/Model-definition-JSON-file.html) 中介绍，但可以在前面的清单中看到一些特性：

**模型关系：** LoopBack 支持大多数数据库系统中通用的大部分（甚至所有）模型间关系（一对一、一对多等）。在第 3 部分中，当我们进一步细化模型时，会更详细地分析如何使用和定义 LoopBack 的模型关系。

- `base` ：此模型的基础类型。我们已讨论了 `Model` 和 `PersistedModel` ，但 LoopBack 模型也可以扩展一些预定义的模型类型，比如 `User` 、 `Email` 等。
- `idInjection` ：表明 LoopBack 是否将负责管理模型中的 **id** 字段，假设该字段是模型的主键。对于大多数模型，该值默认为 **true** 。
- `relations` ：定义此模型与其他模型具有的关系；本系列第 3 部分将更详细地讨论各种关系。
- `validations` ：模型的这部分用于定义对属性的验证。LoopBack 的当前文档明确表明尚未实现此选项。大概在未来的版本中，我们将能设置 LoopBack 将自动执行的最小长度、最大长度和基于正则表达式的验证。就现在而言，所有验证都需要手写代码。
- `acls` ：不同于 NodeJS 世界的大部分类似产品，LoopBack 拥有功能丰富且强大的访问控制模型，允许您定义角色、权限和访问控制列表来管理用户的对象访问权。我们将在未来的文章中探讨 LoopBack 的访问控制。

模型定义文件提供了 API 的基本结构，但 LoopBack 的开发人员认识到，不是所有属性都能以类似 JSON 的数据格式进行捕获；有时需要编写代码。因此，hero.js 文件允许开发人员获取传入的 `Hero` 对象，添加看起来适合模型原型对象的属性。

例如，我们提到过，英雄的 “好” 是由他们的因缘来衡量的：他们每做一件好事，他们的因缘就会升高，每做一件坏事，因缘就会降低。通常，要影响英雄的属性，可以发出 `PUT` 请求，在请求的正文中传入新数据。并不总是这样；我们更喜欢只通过赞同票和反对票来影响因缘。在 LoopBack 中，要实现此目的，可以在 `Hero` 上创建两个方法，然后将它们注册为远程方法，表明它们包含在整个 Hero web API 中。

您将在（空）hero.js 文件中添加这些方法。首先，在 `Hero` 上定义两个方法，然后将每个方法注册为远程方法：

```
    module.exports = function(Hero) {
        Hero.prototype.upvote = function(cb) {
            var response = "Yay! Hero did a good thing";
            this.Karma += 1;
            this.save(function(err, hero) {
                if (err) throw err;

                cb(null, hero);
            });
        };
        Hero.remoteMethod('prototype.upvote', {
            http: { path: '/upvote', verb: 'post' },
            returns: { arg: 'result', type: 'object' }
        });

        Hero.prototype.downvote = function(cb) {
            var response = "Boo! Hero did a bad thing";
            this.Karma -= 1;
            this.save(function(err, hero) {
                if (err) throw err;

                cb(null, hero);
            });
        };
        Hero.remoteMethod('prototype.downvote', {
            http: { path: '/downvote', verb: 'post' },
            returns: { arg: 'result', type: 'object' }
        });
    };

```

Show moreShow more icon

请注意，这些方法是在 `Hero` 的原型上定义的；要将它们定义为实例方法（也就是说，在每个 `Hero` 实例上定义它们），就必须这么做。如果没有该定义，LoopBack 会假设它们是静态方法，与任何特定的 `Hero` 实例都没有关联（因此无法确定哪个 `Hero` 的因缘被投赞同票或反对票）。适当修改 `Karma` 后，我们依靠 `Hero` 从 `PersistedModel` 继承的内置 `save` 方法，将 `Hero` 的新数据写回数据源。

注册远程方法相对比较简单，仅需调用一个 `remoteMethod()` 方法，该方法将获取要注册的方法的名称和用于公开该方法的 URL 的一些数据：HTTP 路径和 verb、参数的一些描述（在本例中没有）和返回值。

这看起来比较麻烦，但最终结果是 LoopBack 现在拥有您希望公开的端点的直接信息，而且有了该信息，它就可以将这两个端点添加到在 LoopBack GUI 的 Explorer 视图中公开的 Swagger UI 中。

## 模型验证

通过在 hero.js 中添加一条验证约束，也可以添加验证代码来确认没有两个英雄选择相同的代号名称：

```
    module.exports = function(Hero) {

        // ... as above

        Hero.validatesUniquenessOf('Codename');
    };

```

Show moreShow more icon

此代码使用了 `validatesUniquenessOf` 方法，该方法是 `PersistedModel` 通过 [Validatable](http://apidocs.strongloop.com/loopback-datasource-juggler/#validatable) 混合类获得的，该混合类包含以下方法：

- `validatesAbsenceOf` ：验证缺少一个或多个指定的属性；也就是说，要被视为有效，模型不应包含某个属性，而且在验证的字段为非空值时，验证会失败。
- `validatesExclusionOf` ：验证该属性的值不在某个值范围内。
- `validatesFormatOf` ：需要一个属性值来匹配某个正则表达式。
- `validatesInclusionOf` ：要求一个属性值在某个值范围内。（这是 `validatesExclusionOf` 的逻辑相反值。）
- `validatesLengthOf` ：要求一个属性值的长度在指定范围内。
- `validatesNumericalityOf` ：验证某个属性是数字。
- `validatesPresenceOf` ：要求某个模型拥有给定属性的值。（这是 `ValidateAbsenceOf` 的逻辑相反值。）
- `validatesUniquenessOf` ：验证某个属性值在所有模型中是唯一的。请注意，不是所有数据存储都支持此方法；截至撰写本文时，只有内存型、Oracle 和 MongoDB 数据源连接器支持此方法。
- `validate` ：允许附加一个自定义验证函数。

此外， `Validatable` 向每个实例添加了一个 `isValid()` 方法，所以您可以随时验证某个对象是否有效（例如，在将它作为结果发回之前），而无需明确尝试存储该对象。

说得明确点， `Validatable` 及其方法可从任何继承自基本 `Model` 类型的模型进行访问 — 对象无需是 `PersistedModel` 也可获得对象的好处。

## 发现模型

对于已有的数据库无法从头开始重新构建的开发人员，LoopBack 没有忘记你们。LoopBack 能利用现有的数据源来发现模型，比如关系数据库表和模式（目前支持 MySQL、PostgreSQL、Oracle 和 SQL Server）。这是一个生成模型定义文件的一次性流程，LoopBack 随后会使用这些文件，就像使用其他任何模型定义一样。

LoopBack 文档建议在一个独立 NodeJS 脚本中完成这项任务，比如：

```
    var loopback = require('loopback');
    var ds = loopback.createDataSource('oracle', {
      "host": "oracle-demo.strongloop.com",
      "port": 1521,
      "database": "XE",
      "username": "demo",
      "password": "L00pBack"
    });

    // Discover and build models from INVENTORY table
    ds.discoverAndBuildModels('INVENTORY', {visited: {}, associations: true},
    function (err, models) {
      // Now we have a list of models keyed by the model name
      // Find the first record from the inventory
      models.Inventory.findOne({}, function (err, inv) {
        if(err) {
          console.error(err);
          return;
        }
        console.log("\nInventory: ", inv);
        // Navigate to the product model
        // Assumes inventory table has a foreign key relationship to product table
        inv.product(function (err, prod) {
          console.log("\nProduct: ", prod);
          console.log("\n ------------- ");
        });
      });
    });

```

Show moreShow more icon

为了能在运行时使用输出，必须将输出写入到一个文件中（通常为 common/models/model-name.json），然后在 server/model-config.json 文件中手动注册。

对于不是关系数据库的数据源，LoopBack 也可以从非结构化 (JSON) 数据实例（也即 MongoDB 数据库、REST 数据源或 SOAP 数据源）推断出模型。我们完成此任务的方式类似于关系数据库：获取一个用于推断的样本实例，将它传递给数据源的 `buildModelFromInstance()` 方法。返回值是生成的模型类型，可通过与之前的 `Hero` 对象相同的方式使用它。LoopBack 使用一个原始 JSON 对象演示了一个相关示例：

```
    module.exports = function(app) {
      var db = app.dataSources.db;

      // Instance JSON document
      var user = {
        name: 'Joe',
        age: 30,
        birthday: new Date(),
        vip: true,
        address: {
          street: '1 Main St',
          city: 'San Jose',
          state: 'CA',
          zipcode: '95131',
          country: 'US'
        },
        friends: ['John', 'Mary'],
        emails: [
          {label: 'work', id: 'x@sample.com'},
          {label: 'home', id: 'x@home.com'}
        ],
        tags: []
      };

      // Create a model from the user instance
      var User = db.buildModelFromInstance('User', user, {idInjection: true});

      // Use the model for create, retrieve, update, and delete
      var obj = new User(user);

      console.log(obj.toObject());

      User.create(user, function (err, u1) {
        console.log('Created: ', u1.toObject());
        User.findById(u1.id, function (err, u2) {
          console.log('Found: ', u2.toObject());
        });
      });
    };

```

Show moreShow more icon

对于非结构化数据样本，在服务器的正常启动序列中以引导脚本形式运行此代码更容易一些。但是，除非数据的结构不断变化（或许甚至在开发人员不知情的情况下），与服务器每次启动时使用资源来重新解析数据实例相比，一次性将此代码转换为模型定义文件更为明智。

## Model API 和 LoopBack GUI

如果再次启动服务器（在 root 目录中运行 `node` ）并浏览到 `http://localhost:3000/explorer`，您将会看到已连接了 `Hero` 类型。可以使用交互式 GUI 来创建新英雄 — 例如，使用 `POST /Heroes` 端点将一个英雄添加到系统，然后使用 `GET /Heroes` 查看已添加的所有英雄的完整列表。

尽管它不会取代使用该 API 的单元测试，但 LoopBack GUI 使得在开发期间交互式地测试该 API 变得更容易。

## 模型引导

任何系统通常都有一组需要引导才会存在的对象。所有超级英雄宇宙中都有一些共同的英雄，所以为了简便起见，我们希望系统在启动时预加载这些英雄。LoopBack 为这种类型的引导代码提供了一个挂钩（hook） — boot 子目录中包含的所有代码都将在服务器启动时执行，所以您只需在这里添加一个文件，用于导出一个要在启动期间调用的函数即可。此函数接收一个参数 — 应用程序对象 — 我们可使用该对象获取模型原型对象，并在新对象不存在时创建它们：

```
    module.exports = function(app) {
        app.models.Hero.count(function(err, count) {
            if (err) throw err;

            console.log("Found",count,"Heroes");

            if (count < 1) {
                app.models.Hero.create([{
                    Codename: "Superman",
                    FirstName: "Clark",
                    LastName: "Kent",
                    Karma: 100
                }], function(err, heroes) {
                    if (err) throw err;

                    console.log('Models created:', heroes);
                });
            }
        });
    };

```

Show moreShow more icon

定义的所有模型类型都将是 `app.models` 对象的成员，所以在代码库中引用应用程序对象的任何地方，所有模型都是一个简单的属性解引用。

## 结束语

无论是我们在本系列中分析 LoopBack 模型时，还是在您自行钻研 LoopBack 时，您都会发现这些模型的更多细节，但您找到的大部分信息可能会与我们这里介绍的主题有所不同；您目前掌握的知识已足以帮助您开始使用 LoopBack 模型。但是，第 3 部分会有一个更大的主题等着我们：如何建立模型之间的关系。LoopBack 对关系提供了一定的结构化支持，但建模关系 — 特别是跨不同类型的数据存储系统 — 并不总是简单和直观的。尽管 LoopBack 尝试抽象化底层数据，但它无法处理所有情况下的所有场景。

但现在是时候说再见了。祝 LoopBack 学习愉快！

本文翻译自： [LoopBack Models](https://developer.ibm.com/articles/wa-get-started-with-loopback-neward-2/)（2017-04-24）

Ted Neward