package cms

import (
	"context"
	"net/http"
)

type UploadCertificateRequest struct {
	Name               *string `json:"name,omitempty"`
	Certificate        *string `json:"certificate,omitempty"`
	CertificateChain   *string `json:"certificateChain,omitempty"`
	PrivateKey         *string `json:"privateKey,omitempty"`
	EncryptionStandard *string `json:"encryptionStandard,omitempty"`
	EncCertificate     *string `json:"encCertificate,omitempty"`
	EncPrivateKey      *string `json:"encPrivateKey,omitempty"`
}

type UploadCertificateResponse struct {
	apiResponseBase
}

func (c *Client) UploadCertificate(req *UploadCertificateRequest) (*UploadCertificateResponse, error) {
	return c.UploadCertificateWithContext(context.Background(), req)
}

func (c *Client) UploadCertificateWithContext(ctx context.Context, req *UploadCertificateRequest) (*UploadCertificateResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/v1/certificate/upload")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &UploadCertificateResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
