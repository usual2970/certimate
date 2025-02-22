package baishansdk

import (
	"net/http"
)

func (c *Client) CreateCertificate(req *CreateCertificateRequest) (*CreateCertificateResponse, error) {
	resp := CreateCertificateResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/v2/domain/certificate", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetDomainConfig(req *GetDomainConfigRequest) (*GetDomainConfigResponse, error) {
	resp := GetDomainConfigResponse{}
	err := c.sendRequestWithResult(http.MethodGet, "/v2/domain/config", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) SetDomainConfig(req *SetDomainConfigRequest) (*SetDomainConfigResponse, error) {
	resp := SetDomainConfigResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/v2/domain/config", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
