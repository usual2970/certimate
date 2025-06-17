package onepanel

import (
	"context"
	"net/http"
)

type UploadWebsiteSSLRequest struct {
	SSLID           int64  `json:"sslID"`
	Type            string `json:"type"`
	Certificate     string `json:"certificate"`
	CertificatePath string `json:"certificatePath"`
	PrivateKey      string `json:"privateKey"`
	PrivateKeyPath  string `json:"privateKeyPath"`
	Description     string `json:"description"`
}

type UploadWebsiteSSLResponse struct {
	apiResponseBase
}

func (c *Client) UploadWebsiteSSL(req *UploadWebsiteSSLRequest) (*UploadWebsiteSSLResponse, error) {
	return c.UploadWebsiteSSLWithContext(context.Background(), req)
}

func (c *Client) UploadWebsiteSSLWithContext(ctx context.Context, req *UploadWebsiteSSLRequest) (*UploadWebsiteSSLResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/websites/ssl/upload")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &UploadWebsiteSSLResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
