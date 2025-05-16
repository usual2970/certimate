package cachefly

import (
	"net/http"
)

func (c *Client) CreateCertificate(req *CreateCertificateRequest) (*CreateCertificateResponse, error) {
	resp := &CreateCertificateResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/certificates", req, resp)
	return resp, err
}
