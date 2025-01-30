// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model



type RenewDnsOrderResponseBody struct {

	// 订单ID
	OrderId string `json:"orderId,omitempty"`

	// 订单项集合
	Products *[]RenewDnsOrderResponseProducts `json:"products,omitempty"`
}
