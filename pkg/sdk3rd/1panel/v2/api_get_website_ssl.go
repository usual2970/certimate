package onepanelv2

import (
	"context"
	"fmt"
	"net/http"
)

type GetWebsiteSSLResponse struct {
	apiResponseBase

	Data *struct {
		ID            int64  `json:"id"`
		Provider      string `json:"provider"`
		Description   string `json:"description"`
		PrimaryDomain string `json:"primaryDomain"`
		Domains       string `json:"domains"`
		Type          string `json:"type"`
		Organization  string `json:"organization"`
		Status        string `json:"status"`
		StartDate     string `json:"startDate"`
		ExpireDate    string `json:"expireDate"`
		CreatedAt     string `json:"createdAt"`
		UpdatedAt     string `json:"updatedAt"`
	} `json:"data,omitempty"`
}

func (c *Client) GetWebsiteSSL(sslId int64) (*GetWebsiteSSLResponse, error) {
	return c.GetWebsiteSSLWithContext(context.Background(), sslId)
}

func (c *Client) GetWebsiteSSLWithContext(ctx context.Context, sslId int64) (*GetWebsiteSSLResponse, error) {
	if sslId == 0 {
		return nil, fmt.Errorf("sdkerr: unset sslId")
	}

	httpreq, err := c.newRequest(http.MethodGet, fmt.Sprintf("/websites/ssl/%d", sslId))
	if err != nil {
		return nil, err
	} else {
		httpreq.SetContext(ctx)
	}

	result := &GetWebsiteSSLResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
