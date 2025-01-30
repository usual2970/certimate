// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ModifyLineGroupOpenapiResponseStateEnum string

// List of State
const (
    ModifyLineGroupOpenapiResponseStateEnumError ModifyLineGroupOpenapiResponseStateEnum = "ERROR"
    ModifyLineGroupOpenapiResponseStateEnumException ModifyLineGroupOpenapiResponseStateEnum = "EXCEPTION"
    ModifyLineGroupOpenapiResponseStateEnumForbidden ModifyLineGroupOpenapiResponseStateEnum = "FORBIDDEN"
    ModifyLineGroupOpenapiResponseStateEnumOk ModifyLineGroupOpenapiResponseStateEnum = "OK"
)

type ModifyLineGroupOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ModifyLineGroupOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *ModifyLineGroupOpenapiResponseBody `json:"body,omitempty"`
}
