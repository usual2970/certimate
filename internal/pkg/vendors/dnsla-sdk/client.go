package dnslasdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	apiId     string
	apiSecret string

	client *resty.Client
}

func NewClient(apiId, apiSecret string) *Client {
	client := resty.New()

	return &Client{
		apiId:     apiId,
		apiSecret: apiSecret,
		client:    client,
	}
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.client.SetTimeout(timeout)
	return c
}

func (c *Client) sendRequest(method string, path string, params interface{}) (*resty.Response, error) {
	req := c.client.R().SetBasicAuth(c.apiId, c.apiSecret)
	req.Method = method
	req.URL = "https://api.dns.la/api" + path
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
		req = req.
			SetHeader("Content-Type", "application/json").
			SetBody(params)
	}

	resp, err := req.Send()
	if err != nil {
		return resp, fmt.Errorf("dnsla api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return resp, fmt.Errorf("dnsla api error: unexpected status code: %d, %s", resp.StatusCode(), resp.Body())
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
		return fmt.Errorf("dnsla api error: failed to parse response: %w", err)
	} else if errcode := result.GetCode(); errcode/100 != 2 {
		return fmt.Errorf("dnsla api error: %d - %s", errcode, result.GetMessage())
	}

	return nil
}
