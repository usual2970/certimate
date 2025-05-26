package unicloud

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	username string
	password string

	serverlessJwtToken    string
	serverlessJwtTokenExp time.Time
	serverlessJwtTokenMtx sync.Mutex

	serverlessClient *resty.Client

	apiUserToken    string
	apiUserTokenMtx sync.Mutex

	apiClient *resty.Client
}

const (
	uniIdentityEndpoint     = "https://account.dcloud.net.cn/client"
	uniIdentityClientSecret = "ba461799-fde8-429f-8cc4-4b6d306e2339"
	uniIdentityAppId        = "__UNI__uniid_server"
	uniIdentitySpaceId      = "uni-id-server"
	uniConsoleEndpoint      = "https://unicloud.dcloud.net.cn/client"
	uniConsoleClientSecret  = "4c1f7fbf-c732-42b0-ab10-4634a8bbe834"
	uniConsoleAppId         = "__UNI__unicloud_console"
	uniConsoleSpaceId       = "dc-6nfabcn6ada8d3dd"
)

func NewClient(username, password string) *Client {
	client := &Client{
		username: username,
		password: password,
	}
	client.serverlessClient = resty.New()
	client.apiClient = resty.New().
		SetBaseURL("https://unicloud-api.dcloud.net.cn/unicloud/api").
		SetPreRequestHook(func(c *resty.Client, req *http.Request) error {
			if client.apiUserToken != "" {
				req.Header.Set("Token", client.apiUserToken)
			}

			return nil
		})

	return client
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.serverlessClient.SetTimeout(timeout)
	return c
}

func (c *Client) generateSignature(params map[string]any, secret string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	canonicalStr := ""
	for i, k := range keys {
		if i > 0 {
			canonicalStr += "&"
		}
		canonicalStr += k + "=" + fmt.Sprintf("%v", params[k])
	}

	mac := hmac.New(md5.New, []byte(secret))
	mac.Write([]byte(canonicalStr))
	sign := mac.Sum(nil)
	signHex := hex.EncodeToString(sign)

	return signHex
}

func (c *Client) buildServerlessClientInfo(appId string) (_clientInfo map[string]any, _err error) {
	return map[string]any{
		"PLATFORM":           "web",
		"OS":                 strings.ToUpper(runtime.GOOS),
		"APPID":              appId,
		"DEVICEID":           "certimate",
		"LOCALE":             "zh-Hans",
		"osName":             runtime.GOOS,
		"appId":              appId,
		"appName":            "uniCloud",
		"deviceId":           "certimate",
		"deviceType":         "pc",
		"uniPlatform":        "web",
		"uniCompilerVersion": "4.45",
		"uniRuntimeVersion":  "4.45",
	}, nil
}

func (c *Client) buildServerlessPayloadInfo(appId, spaceId, target, method, action string, params, data interface{}) (map[string]any, error) {
	clientInfo, err := c.buildServerlessClientInfo(appId)
	if err != nil {
		return nil, err
	}

	functionArgsParams := make([]any, 0)
	if params != nil {
		functionArgsParams = append(functionArgsParams, params)
	}

	functionArgs := map[string]any{
		"clientInfo": clientInfo,
		"uniIdToken": c.serverlessJwtToken,
	}
	if method != "" {
		functionArgs["method"] = method
		functionArgs["params"] = make([]any, 0)
	}
	if action != "" {
		type _obj struct{}
		functionArgs["action"] = action
		functionArgs["data"] = &_obj{}
	}
	if params != nil {
		functionArgs["params"] = []any{params}
	}
	if data != nil {
		functionArgs["data"] = data
	}

	jsonb, err := json.Marshal(map[string]any{
		"functionTarget": target,
		"functionArgs":   functionArgs,
	})
	if err != nil {
		return nil, err
	}

	payload := map[string]any{
		"method":    "serverless.function.runtime.invoke",
		"params":    string(jsonb),
		"spaceId":   spaceId,
		"timestamp": time.Now().UnixMilli(),
	}

	return payload, nil
}

func (c *Client) invokeServerless(endpoint, clientSecret, appId, spaceId, target, method, action string, params, data interface{}) (*resty.Response, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("unicloud api error: endpoint cannot be empty")
	}

	payload, err := c.buildServerlessPayloadInfo(appId, spaceId, target, method, action, params, data)
	if err != nil {
		return nil, fmt.Errorf("unicloud api error: failed to build request: %w", err)
	}

	clientInfo, _ := c.buildServerlessClientInfo(appId)
	clientInfoJsonb, _ := json.Marshal(clientInfo)

	sign := c.generateSignature(payload, clientSecret)

	req := c.serverlessClient.R().
		SetHeader("Origin", "https://unicloud.dcloud.net.cn").
		SetHeader("Referer", "https://unicloud.dcloud.net.cn").
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Client-Info", string(clientInfoJsonb)).
		SetHeader("X-Client-Token", c.serverlessJwtToken).
		SetHeader("X-Serverless-Sign", sign).
		SetBody(payload)
	resp, err := req.Post(endpoint)
	if err != nil {
		return resp, fmt.Errorf("unicloud api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return resp, fmt.Errorf("unicloud api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.String())
	}

	return resp, nil
}

func (c *Client) invokeServerlessWithResult(endpoint, clientSecret, appId, spaceId, target, method, action string, params, data interface{}, result BaseResponse) error {
	resp, err := c.invokeServerless(endpoint, clientSecret, appId, spaceId, target, method, action, params, data)
	if err != nil {
		if resp != nil {
			json.Unmarshal(resp.Body(), &result)
		}
		return err
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("unicloud api error: failed to unmarshal response: %w", err)
	} else if success := result.GetSuccess(); !success {
		return fmt.Errorf("unicloud api error: code='%s', message='%s'", result.GetErrorCode(), result.GetErrorMessage())
	}

	return nil
}

func (c *Client) sendRequest(method string, path string, params interface{}) (*resty.Response, error) {
	req := c.apiClient.R()
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
		req = req.SetHeader("Content-Type", "application/json").SetBody(params)
	}

	resp, err := req.Execute(method, path)
	if err != nil {
		return resp, fmt.Errorf("unicloud api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return resp, fmt.Errorf("unicloud api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.String())
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
		return fmt.Errorf("unicloud api error: failed to unmarshal response: %w", err)
	} else if retcode := result.GetReturnCode(); retcode != 0 {
		return fmt.Errorf("unicloud api error: ret='%d', desc='%s'", retcode, result.GetReturnDesc())
	}

	return nil
}
