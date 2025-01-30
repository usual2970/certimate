// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)

type ModifyDomainBody struct {
    position.Body
	// 实例ID
	InstanceId string `json:"instanceId"`

	// 域名
	DomainName string `json:"domainName"`
}
