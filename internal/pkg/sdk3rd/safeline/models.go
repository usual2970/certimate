package safeline

type BaseResponse interface {
	GetErrCode() *string
	GetErrMsg() *string
}

type baseResponse struct {
	ErrCode *string `json:"err,omitempty"`
	ErrMsg  *string `json:"msg,omitempty"`
}

func (r *baseResponse) GetErrCode() *string {
	return r.ErrCode
}

func (r *baseResponse) GetErrMsg() *string {
	return r.ErrMsg
}

type UpdateCertificateRequest struct {
	Id     int32                              `json:"id"`
	Type   int32                              `json:"type"`
	Manual *UpdateCertificateRequestBodyManul `json:"manual"`
}

type UpdateCertificateRequestBodyManul struct {
	Crt string `json:"crt"`
	Key string `json:"key"`
}

type UpdateCertificateResponse struct {
	baseResponse
}
