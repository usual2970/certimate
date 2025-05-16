package bunny

import (
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) AddCustomCertificate(req *AddCustomCertificateRequest) ([]byte, error) {
	if req.PullZoneId == "" {
		return nil, fmt.Errorf("bunny api error: invalid parameter: PullZoneId")
	}

	resp, err := c.sendRequest(http.MethodPost, fmt.Sprintf("/pullzone/%s/addCertificate", url.PathEscape(req.PullZoneId)), req)
	return resp.Body(), err
}
