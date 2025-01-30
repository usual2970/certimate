// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type GetRecordOpenapiResponseStateEnum string

// List of State
const (
    GetRecordOpenapiResponseStateEnumError GetRecordOpenapiResponseStateEnum = "ERROR"
    GetRecordOpenapiResponseStateEnumException GetRecordOpenapiResponseStateEnum = "EXCEPTION"
    GetRecordOpenapiResponseStateEnumForbidden GetRecordOpenapiResponseStateEnum = "FORBIDDEN"
    GetRecordOpenapiResponseStateEnumOk GetRecordOpenapiResponseStateEnum = "OK"
)

type GetRecordOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State GetRecordOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *GetRecordOpenapiResponseBody `json:"body,omitempty"`
}
