package btwaf

import (
	"context"
	"net/http"
)

type GetSiteListRequest struct {
	Page     *int32  `json:"p,omitempty"`
	PageSize *int32  `json:"p_size,omitempty"`
	SiteName *string `json:"site_name,omitempty"`
}

type GetSiteListResponse struct {
	apiResponseBase

	Result *struct {
		List  []*SiteRecord `json:"list"`
		Total int32         `json:"total"`
	} `json:"res,omitempty"`
}

func (c *Client) GetSiteList(req *GetSiteListRequest) (*GetSiteListResponse, error) {
	return c.GetSiteListWithContext(context.Background(), req)
}

func (c *Client) GetSiteListWithContext(ctx context.Context, req *GetSiteListRequest) (*GetSiteListResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/wafmastersite/get_site_list")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &GetSiteListResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
