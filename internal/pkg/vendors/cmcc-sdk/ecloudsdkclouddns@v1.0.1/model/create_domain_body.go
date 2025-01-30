// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)

type CreateDomainBody struct {
    position.Body
	// 域名创建实体
	DomainList *[]CreateDomainRequestDomainList `json:"domainList"`
}
