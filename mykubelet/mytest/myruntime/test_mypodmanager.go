package main

import (
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/mykubelet/mytest/myclient/myclientlib"
	"k8s.io/kubernetes/mykubelet/mytest/myruntime/mycore"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
	"net/http"
	"sort"
	"time"
)

func main() {
	client := myclientlib.InitClient()
	nodeName := "mylain"
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

		// 设置podCache
		http.HandleFunc("/setcache", func(writer http.ResponseWriter, request *http.Request) {
			podId := request.URL.Query().Get("id")
			if podId == "" {
				writer.Write([]byte("缺少ID"))
				return
			}
			// 随便构建一个status，让managerPodLoop()接收到
			getPod, exist := podCache.PodManager.GetPodByUID(types.UID(podId))
			if !exist {
				writer.Write([]byte("没有找到Pod"))
				return
			}
			podStatus := mycore.SetPodReady(getPod)
			podCache.InnerPodCache.Set(types.UID(podId), podStatus, nil, time.Now())
			writer.Write([]byte("设置成功"))
		})
		http.ListenAndServe(":8080", nil)
	}()

	// 启动pleg 定期和本地cri交互 比对pod状态把变更事件放入podCache
	//fmt.Println("启动pleg")
	//mycore.StartPleg(podCache.Clock, podCache.InnerPodCache)

	// podConfig和apiserver交互 遍历podUpdate channel
	fmt.Println("开始监听")
	for item := range podCache.PodConfig.Updates() {
		pods := item.Pods
		switch item.Op {
		case kubetypes.ADD:
			mycore.HandlePodAdditions(podCache, pods)
		case kubetypes.UPDATE:
			mycore.HandlePodUpdates(podCache, pods)
		case kubetypes.DELETE:
			// 默认情况下不会直接删除，而是设置DeletionTimestamp和DeletionGracePeriodSeconds(优雅删除)
			// 所以kubelet得到的首先是更新事件，并把DeletionGracePeriodSeconds作为停止容器的timeout(还要取出preStop时间)
			// 接下来kubelet会处理一系列的操作(如cri)，接着触发一次Delete操作，最后触发Remove操作(一些清理工作)
			mycore.HandlePodUpdates(podCache, pods)
		case kubetypes.REMOVE:
			mycore.HandlePodRemoves(podCache, pods)
		}
	}
}
