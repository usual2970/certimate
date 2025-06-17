package ratpanel

type apiResponse interface {
	GetMessage() string
}

type apiResponseBase struct {
	Message *string `json:"msg,omitempty"`
}

func (r *apiResponseBase) GetMessage() string {
	if r.Message == nil {
		return ""
	}

	return *r.Message
}

var _ apiResponse = (*apiResponseBase)(nil)
