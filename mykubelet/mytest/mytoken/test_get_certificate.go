package main

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/mykubelet/mylib"
	"k8s.io/kubernetes/mykubelet/mytest/mytoken/mytokenlib"
	"log"
	"os"
)

// 创建并批准csr资源后，获取证书
func main() {
	client := mytokenlib.InitDefaultClient()

	getCsr, err := client.CertificatesV1().CertificateSigningRequests().
		Get(context.Background(), "testcsr", metav1.GetOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	if getCsr.Status.Certificate != nil {
		// 证书保存到文件
		err = os.WriteFile(mylib.TestCertFile, getCsr.Status.Certificate, 0655)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
