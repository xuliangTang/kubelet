package main

import (
	"context"
	"encoding/json"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func main() {
	// 初始化kubeclient
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", "./mykubelet/kubelet.config")
	if err != nil {
		log.Fatalln(err)
	}
	client, err := clientset.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatalln(err)
	}

	// 获取node列表
	nodes, err := client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	for _, node := range nodes.Items {
		fmt.Println(node.Name)
	}

	// patch更新node节点状态为true
	payload := []jsonPatch{
		{
			Op:   "replace",
			Path: "/status/conditions/3",
			Value: &value{
				Type:   "Ready",
				Status: "True",
			},
		},
	}
	payloadByte, _ := json.Marshal(payload)
	node, err := client.CoreV1().Nodes().Patch(context.Background(), "mylain", types.JSONPatchType,
		payloadByte, metav1.PatchOptions{}, "status")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(node.Status.Conditions[3])
}

type jsonPatch struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value *value `json:"value,omitempty"`
}

type value struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}
