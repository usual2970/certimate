package ao

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/usual2970/certimate/internal/pkg/sdk3rd/ctyun/openapi"
)

const endpoint = "https://accessone-global.ctapi.ctyun.cn"

type Client struct {
	client *openapi.Client
}

func NewClient(accessKeyId, secretAccessKey string) (*Client, error) {
	client, err := openapi.NewClient(endpoint, accessKeyId, secretAccessKey)
	if err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.client.SetTimeout(timeout)
	return c
}

func (c *Client) newRequest(method string, path string) (*resty.Request, error) {
	return c.client.NewRequest(method, path)
}

func (c *Client) doRequest(request *resty.Request) (*resty.Response, error) {
	return c.client.DoRequest(request)
}

func (c *Client) doRequestWithResult(request *resty.Request, result baseResultInterface) (*resty.Response, error) {
	return c.client.DoRequestWithResult(request, result)
}
