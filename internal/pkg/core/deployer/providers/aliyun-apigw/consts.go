package aliyunapigw

type ServiceType string

const (
	// 服务类型：原 API 网关。
	SERVICE_TYPE_TRADITIONAL = ServiceType("traditional")
	// 服务类型：云原生 API 网关。
	SERVICE_TYPE_CLOUDNATIVE = ServiceType("cloudnative")
)
