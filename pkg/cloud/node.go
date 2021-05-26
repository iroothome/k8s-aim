package cloud

// Manufacturers 云厂家枚举类型
type Manufacturers string

const (
	AliYun    Manufacturers = "AliYun"    // 阿里云厂家
	Tencent   Manufacturers = "Tencent"   // 腾讯云厂家
	QingCloud Manufacturers = "QingCloud" // 青云
)

// ClusterNode kubernetes cluster node
type ClusterNode struct {
	Name     string   // kubernetes cluster worker node name
	HostName string   // ecs hostname
	Ip       string   // ecs ip address
	Tag      []string // kubernetes node tags
}

// Node kubernetes cluster node
type Node interface {

	// CreateClusterNode 创建k8s集群Node节点
	CreateClusterNode(node ClusterNode) (bool, error)

	// JoinCluster 加入k8s集群
	JoinCluster(node ClusterNode) (bool, error)

	// Monitor 初始化监控
	Monitor(node ClusterNode) error
}
