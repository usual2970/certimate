// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)
type ListDomainForOpenapiBodyStateEnum string

// List of State
const (
    ListDomainForOpenapiBodyStateEnumCabLocked ListDomainForOpenapiBodyStateEnum = "CAB_LOCKED"
    ListDomainForOpenapiBodyStateEnumDisabled ListDomainForOpenapiBodyStateEnum = "DISABLED"
    ListDomainForOpenapiBodyStateEnumEnabled ListDomainForOpenapiBodyStateEnum = "ENABLED"
    ListDomainForOpenapiBodyStateEnumLocked ListDomainForOpenapiBodyStateEnum = "LOCKED"
)
type ListDomainForOpenapiBodyPackageTypeEnum string

// List of PackageType
const (
    ListDomainForOpenapiBodyPackageTypeEnumBasic ListDomainForOpenapiBodyPackageTypeEnum = "BASIC"
    ListDomainForOpenapiBodyPackageTypeEnumPremium ListDomainForOpenapiBodyPackageTypeEnum = "PREMIUM"
    ListDomainForOpenapiBodyPackageTypeEnumStandard ListDomainForOpenapiBodyPackageTypeEnum = "STANDARD"
    ListDomainForOpenapiBodyPackageTypeEnumUltimate ListDomainForOpenapiBodyPackageTypeEnum = "ULTIMATE"
)

type ListDomainForOpenapiBody struct {
    position.Body
	// 是否筛选出不带安全防护的域名
	Security *bool `json:"security,omitempty"`

	// 是否仅筛选出主域名
	Auth *bool `json:"auth,omitempty"`

	// 域名，支持模糊查询
	DomainNameLike string `json:"domainNameLike,omitempty"`

	// 排序(解析记录数)
	OrderBy *[]ListDomainForOpenapiRequestOrderBy `json:"orderBy,omitempty"`

	// 状态，ENABLED-正常，DISABLED-停用，不填则查询全部
	State ListDomainForOpenapiBodyStateEnum `json:"state,omitempty"`

	// 套餐类型
	PackageType ListDomainForOpenapiBodyPackageTypeEnum `json:"packageType,omitempty"`
}
