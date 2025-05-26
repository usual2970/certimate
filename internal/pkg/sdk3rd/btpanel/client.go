package btpanel

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	apiKey string

	client *resty.Client
}

func NewClient(serverUrl, apiKey string) *Client {
	client := resty.New().
		SetBaseURL(strings.TrimRight(serverUrl, "/"))

	return &Client{
		apiKey: apiKey,
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

func (c *Client) generateSignature(timestamp string) string {
	keyMd5 := md5.Sum([]byte(c.apiKey))
	keyMd5Hex := strings.ToLower(hex.EncodeToString(keyMd5[:]))

	signMd5 := md5.Sum([]byte(timestamp + keyMd5Hex))
	signMd5Hex := strings.ToLower(hex.EncodeToString(signMd5[:]))
	return signMd5Hex
}

func (c *Client) sendRequest(path string, params interface{}) (*resty.Response, error) {
	timestamp := time.Now().Unix()

	data := make(map[string]string)
	if params != nil {
		temp := make(map[string]any)
		jsonb, _ := json.Marshal(params)
		json.Unmarshal(jsonb, &temp)
		for k, v := range temp {
			if v != nil {
				switch reflect.Indirect(reflect.ValueOf(v)).Kind() {
				case reflect.String:
					data[k] = v.(string)
				case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
					data[k] = fmt.Sprintf("%v", v)
				default:
					if t, ok := v.(time.Time); ok {
						data[k] = t.Format(time.RFC3339)
					} else {
						jbytes, _ := json.Marshal(v)
						data[k] = string(jbytes)
					}
				}
			}
		}
	}
	data["request_time"] = fmt.Sprintf("%d", timestamp)
	data["request_token"] = c.generateSignature(fmt.Sprintf("%d", timestamp))

	req := c.client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(data)
	resp, err := req.Post(path)
	if err != nil {
		return resp, fmt.Errorf("baota api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return resp, fmt.Errorf("baota api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.String())
	}

	return resp, nil
}

func (c *Client) sendRequestWithResult(path string, params interface{}, result BaseResponse) error {
	resp, err := c.sendRequest(path, params)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("baota api error: failed to unmarshal response: %w", err)
	} else if errstatus := result.GetStatus(); errstatus != nil && !*errstatus {
		if result.GetMessage() == nil {
			return fmt.Errorf("baota api error: unknown error")
		} else {
			return fmt.Errorf("baota api error: message='%s'", *result.GetMessage())
		}
	}

	return nil
}
