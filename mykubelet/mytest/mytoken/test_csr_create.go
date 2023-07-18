package main

import (
	"context"
	"k8s.io/client-go/util/certificate/csr"
	"k8s.io/kubernetes/mykubelet/mylib"
	"k8s.io/kubernetes/mykubelet/mytest/mytoken/mytokenlib"
	"log"
	"os"
	"time"
)

// 创建csr
func main() {
	client := mytokenlib.InitDefaultClient()

	// 创建csr
	csrObj := mylib.CreateCsr(client)

	// 等待csr批准后获取证书并写入文件
	// mylib.WaitSaveCert(client)

	// 也可以使用kubeadm内置方法：pkg/kubelet/certificate/bootstrap/bootstrap.go的356行
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3600)
	defer cancel()
	certData, err := csr.WaitForCertificate(ctx, client, csrObj.Name, csrObj.UID)
	if err != nil {
		log.Fatalln(err)
	}
	// 证书写入文件
	err = os.WriteFile(mylib.TestCertFile, certData, 0655)
	if err != nil {
		log.Fatalln(err)
	}
}
