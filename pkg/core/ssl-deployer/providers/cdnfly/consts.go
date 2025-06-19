package cdnfly

type ResourceType string

const (
	// 资源类型：替换指定网站的证书。
	RESOURCE_TYPE_SITE = ResourceType("site")
	// 资源类型：替换指定证书。
	RESOURCE_TYPE_CERTIFICATE = ResourceType("certificate")
)
