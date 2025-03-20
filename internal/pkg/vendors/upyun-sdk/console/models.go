package console

import (
	"encoding/json"
)

type baseResponse struct {
	Data *baseResponseData `json:"data,omitempty"`
}

func (r *baseResponse) GetData() *baseResponseData {
	return r.Data
}

type baseResponseData struct {
	ErrorCode    json.Number `json:"error_code"`
	ErrorMessage string      `json:"message"`
}

func (r *baseResponseData) GetErrorCode() int {
	if r.ErrorCode.String() == "" {
		return 0
	}

	errcode, err := r.ErrorCode.Int64()
	if err != nil {
		return -1
	}

	return int(errcode)
}

func (r *baseResponseData) GetErrorMessage() string {
	return r.ErrorMessage
}

type signinRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signinResponse struct {
	baseResponse
	Data struct {
		baseResponseData
		Result bool `json:"result"`
	} `json:"data"`
}

type UploadHttpsCertificateRequest struct {
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key"`
}

type UploadHttpsCertificateResponse struct {
	baseResponse
	Data *struct {
		baseResponseData
		Status int `json:"status"`
		Result struct {
			CertificateId string `json:"certificate_id"`
			CommonName    string `json:"commonName"`
			Serial        string `json:"serial"`
		} `json:"result"`
	} `json:"data"`
}

type GetHttpsCertificateManagerRequest struct {
	CertificateId string `json:"certificate_id"`
}

type GetHttpsCertificateManagerResponse struct {
	baseResponse
	Data *struct {
		baseResponseData
		AuthenticateNum     int32                           `json:"authenticate_num"`
		AuthenticateDomains []string                        `json:"authenticate_domain"`
		Domains             []HttpsCertificateManagerDomain `json:"domains"`
	} `json:"data"`
}

type HttpsCertificateManagerDomain struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	BucketId   int64  `json:"bucket_id"`
	BucketName string `json:"bucket_name"`
}

type UpdateHttpsCertificateManagerRequest struct {
	CertificateId string `json:"certificate_id"`
	Domain        string `json:"domain"`
	Https         bool   `json:"https"`
	ForceHttps    bool   `json:"force_https"`
}

type UpdateHttpsCertificateManagerResponse struct {
	baseResponse
	Data *struct {
		baseResponseData
		Status bool `json:"status"`
	} `json:"data"`
}

type GetHttpsServiceManagerRequest struct {
	Domain string `json:"domain"`
}

type GetHttpsServiceManagerResponse struct {
	baseResponse
	Data *struct {
		baseResponseData
		Status  int                         `json:"status"`
		Domains []HttpsServiceManagerDomain `json:"result"`
	} `json:"data"`
}

type HttpsServiceManagerDomain struct {
	CertificateId string `json:"certificate_id"`
	CommonName    string `json:"commonName"`
	Https         bool   `json:"https"`
	ForceHttps    bool   `json:"force_https"`
	PaymentType   string `json:"payment_type"`
	DomainType    string `json:"domain_type"`
	Validity      struct {
		Start int64 `json:"start"`
		End   int64 `json:"end"`
	} `json:"validity"`
}

type MigrateHttpsDomainRequest struct {
	CertificateId string `json:"crt_id"`
	Domain        string `json:"domain_name"`
}

type MigrateHttpsDomainResponse struct {
	baseResponse
	Data *struct {
		baseResponseData
		Status bool `json:"status"`
	} `json:"data"`
}
