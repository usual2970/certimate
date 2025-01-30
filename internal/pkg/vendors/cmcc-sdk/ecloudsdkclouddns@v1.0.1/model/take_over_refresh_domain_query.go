// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)

type TakeOverRefreshDomainQuery struct {
    position.Query
	// 域名名称
	DomainName string `json:"domainName"`
}
