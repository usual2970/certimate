package certificate

import (
	"github.com/usual2970/certimate/internal/pkg/sdk3rd/wangsu/openapi"
)

type baseResponse struct {
	RequestId *string `json:"requestId,omitempty"`
	Code      *string `json:"code,omitempty"`
	Message   *string `json:"message,omitempty"`
}

var _ openapi.Result = (*baseResponse)(nil)

func (r *baseResponse) SetRequestId(requestId string) {
	r.RequestId = &requestId
}

type CreateCertificateRequest struct {
	Name        *string `json:"name,omitempty" required:"true"`
	Certificate *string `json:"certificate,omitempty" required:"true"`
	PrivateKey  *string `json:"privateKey,omitempty"`
	Comment     *string `json:"comment,omitempty" `
}

type CreateCertificateResponse struct {
	baseResponse
	CertificateUrl string `json:"location,omitempty"`
}

type UpdateCertificateRequest struct {
	Name        *string `json:"name,omitempty" required:"true"`
	Certificate *string `json:"certificate,omitempty"`
	PrivateKey  *string `json:"privateKey,omitempty"`
	Comment     *string `json:"comment,omitempty" `
}

type UpdateCertificateResponse struct {
	baseResponse
}

type ListCertificatesResponse struct {
	baseResponse
	Certificates []*struct {
		CertificateId string `json:"certificate-id"`
		Name          string `json:"name"`
		Comment       string `json:"comment"`
		ValidityFrom  string `json:"certificate-validity-from"`
		ValidityTo    string `json:"certificate-validity-to"`
		Serial        string `json:"certificate-serial"`
	} `json:"ssl-certificates,omitempty"`
}
