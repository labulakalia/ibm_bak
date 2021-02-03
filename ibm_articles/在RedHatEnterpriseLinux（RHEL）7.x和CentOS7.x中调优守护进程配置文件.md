# 在 Red Hat Enterprise Linux（RHEL）7.x 和 CentOS 7.x 中调优守护进程配置文件
调优守护程序的 Throughput-performance 配置文件与 Balanced 配置文件之间的差异

**标签:** IBM Power Systems,Linux

[原文链接](https://developer.ibm.com/zh/articles/l-tuned-daemon-profiles/)

Kumuda G

发布: 2020-03-05

* * *

## 简介

经过调优的守护程序使用 **udev** 守护程序来监控互联设备，并根据所选配置文件静态和动态调整系统设置。经过调优的守护程序具有许多预定义的配置文件，比如 Throughput-performance、Balanced 以及 Virtual-guest，适合于常见用例。系统管理员也可以为自己的工作负载创建定制的配置文件。

在 Red Hat Enterprise Linux (RHEL) V7 操作系统中，缺省配置文件为 **Throughput-performance**。 **Virtual-guest** 配置文件选定用于虚拟机， **Balanced** 配置文件则选定用于所有其他案例。 **Balanced** 配置文件是 CentOS 7.x 操作系统中的缺省配置文件。经过调优的守护程序具有一系列推荐或选择某个配置文件作为缺省配置文件的规则。这些规则使用正则表达式来匹配 **/etc/system-release-cpe** 文件中的 **computenode\|server** 字符串，也会使用 `virt-what` 命令的输出。如果第一个表达式为 true，`virt-what` 命令的输出为空，那么推荐使用 **Throughput-performance** 配置文件。如果 `virt-what` 命令的输出显示有关虚拟机的数据，那么推荐使用 **Virtual-guest** 配置文件。这些规则都不符合时，在经过调优的守护程序中选择 **Balanced** 配置文件作为缺省配置文件。

## Throughput-performance 配置文件

带有操作系统且充当计算节点的系统往往旨在达到最佳的吞吐量性能。 **Throughput-performance** 配置文件是专为高吞吐量而优化的服务器配置文件。此配置文件禁用节能机制，启用 `sysctl` 设置来提高磁盘和网络 I/O 的吞吐量性能，并可切换到 Deadline Scheduler。CPU governor 参数设置为 performance。

## Balanced 配置文件

此配置文件的目标是均衡性能和能耗。它旨在成为平衡性能与能耗的解决方案，尽可能使用自动扩展和自动调优功能。对于大部分负载都能取得良好效果。使用该配置文件的缺点就是会增加延迟时间。

该配置文件支持处理器、磁盘、音频和视频插件，并可激活随需应变的调节器策略，即 CPU 频率调节器策略。Radeon 视频图形卡的 **radeon\_powersave** 参数设置为 auto。

您可以通过使用 `tuned-adm` 命令在这些配置文件之间变更或切换。例如，运行 `tuned-adm profile <profile>` 命令。您可以选择最适合自身工作负载的配置文件。 **Throughput-performance** 和 **Balanced** 配置文件是非虚拟系统上的缺省配置文件。

## Virtual-guest 配置文件

该配置文件专为虚拟机而优化。该配置文件的设置可以降低虚拟内存的 swappiness 值，并提高磁盘预读值。

## 配置文件参数比较

下表介绍了这些配置文件之间一些不同系统参数的比较情况。

表 1\. 比较配置文件参数 {: #比较配置文件参数}

系统参数Throughput-performance 配置文件的缺省值Balanced 配置文件的缺省值描述kernel.sched\_wakeup\_granularity\_ns15 微秒/µs2 微秒/µs为抢占当前任务而激活的任务功能。如果将此参数设置为较大的值，那么其他任务就很难实施抢占行为。该参数用于减少过度调度情况。一方面，降低该参数的值可减少唤醒延迟。但是，也可能造成频繁的切换。kernel.sched\_min\_granularity\_ns1 微秒/µs15 微秒/µs某个任务符合被抢占条件之前至少运行的时间。该参数用于控制未抢占时任务可能运行的时间量。对于延迟敏感型任务（或庞大线程数），将该参数设置为较低的值，而对于受计算量限制或面向吞吐量的工作负载，则将其设置为较高的值。vm.dirty\_ratio40%20%表示在所有进程将脏缓存写回磁盘之前，可使用脏页面所占 MemTotal（可用内存总量）的百分比。达到该值后，将会阻止所有新写操作的所有 I/O ，直至清空脏页面为止。如果将该参数设置为较低的值，那么内核会更频繁地清空小型写操作。如果将该参数设置为较高的值，那么小型写操作会堆积在内存中。vm.swappiness1060高参数值可提高文件系统性能，同时积极地将不太活跃的进程换出 RAM。低参数值则避免将进程换出内存，这通常是以 I/O 性能为代价来减少延迟。对于交互式应用，将该参数设置为低值，在减少响应延迟方面有所帮助。CPU governorPerformance（在 scaling\_min\_freq 和 scaling\_max\_freq 参数范围内，将处理器静态设置为最高频率）Ondemand（根据当前系统负载设置处理器频率）处理器根据需求提高或降低频率。设置为 performance 的 governor 参数可保持固定的高频率。它不会根据负载而变化。它最符合性能要求，但能耗也更大。ondemand 参数值可以根据工作负载快速提高频率，并在工作负载完成后慢慢降低。

本文翻译自： [Tuned daemon profiles in Red Hat Enterprise Linux (RHEL) 7.x and CentOS 7.x](https://developer.ibm.com/articles/l-tuned-daemon-profiles/)（2018-01-31）