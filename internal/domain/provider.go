package domain

type AccessProviderType string

/*
授权提供商类型常量值。

	注意：如果追加新的常量值，请保持以 ASCII 排序。
	NOTICE: If you add new constant, please keep ASCII order.
*/
const (
	AccessProviderType1Panel              = AccessProviderType("1panel")
	AccessProviderTypeACMEHttpReq         = AccessProviderType("acmehttpreq")
	AccessProviderTypeAkamai              = AccessProviderType("akamai") // Akamai（预留）
	AccessProviderTypeAliyun              = AccessProviderType("aliyun")
	AccessProviderTypeAWS                 = AccessProviderType("aws")
	AccessProviderTypeAzure               = AccessProviderType("azure")
	AccessProviderTypeBaiduCloud          = AccessProviderType("baiducloud")
	AccessProviderTypeBaishan             = AccessProviderType("baishan")
	AccessProviderTypeBaotaPanel          = AccessProviderType("baotapanel")
	AccessProviderTypeBytePlus            = AccessProviderType("byteplus")
	AccessProviderTypeBunny               = AccessProviderType("bunny")
	AccessProviderTypeBuypass             = AccessProviderType("buypass")
	AccessProviderTypeCacheFly            = AccessProviderType("cachefly")
	AccessProviderTypeCdnfly              = AccessProviderType("cdnfly")
	AccessProviderTypeCloudflare          = AccessProviderType("cloudflare")
	AccessProviderTypeClouDNS             = AccessProviderType("cloudns")
	AccessProviderTypeCMCCCloud           = AccessProviderType("cmcccloud")
	AccessProviderTypeCTCCCloud           = AccessProviderType("ctcccloud") // 联通云（预留）
	AccessProviderTypeCUCCCloud           = AccessProviderType("cucccloud") // 天翼云（预留）
	AccessProviderTypeDeSEC               = AccessProviderType("desec")
	AccessProviderTypeDNSLA               = AccessProviderType("dnsla")
	AccessProviderTypeDogeCloud           = AccessProviderType("dogecloud")
	AccessProviderTypeDynv6               = AccessProviderType("dynv6")
	AccessProviderTypeEdgio               = AccessProviderType("edgio")
	AccessProviderTypeFastly              = AccessProviderType("fastly") // Fastly（预留）
	AccessProviderTypeGname               = AccessProviderType("gname")
	AccessProviderTypeGcore               = AccessProviderType("gcore")
	AccessProviderTypeGoDaddy             = AccessProviderType("godaddy")
	AccessProviderTypeGoEdge              = AccessProviderType("goedge") // GoEdge（预留）
	AccessProviderTypeGoogleTrustServices = AccessProviderType("googletrustservices")
	AccessProviderTypeHuaweiCloud         = AccessProviderType("huaweicloud")
	AccessProviderTypeJDCloud             = AccessProviderType("jdcloud")
	AccessProviderTypeKubernetes          = AccessProviderType("k8s")
	AccessProviderTypeLetsEncrypt         = AccessProviderType("letsencrypt")
	AccessProviderTypeLetsEncryptStaging  = AccessProviderType("letsencryptstaging")
	AccessProviderTypeLocal               = AccessProviderType("local")
	AccessProviderTypeNamecheap           = AccessProviderType("namecheap")
	AccessProviderTypeNameDotCom          = AccessProviderType("namedotcom")
	AccessProviderTypeNameSilo            = AccessProviderType("namesilo")
	AccessProviderTypeNS1                 = AccessProviderType("ns1")
	AccessProviderTypePorkbun             = AccessProviderType("porkbun")
	AccessProviderTypePowerDNS            = AccessProviderType("powerdns")
	AccessProviderTypeQiniu               = AccessProviderType("qiniu")
	AccessProviderTypeQingCloud           = AccessProviderType("qingcloud") // 青云（预留）
	AccessProviderTypeRainYun             = AccessProviderType("rainyun")
	AccessProviderTypeSafeLine            = AccessProviderType("safeline")
	AccessProviderTypeSSH                 = AccessProviderType("ssh")
	AccessProviderTypeSSLCOM              = AccessProviderType("sslcom")
	AccessProviderTypeTencentCloud        = AccessProviderType("tencentcloud")
	AccessProviderTypeUCloud              = AccessProviderType("ucloud")
	AccessProviderTypeUpyun               = AccessProviderType("upyun")
	AccessProviderTypeVercel              = AccessProviderType("vercel")
	AccessProviderTypeVolcEngine          = AccessProviderType("volcengine")
	AccessProviderTypeWangsu              = AccessProviderType("wangsu")
	AccessProviderTypeWebhook             = AccessProviderType("webhook")
	AccessProviderTypeWestcn              = AccessProviderType("westcn")
	AccessProviderTypeZeroSSL             = AccessProviderType("zerossl")
)

