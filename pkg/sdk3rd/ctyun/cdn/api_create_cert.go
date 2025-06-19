package cdn

import (
	"context"
	"net/http"
)

type CreateCertRequest struct {
	Name  *string `json:"name,omitempty"`
	Certs *string `json:"certs,omitempty"`
	Key   *string `json:"key,omitempty"`
}

type CreateCertResponse struct {
	apiResponseBase

	ReturnObj *struct {
		Id int64 `json:"id"`
	} `json:"returnObj,omitempty"`
}

func (c *Client) CreateCert(req *CreateCertRequest) (*CreateCertResponse, error) {
	return c.CreateCertWithContext(context.Background(), req)
}

func (c *Client) CreateCertWithContext(ctx context.Context, req *CreateCertRequest) (*CreateCertResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/v1/cert/creat-cert")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &CreateCertResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
