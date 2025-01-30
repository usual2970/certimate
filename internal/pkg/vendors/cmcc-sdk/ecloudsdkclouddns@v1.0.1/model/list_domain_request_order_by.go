// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListDomainRequestOrderByDirectionEnum string

// List of Direction
const (
    ListDomainRequestOrderByDirectionEnumAsc ListDomainRequestOrderByDirectionEnum = "ASC"
    ListDomainRequestOrderByDirectionEnumDesc ListDomainRequestOrderByDirectionEnum = "DESC"
)

type ListDomainRequestOrderBy struct {

	// 排序字段
	Field string `json:"field"`

	// 顺序
	Direction ListDomainRequestOrderByDirectionEnum `json:"direction"`
}
