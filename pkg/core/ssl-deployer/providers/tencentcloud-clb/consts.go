package tencentcloudclb

type ResourceType string

const (
	// 资源类型：通过 SSL 服务部署到云资源实例。
	RESOURCE_TYPE_VIA_SSLDEPLOY = ResourceType("ssl-deploy")
	// 资源类型：部署到指定负载均衡器。
	RESOURCE_TYPE_LOADBALANCER = ResourceType("loadbalancer")
	// 资源类型：部署到指定监听器。
	RESOURCE_TYPE_LISTENER = ResourceType("listener")
	// 资源类型：部署到指定转发规则域名。
	RESOURCE_TYPE_RULEDOMAIN = ResourceType("ruledomain")
)
