// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListBalanceOpenapiResponseDataStateEnum string

// List of State
const (
    ListBalanceOpenapiResponseDataStateEnumDisabled ListBalanceOpenapiResponseDataStateEnum = "DISABLED"
    ListBalanceOpenapiResponseDataStateEnumEnabled ListBalanceOpenapiResponseDataStateEnum = "ENABLED"
)
type ListBalanceOpenapiResponseDataTypeEnum string

// List of Type
const (
    ListBalanceOpenapiResponseDataTypeEnumA ListBalanceOpenapiResponseDataTypeEnum = "A"
    ListBalanceOpenapiResponseDataTypeEnumAaaa ListBalanceOpenapiResponseDataTypeEnum = "AAAA"
    ListBalanceOpenapiResponseDataTypeEnumCname ListBalanceOpenapiResponseDataTypeEnum = "CNAME"
)

type ListBalanceOpenapiResponseData struct {

	// 主机头
	Rr string `json:"rr,omitempty"`

	// 线路中文名
	LineZh string `json:"lineZh,omitempty"`

	// 负载均衡记录数
	Count *int32 `json:"count,omitempty"`

	// 状态
	State ListBalanceOpenapiResponseDataStateEnum `json:"state,omitempty"`

	// 记录类型
	Type ListBalanceOpenapiResponseDataTypeEnum `json:"type,omitempty"`
}
