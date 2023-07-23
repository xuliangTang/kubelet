package mycore

import (
	"k8s.io/apimachinery/pkg/util/clock"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/kubelet/cri/remote"
	"k8s.io/kubernetes/pkg/kubelet/pleg"
	"log"
	"time"
)

const (
	ContainerdAddress = "tcp://110.188.24.175:8989"
)

// StartPleg 启动pleg，需要传入共享的clock和podCache
func StartPleg(ck clock.RealClock, cache kubecontainer.Cache) {
	// 创建原始cri RuntimeService
	rs, err := remote.NewRemoteRuntimeService(ContainerdAddress, time.Second*3)
	if err != nil {
		log.Fatalln(err)
	}
	// 创建kubelet自己封装的ContainerRuntime
	cr := NewContainerRuntime(rs)

	// 启动pleg 每隔1s会请求cri接口比对
	p := pleg.NewGenericPLEG(cr, 1000, time.Second*1, cache, ck)
	p.Start()
}
