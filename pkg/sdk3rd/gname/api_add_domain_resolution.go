package gname

import (
	"context"
	"encoding/json"
	"net/http"
)

type AddDomainResolutionRequest struct {
	ZoneName    *string `json:"ym,omitempty"`
	RecordType  *string `json:"lx,omitempty"`
	RecordName  *string `json:"zj,omitempty"`
	RecordValue *string `json:"jlz,omitempty"`
	MX          *int32  `json:"mx,omitempty"`
	TTL         *int32  `json:"ttl,omitempty"`
}

type AddDomainResolutionResponse struct {
	apiResponseBase

	Data json.Number `json:"data"`
}

func (c *Client) AddDomainResolution(req *AddDomainResolutionRequest) (*AddDomainResolutionResponse, error) {
	return c.AddDomainResolutionWithContext(context.Background(), req)
}

func (c *Client) AddDomainResolutionWithContext(ctx context.Context, req *AddDomainResolutionRequest) (*AddDomainResolutionResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/api/resolution/add", req)
	if err != nil {
		return nil, err
	} else {
		httpreq.SetContext(ctx)
	}

	result := &AddDomainResolutionResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
