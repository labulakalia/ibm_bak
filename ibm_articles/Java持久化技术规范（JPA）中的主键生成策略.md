# Java 持久化技术规范（JPA）中的主键生成策略
Table，Sequence，Identity，Auto 四种主键生成策略

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-lo-jpaprimarykey/)

王汉敏

发布: 2012-05-07

* * *

## Table 策略 (Table strategy)

这种策略中，持久化引擎 (persistence engine) 使用关系型数据库中的一个表 (Table) 来生成主键。这种策略可移植性比较好，因为所有的关系型数据库都支持这种策略。不同的 J2EE 应用服务器使用不同的持久化引擎。

下面用一个例子来说明这种表生成策略的使用：

##### 清单 1.Table 生成策略

```
@Entity
public class PrimaryKey_Table {

@TableGenerator(name = "PK_SEQ",
table = "SEQUENCE_TABLE",
                pkColumnName  = "SEQUENCE_NAME",
                valueColumnName  = "SEQUENCE_COUNT")

@Id
    @GeneratedValue(strategy =GenerationType.TABLE,generator="PK_SEQ")
    private Long id;

//Getters and Setters
//为了方便，类里面除了一个必需的主键列，没有任何其他列，以后类似

}

```

Show moreShow more icon

首先，清单 1 中使用 @javax.persistence.TableGenerator 这个注解来指定一个用来生成主键的表 (Table)，这个注解可以使用在实体类上，也可以像这个例子一样使用在主键字段上。

其中，在这个例子中，name 属性”PK\_SEQ” 标示了这个生成器，也就是说这个生成器的名字是 PK\_SEQ。这个 Table 属性标示了用哪个表来存贮生成的主键，在这个例子中，用” SEQUENCE\_TABLE” 来存储主键，数据库中有对应的 SEQUENCE\_TABLE 表。其中 pkColumnName 属性用来指定的是生成器那个表中的主键，也就是 SEQUENCE\_TABLE 这个表的主键的名字。属性 valueColumnName 指定列是用来存储最后生成的那个主键的值。

也可以使用持久化引擎提供的缺省得 Table，例如：

##### 清单 2\. 使用确省的表生成器

```
public class PK implements Serializable {
    private static final long serialVersionUID = 1L;

    @Id
    @GeneratedValue(strategy = GenerationType.TABLE)

    private Long id;
// Getters and Setters
}

```

Show moreShow more icon

不同的持久化引擎有不同的缺省值，在 glass fish 中，Table 属性的缺省值是 SEQUENCE, pkColumnName 属性缺省值是 SEQ\_NAME,，valueColumnName 属性的缺省值是 SEQ\_COUNT

## Sequence 策略

一些数据库，比如 Oralce，有一种内置的叫做”序列” （sequence）的机制来生成主键。为了调用这个序列，需要使用 @javax.persistence.SequenceGenerator 这个注解。

例如

##### 清单 3.sequence 策略生成主键

```
@Entity
public class PK_Sequence implements Serializable {
    private static final long serialVersionUID = 1L;
    @SequenceGenerator(name="PK_SEQ_TBL",sequenceName="PK_SEQ_NAME")
    @Id
    @GeneratedValue(strategy = GenerationType.SEQUENCE,generator="PK_SEQ_TBL")
    private Long id;
// Getters and Setters
}

```

Show moreShow more icon

其中的 `@javax.persistence.SequenceGenerator` 定义如下：

##### 清单 4.@SequenceGenerator 注解的定义

```
@Target(value = {ElementType.TYPE, ElementType.METHOD, ElementType.FIELD})
@Retention(value = RetentionPolicy.RUNTIME)
public @interface SequenceGenerator {

    public String name();

    public String sequenceName() default "";

    public String catalog() default "";

    public String schema() default "";

    public int initialValue() default 1;

    public int allocationSize() default 50;
}

```

Show moreShow more icon

