package unicloud

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
)

func (c *Client) ensureServerlessJwtTokenExists() error {
	c.serverlessJwtTokenMtx.Lock()
	defer c.serverlessJwtTokenMtx.Unlock()
	if c.serverlessJwtToken != "" && c.serverlessJwtTokenExp.After(time.Now()) {
		return nil
	}

	params := &loginParams{
		Password: c.password,
	}
	if regexp.MustCompile("^1\\d{10}$").MatchString(c.username) {
		params.Mobile = c.username
	} else if regexp.MustCompile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$").MatchString(c.username) {
		params.Email = c.username
	} else {
		params.Username = c.username
	}

	resp := &loginResponse{}
	if err := c.invokeServerlessWithResult(
		uniIdentityEndpoint, uniIdentityClientSecret, uniIdentityAppId, uniIdentitySpaceId,
		"uni-id-co", "login", "", params, nil,
		resp); err != nil {
		return err
	} else if resp.Data == nil || resp.Data.NewToken == nil || resp.Data.NewToken.Token == "" {
		return fmt.Errorf("unicloud api error: received empty token")
	}

	c.serverlessJwtToken = resp.Data.NewToken.Token
	c.serverlessJwtTokenExp = time.UnixMilli(resp.Data.NewToken.TokenExpired)

	return nil
}

func (c *Client) ensureApiUserTokenExists() error {
	if err := c.ensureServerlessJwtTokenExists(); err != nil {
		return err
	}

	c.apiUserTokenMtx.Lock()
	defer c.apiUserTokenMtx.Unlock()
	if c.apiUserToken != "" {
		return nil
	}

	resp := &getUserTokenResponse{}
	if err := c.invokeServerlessWithResult(
		uniConsoleEndpoint, uniConsoleClientSecret, uniConsoleAppId, uniConsoleSpaceId,
		"uni-cloud-kernel", "", "user/getUserToken", nil, map[string]any{"isLogin": true},
		resp); err != nil {
		return err
	} else if resp.Data == nil || resp.Data.Data == nil || resp.Data.Data.Data == nil || resp.Data.Data.Data.Token == "" {
		return fmt.Errorf("unicloud api error: received empty user token")
	}

	c.apiUserToken = resp.Data.Data.Data.Token

	return nil
}

func (c *Client) CreateDomainWithCert(req *CreateDomainWithCertRequest) (*CreateDomainWithCertResponse, error) {
	if err := c.ensureApiUserTokenExists(); err != nil {
		return nil, err
	}

	resp := &CreateDomainWithCertResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/host/create-domain-with-cert", req, resp)
	return resp, err
}
