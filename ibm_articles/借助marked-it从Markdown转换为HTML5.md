# 借助 marked-it 从 Markdown 转换为 HTML5
使用开源 marked-it 项目将 Markdown 转换为 HTML5

**标签:** Web 开发,前端开发

[原文链接](https://developer.ibm.com/zh/articles/markdown-html5-marked-it/)

Jenifer Schlotfeldt

发布: 2020-09-02

* * *

Markdown 正在迅速成为创作技术内容的首选标记语言之一。Markdown 便于创作，在开发者社区中广泛使用且支持作者和 SME 之间的协作，因此对于内容创作者来说很有吸引力。

随着它越来越受欢迎，人们也越来越需要实现一些增强的创作功能。但是，在扩展 Markdown 以增强创作功能的同时不能破坏作者和合作者所青睐的 Markdown 的简单性。

我们也知道，大多数现代的文档应用程序都需要 HTML5 输出才能满足创建最佳用户体验的设计要求。利用 [marked-it](https://ibm.github.io/marked-it/)，IBM 构建了一个生成器，可将 Markdown 解析为 HTML5，并且处理团队为提供某些高级创作功能而创建的扩展，而不会让 Markdown 标记过度复杂化。

为了尽量保持 Markdown 的简单性，这些扩展可以挂钩到 HTML 生成过程中定义明确的点。这样一来，您就可以定义相应的方式以向元素中添加类、新的 Markdown 标记甚至是定义如何组织一组主题。在本文中，我向您介绍了 marked-it 中的一些关键元素，正是这些元素让 marked-it 成为可在网站上显示通过 Markdown 源生成的 HTML5 内容的有用工具。

## 类

marked-it 中的关键扩展与在 Markdown 中使用元数据属性的功能相关。 **属性** 是可以绑定到 Markdown 元素的名称或名称-值映射。然后，将属性传递到由此生成的 HTML5 元素，并作为类输出。属性通常用于确保生成的 HTML5 输出的 CSS 样式设置正确。在定义属性后，可以将这些值应用于任何 Markdown 元素，例如标题、段落和代码块。

您也可能需要定义版权和其他元数据，以在生成的 HTML5 输出中统一前置或后置。在此扩展中，您可以使用在转换过程中调用的 **标准页眉和页脚文件**。您可能还需要根据自身发布要求为文件添加 HTML 开头和结尾标签。

为了满足 **辅助功能要求**，还有一个扩展用于添加行标头、表摘要和图形标题。

## Markdown 标记

虽然 Markdown 可以包含 HTML，但理想情况下，作者可以完全使用 Markdown。IBM Cloud 团队定义了 marked-it 可以识别的新标记，以防基本 Markdown 标记不够用。

通常，当我们定义新标记时，我们会搜索其他团队已成功扩展基本 Markdown 的行业示例，并加以复制。Conref 是一个很好的例子，我们在此能够发现基于 DITA 标记中流行实现的现有 Markdown 实现，并将其引入 marked-it 中。

marked-it 中新标记的另一个例子适用于 **视频**。尽管编写者可以将 iframe 代码包含在 Markdown 文件中，但是仅将某些内容标记为视频并让生成器完成所有工作，是不是要容易得多？另外，它可以确保所有视频标记都是一致的。该视频扩展使作者可以使用图像 Markdown 标记对视频进行编码。然后，再将其定义为视频。

例如：`![视频标题](https://www.youtube.com/embed/<video-ID>){: video output="iframe"}`

## 组织主题

除了 marked-it 通过每个 Markdown 文件生成的标准 HTML5 输出，生成器还能够以不同格式生成 **目录文件**（TOC file）。为 Markdown 主题生成的每个 TOC 条目都包含一组嵌套的结构化标题（或 topicref），它们与原始 Markdown 源中标题的结构相匹配。在通过 Markdown 源生成的 HTML5 输出文件中，每个 TOC 标题（或 topicref）还包含一个指向对应标题锚定 ID 的链接。

这意味着您只需要按照希望 Markdown 文件在导航中出现的顺序来定义它们，marked-it 就会创建一个完全结构化的目录文件，该文件包含父 Markdown 文件及其中所有标题的条目。

## 谁在使用 marked-it？

在 IBM，以 Markdown 创作内容的例子包括 [IBM Cloud 文档](https://cloud.ibm.com/docs) 和 [IBM Developer](https://developer.ibm.com/zh/) 的内容。在这两种实现中，marked-it 都是自动构建流水线不可或缺的一部分，内容创建者可以进行编写、提交更改，而构建将自动触发以生成可以为各自应用程序所使用的 HTML5 输出。

想要详细了解 marked-it？查看我们的 [规范](https://ibm.github.io/marked-it/)。Marked-it 生成器是一个开源工具，期待大家贡献自己的一份力量！

本文翻译自： [Go from Markdown to HTML5 with marked-it](https://developer.ibm.com/articles/markdown-html5-marked-it/)（2020-08-17）