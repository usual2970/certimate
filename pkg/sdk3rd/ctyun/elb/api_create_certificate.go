package elb

import (
	"context"
	"net/http"
)

type CreateCertificateRequest struct {
	ClientToken *string `json:"clientToken,omitempty"`
	RegionID    *string `json:"regionID,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Type        *string `json:"type,omitempty"`
	Certificate *string `json:"certificate,omitempty"`
	PrivateKey  *string `json:"privateKey,omitempty"`
}

type CreateCertificateResponse struct {
	apiResponseBase

	ReturnObj *struct {
		ID string `json:"id"`
	} `json:"returnObj,omitempty"`
}

func (c *Client) CreateCertificate(req *CreateCertificateRequest) (*CreateCertificateResponse, error) {
	return c.CreateCertificateWithContext(context.Background(), req)
}

func (c *Client) CreateCertificateWithContext(ctx context.Context, req *CreateCertificateRequest) (*CreateCertificateResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/v4/elb/create-certificate")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &CreateCertificateResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
