package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"log"
	"time"
)

const criAddr = "110.188.24.175:8989"

// 连接远程containerd
// 需要开启远程TCP连接配置, 修改/etc/containerd/config.toml
// [grpc]
// tcp_address = "0.0.0.0:8989"
//
// [plugins."io.containerd.grpc.v1.cri"]
// disable_tcp_service = false
func main() {
	client := initGrpcClient()
	runtimeSvc := v1alpha2.NewRuntimeServiceClient(client)

	rsp, err := runtimeSvc.Version(context.Background(), &v1alpha2.VersionRequest{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(rsp.Version)
}

func initGrpcClient() *grpc.ClientConn {
	gopts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(ctx, criAddr, gopts...)
	if err != nil {
		log.Fatalln(err)
	}

	return conn
}
