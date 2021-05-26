package main

import (
	"github.com/eadydb/k8s-aim/config"
	"github.com/eadydb/k8s-aim/pkg/k8s"
	"github.com/eadydb/k8s-aim/pkg/zlog"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {


	// init configuration
	c := &config.Config{}
	c.Init()

	// init kubernetes cluster
	kClient := &k8s.KClient{
		NameSpace:  c.Kubernetes.NameSpace,
		KubeConfig: c.Kubernetes.KubeConfig,
		Token:      c.Kubernetes.Token,
	}
	kClient.Init()


	// 测试kubernetes集群
	deployment, err := kClient.ClientSet.AppsV1().Deployments(c.Kubernetes.KubeConfig).List(kClient.Ctx, metav1.ListOptions{})

	if err != nil || deployment == nil {
		zlog.Errorf("get kubernetes cluster deployment list failed, %s", err)
	}

	for _, deploy := range deployment.Items {
		zlog.Debugw(deploy.Name, zap.String("deployment", deploy.Name))
	}
}
