package common

// Credential 凭证
type Credential struct {
	SecretId  string
	SecretKey string
	Token     string
}

// NewCredential 凭证实例化
func NewCredential(secretId, secretKey string) *Credential {
	return &Credential{
		SecretId:  secretId,
		SecretKey: secretKey,
	}
}

// NewTokenCredential 凭证实例化
func NewTokenCredential(secretId, secretKey, token string) *Credential {
	return &Credential{
		SecretId:  secretId,
		SecretKey: secretKey,
		Token:     token,
	}
}

// GetCredentialParams 凭证参数
func (c *Credential) GetCredentialParams() map[string]string {
	p := map[string]string{
		"SecretId": c.SecretId,
	}
	if c.Token != "" {
		p["Token"] = c.Token
	}

	return p
}
