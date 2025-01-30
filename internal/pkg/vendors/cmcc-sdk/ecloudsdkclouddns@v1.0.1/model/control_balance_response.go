// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ControlBalanceResponseStateEnum string

// List of State
const (
    ControlBalanceResponseStateEnumError ControlBalanceResponseStateEnum = "ERROR"
    ControlBalanceResponseStateEnumException ControlBalanceResponseStateEnum = "EXCEPTION"
    ControlBalanceResponseStateEnumForbidden ControlBalanceResponseStateEnum = "FORBIDDEN"
    ControlBalanceResponseStateEnumOk ControlBalanceResponseStateEnum = "OK"
)

type ControlBalanceResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ControlBalanceResponseStateEnum `json:"state,omitempty"`
}
