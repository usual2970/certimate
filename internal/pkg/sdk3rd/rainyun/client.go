package rainyun

import (
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

func NewClient(apiKey string) *Client {
	client := resty.New().
		SetBaseURL("https://api.v2.rainyun.com").
		SetHeader("x-api-key", apiKey)

	return &Client{
		client: client,
	}
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.client.SetTimeout(timeout)
	return c
}

func (c *Client) sendRequest(method string, path string, params interface{}) (*resty.Response, error) {
	req := c.client.R()
	if strings.EqualFold(method, http.MethodGet) {
		if params != nil {
			jsonb, _ := json.Marshal(params)
			req = req.SetQueryParam("options", string(jsonb))
		}
	} else {
		req = req.SetHeader("Content-Type", "application/json").SetBody(params)
	}

	resp, err := req.Execute(method, path)
	if err != nil {
		return resp, fmt.Errorf("rainyun api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return resp, fmt.Errorf("rainyun api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.String())
	}

	return resp, nil
}

func (c *Client) sendRequestWithResult(method string, path string, params interface{}, result BaseResponse) error {
	resp, err := c.sendRequest(method, path, params)
	if err != nil {
		if resp != nil {
			json.Unmarshal(resp.Body(), &result)
		}
		return err
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("rainyun api error: failed to unmarshal response: %w", err)
	} else if errcode := result.GetCode(); errcode/100 != 2 {
		return fmt.Errorf("rainyun api error: code='%d', message='%s'", errcode, result.GetMessage())
	}

	return nil
}
