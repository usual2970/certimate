package gname

import (
	"context"
	"net/http"
)

type ModifyDomainResolutionRequest struct {
	ID          *int64  `json:"jxid,omitempty"`
	ZoneName    *string `json:"ym,omitempty"`
	RecordType  *string `json:"lx,omitempty"`
	RecordName  *string `json:"zj,omitempty"`
	RecordValue *string `json:"jlz,omitempty"`
	MX          *int32  `json:"mx,omitempty"`
	TTL         *int32  `json:"ttl,omitempty"`
}

type ModifyDomainResolutionResponse struct {
	apiResponseBase
}

func (c *Client) ModifyDomainResolution(req *ModifyDomainResolutionRequest) (*ModifyDomainResolutionResponse, error) {
	return c.ModifyDomainResolutionWithContext(context.Background(), req)
}

func (c *Client) ModifyDomainResolutionWithContext(ctx context.Context, req *ModifyDomainResolutionRequest) (*ModifyDomainResolutionResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/api/resolution/edit", req)
	if err != nil {
		return nil, err
	} else {
		httpreq.SetContext(ctx)
	}

	result := &ModifyDomainResolutionResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
