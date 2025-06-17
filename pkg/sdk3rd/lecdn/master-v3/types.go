package masterv3

type apiResponse interface {
	GetCode() int32
	GetMessage() string
}

type apiResponseBase struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func (r *apiResponseBase) GetCode() int32 {
	return r.Code
}

func (r *apiResponseBase) GetMessage() string {
	return r.Message
}

var _ apiResponse = (*apiResponseBase)(nil)
