package rainyun

type apiResponse interface {
	GetCode() int32
	GetMessage() string
}

type apiResponseBase struct {
	Code    *int32  `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
}

func (r *apiResponseBase) GetCode() int32 {
	if r.Code == nil {
		return 0
	}

	return *r.Code
}

func (r *apiResponseBase) GetMessage() string {
	if r.Message == nil {
		return ""
	}

	return *r.Message
}

var _ apiResponse = (*apiResponseBase)(nil)

type SslRecord struct {
	ID         int32  `json:"ID"`
	UID        int32  `json:"UID"`
	Domain     string `json:"Domain"`
	Issuer     string `json:"Issuer"`
	StartDate  int64  `json:"StartDate"`
	ExpireDate int64  `json:"ExpDate"`
	UploadTime int64  `json:"UploadTime"`
}

type SslDetail struct {
	Cert       string `json:"Cert"`
	Key        string `json:"Key"`
	Domain     string `json:"DomainName"`
	Issuer     string `json:"Issuer"`
	StartDate  int64  `json:"StartDate"`
	ExpireDate int64  `json:"ExpDate"`
	RemainDays int32  `json:"RemainDays"`
}
