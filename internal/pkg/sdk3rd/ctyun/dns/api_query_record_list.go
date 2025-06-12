package dns

import (
	"context"
	"net/http"
)

type QueryRecordListRequest struct {
	Domain   *string `json:"domain,omitempty"`
	Host     *string `json:"host,omitempty"`
	Type     *string `json:"type,omitempty"`
	LineCode *string `json:"lineCode,omitempty"`
	Value    *string `json:"value,omitempty"`
	State    *int32  `json:"state,omitempty"`
}

type QueryRecordListResponse struct {
	baseResult

	ReturnObj *struct {
		Records []*struct {
			RecordId int32  `json:"recordId"`
			Host     string `json:"host"`
			Type     string `json:"type"`
			LineCode string `json:"lineCode"`
			Value    string `json:"value"`
			TTL      int32  `json:"ttl"`
			State    int32  `json:"state"`
			Remark   string `json:"remark"`
		} `json:"records,omitempty"`
	} `json:"returnObj,omitempty"`
}

func (c *Client) QueryRecordList(req *QueryRecordListRequest) (*QueryRecordListResponse, error) {
	return c.QueryRecordListWithContext(context.Background(), req)
}

func (c *Client) QueryRecordListWithContext(ctx context.Context, req *QueryRecordListRequest) (*QueryRecordListResponse, error) {
	request, err := c.newRequest(http.MethodGet, "/v2/queryRecordList")
	if err != nil {
		return nil, err
	} else {
		request.SetContext(ctx)
		request.SetBody(req)
	}

	result := &QueryRecordListResponse{}
	if _, err := c.doRequestWithResult(request, result); err != nil {
		return result, err
	}

	return result, nil
}
