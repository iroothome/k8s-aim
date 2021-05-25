package cloud

// Manufacturers 云厂家枚举类型
type Manufacturers string

const (
	AliYun    Manufacturers = "AliYun"    // 阿里云厂家
	Tencent   Manufacturers = "Tencent"   // 腾讯云厂家
	QingCloud Manufacturers = "QingCloud" // 青云
)

// Cloud 云服务器厂家基础接口
type Cloud interface {
	Image // 镜像

	Instance // 实例

	KeyParis // 密钥对

	// CreateClusterNode 创建k8s集群Node节点
	CreateClusterNode(obj interface{}) error

	// JoinCluster 加入k8s集群
	JoinCluster(obj interface{}) error

	// Monitor 初始化监控
	Monitor(obj interface{}) error
}
