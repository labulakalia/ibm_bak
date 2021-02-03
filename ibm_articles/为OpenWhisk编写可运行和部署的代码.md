# 为 OpenWhisk 编写可运行和部署的代码
如何编写可在 IBM Developer 沙箱中运行和部署的代码

**标签:** JavaScript,Serverless,Web 开发,云计算

[原文链接](https://developer.ibm.com/zh/articles/wa-write-deployable-code-for-openwhisk/)

Doug Tidwell

发布: 2017-08-17

* * *

[IBM Developer 沙箱](https://developer.ibm.com/dwblog/2017/sandbox-interactive-tutorials-developerworks/) 的一个好处是，能够在浏览器中运行代码或将代码部署到 IBM Cloud。您部署的代码是一个 OpenWhisk 操作。但是，OpenWhisk 要求您的代码支持特定的接口，然后您才可以部署它。本教程将解释实现所有这些目的的规则和需求。

在本教程中，您将了解如何编写通过 IBM Developer 沙箱在浏览器中运行，或者作为 OpenWhisk 操作进行部署的 JavaScript 或 Java™ 代码。本教程中的沙箱帮助您直接构建、测试和使用代码。每个沙箱后面都有详细的代码解释。

在本教程中，您将：

1. 看到包含可运行和部署的 JavaScript 和 Java 代码的 IBM Developer 沙箱
2. 了解如何让代码可部署
3. 了解如何让可部署的代码可运行

## 关于 OpenWhisk 或无服务器的操作

简单来讲，OpenWhisk 要求您的代码接受一个 JSON 对象作为其唯一输入，并返回一个 JSON 对象作为其输出。如果使用 JavaScript，则用 JSON 很容易处理；如果使用 Java，需要使用 Google 的 `gson` 库来处理 JSON 对象。下面将介绍处理 JSON 和 `gson` 的所有细节。最后需要注意，构建发送到您的代码的输入 JSON 对象并解释它的 JSON 输出，是调用您的 OpenWhisk 操作的开发人员的职责。

要了解如何在部署之后调用代码，请参阅 IBM Developer 教程 [调用 OpenWhisk 操作](https://www.ibm.com/developerworks/cn/web/wa-invoke-openwhisk-action/index.html) 。

## 一个可运行和部署的简单 JavaScript 示例

兼容 OpenWhisk 的 Hello World 是既能运行又能部署的最简单代码。

```
function main(params) {
var greeting = 'Hello, ' + params.name;
return {greeting};
}

var defaultParameters = {'name': 'Susan'};

if (require.main === module)
console.log(main(defaultParameters));

```

Show moreShow more icon

### 测试代码一

在讨论如何使此代码可运行和部署之前，让我们先测试它。单击上面的 **Run** 按钮。根据针对 `name` 字段传入的值，您会看到类似这样的结果：

```
{ greeting: 'Hello, Susan' }

```

Show moreShow more icon

现在单击 **Deploy** 按钮。沙箱将代码作为一个新操作部署到 OpenWhisk，并显示它的 URL：

```
https://openwhisk.ng.bluemix.net/api/v1/web/devworkssandbox_pub/270003KAD5/0-1001.json

```

Show moreShow more icon

_备注：_ 生成的 URL 对您是唯一的，所以您看到的 URL 会与此处示例中的有所不同。

现在调用我们刚部署到 OpenWhisk 的代码。在下面的沙箱中，将 `[URL]` 替换为上面的 **Deploy** 按钮所生成的 URL。（显然，如果您喜欢的话，还可以更改 `name` 。）按下面的 **Run** 按钮时，下面沙箱中的代码会调用您在部署上述代码时创建的 OpenWhisk 操作。

```
var options = {
// Paste in the URL of the deployed sample here
url: '[URL]',
// Use any name you like...
body: {'name': 'Susan'},
headers: {'Content-Type': 'application/json'},
json: true
};

function owCallback(err, response, body) {
if (!err && response.statusCode == 200) {
    console.log(response.body);
} else {
    console.log('The call failed! ' + err.message);
}
}

var request = require('request');
request.post(options, owCallback);

```

Show moreShow more icon

此沙箱的输出应该看起来很熟悉：

```
{ greeting: 'Hello, Susan' }

```

Show moreShow more icon

尽管结果相同，但上面两个沙箱的内容大相径庭。第一个沙箱包含可在浏览器中运行或部署到 OpenWhisk 的代码。第二个沙箱包含调用 OpenWhisk 操作的代码。可以将第二个沙箱中的 `url` 替换为任何公有 OpenWhisk 操作的 URL，添加该操作所需的 JSON 数据，并使用此代码调用它。

现在继续介绍实际编写可运行和部署的代码的详细信息。

### 编写可部署的 JavaScript 代码

在 JavaScript 中编写的 OpenWhisk 操作的基本结构类似于：

```
function main(params) {
// Do something useful with the parameters here
var output = {};
// Put anything you want your action to return into the output object
return output;
}

```

Show moreShow more icon

调用 JavaScript OpenWhisk 操作时，OpenWhisk 环境会寻找一个名为 `main()` 的方法。如果没有具有该名称的方法，则调用失败。在上面沙箱中的代码中，只需要 `main()` 函数即可部署该代码。对于此示例，向代码传递一个包含字段 `name` 的 JSON 对象，就大功告成了。

### 让 JavaScript 代码可运行

要让代码可运行，需要定义在有人尝试运行此文件时应该发生的行为。在我们的示例中，这意味着创建一个默认 JSON 对象并将它传递给 `main()` 方法。以下代码让它可运行：

```
var defaultParameters = {'name': 'Susan'};

if (require.main === module)
console.log(main(defaultParameters));

```

Show moreShow more icon

适当命名的变量 `defaultParameters` 定义一个 JSON 对象，其中包含在单击 **Run** 按钮时使用的默认参数。下一行更加复杂。如果语句 `require.main === module` 计算为 `true` ，意味着此文件在直接运行，而不是通过 `require()` 运行。如果遇到此情况，node.js 运行时会调用 `main()` 方法。（有关此技术的完整细节，Node.js 文档有 [一段有关访问 `main()` 方法的讨论](https://nodejs.org/api/modules.html#modules_accessing_the_main_module) ）。

在 `console.log()` 调用内调用 `main()` 方法意味着，从 `main()` 方法返回的所有结果都将输出到屏幕上。这适用于示例中的简单 JSON 对象，但您常常拥有更复杂的对象。在那时，您需要使用 `JSON.stringify()` 方法：

```
if (require.main === module)
console.log(JSON.stringify(main(defaultParameters), null, 2));

```

Show moreShow more icon

`JSON.stringify()` 方法是一个格式化 JSON 数据的便捷函数。这里的第一个参数是 `main()` 的实际调用，第二个参数是 `null` （我们可以传入一个函数来操作 JSON 数据，但我们不会这么做），第三个参数将数据缩进两个空格。在处理返回一个包含大量有用信息的复杂 JSON 对象的操作时，使用 `JSON.stringify()` 对代码进行格式化是值得的。

补充一点：您不需要声明 `defaultParameters` 对象。可以将 JSON 数据放在对 `main()` 的调用内：

```
if (require.main === module)
console.log(main({'name': 'Susan'}));

```

Show moreShow more icon

如果用于传入一个复杂的 JSON 对象，这么做是可行的，但不实用。

## 创建异步操作

OpenWhisk 支持异步操作，允许 JavaScript 代码返回一个 `Promise` 而不是基本 JSON 对象。在此示例中，该代码调用 [Watson Language Translator 服务](https://www.ibm.com/cloud/watson-language-translator) 。该服务接受一个字符串，并基于两段表明源语言和目标语言的代码来翻译它。

异步操作必不可少，因为您不希望 `main()` 方法在代码从 Watson 服务获得结果之前结束。通过返回 `Promise` ，只要有结果从服务返回来，客户端就会获得它们。此技术对长期运行的任务很有用，比如数据库访问或文件 I/O。（有关 `Promise` 的更多信息，包括一个允许编写您自己的 Promise 的沙箱，请参阅 [Ted Neward 介绍 ECMAScript 6 的新特性的教程。](https://www.ibm.com/developerworks/web/library/wa-ecmascript6-neward-p4/index.html) ）

```
function main(params) {
const LanguageTranslationV2 =
    require('watson-developer-cloud/language-translation/v2');

var url = params.url ||
            'https://gateway.watsonplatform.net/language-translator/api';
var use_unauthenticated =  params.use_unauthenticated || false;

const language_translator = new LanguageTranslationV2({
    'username': params.username,
    'password': params.password,
    'url' : url,
    'use_unauthenticated': use_unauthenticated
});

return new Promise(function (resolve, reject) {
    language_translator.translate({'text': params.textToTranslate,
                                   'source': 'en', 'target': 'es'},
                                  function(err, res) {
      if (err)
        reject(err);
      else
        resolve(res);
      });
});
}

const defaultParameters = {
'textToTranslate' : 'That that is not confusing is amazing.',
'username'        : '',
'password'        : '',
'url'             : 'https://sandbox-watson-proxy.mybluemix.net/language-translator/api',
'use_unauthenticated' : true
}

if (require.main === module)
main(defaultParameters)
    .then((results) => console.log(results))
    .catch((error) => console.log(error.message));

```

Show moreShow more icon

### 编写可部署的异步 JavaScript 代码

异步 OpenWhisk 操作的基本结构类似于：

```
function main(params) {
return new Promise(function (resolve, reject) {
    invokeSomething(params, function(error, results) {
      if (error)
        reject(error);
      else
        resolve(results);
    });
});
}

```

Show moreShow more icon

在上面沙箱中的代码中，对 WatsonLanguage Translator 服务的调用发生在 `Promise` 内。在某个时刻，该服务返回一个 `Error` 对象和一个包含结果的 JSON 对象。如果任何地方出错（ `Error` 对象不为 `null` ）， `reject()` 方法会返回 `Error` 对象。 否则， `resolve()` 方法返回该服务所返回的 JSON 对象。

_备注：_ 在上面的沙箱中， _从此页_ 调用此代码时， `username` 和 `password` 值可以是空的。要在部署此代码后调用它，可对 `defaultParameters` 对象执行以下更改：

1. 在 `username` 、 `password` 和 `url` 字段中填入用于 Language Translator 服务的 IBM Cloud 凭证。
2. 将 `use_unauthenticated` 字段设置为 `false` 。

参阅 [Language Translator 服务概述](https://www.ibm.com/cloud/watson-language-translator) ，进一步了解如何创建 Translator 服务的实例并生成它的凭证。

### 让异步 JavaScript 代码可运行

让代码可运行更加复杂，因为您需要处理 `Promise` 对象。跟以前一样，使用默认参数定义一个 JSON 对象，但是在调用 `main()` 后，必须使用 `then()` 和 `catch()` 方法：

```
const defaultParameters = {
'textToTranslate' : 'That that is not confusing is amazing.',
'username'        : '',
'password'        : '',
'url'             : 'https://sandbox-watson-proxy.mybluemix.net/language-translator/api',
'use_unauthenticated' : true
}

if (require.main === module)
main(defaultParameters)
    .then((results) => console.log(JSON.stringify(results, null, 2)))
    .catch((error) => console.log(error.message));

```

Show moreShow more icon

`then()` 方法处理 `Promise` 中的 `resolve()` 调用返回的对象。在本例中，这是包含该文本的译文的 JSON 对象。 `catch()` 方法处理错误情况，处理 `reject()` 方法返回的 `Error` 对象。

有关异步操作的完整细节，请参阅 [IBM Cloud Functions JavaScript 文档](https://cloud.ibm.com/docs/openwhisk) 并向下滚动到 “创建异步操作”小节。

## 一个可运行和部署的 Java 示例

现在继续看看 Java。下面是一个可运行和部署的简单示例：

```
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;

public class SimpleSample {
private static String data = "{'name': 'Susan'}";

public static void main(String[] args) {
    JsonObject jsonArgs = new JsonParser().parse(data).getAsJsonObject();
    main(jsonArgs);
}

public static JsonObject main(JsonObject params) {
    String greeting =
      "{'greeting': 'Hello, " + params.get("name").getAsString() + "'}";
    System.out.println(greeting);
    return new JsonParser().parse(greeting).getAsJsonObject();
}
}

```

Show moreShow more icon

### 测试代码二

与我们处理上面的 JavaScript 示例的方式一样，在谈论如何使代码可运行和部署之前，让我们先测试它。单击上面的 **Run** 按钮。根据针对 `name` 字段传入的值，您会看到类似这样的结果：

```
{'greeting': 'Hello, Susan'}

```

Show moreShow more icon

现在单击 **Deploy** 按钮。沙箱将代码作为新操作部署到 OpenWhisk，并显示它的 URL：

```
https://openwhisk.ng.bluemix.net/api/v1/web/devworkssandbox_pub/270003KAD5/0-1003.json

```

Show moreShow more icon

（同样地，您将看到的 URL 对您是唯一的。）

现在我们将调用刚部署到 OpenWhisk 的代码。 在下面的沙箱中，将 `[URL]` 替换为上面的 **Deploy** 按钮所生成的 URL。 按下 **Run** 按钮时，此代码会调用您在上面的 Java 沙箱中部署代码时创建的 OpenWhisk 操作。

```
import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URL;

public class OpenWhiskClient {
    public static void main(String[] args) {
      try {
          URL url = new URL("[URL]");
          HttpURLConnection hConn = (HttpURLConnection) url.openConnection();
        hConn.setRequestProperty("Content-Type", "application/json");
          hConn.setRequestMethod("POST");
        hConn.setDoOutput(true);

          String input = "{\"name\": \"Susan\"}";

          OutputStream os = hConn.getOutputStream();
          os.write(input.getBytes());
        os.flush();
        os.close();

        if (hConn.getResponseCode() == HttpURLConnection.HTTP_OK) {
            BufferedReader br =
          new BufferedReader(new InputStreamReader((hConn.getInputStream())));

          String output;
            while ((output = br.readLine()) != null) {
                System.out.println(output);
            }
      }
    } catch (Exception e) {
          e.printStackTrace();
      }
}
}

```

Show moreShow more icon

此沙箱的输出应看起来很熟悉：

```
{
"greeting": "Hello, Susan"
}

```

Show moreShow more icon

### 编写可部署的 Java 代码

OpenWhisk 要求 Java 类拥有以下基本结构：

```
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;

// Other imports here

public class MyClass {
public static JsonObject main(JsonObject args) {
    // Parse the args JsonObject so we can work with it

    // Use the Google gson library to create a JsonObject with the
    // appropriate output
    JsonObject returnObject = ...;
    return returnObject;
}
}

```

Show moreShow more icon

_备注：_ 请注意，IBM Developer 沙箱不支持 Java `package` 语句。

`main()` 方法的签名与您过去使用的签名不同。因为 OpenWhisk 需要一个 `JsonObject` 作为输入和输出，所以您需要使用这个签名。OpenWhisk 环境安装了著名的 Google `gson` 库，所以这里使用它来处理 JSON 数据。在 JavaScript 中处理 JSON 更容易，因为 JSON 已内置于该语言中。对于 Java，您需要通过 `gson` 来处理输入和输出数据。

### 让 Java 代码可运行

上面第一个沙箱中的示例既可部署，又可运行。让此代码可运行的代码行创建了一些默认数据，并通过创建一个包含传统签名的新方法来重载 `main()` 方法：

```
private static String data = "{'name': 'Susan'}";

public static void main(String[] args) {
JsonObject jsonArgs = new JsonParser().parse(data).getAsJsonObject();
main(jsonArgs);
}

```

Show moreShow more icon

`main(String[])` 方法根据默认数据创建一个 `JsonObject` ，并调用 `main(JsonObject)` 方法。 就这么简单。

## 补充示例：一个更复杂的可运行和部署的 Java 示例

作为补充示例，下面是之前使用 WatsonLanguage Translator 服务的示例的 Java 版本。请注意，Watson Java SDK 允许使用 Java 集合类和各种 Java 对象（例如 `TranslationResult` 和 `Translation` ）作为基础 JSON 的包装器。这使代码看起来更像常规 Java 程序。返回 JSON 数据时，一行代码将来自该服务的结果转换为 JSON 结构。

```
import java.util.List;

import com.google.gson.JsonObject;
import com.google.gson.JsonParser;

import com.ibm.watson.developer_cloud.language_translator.v2.LanguageTranslator;
import com.ibm.watson.developer_cloud.language_translator.v2.model.Language;
import com.ibm.watson.developer_cloud.language_translator.v2.model.Translation;
import com.ibm.watson.developer_cloud.language_translator.v2.model.TranslationResult;

public class LanguageTranslation {

private static String data =
    "{\"textToTranslate\": \"That that is not confusing is amazing.\"," +
    " \"username\"       : \"\"," +
    " \"password\"       : \"\"," +
    " \"endpoint\"       : \"https://sandbox-watson-proxy.mybluemix.net/language-translator/api\"," +
    " \"skip_authentication\" : \"true\"}";

public static void main(String[] args) {
    JsonParser parser = new JsonParser();
    JsonObject jsonArgs = parser.parse(data).getAsJsonObject();
    main(jsonArgs);
}

public static JsonObject main(JsonObject args) {
    JsonParser parser = new JsonParser();

    LanguageTranslator service = new LanguageTranslator();

    String username = args.get("username").getAsString();
    String password = args.get("password").getAsString();
    service.setUsernameAndPassword(username, password);

    if (args.get("endpoint") != null)
        service.setEndPoint(args.get("endpoint").getAsString());

       if (args.get("skip_authentication") != null)
      service.setSkipAuthentication((args.get("skip_authentication")
        .getAsString() == "true") ? true : false);

    String textToTranslate = args.get("textToTranslate").getAsString();

    TranslationResult result = service.translate(
      textToTranslate, Language.ENGLISH, Language.SPANISH).execute();

    List<Translation> translations = result.getTranslations();
    if (translations.size() > 1) {
      System.out.println("There are multiple translations:");
      for (Translation nextTranslation : translations) {
        System.out.println("- " +
          nextTranslation.getTranslation());
      }
    }
    else {
      System.out.println("The default translation: ");
      System.out.println("- " + result.getFirstTranslation());
    }

    JsonObject returnObject = parser.parse(result.toString())
                                    .getAsJsonObject();
    return returnObject;
}
}

```

Show moreShow more icon

可以看到，这个示例使用了前面介绍的技术，所以您能够运行和部署此代码。（Java 代码比 node.js 版本更简单，因为 Java 是同步的。Java 运行时会耐心地等待 Language Translator 服务返回一个结果。）

要在部署此代码后调用它，请对 `data` 对象执行以下更改：

1. 在 `username` 、 `password` 和 `endpoint` 字段中填入用于 Language Translator 服务的 IBM Cloud 凭证。
2. 将 `skip_authentication` 字段设置为 `false` 。

## 结束语

我们介绍了如何编写任何人都可在 IBM Developer 沙箱中运行和部署的 JavaScript 和 Java 代码。

- 如果将 OpenWhisk 操作部署到您自己的 IBM Cloud 帐户，可以使用 [非常有帮助的 `wsk` 命令行工具](https://cloud.ibm.com/functions/learn/cli?cm_sp=ibmdev-_-developer-articles-_-cloudreg) 来调用、管理和监视您部署的任何 OpenWhisk 操作。

您可以创建自己的 [免费 IBM Cloud 帐户](https://cloud.ibm.com/?cm_sp=ibmdev-_-developer-articles-_-cloudreg) ，将本教程中的代码部署到 OpenWhisk。这是一种进一步学习的不错方式，而且 `wsk` 工具和 IBM Cloud 控制台提供了一全套调试功能。

尽情享受无服务器计算的世界吧！