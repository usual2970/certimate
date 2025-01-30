// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ModifyBalanceRateOpenapiResponseStateEnum string

// List of State
const (
    ModifyBalanceRateOpenapiResponseStateEnumError ModifyBalanceRateOpenapiResponseStateEnum = "ERROR"
    ModifyBalanceRateOpenapiResponseStateEnumException ModifyBalanceRateOpenapiResponseStateEnum = "EXCEPTION"
    ModifyBalanceRateOpenapiResponseStateEnumForbidden ModifyBalanceRateOpenapiResponseStateEnum = "FORBIDDEN"
    ModifyBalanceRateOpenapiResponseStateEnumOk ModifyBalanceRateOpenapiResponseStateEnum = "OK"
)

type ModifyBalanceRateOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ModifyBalanceRateOpenapiResponseStateEnum `json:"state,omitempty"`
}
