package main

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/mykubelet/mytest/myclient/myclientlib"
	"k8s.io/kubernetes/pkg/kubelet/configmap"
	kubepod "k8s.io/kubernetes/pkg/kubelet/pod"
	"k8s.io/kubernetes/pkg/kubelet/secret"
	"log"
	"reflect"
)

// 手动调用podManager创建mirrorPod, nodeName必须是当前kubelet节点名称
func main() {
	client := myclientlib.InitClient()

	fact := informers.NewSharedInformerFactory(client, 0)
	fact.Core().V1().Nodes().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{})
	nodeLister := fact.Core().V1().Nodes().Lister()
	ch := make(chan struct{})
	fact.Start(ch)
	if waitMap := fact.WaitForCacheSync(ch); waitMap[reflect.TypeOf(&v1.Node{})] {
		// 创建podManager
		secretManager := secret.NewSimpleSecretManager(client)
		configMapManager := configmap.NewSimpleConfigMapManager(client)
		mirrorPodClient := kubepod.NewBasicMirrorClient(client, "mylain", nodeLister)
		podManager := kubepod.NewBasicPodManager(mirrorPodClient, secretManager, configMapManager)

		// 创建静态pod
		pod := &v1.Pod{}
		pod.Name = "mystatic-mylain"
		pod.Namespace = "default"
		pod.Spec.NodeName = "mylain"
		pod.Spec.Containers = []v1.Container{
			{
				Name:  "mystatic-container",
				Image: "mystatic:alpine",
			},
		}
		if err := podManager.CreateMirrorPod(pod); err != nil {
			log.Fatalln(err)
		}
	}
}
