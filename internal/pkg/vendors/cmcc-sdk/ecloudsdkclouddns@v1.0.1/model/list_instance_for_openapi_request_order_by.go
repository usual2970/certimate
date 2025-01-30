// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListInstanceForOpenapiRequestOrderByDirectionEnum string

// List of Direction
const (
    ListInstanceForOpenapiRequestOrderByDirectionEnumAsc ListInstanceForOpenapiRequestOrderByDirectionEnum = "ASC"
    ListInstanceForOpenapiRequestOrderByDirectionEnumDesc ListInstanceForOpenapiRequestOrderByDirectionEnum = "DESC"
)

type ListInstanceForOpenapiRequestOrderBy struct {

	// 排序字段
	Field string `json:"field,omitempty"`

	// 顺序
	Direction ListInstanceForOpenapiRequestOrderByDirectionEnum `json:"direction,omitempty"`
}
