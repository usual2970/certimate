package domain

var ErrAuthFailed = NewXError(4999, "auth failed")

type XError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewXError(code int, msg string) *XError {
	return &XError{code, msg}
}

func (e *XError) Error() string {
	return e.Msg
}

func (e *XError) GetCode() int {
	if e.Code == 0 {
		return 100
	}
	return e.Code
}
