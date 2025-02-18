package baishansdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	apiToken string
	client   *resty.Client
}

func NewClient(apiToken string) *Client {
	client := resty.New()

	return &Client{
		apiToken: apiToken,
		client:   client,
	}
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.client.SetTimeout(timeout)
	return c
}

func (c *Client) sendRequest(method string, path string, params map[string]any) (*resty.Response, error) {
	req := c.client.R()
	req.Method = method
	req.URL = "https://cdn.api.baishan.com" + path
	if strings.EqualFold(method, http.MethodGet) {
		data := make(map[string]string)
		for k, v := range params {
			data[k] = fmt.Sprintf("%v", v)
		}
		req = req.
			SetQueryParams(data).
			SetQueryParam("token", c.apiToken)
	} else {
		req = req.
			SetHeader("Content-Type", "application/json").
			SetQueryParam("token", c.apiToken).
			SetBody(params)
	}

	resp, err := req.Send()
	if err != nil {
		return nil, fmt.Errorf("baishan api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("baishan api error: unexpected status code: %d, %s", resp.StatusCode(), resp.Body())
	}

	return resp, nil
}

func (c *Client) sendRequestWithResult(method string, path string, params map[string]any, result BaseResponse) error {
	resp, err := c.sendRequest(method, path, params)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("baishan api error: failed to parse response: %w", err)
	} else if errcode := result.GetCode(); errcode != 0 {
		return fmt.Errorf("baishan api error: %d - %s", errcode, result.GetMessage())
	}

	return nil
}
