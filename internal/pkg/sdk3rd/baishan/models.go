package baishansdk

import "encoding/json"

type BaseResponse interface {
	GetCode() int32
	GetMessage() string
}

type baseResponse struct {
	Code    *int32  `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
}

func (r *baseResponse) GetCode() int32 {
	if r.Code != nil {
		return *r.Code
	}
	return 0
}

func (r *baseResponse) GetMessage() string {
	if r.Message != nil {
		return *r.Message
	}
	return ""
}

type CreateCertificateRequest struct {
	CertificateId *string `json:"cert_id,omitempty"`
	Certificate   string  `json:"certificate"`
	Key           string  `json:"key"`
	Name          string  `json:"name"`
}

type CreateCertificateResponse struct {
	baseResponse
	Data *DomainCertificate `json:"data,omitempty"`
}

type GetDomainConfigRequest struct {
	Domains string   `json:"domains"`
	Config  []string `json:"config"`
}

type GetDomainConfigResponse struct {
	baseResponse
	Data []*struct {
		Domain string        `json:"domain"`
		Config *DomainConfig `json:"config"`
	} `json:"data,omitempty"`
}

type SetDomainConfigRequest struct {
	Domains string        `json:"domains"`
	Config  *DomainConfig `json:"config"`
}

type SetDomainConfigResponse struct {
	baseResponse
	Data *struct {
		Config *DomainConfig `json:"config"`
	} `json:"data,omitempty"`
}

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
