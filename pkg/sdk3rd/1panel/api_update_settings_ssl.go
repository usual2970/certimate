package onepanel

import (
	"context"
	"net/http"
)

type UpdateSettingsSSLRequest struct {
	Cert        string `json:"cert"`
	Key         string `json:"key"`
	SSLType     string `json:"sslType"`
	SSL         string `json:"ssl"`
	SSLID       int64  `json:"sslID"`
	AutoRestart string `json:"autoRestart"`
}

type UpdateSettingsSSLResponse struct {
	apiResponseBase
}

func (c *Client) UpdateSettingsSSL(req *UpdateSettingsSSLRequest) (*UpdateSettingsSSLResponse, error) {
	return c.UpdateSettingsSSLWithContext(context.Background(), req)
}

func (c *Client) UpdateSettingsSSLWithContext(ctx context.Context, req *UpdateSettingsSSLRequest) (*UpdateSettingsSSLResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/settings/ssl/update")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &UpdateSettingsSSLResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