type ApplyCAProviderType string

/*
申请证书 CA 提供商常量值。
短横线前的部分始终等于授权提供商类型。

	注意：如果追加新的常量值，请保持以 ASCII 排序。
	NOTICE: If you add new constant, please keep ASCII order.
*/
const (
	ApplyCAProviderTypeBuypass             = ApplyCAProviderType(AccessProviderTypeBuypass)
	ApplyCAProviderTypeGoogleTrustServices = ApplyCAProviderType(AccessProviderTypeGoogleTrustServices)
	ApplyCAProviderTypeLetsEncrypt         = ApplyCAProviderType(AccessProviderTypeLetsEncrypt)
	ApplyCAProviderTypeLetsEncryptStaging  = ApplyCAProviderType(AccessProviderTypeLetsEncryptStaging)
	ApplyCAProviderTypeSSLCom              = ApplyCAProviderType(AccessProviderTypeSSLCOM)
	ApplyCAProviderTypeZeroSSL             = ApplyCAProviderType(AccessProviderTypeZeroSSL)
)

type ApplyDNSProviderType string

/*
申请证书 DNS 提供商常量值。
短横线前的部分始终等于授权提供商类型。

	注意：如果追加新的常量值，请保持以 ASCII 排序。
	NOTICE: If you add new constant, please keep ASCII order.
*/
const (
	ApplyDNSProviderTypeACMEHttpReq     = ApplyDNSProviderType(AccessProviderTypeACMEHttpReq)
	ApplyDNSProviderTypeAliyun          = ApplyDNSProviderType(AccessProviderTypeAliyun) // 兼容旧值，等同于 [ApplyDNSProviderTypeAliyunDNS]
	ApplyDNSProviderTypeAliyunDNS       = ApplyDNSProviderType(AccessProviderTypeAliyun + "-dns")
	ApplyDNSProviderTypeAWS             = ApplyDNSProviderType(AccessProviderTypeAWS) // 兼容旧值，等同于 [ApplyDNSProviderTypeAWSRoute53]
	ApplyDNSProviderTypeAWSRoute53      = ApplyDNSProviderType(AccessProviderTypeAWS + "-route53")
	ApplyDNSProviderTypeAzure           = ApplyDNSProviderType(AccessProviderTypeAzure) // 兼容旧值，等同于 [ApplyDNSProviderTypeAzure]
	ApplyDNSProviderTypeAzureDNS        = ApplyDNSProviderType(AccessProviderTypeAzure + "-dns")
	ApplyDNSProviderTypeBaiduCloud      = ApplyDNSProviderType(AccessProviderTypeBaiduCloud) // 兼容旧值，等同于 [ApplyDNSProviderTypeBaiduCloudDNS]
	ApplyDNSProviderTypeBaiduCloudDNS   = ApplyDNSProviderType(AccessProviderTypeBaiduCloud + "-dns")
	ApplyDNSProviderTypeBunny           = ApplyDNSProviderType(AccessProviderTypeBunny)
	ApplyDNSProviderTypeCloudflare      = ApplyDNSProviderType(AccessProviderTypeCloudflare)
	ApplyDNSProviderTypeClouDNS         = ApplyDNSProviderType(AccessProviderTypeClouDNS)
	ApplyDNSProviderTypeCMCCCloud       = ApplyDNSProviderType(AccessProviderTypeCMCCCloud)
	ApplyDNSProviderTypeDeSEC           = ApplyDNSProviderType(AccessProviderTypeDeSEC)
	ApplyDNSProviderTypeDNSLA           = ApplyDNSProviderType(AccessProviderTypeDNSLA)
	ApplyDNSProviderTypeDynv6           = ApplyDNSProviderType(AccessProviderTypeDynv6)
	ApplyDNSProviderTypeGcore           = ApplyDNSProviderType(AccessProviderTypeGcore)
	ApplyDNSProviderTypeGname           = ApplyDNSProviderType(AccessProviderTypeGname)
	ApplyDNSProviderTypeGoDaddy         = ApplyDNSProviderType(AccessProviderTypeGoDaddy)
	ApplyDNSProviderTypeHuaweiCloud     = ApplyDNSProviderType(AccessProviderTypeHuaweiCloud) // 兼容旧值，等同于 [ApplyDNSProviderTypeHuaweiCloudDNS]
	ApplyDNSProviderTypeHuaweiCloudDNS  = ApplyDNSProviderType(AccessProviderTypeHuaweiCloud + "-dns")
	ApplyDNSProviderTypeJDCloud         = ApplyDNSProviderType(AccessProviderTypeJDCloud) // 兼容旧值，等同于 [ApplyDNSProviderTypeJDCloudDNS]
	ApplyDNSProviderTypeJDCloudDNS      = ApplyDNSProviderType(AccessProviderTypeJDCloud + "-dns")
	ApplyDNSProviderTypeNamecheap       = ApplyDNSProviderType(AccessProviderTypeNamecheap)
	ApplyDNSProviderTypeNameDotCom      = ApplyDNSProviderType(AccessProviderTypeNameDotCom)
	ApplyDNSProviderTypeNameSilo        = ApplyDNSProviderType(AccessProviderTypeNameSilo)
	ApplyDNSProviderTypeNS1             = ApplyDNSProviderType(AccessProviderTypeNS1)
	ApplyDNSProviderTypePorkbun         = ApplyDNSProviderType(AccessProviderTypePorkbun)
	ApplyDNSProviderTypePowerDNS        = ApplyDNSProviderType(AccessProviderTypePowerDNS)
	ApplyDNSProviderTypeRainYun         = ApplyDNSProviderType(AccessProviderTypeRainYun)
	ApplyDNSProviderTypeTencentCloud    = ApplyDNSProviderType(AccessProviderTypeTencentCloud) // 兼容旧值，等同于 [ApplyDNSProviderTypeTencentCloudDNS]
	ApplyDNSProviderTypeTencentCloudDNS = ApplyDNSProviderType(AccessProviderTypeTencentCloud + "-dns")
	ApplyDNSProviderTypeTencentCloudEO  = ApplyDNSProviderType(AccessProviderTypeTencentCloud + "-eo")
	ApplyDNSProviderTypeVercel          = ApplyDNSProviderType(AccessProviderTypeVercel)
	ApplyDNSProviderTypeVolcEngine      = ApplyDNSProviderType(AccessProviderTypeVolcEngine) // 兼容旧值，等同于 [ApplyDNSProviderTypeVolcEngineDNS]
	ApplyDNSProviderTypeVolcEngineDNS   = ApplyDNSProviderType(AccessProviderTypeVolcEngine + "-dns")
	ApplyDNSProviderTypeWestcn          = ApplyDNSProviderType(AccessProviderTypeWestcn)
)

