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

通过relist()函数获取Pod列表并本地缓存，然后定时再取，每次都和之前都缓存比对，从而就知道哪些Pod发生了变化，从而生成相关都Pod生命周期事件和更改后都状态
```
pkg/kubelet/kubelet.go 1499行
```
主要的runtime接口
```
pkg/kubelet/kuberuntime/kuberuntime_manager.go 164行
```

### syncLoop主循环
具体处理在syncLoopIteration()函数里：
1. 监听PodConfig(对静态Pod的变更产生的事件)的变更：configCh
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
```
pkg/kubelet/pod/pod_manager.go 128行
调用在 pkg/kubelet/kubelet.go 2114行
初始化在 pkg/kubelet/kubelet.go 624行
```

### mirrorPod
这种类型的Pod来源于静态Pod
1. 静态Pod不受apiserver管理，而且无法移动调度到别的节点
2. 因此类似这种Pod在podManager会创建一个类似副本来进行查看
