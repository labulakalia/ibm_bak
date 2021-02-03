# Python API 类型系统的设计与演变
基于 Python 2 的在线服务如何拥抱 Python 3 新特性

**标签:** API 管理,Python,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/wa-lo-python-api-type-sysytem-design-evolution/)

李宇飞

发布: 2018-01-24

* * *

## API 与类型系统

由于众所周知的原因，至今仍有大量生产环境的代码跑在 Python 2.7 之上，在 Python 2 的世界里，并没有一个官方的类型系统实现。那么生产环境的类型系统是如何实现的呢，为什么一定要在在线服务上实现类型系统？下文将针对这两个问题进行深入讨论。

### 什么是 API 的类型系统

人们常说一门编程语言的类型系统，通常指一门编程语言在表达式的类型意义上所具有的表达能力。而对于 API 来说，对于其输入（参数）和输出（响应）都能够有完善的类型表达，那么就可以认为它具有了基本的类型系统。

一个包含了方法（Method/HTTP verb）和路径（Path）的 API，常常称之为一个访问点（endpoint）或 API，每一个 API 具有一个描述性质的声明，称之为 Schema，Schema 可以有多种定义方式，但至少会包含参数（请求字段及其类型定义）和响应（状态码，响应字段及类型定义）。比较典型的是 [OpenAPI](https://github.com/OAI/OpenAPI-Specification) 规范的定义，该规范将在下文详细介绍。

那么在线服务上实现类型系统有何意义？如果一个 API Framework 或者 RPC Remote Call 没有类型系统，会出现什么样的问题呢？

### 为什么要在在线服务上实现类型系统

本文认为在线服务上的类型系统至少有以下几种直接的作用：

- 验证参数的可靠性，由于在服务开发时，不能信任用户的输入，应做好最坏的假设，就如同墨菲在静静地看着你。
- 自动生成文档和超文本链接，一个完善的 Schema 系统，可以为 [HATEOAS](https://en.wikipedia.org/wiki/HATEOAS) （Hypertext As The Engine Of Application State） 提供支持。
- 自动生成 Definition 文件（比如 [thrift](https://thrift.apache.org/) ， [protobuf](https://developers.google.com/protocol-buffers/) 等 RPC 定义），用于在服务端提供兼容多种协议的网关，在客户端为终端用户提供本地验证机制。
- 和异常系统结合，可以为异常诊断和 Traceback 提供支持，使用更有针对性的诊断方式。
- 可以和接口测试相结合，推断返回值的类型（但 Python 2 的库实现比较庞杂，很难实现这一点）。

### 安全性和可解释性

API 类型系统的作用，最终可以总结为在 「 安全性 」和「 可解释性」 上的提升。

如果没有一个一致的类型系统，往往要使用大量冗余代码（自定义函数）来进行参数校验，而非通过自定义类型来验证。并且耗费大量的精力人工编写接口文档，在接口变更后还要人工修改和校对。

在类型系统中，安全性和可解释性是互相依存的关系，仅从安全性考虑，如果代码结构合理，使用自定义函数进行参数校验也是可以接受的，但函数在可解释性上是弱于类型系统的，对于接口附加的元信息（比如参数类型，参数是否可选，参数描述）难以自然地表述。

类型系统在提升了安全性的同时，还兼顾了系统的可解释性，这是在服务治理上非常需要的一点。

## 类型系统实践

下面以 Python 2.7 为例，详细介绍下如何在一个在线服务上实现类型系统，以及类型系统可以帮助研发人员做哪些有意义的事情。

### marshmallow

Python 2 中没有一个官方的类型系统实现，所以在 API 参数的验证中，往往是通过外挂第三方 Schema 实现的。

marshmallow 是本文选用的一个对类型系统进行建模的 Python 库，它有着极高的流行程度，提供了基本的类型定义、参数验证功能和序列化 / 反序列化机制。

现在假设研发团队要开发一个用户相关的接口，首先要对用户这个服务资源进行抽象定义，一个基本的 Schema 定义如下：

##### 清单 1\. 一个用户接口参数模式定义

```
# -*- coding: utf-8 -*-

import re
from marshmallow import Schema, fields, validate
from myapp import fields as myfields

class UserSchema(Schema):
    user_id = myfields.UserId(required=True, help=u'用户的唯一 ID')
    nickname = fields.Str(required=True,
                          validate=validate.Length(min=2, max=20),
                          help=u'用户的昵称')
    email = fields.Email(required=True, u'用户的邮箱，不可重复')

```

Show moreShow more icon

marshmallow 自带了许多内建类型，比如 Email，URL，UUID 等，研发人员也可以根据业务来定制自定义类型，比如上文的 UserId 可以像这样定义：

##### 清单 2\. 自定义类型示例

```
# -*- coding: utf-8 -*-

import re

class UserId(fields.Field):
    """ 长度为 10 - 17 的，由字母、数字、下划线组成的 ID """
    pattern = re.compile(r'^[a-zA-Z0-9\_]{10-17}$')

    # 必选的
    default_error_messages = {
        'invalid': u'不是一个有效的用户 ID',
        'format': u'{value} 无法被格式化为 ID 字符串',
    }

    def _serialize(self, value, attr, obj):
        return value

    def _deserialize(self, value, attr, data):
        # 可以使用任何验证方式，而不仅仅是正则表达式
        if not self.pattern.match(value):
            self.fail('invalid', value=value)
        return value

```

Show moreShow more icon

服务开发人员也可以自己写装饰器或使用开源的库,比如 [webargs](https://github.com/sloria/webargs) 来根据这个 Schema 做参数验证（以 Flask 为例）：

##### 清单 3\. Web 框架集成示例

```
# -*- coding: utf-8 -*-

from flask import Flask, jsonify
from webargs.flaskparser import use_args

from myapp.schema import UserSchema

app = Flask(__name__)

@app.route('/', methods=('GET',))
@use_args(UserSchema)
def echo_user(args):
    return jsonify(**args)

if __name__ == '__main__':
    app.run()

```

Show moreShow more icon

在生产环境的服务中，通常会选择重载 API 注册用的装饰器（比如 @app.route 和 @use\_args）来收集 API 的定义存储到一个全局的对象里（可能是远程对象），来实现框架级的 API 反射机制，以允许服务实例在运行时拿到所有已注册的 API 的声明，以给第三方工具 / RPC 客户端提供最新的 Schema。

在上面的代码定义里，大家可以发现 API 类型系统中几个重要的功能都已经存在了：

- Schema 允许以接口为粒度定义类型声明
- fields 允许自定义类型（包括类型的校验规则，描述和错误信息）
- validate 允许自定义校验规则
- webargs 帮助类型系统与框架进行集成

但仅仅有这些就够了吗？

### validator 和枚举

在繁忙的业务系统开发过程中，通常需要一定程度的抽象来增强代码的可重用性，比如正则表达式和枚举等。

枚举是一种特殊的类型，在线服务对它的可描述性有着更多的诉求。在阅读一个 API 的定义时，人们看到枚举字段，不仅仅想看到这个字段期望什么样的枚举值，更想看到每一个枚举值所代表的涵义，这就要求类型系统扩展（或许是约束）枚举值的定义。

Python 内置的枚举类型有它的优势，但枚举值使用了包装类型，取值时需要通过 .value 函数来获取，而本文所描述的服务已经在线上运行许久了，改造工程浩大，于是采用了类似于 Flask Config Object 的定义风格。

##### 清单 4\. 一种可选的枚举声明定义

```
class UserStateEnum(object):
    OK = 0
    PENDING = 1

    __desc__ = {
        OK: u'有效用户',
        PENDING: u'封禁用户'
    }

```

Show moreShow more icon

通过定义一个类，约定类属性名大写为枚举属性，描述信息放在特殊的字段里，以此来表示枚举类型。

这是一个关键的思维模式： **在线服务 在扩展时必须要考虑 API 的可解释性** 。

### 异常和 RFC 4918

在线服务对于异常系统的诉求是将异常按照危重等级进行分离，保证高危异常的可追溯性，以及低危异常的可解释性。

在理想的情况下，可以把异常简单分为三类：

- 系统异常，由于系统故障或程序 Bug 导致的，应及时发送到 Issue Tracking 的系统中并发送警报。
- 业务异常，由于用户的输入不符合业务逻辑导致的异常，比如用户不存在。可以从日志中审计，可能会需要进行 Issue Tracking，无需报警。
- 参数错误，用户的输入不符合文档约定（契约），比如期望参数是一个 URL，但传来一个普通字符串。同样可以从日志中审计，但无需进行 Issue Tracking，无需报警。

在责权划分上，类型系统应该只包含了第三类异常，不涉及业务逻辑和系统异常的处理。

由于本文所描述的 Web 层遵循 REST 语义来进行服务开发，最早的 HTTP Status 使用了 500，随着类型系统的完善，响应状态码也逐渐细分，上面三类异常分别对应 500、400、422 三种 Status Code。

关于 422 状态码的选取，可以参考 [RFC 4918](http://www.ietf.org/rfc/rfc4918.txt) 和参考文献中一些有益的讨论。

### OpenAPI 与可解释性

对于在线服务的描述和定义，本文比较倾向于参考 OpenAPI 规范，原因是它对机器更加友好，有着严谨的 Spec 定义，有利于生成和分析，同时背后有谷歌、微软等商业公司和强大的社区支持。

相对于 [API BluePrint](https://apiblueprint.org/) 、 [RAML](https://raml.org/) 等规范所强调的人类可读性（Human Readable）， [Swagger](https://swagger.io/) 更加注重定义的规范化和通用性，鼓励社区共同推进规范的演进，在本文写作时， [OpenAPI Specification(OAS) 3.0](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.1.md) 已经发布，一个欣欣向荣的社区也是影响本文选型的关键因素。

类型系统在这里的作用是，对在线服务的接口定义进行描述，并生成一个符合 OpenAPI 规范定义的 JSON 文档，以支持文档生成工具（比如 [Swagger](https://swagger.io/) ）、前端 Mock 工具（比如国内的 [Easy-Mock](https://github.com/easy-mock/easy-mock) ）、接口测试工具（比如下文提到的基于 [py.test](https://docs.pytest.org/en/latest/) 的实现）和前端验证库的需要。

在 OpenAPI 规范中，与类型系统相关的部分主要集中在 paths、schema、data types 三个章节，本文主要实现 data types 章节中所描述的类型与 [marshmallow](http://marshmallow.readthedocs.io/en/latest/) 类型之间的映射，这里举几个特殊的例子。

##### 表 1\. OAS Data Type 与 Marshmallow Type 的映射

**OAS Type****OAS Format****Marshmallow****描述**stringemailEmail电子邮件stringuuidUUIDUUIDintegerenumEnum(Int)上文中定义的枚举类型stringList(Str)字符串列表

在 [OpenAPI](https://github.com/OAI/OpenAPI-Specification) 的定义里，每一个类型（type）都有一个可选的格式（format）可以定义，通常是根据业务所需来定制，这里取 fields 类的类名（小写）作为 format 值。

这里有一个特例，对于容器类型，比如 Enum 和 List，它们的类型取决于它所包装的类型，对于在线服务，常常需要类型系统具有确定性，是不允许 Union 类型存在的，这样设计主要是为了减少序列化 / 反序列化的成本，同时简化代码的分支逻辑。

这里举例说明容器类型的类型定义是如何翻译成 [OpenAPI](https://github.com/OAI/OpenAPI-Specification) 的类型定义的：

##### 清单 5\. List(Int) 翻译为 OpenAPI/OAS 示例

```
{
"type": "array",
"items": {
    "type": "integer"
}
}

```

Show moreShow more icon

##### 清单 6\. Enum(Int) 翻译为 OpenAPI/OAS 示例

```
{
"schema": {
    "type": "integer",
    "enum": [
      400
      404
    ]
}
}

```

Show moreShow more icon

### 接口测试与文档生成

在完成了上述基础的工作之后，就要与测试框架进行集成了。

类型系统与测试框架集成的意义是什么呢？可以分两个类别来看待：

- 第一个类别是需要严格限定接口响应字段的类型，这个时候开发人员会在代码中对接口的响应做类型声明，那么在测试用例中，类型系统的作用自然就是对响应字段类型的校验了，本文称之为严格模式。
- 第二个类别是接口响应无类型声明，那么接口的响应定义就不再具备可解释性，而可解释性对自动化的文档生成是最重要的因素。本文所描述的在线业务处于这样一个阶段，所以在类型系统实现中主要解决的就是这个问题。

如果没有响应参数的类型定义，就需要推导响应的类型，类型推导的方式有两种，静态的和动态的（运行时）。

静态分析在 Python 2 中的实现难度比较高，因为大量的第三方库都没有明确的类型信息，同时许多要经过网络的上下游服务也都没有提供严格的定义，难以在这样复杂的环境中通过静态分析拿到接口响应类型信息。

由于团队有写接口测试的习惯，最终选择了在运行接口测试的时候，和 Python 的测试框架 py.test 集成，通过收集接口测试的返回值来做运行时的类型推导。

下面尽可能简单地描述一下一个真实的实现，本文使用 [yaml](http://www.yaml.org/) 来做用例的定义，比如：

##### 清单 7\. 使用 Yaml 描述的测试用例示例

```
- uri: /echo
method: GET
desc: 测试 ECHO 服务
status: 200
params:
    ping: "pong"
responses:
    ping: "pong"

```

Show moreShow more icon

Pytest 提供了参数化的功能可以用来生成用例，apis 是用例定义的列表：

##### 清单 8\. 描述文件与 pytest 集成的示例

```
@pytest.mark.parametrize("case", apis)
def test_api(case, case_manager, mocker):
    case_obj = case_manager.add(case)
    case_obj.run(mocker)
    print(case_obj.real_response)

```

Show moreShow more icon

用例执行后，用例的响应被保存下来，再尝试对每一个响应字段的值做一个简单的类型推导。

##### 清单 9\. 一种响应值类型推导的实现示例

```
pattern_inferer_map = {
    date_pattern: {'type': 'string', 'format': 'date'},
    datetime_pattern: {'type': 'string', 'format': 'date-time'},
    ip_pattern: {'type': 'string', 'format': 'ip'},
    uuid_pattern: {'type': 'string', 'format': 'uuid'},
    base64_pattern: {'type': 'string', 'format': 'byte'},
    // ...
}

def infer_value(value):
    if isinstance(value, string_types):
        for pattern, type_info in pattern_inferer_map.items():
            if pattern.match(value):
                return type_info
        return {'type': 'string'}
    elif isinstance(value, int):
        return {'type': 'number', 'format': 'int64'}
    elif isinstance(value, float):
        return {'type': 'number', 'format': 'double'}
    elif isinstance(value, bool):
        return {'type': 'boolean'}

def inferer_response(response):
    return {k: infer_value(v) for k, v in response.items()}

```

Show moreShow more icon

暴力地对 Python 类型和 OAS 的类型做一个映射，这样就用最简单的办法完成了一个接口响应的类型推断。

很容易看出，这样的类型推断会存在许多问题，比如 int 和 float 类型的精度无法表达，字符串类型的 format 可能会有误判，尤其依赖完备的测试用例等等。

但本文为什么仍然愿意推荐这种方法，因为它可以使用最小的成本，最大限度地满足研发人员的基本诉求——拿到接口相应的基本类型信息，提升可解释性，这是类型系统中非常重要的一部分。

### 小结

就这样，本文通过重载服务框架的路由装饰器来收集 API 的参数类型信息，通过接口测试来收集 API 的响应类型信息，通过注册自定义的枚举类型和业务类型，再配合框架本身的属性，就可以生成定制化的、符合 OpenAPI 规范的文档了。

拥有类型系统的在线服务，在接口校验、异常处理、测试和文档生成等方面都有全方位的提升，满足了工程师们对一个服务在安全性和可解释性上的基本诉求，这是非常值得投入的一件事。

## 新世界的战鼓

上文介绍了过去两年间，我在 Python 2 在线服务类型系统中的一些思考与实践。与此同时 Python 也在迅速发展，包括 Instgram 在内的诸多公司，已将 Python 3 应用于生产环境了。

基于 Python 3 标准库的类型系统已经有了许多成功的尝试，比如 Web 服务框架 [apistar](https://github.com/encode/apistar) ，类型检查设施 [mypy](http://www.mypy-lang.org/) 等等。在线服务的类型系统在 2/3 上的实现并无明显差异，唯一要注意的是标准库的充分使用。参考文献中会介绍 Python 3 类型系统一些重要的参考文献导读，这是我在了解 Python 3 的类型系统时阅读的一些资料，如果有小伙伴也在做 2 to 3 的一些工作，想必会有所裨益。

如果 Python 世界真的能够迎来一个统一的类型系统实现，对于在线业务系统来说意义是非常重大的，我对此充满了期待。