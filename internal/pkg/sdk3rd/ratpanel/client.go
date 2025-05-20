package ratpanelsdk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
}

func NewClient(apiHost string, accessTokenId int32, accessToken string) *Client {
	client := resty.New().
		SetBaseURL(strings.TrimRight(apiHost, "/")+"/api").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetPreRequestHook(func(c *resty.Client, req *http.Request) error {
			var body []byte
			var err error

			if req.Body != nil {
				body, err = io.ReadAll(req.Body)
				if err != nil {
					return err
				}
				req.Body = io.NopCloser(bytes.NewReader(body))
			}

			canonicalPath := req.URL.Path
			if !strings.HasPrefix(canonicalPath, "/api") {
				index := strings.Index(canonicalPath, "/api")
				if index != -1 {
					canonicalPath = canonicalPath[index:]
				}
			}

			canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s",
				req.Method,
				canonicalPath,
				req.URL.Query().Encode(),
				sha256Sum(string(body)))

			timestamp := time.Now().Unix()
			req.Header.Set("X-Timestamp", fmt.Sprintf("%d", timestamp))

			stringToSign := fmt.Sprintf("%s\n%d\n%s",
				"HMAC-SHA256",
				timestamp,
				sha256Sum(canonicalRequest))
			signature := hmacSha256(stringToSign, accessToken)
			req.Header.Set("Authorization", fmt.Sprintf("HMAC-SHA256 Credential=%d, Signature=%s", accessTokenId, signature))

			return nil
		})

	return &Client{
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

func (c *Client) sendRequest(method string, path string, params interface{}) (*resty.Response, error) {
	req := c.client.R()
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
		return resp, fmt.Errorf("ratpanel api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return resp, fmt.Errorf("ratpanel api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.Body())
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

	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("ratpanel api error: failed to unmarshal response: %w", err)
	} else if errmsg := result.GetMessage(); errmsg != "success" {
		return fmt.Errorf("ratpanel api error: message='%s'", errmsg)
	}

	return nil
}

func sha256Sum(str string) string {
	sum := sha256.Sum256([]byte(str))
	dst := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(dst, sum[:])
	return string(dst)
}

func hmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
