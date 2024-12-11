package huaweicloudcdnsdk

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	hcCdn "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2"
	hcCdnModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"
)

type Client struct {
	hcCdn.CdnClient
}

func NewClient(hcClient *core.HcHttpClient) *Client {
	return &Client{
		CdnClient: *hcCdn.NewCdnClient(hcClient),
	}
}

func (c *Client) UploadDomainMultiCertificatesEx(request *UpdateDomainMultiCertificatesExRequest) (*UpdateDomainMultiCertificatesExResponse, error) {
	requestDef := hcCdn.GenReqDefForUpdateDomainMultiCertificates()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		temp := resp.(*hcCdnModel.UpdateDomainMultiCertificatesResponse)
		return &UpdateDomainMultiCertificatesExResponse{UpdateDomainMultiCertificatesResponse: *temp}, nil
	}
}
