package btpanel

import (
	"context"
	"net/http"
)

type SSLCertSaveCertRequest struct {
	PrivateKey  string `json:"key"`
	Certificate string `json:"csr"`
}

type SSLCertSaveCertResponse struct {
	apiResponseBase

	SSLHash string `json:"ssl_hash"`
}

func (c *Client) SSLCertSaveCert(req *SSLCertSaveCertRequest) (*SSLCertSaveCertResponse, error) {
	return c.SSLCertSaveCertWithContext(context.Background(), req)
}

func (c *Client) SSLCertSaveCertWithContext(ctx context.Context, req *SSLCertSaveCertRequest) (*SSLCertSaveCertResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/ssl/cert/save_cert", req)
	if err != nil {
		return nil, err
	} else {
		httpreq.SetContext(ctx)
	}

	result := &SSLCertSaveCertResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
