package goedge

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (c *Client) ensureAccessTokenExists() error {
	c.accessTokenMtx.Lock()
	defer c.accessTokenMtx.Unlock()
	if c.accessToken != "" && c.accessTokenExp.After(time.Now()) {
		return nil
	}

	req := &getAPIAccessTokenRequest{
		Type:        c.apiRole,
		AccessKeyId: c.accessKeyId,
		AccessKey:   c.accessKey,
	}
	res, err := c.sendRequest(http.MethodPost, "/APIAccessTokenService/getAPIAccessToken", req)
	if err != nil {
		return err
	}

	resp := &getAPIAccessTokenResponse{}
	if err := json.Unmarshal(res.Body(), &resp); err != nil {
		return fmt.Errorf("goedge api error: failed to unmarshal response: %w", err)
	} else if resp.GetCode() != 200 {
		return fmt.Errorf("goedge get access token failed: code='%d', message='%s'", resp.GetCode(), resp.GetMessage())
	}

	c.accessToken = resp.Data.Token
	c.accessTokenExp = time.Unix(resp.Data.ExpiresAt, 0)

	return nil
}

func (c *Client) UpdateSSLCert(req *UpdateSSLCertRequest) (*UpdateSSLCertResponse, error) {
	if err := c.ensureAccessTokenExists(); err != nil {
		return nil, err
	}

	resp := &UpdateSSLCertResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/SSLCertService/updateSSLCert", req, resp)
	return resp, err
}
