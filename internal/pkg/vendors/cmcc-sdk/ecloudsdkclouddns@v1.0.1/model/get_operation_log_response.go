// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type GetOperationLogResponseStateEnum string

// List of State
const (
    GetOperationLogResponseStateEnumError GetOperationLogResponseStateEnum = "ERROR"
    GetOperationLogResponseStateEnumException GetOperationLogResponseStateEnum = "EXCEPTION"
    GetOperationLogResponseStateEnumForbidden GetOperationLogResponseStateEnum = "FORBIDDEN"
    GetOperationLogResponseStateEnumOk GetOperationLogResponseStateEnum = "OK"
)

type GetOperationLogResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State GetOperationLogResponseStateEnum `json:"state,omitempty"`

	Body *GetOperationLogResponseBody `json:"body,omitempty"`
}
