// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)

type RemoveCustomLineOpenapiBody struct {
    position.Body
	// 待删除的自定义线路ID列表
	LineIds []string `json:"lineIds"`
}
