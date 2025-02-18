package btpanelsdk

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	apiHost string
	apiKey  string
	client  *resty.Client
}

func NewClient(apiHost, apiKey string) *Client {
	client := resty.New()

	return &Client{
		apiHost: apiHost,
		apiKey:  apiKey,
		client:  client,
	}
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.client.SetTimeout(timeout)
	return c
}

func (c *Client) generateSignature(timestamp string) string {
	keyMd5 := md5.Sum([]byte(c.apiKey))
	keyMd5Hex := strings.ToLower(hex.EncodeToString(keyMd5[:]))

	signMd5 := md5.Sum([]byte(timestamp + keyMd5Hex))
	signMd5Hex := strings.ToLower(hex.EncodeToString(signMd5[:]))
	return signMd5Hex
}

func (c *Client) sendRequest(path string, params map[string]any) (*resty.Response, error) {
	if params == nil {
		params = make(map[string]any)
	}

	timestamp := time.Now().Unix()
	params["request_time"] = timestamp
	params["request_token"] = c.generateSignature(fmt.Sprintf("%d", timestamp))

	url := strings.TrimRight(c.apiHost, "/") + path
	req := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(params)
	resp, err := req.Post(url)
	if err != nil {
		return nil, fmt.Errorf("baota api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("baota api error: unexpected status code: %d, %s", resp.StatusCode(), resp.Body())
	}

	return resp, nil
}

func (c *Client) sendRequestWithResult(path string, params map[string]any, result BaseResponse) error {
	resp, err := c.sendRequest(path, params)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("baota api error: failed to parse response: %w", err)
	} else if result.GetStatus() != nil && !*result.GetStatus() {
		if result.GetMsg() == nil {
			return fmt.Errorf("baota api error: unknown error")
		} else {
			return fmt.Errorf("baota api error: %s", *result.GetMsg())
		}
	}

	return nil
}
