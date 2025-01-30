// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)
type GetBalanceDetailBodyTypeEnum string

// List of Type
const (
    GetBalanceDetailBodyTypeEnumA GetBalanceDetailBodyTypeEnum = "A"
    GetBalanceDetailBodyTypeEnumAaaa GetBalanceDetailBodyTypeEnum = "AAAA"
    GetBalanceDetailBodyTypeEnumCname GetBalanceDetailBodyTypeEnum = "CNAME"
)

type GetBalanceDetailBody struct {
    position.Body
	// 主机头
	Rr string `json:"rr"`

	// 线路名称
	LineZh string `json:"lineZh"`

	// 域名
	DomainName string `json:"domainName"`

	// 记录类型
	Type GetBalanceDetailBodyTypeEnum `json:"type"`
}
