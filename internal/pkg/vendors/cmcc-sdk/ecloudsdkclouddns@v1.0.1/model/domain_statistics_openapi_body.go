// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)
type DomainStatisticsOpenapiBodyPeriodEnum string

// List of Period
const (
    DomainStatisticsOpenapiBodyPeriodEnumDate DomainStatisticsOpenapiBodyPeriodEnum = "date"
    DomainStatisticsOpenapiBodyPeriodEnumHour DomainStatisticsOpenapiBodyPeriodEnum = "hour"
)

type DomainStatisticsOpenapiBody struct {
    position.Body
	// 统计周期：hour 按小时统计，date 按日期统计，默认为date，如果end_date – start_date > 10，只能是date
	Period DomainStatisticsOpenapiBodyPeriodEnum `json:"period"`

	// 截止日期，如:2021-06-10
	EndDate string `json:"endDate"`

	// 域名
	DomainName string `json:"domainName"`

	// 开始日期，如:2021-06-01
	StartDate string `json:"startDate"`

	// 子域名
	SubdomainName string `json:"subdomainName,omitempty"`
}
