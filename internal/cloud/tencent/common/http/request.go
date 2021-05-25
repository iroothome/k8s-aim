package http

import (
	"github.com/eadydb/k8s-aim/pkg/zlog"
	"io"
	"math/rand"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	POST = "POST"
	GET  = "GET"

	HTTP  = "http"
	HTTPS = "https"

	RootDomain = "tencentcloudapi.com"
	Path       = "/"
)

type Request interface {
	GetAction() string              // 获取接口名称
	GetBodyReader() io.Reader       // 请求内容
	GetScheme() string              // 请求协议 HTTP/HTTPS
	GetRootDomain() string          // 顶级域名
	GetServiceDomain(string) string // 服务域名
	GetDomain() string              // 域名
	GetHttpMethod() string          // Http请求方法类型
	GetParams() map[string]string   // 请求参数
	GetPath() string                // 请求路径
	GetService() string             // 服务
	GetUrl() string                 // 请求url地址
	GetVersion() string             // 版本号
	SetScheme(string)               // 设置Scheme
	SetRootDomain(string)           // 设置顶级域名
	SetDomain(string)               // 设置域名
	SetHttpMethod(string)           // 设置Http请求方法类型
}

type BaseRequest struct {
	httpMethod string            // Http方法类型
	scheme     string            // 请求协议 HTTP/HTTPS
	rootDomain string            // 顶级域名
	domain     string            // 域名
	path       string            // 请求路径
	params     map[string]string // 参数
	formParams map[string]string // 表单参数

	service string // 服务
	version string // 版本号
	action  string // 接口名称
}

// GetAction 获取接口名称
func (r *BaseRequest) GetAction() string {
	return r.action
}

// GetHttpMethod Http请求方法类型，POST、GET
func (r *BaseRequest) GetHttpMethod() string {
	return r.httpMethod
}

// GetParams 请求参数
func (r *BaseRequest) GetParams() map[string]string {
	return r.params
}

// GetPath 请求路径
func (r *BaseRequest) GetPath() string {
	return r.path
}

// GetDomain 域名
func (r *BaseRequest) GetDomain() string {
	return r.domain
}

func (r *BaseRequest) GetScheme() string {
	return r.scheme
}

// GetRootDomain 顶级域名
func (r *BaseRequest) GetRootDomain() string {
	return r.rootDomain
}

// GetServiceDomain 服务域名
func (r *BaseRequest) GetServiceDomain(service string) (domain string) {
	rootDomain := r.rootDomain
	if rootDomain == "" {
		rootDomain = RootDomain
	}
	domain = service + "." + rootDomain
	return
}

// SetDomain 设置域名
func (r *BaseRequest) SetDomain(domain string) {
	r.domain = domain
}

func (r *BaseRequest) SetScheme(scheme string) {
	scheme = strings.ToLower(scheme)
	switch scheme {
	case HTTP:
		r.scheme = HTTP
	default:
		r.scheme = HTTPS
	}
}

// SetRootDomain 设置顶级域名
func (r *BaseRequest) SetRootDomain(rootDomain string) {
	r.rootDomain = rootDomain
}

// SetHttpMethod 设置请求方法类型 POST/GET
func (r *BaseRequest) SetHttpMethod(method string) {
	switch strings.ToUpper(method) {
	case POST:
		{
			r.httpMethod = POST
		}
	case GET:
		{
			r.httpMethod = GET
		}
	default:
		{
			r.httpMethod = GET
		}
	}
}

// GetService 服务
func (r *BaseRequest) GetService() string {
	return r.service
}

// GetUrl 请求地址
func (r *BaseRequest) GetUrl() string {
	if r.httpMethod == GET {
		return r.GetScheme() + "://" + r.domain + r.path + "?" + GetUrlQueriesEncoded(r.params)
	} else if r.httpMethod == POST {
		return r.GetScheme() + "://" + r.domain + r.path
	} else {
		return ""
	}
}

// GetVersion 版本号
func (r *BaseRequest) GetVersion() string {
	return r.version
}

// GetUrlQueriesEncoded 请求参数转码
func GetUrlQueriesEncoded(params map[string]string) string {
	values := url.Values{}
	for key, value := range params {
		if value != "" {
			values.Add(key, value)
		}
	}
	return values.Encode()
}

