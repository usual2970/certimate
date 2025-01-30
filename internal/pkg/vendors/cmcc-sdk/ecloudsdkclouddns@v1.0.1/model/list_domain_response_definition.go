// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model



type ListDomainResponseDefinition struct {

	// 静态字段
	StaticFields *[]ListDomainResponseStaticFields `json:"staticFields,omitempty"`

	// 是否显示表头
	ShowHeader *bool `json:"showHeader,omitempty"`

	// 表格的名称
	Name string `json:"name,omitempty"`

	// 表格所拥有的所有的字段
	Fields *[]ListDomainResponseFields `json:"fields,omitempty"`

	// 动作列表
	Actions *[]ListDomainResponseActions `json:"actions,omitempty"`

	// 表格属性
	Properties interface{} `json:"properties,omitempty"`
}
