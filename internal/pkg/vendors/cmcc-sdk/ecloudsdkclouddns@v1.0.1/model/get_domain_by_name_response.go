// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type GetDomainByNameResponseStateEnum string

// List of State
const (
    GetDomainByNameResponseStateEnumError GetDomainByNameResponseStateEnum = "ERROR"
    GetDomainByNameResponseStateEnumException GetDomainByNameResponseStateEnum = "EXCEPTION"
    GetDomainByNameResponseStateEnumForbidden GetDomainByNameResponseStateEnum = "FORBIDDEN"
    GetDomainByNameResponseStateEnumOk GetDomainByNameResponseStateEnum = "OK"
)

type GetDomainByNameResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State GetDomainByNameResponseStateEnum `json:"state,omitempty"`

	Body *GetDomainByNameResponseBody `json:"body,omitempty"`
}
