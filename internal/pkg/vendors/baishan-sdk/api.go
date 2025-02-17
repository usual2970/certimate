package baishansdk

import (
	"encoding/json"
	"net/http"
)

func (c *Client) CreateCertificate(req *CreateCertificateRequest) (*CreateCertificateResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := CreateCertificateResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/v2/domain/certificate", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetDomainConfig(req *GetDomainConfigRequest) (*GetDomainConfigResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := GetDomainConfigResponse{}
	err := c.sendRequestWithResult(http.MethodGet, "/v2/domain/config", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SetDomainConfig(req *SetDomainConfigRequest) (*SetDomainConfigResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := SetDomainConfigResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/v2/domain/config", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
