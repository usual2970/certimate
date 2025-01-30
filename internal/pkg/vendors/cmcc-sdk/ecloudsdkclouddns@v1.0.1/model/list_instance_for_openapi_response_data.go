// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type ListInstanceForOpenapiResponseDataPackageTypeEnum string

// List of PackageType
const (
    ListInstanceForOpenapiResponseDataPackageTypeEnumBasic ListInstanceForOpenapiResponseDataPackageTypeEnum = "BASIC"
    ListInstanceForOpenapiResponseDataPackageTypeEnumPremium ListInstanceForOpenapiResponseDataPackageTypeEnum = "PREMIUM"
    ListInstanceForOpenapiResponseDataPackageTypeEnumStandard ListInstanceForOpenapiResponseDataPackageTypeEnum = "STANDARD"
    ListInstanceForOpenapiResponseDataPackageTypeEnumUltimate ListInstanceForOpenapiResponseDataPackageTypeEnum = "ULTIMATE"
)
type ListInstanceForOpenapiResponseDataStatusEnum string

// List of Status
const (
    ListInstanceForOpenapiResponseDataStatusEnumCanceling ListInstanceForOpenapiResponseDataStatusEnum = "CANCELING"
    ListInstanceForOpenapiResponseDataStatusEnumDisabled ListInstanceForOpenapiResponseDataStatusEnum = "DISABLED"
    ListInstanceForOpenapiResponseDataStatusEnumEnabled ListInstanceForOpenapiResponseDataStatusEnum = "ENABLED"
    ListInstanceForOpenapiResponseDataStatusEnumExpired ListInstanceForOpenapiResponseDataStatusEnum = "EXPIRED"
)

type ListInstanceForOpenapiResponseData struct {

	// 订单ID
	OrderId string `json:"orderId,omitempty"`

	// 自动续订的时长
	AutoRenewDuration *int32 `json:"autoRenewDuration,omitempty"`

	// 套餐类型
	PackageType ListInstanceForOpenapiResponseDataPackageTypeEnum `json:"packageType,omitempty"`

	// 订购时长
	Duration *int32 `json:"duration,omitempty"`

	// 是否达到绑定次数上限
	Bindable *bool `json:"bindable,omitempty"`

	// 实例ID
	InstanceId string `json:"instanceId,omitempty"`

	// 产品购买时间
	OrderTime string `json:"orderTime,omitempty"`

	// 是否修改过到期时间
	ModifyEnd *bool `json:"modifyEnd,omitempty"`

	// 绑定的域名
	DomainName string `json:"domainName,omitempty"`

	// 是否自动续订
	AutoRenew *bool `json:"autoRenew,omitempty"`

	// 订购时长单位:year,month
	DurationUnit string `json:"durationUnit,omitempty"`

	// 产品到期时间
	EndTime string `json:"endTime,omitempty"`

	// 产品编号
	ServiceId string `json:"serviceId,omitempty"`

	// 状态
	Status ListInstanceForOpenapiResponseDataStatusEnum `json:"status,omitempty"`
}
