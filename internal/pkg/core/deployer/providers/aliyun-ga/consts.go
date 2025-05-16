package aliyunga

type ResourceType string

const (
	// 资源类型：部署到指定全球加速器。
	RESOURCE_TYPE_ACCELERATOR = ResourceType("accelerator")
	// 资源类型：部署到指定监听器。
	RESOURCE_TYPE_LISTENER = ResourceType("listener")
)
