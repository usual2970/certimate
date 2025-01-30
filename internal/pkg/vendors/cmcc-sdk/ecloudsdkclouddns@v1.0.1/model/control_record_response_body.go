// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ControlRecordResponseBodyCodeEnum string

// List of Code
const (
    ControlRecordResponseBodyCodeEnumError ControlRecordResponseBodyCodeEnum = "ERROR"
    ControlRecordResponseBodyCodeEnumSuccess ControlRecordResponseBodyCodeEnum = "SUCCESS"
)

type ControlRecordResponseBody struct {

	// 结果说明
	Msg string `json:"msg,omitempty"`

	// 解析记录ID
	RecordId string `json:"recordId,omitempty"`

	// 结果码
	Code ControlRecordResponseBodyCodeEnum `json:"code,omitempty"`

	// 域名
	DomainName string `json:"domainName,omitempty"`
}
