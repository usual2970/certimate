package dns

import (
	"context"
	"net/http"
)

type DeleteRecordRequest struct {
	RecordId *int32 `json:"recordId,omitempty"`
}

type DeleteRecordResponse struct {
	apiResponseBase
}

func (c *Client) DeleteRecord(req *DeleteRecordRequest) (*DeleteRecordResponse, error) {
	return c.DeleteRecordWithContext(context.Background(), req)
}

func (c *Client) DeleteRecordWithContext(ctx context.Context, req *DeleteRecordRequest) (*DeleteRecordResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/v2/deleteRecord")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &DeleteRecordResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
