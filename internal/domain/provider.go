package domain

type AccessProviderType string

/*
授权提供商类型常量值。

	注意：如果追加新的常量值，请保持以 ASCII 排序。
	NOTICE: If you add new constant, please keep ASCII order.
*/
const (
	AccessProviderType1Panel       = AccessProviderType("1panel")
	AccessProviderTypeACMEHttpReq  = AccessProviderType("acmehttpreq")
	AccessProviderTypeAkamai       = AccessProviderType("akamai") // Akamai（预留）
	AccessProviderTypeAliyun       = AccessProviderType("aliyun")
	AccessProviderTypeAWS          = AccessProviderType("aws")
	AccessProviderTypeAzure        = AccessProviderType("azure")
	AccessProviderTypeBaiduCloud   = AccessProviderType("baiducloud")
	AccessProviderTypeBaishan      = AccessProviderType("baishan")
	AccessProviderTypeBaotaPanel   = AccessProviderType("baotapanel")
	AccessProviderTypeBytePlus     = AccessProviderType("byteplus")
	AccessProviderTypeCacheFly     = AccessProviderType("cachefly")
	AccessProviderTypeCdnfly       = AccessProviderType("cdnfly")
	AccessProviderTypeCloudflare   = AccessProviderType("cloudflare")
	AccessProviderTypeClouDNS      = AccessProviderType("cloudns")
	AccessProviderTypeCMCCCloud    = AccessProviderType("cmcccloud")
	AccessProviderTypeCTCCCloud    = AccessProviderType("ctcccloud") // 联通云（预留）
	AccessProviderTypeCUCCCloud    = AccessProviderType("cucccloud") // 天翼云（预留）
	AccessProviderTypeDNSLA        = AccessProviderType("dnsla")
	AccessProviderTypeDogeCloud    = AccessProviderType("dogecloud")
	AccessProviderTypeEdgio        = AccessProviderType("edgio")
	AccessProviderTypeFastly       = AccessProviderType("fastly") // Fastly（预留）
	AccessProviderTypeGname        = AccessProviderType("gname")
	AccessProviderTypeGcore        = AccessProviderType("gcore")
	AccessProviderTypeGoDaddy      = AccessProviderType("godaddy")
	AccessProviderTypeGoEdge       = AccessProviderType("goedge") // GoEdge（预留）
	AccessProviderTypeHuaweiCloud  = AccessProviderType("huaweicloud")
	AccessProviderTypeJDCloud      = AccessProviderType("jdcloud")
	AccessProviderTypeKubernetes   = AccessProviderType("k8s")
	AccessProviderTypeLocal        = AccessProviderType("local")
	AccessProviderTypeNamecheap    = AccessProviderType("namecheap")
	AccessProviderTypeNameDotCom   = AccessProviderType("namedotcom")
	AccessProviderTypeNameSilo     = AccessProviderType("namesilo")
	AccessProviderTypeNS1          = AccessProviderType("ns1")
	AccessProviderTypePowerDNS     = AccessProviderType("powerdns")
	AccessProviderTypeQiniu        = AccessProviderType("qiniu")
	AccessProviderTypeQingCloud    = AccessProviderType("qingcloud") // 青云（预留）
	AccessProviderTypeRainYun      = AccessProviderType("rainyun")
	AccessProviderTypeSafeLine     = AccessProviderType("safeline")
	AccessProviderTypeSSH          = AccessProviderType("ssh")
	AccessProviderTypeTencentCloud = AccessProviderType("tencentcloud")
	AccessProviderTypeUCloud       = AccessProviderType("ucloud")
	AccessProviderTypeVolcEngine   = AccessProviderType("volcengine")
	AccessProviderTypeWebhook      = AccessProviderType("webhook")
	AccessProviderTypeWestcn       = AccessProviderType("westcn")
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
	ApplyDNSProviderTypeAzure           = ApplyDNSProviderType("azure") // 兼容旧值，等同于 [ApplyDNSProviderTypeAzure]
	ApplyDNSProviderTypeAzureDNS        = ApplyDNSProviderType("azure-dns")
	ApplyDNSProviderTypeBaiduCloud      = ApplyDNSProviderType("baiducloud") // 兼容旧值，等同于 [ApplyDNSProviderTypeBaiduCloudDNS]
	ApplyDNSProviderTypeBaiduCloudDNS   = ApplyDNSProviderType("baiducloud-dns")
	ApplyDNSProviderTypeCloudflare      = ApplyDNSProviderType("cloudflare")
	ApplyDNSProviderTypeClouDNS         = ApplyDNSProviderType("cloudns")
	ApplyDNSProviderTypeCMCCCloud       = ApplyDNSProviderType("cmcccloud")
	ApplyDNSProviderTypeDNSLA           = ApplyDNSProviderType("dnsla")
	ApplyDNSProviderTypeGcore           = ApplyDNSProviderType("gcore")
	ApplyDNSProviderTypeGname           = ApplyDNSProviderType("gname")
	ApplyDNSProviderTypeGoDaddy         = ApplyDNSProviderType("godaddy")
	ApplyDNSProviderTypeHuaweiCloud     = ApplyDNSProviderType("huaweicloud") // 兼容旧值，等同于 [ApplyDNSProviderTypeHuaweiCloudDNS]
	ApplyDNSProviderTypeHuaweiCloudDNS  = ApplyDNSProviderType("huaweicloud-dns")
	ApplyDNSProviderTypeJDCloud         = ApplyDNSProviderType("jdcloud") // 兼容旧值，等同于 [ApplyDNSProviderTypeJDCloudDNS]
	ApplyDNSProviderTypeJDCloudDNS      = ApplyDNSProviderType("jdcloud-dns")
	ApplyDNSProviderTypeNamecheap       = ApplyDNSProviderType("namecheap")
	ApplyDNSProviderTypeNameDotCom      = ApplyDNSProviderType("namedotcom")
	ApplyDNSProviderTypeNameSilo        = ApplyDNSProviderType("namesilo")
	ApplyDNSProviderTypeNS1             = ApplyDNSProviderType("ns1")
	ApplyDNSProviderTypePowerDNS        = ApplyDNSProviderType("powerdns")
	ApplyDNSProviderTypeRainYun         = ApplyDNSProviderType("rainyun")
	ApplyDNSProviderTypeTencentCloud    = ApplyDNSProviderType("tencentcloud") // 兼容旧值，等同于 [ApplyDNSProviderTypeTencentCloudDNS]
	ApplyDNSProviderTypeTencentCloudDNS = ApplyDNSProviderType("tencentcloud-dns")
	ApplyDNSProviderTypeVolcEngine      = ApplyDNSProviderType("volcengine") // 兼容旧值，等同于 [ApplyDNSProviderTypeVolcEngineDNS]
	ApplyDNSProviderTypeVolcEngineDNS   = ApplyDNSProviderType("volcengine-dns")
	ApplyDNSProviderTypeWestcn          = ApplyDNSProviderType("westcn")
)

