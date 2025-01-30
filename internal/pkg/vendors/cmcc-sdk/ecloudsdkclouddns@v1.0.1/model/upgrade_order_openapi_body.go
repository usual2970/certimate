// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)
type UpgradeOrderOpenapiBodyPackageTypeEnum string

// List of PackageType
const (
    UpgradeOrderOpenapiBodyPackageTypeEnumBasic UpgradeOrderOpenapiBodyPackageTypeEnum = "BASIC"
    UpgradeOrderOpenapiBodyPackageTypeEnumPremium UpgradeOrderOpenapiBodyPackageTypeEnum = "PREMIUM"
    UpgradeOrderOpenapiBodyPackageTypeEnumStandard UpgradeOrderOpenapiBodyPackageTypeEnum = "STANDARD"
    UpgradeOrderOpenapiBodyPackageTypeEnumUltimate UpgradeOrderOpenapiBodyPackageTypeEnum = "ULTIMATE"
)

type UpgradeOrderOpenapiBody struct {
    position.Body
	// 资源ID
	InstanceId string `json:"instanceId"`

	// 产品套餐类型
	PackageType UpgradeOrderOpenapiBodyPackageTypeEnum `json:"packageType"`
}
