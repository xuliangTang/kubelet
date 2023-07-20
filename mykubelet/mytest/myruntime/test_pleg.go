package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/clock"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"k8s.io/kubernetes/mykubelet/mylib"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/kubelet/pleg"
	"net/http"
	"time"
)

func main() {
	// 真正的cri模拟实现
	rs := &mylib.MyRuntimeService{}
	// 模拟创建 kubelet 封装的 runtime
	var cr kubecontainer.Runtime = mylib.NewContainerRuntime(rs, "containerd")

	// 手动创建pleg
	cache := kubecontainer.NewCache()
	p := pleg.NewGenericPLEG(cr, 1000, time.Second*1, cache, clock.RealClock{})
	go func() {
		for {
			select {
			case v := <-p.Watch(): // 获取Pod生命周期事件
				fmt.Println(v)
			}
		}
	}()
	p.Start()

	// 手动测试pod状态变更
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		// 手动修改第一个pod的状态为notReady
		mylib.MockData_Pods[0].State = runtimeapi.PodSandboxState_SANDBOX_NOTREADY
		writer.Write([]byte("手动修改状态为notReady"))
	})
	http.ListenAndServe(":8080", nil)
}
