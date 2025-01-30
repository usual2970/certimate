// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type RemoveLineGroupOpenapiResponseStateEnum string

// List of State
const (
    RemoveLineGroupOpenapiResponseStateEnumError RemoveLineGroupOpenapiResponseStateEnum = "ERROR"
    RemoveLineGroupOpenapiResponseStateEnumException RemoveLineGroupOpenapiResponseStateEnum = "EXCEPTION"
    RemoveLineGroupOpenapiResponseStateEnumForbidden RemoveLineGroupOpenapiResponseStateEnum = "FORBIDDEN"
    RemoveLineGroupOpenapiResponseStateEnumOk RemoveLineGroupOpenapiResponseStateEnum = "OK"
)

type RemoveLineGroupOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State RemoveLineGroupOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *[]RemoveLineGroupOpenapiResponseBody `json:"body,omitempty"`
}
