package ufile

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

type AddUFileSSLCertRequest struct {
	request.CommonBase

	BucketName      *string `required:"true"`
	Domain          *string `required:"true"`
	CertificateName *string `required:"true"`
	USSLId          *string `required:"false"`
}

type AddUFileSSLCertResponse struct {
	response.CommonBase
}

func (c *UFileClient) NewAddUFileSSLCertRequest() *AddUFileSSLCertRequest {
	req := &AddUFileSSLCertRequest{}

	c.Client.SetupRequest(req)

	req.SetRetryable(false)
	return req
}

func (c *UFileClient) AddUFileSSLCert(req *AddUFileSSLCertRequest) (*AddUFileSSLCertResponse, error) {
	var err error
	var res AddUFileSSLCertResponse

	reqCopier := *req

	err = c.Client.InvokeAction("AddUFileSSLCert", &reqCopier, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}
