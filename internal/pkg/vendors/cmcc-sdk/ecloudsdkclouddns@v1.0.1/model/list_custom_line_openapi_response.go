// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListCustomLineOpenapiResponseStateEnum string

// List of State
const (
    ListCustomLineOpenapiResponseStateEnumError ListCustomLineOpenapiResponseStateEnum = "ERROR"
    ListCustomLineOpenapiResponseStateEnumException ListCustomLineOpenapiResponseStateEnum = "EXCEPTION"
    ListCustomLineOpenapiResponseStateEnumForbidden ListCustomLineOpenapiResponseStateEnum = "FORBIDDEN"
    ListCustomLineOpenapiResponseStateEnumOk ListCustomLineOpenapiResponseStateEnum = "OK"
)

type ListCustomLineOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ListCustomLineOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *ListCustomLineOpenapiResponseBody `json:"body,omitempty"`
}
