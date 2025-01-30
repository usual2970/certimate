// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ControlDomainOpenapiResponseStateEnum string

// List of State
const (
    ControlDomainOpenapiResponseStateEnumError ControlDomainOpenapiResponseStateEnum = "ERROR"
    ControlDomainOpenapiResponseStateEnumException ControlDomainOpenapiResponseStateEnum = "EXCEPTION"
    ControlDomainOpenapiResponseStateEnumForbidden ControlDomainOpenapiResponseStateEnum = "FORBIDDEN"
    ControlDomainOpenapiResponseStateEnumOk ControlDomainOpenapiResponseStateEnum = "OK"
)

type ControlDomainOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ControlDomainOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *[]ControlDomainOpenapiResponseBody `json:"body,omitempty"`
}
