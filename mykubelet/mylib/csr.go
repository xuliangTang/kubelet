package mylib

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	cryptorand "crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	certificatesv1 "k8s.io/api/certificates/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/cert"
	"k8s.io/utils/pointer"
	"log"
	"os"
	"time"
)

const (
	TestPrivateKeyFile = "./mykubelet/certs/kubelet.key" // 测试的私钥key文件
	TestCertFile       = "./mykubelet/certs/kubelet.pem" // 测试的csr批准后的证书文件
)

// CreateCsr 创建csr资源
func CreateCsr(client *clientset.Clientset) *certificatesv1.CertificateSigningRequest {
	csr := &certificatesv1.CertificateSigningRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testcsr",
		},
		Spec: certificatesv1.CertificateSigningRequestSpec{
			Request: GenCSRPEM(),
			Usages: []certificatesv1.KeyUsage{
				certificatesv1.UsageClientAuth,
			},
			ExpirationSeconds: pointer.Int32(int32(time.Second * 3600 / time.Second)),
			SignerName:        certificatesv1.KubeAPIServerClientSignerName,
		},
	}

	newCsr, err := client.CertificatesV1().CertificateSigningRequests().Create(context.Background(), csr, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	return newCsr
}

// GenCSRPEM 生成csr(证书签名请求)的spec.request
func GenCSRPEM() []byte {
	// 生成客户端私钥
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), cryptorand.Reader)
	if err != nil {
		log.Fatalln(err)
	}

	// 保存私钥文件
	if err = savePrivateKeyToFile(privateKey); err != nil {
		log.Fatalln(err)
	}

	// 生成csr
	cr := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   fmt.Sprintf("system:node:%s", "mylain"),
			Organization: []string{"system:nodes"},
		},
	}
	csrPEM, err := cert.MakeCSRFromTemplate(privateKey, cr)
	if err != nil {
		log.Fatalln(err)
	}

	return csrPEM
}

// WaitSaveCert 等待csr批准，获取证书
func WaitSaveCert(client *clientset.Clientset) {
	// 设置超时时间
	stopCh := make(chan struct{})
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3600)
	go func() {
		select {
		case <-ctx.Done():
			log.Println("timeout")
			stopCh <- struct{}{}
		}
	}()

	// 监听csr
	fact := informers.NewSharedInformerFactory(client, 0)
	fact.Certificates().V1().CertificateSigningRequests().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {
			if obj, ok := newObj.(*certificatesv1.CertificateSigningRequest); ok {
				if obj.Name == "testcsr" && obj.Status.Certificate != nil {
					// 获取批准后的证书
					err := os.WriteFile(TestCertFile, obj.Status.Certificate, 0655)
					if err != nil {
						log.Fatalln(err)
					}
					stopCh <- struct{}{}
				}
			}
		},
	})
	fact.Start(stopCh)

	<-stopCh
}

// 保存私钥文件
func savePrivateKeyToFile(key *ecdsa.PrivateKey) error {
	b, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return err
	}
	privatePem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: b,
		},
	)

	_ = os.Remove(TestPrivateKeyFile)
	err = os.WriteFile(TestPrivateKeyFile, privatePem, 0600)
	if err != nil {
		return err
	}

	return nil
}
