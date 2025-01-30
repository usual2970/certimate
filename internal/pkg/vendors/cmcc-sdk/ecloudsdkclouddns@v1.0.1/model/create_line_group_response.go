// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type CreateLineGroupResponseStateEnum string

// List of State
const (
    CreateLineGroupResponseStateEnumError CreateLineGroupResponseStateEnum = "ERROR"
    CreateLineGroupResponseStateEnumException CreateLineGroupResponseStateEnum = "EXCEPTION"
    CreateLineGroupResponseStateEnumForbidden CreateLineGroupResponseStateEnum = "FORBIDDEN"
    CreateLineGroupResponseStateEnumOk CreateLineGroupResponseStateEnum = "OK"
)

type CreateLineGroupResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State CreateLineGroupResponseStateEnum `json:"state,omitempty"`

	Body *CreateLineGroupResponseBody `json:"body,omitempty"`
}
