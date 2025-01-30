// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type TakeOverRefreshDomainOpenapiResponseStateEnum string

// List of State
const (
    TakeOverRefreshDomainOpenapiResponseStateEnumError TakeOverRefreshDomainOpenapiResponseStateEnum = "ERROR"
    TakeOverRefreshDomainOpenapiResponseStateEnumException TakeOverRefreshDomainOpenapiResponseStateEnum = "EXCEPTION"
    TakeOverRefreshDomainOpenapiResponseStateEnumForbidden TakeOverRefreshDomainOpenapiResponseStateEnum = "FORBIDDEN"
    TakeOverRefreshDomainOpenapiResponseStateEnumOk TakeOverRefreshDomainOpenapiResponseStateEnum = "OK"
)

type TakeOverRefreshDomainOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State TakeOverRefreshDomainOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *TakeOverRefreshDomainOpenapiResponseBody `json:"body,omitempty"`
}
