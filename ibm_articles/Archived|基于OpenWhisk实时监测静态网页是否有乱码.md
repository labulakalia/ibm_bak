# Archived | 基于 OpenWhisk 实时监测静态网页是否有乱码
OpenWhisk 实战

**标签:** DevOps,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/wa-lo-base-openwhisk-ontime-monitor-v1/)

傅 智勇

发布: 2017-10-10

* * *

**本文已归档**

**归档日期：:** 2020-02-14

此内容不再被更新或维护。 内容是按“原样”提供。鉴于技术的快速发展，某些内容，步骤或插图可能已经改变。

## 概览

IBM OpenWhisk (注意：OpenWhisk 现已贡献给 Apache, IBM OpenWhisk 现已更名为 [IBM Cloud Founctions](https://cloud.ibm.com/functions/?cm_sp=ibmdev-_-developer-articles-_-cloudreg)) 是一种开源的微服务计算平台，它执行应用程序逻辑来响应事件，或通过直接调用来执行某项特定的任务。本文通过一个实例，来讲述如何应用 IBM OpenWhisk 封装好的服务，和如何在 OpenWhisk 使用第三方服务，以及如何自定义 OpenWhisk 服务来提供一套 DevOps 的静态网站乱码检测解决方案。

## OpenWhisk 简要概述

首先介绍 OpenWhisk 采用了什么机制，以及为我们提供什么服务：

- OpenWhisk 采用 Action，Trigger 和 Rule 的基于事件的编程模型，使我们可以将主要精力集中在业务逻辑的开发上，而不是运行时的软硬件环境的部署和运维上；
- 可以将 IBM Cloud 上的 Service、开源，以及第三方 Service 打成标准的 Catalog 包，以在 OpenWhisk 上使用；
- 细粒度的计费方式，以及透明的自动伸缩，使我们只为需要的资源付费，我们不需要为业务无关的资源买单；
- OpenWhisk 提供的 iOS SDK和Starter应用程序，以及服务端支持 Swift 语言，使我们可以写出非常酷的Mobile应用。

其中OpenWhisk包括 Action，Trigger，Rule，Sequence 以及 Action Runtime 几个部件，下面对这几个部件进行简要说明：

### Actions（行为）

Action就是您需要执行的一个特定任务，比如：本文实例中的Garbage Char Detection就是一个Action，发送检测结果到Slack也是一个Action；可以采用您比较擅长的编程语言来编写Action； 你可以通过在Action代码里调用REST API来完成一项特定的任务，比如该实例调用的前面介绍Garbage Char Detect REST API，和Slack的Post Text API；Action也可以自动响应由 IBM Cloud 或第三方服务的触发器触发的事件，比如Garbage Detect行为和发送扫描结果到Slack行为就是响应GitHub push事件。

### Trigger（触发器）和 Rule（规则）

Trigger是响应一个特定的事件的申明；Trigger既可以人为触发，也可以接受到某个特定事件被自动触发。Rule则是将一个Trigger与具体的Action进行关联，当一个Trigger被特定的事件触发时，与该触发器关联的Action就会执行。

### Sequence（序列）

所谓Sequence就是一个有序的Action序列，Action序列中，前一个Action的输出，会作为后一个Action的输入；本文讲解实例中的Garbage Char Detection Action和发送扫描结果到Slack Action就是通过Sequence串联起来的；Sequence 可以像单个的Action那样被调用。

### Runtimes（运行时环境）

Runtimes即行为的运行时环境，目前支持的有Node.js，Python和Swift3。当然，您也可以采用其他的编程语言，不过需要将你的Action代码放到Docker容器上运行；您可以采用 IBM CLoud OpenWhisk 提供了基于Web接口的IDE来开发和执行Whisk程序，也可以选择你喜欢的IDE，然后通过OpenWhisk CLI将代码上传到Whisk上执行。

## OpenWhisk 应用场景

本文介绍的第一个应用场景是，OpenWhisk通过调用 IBM Cloud 上的天气预报服务去收集气象数据，然后通过 Filter Action 抽取出感兴趣的数据，再通过Waston MT Translation 服务将天气情况翻译成指定语言，发送到对应的终端。其流程图如图1：

#### 图 1\. OpenWhisk 天气预报

![OpenWhisk 天气预报](../ibm_articles_img/wa-lo-base-openwhisk-ontime-monitor-v1_images_image001.png)

第二个应用场景是，用户通过上传图片到社交 APP，图片的路径保存在 Cloudant DB，Cloudant DB 检测到有更新，会触发 Waston的 AlchemyAPI Action，对图片进行分析，并自动添加响应的 tag。其流程图如 图2：

#### 图 2\. OpenWhisk 给图片自动添加 tag

![penWhisk 给图片自动添加 tag](../ibm_articles_img/wa-lo-base-openwhisk-ontime-monitor-v1_images_image002.png)

第三个应用场景是，OpenWhisk上的GitHub WebHook触发器去接收GitHub特定repository中的特定事件，然后对事件返回的Payload数据进行分析过滤，翻译并且post结果到Slack。其流程图如图3：

#### 图 3\. Git2Slack

![Git2Slack](../ibm_articles_img/wa-lo-base-openwhisk-ontime-monitor-v1_images_image003.png)

## 实例：基于OpenWhisk实时监测静态页面是否有乱码

接下来以网页内容是否包含乱码字符的实例讲解第三个应用场景的具体实现流程，不过有点不一样，这里将翻译Action改成了我们自定义的Garbage Char Detection Action。其流程图如图4：

#### 图 4\. OpenWhisk Garbage Char Detection from GitHub to Slack

![OpenWhisk Garbage Char Detection from GitHub to Slack](../ibm_articles_img/wa-lo-base-openwhisk-ontime-monitor-v1_images_image004.png)

实现该解决方案需要以下几个步骤，参考 [github openwhisk-test](https://github.com/icnbrave/openwhisk-test/tree/master/wsk) ：

1. 在GitHub上创建一个Repository，并将静态页面的代码Push到该Repository，然后部署到一个Web服务器，本文选择部署到 IBM Cloud 平台之上；
2. 创建OpenWhisk GitHub触发器；
3. 自定义的Filter Action从GitHub事件中过滤出感兴趣数据；
4. 自定义Garbage Char Detection Action对有改动或新增的页面进行乱码检测；
5. 自定义Slack Post Action将检测结果发送到Slack指定的Channel。

### 实现乱码检测 API

在讲解该实例的具体步骤之前，首先来了解一下Garbage Char Detection API。乱码检测服务代码可以从 [https://github.com/icnbrave/openwhisk-test](https://github.com/icnbrave/openwhisk-test) 下载，其GarbageCodeDetect目录是乱码检测服务对应的代码，SampleApp是基于GarbageCodeDetect服务的演示应用程序。wsk目录包括本实例中会用到的一些Action演示代码。

将 GarbageCodeDetect 服务以java应用程序部署到 IBM Cloud 之后，就可以使用Garbage Code Detect（乱码检测）API。

应用程序 Web UI 的 URL 为 `http://garbagecodedetection.mybluemix.net` ，乱码检测 REST API 则是 `http://garbagecodedetection.mybluemix.net/rest/garbagechar_scan` ，其中测试目录下面也有一些包含乱码的静态页面:

`http://garbagecodedetection.mybluemix.net/test/garbledUTF8-2.html`

`http://garbagecodedetection.mybluemix.net/test/garbledUTF8.html`

`http://garbagecodedetection.mybluemix.net/test/garbledBig5.html`

不包含乱码的静态页面有：

`http://garbagecodedetection.mybluemix.net/test/goodUTF8.html`

`http://garbagecodedetection.mybluemix.net/test/goodBig5.html`

乱码检测 REST API 的HTTP方法是POST，参数是一个JSON数组，其中每个数组元素都包括一个必须参数url，和一个可选参数encoding，例如：

若没有指定页面的编码，则服务会使用页面的默认编码。该API的作用是检测JSON数组中提供所有的页面中所有含有乱码的字串，如果用户传入API的JSON数据如上有乱码字串，API 会返回:

![传入 API 的 JSON 数据如上有乱码字串，API 返回](../ibm_articles_img/wa-lo-base-openwhisk-ontime-monitor-v1_images_image005.png)

### 为静态网页内容创建 GitHub Repository

首先，需要将测试网站内容上传到 GitHub上，并且部署到Web服务器之上。

本实例网站内容托管的GitHub地址是 [https://github.com/icnbrave/garbage-test-app](https://github.com/icnbrave/garbage-test-app) ；可以通过以下几个步骤，将该代码以静态网站的形式部署到 IBM Cloud 平台之上：

```
$ cf api https://api.ng.bluemix.net
$ cf login <with you bluemix ID>
$ git clone https://github.com/icnbrave/garbage-test-app.git
$ cd garbage-test-app
$ touch Staticfile
$ echo "directory:visible" > Staticfile
$ cf push garbage-test-app -b

```

Show moreShow more icon

[https://github.com/cloudfoundry/staticfile-buildpack.git](https://github.com/cloudfoundry/staticfile-buildpack.git)

部署成功之后，在 `https://garbage-test-app.mybluemix.net/` 可以查看到静态文件列表，如图5所示，点击其中某个文件，可以看到具体内容。

#### 图 5\. 静态网页内容列表

![静态网页内容列表](../ibm_articles_img/wa-lo-base-openwhisk-ontime-monitor-v1_images_image006.png)

### 创建 OpenWhisk GitHub 触发器

可以通过以下步骤为garbage-test-app创建GitHub触发器：

1. 首先为测试应用程序创建 package binding，创建 package binding 时需要指定 github repository，accessToken和管理用户；





    ```
    $ wsk package update myGit -p  repository icnbrave/garbage-test-app -p accessToken  <GITHUB_ACCESSTOKEN> -p username <GITHUB_USERNAME>

    ```





    Show moreShow more icon

2. 为 github push 事件创建触发器，即绑定的 github repo 有 push 时，该触发器被触发。





    ```
    $ wsk trigger create gitTrigger --feed myGit/webhook -p events push

    ```





    Show moreShow more icon


当触发器创建成功之后，可以在GitHub Repository的Settings中会创建一个Webhook，GitHub就是通过该Webhook去触发OpenWhisk触发器，如图6所示：

#### 图 6\. github webhooks 设置

![github webhooks设置](../ibm_articles_img/wa-lo-base-openwhisk-ontime-monitor-v1_images_image007.png)

### 自定义的 Filter Action 从 GitHub 事件中过滤出感兴趣数据

GitHub的每个事件都会附带对应的 payload 数据，可以通过 [Event Types & Payloads](https://developer.github.com/v3/activity/events/types/) 查看每GitHub Event对应的payload数据，我们可以从payload中提取我们想要的数据。比如，我们可以从push事件对应的payload数据中提取最近几次commit中更新、添加的文件，然后只对这部分文件进行garbage char检测。

GitHub Push Event Payload部分数据如下：

##### 清单 1\. GitHub Push Event Payload 部分返回数据

```
...
"commits": [
{
"id": "0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c",
"tree_id": "f9d2a07e9488b91af2641b26b9407fe22a451433",
"distinct": true,
"message": "Update README.md",
"timestamp": "2015-05-05T19:40:15-04:00",
"url":
      "https://github.com/baxterthehacker/public-repo/commit/0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c",
"author": {
"name": "baxterthehacker",
"email": "baxterthehacker@users.noreply.github.com",
"username": "baxterthehacker"
},
"committer": {
"name": "baxterthehacker",
"email": "baxterthehacker@users.noreply.github.com",
"username": "baxterthehacker"
},
"added": [
],
"removed": [
],
"modified": [
"README.md"
]
}
],
...

```

Show moreShow more icon

通过如下命令将以下 [getPushPayload.js](https://github.com/icnbrave/openwhisk-test/blob/master/wsk/getPushPayload.js) 创建getPushPayloadAction来接收并且过滤当前push事件中新增或更改的文件，然后得出他们在 IBM Cloud 上的URL并送到下一个 action：

```
$ wsk action create getPushPayloadAction getPushPayload.js

```

Show moreShow more icon

##### 清单 2\. getPushPayload.js

```
function main(params){
var testAppUrl = "https://garbagetestapp.mybluemix.net/";
var head_commit = params["head_commit"] || "";
var files = [];
if(head_commit != ""){
files=files.concat(head_commit["added"]);
files=files.concat(head_commit["modified"]);
}
console.log("files: ", files);
// files map to test app link
var urls=[];
for(var i=0,len=files.length; i<len; i++){
urls.push(testAppUrl + files[i]);
}
return {payload: urls.join(",")};
}

```

Show moreShow more icon

然后通过如下命令创建一个Rule将GitHub Trigger与该Action进行关联：

```
$ wsk rule create git2slackRule gitTrigger getPushPayloadAction

```

Show moreShow more icon

关联之后， garbage-test-app repo 任何 push event 都会触发这个Action。

### 自定义Garbage Char Detection Action对有改动或新增的页面进行乱码检测

第四步主要介绍 Garbage Char Detection API，以及如何在Action中接收getPushPayloadAction传来数据，并调用Garbage Char Detection API。

用 [garbageDetectionAction.js](https://github.com/icnbrave/openwhisk-test/blob/master/wsk/garbageDetectionAction.js) 创建garbageDetectionAction：

```
$ wsk action create garbageDetectionAction garbageDetectionAction.js

```

Show moreShow more icon

garbageDetectionAction.js内容如代码清单2：

##### 清单 3\. garbageDetectionAction.js

```
var request = require('request');
var http = require("http");
var querystring = require('querystring');
var utils = {
map : function(arr, func){
var res = [];
for(var i=0,len=arr.length; i<len; i++){
res.push(func(arr[i]));
}
return res;
},
encodeStrings : function(arr){
return utils.map(arr, querystring.escape);
},
decodeStrings : function(arr){
return utils.map(arr, decodeURI);
},
};
function main(params) {
// payload is urls with ',' split, such as "url1,url2"
var payload = params["payload"].split(',') || [];
var data=[];
for(var i=0,len=payload.length; i<len; i++){
data.push({'url': payload[i]});
}
console.log(data);
var options = {
host: 'garbagecodedetection.mybluemix.net',
port: 80,
path: '/rest/garbagechar_scan',
method: 'POST',
headers: {
'Content-Type': 'application/json',
}
};
var req = http.request(options, function (res) {
res.on('data', function (message) {
var ret= eval('(' + message + ')');
console.log('response : ' ,ret);
console.log('response type: ', typeof ret);
for(var i=0; i<ret.length; i++) {
ret[i]['garbled_lines'] =
     utils.encodeStrings(ret[i]['garbled_lines']);
}
console.log('Encoded response: ', ret)
whisk.done({result: ret});
});
});
req.on('error', function(e) {
console.log('problem with request: ' + e.message);
whisk.done({text: 'Problem with request' + e.message});
});
// write data to request body
req.write(JSON.stringify(data));
req.end();
return whisk.async();
}

```

Show moreShow more icon

通过下面命令来测试garbageDetectionAction：

```
$ wsk action invoke garbageDetectAction --blocking --result -p payload

```

Show moreShow more icon

`http://garbagecodedetection.mybluemix.net/test/garbledUTF8-2.html` 和 `http://garbagecodedetection.mybluemix.net/test/garbledBig5.html`在终端会输出如下结果如图 7：

#### 图 7\. 检测结果

![检测结果](../ibm_articles_img/wa-lo-base-openwhisk-ontime-monitor-v1_images_image008.png)

### 自定义Slack Post Action将检测结果发送到Slack指定的Channel

OpenWhisk Catalog自带了Slack Post Action，但是自带的Slack Post Action功能有限，不能满足本实例的需求，所以需要自定义Slack Post方法来发送检测结果到Slack Channel。自定的基本步骤如下：

1. 为您的团队创建一个 Slack 并且配置 Slack [incoming webhook](https://api.slack.com/incoming-webhooks) ，配置好Slack之后，可以拿到类似如下的 Webhook URL：`https://hooks.slack.com/services/aaaaaaaaa/bbbbbbbbb/ccccccccccccccccccccccc`。
2. 用生成的 Webhook URL 和对应的Channel创建一个 Package





    ```
    $ wsk package create myCustomSlack --param url "https://hooks.slack.com/services/..." --param username Bob --param channel "#MySlackChannel"

    ```





    Show moreShow more icon

3. 需要写一个函数接收 Garbage Char Detection Action 传来的数据，然后发送到 Slack ，并用该函数来创建 post2slack Action





    ```
    $ wsk action create myCustomSlack/post2slack slackPost.js

    ```





    Show moreShow more icon


其中 [slackPost.js](https://github.com/icnbrave/openwhisk-test/blob/master/wsk/slackPost.js) 内容如下：

##### 清单 4\. SlackPost.js

```

var request = require('request');
var utils = {
map : function(arr, func){
var res = [];
for(var i=0,len=arr.length; i<len; i++){
res.push(func(arr[i]));
}
return res;
},
encodeStrings : function(arr){
return utils.map(arr, querystring.escape);
},
decodeStrings : function(arr){
return utils.map(arr, decodeURI);
},
};
/**
* Action to post to slack
* @param {string} url - Slack webhook url
* @param {string} channel - Slack channel to post the message to
* @param {string} username - name to post to the message as
* @param {string} text - message to post
* @param {string} icon_emoji - (optional) emoji to use as the icon for the
     message
* @param {string} as_user - (optional) when the token belongs to a bot,
     whether to post as the bot itself
* @param {object} attachments - (Optional) message attachments
* @return {object} whisk async
*/
function main(params){
if (checkParams(params)){
d = params.result;
console.log('Result: ', d);
for(var i=0,len=d.length; i<len; i++){
d[i]['garbled_lines'] = utils.decodeStrings(d[i]['garbled_lines']);
}
console.log('Decoded Result: ', d);
var body = {
channel: params.channel,
username: params.username || 'Simple Message Bot',
text: format(JSON.stringify(d, null, '\t'))
};
if (params.icon_emoji){
body.icon_emoji = params.icon_emoji;
}
if (params.token){
//
// this allows us to support /api/chat.postMessage
// e.g. users can pass params.url =
     https://slack.com/api/chat.postMessage
// and params.token = \u003ctheir auth token\u003e
//
body.token = params.token;
} else {
//
// the webhook api expects a nested payload
//
// Notice that we need to stringify; this is due to limitations
// of the formData npm: it does not handle nested objects
//
console.log(body);
console.log("to: "+ params.url);
body = {
payload: JSON.stringify(body)
};
}
if (params.as_user === true){
body.as_user = true;
}
if (params.attachments) {
body.attachments = params.attachments;
}
var promise = new Promise(function(resolve, reject){
request.post({
url: params.url,
formData: body,
}, function(err, res, body) {
if (err) {
console.log('Error: ', err, body);
reject(err);
} else {
console.log('Success: ', params.text, ' successfully sent');
resolve();
}
});
});
return promise;
}
}
/**
*
* Checks if all required params are set
*/
function checkParams(params){
console.log('Post2Slack params: ', params);
if (params.result === undefined) {
whisk.error('No post data provided');
return false;
}
if (params.url === undefined) {
whisk.error('No webhook URL provided');
return false;
}
if (params.channel === undefined) {
whisk.error('No channel provided');
return false;
}
return true;
}
/**
* format text to slack
*/
function format(str){
return '\n\`\`\`'+str+'\n\`\`\`';
}

```

Show moreShow more icon

1. 将 post2slack action 加入到 git2slack Sequence

```
$ wsk action create myCustomSlack/post2slack slackPost.js

```

Show moreShow more icon

为了调试方便，我们没有让 GitHub Trigger 来触发，还是先手动传入测试数据

```
$ wsk action create git2slack --sequence garbageDetectionAction,myCustomSlack/post2slack

```

Show moreShow more icon

测试 git2slack 序列：

```
$ wsk action invoke git2slack --blocking --result -p payload

```

Show moreShow more icon

`http://garbagecodedetection.mybluemix.net/test/garbledUTF8-2.html` 和 `http://garbagecodedetection.mybluemix.net/test/garbledBig5.html` 输出结果如图 8：

#### 图 8\. 触发 git2slack action 结果

![触发 git2slack action 结果](../ibm_articles_img/wa-lo-base-openwhisk-ontime-monitor-v1_images_image009.png)

当一切调试通过之后，我们在更新 git2slack 序列，将 getPushPayloadAction 加入到 git2slack 序列，并更新 git2slackRule.

通过以下命令更新 git2slack 序列：

```
$ wsk action update git2slack –sequence getPushPayloadAction,garbageDetectAction,myCustomSlack/post2slack

```

Show moreShow more icon

通过如下命令集更新 git2slackRule：

```
$ wsk rule disable git2slackRule
$ wsk rule update git2slackRule gitTrigger git2slack
$ wsk rule enable git2slackRule

```

Show moreShow more icon

完成以上步骤之后，我们更新测试应用程序中的文档 garbledUTF8-2.html，然后部署到 IBM Cloud 平台，并且 push 到 GitHub Repository，这时我们在Slack 会收到如下通知：

#### 图 9\. Git 仓库更新触发 git2slck action 图

![Git 仓库更新触发 git2slck action 图](../ibm_articles_img/wa-lo-base-openwhisk-ontime-monitor-v1_images_image010.png)

如图10，我们可以看到所有的 Unicode在Slack 上都显示”?”.经过查看 [OpenWhisk](https://github.com/openwhisk/openwhisk) 源代码，发现OpenWhisk 只支持 ASCII，并不支持 Unicode.

#### 图 10\. Unicode 字符显示异常图

![Unicode 字符显示异常图](../ibm_articles_img/wa-lo-base-openwhisk-ontime-monitor-v1_images_image010.png)

#### 解决方案

通过 [https://github.com/openwhisk/openwhisk/issues/252](https://github.com/openwhisk/openwhisk/issues/252) 提供的方法，我们可以在 OpenWhisk 处理或传输 Unicode 之前，将 Unicode 进行编码，然后在发送 Unicode 到 Slack 的 Action 中进行解码就可以让 Unicode 字符正确的在Slack 上显示。相应代码修改如下：

更改 [garbageDetectionAction.js](https://github.com/icnbrave/openwhisk-test/blob/master/wsk/garbageDetectionAction.js)

```
for(var i=0; i<ret.length; i++)
{
    ret[i]['garbled_lines'] = utils.encodeStrings(ret[i]['garbled_lines']);
}

```

Show moreShow more icon

更改 [slackPost.js](https://github.com/icnbrave/openwhisk-test/blob/master/wsk/slackPost.js)

```
for(var i=0,len=d.length; i<len; i++){
d[i]['garbled_lines'] = utils.decodeStrings(d[i]['garbled_lines']);
}

```

Show moreShow more icon

其中，utils方法如下：

```
var utils = {
map : function(arr, func){
var res = [];
for(var i=0,len=arr.length; i<len; i++){
res.push(func(arr[i]));
}
return res;
},
encodeStrings : function(arr){
return utils.map(arr, querystring.escape);
},
decodeStrings : function(arr){
return utils.map(arr, decodeURI);
},
};

```

Show moreShow more icon

更新好代码之后，需要更新Action garbageDetectionAction，myCustomSlack/post2slack, 以及git2slack sequence。然后，按照测试步骤，更新 [garbage-test-app](https://github.com/icnbrave/garbage-test-app) 某个文件，然后deploy到 IBM Cloud，再push到GitHub Repository，可以在Slack收到正常的返回结果，如图11所示：

#### 图 11\. Unicode 显示正常图

![Unicode 显示正常图](../ibm_articles_img/wa-lo-base-openwhisk-ontime-monitor-v1_images_image011.png)

## 结束语

本文阐述了 IBM OpenWhisk 简要概述，和一些应用场景，以及结合一个实例讲解如何编写和使用 OpenWhisk 相关组件来完成一个 DevOps 的解决方案。最后，给出了实例实现过程遇到的一个问题，以及对应的解决方法。