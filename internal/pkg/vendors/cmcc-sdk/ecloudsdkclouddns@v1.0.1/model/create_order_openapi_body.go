// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)
type CreateOrderOpenapiBodyPeriodTypeEnum string

// List of PeriodType
const (
    CreateOrderOpenapiBodyPeriodTypeEnumMonth CreateOrderOpenapiBodyPeriodTypeEnum = "month"
    CreateOrderOpenapiBodyPeriodTypeEnumYear CreateOrderOpenapiBodyPeriodTypeEnum = "year"
)
type CreateOrderOpenapiBodyPackageTypeEnum string

// List of PackageType
const (
    CreateOrderOpenapiBodyPackageTypeEnumBasic CreateOrderOpenapiBodyPackageTypeEnum = "BASIC"
    CreateOrderOpenapiBodyPackageTypeEnumPremium CreateOrderOpenapiBodyPackageTypeEnum = "PREMIUM"
    CreateOrderOpenapiBodyPackageTypeEnumStandard CreateOrderOpenapiBodyPackageTypeEnum = "STANDARD"
    CreateOrderOpenapiBodyPackageTypeEnumUltimate CreateOrderOpenapiBodyPackageTypeEnum = "ULTIMATE"
)

type CreateOrderOpenapiBody struct {
    position.Body
	// 订购周期数
	PeriodNum *int32 `json:"periodNum"`

	// 订购周期类型
	PeriodType CreateOrderOpenapiBodyPeriodTypeEnum `json:"periodType"`

	// 域名
	DomainName string `json:"domainName,omitempty"`

	// 是否开启自动续订
	AutoRenew *bool `json:"autoRenew,omitempty"`

	// 产品套餐类型
	PackageType CreateOrderOpenapiBodyPackageTypeEnum `json:"packageType"`
}
