// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)
type ListInstanceForOpenapiBodyPackageTypeEnum string

// List of PackageType
const (
    ListInstanceForOpenapiBodyPackageTypeEnumBasic ListInstanceForOpenapiBodyPackageTypeEnum = "BASIC"
    ListInstanceForOpenapiBodyPackageTypeEnumPremium ListInstanceForOpenapiBodyPackageTypeEnum = "PREMIUM"
    ListInstanceForOpenapiBodyPackageTypeEnumStandard ListInstanceForOpenapiBodyPackageTypeEnum = "STANDARD"
    ListInstanceForOpenapiBodyPackageTypeEnumUltimate ListInstanceForOpenapiBodyPackageTypeEnum = "ULTIMATE"
)

type ListInstanceForOpenapiBody struct {
    position.Body
	// 实例ID
	InstanceId string `json:"instanceId,omitempty"`

	// 绑定的域名，支持模糊搜索
	DomainNameLike string `json:"domainNameLike,omitempty"`

	// 排序
	OrderBy *[]ListInstanceForOpenapiRequestOrderBy `json:"orderBy,omitempty"`

	// 套餐版本
	PackageType ListInstanceForOpenapiBodyPackageTypeEnum `json:"packageType,omitempty"`
}
