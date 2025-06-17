package ratpanel

import (
	"context"
	"net/http"
)

type SetSettingCertRequest struct {
	Certificate string `json:"cert"`
	PrivateKey  string `json:"key"`
}

type SetSettingCertResponse struct {
	apiResponseBase
}

func (c *Client) SetSettingCert(req *SetSettingCertRequest) (*SetSettingCertResponse, error) {
	return c.SetSettingCertWithContext(context.Background(), req)
}

func (c *Client) SetSettingCertWithContext(ctx context.Context, req *SetSettingCertRequest) (*SetSettingCertResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/setting/cert")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &SetSettingCertResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
