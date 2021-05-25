package profile

type HttpProfile struct {
	ReqMethod  string // 请求方式 POST/GET
	ReqTimeout int    // 请求超时时间
	Scheme     string // HTTP/HTTPS
	RootDomain string // 顶级域名
	Endpoint   string // 节点
	Protocol   string // 协议
}

func NewHttpProfile() *HttpProfile {
	return &HttpProfile{
		ReqMethod:  "POST",
		ReqTimeout: 60,
		Scheme:     "HTTPS",
		RootDomain: "",
		Endpoint:   "",
	}
}
