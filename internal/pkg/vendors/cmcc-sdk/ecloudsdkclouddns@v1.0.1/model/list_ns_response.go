// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListNsResponseStateEnum string

// List of State
const (
    ListNsResponseStateEnumError ListNsResponseStateEnum = "ERROR"
    ListNsResponseStateEnumException ListNsResponseStateEnum = "EXCEPTION"
    ListNsResponseStateEnumForbidden ListNsResponseStateEnum = "FORBIDDEN"
    ListNsResponseStateEnumOk ListNsResponseStateEnum = "OK"
)

type ListNsResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State ListNsResponseStateEnum `json:"state,omitempty"`

	Body *[]ListNsResponseBody `json:"body,omitempty"`
}
