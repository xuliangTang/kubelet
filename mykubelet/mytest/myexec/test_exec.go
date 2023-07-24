package main

import (
	"github.com/emicklei/go-restful"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/mykubelet/mytest/myexec/myexeclib"
	"net/http"
)

func main() {
	container := restful.NewContainer()
	ws := new(restful.WebService)
	ws.Path("/exec")
	{
		ws.Route(ws.GET("/{podNamespace}/{podID}/{containerName}").
			To(myexeclib.GetExec).
			Operation("getExec"))
		ws.Route(ws.POST("/{podNamespace}/{podID}/{containerName}").
			To(myexeclib.GetExec).
			Operation("getExec"))
		ws.Route(ws.GET("/{podNamespace}/{podID}/{uid}/{containerName}").
			To(myexeclib.GetExec).
			Operation("getExec"))
		ws.Route(ws.POST("/{podNamespace}/{podID}/{uid}/{containerName}").
			To(myexeclib.GetExec).
			Operation("getExec"))
	}
	container.Add(ws)

	klog.Info("启动http服务，监听9090端口")
	http.ListenAndServe(":9090", container)
}
