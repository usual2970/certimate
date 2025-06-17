package dogecloud

import (
	"context"
	"net/http"
)

type UploadCdnCertRequest struct {
	Note        string `json:"note"`
	Certificate string `json:"cert"`
	PrivateKey  string `json:"private"`
}

type UploadCdnCertResponse struct {
	apiResponseBase

	Data *struct {
		Id int64 `json:"id"`
	} `json:"data,omitempty"`
}

func (c *Client) UploadCdnCert(req *UploadCdnCertRequest) (*UploadCdnCertResponse, error) {
	return c.UploadCdnCertWithContext(context.Background(), req)
}

func (c *Client) UploadCdnCertWithContext(ctx context.Context, req *UploadCdnCertRequest) (*UploadCdnCertResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/cdn/cert/upload.json")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &UploadCdnCertResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
