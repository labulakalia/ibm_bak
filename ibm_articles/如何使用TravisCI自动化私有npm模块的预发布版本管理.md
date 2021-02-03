# 如何使用 Travis CI 自动化私有 npm 模块的预发布版本管理
Manage multiple streams of development with minimal complication

**标签:** DevOps

[原文链接](https://developer.ibm.com/zh/articles/d-private-npm-modules-travis-ci-presence-insights-trs/)

Alex Weidner

发布: 2016-10-17

* * *

## Presence Insights：背景知识

IBM Presence Insights 包含近 30 个库，以两周为周期进行迭代。在这个两周的周期中，这些库的默认分支会随着多个团队交付特性而快速更改。因为我们依靠 npm 将这些库部署到临时环境，所以，以传统方式处理此过程会导致在每个开发周期砍掉许多版本，其中大部分版本都不重要或无关紧要。这产生了过于复杂的版本历史，在突然出现性能退化时通常很难进行分析。

通过 [`npm dist-tags`](https://docs.npmjs.com/cli/dist-tag) ，可以建立和管理多个开发流，而不会让库的版本历史复杂化。 `dist-tags` 是有效的安装目标，这意味着一个名为 `unstable` 的给定 `dist-tag` 可通过 `npm` 使用 `npm install @pi/library@unstable` 命令来安装。这还意味着它是 `package.json` 的有效目标。

当在迭代期间将代码签入到任何分支并通过回归测试时，Travis CI 会基于当前（”最新”）版本来创建一个版本。最新版本包含更多元数据，包括分支名称和时间戳。此版本在一个基于该分支名称的标签下发布。

## 参考资料

此信息会传递到：

- travis.yml
- package.json
- @pidev/publish

### .travis.yml

在签入代码时，来自 PI 记录器的 Travis.yaml 会控制 Travis 运行的编译版本的配置。下面是一个示例：

```
deploy:
    provider: script
    script: npm run travis-publish
    skip_cleanup: true
    on:
        all_branches: true
        node: '0.12'
        tags: false

```

Show moreShow more icon

这段代码声明部署由我们通过 `npm run travis-publish` 提供的一段自定义脚本来管理。 `on:d` 用于定义要发布哪些模块。因为我们希望管理预发布版本，所以我们从所有分支发布，而 **不** 发布 on 标签（发布脚本在每次迭代结束时推送 on 标签）。我们启用了 `skip_cleanup` ，因为我们的自定义发布脚本依赖于一些 Node 依赖项。

### package.json

定义 `npm run travis-publish` 的含义。来自 Presence Insights 记录器的一个 `package.json` 示例如下：

> `"travis-publish":"publish npm --platform travis"`

publish 是一个使用 Node 编写的 CLI，被指定为一个开发依赖项：

> `"@pidev/publish":"0.x"`

### @pidev/publish

定义我们的部署的自定义逻辑，包括利用分支名和时间戳构成元数据。它是一个简单的 [commander](https://npmjs.org/commander) \+ [shelljs](https://npmjs.org/shelljs) node CLI，通过 `publish npm` 公开一个发不到 `npm` 的命令。当 `--platform travis` 存在时，会使用一个访问 Travis 环境变量的例程。

对于 Presence Insights，我们将临时环境命名为 **`YS1-dev`** 。当更改签入到默认分支时，使用标签 `dev` 来匹配临时区域的名称。对于其他团队或情景，推荐使用 `rc` ，以清楚表明默认分支上的所有内容均可发布。

```
// change on master branch, main development stream
if (branch === 'master') {
    sh.echo(ch.blue('creating dev version and publishing'));
    version += '-dev-' + moment().tz('America/New_york').format('YYYYMMDD-HHmm');
    tag = 'dev';
}

```

Show moreShow more icon

版本是利用 `dist-tag` 创建的，时区的格式为 `YYYYMMDD-HHmm` 。例如： `@pilib/logger@1.0.5-dev-20160416-0945` 。

**来自在 Travis 上运行的一次编译的示例输出：**

```
creating dev version and publishing
exec npm version 1.2.11-dev-20160421-1754
v1.2.11-dev-20160421-1754
exec npm dist-tag add @pilib/logger@1.2.11-dev-20160421-1754 dev
+dev: @pilib/logger@1.2.11-dev-20160421-1754
exec npm publish --tag dev
+ @pilib/logger@1.2.11-dev-20160421-1754

```

Show moreShow more icon

如果它不是默认分支，我们会稍微将其规范化，以确保分支名可在版本元数据中安全地使用。

- 删除 /
- 将 \_ 替换为 –
- 使用 normalize()

```
else {
    // if its not master branch, its a feature branch, create dist-tag for feature
    // based on branch name, publish to that.
    sh.echo(ch.blue('creating a feature version and publishing'));

    var normalizedBranch = branch.split('/').join('-');
    normalizedBranch = normalizedBranch.split('_').join('-');
    normalizedBranch = normalizedBranch.normalize();

    sh.echo(ch.green(normalizedBranch) + ' is the dist-tag for this feature');
    version += '-' + normalizedBranch + '-' + moment().tz('America/New_york').format('YYYYMMDD-HHmm');
    tag = normalizedBranch;
}

```

Show moreShow more icon

在此场景中，一个名为 `feature/geofences_39284` （其中 39284 是工作项 id）的分支被规范化为标签 `feature-geofences-39284` ，并使用完整版本 `@pilib/logger@1.0.5-feature-geofences-39284-20160416-0945` 。此版本可通过 `npm` install `@pilib/logger@feature-geofences-39284` 获得。

解析出相关信息并将标签和预发布版本字符串拼凑在一起后，发布到 `npm` 而不更新 `latest` 标签。

```
// properly exit if version fails
if (sh.exec('npm version ' + version).code !== 0) {
    sh.exit(1);
}

// properly exit if adding the dist-tag fails
if (sh.exec('npm dist-tag add ' + scope + '/' + name + '@' + version + ' ' + tag).code !== 0) {
    sh.exit(1);
}

// properly exit if publish fails
if (sh.exec('npm publish --tag ' + tag).code !== 0) {
    sh.exit(1);
}

```

Show moreShow more icon

首先，剪切掉我们创建的版本。如果这是第一次，则显式添加 `dist-tag` ，然后使用 `--tag` 标志发布，以告诉 `npm` 不要更新最新版本。

这样，签入的代码就可以通过创建的 dist-tag 用于本地和临时环境。指定 dev 作为私有 Presence Insights 依赖项版本的项目，将在整个迭代周期中自动选择对默认分支的更改。指定未来版本的项目将选择对该分支的更改。

## 下一步

1. 目前依靠 `moment-timezone` 等来精减 `@pidev/publish` ，但时区仅需出现一次。它还使用了 commander，这是一个大型框架。您可以精减 `@pidev/publish` 来仅依靠 `shelljs` 、 `yargs` 和 `chalk` 。
2. 在 `@pidev/publish` 中参数化默认分支和表示它的标签。我们使用主分支作为默认分支，但临时环境的名称为 dev，所以容易产生混淆。将名称从 `dev` 更改为 staging 或 `rc` 。

## 结束语

在本文中，您学习了如何结合使用 `npm dist-tags` 和 Travis CI，自动化预发布版本的管理。预发布版本是一种强大的工具，有助于在部署的环境中快速获取反馈，同时维护简洁且容易理解的版本历史，并且支持更轻松地对单个版本执行迭代式工作。