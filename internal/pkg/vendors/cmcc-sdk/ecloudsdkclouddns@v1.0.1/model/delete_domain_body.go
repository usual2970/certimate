// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)

type DeleteDomainBody struct {
    position.Body
	// 待删除的域名列表
	DomainNameList []string `json:"domainNameList"`
}
