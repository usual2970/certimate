package onepanelv2

import (
	"context"
	"net/http"
)

type SearchWebsiteSSLRequest struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"pageSize"`
}

type SearchWebsiteSSLResponse struct {
	apiResponseBase

	Data *struct {
		Items []*struct {
			ID          int64  `json:"id"`
			PEM         string `json:"pem"`
			PrivateKey  string `json:"privateKey"`
			Domains     string `json:"domains"`
			Description string `json:"description"`
			Status      string `json:"status"`
			UpdatedAt   string `json:"updatedAt"`
			CreatedAt   string `json:"createdAt"`
		} `json:"items"`
		Total int32 `json:"total"`
	} `json:"data,omitempty"`
}

func (c *Client) SearchWebsiteSSL(req *SearchWebsiteSSLRequest) (*SearchWebsiteSSLResponse, error) {
	return c.SearchWebsiteSSLWithContext(context.Background(), req)
}

func (c *Client) SearchWebsiteSSLWithContext(ctx context.Context, req *SearchWebsiteSSLRequest) (*SearchWebsiteSSLResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/websites/ssl/search")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &SearchWebsiteSSLResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
