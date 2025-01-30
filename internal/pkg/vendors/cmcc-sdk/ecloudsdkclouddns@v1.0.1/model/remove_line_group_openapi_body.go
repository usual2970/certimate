// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)

type RemoveLineGroupOpenapiBody struct {
    position.Body
	// 待删除的线路分组ID列表
	GroupIds []string `json:"groupIds"`
}
