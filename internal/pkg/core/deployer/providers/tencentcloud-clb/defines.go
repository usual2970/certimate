package tencentcloudclb

type DeployResourceType string

const (
	// 资源类型：通过 SSL 服务部署到云资源实例。
	DEPLOY_RESOURCE_USE_SSLDEPLOY = DeployResourceType("ssl-deploy")
	// 资源类型：部署到指定负载均衡器。
	DEPLOY_RESOURCE_LOADBALANCER = DeployResourceType("loadbalancer")
	// 资源类型：部署到指定监听器。
	DEPLOY_RESOURCE_LISTENER = DeployResourceType("listener")
	// 资源类型：部署到指定转发规则域名。
	DEPLOY_RESOURCE_RULEDOMAIN = DeployResourceType("ruledomain")
)
