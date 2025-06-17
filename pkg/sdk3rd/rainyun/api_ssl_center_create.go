package rainyun

import (
	"context"
	"net/http"
)

type SslCenterCreateRequest struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

type SslCenterCreateResponse struct {
	apiResponseBase
}

func (c *Client) SslCenterCreate(req *SslCenterCreateRequest) (*SslCenterCreateResponse, error) {
	return c.SslCenterCreateWithContext(context.Background(), req)
}

func (c *Client) SslCenterCreateWithContext(ctx context.Context, req *SslCenterCreateRequest) (*SslCenterCreateResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/product/sslcenter/")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &SslCenterCreateResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
