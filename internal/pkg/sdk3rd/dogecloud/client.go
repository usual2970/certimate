package dogecloud

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const dogeHost = "https://api.dogecloud.com"

type Client struct {
	accessKey string
	secretKey string
}

func NewClient(accessKey, secretKey string) *Client {
	return &Client{accessKey: accessKey, secretKey: secretKey}
}

func (c *Client) UploadCdnCert(note, cert, private string) (*UploadCdnCertResponse, error) {
	req := &UploadCdnCertRequest{
		Note:        note,
		Certificate: cert,
		PrivateKey:  private,
	}

	reqBts, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	reqMap := make(map[string]interface{})
	err = json.Unmarshal(reqBts, &reqMap)
	if err != nil {
		return nil, err
	}

	respBts, err := c.sendReq(http.MethodPost, "cdn/cert/upload.json", reqMap, true)
	if err != nil {
		return nil, err
	}

	resp := &UploadCdnCertResponse{}
	err = json.Unmarshal(respBts, resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != nil && *resp.Code != 0 && *resp.Code != 200 {
		return nil, fmt.Errorf("dogecloud api error, code: %d, msg: %s", *resp.Code, *resp.Message)
	}

	return resp, nil
}

func (c *Client) BindCdnCertWithDomain(certId int64, domain string) (*BindCdnCertResponse, error) {
	req := &BindCdnCertRequest{
		CertId: certId,
		Domain: &domain,
	}

	reqBts, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	reqMap := make(map[string]interface{})
	err = json.Unmarshal(reqBts, &reqMap)
	if err != nil {
		return nil, err
	}

	respBts, err := c.sendReq(http.MethodPost, "cdn/cert/bind.json", reqMap, true)
	if err != nil {
		return nil, err
	}

	resp := &BindCdnCertResponse{}
	err = json.Unmarshal(respBts, resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != nil && *resp.Code != 0 && *resp.Code != 200 {
		return nil, fmt.Errorf("dogecloud api error, code: %d, msg: %s", *resp.Code, *resp.Message)
	}

	return resp, nil
}

func (c *Client) BindCdnCertWithDomainId(certId int64, domainId int64) (*BindCdnCertResponse, error) {
	req := &BindCdnCertRequest{
		CertId:   certId,
		DomainId: &domainId,
	}

	reqBts, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	reqMap := make(map[string]interface{})
	err = json.Unmarshal(reqBts, &reqMap)
	if err != nil {
		return nil, err
	}

	respBts, err := c.sendReq(http.MethodPost, "cdn/cert/bind.json", reqMap, true)
	if err != nil {
		return nil, err
	}

	resp := &BindCdnCertResponse{}
	err = json.Unmarshal(respBts, resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != nil && *resp.Code != 0 && *resp.Code != 200 {
		return nil, fmt.Errorf("dogecloud api error, code: %d, msg: %s", *resp.Code, *resp.Message)
	}

	return resp, nil
}

// 调用多吉云的 API。
// https://docs.dogecloud.com/cdn/api-access-token?id=go
//
// 入参：
//   - method：GET 或 POST
//   - path：是调用的 API 接口地址，包含 URL 请求参数 QueryString，例如：/console/vfetch/add.json?url=xxx&a=1&b=2
//   - data：POST 的数据，对象，例如 {a: 1, b: 2}，传递此参数表示不是 GET 请求而是 POST 请求
//   - jsonMode：数据 data 是否以 JSON 格式请求，默认为 false 则使用表单形式（a=1&b=2）
func (c *Client) sendReq(method string, path string, data map[string]interface{}, jsonMode bool) ([]byte, error) {
	body := ""
	mime := ""
	if jsonMode {
		_body, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = string(_body)
		mime = "application/json"
	} else {
		values := url.Values{}
		for k, v := range data {
			values.Set(k, v.(string))
		}
		body = values.Encode()
		mime = "application/x-www-form-urlencoded"
	}

	path = strings.TrimPrefix(path, "/")
	signStr := "/" + path + "\n" + body
	hmacObj := hmac.New(sha1.New, []byte(c.secretKey))
	hmacObj.Write([]byte(signStr))
	sign := hex.EncodeToString(hmacObj.Sum(nil))
	auth := fmt.Sprintf("TOKEN %s:%s", c.accessKey, sign)

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", dogeHost, path), strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", mime)
	req.Header.Add("Authorization", auth)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return r, nil
}
