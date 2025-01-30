// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)
type AddOrderBodyPeriodTypeEnum string

// List of PeriodType
const (
    AddOrderBodyPeriodTypeEnumMonth AddOrderBodyPeriodTypeEnum = "month"
    AddOrderBodyPeriodTypeEnumYear AddOrderBodyPeriodTypeEnum = "year"
)
type AddOrderBodyPackageTypeEnum string

// List of PackageType
const (
    AddOrderBodyPackageTypeEnumBasic AddOrderBodyPackageTypeEnum = "BASIC"
    AddOrderBodyPackageTypeEnumPremium AddOrderBodyPackageTypeEnum = "PREMIUM"
    AddOrderBodyPackageTypeEnumStandard AddOrderBodyPackageTypeEnum = "STANDARD"
    AddOrderBodyPackageTypeEnumUltimate AddOrderBodyPackageTypeEnum = "ULTIMATE"
)

type AddOrderBody struct {
    position.Body
	// 订购周期数
	PeriodNum *int32 `json:"periodNum"`

	// 订购周期类型
	PeriodType AddOrderBodyPeriodTypeEnum `json:"periodType"`

	// 域名
	DomainName string `json:"domainName,omitempty"`

	// 是否开启自动续订
	AutoRenew *bool `json:"autoRenew,omitempty"`

	// 产品套餐类型
	PackageType AddOrderBodyPackageTypeEnum `json:"packageType"`
}
