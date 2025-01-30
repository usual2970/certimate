// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model



type RenewProductOpenapiResponseProducts struct {

	// 资源ID，订购时不返回
	InstanceId string `json:"instanceId,omitempty"`

	// 订单项ID
	OrderExtId string `json:"orderExtId"`

	// 订单项序号
	SequenceId *int32 `json:"sequenceId"`
}
