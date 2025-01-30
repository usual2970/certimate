// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type CreateFreeDomainResponseStateEnum string

// List of State
const (
    CreateFreeDomainResponseStateEnumError CreateFreeDomainResponseStateEnum = "ERROR"
    CreateFreeDomainResponseStateEnumException CreateFreeDomainResponseStateEnum = "EXCEPTION"
    CreateFreeDomainResponseStateEnumForbidden CreateFreeDomainResponseStateEnum = "FORBIDDEN"
    CreateFreeDomainResponseStateEnumOk CreateFreeDomainResponseStateEnum = "OK"
)

type CreateFreeDomainResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State CreateFreeDomainResponseStateEnum `json:"state,omitempty"`

	Body *CreateFreeDomainResponseBody `json:"body,omitempty"`
}
