# iCalendar 编程基础：了解和使用 iCal4j
iCalendar 标准和日历文件格式以及读写/处理 iCalendar 数据流的 API —— iCal4j

**标签:** API 管理,Java

[原文链接](https://developer.ibm.com/zh/articles/j-lo-ical4j/)

刘奇

发布: 2010-03-24

* * *

## iCalendar 简介

iCalendar，简称”iCal”，是”日历数据交换”的标准（ [RFC 2445](http://tools.ietf.org/html/rfc2445) ），该标准提供了一种公共的数据格式用于存储关于日历方面的信息，比如事件、约定、待办事项等。它不仅允许用户通过电子邮件发送会议或者待办事件等，也允许独立使用，而不局限于某种传输协议。

目前，所有流行日历工具比如：Lotus Notes、Outlook、GMail 和 Apple 的 iCal 都支持 iCalendar 标准，其文件扩展名为 .ical、.ics、.ifb 或者 .icalendar。C&S（Calendaring and Scheduling） 核心对象是一系列日历和行程安排信息。通常情况下，这些日历和行程信息仅仅包含一个 iCalendar 组件（iCalendar 组件分为 Events(VEVENT)、To-do(VTODO)、Journal(VJOURNAL)、Free/busy time (VFREEBUSY)、VTIMEZONE (time zones) 和 VALARM (alarms)），但是多个 iCalendar 组件可以被组织在一起。

C&S 核心对象第一行必须是”BEGIN:VCALENDAR”, 并且最后行必须是”END:VCALENDAR”。在这两行之间主要是由一系列日历属性和一个或者多个 iCalendar 组件组成。

下面看一个例子，它表示发生在 1997 年七月十四日下午五点与 1997 年七月十五日四点之间的事件”Bastille Day Party”。

```
BEGIN:VCALENDAR                   ------ 起始
VERSION:2.0                       ------ 版本
PRODID:iCal4j v1.0//EN            ------ 创建该对象的标志符
BEGIN:VEVENT                      ------ 事件开始
DTSTART:19970714T170000Z          ------ 事件起始时间
DTEND:19970715T040000Z            ------ 事件结束时间
SUMMARY:Bastille Day Party        ------ 事件概要
END:VEVENT                        ------ 事件结束
END:VCALENDAR                     ------ 结束

```

Show moreShow more icon

## iCalendar 编程基础

### iCal4j 简介

iCal4j（产生于 2004 年 4 月，目前是 2.0 版本），是一组读写 iCalendar 数据流的 Java API，支持 iCalendar 规范 [RFC 2445](http://tools.ietf.org/html/rfc2445) ，主要包括解析器、对象模型以及生成器。

### 文件读写

任何一种文件格式读写都是最基本的操作，清单 1 向我们演示了这两项基本操作。CalendarBuilder 对象是用来通过输入流解析和构造 iCalendar 模型。值得注意的是， CalendarBuilder 并不是 **线程安全** 的。

##### 清单 1\. iCalendar 文件读写示例

```
public static void readAndWrite(String in, String out)
    throws IOException, ParserException, ValidationException {
    FileInputStream fin = new FileInputStream(in);
    CalendarBuilder builder = new CalendarBuilder();
    Calendar calendar = builder.build(fin);
    //TODO: 对 iCalendar 数据进行处理
......
    FileOutputStream fout = new FileOutputStream(out);
    CalendarOutputter outputter = new CalendarOutputter();
    outputter.output(calendar, fout);
}

```

Show moreShow more icon

### iCalendar 索引

对组件和属性进行索引之后，我们可以更加有效的查找组件和属性，通常情况下大家用索引去不断地检查某事件（或约定等）是否存在。假定一个场景，您经常性的需要更新一些日历，并且需要检查事件是否已经存在。因为您需要不断地检查日历中的事件，此时对日历中的事件进行索引将是一件有意义的事情。

##### 清单 2\. iCalendar 索引

```
// 创建索引列表
IndexedComponentList indexedEvents = new IndexedComponentList(
     myCalendar.getComponents(Component.VEVENT), Property.UID);

// 检查事件
for (Iterator i=inputCalendar.getComponents(Component.VEVENT).iterator(); i.hasNext();){
     VEvent event = (VEvent) i.next();
     Component existing = indexedEvents.getComponent(event.getUid().getValue());
     if (existing == null) {
          myCalendar.getComponents().add(event);
     }
     else if (!event.equals(existing)) {
          // 删除已经存在的事件并添加修改后的事件
          myCalendar.getComponents().remove(existing);
          myCalendar.getComponents().add(event);
     }
}

```

Show moreShow more icon

如清单 2 所示，这里请注意，UID 被用来标示唯一的事件。在得到索引组件列表后，我们就可以通过 UID 来检查某一事件是否存在。如果存在的话，修改该事件；否则，添加该事件。

### 创建日历/事件/会议

iCalendar 向我们提供了一种公共的数据格式用于存储关于日历方面的信息比如事件、约定、待办事项等，那么如何创建事件/约定/待办事件就是一个大家比较关心的问题，如清单 3。一个日历，有一些基本的属性是必须的，例如 prodid 和 version，但是又只能有一次，一些是可选的，例如 calscale，具体的可以参考 [RFC 2445](http://tools.ietf.org/html/rfc2445) 。

##### 清单 3\. 创建一个事件示例

```
// 创建一个时区（TimeZone）
TimeZoneRegistry registry = TimeZoneRegistryFactory.getInstance().createRegistry();
TimeZone timezone = registry.getTimeZone("America/Mexico_City");
VTimeZone tz = timezone.getVTimeZone();

// 起始时间是：2008 年 4 月 1 日 上午 9 点
java.util.Calendar startDate = new GregorianCalendar();
startDate.setTimeZone(timezone);
startDate.set(java.util.Calendar.MONTH, java.util.Calendar.APRIL);
startDate.set(java.util.Calendar.DAY_OF_MONTH, 1);
startDate.set(java.util.Calendar.YEAR, 2008);
startDate.set(java.util.Calendar.HOUR_OF_DAY, 9);
startDate.set(java.util.Calendar.MINUTE, 0);
startDate.set(java.util.Calendar.SECOND, 0);

// 结束时间是：2008 年 4 月 1 日 下午 1 点
java.util.Calendar endDate = new GregorianCalendar();
endDate.setTimeZone(timezone);
endDate.set(java.util.Calendar.MONTH, java.util.Calendar.APRIL);
endDate.set(java.util.Calendar.DAY_OF_MONTH, 1);
endDate.set(java.util.Calendar.YEAR, 2008);
endDate.set(java.util.Calendar.HOUR_OF_DAY, 13);
endDate.set(java.util.Calendar.MINUTE, 0);
endDate.set(java.util.Calendar.SECOND, 0);

// 创建事件
String eventName = "Progress Meeting";
DateTime start = new DateTime(startDate.getTime());
DateTime end = new DateTime(endDate.getTime());
VEvent meeting = new VEvent(start, end, eventName);

// 添加时区信息
meeting.getProperties().add(tz.getTimeZoneId());

// 生成唯一标志符
UidGenerator ug = new UidGenerator("uidGen");
Uid uid = ug.generateUid();
meeting.getProperties().add(uid);

// 添加参加者 .
Attendee dev1 = new Attendee(URI.create("mailto:dev1@mycompany.com"));
dev1.getParameters().add(Role.REQ_PARTICIPANT);
dev1.getParameters().add(new Cn("Developer 1"));
meeting.getProperties().add(dev1);

Attendee dev2 = new Attendee(URI.create("mailto:dev2@mycompany.com"));
dev2.getParameters().add(Role.OPT_PARTICIPANT);
dev2.getParameters().add(new Cn("Developer 2"));
meeting.getProperties().add(dev2);

// 创建日历
net.fortuna.ical4j.model.Calendar icsCalendar = new net.fortuna.ical4j.model.Calendar();
icsCalendar.getProperties().add(new ProdId("-//Events Calendar//iCal4j 1.0//EN"));
icsCalendar.getProperties().add(CalScale.GREGORIAN);

// 添加事件
icsCalendar.getComponents().add(meeting);

```

Show moreShow more icon

[RFC 2445](http://tools.ietf.org/html/rfc2445) 定义了很多关联组件属性，例如 Attendee、Contact、Organizer、Recurrence ID、UID 等。在这里，着重介绍一下组件（VEVENT、VTODO、VJOURNAL、VFREEBUSY）属性 – 唯一标志符 UID。UID 必须是一个全球唯一标志符，在 iCal4j 中可以使用 UidGenerator 来生成唯一标志符。UID 是用来代表一个日历组件，经常被应用程序用来匹配后来的各种请求，例如修改、删除、 回复等，对应用程序之间通讯起着非常重要的作用。举个简单的例子，GMail 用户 A 创建了一个会议，并邀请一个 Lotus Notes 用户 B 参加。这时，B 将收到一封来自 GMail 的邮件，决定是否接受会议邀请并通知用户 A。GMail 系统如何知道来自 B 的通知（Notification）是关于哪个会议的呢？答案是 UID。

## iCalendar 高级编程

### 时间和时区

iCal4j 提供了自己的 Date/Time 和 TimeZone 实现，与之相关的类主要有 Date、DateTime、TimeZone、TimeZoneRegistry、ZoneInfo 和 TzUrl 等。

Date 对象表示一个日期（不包括时间），而 DateTime 对象则表示日期（包含时间）。如果表示一个 UTC 时间，则只需要在时间之后加上一个”Z ”，如 DTSTART:19980119T070000Z。另外，也可以对一个时间指定一个时区，例如 DTSTART;TZID=US-Eastern:19980119T020000。

TimeZoneRegistry 则表示 net.fortuna.ical4j.model.TimeZone 句柄的一个仓库（Repository）。清单 4 演示了如何获取一个 TimeZoneRegistry，并通过 TimeZoneRegistry 来获取一个时区（TimeZone）。

##### 清单 4\. 获取时区示例

```
CalendarBuilder builder = new CalendarBuilder();
Calendar calendar = builder.build(new FileInputStream("mycalendar.ics"));

TimeZoneRegistry registry = builder.getRegistry();
TimeZone tz = registry.getTimeZone("Australia/Melbourne");

```

Show moreShow more icon

### 添加附件和二进制数据

iCal4j 代码库允许用户添加二进制数据，类 Attach 可以被用来实现二进制数据的添加，如代码清单 5。

##### 代码清单 5\. 添加二进制数据示例

```
public static void attachBinaryAttachment(String subject, String out, String binary)
        throws IOException, ParserException, ValidationException, ParseException {
    FileInputStream bin = new FileInputStream(binary);
    ByteArrayOutputStream bout = new ByteArrayOutputStream();
    Calendar calendar = new Calendar();
    DateFormat format = new SimpleDateFormat("MM/dd/yyyy HH:mm");
    DateTime start = new DateTime(format.parse("11/09/2009 08:00").getTime());
    DateTime end = new DateTime(format.parse("11/09/2009 09:00").getTime());
    calendar.getProperties().add(new ProdId("-//Ben Fortuna//iCal4j 1.0//EN"));
    calendar.getProperties().add(Version.VERSION_2_0);
    calendar.getProperties().add(CalScale.GREGORIAN);

    VEvent event = new VEvent(start, end, subject);
    event.getProperties().add(new Uid(new UidGenerator("iCal4j")
        .generateUid().getValue()));
    for (int i = bin.read(); i >= 0;) {
        bout.write(i);
        i = bin.read();
    }
    ParameterList params = new ParameterList();
    params.add(Encoding.BASE64);
    params.add(Value.BINARY);
    Attach attach = new Attach(params, bout.toByteArray());
    event.getProperties().add(attach);
    calendar.getComponents().add(event);

    // 验证
    calendar.validate();

    FileOutputStream fout = new FileOutputStream(out);
    CalendarOutputter outputter = new CalendarOutputter();
    outputter.output(calendar, fout);
}

```

Show moreShow more icon

程序输出结果为：

```
BEGIN:VCALENDAR
PRODID:-//Ben Fortuna//iCal4j 1.0//EN
VERSION:2.0
CALSCALE:GREGORIAN
BEGIN:VEVENT
DTSTAMP:20091119T124443Z
DTSTART:20091109T080000
DTEND:20091109T090000
SUMMARY:iCal4j Attach Binary Example
UID:20091119T124443Z-iCal4j@192.168.1.100
ATTACH;ENCODING=BASE64;VALUE=BINARY:VGhp333cyBhIENsaWVudCBmb3IgRS1CdXNpb
mVzcyBJbWFnZQ==
END:VEVENT
END:VCALENDAR

```

Show moreShow more icon

### 循环事件

iCalendar 规范支持循环事件，即事件不止一次发生，通常使用一系列的日期（RDATE）或者一个循环规则（RRULE）。现在列举一个简单的 RRULE 例子。

- 事件每间隔一月，在每个月的 29 日发生。

`RRULE:FREQ=MONTHLY;INTERVAL=2;BYDAY=29`

_注意：非闰年二月份只有 28 天，因此这条规则将表示事件在上一年的 12 月份和 4 月之间将有一个跳跃。_

下面分别讲解一下程序如何生成一系列的日期（RDATE）和一个循环规则（RRULE）。清单 6 讲解了如何生成一个日期列表。

##### 清单 6\. 创建循环事件（使用 RDate）示例

```
public static void addRDate(String subject, String out)
        throws IOException, ParserException, ValidationException,
        URISyntaxException, ParseException {
    Calendar calendar = new Calendar();
    calendar.getProperties().add(new ProdId("-//Ben Fortuna//iCal4j 1.0//EN"));
    calendar.getProperties().add(Version.VERSION_2_0);
    calendar.getProperties().add(CalScale.GREGORIAN);

    PeriodList periodList = new PeriodList();
    ParameterList paraList = new ParameterList();
    DateFormat format = new SimpleDateFormat("MM/dd/yyyy hh:mm");
    DateTime startDate1 = new DateTime(format.parse(("11/09/2009 08:00")));
    DateTime startDate2 = new DateTime(format.parse(("11/10/2009 09:00")));
    DateTime endDate1 = new DateTime(format.parse(("11/09/2009 09:00")));
    DateTime endDate2 = new DateTime(format.parse(("11/10/2009 11:00")));
    periodList.add(new Period(startDate1, endDate1));
    periodList.add(new Period(startDate2, endDate2));

    VEvent event = new VEvent(startDate1, endDate1, subject);
    event.getProperties().add(new Uid(new UidGenerator("iCal4j").
        generateUid().getValue()));
    paraList.add(ParameterFactoryImpl.getInstance().createParameter(
        Value.PERIOD.getName(), Value.PERIOD.getValue()));
    RDate rdate = new RDate(paraList,periodList);
    event.getProperties().add(rdate);
    calendar.getComponents().add(event);

    // 验证
    calendar.validate();

    FileOutputStream fout = new FileOutputStream(out);
    CalendarOutputter outputter = new CalendarOutputter();
    outputter.output(calendar, fout);
}

```

Show moreShow more icon

程序输出结果为：

```
BEGIN:VCALENDAR
PRODID:-//Ben Fortuna//iCal4j 1.0//EN
VERSION:2.0
CALSCALE:GREGORIAN
BEGIN:VEVENT
DTSTAMP:20091119T100938Z
DTSTART:20091109T000000Z
DTEND:20091109T010000Z
SUMMARY:iCal4j RDate Example
UID:20091119T100938Z-iCal4j@IBM-3C6595B4005.cn.ibm.com
RDATE;VALUE=PERIOD:20091109T000000Z/20091109T010000Z,20091110T010000Z/200 91110T030000Z
END:VEVENT
END:VCALENDAR

```

Show moreShow more icon

清单 7 演示了如何用 RRule 创建一个循环事件，该事件从 2009 年 11 月 9 日开始，间隔一周，持续 4 次。

##### 代码清单 7\. 创建循环事件（使用 RRule）示例

```
public static void addRRule(String subject, String out)
        throws IOException, ParserException, ValidationException,
        URISyntaxException, ParseException {
    Calendar calendar = new Calendar();
    DateFormat format = new SimpleDateFormat("MM/dd/yyyy HH:mm");
    DateTime start = new DateTime(format.parse("11/09/2009 08:00").getTime());
    DateTime end = new DateTime(format.parse("11/09/2009 09:00").getTime());
    calendar.getProperties().add(new ProdId("-//Ben Fortuna//iCal4j 1.0//EN"));
    calendar.getProperties().add(Version.VERSION_2_0);
    calendar.getProperties().add(CalScale.GREGORIAN);

    VEvent event = new VEvent(start, end, subject);
    event.getProperties().add(new Uid(new UidGenerator("iCal4j").generateUid()
        .getValue()));

    Recur recur = new Recur(Recur.WEEKLY, 4);
    recur.setInterval(2);
    RRule rule = new RRule(recur);
    event.getProperties().add(rule);
    calendar.getComponents().add(event);

    // 验证
    calendar.validate();

    FileOutputStream fout = new FileOutputStream(out);
    CalendarOutputter outputter = new CalendarOutputter();
    outputter.output(calendar, fout);
}

```

Show moreShow more icon

程序输出结果为：

```
BEGIN:VCALENDAR
    PRODID:-//Ben Fortuna//iCal4j 1.0//EN
    VERSION:2.0
    CALSCALE:GREGORIAN
    BEGIN:VEVENT
        DTSTAMP:20091119T094724Z
        DTSTART:20091109T080000
        DTEND:20091109T090000
        SUMMARY:iCal4j RRule Usage Example
        UID:20091119T094724Z-iCal4j@IBM-3C6595B4005.cn.ibm.com
        RRULE:FREQ=WEEKLY;INTERVAL=2;COUNT=4
    END:VEVENT
END:VCALENDAR

```

Show moreShow more icon

### 扩展

通常情况下我们把没有在 [RFC2445](http://www.ietf.org/rfc/rfc2445.txt) 中没有定义的组件（components）/属性（properties）/参数（parameters）都称为非标准的或者扩展的对象。iCalendar 允许用户扩展组件/属性/参数，但是名字必须以 X- 开头以兼容 iCalendar 标准（除非使能 ical4j.parsing.relaxed），这些对象分别用 XComponent, XProperty and XParameter 来表示。

下面的例子演示了如何扩展一个事件的属性 ical4j\_extension\_sample，见清单 8。

##### 清单 8\. 扩展事件属性示例

```
public static void extension(String subject, String out)
        throws IOException, ParserException, ValidationException, URISyntaxException
        , ParseException {
    Calendar calendar = new Calendar();
    DateFormat format = new SimpleDateFormat("MM/dd/yyyy HH:mm");
    DateTime start = new DateTime(format.parse("11/09/2009 08:00").getTime());
    DateTime end = new DateTime(format.parse("11/09/2009 09:00").getTime());
    calendar.getProperties().add(new ProdId("-//Ben Fortuna//iCal4j 1.0//EN"));
    calendar.getProperties().add(Version.VERSION_2_0);
    calendar.getProperties().add(CalScale.GREGORIAN);

    VEvent event = new VEvent(start, end, subject);
    event.getProperties().add(new Uid(new UidGenerator("iCal4j")
        .generateUid().getValue()));
    // 注意：扩展的属性必须以 X- 开头
    Property xProperty = PropertyFactoryImpl.getInstance()
        .createProperty("X-iCal4j-extension");
    xProperty.setValue("ical4j_extension_sample");
    event.getProperties().add(xProperty);
    calendar.getComponents().add(event);
    // 验正
    calendar.validate();

    FileOutputStream fout = new FileOutputStream(out);
    CalendarOutputter outputter = new CalendarOutputter();
    outputter.output(calendar, fout);
}

```

Show moreShow more icon

程序输出结果为：

```
BEGIN:VCALENDAR
    PRODID:-//Ben Fortuna//iCal4j 1.0//EN
    VERSION:2.0
    CALSCALE:GREGORIAN
    BEGIN:VEVENT
        DTSTAMP:20091119T093429Z
        DTSTART:20091109T080000
        DTEND:20091109T090000
        SUMMARY:iCal4j Extension Example
        UID:20091119T093429Z-iCal4j@IBM-3C6595B4005.cn.ibm.com
        X-iCal4j-extension:ical4j_extension_sample
    END:VEVENT
END:VCALENDAR

```

Show moreShow more icon

## iCal4j Connector 介绍

iCal4j Connector 是 iCal4j 代码库的一个扩展，它提供连接后端 Calendar/vCard 服务器支持（支持的 HTTP 方法包括 GET, HEAD, POST, PUT, DELETE, TRACE, COPY, MOVE），服务器包括 CalDAV、与 CardDAV 相兼容的服务器和 JCR（Java Content Repository）。