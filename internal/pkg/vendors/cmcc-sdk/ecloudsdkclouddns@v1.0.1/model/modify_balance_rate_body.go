// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)
type ModifyBalanceRateBodyTypeEnum string

// List of Type
const (
    ModifyBalanceRateBodyTypeEnumA ModifyBalanceRateBodyTypeEnum = "A"
    ModifyBalanceRateBodyTypeEnumAaaa ModifyBalanceRateBodyTypeEnum = "AAAA"
    ModifyBalanceRateBodyTypeEnumCname ModifyBalanceRateBodyTypeEnum = "CNAME"
)

type ModifyBalanceRateBody struct {
    position.Body
	// 主机头
	Rr string `json:"rr,omitempty"`

	// 线路中文名
	LineZh string `json:"lineZh,omitempty"`

	// 解析记录权重分配列表
	RateList *[]ModifyBalanceRateRequestRateList `json:"rateList,omitempty"`

	// 域名
	DomainName string `json:"domainName,omitempty"`

	// 记录类型
	Type ModifyBalanceRateBodyTypeEnum `json:"type,omitempty"`
}
