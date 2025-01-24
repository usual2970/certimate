// #region AccessProvider
/*
  注意：如果追加新的常量值，请保持以 ASCII 排序。
  NOTICE: If you add new constant, please keep ASCII order.
 */
export const ACCESS_PROVIDERS = Object.freeze({
  ACMEHTTPREQ: "acmehttpreq",
  ALIYUN: "aliyun",
  AWS: "aws",
  AZURE: "azure",
  BAIDUCLOUD: "baiducloud",
  BYTEPLUS: "byteplus",
  CLOUDFLARE: "cloudflare",
  CLOUDNS: "cloudns",
  DOGECLOUD: "dogecloud",
  GNAME: "gname",
  GODADDY: "godaddy",
  EDGIO: "edgio",
  HUAWEICLOUD: "huaweicloud",
  KUBERNETES: "k8s",
  LOCAL: "local",
  NAMEDOTCOM: "namedotcom",
  NAMESILO: "namesilo",
  NS1: "ns1",
  POWERDNS: "powerdns",
  QINIU: "qiniu",
  RAINYUN: "rainyun",
  SSH: "ssh",
  TENCENTCLOUD: "tencentcloud",
  UCLOUD: "ucloud",
  VOLCENGINE: "volcengine",
  WEBHOOK: "webhook",
  WESTCN: "westcn",
} as const);

export type AccessProviderType = (typeof ACCESS_PROVIDERS)[keyof typeof ACCESS_PROVIDERS];

export const ACCESS_USAGES = Object.freeze({
  ALL: "all",
  APPLY: "apply",
  DEPLOY: "deploy",
} as const);

export type AccessUsageType = (typeof ACCESS_USAGES)[keyof typeof ACCESS_USAGES];

export type AccessProvider = {
  type: AccessProviderType;
  name: string;
  icon: string;
  usage: AccessUsageType;
};

