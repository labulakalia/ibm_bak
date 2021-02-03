# 快速构建 DMS Web 应用程序之利刃: Flask + Celery + Redis + Cloudant
快速搭建、二次开发 DMS Web 平台

**标签:** Python,Web 开发,云计算,数据库

[原文链接](https://developer.ibm.com/zh/articles/wa-build-dms-web-application-using-flask-celery-redis-cloudant/)

王芳, 卢屹, 杨卓

发布: 2018-10-24

* * *

## 系统介绍

本 DMS（Data Migration System）Web 系统采用 Python 语言开发，适用于跨平台跨版本系统之间业务数据的迁移。同时依托 IBM Cloud 平台，能够提供高可靠性和便捷性，使短时间内搭建系统，业务数据快速迁移成为可能。

本系统主要分为以下四大模块：

- **Web UI** ：作为与用户的直接联系，是系统中非常必要的一个模块，本系统采用 Flask 框架实现将后台传回的数据渲染至前端页面。
- **测试模块：** 本系统通过 Jupyter 对测试代码进行测试。
- **数据处理模块** ：因为 Pandas 对不同类型的数据文件如 `.csv` 、 `.txt` 、 `.pickle` 都有强大的支持，本系统采用 Pandas 对数据进行处理。
- **数据存储模块** ：本系统采用 Salesforce 对数据进行存储，由于 Salesforce 的写入时间较长，且数据量较大，为了更好的用户体验，本系统不采用将数据直接写入 Salesforce 的方法，而是通过 Flask 将信息发送给 Redis 和 Cloudant，使用 Redis 作为缓存存储临时变量，并将持久化变量存储在 Cloudant 中，当数据量较小的时候可以直接将数据写入 Salesforce 中。基于 Flask 的机制，数据需要先序列化存储之后，再反序列化才能使用，本系统采用 Pickle 将数据进行序列化及反序列化。通过 Celery 对消息队列以及暂存数据进行任务调度，并最后将数据写入 Salesforce 中完成服务端的通信任务。

下图是一个基于 Flask 的 DMS Web 系统的架构图：

##### 图 1\. 系统架构图

![图 1. 系统架构图](../ibm_articles_img/wa-build-DMS-web-application-using-Flask-Celery-Redis-Cloudant_images_image001.png)

## 相关技术简介

### Flask

Flask 是一个基于 Python 的轻量级 Web 框架，WSGI 工具箱采用 Werkzeug，模板引擎使用 Jinja2。由于其不依赖于特殊的工具或库，并且没有数据抽象层、表单验证或是其他任何已有多种库可以胜任的功能，从而保持核心简单、易于扩展，而被定义为”微”框架。但是，Flask 可以通过扩展来添加应用功能。并且 Flask 具有自带开发用服务器和 debugger、集成单元测试和 RESTful 请求调度 (request dispatching)、支持 secure cookie 的特点。

Jinja2 是基于 Python 的模版引擎，支持 Unicode，具有集成的沙箱执行环境并支持选择自动转义。Jinja2 拥有强大的自动 HTML 转移系统，可以有效的阻止跨站脚本攻击；通过模版继承机制，对所有模版使用相似布局；通过在第一次加载时将源码转化为 Python 字节码从而加快模版执行时间。

### Redis

Redis 是一个使用 ANSIC 语言编写、遵守 BSD 协议、Key-Value 的存储系统。拥有支持数据持久化、支持 string、map、list、set、sorted set 等数据结构和支持数据备份的特点。

Redis 会周期性地把更新的数据写入磁盘或把修改操作写入追加的记录文件，并且在此基础上实现主从（master-slave）同步，因此数据可以从主服务器向任意数量的从服务器上同步，从服务器可以是关联其他从服务器的主服务器。而且由于 Redis 完全实现了发布/订阅机制，使得从数据库在任何地方同步树时，可订阅一个频道并接收主服务器完整的消息发布记录。

### MQ

MQ 消息队列是一种应用程序的通信方法，应用程序可通过读写出入对立的消息进行通信。MQ 是一种消费者-生产者 （Producer-Customer）模式的实现。生产者-消费者模式由生产者、消费者和缓存区三个模块构成。缓存区作为一个中介的存在，生产者将数据放入缓存区，消费者从缓存区取出数据。本系统中，Flask 作为生产者，Salesforce 作为消费者，而 MQ 则是中间的缓存区。应用生产者-消费者模式能够有效的降低两者之间的耦合，减少互相之间的依赖；由于缓存区的存在，消费者无需直接从生产者处获取数据，能够支持并发任务、减少阻塞。

RabbitMQ 则是由 erlang 开发的 AMQP（高级消息队列协议）的开源实现，作为一个消息队列管理工具与 Celery 集成后，负责处理服务器之间的通信任务。RabbitMQ 的使用过程如下：

1. 客端连接到消息队列服务器并打开一个 channel。
2. 客户端声明一个 exchange、一个 queue，并分别设置相关属性。
3. 客户端使用 routing key 在 exchange 与 queue 之间绑定好关系。
4. 客户端投递消息到 exchange，exchange 根据消息的 key 和设置好的 binding，将消息投递到队列中。

RabbitMQ 常用的 Exchange Type 有以下三种：

- Fanout：能够将所有发送到该 exchange 的消息投递到所有与它绑定的队列中。
- Direct：把消息投递到那些 binding key 与 routing key 完全匹配的队列中。
- Topic：将消息路由到 binding key 与 routing key 模式匹配的队列中。

### Celery

Celery 是一个基于 Python 的分布式异步消息任务队列，通过它可以实现任务的异步处理，支持使用任务队列的方式在分布的机器、进程、线程上执行任务调度。并且 Celery 拥有任务失败或执行中断时自动尝试重新执行任务的高可用性以及处理任务的高效性。Celery 的架构是由消息中间件、任务执行单元和任务执行结果存储三部分构成的，所以在执行时需要配合一个消息中间件（Broker）来实现接收和发送任务消息，本文中消息中间件由 Redis + RabbitMQ 共同实现，RabbitMQ 作为 Broker，Redis 作为 Backend。图 2 是 Celery 的工作流程图：

##### 图 2\. Celery 工作流程图

![图 2. Celery 工作流程图](../ibm_articles_img/wa-build-DMS-web-application-using-Flask-Celery-Redis-Cloudant_images_image002.png)

### Cloudant

IBM Cloudant NoSQL Database for Bluemix（简称 Cloudant）在整个系统中作为数据库存在，是以 Apache CouchDB 为基础开发的多租户、独立、安装的服务。Cloudant 是一个 NoSQL 数据库即服务（DBaaS）。该服务是一个运营数据存储，用于存储 JSON 文档并对其建立索引，并且可以处理并行的读写操作。客户可以通过 HTTP API 从应用程序获取（get）、放置（put），查询（query）JSON 文档并对其建立索引（index）。

### Pickle

Pickle 是 Python 的标准模块之一，无需单独安装，主要有以下两类接口：

- 序列化：此操作能够将程序中运行的对象信息保存到文件中，并永久存储。
- 反序列化：此操作能够从文件中创建上一次程序保存的对象。

模块化后的对象以二进制的形式保存在文件中，后缀为 `.pkl` 。

## 快速搭建项目环境

了解了不同技术在系统中充当的角色，我们现在就可以开始安装搭建项目了。

### 搭建 Python 环境

因为 Flask 是基于 Python 的微框架，所以首先要保证计算机拥有 Python 环境。在电脑终端，输入 Python 查看是否有 Python 环境。如若出现版本号，则已安装 Python，若 import 报错，则还未安装 Python。如果是 macOS 则自带 Python 2.7。

### 搭建 Flask 环境

#### 安装 Virtualenv

Python 有多个版本，而且不同的版本间是不兼容的，而不同项目中使用不同版本的 Python 可能会造成外部库不兼容。所以我们就需要根据项目的要求选择对应的 Python 版本进行安装。virtualenv 对 Python 各个版本都能够提供很好的支持，我们可以快捷方便地根据项目的需求进行选择安装。

- **Linux/macOS X** ：通过 `$ sudo easy_install virtualenv` 或使用 pip `$ sudo pip install virtualenv` 进行安装，但使用这种方法我们首先要保证已安装 pip 环境。
- **Ubuntu** ：使用 `$ sudo apt-get install python-virtualenv` 进行安装。

待虚拟环境搭建完毕后，即可创建 Flask 工作环境，如清单 1 所示：

##### 清单 1\. 创建 Flask 工作环境

```
$ mkdir my-project
$ cd my-project
$ virtualenv venv

```

Show moreShow more icon

当显示 `New python executable in venv/bin/python Installing distribute...done.` 则安装成功。当在某项目中工作时，执行以下命令来激活 virtualenv：

```
$ . venv/bin/activate (Mac/Linux)
$ venv\scripts\activate (Windows)

```

Show moreShow more icon

#### 安装 Flask

激活 virtualenv 后，通过命令 `$ pip install Flask` 就可以安装并激活 Flask 了。

#### 搭建一个 Flask Demo

##### 清单 2\. 一个 Flask Demo

```
from flask import Flask
app = Flask(__name__)
@app.route('/')
def hello_world():
    return 'Hello World!'
if __name__ == '__main__':
    app.run()

```

Show moreShow more icon

保存为 `hello.py` ，在终端运行 `$ python hello.py` ，然后在浏览器运行 `localhost:5000/` 便可以启动了。也可以通过 `app.run(host='0.0.0.0)` 修改调用方法。

### 安装使用 Redis

由于数据量很大，我们需要一个缓存来对临时变量进行存储，可以使用 Redis 实现这个愿景。Redis 的安装也可以直接通过命令行实现：

#### macOS

通过 `$ brew install redis` 安装并用 `$ brew services start redis` 或 `$ redis-server /usr/local/etc/redis.conf` 来启动。

#### Linux

通过以下命令安装：

```
$ wget http://download.redis.io/releases/redis-2.8.17.tar.gz
$ tar xzf redis-2.8.17.tar.gz
$ cd redis-2.8.17
$ make

```

Show moreShow more icon

并用以下命令启动：

```
$ cd src
$ ./redis-server

```

Show moreShow more icon

#### Windows

在 [`https://github.com/MSOpenTech/redis/releases`](https://github.com/MSOpenTech/redis/releases) 下载安装包，将 `Redis-x64-xxx.zip` 压缩包到 C 盘，解压后，将文件夹重新命名为 `redis` 。打开一个 CMD 窗口使用命令切换目录到 `C:\redis` 运行 `redis-server.exe redis.windows.conf` 。

为方便起见，可以把 Redis 的路径加到系统的环境变量里，这样就省得再输路径了，后面的那个 `redis.windows.conf` 可以省略，如果省略，会启用默认的。这时候另启一个 CMD 窗口。(注意，原来的不要关闭，不然就无法访问服务端了。) 切换到 redis 目录下运行: `redis-cli.exe -h 127.0.0.1 -p 6379` 。

设置键值对： `set myKey abc`

取出键值对： `get myKey`

#### Ubuntu

通过以下命令安装：

```
$sudo apt-get update
$sudo apt-get install redis-server

```

Show moreShow more icon

并通过以下命令启动：

`$ redis-server`

安装好后，我们需要对 Redis 进行配置，Redis 默认配置文件 `redis.conf` 在根目录下，我们可以通过 Redis 的 config 命令来获取和设置所有的 Redis 配置： `$ redis 127.0.0.1:6379> CONFIG GET CONFIG_SETTING_NAME` 。

配置好后，就差最后一步即可使用 Redis 了，我们要将它引入到程序中来，通过 `iredis.Redis(host=url, port=port, db=0)` 将 Redis 连接至本地数据库，最后通过 `redis_client = redis.from_url(REDIS_URL)` 连接 RPC。

### 安装使用 MQ

有了缓存，还需要通信，本系统中通过消息队列实现通信。通过 `$ brew install rabbitmq` 实现安装。安装好后进入到安装目录，并执行 `$ sbin/rabbitmq-server` 来启动 MQ 的服务。然后在 Celery 中配置 Broker URL，如何配置我们会在 Celery 中详细讨论。

### 安装使用 Celery

有了消息队列，还需要一个角色去调度这些队列，Celery 就可以做到这个效果。通过命令 `$ sudo easy_install Celery` 或 `$ sudo pip install Celery` 进行安装。

同样，Celery 安装好后也需要配置：

- 当 Celery 默认使用的 Broker 时，RabbitMQ 的配置如下：

    `app.conf.broker_url = 'amqp://guest:guest@localhost:5672'`

- 当 Celery 使用 Redis 做 Broker 先通过以下命令安装 Redis 组件：

    `$$ pip install -U celery-with-redis`


然后进行配置：

`app.conf.broker_url = 'rediss://admin:password@host.composedb.com:39379'`

配置好后就可以创建调度任务了，步骤如下：

首先导入 Celery： `from celery import Celery` ，然后通过 `app.config.update` 写入配置，如图 3 所示：

##### 图 3\. 定义 Celery

![图 3. 定义 Celery](../ibm_articles_img/wa-build-DMS-web-application-using-Flask-Celery-Redis-Cloudant_images_image003.png)

然后，创建 Celery 任务，如清单 3 所示：

##### 清单 3\. 创建 Celery 任务

```
@celery.task(bind=True)
def run(x):
return x
@celery.task(bind=True)
def power(x):
return x*x

```

Show moreShow more icon

然后，启动 Celery Worker：

`$ celery -A name.py worker -l info`

最后在 Python 脚本中调用 Celery 中定义的函数就可以运行 Python 脚本：

```
from project.tasks import run, power
run.delay(args)

```

Show moreShow more icon

### 快速使用 Cloudant

一个项目中除了临时变量，还有一些需要持久化保存的数据，这个时候就需要 Cloudant 来对持久化数据进行存储了。

使用 Cloudant 需要先创建一个 IBM Cloud 账户，并将账户连接到 IBM ID，而后可以开始构建数据库。在浏览器中输入网址: [`https://cloud.ibm.com/`]即可进入 IBM Cloud 注册页面，如图 4 所示：

##### 图 4\. IBM Cloud 入口

![图 4. IBM Cloud 入口](../ibm_articles_img/wa-build-DMS-web-application-using-Flask-Celery-Redis-Cloudant_images_image004.png)

进入官网后创建 IBM Cloud 账户，便可以看到账户下的资源信息，如图 5 所示：

##### 图 5\. 仪表盘文件

![图 5. 仪表盘文件](../ibm_articles_img/wa-build-DMS-web-application-using-Flask-Celery-Redis-Cloudant_images_image005.png)

然后点击 **目录** ，进入 **Databases** \> **Cloudant** 。

##### 图 6\. Cloudant 入口

![图 6. Cloudant 入口](../ibm_articles_img/wa-build-DMS-web-application-using-Flask-Celery-Redis-Cloudant_images_image006.png)

拥有 IBM Cloud 的账号后，安装 Python 库来使用 Cloudant： `$ pip freeze` ，查看安装列表中是否安装了 Cloudant 2.8.1，若没有，则运行命令： `$ pip install cloudant==2.8.1` 。然后在 Python 文件中引入 Cloudant: `import cloudant` 。

此时便可以将 Cloudant 服务接入 IBM Cloud，首先登录至 IBM Cloud，并将 Cloudant NoSQL DB 服务添加到 IBM Cloud 空间，然后打开仪表盘，点击 **进入服务** ，点击左侧 **服务凭证** 添加凭证，如图 7 所示：

##### 图 7\. 添加凭证页面

![图 7. 添加凭证页面](../ibm_articles_img/wa-build-DMS-web-application-using-Flask-Celery-Redis-Cloudant_images_image007.png)

在脚本中进行基本认证： `client = Cloudant(USERNAME, PASSWORD, url=url, connect=True)`

然后创建数据库：

1. 在 Python 脚本中定义数据库名： `databaseName = "<DatabaseName>"` 。
2. 创建数据库： `myDatabase = client.create_database(databaseName)` 。
3. 通过以下代码验证数据库是否创建成功：





    ```
    if myDatabase.exists():
    print "'{0}' successfully created.\n".format(databaseName)

    ```





    Show moreShow more icon


如果成功，则可以对数据进行保存并进行一些操作，如：

- 定义一组数据，如清单 4 所示。
- 通过 Python 逐行获取数据，如清单 5 所示。

##### 清单 4\. 定义一组数据

```
sampleData = [
[1, "one", "boiling", 100],
[2, "two", "hot", 40],
[3, "three", "warm", 20],
[4, "four", "cold", 10],
[5, "five", "freezing", 0]
]

```

Show moreShow more icon

##### 清单 5\. 通过 Python 逐行获取数据

```
# Create documents by using the sample data.
# Go through each row in the array
for document in sampleData:
# Retrieve the fields in each row.
number = document[0]
name = document[1]
description = document[2]
temperature = document[3]

# Create a JSON document that represents
# all the data in the row.
jsonDocument = {
     "numberField": number,
     "nameField": name,
     "descriptionField": description,
     "temperatureField": temperature
}

# Create a document by using the database API.
newDocument = myDatabase.create_document(jsonDocument)

# Check that the document exists in the database.
if newDocument.exists():
     print "Document '{0}' successfully created.".format(number)

```

Show moreShow more icon

### 使用 Pickle

处理数据，必然会产生一些中间数据。Web 应用中，不应当将数据存储于 session，这样会导致内存不足问题。将数据序列化并存储到硬盘是一个很好的办法。Pickle 正是 Python 最好的序列化反序列化利器。因 Python 安装时已经集成 Pickle，所以直接引用即可： `import pickle` 。

Pickle 拥有自己的序列化和反序列化 API：

- 存储数据： `pickle.dump(obj, file,[,protocol])`
- 读取字符串： `pickle.load(file)`

Pandas 对 Pickle 的操作也有简单的 API：

- 存储数据： `df.to_pickle('test.pkl')`
- 读取数据： `df.read_pickle('test.pkl')`

## 创建测试脚本

一个完整的系统往往会比较大，更改代码时，可能会牵一发而动全身，所以需要打断点来排除错误。通过创建测试脚本，可以先确定代码的可行性，后移植到工程之中减少 debug 的工作量。下面我们就对 Jupyter 进行安装创建：

首先，安装 Jupyter notebook。macOS 自带 Python 情况下： `$ pip install –-user jupyter` 。

然后，启动 Jupyter： `$ jupyter notebook` 。或通过 anaconda-Navigator 启动，进入程序后，点击 Jupyter notebook 的 **Launch** 按钮即可启动。

##### 图 8\. Jupyter 登录页面

![图 8. Jupyter 登录页面](../ibm_articles_img/wa-build-DMS-web-application-using-Flask-Celery-Redis-Cloudant_images_image008.png)

启动 Jupyter 后，浏览器中的页面会显示当前目录下所有的文件，通过上传（Upload）和新建 （New）进行文件操作。点击 **New** 可以新建一个 `.ipynb` 文件。

##### 图 9\. Jupyter 新建文件页

![图 9. Jupyter 新建文件页](../ibm_articles_img/wa-build-DMS-web-application-using-Flask-Celery-Redis-Cloudant_images_image009.png)

如图 10 所示，进入界面即可添加测试代码，然后点击 **Run** 运行。

##### 图 10\. Jupyter 测试

![图 10. Jupyter 测试](../ibm_articles_img/wa-build-DMS-web-application-using-Flask-Celery-Redis-Cloudant_images_image010.png)

## 处理数据

本项目中因涉及大量的数据操作，采用 Pandas 对数据进行处理。Pandas 是 Python 的一个数据分析包，可用来进行数据分析，并提供了很多标准库和 API，能够快速地处理数据和函数，例如进行数据读取、过滤、排序、分组、选取、清洗等操作。

通过 `$pip install pandas` 命令安装 Pandas，并通过 `import pandas as pd` 引入。Pandas 对文件和数据都有良好的支持 API，如：

- 读取文件： `pd.read_excel(filename); pd.read_pickle(filename); pd.read_csv(filename)`
- 存储文件： `pd.to_csv(filename); pd.to_pickle(filename); pd.to_excel(filename)`

## 存储数据

数据处理完毕后，我们要将数据存入 Salesforce。可以通过 Simple Salesforce 来实现：

1. 运行 `$ sudo pip intall simple_salesforce` 命令安装 Simple Salsesforce。
2. 在在脚本中引入 Salsesforce： `from simple_salesforce import Salesforce` 并进行认证配置。
3. 通过 API 生成 `sf = Salesforce(instance_url=cred['instance_url'], session_id='')` 并使用 `sf.apexecute` 引入 Apex 语言 (一种面向对象的程序设计语言，可在大多数进程中执行编程功能)，即可将数据存储进此对象中。
4. 使用 `sf` 对象保存数据： `sf.Contact.create({'LastName':'Smith','Email':'example@example.com'})` 。

## 部署项目

本项目部署在 IBM Cloud 上。首先登录 IBM Cloud，进入仪表盘界面。点击 **目录** ，出现一系列可使用服务，选择 **Python** 入口，填写信息后，点击 **Create** 按钮，创建应用环境。一段时间后，仪表盘中则会出现新创建的应用环境。之后，则可以上传应用至 IBM Cloud。

##### 免费试用 IBM Cloud

利用 [IBM Cloud Lite](https://cloud.ibm.com/registration?cm_sp=ibmdev-_-developer-articles-_-cloudreg) 快速轻松地构建您的下一个应用程序。您的免费帐户从不过期，而且您会获得 256 MB 的 Cloud Foundry 运行时内存和包含 Kubernetes 集群的 2 GB 存储空间。 [了解所有细节](https://www.ibm.com/cloud/free/) 并确定如何开始。

##### 图 11\. IBM Cloud 项目目录

![图 11. IBM Cloud 项目目录](../ibm_articles_img/wa-build-DMS-web-application-using-Flask-Celery-Redis-Cloudant_images_image011.png)

##### 图 12\. 创建 Cloud Foundry 应用程序

![图 12. 创建 Cloud Foundry 应用程序](../ibm_articles_img/wa-build-DMS-web-application-using-Flask-Celery-Redis-Cloudant_images_image012.png)

上传应用，需要用到 IBM Cloud CLI，可以在刚刚创建的应用入门页面下载。

1. 首先要通过命令告诉 Cloud Foundry 所使用的 URL： `$ bluemix api https://api.ng.bluemix.net` 。
2. 然后通过命令登录至 IBM Cloud 环境。使用命令 `$ bluemix login -u username -o organizationname` 打开项目所在根目录，确保当前目录有 `manifest.yml` 文件。
3. 然后使用命令 `$ bluemix app push` 上传应用并根据 `manifest.yml` 文件中所写的内容进行配置。

`manifest.yml` 文件应放置在工程的根目录下，如清单 6 所示，文件内容通过 YAML 缩进格式来表现层级关系：

##### 清单 6\. Manifest.yml 示例

```
---
command: pip install --upgrade pip
applications:
- name: GMTDEV
memory: 512M
instances: 1
command: gunicorn run:app --workers=4
- name: GMTCeleryDEV
memory: 512M
instances: 1
no-route: true
command: celery -A tasks.celery worker --loglevel=info
health-check-type: process

```

Show moreShow more icon

上传完成后，可以在仪表盘中查看上传状态。显示下图效果，则此时已部署成功：

##### 图 13\. 应用程序管理仪表盘

![图 13. 应用程序管理仪表盘](../ibm_articles_img/wa-build-DMS-web-application-using-Flask-Celery-Redis-Cloudant_images_image013.png)

## 结束语

本文介绍了如何通过 Flask + Celery + Redis + Cloudant 等快速构建一个 Web 工程的全过程。首先介绍了 Flask、Redis、MQ、Celery、Cloudant、Pandas 等技术的基本概念及其安装使用方式，然后介绍了项目的整体技术架构，并完整展示了一个项目的搭建流程，最后展示了如何在 IBM Cloud 上部署应用。希望这篇文章能够帮助您快速、高效地构建一个 DMS Web 系统。