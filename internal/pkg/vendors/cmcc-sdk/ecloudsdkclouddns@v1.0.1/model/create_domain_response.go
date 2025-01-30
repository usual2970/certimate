// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type CreateDomainResponseStateEnum string

// List of State
const (
    CreateDomainResponseStateEnumError CreateDomainResponseStateEnum = "ERROR"
    CreateDomainResponseStateEnumException CreateDomainResponseStateEnum = "EXCEPTION"
    CreateDomainResponseStateEnumForbidden CreateDomainResponseStateEnum = "FORBIDDEN"
    CreateDomainResponseStateEnumOk CreateDomainResponseStateEnum = "OK"
)

type CreateDomainResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State CreateDomainResponseStateEnum `json:"state,omitempty"`

	Body *[]CreateDomainResponseBody `json:"body,omitempty"`
}
