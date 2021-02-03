# 使用 WebSocket 和 SSE 实现 HTTP 服务器推送
Web 和移动应用程序的实时数据传输

**标签:** Web 开发,移动开发

[原文链接](https://developer.ibm.com/zh/articles/wa-http-server-push-with-websocket-sse/)

Daniel Barboza

发布: 2017-06-05

* * *

HTTP 服务器推送也称为 HTTP 流，是一种客户端-服务器通信模式，它将信息从 HTTP 服务器异步发送到客户端，而无需客户端请求。在高度交互的 Web 或移动应用程序中，一个或多个客户端需要连续不断地从服务器接收信息，服务器推送架构对这类应用程序特别有效。在本文中，您将了解 WebSocket 和 SSE（服务器发送的事件），它们是实现 HTTP 服务器推送的两种技术。

我首先将概述两种解决方案之间的技术差异，以及如何在 Web 架构中实现每种解决方案。我将通过一个示例应用程序，向您展示如何设置一个 SSE 实现，然后介绍一个类似的 WebSocket 实现。最后，我将比较两种技术，并提出我关于在不同类型 Web 或移动应用程序中使用它们的结论。

请注意，本文要求熟悉 HTTP 服务器推送的语言和概念。两个应用程序都是在 Python 中使用 CherryPy 编写的。

## 请求-响应的局限性

网络上的客户端-服务器通信在过去曾是一种请求-响应模型，要求客户端（比如 Web 浏览器）向服务器请求资源。服务器通过发送所请求的资源来响应客户端请求。如果资源不可用，或者客户端没有权限访问它，那么服务器会发送一条错误消息。在请求-响应架构中，服务器绝不能向客户端发送未经请求的消息。

随着 Web 应用程序变得更强大和更具交互性，请求-响应模型的局限性也开始显现出来。需要更频繁更新的客户端应用程序被要求更频繁地发送 `GET` 请求。这种技术称为 _轮询_ ，在高峰期间，这可能会使服务器不堪重负，并导致性能问题。该技术效率低下，因为客户端发送的许多请求都没有返回更新。此外，客户端只能按指定间隔进行轮询，这可能减缓客户端的响应速度。

HTTP 服务器推送技术的出现，就是为了解决与频繁轮询相关的性能问题和其他局限。尤其是对于交互式 Web 应用程序，比如游戏和屏幕共享服务，Web 服务器能更高效地在新数据可用时向客户端发送更新。

## 比较 WebSocket 与 SSE

WebSocket 和 SSE 都是传统请求-响应 Web 架构的替代方案，但它们不是完全冲突的技术。WebSocket 架构在客户端与服务器之间打开一个套接字，用于实现全双工（双向）通信。无需发送 `GET` 消息并等待服务器响应，客户端只需监听该套接字，接收服务器更新，并使用收到的数据来发起或支持各种交互。客户端也可以使用套接字与服务器通信，例如在成功收到更新时发送 `ACK` 消息。

SSE 是一种更简单的标准，是作为 HTML5 的扩展而开发的。尽管 SSE 支持从服务器向客户端发送异步消息，但客户端无法向服务器发送消息。对于客户端只需接收从服务器传入的更新的应用程序，SSE 的半双工通信模型最适合。与 WebSocket 相比，SSE 的一个优势是它是基于 HTTP 而运行的，不需要其他组件。

对于需要在客户端与服务器之间频繁通信的多用途 Web 应用程序，显然应该选择 WebSocket。对于希望从服务器向客户端传输异步数据，而不需要响应的应用程序，SSE 更适合一些。

### 浏览器支持

在比较 HTTP 协议时，浏览器支持是一个必要的考虑因素。浏览表 1 中的数据，可以看到所有现代浏览器都支持 [WebSocket](http://caniuse.com/#feat=WebSockets) 协议，包括移动浏览器。Microsoft IE 和 Edge 不支持 [SSE](http://caniuse.com/#feat=eventsource) 。

##### 表 1\. 2017 年 1 月浏览器使用分布

2017ChromeIE/EdgeFirefoxSafariOpera**1 月**73.7%4.9%15.4%3.6%1.0%

对于必须在所有浏览器中运行的应用程序，WebSocket 目前是更好的选择。

### 开发工作量

在比较协议时，尤其是 WebSocket 和 SSE 等较新的协议，工作量是要考虑的另一个因素。这里的 “工作量” 指的是代码行，或者您要花多少时间为使用给定协议的应用程序编写代码。对于具有严格的时间限制或开发预算的项目，这个指标特别重要。

实现 SSE 的工作量比 WebSocket 要少得多。它适合任何使用 HTML5 编写的应用程序，主要负责在从服务器发送到客户端的消息中添加一个 HTTP 标头。如果给定了正确标头，客户端就会自动将消息识别为服务器发送的事件。不同于 WebSocket，SSE 不需要在服务器与客户端之间建立或维护套接字连接。

WebSocket 协议要求配置一个服务器端套接字来监听客户端连接。客户端自动打开一个与服务器的套接字并等待消息，消息可以异步发送。每个应用程序都可以定义自己的消息格式、持久连接（脉动信号）策略等。

表 2 总结了 SSE 和 WebSocket 协议的优劣。

##### 表 2\. 比较 SSE 与 WebSocket

SSEWebSocket**通信类型**半双工（单向）全双工（双向）**浏览器支持**目前在 Microsoft 浏览器中不可用。可用于所有主要浏览器。**开发工作量**小：只需发送一条包含特定标头的 HTTP 消息。中等：需要建立并维护 TCP 套接字通信。在服务器端还需要一个监听器套接字。

现在让我们通过一个简单的 Web 应用程序示例来了解每种技术的工作原理。

## 开发 SSE 应用程序

SSE 是一种仅使用 HTTP 传送异步消息的 HTML5 标准。不同于 WebSocket，SSE 不需要在后端创建服务器套接字，也不需要在前端打开连接。这大大降低了复杂性。

### SSE 前端

清单 1 给出了一个使用 SSE 的简单 HTTP 服务器推送应用程序的 UI 代码：

##### 清单 1\. SSE 前端

```
var source = new EventSource("/user-log-stream”);
source.onmessage = function(event) {
var message = event.data;
// do stuff based on received message
};

```

Show moreShow more icon

`EventSource` 是一个与服务器建立 HTTP 连接的接口。服务器使用 _事件流格式_ 将消息发送到客户端，这种格式是一种使用 UTF-8 编码的简单的文本数据流。建立联系后，HTTP 连接会对给定消息流保持开放，以便服务器能发送更新。在清单 1 中，我们创建了一个 HTTP 连接来接收与 `/#tabs/user-log-stream` URI 相关的事件。

### SSE 后端

在后端，我们为 URL `/user-log-stream` 创建了一个分派器。该分派器将从前端接收连接请求，并发回一条异步消息来发起通信。SSE 客户端不能向服务器发送消息。

清单 2 中的代码演示了后端代码。

##### 清单 2\. SSE 后端

```
import cherrypy

class UserLogStream
    messages = []

    @cherrypy.expose
    def stream(self):
        cherrypy.response.headers["Content-Type"] = "text/event-stream"
        while True:
            if len(messages) > 0:
                for msg in messages:
                    data = "data:” + msg + "\n\n”
                    yield data
                 messages = []

routes_dispatcher = cherrypy.dispatch.RoutesDispatcher()
routes_dispatcher.connect('user-log-stream', '/', controller = UserLogStream(), action='stream')

```

Show moreShow more icon

`UserLogStream` 类的 `stream` 方法被分配给每个与 `/user-log-stream` URI 连接的 `EventSource` 。任何待发送消息都将被发送到已连接的 `EventSource` 。请注意，消息格式不是可选的：SSE 协议要求消息以 `data:` 开头，以 `\n\n` 结尾。尽管这些示例中使用了 HTTP 作为消息格式，但也可以使用 JSON 或另一种格式发送消息。

尽管这是一个非常基本的 SSE 实现示例，但它演示了该协议的简单性。

## 开发 WebSocket 应用程序

像 SSE 示例一样，WebSocket 应用程序基于 CherryPy（一种 Python Web 框架）而构建。 [Project WoK](https://github.com/kimchi-project/wok) 是后端 Web 服务器， `python-websockify` 库插件处理套接字连接。该插件是 [Kimchi](https://github.com/kimchi-project/kimchi) 的一部分。

组件包括：

- Project WoK：后端 Python 逻辑，提供要由推送服务器广播的消息。
- Websockify 代理：此代理由 `python-websockify` 库实现，它使用 Unix 套接字从推送服务器接收消息，并将这些消息传送到 UI 上连接的 WebSocket。
- 推送服务器：一个普通 Unix 套接字服务器，它向已连接的客户端广播后端消息。
- 前端：UI 与 Websockify 代理建立 WebSocket 连接，并监听服务器消息。根据具体的消息，将会执行一个特定操作，例如刷新一个列表或向用户显示一条消息。

### Websockify 代理

Websockify 代理由 `python-websockify` 库实现。它使用 CherryPy Web 服务器引擎来实现初始化，如清单 3 所示。

##### 清单 3\. Websockify 代理

```
params = {'listen_host': '127.0.0.1',
          'listen_port': config.get('server', 'websockets_port'),
          'ssl_only': False}

# old websockify: do not use TokenFile
if not tokenFile:
    params['target_cfg'] = WS_TOKENS_DIR

# websockify 0.7 and higher: use TokenFile
else:
    params['token_plugin'] = TokenFile(src=WS_TOKENS_DIR)

def start_proxy():
    try:
        server = WebSocketProxy(RequestHandlerClass=CustomHandler,
                                **params)
    except TypeError:
        server = CustomHandler(**params)

    server.start_server()

proc = Process(target=start_proxy)
proc.start()
return proc

```

Show moreShow more icon

Websockify 代理将所有已注册的令牌都绑定到它的 WebSocket URI，如下所示：

##### 清单 4\. 绑定已注册的令牌

```
def add_proxy_token(name, port, is_unix_socket=False):
    with open(os.path.join(WS_TOKENS_DIR, name), 'w') as f:
        """
        From python documentation base64.urlsafe_b64encode(s)
        substitutes - instead of + and _ instead of / in the
        standard Base64 alphabet, BUT the result can still
        contain = which is not safe in a URL query component.
        So remove it when needed as base64 can work well without it.
        """
        name = base64.urlsafe_b64encode(name).rstrip('=')
        if is_unix_socket:
            f.write('%s: unix_socket:%s' % (name.encode('utf-8'), port))
        else:
            f.write('%s: localhost:%s' % (name.encode('utf-8'), port))

```

Show moreShow more icon

下面是添加一个名为 `myUnixSocket` 的 Websockify 代理条目的命令。此条目通过代理将 WebSocket 连接传送到 Unix 套接字 `/run/my_socket` ：

```
add_proxy_token('myUnixSocket', '/run/my_socket', True)

```

Show moreShow more icon

在 UI 中，我们使用以下 URI 打开一个与 Unix 套接字的 WebSocket 连接：

```
wss://<server_address>:<port>/websockify?token=<b64encodedtoken>

```

Show moreShow more icon

在 URI 中， `b64encodedtoken` 是 `myUnixSocket` 字符串的 base64 值，不包括 `=` 字符。有关此配置的更多信息，请参阅 [WebSocket 代理模块](https://github.com/kimchi-project/wok/blob/master/src/wok/websocket.py) 。

### 推送服务器

推送服务器有两个要求：

1. 它必须能同时处理多个连接。
2. 它必须能向所有已连接的客户端广播同一条消息。

出于此示例的用途，推送服务器将充当一个广播代理。尽管能够接收来自客户端的消息，但我们不希望这样做。

清单 5 给出了推送服务器的初始版的相关 Python 代码。（可访问 [GitHub 上的 Project WoK](https://github.com/kimchi-project/wok/blob/master/src/wok/pushserver.py) 来获取 `pushserver` 模块的最终版本。）

##### 清单 5\. WebSocket 推送服务器

```
class PushServer(object):

    def __init__(self):
        self.set_socket_file()

        websocket.add_proxy_token(TOKEN_NAME, self.server_addr, True)

        self.connections = []

        self.server_running = True
        self.server_socket = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        self.server_socket.setsockopt(socket.SOL_SOCKET,
                                      socket.SO_REUSEADDR, 1)
        self.server_socket.bind(self.server_addr)
        self.server_socket.listen(10)
        wok_log.info('Push server created on address %s' % self.server_addr)

        self.connections.append(self.server_socket)
        cherrypy.engine.subscribe('stop', self.close_server, 1)

        server_loop = threading.Thread(target=self.listen)
        server_loop.setDaemon(True)
        server_loop.start()

    def listen(self):
        try:
            while self.server_running:
                read_ready, _, _ = select.select(self.connections,
                                                 [], [], 1)
                for sock in read_ready:
                    if not self.server_running:
                        break

                    if sock == self.server_socket:

                        new_socket, addr = self.server_socket.accept()
                        self.connections.append(new_socket)
                    else:
                        try:
                            data = sock.recv(4096)
                        except:
                            try:
                                self.connections.remove(sock)
                            except ValueError:
                                pass

        except Exception as e:
            raise RuntimeError('Exception occurred in listen() of pushserver '
                               'module: %s' % e.message)

    def send_notification(self, message):
        for sock in self.connections:
            if sock != self.server_socket:
                    sock.send(message)

```

Show moreShow more icon

`PushServer()` 类的 `__init__()` 将会初始化 Unix 套接字，以便监听连接。我们在 WebSocket 代理中也添加了一个令牌。此令牌将被前端用在 WebSocket 连接的 URI 中。

可通过许多方法实现服务器推送架构。出于广播服务器推送实现的用途，我认为使用一个监听器线程 (`select.select()`) 来控制已知连接是最简单的方法。在这里， `listen()` 方法在一个守护线程中运行，负责管理现有连接。 `select.select` 方法以非阻塞函数的形式运行，它在来自 `self.connections` 数组的所有套接字都准备好读取时返回，或者在一秒超时后返回。建立新连接后，它会被添加到 `self.connections` 数组中。连接关闭时，会从数组中删除它。

`listen()` 方法可通过两种方法从数组中删除套接字：在 `recv()` 调用的 `except` 代码块中检测到套接字已关闭，或者收到一条 `CLOSE` 消息。 `send_notification()` 方法负责将消息广播到订阅了它们的客户端。在客户端取消订阅时，我们使用了 `except` 代码块关闭套接字连接。不能单独使用 `listen()` 关闭客户端与服务器之间的套接字连接，因为这样做可能导致 `send_notification()` 方法中发生管道损坏错误。

下一节将更详细地解释这一点。

## 排除 WebSocket 前端的故障

后端设置和配置相对比较容易，但编写前端代码就没有那么简单了。从决定打开还是关闭浏览器连接到处理某些 API 的方式上的差异，WebSocket 前端带来了许多挑战。在解决这些挑战时，需要重新利用后端代码的某些元素。

第一个决定是为所有 UI 维护单个 WebSocket 连接，还是设置多个按需打开的 WebSocket 连接。我们可能会先考虑多个按需 WebSocket 连接。我们认为，仅在 UI 期望从后端收到一条特定消息时才应打开一个连接 — 例如等待某个任务完成。如果不想收到异步消息，则不需要打开连接。这会从 UI 和后端节约一些资源，但代价是需要管理每个已打开的 WebSocket 的生命周期。

对于期望从服务器收到大量异步消息的 UI，单个持久 WebSocket 连接可能更合适。这个持久连接被用作通知流，允许任何相关 UI 代码监听收到的消息并执行相关操作。

每个解决方案都各有用武之地，下面将试验每个解决方案。使用最适合您的应用程序的方法很重要。

### 多个 WebSocket 连接

要打开和关闭 WebSocket，只需调用一个构造函数来打开和调用 `close()` 方法。但是，诸如页面刷新之类的浏览器行为上的差异可能让情况变得很复杂。每次打开 WebSocket 都会与 WebSocket 代理建立一个连接，然后与后端的实际推送服务器建立连接。该连接在关闭之前会保持打开。考虑以下代码：

##### 清单 6\. 在每条收到的消息上调用 refreshInterface()

```
notificationsWebSocket = new WebSocket(url, ['base64']);
notificationsWebSocket.onmessage = function(event) {
    refreshInterface();
};

```

Show moreShow more icon

在清单 6 中， `onmessage` 是一个事件处理函数，每次在套接字中收到一条新消息时都会调用它。 `refreshInterface()` 函数在消息到达时被调用。如果打开多个 `notificationsWebSocket` 对象，则会调用 `refreshInterface()` 多次。

对前端和后端而言，打开比我们需要的更多的 WebSocket 连接显然是一种资源浪费。这还会带来意想不到的风险。为了解决这些问题，我们可以避免打开不必要的 WebSocket 连接，并在不再需要现有连接时关闭它们。可以手动或通过浏览器关闭 WebSocket。在我用于测试的浏览器中，在浏览器窗口刷新或关闭时，或者在域更改时，Chrome 和 Firefox 都会关闭已打开的 WebSocket 连接。

WebSocket 关闭时，会发生两件事：

1. 调用一个名为 `onclose` 的事件处理函数来执行任何附加清理。
2. 闭合端（在本例中为 UI WebSocket）启动一次关闭握手来提醒另一端。在此过程中，会发送一个值为 [0x8](https://tools.ietf.org/html/rfc6455#page-9) 的 _关闭帧（ close frame）_ 。收到此帧时，后端会停止发送数据并立即关闭它所在的套接字端。无论 `onclose` 处理函数中配置了何种清理类型，此过程都会发生。

理想情况下，会调用 `onclose` 处理函数，而且推送服务器会收到关闭帧。遗憾的是，我在测试此阶段的实现时发生的过程不是这样的。

#### Websockify 代理和关闭帧

截至编写本文时，Websockify 代理插件 (v0.8.0-2) 不会将关闭帧转发到 Unix 套接字，继而转发到推送服务器。因此，UI 套接字中的连接可能已关闭，但推送服务器毫不知情，仍使它所在的连接端保持打开。

推送服务器尝试将一条新消息广播到所有连接（包括已关闭的套接字）时，它会收到管道已损坏的错误。如果未正确处理，这可能导致推送服务器被终止运行。下面这个套接字发出了一条 `send_notification` 消息：

##### 清单 7\. send\_notification 导致管道损坏错误

```
def send_notification(self, message):
    for sock in self.connections:
        if sock != self.server_socket:
                sock.send(message)

```

Show moreShow more icon

处理这种情况的一个方法是将 `sock.send()` 调用包装在一个 `try/except` 中，然后处理管道损坏异常。一种替代方案是建立一次独特的关闭握手，使用 `onclose` 事件处理函数向后端服务器发送一条消息。推送服务器端然后关闭它自己的套接字连接，保持 `send_notification` 方法不变。

接下来我们将尝试该解决方案。

#### 使用 onclose 发送关闭消息

`onclose` 事件处理函数会在终止 WebSocket 连接之前被调用。但此刻连接仍然有效。我们的想法是允许干净地关闭连接，方法是向推送服务器发送一条关闭消息，让它知道套接字将在前端关闭。基本上讲，我们模拟的是关闭帧在关闭握手中的角色。

在清单 8 中，WebSocket 在关闭连接前向推送服务器发送一条 `CLOSE` 消息：

##### 清单 8\. 在 onclose 事件中发送 CLOSE 消息

```
notificationsWebSocket = new WebSocket(url, ['base64']);
(...)
notificationsWebSocket.onclose = function() {
    notificationsWebSocket.send(window.btoa('CLOSE'));
};

```

Show moreShow more icon

在推送服务器中，我们监听从 UI 传入的所有消息，查看是否收到了 `CLOSE` 消息。如果收到该消息，则立即关闭套接字：

##### 清单 9\. 推送服务器对 CLOSE 的响应

```
def listen(self):
    try:
        while self.server_running:
            read_ready, _, _ = select.select(self.connections,
                                             [], [], 1)
            for sock in read_ready:
                if not self.server_running:
                    break

                if sock == self.server_socket:

                    new_socket, addr = self.server_socket.accept()
                    self.connections.append(new_socket)
                else:
                    try:
                        data = sock.recv(4096)
                    except:
                        try:
                            self.connections.remove(sock)
                        except ValueError:
                            pass

                        continue
                    if data and data == 'CLOSE':
                        try:
                            self.connections.remove(sock)
                        except ValueError:
                            pass
                        sock.close()

```

Show moreShow more icon

在 WebSocket 关闭和发送 `CLOSE` 消息之前调用 `onclose` ，这样就能预防 `send_notification` 例程中发生管道损坏错误。不幸的是，尽管 `onclose` 在 Firefox 中运行良好，但测试表明，我们在 Chrome 中打开该应用程序时，仍在 `send_notification` 方法中获得了管道损坏错误。这是由于 Google Chrome 中的一个未解决的 [问题](https://bugs.chromium.org/p/chromium/issues/detail?id=93698) 。

### 单一、持久的 WebSocket 连接

`onclose` 和关闭帧的问题促使我们采用单一、持久的 WebSocket 连接。对于仅在浏览器离开页面（通过关闭窗口、重新加载或转到另一个 URL）时才会关闭的单一连接，只要求推送服务器在发生管道损坏时关闭连接。一个集中化连接更容易在 Project WoK 中维护，而且 WoK 插件使我们无需从头实现 WebSocket 策略。

清单 10 给出了一个经过修改的 `send_notification` 消息。请注意，它处理管道损坏错误的方式是，从有效连接列表中删除套接字，然后关闭它：

##### 清单 10\. 针对单个持久连接的 send\_notification

```
def send_notification(self, message):
for sock in self.connections:
    if sock != self.server_socket:
        try:
            sock.send(message)
        except IOError as e:
            if 'Broken pipe' in str(e):
                sock.close()
                try:
                    self.connections.remove(sock)
                except ValueError:
                    pass

```

Show moreShow more icon

#### 配置脉动信号消息

在多连接策略中，您会打开一个连接，使用它，然后立即关闭它。单个持久连接应无限期地保持打开，即使没有发送任何消息。但是，Chrome 和 Firefox 中的测试表明，WebSocket 连接会在空闲一段时间后超时，而且连接会终止。

为了避免浏览器超时，我们实现了一条简单的脉动信号消息，定期向推送服务器发送该消息。此消息使 WebSocket 连接保持活动 — 推送服务器不需要响应它。默认超时大约为 300 秒，但我们不能假设所有浏览器都实现相同的超时。为安全起见，我们将脉动信号消息发送间隔设置为 30 秒，如清单 11 所示：

##### 清单 11\. 超时间隔

```
notificationsWebSocket = new WebSocket(url, ['base64']);
var heartbeat = setInterval(function() {
    notificationsWebSocket.send(window.btoa('heartbeat'));
}, 30000);

```

Show moreShow more icon

这个简单的超时会使 WebSocket 连接保持活动，而且不会向推送服务器发送过多脉动信号消息。

#### 添加监听器

建立永久连接后，可以实现一个简单的监听器策略，允许示例应用程序中的其他所有 UI 代码监听消息。下面是该策略的设计：

1. 一个监听器绑定到一种特定消息格式。我们的示例应用程序将发送许多通知消息。消息绑定使我们能仅为特定消息注册监听器。
2. 收到来自推送服务器的消息后，主要 WebSocket 通道验证它的内容，并调用与该消息绑定的所有监听器。
3. WebSocket 通道发起监听器清理。监听器必须在 URI 更改后弃用。例如，在用户浏览到 URI `/#tabs/settings` 中的 **Settings** 选项卡时，应弃用 URI `/#tabs/user-log` 中的 **Activity Log** 选项卡上的侦听器。清理操作被绑定到 `hashchange` 事件， URL (`location.hash`) 每次更改时都会触发该事件。因为我们只需对此过程执行一次，所以可以将 `$.one()' jQuery` 调用设置为仅执行一次。

以下是示例应用程序中处理监听器的代码：

##### 清单 12\. Project WoK 的监听器配置

```
wok.notificationListeners = {};
wok.addNotificationListener = function(msg, func) {
    var listenerArray = wok.notificationListeners[msg];
    if (listenerArray == undefined) {
        listenerArray = [];
    }
    listenerArray.push(func);
    wok.notificationListeners[msg] = listenerArray;
    $(window).one("hashchange", function() {
        var listenerArray = wok.notificationListeners[msg];
        var del_index = listenerArray.indexOf(func);
        listenerArray.splice(del_index, 1);
        wok.notificationListeners[msg] = listenerArray;
    });
};

```

Show moreShow more icon

以下是将消息发送到每个相关监听器的方式：

##### 清单 13\. 监听器通知

```
wok.startNotificationWebSocket = function () {
    wok.notificationsWebSocket = new WebSocket(url, ['base64']);
    wok.notificationsWebSocket.onmessage = function(event) {
        var message  = window.atob(event.data);
         if (message === "") {
                continue;
        }

        var listenerArray = wok.notificationListeners[message];
        if (listenerArray == undefined) {
            continue;
        }
        for (var i = 0; i < listenerArray.length; i++) {
            listenerArray[i](message);
        }
    }
};

```

Show moreShow more icon

有了该通知，我们的操作就基本上完成了。但首先还需要解决一个问题。

### 合并消息

尽管可以使用单一连接，但当前实现对前端如何接收消息做出了不合理的假设。如果推送服务器接连发送了两条消息，先发送 `message1` ，然后发送 `message2` ，结果会怎样？代码要求触发 `onmessage` 事件两次并从 WebSocket 中读取内容，但事实未必如此。事实上，更可能的情况是，收到的消息是 `message1message2` — 一条 _合并消息_ 。

合并消息是一种常见的 TCP 套接字行为：当接连发送两条或更多消息时，接收套接字会对它们进行排队，而且仅为所有消息触发 `onmessage` 事件一次。如果我们未对此制定应对计划，此行为会破坏我们的代码。解决方案是定义一种消息格式，让应用程序确定一条消息在何处结束，另一条消息在何处开始。

一种流行的选择是 [TLV 格式](https://en.wikipedia.org/wiki/Type-length-value) ，表示类型-长度-值 (Type-Length-Value)。在此格式中，消息依次包含 3 个字段：type、length 和 value。类型和长度具有固定大小，值具有 _length_ 字段所声明的可变大小。 _type_ 字段可用于区分消息类型，比如字符串、布尔值、二进制或整数。它也可用于特定于应用程序的用途，比如警告消息、信息消息、错误消息等。

但是，对我们而言，TLV 有点大材小用。我们的消息始终是字符串，所以不需要 _type_ 字段。为了解决消息的串联问题，可以添加一个固定长度的字段，用于声明消息包含多少个 `字符` 。一种更简单的方法是，向发送的每条消息添加一个消息结束标记。收到消息缓存后，前端使用消息标记解析它，将它拆分为单独的前端消息。

清单 14 给出了修改后的 `send_notification` 方法，其中包含一个消息结束标记：

##### 清单 14\. send\_notification 中的消息结束标记

```
END_OF_MESSAGE_MARKER = '//EOM//'

    def send_notification(self, message):
        message += END_OF_MESSAGE_MARKER
        for sock in self.connections:
            if sock != self.server_socket:
                try:
                    sock.send(message)
                except IOError as e:
                    if 'Broken pipe' in str(e):
                        sock.close()
                        try:
                            self.connections.remove(sock)
                        except ValueError:
                            pass

```

Show moreShow more icon

请注意清单 14 中的 `//EOM//` 序列。可以使用任何序列来表示消息结束，只要它不是包含在有效消息中的内容 — 比如 _end_ 。

清单 15 给出了在添加消息结尾解析后的已完成的前端代码。请参阅 [Project WoK](https://github.com/kimchi-project/wok/blob/master/ui/js/src/wok.main.js) 来获取完整源代码

##### 清单 15\. WebSocket 前端

```
wok.notificationListeners = {};
wok.addNotificationListener = function(msg, func) {
    var listenerArray = wok.notificationListeners[msg];
    if (listenerArray == undefined) {
        listenerArray = [];
    }
    listenerArray.push(func);
    wok.notificationListeners[msg] = listenerArray;
    $(window).one("hashchange", function() {
        var listenerArray = wok.notificationListeners[msg];
        var del_index = listenerArray.indexOf(func);
        listenerArray.splice(del_index, 1);
        wok.notificationListeners[msg] = listenerArray;
    });
};

wok.notificationsWebSocket = undefined;
wok.startNotificationWebSocket = function () {
    var addr = window.location.hostname + ':' + window.location.port;
    var token = wok.urlSafeB64Encode('woknotifications').replace(/=*$/g, "");
    var url = 'wss://' + addr + '/websockify?token=' + token;
    wok.notificationsWebSocket = new WebSocket(url, ['base64']);

    wok.notificationsWebSocket.onmessage = function(event) {
        var buffer_rcv = window.atob(event.data);
        var messages = buffer_rcv.split("//EOM//");
        for (var i = 0;  i < messages.length; i++) {
            if (messages[i] === "") {
                continue;
            }
            var listenerArray = wok.notificationListeners[messages[i]];
            if (listenerArray == undefined) {
                continue;
            }
            for (var j = 0; j < listenerArray.length; j++) {
                listenerArray[j](messages[i]);
            }
        }
    };

    var heartbeat = setInterval(function() {
        wok.notificationsWebSocket.send(window.btoa('heartbeat'));
    }, 30000);

};

```

Show moreShow more icon

## 结束语

对于只需要能向客户端传输异步服务器消息的 Web 应用程序，服务器发送的事件是一个优雅且简单的解决方案。作为半双工 HTTP 解决方案，SSE 不允许客户端向服务器回传消息。此外，在编写本文时，所有 Microsoft 浏览器都不支持 SSE。此限制是否是致命弱点，取决于应用程序的目标受众。

WebSocket 更复杂且要求更高，但全双工 TCP 连接使它适用于更广泛的应用场景。WebSocket 受大多数现代 Web 框架支持，而且兼容所有主要的 Web 和移动浏览器。尽管没有演示，但可以使用类似 Tornado 这样的服务器框架，它有助于快速配置推送服务器，而无需从头编写服务器代码。

SSE 是一个更简单且更快的解决方案，但它是不可扩展的。如果 Web 应用程序的需求发生了改变（例如，如果您认为前端应与后端交互），则需要使用 WebSocket 重构应用程序。WebSocket 需要做更多的前期工作，但它是一种更灵活的、可扩展的框架。它更适合会不断添加新功能的复杂应用程序。

本文翻译自： [HTTP server push with WebSocket and SSE](https://developer.ibm.com/articles/wa-http-server-push-with-websocket-sse/)（2017-05-05）