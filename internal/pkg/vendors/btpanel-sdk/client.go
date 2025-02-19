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

func (c *Client) generateSignature(timestamp string) string {
	keyMd5 := md5.Sum([]byte(c.apiKey))
	keyMd5Hex := strings.ToLower(hex.EncodeToString(keyMd5[:]))

	signMd5 := md5.Sum([]byte(timestamp + keyMd5Hex))
	signMd5Hex := strings.ToLower(hex.EncodeToString(signMd5[:]))
	return signMd5Hex
}

func (c *Client) sendRequest(path string, params interface{}) (*resty.Response, error) {
	timestamp := time.Now().Unix()

	data := make(map[string]any)
	if params != nil {
		temp := make(map[string]any)
		jsonData, _ := json.Marshal(params)
		json.Unmarshal(jsonData, &temp)
		for k, v := range temp {
			data[k] = v
		}
	}
	data["request_time"] = timestamp
	data["request_token"] = c.generateSignature(fmt.Sprintf("%d", timestamp))

	url := c.apiHost + path
	req := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data)
	resp, err := req.Post(url)
	if err != nil {
		return nil, fmt.Errorf("baota api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("baota api error: unexpected status code: %d, %s", resp.StatusCode(), resp.Body())
	}

	return resp, nil
}

func (c *Client) sendRequestWithResult(path string, params interface{}, result BaseResponse) error {
	resp, err := c.sendRequest(path, params)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("baota api error: failed to parse response: %w", err)
	} else if errstatus := result.GetStatus(); errstatus != nil && !*errstatus {
		if result.GetMessage() == nil {
			return fmt.Errorf("baota api error: unknown error")
		} else {
			return fmt.Errorf("baota api error: %s", *result.GetMessage())
		}
	}

	return nil
}
