package safeline

import (
	"context"
	"net/http"
)

type UpdateCertificateRequest struct {
	Id     int32             `json:"id"`
	Type   int32             `json:"type"`
	Manual *CertificateManul `json:"manual"`
}

type UpdateCertificateResponse struct {
	apiResponseBase
}

func (c *Client) UpdateCertificate(req *UpdateCertificateRequest) (*UpdateCertificateResponse, error) {
	return c.UpdateCertificateWithContext(context.Background(), req)
}

func (c *Client) UpdateCertificateWithContext(ctx context.Context, req *UpdateCertificateRequest) (*UpdateCertificateResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/api/open/cert")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &UpdateCertificateResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
