package qiniusdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth"
)

const qiniuHost = "https://api.qiniu.com"

type Client struct {
	mac *auth.Credentials
}

func NewClient(mac *auth.Credentials) *Client {
	if mac == nil {
		mac = auth.Default()
	}
	return &Client{mac: mac}
}

func (c *Client) GetDomainInfo(domain string) (*GetDomainInfoResponse, error) {
	respBytes, err := c.sendReq(http.MethodGet, fmt.Sprintf("domain/%s", domain), nil)
	if err != nil {
		return nil, err
	}

	resp := &GetDomainInfoResponse{}
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != nil && *resp.Code != 0 && *resp.Code != 200 {
		return nil, fmt.Errorf("code: %d, error: %s", *resp.Code, *resp.Error)
	}

	return resp, nil
}

func (c *Client) ModifyDomainHttpsConf(domain, certId string, forceHttps, http2Enable bool) (*ModifyDomainHttpsConfResponse, error) {
	req := &ModifyDomainHttpsConfRequest{
		DomainInfoHttpsData: DomainInfoHttpsData{
			CertID:      certId,
			ForceHttps:  forceHttps,
			Http2Enable: http2Enable,
		},
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	respBytes, err := c.sendReq(http.MethodPut, fmt.Sprintf("domain/%s/httpsconf", domain), bytes.NewReader(reqBytes))
	if err != nil {
		return nil, err
	}

	resp := &ModifyDomainHttpsConfResponse{}
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != nil && *resp.Code != 0 && *resp.Code != 200 {
		return nil, fmt.Errorf("code: %d, error: %s", *resp.Code, *resp.Error)
	}

	return resp, nil
}

func (c *Client) EnableDomainHttps(domain, certId string, forceHttps, http2Enable bool) (*EnableDomainHttpsResponse, error) {
	req := &EnableDomainHttpsRequest{
		DomainInfoHttpsData: DomainInfoHttpsData{
			CertID:      certId,
			ForceHttps:  forceHttps,
			Http2Enable: http2Enable,
		},
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	respBytes, err := c.sendReq(http.MethodPut, fmt.Sprintf("domain/%s/sslize", domain), bytes.NewReader(reqBytes))
	if err != nil {
		return nil, err
	}

	resp := &EnableDomainHttpsResponse{}
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != nil && *resp.Code != 0 && *resp.Code != 200 {
		return nil, fmt.Errorf("code: %d, error: %s", *resp.Code, *resp.Error)
	}

	return resp, nil
}

func (c *Client) UploadSslCert(name, commonName, certificate, privateKey string) (*UploadSslCertResponse, error) {
	req := &UploadSslCertRequest{
		Name:        name,
		CommonName:  commonName,
		Certificate: certificate,
		PrivateKey:  privateKey,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	respBytes, err := c.sendReq(http.MethodPost, "sslcert", bytes.NewReader(reqBytes))
	if err != nil {
		return nil, err
	}

	resp := &UploadSslCertResponse{}
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != nil && *resp.Code != 0 && *resp.Code != 200 {
		return nil, fmt.Errorf("qiniu api error, code: %d, error: %s", *resp.Code, *resp.Error)
	}

	return resp, nil
}

func (c *Client) sendReq(method string, path string, body io.Reader) ([]byte, error) {
	path = strings.TrimPrefix(path, "/")

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", qiniuHost, path), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	if err := c.mac.AddToken(auth.TokenQBox, req); err != nil {
		return nil, err
	}

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
