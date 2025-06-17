package cms

import (
	"bytes"
	"encoding/json"
	"strconv"
)

type apiResponse interface {
	GetStatusCode() string
	GetMessage() string
	GetError() string
	GetErrorMessage() string
}

type apiResponseBase struct {
	StatusCode   json.RawMessage `json:"statusCode,omitempty"`
	Message      *string         `json:"message,omitempty"`
	Error        *string         `json:"error,omitempty"`
	ErrorMessage *string         `json:"errorMessage,omitempty"`
	RequestId    *string         `json:"requestId,omitempty"`
}

func (r *apiResponseBase) GetStatusCode() string {
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

func (r *apiResponseBase) GetMessage() string {
	if r.Message == nil {
		return ""
	}

	return *r.Message
}

func (r *apiResponseBase) GetError() string {
	if r.Error == nil {
		return ""
	}

	return *r.Error
}

func (r *apiResponseBase) GetErrorMessage() string {
	if r.ErrorMessage == nil {
		return ""
	}

	return *r.ErrorMessage
}

var _ apiResponse = (*apiResponseBase)(nil)

type CertificateRecord struct {
	Id                  string `json:"id"`
	Origin              string `json:"origin"`
	Type                string `json:"type"`
	ResourceId          string `json:"resourceId"`
	ResourceType        string `json:"resourceType"`
	CertificateId       string `json:"certificateId"`
	CertificateMode     string `json:"certificateMode"`
	Name                string `json:"name"`
	Status              string `json:"status"`
	DetailStatus        string `json:"detailStatus"`
	ManagedStatus       string `json:"managedStatus"`
	Fingerprint         string `json:"fingerprint"`
	IssueTime           string `json:"issueTime"`
	ExpireTime          string `json:"expireTime"`
	DomainType          string `json:"domainType"`
	DomainName          string `json:"domainName"`
	EncryptionStandard  string `json:"encryptionStandard"`
	EncryptionAlgorithm string `json:"encryptionAlgorithm"`
	CreateTime          string `json:"createTime"`
	UpdateTime          string `json:"updateTime"`
}
