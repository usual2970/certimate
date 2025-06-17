package lvdn

import (
	"context"
	"net/http"
)

type QueryDomainDetailRequest struct {
	Domain      *string `json:"domain,omitempty"`
	ProductCode *string `json:"product_code,omitempty"`
}

type QueryDomainDetailResponse struct {
	apiResponseBase

	ReturnObj *struct {
		Domain      string `json:"domain"`
		ProductCode string `json:"product_code"`
		Status      int32  `json:"status"`
		AreaScope   int32  `json:"area_scope"`
		Cname       string `json:"cname"`
		HttpsSwitch int32  `json:"https_switch"`
		CertName    string `json:"cert_name"`
	} `json:"returnObj,omitempty"`
}

func (c *Client) QueryDomainDetail(req *QueryDomainDetailRequest) (*QueryDomainDetailResponse, error) {
	return c.QueryDomainDetailWithContext(context.Background(), req)
}

func (c *Client) QueryDomainDetailWithContext(ctx context.Context, req *QueryDomainDetailRequest) (*QueryDomainDetailResponse, error) {
	httpreq, err := c.newRequest(http.MethodGet, "/live/domain/query-domain-detail")
	if err != nil {
		return nil, err
	} else {
		if req.Domain != nil {
			httpreq.SetQueryParam("domain", *req.Domain)
		}
		if req.ProductCode != nil {
			httpreq.SetQueryParam("product_code", *req.ProductCode)
		}

		httpreq.SetContext(ctx)
	}

	result := &QueryDomainDetailResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
