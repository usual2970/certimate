package btpanel

type apiResponse interface {
	GetStatus() *bool
	GetMessage() *string
}

type apiResponseBase struct {
	Status  *bool   `json:"status,omitempty"`
	Message *string `json:"msg,omitempty"`
}

func (r *apiResponseBase) GetStatus() *bool {
	return r.Status
}

func (r *apiResponseBase) GetMessage() *string {
	return r.Message
}
