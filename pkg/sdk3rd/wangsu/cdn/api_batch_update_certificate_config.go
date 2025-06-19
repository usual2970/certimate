package cdn

import (
	"context"
	"net/http"
)

type BatchUpdateCertificateConfigRequest struct {
	CertificateId int64    `json:"certificateId"`
	DomainNames   []string `json:"domainNames"`
}

type BatchUpdateCertificateConfigResponse struct {
	apiResponseBase
}

func (c *Client) BatchUpdateCertificateConfig(req *BatchUpdateCertificateConfigRequest) (*BatchUpdateCertificateConfigResponse, error) {
	return c.BatchUpdateCertificateConfigWithContext(context.Background(), req)
}

func (c *Client) BatchUpdateCertificateConfigWithContext(ctx context.Context, req *BatchUpdateCertificateConfigRequest) (*BatchUpdateCertificateConfigResponse, error) {
	httpreq, err := c.newRequest(http.MethodPut, "/api/config/certificate/batch")
	if err != nil {
		return nil, err
	} else {
		httpreq.SetBody(req)
		httpreq.SetContext(ctx)
	}

	result := &BatchUpdateCertificateConfigResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
