// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type RemoveLineGroupResponseBodyCodeEnum string

// List of Code
const (
    RemoveLineGroupResponseBodyCodeEnumError RemoveLineGroupResponseBodyCodeEnum = "ERROR"
    RemoveLineGroupResponseBodyCodeEnumSuccess RemoveLineGroupResponseBodyCodeEnum = "SUCCESS"
)

type RemoveLineGroupResponseBody struct {

	// 结果说明
	Msg string `json:"msg,omitempty"`

	// 结果码
	Code RemoveLineGroupResponseBodyCodeEnum `json:"code,omitempty"`

	// 线路分组ID
	GroupId string `json:"groupId,omitempty"`
}
