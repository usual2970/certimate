package baishan

import (
	"context"
	"net/http"
)

type SetDomainCertificateRequest struct {
	CertificateId *string `json:"cert_id,omitempty"`
	Certificate   *string `json:"certificate,omitempty"`
	Key           *string `json:"key,omitempty"`
	Name          *string `json:"name,omitempty"`
}

type SetDomainCertificateResponse struct {
	apiResponseBase

	Data *DomainCertificate `json:"data,omitempty"`
}

func (c *Client) SetDomainCertificate(req *SetDomainCertificateRequest) (*SetDomainCertificateResponse, error) {
	return c.SetDomainCertificateWithContext(context.Background(), req)
}

func (c *Client) SetDomainCertificateWithContext(ctx context.Context, req *SetDomainCertificateRequest) (*SetDomainCertificateResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/v2/domain/certificate")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &SetDomainCertificateResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
