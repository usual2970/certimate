package dnsla

import (
	"context"
	"fmt"
	"net/http"
)

type DeleteRecordResponse struct {
	apiResponseBase
}

func (c *Client) DeleteRecord(recordId string) (*DeleteRecordResponse, error) {
	return c.DeleteRecordWithContext(context.Background(), recordId)
}

func (c *Client) DeleteRecordWithContext(ctx context.Context, recordId string) (*DeleteRecordResponse, error) {
	if recordId == "" {
		return nil, fmt.Errorf("sdkerr: unset recordId")
	}

	httpreq, err := c.newRequest(http.MethodDelete, "/record")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetQueryParam("id", recordId)
		httpreq.SetContext(ctx)
	}

	result := &DeleteRecordResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
