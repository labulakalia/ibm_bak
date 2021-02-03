# 如何使 Cloudant 中的查询更加高效
Cloudant 进阶

**标签:** 分析

[原文链接](https://developer.ibm.com/zh/articles/ba-cn-cloudant-query-efficiency/)

王传阳

发布: 2017-12-13

* * *

随着移动端应用的普及与 Web 端应用的变革，应用程序对数据的灵活性要求越来越高。近年来各种 NoSQL DB 的出现，弥补了传统关系型数据库的不足，拥有更灵活的数据模型，更好地支持应用程序对非结构化和半结构化数据的需求。作为 NoSQL DB 队伍中的一员，IBM Cloudant 面向 Web 和移动应用程序提供全面托管的数据库服务，提供高级索引技术、自定义视图、全文检索和实时查询等丰富的功能特性。本文针对 Cloudant NoSQL DB 中的各种查询，辅以实例，从多个角度阐述如何对查询进行优化的方法和心得，有助于读者更深入地理解各种查询的适用场景以及最佳用法，从而为应用程序提供更加高效的数据服务。

## Cloudant 简介

Cloudant 是 Apache CouchDB 的商业化版本，面向 Web 和移动应用程序提供全 IBM 运维管理数据平台服务。作为搭建在云上的 NoSQL 数据库，它既利用了云的可用性、弹性和影响范围，使应用程序能够扩展到更大规模，又为数据可用性、耐久性和移动性而优化，非常适合于快速增长的移动和 Web 应用。Cloudant 拥有强大的索引功能，以及跨多个数据中心和设备将数据推送到网络边缘的能力，由此加快了访问速度并提高了容错能力。它支持用户随时随地访问数据。

相对于其他数据库产品，Cloudant 具有以下特点：

- 支持副本分发与复制。数据库读写访问权可被赋予云端每个节点中，从而使Cloudant 可以在任何规模的集群中运行，并能保证数据的安全性和一致性。
- 复本之间连续的数据同步。Cloudant 支持在多个复本之间建立同步，并支持实时自动同步。一个复本中数据的持续更新可以自动同步到相关联的其他复本中。
- 提供多种优化查询的方法，包括索引和视图。在对视图的支持方面，Cloudant 支持用户自定义视图，支持的语言是Javascript。
- 提供针对JSON 文档的数据分析和仓储。Cloudant 支持与 IBM DB2 数据仓库集成，进行传统的数据仓库和商业智能分析。同时，Cloudant 也支持与 IBM dashDB 云数据仓库进行集成，进行在线报表展示和商业智能分析。
- 支持以 GeoJSON 格式存储基于位置的数据。

更多相关介绍，读者可通过 [IBM Cloudant 官网](https://www.ibm.com/cn-zh/marketplace/database-management?lnk=hm#product-header-top) 进行了解。

## Cloudant 索引对查询的优化

在数据库中，通常会针对常用数据及相关查询建立索引，从而加快查询速度。在 Cloudant 中亦如此。Cloudant 支持两种类型的索引：

- type”: “text”
- “type”: “json”

这两种类型的索引，在使用目的和使用方式上有明显不同。

从使用目的来看，text 类型的索引侧重于处理 Cloudant 文档的具体内容，而淡化文档本身的结构。因此如果用户对数据库中文档的结构不了解，或者文档结构很多，格式复杂，则可考虑建立 text 类型的索引。相对地，json 类型的索引对文档结构要求较高，它建立在一个或多个具体的域（field）之上。因此，如果用户对数据库中的文档结构很熟悉，则可以选择建立 json 类型的索引，通过明确指定相应的域，使 Cloudant 对所有包含指定域的查询进行优化。在这一点上，json 类型索引的概念更接近于传统的关系型数据库中索引的概念，都是针对具体的列或域进行优化。

从使用方式来看，text 类型的索引只能用在 Cloudant 数据搜索接口（search）中，支持的语法为基于 Lucence 的查询解析语法（ [Lucence Query Parse Syntex](http://lucene.apache.org/core/4_3_0/queryparser/org/apache/lucene/queryparser/classic/package-summary.html#Overview)）。json 类型的索引只能用于 Cloudant 的数据查询（query）接口中，支持使用 Cloudant 查询解析语法通过编写 json 对象进行查询。两种索引所对应的两种查询接口都很强大，支持丰富的自定义查询。即使在 text 类型的索引中，Cloudant 也支持用户精确地定义需要索引的数据范围，例如通过条件判断指明具体的域。因此在大多数场景中（文档结构相对简单、清晰）二者之间可以相互转换，即一种查询可通过建立另一种类型的索引，使用另一种查询接口来完成，但由于二者在使用目的上有所不同，当文档结构复杂且不清晰的情况下，一般只能使用 text 索引和数据搜索接口。

在索引建立的速度方面，相同数据量的前提下，text 类型的索引建立的要比 json 类型的慢，因为在建立 text 类型的索引时，Cloudant 不仅要处理指定的数据结构，还要处理数据结构中的内容。而 json 类型的索引则只关注于结构本身。极端情况下，text 类型的索引可被用于对整个数据库做全文检索，代码非常简单，如下所示：

#### 清单 1\. 使用 text 类型的索引建立全文检索

```

{
    "type": "text",
    "index": { }
}

```

Show moreShow more icon

另外，Cloudant 的一些高级特性，如复杂的聚集计算和基于空间地理位置的计算，也只能在基于 json 类型索引的查询（query）接口中使用。因此，使用这两种类型的索引的关键依然是理解现有数据及对数据的使用方式。

## 基于 Cloudant 视图的查询优化

在传统的关系型数据库中，视图被用于过滤数据、控制权限及优化查询。而在 Cloudant 中，视图（View）则主要被用于过滤数据与优化查询，与权限控制无直接联系。本文将从以下方面展开阐述视图与查询优化相关的内容。

### 视图与索引

在视图创建时，Cloudant 会自动创建相关的索引，因此无需再针对已创建的视图建立索引。与视图关联的索引只对基于视图的查询发生作用，且优先级高于上文谈过的两种索引。

视图中的索引在以下情况下会自动进行增量更新：

- 新文档加入到数据库中
- 已存在的文档从数据库中删除
- 已存在的文档被更新

但有一种情况会导致视图进行全量刷新，即定义该视图的设计文档（design document）发生了更新。一个设计文档可以同时定义多个视图，也可以定义诸如 List 等其他元素。如果设计文档被更新，即使更新的并非是视图本身，也会导致所有定义在其中的视图进行全量刷新。数据量很大时，视图的全量刷新耗时久且会拖慢数据库性能，因此要尽量避免这种情况出现。在定义视图时，力争做到将相互间无关联的视图拆分到不同的设计文档中。

在生产环境中，有时由于业务需要，新建视图或全量刷新视图再所难免。如果数据量很大，为减少长时间占用数据库资源拖慢整个环境，一般不会直接在生产环境的 Cloudant 中直接新建或刷新视图。有一种方法是将生产数据库建立一份备份放在测试环境的 Cloudant 实例中，并建立持续同步机制（continuous replication）。然后在备份数据库中新建或刷新视图，并等待其完成。之后将备份数据库同步复制到生产环境的 Cloudant 实例中。最后使用脚本完成新旧数据库的切换。当然，数据库切换需要生产环境有很短暂的离线时间，可考虑通过应用层缓存或其他技术缩短离线时间。

在一些对数据响应速度要求极高的应用中，视图的增量更新也有可能会产生性能影响，毕竟 Cloudant 对视图中的数据维护了三份拷贝来增加查询性能，并要保持其一致性。事实上 Cloudant 提供了一些参数（如 stable, stale, update）供用户自行决定是否只从一个分片中接收旧数据，是否能接受旧数据，以及视图是否实时刷新等。用户需要结合自己应用程序对数据的具体要求，合理搭配这三个参数。 [Cloudant 官方文档](https://console.bluemix.net/docs/services/Cloudant/api/using_views.html#indexes) 中对此做了详细介绍，在此不再赘述。

### Map 函数

Cloudant 中的视图定义主要包含两部分，map 函数和 reduce 函数。map 函数的输入参数为文档本身，在函数内部可以读取文档每个域的值，加入相关的业务逻辑后，最终通过 emit 函数将变换后的文档输出到 reduce 函数中。

由于视图会默认建立相关索引并实时更新数据，因此一些较为复杂的数据层面的业务逻辑可以放在 map 函数中，借助视图对数据查询的优化提高应用系统处理数据的效率。在使用 emit 函数时，需要注意的是，尽量只输出必需的域，避免输出无用域或整个文档，从而加快视图建立和刷新的速度，以及基于视图的查询速度。 [Cloudant 官方文档](https://console.bluemix.net/docs/services/Cloudant/api/using_views.html#indexes) 对 map 的介绍比较清楚，容易理解，在此不做赘述。

### Reduce 函数

[Cloudant 官方文档](https://console.bluemix.net/docs/services/Cloudant/api/using_views.html#indexes) 对 reduce 函数部分的介绍过于简短，且无实例辅助说明，很难帮助使用者独立实现复杂的 reduce 功能。因此本小节会借助一个例子，着重阐述 reduce 函数的功能，具体实现和底层原理。

Reduce 函数接受三个参数，分别是 keys，values 和 rereduce。通常 reduce 函数需要分别实现 rereduce 为 false 和 true 时的处理逻辑。当 rereduce 是 false 的时候：

- keys 是一个数组，其元素是[key, id]。key 是 map 当中定义的 key，id 则是整个文档的 id。
- values 同样是一个数组，其元素是 keys 数组元素对应的 value。

当 rereduce 是 true 的时候：

- keys 为空。
- values 则是在 rereduce 为 false 时，或上一轮 rereduce 为 true 的迭代计算时返回的数组值。

下面结合一个具体的实例进行阐述。该实例假设数据库里面的文档存储的是用户在某系统中建立的文件，这些文件被放在对应的文件夹中，因此每条记录包含 user\_id，folder\_id，file\_id，还有类似 description，created\_timestamp 等信息。我们的目的是希望通过视图的 map&reduce 函数，来建立一个查询，可以通过指定 folder\_id 和 user\_id，来返回所有的 file\_id（该功能也可通过其他方式在 Cloudant 中实现，此处仅用于演示 map&reduce 函数的自定义实现）。

#### 清单 2\. Map&Reduce 实例 – Map 函数实现

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

清单 2 是 map 函数的实现，主要是将每个文档记录转换为以 user\_id 和 folder\_id 为 key，以 file\_id 为 value 的 key-value 对。

#### 清单 3\. Map&Reduce 实例 – Reduce 函数实现

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

清单 3 是 reduce 函数的实现，代码中对 rereduce 为 true 和 false 的情况分别做了实现。详细的执行过程如图 1 所示。第一步，首先执行的是 rereduce 等于 false 的情况，即对图 1 中 L0 层的输入数据做处理。L0 层的数据是 map 函数的执行结果，在此作为 reduce 函数的输入。L1 层的数据为 rereduce 等于 false 这个分支的计算结果。该分支是将同一个 key 的 values 用逗号连结起来，形成一个字符串，如 L1 层最左侧的 key-value，其 value 经过计算后会变成”File1,File2,File3 “。其次会迭代地执行 rereduce 等于 true 的逻辑，即先将 L1 层输入的字符串以逗号为分隔符拆分，然后仍以 key 为 group，对相同 key 的 value 去重并再次用逗号连结，输出一个新的字符串，放到 L2 层。L2 层如果还有多个父节点，则继续迭代地执行 rereduce 等于 true 的逻辑。

##### 图 1\. Reduce 函数的执行过程

![alt](../ibm_articles_img/ba-cn-cloudant-query-efficiency_images_image001.png)

基于以上 map 和 reduce 函数的业务逻辑，Cloudant 会预先处理相关数据，并把 map 和 reduce 后的结果存放在 Cloudant 索引的 B-tree 中。之后如果执行基于该视图的查询时，相关数据便可快速返回。而在使用基于视图的查询时，通常都会结合 group 和 group\_level 这两个参数。

### 基于视图查询中的 group 和 group\_level

Cloudant 为基于视图的查询提供了很多可选参数，例如可指定起始和终止的 key 值，可一次传入多个 key 值，可指定是否执行 reduce 函数，还可限制或跳过一些记录。在这些参数当中，功能强大且容易用错的参数是 group 和 group\_level。

当一个基于视图的查询不带任何参数时，如果视图中定义了 reduce 函数，则其结果为 reduce 函数最终输出的结果。在上一小节中使用到的例子中，如果直接执行不带参数的基于视图的查询，得到的结果将会如清单 4 所示：

#### 清单 4\. 不带参数的基于视图的查询结果

```

Method: GET /$DATABASE/_design/$DDOC/_view/$VIEW‑NAME
Request: None
Response:
{
    "rows":
        {
            "key": null,
            "value": "F1,F2,F3,F4,F5"
        }
    }

```

Show moreShow more icon

结果正是所有 value 的一个集合字符串，中间用逗号连接。因为 reduce 函数最终执行的是 rereduce 为 true 的情况，所以返回结果中的 key 是 null 值。而我们实现这个视图的目的是想得到在给定的 user\_id 和 folder\_id 的前提下，其所对应的 file\_id。因此该结果与我们的需求并不相符。group 参数默认情况下等于 false。如果在查询中加入 group 参数，指定其为 true, 则结果将会如清单 5 所示：

#### 清单 5\. 带 group 参数的基于视图的查询结果

```

Method: GET /$DATABASE/_design/$DDOC/_view/$VIEW‑NAME?group=true
Request: None
Response:
{
    "rows":
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
    }

```

Show moreShow more icon

结果会以 map 函数输出的 key 为分组，reduce 函数针对每一个分组进行 value 的操作。在此基础上，如果我们传入指定的 key 值，即 user\_id 和 folder\_id 的一个组合，查询结果便只会返回与该 key 组合相对应的 value 值。到这一步才基本满足前面提到的实现该例中视图查询的目的。

group\_level 这个参数在使用时需要有两个前提，一是 group 等于 true，二是 key 必须是复合 key。所谓复合 key 指的是 key 值需要是一个数组而非一个单独的值。前面使用的例子中的 key 便是由 user\_id 和 folder\_id 组成的复合 key。满足这两个前提后，group\_level 可被指定为从 1 到 key 数组长度之间的任意整数值。group\_level 值为 n 意味着 reduce 只会将 key 数组的前 n 个元素相同的 value 做 group 操作。基于之前的例子，清单 6 展示的是 group\_level 为 1 的查询结果：

#### 清单 6\. group\_level 为 1 的基于视图的查询结果

```

Method: GET /$DATABASE/_design/$DDOC/_view/$VIEW‑NAME?group=true&group_level=1
Request: None
Response:
{
    "rows":
        {
            "key": ["U1"],
            "value": "F1,F2,F3,F4,F5"
        },
        {
            "key": ["U2"],
            "value": "F1,F2,F3,F4"
        }
    }

```

Show moreShow more icon

由于我们指定了 group\_level 为 1，因此结果中的 key 数组只包含了复合 key 中的第一个元素，即 user\_id。拥有相同 user\_id 的 value 会被最终放到一起，而不管这些 value 对应的 folder\_id 是否相同。

当 group\_level 为 2 时，得到的查询结果与清单 5 中的相同。即当指定 group 为 true 时，group\_level 默认值相当于是 key 数组的长度值。

在此例中，可能读者会思考如何只用 folder\_id 做为 key 进行 group。目前 Cloudant 的查询接口中并不支持这种 group 的指定，而只能支持顺序地从第一个元素开始，通过增加 group\_level 值进而增加 group 的元素。针对该需求的一种替代方案是新增一个视图或修改该视图的 map 方法，将 folder\_id 放在复合 key 的最前面。

### 利用视图高效处理文档中的引用

在 NoSQL 数据库中并不建议建立文档之间的各种引用关系，而是推荐在需要引用其他文档的文档中以冗余信息的方式存储被引用文档的数据。尽管如此，还是会有一些复杂的业务逻辑需要用到文档引用。一旦文档中出现了类似文档 id 的字段，通常还需要附加多一次查询使用该 id 把被引用文档的数据查询出来。

为提高对被引用文档的查询效率，Cloudant 视图支持通过一次查询，把原文档和被引用文档同时返回。在 map 函数中，最终输出的 value 中如果包含”\_id”字段，并指定相应的文档 id 值，则在进行视图查询时，只要指定 include\_docs 为 true，被指定 id 所对应的文档便会一并被查出并返回。下面举例说明。清单 7 是数据库中的三条数据:

#### 清单 7\. 视图优化文件关联查询样例的源数据

```

{"_id":"2sof3204234u","node_name":"Node 1","parent":"3sldfjsla"}
{"_id":"5ladsfjldd","node_name":"Node 2","parent":"3sldfjsla"}
{"_id":"3sldfjsla","node_name":"Node 3","parent":""}
```

Show moreShow more icon

Node 1 和 Node 2 引用了 Node 3 的 id 作为其 parent。我们现在希望在查询每个 node 时把关联的 parent 一起查出。首先新建视图，编写 map 函数如清单 8 所示：

#### 清单 8\. map 函数处理文档间引用

```

function(doc) {
    if (doc.parent) {
        emit(doc.node_name, { "_id": doc.parent });
    }
}

```

Show moreShow more icon

该视图的 reduce 函数可留空。map 函数中，最终输出的 value 是一个对象，其中只包含一个 id 属性，这里我们传入的是 parent 属性值，即 parent 的 id。然后我们基于该视图进行查询，结果如清单 9 所示：

#### 清单 9\. 视图处理文档间引用后的查询结果

```

Method: GET /$DATABASE/_design/$DDOC/_view/$VIEW‑NAME?include_docs=true
Request: None
Response:
{
"total_rows":3,
"offset":0,
"rows":
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
}

```

Show moreShow more icon

以上返回结果中包含了被引用文档的内容。Node 1 和 Node 2 由于引用了 Node 3 的 id，通过视图对 parent id 的处理，可使 Node 1 和 Node 2 在返回时同时包含 Node 3 的文档数据。可见，Cloudant 视图的这种处理方式优化了对文档间存在引用关系的查询。然而，这种处理对文档格式有一定要求，即存在引用关系的文档间结构要相同或相近。如果结构不同，会导致被引用的文档在视图中被过滤掉而无法被包含在返回结果中。清单 10 对清单 7 中的源数据做了一处改动，即将 Node 3 的 parent 字段移除了：

#### 清单 10\. 视图优化文件关联查询样例的源数据(修改后）

```

{"_id":"2sof3204234u","node_name":"Node 1","parent":"3sldfjsla"}
{"_id":"5ladsfjldd","node_name":"Node 2","parent":"3sldfjsla"}
{"_id":"3sldfjsla","node_name":"Node 3"}
```

Show moreShow more icon

视图的 map 函数定义依然与清单 8 一样。在处理这些数据时，由于 Node 3 不包含 parent 字段，所以它会被 map 函数排除，从而导致其他 Node 的数据在输出时，无法引用到它。所以最终的返回结果会变成下方清单 11 所示：

#### 清单 11\. 视图处理文档间引用后的查询结果（修改后）

```

Method: GET /$DATABASE/_design/$DDOC/_view/$VIEW‑NAME?include_docs=true
Request: None
Response:
{
"total_rows":2,
"offset":0,
"rows":
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
}

```

Show moreShow more icon

此章节主要对视图中的索引与 reduce 函数做了详细阐述，结合其中的例子与 Cloudant 官网的资料，读者可对如何创建复杂的视图有更深入的理解。

## 与 Cloudant 查询相关的其他优化

### 设置查询请求的超时与重试

由于 Cloudant 基于云服务，Cloudant 实例运行所在的节点一般而言是与其他实例共享的。因此，同一节点中的数据库之间难免会出现资源竞争和相互影响的现象。例如，如果节点中某个实例短时间内对节点资源消耗量爆涨，则该节点中另外的实例势必会受到影响，甚至会导致数据请求长时间不返回。

为避免节点中其他实例瞬时资源消耗导致自己实例中的请求延迟或处理失败，我们需要打开 Cloudant 连接中的请求超时与重试机制，当发现请求不能及时返回时，通过超时机制让该请求超时，然后利用重试机制重新发出请求。请求的超时时间与重试次数需要通过多次调试应用系统与 Cloudant 环境间的连接，根据具体性能要求合理设置。

Cloudant 所提供的软件开发工具包中，默认的连接参数中并没有开启请求超时与重试，因此需要开发人员显式地指定。以 NodeJS 为例，清单 12 给出的是设置请求超时的方法：

#### 清单 12\. NodeJS 中设置 Cloudant 请求超时示例

```

var Cloudant = require('cloudant');
var cloudant = Cloudant({account:me,
                         password:password,
                         requestDefaults: { "timeout": 5000 }
                        });

```

Show moreShow more icon

在 Cloudant 初始化函数中，新建属性 requestDefaults，并指定 timeout 为 5000 毫秒。该函数初始化后会建立与 Cloudant 数据库的连接，使用此连接发送的所有请求，其超时时间都会被设置为 5 秒。

需要注意的是，Cloudant 初始化函数支持传入一个叫 retry 的插件，该插件同样支持设置超时时间和重试次数，但其作用与本小节的出发点完全不同。Cloudant 数据库服务的收费指标中有一项规定了当前价格下每秒钟最大的并发查询数量。超过最大值的查询会被拒绝并返回 HTTP 429 消息。而 retry 插件便是只针对该情况下对查询重试次数和超时时间的设置。默认情况下，超过最大值的查询会被重试 3 次，超时时间是 500 毫秒。

### Cloudant 视图与 Join 操作

在传统的关系型数据库中，两个表之间如果有外键关联，则可通过 Join 操作用一次查询来得到来自两个表的信息。然而，在 Cloudant 中，如果两个不同文档类型的数据存在关联，一般情况下需要通过两次 HTTP 查询请求才能拿到相关联的这两个文档类型的数据。

#### 清单 13\. 视图与 Join 操作的样本数据

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

对于清单 13 中的数据，如果想要查出用户的信息以及其关联的文件信息，通常需要两次查询，一次用于查询用户信息，另一次用于查询文件信息。但如果使用 Cloudant 视图，只需要一次查询便可同时返回用户信息与关联的文件信息。清单 14 给出的是视图的定义，该视图只需要 map 函数。

#### 清单 14\. 使用视图处理 Join 操作

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

在该函数中，使用了复合键，第一个元素是 user 的 id。第二个元素分两种情况做了处理。对于 user 类型的文档，第二个元素为 0；对于 file 类型的文档，第二个元素为 1。而复合键对应的值，则包含了具体的文档信息。

由于 Cloudant 对基于视图查询的返回结果默认都是以 key 进行排序的。例如，[“abc”,2]会排在[“abc”]和[“abc”,1]之后，但会排在[“abc”, 2, “xyz”]之前。基于这个规则，我们可通过上方视图实现一次性查出给定用户的用户信息和文件信息。实现方式便是通过指定 startkey 和 endkey。例如，如果想查出用户 id 为 abc 的用户信息和文件信息，则需要指定 startkey 为[“abc”], endkey 为[“abc”,2]。

### 使用 HTTP 长连接或连接池

在 Cloudant 提供的软件开发包中，建立客户端与 Cloudant 之间的连接时，底层是基于 HTTP 协议的。应用程序作为客户端向 Cloudant 发起的请求数量在短时间内可能很大。启用 HTTP 长连接可以对客户端与 Cloudant 之间的连接进行复用，避免每次发出请求时建立连接的过程，从而在一定程度上提高查询速度。清单 15 给出的是在 Node.js 中开启 HTTP 长连接的方法，需要首先使用 npm 安装 agentkeepalive 模块。

#### 清单 15\. Node.js 中启用 HTTP 长连接

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

需要注意的是，在一些客户端中（如 Java 中的 Apache Http Client），HTTP 的长连接是通过 HTTP 连接池实现的，其作用是等效的。

## 总结

Cloudant 作为具有代表性的文档数据库和 NoSQL 数据库，在使用方式上与传统的关系型数据库存在区别，在数据查询和优化方面区别尤为突出。本文所提到的各种方法并非适用于任何情况，需要与具体的数据和使用场景相结合使用。随着数据量的增大与数据复杂度的提高，用户对 Cloudant 的查询效率会有更高要求，优化 Cloudant 查询的方法也在不断增加。在此也欢迎读者朋友就此话题展开讨论，共同使 Cloudant 中的查询更加高效。

## 参考资源

- 参考 [Cloudant 官方文档](https://console.bluemix.net/docs/services/Cloudant/getting-started.html#getting-started-with-cloudant)：全面了解 Cloudant。
- 参考 [CouchDB 官方文档](http://docs.couchdb.org/en/2.1.0/)：了解 Cloudant 相近的底层实现。
- 阅读文章 [Finding your data with Views](http://guide.couchdb.org/draft/views.html)：更好地理解 Reduce 的工作原理。