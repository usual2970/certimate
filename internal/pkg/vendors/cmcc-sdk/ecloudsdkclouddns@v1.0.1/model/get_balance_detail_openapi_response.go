// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type GetBalanceDetailOpenapiResponseStateEnum string

// List of State
const (
    GetBalanceDetailOpenapiResponseStateEnumError GetBalanceDetailOpenapiResponseStateEnum = "ERROR"
    GetBalanceDetailOpenapiResponseStateEnumException GetBalanceDetailOpenapiResponseStateEnum = "EXCEPTION"
    GetBalanceDetailOpenapiResponseStateEnumForbidden GetBalanceDetailOpenapiResponseStateEnum = "FORBIDDEN"
    GetBalanceDetailOpenapiResponseStateEnumOk GetBalanceDetailOpenapiResponseStateEnum = "OK"
)

type GetBalanceDetailOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State GetBalanceDetailOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *[]GetBalanceDetailOpenapiResponseBody `json:"body,omitempty"`
}
