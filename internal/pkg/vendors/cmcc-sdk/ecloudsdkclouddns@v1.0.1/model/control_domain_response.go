// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ControlDomainResponseStateEnum string

// List of State
const (
    ControlDomainResponseStateEnumError ControlDomainResponseStateEnum = "ERROR"
    ControlDomainResponseStateEnumException ControlDomainResponseStateEnum = "EXCEPTION"
    ControlDomainResponseStateEnumForbidden ControlDomainResponseStateEnum = "FORBIDDEN"
    ControlDomainResponseStateEnumOk ControlDomainResponseStateEnum = "OK"
)

type ControlDomainResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ControlDomainResponseStateEnum `json:"state,omitempty"`

	Body *[]ControlDomainResponseBody `json:"body,omitempty"`
}
