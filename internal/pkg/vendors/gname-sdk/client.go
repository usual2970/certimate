package gnamesdk

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/usual2970/certimate/internal/pkg/utils/maps"
)

type GnameClient struct {
	appId  string
	appKey string
	client *resty.Client
}

func NewGnameClient(appId, appKey string) *GnameClient {
	client := resty.New()

	return &GnameClient{
		appId:  appId,
		appKey: appKey,
		client: client,
	}
}

func (c *GnameClient) WithTimeout(timeout time.Duration) *GnameClient {
	c.client.SetTimeout(timeout)
	return c
}

func (c *GnameClient) AddDomainResolution(req *AddDomainResolutionRequest) (*AddDomainResolutionResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := AddDomainResolutionResponse{}
	err := c.sendRequestWithResult("/api/resolution/add", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *GnameClient) ModifyDomainResolution(req *ModifyDomainResolutionRequest) (*ModifyDomainResolutionResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := ModifyDomainResolutionResponse{}
	err := c.sendRequestWithResult("/api/resolution/edit", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *GnameClient) DeleteDomainResolution(req *DeleteDomainResolutionRequest) (*DeleteDomainResolutionResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := DeleteDomainResolutionResponse{}
	err := c.sendRequestWithResult("/api/resolution/delete", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *GnameClient) ListDomainResolution(req *ListDomainResolutionRequest) (*ListDomainResolutionResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := ListDomainResolutionResponse{}
	err := c.sendRequestWithResult("/api/resolution/list", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *GnameClient) generateSignature(params map[string]string) string {
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

func (c *GnameClient) sendRequest(path string, params map[string]any) (*resty.Response, error) {
	if params == nil {
		params = make(map[string]any)
	}

	data := make(map[string]string)
	for k, v := range params {
		data[k] = fmt.Sprintf("%v", v)
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
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("unexpected status code: %d, %s", resp.StatusCode(), resp.Body())
	}

	return resp, nil
}

func (c *GnameClient) sendRequestWithResult(path string, params map[string]any, result BaseResponse) error {
	resp, err := c.sendRequest(path, params)
	if err != nil {
		return err
	}

	jsonResp := make(map[string]any)
	if err := json.Unmarshal(resp.Body(), &jsonResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	if err := maps.Decode(jsonResp, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if result.GetCode() != 1 {
		return fmt.Errorf("API error: %s", result.GetMsg())
	}

	return nil
}
