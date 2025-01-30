// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model



type CreateOrderOpenapiResponseBody struct {

	// 订单ID
	OrderId string `json:"orderId"`

	// 订单项集合
	Products *[]CreateOrderOpenapiResponseProducts `json:"products,omitempty"`
}
