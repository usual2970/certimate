package cachefly

import (
	"context"
	"net/http"
)

type CreateCertificateRequest struct {
	Certificate    *string `json:"certificate,omitempty"`
	CertificateKey *string `json:"certificateKey,omitempty"`
	Password       *string `json:"password,omitempty"`
}

type CreateCertificateResponse struct {
	apiResponseBase

	Id                string   `json:"_id"`
	SubjectCommonName string   `json:"subjectCommonName"`
	SubjectNames      []string `json:"subjectNames"`
	Expired           bool     `json:"expired"`
	Expiring          bool     `json:"expiring"`
	InUse             bool     `json:"inUse"`
	Managed           bool     `json:"managed"`
	Services          []string `json:"services"`
	Domains           []string `json:"domains"`
	NotBefore         string   `json:"notBefore"`
	NotAfter          string   `json:"notAfter"`
	CreatedAt         string   `json:"createdAt"`
}

func (c *Client) CreateCertificate(req *CreateCertificateRequest) (*CreateCertificateResponse, error) {
	return c.CreateCertificateWithContext(context.Background(), req)
}

func (c *Client) CreateCertificateWithContext(ctx context.Context, req *CreateCertificateRequest) (*CreateCertificateResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/certificates")
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
