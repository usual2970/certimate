package ao

import (
	"context"
	"net/http"
	"strconv"
)

type ListCertsRequest struct {
	Page      *int32 `json:"page,omitempty"`
	PerPage   *int32 `json:"per_page,omitempty"`
	UsageMode *int32 `json:"usage_mode,omitempty"`
}

type ListCertsResponse struct {
	apiResponseBase

	ReturnObj *struct {
		Results      []*CertRecord `json:"result,omitempty"`
		Page         int32         `json:"page,omitempty"`
		PerPage      int32         `json:"per_page,omitempty"`
		TotalPage    int32         `json:"total_page,omitempty"`
		TotalRecords int32         `json:"total_records,omitempty"`
	} `json:"returnObj,omitempty"`
}

func (c *Client) ListCerts(req *ListCertsRequest) (*ListCertsResponse, error) {
	return c.ListCertsWithContext(context.Background(), req)
}

func (c *Client) ListCertsWithContext(ctx context.Context, req *ListCertsRequest) (*ListCertsResponse, error) {
	httpreq, err := c.newRequest(http.MethodGet, "/ctapi/v1/accessone/cert/list")
	if err != nil {
		return nil, err
	} else {
		if req.Page != nil {
			httpreq.SetQueryParam("page", strconv.Itoa(int(*req.Page)))
		}
		if req.PerPage != nil {
			httpreq.SetQueryParam("per_page", strconv.Itoa(int(*req.PerPage)))
		}
		if req.UsageMode != nil {
			httpreq.SetQueryParam("usage_mode", strconv.Itoa(int(*req.UsageMode)))
		}

		httpreq.SetContext(ctx)
	}

	result := &ListCertsResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
