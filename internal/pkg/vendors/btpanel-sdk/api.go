package btpanelsdk

import (
	"encoding/json"
)

func (c *Client) ConfigSavePanelSSL(req *ConfigSavePanelSSLRequest) (*ConfigSavePanelSSLResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := ConfigSavePanelSSLResponse{}
	err := c.sendRequestWithResult("/config?action=SavePanelSSL", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SiteSetSSL(req *SiteSetSSLRequest) (*SiteSetSSLResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := SiteSetSSLResponse{}
	err := c.sendRequestWithResult("/site?action=SetSSL", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SystemServiceAdmin(req *SystemServiceAdminRequest) (*SystemServiceAdminResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := SystemServiceAdminResponse{}
	err := c.sendRequestWithResult("/system?action=ServiceAdmin", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SSLCertSaveCert(req *SSLCertSaveCertRequest) (*SSLCertSaveCertResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := SSLCertSaveCertResponse{}
	err := c.sendRequestWithResult("/ssl/cert/save_cert", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SSLSetBatchCertToSite(req *SSLSetBatchCertToSiteRequest) (*SSLSetBatchCertToSiteResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := SSLSetBatchCertToSiteResponse{}
	err := c.sendRequestWithResult("/ssl?action=SetBatchCertToSite", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
