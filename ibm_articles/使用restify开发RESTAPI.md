# 使用 restify 开发 REST API
开发 API 的流行框架

**标签:** API 管理,Node.js,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/wa-lo-use-restify-develop-rest-api/)

成富

发布: 2017-12-05

* * *

随着 Node.js 的流行，JavaScript 在服务器端的开发中也逐渐占有了一席之地。相对于 Java 或 C#等传统后端编程语言，JavaScript 的优势在于语法灵活，浅显易懂，上手简单，用较少的代码就可以完成很多复杂的任务。Node.js 平台上也有很多优质的第三方库。在服务器端应用中，API 有着很重要的地位，是与前端进行交互的基础。本文介绍的 restify 是一个开发 API 的流行框架，被 npm 和 Netflix 等公司所使用。

## restify 入门

restify 是一个 Node.js 模块，只需要通过 npm 或 yarn 来安装即可。本文中使用的是 yarn。所有本文的示例代码采用的都是 ECMAScript 6 的语法。因此需要使用 Babel 来进行转换。代码清单 1 给出了使用 yarn 安装 restify 和 Babel 的命令。除了通过 Babel 转换 JavaScript 代码之后再用 node 命令运行之外，还可以使用 babel-node 命令直接运行。可以把 babel-cli 安装成全局模块以方便调用 babel-node。

##### 清单 1\. 安装 restify

```
$ yarn add restify
$ yarn add --dev babel-cli babel-preset-env

```

Show moreShow more icon

在安装完 Babel 之后，还需要在项目根目录下添加代码清单 2 中的.babelrc 文件。

##### 清单 2\. .babelrc 文件

```
{
"presets": ["env"]
}

```

Show moreShow more icon

完成了相关的设置之后，就可以创建我们的第一个基于 restify 的 REST API，依然是经典的 Hello World 示例。在代码清单 3 中，我们通过 createServer() 方法来创建一个服务器，然后使用 get() 方法来创建一个路径为 /hello 的处理 GET 请求的路由。该路由的处理方法中使用响应对象的 send() 方法来输出 Hello World。最后启动这个服务器，并在 8000 端口监听。这样就完成了整个 API 的代码。

##### 清单 3\. 使用 restify 创建 Hello World 应用

```
const restify = require('restify');

const server = restify.createServer();

server.get('/hello', (req, res, next) => {
res.send('Hello World');
return next();
});

server.listen(8000, () => console.log('%s listening at %s', server.name, server.url));

```

Show moreShow more icon

## 响应处理链

对于每个 HTTP 请求，restify 通过一个响应处理链来对请求进行处理。代码清单 3 中的 next() 方法的作用，就是调用处理链中的下一个处理器。restify 中有三种不同的处理链。

- pre：在确定路由之前执行的处理器链。
- use：在确定路由之后执行的处理器链。
- {httpVerb}：一个路由独有的处理器链。

通过 restify 服务器的 pre()方法可以注册处理器到 pre 处理器链中。对所有接收的 HTTP 请求，都会预先调用该处理器链中的处理器。该处理器链的执行发生在确定路由之前，因此即便是没有路由与请求相对应，也会调用该处理器链。该处理器链适合执行日志记录、性能指标数据采集和 HTTP 请求清理等工作。典型的应用场景包括记录所有请求的信息，类似于 Apache httpd 服务器的访问日志功能，以及添加计数器来记录访问相关的性能指标。

通过 restify 服务器的 use() 方法可以注册处理器到 use 处理器链中。该处理器链在选中了一个路由来处理请求之后被调用，但是发生在实际的路由处理逻辑之前。对于所有的路由，该处理器链中的处理器都会执行。该处理器链适合执行用户认证、应用相关的请求清理和响应格式化等工作。典型的应用场景包括检查用户是否完成认证，对请求和响应的 HTTP 头进行修改等。

在每个处理器的实现中，应该在合适的时机调用 next()方法来把处理流程转交给处理器链中的下一个处理器。具体的时机由每个处理器实现根据需要来决定。这给处理 HTTP 响应带来了极大的灵活性，也使得处理器可以被有效复用。每个处理器的实现逻辑也变得更加简单，只需要专注于完成所设计应有的功能就可以了。在处理完成之后，调用 next()方法即可。在某些情况下，可能不需要由处理器链中的后续处理器来继续进行处理，比如 HTTP 请求的格式是非法的。这个时候可以通过 next(false)来直接终止整个处理器链。在调用 next()方法的时候，也可以传入一个 Error 对象，使得 restify 直接把错误信息发送给客户端并终止处理链。在这种情况下，HTTP 响应的状态码由 Error 对象的属性 statusCode 来确定，默认使用 500。调用 next.ifError(err) 并传入一个 Error 对象可以使得 restify 抛出异常并终止进程，可以用来在出现无法恢复的错误时终止程序。

