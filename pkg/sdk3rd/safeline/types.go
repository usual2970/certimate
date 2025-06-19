package safeline

type apiResponse interface {
	GetErrCode() string
	GetErrMsg() string
}

type apiResponseBase struct {
	ErrCode *string `json:"err,omitempty"`
	ErrMsg  *string `json:"msg,omitempty"`
}

func (r *apiResponseBase) GetErrCode() string {
	if r.ErrCode == nil {
		return ""
	}

	return *r.ErrCode
}

func (r *apiResponseBase) GetErrMsg() string {
	if r.ErrMsg == nil {
		return ""
	}

	return *r.ErrMsg
}

var _ apiResponse = (*apiResponseBase)(nil)

type CertificateManul struct {
	Crt string `json:"crt"`
	Key string `json:"key"`
}
