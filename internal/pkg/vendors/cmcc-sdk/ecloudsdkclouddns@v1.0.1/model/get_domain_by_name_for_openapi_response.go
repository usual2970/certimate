// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type GetDomainByNameForOpenapiResponseStateEnum string

// List of State
const (
    GetDomainByNameForOpenapiResponseStateEnumError GetDomainByNameForOpenapiResponseStateEnum = "ERROR"
    GetDomainByNameForOpenapiResponseStateEnumException GetDomainByNameForOpenapiResponseStateEnum = "EXCEPTION"
    GetDomainByNameForOpenapiResponseStateEnumForbidden GetDomainByNameForOpenapiResponseStateEnum = "FORBIDDEN"
    GetDomainByNameForOpenapiResponseStateEnumOk GetDomainByNameForOpenapiResponseStateEnum = "OK"
)

type GetDomainByNameForOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State GetDomainByNameForOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *GetDomainByNameForOpenapiResponseBody `json:"body,omitempty"`
}
