package unicloud

type apiResponse interface {
	GetSuccess() bool
	GetErrorCode() string
	GetErrorMessage() string

	GetReturnCode() int32
	GetReturnDesc() string
}

type apiResponseBase struct {
	Success *bool              `json:"success,omitempty"`
	Header  *map[string]string `json:"header,omitempty"`
	Error   *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`

	ReturnCode *int32  `json:"ret,omitempty"`
	ReturnDesc *string `json:"desc,omitempty"`
}

func (r *apiResponseBase) GetReturnCode() int32 {
	if r.ReturnCode == nil {
		return 0
	}

	return *r.ReturnCode
}

func (r *apiResponseBase) GetReturnDesc() string {
	if r.ReturnDesc == nil {
		return ""
	}

	return *r.ReturnDesc
}

func (r *apiResponseBase) GetSuccess() bool {
	if r.Success == nil {
		return false
	}

	return *r.Success
}

func (r *apiResponseBase) GetErrorCode() string {
	if r.Error == nil {
		return ""
	}

	return r.Error.Code
}

func (r *apiResponseBase) GetErrorMessage() string {
	if r.Error == nil {
		return ""
	}

	return r.Error.Message
}

var _ apiResponse = (*apiResponseBase)(nil)
