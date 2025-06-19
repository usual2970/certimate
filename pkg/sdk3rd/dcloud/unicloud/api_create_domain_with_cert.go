package unicloud

import (
	"net/http"
)

type CreateDomainWithCertRequest struct {
	Provider string `json:"provider"`
	SpaceId  string `json:"spaceId"`
	Domain   string `json:"domain"`
	Cert     string `json:"cert"`
	Key      string `json:"key"`
}

type CreateDomainWithCertResponse struct {
	apiResponseBase
}

func (c *Client) CreateDomainWithCert(req *CreateDomainWithCertRequest) (*CreateDomainWithCertResponse, error) {
	if err := c.ensureApiUserTokenExists(); err != nil {
		return nil, err
	}

	resp := &CreateDomainWithCertResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/host/create-domain-with-cert", req, resp)
	return resp, err
}
