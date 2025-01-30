// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)
type ModifyOrderBodyPackageTypeEnum string

// List of PackageType
const (
    ModifyOrderBodyPackageTypeEnumPremium ModifyOrderBodyPackageTypeEnum = "PREMIUM"
    ModifyOrderBodyPackageTypeEnumStandard ModifyOrderBodyPackageTypeEnum = "STANDARD"
    ModifyOrderBodyPackageTypeEnumUltimate ModifyOrderBodyPackageTypeEnum = "ULTIMATE"
)

type ModifyOrderBody struct {
    position.Body
	// 资源ID
	InstanceId string `json:"instanceId"`

	// 产品套餐类型
	PackageType ModifyOrderBodyPackageTypeEnum `json:"packageType"`
}
