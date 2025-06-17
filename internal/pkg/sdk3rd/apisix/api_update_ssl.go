package apisix

import (
	"context"
	"fmt"
	"net/http"
)

type UpdateSSLRequest struct {
	Cert   *string   `json:"cert,omitempty"`
	Key    *string   `json:"key,omitempty"`
	SNIs   *[]string `json:"snis,omitempty"`
	Type   *string   `json:"type,omitempty"`
	Status *int32    `json:"status,omitempty"`
}

type UpdateSSLResponse struct {
	apiResponseBase
}

func (c *Client) UpdateSSL(sslId string, req *UpdateSSLRequest) (*UpdateSSLResponse, error) {
	return c.UpdateSSLWithContext(context.Background(), sslId, req)
}

func (c *Client) UpdateSSLWithContext(ctx context.Context, sslId string, req *UpdateSSLRequest) (*UpdateSSLResponse, error) {
	if sslId == "" {
		return nil, fmt.Errorf("sdkerr: unset sslId")
	}

	httpreq, err := c.newRequest(http.MethodPut, fmt.Sprintf("/ssls/%s", sslId))
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &UpdateSSLResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
