# 优化 Cloudant 中的查询
了解如何为应用程序获得最高效的数据查询服务

**标签:** 分析

[原文链接](https://developer.ibm.com/zh/articles/ba-optimize-queries-cloudant/)

Chuan Yang Wang, Tao Liu

发布: 2018-07-11

* * *

## 简介

应用程序需要比以往更大的数据灵活性。近年来各种 NoSQL 数据库的出现弥补了传统关系数据库的不足。它们更灵活的数据模型能够更好地支持非结构化和半结构化数据的应用需求。作为一种 NoSQL DB 产品，IBM® Cloudant 为 Web 和移动应用程序提供了完全托管的数据库服务，还提供了丰富的功能，如高级索引技术、自定义视图、全文搜索和实时查询。本文将从多个角度提供一些方法和经验，供您在运行 Cloudant NoSQL DB 时用来优化查询。这有助于您更深入地了解每个适用场景的最合适查询，从而为您的应用程序提供最高效的数据查询服务。

## Cloudant 简介

Cloudant 是 Apache CouchDB 的商业版本，为 Web 和移动应用程序提供了完整的 IBM 操作和维护数据管理平台。作为构建于云上的 NoSQL 数据库，Cloudant 非常适合 Web 和移动应用程序的快速增长。它使用户能够利用云的可用性、灾备能力和广泛覆盖特征，将应用程序扩展到更高的级别，并通过优化实现更好的数据可用性、耐用性和移动性。Cloudant 拥有强大的索引功能，而且能够跨多个数据中心和设备将数据推送到网络边缘，这些加快了它的访问速度并提高了它的容错能力。它使得用户能够随时随地访问数据。

与其他数据库产品相比，Cloudant 具有以下特性：

- 它支持副本分发和复制。可为云中的每个节点授予数据库读写访问权，这使得 Cloudant 能够在任何规模的群集上运行，并确保数据的安全性和一致性。
- 它能在副本之间执行连续的数据同步。Cloudant 支持在多个副本之间进行同步，还支持实时自动同步。一个副本中的连续数据更新可以自动同步到其他相关副本。
- 它提供了各种方法来优化您的查询，包括索引和视图。在视图支持方面，Cloudant 支持用户定义的视图。它依赖的语言是 JavaScript。
- 它为 JSON 文档提供了数据分析和仓库功能。Cloudant 支持集成 IBM Db2 Warehouse 来实现传统数据仓库和商业智能分析。Cloudant 还支持集成 IBM Db2 Warehouse 来实现在线报告和商业智能分析。
- 它支持采用 GeoJSON 格式存储基于位置的数据。

## 通过索引进行 Cloudant 查询优化

在数据库中，可以通过为经常使用的数据和相关查询建立索引来提高查询速度。Cloudant 支持两种类型的索引：

- “type”: “text”
- “type”: “json”

这两种索引在用途和使用方法上存在显著差异。

从用途角度讲，文本类型的索引关注的是 Cloudant 文档的特定内容，而不是文档本身的结构。因此，如果用户不熟悉数据库中文档的结构，或者文档结构差别较大且格式复杂，那么应该首选文本索引。相较而言，JSON 索引对文档的结构有很高的要求，而且它建立在一个或多个特定字段上。因此，如果用户熟悉数据库中的文档结构，那么可以选择创建 JSON 索引。通过显式指定合适的字段，可以优化所有包含 Cloudant 上的特定字段的查询。从这个角度来看，JSON 索引的概念类似于传统关系数据库中的索引，二者都针对特定的列或字段进行了优化。

从使用方法的角度来看，文本索引只能在 Cloudant 数据搜索接口中使用，其中支持的语法是 [Apache Lucene Query Parse Syntax](https://lucene.apache.org/core/2_9_4/queryparsersyntax.html) 。相较而言，JSON 索引只能用于 Cloudant 的数据查询接口，允许用户根据 Cloudant 查询分析语法使用 JSON 对象进行查询。两种索引类型的接口都很强大，都支持各种自定义查询。甚至对于文本类型的索引，Cloudant 还允许用户精确定义要建立索引的数据范围，例如，通过条件判断来识别特定字段。因此，在大部分场景中（当文档结构相对简单和清晰时），这两种索引可以相互转换。但是，在文档结构复杂且不清楚时，只能使用文本索引和数据搜索接口。

在索引速度方面，当处理相同数量的数据时，文本类型的索引比 JSON 类型的索引慢。原因在于，创建文本索引时，Cloudant 不仅要处理指定的数据结构，还要处理其内容。相较而言，JSON 索引只关心结构本身。在某些情况下，文本索引可用于整个数据库的全文搜索。代码非常简单，如下所示：

#### 示例 1.使用文本索引创建全文搜索

```
{
    "type": "text",
    "index": { }
}

```

Show moreShow more icon

此外，Cloudant 的一些高级功能（比如复杂的聚合和基于地理空间的计算）只能用于基于 JSON 索引的查询接口。因此，选择使用两种索引类型中的哪种类型的关键在于对现有数据的理解。

### 基于 Cloudant 视图的查询优化

在传统关系数据库中，视图用于过滤数据、控制访问和优化查询。在 Cloudant 中，视图主要用于过滤数据和优化查询，与数据访问控制没有多大关系。本节将介绍有关视图和查询优化的相关信息。

### 视图和索引

在创建视图时，Cloudant 会自动为视图中的数据建立索引。因此，不需要为已创建的视图中包含的数据建立索引。与视图相关的索引仅对基于视图的查询有效，而且优先于上述两种索引类型。

当发生以下 3 个事件的任何一个时，索引内容会自动增量更新：

- 向数据库添加新文档。
- 从数据库中删除现有文档。
- 更新数据库中的现有文档。

当更新包含视图定义的设计文档时，该视图的索引会得到全面更新。此外，如果更新设计文档，即使视图本身未更新，也会对它定义的所有视图进行全面更新。在数据量很大时，视图的全面更新会耗费大量时间，而且会降低数据库的性能，应该始终避免这种情况。定义视图时，可以尝试将不相关的视图拆分到不同的设计文档中。

在生产环境中，由于业务需要，新视图或全面更新有时是无法避免的。如果拥有大量数据，则应避免长期占用数据库资源，这会导致减慢整个系统，您不应该直接在 Cloudant 的生产环境中创建或更新视图。

一种方法是在测试环境中创建生产数据库的备份，并建立连续复制。然后在备份数据库中创建或更新视图，并等待操作完成。然后，需要将备份数据库同步回生产 Cloudant 实例。对于需要极快响应速度的应用程序，对视图的增量更新可能也会影响数据库性能，因为 Cloudant 会维护视图的三个副本来提高查询性能，而且它必须保持这些副本的一致性。事实上，Cloudant 提供了一些参数（如稳定、陈旧、更新），让用户来决定是否仅从某个片段接收陈旧的数据，是否接受陈旧的数据，或者是否实时刷新视图。用户需要根据自己的需求在这 3 个选项中选择合适的参数。 [Cloudant 官方文档](https://console.bluemix.net/docs/services/Cloudant/api/using_views.html#indexes) 包含详细的相关介绍，所以这里不再赘述。

### Map 函数

Cloudant 中的视图定义主要包含两个部分：map 函数和 reduce 函数。map 函数的输入是文档。在该函数中，我们可以读取文档中每个字段的值。在添加一些相关的业务逻辑后，我们可以通过 emit 函数将转换后的文档提供给 reduce 函数。由于视图默认情况下会创建相关的索引并实时更新数据，所以可以将一些相对复杂的数据级业务逻辑放入 map 函数中，以便提高数据处理效率。

emit 函数用于输出 map 函数的结果，在使用 emit 函数时，需要重点注意的是，我们应该试着仅输出必要的字段，而不是那些无用的字段或整个文档，从而加速视图创建、刷新和基于视图的查询。 [Cloudant 官方文档](https://console.bluemix.net/docs/services/Cloudant/api/using_views.html#indexes) 对 map 函数进行了清晰易懂的介绍；我们在这里就不再赘述。

### Reduce 函数

[Cloudant 官方文档](https://console.bluemix.net/docs/services/Cloudant/api/using_views.html#indexes) 中对 reduce 函数的介绍过于简单，缺少辅助说明和实际示例。因此，新用户很难独立实现复杂的 reduce 函数。所以本节将提供一个示例来解释 reduce 函数的特性，以及它的具体实现和基本原理。

可以向 reduce 函数传递 3 个参数：”keys”、”values” 和 “rereduce”。通常，reduce 函数必须处理两种情况：

当 rereduce 为 false 时：

- “keys” 是一个数组，它的元素是 [key, id] 格式的数组，其中 key 是 map 函数发出的键，id 标识了用于生成该键的文档。
- “values” 也是一个数组，其中包含为 keys 中的各个元素发出的值。

当 rereduce 为 true 时：

- “keys” 为 null。
- “values” 是一个数组，其中包含之前调用 reduce 函数所返回的值。

下面是一个具体示例。这个示例假设数据库中存储的文档是系统中的用户创建的文件。每个文件都放在相应的文件夹中。因此，每条记录包含 user\_id、folder\_id、file\_id，以及描述和 created\_timestamp 等信息。我们的目标是通过视图的 MapReduce 函数创建一个查询，该查询通过指定特定的 folder\_id 和 user\_id 来返回所有 file\_id（这个特性也可以通过 Cloudant 中的其他方式来实现，下面的示例是为了演示 MapReduce 函数而定制的）。

#### 示例 2.MapReduce 示例 – Map 函数

```
function (doc) {
if((doc.obj_type === "file")
      && doc.file_id && doc.folder_id){
        var user_id = doc.user_id;
        if(user_id){
          emit([user_id, doc.folder_id], doc.file_id);
        }
    }
}

```

Show moreShow more icon

示例 2 是 map 函数的一种实现。它将每个文档转换为一个键值对，其中键是一个由 user\_id 和 folder\_id 组成的列表，值是 file\_id。

#### 示例 3.MapReduce 示例 – Reduce 函数

```
function (keys, values, rereduce) {
if (rereduce) {
    var Rresult = "";
    for(var j in values){
      var sub_value = values[j];
      var sa = values[j].split(",");
      for(var m in sa){
        if (Rresult.indexOf(sa[m])<0){
          if(Rresult == ""){
          Rresult =  sa[m];
        }else{
          Rresult = Rresult + "," + sa[m];
        }
        }
      }
    }
    return Rresult;
} else {
    var result ="";
    for(var i in values){
      if(result.indexOf(values[i])<0){
        if(result == ""){
          result =  values[i];
        }else{
          result = result + "," + values[i];
        }
      }
    }
    return result;
}
}

```

Show moreShow more icon

示例 3 是 reduce 函数的一种实现。该代码分别处理 rereduce 为 true 或 false 时的每种情况。详细的实现过程如图 1 所示。第 1 步，执行 rereduce 为 false 时的过程，即处理图 1 中 L0 层上的数据。L0 层上的数据是 map 函数的输出，也是这里的 reduce 函数的输入。L1 层上的输入数据是在 rereduce 为 false 时 L0 层的输出。在此过程中，具有相同键的值将串联成一个字符串，并使用逗号作为分隔符。例如，对于 L1 层上最左边的键值对，值类似于 “File1, File2, File3″。在此之后，将以迭代方式实现 rereduce 为 true 时的代码。也就是说，我们使用逗号拆分 L1 层上具有相同键的输入字符串，删除重复的字符串，然后以逗号分隔方式将它们串联在一起。新字符串会发送到 L2 层。如果 L2 层上有多个父节点，则重复此过程。

##### Reduce 函数的执行过程

![reduce-函数的执行过程](../ibm_articles_img/ba-optimize-queries-cloudant_images_image001.png)

基于前面的 MapReduce 函数的业务逻辑，Cloudant 将对相关数据进行预处理，然后将 map 和 reduce 函数的结果发送给 Cloudant 索引的 B 树。稍后，如果基于此视图执行查询，相关数据很快将会返回。在使用基于视图的查询时，经常使用 group 和 group\_level 这两个参数。

### 基于视图的查询中的 group 和 group\_level

Cloudant 为基于视图的查询提供了许多可选参数，例如开始和结束键值、一次传递多个键值、指定是否执行 reduce 函数，或者限制或跳过一些记录。在这些参数中，最强大、最容易出错的参数是 group 和 group\_level。

当基于视图的查询没有参数时，如果在视图中定义了 reduce 函数，结果将是 reduce 函数的最终输出的结果。在上一节的示例中，如果直接执行一个基于视图的没有参数的查询，结果将类似于示例 4。

#### 示例 4.没有参数时基于视图的查询结果

```
Method: GET /$DATABASE/_design/$DDOC/_view/$VIEW-NAME
Request: None
Response:
{
    "rows": [
        {
            "key": null,
            "value": "F1,F2,F3,F4,F5"
        }
    ]
}

```

Show moreShow more icon

结果是用逗号相连的所有值字符串的集合。结果中返回的键为 null，因为上一步中的 reduce 函数是在 rereduce 等于 true 时执行的。但是我们的目的是获取所有 file\_id 项，并根据 user\_id 和 folder\_id 对它们进行分组。因此，结果与我们的要求不符。默认情况下，参数 “group” 等于 false。如果在查询中包含此参数并将它指定为 true，那么结果将类似于示例 5。

#### 示例 5.带参数 group 时基于视图的查询结果

```
Method: GET /$DATABASE/_design/$DDOC/_view/$VIEW-NAME?group=true
Request: None
Response:
{
    "rows": [
        {
            "key": ["U1","Fo1"],
            "value": "F1,F2,F3,F4,F5"
        },
        {
            "key": ["U2","Fo1"],
            "value": "F1,F2,F3,F4"
        },
        {
            "key": ["U2","Fo2"],
            "value": "F1,F3,F4"
        }
    ]
}

```

Show moreShow more icon

按 map 函数输出的键对结果进行分组，然后 reduce 函数对每个组中的值执行相关操作。在此基础上，如果我们传入指定的键值，即 user\_id 和 folder\_id 的组合，那么查询结果将只返回与该键对应的值。到目前为止，我们基本上满足了前面提到的要求。

使用参数 group\_level 有两个前提条件。一个是参数 group 应该等于 true。另一个是键必须是一个组合键。所谓的组合键，是指该键需要是一个数组而不是单个值。前面的示例中使用的键是一个由 user\_id 和 folder\_id 组成的组合键。满足这两个前提条件后，可以将 group\_level 指定为从 1 到键数组长度的任意整数值。group\_level 值为 n 意味着 reduce 仅按照键数组的前 n 个元素对值进行分组。基于前面的示例，示例 6 给出了 group\_level 为 1 时的查询：

#### 示例 6.参数 group\_level 为 1 时基于视图的查询结果

```
Method: GET /$DATABASE/_design/$DDOC/_view/$VIEW-NAME?group=true&group_level=1
Request: None
Response:
{
    "rows": [
        {
            "key": ["U1"],
            "value": "F1,F2,F3,F4,F5"
        },
        {
            "key": ["U2"],
            "value": "F1,F2,F3,F4"
        }
    ]
}

```

Show moreShow more icon

因为我们将 group\_level 指定为 1，所以结果中的键数组仅包含组合键的第一个元素 user\_id。共享通用 user\_id 的值将被分组到一起，无论 folder\_id 是什么。

当 group\_level 为 2 时，查询的结果与示例 5 相同。也就是说，当将参数 group 指定为 true 时，group\_level 的默认值是键数组的长度。

在本例中，读者可能想知道如何仅通过 folder\_id 键对值进行分组。目前，Cloudant 的查询接口不支持这种分配。如果想要仅通过 folder\_id 对值进行分组，一种替代方法是修改此视图的 map 函数或添加一个新视图，以便参数 folder\_id 放在组合键的最左侧位置。

### 使用视图有效地操作文档引用

不建议在 NoSQL 数据库中的文档之间建立各种引用，而是建议以冗余信息的形式将数据存储在引用的文档中。尽管如此，一些复杂的业务逻辑仍然需要文档引用。一旦文档中出现类似 ID 的字段，通常就需要使用此 ID 执行一次额外查询来获取引用文档中的数据。

为了提高对引用文档的查询效率，Cloudant 视图支持仅使用一个查询同时返回原始文档和引用文档。在 map 函数中，值的最终输出包含 “\_id” 字段，而且指定了相应的文档 ID 值。因此，在查询视图时，只要将参数 include\_docs 指定为 true，就会检索并返回与指定 ID 对应的文档。下面给出了一个示例。示例 7 是数据库中的 3 段数据。

#### 示例 7.文档引用查询的视图优化样本的源数据

```
{"_id":"2sof3204234u","node_name":"Node 1","parent":"3sldfjsla"}
{"_id":"5ladsfjldd","node_name":"Node 2","parent":"3sldfjsla"}
{"_id":"3sldfjsla","node_name":"Node 3","parent":""}

```

Show moreShow more icon

节点 1 和 2 引用了节点 3 的 ID 作为其父节点。我们现在希望在查询每个节点时找到相关的父节点。首先，我们需要创建一个新的视图。map 函数如示例 8 所示。

#### 示例 8.使用 map 函数处理文档间引用

```
function(doc) {
    if (doc.parent) {
        emit(doc.node_name, { "_id": doc.parent });
    }
}

```

Show moreShow more icon

此视图的 reduce 函数可以留空。在 map 函数中，最终输出的值是一个对象，它仅包含一个 ID 属性，我们在这里传入它的父属性的值，即父属性的 ID。然后，我们基于此视图来查询数据，结果如示例 9 所示。

#### 示例 9.使用视图来处理文档间引用后的查询结果

```
Method: GET /$DATABASE/_design/$DDOC/_view/$VIEW-NAME?include_docs=true
Request: None
Response:
{
"total_rows":3,
"offset":0,
"rows":[
      {
         "id":"2sof3204234u",
         "key":"Node 1",
         "value":{
            "_id":"3sldfjsla"
         },
         "doc":{
            "_id":"3sldfjsla",
            "node_name":"Node 3",
            "parent":""
         }
      },
      {
         "id":"5ladsfjldd",
         "key":"Node 2",
         "value":{
            "_id":"3sldfjsla"
         },
         "doc":{
            "_id":"3sldfjsla",
            "node_name":"Node 3",
            "parent":""
         }
      },
      {
         "id":"3sldfjsla",
         "key":"Node 3",
         "value":{
            "_id":""
         },
         "doc":null
      }
]
}

```

Show moreShow more icon

在上面返回的结果中，字段 “doc” 包含引用文档的内容。节点 1 和节点 2 都引用了节点 3 的 ID。通过使用视图，节点 1 和节点 2 在返回时包含节点 3 的文档数据。

但是，这种处理对文档格式有一定的要求：存在引用关系的文档应该具有相同或类似的结构。如果结构不同，则会导致引用文档被过滤掉，因此不能包含在返回的结果中。示例 10 修改了示例 7 中的源数据，从节点 3 中删除了父字段。

#### 示例 10.文档引用查询的视图优化样本的源数据（修改版）

```
{"_id":"2sof3204234u","node_name":"Node 1","parent":"3sldfjsla"}
{"_id":"5ladsfjldd","node_name":"Node 2","parent":"3sldfjsla"}
{"_id":"3sldfjsla","node_name":"Node 3"}

```

Show moreShow more icon

map 函数的定义与示例 8 中相同。处理此数据时，因为节点 3 不含 parent 字段，所以 map 函数将其排除在外，在输出来自其他节点的数据时导致引用失败。最终结果如示例 11 所示。

#### 示例 11.使用视图来处理文档间引用后的查询结果（经过修改后）

```
Method: GET /$DATABASE/_design/$DDOC/_view/$VIEW-NAME?include_docs=true
Request: None
Response:
{
"total_rows":2,
"offset":0,
"rows":[
      {
         "id":"2sof3204234u",
         "key":"Node 1",
         "value":{
            "_id":"3sldfjsla"
         },
         "doc":null
      },
      {
         "id":"5ladsfjldd",
         "key":"Node 2",
         "value":{
            "_id":"3sldfjsla"
         },
         "doc":null
      }
]
}

```

Show moreShow more icon

本节主要阐述视图中的索引和 reduce 函数。通过结合使用这些示例与 Cloudant 官方网站上的数据，读者可以更深入地了解如何创建复杂的视图。

## 与 Cloudant 查询相关的其他优化

### 用于调优请求和重试超时值的设置

由于 Cloudant 是基于云的服务，所以 Cloudant 实例通常会与其他实例共享它运行时所在的节点。因此，同一节点中的数据库之间必然存在资源竞争和相互影响。例如，如果一个节点的实例在很短时间内大量使用节点资源，同一节点中的其他实例势必会受到影响，导致对数据的请求在相当长的时间内不会被响应。

为了避免由于同一节点中其他实例的瞬时资源消耗而导致自己实例的请求延迟或处理失败，我们需要启用 Cloudant 的请求超时和重试机制。在发现请求无法及时返回时，需要通过超时机制让请求超时。然后，使用重试机制重新发送请求。在多次调优应用程序系统与 Cloudant 环境之间的连接后，需要根据特定的性能要求来设置请求超时和重试次数。

在 Cloudant 提供的软件开发工具包中，默认情况下已禁用请求超时和重试。因此，开发人员需要显式指定它。以 NodeJS 为例，示例 12 展示了如何设置请求超时。

#### 示例 12.在 NodeJS 中设置 Cloudant 请求超时的示例

```
var Cloudant = require('cloudant');
var cloudant = Cloudant({account:me,
password:password,
requestDefaults: { "timeout": 5000 }
});

```

Show moreShow more icon

在 Cloudant 初始化函数中，创建一个 requestDefaults 属性并将超时指定为 5000 毫秒。初始化该函数后，将建立到 Cloudant 数据库的连接。通过此连接发送的所有请求都将其超时设置为 5 秒。

请注意，Cloudant 初始化函数支持传递一个名为 retry 的插件，还支持设置超时和重试次数，但它们与上面的概念有所不同。Cloudant DBaaS 的一个定价细节指定了在特定价格下的每秒并发查询的最大数量。超过此最大值的查询将被拒绝，并返回一条 HTTP 429 消息。retry 插件只能在这种情况下设置查询重试次数和超时。默认情况下，超过最大值的查询将重试 3 次，超时为 500 毫秒。

### Cloudant 视图和连接操作

在传统关系数据库中，如果两个表之间存在外键关联，您可以通过一段查询和连接操作从两个表中获取信息。但在 Cloudant 中，如果两种不同文档类型的数据存在关联，通常需要使用两个 HTTP 查询请求来获取这两种文档类型的数据。

#### 示例 13.Cloudant 视图和连接操作的样本数据

```
{
    "id":"file1",
    "name":"how to learn JS",
    "user_id":"user1",
    "doc_type":"file"
},
{
    "id":"file2",
    "name":"Thinking in Java",
    "user_id":"user1",
    "doc_type":"file"
},
{
    "id":"user1",
    "name":"John",
    "doc_type":"user"
}

```

Show moreShow more icon

对于示例 13 中的数据，如果想要查找用户信息及其关联的文件信息，通常需要使用两个查询，一个用于查询用户信息，另一个用于查询文件信息。但是，通过使用 Cloudant 视图，只需一个查询即可实现此目标。示例 14 给出了这个视图的定义，它只需要 map 函数。

#### 示例 14.使用视图来处理连接操作

```
function(doc){
    if(doc.type=="user"){
        emit([doc.id,0],doc)
    }else if(doc.type=="file"){
        emit(doc.user_id,1),doc
    }
}

```

Show moreShow more icon

在此函数中，使用了一个组合键，第一个元素是用户的 ID。第二个元素在两种情况下是不同的。当文档的类型为 user 时，第二个元素为 0；当它的类型为 file 时，第二个元素为 1。然后，组合键的相应值包含特定文档的信息。

默认情况下，Cloudant 的基于视图的查询结果将按照键进行排序。具体来讲，[“abc”, 2] 位于 [“abc”] 和 [“abc”, 1] 之后，但在 [“abc”, 2, “xyz”] 之前。基于此规则，我们可以通过一个查询来获得某个用户的用户信息和文件信息。这是通过指定 startkey 和 endkey 来实现的。例如，如果想要查找 ID 为 “abc” 的用户的用户信息和文件信息，则需要指定 [“abc”] 作为 startkey，并指定 [“abc”, 2] 作为 endkey。

### 使用 HTTP 持久连接或连接池

在 Cloudant 提供的软件开发工具包中，当在客户端与 Cloudant 之间建立连接时，应用程序级协议为 HTTP 协议。从客户端应用程序对 Cloudant 的请求数量在短时间内可能是巨大的。启用 HTTP 持久连接重用了客户端与 Cloudant 之间的连接，并避免每次发出请求时建立连接的过程，从而在一定程度上加快了查询。示例 15 展示了如何在 Node.js 中启用 HTTP 持久连接。您需要先使用 npm 安装 agentkeepalive 模块。

#### 示例 15.在 Node.js 中启用 HTTP 持久连接

```
// create custom HTTPS agent
var HttpsAgent = require('agentkeepalive').HttpsAgent;
var myagent = new HttpsAgent({
maxSockets: 50,
maxKeepAliveRequests: 0,
maxKeepAliveTime: 30000
});

// Setup Cloudant connection
var cloudant = require('cloudant')({
url: 'https://myusername:mypassword@myaccount.cloudant.com',
"requestDefaults": { "agent" : myagent }
});

```

Show moreShow more icon

请注意，在一些客户端中（比如 Java 中的 Apache Http Client），HTTP 持久连接是通过 HTTP 连接池来实现的，而且具有相同的作用。

## 结束语

作为具有代表性的文档数据库和 NoSQL 数据库，Cloudant 在使用上与传统关系数据库有所不同，尤其是在数据查询和优化方面。本文中提到的方法没有涵盖所有情形，而且需要与特定的数据和使用场景结合使用。随着数据量和数据复杂性的增加，用户对 Cloudant 的查询效率会有更高的要求，优化 Cloudant 查询的方式也会增加。我们也鼓励读者对这一主题展开进一步的讨论，使得 Cloudant 中的查询变得更高效。

本文翻译自： [Optimize queries in Cloudant](https://developer.ibm.com/articles/ba-optimize-queries-cloudant/)（2018-07-11）