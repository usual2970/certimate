package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) ensureAccessTokenExists() error {
	c.accessTokenMtx.Lock()
	defer c.accessTokenMtx.Unlock()
	if c.accessToken != "" {
		return nil
	}

	req := &loginRequest{
		Email:    c.username,
		Username: c.username,
		Password: c.password,
	}
	res, err := c.sendRequest(http.MethodPost, "/login", req)
	if err != nil {
		return err
	}

	resp := &loginResponse{}
	if err := json.Unmarshal(res.Body(), &resp); err != nil {
		return fmt.Errorf("lecdn api error: failed to unmarshal response: %w", err)
	} else if resp.GetCode() != 200 {
		return fmt.Errorf("lecdn get token failed: code='%d', message='%s'", resp.GetCode(), resp.GetMessage())
	}

	c.accessToken = resp.Data.Token

	return nil
}

func (c *Client) UpdateCertificate(certId int64, req *UpdateCertificateRequest) (*UpdateCertificateResponse, error) {
	if certId == 0 {
		return nil, fmt.Errorf("lecdn api error: invalid parameter: CertId")
	}

	if err := c.ensureAccessTokenExists(); err != nil {
		return nil, err
	}

	resp := &UpdateCertificateResponse{}
	err := c.sendRequestWithResult(http.MethodPut, fmt.Sprintf("/certificate/%d", certId), req, resp)
	return resp, err
}
