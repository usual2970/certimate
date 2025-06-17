package cdn

type apiResponse interface {
	GetCode() string
	GetMessage() string
}

type apiResponseBase struct {
	Code    *string `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
}

var _ apiResponse = (*apiResponseBase)(nil)

func (r *apiResponseBase) GetCode() string {
	if r.Code == nil {
		return ""
	}

	return *r.Code
}

func (r *apiResponseBase) GetMessage() string {
	if r.Message == nil {
		return ""
	}

	return *r.Message
}
