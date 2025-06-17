package btpanel

import (
	"context"
	"net/http"
)

type SSLSetBatchCertToSiteRequest struct {
	BatchInfo []*SSLSetBatchCertToSiteRequestBatchInfo `json:"BatchInfo"`
}

type SSLSetBatchCertToSiteRequestBatchInfo struct {
	SSLHash  string `json:"ssl_hash"`
	SiteName string `json:"siteName"`
	CertName string `json:"certName"`
}

type SSLSetBatchCertToSiteResponse struct {
	apiResponseBase

	TotalCount   int32 `json:"total"`
	SuccessCount int32 `json:"success"`
	FailedCount  int32 `json:"faild"`
}

func (c *Client) SSLSetBatchCertToSite(req *SSLSetBatchCertToSiteRequest) (*SSLSetBatchCertToSiteResponse, error) {
	return c.SSLSetBatchCertToSiteWithContext(context.Background(), req)
}

func (c *Client) SSLSetBatchCertToSiteWithContext(ctx context.Context, req *SSLSetBatchCertToSiteRequest) (*SSLSetBatchCertToSiteResponse, error) {
	httpreq, err := c.newRequest(http.MethodPost, "/ssl?action=SetBatchCertToSite", req)
	if err != nil {
		return nil, err
	} else {
		httpreq.SetContext(ctx)
	}

	result := &SSLSetBatchCertToSiteResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
