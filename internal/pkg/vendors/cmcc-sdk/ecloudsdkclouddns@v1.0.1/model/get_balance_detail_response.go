// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type GetBalanceDetailResponseStateEnum string

// List of State
const (
    GetBalanceDetailResponseStateEnumError GetBalanceDetailResponseStateEnum = "ERROR"
    GetBalanceDetailResponseStateEnumException GetBalanceDetailResponseStateEnum = "EXCEPTION"
    GetBalanceDetailResponseStateEnumForbidden GetBalanceDetailResponseStateEnum = "FORBIDDEN"
    GetBalanceDetailResponseStateEnumOk GetBalanceDetailResponseStateEnum = "OK"
)

type GetBalanceDetailResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State GetBalanceDetailResponseStateEnum `json:"state,omitempty"`

	Body *[]GetBalanceDetailResponseBody `json:"body,omitempty"`
}
