// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)
type ListInstanceBodyPackageTypeEnum string

// List of PackageType
const (
    ListInstanceBodyPackageTypeEnumBasic ListInstanceBodyPackageTypeEnum = "BASIC"
    ListInstanceBodyPackageTypeEnumFree ListInstanceBodyPackageTypeEnum = "FREE"
    ListInstanceBodyPackageTypeEnumPremium ListInstanceBodyPackageTypeEnum = "PREMIUM"
    ListInstanceBodyPackageTypeEnumStandard ListInstanceBodyPackageTypeEnum = "STANDARD"
    ListInstanceBodyPackageTypeEnumUltimate ListInstanceBodyPackageTypeEnum = "ULTIMATE"
)

type ListInstanceBody struct {
    position.Body
	// 实例ID
	InstanceId string `json:"instanceId,omitempty"`

	// 绑定的域名，支持模糊搜索
	DomainNameLike string `json:"domainNameLike,omitempty"`

	// 排序
	OrderBy *[]ListInstanceRequestOrderBy `json:"orderBy,omitempty"`

	// 套餐
	PackageType ListInstanceBodyPackageTypeEnum `json:"packageType,omitempty"`
}
