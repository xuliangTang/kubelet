package myclientlib

import (
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func InitClient() *clientset.Clientset {
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", "./mykubelet/kubelet.config")
	if err != nil {
		log.Fatalln(err)
	}

	client, err := clientset.NewForConfig(kubeconfig)
	if err != nil {
		log.Fatalln(err)
	}

	return client
}
