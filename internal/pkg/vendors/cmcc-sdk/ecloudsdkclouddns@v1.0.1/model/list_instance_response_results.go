// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListInstanceResponseResultsPackageTypeEnum string

// List of PackageType
const (
    ListInstanceResponseResultsPackageTypeEnumBasic ListInstanceResponseResultsPackageTypeEnum = "BASIC"
    ListInstanceResponseResultsPackageTypeEnumFree ListInstanceResponseResultsPackageTypeEnum = "FREE"
    ListInstanceResponseResultsPackageTypeEnumPremium ListInstanceResponseResultsPackageTypeEnum = "PREMIUM"
    ListInstanceResponseResultsPackageTypeEnumStandard ListInstanceResponseResultsPackageTypeEnum = "STANDARD"
    ListInstanceResponseResultsPackageTypeEnumUltimate ListInstanceResponseResultsPackageTypeEnum = "ULTIMATE"
)
type ListInstanceResponseResultsStatusEnum string

// List of Status
const (
    ListInstanceResponseResultsStatusEnumCanceling ListInstanceResponseResultsStatusEnum = "CANCELING"
    ListInstanceResponseResultsStatusEnumDisabled ListInstanceResponseResultsStatusEnum = "DISABLED"
    ListInstanceResponseResultsStatusEnumEnabled ListInstanceResponseResultsStatusEnum = "ENABLED"
    ListInstanceResponseResultsStatusEnumExpired ListInstanceResponseResultsStatusEnum = "EXPIRED"
)

type ListInstanceResponseResults struct {

	// 订单项ID
	OrderExtId string `json:"orderExtId,omitempty"`

	// 自动续订的时长
	AutoRenewDuration *int32 `json:"autoRenewDuration,omitempty"`

	// 套餐类型
	PackageType ListInstanceResponseResultsPackageTypeEnum `json:"packageType,omitempty"`

	// 订购时长
	Duration *int32 `json:"duration,omitempty"`

	// 实例ID
	InstanceId string `json:"instanceId,omitempty"`

	// 产品购买时间
	OrderTime *int64 `json:"orderTime,omitempty"`

	// 是否修改过到期时间
	ModifyEnd *bool `json:"modifyEnd,omitempty"`

	// 绑定的域名
	DomainName string `json:"domainName,omitempty"`

	// 是否自动续订
	AutoRenew *bool `json:"autoRenew,omitempty"`

	// 订购时长单位:year,month
	DurationUnit string `json:"durationUnit,omitempty"`

	// 产品到期时间
	EndTime *int64 `json:"endTime,omitempty"`

	// 产品编号
	ServiceId string `json:"serviceId,omitempty"`

	// 状态
	Status ListInstanceResponseResultsStatusEnum `json:"status,omitempty"`
}
