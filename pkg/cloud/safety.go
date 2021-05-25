package cloud

// SecurityGroup 安全组
type SecurityGroup interface {

	// Bind 绑定安全组
	Bind(obj interface{}) error

	// UnBind 解绑安全组
	UnBind(obj interface{}) error
}


// KeyParis 密钥
type KeyParis interface {

	// DescribeKeyPairs 查询密钥对
	DescribeKeyPairs(obj interface{}) (interface{}, error)

	// CreateKeyPair 创建密钥对
	CreateKeyPair(obj interface{}) (interface{}, error)

	// BindKeyPairs 绑定密钥对
	BindKeyPairs(obj interface{}) error

	// UnBindKeyPairs 解绑密钥对
	UnBindKeyPairs(obj interface{}) error

	// DeleteKeyPairs 删除密钥对
	DeleteKeyPairs(obj interface{}) error
}
