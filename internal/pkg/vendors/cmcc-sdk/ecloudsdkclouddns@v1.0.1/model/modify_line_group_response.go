// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ModifyLineGroupResponseStateEnum string

// List of State
const (
    ModifyLineGroupResponseStateEnumError ModifyLineGroupResponseStateEnum = "ERROR"
    ModifyLineGroupResponseStateEnumException ModifyLineGroupResponseStateEnum = "EXCEPTION"
    ModifyLineGroupResponseStateEnumForbidden ModifyLineGroupResponseStateEnum = "FORBIDDEN"
    ModifyLineGroupResponseStateEnumOk ModifyLineGroupResponseStateEnum = "OK"
)

type ModifyLineGroupResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ModifyLineGroupResponseStateEnum `json:"state,omitempty"`

	Body *ModifyLineGroupResponseBody `json:"body,omitempty"`
}
