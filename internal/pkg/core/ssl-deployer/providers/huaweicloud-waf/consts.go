package huaweicloudwaf

type ResourceType string

const (
	// 资源类型：替换指定证书。
	RESOURCE_TYPE_CERTIFICATE = ResourceType("certificate")
	// 资源类型：部署到云模式防护网站。
	RESOURCE_TYPE_CLOUDSERVER = ResourceType("cloudserver")
	// 资源类型：部署到独享模式防护网站。
	RESOURCE_TYPE_PREMIUMHOST = ResourceType("premiumhost")
)
