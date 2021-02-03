# Java API 设计实践
让您的 API 可同时在模块化和非模块化 Java 环境中使用

**标签:** API 管理,Java,Java 平台

[原文链接](https://developer.ibm.com/zh/articles/api-design-practices-for-java/)

[Bentley Hargrave](https://developer.ibm.com/zh/profiles/hargrave)

发布: 2018-08-23

* * *

## 前言

了解在设计 Java API 时应该运用的一些 API 设计实践。这些实践通常很有用，而且可确保 API 能在诸如 OSGi 和 Java Platform Module System (JPMS) 之类的模块化环境中得到正确使用。有些实践是规定性的，有些则是禁止性的。当然，其他良好的 API 设计实践也同样适用。

OSGi 环境提供了一个模块化运行时，使用 Java 类加载器概念来强制实施类型 _可见性_ 封装。每个模块都将有自己的类加载器，该加载器将 _连接_ 到其他模块的类加载器，以共享导出的包并使用导入的包。

Java 9 引入了 JPMS，后者提供了一个模块化平台，使用来自 Java 语言规范的 [访问控制](https://docs.oracle.com/javase/specs/jls/se9/html/jls-6.html#jls-6.6) 概念来强制实施类型 _可访问性_ 封装。每个模块都定义了哪些包将被导出，从而可供其他模块访问。默认情况下，JMPS 层中的模块都位于同一个类加载器中。

一个包可以包含一个 API。这些 API 包的客户端有两种 _角色_： _API 使用者_ 和 _API 提供者_。API 使用者使用 API 提供者实现的 API。

在以下设计实践中，我们将讨论包的公共部分。包的成员和类型不是公共的，或者说是受保护的（即私有或默认可访问的），无法从包的外部访问它们，所以它们是包的实现细节。

## Java 包必须是一个有凝聚力、稳定的单元

Java 包的设计必须确保它是一个 _有凝聚力_ 且 _稳定_ 的单元。在模块化 Java 中，包是模块之间共享的实体。一个模块可以导出一个包，以便其他模块可以使用这个包。因为包是在模块之间共享的单元，所以它必须有凝聚力，因为包中的所有类型必须与包的特定用途相关。不鼓励使用混杂的包，比如 `java.util`，因为这种包中的类型通常是彼此不相关的。这类没有凝聚力的包可能导致大量依赖项，因为包的不相关部分会引用其他不相关的包，而且对包的某个方面的更改会影响依赖于此包的所有模块，即使模块可能并未实际使用该包的被修改部分。

由于包是一个共享单元，所以它的内容必须是众所周知的，而且随着包在未来版本中的演变，只能以兼容方式更改所包含的 API。这意味着包不能支持 API 超集或子集；例如，可以将 `javax.transaction` 视为内容不稳定的包。包的用户必须能够了解包中提供了哪些类型。这也意味着包应由单个实体（例如一个 jar
文件）提供，不得跨多个实体进行拆分，因为该包的用户必须知道整个包的存在。

此外，该包必须以兼容方式实现在未来版本中的演变。因此，应该对包进行版本控制，而且其版本号必须根据 [语义版本控制](https://semver.org/) 规则进行演变。还有一篇 [OSGi 白皮书: 语义版本控制](https://www.osgi.org/wiki/uploads/Links/SemanticVersioning.pdf)。

但最近我意识到，针对包的主版本更改的语义版本控制建议是错误的。包的演变必须是功能的增加。在语义版本控制中，这会增加次要版本。当删除功能时，您会对包执行不兼容的更改，而不是增加主版本，
您必须改用一个新的包名，而让原始包保持可兼容。为了理解这样做的重要性和必要性，请参阅这篇有关 [Go 语义导入版本控制](https://research.swtch.com/vgo-import) 的论文和 [Rich Hickey 在 Clojure/conj 2016](https://www.youtube.com/watch?v=oyLBGkS5ICk) 上发表的精彩主题演讲。二者都表明，在对包执行不兼容的更改时，应该改用新的包名，而不是更改主版本。

## 最小化包耦合

一个包中的类型可以引用其他包中的类型，例如，某个方法的参数类型和返回类型，以及某个字段的类型。这种包间耦合会在包上造成所谓的 _使用_ 限制。这意味着 API 使用者必须使用 API 提供者引用的相同包，以便他们都能了解所引用的类型。

一般而言，我们希望最小化这种包耦合，以便最小化包上的使用限制。这可以简化 OSGi 环境中的连接解析，并最大限度地减少依赖扇出，从而简化部署。

## 接口优先于类

对于 API，接口优先于类。这是一种相当常见的 API 设计实践，对模块化 Java 也很重要。接口的使用提高了实现的自由度，还支持多种实现。接口对于让 API 使用者与 API 提供者分离至关重要。无论是实现接口的 API 提供者，还是调用接口上的方法的 API 使用者，都允许使用包含 API 接口的包。通过这种方式，API 使用者不会直接依赖于 API 提供者。它们都只依赖于 API 包。

除接口外，抽象类有时也是一种有效的设计选择，但接口通常是首选，尤其是考虑到接口的最新改进允许添加 `default` 方法。

最后，API 通常需要许多小的具体类，比如事件类型和异常类型。这没什么问题，但这些类型通常应该是不可变的且不被API 使用者用来创建子类。

## 避免静态

应在 API 中要避免静态。类型不应包含静态成员。静态工厂也应该避免。实例的创建应与 API 分离。例如，API 使用者应通过依赖注入或者 OSGi 服务注册表或 jPMS 中的 `java.util.ServiceLoader` 之类的对象注册表来接收 API 类型的对象实例。

避免静态也是一种创建可测试的 API 的良好实践，因为静态不容易模仿。

## 单例

API 设计中有时存在单例对象。但是，不应通过静态 `getInstance` 方法或静态字段等静态对象来访问单例对象。当需要使用单例对象时，该对象应由 API 定义为单例，并通过上文提到的依赖注入或对象注册表提供给 API 使用者。

## 避免类加载器假设

API 通常具有可扩展性机制，API 使用者可在其中提供 API 提供者必须加载的类名。然后，API 提供者必须使用 `Class.forName`（可能使用线程上下文类加载器）来加载该类。这种机制会假定从 API 提供者（或线程上下文类加载器）到 API 使用者的类可见性。API 设计必须避免类加载器假设。模块化的一个主要特点是类型封装。一个模块（例如 API 提供者）不能对另一个模块（例如 API 使用者）的实现细节具有可见性/可访问性。

API 设计必须避免在 API 使用者与 API 提供者之间传递类名，而且必须避免与类加载器分层结构和类型可视性/可访问性有关的假设。为了提供可扩展性模型，API 设计应让 API 使用者将类对象，或者最好将实例对象，传递给 API 提供者。这可以通过 API 中的一个方法或 OSGi 服务注册表之类的对象注册表来完成。请参见 [白板模式](https://www.osgi.org/wiki/uploads/Links/whiteboard.pdf)。

在 JPMS 模块中不使用 `java.util.ServiceLoader` 类时，也会受到类加载器假设的影响，因为它假设所有提供者都对线程上下文类加载器或所提供的类加载器可见。此假设在模块化环境中通常是不成立的，但 JPMS 允许通过模块声明来声明模块提供或使用了一个
`ServiceLoader` 托管服务。

## 不进行持久性假设

许多 API 设计都假设只有一个构造阶段，对象在该阶段被实例化并添加到 API，但是忽略了动态系统中可能发生的解构阶段。API 设计应考虑到，对象可以添加，也可以删除。例如，大多数监听器 API 都允许添加和删除监听器。但是，许多 API 设计仅假设对象可以添加，并且从未删除这些对象。例如，许多依赖注入系统无法撤销注入的对象。

在 OSGi 环境中，可以添加和删除模块，因此能够兼顾到此类动态的 API 设计非常重要。 [OSGi Declarative Services 规范](https://osgi.org/specification/osgi.cmpn/7.0.0/service.component.html)
为 OSGi 定义了一个依赖注入模型，该模型支持这些动态操作，包括撤销注入对象。

## 明确规定 API 使用者和 API 提供者的类型角色

如前言中所述，API 包的客户端有两种角色：API 使用者和 API 提供者。API 使用者使用 API，而 API 提供者实现 API。对于 API 中的接口（和抽象类）类型，API 设计一定要明确规定哪些类型仅由 API 提供者实现，哪些类型可由 API 使用者实现。例如，监听器接口通常由 API 使用者实现，而实例被传递给 API 提供者。

API 提供者对 API 使用者和 API 提供者实现的类型更改都很敏感。提供者必须实现 API 提供者类型中的任何新更改，而且必须了解且可能调用 API 使用者类型中的任何新更改。API 使用者通常可以忽略 API 提供者类型的（兼容）更改，除非它希望通过更改来调用新功能。但是 API 使用者对 API 使用者类型的更改很敏感，而且可能需要修改才能实现新功能。例如，在 `javax.servlet` 包中，`ServletContext` 类型由 API 提供者（比如 servlet 容器）实现。向 `ServletContext` 添加新方法需要更新所有 API 提供者来实现新方法，但 API 使用者无需执行更改，除非他们希望调用该新方法。但是，`Servlet` 类型由 API 使用者实现，向 `Servlet` 添加新方法需要修改所有 API 使用者来实现新方法，还需要修改所有 API 提供者来使用该新方法。因此，`ServletContext` 类型有一个 API 提供者角色，`Servlet` 类型有一个 API 使用者角色。

由于通常有许多 API 使用者和很少的 API 提供者，所以在考虑更改 API 使用者类型时，必须非常谨慎地执行 API 演变，而 API 提供者类型更改的要求更加宽松。这是因为，您只需要更改少数 API 提供者来支持更新的 API，但您不希望在更新 API 时需要更改许多现有的 API 使用者。仅当 API 使用者希望使用新 API 时，API 使用者才需要执行更改。

OSGi Alliance 定义了 [文档注释](https://osgi.org/specification/osgi.core/7.0.0/framework.module.html#framework.module-semantic.versioning.type.roles)、 [`ProviderType`](https://osgi.org/specification/osgi.core/7.0.0/framework.api.html#org.osgi.annotation.versioning.ProviderType) 和 [`ConsumerType`](https://osgi.org/specification/osgi.core/7.0.0/framework.api.html#org.osgi.annotation.versioning.ConsumerType) 来标记 API 包中的类型角色。这些注释包含在 [osgi.annotation](http://search.maven.org/#artifactdetails%7Corg.osgi%7Cosgi.annotation%7C7.0.0%7Cjar) jar 中供您的 API 使用。

## 结束语

在您下次设计 API 时，请考虑这些 API 设计实践。这样，您的 API 就可以同时在模块化和非模块化的 Java 环境中使用。

本文翻译自: [API design practices for Java](https://developer.ibm.com/articles/api-design-practices-for-java/)（2018-08-23）