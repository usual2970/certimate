// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ModifyBalanceRateResponseStateEnum string

// List of State
const (
    ModifyBalanceRateResponseStateEnumError ModifyBalanceRateResponseStateEnum = "ERROR"
    ModifyBalanceRateResponseStateEnumException ModifyBalanceRateResponseStateEnum = "EXCEPTION"
    ModifyBalanceRateResponseStateEnumForbidden ModifyBalanceRateResponseStateEnum = "FORBIDDEN"
    ModifyBalanceRateResponseStateEnumOk ModifyBalanceRateResponseStateEnum = "OK"
)

type ModifyBalanceRateResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ModifyBalanceRateResponseStateEnum `json:"state,omitempty"`
}
