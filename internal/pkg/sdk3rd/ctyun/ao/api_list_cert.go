package ao

import (
	"context"
	"net/http"
	"strconv"
)

type ListCertRequest struct {
	Page      *int32 `json:"page,omitempty"`
	PerPage   *int32 `json:"per_page,omitempty"`
	UsageMode *int32 `json:"usage_mode,omitempty"`
}

type ListCertResponse struct {
	baseResult

	ReturnObj *struct {
		Results      []*CertRecord `json:"result,omitempty"`
		Page         int32         `json:"page,omitempty"`
		PerPage      int32         `json:"per_page,omitempty"`
		TotalPage    int32         `json:"total_page,omitempty"`
		TotalRecords int32         `json:"total_records,omitempty"`
	} `json:"returnObj,omitempty"`
}

func (c *Client) ListCert(req *ListCertRequest) (*ListCertResponse, error) {
	return c.ListCertWithContext(context.Background(), req)
}

func (c *Client) ListCertWithContext(ctx context.Context, req *ListCertRequest) (*ListCertResponse, error) {
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

	result := &ListCertResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
