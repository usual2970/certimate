package ao

import (
	"context"
	"net/http"
	"strconv"
)

type QueryCertRequest struct {
	Id        *int64  `json:"id,omitempty"`
	Name      *string `json:"name,omitempty"`
	UsageMode *int32  `json:"usage_mode,omitempty"`
}

type QueryCertResponse struct {
	apiResponseBase

	ReturnObj *struct {
		Result *CertDetail `json:"result,omitempty"`
	} `json:"returnObj,omitempty"`
}

func (c *Client) QueryCert(req *QueryCertRequest) (*QueryCertResponse, error) {
	return c.QueryCertWithContext(context.Background(), req)
}

func (c *Client) QueryCertWithContext(ctx context.Context, req *QueryCertRequest) (*QueryCertResponse, error) {
	httpreq, err := c.newRequest(http.MethodGet, "/ctapi/v1/accessone/cert/query")
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

	result := &QueryCertResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
