package rainyun

type BaseResponse interface {
	GetCode() int32
	GetMessage() string
}

type baseResponse struct {
	Code    *int32  `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
}

func (r *baseResponse) GetCode() int32 {
	if r.Code != nil {
		return *r.Code
	}
	return 0
}

func (r *baseResponse) GetMessage() string {
	if r.Message != nil {
		return *r.Message
	}
	return ""
}

type SslCenterListFilters struct {
	Domain *string `json:"Domain,omitempty"`
}

type SslCenterListRequest struct {
	Filters *SslCenterListFilters `json:"columnFilters,omitempty"`
	Sort    []*string             `json:"sort,omitempty"`
	Page    *int32                `json:"page,omitempty"`
	PerPage *int32                `json:"perPage,omitempty"`
}

type SslCenterListResponse struct {
	baseResponse
	Data *struct {
		TotalRecords int32 `json:"TotalRecords"`
		Records      []*struct {
			ID         int32  `json:"ID"`
			UID        int32  `json:"UID"`
			Domain     string `json:"Domain"`
			Issuer     string `json:"Issuer"`
			StartDate  int64  `json:"StartDate"`
			ExpireDate int64  `json:"ExpDate"`
			UploadTime int64  `json:"UploadTime"`
		} `json:"Records"`
	} `json:"data,omitempty"`
}

type SslCenterGetResponse struct {
	baseResponse
	Data *struct {
		Cert       string `json:"Cert"`
		Key        string `json:"Key"`
		Domain     string `json:"DomainName"`
		Issuer     string `json:"Issuer"`
		StartDate  int64  `json:"StartDate"`
		ExpireDate int64  `json:"ExpDate"`
		RemainDays int32  `json:"RemainDays"`
	} `json:"data,omitempty"`
}

type SslCenterCreateRequest struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

type SslCenterCreateResponse struct {
	baseResponse
}

type RcdnInstanceSslBindRequest struct {
	CertId  int32    `json:"cert_id"`
	Domains []string `json:"domains"`
}

type RcdnInstanceSslBindResponse struct {
	baseResponse
}
