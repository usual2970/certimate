package openapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type Client struct {
	client *resty.Client
}

func NewClient(endpoint, accessKeyId, secretAccessKey string) (*Client, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("sdk error: unset endpoint")
	}
	if _, err := url.Parse(endpoint); err != nil {
		return nil, fmt.Errorf("sdk error: invalid endpoint: %w", err)
	}
	if accessKeyId == "" {
		return nil, fmt.Errorf("sdk error: unset accessKey")
	}
	if secretAccessKey == "" {
		return nil, fmt.Errorf("sdk error: unset secretKey")
	}

	client := resty.New().
		SetBaseURL(endpoint).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "certimate").
		SetPreRequestHook(func(c *resty.Client, req *http.Request) error {
			// 生成时间戳及流水号
			now := time.Now()
			eopDate := now.Format("20060102T150405Z")
			eopReqId := uuid.New().String()

			// 获取查询参数
			queryStr := ""
			if req.URL != nil {
				queryStr = req.URL.Query().Encode()
			}

			// 获取请求正文
			payloadStr := ""
			if req.Body != nil {
				reader, err := req.GetBody()
				if err != nil {
					return err
				}

				defer reader.Close()
				payload, err := io.ReadAll(reader)
				if err != nil {
					return err
				}

				payloadStr = string(payload)
			}

			// 构造代签字符串
			payloadHash := sha256.Sum256([]byte(payloadStr))
			payloadHashHex := hex.EncodeToString(payloadHash[:])
			dataToSign := fmt.Sprintf("ctyun-eop-request-id:%s\neop-date:%s\n\n%s\n%s", eopReqId, eopDate, queryStr, payloadHashHex)

			// 生成 ktime
			hasher := hmac.New(sha256.New, []byte(secretAccessKey))
			hasher.Write([]byte(eopDate))
			ktime := hasher.Sum(nil)

			// 生成 kak
			hasher = hmac.New(sha256.New, ktime)
			hasher.Write([]byte(accessKeyId))
			kak := hasher.Sum(nil)

			// 生成 kdata
			hasher = hmac.New(sha256.New, kak)
			hasher.Write([]byte(now.Format("20060102")))
			kdate := hasher.Sum(nil)

			// 构造签名
			hasher = hmac.New(sha256.New, kdate)
			hasher.Write([]byte(dataToSign))
			sign := hasher.Sum(nil)
			signStr := base64.StdEncoding.EncodeToString(sign)

			// 设置请求头
			req.Header.Set("ctyun-eop-request-id", eopReqId)
			req.Header.Set("eop-date", eopDate)
			req.Header.Set("eop-authorization", fmt.Sprintf("%s Headers=ctyun-eop-request-id;eop-date Signature=%s", accessKeyId, signStr))

			return nil
		})

	return &Client{
		client: client,
	}, nil
}

func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.client.SetTimeout(timeout)
	return c
}

func (c *Client) NewRequest(method string, path string) (*resty.Request, error) {
	if method == "" {
		return nil, fmt.Errorf("sdk error: unset method")
	}
	if path == "" {
		return nil, fmt.Errorf("sdk error: unset path")
	}

	req := c.client.R()
	req.Method = method
	req.URL = path
	return req, nil
}

func (c *Client) DoRequest(request *resty.Request) (*resty.Response, error) {
	if request == nil {
		return nil, fmt.Errorf("sdk error: nil request")
	}

	// WARN:
	//   PLEASE DO NOT USE `req.SetResult` or `req.SetError` here.

	resp, err := request.Send()
	if err != nil {
		return resp, fmt.Errorf("sdk error: failed to send request: %w", err)
	} else if resp.IsError() {
		return resp, fmt.Errorf("sdk error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.String())
	}

	return resp, nil
}

func (c *Client) DoRequestWithResult(request *resty.Request, result any) (*resty.Response, error) {
	if request == nil {
		return nil, fmt.Errorf("sdk error: nil request")
	}

	response, err := c.DoRequest(request)
	if err != nil {
		if response != nil {
			json.Unmarshal(response.Body(), &result)
		}
		return response, err
	}

	if len(response.Body()) != 0 {
		if err := json.Unmarshal(response.Body(), &result); err != nil {
			return response, fmt.Errorf("sdk error: failed to unmarshal response: %w", err)
		}
	}

	return response, nil
}
