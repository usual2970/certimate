package dogecloud

import (
	"context"
	"net/http"
)

type BindCdnCertRequest struct {
	CertId int64  `json:"id"`
	Domain string `json:"domain"`
}

type BindCdnCertResponse struct {
	apiResponseBase
}

func (c *Client) BindCdnCert(req *BindCdnCertRequest) (*BindCdnCertResponse, error) {
	return c.BindCdnCertWithContext(context.Background(), req)
}

func (c *Client) BindCdnCertWithContext(ctx context.Context, req *BindCdnCertRequest) (*BindCdnCertResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/cdn/cert/bind.json")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &BindCdnCertResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
