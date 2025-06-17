package icdn

import (
	"context"
	"net/http"
	"strconv"
)

type QueryCertDetailRequest struct {
	Id        *int64  `json:"id,omitempty"`
	Name      *string `json:"name,omitempty"`
	UsageMode *int32  `json:"usage_mode,omitempty"`
}

type QueryCertDetailResponse struct {
	apiResponseBase

	ReturnObj *struct {
		Result *CertDetail `json:"result,omitempty"`
	} `json:"returnObj,omitempty"`
}

func (c *Client) QueryCertDetail(req *QueryCertDetailRequest) (*QueryCertDetailResponse, error) {
	return c.QueryCertDetailWithContext(context.Background(), req)
}

func (c *Client) QueryCertDetailWithContext(ctx context.Context, req *QueryCertDetailRequest) (*QueryCertDetailResponse, error) {
	httpreq, err := c.newRequest(http.MethodGet, "/v1/cert/query-cert-detail")
	if err != nil {
		return nil, err
	} else {
		if req.Id != nil {
			httpreq.SetQueryParam("id", strconv.Itoa(int(*req.Id)))
		}
		if req.Name != nil {
			httpreq.SetQueryParam("name", *req.Name)
		}
		if req.UsageMode != nil {
			httpreq.SetQueryParam("usage_mode", strconv.Itoa(int(*req.UsageMode)))
		}

		httpreq.SetContext(ctx)
	}

	result := &QueryCertDetailResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
