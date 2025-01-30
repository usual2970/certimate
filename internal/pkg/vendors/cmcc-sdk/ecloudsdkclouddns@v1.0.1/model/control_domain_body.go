// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)
type ControlDomainBodyOperateEnum string

// List of Operate
const (
    ControlDomainBodyOperateEnumDisable ControlDomainBodyOperateEnum = "DISABLE"
    ControlDomainBodyOperateEnumEnable ControlDomainBodyOperateEnum = "ENABLE"
)

type ControlDomainBody struct {
    position.Body
	// 操作
	Operate ControlDomainBodyOperateEnum `json:"operate"`

	// 需要启停的域名列表
	DomainNameList []string `json:"domainNameList"`
}
