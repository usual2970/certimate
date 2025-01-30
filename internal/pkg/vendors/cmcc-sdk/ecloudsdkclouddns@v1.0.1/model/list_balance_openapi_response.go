// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListBalanceOpenapiResponseStateEnum string

// List of State
const (
    ListBalanceOpenapiResponseStateEnumError ListBalanceOpenapiResponseStateEnum = "ERROR"
    ListBalanceOpenapiResponseStateEnumException ListBalanceOpenapiResponseStateEnum = "EXCEPTION"
    ListBalanceOpenapiResponseStateEnumForbidden ListBalanceOpenapiResponseStateEnum = "FORBIDDEN"
    ListBalanceOpenapiResponseStateEnumOk ListBalanceOpenapiResponseStateEnum = "OK"
)

type ListBalanceOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ListBalanceOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *ListBalanceOpenapiResponseBody `json:"body,omitempty"`
}
