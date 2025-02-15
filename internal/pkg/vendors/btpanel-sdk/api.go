package btpanelsdk

type BaseResponse interface {
	GetStatus() *bool
	GetMsg() *string
}

type baseResponse struct {
	Status *bool   `json:"status,omitempty"`
	Msg    *string `json:"msg,omitempty"`
}

func (r *baseResponse) GetStatus() *bool {
	return r.Status
}

func (r *baseResponse) GetMsg() *string {
	return r.Msg
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
