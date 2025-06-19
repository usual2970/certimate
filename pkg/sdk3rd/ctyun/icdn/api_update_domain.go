package icdn

import (
	"context"
	"net/http"
)

type UpdateDomainRequest struct {
	Domain      *string `json:"domain,omitempty"`
	HttpsStatus *string `json:"https_status,omitempty"`
	CertName    *string `json:"cert_name,omitempty"`
}

type UpdateDomainResponse struct {
	apiResponseBase
}

func (c *Client) UpdateDomain(req *UpdateDomainRequest) (*UpdateDomainResponse, error) {
	return c.UpdateDomainWithContext(context.Background(), req)
}

func (c *Client) UpdateDomainWithContext(ctx context.Context, req *UpdateDomainRequest) (*UpdateDomainResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/v1/domain/update-domain")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &UpdateDomainResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
