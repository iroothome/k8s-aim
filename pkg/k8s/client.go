package k8s

import (
	"context"
	"flag"
	"github.com/eadydb/k8s-aim/pkg/zlog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

// KClient Kubernetes Cluster client
type KClient struct {
	Config     *rest.Config
	ClientSet  *kubernetes.Clientset
	Token      string
	KubeConfig string
	NameSpace  string
	Ctx        context.Context
}

// Init 初始化Kubernetes Cluster 客户端
func (c *KClient) Init() {
	var err error
	var config *rest.Config
	var kubeConfig *string

	if c.KubeConfig == "" {
		if home := homeDir(); home != "" {
			kubeConfig = flag.String("kubeConfig", filepath.Join(home, ".kube", "config"), "(可选)kubeConfig 文件的绝对路径")
		} else {
			kubeConfig = flag.String("kubeConfig", "", "kubeConfig 文件的绝对路径")
		}
	} else {
		kubeConfig = &c.KubeConfig
	}

	// 首先使用 inCluster 模式(需要区配置对应的RBAC 权限,默认的sa是default-->是没有获取deployment的List权限)
	if config, err = rest.InClusterConfig(); err != nil {
		// 使用KubeConfig文件配置集群配置Config对象
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeConfig); err != nil {
			zlog.Panicf("Load kubernetes cluster config failed, %s", err.Error())
		}
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		zlog.Panicf("Init Kubernetes cluster client set failed, %s", err.Error())
	}

	c.Config = config
	c.ClientSet = clientSet
	c.Ctx = context.Background()
}

// homeDir 当前Home目录
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}
