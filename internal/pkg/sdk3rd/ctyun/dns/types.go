package dns

import "encoding/json"

type baseResult struct {
	StatusCode   json.RawMessage `json:"statusCode,omitempty"`
	Message      *string         `json:"message,omitempty"`
	Error        *string         `json:"error,omitempty"`
	ErrorMessage *string         `json:"errorMessage,omitempty"`
	RequestId    *string         `json:"requestId,omitempty"`
}

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
