package certcenter

import (
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"github.com/volcengine/volcengine-go-sdk/volcengine/request"
)

type CertCenterAPI interface {
	ImportCertificate(*ImportCertificateInput) (*ImportCertificateOutput, error)
	ImportCertificateWithContext(volcengine.Context, *ImportCertificateInput, ...request.Option) (*ImportCertificateOutput, error)
	ImportCertificateRequest(*ImportCertificateInput) (*request.Request, *ImportCertificateOutput)
}

var _ CertCenterAPI = (*CertCenter)(nil)
