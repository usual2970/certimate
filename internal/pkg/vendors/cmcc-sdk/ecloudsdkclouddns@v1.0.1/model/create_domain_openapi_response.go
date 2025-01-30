// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type CreateDomainOpenapiResponseStateEnum string

// List of State
const (
    CreateDomainOpenapiResponseStateEnumError CreateDomainOpenapiResponseStateEnum = "ERROR"
    CreateDomainOpenapiResponseStateEnumException CreateDomainOpenapiResponseStateEnum = "EXCEPTION"
    CreateDomainOpenapiResponseStateEnumForbidden CreateDomainOpenapiResponseStateEnum = "FORBIDDEN"
    CreateDomainOpenapiResponseStateEnumOk CreateDomainOpenapiResponseStateEnum = "OK"
)

type CreateDomainOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State CreateDomainOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *[]CreateDomainOpenapiResponseBody `json:"body,omitempty"`
}
