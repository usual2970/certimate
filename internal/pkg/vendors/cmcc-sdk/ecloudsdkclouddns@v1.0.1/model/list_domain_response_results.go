// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListDomainResponseResultsServiceTypeEnum string

// List of ServiceType
const (
    ListDomainResponseResultsServiceTypeEnumBasic ListDomainResponseResultsServiceTypeEnum = "BASIC"
    ListDomainResponseResultsServiceTypeEnumFree ListDomainResponseResultsServiceTypeEnum = "FREE"
    ListDomainResponseResultsServiceTypeEnumPremium ListDomainResponseResultsServiceTypeEnum = "PREMIUM"
    ListDomainResponseResultsServiceTypeEnumStandard ListDomainResponseResultsServiceTypeEnum = "STANDARD"
    ListDomainResponseResultsServiceTypeEnumUltimate ListDomainResponseResultsServiceTypeEnum = "ULTIMATE"
)
type ListDomainResponseResultsAdoptStateEnum string

// List of AdoptState
const (
    ListDomainResponseResultsAdoptStateEnumUncertain ListDomainResponseResultsAdoptStateEnum = "UNCERTAIN"
    ListDomainResponseResultsAdoptStateEnumUnused ListDomainResponseResultsAdoptStateEnum = "UNUSED"
    ListDomainResponseResultsAdoptStateEnumUsing ListDomainResponseResultsAdoptStateEnum = "USING"
)
type ListDomainResponseResultsInstanceStatusEnum string

// List of InstanceStatus
const (
    ListDomainResponseResultsInstanceStatusEnumCanceling ListDomainResponseResultsInstanceStatusEnum = "CANCELING"
    ListDomainResponseResultsInstanceStatusEnumDisabled ListDomainResponseResultsInstanceStatusEnum = "DISABLED"
    ListDomainResponseResultsInstanceStatusEnumEnabled ListDomainResponseResultsInstanceStatusEnum = "ENABLED"
    ListDomainResponseResultsInstanceStatusEnumExpired ListDomainResponseResultsInstanceStatusEnum = "EXPIRED"
)

type ListDomainResponseResults struct {

	// 套餐类型
	ServiceType ListDomainResponseResultsServiceTypeEnum `json:"serviceType,omitempty"`

	// 域名的最新修改时间
	ModifiedTime *int64 `json:"modifiedTime,omitempty"`

	// 是否是被授权的子域名
	Flag *bool `json:"flag,omitempty"`

	// 接管状态
	AdoptState ListDomainResponseResultsAdoptStateEnum `json:"adoptState,omitempty"`

	// 域名注册商
	Registrar string `json:"registrar,omitempty"`

	// 备注
	Description string `json:"description,omitempty"`

	// 域名到期时间
	DomainExpireDate string `json:"domainExpireDate,omitempty"`

	// 实例状态
	InstanceStatus ListDomainResponseResultsInstanceStatusEnum `json:"instanceStatus,omitempty"`

	// 域名的解析记录数
	RecordNum *int32 `json:"recordNum,omitempty"`

	// 域名TTL
	Ttl *int32 `json:"ttl,omitempty"`

	// 域名的解锁时间
	UnlockDate *int64 `json:"unlockDate,omitempty"`

	// 域名ID
	DomainId string `json:"domainId,omitempty"`

	// 实例ID
	InstanceId string `json:"instanceId,omitempty"`

	// 接管状态最后刷新时间
	AdoptLastRefresh *int64 `json:"adoptLastRefresh,omitempty"`

	// 域名在用NS服务器
	NsInUsing []string `json:"nsInUsing,omitempty"`

	// 到期时间
	ExpirationTime *int64 `json:"expirationTime,omitempty"`

	// 域名
	DomainName string `json:"domainName,omitempty"`

	// 域名状态
	State string `json:"state,omitempty"`

	// 域名之前的状态
	LastStatus string `json:"lastStatus,omitempty"`

	// 域名的锁定时间
	LockDate *int64 `json:"lockDate,omitempty"`
}
