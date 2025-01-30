// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListNewRecordLineOpenapiResponseStateEnum string

// List of State
const (
    ListNewRecordLineOpenapiResponseStateEnumError ListNewRecordLineOpenapiResponseStateEnum = "ERROR"
    ListNewRecordLineOpenapiResponseStateEnumException ListNewRecordLineOpenapiResponseStateEnum = "EXCEPTION"
    ListNewRecordLineOpenapiResponseStateEnumForbidden ListNewRecordLineOpenapiResponseStateEnum = "FORBIDDEN"
    ListNewRecordLineOpenapiResponseStateEnumOk ListNewRecordLineOpenapiResponseStateEnum = "OK"
)

type ListNewRecordLineOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ListNewRecordLineOpenapiResponseStateEnum `json:"state,omitempty"`

	Body []interface{} `json:"body,omitempty"`
}
