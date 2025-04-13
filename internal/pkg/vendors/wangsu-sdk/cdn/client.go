package cdn

import (
	"time"

	"github.com/usual2970/certimate/internal/pkg/vendors/wangsu-sdk/openapi"
)

type Client struct {
	client *openapi.Client
}

func NewClient(accessKey, secretKey string) *Client {
	return &Client{client: openapi.NewClient(accessKey, secretKey)}
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.client.WithTimeout(timeout)
	return c
}
