package btpanelsdk

type BaseResponse interface {
	GetStatus() *bool
	GetMsg() *string
}

type SetSiteSSLRequest struct {
	Type     string `json:"type"`
	SiteName string `json:"siteName"`
	Key      string `json:"key"`
	Csr      string `json:"csr"`
}

type SetSiteSSLResponse struct {
	Status *bool   `json:"status,omitempty"`
	Msg    *string `json:"msg,omitempty"`
}

func (r *SetSiteSSLResponse) GetStatus() *bool {
	return r.Status
}

func (r *SetSiteSSLResponse) GetMsg() *string {
	return r.Msg
}
