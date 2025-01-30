// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type DeleteDomainOpenapiResponseStateEnum string

// List of State
const (
    DeleteDomainOpenapiResponseStateEnumError DeleteDomainOpenapiResponseStateEnum = "ERROR"
    DeleteDomainOpenapiResponseStateEnumException DeleteDomainOpenapiResponseStateEnum = "EXCEPTION"
    DeleteDomainOpenapiResponseStateEnumForbidden DeleteDomainOpenapiResponseStateEnum = "FORBIDDEN"
    DeleteDomainOpenapiResponseStateEnumOk DeleteDomainOpenapiResponseStateEnum = "OK"
)

type DeleteDomainOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State DeleteDomainOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *[]DeleteDomainOpenapiResponseBody `json:"body,omitempty"`
}
