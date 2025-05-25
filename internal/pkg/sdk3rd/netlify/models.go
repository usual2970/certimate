package netlify

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

type ProvisionSiteTLSCertificateParams struct {
	Certificate    string `json:"certificate"`
	CACertificates string `json:"ca_certificates"`
	Key            string `json:"key"`
}

type ProvisionSiteTLSCertificateResponse struct {
	baseResponse
	Domains   []string `json:"domains,omitempty"`
	State     string   `json:"state,omitempty"`
	ExpiresAt string   `json:"expires_at,omitempty"`
	CreatedAt string   `json:"created_at,omitempty"`
	UpdatedAt string   `json:"updated_at,omitempty"`
}
