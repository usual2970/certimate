package apisix

import (
	"crypto/tls"
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
		SetBaseURL(strings.TrimRight(serverUrl, "/")+"/apisix/admin").
		SetHeader("User-Agent", "certimate").
		SetPreRequestHook(func(c *resty.Client, req *http.Request) error {
			req.Header.Set("X-API-KEY", apiKey)

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

func (c *Client) sendRequest(method string, path string, params interface{}) (*resty.Response, error) {
	req := c.client.R()
	if strings.EqualFold(method, http.MethodGet) {
		qs := make(map[string]string)
		if params != nil {
			temp := make(map[string]any)
			jsonb, _ := json.Marshal(params)
			json.Unmarshal(jsonb, &temp)
			for k, v := range temp {
				if v != nil {
					qs[k] = fmt.Sprintf("%v", v)
				}
			}
		}

		req = req.SetQueryParams(qs)
	} else {
		req = req.SetHeader("Content-Type", "application/json").SetBody(params)
	}

	resp, err := req.Execute(method, path)
	if err != nil {
		return resp, fmt.Errorf("apisix api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return resp, fmt.Errorf("apisix api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.String())
	}

	return resp, nil
}

func (c *Client) sendRequestWithResult(method string, path string, params interface{}, result interface{}) error {
	resp, err := c.sendRequest(method, path, params)
	if err != nil {
		if resp != nil {
			json.Unmarshal(resp.Body(), &result)
		}
		return err
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("apisix api error: failed to unmarshal response: %w", err)
	}

	return nil
}
