// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type CreateLineGroupOpenapiResponseStateEnum string

// List of State
const (
    CreateLineGroupOpenapiResponseStateEnumError CreateLineGroupOpenapiResponseStateEnum = "ERROR"
    CreateLineGroupOpenapiResponseStateEnumException CreateLineGroupOpenapiResponseStateEnum = "EXCEPTION"
    CreateLineGroupOpenapiResponseStateEnumForbidden CreateLineGroupOpenapiResponseStateEnum = "FORBIDDEN"
    CreateLineGroupOpenapiResponseStateEnumOk CreateLineGroupOpenapiResponseStateEnum = "OK"
)

type CreateLineGroupOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State CreateLineGroupOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *CreateLineGroupOpenapiResponseBody `json:"body,omitempty"`
}
