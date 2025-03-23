package cacheflysdk

type BaseResponse interface {
	GetMessage() string
}

type baseResponse struct {
	Message *string `json:"message,omitempty"`
}

func (r *baseResponse) GetMessage() string {
	if r.Message != nil {
		return *r.Message
	}
	return ""
}

type CreateCertificateRequest struct {
	Certificate    string  `json:"certificate"`
	CertificateKey string  `json:"certificateKey"`
	Password       *string `json:"password"`
}

type CreateCertificateResponse struct {
	baseResponse
	Id                string   `json:"_id"`
	SubjectCommonName string   `json:"subjectCommonName"`
	SubjectNames      []string `json:"subjectNames"`
	Expired           bool     `json:"expired"`
	Expiring          bool     `json:"expiring"`
	InUse             bool     `json:"inUse"`
	Managed           bool     `json:"managed"`
	Services          []string `json:"services"`
	Domains           []string `json:"domains"`
	NotBefore         string   `json:"notBefore"`
	NotAfter          string   `json:"notAfter"`
	CreatedAt         string   `json:"createdAt"`
}
