# 在 Sails 中建模关系
使用关联来连接模型对象

**标签:** Node.js,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/wa-sailsjs3/)

Ted Neward

发布: 2016-10-26

* * *

**关于本系列**

像同类产品 Ruby on Rails 一样，Sails.js 是一个旨在帮助开发人员构建 Web 应用程序的框架。Rails 默认用于构建在服务器上生成 HTML 并将其发回给浏览器的应用程序，与它不同的是，Sails 为构建 HTTP API 提供了很好的开箱即用支持。这种设计有助于在客户端和服务器之间实现更好的关注点分离，并促进客户端和服务器之间的更多的互操作性。在本系列中，Web 开发人员和教师 Ted Neward 将向您介绍 Sails.js。HTTP API 示例应用程序可以是 AngularJS、React 或几乎其他任何客户端应用程序的后端。

_好久不见_ ，Sails 粉丝们！（你们中的一些人是十足的代码海盗，因此我要为大家准备了一个热烈的海盗欢迎仪式。）在本系列的 [上一篇教程](http://www.ibm.com/developerworks/cn/web/wa-build-deploy-web-app-sailsjs-2-bluemix/index.html) 中，您为 HTTP API 示例应用程序构建了一些基本模型，然后让 Sails 蓝图自动配置该应用程序的 CRUD 操作。但是，在实际的应用中，您需要的不仅仅是一些管理所有操作的简单模型，尤其是因为大部分数据都与其他数据相连。您构建的系统包含数十种数据关系，所有这些关系都需要在您的模型中反映出来。

在这一期的教程中，将学习如何定义反映数据之间的关系的模型。具体地讲，您将：

- 重构 `BlogEntry` 模型，以便它不仅能管理博客文章（blog entry），还能管理更广泛的内容类型。
- 探索使用 Sails 中两种不同类型的数据模型建模作者和作者身份的可能性。
- 向您急速发展的内容管理系统中的文章添加评论和标签。

开始之前，让我快速回顾一下目前为止介绍的内容。

## 在 Sails 中建模

您开发的 [示例应用程序](http://public.dhe.ibm.com/software/dw/web/wa-sailsjs3-source.zip) 是一个 HTTP API。目前为止，我们将它构想为一个博客平台，该平台是一种内容管理系统（CMS）。HTTP API 定义了 CMS 后端，在本例中，这个后端为一个未定义 UI 的博客提供支持。您构建了一个 RESTful API 来存储博客文章、评论和 RSS 提要。该 API 支持并响应查询（例如排序和查找特定的博客或文章）。

在 [上一篇教程](http://www.ibm.com/developerworks/cn/web/wa-build-deploy-web-app-sailsjs-2-bluemix/index.html) 中，您定义了第一个模型 `BlogEntry` ，它位于项目目录中的 `api/models/BlogEntry.js` 文件中。清单 1 给出了 `BlogEntry` 的代码。

##### 清单 1\. BlogEntry.js

```
module.exports = {

attributes: {
    title: {
      type: 'string',
      required: true,
      defaultsTo: ''
    },
    body: {
      type: 'string',
      required: true,
      defaultsTo: ''
    }
}

};

```

Show moreShow more icon

## 重构 Sails 模型

如果您现在想要将博客 API 扩展为一个更加通用的 CMS，该怎么办？现在您还处于开发流程中足够早的阶段，这种情形尚处于可管理状态，而且最终得到的结果会比一个博客平台灵活得多。它不会损害在 Sails.js 中重构相对容易的事实。

您首先要将模型的名称从 _BlogEntry_ 更改为某个用途更广的名称。（顺便说一下，我听说编程过程中最困难的三件事包括命名和差一 (off-by-one) 错误。 _我们拭目以待！_ ）因为该模型的类型是从其文件名中获得的，所以将它从 `BlogEntry.js` 重命名为 `Entry.js` 会告诉 Sails 更新模型类型。对相应的 `api/controllers/BlogEntryController.js` 文件执行相同操作，将它更改为 `EntryController.js` ，就这么简单：您已重构了您的模型。您的 HTTP API 现在能够用于比网络博客更多的用途。

但是回想一下，您使用了 `sails-disk` 作为开发数据库适配器。 `sails-disk` 是一种直接存储到磁盘的序列化格式；所以它没有表、列或其他任何类似数据库的基础架构。这种简单性使 `sails-disk` 在开发期间很容易使用，但您需要在代码接近生产阶段时将它替换为其他格式。您可能想知道在应用程序准备好上线时，这种看似容易的重构将如何进行。

幸运的是，Sails 中的每个模型对象可保留许多模型属性。设置模型属性，使 Sails 能够将模型与底层数据库相匹配。您可以从 Sails 文档中了解模型设置。就目前而言，您只需要关心 `tableName` 。如果您对某个真实的数据库使用了此属性，结果将类似于：

##### 清单 2\. Entry.js

```
module.exports = {

tableName: "blogentry",
    // this would map to a relational table by this name,
    // or a MongoDB collection, and so on

attributes: {
    id: {
      type: 'integer',
      primaryKey: true
    },
    title: {
      type: 'string',
      required: true,
      defaultsTo: ''
    },
    body: {
      type: 'string',
      required: true,
      defaultsTo: ''
    }
}

};

```

Show moreShow more icon

在这里，可以看到 `Entry.js` 指定了它所绑定的表名称，在本例中为”blogentry”。如果模型对象中的特定字段需要与底层数据库中的指定列对应，您可以使用 `columnName` 属性来注释每个字段，命名它应该映射到的表列（或集合中的字段，具体取决于数据存储类型）。清单 3 给出了一个例子。

##### 清单 3\. 将字段映射到列

```
module.exports = {

tableName: "blogentry",
    // this would map to a relational table by this name,
    // or a MongoDB collection, and so on

attributes: {
    id: {
      type: 'integer',
      primaryKey: true,
      columnName: 'blogentry_pk'
    },
    title: {
      type: 'string',
      required: true,
      defaultsTo: '',
      columnName: 'blogtitle'
    },
    body: {
      type: 'string',
      required: true,
      defaultsTo: ''
    }
}

};

```

Show moreShow more icon

在向系统添加更多条目类型时，需要执行一些额外的更改，但就现在而言，这就足够了。

## Sails 中的关联

大多数发表的内容类型都会存储和显示一位或多位作者，所以您需要一个相关的模型。清单 4 展示了第一次尝试。

##### 清单 4\. 在 CMS 中重新表示作者

```
module.exports = {
autoPK: false,
attributes: {
    fullName: {
      type: 'string',
      required: true
    },
    bio: {
      type: 'string'
    },
    username: {
      type: 'string',
      unique: true,
      required: true
    },
    email: {
      type: 'email',
      required: true
    }
}
};

```

Show moreShow more icon

目前而言， `Author.js` 是一种简单易懂的数据类型，而且它能够很好地表示作者的属性：全名、简介、用户名等。该模型中缺少的是 _作者身份_ （authorship）的概念：一位作者创建了一篇文章，因此每篇文章都由一位作者编写。这比您目前处理的概念更难建模。事实上，这时就需要使用 Sails _关联_ ，此概念不同于简单的属性。

### 评论和标签

作者身份不是您唯一需要为此应用程序建模的关联，所以在解决这个大问题之前，让我们看看两个更简单的模型。每篇文章都拥有评论和一组可用于描述它的标签。添加标签会得到一个用于直观显示的”标签云”（tag cloud），生成一种基于元数据的主题分类系统。通过建模这些类型，您可以练习使用关联。只需记住，实际的 CMS 需要十几种关联。

回想一下，我们遵循的一条编程规则（包括使用 Sails 编程）就是 _保持简单_ （Keep it simple）。按照这种编码精神，评论的数据模型基本上应该仅包含评论的正文，以及可选的发表评论的人的电子邮件地址，如清单 5 所示。

##### 清单 5\. Comment.js

```
module.exports = {

attributes: {
    body: {
      type: 'string',
      required: true
    },
    commenterEmail: {
      type: 'email'
    }
}
};

```

Show moreShow more icon

类似地，标签的数据模型仅包含标签的名称（通常是内容的元数据，比如”Java”或”Sails”），不需要额外的修饰。清单 6 给出了在您的 CMS 中实现内容标签的数据模型。

##### 清单 6\. Tag.js

```
module.exports = {

attributes: {
    name: {
      type: 'string'
    }
}
};

```

Show moreShow more icon

在定义了两个数据模型后，您就掌握了定义更多模型所需的基础知识。想要添加一个模型时，只需在 `api/model` 目录中创建一个模型。也可以输入命令： `sails generate model ...` ，Sails 就会为您添加一个占位符。

现在是时候开始定义您的数据模型与它们包含的数据之间的关系了。

## 显式关系

Sails 没有采用一些数据库系统所使用的隐式方法，它使用了显式的关系模型。例如，对于关系数据库系统，Sails 会使用 _数据_ 来建模两个表之间的关联——在一个表中定义的主键，它的值用作另一个表中某一行的外键值——而不是在数据库模式（schema）中结构化地定义它。

关系数据库追随者们会注意到，大部分数据库系统都支持结构化定义的关系。我的意思是说，在 RDBMS 中，您可以使用数据库约束来确保任何作为外键值插入的值也存在于相关的表中。对于 Sails，我们使用数据（而不是某种物理结构）来表示这种关系。与关系模型相反，可以考虑一种面向文档的数据库，比如 MongoDB 或 CouchDB。在面向文档的系统中，您嵌入了一个数组作为文档的成员，而不是将一组值与其他数据元素关联。

在针对一种特定数据结构而建模时，隐式建模非常适合；在您的数据要使用多种数据库类型来组织时，它就不太适合了。Sails 需要显式理解关系，以便知道如何针对给定数据库类型来建模和搭建语句或查询——无论是 RDBMS、NoSQL 或其他某种类型。尽管这可能对您的模型对象提出一些不熟悉的需求，但这些要求不是太严格；您只需要学会更加显式地考虑数据及其连接方式即可。

## 一对多关系

首先，考虑文章与作者的关系。一位作者可以编写多篇文章，而每篇文章只能有一个作者。不出所料，Sails 将此称为 _一对多关系_ （one-to-many relationship）。作者与文章之间的关系也是 _双向的_ （bidirectional），因为应该可以检索给定作者的所有文章，以及查看任何给定文章的作者。（事实证明，Sails 默认情况下将在查询中自动拉取这部分附加数据，并将它发送到客户端。）

定义一对多关系需要修改该关联两端的模型对象。您需要定义作为关联的 “一” 端上的集合的字段，以及将关联的 “多” 端连接到这个 “一” 端的字段。这不太适合用文字描述，但在代码中看到会简单得多。清单 7 给出了修改后的 `Author` 模型。

##### 清单 7\. Author.js 和关联的文章

```
module.exports = {
attributes: {
    fullName: {
      type: 'string',
      required: true
    },
    bio: {
      type: 'string'
    },
    username: {
      type: 'string',
      unique: true,
      required: true
    },
    email: {
      type: 'email',
      required: true
    },
    entries: {
      collection: 'entry',
      via: 'author'
    }
}
};

```

Show moreShow more icon

清单 7 中重构的代码表明， `Author` 拥有一个 entries 字段，该字段包含由若干篇文章形成的一个 `Entry` 对象集合。 `Entry` 类型通过 `Entry` 对象上的”author”字段指向 `Author` 实例。所有这些意味着 `Entry` 类型需要看起来类似于清单 8。

##### 清单 8\. 文章类型

```
module.exports = {

attributes: {
    title: {
      type: 'string',
      required: true,
      defaultsTo: ''
    },
    body: {
      type: 'string',
      required: true,
      defaultsTo: ''
    },
    comments: {
      collection: 'comment',
      via: 'owner'
    },
    author: {
      model: 'author'
    }
}

};

```

Show moreShow more icon

请注意，在 [清单 7](#清单-7-author-js-和关联的文章) 和 [清单 8](#清单-8-文章类型) 中，引用的类型（在 `Author` 的 entries 字段的 _collection_ 字段中，以及 `Entry` 的 author 字段的 _model_ 字段中）使用了小写。

这是因为 Sails 从小写形式的文件名（也称为类型的身份）获取模型类型。小写的类型（前面的清单中的 `entries` 和 `author` ）变成了所生成的蓝图路由的前缀。小写的类型也变成了 Sails 系统中模型的正式名称。

当我在本系列的下一篇教程中讨论控制器时，您还会看到身份的概念。Sails 需要能够在控制器级别确定控制器和模型是否具有相同的身份。它使用该信息生成正确的默认蓝图路由。就目前而言，只需注意 Sails 要求用作 model 字段值的类型应为小写。

## 链接模型对象

我们返回到 `Author` 与 `Entry` 之间的一对多关联上。还需要了解如何在您的代码中使用该关联，包括不仅需要定义它，还需要知道在从数据库检索一个模型对象时期望获得哪些信息。对我们而言，幸运的是，Sails 能非常灵活地链接模型对象。

当您创建一个 `Author` 实例时，Sails 为它生成了一个唯一主键，该主键是在 _id_ 字段中定义的。您可以使用新的 _id_ 作为关联字段的值，Sails 会自动连接两个对象，如清单 9 所示。

##### 清单 9\. 连接对象

```
Author.create({
    fullName: "Fred Flintstone",
    bio: "Lives in Bedrock, blogs in cyberspace",
    username: "fredf",
    email: "fred@flintstone.com"
}).exec(function (err, author) {
    Entry.create({
        title: "Hello",
        body: "Yabba dabba doo!",
        author: author.id
    }).exec(function (err, created) {
        Entry.create({
            title: "Quit",
            body: "Mr Slate is a jerk",
            author: author
        }).exec(function (err, created) {
            return res.send("Database seeded");
        });
    });
});

```

Show moreShow more icon

[清单 9](#清单-9-连接对象) 中的代码是经典的 Node.js，被称为”callbacks galore”。第一个调用使用您传入的值创建了一个 `Author` 实例。在触发 `exec()` 中的回调时，您将获取 `Author` 的 ID 值并设置为您新创建的 `Entry` 对象中的 _author_ 字段的值。或者更简单地讲：通过设置 _author_ 字段来引用正确的 `Author` ，将 `Entry` 链接到 `Author` 。

作为前面提及的”callbacks galore”的 **替代方案**，Sails 支持蓝鸟承诺。我想避免针对不同承诺风格的库的利弊的争议，所以我暂时对代码采用经典回调格式。

### 链接模型对象的另一种方式

如果前一种方法不适合您，Sails 还提供了一种替代方案：无需获取 `Author` 的 _id_ 字段，您可以传递整个 `Author` 对象。无论采用哪种方式，最终结果都是一样的。

对于更习惯于考虑物理存储模型的开发人员，直接传入对象可能更有吸引力，而传入 _id_ 能更准确地反应对象之间的链接。如果您习惯于”对象思维”，传递对象可以确认对象现在已链接，但抽象化了它们的链接方式的细节。

## 揭示关联

无论您如何达到这一步，一旦在数据库中设置对象后，Sails 就会不遗余力地明确显示它们之间的关联。清单 10 展示了在您运行 [清单 9](#清单-9-连接对象) 中的代码，然后访问 **[http://localhost:1337/author](http://localhost:1337/author)** （从系统获取所有作者的蓝图默认路由）时，您会看到的结果。

##### 清单 10\. 返回的作者查询结果

```
[
{
    "entries": [
      {
        "title": "Hello",
        "body": "Yabba dabba doo!",
        "author": 6,
        "createdAt": "2016-02-16T21:15:55.722Z",
        "updatedAt": "2016-02-16T21:15:55.722Z",
        "id": 6
      },
      {
        "title": "Quit",
        "body": "Mr Slate is a jerk",
        "author": 6,
        "createdAt": "2016-02-16T21:15:55.725Z",
        "updatedAt": "2016-02-16T21:15:55.725Z",
        "id": 7
      }
    ],
    "fullName": "Fred Flintstone",
    "bio": "Lives in Bedrock, blogs in cyberspace",
    "username": "fredf",
    "email": "fred@flintstone.com",
    "createdAt": "2016-02-16T21:15:55.716Z",
    "updatedAt": "2016-02-16T21:15:55.716Z",
    "id": 6
}
]

```

Show moreShow more icon

类似地，访问 `Entry` 对象相应路由的过程类似于清单 11。

##### 清单 11\. 返回的文章查询结果

```
[
{
    "comments": [],
    "author": {
      "fullName": "Fred Flintstone",
      "bio": "Lives in Bedrock, blogs in cyberspace",
      "username": "fredf",
      "email": "fred@flintstone.com",
      "createdAt": "2016-02-16T21:15:55.716Z",
      "updatedAt": "2016-02-16T21:15:55.716Z",
      "id": 6
    },
    "title": "Hello",
    "body": "Yabba dabba doo!",
    "createdAt": "2016-02-16T21:15:55.722Z",
    "updatedAt": "2016-02-16T21:15:55.722Z",
    "id": 6
},
{
    "comments": [],
    "author": {
      "fullName": "Fred Flintstone",
      "bio": "Lives in Bedrock, blogs in cyberspace",
      "username": "fredf",
      "email": "fred@flintstone.com",
      "createdAt": "2016-02-16T21:15:55.716Z",
      "updatedAt": "2016-02-16T21:15:55.716Z",
      "id": 6
    },
    "title": "Quit",
    "body": "Mr Slate is a jerk",
    "createdAt": "2016-02-16T21:15:55.725Z",
    "updatedAt": "2016-02-16T21:15:55.725Z",
    "id": 7
}
]

```

Show moreShow more icon

尽管文章和作者是单独存储的，但 Sails 会利用它对这些模型对象关联的了解，”填充”合适的字段并让它们看起来像是一个平面对象。因为 Sails 的目的是用作一个后端 HTTP API 实现库，所以会将从应用程序 UI（移动或 Web）到数据库的往返次数保持到最少。通过”扁平化”数据结构，您可以在一次（可能很漫长的）往返中了解 `Author` 的完整细节。这有助于实现能更有效地执行和扩展的更高效系统。

我们仅探索了一种关联模型（一对多模型），但 Sails 支持所有模型：一对一、多对多，以及这些主题上的一些不太传统的变体。它们在很大程度上具有类似的工作方式：在模型对象上定义合适的字段，将对象或其 ID 分配给关联字段，让 Sails 负责处理剩余工作。

## 结束语

现在，您应该已经很好地掌握了如何在 Sails 中建模实体。在本系列教程的下一期中，我们会通过 Sails 控制器，向系统中引入更多”活动”。默认的蓝图路由非常不错，它们适合基本的 CRUD 式访问，但大部分系统都需要定义更详尽地反映对象模型的自定义路由。我们下一次将深入介绍所有这些主题。现在是时候（再次）说 _旅途愉快_ 了！

本文翻译自： [Modeling relationships in Sails](https://developer.ibm.com/articles/wa-sailsjs3/)（2016-04-04）