package cdnflysdk

import (
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) GetSite(req *GetSiteRequest) (*GetSiteResponse, error) {
	resp := &GetSiteResponse{}
	err := c.sendRequestWithResult(http.MethodGet, fmt.Sprintf("/v1/sites/%s", url.PathEscape(req.Id)), req, resp)
	return resp, err
}

func (c *Client) UpdateSite(req *UpdateSiteRequest) (*UpdateSiteResponse, error) {
	resp := &UpdateSiteResponse{}
	err := c.sendRequestWithResult(http.MethodPut, fmt.Sprintf("/v1/sites/%s", url.PathEscape(req.Id)), req, resp)
	return resp, err
}

func (c *Client) CreateCertificate(req *CreateCertificateRequest) (*CreateCertificateResponse, error) {
	resp := &CreateCertificateResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/v1/certs", req, resp)
	return resp, err
}

func (c *Client) UpdateCertificate(req *UpdateCertificateRequest) (*UpdateCertificateResponse, error) {
	resp := &UpdateCertificateResponse{}
	err := c.sendRequestWithResult(http.MethodPut, fmt.Sprintf("/v1/certs/%s", url.PathEscape(req.Id)), req, resp)
	return resp, err
}
