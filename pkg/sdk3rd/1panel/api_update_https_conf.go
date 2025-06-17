package onepanel

import (
	"context"
	"fmt"
	"net/http"
)

type UpdateHttpsConfRequest struct {
	WebsiteID       int64    `json:"websiteId"`
	Enable          bool     `json:"enable"`
	Type            string   `json:"type"`
	WebsiteSSLID    int64    `json:"websiteSSLId"`
	PrivateKey      string   `json:"privateKey"`
	Certificate     string   `json:"certificate"`
	PrivateKeyPath  string   `json:"privateKeyPath"`
	CertificatePath string   `json:"certificatePath"`
	ImportType      string   `json:"importType"`
	HttpConfig      string   `json:"httpConfig"`
	SSLProtocol     []string `json:"SSLProtocol"`
	Algorithm       string   `json:"algorithm"`
	Hsts            bool     `json:"hsts"`
}

type UpdateHttpsConfResponse struct {
	apiResponseBase
}

func (c *Client) UpdateHttpsConf(websiteId int64, req *UpdateHttpsConfRequest) (*UpdateHttpsConfResponse, error) {
	return c.UpdateHttpsConfWithContext(context.Background(), websiteId, req)
}

func (c *Client) UpdateHttpsConfWithContext(ctx context.Context, websiteId int64, req *UpdateHttpsConfRequest) (*UpdateHttpsConfResponse, error) {
	if websiteId == 0 {
		return nil, fmt.Errorf("sdkerr: unset websiteId")
	}

	httpreq, err := c.newRequest(http.MethodPost, fmt.Sprintf("/websites/%d/https", websiteId))
	if err != nil {
		return nil, err
	} else {
		req.WebsiteID = websiteId
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &UpdateHttpsConfResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
