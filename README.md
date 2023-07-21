## 源码
### 节点状态
```
pkg/kubelet/kubelet_node_status.go 588行
```
kubectl是怎么判断Node是Ready的
```
pkg/printers/internalversion/printers.go
```

### kubeClient初始化
```
cmd/kubelet/app/server.go 604行
```

### kubelet续期Lease对象
[参考文档](https://kubernetes.io/zh-cn/docs/reference/config-api/kubelet-config.v1beta1/) 搜索`nodeLeaseDurationSeconds`
```
pkg/kubelet/kubelet.go 863行
```

### kubeadm BootstrapToken
```
cmd/kubeadm/app/cmd/init.go 190行
```

### kubelet生成私钥和证书请求
```
pkg/kubelet/certificate/kubelet.go 90行
staging/src/k8s.io/client-go/util/certificate/certificate_manager.go 640行
```

### kubelet创建csr资源
```
staging/src/k8s.io/client-go/util/certificate/csr/csr.go 52行
```

### 租约控制器
```
staging/src/k8s.io/component-helpers/apimachinery/lease/controller.go 74行
```

### pleg模块
Pod Lifecycle Event Generator(Pod生命周期生成器): 定期检查节点上Pod运行状态，把Pod的状态变化封装为特有的Event(PodLifecycleEvent)，从而触发kubelet的主同步机制 

pleg是怎么判断容器发生了变化？（如新增删除Pod）

通过relist()函数获取Pod列表并本地缓存(podCache)，然后定时再取，每次都和之前都缓存比对，从而就知道哪些Pod发生了变化，从而生成相关都Pod生命周期事件和更改后都状态
```
pkg/kubelet/kubelet.go 1499行
```
主要的runtime接口
```
pkg/kubelet/kuberuntime/kuberuntime_manager.go 164行
```

### syncLoop主循环
具体处理在syncLoopIteration()函数里：
1. 监听PodConfig(对静态及动态Pod的变更产生的事件)的变更：configCh
2. 监听pleg的事件：plegCh
3. 监听等待同步的Pod：syncCh
4. 监听清理Pod的事件：housekeepingCh
5. 监听来自(如存活、就绪、启动)探针的事件：livenessManager readinessManager startupManager
```
pkg/kubelet/kubelet.go 2022行
```

### 静态Pod初始化
静态Pod对应 /etc/kubernetes/manifests 目录的文件
```
pkg/kubelet/config/config.go 70行
```

### statusManager
主要功能是将Pod状态信息同步到k8s apiserver，它不会直接监控Pod的状态，而是提供接口供其他manager(如probeManager)进行调用，同时syncLoop主循环也会调用到它

它暴露的几个主要方法：
- SetPodStatus(): Pod状态发生变化时调用，更新状态到apiserver
- SetContainerReadiness(): Pod中容器的健康状态发生变化时调用，修改Pod的健康状态
- TerminatePod(): 删除Pod时调用，把Pod中所有的容器设置为terminated状态
- RemoveOrphanedStatuses(): 删除孤儿Pod时调用
```
pkg/kubelet/status/status_manager.go 93行
调用在 pkg/kubelet/kubelet.go 1491行
初始化在 pkg/kubelet/status/status_manager.go 122行
```
初始化的参数：
- kubeClient: 用于和apiserver交互
- podManager: pod内存形式的管理器(管理kubelet对pod的访问)
- podStatuses: pod与状态的对应关系
- podStatusesChannel: 当其他组件调用statusManager更新pod状态时，会调用这个channel
- apiStatusVersions: 维护最新的pod status版本号，每更新一次会加1
- podDeletionSafety: 删除pod的接口

### podManager
用来在本地缓存存储Pod相关资源
```
pkg/kubelet/pod/pod_manager.go 128行
调用在 pkg/kubelet/kubelet.go 2114行
初始化在 pkg/kubelet/kubelet.go 624行
```

### mirrorPod
这种类型的Pod来源于静态Pod
1. 静态Pod不受apiserver管理，而且无法移动调度到别的节点
2. 因此类似这种Pod在podManager会创建一个类似副本来进行查看

### podConfig
podManager里的数据从哪里收集？

podConfig的updates()方法返回里configCh，用来监听file、http(这2个属于静态Pod)和apiserver(动态Pod)的事件更新。得到数据后塞入PodManager
```
初始化在 pkg/kubelet/kubelet.go 434行
pkg/kubelet/kubelet.go 2099行调用了
```
SyncHandler接口用来对PodManager缓存新增、更新、删除等操作
```
pkg/kubelet/kubelet.go 195行
```

### podWorker
PodWorker是kubelet中用于管理每个Pod的协程角色
1. 每创建一个新的Pod，都会为其配置一个专有的podWorker
2. 每个podWorker都是一个协程，它会创建一个类型为UpdatePodOptions(pod更新事件)的channel
3. 获得pod的更新事件后调用podWorker中syncPodFn(Kubelet中的syncPod)函数进行具体的同步工作(syncPod用来将Pod的最新状态上报给apiServer、创建pod的专属目录等)
```
初始化在 pkg/kubelet/kubelet.go 655行
```
其中**managerPodLoop()**函数的基本作用就是监听podUpdates更新事件，从而触发PodSyncFn
```
调用在 pkg/kubelet/pod_workers.go 750行左右 
```
里面有个阻塞函数`p.podCache.GetNewerThan(pod.UID, lastSyncTime)`，它会等待podCache(本地pod和状态的映射关系map)有针对这个Pod的状态数据，才会继续往下执行