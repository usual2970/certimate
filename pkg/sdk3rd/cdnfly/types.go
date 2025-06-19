package cdnfly

import (
	"bytes"
	"encoding/json"
	"strconv"
)

type apiResponse interface {
	GetCode() string
	GetMessage() string
}

type apiResponseBase struct {
	Code    json.RawMessage `json:"code"`
	Message string          `json:"msg"`
}

func (r *apiResponseBase) GetCode() string {
	if r.Code == nil {
		return ""
	}

	decoder := json.NewDecoder(bytes.NewReader(r.Code))
	token, err := decoder.Token()
	if err != nil {
		return ""
	}

	switch t := token.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case json.Number:
		return t.String()
	default:
		return ""
	}
}

func (r *apiResponseBase) GetMessage() string {
	return r.Message
}

var _ apiResponse = (*apiResponseBase)(nil)
