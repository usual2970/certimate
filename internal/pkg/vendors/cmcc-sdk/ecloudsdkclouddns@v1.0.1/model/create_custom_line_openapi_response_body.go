// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model



type CreateCustomLineOpenapiResponseBody struct {

	// 占用自定义线路的主机头
	Rrs string `json:"rrs,omitempty"`

	// 线路名称
	LineZh string `json:"lineZh,omitempty"`

	// IP段数
	Count *int32 `json:"count,omitempty"`

	// 线路ID
	LineId string `json:"lineId,omitempty"`

	// 包含IP段
	Ips string `json:"ips,omitempty"`
}
