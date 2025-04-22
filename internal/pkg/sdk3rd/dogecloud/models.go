package dogecloudsdk

type BaseResponse struct {
	Code    *int    `json:"code,omitempty"`
	Message *string `json:"msg,omitempty"`
}

type UploadCdnCertRequest struct {
	Note        string `json:"note"`
	Certificate string `json:"cert"`
	PrivateKey  string `json:"private"`
}

type UploadCdnCertResponseData struct {
	Id int64 `json:"id"`
}

type UploadCdnCertResponse struct {
	BaseResponse
	Data *UploadCdnCertResponseData `json:"data,omitempty"`
}

type BindCdnCertRequest struct {
	CertId   int64   `json:"id"`
	DomainId *int64  `json:"did,omitempty"`
	Domain   *string `json:"domain,omitempty"`
}

type BindCdnCertResponse struct {
	BaseResponse
}
