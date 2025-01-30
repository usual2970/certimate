// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListInstanceForOpenapiResponseStateEnum string

// List of State
const (
    ListInstanceForOpenapiResponseStateEnumError ListInstanceForOpenapiResponseStateEnum = "ERROR"
    ListInstanceForOpenapiResponseStateEnumException ListInstanceForOpenapiResponseStateEnum = "EXCEPTION"
    ListInstanceForOpenapiResponseStateEnumForbidden ListInstanceForOpenapiResponseStateEnum = "FORBIDDEN"
    ListInstanceForOpenapiResponseStateEnumOk ListInstanceForOpenapiResponseStateEnum = "OK"
)

type ListInstanceForOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ListInstanceForOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *ListInstanceForOpenapiResponseBody `json:"body,omitempty"`
}
