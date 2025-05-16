package gname

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	appId  string
	appKey string

	client *resty.Client
}

func NewClient(appId, appKey string) *Client {
	client := resty.New()

	return &Client{
		appId:  appId,
		appKey: appKey,
		client: client,
	}
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.client.SetTimeout(timeout)
	return c
}

func (c *Client) generateSignature(params map[string]string) string {
	// Step 1: Sort parameters by ASCII order
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Step 2: Create string A with URL-encoded values
	var pairs []string
	for _, k := range keys {
		encodedValue := url.QueryEscape(params[k])
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, encodedValue))
	}
	stringA := strings.Join(pairs, "&")

	// Step 3: Append appkey to create string B
	stringB := stringA + c.appKey

	// Step 4: Calculate MD5 and convert to uppercase
	hash := md5.Sum([]byte(stringB))
	return strings.ToUpper(fmt.Sprintf("%x", hash))
}

func (c *Client) sendRequest(path string, params interface{}) (*resty.Response, error) {
	data := make(map[string]string)
	if params != nil {
		temp := make(map[string]any)
		jsonb, _ := json.Marshal(params)
		json.Unmarshal(jsonb, &temp)
		for k, v := range temp {
			if v != nil {
				data[k] = fmt.Sprintf("%v", v)
			}
		}
	}
	data["appid"] = c.appId
	data["gntime"] = fmt.Sprintf("%d", time.Now().Unix())
	data["gntoken"] = c.generateSignature(data)

	url := "http://api.gname.com" + path
	req := c.client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(data)
	resp, err := req.Post(url)
	if err != nil {
		return resp, fmt.Errorf("gname api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return resp, fmt.Errorf("gname api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.String())
	}

	return resp, nil
}

func (c *Client) sendRequestWithResult(path string, params interface{}, result BaseResponse) error {
	resp, err := c.sendRequest(path, params)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("gname api error: failed to unmarshal response: %w", err)
	} else if errcode := result.GetCode(); errcode != 1 {
		return fmt.Errorf("gname api error: code='%d', message='%s'", errcode, result.GetMessage())
	}

	return nil
}
