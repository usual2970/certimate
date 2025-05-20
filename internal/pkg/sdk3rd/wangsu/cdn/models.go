package cdn

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

type BatchUpdateCertificateConfigRequest struct {
	CertificateId int64    `json:"certificateId" required:"true"`
	DomainNames   []string `json:"domainNames" required:"true"`
}

type BatchUpdateCertificateConfigResponse struct {
	baseResponse
}
