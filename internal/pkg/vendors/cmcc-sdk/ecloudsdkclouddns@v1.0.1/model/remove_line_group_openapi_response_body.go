// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type RemoveLineGroupOpenapiResponseBodyCodeEnum string

// List of Code
const (
    RemoveLineGroupOpenapiResponseBodyCodeEnumError RemoveLineGroupOpenapiResponseBodyCodeEnum = "ERROR"
    RemoveLineGroupOpenapiResponseBodyCodeEnumSuccess RemoveLineGroupOpenapiResponseBodyCodeEnum = "SUCCESS"
)

type RemoveLineGroupOpenapiResponseBody struct {

	// 结果说明
	Msg string `json:"msg,omitempty"`

	// 结果码
	Code RemoveLineGroupOpenapiResponseBodyCodeEnum `json:"code,omitempty"`

	// 线路分组ID
	GroupId string `json:"groupId,omitempty"`
}
