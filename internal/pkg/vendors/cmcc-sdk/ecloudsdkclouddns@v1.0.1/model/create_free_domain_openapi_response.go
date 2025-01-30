// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type CreateFreeDomainOpenapiResponseStateEnum string

// List of State
const (
    CreateFreeDomainOpenapiResponseStateEnumError CreateFreeDomainOpenapiResponseStateEnum = "ERROR"
    CreateFreeDomainOpenapiResponseStateEnumException CreateFreeDomainOpenapiResponseStateEnum = "EXCEPTION"
    CreateFreeDomainOpenapiResponseStateEnumForbidden CreateFreeDomainOpenapiResponseStateEnum = "FORBIDDEN"
    CreateFreeDomainOpenapiResponseStateEnumOk CreateFreeDomainOpenapiResponseStateEnum = "OK"
)

type CreateFreeDomainOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State CreateFreeDomainOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *CreateFreeDomainOpenapiResponseBody `json:"body,omitempty"`
}
