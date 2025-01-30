// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

type BatchHostResolveQuery struct {
	// 指定解析结果IP的类型，可以选择6(IPv6)或4(IPv4)。默认值为4
	Query string `json:"query,omitempty"`

	// 解析域名，多个域名之间以逗号,分隔，单次请求最多允许携带5个域名
	Host string `json:"host"`
}
