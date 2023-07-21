package main

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/kubernetes/mykubelet/mytest/myclient/myclientlib"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/kubelet/config"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
	"reflect"
)

// 手动创建podConfig，监听事件
func main() {
	client := myclientlib.InitClient()
	nodeName := "lain1" // 测试的节点

	fact := informers.NewSharedInformerFactory(client, 0)
	fact.Core().V1().Nodes().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{})
	ch := make(chan struct{})
	fact.Start(ch)
	if waitMap := fact.WaitForCacheSync(ch); waitMap[reflect.TypeOf(&v1.Node{})] {
		// 事件分发器广播(分发给watch它的函数，用channel实现)
		eventBroadcaster := record.NewBroadcaster()
		// 事件记录器(如Pod生命周期事件、各种错误事件)
		eventRecorder := eventBroadcaster.NewRecorder(legacyscheme.Scheme, v1.EventSource{Component: "kubelet", Host: nodeName})
		// 创建PodConfig
		podConfig := config.NewPodConfig(config.PodConfigNotificationIncremental, eventRecorder)
		// 注入clientset
		config.NewSourceApiserver(client, types.NodeName(nodeName), func() bool {
			return fact.Core().V1().Nodes().Informer().HasSynced()
		}, podConfig.Channel(kubetypes.ApiserverSource)) // 关联configCh，会把相关的内容注入到ch里

		fmt.Println("开始监听")
		for item := range podConfig.Updates() { // updates()返回的就是configCh，当pod产生新增、删除等变化会出现新的事件产生
			for _, pod := range item.Pods {
				fmt.Println(pod.Name)
			}
		}
	}
}
