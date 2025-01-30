// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model



type UpgradeOrderOpenapiResponseBody struct {

	// 订单ID
	OrderId string `json:"orderId"`

	// 订单项集合
	Products *[]UpgradeOrderOpenapiResponseProducts `json:"products,omitempty"`
}
