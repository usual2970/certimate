package btwaf

type BaseResponse interface {
	GetCode() int32
}

type baseResponse struct {
	Code *int32 `json:"code,omitempty"`
}

func (r *baseResponse) GetCode() int32 {
	if r.Code != nil {
		return *r.Code
	}
	return 0
}

type GetSiteListRequest struct {
	Page     *int32  `json:"p,omitempty"`
	PageSize *int32  `json:"p_size,omitempty"`
	SiteName *string `json:"site_name,omitempty"`
}

type GetSiteListResponse struct {
	baseResponse
	Result *struct {
		List []*struct {
			SiteId     string `json:"site_id"`
			SiteName   string `json:"site_name"`
			Type       string `json:"types"`
			Status     int32  `json:"status"`
			CreateTime int64  `json:"create_time"`
			UpdateTime int64  `json:"update_time"`
		} `json:"list"`
		Total int32 `json:"total"`
	} `json:"res,omitempty"`
}

type SiteServerInfo struct {
	ListenSSLPorts *[]int32           `json:"listen_ssl_port,omitempty"`
	SSL            *SiteServerSSLInfo `json:"ssl,omitempty"`
}

type SiteServerSSLInfo struct {
	IsSSL      *int32  `json:"is_ssl,omitempty"`
	FullChain  *string `json:"full_chain,omitempty"`
	PrivateKey *string `json:"private_key,omitempty"`
}

type ModifySiteRequest struct {
	SiteId string          `json:"site_id"`
	Type   *string         `json:"types,omitempty"`
	Server *SiteServerInfo `json:"server,omitempty"`
}

type ModifySiteResponse struct {
	baseResponse
}

type ConfigSetSSLRequest struct {
	CertContent string `json:"certContent"`
	KeyContent  string `json:"keyContent"`
}

type ConfigSetSSLResponse struct {
	baseResponse
}
