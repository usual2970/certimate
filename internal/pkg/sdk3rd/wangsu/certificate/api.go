package certificate

import (
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) ListCertificates() (*ListCertificatesResponse, error) {
	resp := &ListCertificatesResponse{}
	_, err := c.client.SendRequestWithResult(http.MethodGet, "/api/ssl/certificate", nil, resp)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (c *Client) CreateCertificate(req *CreateCertificateRequest) (*CreateCertificateResponse, error) {
	resp := &CreateCertificateResponse{}
	rres, err := c.client.SendRequestWithResult(http.MethodPost, "/api/certificate", req, resp)
	if err != nil {
		return resp, err
	}

	resp.CertificateUrl = rres.Header().Get("Location")
	return resp, err
}

func (c *Client) UpdateCertificate(certificateId string, req *UpdateCertificateRequest) (*UpdateCertificateResponse, error) {
	if certificateId == "" {
		return nil, fmt.Errorf("wangsu api error: invalid parameter: certificateId")
	}

	resp := &UpdateCertificateResponse{}
	_, err := c.client.SendRequestWithResult(http.MethodPut, fmt.Sprintf("/api/certificate/%s", url.PathEscape(certificateId)), req, resp)
	if err != nil {
		return resp, err
	}

	return resp, err
}
