// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type DeleteDomainOpenapiResponseBodyCodeEnum string

// List of Code
const (
    DeleteDomainOpenapiResponseBodyCodeEnumError DeleteDomainOpenapiResponseBodyCodeEnum = "ERROR"
    DeleteDomainOpenapiResponseBodyCodeEnumSuccess DeleteDomainOpenapiResponseBodyCodeEnum = "SUCCESS"
)

type DeleteDomainOpenapiResponseBody struct {

	// 结果说明
	Msg string `json:"msg,omitempty"`

	// 结果码
	Code DeleteDomainOpenapiResponseBodyCodeEnum `json:"code,omitempty"`

	// 域名
	DomainName string `json:"domainName,omitempty"`

	// 域名ID
	DomainId string `json:"domainId,omitempty"`
}
