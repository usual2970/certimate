// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model



type AddOrderResponseBody struct {

	// 订单ID
	OrderId string `json:"orderId,omitempty"`

	// 订单项集合
	Products *[]AddOrderResponseProducts `json:"products,omitempty"`
}
