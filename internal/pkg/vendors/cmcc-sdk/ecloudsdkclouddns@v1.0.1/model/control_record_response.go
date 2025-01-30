// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ControlRecordResponseStateEnum string

// List of State
const (
    ControlRecordResponseStateEnumError ControlRecordResponseStateEnum = "ERROR"
    ControlRecordResponseStateEnumException ControlRecordResponseStateEnum = "EXCEPTION"
    ControlRecordResponseStateEnumForbidden ControlRecordResponseStateEnum = "FORBIDDEN"
    ControlRecordResponseStateEnumOk ControlRecordResponseStateEnum = "OK"
)

type ControlRecordResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ControlRecordResponseStateEnum `json:"state,omitempty"`

	Body *[]ControlRecordResponseBody `json:"body,omitempty"`
}
