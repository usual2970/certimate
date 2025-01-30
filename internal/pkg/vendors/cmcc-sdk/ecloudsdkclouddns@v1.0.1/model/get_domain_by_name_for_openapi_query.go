// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)

type GetDomainByNameForOpenapiQuery struct {
    position.Query
	// 域名
	DomainName string `json:"domainName"`
}
