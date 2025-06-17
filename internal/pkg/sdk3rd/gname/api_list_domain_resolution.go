package gname

import (
	"context"
	"net/http"
)

type ListDomainResolutionRequest struct {
	ZoneName *string `json:"ym,omitempty"`
	Page     *int32  `json:"page,omitempty"`
	PageSize *int32  `json:"limit,omitempty"`
}

type ListDomainResolutionResponse struct {
	apiResponseBase

	Count    int32                        `json:"count"`
	Data     []*DomainResolutionRecordord `json:"data"`
	Page     int32                        `json:"page"`
	PageSize int32                        `json:"pagesize"`
}

func (c *Client) ListDomainResolution(req *ListDomainResolutionRequest) (*ListDomainResolutionResponse, error) {
	return c.ListDomainResolutionWithContext(context.Background(), req)
}

func (c *Client) ListDomainResolutionWithContext(ctx context.Context, req *ListDomainResolutionRequest) (*ListDomainResolutionResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/api/resolution/list", req)
	if err != nil {
		return nil, err
	} else {
		httpreq.SetContext(ctx)
	}

	result := &ListDomainResolutionResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
