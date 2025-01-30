// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListLineGroupOpenapiResponseStateEnum string

// List of State
const (
    ListLineGroupOpenapiResponseStateEnumError ListLineGroupOpenapiResponseStateEnum = "ERROR"
    ListLineGroupOpenapiResponseStateEnumException ListLineGroupOpenapiResponseStateEnum = "EXCEPTION"
    ListLineGroupOpenapiResponseStateEnumForbidden ListLineGroupOpenapiResponseStateEnum = "FORBIDDEN"
    ListLineGroupOpenapiResponseStateEnumOk ListLineGroupOpenapiResponseStateEnum = "OK"
)

type ListLineGroupOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ListLineGroupOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *ListLineGroupOpenapiResponseBody `json:"body,omitempty"`
}
