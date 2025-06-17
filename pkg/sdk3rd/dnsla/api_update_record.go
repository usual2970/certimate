package dnsla

import (
	"context"
	"net/http"
)

type UpdateRecordRequest struct {
	Id         string  `json:"id"`
	GroupId    *string `json:"groupId,omitempty"`
	LineId     *string `json:"lineId,omitempty"`
	Type       *int32  `json:"type,omitempty"`
	Host       *string `json:"host,omitempty"`
	Data       *string `json:"data,omitempty"`
	Ttl        *int32  `json:"ttl,omitempty"`
	Weight     *int32  `json:"weight,omitempty"`
	Preference *int32  `json:"preference,omitempty"`
}

type UpdateRecordResponse struct {
	apiResponseBase
}

func (c *Client) UpdateRecord(req *UpdateRecordRequest) (*UpdateRecordResponse, error) {
	return c.UpdateRecordWithContext(context.Background(), req)
}

func (c *Client) UpdateRecordWithContext(ctx context.Context, req *UpdateRecordRequest) (*UpdateRecordResponse, error) {
	httpreq, err := c.newRequest(http.MethodPut, "/record")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &UpdateRecordResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
