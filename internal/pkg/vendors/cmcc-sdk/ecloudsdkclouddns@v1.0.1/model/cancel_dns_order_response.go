// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type CancelDnsOrderResponseStateEnum string

// List of State
const (
    CancelDnsOrderResponseStateEnumError CancelDnsOrderResponseStateEnum = "ERROR"
    CancelDnsOrderResponseStateEnumException CancelDnsOrderResponseStateEnum = "EXCEPTION"
    CancelDnsOrderResponseStateEnumForbidden CancelDnsOrderResponseStateEnum = "FORBIDDEN"
    CancelDnsOrderResponseStateEnumOk CancelDnsOrderResponseStateEnum = "OK"
)

type CancelDnsOrderResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State CancelDnsOrderResponseStateEnum `json:"state,omitempty"`

	Body *bool `json:"body,omitempty"`
}
