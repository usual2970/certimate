package dnsla

import (
	"context"
	"net/http"
	"strconv"
)

type ListDomainsRequest struct {
	GroupId   *string `json:"groupId,omitempty"`
	PageIndex *int32  `json:"pageIndex,omitempty"`
	PageSize  *int32  `json:"pageSize,omitempty"`
}

type ListDomainsResponse struct {
	apiResponseBase
	Data *struct {
		Total   int32           `json:"total"`
		Results []*DomainRecord `json:"results"`
	} `json:"data,omitempty"`
}

func (c *Client) ListDomains(req *ListDomainsRequest) (*ListDomainsResponse, error) {
	return c.ListDomainsWithContext(context.Background(), req)
}

func (c *Client) ListDomainsWithContext(ctx context.Context, req *ListDomainsRequest) (*ListDomainsResponse, error) {
	httpreq, err := c.newRequest(http.MethodGet, "/domainList")
	if err != nil {
		return nil, err
	} else {
		if req.GroupId != nil {
			httpreq.SetQueryParam("groupId", *req.GroupId)
		}
		if req.PageIndex != nil {
			httpreq.SetQueryParam("pageIndex", strconv.Itoa(int(*req.PageIndex)))
		}
		if req.PageSize != nil {
			httpreq.SetQueryParam("pageSize", strconv.Itoa(int(*req.PageSize)))
		}

		httpreq.SetContext(ctx)
	}

	result := &ListDomainsResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
