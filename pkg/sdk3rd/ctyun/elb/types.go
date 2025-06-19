package elb

import (
	"bytes"
	"encoding/json"
	"strconv"
)

type apiResponse interface {
	GetStatusCode() string
	GetMessage() string
	GetError() string
	GetDescription() string
}

type apiResponseBase struct {
	StatusCode  json.RawMessage `json:"statusCode,omitempty"`
	Message     *string         `json:"message,omitempty"`
	Error       *string         `json:"error,omitempty"`
	Description *string         `json:"description,omitempty"`
	RequestId   *string         `json:"requestId,omitempty"`
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

func (r *apiResponseBase) GetDescription() string {
	if r.Description == nil {
		return ""
	}

	return *r.Description
}

var _ apiResponse = (*apiResponseBase)(nil)

type CertificateRecord struct {
	ID          string `json:"ID"`
	RegionID    string `json:"regionID"`
	AzName      string `json:"azName"`
	ProjectID   string `json:"projectID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"privateKey"`
	Status      string `json:"status"`
	CreatedTime string `json:"createdTime"`
	UpdatedTime string `json:"updatedTime"`
}

type ListenerRecord struct {
	ID                  string `json:"ID"`
	RegionID            string `json:"regionID"`
	AzName              string `json:"azName"`
	ProjectID           string `json:"projectID"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	LoadBalancerID      string `json:"loadBalancerID"`
	Protocol            string `json:"protocol"`
	ProtocolPort        int32  `json:"protocolPort"`
	CertificateID       string `json:"certificateID,omitempty"`
	CaEnabled           bool   `json:"caEnabled"`
	ClientCertificateID string `json:"clientCertificateID,omitempty"`
	Status              string `json:"status"`
	CreatedTime         string `json:"createdTime"`
	UpdatedTime         string `json:"updatedTime"`
}
