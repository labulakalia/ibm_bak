# 使用 Micrometer 记录 Java 应用性能指标
Micrometer 与流行监控系统的整合以及与 Spring Boot 的集成

**标签:** Java,Spring

[原文链接](https://developer.ibm.com/zh/articles/j-using-micrometer-to-record-java-metric/)

成 富

发布: 2019-02-13

* * *

运行良好的应用离不开对性能指标的收集。这些性能指标可以有效地对生产系统的各方面行为进行监控，帮助运维人员掌握系统运行状态和查找问题原因。性能指标监控通常由两个部分组成：第一个部分是性能指标数据的收集，需要在应用程序代码中添加相应的代码来完成；另一个部分是后台监控系统，负责对数据进行聚合计算和提供 API 接口。在应用中使用计数器、计量仪和计时器来记录关键的性能指标。在专用的监控系统中对性能指标进行汇总，并生成相应的图表来进行可视化分析。

目前已经有非常多的监控系统，常用的如 Prometheus、New Relic、Influx、Graphite 和 Datadog，每个系统都有自己独特的数据收集方式。这些监控系统有的是需要自主安装的软件，有的则是云服务。它们的后台实现千差万别，数据接口也是各有不同。在指标数据收集方面，大多数时候都是使用与后台监控系统对应的客户端程序。此外，这些监控系统一般都会提供不同语言和平台使用的第三方库，这不可避免的会带来供应商锁定的问题。一旦针对某监控系统的数据收集代码添加到应用程序中，当需要切换监控系统时，也要对应用程序进行大量的修改。Micrometer 的出现恰好解决了这个问题，其作用可以类比于 SLF4J 在 Java 日志记录中的作用。

## Micrometer 简介

Micrometer 为 Java 平台上的性能数据收集提供了一个通用的 API，应用程序只需要使用 Micrometer 的通用 API 来收集性能指标即可。Micrometer 会负责完成与不同监控系统的适配工作。这就使得切换监控系统变得很容易。Micrometer 还支持推送数据到多个不同的监控系统。

在 Java 应用中使用 Micrometer 非常的简单。只需要在 Maven 或 Gradle 项目中添加相应的依赖即可。Micrometer 包含如下三种模块，分组名称都是 io.micrometer：

- 包含数据收集 SPI 和基于内存的实现的核心模块 `micrometer-core` 。
- 针对不同监控系统的实现模块，如针对 Prometheus 的 `micrometer-registry-prometheus` 。
- 与测试相关的模块 `micrometer-test` 。

在 Java 应用中，只需要根据所使用的监控系统，添加所对应的模块即可。比如，使用 Prometheus 的应用只需要添加 `micrometer-registry-prometheus` 模块即可。模块 `micrometer-core` 会作为传递依赖自动添加。本文使用的 Micrometer 版本是 1.1.1。清单 1 给出了使用 Micrometer 的 Maven 项目的示例：

##### 清单 1\. 使用 Micrometer 的 Maven 项目

```
<dependency>
    <groupId>io.micrometer</groupId>
    <artifactId>micrometer-registry-prometheus</artifactId>
    <version>1.1.1</version>
</dependency>

```

Show moreShow more icon

## 计量器注册表

Micrometer 中有两个最核心的概念，分别是计量器（Meter）和计量器注册表（MeterRegistry）。计量器表示的是需要收集的性能指标数据，而计量器注册表负责创建和维护计量器。每个监控系统有自己独有的计量器注册表实现。模块 `micrometer-core` 中提供的类 `SimpleMeterRegistry` 是一个基于内存的计量器注册表实现。 `SimpleMeterRegistry` 不支持导出数据到监控系统，主要用来进行本地开发和测试。

Micrometer 支持多个不同的监控系统。通过计量器注册表实现类 `CompositeMeterRegistry` 可以把多个计量器注册表组合起来，从而允许同时发布数据到多个监控系统。对于由这个类创建的计量器，它们所产生的数据会对 `CompositeMeterRegistry` 中包含的所有计量器注册表都产生影响。在清单 2 中，我创建了一个 `CompositeMeterRegistry` 对象，并在其中添加了两个 `SimpleMeterRegistry` 对象。一个 `SimpleMeterRegistry` 对象在创建时通过实现 SimpleConfig 接口提供了不同的名称前缀。

##### 清单 2\. CompositeMeterRegistry 使用示例

```
public class CompositeMeterRegistryExample {

public static void main(String[] args) {
    CompositeMeterRegistry registry = new CompositeMeterRegistry();
    registry.add(new SimpleMeterRegistry());
    registry.add(new SimpleMeterRegistry(new MyConfig(), Clock.SYSTEM));

    Counter counter = registry.counter("simple");
    counter.increment();
}

private static class MyConfig implements SimpleConfig {

    public String get(final String key) {
      return null;
    }

    public String prefix() {
      return "my";
    }
}
}

```

Show moreShow more icon

Micrometer 本身提供了一个静态的全局计量器注册表对象 Metrics.globalRegistry。该注册表是一个组合注册表。使用 Metrics 类中的静态方法创建的计量器，都会被添加到该全局注册表中。对于大多数应用来说，这个全局注册表对象就可以满足需求，不需要额外创建新的注册表对象。不过由于该对象是静态的，在某些场合，尤其是进行单元测试时，会产生一些问题。在 清单 3 中， `Metrics.addRegistry()` 方法直接在全局注册表对象中添加新的注册表对象，而 `Metrics.counter()` 方法创建的计数器自动添加到全局注册表中。

##### 清单 3\. 使用全局计量器注册表对象

```
public class GlobalRegistryExample {

public static void main(String[] args) {
    Metrics.addRegistry(new SimpleMeterRegistry());
    Counter counter = Metrics.counter("simple");
    counter.increment();
}
}

```

Show moreShow more icon

## 使用计量器

计量器用来收集不同类型的性能指标信息。Micrometer 提供了不同类型的计量器实现。计量器对象由计量器注册表创建并管理。

### 计量器名称和标签

每个计量器都有自己的名称。由于不同的监控系统有自己独有的推荐命名规则，Micrometer 使用句点 . 分隔计量器名称中的不同部分，如 `a.b.c` 。Micrometer 会负责完成所需的转换，以满足不同监控系统的需求。

每个计量器在创建时都可以指定一系列标签。标签以名值对的形式出现。监控系统使用标签对数据进行过滤。除了每个计量器独有的标签之外，每个计量器注册表还可以添加通用标签。所有该注册表导出的数据都会带上这些通用标签。

在清单 4 中，使用 MeterRegistry 的 `config()` 方法可以得到该注册表对象的 MeterRegistry.Config 对象，再使用 `commonTags()` 方法来设置通用标签。多个标签按照名称和值依次排列的方式来指定。在创建计量器时，在提供了名称之后，以同样的方式指定该计量器的标签。

##### 清单 4\. 计量器注册表的通用标签

```
SimpleMeterRegistry registry = new SimpleMeterRegistry();
registry.config().commonTags("tag1", "a", "tag2", "b");
Counter counter = registry.counter("simple", "tag3", "c");
counter.increment();

```

Show moreShow more icon

### 计数器

计数器（Counter）表示的是单个的只允许增加的值。通过 MeterRegistry 的 `counter()` 方法来创建表示计数器的 Counter 对象。还可以使用 `Counter.builder()` 方法来创建 Counter 对象的构建器。Counter 所表示的计数值是 double 类型，其 `increment()` 方法可以指定增加的值。默认情况下增加的值是 1.0。

如果已经有一个方法返回计数值，可以直接从该方法中创建类型为 FunctionCounter 的计数器。在清单 5 中，方法 `counter()` 使用了两种不同的方法来创建 Counter 对象。方法 `functionCounter()` 同样使用了两种不同的方法来创建 FunctionCounter 对象。

##### 清单 5\. 计数器使用示例

```
public class Counters {

private SimpleMeterRegistry registry = new SimpleMeterRegistry();

private double value = 0.0;

public void counter() {
    Counter counter1 = registry.counter("simple1");
    counter1.increment(2.0);
    Counter counter2 = Counter.builder("simple2")
        .description("A simple counter")
        .tag("tag1", "a")
        .register(registry);
    counter2.increment();
}

public void functionCounter() {
    List<Tag> tags = new ArrayList<>();
    registry.more().counter("function1", tags, this, Counters::getValue);

    FunctionCounter functionCounter = FunctionCounter.builder("function2", this, Counters::getValue)
        .description("A function counter")
        .tags(tags)
        .register(registry);
    functionCounter.count();
}

private double getValue() {
    return value++;
}
}

```

Show moreShow more icon

### 计量仪

计量仪（Gauge）表示的是单个的变化的值。与计数器的不同之处在于，计量仪的值并不总是增加的。与创建 Counter 对象类似，Gauge 对象可以从计量器注册表中创建，也可以使用 `Gauge.builder()` 方法返回的构造器来创建。清单 6 中给出了计量仪的使用示例，其中 `gauge()` 方法创建的是记录任意 Number 对象的值， `gaugeCollectionSize()` 方法记录集合的大小， `gaugeMapSize()` 方法记录 Map 的大小。需要注意的是，这 3 个方法返回的并不是 Gauge 对象，而是被记录的对象。这是由于 Gauge 对象一旦被创建，就不能手动对其中的值进行修改。在每次取样时，Gauge 会返回当前值。正因为如此，得到一个 Gauge 对象，除了进行测试之外，没有其他的意义。

##### 清单 6\. 计量仪使用示例

```
public class Gauges {

private SimpleMeterRegistry registry = new SimpleMeterRegistry();

public void gauge() {
    AtomicInteger value = registry.gauge("gauge1", new AtomicInteger(0));
    value.set(1);

    List<String> list = registry.gaugeCollectionSize("list.size", Collections.emptyList(), new ArrayList<>());
    list.add("a");

    Map<String, String> map = registry.gaugeMapSize("map.size", Collections.emptyList(), new HashMap<>());
    map.put("a", "b");

    Gauge.builder("value", this, Gauges::getValue)
        .description("a simple gauge")
        .tag("tag1", "a")
        .register(registry);
}

private double getValue() {
    return ThreadLocalRandom.current().nextDouble();
}
}

```

Show moreShow more icon

### 计时器

计时器（Timer）通常用来记录事件的持续时间。计时器会记录两类数据：事件的数量和总的持续时间。在使用计时器之后，就不再需要单独创建一个计数器。计时器可以从注册表中创建，或者使用 `Timer.builder()` 方法返回的构建器来创建。Timer 提供了不同的方式来记录持续时间。第一种方式是使用 `record()` 方法来记录 Runnable 和 Callable 对象的运行时间；第二种方式是使用 `Timer.Sample` 来保存计时状态。

在清单 7 中，方法 `record()` 使用 Timer 对象的 record() 方法来记录一个 Runnable 对象的运行时间。方法 `sample()` 中首先使用 `Timer.start()` 来创建一个新的 `Timer.Sample` 对象并启动计时。调用 `Timer.Sample` 的 `stop()` 方法把记录的时间保存到 Timer 对象中。

##### 清单 7\. 计时器使用示例

```
public class Timers {
private SimpleMeterRegistry registry = new SimpleMeterRegistry();

public void record() {
    Timer timer = registry.timer("simple");
    timer.record(() -> {
      try {
        Thread.sleep(3000);
      } catch (InterruptedException e) {
        e.printStackTrace();
      }
    });
}

public void sample() {
    Timer.Sample sample = Timer.start();
    new Thread(() -> {
      try {
        Thread.sleep(2000);
      } catch (InterruptedException e) {
        e.printStackTrace();
      }
      sample.stop(registry.timer("sample"));
    }).start();
}
}

```

Show moreShow more icon

如果一个任务的耗时很长，直接使用 Timer 并不是一个好的选择，因为 Timer 只有在任务完成之后才会记录时间。更好的选择是使用 `LongTaskTimer` 。 `LongTaskTimer` 可以在任务进行中记录已经耗费的时间，它通过注册表的 `more().longTaskTimer()` 来创建，如清单 8 所示：

##### 清单 8\. LongTaskTimer 使用示例

```
public void longTask() {
LongTaskTimer timer = registry.more().longTaskTimer("long");
timer.record(() -> {
    try {
      Thread.sleep(3000);
    } catch (InterruptedException e) {
      e.printStackTrace();
    }
});
}

```

Show moreShow more icon

### 分布概要

分布概要（Distribution summary）用来记录事件的分布情况。计时器本质上也是一种分布概要。表示分布概要的类 `DistributionSummary` 可以从注册表中创建，也可以使用 `DistributionSummary.builder()` 提供的构建器来创建。分布概要根据每个事件所对应的值，把事件分配到对应的桶（bucket）中。Micrometer 默认的桶的值从 1 到最大的 long 值。可以通过 `minimumExpectedValue` 和 `maximumExpectedValue` 来控制值的范围。如果事件所对应的值较小，可以通过 scale 来设置一个值来对数值进行放大。与分布概要密切相关的是直方图和百分比（percentile）。大多数时候，我们并不关注具体的数值，而是数值的分布区间。比如在查看 HTTP 服务响应时间的性能指标时，通常关注是的几个重要的百分比，如 50%，75%和 90%等。所关注的是对于这些百分比数量的请求都在多少时间内完成。Micrometer 提供了两种不同的方式来处理百分比。

- 对于 Prometheus 这样本身提供了对百分比支持的监控系统，Micrometer 直接发送收集的直方图数据，由监控系统完成计算。
- 对于其他不支持百分比的系统，Micrometer 会进行计算，并把百分比结果发送到监控系统。

在清单 9 中，创建的 `DistributionSummary` 所发布的百分比包括 `0.5` 、 `0.75` 和 `0.9` 。使用 `record()` 方法来记录数值，而 `takeSnapshot()` 方法返回当前数据的快照。

##### 清单 9\. 分布概要使用示例

```
public class DistributionSummaries {
private SimpleMeterRegistry registry = new SimpleMeterRegistry();

public void summary() {
    DistributionSummary summary = DistributionSummary.builder("simple")
        .description("simple distribution summary")
        .minimumExpectedValue(1L)
        .maximumExpectedValue(10L)
        .publishPercentiles(0.5, 0.75, 0.9)
        .register(registry);
    summary.record(1);
    summary.record(1.3);
    summary.record(2.4);
    summary.record(3.5);
    summary.record(4.1);
    System.out.println(summary.takeSnapshot());
}

public static void main(String[] args) {
    new DistributionSummaries().summary();
}
}

```

Show moreShow more icon

## 集成监控系统

Micrometer 提供了对多种不同的监控系统的支持。

### JMX

JMX 是导出 Micrometer 收集的性能数据的最简单有效的方式。虽然 JMX 所提供的功能比较弱，但是在很多情况下，JMX 就已经可以满足需求了。如果需要导出数据到 JMX，只需要添加对库 `io.micrometer:micrometer-registry-jmx` 的依赖即可。Micrometer 会根据计量器的名称和标签来生成对应的 JMX 对象名称。默认的命名规则是在计量器名称之后，加上按标签名称字母排序的以句点分隔的名值对。

### Prometheus

Prometheus 与其他监控系统的不同在于，Prometheus 采取的是主动抽取数据的方式。因此客户端需要暴露 HTTP 服务，并由 Prometheus 定期来访问以获取数据。Micrometer 的 Prometheus 注册表已经提供了 HTTP 服务所需要返回的内容，只需要使用 Servlet 来提供 HTTP 服务即可。

## 集成 Spring Boot

从 Spring Boot 2.0 开始，Micrometer 就是 Spring Boot 默认提供的性能指标收集库。Spring Boot Actuator 提供了对 Micrometer 的自动配置。Spring Boot 会自动配置一个组合注册表对象，并把 CLASSPATH 上找到的所有支持的注册表实现都添加起来。只需要在 CLASSPATH 上添加相应的依赖库，Spring Boot 会完成所需的配置。这些注册表对象也会被自动添加到全局注册表对象中。如果需要对该注册表进行配置，添加类型为 `MeterRegistryCustomizer` 的 bean 即可。在需要使用注册表的地方，可以通过依赖注入的方式来使用 MeterRegistry 对象。

在清单 10 中，Spring 配置类 `AppConfig` 中声明了一个类型为 `MeterRegistryCustomizer<MeterRegistry>` 的 `bean` ，可以对 MeterRegistry 进行配置。这里使用 `commonTags()` 方法来添加通用标签。

##### 清单 10\. Spring Boot 中 Micrometer 的配置示例

```
@Configuration
@EnableWebMvc
@ComponentScan(basePackageClasses = AppConfig.class)
class AppConfig {
@Bean
MeterRegistryCustomizer<MeterRegistry> meterRegistryCustomizer() {
    return registry -> registry.config().commonTags("tag1", "a", "tag2", "b");
}
}

```

Show moreShow more icon

清单 11 中的 REST 控制器 AppController 通过依赖注入获取到 MeterRegistry 对象，并创建一个计数器。在方法 `greeting()` 中，对计数器进行递增。

##### 清单 11\. 使用 MeterRegistry 对象的示例

```
@RestController
@RequestMapping("/app")
public class AppController {
private final Counter counter;

public AppController(final MeterRegistry registry) {
    this.counter = registry.counter("greeting");
}

@RequestMapping("/greeting")
public String greeting() {
    this.counter.increment();
    return "hello world #" + this.counter.count();
}
}

```

Show moreShow more icon

对于 Prometheus 来说，Spring Boot Actuator 会自动配置一个 URL 为 `/actuator/Prometheus` 的 HTTP 服务来供 Prometheus 抓取数据。不过该 Actuator 服务默认是关闭的，需要通过 Spring Boot 的配置打开。清单 12 中的 application.yml 文件给出了如何打开该服务的示例。

##### 清单 12\. 启用 Prometheus 服务端点的 application.yml 文件

```
management:
endpoints:
    web:
      exposure:
        include: "*"

```

Show moreShow more icon

## 结束语

作为良好 Java 应用中的重要一环，性能指标数据的收集已经是应用中不可或缺的部分。正如 SLF4J 在 Java 日志记录中的作用一样，Micrometer 为 Java 平台上的性能指标数据收集提供了一个通用的可依赖的 API，避免了可能的供应商锁定问题。利用 Micrometer 提供的多种计量器，可以收集多种类型的性能指标数据，并通过计量器注册表发送到不同的监控系统。

## 下载示例源码

本文示例代码可以 [在 GitHub 上下载](https://github.com/VividcodeIO/micrometer-starter)。