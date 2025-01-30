// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

type ListDomainForOpenapiResponseStateEnum string

// List of State
const (
	ListDomainForOpenapiResponseStateEnumError     ListDomainForOpenapiResponseStateEnum = "ERROR"
	ListDomainForOpenapiResponseStateEnumException ListDomainForOpenapiResponseStateEnum = "EXCEPTION"
	ListDomainForOpenapiResponseStateEnumForbidden ListDomainForOpenapiResponseStateEnum = "FORBIDDEN"
	ListDomainForOpenapiResponseStateEnumOk        ListDomainForOpenapiResponseStateEnum = "OK"
)

type ListDomainForOpenapiResponse struct {
	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ListDomainForOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *ListDomainForOpenapiResponseBody `json:"body,omitempty"`
}
