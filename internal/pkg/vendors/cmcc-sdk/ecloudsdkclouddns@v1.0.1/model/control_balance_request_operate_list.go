// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ControlBalanceRequestOperateListOperateEnum string

// List of Operate
const (
    ControlBalanceRequestOperateListOperateEnumDisable ControlBalanceRequestOperateListOperateEnum = "DISABLE"
    ControlBalanceRequestOperateListOperateEnumEnable ControlBalanceRequestOperateListOperateEnum = "ENABLE"
)
type ControlBalanceRequestOperateListTypeEnum string

// List of Type
const (
    ControlBalanceRequestOperateListTypeEnumA ControlBalanceRequestOperateListTypeEnum = "A"
    ControlBalanceRequestOperateListTypeEnumAaaa ControlBalanceRequestOperateListTypeEnum = "AAAA"
    ControlBalanceRequestOperateListTypeEnumCaa ControlBalanceRequestOperateListTypeEnum = "CAA"
    ControlBalanceRequestOperateListTypeEnumCmauth ControlBalanceRequestOperateListTypeEnum = "CMAUTH"
    ControlBalanceRequestOperateListTypeEnumCname ControlBalanceRequestOperateListTypeEnum = "CNAME"
    ControlBalanceRequestOperateListTypeEnumMx ControlBalanceRequestOperateListTypeEnum = "MX"
    ControlBalanceRequestOperateListTypeEnumNs ControlBalanceRequestOperateListTypeEnum = "NS"
    ControlBalanceRequestOperateListTypeEnumPtr ControlBalanceRequestOperateListTypeEnum = "PTR"
    ControlBalanceRequestOperateListTypeEnumRp ControlBalanceRequestOperateListTypeEnum = "RP"
    ControlBalanceRequestOperateListTypeEnumSpf ControlBalanceRequestOperateListTypeEnum = "SPF"
    ControlBalanceRequestOperateListTypeEnumSrv ControlBalanceRequestOperateListTypeEnum = "SRV"
    ControlBalanceRequestOperateListTypeEnumTxt ControlBalanceRequestOperateListTypeEnum = "TXT"
    ControlBalanceRequestOperateListTypeEnumUrl ControlBalanceRequestOperateListTypeEnum = "URL"
)

type ControlBalanceRequestOperateList struct {

	// 主机头
	Rr string `json:"rr"`

	// 操作
	Operate ControlBalanceRequestOperateListOperateEnum `json:"operate"`

	// 线路中文名
	LineZh string `json:"lineZh"`

	// 记录类型
	Type ControlBalanceRequestOperateListTypeEnum `json:"type"`
}
