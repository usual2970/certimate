package ratpanel

import (
	"context"
	"net/http"
)

type SetWebsiteCertRequest struct {
	SiteName    string `json:"name"`
	Certificate string `json:"cert"`
	PrivateKey  string `json:"key"`
}

type SetWebsiteCertResponse struct {
	apiResponseBase
}

func (c *Client) SetWebsiteCert(req *SetWebsiteCertRequest) (*SetWebsiteCertResponse, error) {
	return c.SetWebsiteCertWithContext(context.Background(), req)
}

func (c *Client) SetWebsiteCertWithContext(ctx context.Context, req *SetWebsiteCertRequest) (*SetWebsiteCertResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/website/cert")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &SetWebsiteCertResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
