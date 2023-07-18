package main

import (
	"k8s.io/kubernetes/mykubelet/mylib"
	"k8s.io/kubernetes/mykubelet/mytest/mytoken/mytokenlib"
)

// 创建csr
func main() {
	client := mytokenlib.InitDefaultClient()

	// 创建csr
	mylib.CreateCsr(client)

	// 等待csr批准后获取证书
	mylib.WaitSaveCert(client)
}
