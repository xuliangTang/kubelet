package main

import (
	"encoding/json"
	"fmt"
	"k8s.io/kubernetes/mykubelet/mytest/myclient/myclientlib"
	"k8s.io/kubernetes/mykubelet/mytest/myruntime/mycore"
	"k8s.io/kubernetes/pkg/kubelet/types"
	"net/http"
	"sort"
)

func main() {
	client := myclientlib.InitClient()
	nodeName := "lain1"
	podCache := mycore.NewPodCache(client, nodeName)

	go func() {
		fmt.Println("启动http服务")
		// 获取当前podManager缓存中等所有pod列表
		http.HandleFunc("/pods", func(writer http.ResponseWriter, request *http.Request) {
			var pods []string
			for _, pod := range podCache.PodManager.GetPods() {
				pods = append(pods, fmt.Sprintf("%s/%s", pod.Namespace, pod.Name))
			}
			sort.Strings(pods)
			b, _ := json.Marshal(pods)

			writer.Header().Set("Content-Type", "application/json")
			writer.Write(b)
		})
		http.ListenAndServe(":8080", nil)
	}()

	fmt.Println("开始监听")
	for item := range podCache.PodConfig.Updates() {
		pods := item.Pods
		switch item.Op {
		case types.ADD:
			for _, p := range pods {
				podCache.PodManager.AddPod(p)
			}
		case types.UPDATE:
			for _, p := range pods {
				podCache.PodManager.UpdatePod(p)
			}
		case types.DELETE:
			for _, p := range pods {
				podCache.PodManager.DeletePod(p)
			}
		}
	}
}
