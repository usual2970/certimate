// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)
type ControlDomainOpenapiBodyOperateEnum string

// List of Operate
const (
    ControlDomainOpenapiBodyOperateEnumDisable ControlDomainOpenapiBodyOperateEnum = "DISABLE"
    ControlDomainOpenapiBodyOperateEnumEnable ControlDomainOpenapiBodyOperateEnum = "ENABLE"
)

type ControlDomainOpenapiBody struct {
    position.Body
	// 操作
	Operate ControlDomainOpenapiBodyOperateEnum `json:"operate"`

	// 需要启停的域名列表
	DomainNameList []string `json:"domainNameList"`
}
