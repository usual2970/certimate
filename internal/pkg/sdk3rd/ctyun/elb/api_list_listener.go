package elb

import (
	"context"
	"net/http"
)

type ListListenerRequest struct {
	ClientToken     *string `json:"clientToken,omitempty"`
	RegionID        *string `json:"regionID,omitempty"`
	ProjectID       *string `json:"projectID,omitempty"`
	IDs             *string `json:"IDs,omitempty"`
	Name            *string `json:"name,omitempty"`
	LoadBalancerID  *string `json:"loadBalancerID,omitempty"`
	AccessControlID *string `json:"accessControlID,omitempty"`
}

type ListListenerResponse struct {
	baseResult

	ReturnObj []*ListenerRecord `json:"returnObj,omitempty"`
}

func (c *Client) ListListener(req *ListListenerRequest) (*ListListenerResponse, error) {
	return c.ListListenerWithContext(context.Background(), req)
}

func (c *Client) ListListenerWithContext(ctx context.Context, req *ListListenerRequest) (*ListListenerResponse, error) {
	httpreq, err := c.newRequest(http.MethodGet, "/v4/elb/list-listener")
	if err != nil {
		return nil, err
	} else {
		if req.ClientToken != nil {
			httpreq.SetQueryParam("clientToken", *req.ClientToken)
		}
		if req.RegionID != nil {
			httpreq.SetQueryParam("regionID", *req.RegionID)
		}
		if req.ProjectID != nil {
			httpreq.SetQueryParam("projectID", *req.ProjectID)
		}
		if req.IDs != nil {
			httpreq.SetQueryParam("IDs", *req.IDs)
		}
		if req.Name != nil {
			httpreq.SetQueryParam("name", *req.Name)
		}
		if req.LoadBalancerID != nil {
			httpreq.SetQueryParam("loadBalancerID", *req.LoadBalancerID)
		}
		if req.LoadBalancerID != nil {
			httpreq.SetQueryParam("accessControlID", *req.AccessControlID)
		}

		httpreq.SetContext(ctx)
	}

	result := &ListListenerResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
