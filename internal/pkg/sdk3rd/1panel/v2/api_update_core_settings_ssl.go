package onepanelv2

import (
	"context"
	"net/http"
)

type UpdateCoreSettingsSSLRequest struct {
	Cert        string `json:"cert"`
	Key         string `json:"key"`
	SSLType     string `json:"sslType"`
	SSL         string `json:"ssl"`
	SSLID       int64  `json:"sslID"`
	AutoRestart string `json:"autoRestart"`
}

type UpdateCoreSettingsSSLResponse struct {
	apiResponseBase
}

func (c *Client) UpdateCoreSettingsSSL(req *UpdateCoreSettingsSSLRequest) (*UpdateCoreSettingsSSLResponse, error) {
	return c.UpdateCoreSettingsSSLWithContext(context.Background(), req)
}

func (c *Client) UpdateCoreSettingsSSLWithContext(ctx context.Context, req *UpdateCoreSettingsSSLRequest) (*UpdateCoreSettingsSSLResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/core/settings/ssl/update")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &UpdateCoreSettingsSSLResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
