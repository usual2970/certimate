// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type DeleteDomainResponseBodyCodeEnum string

// List of Code
const (
    DeleteDomainResponseBodyCodeEnumError DeleteDomainResponseBodyCodeEnum = "ERROR"
    DeleteDomainResponseBodyCodeEnumSuccess DeleteDomainResponseBodyCodeEnum = "SUCCESS"
)

type DeleteDomainResponseBody struct {

	// 结果说明
	Msg string `json:"msg,omitempty"`

	// 结果码
	Code DeleteDomainResponseBodyCodeEnum `json:"code,omitempty"`

	// 域名
	DomainName string `json:"domainName,omitempty"`

	// 域名ID
	DomainId string `json:"domainId,omitempty"`
}
