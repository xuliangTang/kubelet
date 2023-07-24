package main

import (
	"crypto/tls"
	"k8s.io/apimachinery/pkg/util/httpstream/spdy"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"net/http"
	"os"
)

func main() {
	req, _ := http.NewRequest("GET", "http://localhost:9090/exec/default/abd/bbb", nil)
	req.Header.Set("Upgrade", "SPDY/3.1")
	req.Header.Set("Connection", "Upgrade")
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	rt := spdy.NewRoundTripper(tlsConfig, true, false)

	executor, err := remotecommand.NewSPDYExecutorForTransports(rt, rt, http.MethodGet, req.URL)
	if err != nil {
		log.Fatal(err)
	}
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		log.Fatalln(err)
	}
}
