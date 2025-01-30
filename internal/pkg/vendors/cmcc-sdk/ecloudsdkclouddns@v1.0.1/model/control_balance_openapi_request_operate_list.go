// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ControlBalanceOpenapiRequestOperateListOperateEnum string

// List of Operate
const (
    ControlBalanceOpenapiRequestOperateListOperateEnumDisable ControlBalanceOpenapiRequestOperateListOperateEnum = "DISABLE"
    ControlBalanceOpenapiRequestOperateListOperateEnumEnable ControlBalanceOpenapiRequestOperateListOperateEnum = "ENABLE"
)
type ControlBalanceOpenapiRequestOperateListTypeEnum string

// List of Type
const (
    ControlBalanceOpenapiRequestOperateListTypeEnumA ControlBalanceOpenapiRequestOperateListTypeEnum = "A"
    ControlBalanceOpenapiRequestOperateListTypeEnumAaaa ControlBalanceOpenapiRequestOperateListTypeEnum = "AAAA"
    ControlBalanceOpenapiRequestOperateListTypeEnumCname ControlBalanceOpenapiRequestOperateListTypeEnum = "CNAME"
)

type ControlBalanceOpenapiRequestOperateList struct {

	// 主机头
	Rr string `json:"rr"`

	// 操作
	Operate ControlBalanceOpenapiRequestOperateListOperateEnum `json:"operate"`

	// 线路中文名
	LineZh string `json:"lineZh"`

	// 记录类型
	Type ControlBalanceOpenapiRequestOperateListTypeEnum `json:"type"`
}
