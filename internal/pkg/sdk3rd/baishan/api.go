package baishan

import (
	"net/http"
)

func (c *Client) CreateCertificate(req *CreateCertificateRequest) (*CreateCertificateResponse, error) {
	resp := &CreateCertificateResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/v2/domain/certificate", req, resp)
	return resp, err
}

func (c *Client) GetDomainConfig(req *GetDomainConfigRequest) (*GetDomainConfigResponse, error) {
	resp := &GetDomainConfigResponse{}
	err := c.sendRequestWithResult(http.MethodGet, "/v2/domain/config", req, resp)
	return resp, err
}

func (c *Client) SetDomainConfig(req *SetDomainConfigRequest) (*SetDomainConfigResponse, error) {
	resp := &SetDomainConfigResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/v2/domain/config", req, resp)
	return resp, err
}
