package dns

import (
	"context"
	"net/http"
)

type UpdateRecordRequest struct {
	RecordId *int32  `json:"recordId,omitempty"`
	Domain   *string `json:"domain,omitempty"`
	Host     *string `json:"host,omitempty"`
	Type     *string `json:"type,omitempty"`
	LineCode *string `json:"lineCode,omitempty"`
	Value    *string `json:"value,omitempty"`
	TTL      *int32  `json:"ttl,omitempty"`
	State    *int32  `json:"state,omitempty"`
	Remark   *string `json:"remark"`
}

type UpdateRecordResponse struct {
	apiResponseBase

	ReturnObj *struct {
		RecordId int32 `json:"recordId"`
	} `json:"returnObj,omitempty"`
}

func (c *Client) UpdateRecord(req *UpdateRecordRequest) (*UpdateRecordResponse, error) {
	return c.UpdateRecordWithContext(context.Background(), req)
}

func (c *Client) UpdateRecordWithContext(ctx context.Context, req *UpdateRecordRequest) (*UpdateRecordResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/v2/updateRecord")
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
