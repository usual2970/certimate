// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model



type RenewProductOpenapiResponseBody struct {

	// 订单ID
	OrderId string `json:"orderId"`

	// 订单项集合
	Products *[]RenewProductOpenapiResponseProducts `json:"products,omitempty"`
}
