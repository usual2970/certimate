package onepanelsdk

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
	apiHost string
	apiKey  string

	client *resty.Client
}

func NewClient(apiHost, apiKey string) *Client {
	client := resty.New()

	return &Client{
		apiHost: strings.TrimRight(apiHost, "/"),
		apiKey:  apiKey,
		client:  client,
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

func (c *Client) generateToken(timestamp string) string {
	tokenMd5 := md5.Sum([]byte("1panel" + c.apiKey + timestamp))
	tokenMd5Hex := hex.EncodeToString(tokenMd5[:])
	return tokenMd5Hex
}

func (c *Client) sendRequest(method string, path string, params interface{}) (*resty.Response, error) {
	req := c.client.R()
	req.Method = method
	req.URL = c.apiHost + "/api/v1" + path
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

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	token := c.generateToken(timestamp)
	req.SetHeader("1Panel-Timestamp", timestamp)
	req.SetHeader("1Panel-Token", token)

	resp, err := req.Send()
	if err != nil {
		return nil, fmt.Errorf("1panel api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("1panel api error: unexpected status code: %d, %s", resp.StatusCode(), resp.Body())
	}

	return resp, nil
}

func (c *Client) sendRequestWithResult(method string, path string, params interface{}, result BaseResponse) error {
	resp, err := c.sendRequest(method, path, params)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("1panel api error: failed to parse response: %w", err)
	} else if errcode := result.GetCode(); errcode/100 != 2 {
		return fmt.Errorf("1panel api error: %d - %s", errcode, result.GetMessage())
	}

	return nil
}
