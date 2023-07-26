package main

import (
	"encoding/json"
	"fmt"
	v1 "k8s.io/api/core/v1"
	apitypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/mykubelet/mytest/myclient/myclientlib"
	"k8s.io/kubernetes/mykubelet/mytest/myruntime/mycore"
	"k8s.io/kubernetes/mykubelet/mytest2/mytest2lib"
	"k8s.io/kubernetes/pkg/kubelet/types"
	"log"
	"net/http"
)

func main() {
	/*runtimeSvc := &mytest2lib.MyRuntimeService{}
	runtime := mycore.NewContainerRuntime(runtimeSvc)

	cache := container.NewCache()

	p := pleg.NewGenericPLEG(runtime, 1000, time.Second, cache, clock.RealClock{})
	go func() {
		for {
			select {
			case v := <-p.Watch():
				fmt.Println(v)
			}
		}
	}()
	p.Start()
	select {}*/

	client := myclientlib.InitClient()
	mykubelet := mytest2lib.NewMyKubelet(client)

	fmt.Println("开始监听")

	go func() {
		http.HandleFunc("/addpod", func(writer http.ResponseWriter, request *http.Request) {
			pod := &v1.Pod{}
			pod.Name = "test-mylain"
			pod.Namespace = "default"
			pod.Spec.NodeName = "mylain"
			pod.Spec.Containers = []v1.Container{
				{
					Name:  "tst",
					Image: "nginx:1.18-alpine",
				},
			}
			if err := mykubelet.PodManager.CreateMirrorPod(pod); err != nil {
				log.Fatalln(err)
			}
		})

		http.HandleFunc("/pods", func(writer http.ResponseWriter, request *http.Request) {
			pods := mykubelet.PodManager.GetPods()
			var podNames []string
			for _, p := range pods {
				podNames = append(podNames, p.Name)
			}
			b, _ := json.Marshal(podNames)
			writer.Write(b)
		})

		http.HandleFunc("/setcache", func(writer http.ResponseWriter, request *http.Request) {
			podId := request.URL.Query().Get("id")
			if len(podId) == 0 {
				writer.Write([]byte("id为空"))
				return
			}

			getPod, exist := mykubelet.PodManager.GetPodByUID(apitypes.UID(podId))
			if !exist {
				writer.Write([]byte("pod没找到"))
				return
			}

			/*podStatus := &container.PodStatus{
				ID: apitypes.UID(podId),
				SandboxStatuses: []*v1alpha2.PodSandboxStatus{
					{
						Id:    podId,
						State: v1alpha2.PodSandboxState_SANDBOX_READY,
					},
				},
			}

			containerStatuses := make([]*container.Status, len(getPod.Spec.Containers))
			for i, c := range getPod.Spec.Containers {
				containerStatuses[i] = &container.Status{
					Name:      c.Name,
					Image:     c.Image,
					State:     container.ContainerStateRunning,
					CreatedAt: time.Now(),
					StartedAt: time.Now().Add(time.Second * 3),
				}
			}
			podStatus.ContainerStatuses = containerStatuses*/
			podStatus := mycore.SetPodReady(getPod)

			mykubelet.PodCache.Set(apitypes.UID(podId), podStatus, nil, mykubelet.Clock.Now())
			writer.Write([]byte("success"))
		})
		http.ListenAndServe(":8080", nil)
	}()

	//mykubelet.Pleg.Start()

	for item := range mykubelet.PodConfig.Updates() {
		switch item.Op {
		case types.ADD:
			mykubelet.HandlePodAdditions(item.Pods)
		case types.UPDATE, types.DELETE:
			mykubelet.HandlePodUpdates(item.Pods)
		case types.REMOVE:
			mykubelet.HandlePodRemoves(item.Pods)
		}
	}
}
