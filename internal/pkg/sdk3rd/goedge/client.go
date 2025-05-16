package goedge

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	apiHost     string
	apiRole     string
	accessKeyId string
	accessKey   string

	accessToken    string
	accessTokenExp time.Time
	accessTokenMtx sync.Mutex

	client *resty.Client
}

func NewClient(apiHost, apiRole, accessKeyId, accessKey string) *Client {
	client := resty.New()

	return &Client{
		apiHost:     strings.TrimRight(apiHost, "/"),
		apiRole:     apiRole,
		accessKeyId: accessKeyId,
		accessKey:   accessKey,
		client:      client,
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
	req := c.client.R().SetBasicAuth(c.accessKeyId, c.accessKey)
	req.Method = method
	req.URL = c.apiHost + path
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
			SetHeader("X-Edge-Access-Token", c.accessToken)
	} else {
		req = req.
			SetHeader("Content-Type", "application/json").
			SetHeader("X-Edge-Access-Token", c.accessToken).
			SetBody(params)
	}

	resp, err := req.Send()
	if err != nil {
		return resp, fmt.Errorf("goedge api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return resp, fmt.Errorf("goedge api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.Body())
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
		return fmt.Errorf("goedge api error: failed to parse response: %w", err)
	} else if errcode := result.GetCode(); errcode != 200 {
		return fmt.Errorf("goedge api error: code='%d', message='%s'", errcode, result.GetMessage())
	}

	return nil
}
