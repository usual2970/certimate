package btpanelsdk

import "net/http"

func (c *Client) SettingCert(req *SettingCertRequest) (*SettingCertResponse, error) {
	resp := &SettingCertResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/setting/cert", req, resp)
	return resp, err
}

func (c *Client) WebsiteCert(req *SiteCertRequest) (*SiteCertResponse, error) {
	resp := &SiteCertResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/website/cert", req, resp)
	return resp, err
}
