// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)

type CreateDomainOpenapiBody struct {
    position.Body
	// 域名创建实体
	DomainList *[]CreateDomainOpenapiRequestDomainList `json:"domainList"`
}