export const accessProvidersMap: Map<AccessProvider["type"] | string, AccessProvider> = new Map(
  /*
   注意：此处的顺序决定显示在前端的顺序。
   NOTICE: The following order determines the order displayed at the frontend.
  */
  [
    [ACCESS_PROVIDERS.LOCAL, "common.provider.local", "/imgs/providers/local.svg", ACCESS_USAGES.DEPLOY],
    [ACCESS_PROVIDERS.SSH, "common.provider.ssh", "/imgs/providers/ssh.svg", ACCESS_USAGES.DEPLOY],
    [ACCESS_PROVIDERS.WEBHOOK, "common.provider.webhook", "/imgs/providers/webhook.svg", ACCESS_USAGES.DEPLOY],
    [ACCESS_PROVIDERS.KUBERNETES, "common.provider.kubernetes", "/imgs/providers/kubernetes.svg", ACCESS_USAGES.DEPLOY],
    [ACCESS_PROVIDERS.ALIYUN, "common.provider.aliyun", "/imgs/providers/aliyun.svg", ACCESS_USAGES.ALL],
    [ACCESS_PROVIDERS.TENCENTCLOUD, "common.provider.tencentcloud", "/imgs/providers/tencentcloud.svg", ACCESS_USAGES.ALL],
    [ACCESS_PROVIDERS.HUAWEICLOUD, "common.provider.huaweicloud", "/imgs/providers/huaweicloud.svg", ACCESS_USAGES.ALL],
    [ACCESS_PROVIDERS.VOLCENGINE, "common.provider.volcengine", "/imgs/providers/volcengine.svg", ACCESS_USAGES.ALL],
    [ACCESS_PROVIDERS.AWS, "common.provider.aws", "/imgs/providers/aws.svg", ACCESS_USAGES.ALL],
    [ACCESS_PROVIDERS.BAIDUCLOUD, "common.provider.baiducloud", "/imgs/providers/baiducloud.svg", ACCESS_USAGES.DEPLOY],
    [ACCESS_PROVIDERS.QINIU, "common.provider.qiniu", "/imgs/providers/qiniu.svg", ACCESS_USAGES.DEPLOY],
    [ACCESS_PROVIDERS.DOGECLOUD, "common.provider.dogecloud", "/imgs/providers/dogecloud.svg", ACCESS_USAGES.DEPLOY],
    [ACCESS_PROVIDERS.BYTEPLUS, "common.provider.byteplus", "/imgs/providers/byteplus.svg", ACCESS_USAGES.DEPLOY],
    [ACCESS_PROVIDERS.UCLOUD, "common.provider.ucloud", "/imgs/providers/ucloud.svg", ACCESS_USAGES.DEPLOY],
    [ACCESS_PROVIDERS.EDGIO, "common.provider.edgio", "/imgs/providers/edgio.svg", ACCESS_USAGES.DEPLOY],
    [ACCESS_PROVIDERS.AZURE, "common.provider.azure", "/imgs/providers/azure.svg", ACCESS_USAGES.APPLY],
    [ACCESS_PROVIDERS.CLOUDFLARE, "common.provider.cloudflare", "/imgs/providers/cloudflare.svg", ACCESS_USAGES.APPLY],
    [ACCESS_PROVIDERS.CLOUDNS, "common.provider.cloudns", "/imgs/providers/cloudns.svg", ACCESS_USAGES.APPLY],
    [ACCESS_PROVIDERS.GNAME, "common.provider.gname", "/imgs/providers/gname.svg", ACCESS_USAGES.APPLY],
    [ACCESS_PROVIDERS.GODADDY, "common.provider.godaddy", "/imgs/providers/godaddy.svg", ACCESS_USAGES.APPLY],
    [ACCESS_PROVIDERS.NAMEDOTCOM, "common.provider.namedotcom", "/imgs/providers/namedotcom.svg", ACCESS_USAGES.APPLY],
    [ACCESS_PROVIDERS.NAMESILO, "common.provider.namesilo", "/imgs/providers/namesilo.svg", ACCESS_USAGES.APPLY],
    [ACCESS_PROVIDERS.NS1, "common.provider.ns1", "/imgs/providers/ns1.svg", ACCESS_USAGES.APPLY],
    [ACCESS_PROVIDERS.RAINYUN, "common.provider.rainyun", "/imgs/providers/rainyun.svg", ACCESS_USAGES.APPLY],
    [ACCESS_PROVIDERS.WESTCN, "common.provider.westcn", "/imgs/providers/westcn.svg", ACCESS_USAGES.APPLY],
    [ACCESS_PROVIDERS.POWERDNS, "common.provider.powerdns", "/imgs/providers/powerdns.svg", ACCESS_USAGES.APPLY],
    [ACCESS_PROVIDERS.ACMEHTTPREQ, "common.provider.acmehttpreq", "/imgs/providers/acmehttpreq.svg", ACCESS_USAGES.APPLY],
  ].map(([type, name, icon, usage]) => [
    type,
    {
      type: type as AccessProviderType,
      name: name,
      icon: icon,
      usage: usage as AccessUsageType,
    },
  ])
);
// #endregion

