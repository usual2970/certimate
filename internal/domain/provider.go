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

type CAProviderType string

/*
证书颁发机构提供商常量值。
短横线前的部分始终等于授权提供商类型。

	注意：如果追加新的常量值，请保持以 ASCII 排序。
	NOTICE: If you add new constant, please keep ASCII order.
*/
const (
	CAProviderTypeBuypass             = CAProviderType(AccessProviderTypeBuypass)
	CAProviderTypeGoogleTrustServices = CAProviderType(AccessProviderTypeGoogleTrustServices)
	CAProviderTypeLetsEncrypt         = CAProviderType(AccessProviderTypeLetsEncrypt)
	CAProviderTypeLetsEncryptStaging  = CAProviderType(AccessProviderTypeLetsEncryptStaging)
	CAProviderTypeSSLCom              = CAProviderType(AccessProviderTypeSSLCOM)
	CAProviderTypeZeroSSL             = CAProviderType(AccessProviderTypeZeroSSL)
)

type AcmeDns01ProviderType string

/*
ACME DNS-01 提供商常量值。
短横线前的部分始终等于授权提供商类型。

	注意：如果追加新的常量值，请保持以 ASCII 排序。
	NOTICE: If you add new constant, please keep ASCII order.
*/
const (
	AcmeDns01ProviderTypeACMEHttpReq     = AcmeDns01ProviderType(AccessProviderTypeACMEHttpReq)
	AcmeDns01ProviderTypeAliyun          = AcmeDns01ProviderType(AccessProviderTypeAliyun) // 兼容旧值，等同于 [AcmeDns01ProviderTypeAliyunDNS]
	AcmeDns01ProviderTypeAliyunDNS       = AcmeDns01ProviderType(AccessProviderTypeAliyun + "-dns")
	AcmeDns01ProviderTypeAWS             = AcmeDns01ProviderType(AccessProviderTypeAWS) // 兼容旧值，等同于 [AcmeDns01ProviderTypeAWSRoute53]
	AcmeDns01ProviderTypeAWSRoute53      = AcmeDns01ProviderType(AccessProviderTypeAWS + "-route53")
	AcmeDns01ProviderTypeAzure           = AcmeDns01ProviderType(AccessProviderTypeAzure) // 兼容旧值，等同于 [AcmeDns01ProviderTypeAzure]
	AcmeDns01ProviderTypeAzureDNS        = AcmeDns01ProviderType(AccessProviderTypeAzure + "-dns")
	AcmeDns01ProviderTypeBaiduCloud      = AcmeDns01ProviderType(AccessProviderTypeBaiduCloud) // 兼容旧值，等同于 [AcmeDns01ProviderTypeBaiduCloudDNS]
	AcmeDns01ProviderTypeBaiduCloudDNS   = AcmeDns01ProviderType(AccessProviderTypeBaiduCloud + "-dns")
	AcmeDns01ProviderTypeBunny           = AcmeDns01ProviderType(AccessProviderTypeBunny)
	AcmeDns01ProviderTypeCloudflare      = AcmeDns01ProviderType(AccessProviderTypeCloudflare)
	AcmeDns01ProviderTypeClouDNS         = AcmeDns01ProviderType(AccessProviderTypeClouDNS)
	AcmeDns01ProviderTypeCMCCCloud       = AcmeDns01ProviderType(AccessProviderTypeCMCCCloud)
	AcmeDns01ProviderTypeDeSEC           = AcmeDns01ProviderType(AccessProviderTypeDeSEC)
	AcmeDns01ProviderTypeDNSLA           = AcmeDns01ProviderType(AccessProviderTypeDNSLA)
	AcmeDns01ProviderTypeDynv6           = AcmeDns01ProviderType(AccessProviderTypeDynv6)
	AcmeDns01ProviderTypeGcore           = AcmeDns01ProviderType(AccessProviderTypeGcore)
	AcmeDns01ProviderTypeGname           = AcmeDns01ProviderType(AccessProviderTypeGname)
	AcmeDns01ProviderTypeGoDaddy         = AcmeDns01ProviderType(AccessProviderTypeGoDaddy)
	AcmeDns01ProviderTypeHuaweiCloud     = AcmeDns01ProviderType(AccessProviderTypeHuaweiCloud) // 兼容旧值，等同于 [AcmeDns01ProviderTypeHuaweiCloudDNS]
	AcmeDns01ProviderTypeHuaweiCloudDNS  = AcmeDns01ProviderType(AccessProviderTypeHuaweiCloud + "-dns")
	AcmeDns01ProviderTypeJDCloud         = AcmeDns01ProviderType(AccessProviderTypeJDCloud) // 兼容旧值，等同于 [AcmeDns01ProviderTypeJDCloudDNS]
	AcmeDns01ProviderTypeJDCloudDNS      = AcmeDns01ProviderType(AccessProviderTypeJDCloud + "-dns")
	AcmeDns01ProviderTypeNamecheap       = AcmeDns01ProviderType(AccessProviderTypeNamecheap)
	AcmeDns01ProviderTypeNameDotCom      = AcmeDns01ProviderType(AccessProviderTypeNameDotCom)
	AcmeDns01ProviderTypeNameSilo        = AcmeDns01ProviderType(AccessProviderTypeNameSilo)
	AcmeDns01ProviderTypeNS1             = AcmeDns01ProviderType(AccessProviderTypeNS1)
	AcmeDns01ProviderTypePorkbun         = AcmeDns01ProviderType(AccessProviderTypePorkbun)
	AcmeDns01ProviderTypePowerDNS        = AcmeDns01ProviderType(AccessProviderTypePowerDNS)
	AcmeDns01ProviderTypeRainYun         = AcmeDns01ProviderType(AccessProviderTypeRainYun)
	AcmeDns01ProviderTypeTencentCloud    = AcmeDns01ProviderType(AccessProviderTypeTencentCloud) // 兼容旧值，等同于 [AcmeDns01ProviderTypeTencentCloudDNS]
	AcmeDns01ProviderTypeTencentCloudDNS = AcmeDns01ProviderType(AccessProviderTypeTencentCloud + "-dns")
	AcmeDns01ProviderTypeTencentCloudEO  = AcmeDns01ProviderType(AccessProviderTypeTencentCloud + "-eo")
	AcmeDns01ProviderTypeVercel          = AcmeDns01ProviderType(AccessProviderTypeVercel)
	AcmeDns01ProviderTypeVolcEngine      = AcmeDns01ProviderType(AccessProviderTypeVolcEngine) // 兼容旧值，等同于 [AcmeDns01ProviderTypeVolcEngineDNS]
	AcmeDns01ProviderTypeVolcEngineDNS   = AcmeDns01ProviderType(AccessProviderTypeVolcEngine + "-dns")
	AcmeDns01ProviderTypeWestcn          = AcmeDns01ProviderType(AccessProviderTypeWestcn)
)

