# Cloud-Provider-Openstack 介绍——CSI 块存储部分
CSI (Container Storage Interface) 在 Openstack 上 cinder 的实现介绍

**标签:** Kubernetes,云计算

[原文链接](https://developer.ibm.com/zh/articles/cl-lo-cloud-provider-openstack-csi/)

纪晨

发布: 2020-09-17

* * *

## 概览

随着 Kubernetes 逐步成为容器 orchestrator 领域的事实标准，标准化成为了非常重要的目标之一，现在有 [CSI](https://kubernetes-csi.github.io/docs/)(Container Storage Interface)、 [CRI](https://github.com/kubernetes/kubernetes/blob/242a97307b34076d5d8f5bbeb154fa4d97c9ef1d/docs/devel/container-runtime-interface.md)(Container Runtime Interface)、 [CNI](https://github.com/containernetworking/cni)(Container Network Interface) 等，在不同的云提供商(Cloud provider)中有各自不同的实现。

利用 OpenStack 的 [cinder](https://wiki.openstack.org/wiki/Cinder) 和 [manila](https://wiki.openstack.org/wiki/Manila) 分别作为块存储和文件存储的后端，提供 Kubernetes API 可以使用的存储后端。由于篇幅有限，本文主要讨论了 OpenStack Cinder 作为后端的实现和使用简介。

## 概念介绍

[Kubernetes Cloud Provider](https://kubernetes.io/zh/docs/concepts/architecture/cloud-controller/) 为了更好地让 Kubernetes 在公有云平台上运行以及提供容器云服务，云厂商（例如 IBM Cloud、AWS、OpenStack 等）需要通过实现 [该接口](https://github.com/kubernetes/kubernetes/blob/master/pkg/cloudprovider/providers/providers.go) 来提供自己的 Cloud Provider，它是 Kubernetes 中开放给各个云厂商的通用接口，便于 Kubernetes 创建/管理/利用各个云服务商提供的资源，这些资源包括但不限于虚拟机/裸金属机器资源、网络资源(VPC)，负载均衡服务(Load balancer)、弹性公网 IP (Floating IP)、存储服务（块，对象，文件存储）等。

随着 Kubernetes 不断地成熟，OpenStack 社区也在和 Kubernetes 不断地加强合作，通过 SIG (special interest group)来进行合作，主要项目是 [Cloud Provider OpenStack](https://github.com/kubernetes/cloud-provider-openstack)，该项目是 Cloud Provider 在 OpenStack 上的实现，这样通过这一层抽象，我们可以在 PaaS (Kubernetes) 和 IaaS (OpenStack) 之间加入一层粘合剂，从而所有的针对云提供商的 API 调用可以通过特定的 OpenStack 一个或者多个子项目来完成，Cloud Provider OpenStack 包含了若干个子项目:

- [OpenStack Cloud Controller Manager](https://github.com/kubernetes/cloud-provider-openstack/blob/master/docs/using-openstack-cloud-controller-manager.md)：OpenStack 的云控制器
- [Octavia Ingress Controller](https://github.com/kubernetes/cloud-provider-openstack/blob/master/docs/using-octavia-ingress-controller.md)：通过 OpenStack 的 Octavia 子项目实现的 Ingress Controller
- [Keystone Webhook Authentication Authorization](https://github.com/kubernetes/cloud-provider-openstack/blob/master/docs/using-octavia-ingress-controller.md)：使用 OpenStack keystone 项目作为认证后端
- [Manila CSI Plugin](https://github.com/kubernetes/cloud-provider-openstack/blob/master/docs/using-manila-csi-plugin.md)：使用 OpenStack 的 manila（文件存储）项目实现的 CSI plugin
- [Barbican KMS Plugin](https://github.com/kubernetes/cloud-provider-openstack/blob/master/docs/using-barbican-kms-plugin.md)：使用 OpenStack 的 barbican 项目实现的 kms provide
- [Cinder CSI Plugin](https://github.com/kubernetes/cloud-provider-openstack/blob/master/docs/using-cinder-csi-plugin.md)：使用 OpenStack 的 Cinder（块存储）项目实现的 CSI plugin，本文会重点介绍该项目

Kubernetes 中有卷和持久卷的概念，代表临时的或者持久的磁盘文件，分别用 Volume 和 PersistentVolume 这两个对象来表示。临时的卷没有办法解决内容的持久化问题，所以 Kubernetes 引入了 PersistentVolume（PV）这个对象，而如何提供这些卷，就是存储提供者（例如云提供商）要考虑的问题了。

CSI 是 Container Storage Interface 的缩写， 在容器编排引擎（例如 Kubernetes）和存储之间定义了一套标准的接口额（概念上与 CRI 和 CNI 类似），这样把 Kubernetes 和存储提供者彻底解耦。CSI 只是一套接口，因此不同的存储提供商都提供了自己的实现，在 Openstack 上有三种存储服务分别是 swift (对象存储)，cinder（块存储），manila（文件存储）。通过这样的方式，无论是虚拟机和容器的 workload，都可以通过 openstack 来管理块设备，从而达到统一管理的目的。

例如我们可以假设我们有 2 个应用负载，一部分需要在虚拟机中运行，另一部分在容器中，他们都有对于 persistent 块存储有一定的需求，我们可以通过 IaaS （OpenStack）来管理所有的该类型存储，而如何使用则由 Cloud Provider OpenStack（针对容器）或者 cinder(针对虚拟机)来决定。

## 创建 Kubernetes 运行环境

由于使用 CSI 需要 Kubernetes 运行环境，为了简单起见，笔者直接使用了源码的开发模式，读者可以直接使用现有的 Kubernetes 环境或者通过 Kubeadm 等工具创建环境；本文假设采用了源码的开发模式。

首先，由于需要编译，所以需要确保 go 编译环境创建完成，可以参考 [创建 golang 环境](https://golang.google.cn/)；接着需要通过 git 下载 [kubernetes](https://github.com/kubernetes/kubernetes)，并且运行如下命令从而创建一个 Kubernetes 运行环境：

```
ALLOW_PRIVILEGED=true RUNTIME_CONFIG="storage.k8s.io/v1=true" LOG_LEVEL=5 hack/l
ocal-up-cluster.sh
WARNING : The kubelet is configured to not fail even if swap is enabled; production deployments should disable swap.
make: Entering directory '/root/go/src/github.com/kubernetes'
make[1]: Entering directory '/root/go/src/github.com/kubernetes'
…..
To start using your cluster, you can open up another terminal/tab and run:

export KUBECONFIG=/var/run/kubernetes/admin.kubeconfig
cluster/kubectl.sh

Alternatively, you can write to the default kubeconfig:

export KUBERNETES_PROVIDER=local

cluster/kubectl.sh config set-cluster local --server=https://localhost:6443 --certificate-authority=/var/run/kubernetes/server-ca.crt
cluster/kubectl.sh config set-credentials myself --client-key=/var/run/kubernetes/client-admin.key --client-certificate=/var/run/kubernetes/client-admin.crt
cluster/kubectl.sh config set-context local --cluster=local --user=myself
cluster/kubectl.sh config use-context local
cluster/kubectl.sh

```

Show moreShow more icon

当我们看到如上的界面的时候，我们的本地 cluster 已经可以使用了。

为了简便起见，笔者创建了一个 alias，这样可以直接使用 kubectl 来完成后续命令：

```
alias 'kubectl=/root/go/src/github.com/kubernetes/cluster/kubectl.sh'

```

Show moreShow more icon

通过如下命令确定 cluster 已经准备完毕， 如果该命令出现错误，请检查上述命令的 log：

```
# kubectl get pods
No resources found in default namespace.

```

Show moreShow more icon

## 创建 CSI 运行环境

接着，我们需要创建 CSI 运行环境，所需的定义都已经在 [manifests/cinder-csi-plugin](https://github.com/kubernetes/cloud-provider-openstack/tree/master/manifests) 这个目录。通过如下命令创建：

```
# kubectl apply -f manifests/cinder-csi-plugin/
serviceaccount/csi-cinder-controller-sa created
clusterrole.rbac.authorization.k8s.io/csi-attacher-role created
clusterrolebinding.rbac.authorization.k8s.io/csi-attacher-binding created
clusterrole.rbac.authorization.k8s.io/csi-provisioner-role created
clusterrolebinding.rbac.authorization.k8s.io/csi-provisioner-binding created
clusterrole.rbac.authorization.k8s.io/csi-snapshotter-role created
clusterrolebinding.rbac.authorization.k8s.io/csi-snapshotter-binding created
clusterrole.rbac.authorization.k8s.io/csi-resizer-role created
clusterrolebinding.rbac.authorization.k8s.io/csi-resizer-binding created
role.rbac.authorization.k8s.io/external-resizer-cfg created
rolebinding.rbac.authorization.k8s.io/csi-resizer-role-cfg created
service/csi-cinder-controller-service created
statefulset.apps/csi-cinder-controllerplugin created
serviceaccount/csi-cinder-node-sa created
clusterrole.rbac.authorization.k8s.io/csi-nodeplugin-role created
clusterrolebinding.rbac.authorization.k8s.io/csi-nodeplugin-binding created
daemonset.apps/csi-cinder-nodeplugin created
csidriver.storage.k8s.io/cinder.csi.openstack.org created
secret/cloud-config created

```

Show moreShow more icon

从上面的输出可以看出为了运行 CSI 我们创建的资源主要分为几类：

- clusterrole/role

    由于 kubernetes 对于权限的控制和要求，需要定义一系列的 role 和 clusterrole，以 csi-nodeplugin-role 为例，定义了该 clusterrole 可以针对 events 这个资源进行 [“get”, “list”, “watch”, “create”, “update”, “patch”] 等操作。





    ```
    kind: ClusterRole
    apiVersion: rbac.authorization.k8s.io/v1
    metadata:
    name: csi-nodeplugin-role
    rules:
      - apiGroups: [""]
        resources: ["events"]
    verbs: ["get", "list", "watch", "create", "update", "patch"]

    ```





    Show moreShow more icon

- serviceaccount

    serviceaccount 针对运行在 Pod 的进程进行授权，以如下 csi-cinder-node-sa 为例，定义了一个 serviceaccount 被 clusterrolebinding 使用。





    ```
    apiVersion: v1
    kind: ServiceAccount
    metadata:
    name: csi-cinder-node-sa
    namespace: kube-system

    ```





    Show moreShow more icon

- clusterrolebinding

    有了 clusterrole 和 serviceaccount 之后，需要通过 clusterrolebinding 把二者结合在一起，以如下 csi-nodeplugin-binding 为例，该 clusterrolebinding 使得拥有 csi-cinder-node-sa serviceaccount 的 pod 拥有 csi-nodeplugin-role 这个 role 来进行 API 操作。





    ```
    kind: ClusterRoleBinding
    apiVersion: rbac.authorization.k8s.io/v1
    metadata:
    name: csi-nodeplugin-binding
    subjects:
      - kind: ServiceAccount
        name: csi-cinder-node-sa
        namespace: kube-system
    roleRef:
    kind: ClusterRole
    name: csi-nodeplugin-role
    apiGroup: rbac.authorization.k8s.io

    ```





    Show moreShow more icon

- csidriver

    定义以 cinder.csi.openstack.org 为名字的 CSIDriver。





    ```
    apiVersion: storage.k8s.io/v1beta1
    kind: CSIDriver
    metadata:
    name: cinder.csi.openstack.org
    spec:
    attachRequired: true
    podInfoOnMount: false

    ```





    Show moreShow more icon

- secret

    cloud-config 这个 secret 里存放了 Openstack cloud 的信息，所有的后续的通过 openstack API 的调用操作都需要该文件完成认证，授权和服务发现；里是通过 base64 转码之后的 cloud 配置内容。





    ```
    apiVersion: v1
    metadata:
    name: cloud-config
    namespace: kube-system
    data:
    cloud.conf: <conf>
    statefulset

    ```





    Show moreShow more icon

    csi-cinder-controllerplugin 定义了 CSI controller，也就是运行了一系列的 CSI 容器来完成 CSI 的控制工作，具体定义请参考 [这里](https://github.com/kubernetes/cloud-provider-openstack/blob/master/manifests/cinder-csi-plugin/cinder-csi-controllerplugin.yaml)。

- daemonset

    csi-cinder-nodeplugin 定义了每一个计算节点（node）上要运行的服务，具体定义请参考 [这里](https://github.com/kubernetes/cloud-provider-openstack/blob/master/manifests/cinder-csi-plugin/cinder-csi-nodeplugin.yaml)。


创建完成之后，可以通过如下命令来查看已经创建的 pod 以及相关 container 的定义：

```
# kubectl get pods --all-namespaces
NAMESPACE     NAME                            READY   STATUS    RESTARTS   AGE
kube-system   csi-cinder-controllerplugin-0   5/5     Running   0          24s
kube-system   csi-cinder-nodeplugin-ssmc6     2/2     Running   0          24s
kube-system   kube-dns-547db76c8f-ckwrp       3/3     Running   0          2m50s

```

Show moreShow more icon

获取容器所用的镜像名字。

```
# kubectl get pods --all-namespaces -o=jsonpath='{range .items[*]}{"\n"}{.metadata.name}{":\t"}{range .spec.containers[*]}{.image}{", "}{end}{end}'

csi-cinder-controllerplugin-0:  quay.io/k8scsi/csi-attacher:v1.2.1, quay.io/k8scsi/csi-provisioner:v1.3.0, quay.io/k8scsi/csi-snapshotter:v1.2.0, quay.io/k8scsi/csi-resizer:v0.2.0, docker.io/k8scloudprovider/cinder-csi-plugin:latest,
csi-cinder-nodeplugin-ssmc6:    quay.io/k8scsi/csi-node-driver-registrar:v1.1.0, docker.io/k8scloudprovider/cinder-csi-plugin:latest

```

Show moreShow more icon

我们可以看出创建了两个 pod，分别是 controller 和 node(笔者的环境是 local 所以 controller 和 node 是在一起的)， controller 这个 pod 运行了 5 个容器而 node 这个 pod 运行了两个容器。为了使得开发 CSI driver 简便，Kubernetes CSI 社区提供了一系列的 sidecar 容器，这样开发者只需要关注支持自己存储的逻辑, 除非有自己的定制需求，可以直接使用社区提供的容器。

- attacher 这是一个边车（sidecar）容器，主要是 watch VolumeAttachment 这个对象并且针对 CSI 对象调用 PublishVolume 和 UnpublisVolume 这个动作，具体可以参考 [这里](https://kubernetes-csi.github.io/docs/external-attacher.html)。
- provisioner 这是一个边车（sidecar）容器，主要是 watch PersistentVolumeClaim 这个对象并且针对 CSI 对象调用 CreateVolume 这个动作，具体可以参考 [这里](https://kubernetes-csi.github.io/docs/external-provisioner.html)。
- snapshotter 这是一个边车（sidecar）容器，主要是 VolumeSnapshot 和 VolumeSnapshotContent 这两个对象，具体可以参考 [这里](https://kubernetes-csi.github.io/docs/external-snapshotter.html)。
- resizer 这是一个边车（sidecar）容器，主要是 watch PersistentVolumeClaim 这个对象并且针对 CSI 对象调用 ControllerExpandVolume 这个动作，具体可以参考 [这里](https://kubernetes-csi.github.io/docs/external-resizer.html)。
- registrar 这是一个边车（sidecar）容器，主要是从 CSI 获取 driver 的信息并向 kubelet 注册，具体可以参考 [这里](https://kubernetes-csi.github.io/docs/node-driver-registrar.html)。
- cinder-csi-plugin 是主要的服务容器，运行 CSI driver，针对 controller 和 node 有不同的配置，具体可以参考上文 csi-cinder-controllerplugin 和 csi-cinder-nodeplugin 的定义。

## 使用 CSI

通过上面的准备工作，我们已经准备好所需要的环境，下面我们开始创建工作容器来通过 CSI 来分配和挂载 openstack cinder 的卷，在执行如下操作之前，请确定上面的 CSI 的 pods 都正常运行。 接着，我们通过命令创建如下定义：

```
# kubectl apply -f examples/cinder-csi-plugin/nginx.yaml
storageclass.storage.k8s.io/csi-sc-cinderplugin created
persistentvolumeclaim/csi-pvc-cinderplugin created
pod/nginx created

```

Show moreShow more icon

创建的资源主要分为如下几类：

- StorageClass

    首先需要创建一个名为 csi-sc-cinderplugin 的 storageclass, 并且它的后端是 cinder.csi.openstack.org 这个 driver， 这和上述的 csidriver 是对应的，如果不指定或者指定错了该 driver 名字，会导致后续错误。





    ```
    apiVersion: storage.k8s.io/v1
    kind: StorageClass
    metadata:
    name: csi-sc-cinderplugin
    provisioner: cinder.csi.openstack.org

    ```





    Show moreShow more icon

- PersistentVolumeClaim

    接着，我们需要创建一个名为 `csi-pvc-cinderplugin` 的 [PersistentVolumeClaim](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims)，从如下定义我们可以看出，该 PersistentVolumeClaim 是 RWO（ReadWriteOnce），每一个 PVC 都有一个 [访问模式](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes)，csi-pvc-cinderplugin 的访问模式就是 RWO，代表该 volume 只能被一个 node 挂载成读写模式，这个 volume 的大小为 1G。当该 PersistentVolumeClaim 被创建之后，Cloud-Provider-Openstack 会通过 cinder 接口来在 OpenStack 中创建 volume 来满足这个需求。





    ```
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
    name: csi-pvc-cinderplugin
    spec:
    accessModes:
      - ReadWriteOnce
    resources:
        requests:
          storage: 1Gi
    storageClassName: csi-sc-cinderplugin

    ```





    Show moreShow more icon

- Pod

    最后，创建一个 Pod 来使用创建的 volume，也就是该 Pod 使用了我们创建的 csi-data-cinderplugin 这个 PersistentVolumeClaim，从 yaml 定义，我们也可以看出 csi-data-cinderplugin 这个 volume 会被挂载到 `/var/lib/www/html` 这个位置。





    ```
    apiVersion: v1
    kind: Pod
    metadata:
    name: nginx
    spec:
    containers:
      - image: nginx
        imagePullPolicy: IfNotPresent
        name: nginx
        ports:
        - containerPort: 80
          protocol: TCP
        volumeMounts:
          - mountPath: /var/lib/www/html
            name: csi-data-cinderplugin
    volumes:
      - name: csi-data-cinderplugin
        persistentVolumeClaim:
          claimName: csi-pvc-cinderplugin
          readOnly: false

    ```





    Show moreShow more icon


现在我们可以通过 openstack 和 Kubernetes 命令分别查看，有一个 1G 的卷已经被分配并且挂载到相应的容器之上了，并且该 volume 已经是 in-use(挂载状态)了。

通过 openstack 命令，这个 volume 大小为 1G，并且 attach 为 jj-cloud-provider 这个机器的 `/dev/vdb` 设备了。

```
# openstack volume list
+--------------------------------------+------------------------------------------+--------+------+-----------------------------------------------------+
| ID                                   | Display Name                             | Status | Size | Attached to                                         |
+--------------------------------------+------------------------------------------+--------+------+-----------------------------------------------------+
| c78334f9-b296-436e-bde8-b8e51e1de458 | pvc-631cef46-a56c-46ad-89ce-60ec23acdf04 | in-use |    1 | Attached to jj-cloud-provider on /dev/vdb           |

```

Show moreShow more icon

接着我们查看 PersistentVolumeClaim 的状态：

```
# kubectl get pvc
NAME                   STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS          AGE
csi-pvc-cinderplugin   Bound    pvc-631cef46-a56c-46ad-89ce-60ec23acdf04   1Gi        RWO            csi-sc-cinderplugin   21s

```

Show moreShow more icon

从操作系统层面，我们也可以在 worker node 之上我们可以看到有一块盘已经创建:

```
# ls /dev/vdb
/dev/vdb

# mount | grep vdb
/dev/vdb on /var/lib/kubelet/plugins/kubernetes.io/csi/pv/pvc-631cef46-a56c-46ad-89ce-60ec23acdf04/globalmount type ext4 (rw,relatime,data=ordered)
/dev/vdb on /var/lib/kubelet/pods/c9097d48-3021-42e6-a5b5-15bc354bbf2b/volumes/kubernetes.io~csi/pvc-631cef46-a56c-46ad-89ce-60ec23acdf04/mount type ext4 (rw,relatime,data=ordered)

```

Show moreShow more icon

## 删除卷

最后，我们清除之前的 Pod 和 PersistentVolumeClaim，可以看出，当 pod 被删除之后，之前的创建的 volume 已经恢复到 available 状态（也就是非挂载状态）。
首先，我们删除创建的工作 Pod：

```
# kubectl delete pod nginx
pod "nginx" deleted

```

Show moreShow more icon

当容器删除之后，该 Volume 已经不再挂载在工作节点上了：

```
# ls /dev/vdb
ls: cannot access '/dev/vdb': No such file or directory

```

Show moreShow more icon

同时 OpenStack 的 Volume 已经变为 available 状态：

```
# openstack volume list
+--------------------------------------+------------------------------------------+-----------+------+-----------------------------------------------------+
| ID                                   | Display Name                             | Status    | Size | Attached to                                         |
+--------------------------------------+------------------------------------------+-----------+------+-----------------------------------------------------+
| c78334f9-b296-436e-bde8-b8e51e1de458 | pvc-631cef46-a56c-46ad-89ce-60ec23acdf04 | available |    1 |

```

Show moreShow more icon

删除 PersistentVolumeClaim 并检查系统使用资源状态：

```
# kubectl get pvc
NAME                   STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS          AGE
csi-pvc-cinderplugin   Bound    pvc-631cef46-a56c-46ad-89ce-60ec23acdf04   1Gi        RWO            csi-sc-cinderplugin   13m

# kubectl delete pvc csi-pvc-cinderplugin
persistentvolumeclaim "csi-pvc-cinderplugin" deleted

# kubectl get pv
No resources found in default namespace.
# kubectl get pvc
No resources found in default namespace.

```

Show moreShow more icon

## 进阶功能

上面我们介绍了最简单的通过 cinder 来创建和删除块设备功能，除此之外，我们还可以使用多种进阶功能。

- [块镜像](https://github.com/kubernetes/cloud-provider-openstack/blob/master/docs/using-cinder-csi-plugin.md#example-snapshot-create-and-restore)：通过 Cinder 的镜像（snapshot）功能，可以创建 snapshot 作为 volumesnapshot.snapshot.storage.k8s.io 这个资源在 OpenStack 上的实现。
- [块扩展](https://github.com/kubernetes/cloud-provider-openstack/blob/master/docs/using-cinder-csi-plugin.md#volume-expansion)：通过 Cinder 的 Expand 功能，可以在线扩展(不能缩小)某块 Volume。
- [块克隆](https://github.com/kubernetes/cloud-provider-openstack/blob/master/docs/using-cinder-csi-plugin.md#volume-cloning)：通过 Cinder 的 Clone 功能，完成对某快 Volume 的克隆操作。

## 结束语

本文首先简介了 CSI 定义以及其使用场景，包括所有的 side car 的容器定义以及如何部署这些容器环境，然后以 OpenStack cinder 为例介绍了如何使用 Kubernetes 通过 CSI 来从 Cloud Provider OpenStack 申请和管理卷块设备资源。