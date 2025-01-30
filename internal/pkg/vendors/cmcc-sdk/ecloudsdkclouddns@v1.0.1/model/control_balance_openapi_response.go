// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ControlBalanceOpenapiResponseStateEnum string

// List of State
const (
    ControlBalanceOpenapiResponseStateEnumError ControlBalanceOpenapiResponseStateEnum = "ERROR"
    ControlBalanceOpenapiResponseStateEnumException ControlBalanceOpenapiResponseStateEnum = "EXCEPTION"
    ControlBalanceOpenapiResponseStateEnumForbidden ControlBalanceOpenapiResponseStateEnum = "FORBIDDEN"
    ControlBalanceOpenapiResponseStateEnumOk ControlBalanceOpenapiResponseStateEnum = "OK"
)

type ControlBalanceOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ControlBalanceOpenapiResponseStateEnum `json:"state,omitempty"`
}
