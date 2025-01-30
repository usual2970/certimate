// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListBalanceResponseResultsStateEnum string

// List of State
const (
    ListBalanceResponseResultsStateEnumDisabled ListBalanceResponseResultsStateEnum = "DISABLED"
    ListBalanceResponseResultsStateEnumEnabled ListBalanceResponseResultsStateEnum = "ENABLED"
)
type ListBalanceResponseResultsTypeEnum string

// List of Type
const (
    ListBalanceResponseResultsTypeEnumA ListBalanceResponseResultsTypeEnum = "A"
    ListBalanceResponseResultsTypeEnumAaaa ListBalanceResponseResultsTypeEnum = "AAAA"
    ListBalanceResponseResultsTypeEnumCname ListBalanceResponseResultsTypeEnum = "CNAME"
)

type ListBalanceResponseResults struct {

	// 主机头
	Rr string `json:"rr,omitempty"`

	// 线路中文名
	LineZh string `json:"lineZh,omitempty"`

	// 负载均衡记录数
	Count *int32 `json:"count,omitempty"`

	// 状态
	State ListBalanceResponseResultsStateEnum `json:"state,omitempty"`

	// 记录类型
	Type ListBalanceResponseResultsTypeEnum `json:"type,omitempty"`
}
