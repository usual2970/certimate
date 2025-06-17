package console

import (
	"encoding/json"
)

type apiResponse interface {
	GetData() *apiResponseBaseData
}

type apiResponseBase struct {
	Data *apiResponseBaseData `json:"data,omitempty"`
}

func (r *apiResponseBase) GetData() *apiResponseBaseData {
	return r.Data
}

var _ apiResponse = (*apiResponseBase)(nil)

type apiResponseBaseData struct {
	ErrorCode json.Number `json:"error_code,omitempty"`
	Message   string      `json:"message,omitempty"`
}

func (r *apiResponseBaseData) GetErrorCode() int32 {
	if r.ErrorCode.String() == "" {
		return 0
	}

	errcode, err := r.ErrorCode.Int64()
	if err != nil {
		return -1
	}

	return int32(errcode)
}

func (r *apiResponseBaseData) GetMessage() string {
	return r.Message
}
