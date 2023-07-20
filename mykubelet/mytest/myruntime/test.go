package main

import (
	"fmt"
	"k8s.io/kubernetes/mykubelet/mylib"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
)

func main() {
	// 真正的cri模拟实现
	rs := &mylib.MyRuntimeService{}
	// 模拟创建kubelet封装的runtime
	var cr kubecontainer.Runtime = mylib.NewContainerRuntime(rs, "containerd")

	fmt.Println(cr.GetPods(true))
}
