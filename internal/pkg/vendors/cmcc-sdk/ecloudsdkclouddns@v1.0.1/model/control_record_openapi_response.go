// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ControlRecordOpenapiResponseStateEnum string

// List of State
const (
    ControlRecordOpenapiResponseStateEnumError ControlRecordOpenapiResponseStateEnum = "ERROR"
    ControlRecordOpenapiResponseStateEnumException ControlRecordOpenapiResponseStateEnum = "EXCEPTION"
    ControlRecordOpenapiResponseStateEnumForbidden ControlRecordOpenapiResponseStateEnum = "FORBIDDEN"
    ControlRecordOpenapiResponseStateEnumOk ControlRecordOpenapiResponseStateEnum = "OK"
)

type ControlRecordOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ControlRecordOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *[]ControlRecordOpenapiResponseBody `json:"body,omitempty"`
}