// #region DNSProvider
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
  AZURE_DNS: `${ACCESS_PROVIDERS.AZURE}-dns`,
  CLOUDFLARE: `${ACCESS_PROVIDERS.CLOUDFLARE}`,
  CLOUDNS: `${ACCESS_PROVIDERS.CLOUDNS}`,
  GNAME: `${ACCESS_PROVIDERS.GNAME}`,
  GODADDY: `${ACCESS_PROVIDERS.GODADDY}`,
  HUAWEICLOUD: `${ACCESS_PROVIDERS.HUAWEICLOUD}`, // 兼容旧值，等同于 `HUAWEICLOUD_DNS`
  HUAWEICLOUD_DNS: `${ACCESS_PROVIDERS.HUAWEICLOUD}-dns`,
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
    [APPLY_DNS_PROVIDERS.ALIYUN_DNS, "common.provider.aliyun.dns"],
    [APPLY_DNS_PROVIDERS.TENCENTCLOUD_DNS, "common.provider.tencentcloud.dns"],
    [APPLY_DNS_PROVIDERS.HUAWEICLOUD_DNS, "common.provider.huaweicloud.dns"],
    [APPLY_DNS_PROVIDERS.VOLCENGINE_DNS, "common.provider.volcengine.dns"],
    [APPLY_DNS_PROVIDERS.AWS_ROUTE53, "common.provider.aws.route53"],
    [APPLY_DNS_PROVIDERS.AZURE_DNS, "common.provider.azure.dns"],
    [APPLY_DNS_PROVIDERS.CLOUDFLARE, "common.provider.cloudflare"],
    [APPLY_DNS_PROVIDERS.CLOUDNS, "common.provider.cloudns"],
    [APPLY_DNS_PROVIDERS.GNAME, "common.provider.gname"],
    [APPLY_DNS_PROVIDERS.GODADDY, "common.provider.godaddy"],
    [APPLY_DNS_PROVIDERS.NAMEDOTCOM, "common.provider.namedotcom"],
    [APPLY_DNS_PROVIDERS.NAMESILO, "common.provider.namesilo"],
    [APPLY_DNS_PROVIDERS.NS1, "common.provider.ns1"],
    [APPLY_DNS_PROVIDERS.RAINYUN, "common.provider.rainyun"],
    [APPLY_DNS_PROVIDERS.WESTCN, "common.provider.westcn"],
    [APPLY_DNS_PROVIDERS.POWERDNS, "common.provider.powerdns"],
    [APPLY_DNS_PROVIDERS.ACMEHTTPREQ, "common.provider.acmehttpreq"],
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
  ALIYUN_ALB: `${ACCESS_PROVIDERS.ALIYUN}-alb`,
  ALIYUN_CDN: `${ACCESS_PROVIDERS.ALIYUN}-cdn`,
  ALIYUN_CLB: `${ACCESS_PROVIDERS.ALIYUN}-clb`,
  ALIYUN_DCDN: `${ACCESS_PROVIDERS.ALIYUN}-dcdn`,
  ALIYUN_LIVE: `${ACCESS_PROVIDERS.ALIYUN}-live`,
  ALIYUN_NLB: `${ACCESS_PROVIDERS.ALIYUN}-nlb`,
  ALIYUN_OSS: `${ACCESS_PROVIDERS.ALIYUN}-oss`,
  ALIYUN_WAF: `${ACCESS_PROVIDERS.ALIYUN}-waf`,
  AWS_CLOUDFRONT: `${ACCESS_PROVIDERS.AWS}-cloudfront`,
  BAIDUCLOUD_CDN: `${ACCESS_PROVIDERS.BAIDUCLOUD}-cdn`,
  BYTEPLUS_CDN: `${ACCESS_PROVIDERS.BYTEPLUS}-cdn`,
  DOGECLOUD_CDN: `${ACCESS_PROVIDERS.DOGECLOUD}-cdn`,
  EDGIO_APPLICATIONS: `${ACCESS_PROVIDERS.EDGIO}-applications`,
  HUAWEICLOUD_CDN: `${ACCESS_PROVIDERS.HUAWEICLOUD}-cdn`,
  HUAWEICLOUD_ELB: `${ACCESS_PROVIDERS.HUAWEICLOUD}-elb`,
  KUBERNETES_SECRET: `${ACCESS_PROVIDERS.KUBERNETES}-secret`,
  LOCAL: `${ACCESS_PROVIDERS.LOCAL}`,
  QINIU_CDN: `${ACCESS_PROVIDERS.QINIU}-cdn`,
  QINIU_PILI: `${ACCESS_PROVIDERS.QINIU}-pili`,
  SSH: `${ACCESS_PROVIDERS.SSH}`,
  TENCENTCLOUD_CDN: `${ACCESS_PROVIDERS.TENCENTCLOUD}-cdn`,
  TENCENTCLOUD_CLB: `${ACCESS_PROVIDERS.TENCENTCLOUD}-clb`,
  TENCENTCLOUD_COS: `${ACCESS_PROVIDERS.TENCENTCLOUD}-cos`,
  TENCENTCLOUD_CSS: `${ACCESS_PROVIDERS.TENCENTCLOUD}-css`,
  TENCENTCLOUD_ECDN: `${ACCESS_PROVIDERS.TENCENTCLOUD}-ecdn`,
  TENCENTCLOUD_EO: `${ACCESS_PROVIDERS.TENCENTCLOUD}-eo`,
  UCLOUD_UCDN: `${ACCESS_PROVIDERS.UCLOUD}-ucdn`,
  UCLOUD_US3: `${ACCESS_PROVIDERS.UCLOUD}-us3`,
  VOLCENGINE_CDN: `${ACCESS_PROVIDERS.VOLCENGINE}-cdn`,
  VOLCENGINE_CLB: `${ACCESS_PROVIDERS.VOLCENGINE}-clb`,
  VOLCENGINE_DCDN: `${ACCESS_PROVIDERS.VOLCENGINE}-dcdn`,
  VOLCENGINE_LIVE: `${ACCESS_PROVIDERS.VOLCENGINE}-live`,
  VOLCENGINE_TOS: `${ACCESS_PROVIDERS.VOLCENGINE}-tos`,
  WEBHOOK: `${ACCESS_PROVIDERS.WEBHOOK}`,
} as const);

