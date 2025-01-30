// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListLineGroupResponseStateEnum string

// List of State
const (
    ListLineGroupResponseStateEnumError ListLineGroupResponseStateEnum = "ERROR"
    ListLineGroupResponseStateEnumException ListLineGroupResponseStateEnum = "EXCEPTION"
    ListLineGroupResponseStateEnumForbidden ListLineGroupResponseStateEnum = "FORBIDDEN"
    ListLineGroupResponseStateEnumOk ListLineGroupResponseStateEnum = "OK"
)

type ListLineGroupResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ListLineGroupResponseStateEnum `json:"state,omitempty"`

	Body *ListLineGroupResponseBody `json:"body,omitempty"`
}
