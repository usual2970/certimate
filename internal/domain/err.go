package domain

var (
	ErrInvalidParams  = NewXError(400, "invalid params")
	ErrRecordNotFound = NewXError(404, "record not found")
)

func IsRecordNotFound(err error) bool {
	if e, ok := err.(*XError); ok {
		return e.GetCode() == ErrRecordNotFound.GetCode()
	}
	return false
}

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