从定义中可以看出这个注解可以用在类上，也可以用在方法和字段上，其中 name 属性指定的是所使用的生成器；sequenceName 指定的是数据库中的序列；initialValue 指定的是序列的初始值，和 @TableGenerator 不同是它的缺省值 1；allocationSize 指定的是持久化引擎 (persistence engine) 从序列 (sequence) 中读取值时的缓存大小，它的缺省值是 50。

## Identity 策略

一些数据库，用一个 Identity 列来生成主键，使用这个策略生成主键的时候，只需要在 @GeneratedValue 中用 strategy 属性指定即可。如下所示：

##### 清单 5.strategy 策略生成主键

```
@Entity
public class PK_Identity implements Serializable {
    private static final long serialVersionUID = 1L;
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
// Getters and Setters
}

```

Show moreShow more icon

## Auto 策略

使用 AUTO 策略就是将主键生成的策略交给持久化引擎 (persistence engine) 来决定，由它自己选择从 Table 策略，Sequence 策略和 Identity 策略三种策略中选择合适的主键生成策略。不同的持久化引擎 (persistence engine) 使用不同的策略，在 galss fish 中使用的是 Table 策略。

使用 AUTO 策略时，我们可以显示使用，如：

##### 清单 6.Auto 策略生成主键

```
@Entity
public class PK_Auto implements Serializable {
    private static final long serialVersionUID = 1L;
    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private Long id;

// Getters and Setters
    }

```

Show moreShow more icon

或则只使用：

@Generated Value

或者干脆什么都不写，因为缺省得主键生成策略就是 AUTO。

## 复合主键

在对象关系映射模型中，使用单独的一个字段作为主键是一种非常好的做法，但是在实际应用中，经常会遇到复合主键的问题，就是使用两个或两个以上的字段作为主键。比如，在一些历史遗留的数据库表中，经常出现复合主键的问题，为了解决这种问题，JPA2.0 中采用的 @EmbeddedId 和 @IdClass 两种方法解决这种问题。它们都需要将用于主键的字段单独放在一个主键类 (primary key class) 里面，并且该主键类必须重写 equals () 和 hashCode () 方法，必须实现 Serializable 接口，必须拥有无参构造函数。

### @EmbeddedId 复合主键

清单 7 中的 NewsId 类被用做主键类，它用 @Embeddable 注解进行了注释，说明这个类可以嵌入到其他类中。之外这个类中还重写了 hashCode () 和 equals () 方法， 因为这个类中的两个属性要用作主键，必须有一种判定它们是否相等并且唯一的途径。

##### 清单 7.@EmbeddedId 中的主键类

```
@Embeddable
public class NewsId implements Serializable {
    private static final long serialVersionUID = 1L;

    private String title;
    private String language;

    public String getLanguage() {
        return language;
    }

    public void setLanguage(String language) {
        this.language = language;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    @Override
    public boolean equals(Object obj) {
        if (obj == null) {
            return false;
        }
        if (getClass() != obj.getClass()) {
            return false;
        }
        final NewsId other = (NewsId) obj;
if ((this.title == null) ? (other.title != null) : !this.title.equals(other.title)) {
            return false;
        }
if ((this.language == null) ? (other.language != null) : !this.language.equals(
    other.language)) {
            return false;
        }
        return true;
    }

    @Override
    public int hashCode() {
        int hash = 5;
        hash = 41 * hash + (this.title != null ? this.title.hashCode() : 0);
        hash = 41 * hash + (this.language != null ? this.language.hashCode() : 0);
        return hash;
    }

}

```

Show moreShow more icon

清单 8 中的 News 类使用了清单 7 中定义的主键类，可以看倒非常简单，只需要使用 @EmbeddedId 指定主键类即可。

##### 清单 8.News 实体类使用定义好的主键类

```
@Entity
public class News implements Serializable {
    private static final long serialVersionUID = 1L;

    @EmbeddedId
    private NewsId id;
    private String content;

// Getters and Setters

}

```

