package mycore

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/kubelet/config"
	"k8s.io/kubernetes/pkg/kubelet/configmap"
	"k8s.io/kubernetes/pkg/kubelet/pod"
	"k8s.io/kubernetes/pkg/kubelet/secret"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
)

type PodCache struct {
	client     *kubernetes.Clientset
	PodManager pod.Manager
	PodConfig  *config.PodConfig
}

// NewPodCache 初始化PodConfig
func NewPodCache(client *kubernetes.Clientset, nodeName string) *PodCache {
	fact := informers.NewSharedInformerFactory(client, 0)
	fact.Core().V1().Nodes().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{})
	ch := make(chan struct{})
	fact.Start(ch)

	// 创建podManager
	nodeLister := fact.Core().V1().Nodes().Lister()
	secretManager := secret.NewSimpleSecretManager(client)
	configMapManager := configmap.NewSimpleConfigMapManager(client)
	mirrorPodClient := pod.NewBasicMirrorClient(client, "mylain", nodeLister)
	podManager := pod.NewBasicPodManager(mirrorPodClient, secretManager, configMapManager)

	return &PodCache{
		client:     client,
		PodManager: podManager,
		PodConfig:  newPodConfig(client, fact, nodeName),
	}
}

// 创建podConfig
func newPodConfig(client *kubernetes.Clientset, fact informers.SharedInformerFactory, nodeName string) *config.PodConfig {
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

	return podConfig
}
