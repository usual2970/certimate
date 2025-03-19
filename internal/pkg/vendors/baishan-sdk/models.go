package baishansdk

type BaseResponse interface {
	GetCode() int
	GetMessage() string
}

type baseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r *baseResponse) GetCode() int {
	return r.Code
}

func (r *baseResponse) GetMessage() string {
	return r.Message
}

type CreateCertificateRequest struct {
	Certificate string `json:"certificate"`
	Key         string `json:"key"`
	Name        string `json:"name"`
}

type CreateCertificateResponse struct {
	baseResponse
	Data *DomainCertificate `json:"data"`
}

type GetDomainConfigRequest struct {
	Domains string `json:"domains"`
	Config  []string `json:"config"`
}

type GetDomainConfigResponse struct {
	baseResponse
	Data []*struct {
		Domain string        `json:"domain"`
		Config *DomainConfig `json:"config"`
	} `json:"data"`
}

type SetDomainConfigRequest struct {
	Domains string        `json:"domains"`
	Config  *DomainConfig `json:"config"`
}

type SetDomainConfigResponse struct {
	baseResponse
	Data *struct {
		Config *DomainConfig `json:"config"`
	} `json:"data"`
}

type DomainCertificate struct {
	CertId         int64  `json:"cert_id"`
	Name           string `json:"name"`
	CertStartTime  string `json:"cert_start_time"`
	CertExpireTime string `json:"cert_expire_time"`
}

type DomainConfig struct {
	Https *DomainConfigHttps `json:"https"`
}

type DomainConfigHttps struct {
	CertId      int64   `json:"cert_id"`
	ForceHttps  *string `json:"force_https,omitempty"`
	EnableHttp2 *string `json:"http2,omitempty"`
	EnableOcsp  *string `json:"ocsp,omitempty"`
}
