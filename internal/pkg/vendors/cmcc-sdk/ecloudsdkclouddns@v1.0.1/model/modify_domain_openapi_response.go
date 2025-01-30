// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ModifyDomainOpenapiResponseStateEnum string

// List of State
const (
    ModifyDomainOpenapiResponseStateEnumError ModifyDomainOpenapiResponseStateEnum = "ERROR"
    ModifyDomainOpenapiResponseStateEnumException ModifyDomainOpenapiResponseStateEnum = "EXCEPTION"
    ModifyDomainOpenapiResponseStateEnumForbidden ModifyDomainOpenapiResponseStateEnum = "FORBIDDEN"
    ModifyDomainOpenapiResponseStateEnumOk ModifyDomainOpenapiResponseStateEnum = "OK"
)

type ModifyDomainOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ModifyDomainOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *ModifyDomainOpenapiResponseBody `json:"body,omitempty"`
}
