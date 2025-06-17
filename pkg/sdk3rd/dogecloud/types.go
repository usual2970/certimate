package dogecloud

type apiResponse interface {
	GetCode() int
	GetMessage() string
}

type apiResponseBase struct {
	Code    *int    `json:"code,omitempty"`
	Message *string `json:"msg,omitempty"`
}

func (r *apiResponseBase) GetCode() int {
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
