// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ModifyOrderResponseStateEnum string

// List of State
const (
    ModifyOrderResponseStateEnumError ModifyOrderResponseStateEnum = "ERROR"
    ModifyOrderResponseStateEnumException ModifyOrderResponseStateEnum = "EXCEPTION"
    ModifyOrderResponseStateEnumForbidden ModifyOrderResponseStateEnum = "FORBIDDEN"
    ModifyOrderResponseStateEnumOk ModifyOrderResponseStateEnum = "OK"
)

type ModifyOrderResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ModifyOrderResponseStateEnum `json:"state,omitempty"`

	Body *ModifyOrderResponseBody `json:"body,omitempty"`
}
