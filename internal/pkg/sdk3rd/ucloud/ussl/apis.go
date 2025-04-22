package ussl

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

type UploadNormalCertificateRequest struct {
	request.CommonBase

	CertificateName *string `required:"true"`
	SslPublicKey    *string `required:"true"`
	SslPrivateKey   *string `required:"true"`
	SslMD5          *string `required:"true"`
	SslCaKey        *string `required:"false"`
}

type UploadNormalCertificateResponse struct {
	response.CommonBase

	CertificateID  int
	LongResourceID string
}

func (c *USSLClient) NewUploadNormalCertificateRequest() *UploadNormalCertificateRequest {
	req := &UploadNormalCertificateRequest{}

	c.Client.SetupRequest(req)

	req.SetRetryable(false)
	return req
}

func (c *USSLClient) UploadNormalCertificate(req *UploadNormalCertificateRequest) (*UploadNormalCertificateResponse, error) {
	var err error
	var res UploadNormalCertificateResponse

	reqCopier := *req

	err = c.Client.InvokeAction("UploadNormalCertificate", &reqCopier, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}

type GetCertificateListRequest struct {
	request.CommonBase

	Mode           *string `required:"true"`
	StateCode      *string `required:"false"`
	Brand          *string `required:"false"`
	CaOrganization *string `required:"false"`
	Domain         *string `required:"false"`
	Sort           *string `required:"false"`
	Page           *int    `required:"false"`
	PageSize       *int    `required:"false"`
}

type GetCertificateListResponse struct {
	response.CommonBase

	CertificateList []*CertificateListItem
	TotalCount      int
}

func (c *USSLClient) NewGetCertificateListRequest() *GetCertificateListRequest {
	req := &GetCertificateListRequest{}

	c.Client.SetupRequest(req)

	req.SetRetryable(false)
	return req
}

func (c *USSLClient) GetCertificateList(req *GetCertificateListRequest) (*GetCertificateListResponse, error) {
	var err error
	var res GetCertificateListResponse

	reqCopier := *req

	err = c.Client.InvokeAction("GetCertificateList", &reqCopier, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}

type GetCertificateDetailInfoRequest struct {
	request.CommonBase

	CertificateID *int `required:"true"`
}

type GetCertificateDetailInfoResponse struct {
	response.CommonBase

	CertificateInfo *CertificateInfo
}

func (c *USSLClient) NewGetCertificateDetailInfoRequest() *GetCertificateDetailInfoRequest {
	req := &GetCertificateDetailInfoRequest{}

	c.Client.SetupRequest(req)

	req.SetRetryable(false)
	return req
}

func (c *USSLClient) GetCertificateDetailInfo(req *GetCertificateDetailInfoRequest) (*GetCertificateDetailInfoResponse, error) {
	var err error
	var res GetCertificateDetailInfoResponse

	reqCopier := *req

	err = c.Client.InvokeAction("GetCertificateDetailInfo", &reqCopier, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}

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
