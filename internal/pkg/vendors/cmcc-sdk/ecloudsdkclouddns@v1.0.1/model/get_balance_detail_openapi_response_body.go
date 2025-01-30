// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type GetBalanceDetailOpenapiResponseBodyStateEnum string

// List of State
const (
    GetBalanceDetailOpenapiResponseBodyStateEnumDisabled GetBalanceDetailOpenapiResponseBodyStateEnum = "DISABLED"
    GetBalanceDetailOpenapiResponseBodyStateEnumEnabled GetBalanceDetailOpenapiResponseBodyStateEnum = "ENABLED"
)
type GetBalanceDetailOpenapiResponseBodyTypeEnum string

// List of Type
const (
    GetBalanceDetailOpenapiResponseBodyTypeEnumA GetBalanceDetailOpenapiResponseBodyTypeEnum = "A"
    GetBalanceDetailOpenapiResponseBodyTypeEnumAaaa GetBalanceDetailOpenapiResponseBodyTypeEnum = "AAAA"
    GetBalanceDetailOpenapiResponseBodyTypeEnumCname GetBalanceDetailOpenapiResponseBodyTypeEnum = "CNAME"
)

type GetBalanceDetailOpenapiResponseBody struct {

	// 解析记录ID
	RecordId string `json:"recordId,omitempty"`

	// 主机头
	Rr string `json:"rr,omitempty"`

	// 线路中文名
	LineZh string `json:"lineZh,omitempty"`

	// 负载均衡权重
	Rate *int32 `json:"rate,omitempty"`

	// 解析记录状态
	State GetBalanceDetailOpenapiResponseBodyStateEnum `json:"state,omitempty"`

	// 记录类型
	Type GetBalanceDetailOpenapiResponseBodyTypeEnum `json:"type,omitempty"`

	// 记录值
	Value string `json:"value,omitempty"`
}
