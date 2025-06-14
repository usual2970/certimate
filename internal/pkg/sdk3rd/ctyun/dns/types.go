package dns

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

type DnsRecord struct {
	RecordId int32  `json:"recordId"`
	Host     string `json:"host"`
	Type     string `json:"type"`
	LineCode string `json:"lineCode"`
	Value    string `json:"value"`
	TTL      int32  `json:"ttl"`
	State    int32  `json:"state"`
	Remark   string `json:"remark"`
}
