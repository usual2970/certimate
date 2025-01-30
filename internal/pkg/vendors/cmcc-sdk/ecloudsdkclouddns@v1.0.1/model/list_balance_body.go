// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)

type ListBalanceBody struct {
    position.Body
	// 域名
	DomainName string `json:"domainName,omitempty"`

	// 排序(解析记录数)
	OrderBy *[]ListBalanceRequestOrderBy `json:"orderBy,omitempty"`

	// 状态
	State string `json:"state,omitempty"`
}
