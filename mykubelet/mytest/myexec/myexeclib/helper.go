package myexeclib

import (
	"context"
	"errors"
	"github.com/emicklei/go-restful"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/apimachinery/pkg/util/proxy"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/kubelet/cri/streaming"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type responder struct{}

func (r *responder) Error(w http.ResponseWriter, req *http.Request, err error) {
	klog.ErrorS(err, "Error while proxying request")
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func proxyStream(w http.ResponseWriter, r *http.Request, url *url.URL) {
	// TODO(random-liu): Set MaxBytesPerSec to throttle the stream.
	handler := proxy.NewUpgradeAwareHandler(url, nil /*transport*/, false /*wrapTransport*/, true /*upgradeRequired*/, &responder{})
	handler.ServeHTTP(w, r)
}

func GetExec(request *restful.Request, response *restful.Response) {
	//   pod.Name + "_" + pod.Namespace  == PodFullName
	url, err := GetUrl()
	if err != nil {
		streaming.WriteError(err, response.ResponseWriter)
		return
	}
	proxyStream(response.ResponseWriter, request.Request, url)
}

const RemoteRuntimeAddress = "110.188.24.175:8989" // 远程的runtime 地址
const RemoteRuntimeIp = "110.188.24.175"
const ContainerId = "95fb38d72a8e0336d458459de8b73d00ae8168971a3aa938b46ad471d9adbb39" // 容器id kubectl get pod -o yaml查看

func initRuntimeClient() runtimeapi.RuntimeServiceClient {
	gopts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(ctx, RemoteRuntimeAddress, gopts...)
	if err != nil {
		log.Fatalln(err)
	}
	return runtimeapi.NewRuntimeServiceClient(conn)
}

func runtimeExec(req *runtimeapi.ExecRequest) (*runtimeapi.ExecResponse, error) {
	runtimeClient := initRuntimeClient()
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()

	resp, err := runtimeClient.Exec(ctx, req)
	if err != nil {
		klog.ErrorS(err, "Exec cmd from runtime service failed", "containerID", req.ContainerId, "cmd", req.Cmd)
		return nil, err
	}
	klog.V(10).InfoS("[RemoteRuntimeService] Exec Response")

	if resp.Url == "" {
		errorMessage := "URL is not set"
		err := errors.New(errorMessage)
		klog.ErrorS(err, "Exec failed")
		return nil, err
	}

	return resp, nil

}
func GetUrl() (*url.URL, error) {
	req := &runtimeapi.ExecRequest{
		ContainerId: ContainerId,
		Cmd:         []string{"ls"},
		Tty:         false,
		Stdin:       true,
		Stdout:      true,
		Stderr:      true,
	}
	resp, err := runtimeExec(req)
	if err != nil {
		return nil, err
	}

	// 需要修改containerd配置：vi /etc/containerd/config.toml
	// stream_server_address改为0.0.0.0 (默认127.0.0.1 只能本地访问)
	// stream_server_port改为一个固定的端口如6595 (默认0 代表随机生成端口)

	// 这个地址 是 容器运行时生成的，每次都不一样，他会启动时监听一个地址用于给我们exec
	klog.Info("得到的URL是：", resp.Url)
	resp.Url = strings.Replace(resp.Url, "[::]", RemoteRuntimeIp, -1)
	klog.Info("修改过后的URL是：", resp.Url)
	return url.Parse(resp.Url)
}
