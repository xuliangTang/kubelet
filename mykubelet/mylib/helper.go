package mylib

import v1 "k8s.io/api/core/v1"

// SetNodeLabels 设置node标签，用于 pkg/kubelet/kubelet_node_status.go 的 311 行
func SetNodeLabels(node *v1.Node) {
	node.Labels["beta.kubernetes.io/os"] = "lain_os"
	node.Labels["kubernetes.io/hostname"] = "lain"
	node.Labels["type"] = "agent"
}
