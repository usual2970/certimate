package bunnysdk

import (
	"fmt"
	"net/http"
)

func (c *Client) AddCustomCertificate(req *AddCustomCertificateRequest) ([]byte, error) {
	resp, err := c.sendRequest(http.MethodPost, fmt.Sprintf("/pullzone/%s/addCertificate", req.PullZoneId), req)
	return resp.Body(), err
}
