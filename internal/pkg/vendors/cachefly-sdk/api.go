package cacheflysdk

import (
	"encoding/json"
	"net/http"
)

func (c *Client) CreateCertificate(req *CreateCertificateRequest) (*CreateCertificateResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := CreateCertificateResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/certificates", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
