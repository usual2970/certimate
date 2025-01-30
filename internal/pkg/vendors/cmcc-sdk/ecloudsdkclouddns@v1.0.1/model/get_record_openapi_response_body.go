// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type GetRecordOpenapiResponseBodyTypeEnum string

// List of Type
const (
    GetRecordOpenapiResponseBodyTypeEnumA GetRecordOpenapiResponseBodyTypeEnum = "A"
    GetRecordOpenapiResponseBodyTypeEnumAaaa GetRecordOpenapiResponseBodyTypeEnum = "AAAA"
    GetRecordOpenapiResponseBodyTypeEnumCname GetRecordOpenapiResponseBodyTypeEnum = "CNAME"
    GetRecordOpenapiResponseBodyTypeEnumMx GetRecordOpenapiResponseBodyTypeEnum = "MX"
    GetRecordOpenapiResponseBodyTypeEnumTxt GetRecordOpenapiResponseBodyTypeEnum = "TXT"
    GetRecordOpenapiResponseBodyTypeEnumNs GetRecordOpenapiResponseBodyTypeEnum = "NS"
    GetRecordOpenapiResponseBodyTypeEnumSpf GetRecordOpenapiResponseBodyTypeEnum = "SPF"
    GetRecordOpenapiResponseBodyTypeEnumSrv GetRecordOpenapiResponseBodyTypeEnum = "SRV"
    GetRecordOpenapiResponseBodyTypeEnumCaa GetRecordOpenapiResponseBodyTypeEnum = "CAA"
    GetRecordOpenapiResponseBodyTypeEnumCmauth GetRecordOpenapiResponseBodyTypeEnum = "CMAUTH"
)
type GetRecordOpenapiResponseBodyTimedStatusEnum string

// List of TimedStatus
const (
    GetRecordOpenapiResponseBodyTimedStatusEnumDisabled GetRecordOpenapiResponseBodyTimedStatusEnum = "DISABLED"
    GetRecordOpenapiResponseBodyTimedStatusEnumEnabled GetRecordOpenapiResponseBodyTimedStatusEnum = "ENABLED"
    GetRecordOpenapiResponseBodyTimedStatusEnumTimed GetRecordOpenapiResponseBodyTimedStatusEnum = "TIMED"
)
type GetRecordOpenapiResponseBodyStateEnum string

// List of State
const (
    GetRecordOpenapiResponseBodyStateEnumDisabled GetRecordOpenapiResponseBodyStateEnum = "DISABLED"
    GetRecordOpenapiResponseBodyStateEnumEnabled GetRecordOpenapiResponseBodyStateEnum = "ENABLED"
)

type GetRecordOpenapiResponseBody struct {

	// 主机头
	Rr string `json:"rr,omitempty"`

	// 修改时间
	ModifiedTime string `json:"modifiedTime,omitempty"`

	// 线路中文名
	LineZh string `json:"lineZh,omitempty"`

	// 备注
	Description string `json:"description,omitempty"`

	// 线路ID
	LineId string `json:"lineId,omitempty"`

	// 权重值
	Weight *int32 `json:"weight,omitempty"`

	// MX优先级
	MxPri *int32 `json:"mxPri,omitempty"`

	// 记录类型
	Type GetRecordOpenapiResponseBodyTypeEnum `json:"type,omitempty"`

	// 缓存的生命周期
	Ttl *int32 `json:"ttl,omitempty"`

	// 标签
	Tags *[]GetRecordOpenapiResponseTags `json:"tags,omitempty"`

	// 解析记录ID
	RecordId string `json:"recordId,omitempty"`

	// 定时状态
	TimedStatus GetRecordOpenapiResponseBodyTimedStatusEnum `json:"timedStatus,omitempty"`

	// 域名名称
	DomainName string `json:"domainName,omitempty"`

	// 线路英文名
	LineEn string `json:"lineEn,omitempty"`

	// 状态
	State GetRecordOpenapiResponseBodyStateEnum `json:"state,omitempty"`

	// 记录值
	Value string `json:"value,omitempty"`

	// 定时发布时间
	Pubdate string `json:"pubdate,omitempty"`
}
