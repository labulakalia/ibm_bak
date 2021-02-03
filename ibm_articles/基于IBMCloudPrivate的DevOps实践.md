# 基于 IBM Cloud Private 的 DevOps 实践
充分利用容器云的优势，构建灵活的 Jenkins 集群

**标签:** DevOps,云计算

[原文链接](https://developer.ibm.com/zh/articles/d-based-ibm-cloud-private/)

吴伟

发布: 2018-05-08

* * *

## 概览

随着容器技术的兴起，一种新型的软件分发方式随即引爆了软件行业，而 Kubernetes 的崛起，更是为容器技术的大规模应用提供了坚实的技术支撑。当企业基于 IBM Cloud Private 构建容器云环境时，大部分应用都会陆续迁移到容器云环境内。在此背景下，DevOps 基础组件、应用团队开发的应用以及应用所依赖的中间件都会共存于容器云环境内，此时如何充分利用容器云的优势，构建灵活的 Jenkins 集群就成为一个迫切需要解决的问题。此外，从物理机、虚拟机环境迁移到容器云环境下，如何基于纯容器云环境实施 DevOps 实践也变成一种挑战。本文将从这两个方面给出参考解决思路。

## IBM Cloud Private 简介

IBM Cloud Private 是以 Docker、Kubernetes、Calico 等技术为核心，内置私有镜像中心、计量、监控、报警、日志等套件以及其他众多可选组件为一体的综合性容器云管理平台。IBM Cloud Private 当前提供 CE 和 EE 两个版本，其中 CE 版本提供了有限的功能集合，用户可以利用 IBM Cloud Private CE 去试用 IBM Cloud Private。有关 IBM Cloud Private 的详细内容，可参考官方文档： [IBM® Cloud Private v2.1.0.2 documentation](https://www.ibm.com/support/knowledgecenter/en/SSBS6K_2.1.0.2/kc_welcome_containers.html) 。

本文将演示的 DevOps 实践方法可以适用于原生 Kubernetes 环境、IBM Cloud Private CE 环境或 IBM Cloud Private EE 环境，文中演示的环境为 IBM Cloud Private V2.1.0.2。

## 在 IBM Cloud Private 中构建动态可伸缩的 Jenkins 集群

Jenkins 是当前主流的、用于 CI/CD 的软件，丰富的 Jenkins 插件生态可以为应用团队提供有力的功能支撑，是团队在实施 DevOps 的过程中的一个核心组件。

### Jenkins 集群形态简介

一般在大型开发团队中会构建 Jenkins 集群来作为 CI/CD 服务器，其中典型的集群形态如下图所示：

##### 图 1\. Jenkin 集群形态

![Jenkin 集群形态](../ibm_articles_img/d-based-ibm-cloud-private_images_image001.png)

Master 节点一般作为管理、配置中心，用户、用户权限、各类凭证、插件以及 Job 的管理中心，一个 Jenkins 集群中只有一个 Master 节点；而 Slave 节点通常作为具体 Job 的执行地点，一般根据项目组或应用划分，不同应用或项目组的应用分别在不同的 Slave 节点上执行；这种一个 Master 节点 N 个 Slave 节点的结构是 Jenkins 常见的集群形态。

在构建 Jenkins 集群时，先构建出 Master 节点，随后构建出 Slave 节点，并与 Master 节点连接，形成集群态。但这种集群结构有其优缺点：

- 优点：应用团队可以通过增减 Slave 节点来调整 Jenkins 集群的规模，来应对资源压力。
- 缺点：在物理机、虚拟机环境下，Slave 节点需要手动调整；其次，当 Slave 节点长时间没有任务执行时，会导致节点闲置、形成资源浪费。

当我们将 Jenkins 集群整体迁入容器云环境下时，整个集群形态没有发生变化。但这种静态集群形态会对容器云的计算资源形成浪费，我们需要一种更能适应容器云环境下的解决方法。

### 容器云环境下的 Jenkins 集群形态

在容器云环境下，整个 Jenkins 集群形态没有发生变化，变化的地方在于：Slave 节点不再是静态的，而是按需创建出来的。其形态特征如下图所示：

##### 图 2\. Jenkin 集群形态

![Jenkin 集群形态](../ibm_articles_img/d-based-ibm-cloud-private_images_image002.png)

在容器云环境下，Master 节点作为集群主节点以 Kubernetes 其中的一种资源形态存在（Deployment 或 Stateful Set），作为集群的配置中心。Slave 节点在有任务被触发时由 Jenkins 中的 [Kubernetes 插件](https://github.com/jenkinsci/kubernetes-plugin) 向 Kubernetes API Server 发送创建 Pod 请求动态生成，生成后与 Master 节点连接，构成集群态，最后在 Slave 节点内执行下发任务，任务完成后 Slave 节点随即被销毁。

这种动态可伸缩的集群形态有其优缺点：

- 优点：可最大化利用容器云的优势，实现 Slave 节点资源的最大化利用，同时又避免了因 Slave 节点空闲导致占用容器云计算资源的缺陷。
- 缺点：每个 Slave 节点是按需创建的，在 Slave 节点执行任务之前存在一个创建时间窗口（等待时间目前是在 1 分钟之内）；

在下文中将详细阐述如何基于 IBM Cloud Private 去构建这种动态 Jenkins 集群。

### 制作 Jenkins Master 镜像

在 [Docker Hub](https://hub.docker.com/) 中，Jenkins 具有官方镜像，其地址为： [OFFICIAL REPOSITORY – Jenkins](https://hub.docker.com/_/jenkins/) ，但是该镜像中只有裸的 Jenkins，并不具备一些必备的插件。而在实际应用场景中，IBM Cloud Private 通常位于企业内部，极端情况下甚至无法与外网互联，此时则需要预先定制好 Jenkins Master 镜像，其目的是使 Master 启动时即具有 CI/CD 中所必备的插件，实现开箱即用的效果。

制作 Jenkins Master 镜像的流程如下：

- 首先根据官方 Jenkins 镜像在本地生成 Jenkins 容器；
- 在本地 Jenkins 容器中安装 CI/CD 中所必的插件，并从中导出插件列表；
- 以 Jenkins 官方镜像为基础镜像，在其内通过脚本安装插件，最后封装成定制的 Jenkins Master 镜像；

其中根据 Jenkins 官方镜像在本地生成 Jenkins 容器，安装插件并导出插件列表的具体步骤参考如下：

##### 清单 1\. 获取预安装插件列表

```
# (1). 在本地下载官方 Jenkins 镜像
         docker pull jenkins:2.60.3
# (2). 在本地创建目录
         mkdir -p ~/tmp/volume/Jenkins && chown -R 1000:1000 ~/tmp/volume/jenkins

# (3). 在本地启动 Jenkins 容器
         docker run -d --name=jenkins -p 8080:8080 -v
                ~/tmp/volume/jenkins:/var/jenkins_home jenkins:2.60.3
# (4). 查看 Jenkins 容器日志获取初始随机密码
          docker logs -f jenkins

# (5). 当出现下列信息时复制随机密码并输入到 Jenkins 页面中完成 Jenkins 初始化
        *************************************************************
        Jenkins initial setup is required. An admin user has been created and a
            password generated.
        Please use the following password to proceed to installation:
        01a96b6c98084aa69be95a4b7703728a
        This may also be found at:
            /var/jenkins_home/secrets/initialAdminPassword
        *************************************************************
# (6). 安装推荐插件
        # Install suggested plugins
# 为 Jenkins 设置初始管理员账号与密码

# (7). 安装 Kubernetes 插件
# 在[插件管理]>[可选插件]中安装 Kubernetes 插件，本文演示的插件版本为 1.5.1

# (8). 安装其他 CI/CD 所需要的插件
# 在[插件管理]>[可选插件]中安装其他 CI/CD 所需要的插件

# (9). 导出插件列表到 plugins.txt 文件
        #curl -sSL
        "http://username:password@host:port/pluginManager/api/xml?depth=1&xpath=/*/*/shortName|/*/*/version&wrapper=plugins"
        | perl -pe
        's/.*?<shortName>([\w-]+).*?<version>([^<]+)()(<\/\w+>)+/\1
        \2\n/g'|sed 's/ /:/' > plugins.txt

#例如：

        curl -sSL
        "http://root:root@localhost:8080/pluginManager/api/xml?depth=1&xpath=/*/*/shortName|/*/*/version&wrapper=plugins"
        | perl -pe
        's/.*?<shortName>([\w-]+).*?<version>([^<]+)()(<\/\w+>)+/\1
        \2\n/g'|sed 's/ /:/' > plugins.txt

```

Show moreShow more icon

最终得到的插件列表 plugins.txt 如下所示：

##### 清单 2\. 获取预安装插件列表

```
cloudbees-folder:6.4
bouncycastle-api:2.16.2
structs:1.14
script-security:1.43
workflow-step-api:2.14
scm-api:2.2.6
workflow-api:2.27
junit:1.24
antisamy-markup-formatter:1.5
token-macro:2.5
build-timeout:1.19
credentials:2.1.16
ssh-credentials:1.13
plain-credentials:1.4
credentials-binding:1.16
timestamper:1.8.9
workflow-support:2.18
durable-task:1.22
workflow-durable-task-step:2.19
matrix-project:1.13
resource-disposer:0.8
ws-cleanup:0.34
ant:1.8
gradle:1.28
pipeline-milestone-step:1.3.1
jquery-detached:1.2.1
jackson2-api:2.8.11.1
ace-editor:1.1
workflow-scm-step:2.6
workflow-cps:2.45
pipeline-input-step:2.8
pipeline-stage-step:2.3
workflow-job:2.12.2
pipeline-graph-analysis:1.6
pipeline-rest-api:2.10
handlebars:1.1.1
momentjs:1.1.1
pipeline-stage-view:2.10
pipeline-build-step:2.7
pipeline-model-api:1.2.7
pipeline-model-extensions:1.2.7
apache-httpcomponents-client-4-api:4.5.3-2.1
jsch:0.1.54.2
git-client:2.7.1
git-server:1.7
workflow-cps-global-lib:2.9
display-url-api:2.2.0
mailer:1.21
branch-api:2.0.19
workflow-multibranch:2.16
authentication-tokens:1.3
docker-commons:1.11
docker-workflow:1.15.1
pipeline-stage-tags-metadata:1.2.7
pipeline-model-declarative-agent:1.1.1
workflow-basic-steps:2.6
pipeline-model-definition:1.2.7
workflow-aggregator:2.5
github-api:1.90
git:3.8.0
github:1.29.0
github-branch-source:2.3.3
pipeline-github-lib:1.0
mapdb-api:1.0.9.0
subversion:2.10.5
ssh-slaves:1.26
matrix-auth:2.2
pam-auth:1.3
ldap:1.20
email-ext:2.62
kubernetes-credentials:0.3.1
variant:1.1
kubernetes:1.5.1
javadoc:1.4
maven-plugin:3.1.2
analysis-core:1.95
checkstyle:3.50
performance:3.8
findbugs:4.72
jira:2.5
role-strategy:2.7.0

```

Show moreShow more icon

随后，我们可以开始定制 Jenkins Master 镜像，具体步骤参考如下：

##### 清单 3\. 定制 Jenkins Master 镜像

```
# (1). 在本地创建文件夹
        mkdir -p ~/tmp/masterbuildtxt && cd ~/tmp/masterbuildtxt

# (2). 在文件夹内放置 plugins.txt 文件
        root@kvm-015414:~/tmp/masterbuildtxt# mv ~/plugins.txt .
        root@kvm-015414:~/tmp/masterbuildtxt# ls -l
        total 4
        -rw-r--r-- 1 root root 1652 4 月  14 11:56 plugins.txt

# (3). 在文件夹内创建 Dockerfile
        cat >> Dockerfile <<EOF
        FROM jenkins:2.60.3
        LABEL maintainer="weiwunb@cn.ibm.com"
        # 预安装插件
        COPY plugins.txt /usr/share/jenkins/ref/plugins.txt
        RUN /usr/local/bin/install-plugins.sh <
            /usr/share/jenkins/ref/plugins.txt
        EOF

# (4). 文件夹内文件列表
        root@kvm-015414:~/tmp/masterbuildtxt# ls -l
        total 8
        -rw-r--r-- 1 root root 203 4 月 14 12:04 Dockerfile
        -rw-r--r-- 1 root root 1652 4 月  14 11:56 plugins.txt

# (5). 构建 Jenkins Master 镜像
#执行构建之后会在线下载插件，时间较长
        #此处以 IBM Cloud Private 默认的 CA Domain 为前缀，请以实际情况为准做调整
        docker build -t mycluster.icp:8500/default/jenkins:2.60.3 .
# (6). 查看本地镜像
        docker images | grep jenkins

```

Show moreShow more icon

至此，我们已经基于 [Jenkins 官方镜像](https://github.com/jenkinsci/docker/blob/587b2856cd225bb152c4abeeaaa24934c75aa460/Dockerfile) ，定制了内置 [Kubernetes 插件](https://github.com/jenkinsci/kubernetes-plugin) 以及其他 CI/CD 所需要插件的 Jenkins Master 镜像。

### 制作 Jenkins Slave 镜像

Jenkins Slave 作为承载实际任务执行的场地，其通常需要执行：应用程序的单元测试、应用程序打包、应用程序镜像制作及推送、以及最后在 IBM Cloud Private（或原生 Kubernetes 集群）中部署/更新应用等任务。

未来 [Kubernetes 插件 Kubernetes 插件](https://github.com/jenkinsci/kubernetes-plugin) 将利用 Jenkins Slave 镜像动态创建 Slave 节点，因此 Jenkins Slave 镜像需要满足的基础需求，其基础镜像为： [jenkinsci/docker-jnlp-slave](https://github.com/jenkinsci/docker-jnlp-slave) 。我们在根据业务场景和需求定制 Jenkins Slave 镜像时，需要以该镜像为基础镜像，增加所需的应用软件到其内。因此，在制作 Jenkins Slave 镜像之前，需要应用团队做好规划，按场景分别定制不同的 Jenkins Slave 镜像。为便于演示，本文此处以一个特定的场景为例，举例说明如何定制 Jenkins Slave 镜像。

此处以常见的 Java 应用开发测试场景为例，在 Jenkins Slave 内，需要完成应用的构建、自动化测试、应用镜像构建、应用在在 IBM Cloud Private（或原生 Kubernetes 集群）中部署/更新。总结起来，该 Slave 镜像内应该含有以下软件：

##### 表 1\. Slave 样例软件列表

**名称****解释**Maven适用于 Maven 应用的构建Gradle适用于 Gradle 应用的构建Docker用于构建应用镜像，Docker in Docker 与宿主机共享/var/run/docker.sockkubectl用于与 Kubernetes API Server 通信，在集群内部署/更新应用

随后，我们可以开始定制 Jenkins Slave 镜像，具体步骤参考如下：

##### 清单 4\. 定制 Jenkins Slave 镜像样例

```
# (1). 在本地创建文件夹
        mkdir -p ~/tmp/slavebuildtxt && cd ~/tmp/slavebuildtxt

# (2). 在文件夹内放置 Maven、Gradle、docker、kubectl
        root@kvm-015414:~/tmp/slavebuildtxt# ls -l
        total 79024
        -rwxr-xr-x 1 root root 13511128 4 月 14 21:06 docker
        -rw-r--r-- 1 root root 1566 4 月 14 21:16 Dockerfile
        drwxr-xr-x 5 root root 4096 4 月 14 13:54 gradle
        -rw-r--r-- 1 root root 67390556 4 月 14 14:03 kubectl
        drwxr-xr-x 6 root root 4096 4 月 14 13:54 maven

# (3). 在文件夹内创建 config
        cat >> config <<EOF
        apiVersion: v1
        kind: Config
        preferences: {}
        clusters:
        - cluster:
        certificate-authority: /etc/kubernetes/conf/ca.crt
        server: https://kubernetes.default
        name: mycluster.icp
        users:
        - name: admin
        user:
        client-certificate: /etc/kubernetes/conf/kubecfg.crt
        client-key: /etc/kubernetes/conf/kubecfg.key
        contexts:
        - context:
        cluster: mycluster.icp
        user: admin
        name: mycluster.icp.cxt
        current-context: mycluster.icp.cxt
        EOF

# (4). 文件夹内文件列表
        root@kvm-015414:~/tmp/slavebuildtxt# ls -l
        total 79024
        -rw-r--r-- 1 root root 441 4 月 14 14:06 config
        -rwxr-xr-x 1 root root 13511128 4 月 14 21:06 docker
        -rw-r--r-- 1 root root 1566 4 月 14 21:16 Dockerfile
        drwxr-xr-x 5 root root 4096 4 月 14 13:54 gradle
        -rw-r--r-- 1 root root 67390556 4 月 14 14:03 kubectl
        drwxr-xr-x 6 root root 4096 4 月 14 13:54 maven

# (5). 修改 Gradle 默认缓存位置
#编辑文件：~/tmp/slavebuildtxt/gradle/bin/gradle
#修改 GRADLE_OPTS 配置，将缓存位置设置为 HOME=/home/jenkins 目录下
        GRADLE_OPTS=-Dgradle.user.home=$HOME/.gradle/

# (6). 修改 Maven 默认缓存位置
#编辑文件：~/tmp/slavebuildtxt/maven/conf/settings.xml
#修改 localRepository 配置，将缓存位置设置为 HOME=/home/jenkins 目录下
        <localRepository>$HOME/.m2/repository</localRepository>

# (7). 在文件夹内创建 Dockerfile
        cat >> Dockerfile<<EOF
        FROM jenkins/jnlp-slave:latest
        LABEL maintainer="weiwunb@cn.ibm.com"
        USER root
        #------------------------------------------------#
        #安装 Maven3.5.2
        COPY maven /usr/share/maven/
        RUN chmod +x /usr/share/maven/bin/mvn && ln -s
        /usr/share/maven/bin/mvn /usr/local/bin/mvn
        #------------------------------------------------#
        #安装 Gradle4.4.1
        COPY gradle /usr/share/gradle/
        RUN chmod +x /usr/share/gradle/bin/gradle && ln -s
        /usr/share/gradle/bin/gradle /usr/local/bin/gradle
        #------------------------------------------------#
        #安装 kubectl 1.9.2
        COPY kubectl /usr/local/bin/
        RUN chmod +x /usr/local/bin/kubectl
        #------------------------------------------------#
        #预配置 kubectl
        #此处文件均为空文件将在运行时由 ConfigMap 挂载为 Volume 填充真实证书文件
        ENV CERT_DIR /etc/kubernetes/conf
        ARG DOCKER_SOCK_DIR=/var/run/docker.sock
        COPY config ${CERT_DIR}/
        RUN mkdir -p ${CERT_DIR} \
        mkidr -p /root/.kube \
        touch ${CERT_DIR}/ca.crt \
        touch ${CERT_DIR}/kubecfg.crt \
        touch ${CERT_DIR}/kubecfg.key
        RUN export KUBECONFIG=${CERT_DIR}/config && kubectl config view
        #------------------------------------------------#
        #仅提取 docker-ce_17.03.0_ce-0_debian-stretch_amd64 内的 docker 二进制文件
        copy docker /usr/bin/docker
        RUN apt-get -y update && apt-get install -y apt-utils iptables
        libdevmapper1.02.1 libltdl7 libseccomp2 \
        && apt-get -y autoremove && chmod +x /usr/bin/docker
        #------------------------------------------------#
        #暴露证书文件所在文件夹，在运行时由 ConfigMap 挂载为 Volume 填充真实证书文件
        VOLUME ${CERT_DIR}
        #Docker in Docker 与宿主机共享/var/run/docker.sock
        VOLUME ${DOCKER_SOCK_DIR}
        EOF
# (8). 查看文件夹内文件列表
        root@kvm-015414:~/tmp/slavebuildtxt# ls -l
        total 79024
        -rw-r--r-- 1 root root 441 4 月 14 14:06 config
        -rwxr-xr-x 1 root root 13511128 4 月 14 21:06 docker
        -rw-r--r-- 1 root root 1566 4 月 14 21:16 Dockerfile
        drwxr-xr-x 5 root root 4096 4 月 14 13:54 gradle
        -rw-r--r-- 1 root root 67390556 4 月 14 14:03 kubectl
        drwxr-xr-x 6 root root 4096 4 月 14 13:54 maven

# (9). 构建样例 Jenkins Slave 镜像
        docker build -t mycluster.icp:8500/default/jnlp-slave-java:latest .

# (9). 查看 Jenkins Slave 镜像
        docker images | grep jnlp-slave

```

Show moreShow more icon

至此，我们已经基于 [jenkinsci/docker-slave](https://github.com/jenkinsci/docker-slave/) 和 [jenkinsci/docker-jnlp-slave](https://github.com/jenkinsci/docker-jnlp-slave/) ，定制了能够为 Java 应用开发打包、构建镜像与部署应用的 Slave 镜像。

在完成了 Jenkins Master 与 Jenkins Slave 的镜像之后，我们将转到 IBM Cloud Private 环境内，在其内用这些镜像构建出动态 Jenkins 集群。

### 将镜像推送至 IBM Cloud Private 私有镜像中心

应用开发团队在完成 Jenkins Master、Slave 的镜像之后，需要将镜像推送到 IBM Cloud Private 私有镜像中心，以获得本地拉取镜像的速度优势。具体过程如下：

##### 清单 5\. 将 Jenkins Master、Slave 镜像推送到 IBM Cloud Private 私有镜像中心

```
# (1). 在本地创建目录
        mkdir -p /etc/docker/certs.d/mycluster.icp\:8500/
        cd /etc/docker/certs.d/mycluster.icp\:8500/

# (2). 从任意 IBM Cloud Private Master 节点中获取私有镜像证书
        # scp root@<master_node>:/etc/docker/certs.d/<cluster_CA_domain>\:8500/ca.crt.
#例如
        scp root@9.110.86.34:/etc/docker/certs.d/mycluster.icp\:8500/ca.crt .
# (3). 查看证书是否拷贝至本地
        root@kvm-015414:/etc/docker/certs.d/mycluster.icp:8500# pwd
        /etc/docker/certs.d/mycluster.icp:8500
        root@kvm-015414:/etc/docker/certs.d/mycluster.icp:8500# ls -l
        total 4
        -rw-r--r-- 1 root root 1108 4 月 14 15:06 ca.crt

# (4). 在修改本地 hosts 文件映射<master_ip> <cluster_CA_domain>
#注意:如果是 IBM Cloud Private 是 HA 结构 请用<cluster_vip> <cluster_CA_domain>
#此处以<cluster_vip> <cluster_CA_domain>做映射
        echo "9.110.86.214 mycluster.icp" >> /etc/hosts

# (5). 在本地重启 Docker 服务
        service docker restart

# (6). 尝试登陆 IBM Cloud Private 私有镜像中心
       # docker login <cluster_CA_domain>:8500

# 例如：
        root@kvm-015414:# docker login mycluster.icp:8500
        Username: admin
        Password:
        Login Succeeded

# (7). 将 Jenkins Master 镜像推送至 IBM Cloud Private 私有镜像中心
        docker push mycluster.icp:8500/default/jenkins:2.60.3

# (8). 将 Jenkins Slave 样例镜像推送至 IBM Cloud Private 私有镜像中心
        docker push mycluster.icp:8500/default/jnlp-slave-java:latest

```

Show moreShow more icon

在将镜像推送到 IBM Cloud Private 私有镜像中心后，可以在 IBM Cloud Private Console > Catalog > Images 中查看镜像：

##### 图 3\. IBM Cloud Private私有镜像中心

![IBM Cloud Private私有镜像中心](../ibm_articles_img/d-based-ibm-cloud-private_images_image003.png)

在将镜像推送到 IBM Cloud Private 私有镜像中心之后，我们可以着手在集群内部署 Jenkins Master 节点。

### 在 IBM Cloud Private 内部署 Jenkins Master

如果需要在 IBM Cloud Private 内部署应用，首先我们需要能和 Kubernetes API Server 通信，我们选择用 kubectl 直接与 Kubernetes API Server 通信，步骤如下：

##### 清单 6\. 配置 kubectl 与 IBM Cloud Private 进行通信

```
# (1). 在本地准备 kubectl 客户端
        apt-get update && apt-get install -y apt-transport-https curl
        curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add
        -
        cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
        deb http://apt.kubernetes.io/ kubernetes-xenial main
        EOF
        apt-get update
        apt-get install -y kubectl

# (2). 在本地创建文件夹~/tmp/cert/
        mkdir –p /etc/kubernetes/conf/
# (3). 在任意 IBM Cloud Private 的 master 节点内获取以下三个证书并拷贝至本地文件夹
        scp <user>@<master_ip>:/etc/cfc/conf/ ca.crt /etc/kubernetes/conf/
        scp <user>@<master_ip>:/etc/cfc/conf/kubecfg.crt
        /etc/kubernetes/conf/
        scp <user>@<master_ip>:/etc/cfc/conf/kubecfg.key
        /etc/kubernetes/conf/

# (4). 执行以下指令配置 kubectl
        kubectl config set-cluster cfc --server=https://<master_ip or
        cluster_vip>:8001 --certificate-authority=/kubernetes/conf/ca.crt
        kubectl config set-credentials user
        --client-certificate=/kubernetes/conf/kubecfg.crt
        --client-key=/kubernetes/conf/kubecfg.key
        bkubectl config set-context cfc --cluster=cfc --user=user
        --namespace=kube-system
        kubectl config use-context cfc

# (5). 执行以下指令配置确认与 IBM Cloud Private 内的 API Server 通信正常
        kubectl cluster-info

```

Show moreShow more icon

##### 清单 7\. 在 IBM Cloud Private 内部署 Jenkins Master

```
# (1). 在本地准备文件
        cat >> jenkins-master.yaml <<EOF
        apiVersion: v1
        kind: ServiceAccount
        metadata:
        name: jenkins #创建服务账号
        namespace: default #命名空间
        ---
        apiVersion: rbac.authorization.k8s.io/v1beta1
        kind: ClusterRoleBinding
        metadata:
        name: jenkins-crb
        roleRef:
        apiGroup: rbac.authorization.k8s.io
        kind: ClusterRole
        name: cluster-admin #角色：Cluster Master
        subjects:
        - kind: ServiceAccount
        name: jenkins #将角色与服务账号绑定
        namespace: default #命名空间
        ---
        apiVersion: v1
        kind: Service
        metadata:
        name: jenkins #服务名
        namespace: default #所在命名空间
        spec:
        clusterIP: None
        ports:
        - name: http #Jenkins Master Web 服务端口
        port: 8080
        targetPort: 8080
        - name: slavelistener #Jenkins Master slave 监听端口
        port: 50000
        targetPort: 50000
        selector:
        devops: jenkins-master #Pod 标签
        ---
        apiVersion: apps/v1
        kind: StatefulSet
        metadata:
        name: jenkins-master
        namespace: default
        spec:
        replicas: 1
        serviceName: jenkins
        selector:
        matchLabels:
        devops: jenkins-master #Pod 标签
        template:
        metadata:
        labels:
        devops: jenkins-master #Pod 标签
        spec:
        serviceAccount: jenkins #Pod 需要使用的服务账号
        containers:
        - name: jenkins
        image: mycluster.icp:8500/default/jenkins:2.60.3 #Jenkins Docker 镜像名称
        imagePullPolicy: IfNotPresent
        ports:
        - name: http #Jenkins Master Web 服务端口
        containerPort: 8080
        - name: slavelistener #Jenkins Master 供未来 Slave 连接的端口
        containerPort: 50000
        volumeMounts:
        - name: storage #存储卷名称
        mountPath: /var/jenkins_home #挂载目录
        env:
        - name: JAVA_OPTS
        value: "-Xms512m -Xmx896m -Duser.timezone=Asia/Shanghai" #配置 Jenkins JVM
        最小、最大堆内存 所在时区
        resources:
        #https://jenkins.io/doc/book/hardware-recommendations/
        limits:
        cpu: "2000m"
        memory: 1Gi
        requests:
        cpu: "2000m"
        memory: 1Gi
        volumeClaimTemplates:
        - metadata:
        name: storage
        namespace: default
        spec:
        storageClassName: glusterfs
        accessModes:
        - ReadWriteMany
        resources:
        requests:
        storage: 2Gi
        ---
        apiVersion: extensions/v1beta1
        kind: Ingress
        metadata:
        name: jenkins-ing #Ingress 名
        namespace: default
        spec:
        rules:
        - host: default.jenkins.mycluster.icp #Jenkins Master 的路由地址
        http:
        paths:
        - path:
        backend:
        serviceName: jenkins #服务名
        servicePort: 8080
        EOF

# (2). 执行部署
       kubectl apply -f jenkins-master.yaml

# (4). 查看 Jenkins Master 的 Pod 状态
       kubectl get pod -o wide -n=default
# (5). 查看 Jenkins Master 的 Pod 日志
       kubectl logs jenkins-master-0 -n=default
# (6). 修改本地 hosts，增加 Jenkins Master 的 Ingress 路由与 IP 映射
       <proxy_vip> <Jenkins Master Ingress>

#例如
       9.110.86.211 default.jenkins.mycluster.icp

```

Show moreShow more icon

当观察 Jenkins Master 的 Pod 已经正常启动之后，在本地浏览器中输入 Jenkins Master 的 Ingress 地址，进入初始化阶段：

##### 图 4\. Jenkins 初始化界面

![Jenkins 初始化界面](../ibm_articles_img/d-based-ibm-cloud-private_images_image004.png)

在 Administrator password 中输入从 jenkins-master-0 日志中获取的随机初始密码，单击 Continue。当 IBM Cloud Private 所出的环境不可连接外网时，会出现以下界面：

##### 图 5\. Jenkins 初始化界面

![Jenkins 初始化界面](../ibm_articles_img/d-based-ibm-cloud-private_images_image005.png)

直接单击 Skip Plugin Installations，跳过插件安装环节（因为我们已经在 Jenkins Master 镜像中预安装了所有必须的插件），进入设置初始管理员账号信息界面：

##### 图 6\. Jenkins 初始化管理员界面

![Jenkins 初始化管理员界面](../ibm_articles_img/d-based-ibm-cloud-private_images_image006.png)

设置初始管理员账号信息，进入 Jenkins。

### 创建 kubectl-cert-cm

我们在制作 Jenkins Slave 镜像时，在 Slave 镜像内封装了 kubectl 命令，并提供了一个 config 文件作为 kubectl config 的参数文件，但 config 内所使用的三个证书文件在制作镜像时均为空文件，这样做的目的在于：我们希望可以在不同的 IBM Cloud Private 环境内复用 Slave 镜像。

如果要想在 Slave 镜像内利用 kubectl 与 IBM Cloud Private 的 API Server 进行通信，还需要在运行时为其指定真实的证书文件，本文中就选取通过将证书文件封装成 ConfigMap 资源对象，在运行时挂载到 Slave Pod 的${CERT\_DIR}目录下来实现将认证证书和 Slave 镜像相分离的效果。

##### 清单 7\. 创建 kubectl-cert-cm 资源对象

```
# (1). 创建 kubectl-cert-cm 资源对象
       kubectl create configmap kubectl-cert-cm --from-file=/etc/kubernetes/conf
            -n=default

# (2). 查看所创建的 kubectl-cert-cm 资源对象
       kubectl describe configmap kubectl-cert-cm -n=default

```

Show moreShow more icon

### 创建 jenkins-slave -pvc

我们在创建 Jenkins Master 时为其指定了 PVC，该 PVC 将自动从 IBM Cloud Private 集群内的 GlusterFS 集群内分配 PV 用于 Master 节点的数据持久化。同样的，我们希望未来 Jenkins Slave 在执行时，Slave 内执行任务所生成的数据也能被持久化，因此我们需要先在 IBM Cloud Private 集群内创建 PVC，未来该 PVC 会挂载到 Slave 节点的/home/jenkins 目录下。

##### 清单 8\. 创建 jenkins-slave-pvc 资源对象

```
# (1). 在本地创建 jenkins-slave-pvc.yaml 文件
        cat >> jenkins-slave-pvc.yaml << EOF
        kind: PersistentVolumeClaim
        apiVersion: v1
        metadata:
        name: jenkins-slave-pvc #PVC 名称 挂载到 Jenkins Slave 工作目录
        namespace: default #所在命名空间
        spec:
        storageClassName: glusterfs
        accessModes:
        - ReadWriteMany
        resources:
        requests:
        storage: 2Gi #样例容量，请按实际需求设定容量
        EOF
# (2). 创建 jenkins-slave-pvc 资源对象
        kubectl apply -f jenkins-slave-pvc.yaml
# (3). 查看 pvc 状态
#如果表示 PVC 的 Status 显示 Bound，则 PVC 已经和 PV 绑定成功，可以挂载到 Slave 目录下
        kubectl describe –f jenkins-slave-pvc.yaml

```

Show moreShow more icon

### 在 Jenkins Master 完成 Kubernetes 插件配置

在正式进入 Jenkins Master 中，利用 [Kubernetes 插件](https://github.com/jenkinsci/kubernetes-plugin) 进行配置时，我们先列举将要配置的信息：

##### 表 2\. 配置信息清单

**名称****解释**Kubernetes URLIBM Cloud Private 内的 API Server Service 地址: [https://kubernetes.default](https://kubernetes.default)Jenkins URLJenkins Master 的 Service 地址及端口号: [http://jenkins.default:8080/](http://jenkins.default:8080/)Container Cap最大并发 Slave Pod 数量：30Pod.Templage.Namejnlp-slave-javaPod.Templage.NamespacedefaultPod.Templage.Labelsjnlp-slave-javaPod.Templage.Containers.NamejnlpPod.Templage.Containers.Docker Imagemycluster.icp:8500/default/jnlp-slave-java:latestPod.Templage.Containers.Always pull imageyesPod.Templage.Containers.EnvVarsKUBECONFIG=/etc/kubernetes/conf/configPod.Templage.VolumesPersistent Volume Claim：jenkins-slave-pvc，挂载目录：/home/jenkinsPod.Templage.VolumesConfig Map Volume：kubectl-cert-cm，挂载目录：/etc/kubernetes/confPod.Templage.VolumesHost Path Volume：Host Path：/var/run/docker.sock Mount Path：/var/run/docker.sock

在列出清单之后，我们可以开始进入 Jenkins Master 利用 [Kubernetes 插件](https://github.com/jenkinsci/kubernetes-plugin) 进行配置。

- 登录 Jenkins Master，进入[系统管理] > [系统设置] > [云]：

##### 图 7\. Kubernetes-plugin 新增云

![Kubernetes-plugin 新增云](../ibm_articles_img/d-based-ibm-cloud-private_images_image007.png)

- 单击[新增一个云]，在 Kubernetes URL 中输入 IBM Cloud Private 内 API Server 的服务名，单击 Test Connection，测试 Jenkins Master 是否可以与 API Server 进行通信。

注意：确保测试结果为 Connection test successful，否则 Jenkins Master 将无法向 API Server 发送创建 Slave Pod 的请求。

##### 图 8\. Kubernetes-plugin 配置 Kubernetes URL

![Kubernetes-plugin 配置 Kubernetes URL](../ibm_articles_img/d-based-ibm-cloud-private_images_image008.png)

- 在 Jenkins URL 中输入 Jenkins Master 的 Service 地址及端口号，在 Container Cap 中设置最大并发 Slave Pod 的数量：

注意：Container Cap 数量决定了未来会有多少 Slave 节点的数量，请按实际任务并发数量需求进行设定。此处演示的为 30，意为 Jenkins 集群内最大可同时存在 30 个 Slave 节点。

##### 图 9\. Kubernetes-plugin 配置 Jenkins URL

![Kubernetes-plugin 配置 Jenkins URL](../ibm_articles_img/d-based-ibm-cloud-private_images_image009.png)

- 单击 Add Pod Templage，添加 Slave 节点的 Pod 模板信息：

说明：在 Name 和 Labels 中输入 Jenkins Slave 的标签；在 Namespace 中指定未来 Slave 节点会在哪个 Namespace 内被创建；Usage 中指定：只允许运行绑定到这台机器的 Job；

##### 图 10\. Kubernetes-plugin 配置 Pod Template

![Kubernetes-plugin 配置 Pod Template](../ibm_articles_img/d-based-ibm-cloud-private_images_image010.png)

- 单击 Add Containers，添加 Pod 内的容器模板：

说明：Name 必须为 jnlp；Docker Image 为我们定制的 Jenkins Slave 镜像；Working directory 为：/home/jenkins；

##### 图 11\. Kubernetes-plugin 配置 Container Template

![Kubernetes-plugin配置Container Template](../ibm_articles_img/d-based-ibm-cloud-private_images_image011.png)

- 单击[Add Environment Variable] > [Environment Variable]：

说明：在 key 中输入 KUBECONFIG；在 value 中输入：/etc/kubernetes/conf/config

##### 图 12\. Kubernetes-plugin 配置环境变量

![Kubernetes-plugin配置环境变量](../ibm_articles_img/d-based-ibm-cloud-private_images_image012.png)

- 单击[Add Volume] > [Persistent Volume Claim]：

说明：在 Claim Name 中输入 jenkins-slave-pvc，在 Mount Path 中输入：/home/Jenkins；

##### 图 13\. Kubernetes-plugin 配置 Persistent Volume Claim

![alt](../ibm_articles_img/d-based-ibm-cloud-private_images_image013.png)

- 单击[Add Volume] > [Config Map Volume]：

说明：在 Config Map name 中输入 kubectl-cert-cm，在 Mount Path 中输入：/etc/kubernetes/conf；

##### 图 13\. Kubernetes-plugin 配置 Config Map Volume

![alt](../ibm_articles_img/d-based-ibm-cloud-private_images_image014.png)

- 单击[Add Volume] > [Host Path Volume]：

说明：在 Host Path 中输入：/var/run/docker.sock，在 Mount Path 中输入：/var/run/docker.sock

##### 图 14\. Kubernetes-plugin 配置 Host Path Volume

![alt](../ibm_articles_img/d-based-ibm-cloud-private_images_image015.png)

- 单击[保存]，完成配置。

### 测试动态 Jenkins 集群

在利用 [Kubernetes 插件](https://github.com/jenkinsci/kubernetes-plugin) 完成配置后，我们可以新建一个 Job 来测试动态 Slave 的生成情况。测试目标：

- 当任务被触发时， [Kubernetes 插件](https://github.com/jenkinsci/kubernetes-plugin) 应该能动态生成 Jenkins Slave 节点并执行任务；
- Jenkins Slave 内的 Gradle、Maven、Docker 能正常工作并输出版本号；
- Kubectl 能输出 IBM Cloud Private 集群信息；
- 连续测试多个任务，观察是否有多个 Slave 节点被生成，形成多 Slave 节点状态；

具体测试过程如下：

- 创建一个任务，类型为：自由软件风格，在任务内设置 Restrict where this project can be run：

说明：在 Label Expression 中设置 Kubernetes 插件中配置的 Slave Pod Templage Labels。

##### 图 15\. 设置 Label Expression

![alt](../ibm_articles_img/d-based-ibm-cloud-private_images_image016.png)

- 在构建中设置测试脚本：

##### 清单 9\. 创建 jenkins-slave-pvc 资源对象

```

ls -l /etc/kubernetes/conf && ls -l /home/jenkins
gradle -v && mvn -v && docker info
kubectl cluster-info
sleep 600

```

Show moreShow more icon

##### 图 16\. 设置构建 Shell

![alt](../ibm_articles_img/d-based-ibm-cloud-private_images_image017.png)

- 单击[保存]，执行[立即构建]，开始测试任务：

测试结果：Slave 节点由任务触发而被生成，并与 Master 节点连接，能执行下发任务；

##### 图 17\. 测试 Job

![alt](../ibm_articles_img/d-based-ibm-cloud-private_images_image018.png)

查看日志，从日志中可以看出 Gradle、Maven、Docker 均能正常工作，kubectl 能与 API Server 通信：

##### 图 18\. 测试 Job 日志

![alt](../ibm_articles_img/d-based-ibm-cloud-private_images_image019.png)

- 新建多个 Job，依次单击[立即构建]，测试并发任务执行场景：

注意：在本文所测试的 [Kubernetes 插件](https://github.com/jenkinsci/kubernetes-plugin) 版本为 1.5.1，Jenkins 版本为 2.60.3，在此背景下，如果是同一个 Job 被多次触发，则后续被触发的任务则会一直排队，直到前一个执行完毕；如果是不同的任务被同时触发，则会构建不同的 Slave 节点分别去执行这些任务。

##### 图 19\. 测试并行执行多个 Job

![alt](../ibm_articles_img/d-based-ibm-cloud-private_images_image020.png)

至此，我们已经利用 [Kubernetes 插件](https://github.com/jenkinsci/kubernetes-plugin) 在 IBM Cloud Private 内成功构建出了一种以静态 Jenkins Master 为核心，按需创建的动态 Jenkins Slave 节点为特征的 Jenkins 集群。

## 一种在 IBM Cloud Private 中的 DevOps 实施方法

当我们在 IBM Cloud Private 内构建出动态 Jenkins 集群后，我们即可开始着手实施 DevOps 实践。我们以常见的 Java 应用开发测试场景为例举例说明。下图中展示了我们在团队内部实施的 DevOps 实践方法：

##### 图 20\. DevOps 实施流程示意

![alt](../ibm_articles_img/d-based-ibm-cloud-private_images_image021.png)

鉴于团队对 Kubernetes 技术掌握尚不成熟，且团队希望将用于在 Kubernetes 内构件资源的 yaml 文件随同源码一同管理，因此我们团队在实际应用中没有使用 Helm。在实际操作过程中，我们通过动态替换镜像 TAG 的方式实现了模板与数据的分离。简要流程介绍如下：

- 应用开发团队需要在各自项目中自编写 Dockerfile、Jenkinsfile、Kubernetes.yaml 以及.dockeringnore 等文件，这些文件的作用及说明如下：

##### 表 3\. 文件清单说明

**名称****解释**Dockerfile用于构建应用镜像.dockeringnore用于排除构建上下文中的非必要文件Kubernetes.yaml用于构建 Kubernetes 资源Jenkinsfile用于设定自定义 CI/CD 流程

- 其中 Dockerfile 样例如下所示：

##### 清单 10\. Dockerfile样例

```

#
    https://github.com/docker-library/tomcat/blob/master/8.5/jre8/Dockerfile
# 位于 IBM Cloud Private Docker Registry 内的 tomcat 8.5 镜像
FROM mycluster.icp:8500/devops/tomcat:8.5
# 维护者信息
LABEL maintainer="weiwunb@cn.ibm.com"
# 定义 tomcat webapps 所在目录
ARG WEBAPPS_DIR=/usr/local/tomcat/webapps
# 将项目 build 的 war 包添加至 WEBAPPS_DIR 目录下
COPY target/spring-mvc-showcase.war ${WEBAPPS_DIR}/

```

Show moreShow more icon

- 其中 Kubernetes.yaml 文件样例如下所示：

说明：我们在 Kubernetes.yaml 内部为镜像的 TAG 部分设定了占位符，在实际构建过程中会将 TAG 占位符替换为实际的 TAG 值（一般为镜像版本号）；

##### 清单 11\. Kubernetes.yaml 样例

```

#example 项目的 Deployment API 描述#
apiVersion: apps/v1beta2
kind: Deployment
metadata:
name: springmvc-showcase
namespace: default #所属的 namespace
labels:
app: springmvc-showcase #Pod 的标签，将在 Service 中使用
version: IMAGE_TAG #版本模板占位符，将在实际更新时被替换为实际的版本号
spec:
replicas: 1 #初始化副本数量
selector:
matchLabels:
app: springmvc-showcase
template:
metadata:
labels:
app: springmvc-showcase #Pod 的标签，将在 Service 中使用
spec:
containers:
- name: springmvc-showcase #容器的名字
image: mycluster.icp:8500/default/spring-mvc-showcase:IMAGE_TAG
#容器所使用的镜像及其版本;版本模板占位符，将在实际更新时被替换为实际的版本号
ports:
- name: web
containerPort: 8080 #容器内暴露的端口
env: #容器内提供的可修改环境变量
- name: JAVA_OPTS
value: '-Xmx2048m'
---
#example 项目的 Service API 描述#
apiVersion: v1
kind: Service
metadata:
name: springmvc-showcase-svc
namespace: default #所属的 namespace
spec:
ports:
- name: web
port: 8080
targetPort: 8080 #example 项目的容器所暴露的端口
selector:
app: springmvc-showcase #example 项目的 Pod 的标签
---
#example 项目的 Ingress API 描述#
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
name: springmvc-showcase-ingress
namespace: default #所属的 namespace
spec:
rules:
- host: example.trtjk.dev.in #路由的 host
http:
paths:
- path: / #路由的路径
backend:
serviceName: springmvc-showcase-svc #所指向的 Service Name
servicePort: 8080 #Service 所使用的端口

```

Show moreShow more icon

- 其中 Jenkinsfile 内定义自定义的 CI/CD 流程，其样例如下所示（此处举例的是在一个 Pipeline 流程中完成项目构建、镜像打包、推送镜像以及部署应用）：

##### 清单 12\. Jenkinsfile样例

```

pipeline {
agent {
//指定此 Pipeline 将在动态 Jenkins Slave 中执行
label 'jnlp-slave-java'
}
options {
//仅保留最近 10 次的构建记录
buildDiscarder(logRotator(numToKeepStr: '10'))
//输出时间戳
timestamps()
}
environment {
//定义 build 的源码分支
BRANCH_NAME = 'master'
//定义应用名称
APP_NAME = 'spring-mvc-showcase'
//定义 IBM Cloud Private 内的 Docker Registry 认证信息
DOCKERREG_ACCESS_KEY = credentials('icp-docker-registry-credentials')
//定义 IBM Cloud Private 内的 Docker Registry 地址
DOCKERREG_DOMAIN = 'mycluster.icp:8500'
//定义镜像所归属的 Namespace
TARGET_NAMESPACE = 'default'
//定义镜像的 Tag(即版本号由 VersionNumber 插件动态生成)
//格式：<用于 build 的源码分支>-<build 的所属年月日>-<当天的第 N 次
build>
IMAGE_TAG = VersionNumber(projectStartDate: '1970-01-01',
versionNumberString: '${BUILD_DATE_FORMATTED, "yyyyMMdd"}-${BUILDS_TODAY,
XXX}', versionPrefix: "${BRANCH_NAME}-")
//应用镜像的完整构成
APP_DOCKER_IMAGE
="${DOCKERREG_DOMAIN}/${TARGET_NAMESPACE}/${APP_NAME}:${IMAGE_TAG}"
}
stages {
stage('Build project package') {
steps {
//利用 maven 对项目进行 build
sh 'mvn clean package -Dmaven.test.skip'
//检查应用的包
sh 'ls -l target | grep war'
}
}
stage('Build Docker Image') {
steps {
//先登录 IBM Cloud Private 内的 Docker Registry
sh "docker login -u ${DOCKERREG_ACCESS_KEY_USR} -p
${DOCKERREG_ACCESS_KEY_PSW} https://${DOCKERREG_DOMAIN}"
//以项目根目录为构建上下文 根据 Dockerfile 构建应用镜像
sh "docker build -t ${APP_DOCKER_IMAGE} ."
}
}
stage('Push Docker Image') {
steps {
//先登录 IBM Cloud Private 内的 Docker Registry
sh "docker login -u ${DOCKERREG_ACCESS_KEY_USR} -p
${DOCKERREG_ACCESS_KEY_PSW} https://${DOCKERREG_DOMAIN}"
//将所构建的镜像推送到 IBM Cloud Private 内的 Docker Registry
sh "docker push ${APP_DOCKER_IMAGE}"
}
}
tage('Update Application') {
steps {
//先通过 sed 命令将 kubernetes.yaml 中的 IMAGE_TAG 替换为实际的版本号
sh "sed -i 's/IMAGE_TAG/${IMAGE_TAG}/g' kubernetes.yaml"
//然后调用 kubectl apply 命令根据更新后的 kubernetes.yaml 去更新应用
sh 'kubectl apply -f kubernetes.yaml --record'
}
}
}
}

```

Show moreShow more icon

每个项目都会在 Jenkins Master 中对应一个 Pipeline 类型的 Job，构建结果如下图所示：

##### 图 21\. Pipeline 执行流程样例

![Pipeline执行流程样例](../ibm_articles_img/d-based-ibm-cloud-private_images_image022.png)

除了在 Jenkinsfile 中设定以上流程之外，应用开发团队可以按需定制代码规范检查、代码漏洞扫描以及自动化测试等任务，在此不再赘述。

## 结束语

Jenkins 作为主流的 CI/CD 软件，伴随开发人员从物理机、虚拟机时代走进了容器化时代，随之而来的调整即是 Jenkins 集群也以容器化形式呈现出来。本文基于 IBM Cloud Private 环境，结合 [Kubernetes 插件](https://github.com/jenkinsci/kubernetes-plugin) ，构建出一种以单一静态 Master 节点，按需创建的动态 Slave 为辅助的容器化 Jenkins 集群，该集群能较好的适应容器云环境，实现资源利用的最大化。此外，本文在最后简要提出一种在团队内部实施的 DevOps 实践，团队在直接使用 yaml 文件表述 Kubernetes 资源形态，还是通过 Helm 来管理 Kubernetes 资源形态的决策中，我们选择了前者，并通过在构建过程中动态替换镜像 TAG 的方式实现了模板与实际数据的分离。

## 参考资料

- 参考 [https://jenkins.io/](https://jenkins.io/)，了解 Jenkins 相关信息。
- 参考 [jenkinsci/docker/Dockerfile](https://github.com/jenkinsci/docker/blob/587b2856cd225bb152c4abeeaaa24934c75aa460/Dockerfile)，查看 Jenkins 官方镜像源码。
- 参考 [jenkinsci/slave](https://github.com/jenkinsci/docker-slave/)，查看 jenkinsci/docker-slave 官方镜像源码。
- 参考 [jenkinsci/docker-jnlp-slave](https://github.com/jenkinsci/docker-jnlp-slave)，查看 Jenkins jnlp-slave 官方镜像源码。