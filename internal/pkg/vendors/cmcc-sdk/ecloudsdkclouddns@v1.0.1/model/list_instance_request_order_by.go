// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListInstanceRequestOrderByDirectionEnum string

// List of Direction
const (
    ListInstanceRequestOrderByDirectionEnumAsc ListInstanceRequestOrderByDirectionEnum = "ASC"
    ListInstanceRequestOrderByDirectionEnumDesc ListInstanceRequestOrderByDirectionEnum = "DESC"
)

type ListInstanceRequestOrderBy struct {

	// 排序字段
	Field string `json:"field,omitempty"`

	// 顺序
	Direction ListInstanceRequestOrderByDirectionEnum `json:"direction,omitempty"`
}
