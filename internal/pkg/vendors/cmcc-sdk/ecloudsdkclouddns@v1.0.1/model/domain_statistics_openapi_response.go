// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type DomainStatisticsOpenapiResponseStateEnum string

// List of State
const (
    DomainStatisticsOpenapiResponseStateEnumError DomainStatisticsOpenapiResponseStateEnum = "ERROR"
    DomainStatisticsOpenapiResponseStateEnumException DomainStatisticsOpenapiResponseStateEnum = "EXCEPTION"
    DomainStatisticsOpenapiResponseStateEnumForbidden DomainStatisticsOpenapiResponseStateEnum = "FORBIDDEN"
    DomainStatisticsOpenapiResponseStateEnumOk DomainStatisticsOpenapiResponseStateEnum = "OK"
)

type DomainStatisticsOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State DomainStatisticsOpenapiResponseStateEnum `json:"state,omitempty"`

	Body interface{} `json:"body,omitempty"`
}
