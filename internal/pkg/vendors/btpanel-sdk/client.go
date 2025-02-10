package btpanelsdk

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/usual2970/certimate/internal/pkg/utils/maps"
)

type BaoTaPanelClient struct {
	apiHost string
	apiKey  string
	client  *resty.Client
}

func NewBaoTaPanelClient(apiHost, apiKey string) *BaoTaPanelClient {
	client := resty.New()

	return &BaoTaPanelClient{
		apiHost: apiHost,
		apiKey:  apiKey,
		client:  client,
	}
}

func (c *BaoTaPanelClient) WithTimeout(timeout time.Duration) *BaoTaPanelClient {
	c.client.SetTimeout(timeout)
	return c
}

func (c *BaoTaPanelClient) SetSiteSSL(req *SetSiteSSLRequest) (*SetSiteSSLResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := SetSiteSSLResponse{}
	err := c.sendRequestWithResult("/site?action=SetSSL", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *BaoTaPanelClient) generateSignature(timestamp string) string {
	keyMd5 := md5.Sum([]byte(c.apiKey))
	keyMd5Hex := strings.ToLower(hex.EncodeToString(keyMd5[:]))

	signMd5 := md5.Sum([]byte(timestamp + keyMd5Hex))
	signMd5Hex := strings.ToLower(hex.EncodeToString(signMd5[:]))
	return signMd5Hex
}

func (c *BaoTaPanelClient) sendRequest(path string, params map[string]any) (*resty.Response, error) {
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
		return nil, fmt.Errorf("baota: failed to send request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("baota: unexpected status code: %d, %s", resp.StatusCode(), resp.Body())
	}

	return resp, nil
}

func (c *BaoTaPanelClient) sendRequestWithResult(path string, params map[string]any, result BaseResponse) error {
	resp, err := c.sendRequest(path, params)
	if err != nil {
		return err
	}

	jsonResp := make(map[string]any)
	if err := json.Unmarshal(resp.Body(), &jsonResp); err != nil {
		return fmt.Errorf("baota: failed to parse response: %w", err)
	}
	if err := maps.Decode(jsonResp, &result); err != nil {
		return fmt.Errorf("baota: failed to parse response: %w", err)
	}

	if result.GetStatus() != nil && !*result.GetStatus() {
		if result.GetMsg() == nil {
			return fmt.Errorf("baota api error: unknown error")
		} else {
			return fmt.Errorf("baota api error: %s", result.GetMsg())
		}
	}

	return nil
}
