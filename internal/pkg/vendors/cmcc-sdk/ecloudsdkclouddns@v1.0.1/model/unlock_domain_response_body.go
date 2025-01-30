// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type UnlockDomainResponseBodyCodeEnum string

// List of Code
const (
    UnlockDomainResponseBodyCodeEnumError UnlockDomainResponseBodyCodeEnum = "ERROR"
    UnlockDomainResponseBodyCodeEnumSuccess UnlockDomainResponseBodyCodeEnum = "SUCCESS"
)

type UnlockDomainResponseBody struct {

	// 结果说明
	Msg string `json:"msg,omitempty"`

	// 结果码
	Code UnlockDomainResponseBodyCodeEnum `json:"code,omitempty"`

	// 域名
	DomainName string `json:"domainName,omitempty"`

	// 域名ID
	DomainId string `json:"domainId,omitempty"`
}
