package ussl

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

type DownloadCertificateRequest struct {
	request.CommonBase

	CertificateID *int `required:"true"`
}

type DownloadCertificateResponse struct {
	response.CommonBase

	CertificateUrl string
	CertCA         *CertificateDownloadInfo
	Certificate    *CertificateDownloadInfo
}

func (c *USSLClient) NewDownloadCertificateRequest() *DownloadCertificateRequest {
	req := &DownloadCertificateRequest{}

	c.Client.SetupRequest(req)

	req.SetRetryable(false)
	return req
}

func (c *USSLClient) DownloadCertificate(req *DownloadCertificateRequest) (*DownloadCertificateResponse, error) {
	var err error
	var res DownloadCertificateResponse

	reqCopier := *req

	err = c.Client.InvokeAction("DownloadCertificate", &reqCopier, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}
