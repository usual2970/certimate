package baiducloudappblb

type ResourceType string

const (
	// 资源类型：部署到指定负载均衡器。
	RESOURCE_TYPE_LOADBALANCER = ResourceType("loadbalancer")
	// 资源类型：部署到指定监听器。
	RESOURCE_TYPE_LISTENER = ResourceType("listener")
)