type DeploymentProviderType string

/*
部署证书主机提供商常量值。
短横线前的部分始终等于授权提供商类型。

	注意：如果追加新的常量值，请保持以 ASCII 排序。
	NOTICE: If you add new constant, please keep ASCII order.
*/
const (
	DeploymentProviderType1PanelConsole         = DeploymentProviderType(AccessProviderType1Panel + "-console")
	DeploymentProviderType1PanelSite            = DeploymentProviderType(AccessProviderType1Panel + "-site")
	DeploymentProviderTypeAliyunALB             = DeploymentProviderType(AccessProviderTypeAliyun + "-alb")
	DeploymentProviderTypeAliyunAPIGW           = DeploymentProviderType(AccessProviderTypeAliyun + "-apigw")
	DeploymentProviderTypeAliyunCAS             = DeploymentProviderType(AccessProviderTypeAliyun + "-cas")
	DeploymentProviderTypeAliyunCASDeploy       = DeploymentProviderType(AccessProviderTypeAliyun + "-casdeploy")
	DeploymentProviderTypeAliyunCDN             = DeploymentProviderType(AccessProviderTypeAliyun + "-cdn")
	DeploymentProviderTypeAliyunCLB             = DeploymentProviderType(AccessProviderTypeAliyun + "-clb")
	DeploymentProviderTypeAliyunDCDN            = DeploymentProviderType(AccessProviderTypeAliyun + "-dcdn")
	DeploymentProviderTypeAliyunESA             = DeploymentProviderType(AccessProviderTypeAliyun + "-esa")
	DeploymentProviderTypeAliyunFC              = DeploymentProviderType(AccessProviderTypeAliyun + "-fc")
	DeploymentProviderTypeAliyunLive            = DeploymentProviderType(AccessProviderTypeAliyun + "-live")
	DeploymentProviderTypeAliyunNLB             = DeploymentProviderType(AccessProviderTypeAliyun + "-nlb")
	DeploymentProviderTypeAliyunOSS             = DeploymentProviderType(AccessProviderTypeAliyun + "-oss")
	DeploymentProviderTypeAliyunVOD             = DeploymentProviderType(AccessProviderTypeAliyun + "-vod")
	DeploymentProviderTypeAliyunWAF             = DeploymentProviderType(AccessProviderTypeAliyun + "-waf")
	DeploymentProviderTypeAWSACM                = DeploymentProviderType(AccessProviderTypeAWS + "-acm")
	DeploymentProviderTypeAWSCloudFront         = DeploymentProviderType(AccessProviderTypeAWS + "-cloudfront")
	DeploymentProviderTypeAzureKeyVault         = DeploymentProviderType(AccessProviderTypeAzure + "-keyvault")
	DeploymentProviderTypeBaiduCloudAppBLB      = DeploymentProviderType(AccessProviderTypeBaiduCloud + "-appblb")
	DeploymentProviderTypeBaiduCloudBLB         = DeploymentProviderType(AccessProviderTypeBaiduCloud + "-blb")
	DeploymentProviderTypeBaiduCloudCDN         = DeploymentProviderType(AccessProviderTypeBaiduCloud + "-cdn")
	DeploymentProviderTypeBaiduCloudCert        = DeploymentProviderType(AccessProviderTypeBaiduCloud + "-cert")
	DeploymentProviderTypeBaishanCDN            = DeploymentProviderType(AccessProviderTypeBaishan + "-cdn")
	DeploymentProviderTypeBaotaPanelConsole     = DeploymentProviderType(AccessProviderTypeBaotaPanel + "-console")
	DeploymentProviderTypeBaotaPanelSite        = DeploymentProviderType(AccessProviderTypeBaotaPanel + "-site")
	DeploymentProviderTypeBunnyCDN              = DeploymentProviderType(AccessProviderTypeBunny + "-cdn")
	DeploymentProviderTypeBytePlusCDN           = DeploymentProviderType(AccessProviderTypeBytePlus + "-cdn")
	DeploymentProviderTypeCacheFly              = DeploymentProviderType(AccessProviderTypeCacheFly)
	DeploymentProviderTypeCdnfly                = DeploymentProviderType(AccessProviderTypeCdnfly)
	DeploymentProviderTypeDogeCloudCDN          = DeploymentProviderType(AccessProviderTypeDogeCloud + "-cdn")
	DeploymentProviderTypeEdgioApplications     = DeploymentProviderType(AccessProviderTypeEdgio + "-applications")
	DeploymentProviderTypeGcoreCDN              = DeploymentProviderType(AccessProviderTypeGcore + "-cdn")
	DeploymentProviderTypeHuaweiCloudCDN        = DeploymentProviderType(AccessProviderTypeHuaweiCloud + "-cdn")
	DeploymentProviderTypeHuaweiCloudELB        = DeploymentProviderType(AccessProviderTypeHuaweiCloud + "-elb")
	DeploymentProviderTypeHuaweiCloudSCM        = DeploymentProviderType(AccessProviderTypeHuaweiCloud + "-scm")
	DeploymentProviderTypeHuaweiCloudWAF        = DeploymentProviderType(AccessProviderTypeHuaweiCloud + "-waf")
	DeploymentProviderTypeJDCloudALB            = DeploymentProviderType(AccessProviderTypeJDCloud + "-alb")
	DeploymentProviderTypeJDCloudCDN            = DeploymentProviderType(AccessProviderTypeJDCloud + "-cdn")
	DeploymentProviderTypeJDCloudLive           = DeploymentProviderType(AccessProviderTypeJDCloud + "-live")
	DeploymentProviderTypeJDCloudVOD            = DeploymentProviderType(AccessProviderTypeJDCloud + "-vod")
	DeploymentProviderTypeKubernetesSecret      = DeploymentProviderType(AccessProviderTypeKubernetes + "-secret")
	DeploymentProviderTypeLocal                 = DeploymentProviderType(AccessProviderTypeLocal)
	DeploymentProviderTypeQiniuCDN              = DeploymentProviderType(AccessProviderTypeQiniu + "-cdn")
	DeploymentProviderTypeQiniuKodo             = DeploymentProviderType(AccessProviderTypeQiniu + "-kodo")
	DeploymentProviderTypeQiniuPili             = DeploymentProviderType(AccessProviderTypeQiniu + "-pili")
	DeploymentProviderTypeRainYunRCDN           = DeploymentProviderType(AccessProviderTypeRainYun + "-rcdn")
	DeploymentProviderTypeSafeLine              = DeploymentProviderType(AccessProviderTypeSafeLine)
	DeploymentProviderTypeSSH                   = DeploymentProviderType(AccessProviderTypeSSH)
	DeploymentProviderTypeTencentCloudCDN       = DeploymentProviderType(AccessProviderTypeTencentCloud + "-cdn")
	DeploymentProviderTypeTencentCloudCLB       = DeploymentProviderType(AccessProviderTypeTencentCloud + "-clb")
	DeploymentProviderTypeTencentCloudCOS       = DeploymentProviderType(AccessProviderTypeTencentCloud + "-cos")
	DeploymentProviderTypeTencentCloudCSS       = DeploymentProviderType(AccessProviderTypeTencentCloud + "-css")
	DeploymentProviderTypeTencentCloudECDN      = DeploymentProviderType(AccessProviderTypeTencentCloud + "-ecdn")
	DeploymentProviderTypeTencentCloudEO        = DeploymentProviderType(AccessProviderTypeTencentCloud + "-eo")
	DeploymentProviderTypeTencentCloudSCF       = DeploymentProviderType(AccessProviderTypeTencentCloud + "-scf")
	DeploymentProviderTypeTencentCloudSSL       = DeploymentProviderType(AccessProviderTypeTencentCloud + "-ssl")
	DeploymentProviderTypeTencentCloudSSLDeploy = DeploymentProviderType(AccessProviderTypeTencentCloud + "-ssldeploy")
	DeploymentProviderTypeTencentCloudVOD       = DeploymentProviderType(AccessProviderTypeTencentCloud + "-vod")
	DeploymentProviderTypeTencentCloudWAF       = DeploymentProviderType(AccessProviderTypeTencentCloud + "-waf")
	DeploymentProviderTypeUCloudUCDN            = DeploymentProviderType(AccessProviderTypeUCloud + "-ucdn")
	DeploymentProviderTypeUCloudUS3             = DeploymentProviderType(AccessProviderTypeUCloud + "-us3")
	DeploymentProviderTypeUpyunCDN              = DeploymentProviderType(AccessProviderTypeUpyun + "-cdn")
	DeploymentProviderTypeUpyunFile             = DeploymentProviderType(AccessProviderTypeUpyun + "-file")
	DeploymentProviderTypeVolcEngineALB         = DeploymentProviderType(AccessProviderTypeVolcEngine + "-alb")
	DeploymentProviderTypeVolcEngineCDN         = DeploymentProviderType(AccessProviderTypeVolcEngine + "-cdn")
	DeploymentProviderTypeVolcEngineCertCenter  = DeploymentProviderType(AccessProviderTypeVolcEngine + "-certcenter")
	DeploymentProviderTypeVolcEngineCLB         = DeploymentProviderType(AccessProviderTypeVolcEngine + "-clb")
	DeploymentProviderTypeVolcEngineDCDN        = DeploymentProviderType(AccessProviderTypeVolcEngine + "-dcdn")
	DeploymentProviderTypeVolcEngineImageX      = DeploymentProviderType(AccessProviderTypeVolcEngine + "-imagex")
	DeploymentProviderTypeVolcEngineLive        = DeploymentProviderType(AccessProviderTypeVolcEngine + "-live")
	DeploymentProviderTypeVolcEngineTOS         = DeploymentProviderType(AccessProviderTypeVolcEngine + "-tos")
	DeploymentProviderTypeWangsuCDNPro          = DeploymentProviderType(AccessProviderTypeWangsu + "-cdnpro")
	DeploymentProviderTypeWebhook               = DeploymentProviderType(AccessProviderTypeWebhook)
)

type NotificationProviderType string

/*
消息通知提供商常量值。
短横线前的部分始终等于授权提供商类型。

	注意：如果追加新的常量值，请保持以 ASCII 排序。
	NOTICE: If you add new constant, please keep ASCII order.
*/
const (
	NotificationProviderTypeWebhook = NotificationProviderType(AccessProviderTypeWebhook)
)
