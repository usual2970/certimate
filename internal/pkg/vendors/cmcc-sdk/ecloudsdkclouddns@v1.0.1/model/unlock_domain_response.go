// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type UnlockDomainResponseStateEnum string

// List of State
const (
    UnlockDomainResponseStateEnumError UnlockDomainResponseStateEnum = "ERROR"
    UnlockDomainResponseStateEnumException UnlockDomainResponseStateEnum = "EXCEPTION"
    UnlockDomainResponseStateEnumForbidden UnlockDomainResponseStateEnum = "FORBIDDEN"
    UnlockDomainResponseStateEnumOk UnlockDomainResponseStateEnum = "OK"
)

type UnlockDomainResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State UnlockDomainResponseStateEnum `json:"state,omitempty"`

	Body *UnlockDomainResponseBody `json:"body,omitempty"`
}