type DeployProviderType string

/*
部署证书主机提供商常量值。
短横线前的部分始终等于授权提供商类型。

	注意：如果追加新的常量值，请保持以 ASCII 排序。
	NOTICE: If you add new constant, please keep ASCII order.
*/
const (
	DeployProviderType1PanelConsole         = DeployProviderType(AccessProviderType1Panel + "-console")
	DeployProviderType1PanelSite            = DeployProviderType(AccessProviderType1Panel + "-site")
	DeployProviderTypeAliyunALB             = DeployProviderType(AccessProviderTypeAliyun + "-alb")
	DeployProviderTypeAliyunAPIGW           = DeployProviderType(AccessProviderTypeAliyun + "-apigw")
	DeployProviderTypeAliyunCAS             = DeployProviderType(AccessProviderTypeAliyun + "-cas")
	DeployProviderTypeAliyunCASDeploy       = DeployProviderType(AccessProviderTypeAliyun + "-casdeploy")
	DeployProviderTypeAliyunCDN             = DeployProviderType(AccessProviderTypeAliyun + "-cdn")
	DeployProviderTypeAliyunCLB             = DeployProviderType(AccessProviderTypeAliyun + "-clb")
	DeployProviderTypeAliyunDCDN            = DeployProviderType(AccessProviderTypeAliyun + "-dcdn")
	DeployProviderTypeAliyunESA             = DeployProviderType(AccessProviderTypeAliyun + "-esa")
	DeployProviderTypeAliyunFC              = DeployProviderType(AccessProviderTypeAliyun + "-fc")
	DeployProviderTypeAliyunLive            = DeployProviderType(AccessProviderTypeAliyun + "-live")
	DeployProviderTypeAliyunNLB             = DeployProviderType(AccessProviderTypeAliyun + "-nlb")
	DeployProviderTypeAliyunOSS             = DeployProviderType(AccessProviderTypeAliyun + "-oss")
	DeployProviderTypeAliyunVOD             = DeployProviderType(AccessProviderTypeAliyun + "-vod")
	DeployProviderTypeAliyunWAF             = DeployProviderType(AccessProviderTypeAliyun + "-waf")
	DeployProviderTypeAWSACM                = DeployProviderType(AccessProviderTypeAWS + "-acm")
	DeployProviderTypeAWSCloudFront         = DeployProviderType(AccessProviderTypeAWS + "-cloudfront")
	DeployProviderTypeAzureKeyVault         = DeployProviderType(AccessProviderTypeAzure + "-keyvault")
	DeployProviderTypeBaiduCloudAppBLB      = DeployProviderType(AccessProviderTypeBaiduCloud + "-appblb")
	DeployProviderTypeBaiduCloudBLB         = DeployProviderType(AccessProviderTypeBaiduCloud + "-blb")
	DeployProviderTypeBaiduCloudCDN         = DeployProviderType(AccessProviderTypeBaiduCloud + "-cdn")
	DeployProviderTypeBaiduCloudCert        = DeployProviderType(AccessProviderTypeBaiduCloud + "-cert")
	DeployProviderTypeBaishanCDN            = DeployProviderType(AccessProviderTypeBaishan + "-cdn")
	DeployProviderTypeBaotaPanelConsole     = DeployProviderType(AccessProviderTypeBaotaPanel + "-console")
	DeployProviderTypeBaotaPanelSite        = DeployProviderType(AccessProviderTypeBaotaPanel + "-site")
	DeployProviderTypeBunnyCDN              = DeployProviderType(AccessProviderTypeBunny + "-cdn")
	DeployProviderTypeBytePlusCDN           = DeployProviderType(AccessProviderTypeBytePlus + "-cdn")
	DeployProviderTypeCacheFly              = DeployProviderType(AccessProviderTypeCacheFly)
	DeployProviderTypeCdnfly                = DeployProviderType(AccessProviderTypeCdnfly)
	DeployProviderTypeDogeCloudCDN          = DeployProviderType(AccessProviderTypeDogeCloud + "-cdn")
	DeployProviderTypeEdgioApplications     = DeployProviderType(AccessProviderTypeEdgio + "-applications")
	DeployProviderTypeGcoreCDN              = DeployProviderType(AccessProviderTypeGcore + "-cdn")
	DeployProviderTypeHuaweiCloudCDN        = DeployProviderType(AccessProviderTypeHuaweiCloud + "-cdn")
	DeployProviderTypeHuaweiCloudELB        = DeployProviderType(AccessProviderTypeHuaweiCloud + "-elb")
	DeployProviderTypeHuaweiCloudSCM        = DeployProviderType(AccessProviderTypeHuaweiCloud + "-scm")
	DeployProviderTypeHuaweiCloudWAF        = DeployProviderType(AccessProviderTypeHuaweiCloud + "-waf")
	DeployProviderTypeJDCloudALB            = DeployProviderType(AccessProviderTypeJDCloud + "-alb")
	DeployProviderTypeJDCloudCDN            = DeployProviderType(AccessProviderTypeJDCloud + "-cdn")
	DeployProviderTypeJDCloudLive           = DeployProviderType(AccessProviderTypeJDCloud + "-live")
	DeployProviderTypeJDCloudVOD            = DeployProviderType(AccessProviderTypeJDCloud + "-vod")
	DeployProviderTypeKubernetesSecret      = DeployProviderType(AccessProviderTypeKubernetes + "-secret")
	DeployProviderTypeLocal                 = DeployProviderType(AccessProviderTypeLocal)
	DeployProviderTypeQiniuCDN              = DeployProviderType(AccessProviderTypeQiniu + "-cdn")
	DeployProviderTypeQiniuKodo             = DeployProviderType(AccessProviderTypeQiniu + "-kodo")
	DeployProviderTypeQiniuPili             = DeployProviderType(AccessProviderTypeQiniu + "-pili")
	DeployProviderTypeRainYunRCDN           = DeployProviderType(AccessProviderTypeRainYun + "-rcdn")
	DeployProviderTypeSafeLine              = DeployProviderType(AccessProviderTypeSafeLine)
	DeployProviderTypeSSH                   = DeployProviderType(AccessProviderTypeSSH)
	DeployProviderTypeTencentCloudCDN       = DeployProviderType(AccessProviderTypeTencentCloud + "-cdn")
	DeployProviderTypeTencentCloudCLB       = DeployProviderType(AccessProviderTypeTencentCloud + "-clb")
	DeployProviderTypeTencentCloudCOS       = DeployProviderType(AccessProviderTypeTencentCloud + "-cos")
	DeployProviderTypeTencentCloudCSS       = DeployProviderType(AccessProviderTypeTencentCloud + "-css")
	DeployProviderTypeTencentCloudECDN      = DeployProviderType(AccessProviderTypeTencentCloud + "-ecdn")
	DeployProviderTypeTencentCloudEO        = DeployProviderType(AccessProviderTypeTencentCloud + "-eo")
	DeployProviderTypeTencentCloudSCF       = DeployProviderType(AccessProviderTypeTencentCloud + "-scf")
	DeployProviderTypeTencentCloudSSL       = DeployProviderType(AccessProviderTypeTencentCloud + "-ssl")
	DeployProviderTypeTencentCloudSSLDeploy = DeployProviderType(AccessProviderTypeTencentCloud + "-ssldeploy")
	DeployProviderTypeTencentCloudVOD       = DeployProviderType(AccessProviderTypeTencentCloud + "-vod")
	DeployProviderTypeTencentCloudWAF       = DeployProviderType(AccessProviderTypeTencentCloud + "-waf")
	DeployProviderTypeUCloudUCDN            = DeployProviderType(AccessProviderTypeUCloud + "-ucdn")
	DeployProviderTypeUCloudUS3             = DeployProviderType(AccessProviderTypeUCloud + "-us3")
	DeployProviderTypeUpyunCDN              = DeployProviderType(AccessProviderTypeUpyun + "-cdn")
	DeployProviderTypeUpyunFile             = DeployProviderType(AccessProviderTypeUpyun + "-file")
	DeployProviderTypeVolcEngineALB         = DeployProviderType(AccessProviderTypeVolcEngine + "-alb")
	DeployProviderTypeVolcEngineCDN         = DeployProviderType(AccessProviderTypeVolcEngine + "-cdn")
	DeployProviderTypeVolcEngineCertCenter  = DeployProviderType(AccessProviderTypeVolcEngine + "-certcenter")
	DeployProviderTypeVolcEngineCLB         = DeployProviderType(AccessProviderTypeVolcEngine + "-clb")
	DeployProviderTypeVolcEngineDCDN        = DeployProviderType(AccessProviderTypeVolcEngine + "-dcdn")
	DeployProviderTypeVolcEngineImageX      = DeployProviderType(AccessProviderTypeVolcEngine + "-imagex")
	DeployProviderTypeVolcEngineLive        = DeployProviderType(AccessProviderTypeVolcEngine + "-live")
	DeployProviderTypeVolcEngineTOS         = DeployProviderType(AccessProviderTypeVolcEngine + "-tos")
	DeployProviderTypeWangsuCDNPro          = DeployProviderType(AccessProviderTypeWangsu + "-cdnpro")
	DeployProviderTypeWebhook               = DeployProviderType(AccessProviderTypeWebhook)
)

type NotifyProviderType string

/*
消息通知提供商常量值。
短横线前的部分始终等于授权提供商类型。

	注意：如果追加新的常量值，请保持以 ASCII 排序。
	NOTICE: If you add new constant, please keep ASCII order.
*/
const (
	NotifyProviderTypeWebhook = NotifyProviderType(AccessProviderTypeWebhook)
)
