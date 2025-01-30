// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type LockDomainOpenapiResponseStateEnum string

// List of State
const (
    LockDomainOpenapiResponseStateEnumError LockDomainOpenapiResponseStateEnum = "ERROR"
    LockDomainOpenapiResponseStateEnumException LockDomainOpenapiResponseStateEnum = "EXCEPTION"
    LockDomainOpenapiResponseStateEnumForbidden LockDomainOpenapiResponseStateEnum = "FORBIDDEN"
    LockDomainOpenapiResponseStateEnumOk LockDomainOpenapiResponseStateEnum = "OK"
)

type LockDomainOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State LockDomainOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *LockDomainOpenapiResponseBody `json:"body,omitempty"`
}
