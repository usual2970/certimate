// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)

type CreateLineGroupBody struct {
    position.Body
	// 域名
	DomainName string `json:"domainName"`

	// 线路分组名称
	Name string `json:"name"`

	// 线路分组中的线路ID集合
	LineIds string `json:"lineIds"`
}
