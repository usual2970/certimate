// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type RemoveCustomLineOpenapiResponseBodyCodeEnum string

// List of Code
const (
    RemoveCustomLineOpenapiResponseBodyCodeEnumError RemoveCustomLineOpenapiResponseBodyCodeEnum = "ERROR"
    RemoveCustomLineOpenapiResponseBodyCodeEnumSuccess RemoveCustomLineOpenapiResponseBodyCodeEnum = "SUCCESS"
)

type RemoveCustomLineOpenapiResponseBody struct {

	// 结果说明
	Msg string `json:"msg,omitempty"`

	// 结果码
	Code RemoveCustomLineOpenapiResponseBodyCodeEnum `json:"code,omitempty"`

	// 自定义线路ID
	LineId string `json:"lineId,omitempty"`
}
