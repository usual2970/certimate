package qiniusdk

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/client"
)

const qiniuHost = "https://api.qiniu.com"

type Client struct {
	client *client.Client
}

func NewClient(mac *auth.Credentials) *Client {
	if mac == nil {
		mac = auth.Default()
	}

	client := client.DefaultClient
	client.Transport = newTransport(mac, nil)
	return &Client{client: &client}
}

func (c *Client) GetDomainInfo(domain string) (*GetDomainInfoResponse, error) {
	resp := new(GetDomainInfoResponse)
	if err := c.client.Call(context.Background(), resp, http.MethodGet, c.urlf("domain/%s", domain), nil); err != nil {
		return nil, err
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
	resp := new(ModifyDomainHttpsConfResponse)
	if err := c.client.CallWithJson(context.Background(), resp, http.MethodPut, c.urlf("domain/%s/httpsconf", domain), nil, req); err != nil {
		return nil, err
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
	resp := new(EnableDomainHttpsResponse)
	if err := c.client.CallWithJson(context.Background(), resp, http.MethodPut, c.urlf("domain/%s/sslize", domain), nil, req); err != nil {
		return nil, err
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
	resp := new(UploadSslCertResponse)
	if err := c.client.CallWithJson(context.Background(), resp, http.MethodPost, c.urlf("sslcert"), nil, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) urlf(pathf string, pathargs ...any) string {
	path := fmt.Sprintf(pathf, pathargs...)
	path = strings.TrimPrefix(path, "/")
	return qiniuHost + "/" + path
}