在代码清单 4 中，通过 pre() 方法注册的处理器会记录请求的完整路径。第一个 use()方法注册了 restify 的插件 queryParser，其作用是把请求的查询字符串解析成 JavaScript 对象。第二个 use() 方法注册的处理器把请求的 HTTP 头 Accept 设置为 application/json，也就是 API 只接受 JSON 格式的请求。最后通过 get() 方法注册了两个对于 GET 请求的处理器，第一个设置了响应的额外 HTTP 头 X-Test，第二个设置响应的内容。当请求包含了查询参数 boom 时，服务器会直接返回 500 错误。

##### 清单 4\. 响应处理器链示例

```
const restify = require('restify');

const server = restify.createServer();

server.pre((req, res, next) => {
console.log('req: %s', req.href());
return next();
});

server.use(restify.plugins.queryParser());

server.use((req, res, next) => {
req.headers.accept = 'application/json';
return next();
});

server.get('/handler_chain', [(req, res, next) => {
res.header('X-Test', 'test');
return next();
}, (req, res, next) => {
if (req.query.boom) {
return next(new Error('boom!'));
}
res.send({
msg: 'handled!'
});
return next();
}]);

server.listen(8000, () => console.log('%s listening at %s', server.name, server.url));

```

Show moreShow more icon

## 路由

restify 的路由表示的是对 HTTP 请求的处理逻辑。一个路由有 3 个部分：分别是 HTTP 动词、匹配条件和处理方法。restify 服务器对象可以使用方法 get、put、post、del、head、patch 和 opts，分别与名称相同的 HTTP 请求动词相对应。这些方法用来创建路由。这些方法的第一个参数定义了路由的匹配条件。该参数的值可以是字符串或正则表达式 Regex 对象，也可以是包含了属性 name、path 和 version 的 JavaScript 对象。方法的第二个参数是进行实际处理的函数，接受 req、res 和 next 三个参数，分别表示 HTTP 请求、HTTP 响应和处理链的下一个处理器。

在代码清单 5 中，logHandler 是通用的路由处理方法，其作用是记录请求的路径和参数，并把解析之后的参数对象作为响应返回。代码清单 5 中一共定义了 4 个使用 get 方法的路由。第一个路由使用的是完全匹配的路径。第二个路由中包含了参数 id，在路径中以 “:id” 来表示。当访问路径 /route/user/1 时，返回的结果是 {“id”:”1″}。第三个路由使用正则表达式。只有在 /route/order/ 后全部为数字的路径才能满足匹配。当访问路径 /route/order/123 时，返回的结果为{“0″:”123”}。其中 0 表示是对应于正则表达式中的第一个匹配分组。当尝试访问路径 /route/order/xyz 时，服务器会返回 404 错误，因为 xyz 不匹配正则表达式。最后一个路由使用一个 JavaScript 对象声明了路径和版本号。关于版本号的使用，在下一节会提到。

##### 清单 5\. 路由示例

```
const restify = require('restify');

const server = restify.createServer();

const logHandler = (req, res, next) => {
console.log('req: %s, params: %s', req.href(), JSON.stringify(req.params));
res.send(req.params);
return next();
};

server.get('/route/simple', logHandler);
server.get('/route/user/:id', logHandler);
server.get(/^\/route\/order\/(\d+)/, logHandler);
server.get({
path: '/route/versioned',
version: '1.0.0'
}, logHandler);

server.listen(8000, () => console.log('%s listening at %s', server.name, server.url));

```

Show moreShow more icon

## 多版本路由

