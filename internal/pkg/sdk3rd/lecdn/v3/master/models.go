package master

type BaseResponse interface {
	GetCode() int32
	GetMessage() string
}

type baseResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func (r *baseResponse) GetCode() int32 {
	return r.Code
}

func (r *baseResponse) GetMessage() string {
	return r.Message
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	baseResponse
	Data *struct {
		UserId   int64  `json:"user_id"`
		Username string `json:"username"`
		Token    string `json:"token"`
	} `json:"data,omitempty"`
}

type UpdateCertificateRequest struct {
	ClientId    int64  `json:"client_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	SSLPEM      string `json:"ssl_pem"`
	SSLKey      string `json:"ssl_key"`
	AutoRenewal bool   `json:"auto_renewal"`
}

type UpdateCertificateResponse struct {
	baseResponse
}
