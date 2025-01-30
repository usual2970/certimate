// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)
type ListNewRecordLineOpenapiQueryPackageTypeEnum string

// List of PackageType
const (
    ListNewRecordLineOpenapiQueryPackageTypeEnumBasic ListNewRecordLineOpenapiQueryPackageTypeEnum = "BASIC"
    ListNewRecordLineOpenapiQueryPackageTypeEnumFree ListNewRecordLineOpenapiQueryPackageTypeEnum = "FREE"
    ListNewRecordLineOpenapiQueryPackageTypeEnumPremium ListNewRecordLineOpenapiQueryPackageTypeEnum = "PREMIUM"
    ListNewRecordLineOpenapiQueryPackageTypeEnumStandard ListNewRecordLineOpenapiQueryPackageTypeEnum = "STANDARD"
    ListNewRecordLineOpenapiQueryPackageTypeEnumUltimate ListNewRecordLineOpenapiQueryPackageTypeEnum = "ULTIMATE"
)

type ListNewRecordLineOpenapiQuery struct {
    position.Query
	// 套餐类型
	PackageType ListNewRecordLineOpenapiQueryPackageTypeEnum `json:"packageType"`
}
