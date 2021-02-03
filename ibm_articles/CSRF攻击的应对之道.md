# CSRF 攻击的应对之道
如何在网站中防御 CSRF 的攻击

**标签:** Web 开发

[原文链接](https://developer.ibm.com/zh/articles/1102-niugang-csrf/)

牛刚, 童强国

发布: 2011-02-24

* * *

## CSRF 背景与介绍

CSRF（Cross Site Request Forgery，跨站域请求伪造）是一种网络的攻击方式，它在 2007 年曾被列为互联网 20 大安全隐患之一。其他安全隐患，比如 SQL 脚本注入，跨站域脚本攻击等在近年来已经逐渐为众人熟知，很多网站也都针对他们进行了防御。然而，对于大多数人来说，CSRF 却依然是一个陌生的概念。即便是大名鼎鼎的 Gmail, 在 2007 年底也存在着 CSRF 漏洞，从而被黑客攻击而使 Gmail 的用户造成巨大的损失。

### CSRF 攻击实例

CSRF 攻击可以在受害者毫不知情的情况下以受害者名义伪造请求发送给受攻击站点，从而在并未授权的情况下执行在权限保护之下的操作。比如说，受害者 Bob 在银行有一笔存款，通过对银行的网站发送请求 `http://bank.example/withdraw?account=bob&amount=1000000&for=bob2` 可以使 Bob 把 1000000 的存款转到 bob2 的账号下。通常情况下，该请求发送到网站后，服务器会先验证该请求是否来自一个合法的 session，并且该 session 的用户 Bob 已经成功登陆。黑客 Mallory 自己在该银行也有账户，他知道上文中的 URL 可以把钱进行转帐操作。Mallory 可以自己发送一个请求给银行：`http://bank.example/withdraw?account=bob&amount=1000000&for=Mallory`。但是这个请求来自 Mallory 而非 Bob，他不能通过安全认证，因此该请求不会起作用。这时，Mallory 想到使用 CSRF 的攻击方式，他先自己做一个网站，在网站中放入如下代码：`src=”http://bank.example/withdraw?account=bob&amount=1000000&for=Mallory ”`，并且通过广告等诱使 Bob 来访问他的网站。当 Bob 访问该网站时，上述 url 就会从 Bob 的浏览器发向银行，而这个请求会附带 Bob 浏览器中的 cookie 一起发向银行服务器。大多数情况下，该请求会失败，因为他要求 Bob 的认证信息。但是，如果 Bob 当时恰巧刚访问他的银行后不久，他的浏览器与银行网站之间的 session 尚未过期，浏览器的 cookie 之中含有 Bob 的认证信息。这时，悲剧发生了，这个 url 请求就会得到响应，钱将从 Bob 的账号转移到 Mallory 的账号，而 Bob 当时毫不知情。等以后 Bob 发现账户钱少了，即使他去银行查询日志，他也只能发现确实有一个来自于他本人的合法请求转移了资金，没有任何被攻击的痕迹。而 Mallory 则可以拿到钱后逍遥法外。

### CSRF 攻击的对象

在讨论如何抵御 CSRF 之前，先要明确 CSRF 攻击的对象，也就是要保护的对象。从以上的例子可知，CSRF 攻击是黑客借助受害者的 cookie 骗取服务器的信任，但是黑客并不能拿到 cookie，也看不到 cookie 的内容。另外，对于服务器返回的结果，由于浏览器同源策略的限制，黑客也无法进行解析。因此，黑客无法从返回的结果中得到任何东西，他所能做的就是给服务器发送请求，以执行请求中所描述的命令，在服务器端直接改变数据的值，而非窃取服务器中的数据。所以，我们要保护的对象是那些可以直接产生数据改变的服务，而对于读取数据的服务，则不需要进行 CSRF 的保护。比如银行系统中转账的请求会直接改变账户的金额，会遭到 CSRF 攻击，需要保护。而查询余额是对金额的读取操作，不会改变数据，CSRF 攻击无法解析服务器返回的结果，无需保护。

## 当前防御 CSRF 的几种策略

在业界目前防御 CSRF 攻击主要有三种策略：验证 HTTP Referer 字段；在请求地址中添加 token 并验证；在 HTTP 头中自定义属性并验证。下面就分别对这三种策略进行详细介绍。

### 验证 HTTP Referer 字段

根据 HTTP 协议，在 HTTP 头中有一个字段叫 Referer，它记录了该 HTTP 请求的来源地址。在通常情况下，访问一个安全受限页面的请求来自于同一个网站，比如需要访问 `http://bank.example/withdraw?account=bob&amount=1000000&for=Mallory`，用户必须先登陆 bank.example，然后通过点击页面上的按钮来触发转账事件。这时，该转帐请求的 Referer 值就会是转账按钮所在的页面的 URL，通常是以 bank.example 域名开头的地址。而如果黑客要对银行网站实施 CSRF 攻击，他只能在他自己的网站构造请求，当用户通过黑客的网站发送请求到银行时，该请求的 Referer 是指向黑客自己的网站。因此，要防御 CSRF 攻击，银行网站只需要对于每一个转账请求验证其 Referer 值，如果是以 bank.example 开头的域名，则说明该请求是来自银行网站自己的请求，是合法的。如果 Referer 是其他网站的话，则有可能是黑客的 CSRF 攻击，拒绝该请求。

这种方法的显而易见的好处就是简单易行，网站的普通开发人员不需要操心 CSRF 的漏洞，只需要在最后给所有安全敏感的请求统一增加一个拦截器来检查 Referer 的值就可以。特别是对于当前现有的系统，不需要改变当前系统的任何已有代码和逻辑，没有风险，非常便捷。

然而，这种方法并非万无一失。Referer 的值是由浏览器提供的，虽然 HTTP 协议上有明确的要求，但是每个浏览器对于 Referer 的具体实现可能有差别，并不能保证浏览器自身没有安全漏洞。使用验证 Referer 值的方法，就是把安全性都依赖于第三方（即浏览器）来保障，从理论上来讲，这样并不安全。事实上，对于某些浏览器，比如 IE6 或 FF2，目前已经有一些方法可以篡改 Referer 值。如果 bank.example 网站支持 IE6 浏览器，黑客完全可以把用户浏览器的 Referer 值设为以 bank.example 域名开头的地址，这样就可以通过验证，从而进行 CSRF 攻击。

即便是使用最新的浏览器，黑客无法篡改 Referer 值，这种方法仍然有问题。因为 Referer 值会记录下用户的访问来源，有些用户认为这样会侵犯到他们自己的隐私权，特别是有些组织担心 Referer 值会把组织内网中的某些信息泄露到外网中。因此，用户自己可以设置浏览器使其在发送请求时不再提供 Referer。当他们正常访问银行网站时，网站会因为请求没有 Referer 值而认为是 CSRF 攻击，拒绝合法用户的访问。

### 在请求地址中添加 token 并验证

CSRF 攻击之所以能够成功，是因为黑客可以完全伪造用户的请求，该请求中所有的用户验证信息都是存在于 cookie 中，因此黑客可以在不知道这些验证信息的情况下直接利用用户自己的 cookie 来通过安全验证。要抵御 CSRF，关键在于在请求中放入黑客所不能伪造的信息，并且该信息不存在于 cookie 之中。可以在 HTTP 请求中以参数的形式加入一个随机产生的 token，并在服务器端建立一个拦截器来验证这个 token，如果请求中没有 token 或者 token 内容不正确，则认为可能是 CSRF 攻击而拒绝该请求。

这种方法要比检查 Referer 要安全一些，token 可以在用户登陆后产生并放于 session 之中，然后在每次请求时把 token 从 session 中拿出，与请求中的 token 进行比对，但这种方法的难点在于如何把 token 以参数的形式加入请求。对于 GET 请求，token 将附在请求地址之后，这样 URL 就变成 `http://url?csrftoken=tokenvalue`。 而对于 POST 请求来说，要在 form 的最后加上 `<input type=”hidden” name=”csrftoken” value=”tokenvalue”/>`，这样就把 token 以参数的形式加入请求了。但是，在一个网站中，可以接受请求的地方非常多，要对于每一个请求都加上 token 是很麻烦的，并且很容易漏掉，通常使用的方法就是在每次页面加载时，使用 javascript 遍历整个 dom 树，对于 dom 中所有的 a 和 form 标签后加入 token。这样可以解决大部分的请求，但是对于在页面加载之后动态生成的 html 代码，这种方法就没有作用，还需要程序员在编码时手动添加 token。

该方法还有一个缺点是难以保证 token 本身的安全。特别是在一些论坛之类支持用户自己发表内容的网站，黑客可以在上面发布自己个人网站的地址。由于系统也会在这个地址后面加上 token，黑客可以在自己的网站上得到这个 token，并马上就可以发动 CSRF 攻击。为了避免这一点，系统可以在添加 token 的时候增加一个判断，如果这个链接是链到自己本站的，就在后面添加 token，如果是通向外网则不加。不过，即使这个 csrftoken 不以参数的形式附加在请求之中，黑客的网站也同样可以通过 Referer 来得到这个 token 值以发动 CSRF 攻击。这也是一些用户喜欢手动关闭浏览器 Referer 功能的原因。

### 在 HTTP 头中自定义属性并验证

这种方法也是使用 token 并进行验证，和上一种方法不同的是，这里并不是把 token 以参数的形式置于 HTTP 请求之中，而是把它放到 HTTP 头中自定义的属性里。通过 XMLHttpRequest 这个类，可以一次性给所有该类请求加上 csrftoken 这个 HTTP 头属性，并把 token 值放入其中。这样解决了上种方法在请求中加入 token 的不便，同时，通过 XMLHttpRequest 请求的地址不会被记录到浏览器的地址栏，也不用担心 token 会透过 Referer 泄露到其他网站中去。

然而这种方法的局限性非常大。XMLHttpRequest 请求通常用于 Ajax 方法中对于页面局部的异步刷新，并非所有的请求都适合用这个类来发起，而且通过该类请求得到的页面不能被浏览器所记录下，从而进行前进，后退，刷新，收藏等操作，给用户带来不便。另外，对于没有进行 CSRF 防护的遗留系统来说，要采用这种方法来进行防护，要把所有请求都改为 XMLHttpRequest 请求，这样几乎是要重写整个网站，这代价无疑是不能接受的。

## Java 代码示例

下文将以 Java 为例，对上述三种方法分别用代码进行示例。无论使用何种方法，在服务器端的拦截器必不可少，它将负责检查到来的请求是否符合要求，然后视结果而决定是否继续请求或者丢弃。在 Java 中，拦截器是由 Filter 来实现的。我们可以编写一个 Filter，并在 web.xml 中对其进行配置，使其对于访问所有需要 CSRF 保护的资源的请求进行拦截。

在 filter 中对请求的 Referer 验证代码如下

##### 清单 1\. 在 Filter 中验证 Referer

```
// 从 HTTP 头中取得 Referer 值
String referer=request.getHeader("Referer");
// 判断 Referer 是否以 bank.example 开头
if((referer!=null) &&(referer.trim().startsWith("bank.example”))){
    chain.doFilter(request, response);
}else{
    request.getRequestDispatcher("error.jsp”).forward(request,response);
}

```

Show moreShow more icon

以上代码先取得 Referer 值，然后进行判断，当其非空并以 bank.example 开头时，则继续请求，否则的话可能是 CSRF 攻击，转到 error.jsp 页面。

如果要进一步验证请求中的 token 值，代码如下

##### 清单 2\. 在 filter 中验证请求中的 token

```
HttpServletRequest req = (HttpServletRequest)request;
HttpSession s = req.getSession();

// 从 session 中得到 csrftoken 属性
String sToken = (String)s.getAttribute("csrftoken”);
if(sToken == null){

    // 产生新的 token 放入 session 中
    sToken = generateToken();
    s.setAttribute("csrftoken”,sToken);
    chain.doFilter(request, response);
} else{

    // 从 HTTP 头中取得 csrftoken
    String xhrToken = req.getHeader("csrftoken”);

    // 从请求参数中取得 csrftoken
    String pToken = req.getParameter("csrftoken”);
    if(sToken != null && xhrToken != null && sToken.equals(xhrToken)){
        chain.doFilter(request, response);
    }else if(sToken != null && pToken != null && sToken.equals(pToken)){
        chain.doFilter(request, response);
    }else{
        request.getRequestDispatcher("error.jsp”).forward(request,response);
    }
}

```

Show moreShow more icon

首先判断 session 中有没有 csrftoken，如果没有，则认为是第一次访问，session 是新建立的，这时生成一个新的 token，放于 session 之中，并继续执行请求。如果 session 中已经有 csrftoken，则说明用户已经与服务器之间建立了一个活跃的 session，这时要看这个请求中有没有同时附带这个 token，由于请求可能来自于常规的访问或是 XMLHttpRequest 异步访问，我们分别尝试从请求中获取 csrftoken 参数以及从 HTTP 头中获取 csrftoken 自定义属性并与 session 中的值进行比较，只要有一个地方带有有效 token，就判定请求合法，可以继续执行，否则就转到错误页面。生成 token 有很多种方法，任何的随机算法都可以使用，Java 的 UUID 类也是一个不错的选择。

除了在服务器端利用 filter 来验证 token 的值以外，我们还需要在客户端给每个请求附加上这个 token，这是利用 js 来给 html 中的链接和表单请求地址附加 csrftoken 代码，其中已定义 token 为全局变量，其值可以从 session 中得到。

##### 清单 3\. 在客户端对于请求附加 token

```
function appendToken(){
    updateForms();
    updateTags();
}

function updateForms() {
    // 得到页面中所有的 form 元素
    var forms = document.getElementsByTagName('form');
    for(i=0; i<forms.length; i++) {
        var url = forms[i].action;

        // 如果这个 form 的 action 值为空，则不附加 csrftoken
        if(url == null || url == "" ) continue;

        // 动态生成 input 元素，加入到 form 之后
        var e = document.createElement("input");
        e.name = "csrftoken";
        e.value = token;
        e.type="hidden";
        forms[i].appendChild(e);
    }
}

function updateTags() {
    var all = document.getElementsByTagName('a');
    var len = all.length;

    // 遍历所有 a 元素
    for(var i=0; i<len; i++) {
        var e = all[i];
        updateTag(e, 'href', token);
    }
}

function updateTag(element, attr, token) {
    var location = element.getAttribute(attr);
    if(location != null && location != '' '' ) {
        var fragmentIndex = location.indexOf('#');
        var fragment = null;
        if(fragmentIndex != -1){

            //url 中含有只相当页的锚标记
            fragment = location.substring(fragmentIndex);
            location = location.substring(0,fragmentIndex);
        }

        var index = location.indexOf('?');

        if(index != -1) {
            //url 中已含有其他参数
            location = location + '&csrftoken=' + token;
        } else {
            //url 中没有其他参数
            location = location + '?csrftoken=' + token;
        }
        if(fragment != null){
            location += fragment;
        }

        element.setAttribute(attr, location);
    }
}

```

Show moreShow more icon

在客户端 html 中，主要是有两个地方需要加上 token，一个是表单 form，另一个就是链接 a。这段代码首先遍历所有的 form，在 form 最后添加一隐藏字段，把 csrftoken 放入其中。然后，代码遍历所有的链接标记 a，在其 href 属性中加入 csrftoken 参数。注意对于 a.href 来说，可能该属性已经有参数，或者有锚标记。因此需要分情况讨论，以不同的格式把 csrftoken 加入其中。

如果你的网站使用 XMLHttpRequest，那么还需要在 HTTP 头中自定义 csrftoken 属性，利用 dojo.xhr 给 XMLHttpRequest 加上自定义属性代码如下：

##### 清单 4\. 在 HTTP 头中自定义属性

```
var plainXhr = dojo.xhr;

// 重写 dojo.xhr 方法
dojo.xhr = function(method,args,hasBody) {
    // 确保 header 对象存在
    args.headers = args.header || {};

    tokenValue = '<%=request.getSession(false).getAttribute("csrftoken")%>';
    var token = dojo.getObject("tokenValue");

    // 把 csrftoken 属性放到头中
    args.headers["csrftoken"] = (token) ? token : "  ";
    return plainXhr(method,args,hasBody);
};

```

Show moreShow more icon

这里改写了 dojo.xhr 的方法，首先确保 dojo.xhr 中存在 HTTP 头，然后在 args.headers 中添加 csrftoken 字段，并把 token 值从 session 里拿出放入字段中。

## CSRF 防御方法选择之道

通过上文讨论可知，目前业界应对 CSRF 攻击有一些克制方法，但是每种方法都有利弊，没有一种方法是完美的。如何选择合适的方法非常重要。如果网站是一个现有系统，想要在最短时间内获得一定程度的 CSRF 的保护，那么验证 Referer 的方法是最方便的，要想增加安全性的话，可以选择不支持低版本浏览器，毕竟就目前来说，IE7+, FF3+ 这类高版本浏览器的 Referer 值还无法被篡改。

如果系统必须支持 IE6，并且仍然需要高安全性。那么就要使用 token 来进行验证，在大部分情况下，使用 XmlHttpRequest 并不合适，token 只能以参数的形式放于请求之中，若你的系统不支持用户自己发布信息，那这种程度的防护已经足够，否则的话，你仍然难以防范 token 被黑客窃取并发动攻击。在这种情况下，你需要小心规划你网站提供的各种服务，从中间找出那些允许用户自己发布信息的部分，把它们与其他服务分开，使用不同的 token 进行保护，这样可以有效抵御黑客对于你关键服务的攻击，把危害降到最低。毕竟，删除别人一个帖子比直接从别人账号中转走大笔存款严重程度要轻的多。

如果是开发一个全新的系统，则抵御 CSRF 的选择要大得多。笔者建议对于重要的服务，可以尽量使用 XMLHttpRequest 来访问，这样增加 token 要容易很多。另外尽量避免在 js 代码中使用复杂逻辑来构造常规的同步请求来访问需要 CSRF 保护的资源，比如 window.location 和 document.createElement(“a”) 之类，这样也可以减少在附加 token 时产生的不必要的麻烦。

最后，要记住 CSRF 不是黑客唯一的攻击手段，无论你 CSRF 防范有多么严密，如果你系统有其他安全漏洞，比如跨站域脚本攻击 XSS，那么黑客就可以绕过你的安全防护，展开包括 CSRF 在内的各种攻击，你的防线将如同虚设。

## 结束语

可见，CSRF 是一种危害非常大的攻击，又很难以防范。目前几种防御策略虽然可以很大程度上抵御 CSRF 的攻击，但并没有一种完美的解决方案。一些新的方案正在研究之中，比如对于每次请求都使用不同的动态口令，把 Referer 和 token 方案结合起来，甚至尝试修改 HTTP 规范，但是这些新的方案尚不成熟，要正式投入使用并被业界广为接受还需时日。在这之前，我们只有充分重视 CSRF，根据系统的实际情况选择最合适的策略，这样才能把 CSRF 的危害降到最低。