package btpanel

type BaseResponse interface {
	GetStatus() *bool
	GetMessage() *string
}

type baseResponse struct {
	Status  *bool   `json:"status,omitempty"`
	Message *string `json:"msg,omitempty"`
}

func (r *baseResponse) GetStatus() *bool {
	return r.Status
}

func (r *baseResponse) GetMessage() *string {
	return r.Message
}

type ConfigSavePanelSSLRequest struct {
	PrivateKey  string `json:"privateKey"`
	Certificate string `json:"certPem"`
}

type ConfigSavePanelSSLResponse struct {
	baseResponse
}

type SiteSetSSLRequest struct {
	Type        string `json:"type"`
	SiteName    string `json:"siteName"`
	PrivateKey  string `json:"key"`
	Certificate string `json:"csr"`
}

type SiteSetSSLResponse struct {
	baseResponse
}

type SystemServiceAdminRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type SystemServiceAdminResponse struct {
	baseResponse
}

type SSLCertSaveCertRequest struct {
	PrivateKey  string `json:"key"`
	Certificate string `json:"csr"`
}

type SSLCertSaveCertResponse struct {
	baseResponse
	SSLHash string `json:"ssl_hash"`
}

type SSLSetBatchCertToSiteRequest struct {
	BatchInfo []*SSLSetBatchCertToSiteRequestBatchInfo `json:"BatchInfo"`
}

type SSLSetBatchCertToSiteRequestBatchInfo struct {
	SSLHash  string `json:"ssl_hash"`
	SiteName string `json:"siteName"`
	CertName string `json:"certName"`
}

type SSLSetBatchCertToSiteResponse struct {
	baseResponse
	TotalCount   int32 `json:"total"`
	SuccessCount int32 `json:"success"`
	FailedCount  int32 `json:"faild"`
}
