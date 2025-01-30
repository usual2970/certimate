// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)
type ModifyBalanceRateOpenapiBodyTypeEnum string

// List of Type
const (
    ModifyBalanceRateOpenapiBodyTypeEnumA ModifyBalanceRateOpenapiBodyTypeEnum = "A"
    ModifyBalanceRateOpenapiBodyTypeEnumAaaa ModifyBalanceRateOpenapiBodyTypeEnum = "AAAA"
    ModifyBalanceRateOpenapiBodyTypeEnumCname ModifyBalanceRateOpenapiBodyTypeEnum = "CNAME"
)

type ModifyBalanceRateOpenapiBody struct {
    position.Body
	// 主机头
	Rr string `json:"rr"`

	// 线路中文名
	LineZh string `json:"lineZh"`

	// 解析记录权重分配列表
	RateList *[]ModifyBalanceRateOpenapiRequestRateList `json:"rateList"`

	// 域名
	DomainName string `json:"domainName"`

	// 记录类型
	Type ModifyBalanceRateOpenapiBodyTypeEnum `json:"type"`
}
