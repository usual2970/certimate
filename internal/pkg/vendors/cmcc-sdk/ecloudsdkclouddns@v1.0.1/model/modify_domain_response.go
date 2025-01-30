// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ModifyDomainResponseStateEnum string

// List of State
const (
    ModifyDomainResponseStateEnumError ModifyDomainResponseStateEnum = "ERROR"
    ModifyDomainResponseStateEnumException ModifyDomainResponseStateEnum = "EXCEPTION"
    ModifyDomainResponseStateEnumForbidden ModifyDomainResponseStateEnum = "FORBIDDEN"
    ModifyDomainResponseStateEnumOk ModifyDomainResponseStateEnum = "OK"
)

type ModifyDomainResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ModifyDomainResponseStateEnum `json:"state,omitempty"`

	Body *ModifyDomainResponseBody `json:"body,omitempty"`
}
