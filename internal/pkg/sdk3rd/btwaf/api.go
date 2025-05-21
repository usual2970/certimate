package btwaf

func (c *Client) GetSiteList(req *GetSiteListRequest) (*GetSiteListResponse, error) {
	resp := &GetSiteListResponse{}
	err := c.sendRequestWithResult("/wafmastersite/get_site_list", req, resp)
	return resp, err
}

func (c *Client) ModifySite(req *ModifySiteRequest) (*ModifySiteResponse, error) {
	resp := &ModifySiteResponse{}
	err := c.sendRequestWithResult("/wafmastersite/modify_site", req, resp)
	return resp, err
}

func (c *Client) ConfigSetSSL(req *ConfigSetSSLRequest) (*ConfigSetSSLResponse, error) {
	resp := &ConfigSetSSLResponse{}
	err := c.sendRequestWithResult("/config/set_cert", req, resp)
	return resp, err
}
