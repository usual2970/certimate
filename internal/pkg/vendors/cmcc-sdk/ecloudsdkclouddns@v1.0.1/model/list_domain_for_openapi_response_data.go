// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListDomainForOpenapiResponseDataServiceTypeEnum string

// List of ServiceType
const (
    ListDomainForOpenapiResponseDataServiceTypeEnumBasic ListDomainForOpenapiResponseDataServiceTypeEnum = "BASIC"
    ListDomainForOpenapiResponseDataServiceTypeEnumPremium ListDomainForOpenapiResponseDataServiceTypeEnum = "PREMIUM"
    ListDomainForOpenapiResponseDataServiceTypeEnumStandard ListDomainForOpenapiResponseDataServiceTypeEnum = "STANDARD"
    ListDomainForOpenapiResponseDataServiceTypeEnumUltimate ListDomainForOpenapiResponseDataServiceTypeEnum = "ULTIMATE"
)
type ListDomainForOpenapiResponseDataAdoptStateEnum string

// List of AdoptState
const (
    ListDomainForOpenapiResponseDataAdoptStateEnumUncertain ListDomainForOpenapiResponseDataAdoptStateEnum = "UNCERTAIN"
    ListDomainForOpenapiResponseDataAdoptStateEnumUnused ListDomainForOpenapiResponseDataAdoptStateEnum = "UNUSED"
    ListDomainForOpenapiResponseDataAdoptStateEnumUsing ListDomainForOpenapiResponseDataAdoptStateEnum = "USING"
)
type ListDomainForOpenapiResponseDataInstanceStatusEnum string

// List of InstanceStatus
const (
    ListDomainForOpenapiResponseDataInstanceStatusEnumCanceling ListDomainForOpenapiResponseDataInstanceStatusEnum = "CANCELING"
    ListDomainForOpenapiResponseDataInstanceStatusEnumDisabled ListDomainForOpenapiResponseDataInstanceStatusEnum = "DISABLED"
    ListDomainForOpenapiResponseDataInstanceStatusEnumEnabled ListDomainForOpenapiResponseDataInstanceStatusEnum = "ENABLED"
    ListDomainForOpenapiResponseDataInstanceStatusEnumExpired ListDomainForOpenapiResponseDataInstanceStatusEnum = "EXPIRED"
)

type ListDomainForOpenapiResponseData struct {

	// 套餐类型
	ServiceType ListDomainForOpenapiResponseDataServiceTypeEnum `json:"serviceType,omitempty"`

	// 域名的最新修改时间
	ModifiedTime string `json:"modifiedTime,omitempty"`

	// 是否是被授权的子域名
	Flag *bool `json:"flag,omitempty"`

	// 接管状态
	AdoptState ListDomainForOpenapiResponseDataAdoptStateEnum `json:"adoptState,omitempty"`

	// 域名注册商
	Registrar string `json:"registrar,omitempty"`

	// 备注
	Description string `json:"description,omitempty"`

	// 域名到期时间
	DomainExpireDate string `json:"domainExpireDate,omitempty"`

	// 实例状态
	InstanceStatus ListDomainForOpenapiResponseDataInstanceStatusEnum `json:"instanceStatus,omitempty"`

	// 域名的解析记录数
	RecordNum *int32 `json:"recordNum,omitempty"`

	// 域名TTL
	Ttl *int32 `json:"ttl,omitempty"`

	// 域名的解锁时间
	UnlockDate string `json:"unlockDate,omitempty"`

	// 域名ID
	DomainId string `json:"domainId,omitempty"`

	// 实例ID
	InstanceId string `json:"instanceId,omitempty"`

	// 接管状态最后刷新时间
	AdoptLastRefresh string `json:"adoptLastRefresh,omitempty"`

	// 域名在用NS服务器
	NsInUsing []string `json:"nsInUsing,omitempty"`

	// 到期时间
	ExpirationTime string `json:"expirationTime,omitempty"`

	// 域名
	DomainName string `json:"domainName,omitempty"`

	// 是否在dnspod上
	InDnspod *bool `json:"inDnspod,omitempty"`

	// 域名状态
	State string `json:"state,omitempty"`

	// 域名之前的状态
	LastStatus string `json:"lastStatus,omitempty"`

	// 域名的锁定时间
	LockDate string `json:"lockDate,omitempty"`
}
