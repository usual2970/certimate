// #region AccessProvider
/*
  注意：如果追加新的常量值，请保持以 ASCII 排序。
  NOTICE: If you add new constant, please keep ASCII order.
 */
export const ACCESS_PROVIDERS = Object.freeze({
  ["1PANEL"]: "1panel",
  ACMEHTTPREQ: "acmehttpreq",
  ALIYUN: "aliyun",
  AWS: "aws",
  AZURE: "azure",
  BAIDUCLOUD: "baiducloud",
  BAISHAN: "baishan",
  BAOTAPANEL: "baotapanel",
  BYTEPLUS: "byteplus",
  CACHEFLY: "cachefly",
  CDNFLY: "cdnfly",
  CLOUDFLARE: "cloudflare",
  CLOUDNS: "cloudns",
  CMCCCLOUD: "cmcccloud",
  DNSLA: "dnsla",
  DOGECLOUD: "dogecloud",
  DYNV6: "dynv6",
  GCORE: "gcore",
  GNAME: "gname",
  GODADDY: "godaddy",
  EDGIO: "edgio",
  HUAWEICLOUD: "huaweicloud",
  JDCLOUD: "jdcloud",
  KUBERNETES: "k8s",
  LOCAL: "local",
  NAMECHEAP: "namecheap",
  NAMEDOTCOM: "namedotcom",
  NAMESILO: "namesilo",
  NS1: "ns1",
  POWERDNS: "powerdns",
  QINIU: "qiniu",
  RAINYUN: "rainyun",
  SAFELINE: "safeline",
  SSH: "ssh",
  TENCENTCLOUD: "tencentcloud",
  UCLOUD: "ucloud",
  UPYUN: "upyun",
  VOLCENGINE: "volcengine",
  WEBHOOK: "webhook",
  WESTCN: "westcn",
} as const);

export type AccessProviderType = (typeof ACCESS_PROVIDERS)[keyof typeof ACCESS_PROVIDERS];

export const ACCESS_USAGES = Object.freeze({
  APPLY: "apply",
  DEPLOY: "deploy",
} as const);

export type AccessUsageType = (typeof ACCESS_USAGES)[keyof typeof ACCESS_USAGES];

export type AccessProvider = {
  type: AccessProviderType;
  name: string;
  icon: string;
  usages: AccessUsageType[];
};

