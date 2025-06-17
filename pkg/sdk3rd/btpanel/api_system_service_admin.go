package btpanel

import (
	"context"
	"net/http"
)

type SystemServiceAdminRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type SystemServiceAdminResponse struct {
	apiResponseBase
}

func (c *Client) SystemServiceAdmin(req *SystemServiceAdminRequest) (*SystemServiceAdminResponse, error) {
	return c.SystemServiceAdminWithContext(context.Background(), req)
}

func (c *Client) SystemServiceAdminWithContext(ctx context.Context, req *SystemServiceAdminRequest) (*SystemServiceAdminResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/system?action=ServiceAdmin", req)
	if err != nil {
		return nil, err
	} else {
		httpreq.SetContext(ctx)
	}

	result := &SystemServiceAdminResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
