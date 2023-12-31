package main

import (
	"k8s.io/apimachinery/pkg/util/proxy"
	genericfeatures "k8s.io/apiserver/pkg/features"
	"k8s.io/apiserver/pkg/registry/rest"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/kubernetes/pkg/capabilities"
	"log"
	"net/http"
	"net/url"
)

func main() {
	urlStr := "http://localhost:9090/exec/default/mypod/mycontainer"
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	proxyHandler := proxy.NewUpgradeAwareHandler(urlObj, http.DefaultTransport, false, true, proxy.NewErrorResponder(nil))

	log.Println("启动假的apiserver")
	http.ListenAndServe(":6443", proxyHandler)
}

func newThrottledUpgradeAwareProxyHandler(location *url.URL, transport http.RoundTripper, wrapTransport, upgradeRequired, interceptRedirects bool, responder rest.Responder) *proxy.UpgradeAwareHandler {
	handler := proxy.NewUpgradeAwareHandler(location, transport, wrapTransport, upgradeRequired, proxy.NewErrorResponder(responder))
	handler.InterceptRedirects = interceptRedirects && utilfeature.DefaultFeatureGate.Enabled(genericfeatures.StreamingProxyRedirects)
	handler.RequireSameHostRedirects = utilfeature.DefaultFeatureGate.Enabled(genericfeatures.ValidateProxyRedirects)
	handler.MaxBytesPerSec = capabilities.Get().PerConnectionBandwidthLimitBytesPerSec
	return handler
}
