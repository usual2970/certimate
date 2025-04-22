package console

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (c *Client) getCookie() error {
	req := &signinRequest{Username: c.username, Password: c.password}
	res, err := c.sendRequest(http.MethodPost, "/accounts/signin/", req)
	if err != nil {
		return err
	}

	resp := &signinResponse{}
	if err := json.Unmarshal(res.Body(), &resp); err != nil {
		return fmt.Errorf("upyun api error: failed to parse response: %w", err)
	} else if !resp.Data.Result {
		return errors.New("upyun console signin failed")
	}

	c.loginCookie = res.Header().Get("Set-Cookie")

	return nil
}

func (c *Client) UploadHttpsCertificate(req *UploadHttpsCertificateRequest) (*UploadHttpsCertificateResponse, error) {
	if c.loginCookie == "" {
		if err := c.getCookie(); err != nil {
			return nil, err
		}
	}

	resp := &UploadHttpsCertificateResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/api/https/certificate/", req, resp)
	return resp, err
}

func (c *Client) GetHttpsCertificateManager(certificateId string) (*GetHttpsCertificateManagerResponse, error) {
	if c.loginCookie == "" {
		if err := c.getCookie(); err != nil {
			return nil, err
		}
	}

	req := &GetHttpsCertificateManagerRequest{CertificateId: certificateId}
	resp := &GetHttpsCertificateManagerResponse{}
	err := c.sendRequestWithResult(http.MethodGet, "/api/https/certificate/manager/", req, resp)
	return resp, err
}

func (c *Client) UpdateHttpsCertificateManager(req *UpdateHttpsCertificateManagerRequest) (*UpdateHttpsCertificateManagerResponse, error) {
	if c.loginCookie == "" {
		if err := c.getCookie(); err != nil {
			return nil, err
		}
	}

	resp := &UpdateHttpsCertificateManagerResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/api/https/certificate/manager", req, resp)
	return resp, err
}

func (c *Client) GetHttpsServiceManager(domain string) (*GetHttpsServiceManagerResponse, error) {
	if c.loginCookie == "" {
		if err := c.getCookie(); err != nil {
			return nil, err
		}
	}

	req := &GetHttpsServiceManagerRequest{Domain: domain}
	resp := &GetHttpsServiceManagerResponse{}
	err := c.sendRequestWithResult(http.MethodGet, "/api/https/services/manager", req, resp)
	return resp, err
}

func (c *Client) MigrateHttpsDomain(req *MigrateHttpsDomainRequest) (*MigrateHttpsDomainResponse, error) {
	if c.loginCookie == "" {
		if err := c.getCookie(); err != nil {
			return nil, err
		}
	}

	resp := &MigrateHttpsDomainResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/api/https/migrate/domain", req, resp)
	return resp, err
}
