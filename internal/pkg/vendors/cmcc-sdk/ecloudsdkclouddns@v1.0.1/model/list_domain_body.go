// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)
type ListDomainBodyStateEnum string

// List of State
const (
    ListDomainBodyStateEnumCabLocked ListDomainBodyStateEnum = "CAB_LOCKED"
    ListDomainBodyStateEnumDisabled ListDomainBodyStateEnum = "DISABLED"
    ListDomainBodyStateEnumEnabled ListDomainBodyStateEnum = "ENABLED"
    ListDomainBodyStateEnumLocked ListDomainBodyStateEnum = "LOCKED"
)
type ListDomainBodyPackageTypeEnum string

// List of PackageType
const (
    ListDomainBodyPackageTypeEnumBasic ListDomainBodyPackageTypeEnum = "BASIC"
    ListDomainBodyPackageTypeEnumFree ListDomainBodyPackageTypeEnum = "FREE"
    ListDomainBodyPackageTypeEnumPremium ListDomainBodyPackageTypeEnum = "PREMIUM"
    ListDomainBodyPackageTypeEnumStandard ListDomainBodyPackageTypeEnum = "STANDARD"
    ListDomainBodyPackageTypeEnumUltimate ListDomainBodyPackageTypeEnum = "ULTIMATE"
)

type ListDomainBody struct {
    position.Body
	// 域名，支持模糊查询
	DomainNameLike string `json:"domainNameLike,omitempty"`

	// 排序(解析记录数)
	OrderBy *[]ListDomainRequestOrderBy `json:"orderBy,omitempty"`

	// 状态，ENABLED-正常，DISABLED-停用，不填则查询全部
	State ListDomainBodyStateEnum `json:"state,omitempty"`

	// 套餐类型
	PackageType ListDomainBodyPackageTypeEnum `json:"packageType,omitempty"`
}
