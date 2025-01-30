// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model



type ListLineGroupOpenapiResponseData struct {

	// 占用线路分组的主机头
	Rrs string `json:"rrs,omitempty"`

	// 域名名称
	DomainName string `json:"domainName,omitempty"`

	// 线路分组ID
	GroupId string `json:"groupId,omitempty"`

	// 线路数
	Count *int32 `json:"count,omitempty"`

	// 线路分组名称
	Name string `json:"name,omitempty"`

	// 线路名称集合
	LineZhs string `json:"lineZhs,omitempty"`

	// 线路ID集合
	LineIds string `json:"lineIds,omitempty"`
}
