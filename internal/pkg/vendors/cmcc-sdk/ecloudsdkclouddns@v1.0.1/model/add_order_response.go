// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type AddOrderResponseStateEnum string

// List of State
const (
    AddOrderResponseStateEnumError AddOrderResponseStateEnum = "ERROR"
    AddOrderResponseStateEnumException AddOrderResponseStateEnum = "EXCEPTION"
    AddOrderResponseStateEnumForbidden AddOrderResponseStateEnum = "FORBIDDEN"
    AddOrderResponseStateEnumOk AddOrderResponseStateEnum = "OK"
)

type AddOrderResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State AddOrderResponseStateEnum `json:"state,omitempty"`

	Body *AddOrderResponseBody `json:"body,omitempty"`
}
