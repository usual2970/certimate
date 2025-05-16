package safeline

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	apiHost  string
	apiToken string

	client *resty.Client
}

func NewClient(apiHost, apiToken string) *Client {
	client := resty.New()

	return &Client{
		apiHost:  strings.TrimRight(apiHost, "/"),
		apiToken: apiToken,
		client:   client,
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
	url := c.apiHost + path
	req := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-SLCE-API-TOKEN", c.apiToken).
		SetBody(params)
	resp, err := req.Post(url)
	if err != nil {
		return resp, fmt.Errorf("safeline api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return resp, fmt.Errorf("safeline api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.Body())
	}

	return resp, nil
}

func (c *Client) sendRequestWithResult(path string, params interface{}, result BaseResponse) error {
	resp, err := c.sendRequest(path, params)
	if err != nil {
		if resp != nil {
			json.Unmarshal(resp.Body(), &result)
		}
		return err
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("safeline api error: failed to parse response: %w", err)
	} else if errcode := result.GetErrCode(); errcode != nil && *errcode != "" {
		if result.GetErrMsg() == nil {
			return fmt.Errorf("safeline api error: code='%s'", *errcode)
		} else {
			return fmt.Errorf("safeline api error: code='%s', message='%s'", *errcode, *result.GetErrMsg())
		}
	}

	return nil
}
