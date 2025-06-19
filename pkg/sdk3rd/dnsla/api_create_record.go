package dnsla

import (
	"context"
	"net/http"
)

type CreateRecordRequest struct {
	DomainId   string  `json:"domainId"`
	GroupId    *string `json:"groupId,omitempty"`
	LineId     *string `json:"lineId,omitempty"`
	Type       int32   `json:"type"`
	Host       string  `json:"host"`
	Data       string  `json:"data"`
	Ttl        int32   `json:"ttl"`
	Weight     *int32  `json:"weight,omitempty"`
	Preference *int32  `json:"preference,omitempty"`
}

type CreateRecordResponse struct {
	apiResponseBase
	Data *struct {
		Id string `json:"id"`
	} `json:"data,omitempty"`
}

func (c *Client) CreateRecord(req *CreateRecordRequest) (*CreateRecordResponse, error) {
	return c.CreateRecordWithContext(context.Background(), req)
}

func (c *Client) CreateRecordWithContext(ctx context.Context, req *CreateRecordRequest) (*CreateRecordResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/record")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &CreateRecordResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
