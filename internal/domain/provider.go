package domain

type AccessProviderType string

/*
授权提供商类型常量值。

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

type ApplyDNSProviderType string

/*
申请证书 DNS 提供商常量值。
短横线前的部分始终等于授权提供商类型。

	注意：如果追加新的常量值，请保持以 ASCII 排序。
	NOTICE: If you add new constant, please keep ASCII order.
*/
const (
	ApplyDNSProviderTypeACMEHttpReq     = ApplyDNSProviderType("acmehttpreq")
	ApplyDNSProviderTypeAliyun          = ApplyDNSProviderType("aliyun") // 兼容旧值，等同于 [ApplyDNSProviderTypeAliyunDNS]
	ApplyDNSProviderTypeAliyunDNS       = ApplyDNSProviderType("aliyun-dns")
	ApplyDNSProviderTypeAWS             = ApplyDNSProviderType("aws") // 兼容旧值，等同于 [ApplyDNSProviderTypeAWSRoute53]
	ApplyDNSProviderTypeAWSRoute53      = ApplyDNSProviderType("aws-route53")
	ApplyDNSProviderTypeCloudflare      = ApplyDNSProviderType("cloudflare")
	ApplyDNSProviderTypeGoDaddy         = ApplyDNSProviderType("godaddy")
	ApplyDNSProviderTypeHuaweiCloud     = ApplyDNSProviderType("huaweicloud") // 兼容旧值，等同于 [ApplyDNSProviderTypeHuaweiCloudDNS]
	ApplyDNSProviderTypeHuaweiCloudDNS  = ApplyDNSProviderType("huaweicloud-dns")
	ApplyDNSProviderTypeNameDotCom      = ApplyDNSProviderType("namedotcom")
	ApplyDNSProviderTypeNameSilo        = ApplyDNSProviderType("namesilo")
	ApplyDNSProviderTypePowerDNS        = ApplyDNSProviderType("powerdns")
	ApplyDNSProviderTypeTencentCloud    = ApplyDNSProviderType("tencentcloud") // 兼容旧值，等同于 [ApplyDNSProviderTypeTencentCloudDNS]
	ApplyDNSProviderTypeTencentCloudDNS = ApplyDNSProviderType("tencentcloud-dns")
	ApplyDNSProviderTypeVolcEngine      = ApplyDNSProviderType("volcengine") // 兼容旧值，等同于 [ApplyDNSProviderTypeVolcEngineDNS]
	ApplyDNSProviderTypeVolcEngineDNS   = ApplyDNSProviderType("volcengine-dns")
)

type DeployProviderType string

/*
部署目标提供商常量值。
短横线前的部分始终等于授权提供商类型。

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
	DeployProviderTypeKubernetesSecret = DeployProviderType("k8s-secret")
	DeployProviderTypeLocal            = DeployProviderType("local")
	DeployProviderTypeQiniuCDN         = DeployProviderType("qiniu-cdn")
	DeployProviderTypeSSH              = DeployProviderType("ssh")
	DeployProviderTypeTencentCloudCDN  = DeployProviderType("tencentcloud-cdn")
	DeployProviderTypeTencentCloudCLB  = DeployProviderType("tencentcloud-clb")
	DeployProviderTypeTencentCloudCOS  = DeployProviderType("tencentcloud-cos")
	DeployProviderTypeTencentCloudECDN = DeployProviderType("tencentcloud-ecdn")
	DeployProviderTypeTencentCloudEO   = DeployProviderType("tencentcloud-eo")
	DeployProviderTypeVolcEngineCDN    = DeployProviderType("volcengine-cdn")
	DeployProviderTypeVolcEngineDCDN   = DeployProviderType("volcengine-dcdn")
	DeployProviderTypeVolcEngineLive   = DeployProviderType("volcengine-live")
	DeployProviderTypeVolcEngineTOS    = DeployProviderType("volcengine-tos")
	DeployProviderTypeWebhook          = DeployProviderType("webhook")
)
