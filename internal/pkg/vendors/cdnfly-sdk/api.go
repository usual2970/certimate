package cdnflysdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetSite(req *GetSiteRequest) (*GetSiteResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := GetSiteResponse{}
	err := c.sendRequestWithResult(http.MethodGet, fmt.Sprintf("/v1/sites/%s", req.Id), params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateSite(req *UpdateSiteRequest) (*UpdateSiteResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := UpdateSiteResponse{}
	err := c.sendRequestWithResult(http.MethodPut, fmt.Sprintf("/v1/sites/%s", req.Id), params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateCertificate(req *CreateCertificateRequest) (*CreateCertificateResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := CreateCertificateResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/v1/certs", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateCertificate(req *UpdateCertificateRequest) (*UpdateCertificateResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := UpdateCertificateResponse{}
	err := c.sendRequestWithResult(http.MethodPut, fmt.Sprintf("/v1/certs/%s", req.Id), params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
