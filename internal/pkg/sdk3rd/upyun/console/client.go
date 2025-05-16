package console

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	username string
	password string

	loginCookie string

	client *resty.Client
}

func NewClient(username, password string) *Client {
	client := resty.New()

	return &Client{
		username: username,
		password: password,
		client:   client,
	}
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.client.SetTimeout(timeout)
	return c
}

func (c *Client) sendRequest(method string, path string, params interface{}) (*resty.Response, error) {
	req := c.client.R().SetBasicAuth(c.username, c.password)
	req.Method = method
	req.URL = "https://console.upyun.com" + path
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

		req = req.
			SetQueryParams(qs).
			SetHeader("Cookie", c.loginCookie)
	} else {
		req = req.
			SetHeader("Content-Type", "application/json").
			SetHeader("Cookie", c.loginCookie).
			SetBody(params)
	}

	resp, err := req.Send()
	if err != nil {
		return resp, fmt.Errorf("upyun api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return resp, fmt.Errorf("upyun api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.Body())
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
		return fmt.Errorf("upyun api error: failed to parse response: %w", err)
	}

	tresp := &baseResponse{}
	if err := json.Unmarshal(resp.Body(), &tresp); err != nil {
		return fmt.Errorf("upyun api error: failed to parse response: %w", err)
	} else if tdata := tresp.GetData(); tdata == nil {
		return fmt.Errorf("upyun api error: empty data")
	} else if errcode := tdata.GetErrorCode(); errcode > 0 {
		return fmt.Errorf("upyun api error: code='%d', message='%s'", errcode, tdata.GetErrorMessage())
	}

	return nil
}
