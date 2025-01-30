// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)

type ListInstanceForOpenapiQuery struct {
    position.Query
	// 页大小，请输入1-100
	PageSize *int32 `json:"pageSize,omitempty"`

	// 当前页
	Page *int32 `json:"page,omitempty"`
}
