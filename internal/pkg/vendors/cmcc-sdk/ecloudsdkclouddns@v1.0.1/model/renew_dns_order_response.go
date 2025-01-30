// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type RenewDnsOrderResponseStateEnum string

// List of State
const (
    RenewDnsOrderResponseStateEnumError RenewDnsOrderResponseStateEnum = "ERROR"
    RenewDnsOrderResponseStateEnumException RenewDnsOrderResponseStateEnum = "EXCEPTION"
    RenewDnsOrderResponseStateEnumForbidden RenewDnsOrderResponseStateEnum = "FORBIDDEN"
    RenewDnsOrderResponseStateEnumOk RenewDnsOrderResponseStateEnum = "OK"
)

type RenewDnsOrderResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State RenewDnsOrderResponseStateEnum `json:"state,omitempty"`

	Body *RenewDnsOrderResponseBody `json:"body,omitempty"`
}
