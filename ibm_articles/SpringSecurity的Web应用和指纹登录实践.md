# Spring Security 的 Web 应用和指纹登录实践
Spring Security 的架构设计、核心组件以及在 Web 应用中的开发方式

**标签:** Spring,Web 开发

[原文链接](https://developer.ibm.com/zh/articles/wa-spring-security-web-application-and-fingerprint-login/)

刘 少飞

更新: 2018-12-04 \| 发布: 2018-11-06

* * *

## 前言

Java 开发人员在解决 Web 应用安全相关的问题时，通常会采用两个非常流行的安全框架，Shiro 和 Spring Security。Shiro 配置简单，上手快，满足一般应用的安全需求，但是功能相对单一。Spring Security 安全粒度细，与 Spring Framework 无缝集成，满足绝大多数企业级应用的安全需求，但是配置复杂，学习曲线陡峭。

Spring Security 相对 Shiro 功能强大，并且 Spring Framework，Spring Boot，Spring Cloud 对 Spring Security 的支持更加友好 (毕竟是 “亲儿子”)。本文将介绍 Spring Security 的架构设计、核心组件，在 Web 应用中的开发方式，最后以一个指纹登录的实例收尾。

## Spring Security 核心设计

Spring Security 有五个核心组件：SecurityContext、SecurityContextHolder、Authentication、Userdetails 和 AuthenticationManager。下面分别介绍一下各个组件。

### SecurityContext

SecurityContext 即安全上下文，关联当前用户的安全信息。用户通过 Spring Security 的校验之后，SecurityContext 会存储验证信息，下文提到的 Authentication 对象包含当前用户的身份信息。SecurityContext 的接口签名如清单 1 所示:

##### 清单 1\. SecurityContext 的接口签名

```
public interface SecurityContext extends Serializable {
       Authentication getAuthentication();
       void setAuthentication(Authentication authentication);
}

```

Show moreShow more icon

SecurityContext 存储在 SecurityContextHolder 中。

### SecurityContextHolder

SecurityContextHolder 存储 SecurityContext 对象。SecurityContextHolder 是一个存储代理，有三种存储模式分别是：

- MODE\_THREADLOCAL：SecurityContext 存储在线程中。
- MODE\_INHERITABLETHREADLOCAL：SecurityContext 存储在线程中，但子线程可以获取到父线程中的 SecurityContext。
- MODE\_GLOBAL：SecurityContext 在所有线程中都相同。

SecurityContextHolder 默认使用 MODE\_THREADLOCAL 模式，SecurityContext 存储在当前线程中。调用 SecurityContextHolder 时不需要显示的参数传递，在当前线程中可以直接获取到 SecurityContextHolder 对象。但是对于很多 C 端的应用（音乐播放器，游戏等等），用户登录完毕，在软件的整个生命周期中只有当前登录用户，面对这种情况 SecurityContextHolder 更适合采用 MODE\_GLOBAL 模式，SecurityContext 相当于存储在应用的进程中，SecurityContext 在所有线程中都相同。

### Authentication

Authentication 即验证，表明当前用户是谁。什么是验证，比如一组用户名和密码就是验证，当然错误的用户名和密码也是验证，只不过 Spring Security 会校验失败。Authentication 接口签名如清单 2 所示:

##### 清单 2\. Authentication 的接口签名

```
public interface Authentication extends Principal, Serializable {
       Collection<? extends GrantedAuthority> getAuthorities();
       Object getCredentials();
       Object getDetails();
       Object getPrincipal();
       boolean isAuthenticated();
       void setAuthenticated(boolean isAuthenticated);
}

```

Show moreShow more icon

Authentication 是一个接口，实现类都会定义 authorities，credentials，details，principal，authenticated 等字段，具体含义如下：

- `getAuthorities`: 获取用户权限，一般情况下获取到的是用户的角色信息。
- `getCredentials`: 获取证明用户认证的信息，通常情况下获取到的是密码等信息。
- `getDetails`: 获取用户的额外信息，比如 IP 地址、经纬度等。
- `getPrincipal`: 获取用户身份信息，在未认证的情况下获取到的是用户名，在已认证的情况下获取到的是 UserDetails (暂时理解为，当前应用用户对象的扩展)。
- `isAuthenticated`: 获取当前 Authentication 是否已认证。
- `setAuthenticated`: 设置当前 Authentication 是否已认证。

在验证前，principal 填充的是用户名，credentials 填充的是密码，detail 填充的是用户的 IP 或者经纬度之类的信息。通过验证后，Spring Security 对 Authentication 重新注入，principal 填充用户信息（包含用户名、年龄等）, authorities 会填充用户的角色信息，authenticated 会被设置为 true。重新注入的 Authentication 会被填充到 SecurityContext 中。

### UserDetails

UserDetails 提供 Spring Security 需要的用户核心信息。UserDetails 的接口签名如清单 3 所示:

##### 清单 3\. UserDetails 的接口签名

```
public interface UserDetails extends Serializable {
       Collection<? extends GrantedAuthority> getAuthorities();
       String getPassword();
       String getUsername();
       boolean isAccountNonExpired();
       boolean isAccountNonLocked();
       boolean isCredentialsNonExpired();
       boolean isEnabled();
}

```

Show moreShow more icon

UserDetails 用 `isAccountNonExpired`, `isAccountNonLocked` ， `isCredentialsNonExpired` ， `isEnabled` 表示用户的状态（与下文中提到的 `DisabledException` ， `LockedException` ， `BadCredentialsException` 相对应），具体含义如下：

- `getAuthorites` ：获取用户权限，本质上是用户的角色信息。
- `getPassword`: 获取密码。
- `getUserName`: 获取用户名。
- `isAccountNonExpired`: 账户是否过期。
- `isAccountNonLocked`: 账户是否被锁定。
- `isCredentialsNonExpired`: 密码是否过期。
- `isEnabled`: 账户是否可用。

UserDetails 也是一个接口，实现类都会继承当前应用的用户信息类，并实现 UserDetails 的接口。假设应用的用户信息类是 User，自定义的 CustomUserdetails 继承 User 类并实现 UserDetails 接口。

### AuthenticationManager

AuthenticationManager 负责校验 Authentication 对象。在 AuthenticationManager 的 authenticate 函数中，开发人员实现对 Authentication 的校验逻辑。如果 authenticate 函数校验通过，正常返回一个重新注入的 Authentication 对象；校验失败，则抛出 AuthenticationException 异常。authenticate 函数签名如清单 4 所示:

##### 清单 4\. authenticate 函数签名

```
Authentication authenticate(Authentication authentication)throws AuthenticationException;

```

Show moreShow more icon

AuthenticationManager 可以将异常抛出的更加明确：

- 当用户不可用时抛出 `DisabledException` 。
- 当用户被锁定时抛出 `LockedException` 。
- 当用户密码错误时抛出 `BadCredentialsException` 。

重新注入的 Authentication 会包含当前用户的详细信息，并且被填充到 SecurityContext 中，这样 Spring Security 的验证流程就完成了，Spring Security 可以识别到 “你是谁”。

### 基本校验流程示例

下面采用 Spring Security 的核心组件写一个最基本的用户名密码校验示例，如清单 5 所示:

##### 清单 5\. Spring Security 核心组件伪代码

```
AuthenticationManager amanager = new CustomAuthenticationManager();
Authentication namePwd = new CustomAuthentication("name”, "password”);
try {
       Authentication result = amanager.authenticate(namePwd);
       SecurityContextHolder.getContext.setAuthentication(result);
} catch(AuthenticationException e) {
       // TODO 验证失败
}

```

Show moreShow more icon

Spring Security 的核心组件易于理解，其基本校验流程是: 验证信息传递过来，验证通过，将验证信息存储到 SecurityContext 中；验证失败，做出相应的处理。

## Spring Security 在 Web 中的设计

Spring Security 的一个常见应用场景就是 Web。下面讨论 Spring Security 在 Web 中的使用方式。

### Spring Security 最简登录实例

Spring Security 在 Web 中的使用相对要复杂一点，会涉及到很多组件。现在给出自定义登录的伪代码，如清单 6 所示。您可以 [点击这里](https://github.com/springAppl/rachel/blob/master/src/main/java/org/security/web/RachelApplication.java) ，查看完整的代码。

##### 清单 6\. Web 登录伪代码

```
@Controller
public class UserController {

       @PostMapping("/login”)
       public void login(String name, String password){
              matchNameAndPassword(name, password);
              User user = getUser(name);
              Authentication auth = new CustomAuthentication(user, password);
              auth.setAuthenticated(true);
              SecurityContextHolder.getContext.setAuthentication(auth);
       }
}

```

Show moreShow more icon

观察代码会发现，如果用 Spring Security 来集成已存在的登录逻辑，真正和 Spring Security 关联的代码只有短短 3 行。验证逻辑可以不经过 AuthenticationManager，真正需要做的就是把经过验证的用户信息注入到 Authentication 中，并将 Authentication 填充到 SecurityContext 中。在实际情况中，登录逻辑的确可以这样写，尤其是已经存在登录逻辑的时候，通常会这样写。这样写虽然方便，但是不符合 Spring Security 在 Web 中的架构设计。

下面视频中会介绍已存在的项目如何与 Spring Security 进行集成，要求对已存在的登录验证逻辑不变，但可以使用 Spring Security 的优秀特性和功能。

### Spring Security 在 Web 中的核心组件

下面介绍在 Web 环境中 Spring Security 的核心组件。

#### FilterChainProxy

FilterChaniProxy 是 FilterChain 代理。FilterChain 维护了一个 Filter 队列，这些 Filter 为 Spring Security 提供了强大的功能。一个很常见的问题是：Spring Security 在 Web 中的入口是哪里？答案是 Filter。Spring Security 在 Filter 中创建 Authentication 对象，并调用 AuthenticationManager 进行校验。Spring Security 选择 Filter，而没有采用上文中 Controller 的方式有以下优点。Spring Security 依赖 J2EE 标准，无需依赖特定的 MVC 框架。另一方面 Spring MVC 通过 Servlet 做请求转发，如果 Spring Security 采用 Servlet，那么 Spring Security 和 Spring MVC 的集成会存在问题。FilterChain 维护了很多 Filter，每个 Filter 都有自己的功能，因此在 Spring Security 中添加新功能时，推荐通过 Filter 的方式来实现。

#### ProviderManager

ProviderManager 是 AuthenticationManager 的实现类。ProviderManager 并没有实现对 Authentication 的校验功能，而是采用代理模式将校验功能交给 AuthenticationProvider 去实现。这样设计是因为在 Web 环境中可能会支持多种不同的验证方式，比如用户名密码登录、短信登录、指纹登录等等，如果每种验证方式的代码都写在 ProviderManager 中，想想都是灾难。因此为每种验证方式提供对应的 AuthenticationProvider，ProviderManager 将验证任务代理给对应的 AuthenticationProvider，这是一种不错的解决方案。在 ProviderManager 中可以找到以下代码，如清单 7 所示:

##### 清单 7\. ProviderManager 代码片段

```
private List<AuthenticationProvider> providers;
public Authentication authenticate(Authentication authentication)
                      throws AuthenticationException {
       ......
       for (AuthenticationProvider provider : getProviders()) {
              if (!provider.supports(toTest)) {
                      continue;
              }
              try {
                      result = provider.authenticate(authentication);
                      if (result != null) {
                             copyDetails(authentication, result);
                             break;
                      }
              }
       }
}

```

Show moreShow more icon

ProviderManager 维护了一个 AuthenticationProvider 队列。当 Authentication 传递进来时，ProviderManager 通过 supports 函数查找支持校验的 AuthenticationProvider。如果没有找到支持的 AuthenticationProvider 将抛出 `ProviderNotFoundException` 异常。

#### AuthenticationProvider

AuthenticationProvider 是在 Web 环境中真正对 Authentication 进行校验的组件。其接口签名如清单 8 所示:

##### 清单 8\. AuthenticationProvider 的接口签名

```
public interface AuthenticationProvider {
       Authentication authenticate(Authentication authentication)
                      throws AuthenticationException;
       boolean supports(Class<?> authentication);
}

```

Show moreShow more icon

其中，authenticate 函数用于校验 Authentication 对象；supports 函数用于判断 provider 是否支持校验 Authentication 对象。

当应用添加新的验证方式时，验证逻辑需要写在对应 AuthenticationProvider 中的 authenticate 函数中。验证通过返回一个重新注入的 Authentication，验证失败抛出 `AuthenticationException` 异常。

### Spring Security 在 Web 中的认证示例

下面的视频中会介绍采用 Spring Security 提供的 `UsernamePasswordAuthenticationFilter` 实现登录验证。

下面以用户名密码登录为例来梳理 Spring Security 在 Web 中的认证流程。上文提到 Spring Security 是以 Filter 来作为校验的入口点。在用户名密码登录中对应的 Filter 是 UsernamePasswordAuthenticationFilter。attemptAuthentication 函数会执行调用校验的逻辑。在 attemptAuthentication 函数中，可以找到以下代码，如清单 9 所示：

##### 清单 9\. attemptAuthentication 函数代码片段

```
public Authentication attemptAuthentication(HttpServletRequest request,HttpServletResponse
response) throws AuthenticationException {
       ......
       UsernamePasswordAuthenticationToken authRequest = new
UsernamePasswordAuthenticationToken(username, password);
       setDetails(request, authRequest);
       return this.getAuthenticationManager().authenticate(authRequest);
}

```

Show moreShow more icon

attemptAuthentication 函数会调用 AuthenticationManager 执行校验逻辑，并获取到重新注入后的 Authentication。在 UsernamePasswordAuthenticationFilter 父类 AbstractAuthenticationProcessingFilter 的 successfulAuthentication 函数中发现以下代码，如清单 10 所示:

##### 清单 10\. successAuthentication 函数

```
protected void successfulAuthentication(HttpServletRequest request,
                      HttpServletResponse response, FilterChain chain, Authentication
authResult)throws IOException, ServletException {
       ......         SecurityContextHolder.getContext().setAuthentication(authResult);
       ......
}

```

Show moreShow more icon

successfulAuthentication 函数会把重新注入的 Authentication 填充到 SecurityContext 中，完成验证。

在 Web 中，AuthenticationManager 的实现类 ProviderManager 并没有实现校验逻辑，而是代理给 AuthenticationProvider, 在用户名密码登录中就是 DaoAuthenticationProvider。DaoAuthenticationProvider 主要完成 3 个功能：获取 UserDetails、校验密码、重新注入 Authentication。在 authenticate 函数中发现以下代码，如清单 11 所示:

##### 清单 11\. DaoAuthenticationProvider.authenticate 函数签名

```
public Authentication authenticate(Authentication authentication)
                      throws AuthenticationException {
       ......
       // 获取 UserDetails
       UserDetails user = this.userCache.getUserFromCache(username);
       if (user == null) {
              cacheWasUsed = false;
              try {
                      user = retrieveUser(username,
       (UsernamePasswordAuthenticationToken) authentication);
              }
              ......
       }
       ......
       try {
              ......
              //校验密码
              additionalAuthenticationChecks(
                      user,
              (UsernamePasswordAuthenticationToken) authentication
              );
       }
       ......
       // 从新注入 Authentication
       return createSuccessAuthentication(
                      principalToReturn,
                      authentication,
                      user
               );
}

```

Show moreShow more icon

首先从 userCache 缓存中查找 UserDetails, 如果缓存中没有获取到，调用 retrieveUser 函数获取 UserDetails。retrieveUser 函数签名如清单 12 所示:

##### 清单 12\. retrieveUser 函数签名

```
protected final UserDetails retrieveUser(String username,
       UsernamePasswordAuthenticationToken authentication)
                      throws AuthenticationException {
       UserDetails loadedUser;
       try {
               loadedUser = this.getUserDetailsService().loadUserByUsername(username);
       }
       ......
       return loadedUser;
}

```

Show moreShow more icon

retrieveUser 函数调用 UserDetailsService 获取 UserDetails 对象。UserDetailsService 接口签名如清单 13 所示：

##### 清单 13\. UserDetailsService 接口签名

```
public interface UserDetailsService {
       UserDetails loadUserByUsername(String username) throws UsernameNotFoundException;
}

```

Show moreShow more icon

UserDetailsService 非常简单，只有一个 loadUserByUserName 函数，函数参数虽然名为 username，但只要是用户的唯一标识符即可。下面是基于数据库存储的简单示例, 如清单 14 所示:

##### 清单 14\. CustomUserDetailsService 类签名

```
public class CustomUserDetailsService implements UserDetailsService {

    @Autowired
    private UserDao userDao;

    @Override
    public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException
{
        User user = userDao.findByName(username);
        if(Objects.isNull(user)) {
           throw new UsernameNotFoundException();
        }
        UserDetails details = new CustomUserDetails(user);
        return details;
    }
}

```

Show moreShow more icon

调用 UserDao 获取 User 对象，将 User 对象包装成 UserDetails 对象。如果没有找到 User 对象，需要抛出 `UsernameNotFoundException` 异常。

DaoAuthenticationProvider 密码校验调用 additionalAuthenticationChecks 函数，具体通过 PasswordEncoder 比对用户输入的密码和存储在应用中的密码是否相等，如果不相等，抛出 `BadCredentialsException` 异常。

DaoAuthenticationProvider 对 Authentication 对象的重新注入通过调用 createSuccessAuthentication 函数, 如清单 15 所示:

##### 清单 15\. createSuccessAuthentication 函数签名

```
protected Authentication createSuccessAuthentication(Object principal,
                      Authentication authentication, UserDetails user) {
        UsernamePasswordAuthenticationToken result = new
        UsernamePasswordAuthenticationToken(
               principal,
               authentication.getCredentials(),
               authoritiesMapper.mapAuthorities(user.getAuthorities())
        );
        result.setDetails(authentication.getDetails());
        return result;
}

```

Show moreShow more icon

以上就是 Spring Security 在 Web 环境中对于用户名密码校验的整个流程，简言之：

1. UsernamePasswordAuthenticationFilter 接受用户名密码登录请求，将 Authentication 传递给 ProviderManager 进行校验。
2. ProviderManager 将校验任务代理给 DaoAuthenticationProvider。
3. DaoAuthenticationProvider 对 Authentication 的用户名和密码进行校验，校验通过后返回重新注入的 Authentication 对象。
4. UsernamePasswordAuthenticationFilter 将重新注入的 Authentication 对象填充到 SecurityContext 中。

## 指纹登录实践

指纹登录和用户名密码登录区别很小，只是将密码换成了指纹特征值。下面采用 Spring Security 推荐写法 Filter-AuthenticationProvider 的形式来定义相关组件以实现指纹登录。完整的项目地址： [https://github.com/springAppl/rachel](https://github.com/springAppl/rachel) 。

### FingerPrintToken

FingerPrintToken 增加 name 和 fingerPrint 字段，分别代表用户名和指纹特征值，如清单 16 所示:

##### 清单 16\. FingerPrintToken 函数签名

```
public class FingerPrintToken implements Authentication {
       private String name;
       private String fingerPrint;
       ......
}

```

Show moreShow more icon

### FingerPrintFilter

FingerPrintFilter 处理指纹登录请求，调用 AuthenticationManager 进行验证，验证通过后调用 SecurityContextHolder 将重新注入的 Authentication 填充到 SecurityContext 中，如清单 17 所示:

##### 清单 17\. doFilter 函数签名

```
public void doFilter(ServletRequest servletRequest, ServletResponse servletResponse,
FilterChain filterChain) throws IOException, ServletException {
       if (Objects.equals(httpServletRequest.getRequestURI(), "/api/finger-print")) {
               // 调用 AuthenticationManager， 并填充 SecurityContext
        }
}

```

Show moreShow more icon

### FingerPrintProvider

FingerPrintProvider 负责处理 FingerPrintToken，需要在 supports 函数中支持处理 FingerPrintToken。authenticate 函数负责 UserDetails 获取，指纹校验，FingerPrintToken 的重新注入。

### FingerPrintUserDetails

FingerPrintUserDetails 继承 User 并实现 UserDetails 的方法，应用的用户信息可以加载到 Spring Security 中使用。

### FingerPrintUserDetailsService

FingerPrintUserDetailsService 获取 FingerUserDetails。通过 UserDao 查找到 User，并将 User 转换为 Spring Security 可识别 UserDetails。

### SecurityConfig

SecurityConfig 继承 WebSecurityConfigurerAdapter，需要定义 Spring Security 配置类。Spring Security 的配置不是本文的重点，配置时只需要注意以下几点：

1. 将 FingerPrintFilter、FingerPrintProvider 添加进去。
2. 将 FingerPrintFilter 的执行顺序放置在 SecurityContextPersistenceFilter 之后即可。Spring Security 维护了一个 Filter 的 list，因此每个 Filter 是有顺序的。
3. 将 “/api/test” 请求设置为用户验证成功后才允许方问。

配置代码在 configure 函数中，如清单 18 所示:

##### 清单 18\. configure 函数

```
protected void configure(HttpSecurity http) throws Exception {
       http
              .userDetailsService(userDetailsService())
              .addFilterAfter(fingerPrintFilter(), SecurityContextPersistenceFilter.class)
              .authenticationProvider(fingerPrintProvider())
              .authorizeRequests()
              .mvcMatchers(HttpMethod.GET, "/api/test").authenticated()
}

```

Show moreShow more icon

## 结束语

在 Web 时代，用户和应用的耦合度越来越高，应用中存储了大量用户的私密信息。随着各种用户信息泄露事件的爆发，安全成为了 Web 应用重要的一个环。Spring Security 由于其强大的功能和 Spring Framework 的高度集成，赢得了开发人员的青睐。本文对 Spring Security 的架构设计与核心组件进行了深入浅出的介绍，分析了 Spring Security 在 Web 应用的集成方式，并展示了一个指纹登录的实例。