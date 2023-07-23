package mycore

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"k8s.io/kubernetes/pkg/kubelet/container"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
	"time"
)

// SetPodReady 构建PodStatus的sandbox为ready containers为running
func SetPodReady(pod *v1.Pod) *container.PodStatus {
	status := &container.PodStatus{
		ID:        types.UID(pod.UID),
		Name:      pod.Name,
		Namespace: pod.Namespace,
		SandboxStatuses: []*runtimeapi.PodSandboxStatus{
			{
				Id:    string(pod.UID),
				State: runtimeapi.PodSandboxState_SANDBOX_READY,
			},
		},
	}

	var containerStatus []*container.Status
	for _, c := range pod.Spec.Containers {
		cs := &container.Status{
			Name:      c.Name,
			Image:     c.Image,
			State:     container.ContainerStateRunning,
			CreatedAt: time.Now(),
			StartedAt: time.Now().Add(time.Second * 3),
		}
		containerStatus = append(containerStatus, cs)
	}
	status.ContainerStatuses = containerStatus

	return status
}

// HandlePodAdditions podConfig获取到add后的处理函数
func HandlePodAdditions(podCache *PodCache, pods []*v1.Pod) {
	for _, pod := range pods {
		// 加入缓存
		podCache.PodManager.AddPod(pod)

		// 模拟执行dispatchWork 创建podWorker，开启协程监听pod状态
		// 它会等待podCache有针对这个pod的数据，然后执行syncPod或syncTerminate把本地pod状态同步给apiserver
		podCache.PodWorkers.UpdatePod(UpdatePodOptions{
			Pod:        pod,
			MirrorPod:  nil,
			UpdateType: kubetypes.SyncPodCreate,
			StartTime:  podCache.Clock.Now(),
		})
	}
}

// HandlePodUpdates podConfig获取到update和delete后的处理函数
func HandlePodUpdates(podCache *PodCache, pods []*v1.Pod) {
	for _, pod := range pods {
		// 更新缓存
		podCache.PodManager.UpdatePod(pod)

		podCache.PodWorkers.UpdatePod(UpdatePodOptions{
			Pod:        pod,
			MirrorPod:  nil,
			UpdateType: kubetypes.SyncPodUpdate,
			StartTime:  podCache.Clock.Now(),
		})
	}
}

// HandlePodRemoves podConfig获取到remove后的处理函数
func HandlePodRemoves(podCache *PodCache, pods []*v1.Pod) {
	for _, pod := range pods {
		// 移除缓存
		podCache.PodManager.DeletePod(pod)

		podCache.PodWorkers.UpdatePod(UpdatePodOptions{
			Pod:        pod,
			UpdateType: kubetypes.SyncPodKill,
		})
	}
}
