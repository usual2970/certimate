// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model



type GetOperationLogResponseResults struct {

	// 域名
	DomainName string `json:"domainName,omitempty"`

	// ID
	Id string `json:"id,omitempty"`

	// 操作内容
	Operation string `json:"operation,omitempty"`

	// 操作人
	Operator string `json:"operator,omitempty"`

	// 操作时间
	OperationTime *int64 `json:"operationTime,omitempty"`
}
