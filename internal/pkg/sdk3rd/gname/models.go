package gnamesdk

import "encoding/json"

type BaseResponse interface {
	GetCode() int32
	GetMessage() string
}

type baseResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"msg"`
}

func (r *baseResponse) GetCode() int32 {
	return r.Code
}

func (r *baseResponse) GetMessage() string {
	return r.Message
}

type AddDomainResolutionRequest struct {
	ZoneName    string `json:"ym"`
	RecordType  string `json:"lx"`
	RecordName  string `json:"zj"`
	RecordValue string `json:"jlz"`
	MX          int32  `json:"mx"`
	TTL         int32  `json:"ttl"`
}

type AddDomainResolutionResponse struct {
	baseResponse
	Data json.Number `json:"data"`
}

type ModifyDomainResolutionRequest struct {
	ID          int64  `json:"jxid"`
	ZoneName    string `json:"ym"`
	RecordType  string `json:"lx"`
	RecordName  string `json:"zj"`
	RecordValue string `json:"jlz"`
	MX          int32  `json:"mx"`
	TTL         int32  `json:"ttl"`
}

type ModifyDomainResolutionResponse struct {
	baseResponse
}

type DeleteDomainResolutionRequest struct {
	ZoneName string `json:"ym"`
	RecordID int64  `json:"jxid"`
}

type DeleteDomainResolutionResponse struct {
	baseResponse
}

type ListDomainResolutionRequest struct {
	ZoneName string `json:"ym"`
	Page     *int32 `json:"page,omitempty"`
	PageSize *int32 `json:"limit,omitempty"`
}

type ListDomainResolutionResponse struct {
	baseResponse
	Count    int32               `json:"count"`
	Data     []*ResolutionRecord `json:"data"`
	Page     int32               `json:"page"`
	PageSize int32               `json:"pagesize"`
}

type ResolutionRecord struct {
	ID          json.Number `json:"id"`
	ZoneName    string      `json:"ym"`
	RecordType  string      `json:"lx"`
	RecordName  string      `json:"zjt"`
	RecordValue string      `json:"jxz"`
	MX          int32       `json:"mx"`
}
