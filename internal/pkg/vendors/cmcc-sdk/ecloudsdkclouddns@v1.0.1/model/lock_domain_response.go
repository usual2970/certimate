// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type LockDomainResponseStateEnum string

// List of State
const (
    LockDomainResponseStateEnumError LockDomainResponseStateEnum = "ERROR"
    LockDomainResponseStateEnumException LockDomainResponseStateEnum = "EXCEPTION"
    LockDomainResponseStateEnumForbidden LockDomainResponseStateEnum = "FORBIDDEN"
    LockDomainResponseStateEnumOk LockDomainResponseStateEnum = "OK"
)

type LockDomainResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State LockDomainResponseStateEnum `json:"state,omitempty"`

	Body *LockDomainResponseBody `json:"body,omitempty"`
}
