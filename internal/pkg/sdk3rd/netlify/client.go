package netlify

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

	client *resty.Client
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

func (c *Client) sendRequest(method string, path string, queryParams interface{}, payloadParams interface{}) (*resty.Response, error) {
	req := c.client.R().SetHeader("Authorization", "Bearer "+c.apiToken)
	req.Method = method
	req.URL = "https://api.netlify.com/api/v1" + path

	if queryParams != nil {
		qs := make(map[string]string)
		temp := make(map[string]any)
		jsonb, _ := json.Marshal(queryParams)
		json.Unmarshal(jsonb, &temp)
		for k, v := range temp {
			if v != nil {
				qs[k] = fmt.Sprintf("%v", v)
			}
		}
		req = req.SetQueryParams(qs)
	}

	if strings.EqualFold(method, http.MethodGet) {
		qs := make(map[string]string)
		if payloadParams != nil {
			temp := make(map[string]any)
			jsonb, _ := json.Marshal(payloadParams)
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
			SetBody(payloadParams)
	}

	resp, err := req.Send()
	if err != nil {
		return resp, fmt.Errorf("netlify api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return resp, fmt.Errorf("netlify api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.Body())
	}

	return resp, nil
}

func (c *Client) sendRequestWithResult(method string, path string, queryParams interface{}, payloadParams interface{}, result BaseResponse) error {
	resp, err := c.sendRequest(method, path, queryParams, payloadParams)
	if err != nil {
		if resp != nil {
			json.Unmarshal(resp.Body(), &result)
		}
		return err
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("netlify api error: failed to parse response: %w", err)
	} else if errcode := result.GetCode(); errcode != 0 {
		return fmt.Errorf("netlify api error: code='%d', message='%s'", errcode, result.GetMessage())
	}

	return nil
}
