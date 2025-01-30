// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type RenewProductOpenapiResponseStateEnum string

// List of State
const (
    RenewProductOpenapiResponseStateEnumError RenewProductOpenapiResponseStateEnum = "ERROR"
    RenewProductOpenapiResponseStateEnumException RenewProductOpenapiResponseStateEnum = "EXCEPTION"
    RenewProductOpenapiResponseStateEnumForbidden RenewProductOpenapiResponseStateEnum = "FORBIDDEN"
    RenewProductOpenapiResponseStateEnumOk RenewProductOpenapiResponseStateEnum = "OK"
)

type RenewProductOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State RenewProductOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *RenewProductOpenapiResponseBody `json:"body,omitempty"`
}
