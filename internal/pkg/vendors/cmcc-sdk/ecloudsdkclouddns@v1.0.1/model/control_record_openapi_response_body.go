// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ControlRecordOpenapiResponseBodyCodeEnum string

// List of Code
const (
    ControlRecordOpenapiResponseBodyCodeEnumError ControlRecordOpenapiResponseBodyCodeEnum = "ERROR"
    ControlRecordOpenapiResponseBodyCodeEnumSuccess ControlRecordOpenapiResponseBodyCodeEnum = "SUCCESS"
)

type ControlRecordOpenapiResponseBody struct {

	// 结果说明
	Msg string `json:"msg,omitempty"`

	// 解析记录ID
	RecordId string `json:"recordId,omitempty"`

	// 结果码
	Code ControlRecordOpenapiResponseBodyCodeEnum `json:"code,omitempty"`

	// 域名
	DomainName string `json:"domainName,omitempty"`
}
