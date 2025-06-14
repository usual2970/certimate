package icdn

import (
	"bytes"
	"encoding/json"
	"strconv"
)

type baseResultInterface interface {
	GetStatusCode() string
	GetMessage() string
	GetError() string
	GetErrorMessage() string
}

type baseResult struct {
	StatusCode   json.RawMessage `json:"statusCode,omitempty"`
	Message      *string         `json:"message,omitempty"`
	Error        *string         `json:"error,omitempty"`
	ErrorMessage *string         `json:"errorMessage,omitempty"`
	RequestId    *string         `json:"requestId,omitempty"`
}

func (r *baseResult) GetStatusCode() string {
	if r.StatusCode == nil {
		return ""
	}

	decoder := json.NewDecoder(bytes.NewReader(r.StatusCode))
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

func (r *baseResult) GetMessage() string {
	if r.Message == nil {
		return ""
	}

	return *r.Message
}

func (r *baseResult) GetError() string {
	if r.Error == nil {
		return ""
	}

	return *r.Error
}

func (r *baseResult) GetErrorMessage() string {
	if r.ErrorMessage == nil {
		return ""
	}

	return *r.ErrorMessage
}

var _ baseResultInterface = (*baseResult)(nil)

type CertRecord struct {
	Id          int64    `json:"id"`
	Name        string   `json:"name"`
	CN          string   `json:"cn"`
	SANs        []string `json:"sans"`
	UsageMode   int32    `json:"usage_mode"`
	State       int32    `json:"state"`
	ExpiresTime int64    `json:"expires"`
	IssueTime   int64    `json:"issue"`
	Issuer      string   `json:"issuer"`
	CreatedTime int64    `json:"created"`
}

type CertDetail struct {
	CertRecord
	Certs string `json:"certs"`
	Key   string `json:"key"`
}
