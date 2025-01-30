// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type RemoveCustomLineOpenapiResponseStateEnum string

// List of State
const (
    RemoveCustomLineOpenapiResponseStateEnumError RemoveCustomLineOpenapiResponseStateEnum = "ERROR"
    RemoveCustomLineOpenapiResponseStateEnumException RemoveCustomLineOpenapiResponseStateEnum = "EXCEPTION"
    RemoveCustomLineOpenapiResponseStateEnumForbidden RemoveCustomLineOpenapiResponseStateEnum = "FORBIDDEN"
    RemoveCustomLineOpenapiResponseStateEnumOk RemoveCustomLineOpenapiResponseStateEnum = "OK"
)

type RemoveCustomLineOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State RemoveCustomLineOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *[]RemoveCustomLineOpenapiResponseBody `json:"body,omitempty"`
}
