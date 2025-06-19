package baishan

import (
	"context"
	"net/http"
)

type SetDomainConfigRequest struct {
	Domains *string       `json:"domains,omitempty"`
	Config  *DomainConfig `json:"config,omitempty"`
}

type SetDomainConfigResponse struct {
	apiResponseBase

	Data *struct {
		Config *DomainConfig `json:"config"`
	} `json:"data,omitempty"`
}

func (c *Client) SetDomainConfig(req *SetDomainConfigRequest) (*SetDomainConfigResponse, error) {
	return c.SetDomainConfigWithContext(context.Background(), req)
}

func (c *Client) SetDomainConfigWithContext(ctx context.Context, req *SetDomainConfigRequest) (*SetDomainConfigResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/v2/domain/config")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &SetDomainConfigResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