type DeployProviderType string

/*
部署目标提供商常量值。
短横线前的部分始终等于授权提供商类型。

	注意：如果追加新的常量值，请保持以 ASCII 排序。
	NOTICE: If you add new constant, please keep ASCII order.
*/
const (
	DeployProviderType1PanelConsole         = DeployProviderType("1panel-console")
	DeployProviderType1PanelSite            = DeployProviderType("1panel-site")
	DeployProviderTypeAliyunALB             = DeployProviderType("aliyun-alb")
	DeployProviderTypeAliyunCASDeploy       = DeployProviderType("aliyun-casdeploy")
	DeployProviderTypeAliyunCDN             = DeployProviderType("aliyun-cdn")
	DeployProviderTypeAliyunCLB             = DeployProviderType("aliyun-clb")
	DeployProviderTypeAliyunDCDN            = DeployProviderType("aliyun-dcdn")
	DeployProviderTypeAliyunESA             = DeployProviderType("aliyun-esa")
	DeployProviderTypeAliyunFC              = DeployProviderType("aliyun-fc")
	DeployProviderTypeAliyunLive            = DeployProviderType("aliyun-live")
	DeployProviderTypeAliyunNLB             = DeployProviderType("aliyun-nlb")
	DeployProviderTypeAliyunOSS             = DeployProviderType("aliyun-oss")
	DeployProviderTypeAliyunVOD             = DeployProviderType("aliyun-vod")
	DeployProviderTypeAliyunWAF             = DeployProviderType("aliyun-waf")
	DeployProviderTypeAWSCloudFront         = DeployProviderType("aws-cloudfront")
	DeployProviderTypeBaiduCloudCDN         = DeployProviderType("baiducloud-cdn")
	DeployProviderTypeBaishanCDN            = DeployProviderType("baishan-cdn")
	DeployProviderTypeBaotaPanelConsole     = DeployProviderType("baotapanel-console")
	DeployProviderTypeBaotaPanelSite        = DeployProviderType("baotapanel-site")
	DeployProviderTypeBytePlusCDN           = DeployProviderType("byteplus-cdn")
	DeployProviderTypeCacheFly              = DeployProviderType("cachefly")
	DeployProviderTypeCdnfly                = DeployProviderType("cdnfly")
	DeployProviderTypeDogeCloudCDN          = DeployProviderType("dogecloud-cdn")
	DeployProviderTypeEdgioApplications     = DeployProviderType("edgio-applications")
	DeployProviderTypeGcoreCDN              = DeployProviderType("gcore-cdn")
	DeployProviderTypeHuaweiCloudCDN        = DeployProviderType("huaweicloud-cdn")
	DeployProviderTypeHuaweiCloudELB        = DeployProviderType("huaweicloud-elb")
	DeployProviderTypeHuaweiCloudWAF        = DeployProviderType("huaweicloud-waf")
	DeployProviderTypeJDCloudALB            = DeployProviderType("jdcloud-alb")
	DeployProviderTypeJDCloudCDN            = DeployProviderType("jdcloud-cdn")
	DeployProviderTypeJDCloudLive           = DeployProviderType("jdcloud-live")
	DeployProviderTypeJDCloudVOD            = DeployProviderType("jdcloud-vod")
	DeployProviderTypeKubernetesSecret      = DeployProviderType("k8s-secret")
	DeployProviderTypeLocal                 = DeployProviderType("local")
	DeployProviderTypeQiniuCDN              = DeployProviderType("qiniu-cdn")
	DeployProviderTypeQiniuPili             = DeployProviderType("qiniu-pili")
	DeployProviderTypeSafeLine              = DeployProviderType("safeline")
	DeployProviderTypeSSH                   = DeployProviderType("ssh")
	DeployProviderTypeTencentCloudCDN       = DeployProviderType("tencentcloud-cdn")
	DeployProviderTypeTencentCloudCLB       = DeployProviderType("tencentcloud-clb")
	DeployProviderTypeTencentCloudCOS       = DeployProviderType("tencentcloud-cos")
	DeployProviderTypeTencentCloudCSS       = DeployProviderType("tencentcloud-css")
	DeployProviderTypeTencentCloudECDN      = DeployProviderType("tencentcloud-ecdn")
	DeployProviderTypeTencentCloudEO        = DeployProviderType("tencentcloud-eo")
	DeployProviderTypeTencentCloudSCF       = DeployProviderType("tencentcloud-scf")
	DeployProviderTypeTencentCloudSSLDeploy = DeployProviderType("tencentcloud-ssldeploy")
	DeployProviderTypeTencentCloudVOD       = DeployProviderType("tencentcloud-vod")
	DeployProviderTypeTencentCloudWAF       = DeployProviderType("tencentcloud-waf")
	DeployProviderTypeUCloudUCDN            = DeployProviderType("ucloud-ucdn")
	DeployProviderTypeUCloudUS3             = DeployProviderType("ucloud-us3")
	DeployProviderTypeVolcEngineCDN         = DeployProviderType("volcengine-cdn")
	DeployProviderTypeVolcEngineCLB         = DeployProviderType("volcengine-clb")
	DeployProviderTypeVolcEngineDCDN        = DeployProviderType("volcengine-dcdn")
	DeployProviderTypeVolcEngineImageX      = DeployProviderType("volcengine-imagex")
	DeployProviderTypeVolcEngineLive        = DeployProviderType("volcengine-live")
	DeployProviderTypeVolcEngineTOS         = DeployProviderType("volcengine-tos")
	DeployProviderTypeWebhook               = DeployProviderType("webhook")
)
