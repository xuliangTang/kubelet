package main

import (
	"context"
	"encoding/json"
	"fmt"
	coordinationv1 "k8s.io/api/coordination/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"time"
)

var lease *coordinationv1.Lease

const (
	leaseNS  = "kube-node-lease"
	nodeName = "mylain"
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

	// 获取lease
	lease, err = client.CoordinationV1().Leases(leaseNS).Get(context.Background(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	// 模拟kubelet心跳传递给apiServer，每隔 NodeLeaseDurationSeconds*0.25s 续租一次（40*0.25=10s）
	// controller-manager会监控节点状态，检查租约是否过期，过期会更新node状态为unknown
	go func() {
		for {
			if err := renew(client); err != nil {
				log.Fatalln(err)
			}

			leaseDuration := time.Duration(40) * time.Second
			renewInterval := time.Duration(float64(leaseDuration) * 0.25)
			time.Sleep(time.Second * renewInterval)
		}
	}()

	// patch手动修改node节点状态为true
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
	node, err := client.CoreV1().Nodes().Patch(context.Background(), nodeName, types.JSONPatchType,
		payloadByte, metav1.PatchOptions{}, "status")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(node.Status.Conditions[3])

	select {}
}

// 续租 更新节点lease的spec.renewTime
func renew(client *clientset.Clientset) error {
	now := metav1.NewMicroTime(time.Now())
	lease.Spec.RenewTime = &now
	newLease, err := client.CoordinationV1().Leases(leaseNS).Update(context.Background(), lease, metav1.UpdateOptions{})
	lease = newLease
	return err
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
