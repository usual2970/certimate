// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)

type CreateCustomLineOpenapiBody struct {
    position.Body
	// 自定义线路名称
	LineZh string `json:"lineZh"`

	// 域名
	DomainName string `json:"domainName"`

	// IP段
	Ips string `json:"ips"`
}
