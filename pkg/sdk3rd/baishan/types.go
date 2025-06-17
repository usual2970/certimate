package baishan

import "encoding/json"

type apiResponse interface {
	GetCode() int32
	GetMessage() string
}

type apiResponseBase struct {
	Code    *int32  `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
}

func (r *apiResponseBase) GetCode() int32 {
	if r.Code == nil {
		return 0
	}

	return *r.Code
}

func (r *apiResponseBase) GetMessage() string {
	if r.Message == nil {
		return ""
	}

	return *r.Message
}

var _ apiResponse = (*apiResponseBase)(nil)

type DomainCertificate struct {
	CertId         json.Number `json:"cert_id"`
	Name           string      `json:"name"`
	CertStartTime  string      `json:"cert_start_time"`
	CertExpireTime string      `json:"cert_expire_time"`
}

type DomainConfig struct {
	Https *DomainConfigHttps `json:"https"`
}

type DomainConfigHttps struct {
	CertId      json.Number `json:"cert_id"`
	ForceHttps  *string     `json:"force_https,omitempty"`
	EnableHttp2 *string     `json:"http2,omitempty"`
	EnableOcsp  *string     `json:"ocsp,omitempty"`
}
