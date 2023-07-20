package mylib

import (
	cri "k8s.io/cri-api/pkg/apis"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"time"
)

type MyRuntimeService struct{}

func (m MyRuntimeService) Version(apiVersion string) (*runtimeapi.VersionResponse, error) {
	return &runtimeapi.VersionResponse{
		Version:     "0.1.0", // 只支持0.1.0
		RuntimeName: "lain",
	}, nil
}

func (m MyRuntimeService) CreateContainer(podSandboxID string, config *runtimeapi.ContainerConfig, sandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {
	return "", nil
}

func (m MyRuntimeService) StartContainer(containerID string) error {
	return nil
}

func (m MyRuntimeService) StopContainer(containerID string, timeout int64) error {
	return nil
}

func (m MyRuntimeService) RemoveContainer(containerID string) error {
	return nil
}

func (m MyRuntimeService) ListContainers(filter *runtimeapi.ContainerFilter) ([]*runtimeapi.Container, error) {
	// return []*runtimeapi.Container{}, nil
	return MockContainers(), nil
}

func (m MyRuntimeService) ContainerStatus(containerID string) (*runtimeapi.ContainerStatus, error) {
	return &runtimeapi.ContainerStatus{}, nil
}

func (m MyRuntimeService) UpdateContainerResources(containerID string, resources *runtimeapi.LinuxContainerResources) error {
	return nil
}

func (m MyRuntimeService) ExecSync(containerID string, cmd []string, timeout time.Duration) (stdout []byte, stderr []byte, err error) {
	return []byte{}, []byte{}, nil
}

func (m MyRuntimeService) Exec(request *runtimeapi.ExecRequest) (*runtimeapi.ExecResponse, error) {
	return &runtimeapi.ExecResponse{}, nil
}

func (m MyRuntimeService) Attach(req *runtimeapi.AttachRequest) (*runtimeapi.AttachResponse, error) {
	return &runtimeapi.AttachResponse{}, nil
}

func (m MyRuntimeService) ReopenContainerLog(ContainerID string) error {
	return nil
}

func (m MyRuntimeService) RunPodSandbox(config *runtimeapi.PodSandboxConfig, runtimeHandler string) (string, error) {
	return "", nil
}

func (m MyRuntimeService) StopPodSandbox(podSandboxID string) error {
	return nil
}

func (m MyRuntimeService) RemovePodSandbox(podSandboxID string) error {
	return nil
}

func (m MyRuntimeService) PodSandboxStatus(podSandboxID string) (*runtimeapi.PodSandboxStatus, error) {
	return &runtimeapi.PodSandboxStatus{}, nil
}

func (m MyRuntimeService) ListPodSandbox(filter *runtimeapi.PodSandboxFilter) ([]*runtimeapi.PodSandbox, error) {
	// return []*runtimeapi.PodSandbox{}, nil
	return MockSandbox(), nil
}

func (m MyRuntimeService) PortForward(request *runtimeapi.PortForwardRequest) (*runtimeapi.PortForwardResponse, error) {
	return &runtimeapi.PortForwardResponse{}, nil
}

func (m MyRuntimeService) ContainerStats(containerID string) (*runtimeapi.ContainerStats, error) {
	return &runtimeapi.ContainerStats{}, nil
}

func (m MyRuntimeService) ListContainerStats(filter *runtimeapi.ContainerStatsFilter) ([]*runtimeapi.ContainerStats, error) {
	return []*runtimeapi.ContainerStats{}, nil
}

func (m MyRuntimeService) UpdateRuntimeConfig(runtimeConfig *runtimeapi.RuntimeConfig) error {
	return nil
}

func (m MyRuntimeService) Status() (*runtimeapi.RuntimeStatus, error) {
	return &runtimeapi.RuntimeStatus{
		Conditions: []*runtimeapi.RuntimeCondition{ //必须有下面2个配置，才会认为node是ready状态
			{
				Type:   "RuntimeReady",
				Status: true,
			},
			{
				Type:   "NetworkReady",
				Status: true,
			},
		},
	}, nil
}

var _ cri.RuntimeService = &MyRuntimeService{}
