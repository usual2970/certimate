package onepanel

import (
	"context"
	"fmt"
	"net/http"
)

type GetHttpsConfResponse struct {
	apiResponseBase

	Data *struct {
		Enable      bool     `json:"enable"`
		HttpConfig  string   `json:"httpConfig"`
		SSLProtocol []string `json:"SSLProtocol"`
		Algorithm   string   `json:"algorithm"`
		Hsts        bool     `json:"hsts"`
	} `json:"data,omitempty"`
}

func (c *Client) GetHttpsConf(websiteId int64) (*GetHttpsConfResponse, error) {
	return c.GetHttpsConfWithContext(context.Background(), websiteId)
}

func (c *Client) GetHttpsConfWithContext(ctx context.Context, websiteId int64) (*GetHttpsConfResponse, error) {
	if websiteId == 0 {
		return nil, fmt.Errorf("sdkerr: unset websiteId")
	}

	httpreq, err := c.newRequest(http.MethodGet, fmt.Sprintf("/websites/%d/https", websiteId))
	if err != nil {
		return nil, err
	} else {
		httpreq.SetContext(ctx)
	}

	result := &GetHttpsConfResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
