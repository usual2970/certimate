// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ControlDomainOpenapiResponseBodyCodeEnum string

// List of Code
const (
    ControlDomainOpenapiResponseBodyCodeEnumError ControlDomainOpenapiResponseBodyCodeEnum = "ERROR"
    ControlDomainOpenapiResponseBodyCodeEnumSuccess ControlDomainOpenapiResponseBodyCodeEnum = "SUCCESS"
)

type ControlDomainOpenapiResponseBody struct {

	// 结果说明
	Msg string `json:"msg,omitempty"`

	// 结果码
	Code ControlDomainOpenapiResponseBodyCodeEnum `json:"code,omitempty"`

	// 域名
	DomainName string `json:"domainName,omitempty"`

	// 域名ID
	DomainId string `json:"domainId,omitempty"`
}
