# SWT、Swing 或 AWT：哪个更适合您？
在为新应用程序选择 GUI 工具包时应该考虑哪些因素？

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/os-swingswt/)

Barry Feigenbaum

发布: 2006-04-27

* * *

## 简介

IBM Developer 上另外一些作者已经展示了如何在 Swing 和 SWT 之间很好地进行迁移（参见参考资料 ）。本文的目标是帮助您在开始开发项目之前确定选择使用哪个 GUI 工具包。

但是首先我们要弄清一个问题：为什么会有多个 Java™ GUI 工具包呢？最好的答案是，一个工具包并不能满足所有的要求，最近也不会开发一个可以满足所有要求的 GUI 工具包。每个工具包都有各自的优缺点，这样就可以根据自己的需求和目标用户来选择适当的工具包。

下面就让我们来学习有关这些工具包的知识。

## AWT 概述

Abstract Windows Toolkit（AWT）是最原始的 Java GUI 工具包。AWT 的主要优点是，它在 Java 技术的每个版本上都成为了一种标准配置，包括早期的 Web 浏览器中的 Java 实现；另外它也非常稳定。这意味着我们不需要单独安装这个工具包，在任何一个 Java 运行环境中都可以使用它，这一点正是我们所希望的特性。

AWT 是一个非常简单的具有有限 GUI 组件、布局管理器和事件的工具包（参见 [清单 1](#清单-1-基本的-awt-class-树) 、 [清单 2](#清单-2-awt-提供了下面的布局管理器) 和 [清单 3](#清单-3-awt-提供了以下事件) ）。这是因为 Sun 公司决定为 AWT 使用一种最小公分母（LCD）的方法。因此它只会使用为所有 Java 主机环境定义的 GUI 组件。最终的结果非常不幸，有些经常使用的组件，例如表、树、进度条等，都不支持。对于需要更多组件类型的应用程序来说，我们需要从头开始创建这些组件。这是一个很大的负担。

##### 清单 1\. 基本的 AWT Class 树

全部在 java.awt 包中， `*` 表示抽象。

```
Object
    CheckboxGroup
    *Component
        Button
        Canvas
        CheckBox
        Choice
        Container
            Panel
                Applet
            ScrollPane
            Window
                Dialog
                Frame
        Label
        List
        TextComponent
            TextArea
            TextField
    MenuComponent
        MenuItem
            CheckboxMenuItem
            Menu
                PopupMenu

```

Show moreShow more icon

**注意** ：另外几个包中还有其他一些 AWT 组件，但是这是基本的组件集。

##### 清单 2\. AWT 提供了下面的布局管理器

全部在 java.awt 包中，`*` 表示接口。

```
*LayoutManager
    FlowLayout
    GridLayout
    *LayoutManager2
        BorderLayout
        CardLayout
        GridBagLayout

```

Show moreShow more icon

**注意** ：另外几个包中还有一些 AWT 布局管理器，很多都是为它们进行布局的容器专门定制的，但是这是基本的布局管理器集。

##### 清单 3\. AWT 提供了以下事件

大部分在 java.awt.events 包中。

```
Object
    EventObject
        AWTEvent
            ActionEvent
            AdjustmentEvent
            ComponentEvent
                ContainerEvent
                FocusEvent
                InputEvent
                    KeyEvent
                    MouseEvent
                        MouseWheelEvent
                PaintEvent
                WindowEvent
            HierarchyEvent
            InputMethodEvent
            InvocationEvent
            ItemEvent
            TextEvent

```

Show moreShow more icon

**注意** ：其他几个包中还有另外一些 AWT 事件，但是这是基本的事件集。这些是从更通用的事件生成的具体事件。

通常对于 AWT 来说（也适用于 Swing 和 SWT），每个事件类型都有一个相关的 _XxxListener_ 接口（ _XxxAdapter_ 的实现可能为空），其中 Xxx 是去掉 Event 后缀的事件名（例如，KeyEvent 事件的接口是 KeyListener），用来把事件传递给处理程序。应用程序会为自己感兴趣处理的事件的事件源（GUI 组件或部件）进行注册。有时监听接口要处理多个事件。

AWT 的一个很好的特性是它通常可以对 GUI 组件自动进行销毁。这意味着您几乎不需要对组件进行销毁。一个例外是高级组件，例如对话框和框架。如果您创建了耗费大量主机资源的资源，就需要手动对其进行销毁。

AWT 组件是 “线程安全的（thread-safe）”，这意味着我们不需要关心在应用程序中是哪一个线程对 GUI 进行了更新。这个特性可以减少很多 GUI 更新的问题，不过使 AWT GUI 运行的速度更慢了。

AWT 让我们可以以 _自顶向下（top-down）_ 或 _自底向上（bottom-up）_ 或以任意组合顺序来构建 GUI。自顶向下的意思是在创建子组件之前首先创建容器组件；自底向上的意思是在创建容器（或父）组件之前创建子组件。在后一种情况中，组件的存在并不依赖于父容器，其父容器可以随时改变。

通常来说，AWT GUI 都是不可访问的。系统并没有为 AWT 程序员提供 API 来指定可访问性信息。可访问性（accessibility）处理的是残疾人可以怎样使用应用程序的问题。一个应用程序要想有很好的可访问性，必须与运行平台一起，让残疾人可以通过使用适当的辅助技术（提供其他用户接口的工具）来使用这些应用程序。很多政府和企业都有一些强制要求应用程序为实现可访问性而采用的标准。

Sun 希望 Java 语言能够成为一种 “编写一次就可以随处运行（write once, run everywhere，即 WORE）” 的环境。这意味着可以在一台机器上开发和测试 Java 代码（例如在 Windows® 上），然后不经测试就可以在另外一个 Java 主机上运行同样的 Java 代码。对于大部分情况来说，Java 技术都可以成功实现这种功能，但是 AWT 却是一个弱点。由于 AWT 要依赖于主机 GUI 的对等体（peer）控件（其中每个 AWT 组件都有一个并行的主机控件或者对等体）来实现这个 GUI，这个 GUI 的外观和行为（这一点更重要）在不同的主机上会有所不同。这会导致出现 “编写一次随处测试（write once, test everywhere，即 WOTE）” 的情况，这样就远远不能满足我们的要求了。

AWT 提供了一个丰富的图形环境，尤其是在 Java V1.2 及其以后版本中更是如此。通过 _Graphics2D_ 对象和 _Java2D_ 、 _Java3D_ 服务，我们可以创建很多功能强大的图形应用程序，例如画图和制表包；结合使用 _JavaSound_ ，我们还可以创建非常有竞争力的交互式游戏。

## Swing 概述

Java Swing 是 _Java Foundation Classes_ （JFC）的一部分，它是试图解决 AWT 缺点的一个尝试。在 Swing 中，Sun 开发了一个经过仔细设计的、灵活而强大的 GUI 工具包。不幸的是，这意味着我们又要花一些时间来学习 Swing 了，对于常见的情况来说，Swing 有些太复杂了。

Swing 是在 AWT 组件基础上构建的。所有 Swing 组件实际上也是 AWT 的一部分。Swing 使用了 AWT 的事件模型和支持类，例如 Colors、Images 和 Graphics。Swing 组件、布局管理器以及事件总结如下（参见 [清单 4](#清单-4-基本的-swing-class-树) 、 [清单 5](#清单-5-swing-提供了以下-layoutmanagers) 和 [清单 6](#清单-6-swing-提供了以下事件) ）。正如您可以看到的一样，这些组件集比 AWT 提供的组件集更为广泛，与 SWT 组件集相比也毫不逊色。

##### 清单 4\. 基本的 Swing Class 树

全部在 javax.swing 包或其子包中，`*` 表示抽象类。

```
Object
    *Component
        Container
            *JComponent
                *AbstractButton
                    JButton
                    JMenuItem
                        JCheckBonMenuItem
                        JMenu
                        JRadioButonMenuItem
                    *JToggleButton
                        JCheckBox
                        JRadioButton
                Box
                Filler
                JColorChooser
                JComboBox
                JDesktopIcon
                JFileChooser
                JInternalFrame
                JLabel
                JLayeredPane
                    JDesktopPane
                JList
                JMenuBar
                JOptionPane
                JPanel
                JPopupMenu
                JProgressBar
                JRootPane
                JScrollBar
                JScrollPane
                JSeparator
                JSlider
                JSplitPane
                JTabbedPane
                JTable
                JTableHeader
                *JTextComponent
                    JEditorPane
                        FrameEditorPane
                        JTextPane
                    JTextArea
                    JtextField
                        JPasswordField
                JToolBar
                JToolTip
                JTree
                JViewport
                    ScrollableTabViewport
            Panel
                Applet
                    JApplet
            Window
                Dialog
                    JDialog
                Frame
                    JFrame
                JWindow

```

Show moreShow more icon

**注意** ：在另外几个包中还有其他一些 Swing 组件，但是这是基本的组件集。

##### 清单 5\. Swing 提供了以下 LayoutManagers

全部在 javax.swing 包或其子包中，`*` 表示接口。

```
*LayoutManager
    CenterLayout
    *LayoutManager2
        BoxLayout
        OverlayLayout
        SpringLayout

```

Show moreShow more icon

**注意** ：在另外几个包中还有其他一些 Swing 布局管理器，很多都是为它们所布局的容器而专门定制的，但是这是基本的布局管理器集。

##### 清单 6\. Swing 提供了以下事件

大部分在 javax.swing.events 包及其子包中。

```
Object
    EventObject
        AWTEvent
            AncestorEvent
            ComponentEvent
                InputEvent
                    KeyEvent
                        MenuKeyEvent
                    MouseEvent
                        MenuDragMouseEvent
            InternalFrameEvent

```

Show moreShow more icon

**注意** ：在另外几个包中还有其他一些 AWT 事件，但是这是基本的事件集。这些是从更通用的事件生成的 “高级” 事件。

为了克服在不同主机上行为也会不同的缺点，Swing 将对主机控件的依赖性降至了最低。实际上，Swing 只为诸如窗口和框架之类的 _顶层_ 组件使用对等体。大部分组件（JComponent 及其子类）都是使用纯 Java 代码来模拟的。这意味着 Swing 天生就可以在所有主机之间很好地进行移植。因此，Swing 通常看起来并不像是本地程序。实际上，它有很多外观，有些模拟（尽管通常并不精确）不同主机的外观，有些则提供了独特的外观。

Swing 对基于对等体的组件使用的术语是 _重量级（heavyweight）_ ，对于模拟的组件使用的术语是 _轻量级（lightweight）_ 。实际上，Swing 可以支持在一个 GUI 中混合使用重量级组件和轻量级组件，例如在一个 JContainer 中混合使用 AWT 和 Swing 控件，但是如果组件产生了重叠，就必须注意绘制这些控件的顺序。

Swing 无法充分利用硬件 GUI 加速器和专用主机 GUI 操作的优点。结果是 Swing 应用程序可能比本地 GUI 的程序运行速度都慢。Sun 花费了大量的精力来改进最近版本的 Swing （Java V1.4 和 1.5）的性能，这种缺点正在变得日益微弱。由于 Swing 的设计更加健壮，因此其代码基础也更坚实。这意味着它可以在一台健壮的机器上比 AWT 和 SWT 上运行得更好。

除了具有更多的组件、布局管理器和事件之外，Swing 还有很多特性使得自己比 AWT 的功能更加强大。下面是更为重要的几个特性：

- **模型与视图和控件分离**: 对于这个模型中的所有组件（例如按钮、列表、表、树、富文本）来说，模型都是与组件分离的。这样可以根据应用程序的需求来采用模型，并在多个视图之间进行共享。为了方便起见，每个组件类型都提供有默认的模型。

- **可编程外观**: 每个组件的外观（外表以及如何处理输入事件）都是由一个单独的、可动态替换的实现来进行控制的。这样我们就可以改变基于 Swing 的 GUI 的部分或全部外观。

- **呈现器和编辑器**: 大部分显示模型内容的组件，例如列表、表和树，都可以处理几乎所有类型的模型元素。这可以通过为每种组件类型和模型类型映射一个渲染器或编辑器来实现。例如，一个具有包含 `java.util.Date` 值的列的表可以有一些专用的代码来呈现数据值和编辑数据值。每一列都可以有不同的类型。

- **可访问性**: 创建一个残疾人可以访问的 GUI 是非常重要的。Swing 为实现具有可访问性的 GUI 提供了丰富的基础设施和 API。这种支持是单独的，但是如果主机上具有可访问性支持，那么它们应该集成在一起。


与 AWT 一样，Swing 可以支持 GUI 组件的自动销毁。Swing 还可以支持 AWT 的自底向上和自顶向下的构建方法。

与 AWT 不同，Swing 组件不是 _线程安全的_ ，这意味着您需要关心在应用程序中是哪个线程在更新 GUI。如果在使用线程时出现了错误，就可能会出现不可预测的行为，包括用户界面故障。有一些工具可以帮助管理线程的问题。

与 AWT 类似，Swing 的一个优点是，它也是 Java 技术的一种标准配置。这意味着您不需要自己来安装它了。不幸的是，Swing 已经有了很大的变化，因此它很容易变得依赖于最新版本的 Java 语言所提供的特性，这可能会强制用户更新自己的 Java 运行时环境。

## SWT 概述

与 AWT 的概念相比，SWT 是一个低级的 GUI 工具包。JFace 是一组用来简化使用 SWT 构建 GUI 的增强组件和工具服务。SWT 的构建者从 AWT 和 Swing 实现中学习了很多经验，他们试图构建一个集二者优点于一体而没有二者的缺点的系统。从很多方面来说，他们已经成功了。

SWT 也是基于一个对等体实现的，在这一点上它与 AWT 非常类似。它克服了 AWT 所面临的 LCD 的问题，方法如下：定义了一组控件，它们可以用来构建大部分办公应用程序或开发者工具，然后可以按照逐个主机的原则，为特定主机所没有提供的控件创建模拟控件（这与 Swing 类似）。对于大部分现代主机来说，几乎所有的控件都是基于本地对等体的。这意味着基于 SWT 的 GUI 既具有主机外观，又具有主机的性能。这样就避免了使用 AWT 和 Swing 而引起的大部分问题。特定的主机具有一些低级功能控件，因此 SWT 提供了扩充（通常是模拟的）版本（通常使用 “C” 作为名字中的第一个字母），从而可以产生更一致的行为。

在对等体工作方式上，SWT 与 AWT 不同。在 SWT 中，对等体只是主机控件上的一些封装程序而已。在 AWT 中，对等体可以提供服务来最小化主机之间的差异（就是在这里，AWT 碰到了很多行为不一致的问题）。这意味着 SWT 应用程序实际上就是一个主机应用程序，它必然会全部继承主机的优点和缺点。这还意味着 SWT 不能完全实现 WORE 的目标；它更像是一种 WOTE 解决方案。这就是说，SWT 尽管不如 Swing 那么优秀，但是它在创建可移植解决方案方面是很杰出的。

SWT 部件、布局和事件总结如下（参见 [清单 7](#清单-7-基本的-swt-class-树) 、 [清单 8](#清单-8-swt-提供了以下布局管理器) 和 [清单 9](#清单-9-swt-提供了以下事件) ）。正如您可以看到的一样，这些组件集比 AWT 提供的组件集更为广泛，与 Swing 组件集相比也毫不逊色。

##### 清单 7\. 基本的 SWT Class 树

大部分在 org.ecipse.swt.widgets 或 org.eclipse.swt.custom 包或子包中，`*` 表示抽象类，`!` 表示在 custom 包中，`~` 表示在其他包中。

```
Object
    *Dialog
         ColorDialog
         DirectoryDialog
         FileDialog
         FontDialog
         MessageDialog
         PrintDialog
    *Widget
        Menu
        *Item
            CoolItem
            !CTabItem
            MenuItem
            TabItem
            TableColumn
            TableItem
            TableTreeItem
            ToolItem
            TrayItem
            TreeColumn
            TreeItem
        *Control
            Button
            Label
            ProgressBar
            Sash
            Scale
            Scrollable
                Composite
                    ~Browser
                    Canvas
                        *~AbstractHyperlink
                            ~Hyperlink
                                ~ImageHyperlink
                            *~ToggleHyperline
                                ~TreeNode
                                ~Twistie
                        AnimatedProgress
                        !CLabel
                        Decorations
                            Shell
                        FormText
                        StyledText
                        TableCursor
                    !CBanner
                    !CCombo
                    Combo
                    CoolBar
                    !CTabFolder
                    ~ExpandableComposite
                        ~Section
                    ~FilteredList
                    ~FilteredTree
                    ~Form
                    Group
                    ~PageBook
                    ProgressIndicator
                    !SashForm
                    !ScrolledComposite
                    TabFolder
                    Table
                    TableTree
                    ToolBar
                    Tray
                    Tree
                    ViewForm
                List
                Text
            Slider

```

Show moreShow more icon

**注意** ：在另外几个包中还有其他一些 SWT 部件，但是这是基本的部件集。

与 AWT 和 Swing 布局管理器类似，SWT 也提供了非常丰富的布局部件集。布局系统与嵌套容器一起使用，可以生成所需要的任何布局算法。所有这 3 个 GUI 库也可以支持对部件的定位实现绝对控制。SWT 没有等效的 BorderLayout 部件，这一点非常令人失望。FormLayout 对于创建表单基本输入来说非常好用。我认为 SWT 的布局机制比 AWT/Swing 布局部件集的使用更难学习。

##### 清单 8\. SWT 提供了以下布局管理器

大部分在 org.eclipse.swt.layout 和 org.eclipse.swt.custom 包或子包中，`*` 表示接口，`!` 表示在 custom 包中。

```
*Layout
    FillLayout
    FormLayout
    GridLayout
    RowLayout
    !StackLayout

```

Show moreShow more icon

**注意** ：在另外几个包中还有其他一些 SWT 布局管理器，很多都是为它们所布局的容器而专门定制的，但是这是基本的布局管理器集。

与 AWT 和 Swing 事件系统一样，SWT 提供了非常丰富的事件集。尽管这些事件并不能与 AWT/Swing 的事件一一对应（例如 AWT 和 Swing 的按钮都会产生 _ActionEvent_ 事件，而 SWT 的按钮产生的则是 _SelectionEvent_ 事件），但是它们通常都是等价的。

##### 清单 9\. SWT 提供了以下事件

大部分在 org.eclipse.swt.events 包或 org.eclipse.swt.custom 包或其子包中，`*` 表示抽象，`!` 表示在 custom 包中。

```
Object
    EventObject
        SWTEventObject
            TypedEvent
                AimEvent
                !BidiSegmentEvent
                ControlEvent
                !CTabFlolderEvent
                DisposeEvent
                DragSourceEvent
                DragTargetEvent
                !ExtendedModifyEvent
                focusEvent
                HelpEvent
                KeyEvent
                    TraverseEvent
                    VerifyEvent
                !LineBackgroundEvent
                !LineStyleEvent
                MenuEvent
                ModifyEvent
                MouseEvent
                PaintEvent
                SelectionEvent
                    TreeEvent
                ShellEvent
                !TextChangedEvent
                !TextChangingEvent

```

Show moreShow more icon

**注意** ：在另外几个包中还有其他一些 SWT 事件，但是这是基本的事件集。这些是从更通用的事件生成的具体事件。

很多 Swing 组件，例如 JTable，都有自己的模型。对应的 SWT 控件（例如 Table）则没有；不过它们有自己的条目。条目通常用来限制显示文本或通常很小的图像（例如图标）。为了提供一种类 Swing 的模型接口，SWT 使用了 JFace _ContentProviders_ 。这些组件可以在应用程序提供的模型（例如 List 或 Table 使用的 `java.util.Array` ）和用作视图的控件之间充当一个桥梁。为了将任意模型对象格式化成条目，SWT 使用了 JFace _LabelProviders_ ，它们可以为任何模型对象生成一个文本或图标格式。这可以对复杂模型对象的混合显示进行限制。其他类似组件，例如 _ColorProviders_ 和 _LabelDecorators_ ，可以增强对这些条目的显示。对于 Tables 的特例来说，SWT 提供了 _CellEditor_ ，它可以临时将任意 SWT 控件链接到一个 Table 单元格上，从而当作这个单元格的编辑器使用。

SWT 不支持 GUI 控件的自动销毁。这意味着我们必须显式地销毁所创建的任何控件和资源，例如颜色和字体，而不能利用 API 调用来实现这种功能。这种工作从某种程度上来说得到了简化，因为容器控制了其子控件的自动销毁功能。

使用 SWT 只能自顶向下地构建 GUI。因此，如果没有父容器，子控件也就不存在了；通常父容器都不能在以后任意改变。这种方法不如 AWT/Swing 灵活。控件是在创建时被添加到父容器中的，在销毁时被从父容器中删除的。而且 SWT 对于 _style_ 位的使用只会在构建时进行，这限制了有些 GUI 控件的灵活性。有些风格只是一些提示性的，它们在所有平台上的行为可能并不完全相同。

与 Swing 类似，SWT 组件也不是线程安全的，这意味着您必须要关心在应用程序中是哪个线程对 GUI 进行了更新。如果在使用线程时发生了错误，就会抛出异常。我认为这比不确定的 Swing 方法要好。有一些工具可以帮助管理线程的问题。

如果所支持的操作系统提供了可访问性服务，那么 SWT GUI 通常也就具有很好的可访问性。当默认信息不够时，SWT 为程序员提供了一个基本的 API 来指定可访问性信息。

SWT 提供了一个有限的图形环境。到目前为止，它对于 Java2D 和 Java3D 的支持还不怎么好。Eclipse 使用一个名为 _Draw2D_ 的组件提供了另外一种单独的 _图形编辑框架（Graphical Editing Framework，GEF）_ ，它可以用来创建一些绘图应用程序，例如 UML 建模工具。不幸的是，GEF 难以单独（即在整个 Eclipse 环境之外）使用。

与 AWT 和 Swing 不同，SWT 和 JFace 并不是 Java 技术的标准配置。它们必须单独进行安装，这可以当作是 Eclipse 安装的一部分，也可以当作是单独的库进行安装。Eclipse 小组已经使它的安装变得非常简单，并且 SWT 可以与 Eclipse 分开单独运行。所需要的 Java 档案文件（JAR）和动态链接库（DLL）以及 UNIX® 和 Macintosh 上使用的类似库可以从 Eclipse Web 站点上单独下载。JFace 库需要您下载所有的 Eclipse 文件，并拷贝所需要的 JAR 文件。在下载所需要的文件之后，我们还需要将这些 JAR 文件放到 Java CLASSPATH 中，并将 DLL 文件放到系统 PATH 中。

## 特性的比较

下表对 AWT、SWT 和 Swing 库的很多特性进行了比较，这种比较并没有按照任何特定顺序来进行。尽管没有完全列出所有特性，但是列出了很多最重要的特性。

##### 表 1\. SWT 、AWT 和 Swing 特性的比较

功能/角色/外表AWTSwingSWT（风格）显示静态文本LabelJLabelLabel, CLabel显示多行静态文本Multiple Labels具有 HTML 内容的 Multiple JLabels 或 JLabel具有新行的 Multiple Labels 或 Label显示多行格式化静态文本具有不同字体的 Multiple Labels具有 HTML 内容的 JLabel具有不同字体的 Multiple Labels单行文本输入TextFieldJTextFieldText(SWT.SINGLE)多行文本输入TextAreaJTextAreaText(SWT.MULTI)显示图像N/AJLabelLabel显示文本和图像N/AJLabelCLabel提示弹出帮助N/A组件的 setToolTip，JToolTip 子类控件的 setToolTip风格化的文本输入N/AJEditorPaneStyledText从条目列表中进行选择ListJListList简单按下具有文本的按钮ButtonJButtonButton(SWT.PUSH)简单按下具有文本或图像的按钮N/AJButtonButton(SWT.PUSH)绘图区域；可能用于定制控件CanvasJPanelCanvas选中/取消复选框CheckBoxJCheckBoxButton(SWT.CHECK)单选按钮选择CheckBoxGroupButtonGroup 和 MenuGroup 和 Menu从一个下拉列表中选择ChoiceJComboBoxCombo、CCombo输入文本或从下拉列表中选择N/AJComboBoxCombo、CCombo可滚动区域ScrollPaneJScrollPane创建 Scrollable 子类顶层窗口Dialog、Frame、WindowJDialog、JFrame、JWindow具有不同风格的 Shell通用窗口WindowJWindowShell框架窗口FrameJFrameShell(SWT.SHELL\_TRIM)对话框窗口DialogJDialogShell(SWT.DIALOG\_TRIM)菜单MenuJMenuMenuMenuItemMenuItemJMenuItemMenuItem菜单快捷键通用击键与 AWT 相同依赖于主机的快捷键弹出菜单PopupMenuJPopupMenuMenu(SWT.POPUP)菜单条MenuBarJMenuBarMenu(SWT.BAR)显示插入符号N/ACaretCaretWeb 浏览器N/AJTextPane（HTML 3.2）Browser（通过嵌入式浏览器）Web 页面中的嵌入式控件AppletJApplet主机控件（例如 OLE）其他控件的通用容器PanelJPanelComposite其他控件的有边界通用容器Panel（如果是手工画的）具有 Border 的 JPanelComposite(SWT.BORDER)其他控件的有边界和标题的通用容器N/A具有 TitledBorder 的 JPanelGroup单选按钮（一个被选中）CheckboxJRadioButtonButton(SWT.RADIO)单选按钮的控件扩充CheckboxGroupRadioButtonGroupGroup箭头按钮N/A具有图像的 JButtonButton(SWT.ARROW)支持文本显示方向通过 ComponentOrientation与 AWT 相同很多组件都可以支持这种风格焦点切换Policy 和 Manager 对象与 AWT 相同下一个控件定制对话框Dialog 子类JDialog 子类Dialog 子类访问系统事件EventQueue 服务与 AWT 相同Display 服务（不如 AWT 健壮）系统访问对话框FileDialogJColorChooser、JFileChooserColorDialog、DirectoryDialog、FileDialog、FontDialog、PrintDialog显示简单消息对话框N/A（必须是 Dialog 子类）JOptionPane 静态方法具有很多风格的 MessageBox显示简单提示对话框N/A（必须是 Dialog 子类）JOptionPane 静态方法N/A（JFace 中用来实现这种功能的子类）布局管理器BorderLayout、CardLayout、FlowLayout、GridLayout、GridBagLayoutAWT 加上 BoxLayout、CenterLayout、SpringLayout FillLayout、FormLayout、GridLayout、RowLayout、StackLayout基本的绘图控件CanvasJPanelCanvas基本绘图Graphics 和 Graphics2D 对象 —— 基本形状和文本，任意 Shapes 和 Strokes、Bezier 以及文件与 AWT 相同GC 对象 —— 基本形状和文本绘图转换Affine，合成与 AWT 相同N/A离屏绘图（Off screen drawing）BufferedImage、drawImage与 AWT 相同Image、drawImage双缓冲区手工自动或手工除非由主机控件提供，否则就是手工打印PrintJob 和 PrintGraphics与 AWT 相同向 Printer 设备绘图定制颜色Color与 AWT 相同Color定制字体Font、FontMetrics与 AWT 相同Font光标选择Cursor与 AWT 相同Cursor图像特性从文件中加载，动态创建，可扩充地编辑与 AWT 相同从文件中加载，动态创建，基本编辑输入自动化Robot与 AWT 相同N/A显示工具条N/AJToolBarToolBar、CoolBar显示进度条N/AJProgressBarProgressBar将空间划分成区域N/AJSplitPaneSash 或 SashForm显示一个分标签页的区域N/AJTabbedPaneTabFolder、CTabFolder显示制表信息N/AJTableTable格式化表的列N/ATableColumnTableColumn显示层次化信息N/AJTreeTree从一定范围的值中进行选择N/AJSliderSlider从一组离散范围的值中进行选择N/AJSpinnerScale对于基本显示的访问Toolkit、GraphicsConfiguration、GraphicsDevice与 AWT 相同Display将条目添加到系统托盘（system tray）中N/AN/ATray

**关键** ：N/A —— 不适用。在很多情况中，这种特性都可以通过创建定制控件或控件容器或利用其他定制编程来实现，不过实现的难度会有所不同。

## 结束语

本文对 Eclipse 的 Standard Windows Toolkit with JFace、Java 的 Swing 和 Abstract Windows Toolkit GUI 工具包进行了比较。通过此处提供的比较，您可以确定在自己的新应用程序中应该使用哪个 GUI 工具包。

在大部分情况中，决定都是在 Swing 与结合了 JFace 的 SWT 之间进行的。通常来说，每个工具包都非常完整且功能强大，足以构建功能完善的 GUI，但是 Swing 通常要比单独使用 SWT（不使用 JFace 时）更好。Swing 具有内嵌于 Java 技术的优点，是完全可移植的，无可争议地是一种更好的架构。Swing 也具有高级图形应用程序所需要的优点。SWT 具有可以作为本地应用程序实现的优点，这可以提高性能，并利用基于 SWT 的 GUI 来实现本地兼容性。

如果您只为一种平台来开发系统，那么 SWT 就具有主机兼容性方面的优点，包括与主机特性的集成，例如在 Windows 上对 ActiveX 控件的使用。