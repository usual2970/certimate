package btpanel

import (
	"context"
	"net/http"
)

type ConfigSavePanelSSLRequest struct {
	PrivateKey  string `json:"privateKey"`
	Certificate string `json:"certPem"`
}

type ConfigSavePanelSSLResponse struct {
	apiResponseBase
}

func (c *Client) ConfigSavePanelSSL(req *ConfigSavePanelSSLRequest) (*ConfigSavePanelSSLResponse, error) {
	return c.ConfigSavePanelSSLWithContext(context.Background(), req)
}

func (c *Client) ConfigSavePanelSSLWithContext(ctx context.Context, req *ConfigSavePanelSSLRequest) (*ConfigSavePanelSSLResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/config?action=SavePanelSSL", req)
	if err != nil {
		return nil, err
	} else {
		httpreq.SetContext(ctx)
	}

	result := &ConfigSavePanelSSLResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
