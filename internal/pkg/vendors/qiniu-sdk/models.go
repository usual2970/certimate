package qiniusdk

type UploadSslCertRequest struct {
	Name       string `json:"name"`
	CommonName string `json:"common_name"`
	Pri        string `json:"pri"`
	Ca         string `json:"ca"`
}

type UploadSslCertResponse struct {
	Code   *int    `json:"code,omitempty"`
	Error  *string `json:"error,omitempty"`
	CertID string  `json:"certID"`
}

type DomainInfoHttpsData struct {
	CertID      string `json:"certId"`
	ForceHttps  bool   `json:"forceHttps"`
	Http2Enable bool   `json:"http2Enable"`
}

type GetDomainInfoResponse struct {
	Code               *int                 `json:"code,omitempty"`
	Error              *string              `json:"error,omitempty"`
	Name               string               `json:"name"`
	Type               string               `json:"type"`
	CName              string               `json:"cname"`
	Https              *DomainInfoHttpsData `json:"https"`
	PareDomain         string               `json:"pareDomain"`
	OperationType      string               `json:"operationType"`
	OperatingState     string               `json:"operatingState"`
	OperatingStateDesc string               `json:"operatingStateDesc"`
	CreateAt           string               `json:"createAt"`
	ModifyAt           string               `json:"modifyAt"`
}

type ModifyDomainHttpsConfRequest struct {
	DomainInfoHttpsData
}

type ModifyDomainHttpsConfResponse struct {
	Code  *int    `json:"code,omitempty"`
	Error *string `json:"error,omitempty"`
}

type EnableDomainHttpsRequest struct {
	DomainInfoHttpsData
}

type EnableDomainHttpsResponse struct {
	Code  *int    `json:"code,omitempty"`
	Error *string `json:"error,omitempty"`
}
