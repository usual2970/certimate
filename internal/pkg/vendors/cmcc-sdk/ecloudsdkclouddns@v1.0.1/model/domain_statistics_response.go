// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type DomainStatisticsResponseStateEnum string

// List of State
const (
    DomainStatisticsResponseStateEnumError DomainStatisticsResponseStateEnum = "ERROR"
    DomainStatisticsResponseStateEnumException DomainStatisticsResponseStateEnum = "EXCEPTION"
    DomainStatisticsResponseStateEnumForbidden DomainStatisticsResponseStateEnum = "FORBIDDEN"
    DomainStatisticsResponseStateEnumOk DomainStatisticsResponseStateEnum = "OK"
)

type DomainStatisticsResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State DomainStatisticsResponseStateEnum `json:"state,omitempty"`

	Body interface{} `json:"body,omitempty"`
}
