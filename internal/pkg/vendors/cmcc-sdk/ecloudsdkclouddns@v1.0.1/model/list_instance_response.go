// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListInstanceResponseStateEnum string

// List of State
const (
    ListInstanceResponseStateEnumError ListInstanceResponseStateEnum = "ERROR"
    ListInstanceResponseStateEnumException ListInstanceResponseStateEnum = "EXCEPTION"
    ListInstanceResponseStateEnumForbidden ListInstanceResponseStateEnum = "FORBIDDEN"
    ListInstanceResponseStateEnumOk ListInstanceResponseStateEnum = "OK"
)

type ListInstanceResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ListInstanceResponseStateEnum `json:"state,omitempty"`

	Body *ListInstanceResponseBody `json:"body,omitempty"`
}
