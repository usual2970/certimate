package onepanelsite

type ResourceType string

const (
	// 资源类型：替换指定网站的证书。
	RESOURCE_TYPE_WEBSITE = ResourceType("website")
	// 资源类型：替换指定证书。
	RESOURCE_TYPE_CERTIFICATE = ResourceType("certificate")
)
