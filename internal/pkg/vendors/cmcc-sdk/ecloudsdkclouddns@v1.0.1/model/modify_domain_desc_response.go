// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ModifyDomainDescResponseStateEnum string

// List of State
const (
    ModifyDomainDescResponseStateEnumError ModifyDomainDescResponseStateEnum = "ERROR"
    ModifyDomainDescResponseStateEnumException ModifyDomainDescResponseStateEnum = "EXCEPTION"
    ModifyDomainDescResponseStateEnumForbidden ModifyDomainDescResponseStateEnum = "FORBIDDEN"
    ModifyDomainDescResponseStateEnumOk ModifyDomainDescResponseStateEnum = "OK"
)

type ModifyDomainDescResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ModifyDomainDescResponseStateEnum `json:"state,omitempty"`

	Body *ModifyDomainDescResponseBody `json:"body,omitempty"`
}
