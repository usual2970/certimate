package apisix

import (
	"fmt"
	"net/http"
)

func (c *Client) UpdateSSL(req *UpdateSSLRequest) (*UpdateSSLResponse, error) {
	if req.ID == "" {
		return nil, fmt.Errorf("1panel api error: invalid parameter: ID")
	}

	resp := &UpdateSSLResponse{}
	err := c.sendRequestWithResult(http.MethodGet, fmt.Sprintf("/ssls/%s", req.ID), req, resp)
	return resp, err
}
