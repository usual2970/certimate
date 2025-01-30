// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type UnlockDomainOpenapiResponseStateEnum string

// List of State
const (
    UnlockDomainOpenapiResponseStateEnumError UnlockDomainOpenapiResponseStateEnum = "ERROR"
    UnlockDomainOpenapiResponseStateEnumException UnlockDomainOpenapiResponseStateEnum = "EXCEPTION"
    UnlockDomainOpenapiResponseStateEnumForbidden UnlockDomainOpenapiResponseStateEnum = "FORBIDDEN"
    UnlockDomainOpenapiResponseStateEnumOk UnlockDomainOpenapiResponseStateEnum = "OK"
)

type UnlockDomainOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State UnlockDomainOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *UnlockDomainOpenapiResponseBody `json:"body,omitempty"`
}
