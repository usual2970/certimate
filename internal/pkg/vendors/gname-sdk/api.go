package gnamesdk

import (
	"encoding/json"
)

func (c *Client) AddDomainResolution(req *AddDomainResolutionRequest) (*AddDomainResolutionResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := AddDomainResolutionResponse{}
	err := c.sendRequestWithResult("/api/resolution/add", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) ModifyDomainResolution(req *ModifyDomainResolutionRequest) (*ModifyDomainResolutionResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := ModifyDomainResolutionResponse{}
	err := c.sendRequestWithResult("/api/resolution/edit", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteDomainResolution(req *DeleteDomainResolutionRequest) (*DeleteDomainResolutionResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := DeleteDomainResolutionResponse{}
	err := c.sendRequestWithResult("/api/resolution/delete", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) ListDomainResolution(req *ListDomainResolutionRequest) (*ListDomainResolutionResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := ListDomainResolutionResponse{}
	err := c.sendRequestWithResult("/api/resolution/list", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
