package cachefly

type apiResponse interface {
	GetMessage() string
}

type apiResponseBase struct {
	Message *string `json:"message,omitempty"`
}

func (r *apiResponseBase) GetMessage() string {
	if r.Message == nil {
		return ""
	}

	return *r.Message
}

var _ apiResponse = (*apiResponseBase)(nil)
