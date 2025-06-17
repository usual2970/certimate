package btpanel

import (
	"context"
	"net/http"
)

type SiteSetSSLRequest struct {
	Type        string `json:"type"`
	SiteName    string `json:"siteName"`
	PrivateKey  string `json:"key"`
	Certificate string `json:"csr"`
}

type SiteSetSSLResponse struct {
	apiResponseBase
}

func (c *Client) SiteSetSSL(req *SiteSetSSLRequest) (*SiteSetSSLResponse, error) {
	return c.SiteSetSSLWithContext(context.Background(), req)
}

func (c *Client) SiteSetSSLWithContext(ctx context.Context, req *SiteSetSSLRequest) (*SiteSetSSLResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/site?action=SetSSL", req)
	if err != nil {
		return nil, err
	} else {
		httpreq.SetContext(ctx)
	}

	result := &SiteSetSSLResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
