package cloud

import (
	"github.com/eadydb/k8s-aim/config"
	"github.com/eadydb/k8s-aim/pkg/cloud"
	"github.com/eadydb/k8s-aim/pkg/k8s"
)

// NodeServer New Node Server
type NodeServer struct {
	config.Config // Configuration
	k8s.KClient   // kubernetes cluster client
}

func (c *NodeServer) CreateClusterNode(node cloud.ClusterNode) (bool, error) {
	return true, nil
}

func (c *NodeServer) JoinCluster(node cloud.ClusterNode) (bool, error) {
	return true, nil
}

func (c *NodeServer) Monitor(node cloud.ClusterNode) error {
	return nil
}
