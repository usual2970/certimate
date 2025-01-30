// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type TakeOverRefreshDomainResponseStateEnum string

// List of State
const (
    TakeOverRefreshDomainResponseStateEnumError TakeOverRefreshDomainResponseStateEnum = "ERROR"
    TakeOverRefreshDomainResponseStateEnumException TakeOverRefreshDomainResponseStateEnum = "EXCEPTION"
    TakeOverRefreshDomainResponseStateEnumForbidden TakeOverRefreshDomainResponseStateEnum = "FORBIDDEN"
    TakeOverRefreshDomainResponseStateEnumOk TakeOverRefreshDomainResponseStateEnum = "OK"
)

type TakeOverRefreshDomainResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State TakeOverRefreshDomainResponseStateEnum `json:"state,omitempty"`

	Body *TakeOverRefreshDomainResponseBody `json:"body,omitempty"`
}
