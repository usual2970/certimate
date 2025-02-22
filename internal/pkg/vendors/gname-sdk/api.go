package gnamesdk

func (c *Client) AddDomainResolution(req *AddDomainResolutionRequest) (*AddDomainResolutionResponse, error) {
	result := AddDomainResolutionResponse{}
	err := c.sendRequestWithResult("/api/resolution/add", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) ModifyDomainResolution(req *ModifyDomainResolutionRequest) (*ModifyDomainResolutionResponse, error) {
	resp := ModifyDomainResolutionResponse{}
	err := c.sendRequestWithResult("/api/resolution/edit", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) DeleteDomainResolution(req *DeleteDomainResolutionRequest) (*DeleteDomainResolutionResponse, error) {
	resp := DeleteDomainResolutionResponse{}
	err := c.sendRequestWithResult("/api/resolution/delete", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) ListDomainResolution(req *ListDomainResolutionRequest) (*ListDomainResolutionResponse, error) {
	resp := ListDomainResolutionResponse{}
	err := c.sendRequestWithResult("/api/resolution/list", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
