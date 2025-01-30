// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type CreateCustomLineOpenapiResponseStateEnum string

// List of State
const (
    CreateCustomLineOpenapiResponseStateEnumError CreateCustomLineOpenapiResponseStateEnum = "ERROR"
    CreateCustomLineOpenapiResponseStateEnumException CreateCustomLineOpenapiResponseStateEnum = "EXCEPTION"
    CreateCustomLineOpenapiResponseStateEnumForbidden CreateCustomLineOpenapiResponseStateEnum = "FORBIDDEN"
    CreateCustomLineOpenapiResponseStateEnumOk CreateCustomLineOpenapiResponseStateEnum = "OK"
)

type CreateCustomLineOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State CreateCustomLineOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *CreateCustomLineOpenapiResponseBody `json:"body,omitempty"`
}
