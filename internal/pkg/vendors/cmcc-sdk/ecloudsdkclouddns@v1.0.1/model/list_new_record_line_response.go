// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListNewRecordLineResponseStateEnum string

// List of State
const (
    ListNewRecordLineResponseStateEnumError ListNewRecordLineResponseStateEnum = "ERROR"
    ListNewRecordLineResponseStateEnumException ListNewRecordLineResponseStateEnum = "EXCEPTION"
    ListNewRecordLineResponseStateEnumForbidden ListNewRecordLineResponseStateEnum = "FORBIDDEN"
    ListNewRecordLineResponseStateEnumOk ListNewRecordLineResponseStateEnum = "OK"
)

type ListNewRecordLineResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ListNewRecordLineResponseStateEnum `json:"state,omitempty"`

	Body []interface{} `json:"body,omitempty"`
}
