package certificate

import (
	"context"
	"net/http"
)

type CreateCertificateRequest struct {
	Name        *string `json:"name,omitempty"`
	Certificate *string `json:"certificate,omitempty"`
	PrivateKey  *string `json:"privateKey,omitempty"`
	Comment     *string `json:"comment,omitempty" `
}

type CreateCertificateResponse struct {
	apiResponseBase

	CertificateLocation string `json:"location,omitempty"`
}

func (c *Client) CreateCertificate(req *CreateCertificateRequest) (*CreateCertificateResponse, error) {
	return c.CreateCertificateWithContext(context.Background(), req)
}

func (c *Client) CreateCertificateWithContext(ctx context.Context, req *CreateCertificateRequest) (*CreateCertificateResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/api/certificate")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &CreateCertificateResponse{}
	if httpresp, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	} else {
		result.CertificateLocation = httpresp.Header().Get("Location")
	}

	return result, nil
}
