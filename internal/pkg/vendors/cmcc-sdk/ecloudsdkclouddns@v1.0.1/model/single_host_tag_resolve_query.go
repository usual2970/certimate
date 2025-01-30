// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

type SingleHostTagResolveQuery struct {
	// 指定解析结果IP的类型，可以选择6(IPv6)或4(IPv4)。默认值为4
	Query string `json:"query,omitempty"`

	// 解析域名
	Host string `json:"host"`

	// 标签名称
	Tag string `json:"tag,omitempty"`
}
