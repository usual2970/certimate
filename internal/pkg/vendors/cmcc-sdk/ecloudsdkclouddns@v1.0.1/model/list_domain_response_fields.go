// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model



type ListDomainResponseFields struct {

	// 该字段是否可见，默认可见
	Visible *bool `json:"visible,omitempty"`

	// 字段名称
	Name string `json:"name,omitempty"`

	// 字段描述，该字段可为空
	Description string `json:"description,omitempty"`

	// 字段标签
	Label string `json:"label,omitempty"`

	// 该字段是否可排序，该字段可为空，空表示不排序
	Sortable *bool `json:"sortable,omitempty"`

	// 字段类型
	Type string `json:"type,omitempty"`
}
