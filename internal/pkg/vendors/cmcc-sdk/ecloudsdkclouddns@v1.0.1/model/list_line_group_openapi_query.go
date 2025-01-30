// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)

type ListLineGroupOpenapiQuery struct {
    position.Query
	// 域名名称
	DomainName string `json:"domainName"`

	// 页大小
	PageSize *int32 `json:"pageSize,omitempty"`

	// 当前页
	Page *int32 `json:"page,omitempty"`
}
