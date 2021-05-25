package cloud

// Instance 云厂商实例
type Instance interface {
	NetWork // 网络

	// CreateInstance 创建实例
	CreateInstance(obj interface{}) (interface{}, error)

	// StartInstance 启动实例
	StartInstance(obj interface{}) error

	// StopInstance 停止实例
	StopInstance(obj interface{}) error

	// RestartInstance 重启实例
	RestartInstance(obj interface{}) error
}

// Image 云厂商镜像
type Image interface {

	// GetImage 加载镜像
	GetImage(obj interface{}) (interface{}, error)
}
