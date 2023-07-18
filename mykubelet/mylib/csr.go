package mylib

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	cryptorand "crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"k8s.io/client-go/util/cert"
	"log"
	"os"
)

const (
	TestPrivateKeyFile = "./mykubelet/certs/kubelet.key" // 测试的私钥key文件
	TestCertFile       = "./mykubelet/certs/kubelet.pem" // 测试的csr批准后的证书文件
)

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
			CommonName:   fmt.Sprintf("system:node:%s", "txl"),
			Organization: []string{"system:nodes"},
		},
	}
	csrPEM, err := cert.MakeCSRFromTemplate(privateKey, cr)
	if err != nil {
		log.Fatalln(err)
	}

	return csrPEM
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