REST API 通常有同时运行多个版本的要求，以支持 API 的演化。restify 内置提供了基于语义化版本号（semver）规范的多版本的支持。在 HTTP 请求中可以使用 HTTP 头 Accept-Version 来指定版本号。每个路由可以按照代码清单 6 中的方式，在属性 version 中指定该路由的一个或多个版本号。如果请求中不包含 HTTP 头 Accept-Version，那么会匹配同一路由中版本最高的那一个。否则，就按照由 Accept-Version 指定的版本号来进行匹配，并调用匹配版本的路由。通过请求的 version() 方法可以获取到 Accept-Version 头的值，matchedVersion() 方法可以获取到匹配到的版本号。

在代码清单 6 中，同一个路由/versioned 有多个版本。版本 1.0.0 的处理方法返回字符串 V1。第二个路由同时支持 2.0.0、2.1.0 和 2.2.0 等 3 个版本，返回的是请求的版本和实际匹配的版本。如果直接访问/versioned，返回的结果是 `{"requestedVersion":"*","matchedVersion":"2.2.0"}`，因为默认匹配最高版本。如果运行 `"curl -s -H 'accept-version: ~1' http://localhost:8000/versioned"`，由于请求中 Accept-Version 头的值为 ~1，会匹配到 1.0.0 版本的路由，返回结果为 V1。如果请求中 Accept-Version 头的值为 ~3， 则返回结果为 `{"code":"InvalidVersion","message":"~3 is not supported by GET /versioned"}`，因为并没有与 ~3 匹配的版本的路由。

##### 清单 6\. 多版本路由示例

```
const restify = require('restify');

const server = restify.createServer();

server.get({
path: '/versioned',
version: '1.0.0',
}, (req, res, next) => {
res.send('V1');
return next();
});

server.get({
path: '/versioned',
version: ['2.0.0', '2.1.0', '2.2.0'],
}, (req, res, next) => {
res.send({
requestedVersion: req.version(),
matchedVersion: req.matchedVersion(),
});
return next();
});

server.listen(8000, () => console.log('%s listening at %s', server.name, server.url));

```

Show moreShow more icon

## WebSocket

在 restify 中也可以使用 WebSocket，不过需要第三方库的支持。本文使用 Socket.IO 来展示 WebSocket 的使用。代码清单 7 中给出了使用 Socket.IO 的服务器端实现。在创建了 restify 的服务器对象之后，可以直接把底层的 server 对象直接由 Socket.IO 来使用。路径为”/”的路由的作用是发送 HTML 页面。接下来的代码使用 Socket.IO 进行数据发送和接收。

##### 清单 7\. 使用 Socket.IO 的服务器端

```
const fs = require('fs');
const restify = require('restify');
const server = restify.createServer();
const io = require('socket.io')(server.server);

server.get('/', (req, res, next) => {
fs.readFile(__dirname + '/index.html', function (err, data) {
    if (err) {
      next(err);
      return;
    }

    res.setHeader('Content-Type', 'text/html');
    res.writeHead(200);
    res.end(data);
    next();
});
});

io.on('connection',(socket) => {
socket.emit('news', { hello: 'world' });
socket.on('my other event', console.log);
});

server.listen(8000, () => console.log('socket.io server listening at %s', server.url));

```

Show moreShow more icon

代码清单 8 中给出了 HTML 页面 index.html 的内容。

##### 清单 8\. 使用 Socket.IO 的客户端

```
<script src="/socket.io/socket.io.js"></script>
<script>
var socket = io.connect('http://localhost:8000');
socket.on('news', function (data) {
    console.log(data);
    socket.emit('my other event', { my: 'data' });
});
</script>

```

Show moreShow more icon

## 内容协商

在之前的示例中，我们都是使用 send() 方法来直接发送响应内容。如果传入的是 JavaScript 对象，restify 会自动转换成 JSON 格式。这是由于 restify 内置提供了对于不同响应内容类型的格式化实现。内置支持的响应内容类型包括 application/json、text/plain 和 application/octet-stream。restify 会根据请求的 Accept 头来确定响应的内容类型。如果无法确定，则默认使用 application/octet-stream。可以在创建 restify 服务器时，添加额外的响应内容类型的支持。

在代码清单 9 中，我们创建了一个对于内容类型 application/base64 的格式化实现。在该实现中，我们会把 String 类型的内容转换成 Base64 编码之后的格式。

##### 清单 9\. 内容协商示例

