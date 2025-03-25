package gnamesdk

func (c *Client) AddDomainResolution(req *AddDomainResolutionRequest) (*AddDomainResolutionResponse, error) {
	resp := &AddDomainResolutionResponse{}
	err := c.sendRequestWithResult("/api/resolution/add", req, resp)
	return resp, err
}

func (c *Client) ModifyDomainResolution(req *ModifyDomainResolutionRequest) (*ModifyDomainResolutionResponse, error) {
	resp := &ModifyDomainResolutionResponse{}
	err := c.sendRequestWithResult("/api/resolution/edit", req, resp)
	return resp, err
}

func (c *Client) DeleteDomainResolution(req *DeleteDomainResolutionRequest) (*DeleteDomainResolutionResponse, error) {
	resp := &DeleteDomainResolutionResponse{}
	err := c.sendRequestWithResult("/api/resolution/delete", req, resp)
	return resp, err
}

func (c *Client) ListDomainResolution(req *ListDomainResolutionRequest) (*ListDomainResolutionResponse, error) {
	resp := &ListDomainResolutionResponse{}
	err := c.sendRequestWithResult("/api/resolution/list", req, resp)
	return resp, err
}
