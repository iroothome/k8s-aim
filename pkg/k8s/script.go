package k8s

import (
	"bytes"
	"github.com/eadydb/k8s-aim/pkg/zlog"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const (
	k8sClusterAddress = "K8S_CLUSTER_ADDRESS"
	k8sToken          = "K8S_TOKEN"
	k8sCert           = "K8S_CERT"
)

// NodeInfo 云实例信息
type NodeInfo struct {
	Region      string // 实例
	Ip          string // ip
	PrivateKey  string // 私钥
	PublicKey   string // 公钥
	K8sNodeName string // Kubernetes Node name
}

// ClusterInfo Kubernetes cluster info
type ClusterInfo struct {
	ClusterAddress string // kubernetes cluster address
	Token          string // token
	CertHash       string // discovery token ca cert hash
}

// NewNodeInfo 实例化
func NewNodeInfo(region, ip, privateKey, publicKey string) *NodeInfo {
	return &NodeInfo{
		Region:     region,
		Ip:         ip,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}

// InstanceClusterScript 安装k8s准备包
// kube-proxy 、 kubelet 等
func (n *NodeInfo) InstanceClusterScript() bool {

	return true
}

// JoinClusterScript 加入kubernetes集群脚本
func (n *NodeInfo) JoinClusterScript(info *ClusterInfo) bool {
	script := n.readScript("/script/k8s/join_k8s.sh")
	if script == "" {
		return false
	}
	script = strings.ReplaceAll(script, k8sClusterAddress, info.ClusterAddress)
	script = strings.ReplaceAll(script, k8sToken, info.Token)
	script = strings.ReplaceAll(script, k8sCert, info.CertHash)

	cmd := exec.Command("bash", "-c", script)
	n.cmdOutPut(cmd)

	return true
}

// RemoveClusterScript 移除节点
func (n *NodeInfo) RemoveClusterScript() bool {

	return true
}

// cmdOutPut 脚本执行结果输出
func (n *NodeInfo) cmdOutPut(cmd *exec.Cmd) {
	var stdoutBuf, stderrBuf bytes.Buffer

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Start()
	if err != nil {
		zlog.Errorf("cmd.Start() failed with '%s'", err)
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	err = cmd.Wait()
	if err != nil {
		zlog.Errorf("cmd.Run() failed with %s", err)
	}

	if errStdout != nil || errStderr != nil {
		zlog.Errorf("failed to capture stdout or stderr")
	}

	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	zlog.Infof("out : %s", outStr)
	zlog.Errorf("error: %s", errStr)

}

// readScript 读取脚本文件
func (n *NodeInfo) readScript(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		zlog.Errorf("load kubernetes script file failed, %s", err)
		return ""
	}
	return string(content)
}
