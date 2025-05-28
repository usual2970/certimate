package btwaf

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
}

func NewClient(serverUrl, apiKey string) *Client {
	client := resty.New().
		SetBaseURL(strings.TrimRight(serverUrl, "/")+"/api").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "certimate").
		SetPreRequestHook(func(c *resty.Client, req *http.Request) error {
			timestamp := fmt.Sprintf("%d", time.Now().Unix())
			keyMd5 := md5.Sum([]byte(apiKey))
			keyMd5Hex := strings.ToLower(hex.EncodeToString(keyMd5[:]))
			signMd5 := md5.Sum([]byte(timestamp + keyMd5Hex))
			signMd5Hex := strings.ToLower(hex.EncodeToString(signMd5[:]))
			req.Header.Set("waf_request_time", timestamp)
			req.Header.Set("waf_request_token", signMd5Hex)

			return nil
		})

	return &Client{
		client: client,
	}
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.client.SetTimeout(timeout)
	return c
}

func (c *Client) WithTLSConfig(config *tls.Config) *Client {
	c.client.SetTLSClientConfig(config)
	return c
}

func (c *Client) sendRequest(path string, params interface{}) (*resty.Response, error) {
	req := c.client.R().SetBody(params)
	resp, err := req.Post(path)
	if err != nil {
		return resp, fmt.Errorf("baota api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return resp, fmt.Errorf("baota api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.String())
	}

	return resp, nil
}

func (c *Client) sendRequestWithResult(path string, params interface{}, result BaseResponse) error {
	resp, err := c.sendRequest(path, params)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("baota api error: failed to unmarshal response: %w", err)
	} else if errcode := result.GetCode(); errcode != 0 {
		return fmt.Errorf("baota api error: code='%d'", errcode)
	}

	return nil
}
