# 探索 Python、机器学习和 NLTK 库
开发一个应用程序，使用 Python、NLTK 和机器学习对 RSS 提要进行分类

**标签:** Python

[原文链接](https://developer.ibm.com/zh/articles/os-pythonnltk/)

Chris Joakim

发布: 2012-12-17

* * *

##### Python 的初学者

本文适用于面临其第一个机器学习实现的软件开发人员，尤其是那些具有 Ruby 或 Java 语言背景的开发人员。

## 挑战：使用机器学习对 RSS 提要进行分类

最近，我接到一项任务，要求为客户创建一个 RSS 提要分类子系统。目标是读取几十个甚至几百个 RSS 提要，将它们的许多文章自动分类到几十个预定义的主题领域当中。客户网站的内容、导航和搜索功能都将由这个每日自动提要检索和分类结果驱动。

客户建议使用机器学习，或许还会使用 Apache Mahout 和 Hadoop 来实现该任务，因为客户最近阅读了有关这些技术的文章。但是，客户的开发团队和我们的开发团队都更熟悉 Ruby，而不是 Java™ 技术。本文将介绍解决方案的技术之旅、学习过程和最终实现。

### 什么是机器学习？

我的第一个问题是，”究竟什么是机器学习？” 我听说过这个术语，并且隐约知道超级计算机 IBM® Watson 最近使用该技术在一场 Jeopardy 比赛中击败了人类竞争者。作为购物者和社交网络活动参与者，我也知道 Amazon.com 和 Facebook 根据其购物者数据在提供建议（如产品和人）方面表现良好。总之，机器学习取决于 IT、数学和自然语言的交集。它主要关注以下三个主题，但客户的解决方案最终仅涉及前两个主题：

- **分类：** 根据类似项目的一组训练数据，将相关的项分配到任意预定义的类别。
- **建议：** 根据类似项目的观察来建议采用的项。
- **集群：** 在一组数据内确定子组。

### Mahout 和 Ruby 的选择

理解了机器学习是什么之后，下一步是确定如何实现它。根据客户的建议，Mahout 是一个合适的起点。我从 Apache 下载了代码，并开始了学习使用 Mahout 及其兄弟 Hadoop 实现机器学习的过程。不幸的是，我发现即使对于有经验的 Java 开发人员而言，Mahout 的学习曲线也很陡峭，并且不存在可用的样例代码。同样不幸的是，机器学习缺乏基于 Ruby 的框架或 gem。

### 发现 Python 和 NLTK

我继续搜索解决方案，并且在结果集中一直遇到 “Python”。作为一名 Ruby 开发人员，虽然我还没有学过该语言，但我也知道 Python 是一个面向相似对象的、基于文本的、可理解和动态的编程语言。尽管两种语言之间存在一些相似之处，但我多年来都忽视了学习 Python，将它视为一项多余的技能集。Python 是我的 “盲点”，我怀疑许多 Ruby 开发人员同行都是这样认为的。

搜索机器学习的书籍，并更深入研究它们的目录，我发现，有相当高比例的此类系统在使用 Python 作为其实现语言，并使用了一个被称为 Natural Language Toolkit（NLTK，自然语言工具包）的库。通过进一步的搜索，我发现 Python 的应用比我意识到的还要广泛，如 Google App Engine、YouTube 和使用 Django 框架构建的网站。它甚至还预安装在我每天都使用的 Mac OS X 工作站上！此外，Python 为数学、科学和工程提供了有趣的标准库（例如，NumPy 和 SciPy）。

我决定推行一个 Python 解决方案，因为我找到了非常好的编码示例。例如，下面这一行代码就是通过 HTTP 读取 RSS 提要并打印其内容所需的所有代码：

```
print feedparser.parse("http://feeds.nytimes.com/nyt/rss/Technology")

```

Show moreShow more icon

## 快速掌握 Python

在学习一门新的编程语言时，最容易的部分往往是学习语言本身。较难的部分是了解它的生态系统：如何安装它、添加库、编写代码、构造代码文件、执行它、调试它并编写单元测试。本节将简要介绍这些主题；请务必参阅 参考资料 ，以获得有关详细信息的链接。

### pip

Python Package Index (`pip`) 是 Python 的标准软件包管理器。您可以使用该程序将库添加到您的系统。它类似于 Ruby 库的 gem。为了将 NLTK 库添加到您的系统，您可以输入以下命令：

```
pip install nltk

```

Show moreShow more icon

为了显示在您的系统上已安装的 Python 库的列表，请运行以下命令：

```
pip freeze

```

Show moreShow more icon

### 运行程序

执行 Python 程序同样很简单。获得一个名称为 _locomotive\_main.py_ 的程序和三个参数，然后您就可以使用 Python 程序编译并执行它：

```
python locomotive_main.py arg1 arg2 arg3

```

Show moreShow more icon

Python 使用 [清单 1](#清单-1-main-检测) 中的 `if __name__ == "__main__":` 语法来确定文件本身是从命令行执行的还是从其他代码导入的。为了让文件变得可以执行，需要添加 `"__main__"` 检测。

##### 清单 1\. Main 检测

```
import sys
import time
import locomotive

if __name__ == "__main__":
    start_time = time.time()
    if len(sys.argv) > 1:
        app = locomotive.app.Application()
        ... additional logic ...

```

Show moreShow more icon

### virtualenv

大多数 Ruby 开发人员熟悉系统范围的库或 gem 的问题。使用一组系统范围内的库的做法一般是不可取的，因为您的其中一个项目可能依赖于某个给定的库的版本 1.0.0，而另一个项目则依赖于版本 1.2.7。同样，Java 开发人员都知道系统范围的 CLASSPATH 存在同样的问题。就像 Ruby 社区使用其 `rvm` 工具，而 Python 社区使用 `virtualenv` 工具（请参阅 参考资料 ，以获得相关链接）来创建独立的执行环境，其中包含特定版本的 Python 和一组库。 [清单 2](#清单-2-使用-virualenv-创建一个虚拟环境的命令) 中的命令显示了如何为您 p1 项目创建一个名为 _p1\_env_ 的虚拟环境，其中包含 `feedparser` 、 `numpy` 、 `scipy` 和 `nltk` 库。

##### 清单 2\. 使用 virualenv 创建一个虚拟环境的命令

```
$ sudo pip install virtualenv
$ cd ~
$ mkdir p1
$ cd p1
$ virtualenv p1_env --distribute
$ source p1_env/bin/activate
(p1_env)[~/p1]$ pip install feedparser
(p1_env)[~/p1]$ pip install numpy
(p1_env)[~/p1]$ pip install scipy
(p1_env)[~/p1]$ pip install nltk
(p1_env)[~/p1]$ pip freeze

```

Show moreShow more icon

每次在一个 shell 窗口使用您的项目时，都需要 “获得” 您的虚拟环境激活脚本。请注意，在激活脚本被获得后，shell 提示符会改变。当在您的系统上创建和使用 shell 窗口，轻松地导航到您的项目目录，并启动其虚拟环境时，您可能想在您的 ~/.bash\_profile 文件中添加以下条目：

```
alias p1="cd ~/p1 ; source p1_env/bin/activate"

```

Show moreShow more icon

### 代码库结构

在完成简单的单文件 “Hello World” 程序的编写之后，Python 开发人员需要理解如何正确地组织其代码库的目录和文件名。Java 和 Ruby 语言在这方面都有各自的要求，Python 也没有什么不同。简单来说，Python 使用 _包_ 的概念对相关的代码进行分组，并提供了明确的名称空间。出于演示目的，在本文中，代码存在于某个给定项目的根目录中，例如 ~/p1。在这个目录中，存在一个用于相同名称的 Python 包的 locomotive 目录。 [清单 3](#清单-3-示例目录结构) 显示了这个目录结构。

##### 清单 3\. 示例目录结构

```
locomotive_main.py
locomotive_tests.py

locomotive/
    __init__.py
    app.py
    capture.py
    category_associations.py
    classify.py
    news.py
    recommend.py
    rss.py

locomotive_tests/
    __init__.py
    app_test.py
    category_associations_test.py
    feed_item_test.pyc
    rss_item_test.py

```

Show moreShow more icon

请注意名称古怪的 **\_init**.py _文件。这些文件指示 Python 为您的包加载必要的库和特定的应用程序代码文件，它们都位于相同的目录中。 [清单 4](#清单-4-locomotive-init-py) 显示了文件 locomotive/\_init_.py 的内容。

##### 清单 4\. locomotive/ **init**.py

```
# system imports; loads installed packages
    import codecs
    import locale
    import sys

    # application imports; these load your specific *.py files
    import app
    import capture
    import category_associations
    import classify
    import rss
    import news
    import recommend

```

Show moreShow more icon

有了结构如 [清单 4](#清单-4-locomotive-init-py) 所示的 locomotive 包之后，在项目的根目录中的主程序就可以导入并使用它。例如，文件 locomotive\_main.py 包含以下导入：

```
import sys         # >-- system library
    import time        # >-- system library
    import locomotive  # >-- custom application code library in the "locomotive" directory

```

Show moreShow more icon

### 测试

Python `unittest` 标准库提供一个非常好的测试解决方案。熟悉 JUnit 的 Java 开发人员和熟悉 Test::Unit 框架的 Ruby 开发人员应该会觉得 [清单 5](#清单-5-python-unittest) 中的 Python `unittest` 代码很容易理解。

##### 清单 5\. Python unittest

```
class AppTest(unittest.TestCase):

      def setUp(self):
          self.app = locomotive.app.Application()

      def tearDown(self):
          pass

      def test_development_feeds_list(self):
          feeds_list = self.app.development_feeds_list()
          self.assertTrue(len(feeds_list) == 15)
          self.assertTrue('feed://news.yahoo.com/rss/stock-markets' in feeds_list)

```

Show moreShow more icon

[清单 5](#清单-5-python-unittest) 中的代码还演示了 Python 的一个显著的特点：所有的代码必须一致缩进，否则无法成功编译。 `tearDown(self)` 方法可能在开始时看起来有点古怪。您可能会问，为什么测试总是被硬编码为通过？事实上并非如此。这只是在 Python 中编写空方法的一种方式。

### 工具

我真正需要的是一个具备语法突出显示、代码完成和断点调试功能的集成开发环境 (IDE)，用该环境帮助我掌握我的 Python 学习曲线。作为使用 Eclipse IDE 进行 Java 开发的一名用户， `pyeclipse` 插件是我考虑的下一个工具。虽然该插件有时比较慢，但它工作得相当不错。我最终投资了 PyCharm IDE，它满足了我的所有 IDE 要求。

在掌握了 Python 及其生态系统的基本知识之后，终于来到开始实现机器学习解决方案的时候。

## 使用 Python 和 NLTK 实现分类

实现解决方案涉及捕获模拟的 RSS 提要、整理其文本、使用一个 `NaiveBayesClassifier` 和 kNN 算法对类别进行分类。下面将会介绍这些操作中的每一个。

### 捕获和解析提要

该项目特别具有挑战性，因为客户还没有定义目标 RSS 提要列表。因此，也不存在 “训练数据”。所以，在初始开发期间必须模拟提要和训练数据。

我用来获得示例提要数据的第一个方法是只提取在某个文本文件中指定的列表中的 RSS 提要。Python 提供了一个很好的 RSS 提要解析库，其名称为 `feedparser` ，它抽象不同的 RSS 和 Atom 格式之间的差异。简单的基于文本的对象序列化的另一个有用的库被幽默地称为 `pickle` （泡菜）。这两个库在 [清单 6](#清单-6-capturefeeds-类) 的代码中均有使用，清单 6 中的代码将每一个 RSS 提要捕获为 “腌制过的” 对象文件，以备后用。如您所见，Python 代码非常简洁，且功能强大。

##### 清单 6\. CaptureFeeds 类

```
import feedparser
import pickle

class CaptureFeeds:

    def __init__(self):
        for (i, url) in enumerate(self.rss_feeds_list()):
            self.capture_as_pickled_feed(url.strip(), i)

    def rss_feeds_list(self):
        f = open('feeds_list.txt', 'r')
        list = f.readlines()
        f.close
        return list

    def capture_as_pickled_feed(self, url, feed_index):
        feed = feedparser.parse(url)
        f = open('data/feed_' + str(feed_index) + '.pkl', 'w')
        pickle.dump(feed, f)
        f.close()

if __name__ == "__main__":
    cf = CaptureFeeds()

```

Show moreShow more icon

下一步的挑战性之大是出乎意料的。现在，我有了样例提要数据，必须对它进行分类，以便将它用作训练数据。 _训练数据_ 是向您的分类算法提供的数据集，以便您能从中进行学习。

例如，我使用的样例提要包括了体育电视网络公司 ESPN。提要的项目之一是关于 Denver Broncos 橄榄球队的 Tim Tebow 被转会到 New York Jets 橄榄球队，在同一时间，Broncos 签了他们新的四分卫 Peyton Manning。提要结果中的另一个项目是 Boeing Company 和它的新喷气式飞机 (jet)。所以，这里的问题是，应该将哪些具体的类别值分配给第一个故事？ `tebow` 、 `broncos` 、 `manning` 、 `jets` 、 `quarterback` 、 `trade` 和 `nfl` 这些值都是合适的。但只有一个值可以在训练数据中被指定为训练数据类别。同样，在第二个故事中，类别应该是 `boeing` 还是 `jet` ？困难的部分在于这些细节。如果您的算法要产生精确的结果，那么大型训练数据集的准确手工分类非常关键。要做到这一点，不应该低估所需的时间。

我需要使用更多的数据，而且这些数据必须已进行了准确的分类，这种情况很快就变得明显。我可以在哪里找到这样的数据呢？进入 Python NLTK。除了是一个出色的语言文本处理库之外，它甚至还带有可下载的示例数据集，或是其术语中的 _文集_ ，以及可以轻松访问此下载数据的应用程序编程接口。要安装 Reuters 文集，可以运行如下所示的命令。会有超过 10,000 篇新闻文章将下载到您的 ~/nltk\_data/corpora/reuters/ 目录中。与 RSS 提要项目一样，每篇 Reuters 新闻文章中都包含一个标题和一个正文，所以这个 NLTK 预分类的数据非常适合于模拟 RSS 提要。

```
$ python               # enter an interactive Python shell
>>> import nltk        # import the nltk library
>>> nltk.download()    # run the NLTK Downloader, then enter 'd' Download
Identifier> reuters    # specify the 'reuters' corpus

```

Show moreShow more icon

特别令人感兴趣的是文件 ~/nltk\_data/corpora/reuters/cats.txt。它包含了一个列表，其中包含文章文件名称，以及为每个文章文件分配的类别。文件看起来如下所示，所以，子目录 test 中的文件 14828 中的文章与主题 `grain` 有关。

```
test/14826 trade
test/14828 grain

```

Show moreShow more icon

### 自然语言是混乱的

RSS 提要分类算法的原始输入，当然是以英语书写的文本。原始，确实如此。

从计算机处理的角度来看，英语或任何自然语言（口语或普通的语言）都是极不规范和不准确的。首先，存在大小写的问题。单词 _Bronco_ 是否等于 _bronco_ ？答案是，也许是。接下来，您要应付标点和空格。 _bronco._ 是否等于 _bronco_ 或 _bronco,_ ？算是吧。然后，有复数形式和相似的单词。 _run、running_ 和 _ran_ 是否相等？这取决于不同的情况。这三个词有一个共同的 _词根_ 。如果将自然语言词汇嵌入在标记语言（如 HTML）中，情况会怎么样呢？在这种情况下，您必须处理像 `<strong>bronco</strong>` 这样的文本。最后，还有一个问题，就是那些经常使用但基本上毫无意义的单词，像 _a、and_ 和 _the_ 。这些所谓的停用词非常碍事。自然语言非常凌乱；在处理之前，需要对它们进行整理。

幸运的是，Python 和 NLTK 让您可以收拾这个烂摊子。在 [清单 7](#清单-7-rssitem-类) 中， `RssItem` 类的 `normalized_words` 方法可以处理所有这些问题。请特别注意 NLTK 如何只使用一行代码就能够清洁嵌入式 HTML 标记的原始文章文本！使用一个正则表达式删除标点，然后每个单词被拆分，并规范化为小写。

##### 清单 7\. RssItem 类

```
class RssItem:
    ...
    regex = re.compile('[%s]' % re.escape(string.punctuation))
    ...
    def normalized_words(self, article_text):
        words   = []
        oneline = article_text.replace('\n', ' ')
        cleaned = nltk.clean_html(oneline.strip())
        toks1   = cleaned.split()
        for t1 in toks1:
            translated = self.regex.sub('', t1)
            toks2 = translated.split()
            for t2 in toks2:
                t2s = t2.strip().lower()
                if self.stop_words.has_key(t2s):
                    pass
                else:
                    words.append(t2s)
        return words

```

Show moreShow more icon

只需这一行代码就可以从 NLTK 获得停用词列表；并且还支持其他自然语言：

```
nltk.corpus.stopwords.words('english')

```

Show moreShow more icon

NLTK 还提供了一些 “词干分析器” 类，以便进一步规范化单词。请查看有关词干、词形归并、句子结构和语法的 NLTK 文档，了解有关的更多信息。

### 使用 Naive Bayes 算法进行分类

算法在 NLTK 中被广泛使用并利用 `nltk.NaiveBayesClassifier` 类实现。Bayes 算法根据特性在其数据集中的每个存在或不存在对项目进行分类。在 RSS 提要项目的情况下，每一个特性都是自然语言的一个给定的（清洁过的）单词。该算法是 “朴实” 的，因为它假设特性（在本例中，单词）之间没有任何关系。

然而，英语这种语言包含超过 250,000 个单词。当然，我不希望为了将 RSS 提要项目传递给算法就要为每个 RSS 提要项目创建一个包含 250,000 个布尔值的对象。那么，我会使用哪些单词？简单来说，答案是在培训数据组中除了停用词之外最常见的单词。NLTK 提供了一个优秀的类，即 `nltk.probability.FreqDist` ，我可以用它来识别这些最常用的单词。在 [清单 8](#清单-8-使用-nltk-freqdist-类) 中， `collect_all_words` 方法返回来自所有培训文章的所有单词的一个数组。

然后，此数组被传递给 `identify_top_words` 方法，以确定最频繁的单词。 `nltk.FreqDist` 类的一个有用的特性是，它实质上是一个散列，但是它的键按其对应的值或 _计数_ 排序。因此，使用 `[:1000]` Python 语法可以轻松获得最频繁的 1000 个单词。

##### 清单 8\. 使用 nltk.FreqDist 类

```
def collect_all_words(self, items):
      all_words = []
      for item in items:
          for w in item.all_words:
              words.append(w)
      return all_words

def identify_top_words(self, all_words):
      freq_dist = nltk.FreqDist(w.lower() for w in all_words)
      return freq_dist.keys()[:1000]

```

Show moreShow more icon

对于利用 NLTK Reuters 文章数据模拟的 RSS 提要项目，我需要确定每个项目的类别。为此，我读取前面提到的 ~/nltk\_data/corpora/reuters/cats.txt 文件。用 Python 读取一个文件非常简单，如下所示：

```
def read_reuters_metadata(self, cats_file):
      f = open(cats_file, 'r')
      lines = f.readlines()
      f.close()
      return lines

```

Show moreShow more icon

接下来的步骤是获得每个 RSS 提要项目的特性。 `RssItem` 类的 `features` 方法（如下所示）可以做到这一点。在该方法中，在文章中的 `all_words` 数组首先被减少到一个较小的 `set` 对象，以消除重复的单词。然后会遍历 `top_words` ，并在该 set 中进行比较，确定是否存在重复的单词。随后返回 1000 个布尔值组成的一个散列，以 `w_` 为键，后面是单词本身。这个 Python 非常简洁。

```
def features(self, top_words):
      word_set = set(self.all_words)
      features = {}
      for w in top_words:
          features["w_%s" % w] = (w in word_set)
      return features

```

Show moreShow more icon

接下来，我收集了训练集的 RSS 提要项目和它们各自的特性，并将它们传递给算法。 [清单 9](#清单-9-训练-nltk-naivebayesclassifier) 中的代码演示了这个任务。请注意，分类器被训练成为只有一行代码。

##### 清单 9\. 训练 nltk.NaiveBayesClassifier

```
def classify_reuters(self):
        ...
        training_set = []
        for item in rss_items:
            features = item.features(top_words)
            tup = (features, item.category)  # tup is a 2-element tuple
            featuresets.append(tup)
        classifier = nltk.NaiveBayesClassifier.train(training_set)

```

Show moreShow more icon

`NaiveBayesClassifier` 在运行中的 Python 程序的内存中，它现在是经过训练的。现在，我只需遍历需要进行分类的 RSS 提要项目集，并要求分类器猜测每个项目的类别。这很简单。

```
for item in rss_items_to_classify:
      features = item.features(top_words)
      category = classifier.classify(feat)

```

Show moreShow more icon

### 变得不那么朴实

如前所述，算法假设每个特性之间是没有关系的。因此，像 “machine learning” 和 “learning machine”，或者 “New York Jet” 和 “jet to New York” 这样的短语是等效的（ _to_ 是一个停用词）。在自然的语言上下文中，这些单词之间有明显的关系。所以，我怎么会让算法变得 “不那么天真”，并识别这些单词的关系？

其中一个技巧是在特性集内包括常见的 _双字词_ （两个单词为一组）和 _三字词_ （三个单词为一组）。NLTK 以 `nltk.bigrams(...)` 和 `nltk.trigrams(...)` 的形式对此提供了支持，现在我们对此应该不再感到惊讶了。正如可以从训练数据组收集最常用的 _n_ 个单词那样，也可以识别最常用的双字词和三字词，并将它们用作特性。

### 您的结果会有所不同

对数据和算法进行完善是一门艺术。您是否应该进一步规范化单词集，也许应该包括词根？或者包括超过 1000 个最常用单词？少一点是否合适？或者是否应该使用更大的训练数据集？是否应该添加更多信用词或 “停用词根”？这些都是您要问自己的正确问题。使用它们进行实验，通过试错法，您可以会为您的数据实现最佳算法。我发现，85% 是一个很好的分类成功率。

### 利用 k-Nearest Neighbors 算法提出建议

客户希望显示在选定类别或相似类别中的 RSS 提要项目。现在，这些项目已经用 Naive Bayes 算法进行分类，这一要求的第一部分已得到了满足。较难的部分是实现 “或相似类别” 的要求。这是机器学习建议器系统开始发挥作用的地方。 _建议器系统_ 根据其他项目的相似性来建议一个项目。Amazon.com 的产品建议和 Facebook 的朋友建议就是此功能的很好的示例。

k-Nearest Neighbors (kNN) 是最常用的建议算法。思路是向它提供一组 _标签_ （即 _类别_ ），并且每个标签都对应一个数据集。然后，该算法对各数据集进行了比较，以识别相似的项目。数据集由多个数值数组构成，数值的范围往往被规范化为从 0 到 1。然后，它可以从数据集识别相似的标签。与只产生一个结果的 Naive Bayes 不同，kNN 可以产生一个有排名的列表，其中包含若干（即， _k_ 的值）个建议。

我发现，建议器算法比分类算法更容易理解和实现，但对于本文来说，其代码过于冗长，并且有复杂的数学，无法在这里详述。请参阅由 Manning 出版的一本很好的新书 _Machine Learning in Action_ ，获取 kNN 编码示例（请参阅 参考资料 中的链接）。在 RSS 提要项目实现的过程中，标签值是项目类别，而数据集是最常用的 1000 个单词的值数组。同样，在构建这个数组时，一部分属于科学范畴，一部分属于数学范畴，还有一部分属于艺术范畴。在数组中，每个单词的值都可以是简单的 0 或 1 的布尔值、文章中单词出现次数的百分比、该百分比的指数值，或一些其他值。

## 结束语

探索 Python、NLTK 和机器学习一直是一个有趣的、令人愉快的经验。Python 语言强大而又简洁，现在已成为我的开发工具包的核心部分。它非常适合于机器学习、自然语言和数学/科学应用程序。虽然本文中并没有提到，但我还发现 Python 对于图表和绘图非常有用。如果 Python 同样是您的盲点，我建议您了解一下它。

本文翻译自： [Explore Python, machine learning, and the NLTK library](https://developer.ibm.com/articles/os-pythonnltk/)（2012-10-09）