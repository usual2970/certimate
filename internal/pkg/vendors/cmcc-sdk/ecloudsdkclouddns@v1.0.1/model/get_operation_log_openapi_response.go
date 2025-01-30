// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type GetOperationLogOpenapiResponseStateEnum string

// List of State
const (
    GetOperationLogOpenapiResponseStateEnumError GetOperationLogOpenapiResponseStateEnum = "ERROR"
    GetOperationLogOpenapiResponseStateEnumException GetOperationLogOpenapiResponseStateEnum = "EXCEPTION"
    GetOperationLogOpenapiResponseStateEnumForbidden GetOperationLogOpenapiResponseStateEnum = "FORBIDDEN"
    GetOperationLogOpenapiResponseStateEnumOk GetOperationLogOpenapiResponseStateEnum = "OK"
)

type GetOperationLogOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State GetOperationLogOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *GetOperationLogOpenapiResponseBody `json:"body,omitempty"`
}
