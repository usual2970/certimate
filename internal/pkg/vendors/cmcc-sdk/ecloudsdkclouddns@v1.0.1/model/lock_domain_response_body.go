// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type LockDomainResponseBodyCodeEnum string

// List of Code
const (
    LockDomainResponseBodyCodeEnumError LockDomainResponseBodyCodeEnum = "ERROR"
    LockDomainResponseBodyCodeEnumSuccess LockDomainResponseBodyCodeEnum = "SUCCESS"
)

type LockDomainResponseBody struct {

	// 结果说明
	Msg string `json:"msg,omitempty"`

	// 结果码
	Code LockDomainResponseBodyCodeEnum `json:"code,omitempty"`

	// 域名
	DomainName string `json:"domainName,omitempty"`

	// 域名ID
	DomainId string `json:"domainId,omitempty"`
}