export const accessProvidersMap: Map<AccessProvider["type"] | string, AccessProvider> = new Map(
  /*
     注意：此处的顺序决定显示在前端的顺序。
     NOTICE: The following order determines the order displayed at the frontend.
    */
  [
    [ACCESS_PROVIDERS.LOCAL, "provider.local", "/imgs/providers/local.svg", [ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.SSH, "provider.ssh", "/imgs/providers/ssh.svg", [ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.WEBHOOK, "provider.webhook", "/imgs/providers/webhook.svg", [ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.KUBERNETES, "provider.kubernetes", "/imgs/providers/kubernetes.svg", [ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.ALIYUN, "provider.aliyun", "/imgs/providers/aliyun.svg", [ACCESS_USAGES.APPLY, ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.TENCENTCLOUD, "provider.tencentcloud", "/imgs/providers/tencentcloud.svg", [ACCESS_USAGES.APPLY, ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.BAIDUCLOUD, "provider.baiducloud", "/imgs/providers/baiducloud.svg", [ACCESS_USAGES.APPLY, ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.HUAWEICLOUD, "provider.huaweicloud", "/imgs/providers/huaweicloud.svg", [ACCESS_USAGES.APPLY, ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.VOLCENGINE, "provider.volcengine", "/imgs/providers/volcengine.svg", [ACCESS_USAGES.APPLY, ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.JDCLOUD, "provider.jdcloud", "/imgs/providers/jdcloud.svg", [ACCESS_USAGES.APPLY, ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.AWS, "provider.aws", "/imgs/providers/aws.svg", [ACCESS_USAGES.APPLY, ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.AZURE, "provider.azure", "/imgs/providers/azure.svg", [ACCESS_USAGES.APPLY, ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.GCORE, "provider.gcore", "/imgs/providers/gcore.png", [ACCESS_USAGES.APPLY, ACCESS_USAGES.DEPLOY]],

    [ACCESS_PROVIDERS.QINIU, "provider.qiniu", "/imgs/providers/qiniu.svg", [ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.UPYUN, "provider.upyun", "/imgs/providers/upyun.svg", [ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.BAISHAN, "provider.baishan", "/imgs/providers/baishan.png", [ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.DOGECLOUD, "provider.dogecloud", "/imgs/providers/dogecloud.png", [ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.BYTEPLUS, "provider.byteplus", "/imgs/providers/byteplus.svg", [ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.UCLOUD, "provider.ucloud", "/imgs/providers/ucloud.svg", [ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.SAFELINE, "provider.safeline", "/imgs/providers/safeline.svg", [ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS["1PANEL"], "provider.1panel", "/imgs/providers/1panel.svg", [ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.BAOTAPANEL, "provider.baotapanel", "/imgs/providers/baotapanel.svg", [ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.CACHEFLY, "provider.cachefly", "/imgs/providers/cachefly.png", [ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.CDNFLY, "provider.cdnfly", "/imgs/providers/cdnfly.png", [ACCESS_USAGES.DEPLOY]],
    [ACCESS_PROVIDERS.EDGIO, "provider.edgio", "/imgs/providers/edgio.svg", [ACCESS_USAGES.DEPLOY]],

    [ACCESS_PROVIDERS.CLOUDFLARE, "provider.cloudflare", "/imgs/providers/cloudflare.svg", [ACCESS_USAGES.APPLY]],
    [ACCESS_PROVIDERS.CLOUDNS, "provider.cloudns", "/imgs/providers/cloudns.png", [ACCESS_USAGES.APPLY]],
    [ACCESS_PROVIDERS.DNSLA, "provider.dnsla", "/imgs/providers/dnsla.svg", [ACCESS_USAGES.APPLY]],
    [ACCESS_PROVIDERS.DYNV6, "provider.dynv6", "/imgs/providers/dynv6.svg", [ACCESS_USAGES.APPLY]],
    [ACCESS_PROVIDERS.GNAME, "provider.gname", "/imgs/providers/gname.png", [ACCESS_USAGES.APPLY]],
    [ACCESS_PROVIDERS.GODADDY, "provider.godaddy", "/imgs/providers/godaddy.svg", [ACCESS_USAGES.APPLY]],
    [ACCESS_PROVIDERS.NAMECHEAP, "provider.namecheap", "/imgs/providers/namecheap.svg", [ACCESS_USAGES.APPLY]],
    [ACCESS_PROVIDERS.NAMEDOTCOM, "provider.namedotcom", "/imgs/providers/namedotcom.svg", [ACCESS_USAGES.APPLY]],
    [ACCESS_PROVIDERS.NAMESILO, "provider.namesilo", "/imgs/providers/namesilo.svg", [ACCESS_USAGES.APPLY]],
    [ACCESS_PROVIDERS.NS1, "provider.ns1", "/imgs/providers/ns1.svg", [ACCESS_USAGES.APPLY]],
    [ACCESS_PROVIDERS.CMCCCLOUD, "provider.cmcccloud", "/imgs/providers/cmcccloud.svg", [ACCESS_USAGES.APPLY]],
    [ACCESS_PROVIDERS.RAINYUN, "provider.rainyun", "/imgs/providers/rainyun.svg", [ACCESS_USAGES.APPLY]],
    [ACCESS_PROVIDERS.WESTCN, "provider.westcn", "/imgs/providers/westcn.svg", [ACCESS_USAGES.APPLY]],
    [ACCESS_PROVIDERS.POWERDNS, "provider.powerdns", "/imgs/providers/powerdns.svg", [ACCESS_USAGES.APPLY]],
    [ACCESS_PROVIDERS.ACMEHTTPREQ, "provider.acmehttpreq", "/imgs/providers/acmehttpreq.svg", [ACCESS_USAGES.APPLY]],
  ].map((e) => [
    e[0] as string,
    {
      type: e[0] as AccessProviderType,
      name: e[1] as string,
      icon: e[2] as string,
      usages: e[3] as AccessUsageType[],
    },
  ])
);
// #endregion

// #region ApplyProvider
/*
    注意：如果追加新的常量值，请保持以 ASCII 排序。
    NOTICE: If you add new constant, please keep ASCII order.
   */
export const APPLY_DNS_PROVIDERS = Object.freeze({
  ACMEHTTPREQ: `${ACCESS_PROVIDERS.ACMEHTTPREQ}`,
  ALIYUN: `${ACCESS_PROVIDERS.ALIYUN}`, // 兼容旧值，等同于 `ALIYUN_DNS`
  ALIYUN_DNS: `${ACCESS_PROVIDERS.ALIYUN}-dns`,
  AWS: `${ACCESS_PROVIDERS.AWS}`, // 兼容旧值，等同于 `AWS_ROUTE53`
  AWS_ROUTE53: `${ACCESS_PROVIDERS.AWS}-route53`,
  AZURE: `${ACCESS_PROVIDERS.AZURE}`, // 兼容旧值，等同于 `AZURE_DNS`
  AZURE_DNS: `${ACCESS_PROVIDERS.AZURE}-dns`,
  BAIDUCLOUD: `${ACCESS_PROVIDERS.BAIDUCLOUD}`, // 兼容旧值，等同于 `BAIDUCLOUD_DNS`
  BAIDUCLOUD_DNS: `${ACCESS_PROVIDERS.BAIDUCLOUD}-dns`,
  CLOUDFLARE: `${ACCESS_PROVIDERS.CLOUDFLARE}`,
  CLOUDNS: `${ACCESS_PROVIDERS.CLOUDNS}`,
  CMCCCLOUD: `${ACCESS_PROVIDERS.CMCCCLOUD}`,
  DNSLA: `${ACCESS_PROVIDERS.DNSLA}`,
  DYNV6: `${ACCESS_PROVIDERS.DYNV6}`,
  GCORE: `${ACCESS_PROVIDERS.GCORE}`,
  GNAME: `${ACCESS_PROVIDERS.GNAME}`,
  GODADDY: `${ACCESS_PROVIDERS.GODADDY}`,
  HUAWEICLOUD: `${ACCESS_PROVIDERS.HUAWEICLOUD}`, // 兼容旧值，等同于 `HUAWEICLOUD_DNS`
  HUAWEICLOUD_DNS: `${ACCESS_PROVIDERS.HUAWEICLOUD}-dns`,
  JDCLOUD: `${ACCESS_PROVIDERS.JDCLOUD}`, // 兼容旧值，等同于 `JDCLOUD_DNS`
  JDCLOUD_DNS: `${ACCESS_PROVIDERS.JDCLOUD}-dns`,
  NAMECHEAP: `${ACCESS_PROVIDERS.NAMECHEAP}`,
  NAMEDOTCOM: `${ACCESS_PROVIDERS.NAMEDOTCOM}`,
  NAMESILO: `${ACCESS_PROVIDERS.NAMESILO}`,
  NS1: `${ACCESS_PROVIDERS.NS1}`,
  POWERDNS: `${ACCESS_PROVIDERS.POWERDNS}`,
  RAINYUN: `${ACCESS_PROVIDERS.RAINYUN}`,
  TENCENTCLOUD: `${ACCESS_PROVIDERS.TENCENTCLOUD}`, // 兼容旧值，等同于 `TENCENTCLOUD_DNS`
  TENCENTCLOUD_DNS: `${ACCESS_PROVIDERS.TENCENTCLOUD}-dns`,
  VOLCENGINE: `${ACCESS_PROVIDERS.VOLCENGINE}`, // 兼容旧值，等同于 `VOLCENGINE_DNS`
  VOLCENGINE_DNS: `${ACCESS_PROVIDERS.VOLCENGINE}-dns`,
  WESTCN: `${ACCESS_PROVIDERS.WESTCN}`,
} as const);

export type ApplyDNSProviderType = (typeof APPLY_DNS_PROVIDERS)[keyof typeof APPLY_DNS_PROVIDERS];

export type ApplyDNSProvider = {
  type: ApplyDNSProviderType;
  name: string;
  icon: string;
  provider: AccessProviderType;
};

export const applyDNSProvidersMap: Map<ApplyDNSProvider["type"] | string, ApplyDNSProvider> = new Map(
  /*
     注意：此处的顺序决定显示在前端的顺序。
     NOTICE: The following order determines the order displayed at the frontend.
    */
  [
    [APPLY_DNS_PROVIDERS.ALIYUN_DNS, "provider.aliyun.dns"],
    [APPLY_DNS_PROVIDERS.TENCENTCLOUD_DNS, "provider.tencentcloud.dns"],
    [APPLY_DNS_PROVIDERS.BAIDUCLOUD_DNS, "provider.baiducloud.dns"],
    [APPLY_DNS_PROVIDERS.HUAWEICLOUD_DNS, "provider.huaweicloud.dns"],
    [APPLY_DNS_PROVIDERS.VOLCENGINE_DNS, "provider.volcengine.dns"],
    [APPLY_DNS_PROVIDERS.JDCLOUD_DNS, "provider.jdcloud.dns"],
    [APPLY_DNS_PROVIDERS.AWS_ROUTE53, "provider.aws.route53"],
    [APPLY_DNS_PROVIDERS.AZURE_DNS, "provider.azure.dns"],
    [APPLY_DNS_PROVIDERS.CLOUDFLARE, "provider.cloudflare"],
    [APPLY_DNS_PROVIDERS.CLOUDNS, "provider.cloudns"],
    [APPLY_DNS_PROVIDERS.DNSLA, "provider.dnsla"],
    [APPLY_DNS_PROVIDERS.DYNV6, "provider.dynv6"],
    [APPLY_DNS_PROVIDERS.GCORE, "provider.gcore"],
    [APPLY_DNS_PROVIDERS.GNAME, "provider.gname"],
    [APPLY_DNS_PROVIDERS.GODADDY, "provider.godaddy"],
    [APPLY_DNS_PROVIDERS.NAMECHEAP, "provider.namecheap"],
    [APPLY_DNS_PROVIDERS.NAMEDOTCOM, "provider.namedotcom"],
    [APPLY_DNS_PROVIDERS.NAMESILO, "provider.namesilo"],
    [APPLY_DNS_PROVIDERS.NS1, "provider.ns1"],
    [APPLY_DNS_PROVIDERS.CMCCCLOUD, "provider.cmcc"],
    [APPLY_DNS_PROVIDERS.RAINYUN, "provider.rainyun"],
    [APPLY_DNS_PROVIDERS.WESTCN, "provider.westcn"],
    [APPLY_DNS_PROVIDERS.POWERDNS, "provider.powerdns"],
    [APPLY_DNS_PROVIDERS.ACMEHTTPREQ, "provider.acmehttpreq"],
  ].map(([type, name]) => [
    type,
    {
      type: type as ApplyDNSProviderType,
      name: name,
      icon: accessProvidersMap.get(type.split("-")[0])!.icon,
      provider: type.split("-")[0] as AccessProviderType,
    },
  ])
);
// #endregion

// #region DeployProvider
/*
    注意：如果追加新的常量值，请保持以 ASCII 排序。
    NOTICE: If you add new constant, please keep ASCII order.
   */
export const DEPLOY_PROVIDERS = Object.freeze({
  ["1PANEL_CONSOLE"]: `${ACCESS_PROVIDERS["1PANEL"]}-console`,
  ["1PANEL_SITE"]: `${ACCESS_PROVIDERS["1PANEL"]}-site`,
  ALIYUN_ALB: `${ACCESS_PROVIDERS.ALIYUN}-alb`,
  ALIYUN_CAS: `${ACCESS_PROVIDERS.ALIYUN}-cas`,
  ALIYUN_CAS_DEPLOY: `${ACCESS_PROVIDERS.ALIYUN}-casdeploy`,
  ALIYUN_CDN: `${ACCESS_PROVIDERS.ALIYUN}-cdn`,
  ALIYUN_CLB: `${ACCESS_PROVIDERS.ALIYUN}-clb`,
  ALIYUN_DCDN: `${ACCESS_PROVIDERS.ALIYUN}-dcdn`,
  ALIYUN_ESA: `${ACCESS_PROVIDERS.ALIYUN}-esa`,
  ALIYUN_FC: `${ACCESS_PROVIDERS.ALIYUN}-fc`,
  ALIYUN_LIVE: `${ACCESS_PROVIDERS.ALIYUN}-live`,
  ALIYUN_NLB: `${ACCESS_PROVIDERS.ALIYUN}-nlb`,
  ALIYUN_OSS: `${ACCESS_PROVIDERS.ALIYUN}-oss`,
  ALIYUN_VOD: `${ACCESS_PROVIDERS.ALIYUN}-vod`,
  ALIYUN_WAF: `${ACCESS_PROVIDERS.ALIYUN}-waf`,
  AWS_ACM: `${ACCESS_PROVIDERS.AWS}-acm`,
  AWS_CLOUDFRONT: `${ACCESS_PROVIDERS.AWS}-cloudfront`,
  AZURE_KEYVAULT: `${ACCESS_PROVIDERS.AZURE}-keyvault`,
  BAIDUCLOUD_CDN: `${ACCESS_PROVIDERS.BAIDUCLOUD}-cdn`,
  BAISHAN_CDN: `${ACCESS_PROVIDERS.BAISHAN}-cdn`,
  BAOTAPANEL_CONSOLE: `${ACCESS_PROVIDERS.BAOTAPANEL}-console`,
  BAOTAPANEL_SITE: `${ACCESS_PROVIDERS.BAOTAPANEL}-site`,
  BYTEPLUS_CDN: `${ACCESS_PROVIDERS.BYTEPLUS}-cdn`,
  CACHEFLY: `${ACCESS_PROVIDERS.CACHEFLY}`,
  CDNFLY: `${ACCESS_PROVIDERS.CDNFLY}`,
  DOGECLOUD_CDN: `${ACCESS_PROVIDERS.DOGECLOUD}-cdn`,
  EDGIO_APPLICATIONS: `${ACCESS_PROVIDERS.EDGIO}-applications`,
  GCORE_CDN: `${ACCESS_PROVIDERS.GCORE}-cdn`,
  HUAWEICLOUD_CDN: `${ACCESS_PROVIDERS.HUAWEICLOUD}-cdn`,
  HUAWEICLOUD_ELB: `${ACCESS_PROVIDERS.HUAWEICLOUD}-elb`,
  HUAWEICLOUD_WAF: `${ACCESS_PROVIDERS.HUAWEICLOUD}-waf`,
  JDCLOUD_ALB: `${ACCESS_PROVIDERS.JDCLOUD}-alb`,
  JDCLOUD_CDN: `${ACCESS_PROVIDERS.JDCLOUD}-cdn`,
  JDCLOUD_LIVE: `${ACCESS_PROVIDERS.JDCLOUD}-live`,
  JDCLOUD_VOD: `${ACCESS_PROVIDERS.JDCLOUD}-vod`,
  KUBERNETES_SECRET: `${ACCESS_PROVIDERS.KUBERNETES}-secret`,
  LOCAL: `${ACCESS_PROVIDERS.LOCAL}`,
  QINIU_CDN: `${ACCESS_PROVIDERS.QINIU}-cdn`,
  QINIU_KODO: `${ACCESS_PROVIDERS.QINIU}-kodo`,
  QINIU_PILI: `${ACCESS_PROVIDERS.QINIU}-pili`,
  SAFELINE: `${ACCESS_PROVIDERS.SAFELINE}`,
  SSH: `${ACCESS_PROVIDERS.SSH}`,
  TENCENTCLOUD_CDN: `${ACCESS_PROVIDERS.TENCENTCLOUD}-cdn`,
  TENCENTCLOUD_CLB: `${ACCESS_PROVIDERS.TENCENTCLOUD}-clb`,
  TENCENTCLOUD_COS: `${ACCESS_PROVIDERS.TENCENTCLOUD}-cos`,
  TENCENTCLOUD_CSS: `${ACCESS_PROVIDERS.TENCENTCLOUD}-css`,
  TENCENTCLOUD_ECDN: `${ACCESS_PROVIDERS.TENCENTCLOUD}-ecdn`,
  TENCENTCLOUD_EO: `${ACCESS_PROVIDERS.TENCENTCLOUD}-eo`,
  TENCENTCLOUD_SCF: `${ACCESS_PROVIDERS.TENCENTCLOUD}-scf`,
  TENCENTCLOUD_SSL: `${ACCESS_PROVIDERS.TENCENTCLOUD}-ssl`,
  TENCENTCLOUD_SSL_DEPLOY: `${ACCESS_PROVIDERS.TENCENTCLOUD}-ssldeploy`,
  TENCENTCLOUD_VOD: `${ACCESS_PROVIDERS.TENCENTCLOUD}-vod`,
  TENCENTCLOUD_WAF: `${ACCESS_PROVIDERS.TENCENTCLOUD}-waf`,
  UCLOUD_UCDN: `${ACCESS_PROVIDERS.UCLOUD}-ucdn`,
  UCLOUD_US3: `${ACCESS_PROVIDERS.UCLOUD}-us3`,
  UPYUN_CDN: `${ACCESS_PROVIDERS.UPYUN}-cdn`,
  UPYUN_FILE: `${ACCESS_PROVIDERS.UPYUN}-file`,
  VOLCENGINE_CDN: `${ACCESS_PROVIDERS.VOLCENGINE}-cdn`,
  VOLCENGINE_CLB: `${ACCESS_PROVIDERS.VOLCENGINE}-clb`,
  VOLCENGINE_DCDN: `${ACCESS_PROVIDERS.VOLCENGINE}-dcdn`,
  VOLCENGINE_IMAGEX: `${ACCESS_PROVIDERS.VOLCENGINE}-imagex`,
  VOLCENGINE_LIVE: `${ACCESS_PROVIDERS.VOLCENGINE}-live`,
  VOLCENGINE_TOS: `${ACCESS_PROVIDERS.VOLCENGINE}-tos`,
  WEBHOOK: `${ACCESS_PROVIDERS.WEBHOOK}`,
} as const);

export type DeployProviderType = (typeof DEPLOY_PROVIDERS)[keyof typeof DEPLOY_PROVIDERS];

export const DEPLOY_CATEGORIES = Object.freeze({
  ALL: "all",
  CDN: "cdn",
  STORAGE: "storage",
  LOADBALANCE: "loadbalance",
  FIREWALL: "firewall",
  AV: "av",
  SERVERLESS: "serverless",
  WEBSITE: "website",
  OTHER: "other",
} as const);

export type DeployCategoryType = (typeof DEPLOY_CATEGORIES)[keyof typeof DEPLOY_CATEGORIES];

export type DeployProvider = {
  type: DeployProviderType;
  name: string;
  icon: string;
  provider: AccessProviderType;
  category: DeployCategoryType;
};

export const deployProvidersMap: Map<DeployProvider["type"] | string, DeployProvider> = new Map(
  /*
     注意：此处的顺序决定显示在前端的顺序。
     NOTICE: The following order determines the order displayed at the frontend.
    */
  [
    [DEPLOY_PROVIDERS.LOCAL, "provider.local", DEPLOY_CATEGORIES.OTHER],
    [DEPLOY_PROVIDERS.SSH, "provider.ssh", DEPLOY_CATEGORIES.OTHER],
    [DEPLOY_PROVIDERS.WEBHOOK, "provider.webhook", DEPLOY_CATEGORIES.OTHER],
    [DEPLOY_PROVIDERS.KUBERNETES_SECRET, "provider.kubernetes.secret", DEPLOY_CATEGORIES.OTHER],
    [DEPLOY_PROVIDERS.ALIYUN_OSS, "provider.aliyun.oss", DEPLOY_CATEGORIES.STORAGE],
    [DEPLOY_PROVIDERS.ALIYUN_CDN, "provider.aliyun.cdn", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.ALIYUN_DCDN, "provider.aliyun.dcdn", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.ALIYUN_ESA, "provider.aliyun.esa", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.ALIYUN_CLB, "provider.aliyun.clb", DEPLOY_CATEGORIES.LOADBALANCE],
    [DEPLOY_PROVIDERS.ALIYUN_ALB, "provider.aliyun.alb", DEPLOY_CATEGORIES.LOADBALANCE],
    [DEPLOY_PROVIDERS.ALIYUN_NLB, "provider.aliyun.nlb", DEPLOY_CATEGORIES.LOADBALANCE],
    [DEPLOY_PROVIDERS.ALIYUN_WAF, "provider.aliyun.waf", DEPLOY_CATEGORIES.FIREWALL],
    [DEPLOY_PROVIDERS.ALIYUN_LIVE, "provider.aliyun.live", DEPLOY_CATEGORIES.AV],
    [DEPLOY_PROVIDERS.ALIYUN_VOD, "provider.aliyun.vod", DEPLOY_CATEGORIES.AV],
    [DEPLOY_PROVIDERS.ALIYUN_FC, "provider.aliyun.fc", DEPLOY_CATEGORIES.SERVERLESS],
    [DEPLOY_PROVIDERS.ALIYUN_CAS, "provider.aliyun.cas", DEPLOY_CATEGORIES.OTHER],
    [DEPLOY_PROVIDERS.ALIYUN_CAS_DEPLOY, "provider.aliyun.cas_deploy", DEPLOY_CATEGORIES.OTHER],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_COS, "provider.tencentcloud.cos", DEPLOY_CATEGORIES.STORAGE],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_CDN, "provider.tencentcloud.cdn", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_ECDN, "provider.tencentcloud.ecdn", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_EO, "provider.tencentcloud.eo", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_CLB, "provider.tencentcloud.clb", DEPLOY_CATEGORIES.LOADBALANCE],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_WAF, "provider.tencentcloud.waf", DEPLOY_CATEGORIES.FIREWALL],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_CSS, "provider.tencentcloud.css", DEPLOY_CATEGORIES.AV],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_VOD, "provider.tencentcloud.vod", DEPLOY_CATEGORIES.AV],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_SCF, "provider.tencentcloud.scf", DEPLOY_CATEGORIES.SERVERLESS],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_SSL, "provider.tencentcloud.ssl", DEPLOY_CATEGORIES.OTHER],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_SSL_DEPLOY, "provider.tencentcloud.ssl_deploy", DEPLOY_CATEGORIES.OTHER],
    [DEPLOY_PROVIDERS.HUAWEICLOUD_CDN, "provider.huaweicloud.cdn", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.HUAWEICLOUD_ELB, "provider.huaweicloud.elb", DEPLOY_CATEGORIES.LOADBALANCE],
    [DEPLOY_PROVIDERS.HUAWEICLOUD_WAF, "provider.huaweicloud.waf", DEPLOY_CATEGORIES.FIREWALL],
    [DEPLOY_PROVIDERS.BAIDUCLOUD_CDN, "provider.baiducloud.cdn", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.VOLCENGINE_TOS, "provider.volcengine.tos", DEPLOY_CATEGORIES.STORAGE],
    [DEPLOY_PROVIDERS.VOLCENGINE_CDN, "provider.volcengine.cdn", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.VOLCENGINE_DCDN, "provider.volcengine.dcdn", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.VOLCENGINE_CLB, "provider.volcengine.clb", DEPLOY_CATEGORIES.LOADBALANCE],
    [DEPLOY_PROVIDERS.VOLCENGINE_IMAGEX, "provider.volcengine.imagex", DEPLOY_CATEGORIES.STORAGE],
    [DEPLOY_PROVIDERS.VOLCENGINE_LIVE, "provider.volcengine.live", DEPLOY_CATEGORIES.AV],
    [DEPLOY_PROVIDERS.JDCLOUD_ALB, "provider.jdcloud.alb", DEPLOY_CATEGORIES.LOADBALANCE],
    [DEPLOY_PROVIDERS.JDCLOUD_CDN, "provider.jdcloud.cdn", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.JDCLOUD_LIVE, "provider.jdcloud.live", DEPLOY_CATEGORIES.AV],
    [DEPLOY_PROVIDERS.JDCLOUD_VOD, "provider.jdcloud.vod", DEPLOY_CATEGORIES.AV],
    [DEPLOY_PROVIDERS.QINIU_KODO, "provider.qiniu.kodo", DEPLOY_CATEGORIES.STORAGE],
    [DEPLOY_PROVIDERS.QINIU_CDN, "provider.qiniu.cdn", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.QINIU_PILI, "provider.qiniu.pili", DEPLOY_CATEGORIES.AV],
    [DEPLOY_PROVIDERS.UPYUN_FILE, "provider.upyun.file", DEPLOY_CATEGORIES.STORAGE],
    [DEPLOY_PROVIDERS.UPYUN_CDN, "provider.upyun.cdn", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.BAISHAN_CDN, "provider.baishan.cdn", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.DOGECLOUD_CDN, "provider.dogecloud.cdn", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.BYTEPLUS_CDN, "provider.byteplus.cdn", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.UCLOUD_US3, "provider.ucloud.us3", DEPLOY_CATEGORIES.STORAGE],
    [DEPLOY_PROVIDERS.UCLOUD_UCDN, "provider.ucloud.ucdn", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.AWS_CLOUDFRONT, "provider.aws.cloudfront", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.AWS_ACM, "provider.aws.acm", DEPLOY_CATEGORIES.OTHER],
    [DEPLOY_PROVIDERS.AZURE_KEYVAULT, "provider.azure.keyvault", DEPLOY_CATEGORIES.OTHER],
    [DEPLOY_PROVIDERS.CACHEFLY, "provider.cachefly", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.CDNFLY, "provider.cdnfly", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS.EDGIO_APPLICATIONS, "provider.edgio.applications", DEPLOY_CATEGORIES.WEBSITE],
    [DEPLOY_PROVIDERS.GCORE_CDN, "provider.gcore.cdn", DEPLOY_CATEGORIES.CDN],
    [DEPLOY_PROVIDERS["1PANEL_SITE"], "provider.1panel.site", DEPLOY_CATEGORIES.WEBSITE],
    [DEPLOY_PROVIDERS["1PANEL_CONSOLE"], "provider.1panel.console", DEPLOY_CATEGORIES.OTHER],
    [DEPLOY_PROVIDERS.BAOTAPANEL_SITE, "provider.baotapanel.site", DEPLOY_CATEGORIES.WEBSITE],
    [DEPLOY_PROVIDERS.BAOTAPANEL_CONSOLE, "provider.baotapanel.console", DEPLOY_CATEGORIES.OTHER],
    [DEPLOY_PROVIDERS.SAFELINE, "provider.safeline", DEPLOY_CATEGORIES.FIREWALL],
  ].map(([type, name, category]) => [
    type,
    {
      type: type as DeployProviderType,
      name: name,
      icon: accessProvidersMap.get(type.split("-")[0])!.icon,
      provider: type.split("-")[0] as AccessProviderType,
      category: category as DeployCategoryType,
    },
  ])
);
// #endregion
