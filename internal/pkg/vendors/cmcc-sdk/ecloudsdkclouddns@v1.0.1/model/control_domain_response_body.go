// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ControlDomainResponseBodyCodeEnum string

// List of Code
const (
    ControlDomainResponseBodyCodeEnumError ControlDomainResponseBodyCodeEnum = "ERROR"
    ControlDomainResponseBodyCodeEnumSuccess ControlDomainResponseBodyCodeEnum = "SUCCESS"
)

type ControlDomainResponseBody struct {

	// 结果说明
	Msg string `json:"msg,omitempty"`

	// 结果码
	Code ControlDomainResponseBodyCodeEnum `json:"code,omitempty"`

	// 域名
	DomainName string `json:"domainName,omitempty"`

	// 域名ID
	DomainId string `json:"domainId,omitempty"`
}
