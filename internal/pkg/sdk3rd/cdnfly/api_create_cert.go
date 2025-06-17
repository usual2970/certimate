package cdnfly

import (
	"context"
	"net/http"
)

type CreateCertRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"des,omitempty"`
	Type        *string `json:"type,omitempty"`
	Cert        *string `json:"cert,omitempty"`
	Key         *string `json:"key,omitempty"`
}

type CreateCertResponse struct {
	apiResponseBase

	Data string `json:"data"`
}

func (c *Client) CreateCert(req *CreateCertRequest) (*CreateCertResponse, error) {
	return c.CreateCertWithContext(context.Background(), req)
}

func (c *Client) CreateCertWithContext(ctx context.Context, req *CreateCertRequest) (*CreateCertResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/certs")
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
