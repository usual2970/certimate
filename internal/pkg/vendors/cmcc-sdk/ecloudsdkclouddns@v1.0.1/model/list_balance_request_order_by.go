// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListBalanceRequestOrderByDirectionEnum string

// List of Direction
const (
    ListBalanceRequestOrderByDirectionEnumAsc ListBalanceRequestOrderByDirectionEnum = "ASC"
    ListBalanceRequestOrderByDirectionEnumDesc ListBalanceRequestOrderByDirectionEnum = "DESC"
)

type ListBalanceRequestOrderBy struct {

	// 排序字段
	Field string `json:"field,omitempty"`

	// 顺序
	Direction ListBalanceRequestOrderByDirectionEnum `json:"direction,omitempty"`
}
