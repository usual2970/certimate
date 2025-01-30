// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type DeleteDomainResponseStateEnum string

// List of State
const (
    DeleteDomainResponseStateEnumError DeleteDomainResponseStateEnum = "ERROR"
    DeleteDomainResponseStateEnumException DeleteDomainResponseStateEnum = "EXCEPTION"
    DeleteDomainResponseStateEnumForbidden DeleteDomainResponseStateEnum = "FORBIDDEN"
    DeleteDomainResponseStateEnumOk DeleteDomainResponseStateEnum = "OK"
)

type DeleteDomainResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State DeleteDomainResponseStateEnum `json:"state,omitempty"`

	Body *[]DeleteDomainResponseBody `json:"body,omitempty"`
}