Show moreShow more icon

清单 9 是持久化后生成的数据库表的结构，可以看出来这个表如我们预想的一样是 Title 和 Language 的联合主键。

##### 清单 9\. 使用主键类生成的表结构

```
CREATE TABLE `news` (
`CONTENT` varchar(255) default NULL,
`TITLE` varchar(255) NOT NULL,
`LANGUAGE` varchar(255) NOT NULL,
PRIMARY KEY  (`TITLE`,`LANGUAGE`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

```

Show moreShow more icon

## IdClass 复合主键

IdClass 这种复合主键策略，在主键类上和 Embeddable 这种复合主键策略稍有不同。如清单 10，这个策略中的主键类不需要使用任何注解 (annotation)，但是仍然必须重写 hashCode() 和 equals() 两个方法。其实也就是将 Embeddable 这种复合主键策略中的主键类的 @Embeddable 注解去掉就可以了。

##### 清单 10\. IdClass 复合主键策略中的主键类

```
public class NewsId implements Serializable {
    private static final long serialVersionUID = 1L;

    private String title;
    private String language;

    public String getLanguage() {
        return language;
    }

    public void setLanguage(String language) {
        this.language = language;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    @Override
    public boolean equals(Object obj) {
        if (obj == null) {
            return false;
        }
        if (getClass() != obj.getClass()) {
            return false;
        }
        final NewsId other = (NewsId) obj;
if ((this.title == null) ? (other.title != null) : !this.title.equals(other.title)) {
            return false;
        }
if ((this.language == null) ? (other.language != null) : !this.language.equals(
    other.language)) {
            return false;
        }
        return true;
    }

    @Override
    public int hashCode() {
        int hash = 5;
        hash = 41 * hash + (this.title != null ? this.title.hashCode() : 0);
        hash = 41 * hash + (this.language != null ? this.language.hashCode() : 0);
        return hash;
    }

}

```

Show moreShow more icon

从清单 11 中可以看出这个 News 实体类首先使用 @IdClass (NewsId.class) 注解指定了主键类。同时在类中也复写了主键类中的两个字段，并分别用 @Id 进行了注解。

##### 清单 11\. IdClass 策略中使用复合主键的 News 类

```
@Entity
@IdClass(NewsId.class)
public class News implements Serializable {
    private static final long serialVersionUID = 1L;

    @Id
    private String title;
    @Id
    private String language;
private String content;
// Getters and Setters
}

```

Show moreShow more icon

从清单 12 中可以看出，两种复合主键的映射策略，持久化后映射到数据库中的表的结构是相同的，一个明显的区别就是在查询的时候稍有不同。如在使用 @EmbeddableId 策略的时候，要使用如下查询语句：

```
Select n.newsId.title from news n

```

Show moreShow more icon

而使用 @IdClass 策略的时候，要使用如下查询语句：

```
Select n.title from news n

```

Show moreShow more icon

另外一点就是使用 @IdClass 这种策略的时候，在复写主键类中的字段的时候务必要保证和主键类中的定义完全一样。

##### 清单 12\. @IdClass 策略生成的表结构

```
CREATE TABLE `news` (
`CONTENT` varchar(255) default NULL,
`TITLE` varchar(255) NOT NULL,
`LANGUAGE` varchar(255) NOT NULL,
PRIMARY KEY  (`TITLE`,`LANGUAGE`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

```

Show moreShow more icon

## 结束语

Java EE 项目开发中的持久层，虽然具体的实现方式，也就是持久化引擎会随着你选择的 Java EE 服务器的不同而有所不同，但是在 JPA(java persistence API) 这个规范之下，每个实体的主键生成策略却只有上面几种，也就是说我们主要掌握了上面几种主键生成策略，就可以在以后 Java EE 项目持久层开发中以不变应万变的姿态来面对纷繁复杂的具体情况了。