// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)
type ListNewRecordLineQueryPackageTypeEnum string

// List of PackageType
const (
    ListNewRecordLineQueryPackageTypeEnumBasic ListNewRecordLineQueryPackageTypeEnum = "BASIC"
    ListNewRecordLineQueryPackageTypeEnumFree ListNewRecordLineQueryPackageTypeEnum = "FREE"
    ListNewRecordLineQueryPackageTypeEnumPremium ListNewRecordLineQueryPackageTypeEnum = "PREMIUM"
    ListNewRecordLineQueryPackageTypeEnumStandard ListNewRecordLineQueryPackageTypeEnum = "STANDARD"
    ListNewRecordLineQueryPackageTypeEnumUltimate ListNewRecordLineQueryPackageTypeEnum = "ULTIMATE"
)

type ListNewRecordLineQuery struct {
    position.Query
	// 套餐类型
	PackageType ListNewRecordLineQueryPackageTypeEnum `json:"packageType"`
}