export type DeployProviderType = (typeof DEPLOY_PROVIDERS)[keyof typeof DEPLOY_PROVIDERS];

export type DeployProvider = {
  type: DeployProviderType;
  name: string;
  icon: string;
  provider: AccessProviderType;
};

export const deployProvidersMap: Map<DeployProvider["type"] | string, DeployProvider> = new Map(
  /*
   注意：此处的顺序决定显示在前端的顺序。
   NOTICE: The following order determines the order displayed at the frontend.
  */
  [
    [DEPLOY_PROVIDERS.LOCAL, "common.provider.local"],
    [DEPLOY_PROVIDERS.SSH, "common.provider.ssh"],
    [DEPLOY_PROVIDERS.WEBHOOK, "common.provider.webhook"],
    [DEPLOY_PROVIDERS.KUBERNETES_SECRET, "common.provider.kubernetes.secret"],
    [DEPLOY_PROVIDERS.ALIYUN_OSS, "common.provider.aliyun.oss"],
    [DEPLOY_PROVIDERS.ALIYUN_CDN, "common.provider.aliyun.cdn"],
    [DEPLOY_PROVIDERS.ALIYUN_DCDN, "common.provider.aliyun.dcdn"],
    [DEPLOY_PROVIDERS.ALIYUN_CLB, "common.provider.aliyun.clb"],
    [DEPLOY_PROVIDERS.ALIYUN_ALB, "common.provider.aliyun.alb"],
    [DEPLOY_PROVIDERS.ALIYUN_NLB, "common.provider.aliyun.nlb"],
    [DEPLOY_PROVIDERS.ALIYUN_WAF, "common.provider.aliyun.waf"],
    [DEPLOY_PROVIDERS.ALIYUN_LIVE, "common.provider.aliyun.live"],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_COS, "common.provider.tencentcloud.cos"],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_CDN, "common.provider.tencentcloud.cdn"],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_ECDN, "common.provider.tencentcloud.ecdn"],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_EO, "common.provider.tencentcloud.eo"],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_CLB, "common.provider.tencentcloud.clb"],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_CSS, "common.provider.tencentcloud.css"],
    [DEPLOY_PROVIDERS.HUAWEICLOUD_CDN, "common.provider.huaweicloud.cdn"],
    [DEPLOY_PROVIDERS.HUAWEICLOUD_ELB, "common.provider.huaweicloud.elb"],
    [DEPLOY_PROVIDERS.BAIDUCLOUD_CDN, "common.provider.baiducloud.cdn"],
    [DEPLOY_PROVIDERS.VOLCENGINE_TOS, "common.provider.volcengine.tos"],
    [DEPLOY_PROVIDERS.VOLCENGINE_CDN, "common.provider.volcengine.cdn"],
    [DEPLOY_PROVIDERS.VOLCENGINE_DCDN, "common.provider.volcengine.dcdn"],
    [DEPLOY_PROVIDERS.VOLCENGINE_CLB, "common.provider.volcengine.clb"],
    [DEPLOY_PROVIDERS.VOLCENGINE_LIVE, "common.provider.volcengine.live"],
    [DEPLOY_PROVIDERS.QINIU_CDN, "common.provider.qiniu.cdn"],
    [DEPLOY_PROVIDERS.QINIU_PILI, "common.provider.qiniu.pili"],
    [DEPLOY_PROVIDERS.DOGECLOUD_CDN, "common.provider.dogecloud.cdn"],
    [DEPLOY_PROVIDERS.BYTEPLUS_CDN, "common.provider.byteplus.cdn"],
    [DEPLOY_PROVIDERS.UCLOUD_US3, "common.provider.ucloud.us3"],
    [DEPLOY_PROVIDERS.UCLOUD_UCDN, "common.provider.ucloud.ucdn"],
    [DEPLOY_PROVIDERS.AWS_CLOUDFRONT, "common.provider.aws.cloudfront"],
    [DEPLOY_PROVIDERS.EDGIO_APPLICATIONS, "common.provider.edgio.applications"],
  ].map(([type, name]) => [
    type,
    {
      type: type as DeployProviderType,
      name: name,
      icon: accessProvidersMap.get(type.split("-")[0])!.icon,
      provider: type.split("-")[0] as AccessProviderType,
    },
  ])
);
// #endregion
