package cdn

import (
	"net/http"
)

func (c *Client) BatchUpdateCertificateConfig(req *BatchUpdateCertificateConfigRequest) (*BatchUpdateCertificateConfigResponse, error) {
	resp := &BatchUpdateCertificateConfigResponse{}
	_, err := c.client.SendRequestWithResult(http.MethodPut, "/api/config/certificate/batch", req, resp)
	if err != nil {
		return resp, err
	}

	return resp, err
}
