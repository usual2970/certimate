// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

type ListDomainResponseStateEnum string

// List of State
const (
	ListDomainResponseStateEnumError     ListDomainResponseStateEnum = "ERROR"
	ListDomainResponseStateEnumException ListDomainResponseStateEnum = "EXCEPTION"
	ListDomainResponseStateEnumForbidden ListDomainResponseStateEnum = "FORBIDDEN"
	ListDomainResponseStateEnumOk        ListDomainResponseStateEnum = "OK"
)

type ListDomainResponse struct {
	RequestId string `json:"requestId,omitempty"`

	State ListDomainResponseStateEnum `json:"state,omitempty"`

	Body *ListDomainResponseBody `json:"body,omitempty"`
}
