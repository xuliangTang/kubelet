go run cmd/kubelet/kubelet.go --v=1 \
  --hostname-override=mylain \
  --kubeconfig=./mykubelet/kubelet.config \
  --bootstrap-kubeconfig=./mykubelet/bootstrap.yaml \
  --config=./mykubelet/kubelet.config.yaml \
  --cert-dir=/Users/zhoujiahong/kubelet/certs/kubelet
