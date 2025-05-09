package goedge

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (c *Client) getAccessToken() error {
	req := &getAPIAccessTokenRequest{
		Type:        c.apiUserType,
		AccessKeyId: c.accessKeyId,
		AccessKey:   c.accessKey,
	}
	res, err := c.sendRequest(http.MethodPost, "/APIAccessTokenService/getAPIAccessToken", req)
	if err != nil {
		return err
	}

	resp := &getAPIAccessTokenResponse{}
	if err := json.Unmarshal(res.Body(), &resp); err != nil {
		return fmt.Errorf("goedge api error: failed to parse response: %w", err)
	} else if resp.GetCode() != 200 {
		return fmt.Errorf("goedge get access token failed: code: %d, message: %s", resp.GetCode(), resp.GetMessage())
	}

	c.accessTokenMtx.Lock()
	c.accessToken = resp.Data.Token
	c.accessTokenExp = time.Unix(resp.Data.ExpiresAt, 0)
	c.accessTokenMtx.Unlock()

	return nil
}

func (c *Client) UpdateSSLCert(req *UpdateSSLCertRequest) (*UpdateSSLCertResponse, error) {
	if c.accessToken == "" || c.accessTokenExp.Before(time.Now()) {
		if err := c.getAccessToken(); err != nil {
			return nil, err
		}
	}

	resp := &UpdateSSLCertResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/SSLCertService/updateSSLCert", req, resp)
	return resp, err
}
