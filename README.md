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