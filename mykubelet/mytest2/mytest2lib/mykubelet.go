package mytest2lib

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/kubelet/config"
	"k8s.io/kubernetes/pkg/kubelet/configmap"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/kubelet/pleg"
	"k8s.io/kubernetes/pkg/kubelet/pod"
	"k8s.io/kubernetes/pkg/kubelet/prober"
	"k8s.io/kubernetes/pkg/kubelet/prober/results"
	"k8s.io/kubernetes/pkg/kubelet/secret"
	"k8s.io/kubernetes/pkg/kubelet/status"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
	"k8s.io/kubernetes/pkg/kubelet/util/queue"
	"time"
)

const NodeName = "mylain"

type MyKubelet struct {
	KubeClient *kubernetes.Clientset
	PodManager pod.Manager
	PodConfig  *config.PodConfig
	PodWorkers PodWorkers
	Pleg       pleg.PodLifecycleEventGenerator

	Clock    clock.Clock
	PodCache kubecontainer.Cache

	statusManager status.Manager
	reasonCache   *ReasonCache
	probeManager  prober.Manager
}

func NewMyKubelet(client *kubernetes.Clientset) *MyKubelet {
	fact := informers.NewSharedInformerFactory(client, 0)
	//fact.Core().V1().Nodes().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{})
	nodeLister := fact.Core().V1().Nodes().Lister()
	ch := make(chan struct{})
	fact.Start(ch)
	fact.WaitForCacheSync(ch)

	// 初始化podManager
	mirrorPodClient := pod.NewBasicMirrorClient(client, NodeName, nodeLister)
	secretManager := secret.NewSimpleSecretManager(client)
	configMapManager := configmap.NewSimpleConfigMapManager(client)
	podManager := pod.NewBasicPodManager(mirrorPodClient, secretManager, configMapManager)

	// 初始化podConfig
	eventBroadcaster := record.NewBroadcaster()
	eventRecorder := eventBroadcaster.NewRecorder(legacyscheme.Scheme, v1.EventSource{Component: "kubelet", Host: NodeName})
	podConfig := config.NewPodConfig(config.PodConfigNotificationIncremental, eventRecorder)
	// 注入clientset
	config.NewSourceApiserver(client, types.NodeName(NodeName), func() bool {
		return fact.Core().V1().Nodes().Informer().HasSynced()
	}, podConfig.Channel(kubetypes.ApiserverSource)) // 关联configCh，会把相关的内容注入到ch里

	// 初始化pleg(mac不兼容)
	//runtimeSvc, err := remote.NewRemoteRuntimeService("tcp://110.188.24.175:8989", time.Second*3)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//runtime := mycore.NewContainerRuntime(runtimeSvc)
	podCache := kubecontainer.NewCache()
	ck := &clock.RealClock{}
	//plg := pleg.NewGenericPLEG(runtime, 1000, time.Second, podCache, ck)

	mykubelet := &MyKubelet{
		KubeClient: client,
		PodManager: podManager,
		PodConfig:  podConfig,
		//Pleg:       plg,
		PodCache:    podCache,
		Clock:       ck,
		reasonCache: NewReasonCache(),
	}

	// 初始化podWorker
	workQueue := queue.NewBasicWorkQueue(ck)
	podWorker := NewPodWorkers(
		mykubelet.syncPod,
		mykubelet.syncTerminatingPod,
		mykubelet.syncTerminatedPod,
		eventRecorder,
		workQueue,
		time.Second*1,
		time.Second*10,
		podCache,
	)
	mykubelet.PodWorkers = podWorker

	// 初始化statusManager
	mykubelet.statusManager = status.NewManager(client, mykubelet.PodManager, mykubelet)
	// 启动
	mykubelet.statusManager.Start()

	// 初始化probeManager
	lm, rm, sm := results.NewManager(), results.NewManager(), results.NewManager()
	mykubelet.probeManager = prober.NewManager(
		mykubelet.statusManager,
		lm,
		rm,
		sm,
		&ContainerCommandRunner{},
		eventRecorder)

	return mykubelet
}

func (m MyKubelet) HandlePodAdditions(pods []*v1.Pod) {
	for _, p := range pods {
		m.PodManager.AddPod(p)
		m.dispatchWork(kubetypes.SyncPodCreate, p, m.Clock.Now())
	}
}

func (m MyKubelet) HandlePodUpdates(pods []*v1.Pod) {
	for _, p := range pods {
		m.PodManager.UpdatePod(p)
		m.dispatchWork(kubetypes.SyncPodUpdate, p, m.Clock.Now())
	}
}

func (m MyKubelet) HandlePodRemoves(pods []*v1.Pod) {
	for _, p := range pods {
		m.PodManager.DeletePod(p)
		m.dispatchWork(kubetypes.SyncPodKill, p, m.Clock.Now())
	}
}

func (m MyKubelet) dispatchWork(updateType kubetypes.SyncPodType, pod *v1.Pod, start time.Time) {
	m.PodWorkers.UpdatePod(UpdatePodOptions{
		UpdateType: updateType,
		Pod:        pod,
		StartTime:  start,
	})
}

func (m *MyKubelet) syncPod(ctx context.Context, updateType kubetypes.SyncPodType, pod, mirrorPod *v1.Pod, podStatus *kubecontainer.PodStatus) (isTerminal bool, err error) {
	fmt.Println("测试的syncPod")

	isTerminal = false
	apiPodStatus := m.generateAPIPodStatus(pod, podStatus)
	m.statusManager.SetPodStatus(pod, apiPodStatus)
	if apiPodStatus.Phase == v1.PodSucceeded || apiPodStatus.Phase == v1.PodFailed {
		isTerminal = true
	}
	return isTerminal, nil
}

func (m *MyKubelet) syncTerminatingPod(ctx context.Context, pod *v1.Pod, podStatus *kubecontainer.PodStatus, runningPod *kubecontainer.Pod, gracePeriod *int64, podStatusFn func(*v1.PodStatus)) error {
	fmt.Println("测试的syncTerminatingPod")
	return nil
}

func (m *MyKubelet) syncTerminatedPod(ctx context.Context, pod *v1.Pod, podStatus *kubecontainer.PodStatus) error {
	fmt.Println("测试的syncTerminatingPod")
	return nil
}

type SyncHandler interface {
	HandlePodAdditions(pods []*v1.Pod)
	HandlePodUpdates(pods []*v1.Pod)
	HandlePodRemoves(pods []*v1.Pod)
}

func (m MyKubelet) PodResourcesAreReclaimed(pod *v1.Pod, status v1.PodStatus) bool {
	return true
}

func (m MyKubelet) PodCouldHaveRunningContainers(pod *v1.Pod) bool {
	return true
}

type ContainerCommandRunner struct{}

func (c ContainerCommandRunner) RunInContainer(id kubecontainer.ContainerID, cmd []string, timeout time.Duration) ([]byte, error) {
	return []byte(""), nil
}

var _ SyncHandler = &MyKubelet{}
var _ status.PodDeletionSafetyProvider = &MyKubelet{}
var _ kubecontainer.CommandRunner = &ContainerCommandRunner{}
