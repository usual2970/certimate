// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type GetRecordResponseBodyTypeEnum string

// List of Type
const (
    GetRecordResponseBodyTypeEnumA GetRecordResponseBodyTypeEnum = "A"
    GetRecordResponseBodyTypeEnumAaaa GetRecordResponseBodyTypeEnum = "AAAA"
    GetRecordResponseBodyTypeEnumCaa GetRecordResponseBodyTypeEnum = "CAA"
    GetRecordResponseBodyTypeEnumCmauth GetRecordResponseBodyTypeEnum = "CMAUTH"
    GetRecordResponseBodyTypeEnumCname GetRecordResponseBodyTypeEnum = "CNAME"
    GetRecordResponseBodyTypeEnumMx GetRecordResponseBodyTypeEnum = "MX"
    GetRecordResponseBodyTypeEnumNs GetRecordResponseBodyTypeEnum = "NS"
    GetRecordResponseBodyTypeEnumPtr GetRecordResponseBodyTypeEnum = "PTR"
    GetRecordResponseBodyTypeEnumRp GetRecordResponseBodyTypeEnum = "RP"
    GetRecordResponseBodyTypeEnumSpf GetRecordResponseBodyTypeEnum = "SPF"
    GetRecordResponseBodyTypeEnumSrv GetRecordResponseBodyTypeEnum = "SRV"
    GetRecordResponseBodyTypeEnumTxt GetRecordResponseBodyTypeEnum = "TXT"
    GetRecordResponseBodyTypeEnumUrl GetRecordResponseBodyTypeEnum = "URL"
)
type GetRecordResponseBodyTimedStatusEnum string

// List of TimedStatus
const (
    GetRecordResponseBodyTimedStatusEnumDisabled GetRecordResponseBodyTimedStatusEnum = "DISABLED"
    GetRecordResponseBodyTimedStatusEnumEnabled GetRecordResponseBodyTimedStatusEnum = "ENABLED"
    GetRecordResponseBodyTimedStatusEnumTimed GetRecordResponseBodyTimedStatusEnum = "TIMED"
)
type GetRecordResponseBodyStateEnum string

// List of State
const (
    GetRecordResponseBodyStateEnumDisabled GetRecordResponseBodyStateEnum = "DISABLED"
    GetRecordResponseBodyStateEnumEnabled GetRecordResponseBodyStateEnum = "ENABLED"
)

type GetRecordResponseBody struct {

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
	Type GetRecordResponseBodyTypeEnum `json:"type,omitempty"`

	// 缓存的生命周期
	Ttl *int32 `json:"ttl,omitempty"`

	// 标签
	Tags *[]GetRecordResponseTags `json:"tags,omitempty"`

	// 解析记录ID
	RecordId string `json:"recordId,omitempty"`

	// 定时状态
	TimedStatus GetRecordResponseBodyTimedStatusEnum `json:"timedStatus,omitempty"`

	// 域名名称
	DomainName string `json:"domainName,omitempty"`

	// 线路英文名
	LineEn string `json:"lineEn,omitempty"`

	// 状态
	State GetRecordResponseBodyStateEnum `json:"state,omitempty"`

	// 记录值
	Value string `json:"value,omitempty"`

	// 定时发布时间
	Pubdate string `json:"pubdate,omitempty"`
}
