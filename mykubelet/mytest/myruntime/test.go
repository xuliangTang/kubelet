package main

import (
	"fmt"
	"k8s.io/kubernetes/mykubelet/mylib"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"log"
)

func main() {
	// 真正的cri模拟实现
	rs := &mylib.MyRuntimeService{}
	// 模拟创建kubelet封装的runtime
	var cr kubecontainer.Runtime = mylib.NewContainerRuntime(rs, "containerd")

	pods, err := cr.GetPods(true)
	if err != nil {
		log.Fatalln(err)
	}

	// 遍历pod和容器
	for _, pod := range pods {
		fmt.Print("pod:", pod.Name, "的容器有：")
		for _, c := range pod.Containers {
			fmt.Print(c.Name, " ")
		}
		fmt.Println()
	}
}
