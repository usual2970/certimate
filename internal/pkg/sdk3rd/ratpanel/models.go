package ratpanelsdk

type BaseResponse interface {
	GetMessage() string
}

type baseResponse struct {
	Message *string `json:"msg,omitempty"`
}

func (r *baseResponse) GetMessage() string {
	if r.Message != nil {
		return *r.Message
	}
	return ""
}

type SettingCertRequest struct {
	Certificate string `json:"cert"`
	PrivateKey  string `json:"key"`
}

type SettingCertResponse struct {
	baseResponse
}

type WebsiteCertRequest struct {
	SiteName    string `json:"name"`
	Certificate string `json:"cert"`
	PrivateKey  string `json:"key"`
}

type WebsiteCertResponse struct {
	baseResponse
}
