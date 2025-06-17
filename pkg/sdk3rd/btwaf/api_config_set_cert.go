package btwaf

import (
	"context"
	"net/http"
)

type ConfigSetCertRequest struct {
	CertContent *string `json:"certContent,omitempty"`
	KeyContent  *string `json:"keyContent,omitempty"`
}

type ConfigSetCertResponse struct {
	apiResponseBase
}

func (c *Client) ConfigSetCert(req *ConfigSetCertRequest) (*ConfigSetCertResponse, error) {
	return c.ConfigSetCertWithContext(context.Background(), req)
}

func (c *Client) ConfigSetCertWithContext(ctx context.Context, req *ConfigSetCertRequest) (*ConfigSetCertResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/config/set_cert")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &ConfigSetCertResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
