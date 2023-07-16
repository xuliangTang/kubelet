package myclientlib

import (
	"context"
	"encoding/json"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
	"log"
)

const (
	nodeName = "mylain"
)

// SetNodeReady 手动修改节点的ready为true
func SetNodeReady(client *clientset.Clientset) {
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
	payloadJson, _ := json.Marshal(payload)

	_, err := client.CoreV1().Nodes().Patch(context.Background(), nodeName, types.JSONPatchType,
		payloadJson, v1.PatchOptions{}, "status")
	if err != nil {
		log.Fatalln(err)
	}
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
