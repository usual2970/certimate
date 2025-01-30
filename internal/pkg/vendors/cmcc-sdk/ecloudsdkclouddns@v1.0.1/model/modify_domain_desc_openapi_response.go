// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ModifyDomainDescOpenapiResponseStateEnum string

// List of State
const (
    ModifyDomainDescOpenapiResponseStateEnumError ModifyDomainDescOpenapiResponseStateEnum = "ERROR"
    ModifyDomainDescOpenapiResponseStateEnumException ModifyDomainDescOpenapiResponseStateEnum = "EXCEPTION"
    ModifyDomainDescOpenapiResponseStateEnumForbidden ModifyDomainDescOpenapiResponseStateEnum = "FORBIDDEN"
    ModifyDomainDescOpenapiResponseStateEnumOk ModifyDomainDescOpenapiResponseStateEnum = "OK"
)

type ModifyDomainDescOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ModifyDomainDescOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *ModifyDomainDescOpenapiResponseBody `json:"body,omitempty"`
}
