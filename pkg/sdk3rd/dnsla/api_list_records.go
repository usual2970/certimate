package dnsla

import (
	"context"
	"net/http"
	"strconv"
)

type ListRecordsRequest struct {
	DomainId  *string `json:"domainId,omitempty"`
	GroupId   *string `json:"groupId,omitempty"`
	LineId    *string `json:"lineId,omitempty"`
	Type      *int32  `json:"type,omitempty"`
	Host      *string `json:"host,omitempty"`
	Data      *string `json:"data,omitempty"`
	PageIndex *int32  `json:"pageIndex,omitempty"`
	PageSize  *int32  `json:"pageSize,omitempty"`
}

type ListRecordsResponse struct {
	apiResponseBase
	Data *struct {
		Total   int32        `json:"total"`
		Results []*DnsRecord `json:"results"`
	} `json:"data,omitempty"`
}

func (c *Client) ListRecords(req *ListRecordsRequest) (*ListRecordsResponse, error) {
	return c.ListRecordsWithContext(context.Background(), req)
}

func (c *Client) ListRecordsWithContext(ctx context.Context, req *ListRecordsRequest) (*ListRecordsResponse, error) {
	httpreq, err := c.newRequest(http.MethodGet, "/recordList")
	if err != nil {
		return nil, err
	} else {
		if req.DomainId != nil {
			httpreq.SetQueryParam("domainId", *req.DomainId)
		}
		if req.GroupId != nil {
			httpreq.SetQueryParam("groupId", *req.GroupId)
		}
		if req.LineId != nil {
			httpreq.SetQueryParam("lineId", *req.LineId)
		}
		if req.Type != nil {
			httpreq.SetQueryParam("type", strconv.Itoa(int(*req.Type)))
		}
		if req.Host != nil {
			httpreq.SetQueryParam("host", *req.Host)
		}
		if req.Data != nil {
			httpreq.SetQueryParam("data", *req.Data)
		}
		if req.PageIndex != nil {
			httpreq.SetQueryParam("pageIndex", strconv.Itoa(int(*req.PageIndex)))
		}
		if req.PageSize != nil {
			httpreq.SetQueryParam("pageSize", strconv.Itoa(int(*req.PageSize)))
		}

		httpreq.SetContext(ctx)
	}

	result := &ListRecordsResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
