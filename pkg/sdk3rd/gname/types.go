package gname

import "encoding/json"

type apiResponse interface {
	GetCode() int32
	GetMessage() string
}

type apiResponseBase struct {
	Code    int32  `json:"code"`
	Message string `json:"msg"`
}

func (r *apiResponseBase) GetCode() int32 {
	return r.Code
}

func (r *apiResponseBase) GetMessage() string {
	return r.Message
}

var _ apiResponse = (*apiResponseBase)(nil)

type DomainResolutionRecordord struct {
	ID          json.Number `json:"id"`
	ZoneName    string      `json:"ym"`
	RecordType  string      `json:"lx"`
	RecordName  string      `json:"zjt"`
	RecordValue string      `json:"jxz"`
	MX          int32       `json:"mx"`
}
