// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type RemoveLineGroupResponseStateEnum string

// List of State
const (
    RemoveLineGroupResponseStateEnumError RemoveLineGroupResponseStateEnum = "ERROR"
    RemoveLineGroupResponseStateEnumException RemoveLineGroupResponseStateEnum = "EXCEPTION"
    RemoveLineGroupResponseStateEnumForbidden RemoveLineGroupResponseStateEnum = "FORBIDDEN"
    RemoveLineGroupResponseStateEnumOk RemoveLineGroupResponseStateEnum = "OK"
)

type RemoveLineGroupResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State RemoveLineGroupResponseStateEnum `json:"state,omitempty"`

	Body *[]RemoveLineGroupResponseBody `json:"body,omitempty"`
}
