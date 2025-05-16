package ratpanelsdk

import "net/http"

func (c *Client) SettingCert(req *SettingCertRequest) (*SettingCertResponse, error) {
	resp := &SettingCertResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/setting/cert", req, resp)
	return resp, err
}

func (c *Client) WebsiteCert(req *WebsiteCertRequest) (*WebsiteCertResponse, error) {
	resp := &WebsiteCertResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/website/cert", req, resp)
	return resp, err
}
