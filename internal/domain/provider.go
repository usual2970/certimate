package domain

type AccessProviderType string

/*
提供商类型常量值。

	注意：如果追加新的常量值，请保持以 ASCII 排序。
	NOTICE: If you add new constant, please keep ASCII order.
*/
const (
	AccessProviderTypeACMEHttpReq  = AccessProviderType("acmehttpreq")
	AccessProviderTypeAliyun       = AccessProviderType("aliyun")
	AccessProviderTypeAWS          = AccessProviderType("aws")
	AccessProviderTypeBaiduCloud   = AccessProviderType("baiducloud")
	AccessProviderTypeBytePlus     = AccessProviderType("byteplus")
	AccessProviderTypeCloudflare   = AccessProviderType("cloudflare")
	AccessProviderTypeDogeCloud    = AccessProviderType("dogecloud")
	AccessProviderTypeGoDaddy      = AccessProviderType("godaddy")
	AccessProviderTypeHuaweiCloud  = AccessProviderType("huaweicloud")
	AccessProviderTypeKubernetes   = AccessProviderType("k8s")
	AccessProviderTypeLocal        = AccessProviderType("local")
	AccessProviderTypeNameDotCom   = AccessProviderType("namedotcom")
	AccessProviderTypeNameSilo     = AccessProviderType("namesilo")
	AccessProviderTypePowerDNS     = AccessProviderType("powerdns")
	AccessProviderTypeQiniu        = AccessProviderType("qiniu")
	AccessProviderTypeSSH          = AccessProviderType("ssh")
	AccessProviderTypeTencentCloud = AccessProviderType("tencentcloud")
	AccessProviderTypeVolcEngine   = AccessProviderType("volcengine")
	AccessProviderTypeWebhook      = AccessProviderType("webhook")
)

type DeployProviderType string

/*
提供商部署目标常量值。
短横线前的部分始终等于提供商类型。

	注意：如果追加新的常量值，请保持以 ASCII 排序。
	NOTICE: If you add new constant, please keep ASCII order.
*/
const (
	DeployProviderTypeAliyunALB        = DeployProviderType("aliyun-alb")
	DeployProviderTypeAliyunCDN        = DeployProviderType("aliyun-cdn")
	DeployProviderTypeAliyunCLB        = DeployProviderType("aliyun-clb")
	DeployProviderTypeAliyunDCDN       = DeployProviderType("aliyun-dcdn")
	DeployProviderTypeAliyunNLB        = DeployProviderType("aliyun-nlb")
	DeployProviderTypeAliyunOSS        = DeployProviderType("aliyun-oss")
	DeployProviderTypeBaiduCloudCDN    = DeployProviderType("baiducloud-cdn")
	DeployProviderTypeBytePlusCDN      = DeployProviderType("byteplus-cdn")
	DeployProviderTypeDogeCloudCDN     = DeployProviderType("dogecloud-cdn")
	DeployProviderTypeHuaweiCloudCDN   = DeployProviderType("huaweicloud-cdn")
	DeployProviderTypeHuaweiCloudELB   = DeployProviderType("huaweicloud-elb")
	DeployProviderTypeK8sSecret        = DeployProviderType("k8s-secret")
	DeployProviderTypeLocal            = DeployProviderType("local")
	DeployProviderTypeQiniuCDN         = DeployProviderType("qiniu-cdn")
	DeployProviderTypeSSH              = DeployProviderType("ssh")
	DeployProviderTypeTencentCloudCDN  = DeployProviderType("tencentcloud-cdn")
	DeployProviderTypeTencentCloudCLB  = DeployProviderType("tencentcloud-clb")
	DeployProviderTypeTencentCloudCOS  = DeployProviderType("tencentcloud-cos")
	DeployProviderTypeTencentCloudECDN = DeployProviderType("tencentcloud-ecdn")
	DeployProviderTypeTencentCloudEO   = DeployProviderType("tencentcloud-eo")
	DeployProviderTypeVolcEngineCDN    = DeployProviderType("volcengine-cdn")
	DeployProviderTypeVolcEngineLive   = DeployProviderType("volcengine-live")
	DeployProviderTypeWebhook          = DeployProviderType("webhook")
)
