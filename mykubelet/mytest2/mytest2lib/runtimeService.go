package mytest2lib

import (
	cri "k8s.io/cri-api/pkg/apis"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"k8s.io/kubernetes/mykubelet/mylib"
	"time"
)

type MyRuntimeService struct {
}

func NewMyRuntimeService() *MyRuntimeService {
	return &MyRuntimeService{}
}

func (m MyRuntimeService) Version(apiVersion string) (*runtimeapi.VersionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m MyRuntimeService) CreateContainer(podSandboxID string, config *runtimeapi.ContainerConfig, sandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m MyRuntimeService) StartContainer(containerID string) error {
	//TODO implement me
	panic("implement me")
}

func (m MyRuntimeService) StopContainer(containerID string, timeout int64) error {
	//TODO implement me
	panic("implement me")
}

func (m MyRuntimeService) RemoveContainer(containerID string) error {
	//TODO implement me
	panic("implement me")
}

func (m MyRuntimeService) ListContainers(filter *runtimeapi.ContainerFilter) ([]*runtimeapi.Container, error) {
	//TODO implement me
	return mylib.MockContainers(), nil
}

func (m MyRuntimeService) ContainerStatus(containerID string) (*runtimeapi.ContainerStatus, error) {
	//TODO implement me
	panic("implement me")
}

func (m MyRuntimeService) UpdateContainerResources(containerID string, resources *runtimeapi.LinuxContainerResources) error {
	//TODO implement me
	panic("implement me")
}

func (m MyRuntimeService) ExecSync(containerID string, cmd []string, timeout time.Duration) (stdout []byte, stderr []byte, err error) {
	//TODO implement me
	panic("implement me")
}

func (m MyRuntimeService) Exec(request *runtimeapi.ExecRequest) (*runtimeapi.ExecResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m MyRuntimeService) Attach(req *runtimeapi.AttachRequest) (*runtimeapi.AttachResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m MyRuntimeService) ReopenContainerLog(ContainerID string) error {
	//TODO implement me
	panic("implement me")
}

func (m MyRuntimeService) RunPodSandbox(config *runtimeapi.PodSandboxConfig, runtimeHandler string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m MyRuntimeService) StopPodSandbox(podSandboxID string) error {
	//TODO implement me
	panic("implement me")
}

func (m MyRuntimeService) RemovePodSandbox(podSandboxID string) error {
	//TODO implement me
	panic("implement me")
}

func (m MyRuntimeService) PodSandboxStatus(podSandboxID string) (*runtimeapi.PodSandboxStatus, error) {
	//TODO implement me
	return &runtimeapi.PodSandboxStatus{}, nil
}

func (m MyRuntimeService) ListPodSandbox(filter *runtimeapi.PodSandboxFilter) ([]*runtimeapi.PodSandbox, error) {
	//TODO implement me
	return mylib.MockSandbox(), nil
}

func (m MyRuntimeService) PortForward(request *runtimeapi.PortForwardRequest) (*runtimeapi.PortForwardResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m MyRuntimeService) ContainerStats(containerID string) (*runtimeapi.ContainerStats, error) {
	//TODO implement me
	panic("implement me")
}

func (m MyRuntimeService) ListContainerStats(filter *runtimeapi.ContainerStatsFilter) ([]*runtimeapi.ContainerStats, error) {
	//TODO implement me
	panic("implement me")
}

func (m MyRuntimeService) UpdateRuntimeConfig(runtimeConfig *runtimeapi.RuntimeConfig) error {
	//TODO implement me
	panic("implement me")
}

func (m MyRuntimeService) Status() (*runtimeapi.RuntimeStatus, error) {
	//TODO implement me
	panic("implement me")
}

var _ cri.RuntimeService = &MyRuntimeService{}
