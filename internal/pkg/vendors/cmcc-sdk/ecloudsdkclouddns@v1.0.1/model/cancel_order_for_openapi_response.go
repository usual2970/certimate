// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type CancelOrderForOpenapiResponseStateEnum string

// List of State
const (
    CancelOrderForOpenapiResponseStateEnumError CancelOrderForOpenapiResponseStateEnum = "ERROR"
    CancelOrderForOpenapiResponseStateEnumException CancelOrderForOpenapiResponseStateEnum = "EXCEPTION"
    CancelOrderForOpenapiResponseStateEnumForbidden CancelOrderForOpenapiResponseStateEnum = "FORBIDDEN"
    CancelOrderForOpenapiResponseStateEnumOk CancelOrderForOpenapiResponseStateEnum = "OK"
)

type CancelOrderForOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State CancelOrderForOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *bool `json:"body,omitempty"`
}
