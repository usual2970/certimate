package aliyunnlb

type DeployResourceType string

const (
	// 资源类型：部署到指定负载均衡器。
	DEPLOY_RESOURCE_LOADBALANCER = DeployResourceType("loadbalancer")
	// 资源类型：部署到指定监听器。
	DEPLOY_RESOURCE_LISTENER = DeployResourceType("listener")
)
