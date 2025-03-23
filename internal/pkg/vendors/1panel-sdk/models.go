package onepanelsdk

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

type UpdateSystemSSLRequest struct {
	Cert        string `json:"cert"`
	Key         string `json:"key"`
	SSLType     string `json:"sslType"`
	SSL         string `json:"ssl"`
	SSLID       int64  `json:"sslID"`
	AutoRestart string `json:"autoRestart"`
}

type UpdateSystemSSLResponse struct {
	baseResponse
}

type SearchWebsiteSSLRequest struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"pageSize"`
}

type SearchWebsiteSSLResponse struct {
	baseResponse
	Data *struct {
		Items []*struct {
			ID          int64  `json:"id"`
			PEM         string `json:"pem"`
			PrivateKey  string `json:"privateKey"`
			Domains     string `json:"domains"`
			Description string `json:"description"`
			Status      string `json:"status"`
			UpdatedAt   string `json:"updatedAt"`
			CreatedAt   string `json:"createdAt"`
		} `json:"items"`
		Total int32 `json:"total"`
	} `json:"data,omitempty"`
}

type UploadWebsiteSSLRequest struct {
	Type            string `json:"type"`
	SSLID           int64  `json:"sslID"`
	Certificate     string `json:"certificate"`
	CertificatePath string `json:"certificatePath"`
	PrivateKey      string `json:"privateKey"`
	PrivateKeyPath  string `json:"privateKeyPath"`
	Description     string `json:"description"`
}

type UploadWebsiteSSLResponse struct {
	baseResponse
}

type GetHttpsConfRequest struct {
	WebsiteID int64 `json:"-"`
}

type GetHttpsConfResponse struct {
	baseResponse
	Data *struct {
		Enable      bool     `json:"enable"`
		HttpConfig  string   `json:"httpConfig"`
		SSLProtocol []string `json:"SSLProtocol"`
		Algorithm   string   `json:"algorithm"`
		Hsts        bool     `json:"hsts"`
	} `json:"data,omitempty"`
}

type UpdateHttpsConfRequest struct {
	WebsiteID       int64    `json:"websiteId"`
	Enable          bool     `json:"enable"`
	Type            string   `json:"type"`
	WebsiteSSLID    int64    `json:"websiteSSLId"`
	PrivateKey      string   `json:"privateKey"`
	Certificate     string   `json:"certificate"`
	PrivateKeyPath  string   `json:"privateKeyPath"`
	CertificatePath string   `json:"certificatePath"`
	ImportType      string   `json:"importType"`
	HttpConfig      string   `json:"httpConfig"`
	SSLProtocol     []string `json:"SSLProtocol"`
	Algorithm       string   `json:"algorithm"`
	Hsts            bool     `json:"hsts"`
}

type UpdateHttpsConfResponse struct {
	baseResponse
}
