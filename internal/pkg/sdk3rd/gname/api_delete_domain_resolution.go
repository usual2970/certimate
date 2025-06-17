package gname

import (
	"context"
	"net/http"
)

type DeleteDomainResolutionRequest struct {
	ZoneName *string `json:"ym,omitempty"`
	RecordID *int64  `json:"jxid,omitempty"`
}

type DeleteDomainResolutionResponse struct {
	apiResponseBase
}

func (c *Client) DeleteDomainResolution(req *DeleteDomainResolutionRequest) (*DeleteDomainResolutionResponse, error) {
	return c.DeleteDomainResolutionWithContext(context.Background(), req)
}

func (c *Client) DeleteDomainResolutionWithContext(ctx context.Context, req *DeleteDomainResolutionRequest) (*DeleteDomainResolutionResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/api/resolution/delete", req)
	if err != nil {
		return nil, err
	} else {
		httpreq.SetContext(ctx)
	}

	result := &DeleteDomainResolutionResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