```
const util = require('util');
const restify = require('restify');

const server = restify.createServer({
formatters: {
'application/base64': (req, res, body) => {
     if (body instanceof Error) {
       return body.stack;
     }
     if (Buffer.isBuffer(body)) {
       return body.toString('base64');
     }
     if (typeof body === 'string') {
       return new Buffer(body).toString('base64');
     }
     return util.inspect(body);
}
}
});

server.get('/content', (req, res, next) => {
res.send('Hello World');
next();
});

server.listen(8000, () => console.log('%s listening at %s', server.name, server.url));

```

Show moreShow more icon

如果直接访问 /content，返回的结果是 SGVsbG8gV29ybGQ=，是 Base64 编码之后的结果。

这是因为默认的内容类型变成了 application/base64。如果使用 `"curl -s -H 'accept: text/plain' localhost:8000/content"` 来访问，则返回的结果是 Hello World，内容类型变为纯文本。如果使用 `"curl -s -H 'accept: application/json' localhost:8000/content"` 来访问，则返回的结果是 “Hello World”。

## 错误处理

在 REST API 的实现中，错误处理是很重要的一部分。在前面的示例中，我们使用 send() 方法发送 Error 对象来使得 restify 产生错误响应。由于 HTTP 的状态码是标准的，restify 提供了一个专门的模块 restify-errors 来创建对应不同状态码的 Error 对象，可以直接在 send() 方法中使用。restify 中产生的错误会作为事件来发送，可以使用 Node.js 标准的事件处理机制来进行处理。需要注意的是，只有使用 next() 方法发送的错误会被作为事件来发送，使用响应对象的 send() 方法发送的则不会。

在代码清单 10 中，路由 /error/500 使用 send() 发送了一个 InternalServerError 错误对象，而/error/400 和/error/404 使用 next() 分别发送了 BadRequestError 和 NotFoundError 错误对象。可以使用 server.on() 来添加对于 NotFound 错误的处理逻辑，但是对于 InternalServer 错误的处理逻辑不会被触发，因为该错误是通过 send() 方法来发送的。

##### 清单 10\. 错误处理示例

```
const restify = require('restify');
const errors = require('restify-errors');

const server = restify.createServer();

server.get('/error/500', (req, res, next) => {
res.send(new errors.InternalServerError('boom!'));
return next();
});
server.get('/error/400', (req, res, next) => next(new errors.BadRequestError('bad request')));
server.get('/error/404', (req, res, next) => next(new errors.NotFoundError('not found')));

server.on('NotFound', (req, res, err, cb) => {
console.error('404 %s', req.href());
return cb();
});

server.on('InternalServer', (req, res, err, cb) => {
console.error('should not appear');
return cb();
});

server.listen(8000, () => console.log('%s listening at %s', server.name, server.url));

```

Show moreShow more icon

## 插件

在 REST API 的开发中，某些任务是比较常见的。restify 提供了一系列插件来满足这些通用的需求。这些插件可以通过 restify.plugins 来访问，并使用 use() 方法来注册。在代码清单 11 中，我们使用了 restify 中的若干个常用插件。acceptParser 用来解析请求的 Accept 头，以确保是服务器端可以处理的类型。如果是服务器端不支持的类型，该插件会返回 406 错误。authorizationParser 用来解析请求中的 Authorization 头，并把解析的结果保存在请求对象的属性 authorization 中。queryParser 之前已经介绍过，用来解析请求中的查询字符串。gzipResponse 用来发送 GZIP 压缩之后的响应。bodyParser 用来解析请求的内容，并把结果保存在请求对象的属性 body 中。目前支持的请求内容类型包括 application/json、application/x-www-form-urlencoded 和 multipart/form-data。

##### 清单 11\. 插件示例

```
const restify = require('restify');

const server = restify.createServer();
server.use(restify.plugins.acceptParser(server.acceptable));
server.use(restify.plugins.authorizationParser());
server.use(restify.plugins.queryParser());
server.use(restify.plugins.gzipResponse());
server.use(restify.plugins.bodyParser());

server.post('/plugins', (req, res, next) => {
console.log(req.body);
res.send({a: 1});
return next();
});

server.listen(8000, () => console.log('%s listening at %s', server.name, server.url));

```

Show moreShow more icon

## 结束语

当需要在 Node.js 上开发 REST API 时，restify 是一个很好的选择。restify 的响应处理链使得处理 HTTP 请求变得非常简单。可以使用不同的方式来定义路由，也提供了对多版本的支持。本文对 restify 进行了详细介绍，包括 WebSocket 支持，内容协商、错误处理和插件等。