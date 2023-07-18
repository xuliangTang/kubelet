package mylib

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	cryptorand "crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"k8s.io/client-go/util/cert"
	"log"
)

// GenCSRPEM 生成csr(证书签名请求)的spec.request
func GenCSRPEM() []byte {
	// 生成客户端私钥
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), cryptorand.Reader)
	if err != nil {
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