// GetBodyReader 请求参数内容
func (r *BaseRequest) GetBodyReader() io.Reader {
	if r.httpMethod == POST {
		s := GetUrlQueriesEncoded(r.params)
		return strings.NewReader(s)
	} else {
		return strings.NewReader("")
	}
}

// Init 初始化
func (r *BaseRequest) Init() *BaseRequest {
	r.domain = ""
	r.path = Path
	r.params = make(map[string]string)
	r.formParams = make(map[string]string)
	return r
}

func (r *BaseRequest) WithApiInfo(service, version, action string) *BaseRequest {
	r.service = service
	r.version = version
	r.action = action
	return r
}

// GetServiceDomain Deprecated, use request.GetServiceDomain instead
func GetServiceDomain(service string) (domain string) {
	domain = service + "." + RootDomain
	return
}

// CompleteCommonParams 组装通用请求参数
func CompleteCommonParams(request Request, region string) {
	params := request.GetParams()
	params["Region"] = region
	if request.GetVersion() != "" {
		params["Version"] = request.GetVersion()
	}
	params["Action"] = request.GetAction()
	params["Timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	params["Nonce"] = strconv.Itoa(rand.Int())
}

func ConstructParams(req Request) (err error) {
	value := reflect.ValueOf(req).Elem()
	err = flatStructure(value, req, "")
	zlog.Debugf("[DEBUG] params=%s", req.GetParams())
	return
}

func flatStructure(value reflect.Value, request Request, prefix string) (err error) {
	zlog.Debugf("[DEBUG] reflect value: %v", value.Type())
	valueType := value.Type()
	for i := 0; i < valueType.NumField(); i++ {
		tag := valueType.Field(i).Tag
		nameTag, hasNameTag := tag.Lookup("name")
		if !hasNameTag {
			continue
		}
		field := value.Field(i)
		kind := field.Kind()
		if kind == reflect.Ptr && field.IsNil() {
			continue
		}
		if kind == reflect.Ptr {
			field = field.Elem()
			kind = field.Kind()
		}
		key := prefix + nameTag
		if kind == reflect.String {
			s := field.String()
			if s != "" {
				request.GetParams()[key] = s
			}
		} else if kind == reflect.Bool {
			request.GetParams()[key] = strconv.FormatBool(field.Bool())
		} else if kind == reflect.Int || kind == reflect.Int64 {
			request.GetParams()[key] = strconv.FormatInt(field.Int(), 10)
		} else if kind == reflect.Uint || kind == reflect.Uint64 {
			request.GetParams()[key] = strconv.FormatUint(field.Uint(), 10)
		} else if kind == reflect.Float64 {
			request.GetParams()[key] = strconv.FormatFloat(field.Float(), 'f', -1, 64)
		} else if kind == reflect.Slice {
			list := value.Field(i)
			for j := 0; j < list.Len(); j++ {
				vj := list.Index(j)
				key := prefix + nameTag + "." + strconv.Itoa(j)
				kind = vj.Kind()
				if kind == reflect.Ptr && vj.IsNil() {
					continue
				}
				if kind == reflect.Ptr {
					vj = vj.Elem()
					kind = vj.Kind()
				}
				if kind == reflect.String {
					request.GetParams()[key] = vj.String()
				} else if kind == reflect.Bool {
					request.GetParams()[key] = strconv.FormatBool(vj.Bool())
				} else if kind == reflect.Int || kind == reflect.Int64 {
					request.GetParams()[key] = strconv.FormatInt(vj.Int(), 10)
				} else if kind == reflect.Uint || kind == reflect.Uint64 {
					request.GetParams()[key] = strconv.FormatUint(vj.Uint(), 10)
				} else if kind == reflect.Float64 {
					request.GetParams()[key] = strconv.FormatFloat(vj.Float(), 'f', -1, 64)
				} else {
					if err = flatStructure(vj, request, key+"."); err != nil {
						return
					}
				}
			}
		} else {
			if err = flatStructure(reflect.ValueOf(field.Interface()), request, prefix+nameTag+"."); err != nil {
				return
			}
		}
	}
	return
}
