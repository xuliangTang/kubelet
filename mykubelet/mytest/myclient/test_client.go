package main

import (
	"context"
	"fmt"
	coordinationv1 "k8s.io/api/coordination/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	myclientlib2 "k8s.io/kubernetes/mykubelet/mytest/myclient/myclientlib"
	"log"
)

var lease *coordinationv1.Lease

const (
	leaseNS  = "kube-node-lease"
	nodeName = "mylain"
)

func main() {
	// 初始化kubeclient
	client := myclientlib2.InitClient()

	// 获取node列表
	nodes, err := client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	for _, node := range nodes.Items {
		fmt.Println(node.Name)
	}

	// 定期续租节点lease
	myclientlib2.Renew(client)

	// patch手动修改node节点状态为true
	myclientlib2.SetNodeReady(client)

	select {}
}
