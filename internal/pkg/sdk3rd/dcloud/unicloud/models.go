package unicloud

type BaseResponse interface {
	GetSuccess() bool
	GetErrorCode() string
	GetErrorMessage() string

	GetReturnCode() int32
	GetReturnDesc() string
}

type baseResponse struct {
	Success *bool              `json:"success,omitempty"`
	Header  *map[string]string `json:"header,omitempty"`
	Error   *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`

	ReturnCode *int32  `json:"ret,omitempty"`
	ReturnDesc *string `json:"desc,omitempty"`
}

func (r *baseResponse) GetReturnCode() int32 {
	if r.ReturnCode != nil {
		return *r.ReturnCode
	}
	return 0
}

func (r *baseResponse) GetReturnDesc() string {
	if r.ReturnDesc != nil {
		return *r.ReturnDesc
	}
	return ""
}

func (r *baseResponse) GetSuccess() bool {
	if r.Success != nil {
		return *r.Success
	}
	return false
}

func (r *baseResponse) GetErrorCode() string {
	if r.Error != nil {
		return r.Error.Code
	}
	return ""
}

func (r *baseResponse) GetErrorMessage() string {
	if r.Error != nil {
		return r.Error.Message
	}
	return ""
}

type loginParams struct {
	Email    string `json:"email,omitempty"`
	Mobile   string `json:"mobile,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
}

type loginResponse struct {
	baseResponse
	Data *struct {
		Code     int32  `json:"errCode"`
		UID      string `json:"uid"`
		NewToken *struct {
			Token        string `json:"token"`
			TokenExpired int64  `json:"tokenExpired"`
		} `json:"newToken,omitempty"`
	} `json:"data,omitempty"`
}

type getUserTokenResponse struct {
	baseResponse
	Data *struct {
		Code int32 `json:"code"`
		Data *struct {
			Result      int32  `json:"ret"`
			Description string `json:"desc"`
			Data        *struct {
				Email string `json:"email"`
				Token string `json:"token"`
			} `json:"data,omitempty"`
		} `json:"data,omitempty"`
	} `json:"data,omitempty"`
}

type CreateDomainWithCertRequest struct {
	Provider string `json:"provider"`
	SpaceId  string `json:"spaceId"`
	Domain   string `json:"domain"`
	Cert     string `json:"cert"`
	Key      string `json:"key"`
}

type CreateDomainWithCertResponse struct {
	baseResponse
}
