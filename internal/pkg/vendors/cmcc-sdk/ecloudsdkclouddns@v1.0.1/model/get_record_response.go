// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type GetRecordResponseStateEnum string

// List of State
const (
    GetRecordResponseStateEnumError GetRecordResponseStateEnum = "ERROR"
    GetRecordResponseStateEnumException GetRecordResponseStateEnum = "EXCEPTION"
    GetRecordResponseStateEnumForbidden GetRecordResponseStateEnum = "FORBIDDEN"
    GetRecordResponseStateEnumOk GetRecordResponseStateEnum = "OK"
)

type GetRecordResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State GetRecordResponseStateEnum `json:"state,omitempty"`

	Body *GetRecordResponseBody `json:"body,omitempty"`
}
