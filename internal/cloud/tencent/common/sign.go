package common

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	tcHttp "github.com/eadydb/k8s-aim/internal/cloud/tencent/common/http"
	"sort"
)

const (
	SHA256 = "HmacSHA256"
	SHA1   = "HmacSHA1"
)

// Sign 签名
func Sign(s, secretKey, method string) string {
	hashed := hmac.New(sha1.New, []byte(secretKey))
	if method == SHA256 {
		hashed = hmac.New(sha256.New, []byte(secretKey))
	}
	hashed.Write([]byte(s))

	return base64.StdEncoding.EncodeToString(hashed.Sum(nil))
}

func sha256hex(s string) string {
	b := sha256.Sum256([]byte(s))
	return hex.EncodeToString(b[:])
}

func hmacsha256(s, key string) string {
	hashed := hmac.New(sha256.New, []byte(key))
	hashed.Write([]byte(s))
	return string(hashed.Sum(nil))
}

// signRequest 请求设置签名参数
func signRequest(request tcHttp.Request, credential *Credential, method string) (err error) {
	if method != SHA256 {
		method = SHA1
	}
	checkAuthParams(request, credential, method)
	s := getStringToSign(request)
	signature := Sign(s, credential.SecretKey, method)
	request.GetParams()["Signature"] = signature
	return
}

// checkAuthParams 检查认证参数
func checkAuthParams(request tcHttp.Request, credential *Credential, method string) {
	params := request.GetParams()
	credentialParams := credential.GetCredentialParams()
	for key, value := range credentialParams {
		params[key] = value
	}
	params["SignatureMethod"] = method
	delete(params, "Signature")
}

func getStringToSign(request tcHttp.Request) string {
	method := request.GetHttpMethod()
	domain := request.GetDomain()
	path := request.GetPath()

	var buf bytes.Buffer
	buf.WriteString(method)
	buf.WriteString(domain)
	buf.WriteString(path)
	buf.WriteString("?")

	params := request.GetParams()
	// sort params
	keys := make([]string, 0, len(params))
	for k, _ := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i := range keys {
		k := keys[i]
		// TODO: check if server side allows empty value in url.
		if params[k] == "" {
			continue
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(params[k])
		buf.WriteString("&")
	}
	buf.Truncate(buf.Len() - 1)
	return buf.String()
}
