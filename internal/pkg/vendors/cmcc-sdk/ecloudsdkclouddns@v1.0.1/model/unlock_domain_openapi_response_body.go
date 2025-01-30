// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type UnlockDomainOpenapiResponseBodyCodeEnum string

// List of Code
const (
    UnlockDomainOpenapiResponseBodyCodeEnumError UnlockDomainOpenapiResponseBodyCodeEnum = "ERROR"
    UnlockDomainOpenapiResponseBodyCodeEnumSuccess UnlockDomainOpenapiResponseBodyCodeEnum = "SUCCESS"
)

type UnlockDomainOpenapiResponseBody struct {

	// 结果说明
	Msg string `json:"msg,omitempty"`

	// 结果码
	Code UnlockDomainOpenapiResponseBodyCodeEnum `json:"code,omitempty"`

	// 域名
	DomainName string `json:"domainName,omitempty"`

	// 域名ID
	DomainId string `json:"domainId,omitempty"`
}
