// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListNsOpenapiResponseStateEnum string

// List of State
const (
    ListNsOpenapiResponseStateEnumError ListNsOpenapiResponseStateEnum = "ERROR"
    ListNsOpenapiResponseStateEnumException ListNsOpenapiResponseStateEnum = "EXCEPTION"
    ListNsOpenapiResponseStateEnumForbidden ListNsOpenapiResponseStateEnum = "FORBIDDEN"
    ListNsOpenapiResponseStateEnumOk ListNsOpenapiResponseStateEnum = "OK"
)

type ListNsOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ListNsOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *[]ListNsOpenapiResponseBody `json:"body,omitempty"`
}
