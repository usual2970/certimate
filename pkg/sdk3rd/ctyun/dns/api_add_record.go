package dns

import (
	"context"
	"net/http"
)

type AddRecordRequest struct {
	Domain   *string `json:"domain,omitempty"`
	Host     *string `json:"host,omitempty"`
	Type     *string `json:"type,omitempty"`
	LineCode *string `json:"lineCode,omitempty"`
	Value    *string `json:"value,omitempty"`
	TTL      *int32  `json:"ttl,omitempty"`
	State    *int32  `json:"state,omitempty"`
	Remark   *string `json:"remark"`
}

type AddRecordResponse struct {
	apiResponseBase

	ReturnObj *struct {
		RecordId int32 `json:"recordId"`
	} `json:"returnObj,omitempty"`
}

func (c *Client) AddRecord(req *AddRecordRequest) (*AddRecordResponse, error) {
	return c.AddRecordWithContext(context.Background(), req)
}

func (c *Client) AddRecordWithContext(ctx context.Context, req *AddRecordRequest) (*AddRecordResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/v2/addRecord")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &AddRecordResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
