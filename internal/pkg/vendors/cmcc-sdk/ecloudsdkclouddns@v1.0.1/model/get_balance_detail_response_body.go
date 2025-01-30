// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type GetBalanceDetailResponseBodyStateEnum string

// List of State
const (
    GetBalanceDetailResponseBodyStateEnumDisabled GetBalanceDetailResponseBodyStateEnum = "DISABLED"
    GetBalanceDetailResponseBodyStateEnumEnabled GetBalanceDetailResponseBodyStateEnum = "ENABLED"
)
type GetBalanceDetailResponseBodyTypeEnum string

// List of Type
const (
    GetBalanceDetailResponseBodyTypeEnumA GetBalanceDetailResponseBodyTypeEnum = "A"
    GetBalanceDetailResponseBodyTypeEnumAaaa GetBalanceDetailResponseBodyTypeEnum = "AAAA"
    GetBalanceDetailResponseBodyTypeEnumCaa GetBalanceDetailResponseBodyTypeEnum = "CAA"
    GetBalanceDetailResponseBodyTypeEnumCmauth GetBalanceDetailResponseBodyTypeEnum = "CMAUTH"
    GetBalanceDetailResponseBodyTypeEnumCname GetBalanceDetailResponseBodyTypeEnum = "CNAME"
    GetBalanceDetailResponseBodyTypeEnumMx GetBalanceDetailResponseBodyTypeEnum = "MX"
    GetBalanceDetailResponseBodyTypeEnumNs GetBalanceDetailResponseBodyTypeEnum = "NS"
    GetBalanceDetailResponseBodyTypeEnumPtr GetBalanceDetailResponseBodyTypeEnum = "PTR"
    GetBalanceDetailResponseBodyTypeEnumRp GetBalanceDetailResponseBodyTypeEnum = "RP"
    GetBalanceDetailResponseBodyTypeEnumSpf GetBalanceDetailResponseBodyTypeEnum = "SPF"
    GetBalanceDetailResponseBodyTypeEnumSrv GetBalanceDetailResponseBodyTypeEnum = "SRV"
    GetBalanceDetailResponseBodyTypeEnumTxt GetBalanceDetailResponseBodyTypeEnum = "TXT"
    GetBalanceDetailResponseBodyTypeEnumUrl GetBalanceDetailResponseBodyTypeEnum = "URL"
)

type GetBalanceDetailResponseBody struct {

	// 解析记录ID
	RecordId string `json:"recordId,omitempty"`

	// 主机头
	Rr string `json:"rr,omitempty"`

	// 线路中文名
	LineZh string `json:"lineZh,omitempty"`

	// 负载均衡权重
	Rate *int32 `json:"rate,omitempty"`

	// 解析记录状态
	State GetBalanceDetailResponseBodyStateEnum `json:"state,omitempty"`

	// 记录类型
	Type GetBalanceDetailResponseBodyTypeEnum `json:"type,omitempty"`

	// 记录值
	Value string `json:"value,omitempty"`
}
