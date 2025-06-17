package elb

import (
	"context"
	"net/http"
)

type ShowListenerRequest struct {
	ClientToken *string `json:"clientToken,omitempty"`
	RegionID    *string `json:"regionID,omitempty"`
	ListenerID  *string `json:"listenerID,omitempty"`
}

type ShowListenerResponse struct {
	apiResponseBase

	ReturnObj []*ListenerRecord `json:"returnObj,omitempty"`
}

func (c *Client) ShowListener(req *ShowListenerRequest) (*ShowListenerResponse, error) {
	return c.ShowListenerWithContext(context.Background(), req)
}

func (c *Client) ShowListenerWithContext(ctx context.Context, req *ShowListenerRequest) (*ShowListenerResponse, error) {
	httpreq, err := c.newRequest(http.MethodGet, "/v4/elb/show-listener")
	if err != nil {
		return nil, err
	} else {
		if req.ClientToken != nil {
			httpreq.SetQueryParam("clientToken", *req.ClientToken)
		}
		if req.RegionID != nil {
			httpreq.SetQueryParam("regionID", *req.RegionID)
		}
		if req.ListenerID != nil {
			httpreq.SetQueryParam("listenerID", *req.ListenerID)
		}

		httpreq.SetContext(ctx)
	}

	result := &ShowListenerResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
