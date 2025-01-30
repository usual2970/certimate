// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type CreateOrderOpenapiResponseStateEnum string

// List of State
const (
    CreateOrderOpenapiResponseStateEnumError CreateOrderOpenapiResponseStateEnum = "ERROR"
    CreateOrderOpenapiResponseStateEnumException CreateOrderOpenapiResponseStateEnum = "EXCEPTION"
    CreateOrderOpenapiResponseStateEnumForbidden CreateOrderOpenapiResponseStateEnum = "FORBIDDEN"
    CreateOrderOpenapiResponseStateEnumOk CreateOrderOpenapiResponseStateEnum = "OK"
)

type CreateOrderOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State CreateOrderOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *CreateOrderOpenapiResponseBody `json:"body,omitempty"`
}
