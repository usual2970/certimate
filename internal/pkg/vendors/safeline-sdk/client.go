package safelinesdk

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	apiHost  string
	apiToken string
	client   *resty.Client
}

func NewClient(apiHost, apiToken string) *Client {
	client := resty.New()

	return &Client{
		apiHost:  apiHost,
		apiToken: apiToken,
		client:   client,
	}
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.client.SetTimeout(timeout)
	return c
}

func (c *Client) sendRequest(path string, params map[string]any) (*resty.Response, error) {
	if params == nil {
		params = make(map[string]any)
	}

	url := strings.TrimRight(c.apiHost, "/") + path
	req := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-SLCE-API-TOKEN", c.apiToken).
		SetBody(params)
	resp, err := req.Post(url)
	if err != nil {
		return nil, fmt.Errorf("safeline: failed to send request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("safeline: unexpected status code: %d, %s", resp.StatusCode(), resp.Body())
	}

	return resp, nil
}

func (c *Client) sendRequestWithResult(path string, params map[string]any, result BaseResponse) error {
	resp, err := c.sendRequest(path, params)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("safeline: failed to parse response: %w", err)
	}

	if result.GetErrCode() != nil && *result.GetErrCode() != "" {
		if result.GetErrMsg() == nil {
			return fmt.Errorf("safeline api error: %s", *result.GetErrCode())
		} else {
			return fmt.Errorf("safeline api error: %s, %s", *result.GetErrCode(), *result.GetErrMsg())
		}
	}

	return nil
}
