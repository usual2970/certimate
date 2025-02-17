package safelinesdk

import (
	"encoding/json"
)

func (c *Client) UpdateCertificate(req *UpdateCertificateRequest) (*UpdateCertificateResponse, error) {
	params := make(map[string]any)
	jsonData, _ := json.Marshal(req)
	json.Unmarshal(jsonData, &params)

	result := UpdateCertificateResponse{}
	err := c.sendRequestWithResult("/api/open/cert", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
