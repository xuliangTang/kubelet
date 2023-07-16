package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	bootstraputil "k8s.io/cluster-bootstrap/token/util"
	bootstraptokenv1 "k8s.io/kubernetes/cmd/kubeadm/app/apis/bootstraptoken/v1"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/options"
	"k8s.io/kubernetes/mykubelet/mytest/mytoken/mytokenlib"
	"log"
)

// 模拟kubeadm token create 手工创建token
// https://kubernetes.io/zh-cn/docs/reference/access-authn-authz/bootstrap-tokens/#bootstrap-token-secret-format
func main() {
	// 使用内置包生成一个token
	token, _ := bootstraputil.GenerateBootstrapToken()
	// 生成BootstrapTokenString对象
	bts, _ := bootstraptokenv1.NewBootstrapTokenString(token)

	// 生成带有默认值的bootstrapToken的配置
	opt := options.NewBootstrapTokenOptions()
	// 设置token值
	opt.Token = bts

	// 创建secret
	bootSecret := bootstraptokenv1.BootstrapTokenToSecret(opt.BootstrapToken)
	// 使用clientset创建
	client := mytokenlib.InitDefaultClient()
	newSecret, err := client.CoreV1().Secrets(metav1.NamespaceSystem).Create(context.Background(), bootSecret, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(newSecret.Name)
}
