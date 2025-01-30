// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListDomainForOpenapiRequestOrderByDirectionEnum string

// List of Direction
const (
    ListDomainForOpenapiRequestOrderByDirectionEnumAsc ListDomainForOpenapiRequestOrderByDirectionEnum = "ASC"
    ListDomainForOpenapiRequestOrderByDirectionEnumDesc ListDomainForOpenapiRequestOrderByDirectionEnum = "DESC"
)

type ListDomainForOpenapiRequestOrderBy struct {

	// 排序字段
	Field string `json:"field"`

	// 顺序
	Direction ListDomainForOpenapiRequestOrderByDirectionEnum `json:"direction"`
}
