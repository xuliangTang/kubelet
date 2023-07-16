package myclientlib

import (
	"context"
	coordinationv1 "k8s.io/api/coordination/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"log"
	"time"
)

var lease *coordinationv1.Lease

const (
	leaseNS = "kube-node-lease"
)

// Renew 每隔10s续租
func Renew(client *clientset.Clientset) {
	// 获取lease
	getLease, err := client.CoordinationV1().Leases(leaseNS).Get(context.Background(), nodeName, metav1.GetOptions{})
	lease = getLease
	if err != nil {
		log.Fatalln(err)
	}

	// 模拟kubelet心跳传递给apiServer，每隔 NodeLeaseDurationSeconds*0.25s 续租一次（40*0.25=10s）
	// controller-manager会监控节点状态，检查租约是否过期，过期会更新node状态为unknown
	leaseDuration := time.Duration(40) * time.Second
	renewInterval := time.Duration(float64(leaseDuration) * 0.25)
	go func() {
		for {
			renewLease(client)
			time.Sleep(renewInterval)
		}
	}()
}

// 续租 更新节点lease的spec.renewTime
func renewLease(client *clientset.Clientset) {
	now := metav1.NewMicroTime(time.Now())
	lease.Spec.RenewTime = &now
	newLease, err := client.CoordinationV1().Leases(leaseNS).Update(context.Background(), lease, metav1.UpdateOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	lease = newLease
}
