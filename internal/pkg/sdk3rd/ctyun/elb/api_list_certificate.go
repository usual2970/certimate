package elb

import (
	"context"
	"net/http"
)

type ListCertificateRequest struct {
	ClientToken *string `json:"clientToken,omitempty"`
	RegionID    *string `json:"regionID,omitempty"`
	IDs         *string `json:"IDs,omitempty"`
	Name        *string `json:"name,omitempty"`
	Type        *string `json:"type,omitempty"`
}

type ListCertificateResponse struct {
	baseResult

	ReturnObj []*CertificateRecord `json:"returnObj,omitempty"`
}

func (c *Client) ListCertificate(req *ListCertificateRequest) (*ListCertificateResponse, error) {
	return c.ListCertificateWithContext(context.Background(), req)
}

func (c *Client) ListCertificateWithContext(ctx context.Context, req *ListCertificateRequest) (*ListCertificateResponse, error) {
	httpreq, err := c.newRequest(http.MethodGet, "/v4/elb/list-certificate")
	if err != nil {
		return nil, err
	} else {
		if req.ClientToken != nil {
			httpreq.SetQueryParam("clientToken", *req.ClientToken)
		}
		if req.RegionID != nil {
			httpreq.SetQueryParam("regionID", *req.RegionID)
		}
		if req.IDs != nil {
			httpreq.SetQueryParam("IDs", *req.IDs)
		}
		if req.Name != nil {
			httpreq.SetQueryParam("name", *req.Name)
		}
		if req.Type != nil {
			httpreq.SetQueryParam("type", *req.Type)
		}

		httpreq.SetContext(ctx)
	}

	result := &ListCertificateResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
