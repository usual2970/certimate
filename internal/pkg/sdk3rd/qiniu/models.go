package qiniusdk

type BaseResponse struct {
	Code  *int    `json:"code,omitempty"`
	Error *string `json:"error,omitempty"`
}

type UploadSslCertRequest struct {
	Name        string `json:"name"`
	CommonName  string `json:"common_name"`
	Certificate string `json:"ca"`
	PrivateKey  string `json:"pri"`
}

type UploadSslCertResponse struct {
	BaseResponse
	CertID string `json:"certID"`
}

type DomainInfoHttpsData struct {
	CertID      string `json:"certId"`
	ForceHttps  bool   `json:"forceHttps"`
	Http2Enable bool   `json:"http2Enable"`
}

type GetDomainInfoResponse struct {
	BaseResponse
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
	BaseResponse
}

type EnableDomainHttpsRequest struct {
	DomainInfoHttpsData
}

type EnableDomainHttpsResponse struct {
	BaseResponse
}
