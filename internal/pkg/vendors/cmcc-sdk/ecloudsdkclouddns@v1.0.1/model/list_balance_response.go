// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListBalanceResponseStateEnum string

// List of State
const (
    ListBalanceResponseStateEnumError ListBalanceResponseStateEnum = "ERROR"
    ListBalanceResponseStateEnumException ListBalanceResponseStateEnum = "EXCEPTION"
    ListBalanceResponseStateEnumForbidden ListBalanceResponseStateEnum = "FORBIDDEN"
    ListBalanceResponseStateEnumOk ListBalanceResponseStateEnum = "OK"
)

type ListBalanceResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ListBalanceResponseStateEnum `json:"state,omitempty"`

	Body *ListBalanceResponseBody `json:"body,omitempty"`
}
