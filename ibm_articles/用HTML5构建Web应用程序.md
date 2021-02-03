# 用 HTML 5 构建 Web 应用程序
创建未来的 Web 应用程序

**标签:** Web 开发,移动开发

[原文链接](https://developer.ibm.com/zh/articles/wa-html5webapp/)

Michael Galpin

发布: 2010-10-13

* * *

## 简介

基于 HTML 5 已经涌现出了很多新的特性和标准。一旦您发现了当今浏览器中的一些可用的特性，就可以在您的应用程序中充分利用这些特性。在本文中，通过开发一些示例应用程序来了解如何探寻和使用最新的 Web 技术。本文的绝大多数代码都是 HTML、JavaScript 和 CSS — Web 开发人员的核心技术。

### 立即开始

为了能跟随我们的示例进行学习，最重要的一件事情是要有多个浏览器供测试使用。建议使用 Mozilla Firefox、Apple Safari 和 Google Chrome 最新版本。本文中我们使用的是 Mozilla Firefox 3.6、Apple Safari 4.04 和 Google Chrome 5.0.322。您可能还希望在移动浏览器上进行测试。例如，最新的 Android 和 iPhone SDK 被用于在模拟机上测试它们的浏览器。

您可以 [下载](http://public.dhe.ibm.com/software/dw/web/wa-html5webapp/FutureWeb.zip) 本文中所使用的源代码。

这些示例包括非常小的用 Java™ 写的后端组件。本文使用的是 JDK 1.6.0\_17 和 Apache Tomcat 6.0.14。有关下载这些工具的链接，请参见参考资源 。

## 探测的能力

有一个关于 Web 开发人员的笑谈，Web 开发人员用他们 20% 的时间写代码，然后用剩下的 80% 的时间让这些代码能够在所有浏览器中运转起来并获得相同的效果。要说 Web 开发人员已习惯于处理不同浏览器间的差别多少有点不切实际。随着浏览器新一轮创新浪潮的出现，这种低效的方法仍然没有改善。最新、最好的浏览器所支持的特性在不断变化。

不过，也有好的一面，这些新特性均集中在 Web 标准，这就让您能够现在就开始使用这些新特性。您可以采用渐进增强的老技巧、提供一些基线特性、查找高级特性，然后用出现的额外特性增强您的应用程序。要做到这一点，我们需要学习如何探测新功能。清单 1 显示了一个简单的探测脚本。

##### 清单 1\. 探测脚本

```
function detectBrowserCapabilities(){
    $("userAgent").innerHTML = navigator.userAgent;
    var hasWebWorkers = !!window.Worker;
    $("workersFlag").innerHTML = "" + hasWebWorkers;
    var hasGeolocation = !!navigator.geolocation;
    $("geoFlag").innerHTML = "" + hasGeolocation;
    if (hasGeolocation){
        document.styleSheets[0].cssRules[1].style.display = "block";
        navigator.geolocation.getCurrentPosition(function(location) {
            $("geoLat").innerHTML = location.coords.latitude;
            $("geoLong").innerHTML = location.coords.longitude;
        });
    }
    var hasDb = !!window.openDatabase;
    $("dbFlag").innerHTML = "" + hasDb;
    var videoElement = document.createElement("video");
    var hasVideo = !!videoElement["canPlayType"];
    var ogg = false;
    var h264 = false;
    if (hasVideo) {
        ogg = videoElement.canPlayType('video/ogg; codecs="theora, vorbis"') || "no";
           h264 = videoElement.canPlayType('video/mp4;
codecs="avc1.42E01E, mp4a.40.2"') || "no";
    }
    $("videoFlag").innerHTML = "" + hasVideo;
    if (hasVideo){
        var vStyle = document.styleSheets[0].cssRules[0].style;
        vStyle.display = "block";
    }
    $("h264Flag").innerHTML = "" + h264;
    $("oggFlag").innerHTML = "" + ogg;
}

```

Show moreShow more icon

目前已经出现了大量新特性和标准，成为了 HTML 5 标准的一部分。本文重点将放在几个最有用的特性上。 [清单 1](#清单-1-探测脚本) 中的脚本探测到了四个新特性：

- Web worker（多线程）
- 地理定位
- 数据库存储
- 本地视频回放

这个脚本的开头显示了用户浏览器的用户代理。它通常是一个惟一标识此浏览器的字符串，但它很容易被篡改。对于这个应用程序它已经足够好了。下一步是开始检测特性。首先要通过在全局范围（视窗）中查找 `Worker` 函数来检测 Web worker。这里用到了一些符合语言习惯的 JavaScript：双重否定。如果 `Worker` 函数不存在，那么 `window.Worker` 的求值结果为未定义，这是 JavaScript 中的一个 “伪” 值。如果在它前面放上一个单否定，就会被求值为真，因此放上一个双否定将被求值为假。检测完该值后，脚本会通过修改清单 2 中显示的 DOM 结构来将这个评估结果显示在屏幕上。

##### 清单 2\. 检测 DOM

```
<input type="button" value="Begin detection"
    onclick="detectBrowserCapabilities()"/>
<div>Your browser's user-agent: <span id="userAgent">
</span></div>
<div>Web Workers? <span id="workersFlag"></span></div>

<div>Database? <span id="dbFlag"></span></div>
<div>Video? <span id="videoFlag"></span></div>
<div class="videoTypes">Can play H.264? <span id="h264Flag">
</span></div>
<div class="videoTypes">Can play OGG? <span id="oggFlag">
</span></div>
<div>Geolocation? <span id="geoFlag"></span></div>
<div class="location">
    <div>Latitude: <span id="geoLat"></span></div>
    <div>Longitude: <span id="geoLong"></span></div>
</div>

```

Show moreShow more icon

清单 2 是一个简单的 HTML 结构，用来显示由检测脚本收集到的诊断信息。如 [清单 1](#清单-1-探测脚本) 中所示，下面要检测的是地理定位。这里再次使用了双重否定，但这次您要检测一个名为 `geolocation` 的对象，它应该是 `navigator` 对象的一个属性。如果找到该对象，那么就用 `geolocation` 对象的 `getCurrentPosition` 函数来获取当前的位置。获取位置可能是一个很慢的过程，因为通常会涉及到扫描 Wi-Fi 网络。在移动设备上，可能还会涉及到扫描无线发射塔和 ping GPS 卫星。由于花费时间很长，因此 `getCurrentPosition` 是异步的，并会将 `callback` 函数作为一个参数。在本例中，我们为 `callback` 使用一个闭包，它显示位置字段（通过启用其 CSS），然后将经度和纬度写到此 DOM。

下一步是检测数据库存储。检测是否有全局函数 `openDatabase` ，这个函数用于创建和访问客户端数据库。

最后，检测本地视频回放。用 DOM API 创建一个视频元素。今天的任何浏览器都能创建这样一个元素。在较老的浏览器中，这会是一个有效的 DOM 元素，但它没有任何的特别含义。这就如同是创建一个名为 `foo` 的元素。在一个较现代的浏览器中，它将是一个专用元素，就像是创建一个 `div` 或 `canPlayType` 元素。它将具有一个名为 `canPlayType` 的函数，所以只检测它的存在就可以了。

即使一个浏览器已经有了本地视频回放的功能，但它所支持的视频类型或它能回放的编解码并不是标准化了的。最容易想到的是要检测这个浏览器所支持的编解码。目前没有编解码的标准列表，但是有两个最常见的是 H.264 和 Ogg Vorbis。为了检测对特定的一个编解码的支持，可以向 `canPlayType` 函数传递一个识别字符串。如果浏览器能支持这个编解码，该函数将返回 `probably` （说真的 — 不是开玩笑的）。如果不支持，该函数将返回 null。在这个检测代码中，只针对这些值进行了检测并将结果显示在此 DOM 中。在一些常用浏览器中测试过此代码后，就会出现清单 3 中所示的综合结果。

##### 清单 3\. 不同浏览器的功能

```
#Firefox 3.6
Your browser's user-agent: Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10.6;
en-US; rv:1.9.2) Gecko/20100115 Firefox/3.6
Web Workers? true
Database? false
Video? true
Can play H.264? no
Can play OGG? probably
Geolocation? true
Latitude: 37.2502812
Longitude: -121.9059866

#Safari 4.0.4
Your browser's user-agent: Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_2;
en-us) AppleWebKit/531.21.8 (KHTML, like Gecko) Version/4.0.4 Safari/531.21.10
Web Workers? true
Database? true
Video? true
Can play H.264? probably
Can play OGG? no
Geolocation? false

#Chrome 5.0.322
Your browser's user-agent: Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_2;
en-US) AppleWebKit/533.1 (KHTML, like Gecko) Chrome/5.0.322.2 Safari/533.1
Web Workers? true
Database? true
Video? true
Can play H.264? no
Can play OGG? no
Geolocation? false

```

Show moreShow more icon

上面所列的所有常用桌面浏览器所支持的特性都不少。

- Firefox 惟一不支持的特性是数据库。对于视频，它只支持 Ogg。
- Safari 惟一不支持的特性是地理定位。
- Chrome 惟一不支持的特性是地理定位，尽管它声称不支持 H.264 或 Ogg。这可能是一个 bug，问题或是出在用于此次测试的 Chrome 的构建，或是出在检测代码。Chrome 实际上是支持 H.264 的。

地理定位在桌面浏览器上并不受广泛支持，但它在移动浏览器上却是受广泛支持的。清单 4 显示了移动浏览器的综合测试结果。

##### 清单 4\. 移动浏览器

```
#iPhone 3.1.3 Simulator
Your browser's user-agent: Mozilla/5.0 (iPhone Simulator; U; CPU iPhone OS 3.1.3
like Mac OS X; en-us) AppleWebKit/528.18 (KHTML, like Gecko)
Version/4.0 Mobile/7E18 Safari/528.16
Web Workers? false
Database? true
Video? true
Can play H.264? maybe
Can play OGG? no
Geolocation? true
Latitude: 37.331689
Longitude: -122.030731

#Android 1.6 Emulator
Your browser's user-agent: Mozilla/5.0 (Linux; Android 1.6; en-us;
sdk Build/Donut) AppleWebKit/528.5+ (KHTML, like Gecko) Version/3.1.2
Mobile Safari/525.20.1
Web Workers? false
Database? false
Video? false
Geolocation? false

#Android 2.1 Emulator
Your browser's user-agent: Mozilla/5.0 (Linux; U; Android 2.1; en-us;
sdk Build/ERD79) AppleWebKit/530.17 (KHTML, like Gecko) Version/4.0
Mobile Safari/530.17
Web Workers? true
Database? true
Video? true
Can play H.264? no
Can play OGG? no
Geolocation? true
Latitude:
Longitude:

```

Show moreShow more icon

上述代码中显示了最新的 iPhone 模拟器之一以及两种 Android。Android 1.6 不支持我们上述的这些检测。但实际上只要使用 Google Gear 它就能支持除视频之外的所有这些特性。它们就相当于是 API（就功能而言），但它们并不符合 Web 标准，因此会得到 [清单 4](#清单-4-移动浏览器) 中所显示的结果。将它与 Android 2.1 做个对照，后者支持所有这些特性。

请注意，iPhone 惟一不支持 Web worker。 [清单 3](#清单-3-不同浏览器的功能) 显示出 Safari 的桌面版本支持 Web worker，因此有理由相信这个特性不久也将会出现在 iPhone 中。

知道了该如何检测用户浏览器的这些特性之后，现在，让我们来探究一个简单的应用程序，这个应用程序将会综合使用这些特性 — 这取决于用户浏览器能处理什么。我们要构建的这个应用程序的功能是使用 Foursquare API 搜索某用户所在地周边的热点场所。

## 构建未来的应用程序

这个例子的重点是如何在移动设备上使用地理定位，但请记住 Firefox 3.5+ 也支持地理定位。这个应用程序首先查找用户当前位置附近的称为 _场所_ 的 Foursquare。场所可以是任何东西，但通常是指饭馆、酒吧、商店等。作为一个 Web 应用程序，我们的示例也受限于目前所有浏览器均执行的同源策略。它不能直接调用 Foursquare 的 API。而是使用一个 Java servlet 来实际代理这些调用。之所以采用 Java 并没有任何特别之处；您也可以用 PHP、Python、Ruby 等轻松编写一个类似的代理。清单 5 显示了一个代理 servlet。

##### 清单 5\. Foursquare 代理 servlet

```
public class FutureWebServlet extends HttpServlet {
    protected void doGet(HttpServletRequest request,
            HttpServletResponse response) throws ServletException, IOException {
        String operation = request.getParameter("operation");
        if (operation != null && operation.equalsIgnoreCase("getDetails")){

            getDetails(request,response);
        }
        String geoLat = request.getParameter("geoLat");
        String geoLong = request.getParameter("geoLong");
        String baseUrl = "http://api.foursquare.com/v1/venues.json?";
        String urlStr = baseUrl + "geolat=" + geoLat + "&geolong=" + geoLong;
        PrintWriter out = response.getWriter();
        proxyRequest(urlStr, out);
    }

    private void proxyRequest(String urlStr, PrintWriter out) throws IOException{
        try {
            URL url = new URL(urlStr);
            InputStream stream = url.openStream();
            BufferedReader reader = new BufferedReader( new InputStreamReader(stream));
            String line = "";
            while (line != null){
                line = reader.readLine();
                if (line != null){
                    out.append(line);
                }
            }
            out.flush();
            stream.close();
        } catch (MalformedURLException e) {
            e.printStackTrace();
        }
    }

    private void getDetails(HttpServletRequest request, HttpServletResponse response)
throws IOException{
        String venueId = request.getParameter("venueId");
        String urlStr = "http://api.foursquare.com/v1/venue.json?vid="+venueId;
        proxyRequest(urlStr, response.getWriter());
    }
}

```

Show moreShow more icon

这里需要注意的重要一点是代理了两个 Foursquare API。一个用来搜索，另一个用来获得这个场所的细节。要辨别这二者，细节 API 添加了一个操作参数。此外，将返回类型指定为 JSON，这会使解析来自 JavaScript 的数据变得十分简单。知道了应用程序代码所能进行哪种类型的调用之后，接下来让我们看看它如何进行这些调用以及如何使用来自 Foursquare 的数据。

### 使用地理定位

第一个调用是一个搜索。 [清单 5](#清单-5-foursquare-代理-servlet) 显示了对于纬度和经度需要两个参数：`geoLat` 和 `geoLong`。如下的清单 6 显示了如何获得应用程序内的这两个参数以及如何调用这个 servlet。

##### 清单 6\. 用位置进行搜索

```
if (!!navigator.geolocation){
    navigator.geolocation.getCurrentPosition(function(location) {
        venueSearch(location.coords.latitude, location.coords.longitude);
    });
}
var allVenues = [];
function venueSearch(geoLat, geoLong){
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if (this.readyState == 4 && this.status == 200){
            var responseObj = eval('(' + this.responseText + ')');
            var venues = responseObj.groups[0].venues;
            allVenues = venues;
            buildVenuesTable(venues);
        }
    }
    xhr.open("GET", "api?geoLat=" + geoLat + "&geoLong="+geoLong);
    xhr.send(null);
}

```

Show moreShow more icon

上述代码会查找浏览器的地理定位功能。如果浏览器有此功能，就会获得这个位置并用纬度和经度调用 `venueSearch` 函数。此函数使用了 Ajax（一个 `XMLHttpRequest` 对象调用 [清单 5](#清单-5-foursquare-代理-servlet) 内的 servlet）。它为 `callback` 函数使用一个闭包、解析来自 Foursquare 的 JSON 数据并将一个由场所对象组成的数组传递给一个名为 `buildVenuesTable` 的函数，如下所示。

##### 清单 7\. 使用场所构建 UI

```
function buildVenuesTable(venues){
    var rows = venues.map(function (venue) {
        var row = document.createElement("tr");
        var nameTd = document.createElement("td");
        nameTd.appendChild(document.createTextNode(venue.name));
        row.appendChild(nameTd);
        var addrTd = document.createElement("td");
        var addrStr = venue.address + " " + venue.city + "," + venue.state;
        addrTd.appendChild(document.createTextNode(addrStr));
        row.appendChild(addrTd);
        var distTd = document.createElement("td");
        distTd.appendChild(document.createTextNode("" + venue.distance));
        row.appendChild(distTd);
        return row;
    });
    var vTable = document.createElement("table");
    vTable.border = 1;
    var header = document.createElement("thead");
    var nameLabel = document.createElement("td");
    nameLabel.appendChild(document.createTextNode("Venue Name"));
    header.appendChild(nameLabel);
    var addrLabel = document.createElement("td");
    addrLabel.appendChild(document.createTextNode("Address"));
    header.appendChild(addrLabel);
    var distLabel = document.createElement("td");
    distLabel.appendChild(document.createTextNode("Distance (m)"));
    header.appendChild(distLabel);
    vTable.appendChild(header);
    var body = document.createElement("tbody");
    rows.forEach(function(row) {
        body.appendChild(row);
    });
    vTable.appendChild(body);
    $("searchResults").appendChild(vTable);
    if (!!window.openDatabase){
        $("saveBtn").style.display = "block";
    }
}

```

Show moreShow more icon

[清单 7](#清单-7-使用场所构建-ui) 内的代码主要是用来创建内含场所信息的数据表的 DOM 代码。但其中也不乏一些亮点。请注意高级 JavaScript 特性（比如数组对象的映射以及 `forEach` 函数）的使用。这里还有在支持地理定位的所有浏览器上均可用的一些特性。最后的两行代码也很有趣。对数据库支持的检测在此执行。如果支持，那么就启用一个 Save 按钮，用户单击此按钮就可将所有该场所的数据保存到一个本地数据库。下一节将讨论这一点是如何实现的。

### 结构化存储

[清单 7](#清单-7-使用场所构建-ui) 展示了典型的渐进增强策略。这个示例测试了是否有数据库支持。如果有此支持，就会添加一个 UI 元素以便向这个应用程序添加一个新特性供其使用。在本例中，它启用了一个按钮。单击此按钮会调用函数 `saveAll` ，如清单 8 所示。

##### 清单 8\. 保存到数据库

```
var db = {};
function saveAll(){
    db = window.openDatabase("venueDb", "1.0", "Venue Database",1000000);
    db.transaction(function(txn){
        txn.executeSql("CREATE TABLE venue (id INTEGER NOT NULL PRIMARY KEY, "+
                "name NVARCHAR(200) NOT NULL, address NVARCHAR(100),
cross_street NVARCHAR(100), "+
                "city NVARCHAR(100), state NVARCHAR(20), geolat TEXT NOT NULL, "+
                "geolong TEXT NOT NULL);");
    });
    allVenues.forEach(saveVenue);
    countVenues();
}
function saveVenue(venue){
    // check if we already have the venue
    db.transaction(function(txn){
        txn.executeSql("SELECT name FROM venue WHERE id = ?", [venue.id],
            function(t, results){
                if (results.rows.length == 1 && results.rows.item(0)['name']){
                    console.log("Already have venue id=" + venue.id);
                } else {
                    insertVenue(venue);
                }
            })
    });
}
function insertVenue(venue){
    db.transaction(function(txn){
        txn.executeSql("INSERT INTO venue (id, name, address, cross_street, "+
                "city, state, geolat, geolong) VALUES (?, ?, ?, ?, "+
                "?, ?, ?, ?);", [venue.id, venue.name,
                 venue.address, venue.crossstreet, venue.city, venue.state,
                 venue.geolat, venue.geolong], null, errHandler);
    });
}
function countVenues(){
    db.transaction(function(txn){
        txn.executeSql("SELECT COUNT(*) FROM venue;",[], function(transaction, results){
            var numRows = results.rows.length;
            var row = results.rows.item(0);
            var cnt = row["COUNT(*)"];
            alert(cnt + " venues saved locally");
        }, errHandler);
    });
}

```

Show moreShow more icon

要将场所数据保存到数据库，先要创建一个用来存储数据的表。创建表的语法是非常标准的 SQL 语法。（所有支持数据库的浏览器均使用 SQLite。查阅 SQLite 文档获得受支持的数据类型、限制等。）SQL 执行异步完成。此外，还会调用事务函数并向其传递一个回调函数。`callback` 函数获得一个事务对象，用来执行 SQL。`executeSQL` 函数接受一个 SQL 字符串，然后是一个可选的参数列表，外加成功和错误处理器函数。如果没有错误处理器，错误就被 “吃掉”。对于 `create table` 语句而言，这是一种理想状况。脚本首次执行时，此表将会被成功创建。当再次执行时，脚本将会失败，因为表已经存在 — 但这也问题不大。因为我们只需要确保在向表内插入行之前，此表已经存在。

此表创建后，通过 `forEach` 函数用从 Foursquare 返回的每个场所调用 `saveVenue` 函数。此函数先是通过查询这个场所来查证这个场所是否已经被存储在本地。这里，使用了一个成功处理器。查询的结果集将被传递给这个处理器。如果没有结果或场所尚未被存储于本地，就会调用 `insertVenue` 函数，由它执行一个插入语句。

借助 `saveAll`，在所有的保存/插入完成后，调用 `countVenues`。目的是查询插入到这个场所表内的行的总数。这里的语法（`row["COUNT(*)"]`）从查询的结果集中拉出该计数。

了解了如何使用数据库支持（如果有的话）后，接下来的一节将会探讨如何使用 Web worker 支持。

### Web worker 的后台处理

回到 [清单 6](#清单-6-用位置进行搜索) ，让我们对它进行稍许修改。如下面的清单 9 所示，检测是否有 Web worker 支持。如果有，就用它来获得从 Foursquare 检索来的每个场所的更多信息。

##### 清单 9\. 修改后的场所搜索

```
function venueSearch(geoLat, geoLong){
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if (this.readyState == 4 && this.status == 200){
            var responseObj = eval('(' + this.responseText + ')');
            var venues = responseObj.groups[0].venues;
            allVenues = venues;
            buildVenuesTable(venues);
            if (!!window.Worker){
                var worker = new Worker("details.js");
                worker.onmessage = function(message){
                    var tips = message.data;
                    displayTips(tips);
                };
                worker.postMessage(allVenues);
            }
        }
    }
    xhr.open("GET", "api?geoLat=" + geoLat + "&geoLong="+geoLong);
    xhr.send(null);
}

```

Show moreShow more icon

上述代码使用了与之前相同的检测方法。如果存在对 Web worker 支持，就会创建一个新的 worker。为了创建一个新的 worker，需要指向 worker 将要运行的脚本的 URL — 在本例中，即 details.js 文件。当此 worker 完成其作业后，它会向主线程发送回一个消息。 `onmessage` 处理器将接收这个消息；我们为它使用了一个简单的闭包。最后，为了初始化这个 worker，用一些数据调用 `postMessage` 以便它能工作起来。将所有检索自 Foursquare 的场所信息传递进来。清单 10 显示了 details.js 的内容，该脚本将由此 worker 执行。

##### 清单 10\. worker 的脚本 details.js

```
var tips = [];
onmessage = function(message){
    var venues = message.data;
    venues.foreach(function(venue){
        var xhr = new XMLHttpRequest();
        xhr.onreadystatechange = function(){
            if (this.readyState == 4 && this.status == 200){
                var venueDetails = eval('(' + this.responseText + ')');
                venueDetails.tips.forEach(function(tip){
                    tip.venueId = venue.id;
                    tips.push(tip);
                });
            }
        };
        xhr.open("GET", "api?operation=getDetails&venueId=" + venueId, true);
        xhr.send(null);
    });
    postMessage(tips);
}

```

Show moreShow more icon

这个细节脚本迭代所有场所。对于每个场所，此脚本会使用 `XMLHttpRequest` 调用 Foursquare 代理以便获得场景的细节。不过，请注意在使用它的 `open` 函数打开连接时，传递进第三个参数（`true`）。这会使调用成为同步的，而不再是通常的异步的。在 worker 中这么做是可以的，因为没有处于主 UI 线程，并且不会冻结整个应用程序。把它变成同步的，意味着一个调用必须结束后，下一个才能开始。处理程序只简单从场所细节中抽取 tips 并收集所有这些 tips 后一并传递回主 UI 线程。为了将此数据传递回去，调用了 `postMessage` 函数，由该函数在这个 worker 上调用 `onmessage` 回调函数，如 [清单 9](#清单-9-修改后的场所搜索) 中所示。

默认地，这个场所搜索返回 10 个场所。不难想象为了获得细节而进行 10 个额外的调用将要花费多长时间。因而使用 Web worker 在后台线程中完成这类任务是很有意义的。

## 结束语

本文介绍了现代浏览器内的一些 HTML 5 的新功能。您了解了如何检测新的特性以及如何将这些新特性添加到您的应用程序内。这些特性的绝大多数现在都已经广受流行浏览器的支持 — 特别是移动浏览器。现在您就可以充分利用地理定位和 Web worker 等特性来创建具有创新性的 Web 应用程序了。

本文翻译自： [Build Web applications with HTML 5](https://developer.ibm.com/articles/wa-html5webapp/)（2010-03-30）