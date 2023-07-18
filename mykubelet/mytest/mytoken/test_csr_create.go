package main

import (
	"context"
	certificatesv1 "k8s.io/api/certificates/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/mykubelet/mylib"
	"k8s.io/kubernetes/mykubelet/mytest/mytoken/mytokenlib"
	"k8s.io/utils/pointer"
	"log"
	"time"
)

// 创建csr
func main() {
	client := mytokenlib.InitDefaultClient()

	csr := &certificatesv1.CertificateSigningRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testcsr",
		},
		Spec: certificatesv1.CertificateSigningRequestSpec{
			Request: mylib.GenCSRPEM(),
			Usages: []certificatesv1.KeyUsage{
				certificatesv1.UsageClientAuth,
			},
			ExpirationSeconds: pointer.Int32(int32(time.Second * 3600 / time.Second)),
			SignerName:        certificatesv1.KubeAPIServerClientSignerName,
		},
	}

	_, err := client.CertificatesV1().CertificateSigningRequests().Create(context.Background(), csr, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln(err)
	}
}
