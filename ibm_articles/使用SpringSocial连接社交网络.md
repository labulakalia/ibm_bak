# 使用 Spring Social 连接社交网络
详细介绍 Spring Social 项目以及如何在项目中使用它

**标签:** Java

[原文链接](https://developer.ibm.com/zh/articles/j-lo-spring-social/)

成富

发布: 2014-07-22

* * *

社交网络已经成为目前人们生活中不可或缺的一部分。流行的社交网站，如 Facebook、LinkedIn、Twitter 和新浪微博等都聚集了非常高的人气，成为很多人每天必上的网站。大量的信息被分享到社交网络中。对于一个应用来说，提供对社交网络的支持是很有必要的。应用可以通过与社交网络的集成来迅速积累人气。与社交网络的集成主要有两个方面的功能：第一个是比较简单的用户登录集成，即允许用户使用社交网络网站上的已有账户来登录应用。这样做的好处是可以省去要用户重新注册的流程，同时也可以与用户已有的社交网络建立连接；第二个是深度的集成方式，即允许用户把应用中的相关内容分享到社交网络中。这样做的好处是可以保持用户的粘性，同时推广应用本身。

在应用中集成社交网络并不是一件复杂的工作。很多社交网络服务都提供开放 API 允许第三方应用来进行集成。只需要参考相关的文档就可以完成。不过集成本身需要一定的工作量。所需要集成的网站越多，相应的工作量越大。Spring Social 作为 Spring 框架组件家族中的一部分，提供了对社交网络服务的抽象，并提供了相关的 API 来进行操作。

使用 Spring Social 的最大好处在于它已经提供了对主流社交网站的支持，只需要简单配置即可。对于一些不太常用的社交网站，也可以从社区中找到相应的组件。本文以一个示例应用来说明 Spring Social 的用法。

## 示例应用

该示例允许用户登录社交网站查看他们的好友信息。这是一个使用 Spring MVC 框架构建的 Web 应用。该应用使用嵌入式的 H2 数据库来存储数据。前端使用 JSP、JSTL 和 Spring 提供的标签库，利用 Apache Tiles 进行页面布局。使用 Spring Security 管理用户认证。页面的样式使用 Twitter Bootstrap 库。项目使用 Gradle 进行构建。前端库通过 Bower 进行管理。

示例所提供的功能比较简单，主要包括登录和注册页面，以及用户的首页。在登录时，用户可以选择使用第三方社交网络的账号。登录之后可以查看社交网络的好友信息。

## 使用社交网络已有账号登录

与社交网络集成的最基本的功能是允许用户使用已有的账号进行登录。用户除了注册新的账号之外，还可以使用已有的其他网站的账号来登录。这种方式通常称为连接第三方网站。对于用户来说，通常的场景是在注册或登录页面，选择第三方社交网络来进行连接。然后跳转到第三方网站进行登录和授权。完成上述步骤之后，新的用户会被创建。在连接过程中，应用可以从第三方社交网站获取到用户的概要信息，来完成新用户的创建。当新用户创建完成之后，下次用户可以直接使用第三方社交网络的账号进行登录。

### 基本配置

Spring Social 提供了处理第三方网站用户登录和注册的 Spring MVC 控制器（Controller）的实现，可以自动完成使用社交网络账号登录和注册的功能。在使用这些控制器之前，需要进行一些基本的配置。Spring Social 的配置是通过实现 org.springframework.social.config.annotation.SocialConfigurer 接口来完成的。 [代码清单 1](#_清单1.使用SpringSocial的基本配置) 给出了实现 SocialConfigurer 接口的 SocialConfig 类的部分内容。注解 “@EnableSocial” 用来启用 Spring Social 的相关功能。注解 “@Configuration” 表明该类也同样包含 Spring 框架的相关配置信息。SocialConfigurer 接口有 3 个方法需要实现：

- addConnectionFactories：该回调方法用来允许应用添加需要支持的社交网络对应的连接工厂的实现。
- getUserIdSource：该回调方法返回一个 org.springframework.social.UserIdSource 接口的实现对象，用来惟一标识当前用户。
- getUsersConnectionRepository：该回调方法返回一个 org.springframework.social.connect.UsersConnectionRepository 接口的实现对象，用来管理用户与社交网络服务提供者之间的对应关系。

具体到示例应用来说，addConnectionFactories 方法的实现中添加了由 org.springframework.social.linkedin.connect.LinkedInConnectionFactory 类表示的 LinkedIn 的连接工厂实现。getUserIdSource 方法的实现中通过 Spring Security 来获取当前登录用户的信息。getUsersConnectionRepository 方法中创建了一个基于数据库的 JdbcUsersConnectionRepository 类的实现对象。LinkedInConnectionFactory 类的创建方式体现了 Spring Social 的强大之处。只需要提供在 LinkedIn 申请的 API 密钥，就可以直接使用 LinkedIn 所提供的功能。与 OAuth 相关的细节都被 Spring Social 所封装。

##### 清单 1\. 使用 Spring Social 的基本配置

```
@Configuration
@EnableSocial
public class SocialConfig implements SocialConfigurer {
@Inject
private DataSource dataSource;
@Override
public void addConnectionFactories(ConnectionFactoryConfigurer cfConfig, Environment env) {
cfConfig.addConnectionFactory(new LinkedInConnectionFactory(env.getProperty("linkedin.consumerKey"), env.getProperty("linkedin.consumerSecret")));
}
@Override
public UserIdSource getUserIdSource() {
return new UserIdSource() {
@Override
public String getUserId() {
Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
if (authentication == null) {
throw new IllegalStateException("Unable to get a ConnectionRepository: no user signed in");
}
return authentication.getName();
}
};
}
@Override
public UsersConnectionRepository getUsersConnectionRepository(ConnectionFactoryLocator connectionFactoryLocator) {
return new JdbcUsersConnectionRepository(dataSource, connectionFactoryLocator, Encryptors.noOpText());
}
}

```

Show moreShow more icon

### 登录控制器

配置了 Spring Social 的下一步是创建相应的 Spring MVC 控制器来处理用户的注册和登录。Spring Social 已经提供相应的控制器实现，只需要创建界面来发出请求即可。首先要创建的是通过第三方社交网络进行登录的控制器。如 [清单 2\. 通过第三方社交网络进行登录的控制器](#清单-2-通过第三方社交网络进行登录的控制器) 所示，只需要创建一个 Spring Bean 就可以创建所需要的控制器。该控制器默认绑定在 “/signin” 的 URL 上面。每个第三方服务提供者都有一个对应的标识符，如 LinkedIn 的标识符是 “linkedin”。服务提供者的标识符通常作为不同控制器 URL 的一部分。比如当用户需要通过 LinkedIn 账户来登录时，对应的 URL 就是 “/signin/linkedin”。如果通过 Twitter 账户来登录，则对应的 URL 是 “/signin/twitter”。

##### 清单 2\. 通过第三方社交网络进行登录的控制器

```
@Bean
public ProviderSignInController providerSignInController(ConnectionFactoryLocator connectionFactoryLocator, UsersConnectionRepository usersConnectionRepository) {
return new ProviderSignInController(connectionFactoryLocator, usersConnectionRepository, new SimpleSignInAdapter(new HttpSessionRequestCache()));
}

```

Show moreShow more icon

在上面的 [清单 2\. 通过第三方社交网络进行登录的控制器](#清单-2-通过第三方社交网络进行登录的控制器) 中，SimpleSignInAdapter 类实现了 org.springframework.social.connect.web.SignInAdapter 接口，其作用是 ProviderSignInController 类和应用本身的登录逻辑之间的桥梁。当用户完成第三方网站的登录之后，通过 SignInAdapter 接口来使得该用户自动登录到应用中。SimpleSignInAdapter 类的实现如 [清单 3\. SimpleSignInAdapter 类的实现](#清单-3-simplesigninadapter-类的实现) 所示。

##### 清单 3\. SimpleSignInAdapter 类的实现

```
public class SimpleSignInAdapter implements SignInAdapter {
private final RequestCache requestCache;
@Inject
public SimpleSignInAdapter(RequestCache requestCache) {
this.requestCache = requestCache;
}
@Override
public String signIn(String localUserId, Connection<?> connection, NativeWebRequest request) {
SignInUtils.signin(localUserId);
return extractOriginalUrl(request);
}
private String extractOriginalUrl(NativeWebRequest request) {
HttpServletRequest nativeReq = request.getNativeRequest(HttpServletRequest.class);
HttpServletResponse nativeRes = request.getNativeResponse(HttpServletResponse.class);
SavedRequest saved = requestCache.getRequest(nativeReq, nativeRes);
if (saved == null) {
return null;
}
requestCache.removeRequest(nativeReq, nativeRes);
removeAutheticationAttributes(nativeReq.getSession(false));
return saved.getRedirectUrl();
}
private void removeAutheticationAttributes(HttpSession session) {
if (session == null) {
return;
}
session.removeAttribute(WebAttributes.AUTHENTICATION_EXCEPTION);
}
}

```

Show moreShow more icon

在 [清单 3\. SimpleSignInAdapter 类的实现](#清单-3-simplesigninadapter-类的实现) 中，SimpleSignInAdapter 类的 signIn 方法调用 SignInUtils 类的 signin 方法来进行登录，其实现见 [清单 4\. SignInUtils 类的 signIn 方法的实现](#清单-4-signinutils-类的-signin-方法的实现) 。signIn 方法的返回值是登录成功之后的跳转 URL。

##### 清单 4\. SignInUtils 类的 signIn 方法的实现

```
public class SignInUtils {
public static void signin(String userId) {
SecurityContextHolder.getContext().setAuthentication(new UsernamePasswordAuthenticationToken(userId, null, null));
}
}

```

Show moreShow more icon

### 前端页面

只需要通过页面向登录控制器的 URL 发出一个 POST 请求，就可以启动通过 LinkedIn 进行登录的过程。如 [清单 5\. 登录页面](#清单-5-登录页面) 中，通过表单提交就可以启动登录过程。用户会被首先转到 LinkedIn 网站的授权页面给应用授权，授权完成之后就可以进行注册。

##### 清单 5\. 登录页面

```
<form id="linkedin-signin-form" action="signin/linkedin" method="POST" class="form-signin" role="form">
<h2 class="form-signin-heading">Or Connect by</h2>
<input type="hidden" name="${_csrf.parameterName}" value="${_csrf.token}"/>
<button type="submit" class="btn btn-primary">LinkedIn</button>
</form>

```

Show moreShow more icon

默认情况下，当第三方网站授权完成之后，用户会被转到 URL”/signup”对应的页面。在这个页面，用户可以补充一些相关的注册信息。与此同时，从第三方网站获取的用户概要信息，如用户的姓名等，可以被预先填充好。该页面在表单提交之后的操作，由 [清单 6\. 用户注册控制器的实现](#清单-6-用户注册控制器的实现) 所示的控制器来处理。

##### 清单 6\. 用户注册控制器的实现

```
@RequestMapping(value="/signup", method=RequestMethod.POST)
public String signup(@Valid @ModelAttribute("signup") SignupForm form, BindingResult formBinding, WebRequest request) {
if (formBinding.hasErrors()) {
return null;
}
Account account = createAccount(form, formBinding);
if (account != null) {
SignInUtils.signin(account.getUsername());
ProviderSignInUtils.handlePostSignUp(account.getUsername(), request);
return "redirect:/";
}
return null;
}

```

Show moreShow more icon

当用户提交注册表单之后，根据用户填写的信息创建相应的账号，然后登录新注册的账号并跳转回首页。

## 使用 API

Spring Social 的另外一个重要的作用在于提供了很多社交网络服务的 API 的封装。当用户通过社交网络的账号登录之后，可以通过相应的 API 获取到相关的信息，或执行一些操作。比如获取到用户的好友信息，或是根据用户的要求发布新的消息。一般来说，很多社交网络服务提供商都有自己的 API 来提供给开发人员使用。不过一般只提供 REST API，并没有 Java 封装。Spring Social 的一个目标是为主流的社交网络提供相应的 Java 封装的 API。对于社区已经有良好的 Java 库的社交网络，如 Twitter，Spring Social 会进行整合与封装；而对于暂时没有 Java 实现的，Spring Social 的组件会提供相应的支持。如示例应用中使用的 LinkedIn 的 Java API 是由 Spring Social 开发的。

在使用 LinkedIn 的 Java API 之前，首先要创建相应的 Spring Bean，如 [清单 7\. 创建 LinkedIn API 对应的 Bean](#清单-7-创建-linkedin-api-对应的-bean) 所示。

##### 清单 7\. 创建 LinkedIn API 对应的 Bean

```
@Bean
@Scope(value="request", proxyMode=ScopedProxyMode.INTERFACES)
public LinkedIn linkedin(ConnectionRepository repository) {
Connection<LinkedIn> connection = repository.findPrimaryConnection(LinkedIn.class);
return connection != null ? connection.getApi() : null;
}

```

Show moreShow more icon

LinkedIn Java API 的核心接口是 org.springframework.social.linkedin.api.LinkedIn。如果当前用户已经连接了 LinkedIn 的账号，就可以通过 ConnectionRepository 接口找到该用户对应的连接的信息，并创建出 API 接口的实现。 [清单 8\. 展示 LinkedIn 上好友信息的控制器的实现](#清单-8-展示-linkedin-上好友信息的控制器的实现) 是展示 LinkedIn 上好友信息的控制器的实现。

##### 清单 8\. 展示 LinkedIn 上好友信息的控制器的实现

```
@Controller
public class LinkedInConnectionsController {
private LinkedIn linkedIn;
@Inject
public LinkedInConnectionsController(LinkedIn linkedIn) {
this.linkedIn = linkedIn;
}
@RequestMapping(value="/connections", method= RequestMethod.GET)
public String connections(Model model) {
model.addAttribute("connections", linkedIn.connectionOperations().getConnections());
return "connections";
}
}

```

Show moreShow more icon

通过 LinkedIn 接口的 connectionOperations 方法可以获取到 org.springframework.social.linkedin.api.ConnectionOperations 接口的实现，然后就可以再通过 getConnections 方法获取到包含好友信息的 org.springframework.social.linkedin.api.LinkedInProfile 类的对象列表。然后就可以在网页上展示，如 [清单 9\. 展示 LinkedIn 好友信息的页面](#清单-9-展示-linkedin-好友信息的页面) 所示。

##### 清单 9\. 展示 LinkedIn 好友信息的页面

```
<ul class="media-list">
<c:forEach items="${connections}" var="connection">
<li class="media">
<a href="" class="pull-left">
<img src="${connection.profilePictureUrl}" alt="" class="media-object">
</a>
<div class="media-body">
<h4 class="media-heading">${connection.firstName} ${connection.lastName}</h4>
<p><c:out value="${connection.headline}" /></p>
</div>
</li>
</c:forEach>
</ul>

```

Show moreShow more icon

## 结束语

在社交网络日益流行的今天，Web 应用都应该考虑增加对社交网络的支持。这样既方便了用户，让他们省去了重新注册的麻烦，又可以增加用户粘性，同时可以通过社交网络进行推广。虽然很多社交网站都提供了开放 API 来让应用使用其服务，但是集成所需的工作量是不小的。Spring Social 作为 Spring 框架中的一部分，为在 Web 应用中集成社交网络提供了良好支持。本文详细介绍了如何通过 Spring Social 来实现通过社交网络账号进行登录以及使用社交网络提供的 API。

## Download

[spring-social-code-sample.zip](http://www.ibm.com/developerWorks/cn/java/j-lo-spring-social/spring-social-code-sample.zip): 代码示例