package baishansdk

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	url := "https://cdn.api.baishan.com" + path

	req := c.client.R()
	if method == http.MethodGet {
		data := make(map[string]string)
		for k, v := range params {
			data[k] = fmt.Sprintf("%v", v)
		}
		req = req.
			SetQueryParams(data).
			SetQueryParam("token", c.apiToken)
	} else if method == http.MethodPost {
		req = req.
			SetHeader("Content-Type", "application/json").
			SetQueryParam("token", c.apiToken).
			SetBody(params)
	}

	var resp *resty.Response
	var err error
	if method == http.MethodGet {
		resp, err = req.Get(url)
	} else if method == http.MethodPost {
		resp, err = req.Post(url)
	} else {
		return nil, fmt.Errorf("baishan: unsupported method: %s", method)
	}

	if err != nil {
		return nil, fmt.Errorf("baishan: failed to send request: %w", err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("baishan: unexpected status code: %d, %s", resp.StatusCode(), resp.Body())
	}

	return resp, nil
}

func (c *Client) sendRequestWithResult(method string, path string, params map[string]any, result BaseResponse) error {
	resp, err := c.sendRequest(method, path, params)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("baishan: failed to parse response: %w", err)
	}

	if result.GetCode() != 0 {
		return fmt.Errorf("baishan api error: %d, %s", result.GetCode(), result.GetMessage())
	}

	return nil
}
