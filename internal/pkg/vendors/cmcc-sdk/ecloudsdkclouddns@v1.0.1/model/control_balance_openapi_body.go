// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)

type ControlBalanceOpenapiBody struct {
    position.Body
	// 操作列表
	OperateList *[]ControlBalanceOpenapiRequestOperateList `json:"operateList"`

	// 域名
	DomainName string `json:"domainName"`
}
