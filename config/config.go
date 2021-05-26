package config

import (
	"github.com/eadydb/k8s-aim/pkg/zlog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// Tencent 腾讯云配置
type Tencent struct {
	SecretId  string `yaml:"secret_id"`  // 用于标识 API 调用者身份
	SecretKey string `yaml:"secret_key"` // 用于加密签名字符串和服务器端验证签名字符串的密钥
}

// Kubernetes kubernetes 相关配置
type Kubernetes struct {
	NameSpace  string `yaml:"namespace"`   // 命名空间
	KubeConfig string `yaml:"kube_config"` // kubeConfig路径
	Token      string `yaml:"token"`       // kubernetes Token
}

// Config 配置文件
type Config struct {
	Manufacturers string      `yaml:"manufacturers"` // 云厂商
	Tencent       *Tencent    `yaml:"tencent"`       // 腾讯云配置
	Kubernetes    *Kubernetes `yaml:"kubernetes"`    // Kubernetes相关配置
}

// loadConfig 加载配置文件
func (c *Config) loadConfig() {
	yamlFile, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		zlog.Errorf("load config file failed, %v", err)
		os.Exit(0)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		zlog.Errorf("unmarshal config file failed, %v", err)
		os.Exit(0)
	}
}

// Init 初始化配置文件
func (c *Config) Init() *Config {
	c.loadConfig()
	return c
}
